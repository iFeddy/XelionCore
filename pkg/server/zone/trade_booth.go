package zone

import (
	"encoding/binary"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// === Dept 19: SafeTrade (`NC_TRADE`) ===

const maxTradeBoardSlots = 12

type TradeBoardSlot struct {
	Occupied  bool
	InvenSlot uint8
	ItemKey   int64
	ItemID    uint16
	Lot       uint32
}

type SafeTradeSide struct {
	CharID   int
	Handle   uint16
	Cen      uint64
	Locked   bool
	Decided  bool
	Board    [maxTradeBoardSlots]TradeBoardSlot
}

type SafeTradeSession struct {
	mu       sync.Mutex
	SideA    *SafeTradeSide
	SideB    *SafeTradeSide
	Finished bool
}

func (s *SafeTradeSession) getSides(charID int) (mine, peer *SafeTradeSide) {
	if s.SideA.CharID == charID {
		return s.SideA, s.SideB
	}
	return s.SideB, s.SideA
}

func (s *SafeTradeSession) unlockBoth(h *Handler, pA, pB *ZonePlayer) {
	s.SideA.Locked = false
	s.SideA.Decided = false
	s.SideB.Locked = false
	s.SideB.Decided = false
	if pA != nil && pA.Client != nil {
		_ = pA.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_BOARDUNLOCK_CMD))
		_ = pA.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_OPPOSITBOARDUNLOCK_CMD))
	}
	if pB != nil && pB.Client != nil {
		_ = pB.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_BOARDUNLOCK_CMD))
		_ = pB.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_OPPOSITBOARDUNLOCK_CMD))
	}
}

func (h *Handler) findPlayerByHandle(handle uint16) *ZonePlayer {
	if handle == 0 {
		return nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, p := range h.players {
		if p != nil && p.MapObjectID == handle {
			return p
		}
	}
	return nil
}

func (h *Handler) findPlayerByCharID(charID int) *ZonePlayer {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, p := range h.players {
		if p != nil && p.Char != nil && p.Char.ID == charID {
			return p
		}
	}
	return nil
}

func (h *Handler) inRangeSq(a, b *ZonePlayer, maxSq int64) bool {
	if a == nil || b == nil || a.Char == nil || b.Char == nil || a.CurrentMap == nil || a.CurrentMap != b.CurrentMap {
		return false
	}
	dx := int64(a.Char.XPos) - int64(b.Char.XPos)
	dy := int64(a.Char.YPos) - int64(b.Char.YPos)
	return dx*dx+dy*dy <= maxSq
}

func sendTradeErrPacket(client *network.Client, op uint16, errCode uint16) {
	if client == nil {
		return
	}
	pw := protocol.NewPacketWriter(op)
	pw.WriteUint16(errCode)
	_ = client.SendPacket(pw)
}

// 1. 0x4C01 PROTO_NC_TRADE_PROPOSE_REQ
// Verificado en Zone.exe sp_NC_TRADE_PROPOSE_REQ (0x005925A0):
// errores vía PROTO_NC_TRADE_PROPOSENO_ACK (0x4C04): estado inválido 0x06C2, distancia>1000 0x06C3, objetivo ocupado/booth/muerto 0x06CA, sin objetivo 0x06C8.
func (h *Handler) handleTradeProposeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	targetHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	if player.State == 3 || player.State == 5 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_PROPOSENO_ACK, 0x06C2)
		return
	}
	target := h.findPlayerByHandle(targetHandle)
	if target == nil || target == player || target.Char == nil || target.Client == nil {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_PROPOSENO_ACK, 0x06C8)
		return
	}
	if !h.inRangeSq(player, target, 1000000) {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_PROPOSENO_ACK, 0x06C3)
		return
	}
	if target.ActiveTrade != nil || target.ActiveBooth != nil || target.State == 3 || target.State == 5 {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_PROPOSENO_ACK, 0x06CA)
		return
	}
	target.PendingTradeFrom = player.MapObjectID
	ask := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_PROPOSE_ASK_REQ)
	ask.WriteUint16(player.MapObjectID)
	_ = target.Client.SendPacket(ask)
}

// 2. 0x4C03 PROTO_NC_TRADE_PROPOSE_ASKNO_ACK
func (h *Handler) handleTradeProposeAskNoAck(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	proposer := h.findPlayerByHandle(player.PendingTradeFrom)
	player.PendingTradeFrom = 0
	if proposer != nil && proposer.Client != nil {
		sendTradeErrPacket(proposer.Client, protocol.PROTO_NC_TRADE_PROPOSENO_ACK, 0x06C7)
	}
}

// 3. 0x4C06 PROTO_NC_TRADE_PROPOSEYES_ACK
func (h *Handler) handleTradeProposeYesAck(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	proposer := h.findPlayerByHandle(player.PendingTradeFrom)
	player.PendingTradeFrom = 0
	if proposer == nil || proposer.Char == nil || proposer.Client == nil || proposer.ActiveTrade != nil || player.ActiveTrade != nil ||
		player.State == 3 || player.State == 5 || proposer.State == 3 || proposer.State == 5 {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_PROPOSENO_ACK, 0x06C8)
		return
	}
	if !h.inRangeSq(player, proposer, 1000000) {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_PROPOSENO_ACK, 0x06C3)
		return
	}

	sess := &SafeTradeSession{
		SideA: &SafeTradeSide{CharID: proposer.Char.ID, Handle: proposer.MapObjectID},
		SideB: &SafeTradeSide{CharID: player.Char.ID, Handle: player.MapObjectID},
	}
	proposer.ActiveTrade = sess
	player.ActiveTrade = sess

	_ = proposer.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_PROPOSE_ASKYES_ACK))

	startA := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_START_CMD)
	startA.WriteUint16(player.MapObjectID)
	_ = proposer.Client.SendPacket(startA)

	startB := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_START_CMD)
	startB.WriteUint16(proposer.MapObjectID)
	_ = player.Client.SendPacket(startB)
}

// 4. 0x4C07 PROTO_NC_TRADE_PROPOSE_CANCEL_CMD
func (h *Handler) handleTradeProposeCancelCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	// Limpiar propuestas entrantes hacia el jugador que cancela
	h.mu.RLock()
	for _, other := range h.players {
		if other != nil && other.PendingTradeFrom == player.MapObjectID {
			other.PendingTradeFrom = 0
			if other.Client != nil {
				_ = other.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_PROPOSE_CANCELED_CMD))
			}
		}
	}
	h.mu.RUnlock()
}

func (h *Handler) cancelTradeSession(player *ZonePlayer) {
	if player == nil || player.Char == nil || player.ActiveTrade == nil {
		return
	}
	sess := player.ActiveTrade
	sess.mu.Lock()
	if sess.Finished {
		sess.mu.Unlock()
		return
	}
	sess.Finished = true
	_, peerSide := sess.getSides(player.Char.ID)
	sess.mu.Unlock()

	player.ActiveTrade = nil
	if peerSide != nil {
		if peer := h.findPlayerByCharID(peerSide.CharID); peer != nil {
			peer.ActiveTrade = nil
			if peer.Client != nil {
				_ = peer.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_CANCEL_CMD))
			}
		}
	}
	if player.Client != nil {
		_ = player.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_CANCEL_ACK))
	}
}

// 5. 0x4C0A PROTO_NC_TRADE_CANCEL_REQ & 0x200B PROTO_NC_ACT_ENDOFTRADE_CMD
func (h *Handler) handleTradeCancelReq(client *network.Client, pkt *protocol.Packet) {
	if player, ok := h.getPlayer(client); ok {
		h.cancelTradeSession(player)
	}
}

func (h *Handler) handleEndOfTradeCmd(client *network.Client, pkt *protocol.Packet) {
	if player, ok := h.getPlayer(client); ok && player != nil {
		if player.ActiveTrade != nil {
			h.cancelTradeSession(player)
		}
		h.resetPlayerNPCState(player)
	}
}

// 6. 0x4C0D PROTO_NC_TRADE_UPBOARD_REQ
// Verificado en Zone.exe sp_NC_TRADE_UPBOARD_REQ (0x00595280):
// errores PROTO_NC_TRADE_UPBOARDFAIL_ACK (0x4C0E): estado 0x0842, faltante 0x0843, slot lleno 0x0844, duplicado 0x0846, NoTrade/Ligado 0x0851; distancia>1000 -> sp_tradebreak.
func (h *Handler) handleTradeUpboardReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	if player.ActiveTrade == nil || player.State == 3 || player.State == 5 {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_UPBOARDFAIL_ACK, 0x0842)
		return
	}
	slotInven, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	sess := player.ActiveTrade
	sess.mu.Lock()
	if sess.Finished {
		sess.mu.Unlock()
		return
	}
	mine, peerSide := sess.getSides(player.Char.ID)
	peer := h.findPlayerByCharID(peerSide.CharID)
	if !h.inRangeSq(player, peer, 1000000) {
		sess.mu.Unlock()
		h.cancelTradeSession(player)
		return
	}

	it := h.findItemAt(player.Char.ID, 9, int16(slotInven))
	if it == nil {
		sess.mu.Unlock()
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_UPBOARDFAIL_ACK, 0x0843)
		return
	}
	if h.getBelongRule(it.ItemID).NoTrade || h.isItemBound(it) {
		sess.mu.Unlock()
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_UPBOARDFAIL_ACK, 0x0851)
		return
	}

	freeBoard := -1
	for i := 0; i < maxTradeBoardSlots; i++ {
		if mine.Board[i].Occupied && mine.Board[i].ItemKey == it.ItemKey {
			sess.mu.Unlock()
			sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_UPBOARDFAIL_ACK, 0x0846)
			return
		}
		if !mine.Board[i].Occupied && freeBoard < 0 {
			freeBoard = i
		}
	}
	if freeBoard < 0 {
		sess.mu.Unlock()
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_UPBOARDFAIL_ACK, 0x0844)
		return
	}

	lot := uint32(it.Flags)
	if lot == 0 {
		lot = 1
	}
	mine.Board[freeBoard] = TradeBoardSlot{
		Occupied:  true,
		InvenSlot: slotInven,
		ItemKey:   it.ItemKey,
		ItemID:    it.ItemID,
		Lot:       lot,
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_UPBOARD_ACK)
	ack.WriteByte(slotInven)
	ack.WriteByte(byte(freeBoard))
	_ = client.SendPacket(ack)

	if peer != nil && peer.Client != nil {
		// Zone.exe 0x00595610: PROTO_NC_TRADE_OPPOSITUPBOARD_CMD = slotboard(1B) + iti_fill2buffer_union(itemid+attr)
		opp := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_OPPOSITUPBOARD_CMD)
		opp.WriteByte(byte(freeBoard))
		opp.WriteBytes(h.buildItemInlineAttr(it))
		_ = peer.Client.SendPacket(opp)
	}
	sess.unlockBoth(h, player, peer)
	sess.mu.Unlock()
}

// 7. 0x4C11 PROTO_NC_TRADE_DOWNBOARD_REQ
func (h *Handler) handleTradeDownboardReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	if player.ActiveTrade == nil || player.State == 3 || player.State == 5 {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_DOWNBOARDFAIL_ACK, 0x0842)
		return
	}
	slotBoard, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	sess := player.ActiveTrade
	sess.mu.Lock()
	if sess.Finished {
		sess.mu.Unlock()
		return
	}
	mine, peerSide := sess.getSides(player.Char.ID)
	peer := h.findPlayerByCharID(peerSide.CharID)
	if !h.inRangeSq(player, peer, 1000000) {
		sess.mu.Unlock()
		h.cancelTradeSession(player)
		return
	}
	if int(slotBoard) >= maxTradeBoardSlots {
		sess.mu.Unlock()
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_DOWNBOARDFAIL_ACK, 0x0843)
		return
	}
	if !mine.Board[slotBoard].Occupied {
		sess.mu.Unlock()
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_DOWNBOARDFAIL_ACK, 0x0847)
		return
	}
	mine.Board[slotBoard] = TradeBoardSlot{}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_DOWNBOARD_ACK)
	ack.WriteByte(slotBoard)
	_ = client.SendPacket(ack)

	if peer != nil && peer.Client != nil {
		opp := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_OPPOSITDOWNBOARD_CMD)
		opp.WriteByte(slotBoard)
		_ = peer.Client.SendPacket(opp)
	}
	sess.unlockBoth(h, player, peer)
	sess.mu.Unlock()
}

// 8. 0x4C15 PROTO_NC_TRADE_CENBOARDING_REQ
func (h *Handler) handleTradeCenBoardingReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	if player.ActiveTrade == nil || player.State == 3 || player.State == 5 {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_CENBOARDINGFAIL_ACK, 0x0842)
		return
	}
	cen, err := pkt.Reader.ReadUint64()
	if err != nil {
		return
	}

	sess := player.ActiveTrade
	sess.mu.Lock()
	if sess.Finished {
		sess.mu.Unlock()
		return
	}
	mine, peerSide := sess.getSides(player.Char.ID)
	peer := h.findPlayerByCharID(peerSide.CharID)
	if !h.inRangeSq(player, peer, 1000000) {
		sess.mu.Unlock()
		h.cancelTradeSession(player)
		return
	}
	if player.Char.TCharacterState.Money < 0 || uint64(player.Char.TCharacterState.Money) < cen {
		sess.mu.Unlock()
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_CENBOARDINGFAIL_ACK, 0x0848)
		return
	}

	mine.Cen = cen

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_CENBOARDING_ACK)
	ack.WriteUint64(cen)
	_ = client.SendPacket(ack)

	if peer != nil && peer.Client != nil {
		opp := protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_OPPOSITCENBOARDING_CMD)
		opp.WriteUint64(cen)
		_ = peer.Client.SendPacket(opp)
	}
	sess.unlockBoth(h, player, peer)
	sess.mu.Unlock()
}

// 9. 0x4C19 PROTO_NC_TRADE_BOARDLOCK_REQ
func (h *Handler) handleTradeBoardLockReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	if player.ActiveTrade == nil || player.State == 3 || player.State == 5 {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_BOARDLOCKFAIL_ACK, 0x0842)
		return
	}
	sess := player.ActiveTrade
	sess.mu.Lock()
	if sess.Finished {
		sess.mu.Unlock()
		return
	}
	mine, peerSide := sess.getSides(player.Char.ID)
	peer := h.findPlayerByCharID(peerSide.CharID)
	if !h.inRangeSq(player, peer, 1000000) {
		sess.mu.Unlock()
		h.cancelTradeSession(player)
		return
	}
	mine.Locked = true

	_ = client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_BOARDLOCK_ACK))
	if peer != nil && peer.Client != nil {
		_ = peer.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_OPPOSITBOARDLOCK_CMD))
	}
	sess.mu.Unlock()
}

// 10. 0x4C1F PROTO_NC_TRADE_DECIDE_REQ — intercambio atómico libre de duplicación
func (h *Handler) handleTradeDecideReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	if player.ActiveTrade == nil || player.State == 3 || player.State == 5 {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_DECIDEFAIL_ACK, 0x0842)
		return
	}
	sess := player.ActiveTrade
	sess.mu.Lock()
	defer sess.mu.Unlock()
	if sess.Finished {
		return
	}
	mine, peerSide := sess.getSides(player.Char.ID)
	if !mine.Locked {
		sendTradeErrPacket(client, protocol.PROTO_NC_TRADE_DECIDEFAIL_ACK, 0x084A)
		return
	}
	mine.Decided = true
	_ = client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_DECIDE_ACK))

	peer := h.findPlayerByCharID(peerSide.CharID)
	if peer != nil && peer.Client != nil {
		_ = peer.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_OPPOSITDECIDE_CMD))
	}

	if !peerSide.Locked || !peerSide.Decided || peer == nil || peer.Char == nil {
		return
	}

	// Ambos jugadores bloquearon y confirmaron: validar inventario, distancia, Cen e ítems anti-duplicación
	pA := h.findPlayerByCharID(sess.SideA.CharID)
	pB := h.findPlayerByCharID(sess.SideB.CharID)
	if pA == nil || pB == nil || pA.Char == nil || pB.Char == nil ||
		!h.inRangeSq(pA, pB, 1000000) ||
		pA.Char.TCharacterState.Money < 0 || pB.Char.TCharacterState.Money < 0 ||
		uint64(pA.Char.TCharacterState.Money) < sess.SideA.Cen ||
		uint64(pB.Char.TCharacterState.Money) < sess.SideB.Cen {
		sess.Finished = true
		if pA != nil {
			pA.ActiveTrade = nil
			_ = pA.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADEFAIL_CMD))
		}
		if pB != nil {
			pB.ActiveTrade = nil
			_ = pB.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADEFAIL_CMD))
		}
		return
	}

	itemsFromA := h.collectValidatedTradeItems(pA, &sess.SideA.Board)
	itemsFromB := h.collectValidatedTradeItems(pB, &sess.SideB.Board)
	freeA := h.listFreeInvenSlots(pA.Char.ID, len(itemsFromB)+len(itemsFromA))
	freeB := h.listFreeInvenSlots(pB.Char.ID, len(itemsFromA)+len(itemsFromB))

	if len(freeA) < len(itemsFromB) || len(freeB) < len(itemsFromA) {
		sess.Finished = true
		pA.ActiveTrade = nil
		pB.ActiveTrade = nil
		_ = pA.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADEFAIL_CMD))
		_ = pB.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADEFAIL_CMD))
		return
	}

	// Ejecutar transferencias atómicas de ítems (CALL usp_SafeTrade_setAll + mover en RAM/UI)
	for i, it := range itemsFromA {
		dstSlot := freeB[i]
		if err := h.charRepo.SafeTradeAll(pB.Char.ID, pA.Char.ID, 0, it.ItemKey, 9, dstSlot); err != nil {
			sess.Finished = true
			pA.ActiveTrade = nil
			pB.ActiveTrade = nil
			_ = pA.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADEFAIL_CMD))
			_ = pB.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADEFAIL_CMD))
			return
		}
		oldSlot := it.Storage
		it.Owner = pB.Char.ID
		it.Storage = dstSlot
		h.sendCellChangeFull(pA.Client, 9, oldSlot, nil, 0xFFFF)
		h.sendCellChangeFull(pB.Client, 9, dstSlot, it, 0xFFFF)
	}
	for i, it := range itemsFromB {
		dstSlot := freeA[i]
		if err := h.charRepo.SafeTradeAll(pA.Char.ID, pB.Char.ID, 0, it.ItemKey, 9, dstSlot); err != nil {
			sess.Finished = true
			pA.ActiveTrade = nil
			pB.ActiveTrade = nil
			_ = pA.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADEFAIL_CMD))
			_ = pB.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADEFAIL_CMD))
			return
		}
		oldSlot := it.Storage
		it.Owner = pA.Char.ID
		it.Storage = dstSlot
		h.sendCellChangeFull(pB.Client, 9, oldSlot, nil, 0xFFFF)
		h.sendCellChangeFull(pA.Client, 9, dstSlot, it, 0xFFFF)
	}

	// Transferir Cen neta
	if sess.SideA.Cen > 0 || sess.SideB.Cen > 0 {
		pA.Char.TCharacterState.Money = pA.Char.TCharacterState.Money - int64(sess.SideA.Cen) + int64(sess.SideB.Cen)
		pB.Char.TCharacterState.Money = pB.Char.TCharacterState.Money - int64(sess.SideB.Cen) + int64(sess.SideA.Cen)
		h.syncPlayerMoney(pA)
		h.syncPlayerMoney(pB)
	}

	sess.Finished = true
	pA.ActiveTrade = nil
	pB.ActiveTrade = nil
	_ = pA.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADECOMPLETE_CMD))
	_ = pB.Client.SendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_TRADE_TRADECOMPLETE_CMD))

	if h.gameLogDispatcher != nil && pA.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType4(1, pA.Char.ID, pA.CurrentMap.MapID, pA.Char.XPos, pA.Char.YPos, 0, pB.Char.ID, int(sess.SideA.Cen), len(itemsFromA))
		h.gameLogDispatcher.EnqueueType4(1, pB.Char.ID, pA.CurrentMap.MapID, pB.Char.XPos, pB.Char.YPos, 0, pA.Char.ID, int(sess.SideB.Cen), len(itemsFromB))
	}
}

func (h *Handler) collectValidatedTradeItems(player *ZonePlayer, board *[maxTradeBoardSlots]TradeBoardSlot) []*database.Item {
	var out []*database.Item
	for _, slot := range board {
		if !slot.Occupied {
			continue
		}
		it := h.findItemAt(player.Char.ID, 9, int16(slot.InvenSlot))
		if it != nil && it.ItemKey == slot.ItemKey {
			out = append(out, it)
		}
	}
	return out
}

func (h *Handler) listFreeInvenSlots(charID int, limit int) []int16 {
	items, _ := h.charRepo.GetCharacterItems(charID)
	used := make(map[int16]bool)
	for _, it := range items {
		if it != nil && it.StorageType == 9 {
			used[it.Storage] = true
		}
	}
	var out []int16
	for s := int16(0); s < 96 && len(out) < limit; s++ {
		if !used[s] {
			out = append(out, s)
		}
	}
	return out
}

// === Dept 26: Tiendas Personales / Kioscos (`NC_BOOTH`) ===

type BoothItemSlot struct {
	BoothSlot  uint8
	InvenSlot  uint8
	ItemKey    int64
	ItemID     uint16
	UnitCost   uint64
	Lot        uint16
	ItemHandle uint32
}

type PlayerBooth struct {
	mu          sync.Mutex
	OwnerHandle uint16
	OwnerCharID int
	IsSell      bool
	Signboard   [30]byte
	Slots       map[uint8]*BoothItemSlot
}

// 1. 0x6801 PROTO_NC_BOOTH_OPEN_REQ
// Verificado en Zone.exe sp_NC_BOOTH_OPEN_REQ (0x0058E900):
// errores PROTO_NC_BOOTH_OPEN_ACK (0x6802): estado 0x0A42, ítem faltante 0x0A43, NoTrade/Ligado 0x0A44, cerca de otro booth (<400u) 0x0A53; éxito 0x0A41.
func (h *Handler) handleBoothOpenReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || player.CurrentMap == nil {
		return
	}
	sendOpenAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_OPEN_ACK)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}
	if player.State == 3 || player.State == 5 || player.ActiveTrade != nil || player.ActiveBooth != nil {
		sendOpenAck(0x0A42)
		return
	}
	signBytes, err := pkt.Reader.ReadBytes(30)
	if err != nil {
		return
	}
	flag, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	// Distancia mínima a otro kiosco en el mismo mapa (Zone.exe 0x58EA9F range 0x190 = 400)
	h.mu.RLock()
	tooClose := false
	for _, other := range h.players {
		if other != nil && other != player && other.ActiveBooth != nil && other.CurrentMap == player.CurrentMap {
			if h.inRangeSq(player, other, 160000) {
				tooClose = true
				break
			}
		}
	}
	h.mu.RUnlock()
	if tooClose {
		sendOpenAck(0x0A53)
		return
	}

	isSell := (flag & 0x01) != 0
	numItems := int(flag >> 1)

	booth := &PlayerBooth{
		OwnerHandle: player.MapObjectID,
		OwnerCharID: player.Char.ID,
		IsSell:      isSell,
		Slots:       make(map[uint8]*BoothItemSlot),
	}
	copy(booth.Signboard[:], signBytes)

	for i := 0; i < numItems; i++ {
		invSlot, err := pkt.Reader.ReadByte()
		if err != nil {
			break
		}
		boothSlot, _ := pkt.Reader.ReadByte()
		unitCost, _ := pkt.Reader.ReadUint64()
		lot, _ := pkt.Reader.ReadUint16()
		if lot == 0 {
			lot = 1
		}
		entry := &BoothItemSlot{
			BoothSlot:  boothSlot,
			InvenSlot:  invSlot,
			UnitCost:   unitCost,
			Lot:        lot,
			ItemHandle: (uint32(player.MapObjectID) << 16) | uint32(boothSlot),
		}
		if isSell {
			it := h.findItemAt(player.Char.ID, 9, int16(invSlot))
			if it == nil {
				sendOpenAck(0x0A43)
				return
			}
			if h.getBelongRule(it.ItemID).NoTrade || h.isItemBound(it) {
				sendOpenAck(0x0A44)
				return
			}
			entry.ItemKey = it.ItemKey
			entry.ItemID = it.ItemID
		} else {
			// Modo Compra: InvenSlot lleva ID o lookup opcional
			entry.ItemID = uint16(invSlot)
		}
		booth.Slots[boothSlot] = entry
	}

	player.ActiveBooth = booth
	player.State = 5 // Vendor

	sendOpenAck(0x0A41)

	var tent [12]byte
	copy(tent[:], "Tent01")
	broad := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_SOMEONEOPEN_CMD)
	broad.WriteUint16(player.MapObjectID)
	broad.WriteBytes(tent[:])
	if isSell {
		broad.WriteByte(1)
	} else {
		broad.WriteByte(0)
	}
	broad.WriteBytes(booth.Signboard[:])
	player.CurrentMap.Broadcast(broad, 0)
}

func (h *Handler) closePlayerBooth(player *ZonePlayer) {
	if player == nil || player.ActiveBooth == nil {
		return
	}
	player.ActiveBooth = nil
	player.State = 1
	if player.CurrentMap != nil {
		closeBroad := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_SOMEONECLOSE_CMD)
		closeBroad.WriteUint16(player.MapObjectID)
		var shape [43]byte
		closeBroad.WriteBytes(shape[:])
		player.CurrentMap.Broadcast(closeBroad, 0)

		closedSearch := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD)
		closedSearch.WriteUint16(player.MapObjectID)
		player.CurrentMap.Broadcast(closedSearch, 0)
	}
}

// 2. 0x6804 PROTO_NC_BOOTH_CLOSE_REQ
func (h *Handler) handleBoothCloseReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_CLOSE_ACK)
	if player.ActiveBooth == nil {
		ack.WriteUint16(0x0A44)
		_ = client.SendPacket(ack)
		return
	}
	h.closePlayerBooth(player)
	ack.WriteUint16(0x0A41)
	_ = client.SendPacket(ack)
}

func (h *Handler) sendBoothEntry(client *network.Client, owner *ZonePlayer, isRefresh bool) {
	if owner == nil || owner.ActiveBooth == nil {
		return
	}
	b := owner.ActiveBooth
	b.mu.Lock()
	defer b.mu.Unlock()

	var op uint16
	if b.IsSell {
		op = protocol.PROTO_NC_BOOTH_ENTRY_SELL_ACK
		if isRefresh {
			op = protocol.PROTO_NC_BOOTH_REFRESH_SELL_ACK
		}
	} else {
		op = protocol.PROTO_NC_BOOTH_ENTRY_BUY_ACK
		if isRefresh {
			op = protocol.PROTO_NC_BOOTH_REFRESH_BUY_ACK
		}
	}

	// Zone.exe StreetBoothBuy/Sell::sb_Show2Client (0x0058D7FA / 0x0058F2E9): nError = 0x0A41
	pw := protocol.NewPacketWriter(op)
	pw.WriteUint16(0x0A41)
	pw.WriteUint16(owner.MapObjectID)
	pw.WriteByte(byte(len(b.Slots)))
	for _, s := range b.Slots {
		if b.IsSell {
			// sbs_FillBufferItem (0x0058D900): datasize(1B = 9+len(inline)) | slot(1B) | unitcost(8B) | inline(itemid+attr)
			it := h.findItemAt(owner.Char.ID, 9, int16(s.InvenSlot))
			if it == nil {
				it = &database.Item{ItemID: s.ItemID, Flags: int32(s.Lot)}
			}
			inline := h.buildItemInlineAttr(it)
			pw.WriteByte(byte(9 + len(inline)))
			pw.WriteByte(s.BoothSlot)
			pw.WriteUint64(s.UnitCost)
			pw.WriteBytes(inline)
		} else {
			// StreetBoothBuy::sb_Show2Client (0x0058D82C, stride 13B): slot(1B) | itemid(2B) | unitcost(8B) | restquantity(2B)
			pw.WriteByte(s.BoothSlot)
			pw.WriteUint16(s.ItemID)
			pw.WriteUint64(s.UnitCost)
			pw.WriteUint16(s.Lot)
		}
	}
	_ = client.SendPacket(pw)
}

// 3. 0x6807 PROTO_NC_BOOTH_ENTRY_REQ & 4. 0x680A PROTO_NC_BOOTH_REFRESH_REQ
func (h *Handler) handleBoothEntryReq(client *network.Client, pkt *protocol.Packet) {
	handle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	owner := h.findPlayerByHandle(handle)
	if owner == nil || owner.ActiveBooth == nil || owner.State != 5 {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_ENTRY_SELL_ACK)
		pw.WriteUint16(0x0A5B)
		pw.WriteUint16(handle)
		pw.WriteByte(0)
		_ = client.SendPacket(pw)
		return
	}
	h.sendBoothEntry(client, owner, false)
}

func (h *Handler) handleBoothRefreshReq(client *network.Client, pkt *protocol.Packet) {
	handle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	owner := h.findPlayerByHandle(handle)
	if owner == nil || owner.ActiveBooth == nil || owner.State != 5 {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_REFRESH_SELL_ACK)
		pw.WriteUint16(0x0A5B)
		pw.WriteUint16(handle)
		pw.WriteByte(0)
		_ = client.SendPacket(pw)
		return
	}
	h.sendBoothEntry(client, owner, true)
}

// 5. 0x680D PROTO_NC_BOOTH_ITEMTRADE_REQ
// Verificado en Zone.exe sp_NC_BOOTH_ITEMTRADE_REQ (0x0058DA30):
// errores PROTO_NC_BOOTH_ITEMTRADE_ACK (0x680E): estado 0x0A42, booth cerrado 0x0A5B, distancia>=300 0x0A4C, sin cen 0x020E, inventario lleno 0x0203; éxito 0x0A41.
func (h *Handler) handleBoothItemTradeReq(client *network.Client, pkt *protocol.Packet) {
	buyer, ok := h.getPlayer(client)
	if !ok || buyer == nil || buyer.Char == nil || h.charRepo == nil {
		return
	}
	boothHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	reqNum, err := pkt.Reader.ReadUint16()
	if err != nil || reqNum == 0 {
		return
	}

	sendAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_ITEMTRADE_ACK)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	if buyer.State == 3 || buyer.State == 5 || buyer.ActiveTrade != nil {
		sendAck(0x0A42)
		return
	}

	seller := h.findPlayerByHandle(boothHandle)
	if seller == nil || seller.Char == nil || seller == buyer {
		sendAck(0x0A42)
		return
	}
	if seller.ActiveBooth == nil || seller.State != 5 {
		sendAck(0x0A5B)
		return
	}
	if !h.inRangeSq(buyer, seller, 89999) {
		sendAck(0x0A4C)
		return
	}

	booth := seller.ActiveBooth
	booth.mu.Lock()
	defer booth.mu.Unlock()

	entry, exists := booth.Slots[slot]
	if !exists || entry == nil || reqNum > entry.Lot || !booth.IsSell {
		sendAck(0x0A42)
		return
	}
	totalCost := entry.UnitCost * uint64(reqNum)
	if buyer.Char.TCharacterState.Money < 0 || uint64(buyer.Char.TCharacterState.Money) < totalCost {
		sendAck(0x020E)
		return
	}

	srcItem := h.findItemAt(seller.Char.ID, 9, int16(entry.InvenSlot))
	srcLot := uint32(0)
	if srcItem != nil {
		srcLot = uint32(srcItem.Flags)
		if srcLot == 0 {
			srcLot = 1
		}
	}
	if srcItem == nil || srcItem.ItemKey != entry.ItemKey || srcLot < uint32(reqNum) {
		sendAck(0x0A43)
		return
	}

	freeBuyerSlot, ok := h.findFreeInvenSlot(buyer.Char.ID)
	if !ok {
		sendAck(0x0203)
		return
	}

	if srcLot == uint32(reqNum) {
		if err := h.charRepo.SafeTradeAll(buyer.Char.ID, seller.Char.ID, 0, srcItem.ItemKey, 9, freeBuyerSlot); err != nil {
			sendAck(0x0A42)
			return
		}
		h.sendCellChangeFull(seller.Client, 9, srcItem.Storage, nil, 0xFFFF)
		srcItem.Owner = buyer.Char.ID
		srcItem.Storage = freeBuyerSlot
		h.sendCellChangeFull(buyer.Client, 9, freeBuyerSlot, srcItem, 0xFFFF)
		delete(booth.Slots, slot)
	} else {
		if _, added := h.addItemToInventory(buyer, srcItem.ItemID, uint32(reqNum), nil); !added {
			sendAck(0x0203)
			return
		}
		h.consumeItemUnits(seller, srcItem, uint32(reqNum))
		entry.Lot -= reqNum
	}

	buyer.Char.TCharacterState.Money -= int64(totalCost)
	seller.Char.TCharacterState.Money += int64(totalCost)
	h.syncPlayerMoney(buyer)
	h.syncPlayerMoney(seller)

	rem := uint16(0)
	if s2, ok := booth.Slots[slot]; ok {
		rem = s2.Lot
	}
	ref := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_BUYREFRESH_CMD)
	ref.WriteByte(slot)
	ref.WriteUint16(rem)
	if seller.Client != nil {
		_ = seller.Client.SendPacket(ref)
	}
	_ = client.SendPacket(ref)

	sendAck(0x0A41)

	if h.gameLogDispatcher != nil && buyer.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueTypeC(1, buyer.Char.ID, buyer.CurrentMap.MapID, buyer.Char.XPos, buyer.Char.YPos, 0, seller.Char.ID, entry.ItemKey, int(entry.ItemID), int(reqNum), int(totalCost))
	}
}

// 6. 0x6810 PROTO_NC_BOOTH_INTERIORSTART_REQ
func (h *Handler) handleBoothInteriorStartReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	sign, _ := pkt.Reader.ReadBytes(30)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_INTERIORSTART_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	var sb [30]byte
	copy(sb[:], sign)
	broad := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_SOMEONEINTERIORSTART_CMD)
	broad.WriteUint16(player.MapObjectID)
	broad.WriteBytes(sb[:])
	player.CurrentMap.Broadcast(broad, 0)
}

// 7. 0x6813 PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ (contra MarketSearchInfo.shn / ItemInfoServer.shn)
func (h *Handler) handleBoothSearchItemListCategorizedReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	category, err := pkt.Reader.ReadUint32()
	if err != nil {
		return
	}
	onlyMyClass, _ := pkt.Reader.ReadByte()
	_, _ = pkt.Reader.ReadByte()

	type searchMatch struct {
		handle    uint32
		owner     uint16
		unitPrice uint64
		inform    []byte
	}
	var matches []searchMatch

	h.mu.RLock()
	for _, p := range h.players {
		if p == nil || p.ActiveBooth == nil || p.CurrentMap != player.CurrentMap || !p.ActiveBooth.IsSell {
			continue
		}
		b := p.ActiveBooth
		b.mu.Lock()
		for _, s := range b.Slots {
			if category != 0 && h.itemStore != nil {
				if info, ok := h.itemStore.GetItemInfo(s.ItemID); ok && info != nil {
					if cat := h.getMarketSearchType(info.InxName); cat != 0 && cat != category {
						continue
					}
					if onlyMyClass != 0 {
						if can, know := h.itemStore.CanUseItemByClassID(s.ItemID, player.Char.Job); know && !can {
							continue
						}
					}
				}
			}
			it := h.findItemAt(p.Char.ID, 9, int16(s.InvenSlot))
			if it == nil {
				it = &database.Item{ItemID: s.ItemID, Flags: int32(s.Lot)}
			}
			inform := h.buildItemInlineAttr(it)
			matches = append(matches, searchMatch{
				handle:    s.ItemHandle,
				owner:     p.MapObjectID,
				unitPrice: s.UnitCost,
				inform:    inform,
			})
			if len(matches) >= 50 {
				break
			}
		}
		b.mu.Unlock()
		if len(matches) >= 50 {
			break
		}
	}
	h.mu.RUnlock()

	// Zone.exe sp_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ (0x00439D06 / 0x00439F42..0x439FB7):
	// nError = 0x31C8, nNumOfItems, cada BoothItemInfo = datasize(1B = 14+len(inline)) | handle(4) | owner(2) | price(8) | inline
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK)
	ack.WriteUint16(0x31C8)
	ack.WriteByte(byte(len(matches)))
	for _, m := range matches {
		ack.WriteByte(byte(14 + len(m.inform)))
		ack.WriteUint32(m.handle)
		ack.WriteUint16(m.owner)
		ack.WriteUint64(m.unitPrice)
		ack.WriteBytes(m.inform)
	}
	_ = client.SendPacket(ack)
	_ = strings.Builder{}
	_ = binary.LittleEndian
}

// 8. 0x6815 PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ
func (h *Handler) handleBoothSearchBoothPositionReq(client *network.Client, pkt *protocol.Packet) {
	_, err := pkt.Reader.ReadUint32() // nBoothItemHandle
	if err != nil {
		return
	}
	ownerHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	owner := h.findPlayerByHandle(ownerHandle)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_ACK)
	if owner == nil || owner.Char == nil || owner.ActiveBooth == nil {
		ack.WriteUint16(0x31D3)
		ack.WriteUint16(ownerHandle)
		ack.WriteUint32(0)
		ack.WriteUint32(0)
	} else {
		ack.WriteUint16(0x31C8)
		ack.WriteUint16(ownerHandle)
		ack.WriteUint32(uint32(owner.Char.XPos))
		ack.WriteUint32(uint32(owner.Char.YPos))
	}
	_ = client.SendPacket(ack)
}
