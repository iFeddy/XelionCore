package zone

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	glua "github.com/yuin/gopher-lua"
)

// miniHouseState mantiene el estado vivo en memoria de la MiniHouse de un dueño.
type miniHouseState struct {
	OwnerCharNo   int
	OwnerName     string
	HouseHandle   uint16
	HouseID       uint16
	Level         uint8
	MaxGuests     uint8
	Password      string
	Title         string
	Notify        string
	ItemInfoClose uint8
	ArrangeMode   bool
	PortalOpen    bool
	Active        bool
	Guests        map[uint16]bool
	Kicked        map[uint32]bool // charNo
	PrevLocation  map[uint32]struct {
		MapID string
		X, Y  int
	}
	FurnitureLocs map[uint16][16]byte // furnitureHandle (slot/key low16) -> 16-byte transform
}

var (
	miniHousesMu sync.RWMutex
	miniHouses   = make(map[int]*miniHouseState) // OwnerCharNo -> state
)

func (h *Handler) getOrCreateMiniHouse(owner *ZonePlayer) *miniHouseState {
	miniHousesMu.Lock()
	defer miniHousesMu.Unlock()
	if st, ok := miniHouses[owner.Char.ID]; ok {
		return st
	}
	maxG := owner.MiniHouseMaxGuests
	if maxG == 0 {
		maxG = 10
	}
	title := owner.MiniHouseTitle
	if title == "" {
		title = fmt.Sprintf("%s's House", owner.Char.Name)
	}
	st := &miniHouseState{
		OwnerCharNo:   owner.Char.ID,
		OwnerName:     owner.Char.Name,
		HouseHandle:   owner.MapObjectID,
		HouseID:       1,
		Level:         1,
		MaxGuests:     maxG,
		Password:      owner.MiniHousePassword,
		Title:         title,
		Guests:        make(map[uint16]bool),
		Kicked:        make(map[uint32]bool),
		PrevLocation: make(map[uint32]struct {
			MapID string
			X, Y  int
		}),
		FurnitureLocs: make(map[uint16][16]byte),
	}
	miniHouses[owner.Char.ID] = st
	return st
}

func (h *Handler) findMiniHouseByHandle(mhHnd uint16) *miniHouseState {
	miniHousesMu.RLock()
	defer miniHousesMu.RUnlock()
	for _, st := range miniHouses {
		if st.HouseHandle == mhHnd && st.Active {
			return st
		}
	}
	return nil
}

func (h *Handler) currentMiniHouseOf(player *ZonePlayer) *miniHouseState {
	miniHousesMu.RLock()
	defer miniHousesMu.RUnlock()
	if player.InsideMiniHouseOwner > 0 {
		if st, ok := miniHouses[player.InsideMiniHouseOwner]; ok {
			return st
		}
	}
	if st, ok := miniHouses[player.Char.ID]; ok && st.Active {
		return st
	}
	return nil
}

func (h *Handler) broadcastMiniHouse(st *miniHouseState, pkt *protocol.PacketWriter) {
	if st == nil {
		return
	}
	if owner := h.findPlayerByCharID(st.OwnerCharNo); owner != nil && owner.Client != nil {
		_ = owner.Client.SendPacket(pkt)
	}
	miniHousesMu.RLock()
	guestHnds := make([]uint16, 0, len(st.Guests))
	for gh := range st.Guests {
		guestHnds = append(guestHnds, gh)
	}
	miniHousesMu.RUnlock()
	for _, gh := range guestHnds {
		if g := h.findPlayerByHandle(gh); g != nil && g.Char.ID != st.OwnerCharNo && g.Client != nil {
			_ = g.Client.SendPacket(pkt)
		}
	}
}

func writeFixedStr(pw *protocol.PacketWriter, s string, size int) {
	buf := make([]byte, size)
	copy(buf, s)
	pw.WriteBytes(buf)
}

func readFixedStr(pr *protocol.PacketReader, size int) string {
	b, _ := pr.ReadBytes(size)
	return strings.TrimRight(string(b), "\x00")
}

// ============================================================================
// Dept 35: NC_MINIHOUSE (23+ C2S Handlers + Box 12 + tMinihouseGate + GameLog F)
// ============================================================================

// handleMiniHouseActivReq maneja PROTO_NC_MINIHOUSE_ACTIV_REQ (0x8C01).
func (h *Handler) handleMiniHouseActivReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	slot, _ := pkt.Reader.ReadByte()
	var mhDef *MiniHouseDef
	if h.dataStore != nil && h.dataStore.MiniHouseCatalog != nil && h.charRepo != nil {
		if items, err := h.charRepo.GetItemsByStorage(player.Char.ID, 9); err == nil {
			for _, it := range items {
				if it.Storage == int16(slot) {
					mhDef = h.dataStore.MiniHouseCatalog.ByItemID[uint16(it.ItemID)]
					break
				}
			}
		}
	}
	st := h.getOrCreateMiniHouse(player)
	if mhDef != nil {
		st.HouseID = mhDef.Handle
		if mhDef.Slots > 0 {
			st.MaxGuests = uint8(mhDef.Slots)
		}
	}
	st.Active = true
	st.HouseHandle = player.MapObjectID

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_ACTIV_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	if h.gameLogDispatcher != nil {
		mapID := ""
		if player.CurrentMap != nil {
			mapID = player.CurrentMap.MapID
		}
		h.gameLogDispatcher.EnqueueTypeF(1, player.Char.ID, mapID, player.Char.XPos, player.Char.YPos, 0, int(st.HouseID), 0, int(slot), 1, int(st.MaxGuests), 0)
	}
}

// handleMiniHouseBuildingReq maneja PROTO_NC_MINIHOUSE_BUILDING_REQ (0x8C03).
func (h *Handler) handleMiniHouseBuildingReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 32 {
		return
	}
	maxVis, _ := pkt.Reader.ReadByte()
	pw := readFixedStr(pkt.Reader, 9)
	title := readFixedStr(pkt.Reader, 21)
	bInfoClose, _ := pkt.Reader.ReadByte()
	notify := ""
	if pkt.Reader.Remaining() >= 101 {
		notify = readFixedStr(pkt.Reader, 101)
	} else if pkt.Reader.Remaining() > 0 {
		rem, _ := pkt.Reader.ReadBytes(pkt.Reader.Remaining())
		notify = strings.TrimRight(string(rem), "\x00")
	}

	st := h.getOrCreateMiniHouse(player)
	if maxVis > 0 {
		st.MaxGuests = maxVis
		player.MiniHouseMaxGuests = maxVis
	}
	st.Password = pw
	player.MiniHousePassword = pw
	player.MiniHouseOpen = (pw == "")
	if title != "" {
		st.Title = title
		player.MiniHouseTitle = title
	}
	st.ItemInfoClose = bInfoClose
	st.Notify = notify
	st.Active = true
	st.HouseHandle = player.MapObjectID

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_BUILDING_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	// PROTO_NC_MINIHOUSE_BUILDING_CMD (66 bytes: player(2) + PROTO_NC_BRIEFINFO_MINIHOUSEBUILD_CMD(64))
	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_BUILDING_CMD)
	cmd.WriteUint16(player.MapObjectID)
	cmd.WriteUint16(st.HouseHandle)
	// CHARBRIEFINFO_CAMP (12B): minihouse(u16) + dummy[10] (0xFF como en sp_MakeTentBriefInfo)
	cmd.WriteUint16(st.HouseID)
	cmd.Fill(10, 0xFF)
	// Name5 charid (20B)
	writeFixedStr(cmd, player.Char.Name, 20)
	// SHINE_COORD_TYPE coord (9B: x u32, y u32, dir u8)
	cmd.WriteUint32(uint32(player.Char.XPos))
	cmd.WriteUint32(uint32(player.Char.YPos))
	cmd.WriteByte(player.Rotation)
	// title[21]
	writeFixedStr(cmd, st.Title, 21)
	if player.CurrentMap != nil {
		player.CurrentMap.Broadcast(cmd, 0)
	} else {
		_ = client.SendPacket(cmd)
	}

	if h.gameLogDispatcher != nil {
		mapID := ""
		if player.CurrentMap != nil {
			mapID = player.CurrentMap.MapID
		}
		h.gameLogDispatcher.EnqueueTypeF(2, player.Char.ID, mapID, player.Char.XPos, player.Char.YPos, 0, int(st.HouseID), 0, int(st.MaxGuests), int(bInfoClose), 0, 0)
	}
}

// handleMiniHouseLoginCompleteReq maneja PROTO_NC_MINIHOUSE_LOGINCOMPLETE_REQ (0x8C0A).
func (h *Handler) handleMiniHouseLoginCompleteReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil {
		return
	}
	st := h.currentMiniHouseOf(player)
	if st == nil {
		st = h.getOrCreateMiniHouse(player)
	}
	st.Guests[player.MapObjectID] = true
	player.InsideMiniHouseOwner = st.OwnerCharNo

	if h.charRepo != nil && player.Char.ID != st.OwnerCharNo {
		cnt, _ := h.charRepo.GetMiniHouseVisitors(st.OwnerCharNo)
		_ = h.charRepo.SetMiniHouseVisitors(st.OwnerCharNo, cnt+1)
	}

	// PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK (140B header + guests*3B)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_LOGINCOMPLETE_ACK)
	ack.WriteInt32(0)
	ack.WriteByte(st.MaxGuests)
	ack.WriteByte(st.ItemInfoClose)
	writeFixedStr(ack, st.Notify, 101)
	writeFixedStr(ack, fmt.Sprintf("MH_%d", st.OwnerCharNo), 32)
	ack.WriteByte(uint8(len(st.Guests)))
	for gh := range st.Guests {
		ack.WriteByte(1)
		ack.WriteUint16(gh)
	}
	_ = client.SendPacket(ack)

	// Enviar muebles colocados en Box 12 (StorageType = 12)
	if h.charRepo != nil {
		if items, err := h.charRepo.GetItemsByStorage(st.OwnerCharNo, 12); err == nil {
			furPkt := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_FUNICHER_CMD)
			furPkt.WriteByte(uint8(len(items)))
			furPkt.WriteByte(0) // Next
			for _, it := range items {
				hnd := uint16(it.ItemKey & 0xFFFF)
				if hnd == 0 {
					hnd = uint16(it.Storage + 1)
				}
				furPkt.WriteUint16(hnd)
				furPkt.WriteUint16(uint16(it.ItemID))
				if loc, okLoc := st.FurnitureLocs[hnd]; okLoc {
					furPkt.WriteBytes(loc[:])
				} else {
					var loc [16]byte
					binary.LittleEndian.PutUint16(loc[0:2], uint16(it.Storage))
					furPkt.WriteBytes(loc[:])
				}
				furPkt.WriteByte(100) // NEndureGrade
			}
			_ = client.SendPacket(furPkt)
		}

		// Enviar portales activos de tMinihouseGate
		if gates, err := h.charRepo.GetMiniHouseGates(st.OwnerCharNo); err == nil {
			portCmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_PORTAL_LIST_CMD)
			portCmd.WriteUint16(uint16(len(gates)))
			for _, g := range gates {
				portCmd.WriteUint16(uint16(g.ItemID))
				portCmd.WriteInt32(int32(g.EndDate.Unix()))
			}
			_ = client.SendPacket(portCmd)
		}
	}

	// Notificar entrada al resto en la casa
	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_LOGINCOMPLETE_CMD)
	cmd.WriteByte(player.Char.Level)
	cmd.WriteUint16(player.MapObjectID)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseVisitReadyReq maneja PROTO_NC_MINIHOUSE_VISITREADY_REQ (0x8C0E).
func (h *Handler) handleMiniHouseVisitReadyReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 2 {
		return
	}
	mhHnd, _ := pkt.Reader.ReadUint16()
	st := h.findMiniHouseByHandle(mhHnd)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_VISITREADY_ACK)
	if st == nil {
		ack.WriteInt32(1) // no encontrado
		ack.WriteByte(0)
		ack.WriteByte(0)
		ack.WriteByte(0)
		_ = client.SendPacket(ack)
		return
	}
	if st.Kicked[uint32(player.Char.ID)] {
		ack.WriteInt32(2) // expulsado
		ack.WriteByte(0)
		ack.WriteByte(uint8(len(st.Guests)))
		ack.WriteByte(st.MaxGuests)
		_ = client.SendPacket(ack)
		return
	}
	ack.WriteInt32(0)
	if st.Password == "" {
		ack.WriteByte(1)
	} else {
		ack.WriteByte(0)
	}
	ack.WriteByte(uint8(len(st.Guests)))
	ack.WriteByte(st.MaxGuests)
	_ = client.SendPacket(ack)
}

// handleMiniHouseVisitReq maneja PROTO_NC_MINIHOUSE_VISIT_REQ (0x8C10).
func (h *Handler) handleMiniHouseVisitReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 11 {
		return
	}
	mhHnd, _ := pkt.Reader.ReadUint16()
	pw := readFixedStr(pkt.Reader, 9)
	st := h.findMiniHouseByHandle(mhHnd)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_VISIT_ACK)
	if st == nil || !st.Active {
		ack.WriteUint16(0x1187) // MiniHouse not found / inactive
		_ = client.SendPacket(ack)
		return
	}
	if player.InsideMiniHouseOwner != 0 {
		ack.WriteUint16(0x1185) // Already inside MiniHouse
		_ = client.SendPacket(ack)
		return
	}
	if st.Kicked[uint32(player.Char.ID)] && player.Char.ID != st.OwnerCharNo {
		ack.WriteUint16(0x118d) // Kicked out
		_ = client.SendPacket(ack)
		return
	}
	if st.Password != "" && pw != st.Password && player.Char.ID != st.OwnerCharNo {
		ack.WriteUint16(0x1182)
		_ = client.SendPacket(ack)
		return
	}
	if uint8(len(st.Guests)) >= st.MaxGuests && player.Char.ID != st.OwnerCharNo {
		ack.WriteUint16(0x1195)
		_ = client.SendPacket(ack)
		return
	}

	if player.CurrentMap != nil {
		st.PrevLocation[uint32(player.Char.ID)] = struct {
			MapID string
			X, Y  int
		}{MapID: player.CurrentMap.MapID, X: player.Char.XPos, Y: player.Char.YPos}
	}
	st.Guests[player.MapObjectID] = true
	player.InsideMiniHouseOwner = st.OwnerCharNo

	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	if h.gameLogDispatcher != nil {
		mapID := ""
		if player.CurrentMap != nil {
			mapID = player.CurrentMap.MapID
		}
		h.gameLogDispatcher.EnqueueTypeF(3, player.Char.ID, mapID, player.Char.XPos, player.Char.YPos, 0, st.OwnerCharNo, 0, len(st.Guests), 0, 0, 0)
	}
}

// handleMiniHouseKickoutReq maneja PROTO_NC_MINIHOUSE_KICKOUT_REQ (0x8C12).
func (h *Handler) handleMiniHouseKickoutReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	targetHnd, _ := pkt.Reader.ReadUint16()
	st := h.getOrCreateMiniHouse(player)

	var errCode uint16
	target := h.findPlayerByHandle(targetHnd)
	if st.OwnerCharNo != player.Char.ID {
		errCode = 0x1188 // not owner
	} else if targetHnd == player.MapObjectID {
		errCode = 0x1191 // self-kick invalid
	} else if target == nil || target.Char == nil {
		errCode = 0x1194 // target invalid
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_KICKOUT_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)

	if errCode == 0 && target != nil && target.Char != nil {
		st.Kicked[uint32(target.Char.ID)] = true
		delete(st.Guests, targetHnd)
		target.InsideMiniHouseOwner = 0
		if target.Client != nil {
			kCmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_KICKOUT_CMD)
			kCmd.WriteByte(1)
			_ = target.Client.SendPacket(kCmd)
		}
		if h.gameLogDispatcher != nil {
			h.gameLogDispatcher.EnqueueTypeF(4, player.Char.ID, "", player.Char.XPos, player.Char.YPos, 0, target.Char.ID, 0, int(targetHnd), 0, 0, 0)
		}
	}
}

// handleMiniHouseKickoutCancelReq maneja PROTO_NC_MINIHOUSE_KICKOUTCANCEL_REQ (0x8C15).
func (h *Handler) handleMiniHouseKickoutCancelReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	targetHnd, _ := pkt.Reader.ReadUint16()
	st := h.getOrCreateMiniHouse(player)
	if target := h.findPlayerByHandle(targetHnd); target != nil && target.Char != nil {
		delete(st.Kicked, uint32(target.Char.ID))
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_KICKOUTCANCEL_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// handleMiniHouseExitReq maneja PROTO_NC_MINIHOUSE_EXIT_REQ (0x8C18).
func (h *Handler) handleMiniHouseExitReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil {
		return
	}
	st := h.currentMiniHouseOf(player)
	if st != nil {
		delete(st.Guests, player.MapObjectID)
		outCmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_LOGOUTCOMPLETE_CMD)
		outCmd.WriteUint16(player.MapObjectID)
		h.broadcastMiniHouse(st, outCmd)
		if prev, okPrev := st.PrevLocation[uint32(player.Char.ID)]; okPrev && prev.MapID != "" {
			if player.CurrentMap != nil && player.CurrentMap.MapID != prev.MapID {
				h.WarpPlayer(player, prev.MapID, prev.X, prev.Y, player.Rotation)
			}
		}
	}
	player.InsideMiniHouseOwner = 0
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_EXIT_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// handleMiniHouseArrangeModeReq maneja PROTO_NC_MINIHOUSE_ARRANGEMODE_REQ (0x8C1A).
func (h *Handler) handleMiniHouseArrangeModeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 1 {
		return
	}
	mode, _ := pkt.Reader.ReadByte()
	st := h.getOrCreateMiniHouse(player)
	st.ArrangeMode = mode != 0

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_ARRANGEMODE_ACK)
	ack.WriteByte(mode)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_ARRANGEMODE_CMD)
	cmd.WriteByte(mode)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseFurnitureInfoCompleteReq maneja PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ (0x8C1D).
func (h *Handler) handleMiniHouseFurnitureInfoCompleteReq(client *network.Client, pkt *protocol.Packet) {
	_, ok := h.getPlayer(client)
	if !ok {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_ACK)
	ack.WriteUint16(0) // so_SendErrorCode(35, 0x1e, err) -> uint16
	_ = client.SendPacket(ack)
}

// handleMiniHouseRearrangeReq maneja PROTO_NC_MINIHOUSE_REARRANGE_REQ (0x8C1F).
func (h *Handler) handleMiniHouseRearrangeReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	cnt, _ := pkt.Reader.ReadByte()
	st := h.getOrCreateMiniHouse(player)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_REARRANGE_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_REARRANGE_CMD)
	cmd.WriteByte(cnt)
	for i := 0; i < int(cnt) && pkt.Reader.Remaining() >= 18; i++ {
		fHnd, _ := pkt.Reader.ReadUint16()
		locBytes, _ := pkt.Reader.ReadBytes(16)
		var loc [16]byte
		copy(loc[:], locBytes)
		st.FurnitureLocs[fHnd] = loc
		cmd.WriteUint16(fHnd)
		cmd.WriteBytes(loc[:])
	}
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseCreateFurnitureReq maneja PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ (0x8C22).
func (h *Handler) handleMiniHouseCreateFurnitureReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	cnt, _ := pkt.Reader.ReadByte()
	st := h.getOrCreateMiniHouse(player)

	type placedItem struct {
		hnd    uint16
		itemID uint16
		loc    [16]byte
	}
	var placed []placedItem

	var box9 []*database.Item
	if h.charRepo != nil {
		box9, _ = h.charRepo.GetItemsByStorage(player.Char.ID, 9)
	}

	for i := 0; i < int(cnt) && pkt.Reader.Remaining() >= 18; i++ {
		fHnd, _ := pkt.Reader.ReadUint16()
		locBytes, _ := pkt.Reader.ReadBytes(16)
		var loc [16]byte
		copy(loc[:], locBytes)

		var furItemID uint16 = fHnd
		for _, it := range box9 {
			if uint16(it.Storage) == fHnd || uint16(it.ItemKey&0xFFFF) == fHnd {
				furItemID = uint16(it.ItemID)
				if h.charRepo != nil {
					_ = h.charRepo.SetItemStorage(it.ItemKey, 12, it.Storage)
				}
				break
			}
		}
		st.FurnitureLocs[fHnd] = loc
		placed = append(placed, placedItem{hnd: fHnd, itemID: furItemID, loc: loc})

		if h.gameLogDispatcher != nil {
			h.gameLogDispatcher.EnqueueTypeF(5, player.Char.ID, "", player.Char.XPos, player.Char.YPos, 0, int(furItemID), int64(fHnd), 1, 0, 0, 0)
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_CREATE_FURNITURE_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	// PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD (1B + N*21B)
	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_CREATE_FURNITURE_CMD)
	cmd.WriteByte(uint8(len(placed)))
	for _, p := range placed {
		cmd.WriteUint16(p.hnd)
		cmd.WriteUint16(p.itemID)
		cmd.WriteBytes(p.loc[:])
		cmd.WriteByte(100)
	}
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseDeleteFurnitureReq maneja PROTO_NC_MINIHOUSE_DELETE_FURNITURE_REQ (0x8C25).
func (h *Handler) handleMiniHouseDeleteFurnitureReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	cnt, _ := pkt.Reader.ReadByte()
	st := h.getOrCreateMiniHouse(player)

	var box12 []*database.Item
	if h.charRepo != nil {
		box12, _ = h.charRepo.GetItemsByStorage(player.Char.ID, 12)
	}

	var deleted []uint16
	for i := 0; i < int(cnt) && pkt.Reader.Remaining() >= 2; i++ {
		fHnd, _ := pkt.Reader.ReadUint16()
		delete(st.FurnitureLocs, fHnd)
		for _, it := range box12 {
			if uint16(it.Storage) == fHnd || uint16(it.ItemKey&0xFFFF) == fHnd {
				if h.charRepo != nil {
					_ = h.charRepo.SetItemStorage(it.ItemKey, 9, it.Storage)
				}
				break
			}
		}
		deleted = append(deleted, fHnd)
		if h.gameLogDispatcher != nil {
			h.gameLogDispatcher.EnqueueTypeF(6, player.Char.ID, "", player.Char.XPos, player.Char.YPos, 0, int(fHnd), 0, 1, 0, 0, 0)
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_DELETE_FURNITURE_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_DELETE_FURNITURE_CMD)
	cmd.WriteByte(uint8(len(deleted)))
	for _, d := range deleted {
		cmd.WriteUint16(d)
	}
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseCompulsionMoveToReq maneja PROTO_NC_MINIHOUSE_COMPULSIONMOVETO_REQ (0x8C28).
func (h *Handler) handleMiniHouseCompulsionMoveToReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 10 {
		return
	}
	hnd, _ := pkt.Reader.ReadUint16()
	loc, _ := pkt.Reader.ReadUint64()
	st := h.getOrCreateMiniHouse(player)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_COMPULSIONMOVETO_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_COMPULSIONMOVETO_CMD)
	cmd.WriteUint16(hnd)
	cmd.WriteUint64(loc)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseModifyPasswordReq maneja PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ (0x8C2B).
func (h *Handler) handleMiniHouseModifyPasswordReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 9 {
		return
	}
	pw := readFixedStr(pkt.Reader, 9)
	st := h.getOrCreateMiniHouse(player)
	var errCode uint16
	if player.Char == nil || st.OwnerCharNo != player.Char.ID {
		errCode = 0x1188 // not owner
	} else {
		st.Password = pw
		player.MiniHousePassword = pw
		player.MiniHouseOpen = (pw == "")
		if pw == "" {
			errCode = 0x11d6 // password cleared
		} else {
			errCode = 0x11d7 // password set
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_CMD)
	if st.Password == "" {
		cmd.WriteByte(1)
	} else {
		cmd.WriteByte(0)
	}
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseModifyTitleReq maneja PROTO_NC_MINIHOUSE_MODIFY_TITLE_REQ (0x8C2E).
func (h *Handler) handleMiniHouseModifyTitleReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 21 {
		return
	}
	title := readFixedStr(pkt.Reader, 21)
	st := h.getOrCreateMiniHouse(player)
	st.Title = title
	player.MiniHouseTitle = title

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_TITLE_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_TITLE_CMD)
	writeFixedStr(cmd, title, 21)
	h.broadcastMiniHouse(st, cmd)

	outCmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_OUTSIDE_TITLE_CMD)
	outCmd.WriteUint16(st.HouseHandle)
	writeFixedStr(outCmd, title, 21)
	if player.CurrentMap != nil {
		player.CurrentMap.Broadcast(outCmd, 0)
	}
}

// handleMiniHouseModifyMaxEnterNumReq maneja PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ (0x8C32).
func (h *Handler) handleMiniHouseModifyMaxEnterNumReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 1 {
		return
	}
	maxVis, _ := pkt.Reader.ReadByte()
	st := h.getOrCreateMiniHouse(player)
	st.MaxGuests = maxVis
	player.MiniHouseMaxGuests = maxVis

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_ACK)
	ack.WriteUint16(0) // so_SendErrorCode(35, 0x33, err) -> uint16
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_CMD)
	cmd.WriteByte(maxVis)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseModifyOwnerBlogReq maneja PROTO_NC_MINIHOUSE_MODIFY_OWNERBLOG_REQ (0x8C35).
func (h *Handler) handleMiniHouseModifyOwnerBlogReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil {
		return
	}
	raw, _ := pkt.Reader.ReadBytes(pkt.Reader.Remaining())
	url := strings.TrimRight(string(raw), "\x00")
	if h.charRepo != nil {
		_ = h.charRepo.SetMiniHouseURL(player.Char.ID, url)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_OWNERBLOG_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// handleMiniHouseOwnerBlogReq maneja PROTO_NC_MINIHOUSE_OWNERBLOG_REQ (0x8C37).
func (h *Handler) handleMiniHouseOwnerBlogReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil {
		return
	}
	url := ""
	if h.charRepo != nil {
		url, _ = h.charRepo.GetMiniHouseURL(player.Char.ID)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_OWNERBLOG_ACK)
	ack.WriteUint16(0)
	writeFixedStr(ack, url, 128)
	_ = client.SendPacket(ack)
}

// handleMiniHouseCharActionReq maneja PROTO_NC_MINIHOUSE_CHAR_ACTION_REQ (0x8C3D).
func (h *Handler) handleMiniHouseCharActionReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 4 {
		return
	}
	hnd, _ := pkt.Reader.ReadUint16()
	furHnd, _ := pkt.Reader.ReadUint16()
	st := h.currentMiniHouseOf(player)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_CHAR_ACTION_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_CHAR_ACTION_CMD)
	cmd.WriteUint16(player.MapObjectID)
	cmd.WriteUint16(hnd)
	cmd.WriteUint16(furHnd)
	cmd.WriteUint16(player.MapObjectID)
	cmd.WriteByte(1)
	cmd.WriteByte(1)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseFurnitureEffectReq maneja PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_REQ (0x8C40).
func (h *Handler) handleMiniHouseFurnitureEffectReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 4 {
		return
	}
	furHnd, _ := pkt.Reader.ReadUint16()
	effHnd, _ := pkt.Reader.ReadUint16()
	st := h.currentMiniHouseOf(player)

	if h.charRepo != nil {
		_ = h.charRepo.DeleteMiniHouseEffectResource(player.Char.ID, int64(furHnd), 0, 0, 0)
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_ACK)
	ack.WriteUint16(effHnd)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_CMD)
	cmd.WriteUint16(effHnd)
	cmd.WriteUint16(1)
	cmd.WriteUint16(player.MapObjectID)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseModifyItemInfoOpenReq maneja PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ (0x8C43).
func (h *Handler) handleMiniHouseModifyItemInfoOpenReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 1 {
		return
	}
	flag, _ := pkt.Reader.ReadByte()
	st := h.getOrCreateMiniHouse(player)
	st.ItemInfoClose = flag

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_ACK)
	ack.WriteUint16(0) // so_SendErrorCode(35, 0x44, err) -> uint16
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_CMD)
	cmd.WriteByte(flag)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseModifyNotifyReq maneja PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ (0x8C48).
func (h *Handler) handleMiniHouseModifyNotifyReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 1 {
		return
	}
	nLen, _ := pkt.Reader.ReadByte()
	msgBytes, _ := pkt.Reader.ReadBytes(int(math.Min(float64(nLen), float64(pkt.Reader.Remaining()))))
	notify := strings.TrimRight(string(msgBytes), "\x00")
	st := h.getOrCreateMiniHouse(player)
	st.Notify = notify

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_ACK)
	ack.WriteUint16(0) // so_SendErrorCode(35, 0x49, err) -> uint16
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD)
	cmd.WriteByte(uint8(len(notify)))
	cmd.WriteBytes([]byte(notify))
	cmd.WriteByte(0) // trailing NUL en Zone.exe (0x5738c3)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHousePortalDelReq maneja PROTO_NC_MINIHOUSE_PORTAL_DEL_REQ (0x8C4D).
func (h *Handler) handleMiniHousePortalDelReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	portalItemID, _ := pkt.Reader.ReadUint16()
	st := h.getOrCreateMiniHouse(player)
	if h.charRepo != nil {
		_ = h.charRepo.DelMiniHouseGate(player.Char.ID, int(portalItemID))
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_PORTAL_DEL_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_PORTAL_DEL_CMD)
	cmd.WriteUint16(portalItemID)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHousePortalOpenReq maneja PROTO_NC_MINIHOUSE_PORTAL_OPEN_REQ (0x8C50).
func (h *Handler) handleMiniHousePortalOpenReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	portalItemID, _ := pkt.Reader.ReadUint16()
	st := h.getOrCreateMiniHouse(player)
	st.PortalOpen = true
	if h.charRepo != nil {
		_ = h.charRepo.SetMiniHouseGate(player.Char.ID, int(portalItemID), 1, time.Now().Add(24*time.Hour), int64(portalItemID))
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_PORTAL_OPEN_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// handleMiniHousePortalCloseReq maneja PROTO_NC_MINIHOUSE_PORTAL_CLOSE_REQ (0x8C53).
func (h *Handler) handleMiniHousePortalCloseReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	st := h.getOrCreateMiniHouse(player)
	st.PortalOpen = false
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_PORTAL_CLOSE_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// handleMiniHousePortalEffectReq maneja PROTO_NC_MINIHOUSE_PORTAL_EFFECT_REQ (0x8C56).
func (h *Handler) handleMiniHousePortalEffectReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	st := h.currentMiniHouseOf(player)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_PORTAL_EFFECT_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_PORTAL_EFFECT_CMD)
	cmd.WriteUint16(1)
	cmd.WriteUint16(1)
	cmd.WriteUint16(player.MapObjectID)
	h.broadcastMiniHouse(st, cmd)
}

// handleMiniHouseFurnitureEmotionReq maneja PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_REQ (0x8C62).
func (h *Handler) handleMiniHouseFurnitureEmotionReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || pkt.Reader.Remaining() < 4 {
		return
	}
	furHnd, _ := pkt.Reader.ReadUint16()
	pos, _ := pkt.Reader.ReadByte()
	emoNo, _ := pkt.Reader.ReadByte()
	st := h.currentMiniHouseOf(player)

	var errCode uint16
	if h.dataStore != nil && h.dataStore.MiniHouseCatalog != nil && len(h.dataStore.MiniHouseCatalog.ValidEmotes) > 0 {
		if !h.dataStore.MiniHouseCatalog.ValidEmotes[uint16(emoNo)] && emoNo == 0 {
			errCode = 1
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_ACK)
	ack.WriteUint16(furHnd)
	ack.WriteByte(pos)
	ack.WriteByte(emoNo)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)

	if errCode == 0 {
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CMD)
		cmd.WriteUint16(1)
		cmd.WriteUint16(player.MapObjectID)
		cmd.WriteUint16(furHnd)
		cmd.WriteByte(pos)
		cmd.WriteByte(emoNo)
		h.broadcastMiniHouse(st, cmd)
	}
}

// handleMiniHouseFurnitureEmotionCancelReq maneja PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ (0x8C65).
func (h *Handler) handleMiniHouseFurnitureEmotionCancelReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok {
		return
	}
	st := h.currentMiniHouseOf(player)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_CMD)
	cmd.WriteUint16(1)
	cmd.WriteUint16(player.MapObjectID)
	h.broadcastMiniHouse(st, cmd)
}

// ============================================================================
// Dept 53: NC_PET (Mascotas + Invocación, Hambre, Humor, Tendencia, GameLog G)
// ============================================================================

// handlePetAskNewNameReq maneja PROTO_NC_PET_ASK_NEW_NAME_REQ (0xD405).
func (h *Handler) handlePetAskNewNameReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 7 {
		return
	}
	petHnd, _ := pkt.Reader.ReadUint16()
	petID, _ := pkt.Reader.ReadUint32()
	firstSummon, _ := pkt.Reader.ReadByte()

	if firstSummon != 0 && h.charRepo != nil {
		petNo, _ := h.charRepo.CreatePet(player.Char.ID, fmt.Sprintf("Pet%d", petID), 100, 0)
		player.ActivePet = &ActivePetState{
			Handle:     petHnd,
			PetID:      petNo,
			PetType:    uint16(petID),
			Mind:       100,
			Hunger:     100,
			Stress:     0,
			Tendency:   0,
			SummonedAt: time.Now(),
		}
	} else if player.ActivePet == nil {
		petNo := 1
		if h.charRepo != nil {
			if pd, err := h.charRepo.GetPetInfo(player.Char.ID); err == nil && pd != nil {
				petNo = pd.PetNo
			}
		}
		player.ActivePet = &ActivePetState{
			Handle:     petHnd,
			PetID:      petNo,
			PetType:    uint16(petID),
			Mind:       100,
			Hunger:     100,
			Stress:     0,
			SummonedAt: time.Now(),
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_PET_ASK_NEW_NAME_ACK)
	ack.WriteUint16(0)
	ack.WriteUint16(petHnd)
	_ = client.SendPacket(ack)

	if h.gameLogDispatcher != nil {
		mapID := ""
		if player.CurrentMap != nil {
			mapID = player.CurrentMap.MapID
		}
		h.gameLogDispatcher.EnqueueTypeG(1, player.Char.ID, mapID, player.Char.XPos, player.Char.YPos, 0, int(petID), int64(petHnd), int(firstSummon), 0, 0)
	}
}

// handlePetSetNameReq maneja PROTO_NC_PET_SET_NAME_REQ (0xD407).
func (h *Handler) handlePetSetNameReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 18 {
		return
	}
	petHnd, _ := pkt.Reader.ReadUint16()
	petName := readFixedStr(pkt.Reader, 16)

	var petID uint32
	if player.ActivePet != nil {
		petID = uint32(player.ActivePet.PetID)
	}
	cleanName := strings.TrimSpace(petName)
	var errCode uint16
	if len(cleanName) < 4 || len(cleanName) > 16 || strings.ContainsAny(cleanName, " \t\r\n") {
		errCode = 0x34c0 // Zone.exe sp_NC_PET_SET_NAME_REQ (0x4d1d59): len 4..16 / GetIsBadName
	} else if h.charRepo != nil {
		if petID == 0 {
			newNo, _ := h.charRepo.CreatePet(player.Char.ID, cleanName, 100, 0)
			petID = uint32(newNo)
		} else {
			_ = h.charRepo.SetPetName(int(petID), player.Char.ID, cleanName)
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_PET_SET_NAME_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint32(petID)
	_ = client.SendPacket(ack)

	if errCode == 0 {
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_PET_SET_NAME_CMD)
		cmd.WriteUint16(petHnd)
		writeFixedStr(cmd, petName, 16)
		if player.CurrentMap != nil {
			player.CurrentMap.Broadcast(cmd, 0)
		} else {
			_ = client.SendPacket(cmd)
		}
		if h.gameLogDispatcher != nil {
			h.gameLogDispatcher.EnqueueTypeG(2, player.Char.ID, "", player.Char.XPos, player.Char.YPos, 0, int(petID), int64(petHnd), 0, 0, 0)
		}
	}
}

// handlePetSetNameCancelReq maneja PROTO_NC_PET_SET_NAME_CANCEL_REQ (0xD409).
func (h *Handler) handlePetSetNameCancelReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	_, _ = pkt.Reader.ReadUint16()
	var petID uint32
	if player.ActivePet != nil {
		petID = uint32(player.ActivePet.PetID)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_PET_SET_NAME_CANCEL_ACK)
	ack.WriteUint16(0)
	ack.WriteUint32(petID)
	_ = client.SendPacket(ack)
}

// FeedActivePet incrementa hambre/humor de la mascota invocada y persiste tendencia.
func (h *Handler) FeedActivePet(player *ZonePlayer, mindDelta, stressDelta int) bool {
	if player == nil || player.ActivePet == nil {
		return false
	}
	newMind := int(player.ActivePet.Mind) + mindDelta
	if newMind > 100 {
		newMind = 100
	} else if newMind < 0 {
		newMind = 0
	}
	newStress := int(player.ActivePet.Stress) + stressDelta
	if newStress > 100 {
		newStress = 100
	} else if newStress < 0 {
		newStress = 0
	}
	player.ActivePet.Mind = uint8(newMind)
	player.ActivePet.Stress = uint8(newStress)
	player.ActivePet.Hunger = uint8(newMind)
	if mindDelta > 0 && int(player.ActivePet.Tendency) < 100 {
		player.ActivePet.Tendency++
	}
	if h.charRepo != nil && player.ActivePet.PetID > 0 {
		_ = h.charRepo.SetPetTendency(player.ActivePet.PetID, byte(player.ActivePet.Tendency))
	}
	if h.gameLogDispatcher != nil && player.Char != nil {
		h.gameLogDispatcher.EnqueueTypeG(3, player.Char.ID, "", player.Char.XPos, player.Char.YPos, 0, player.ActivePet.PetID, int64(player.ActivePet.Handle), newMind, newStress, 0)
	}
	return true
}

// TickActivePet actualiza el hambre periódica de la mascota y dispara eventos Lua GTI.
func (h *Handler) TickActivePet(player *ZonePlayer) {
	if player == nil || player.ActivePet == nil {
		return
	}
	if player.ActivePet.Mind > 0 {
		player.ActivePet.Mind--
	}
	if player.ActivePet.Hunger > 0 {
		player.ActivePet.Hunger--
	}
	if player.CurrentMap != nil {
		player.CurrentMap.mu.RLock()
		eng := player.CurrentMap.luaEngine
		player.CurrentMap.mu.RUnlock()
		if eng != nil {
			_, _ = eng.CallFunction("Pet_AI", glua.LNumber(player.ActivePet.Handle), glua.LNumber(1000))
		}
	}
}

// Métodos ZoneMap bridgePet* requeridos por ZoneWorldBridge en mob_ai.go:

func (m *ZoneMap) findPlayerWithPet(petHandle uint32) *ZonePlayer {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.players {
		if p != nil && p.ActivePet != nil && uint32(p.ActivePet.Handle) == petHandle {
			return p
		}
	}
	return nil
}

func (m *ZoneMap) bridgePetGetMind(petHandle uint32) int {
	p := m.findPlayerWithPet(petHandle)
	if p == nil || p.ActivePet == nil {
		return 0
	}
	return int(p.ActivePet.Mind)
}

func (m *ZoneMap) bridgePetChangeMind(petHandle uint32, newMind int) {
	p := m.findPlayerWithPet(petHandle)
	if p == nil || p.ActivePet == nil {
		return
	}
	if newMind < 0 {
		newMind = 0
	} else if newMind > 100 {
		newMind = 100
	}
	p.ActivePet.Mind = uint8(newMind)
}

func (m *ZoneMap) bridgePetUnsummon(petHandle uint32) {
	p := m.findPlayerWithPet(petHandle)
	if p == nil || p.ActivePet == nil {
		return
	}
	petID := p.ActivePet.PetID
	p.ActivePet = nil
	if m.handler != nil && m.handler.gameLogDispatcher != nil && p.Char != nil {
		m.handler.gameLogDispatcher.EnqueueTypeG(4, p.Char.ID, m.MapID, p.Char.XPos, p.Char.YPos, 0, petID, int64(petHandle), 0, 0, 0)
	}
}

func (m *ZoneMap) bridgePetGetActionRecord(petHandle uint32, action string) (string, string, string, string, bool) {
	if m == nil || m.dataStore == nil || m.dataStore.PetCatalog == nil {
		return "", "", "", "", false
	}
	def, ok := m.dataStore.PetCatalog.ActionsByInx[strings.ToLower(strings.TrimSpace(action))]
	if !ok || def == nil {
		return "", "", "", "", false
	}
	return def.InxName, fmt.Sprintf("%d", def.ID), fmt.Sprintf("%d", def.ActionType), fmt.Sprintf("%d", def.TargetValue), true
}

func (m *ZoneMap) bridgePetGetStress(petHandle uint32) int {
	p := m.findPlayerWithPet(petHandle)
	if p == nil || p.ActivePet == nil {
		return 0
	}
	return int(p.ActivePet.Stress)
}

func (m *ZoneMap) bridgePetSaveTendency(petHandle uint32) bool {
	p := m.findPlayerWithPet(petHandle)
	if p == nil || p.ActivePet == nil {
		return false
	}
	if m.handler != nil && m.handler.charRepo != nil && p.ActivePet.PetID > 0 {
		_ = m.handler.charRepo.SetPetTendency(p.ActivePet.PetID, byte(p.ActivePet.Tendency))
	}
	return true
}

// ============================================================================
// Dept 23: NC_WT (Weapon Title - Títulos de Armas y Conteo de Mobs)
// ============================================================================

// handleWTLicenseReq maneja PROTO_NC_WT_LICENSE_REQ (0x5C01).
func (h *Handler) handleWTLicenseReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player.Char == nil || pkt.Reader.Remaining() < 7 {
		return
	}
	licInv, _ := pkt.Reader.ReadUint16()
	wepInv, _ := pkt.Reader.ReadUint16()
	mobID, _ := pkt.Reader.ReadUint16()
	bSet, _ := pkt.Reader.ReadByte()

	licStorage := byte(licInv >> 8)
	licSlot := int16(licInv & 0xFF)
	wepStorage := byte(wepInv >> 8)
	wepSlot := int16(wepInv & 0xFF)

	var errCode uint16
	if h.charRepo != nil {
		var wepItem, licItem *database.Item
		itemsWep, _ := h.charRepo.GetItemsByStorage(player.Char.ID, wepStorage)
		for _, it := range itemsWep {
			if it.Storage == wepSlot {
				wepItem = it
				break
			}
		}
		itemsLic, _ := h.charRepo.GetItemsByStorage(player.Char.ID, licStorage)
		for _, it := range itemsLic {
			if it.Storage == licSlot {
				licItem = it
				break
			}
		}
		if wepItem == nil {
			errCode = 1
		} else if bSet != 0 {
			if licItem == nil {
				errCode = 2
			} else {
				_ = h.charRepo.SetWeaponTitle(wepItem.ItemKey, licItem.ItemKey, int16(mobID), player.Char.ID, player.Char.Name)
				_ = h.charRepo.SetItemOption(wepItem.ItemKey, 200, int64(mobID))
				_ = h.charRepo.SetItemOption(wepItem.ItemKey, 201, 0)
			}
		} else {
			var clearKey int64
			if licItem != nil {
				clearKey = licItem.ItemKey
			}
			_ = h.charRepo.ClearWeaponTitle(wepItem.ItemKey, int16(mobID), clearKey, 1, 0, 0)
			_ = h.charRepo.SetItemOption(wepItem.ItemKey, 200, 0)
			_ = h.charRepo.SetItemOption(wepItem.ItemKey, 201, 0)
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_WT_LICENSE_ACK)
	ack.WriteUint16(licInv)
	ack.WriteUint16(wepInv)
	ack.WriteByte(bSet)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)

	if errCode == 0 {
		titleCmd := protocol.NewPacketWriter(protocol.PROTO_NC_WT_TITLE_CMD)
		titleCmd.WriteUint16(player.MapObjectID)
		titleCmd.WriteUint16(mobID)
		titleCmd.WriteByte(bSet)
		_ = client.SendPacket(titleCmd)
	}
}

// OnWeaponTitleMobKilled incrementa el contador de bajas del arma equipada contra WeaponTitleData.shn.
func (h *Handler) OnWeaponTitleMobKilled(player *ZonePlayer, mobID uint16) {
	if h == nil || player == nil || player.Char == nil || h.charRepo == nil || mobID == 0 {
		return
	}
	equips, err := h.charRepo.GetItemsByStorage(player.Char.ID, 0)
	if err != nil {
		return
	}
	for _, eq := range equips {
		opts, _ := h.charRepo.GetItemOptions(eq.ItemKey)
		if len(opts) == 0 {
			continue
		}
		targetMob := uint16(opts[200])
		if targetMob != mobID {
			continue
		}
		newKills := opts[201] + 1
		_ = h.charRepo.SetItemOption(eq.ItemKey, 201, newKills)

		if player.Client != nil {
			incCmd := protocol.NewPacketWriter(protocol.PROTO_NC_WT_MOBINC_CMD)
			incCmd.WriteUint16(mobID)
			_ = player.Client.SendPacket(incCmd)
		}

		// Evaluar grado alcanzado según WeaponTitleData.shn
		if h.dataStore != nil && h.dataStore.WeaponTitleCatalog != nil {
			var bestGrade uint8
			for _, def := range h.dataStore.WeaponTitleCatalog.Rows {
				if def.MobID == mobID && uint32(newKills) >= def.MobKillCount && def.Level > bestGrade {
					bestGrade = def.Level
				}
			}
			if bestGrade > 0 && player.Client != nil {
				grCmd := protocol.NewPacketWriter(protocol.PROTO_NC_WT_GRADE_CMD)
				grCmd.WriteUint16(uint16(eq.Storage))
				grCmd.WriteByte(bestGrade)
				_ = player.Client.SendPacket(grCmd)
			}
		}
	}
}
