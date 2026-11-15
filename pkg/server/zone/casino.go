package zone

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

type diceTaisaiBet struct {
	CharNo    int
	Handle    uint16
	BetIndex  uint16
	BetInfo   uint64
	BetType   uint8
	BetCoins  uint64
}

type diceTaisaiTable struct {
	mu           sync.Mutex
	Players      map[uint16]*ZonePlayer
	Bets         map[uint32]diceTaisaiBet // (charNo<<16)|betIdx
	DiceHistory  [21]byte
	ModeType     uint32
	GameStatus   uint32
	MaxBet       uint64
	TaxRate      uint16
}

var globalDiceTable = &diceTaisaiTable{
	Players: make(map[uint16]*ZonePlayer),
	Bets:    make(map[uint32]diceTaisaiBet),
	MaxBet:  100000,
	TaxRate: 5,
}

func (t *diceTaisaiTable) broadcast(pkt *protocol.PacketWriter) {
	t.mu.Lock()
	cls := make([]*network.Client, 0, len(t.Players))
	for _, p := range t.Players {
		if p != nil && p.Client != nil {
			cls = append(cls, p.Client)
		}
	}
	t.mu.Unlock()
	for _, c := range cls {
		_ = c.SendPacket(pkt)
	}
}

// ============================================================================
// Dept 12: PROTO_NC_ITEM_DICE_GAME_START_REQ (0x303C) - BelongDice / PartyDice
// ============================================================================

func (h *Handler) handleDiceGameStartReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 3 {
		return
	}
	itemHandle, _ := pkt.Reader.ReadUint16()
	isYes, _ := pkt.Reader.ReadByte()

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DICE_GAME_START_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	var roll int32
	if isYes != 0 {
		roll = int32(rand.Intn(100) + 1)
	}
	if h.charRepo != nil {
		if gameNo, err := h.charRepo.SetDiceGame(player.Char.ID, int64(itemHandle), fmt.Sprintf("%d", roll)); err == nil {
			_ = h.charRepo.SetDiceLog(gameNo, player.Char.ID, 0, 0, 1, isYes, 0, int64(roll))
		}
	}

	if player.CurrentMap != nil {
		startCmd := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DICE_GAME_START_CMD)
		startCmd.WriteUint16(player.MapObjectID)
		startCmd.WriteInt32(roll)
		player.CurrentMap.Broadcast(startCmd, 0)

		resCmd := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DICE_GAME_RESULT_CMD)
		resCmd.WriteUint16(0)
		writeFixedStr(resCmd, player.Char.Name, 20)
		resCmd.WriteInt32(roll)
		resCmd.WriteUint16(itemHandle)
		player.CurrentMap.Broadcast(resCmd, 0)
	}

	if h.gameLogDispatcher != nil {
		mapID := ""
		if player.CurrentMap != nil {
			mapID = player.CurrentMap.MapID
		}
		h.gameLogDispatcher.EnqueueTypeI(1, player.Char.ID, mapID, player.Char.XPos, player.Char.YPos, 0, int(roll), int(itemHandle))
	}
}

// ============================================================================
// Dept 43: NC_DICE_TAISAI (9 C2S)
// ============================================================================

// handleDiceTaisaiGameJoinReq maneja PROTO_NC_DICE_TAISAI_GAME_JOIN_REQ (0xAC01).
func (h *Handler) handleDiceTaisaiGameJoinReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil {
		return
	}
	if h.dataStore != nil && h.dataStore.CasinoCatalog != nil {
		globalDiceTable.TaxRate = uint16(h.dataStore.CasinoCatalog.TaxRatePct)
	}
	globalDiceTable.mu.Lock()
	globalDiceTable.Players[player.MapObjectID] = player
	player.CasinoDiceJoined = true
	cnt := uint8(len(globalDiceTable.Players))
	hnds := make([]uint16, 0, cnt)
	for hnd := range globalDiceTable.Players {
		hnds = append(hnds, hnd)
	}
	globalDiceTable.mu.Unlock()

	// PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK (50B + gamers*2B)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK)
	ack.WriteUint16(0)
	ack.WriteBytes(globalDiceTable.DiceHistory[:])
	ack.WriteUint32(globalDiceTable.GameStatus)
	ack.WriteUint32(globalDiceTable.ModeType)
	ack.WriteUint16(30)
	ack.WriteUint32(30000)
	ack.WriteUint16(globalDiceTable.TaxRate)
	ack.WriteUint64(globalDiceTable.MaxBet)
	ack.WriteUint16(1)
	ack.WriteByte(cnt)
	for _, gh := range hnds {
		ack.WriteUint16(gh)
	}
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_CMD)
	cmd.WriteUint16(player.MapObjectID)
	globalDiceTable.broadcast(cmd)
}

// handleDiceTaisaiGameLeaveReq maneja PROTO_NC_DICE_TAISAI_GAME_LEAVE_REQ (0xAC05).
func (h *Handler) handleDiceTaisaiGameLeaveReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	globalDiceTable.mu.Lock()
	delete(globalDiceTable.Players, player.MapObjectID)
	globalDiceTable.mu.Unlock()
	player.CasinoDiceJoined = false

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_GAME_LEAVE_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_GAME_LEAVE_CMD)
	cmd.WriteByte(0)
	cmd.WriteUint16(player.MapObjectID)
	globalDiceTable.broadcast(cmd)
}

// handleDiceTaisaiExpectIncomeMoneyReq maneja PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_REQ (0xAC0B).
func (h *Handler) handleDiceTaisaiExpectIncomeMoneyReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	invenSlot, _ := pkt.Reader.ReadUint16()
	var gameCnt int
	var expect int64
	if h.charRepo != nil {
		gameCnt, expect, _ = h.charRepo.GetDiceExpect(player.Char.ID, int64(invenSlot))
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_ACK)
	ack.WriteUint16(0)
	ack.WriteUint16(invenSlot)
	ack.WriteUint32(uint32(gameCnt))
	ack.WriteUint64(uint64(expect))
	_ = client.SendPacket(ack)
}

// handleDiceTaisaiReceiptIncomeMoneyReq maneja PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_REQ (0xAC0F).
func (h *Handler) handleDiceTaisaiReceiptIncomeMoneyReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	invenSlot, _ := pkt.Reader.ReadUint16()
	var expect int64
	if h.charRepo != nil {
		_, expect, _ = h.charRepo.GetDiceExpect(player.Char.ID, int64(invenSlot))
		if expect > 0 {
			_ = h.charRepo.SetDicePay(player.Char.ID, int64(invenSlot), expect)
			player.Char.TCharacterState.Money += expect
			_ = h.charRepo.SaveCharMoney(player.Char.ID, player.Char.TCharacterState.Money)
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_ACK)
	ack.WriteUint16(0)
	ack.WriteUint16(invenSlot)
	ack.WriteUint64(uint64(expect))
	_ = client.SendPacket(ack)
}

// handleDiceTaisaiGameModeChangeReq maneja PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_REQ (0xAC13).
func (h *Handler) handleDiceTaisaiGameModeChangeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	var mode uint32
	if pkt.Reader.Remaining() >= 4 {
		mode, _ = pkt.Reader.ReadUint32()
	} else {
		mode = 1 - globalDiceTable.ModeType
	}
	globalDiceTable.ModeType = mode

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_ACK)
	ack.WriteUint16(0)
	ack.WriteUint32(mode)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_CMD)
	cmd.WriteUint32(mode)
	globalDiceTable.broadcast(cmd)
	_ = player
}

// handleDiceTaisaiRankReq maneja PROTO_NC_DICE_TAISAI_RANK_REQ (0xAC16).
func (h *Handler) handleDiceTaisaiRankReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil {
		return
	}
	var ranks []*struct {
		Name string
		Wins int
	}
	if h.charRepo != nil {
		if rows, err := h.charRepo.GetDiceRank(player.Char.ID, 0); err == nil {
			for _, r := range rows {
				ranks = append(ranks, &struct {
					Name string
					Wins int
				}{Name: r.CharName, Wins: r.WinCount})
			}
		}
	}
	// PROTO_NC_DICE_TAISAI_RANK_ACK (484B + nRankCurrCnt*24B)
	// DiceRankInfo = 24B (Name 20B + uint32 wins) verificado en so_mh_DiceTaiSai_SendRank @ 0x00514060
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_RANK_ACK)
	ack.WriteUint16(0)
	histCnt := len(ranks)
	if histCnt > 20 {
		histCnt = 20
	}
	ack.WriteByte(uint8(histCnt))
	histBuf := make([]byte, 480)
	for i := 0; i < histCnt; i++ {
		r := ranks[i]
		copy(histBuf[i*24:i*24+20], r.Name)
		binary.LittleEndian.PutUint32(histBuf[i*24+20:i*24+24], uint32(r.Wins))
	}
	ack.WriteBytes(histBuf)
	ack.WriteByte(uint8(histCnt))
	for i := 0; i < histCnt; i++ {
		r := ranks[i]
		writeFixedStr(ack, r.Name, 20)
		ack.WriteUint32(uint32(r.Wins))
	}
	_ = client.SendPacket(ack)
}

func resolveTaisaiBet(betInfo uint64) (betType uint8, coins uint64) {
	betType = uint8((betInfo >> 56) & 0xFF)
	coins = betInfo & 0x00FFFFFFFFFFFFFF
	if coins == 0 {
		coins = 10
	}
	return betType, coins
}

// handleDiceTaisaiBettingReq maneja PROTO_NC_DICE_TAISAI_BETTING_REQ (0xAC1A).
func (h *Handler) handleDiceTaisaiBettingReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 10 {
		return
	}
	betInfo, _ := pkt.Reader.ReadUint64()
	betIdx, _ := pkt.Reader.ReadUint16()
	bType, coins := resolveTaisaiBet(betInfo)

	var errCode uint16
	if player.CasinoCoins < coins {
		// Permitir apostar directamente con Coins o fallback a Money en mesa personal
		if uint64(player.Char.TCharacterState.Money) < coins {
			errCode = 1
		} else {
			player.Char.TCharacterState.Money -= int64(coins)
			if h.charRepo != nil {
				_ = h.charRepo.SaveCharMoney(player.Char.ID, player.Char.TCharacterState.Money)
			}
		}
	} else {
		player.CasinoCoins -= coins
		if h.charRepo != nil {
			_ = h.charRepo.SubGambleCoin(player.Char.ID, int64(coins))
		}
	}

	if errCode == 0 {
		key := (uint32(player.Char.ID) << 16) | uint32(betIdx)
		globalDiceTable.mu.Lock()
		globalDiceTable.Bets[key] = diceTaisaiBet{
			CharNo:   player.Char.ID,
			Handle:   player.MapObjectID,
			BetIndex: betIdx,
			BetInfo:  betInfo,
			BetType:  bType,
			BetCoins: coins,
		}
		globalDiceTable.mu.Unlock()
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_BETTING_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint16(betIdx)
	_ = client.SendPacket(ack)

	if errCode == 0 {
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_BETTING_CMD)
		cmd.WriteUint16(player.MapObjectID)
		cmd.WriteUint64(betInfo)
		globalDiceTable.broadcast(cmd)
	}
}

// handleDiceTaisaiBettingCancelReq maneja PROTO_NC_DICE_TAISAI_BETTING_CANCEL_REQ (0xAC1D).
func (h *Handler) handleDiceTaisaiBettingCancelReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 10 {
		return
	}
	betInfo, _ := pkt.Reader.ReadUint64()
	betIdx, _ := pkt.Reader.ReadUint16()

	key := (uint32(player.Char.ID) << 16) | uint32(betIdx)
	globalDiceTable.mu.Lock()
	bet, exists := globalDiceTable.Bets[key]
	if exists {
		delete(globalDiceTable.Bets, key)
	}
	globalDiceTable.mu.Unlock()

	var errCode uint16
	if !exists {
		errCode = 1
	} else {
		player.CasinoCoins += bet.BetCoins
		if h.charRepo != nil {
			_ = h.charRepo.AddGambleCoin(player.Char.ID, 0, int64(bet.BetCoins))
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_BETTING_CANCEL_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint16(betIdx)
	_ = client.SendPacket(ack)

	if errCode == 0 {
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_BETTING_CANCEL_CMD)
		cmd.WriteUint16(player.MapObjectID)
		cmd.WriteUint64(betInfo)
		globalDiceTable.broadcast(cmd)
	}
}

func (h *Handler) rollThreeDiceAndResolve(dept47 bool) [3]byte {
	d := [3]byte{
		byte(rand.Intn(6) + 1),
		byte(rand.Intn(6) + 1),
		byte(rand.Intn(6) + 1),
	}
	sum := int(d[0] + d[1] + d[2])
	isTriple := d[0] == d[1] && d[1] == d[2]
	diceStr := fmt.Sprintf("%d,%d,%d", d[0], d[1], d[2])

	var gameNo int
	if h.charRepo != nil {
		if dept47 {
			gameNo, _ = h.charRepo.SetGambleDice(diceStr)
		} else {
			gameNo, _ = h.charRepo.SetDiceGame(0, 0, diceStr)
		}
	}

	globalDiceTable.mu.Lock()
	copy(globalDiceTable.DiceHistory[:18], globalDiceTable.DiceHistory[3:])
	copy(globalDiceTable.DiceHistory[18:21], d[:])
	bets := make([]diceTaisaiBet, 0, len(globalDiceTable.Bets))
	for k, b := range globalDiceTable.Bets {
		bets = append(bets, b)
		delete(globalDiceTable.Bets, k)
	}
	globalDiceTable.mu.Unlock()

	taxPct := 5
	if h.dataStore != nil && h.dataStore.CasinoCatalog != nil {
		taxPct = h.dataStore.CasinoCatalog.TaxRatePct
	}

	for _, b := range bets {
		won := false
		mult := 2.0
		switch {
		case b.BetType == 0 && !isTriple && sum >= 4 && sum <= 10: // Small
			won = true
		case b.BetType == 1 && !isTriple && sum >= 11 && sum <= 17: // Big
			won = true
		case b.BetType == 2 && isTriple: // Triple
			won = true
			mult = 24.0
		default:
			if h.dataStore != nil && h.dataStore.CasinoCatalog != nil {
				for _, rule := range h.dataStore.CasinoCatalog.DiceGames {
					if rule.BetType == b.BetType {
						mult = rule.RateMult
						won = (sum%2 == int(b.BetType)%2) && !isTriple
						break
					}
				}
			}
		}
		var dividend int64
		if won {
			gross := int64(float64(b.BetCoins) * mult)
			dividend = gross - (gross*int64(taxPct))/100
			if p := h.findPlayerByCharID(b.CharNo); p != nil {
				p.CasinoCoins += uint64(dividend)
				coinCmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_COIN_CHANGE_CMD)
				coinCmd.WriteUint64(p.CasinoCoins)
				if p.Client != nil {
					_ = p.Client.SendPacket(coinCmd)
				}
			}
		}
		if h.charRepo != nil {
			if dept47 {
				_ = h.charRepo.SetGambleDiceLog(gameNo, b.CharNo, int(b.BetIndex), int64(b.BetCoins), int(mult), b.BetType, taxPct, dividend)
				if dividend > 0 {
					_, _ = h.charRepo.SetGambleDiceResult(gameNo, b.CharNo, dividend)
				}
			} else {
				_ = h.charRepo.SetDiceLog(gameNo, b.CharNo, int(b.BetIndex), int64(b.BetCoins), int(mult), b.BetType, taxPct, dividend)
				if dividend > 0 {
					_ = h.charRepo.SetDiceResult(gameNo, b.CharNo, dividend)
				}
			}
		}
		if h.gameLogDispatcher != nil {
			_ = h.gameLogDispatcher.EnqueueTypeI(2, b.CharNo, "", 0, 0, 0, int(b.BetCoins), int(dividend))
		}
	}
	return d
}

// handleDiceTaisaiDiceRollReq maneja PROTO_NC_DICE_TAISAI_DICE_ROLL_REQ (0xAC20).
func (h *Handler) handleDiceTaisaiDiceRollReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_DICE_ROLL_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	d := h.rollThreeDiceAndResolve(false)
	resCmd := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_CMD)
	resCmd.WriteUint16(0)
	resCmd.WriteBytes(d[:])
	globalDiceTable.broadcast(resCmd)
}

// handleDiceTaisaiLargeAmountReq maneja PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_REQ (0xAC28).
func (h *Handler) handleDiceTaisaiLargeAmountReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_CMD)
	cmd.WriteByte(0)
	globalDiceTable.broadcast(cmd)
}

// ============================================================================
// Dept 47: NC_GAMBLE (18 C2S: House, Coin, Taisai Casino & SlotMachine)
// ============================================================================

// handleGambleHouseExitReq maneja PROTO_NC_GAMBLE_GAMBLEHOUSE_EXIT_REQ (0xBC02).
func (h *Handler) handleGambleHouseExitReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_GAMBLEHOUSE_EXIT_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
	if player.Char != nil && player.Char.TownName != "" {
		h.WarpPlayer(player, player.Char.TownName, player.Char.XPos, player.Char.YPos, player.Rotation)
	}
}

// handleGambleCoinVipCardUIOpenReq maneja PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ (0xBC0D).
func (h *Handler) handleGambleCoinVipCardUIOpenReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil {
		return
	}
	var coins, exchanged int64
	if h.charRepo != nil {
		coins, _ = h.charRepo.GetGambleCoin(player.Char.ID)
		exchanged, _ = h.charRepo.GetGambleExchangedCoin(player.Char.ID, time.Now().Add(-24*time.Hour))
	}
	player.CasinoCoins = uint64(coins)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	coinCmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_COIN_CHANGE_CMD)
	coinCmd.WriteUint64(player.CasinoCoins)
	_ = client.SendPacket(coinCmd)

	exCmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_EXCHANGECOIN_CHANGE_CMD)
	exCmd.WriteUint64(uint64(exchanged))
	_ = client.SendPacket(exCmd)
}

// handleGambleCoinBuyReq maneja PROTO_NC_GAMBLE_COIN_BUY_REQ (0xBC0F).
func (h *Handler) handleGambleCoinBuyReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 8 {
		return
	}
	coins, _ := pkt.Reader.ReadUint64()
	costCen := int64(coins * 100) // 1 Coin = 100 Copper

	var errCode uint16
	if coins == 0 || player.Char.TCharacterState.Money < costCen {
		errCode = 1
	} else {
		player.Char.TCharacterState.Money -= costCen
		player.CasinoCoins += coins
		if h.charRepo != nil {
			_ = h.charRepo.SaveCharMoney(player.Char.ID, player.Char.TCharacterState.Money)
			_ = h.charRepo.AddGambleCoin(player.Char.ID, 0, int64(coins))
		}
		if h.gameLogDispatcher != nil {
			_ = h.gameLogDispatcher.EnqueueTypeI(3, player.Char.ID, "", player.Char.XPos, player.Char.YPos, 0, int(coins), int(costCen))
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_COIN_BUY_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)

	if errCode == 0 {
		cCmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_COIN_CHANGE_CMD)
		cCmd.WriteUint64(player.CasinoCoins)
		_ = client.SendPacket(cCmd)
	}
}

// handleGambleCoinSellReq maneja PROTO_NC_GAMBLE_COIN_SELL_REQ (0xBC13).
func (h *Handler) handleGambleCoinSellReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 8 {
		return
	}
	coins, _ := pkt.Reader.ReadUint64()
	gainCen := int64(coins * 100)

	var errCode uint16
	if coins == 0 || player.CasinoCoins < coins {
		errCode = 1
	} else {
		player.CasinoCoins -= coins
		player.Char.TCharacterState.Money += gainCen
		if h.charRepo != nil {
			_ = h.charRepo.SubGambleCoin(player.Char.ID, int64(coins))
			_ = h.charRepo.SaveCharMoney(player.Char.ID, player.Char.TCharacterState.Money)
		}
		if h.gameLogDispatcher != nil {
			_ = h.gameLogDispatcher.EnqueueTypeI(4, player.Char.ID, "", player.Char.XPos, player.Char.YPos, 0, int(coins), int(gainCen))
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_COIN_SELL_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)

	if errCode == 0 {
		cCmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_COIN_CHANGE_CMD)
		cCmd.WriteUint64(player.CasinoCoins)
		_ = client.SendPacket(cCmd)
	}
}

// handleGambleTypeAndWhereStandReq maneja PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ (0xBC17).
func (h *Handler) handleGambleTypeAndWhereStandReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 2 {
		return
	}
	npcHnd, _ := pkt.Reader.ReadUint16()
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_ACK)
	ack.WriteUint16(0)
	ack.WriteUint16(npcHnd)
	ack.WriteUint32(1) // Slot/Taisai
	ack.WriteByte(0)
	_ = client.SendPacket(ack)
}

// handleGambleDiceTaisaiGameJoinReq maneja PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ (0xBC64).
func (h *Handler) handleGambleDiceTaisaiGameJoinReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	globalDiceTable.mu.Lock()
	globalDiceTable.Players[player.MapObjectID] = player
	player.CasinoDiceJoined = true
	cnt := uint8(len(globalDiceTable.Players))
	hnds := make([]uint16, 0, cnt)
	for hnd := range globalDiceTable.Players {
		hnds = append(hnds, hnd)
	}
	globalDiceTable.mu.Unlock()

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_ACK)
	ack.WriteUint16(0)
	ack.WriteBytes(globalDiceTable.DiceHistory[:])
	ack.WriteUint32(globalDiceTable.GameStatus)
	ack.WriteUint32(globalDiceTable.ModeType)
	ack.WriteUint16(30)
	ack.WriteUint32(30000)
	ack.WriteUint16(globalDiceTable.TaxRate)
	ack.WriteUint64(globalDiceTable.MaxBet)
	ack.WriteUint16(1)
	ack.WriteByte(cnt)
	for _, gh := range hnds {
		ack.WriteUint16(gh)
	}
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_CMD)
	cmd.WriteUint16(player.MapObjectID)
	globalDiceTable.broadcast(cmd)
}

// handleGambleDiceTaisaiGameLeaveReq maneja PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_REQ (0xBC68).
func (h *Handler) handleGambleDiceTaisaiGameLeaveReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	globalDiceTable.mu.Lock()
	delete(globalDiceTable.Players, player.MapObjectID)
	globalDiceTable.mu.Unlock()
	player.CasinoDiceJoined = false

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_CMD)
	cmd.WriteByte(0)
	cmd.WriteUint16(player.MapObjectID)
	globalDiceTable.broadcast(cmd)
}

// handleGambleDiceTaisaiAllRankReq maneja PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ (0xBC6B).
func (h *Handler) handleGambleDiceTaisaiAllRankReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK)
	ack.WriteUint16(0)
	if h.charRepo != nil {
		rows, _ := h.charRepo.GetGambleDiceRank()
		if len(rows) > 20 {
			rows = rows[:20]
		}
		ack.WriteByte(uint8(len(rows)))
		for _, r := range rows {
			writeFixedStr(ack, r.CharName, 20)
			ack.WriteUint32(uint32(r.WinCount))
		}
	} else {
		ack.WriteByte(0)
	}
	_ = client.SendPacket(ack)
}

// handleGambleDiceTaisaiCurrRankReq maneja PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ (0xBC6F).
func (h *Handler) handleGambleDiceTaisaiCurrRankReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_ACK)
	ack.WriteUint16(0)
	ack.WriteByte(0)
	_ = client.SendPacket(ack)
}

// handleGambleDiceTaisaiBettingReq maneja PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_REQ (0xBC71).
func (h *Handler) handleGambleDiceTaisaiBettingReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 10 {
		return
	}
	betInfo, _ := pkt.Reader.ReadUint64()
	betIdx, _ := pkt.Reader.ReadUint16()
	bType, coins := resolveTaisaiBet(betInfo)

	var errCode uint16
	if player.CasinoCoins < coins {
		errCode = 1
	} else {
		player.CasinoCoins -= coins
		if h.charRepo != nil {
			_ = h.charRepo.SubGambleCoin(player.Char.ID, int64(coins))
		}
		key := (uint32(player.Char.ID) << 16) | uint32(betIdx)
		globalDiceTable.mu.Lock()
		globalDiceTable.Bets[key] = diceTaisaiBet{
			CharNo:   player.Char.ID,
			Handle:   player.MapObjectID,
			BetIndex: betIdx,
			BetInfo:  betInfo,
			BetType:  bType,
			BetCoins: coins,
		}
		globalDiceTable.mu.Unlock()
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint16(betIdx)
	_ = client.SendPacket(ack)

	if errCode == 0 {
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CMD)
		cmd.WriteUint16(player.MapObjectID)
		cmd.WriteUint64(betInfo)
		globalDiceTable.broadcast(cmd)

		// Resolver tirada de dados de la mesa cuando alguien apuesta
		d := h.rollThreeDiceAndResolve(true)
		res := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_DICE_ROLL_RESULT_CMD)
		res.WriteUint16(0)
		res.WriteBytes(d[:])
		globalDiceTable.broadcast(res)
	}
}

// handleGambleDiceTaisaiBettingCancelReq maneja PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ (0xBC74).
func (h *Handler) handleGambleDiceTaisaiBettingCancelReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 10 {
		return
	}
	betInfo, _ := pkt.Reader.ReadUint64()
	betIdx, _ := pkt.Reader.ReadUint16()

	key := (uint32(player.Char.ID) << 16) | uint32(betIdx)
	globalDiceTable.mu.Lock()
	bet, exists := globalDiceTable.Bets[key]
	if exists {
		delete(globalDiceTable.Bets, key)
	}
	globalDiceTable.mu.Unlock()

	var errCode uint16
	if !exists {
		errCode = 1
	} else {
		player.CasinoCoins += bet.BetCoins
		if h.charRepo != nil {
			_ = h.charRepo.AddGambleCoin(player.Char.ID, 0, int64(bet.BetCoins))
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint16(betIdx)
	_ = client.SendPacket(ack)

	if errCode == 0 {
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_CMD)
		cmd.WriteUint16(player.MapObjectID)
		cmd.WriteUint64(betInfo)
		globalDiceTable.broadcast(cmd)
	}
}

// handleGambleDiceTaisaiLargeAmountReq maneja PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_REQ (0xBC7D).
func (h *Handler) handleGambleDiceTaisaiLargeAmountReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_CMD)
	cmd.WriteByte(0)
	globalDiceTable.broadcast(cmd)
}

// handleGambleSlotMachineJoinReq maneja PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ (0xBCC8).
func (h *Handler) handleGambleSlotMachineJoinReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	var betUnit uint16 = 10
	if h.dataStore != nil && h.dataStore.CasinoCatalog != nil && len(h.dataStore.CasinoCatalog.SlotBetCoins) > 0 {
		betUnit = uint16(h.dataStore.CasinoCatalog.SlotBetCoins[0].Coins)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_ACK)
	ack.WriteUint16(betUnit)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	// Enviar Jackpot actual (PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD 0xBCD8)
	if h.charRepo != nil {
		if jps, err := h.charRepo.GetJackpotCoin("SLOT"); err == nil && len(jps) > 0 {
			jpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTINFO_CMD)
			jpCmd.WriteByte(uint8(len(jps)))
			for _, jp := range jps {
				jpCmd.WriteUint16(uint16(jp.BattingUnit))
				jpCmd.WriteUint32(uint32(jp.JackpotCoin))
			}
			_ = client.SendPacket(jpCmd)
		}
	}
	_ = player
}

// handleGambleSlotMachineLeaveReq maneja PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_REQ (0xBCCA).
func (h *Handler) handleGambleSlotMachineLeaveReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// handleGambleSlotMachineStartReq maneja PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ (0xBCCC).
func (h *Handler) handleGambleSlotMachineStartReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 4 {
		return
	}
	betting, _ := pkt.Reader.ReadUint32()
	if betting == 0 {
		betting = 10
	}
	var errCode uint16
	if player.CasinoCoins < uint64(betting) {
		errCode = 1
	} else {
		player.CasinoCoins -= uint64(betting)
		if h.charRepo != nil {
			_ = h.charRepo.SubGambleCoin(player.Char.ID, int64(betting))
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_START_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

func (h *Handler) evaluateSlotSpin(player *ZonePlayer, betUnit uint16) ([9]byte, uint32, uint16, bool) {
	var screen [9]byte
	for i := 0; i < 9; i++ {
		screen[i] = byte(rand.Intn(7))
	}
	var winCoins uint32
	var winLine uint16
	isJackpot := false

	// Filas horizontales (0-2, 3-5, 6-8)
	for row := 0; row < 3; row++ {
		base := row * 3
		if screen[base] == screen[base+1] && screen[base+1] == screen[base+2] {
			winLine |= 1 << uint(row)
			mult := uint32(10)
			if h.dataStore != nil && h.dataStore.CasinoCatalog != nil {
				if row == 1 && h.dataStore.CasinoCatalog.SlotCenterRate > 0 {
					mult = h.dataStore.CasinoCatalog.SlotCenterRate
				} else if h.dataStore.CasinoCatalog.SlotLineRate > 0 {
					mult = h.dataStore.CasinoCatalog.SlotLineRate
				}
			}
			winCoins += uint32(betUnit) * mult
			if screen[base] == 6 {
				isJackpot = true
			}
		}
	}

	if winCoins > 0 && player.Char != nil {
		player.CasinoCoins += uint64(winCoins)
		if h.charRepo != nil {
			_ = h.charRepo.AddGambleCoin(player.Char.ID, 0, int64(winCoins))
		}
	}
	if h.charRepo != nil && player.Char != nil {
		jpRate := 5
		if h.dataStore != nil && h.dataStore.CasinoCatalog != nil {
			jpRate = int(h.dataStore.CasinoCatalog.JackpotBasePct)
		}
		_ = h.charRepo.SetSlotMachineResult(player.Char.ID, int64(betUnit), byte(winLine),
			fmt.Sprintf("%v", screen), winCoins > 0, int64(winCoins), isJackpot, jpRate)
	}
	if h.gameLogDispatcher != nil && player.Char != nil {
		_ = h.gameLogDispatcher.EnqueueTypeI(5, player.Char.ID, "", player.Char.XPos, player.Char.YPos, 0, int(betUnit), int(winCoins))
	}
	return screen, winCoins, winLine, isJackpot
}

// handleGambleSlotMachineStopButtonReq maneja PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_REQ (0xBCCE).
func (h *Handler) handleGambleSlotMachineStopButtonReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	screen, winCoins, winLine, isJackpot := h.evaluateSlotSpin(player, 10)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK)
	ack.WriteUint16(0)
	ack.WriteBytes(screen[:])
	_ = client.SendPacket(ack)

	// Si hubo premio o jackpot emitir también WHEELSTOP_ACK y SOMEONE_GET_JACKPOT_CMD
	if winCoins > 0 {
		wAck := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK)
		wAck.WriteUint16(0)
		wAck.WriteUint16(10)
		wAck.WriteUint32(winCoins)
		wAck.WriteUint16(winLine)
		_ = client.SendPacket(wAck)

		cCmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_COIN_CHANGE_CMD)
		cCmd.WriteUint64(player.CasinoCoins)
		_ = client.SendPacket(cCmd)

		if isJackpot && player.CurrentMap != nil {
			jpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_SOMEONE_GET_JACKPOT_CMD)
			writeFixedStr(jpCmd, player.Char.Name, 20)
			jpCmd.WriteUint64(uint64(winCoins))
			player.CurrentMap.Broadcast(jpCmd, 0)
		}
	}
}

// handleGambleSlotMachineWheelStopReq maneja PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_REQ (0xBCD0).
func (h *Handler) handleGambleSlotMachineWheelStopReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	_, winCoins, winLine, _ := h.evaluateSlotSpin(player, 10)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK)
	ack.WriteUint16(0)
	ack.WriteUint16(10)
	ack.WriteUint32(winCoins)
	ack.WriteUint16(winLine)
	_ = client.SendPacket(ack)
}

// handleGambleSlotMachineJackpotRankingReq maneja PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_REQ (0xBCD9).
// SLOTMACHINE_JackPotRank = 24B (Charid 20B + WinCoin uint32) verificado en smjpr_PacketMake @ 0x0061FF80 y 0x0061EC60 (3 + 24*N).
func (h *Handler) handleGambleSlotMachineJackpotRankingReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK)
	if h.charRepo != nil {
		rows, _ := h.charRepo.GetJackpotRanking("SLOT", 10)
		if len(rows) > 50 {
			rows = rows[:50]
		}
		ack.WriteByte(uint8(len(rows)))
		for _, r := range rows {
			writeFixedStr(ack, r.CharName, 20)
			ack.WriteUint32(uint32(r.JackpotCoin))
		}
	} else {
		ack.WriteByte(0)
	}
	_ = client.SendPacket(ack)
}

// handleGambleSlotMachineWinRankingReq maneja PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_REQ (0xBCDB).
// SLOTMACHINE_JackPotRank = 24B (Charid 20B + WinCoin uint32) verificado en go_SMS_PacketMake_WintRank @ 0x00624520 y 0x0061ED90 (3 + 24*N).
func (h *Handler) handleGambleSlotMachineWinRankingReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK)
	if h.charRepo != nil {
		rows, _ := h.charRepo.GetJackpotRanking("SLOT", 0)
		if len(rows) > 50 {
			rows = rows[:50]
		}
		ack.WriteByte(uint8(len(rows)))
		for _, r := range rows {
			writeFixedStr(ack, r.CharName, 20)
			ack.WriteUint32(uint32(r.JackpotCoin))
		}
	} else {
		ack.WriteByte(0)
	}
	_ = client.SendPacket(ack)
}
