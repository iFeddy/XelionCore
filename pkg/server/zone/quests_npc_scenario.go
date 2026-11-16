package zone

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/script/lua"
	glua "github.com/yuin/gopher-lua"
)

// Offsets oficiales dentro del blob de 30B (Status u8 + PLAYER_QUEST_DATA 29B) de PLAYER_QUEST_INFO (sizeof=32B en Zone.pdb).
const (
	qdataStatus      = 0  // unsigned char Status
	qdataStartTime   = 1  // __int64 StartTime (+0 en PLAYER_QUEST_DATA)
	qdataEndTime     = 9  // __int64 EndTime (+8 en PLAYER_QUEST_DATA)
	qdataRepeatCount = 17 // unsigned int RepeatCount (+16 en PLAYER_QUEST_DATA)
	qdataProgress    = 21 // unsigned char ProgressStep (+20 en PLAYER_QUEST_DATA)
	qdataMobCount    = 22 // unsigned char[5] End_NPCMobCount (+21 en PLAYER_QUEST_DATA)
	qdataEndLocScen  = 27 // bitfield End_Location / End_Scenario (+26 en PLAYER_QUEST_DATA)
	qdataRunTimeSec  = 28 // unsigned short End_RunningTimeSec (+27 en PLAYER_QUEST_DATA)
)

// Comandos oficiales enum QUEST_SCRIPT_COMMAND en Zone.pdb (STRUCT_QSC sizeof=101).
const (
	qscNone     uint32 = 0
	qscSay      uint32 = 2
	qscScenario uint32 = 3
	qscAccept   uint32 = 6
	qscCancel   uint32 = 7
	qscProgress uint32 = 8
	qscFailed   uint32 = 9
	qscDone     uint32 = 10
	qscLink     uint32 = 11
)

func writeStructQSC(pw *protocol.PacketWriter, cmd uint32, isPigeon byte, data []byte) {
	var qsc [101]byte
	binary.LittleEndian.PutUint32(qsc[0:4], cmd)
	qsc[4] = isPigeon
	if len(data) > 0 {
		copy(qsc[5:101], data)
	}
	pw.WriteBytes(qsc[:])
}

// ============================================================================
// ZoneMap Bridge Implementations (Quests, Mensajes, NPCs y Tutorial)
// ============================================================================

func (m *ZoneMap) bridgeCastTeleport(handle uint32, mapName string, x, y float64) {
	if m == nil || m.handler == nil {
		return
	}
	if mapName == "" {
		mapName = m.MapID
	}
	if handle == 0 {
		for _, p := range m.GetPlayers() {
			m.handler.WarpPlayer(p, mapName, int(x), int(y), p.Rotation)
		}
		return
	}
	if p := m.GetPlayer(uint16(handle)); p != nil {
		m.handler.WarpPlayer(p, mapName, int(x), int(y), p.Rotation)
	}
}

// AddNPC registra un nuevo NPC dinámico en el mapa y lo inserta en la grilla espacial.
func (m *ZoneMap) AddNPC(npc *ZoneNPC) {
	if m == nil || npc == nil {
		return
	}
	m.mu.Lock()
	if npc.Handle == 0 {
		m.nextObjID++
		npc.Handle = m.nextObjID
	}
	if npc.Map == "" {
		npc.Map = m.MapID
	}
	m.npcs[npc.Handle] = npc
	m.mu.Unlock()
	if m.grid != nil {
		m.grid.AddEntity(npc)
	}
}

// RemoveNPC elimina un NPC del mapa y su referencia espacial.
func (m *ZoneMap) RemoveNPC(handle uint16) {
	if m == nil {
		return
	}
	m.mu.Lock()
	npc := m.npcs[handle]
	delete(m.npcs, handle)
	m.mu.Unlock()
	if npc != nil && m.grid != nil {
		m.grid.RemoveEntity(handle)
	}
}

func (m *ZoneMap) bridgeNPCRegen(npcID uint32, mapName string, x, y float64) uint32 {
	if m == nil {
		return 0
	}
	zm := m
	if mapName != "" && m.handler != nil && m.handler.mapManager != nil && !strings.EqualFold(mapName, m.MapID) {
		if target := m.handler.mapManager.GetMap(mapName); target != nil {
			zm = target
		}
	}
	npc := &ZoneNPC{
		MobID:  uint16(npcID),
		X:      int32(x),
		Y:      int32(y),
		Direct: 0,
		Menu:   1,
	}
	zm.AddNPC(npc)
	return uint32(npc.Handle)
}

func (m *ZoneMap) bridgeSpawnNamedNPC(mobInx string, x, y float64, dir int) uint32 {
	if m == nil {
		return 0
	}
	var npcID uint16
	if m.dataStore != nil && m.dataStore.MobCatalog != nil {
		if tmpl := m.dataStore.MobCatalog.LookupTemplate(mobInx); tmpl != nil {
			npcID = tmpl.ID
		}
	}
	if npcID == 0 {
		npcID = 100
	}
	npc := &ZoneNPC{
		MobID:  npcID,
		Name:   mobInx,
		X:      int32(x),
		Y:      int32(y),
		Direct: int16(dir),
		Menu:   1,
	}
	m.AddNPC(npc)
	return uint32(npc.Handle)
}

func (m *ZoneMap) bridgeNPCVanish(handle uint32) {
	if m == nil {
		return
	}
	m.RemoveNPC(uint16(handle))
}

func (m *ZoneMap) bridgeGroupRegen(groupID uint32, mapName string) int {
	if m == nil || m.dataStore == nil || m.mobCatalog == nil {
		return 0
	}
	mr := m.dataStore.MobRegen[m.MapID]
	if mr == nil {
		return 0
	}
	count := 0
	grpStr := fmt.Sprintf("%d", groupID)
	for gKey, g := range mr.Groups {
		if gKey == grpStr || strings.Contains(strings.ToLower(gKey), strings.ToLower(grpStr)) {
			for _, me := range g.Mobs {
				if tmpl := m.mobCatalog.LookupTemplate(me.Mob); tmpl != nil {
					cnt := me.Count
					if cnt <= 0 {
						cnt = 1
					}
					for i := 0; i < cnt; i++ {
						m.SpawnMobFromTemplate(tmpl, int32(g.CenterX), int32(g.CenterY), 0, gKey, me.RegStandard, true)
						count++
					}
				}
			}
		}
	}
	return count
}

func (m *ZoneMap) bridgeBroadcastNotice(msg string) {
	if m == nil {
		return
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NOTICE_CMD)
	msgBytes := []byte(msg)
	if len(msgBytes) > 250 {
		msgBytes = msgBytes[:250]
	}
	pw.WriteByte(0)
	pw.WriteByte(byte(len(msgBytes)))
	pw.WriteBytes(msgBytes)
	for _, p := range m.GetPlayers() {
		if p.Client != nil {
			_ = p.Client.SendPacket(pw)
		}
	}
}

func (m *ZoneMap) bridgeScriptMsg(handle uint32, msg string) {
	if m == nil {
		return
	}
	// PROTO_NC_ACT_SCRIPT_MSG_CMD (0x2079, sizeof = 39 + sum(nlen) -> wire = 41 + sum(nlen))
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SCRIPT_MSG_CMD)
	pw.WriteUint16(uint16(handle))
	writeFixedStr(pw, msg, 32) // sScriptMsgIndex
	msgBytes := []byte(msg)
	if len(msgBytes) > 200 {
		msgBytes = msgBytes[:200]
	}
	var nlen [5]byte
	nlen[0] = byte(len(msgBytes))
	pw.WriteBytes(nlen[:])
	pw.WriteBytes(msgBytes)
	if handle != 0 {
		if p := m.GetPlayer(uint16(handle)); p != nil && p.Client != nil {
			_ = p.Client.SendPacket(pw)
			return
		}
	}
	for _, p := range m.GetPlayers() {
		if p.Client != nil {
			_ = p.Client.SendPacket(pw)
		}
	}
}

func (m *ZoneMap) bridgeSysMsg(handle uint32, code uint16) {
	if m == nil {
		return
	}
	if p := m.GetPlayer(uint16(handle)); p != nil && p.Client != nil {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_MISCERROR_CMD)
		pw.WriteUint16(code)
		_ = p.Client.SendPacket(pw)
	}
}

func (m *ZoneMap) bridgeServerMenu(playerHandle, npcHandle uint32, title string, options []lua.MenuOptionEntry) {
	if m == nil {
		return
	}
	p := m.GetPlayer(uint16(playerHandle))
	if p == nil {
		return
	}
	p.PendingMenuNPC = uint16(npcHandle)
	p.PendingMenuCallbacks = make([]string, len(options))
	for i, opt := range options {
		p.PendingMenuCallbacks[i] = opt.Callback
	}
	if p.Client == nil {
		return
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MENU_SERVERMENU_REQ)
	writeFixedStr(pw, title, 128)
	pw.WriteByte(1) // Priority
	pw.WriteUint16(uint16(npcHandle))
	pw.WriteUint32(uint32(p.Char.XPos))
	pw.WriteUint32(uint32(p.Char.YPos))
	pw.WriteUint16(300) // LimitRange
	pw.WriteByte(byte(len(options)))
	for i, opt := range options {
		// SERVERMENU (sizeof = 33): +0 reply (u8), +1 string ([32]u8)
		pw.WriteByte(byte(i))
		writeFixedStr(pw, opt.Text, 32)
	}
	_ = p.Client.SendPacket(pw)
}

func (m *ZoneMap) bridgeNPCMenuOpen(npcHandle, playerHandle uint32) {
	if m == nil {
		return
	}
	p := m.GetPlayer(uint16(playerHandle))
	if p == nil || p.Client == nil {
		return
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCMENUOPEN_ACK)
	pw.WriteByte(1)
	_ = p.Client.SendPacket(pw)
}

func (m *ZoneMap) bridgeQuestResult(mapName, result string) {
	if m == nil {
		return
	}
	isSuccess := strings.EqualFold(result, "Success") || result == "1"
	for _, p := range m.GetPlayers() {
		m.bridgeQuestResultIndividual(uint32(p.MapObjectID), result)
		if m.InstanceKind == "KQ" {
			m.bridgeEndOfKingdomQuest(isSuccess)
		}
	}
}

func (m *ZoneMap) bridgeQuestResultIndividual(playerHandle uint32, result string) {
	if m == nil || m.handler == nil {
		return
	}
	p := m.GetPlayer(uint16(playerHandle))
	if p == nil || m.handler.charRepo == nil || p.Char == nil {
		return
	}
	isSuccess := strings.EqualFold(result, "Success") || result == "1"
	if !isSuccess {
		return
	}
	quests, _ := m.handler.charRepo.GetCharacterQuests(p.Char.ID)
	for _, q := range quests {
		if q.Status == 1 {
			if m.dataStore != nil && m.dataStore.QuestCatalog != nil {
				if qDef := m.dataStore.QuestCatalog.GetQuest(uint16(q.QuestNo)); qDef != nil {
					m.handler.completeQuestWithRewards(p, qDef, 0)
					break
				}
			}
		}
	}
}

func (m *ZoneMap) bridgeQuestMobKill(playerHandle uint32, questID uint16, mobID uint32, count int) {
	if m == nil || m.handler == nil {
		return
	}
	p := m.GetPlayer(uint16(playerHandle))
	if p == nil {
		return
	}
	for i := 0; i < count; i++ {
		m.handler.advanceSpecificQuestMobKill(p, questID, uint16(mobID))
	}
}

func (m *ZoneMap) bridgeQuestMobKillAllInMap(mapName string, questID uint16, mobID uint32, maxCount int) {
	if m == nil || m.handler == nil {
		return
	}
	for _, p := range m.GetPlayers() {
		m.handler.advanceSpecificQuestMobKill(p, questID, uint16(mobID))
	}
}

func (m *ZoneMap) bridgeGetTutorialInfo(playerHandle uint32) int {
	if m == nil {
		return 0
	}
	if p := m.GetPlayer(uint16(playerHandle)); p != nil && p.Char != nil {
		return int(p.Char.TutorialStep)
	}
	return 0
}

func (m *ZoneMap) bridgeSaveTutorialInfo(playerHandle uint32, saveType, step int) {
	if m == nil || m.handler == nil {
		return
	}
	p := m.GetPlayer(uint16(playerHandle))
	if p == nil || p.Char == nil {
		return
	}
	p.Char.TutorialStep = byte(step)
	if m.handler.charRepo != nil {
		_ = m.handler.charRepo.UpdateCharacter(p.Char)
	}
	if p.Client != nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_ACK)
		ack.WriteUint32(uint32(p.Char.ID))
		ack.WriteByte(byte(saveType))
		ack.WriteByte(byte(step))
		ack.Fill(3, 0)
		_ = p.Client.SendPacket(ack)
	}
}

func (m *ZoneMap) bridgeProgressTutorial(playerHandle uint32, step int) {
	if m == nil {
		return
	}
	p := m.GetPlayer(uint16(playerHandle))
	if p == nil || p.Client == nil {
		return
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_TUTORIAL_STEP_REQ)
	pw.WriteByte(byte(step))
	_ = p.Client.SendPacket(pw)
}

// ============================================================================
// Dept 17: Motor Completo de Misiones (Quests)
// ============================================================================

// GetPlayer devuelve el ZonePlayer asociado al cliente o nil.
func (h *Handler) GetPlayer(client *network.Client) *ZonePlayer {
	p, ok := h.getPlayer(client)
	if !ok {
		return nil
	}
	return p
}

func (h *Handler) getQuestCatalog() *QuestCatalog {
	if h.dataStore != nil {
		return h.dataStore.QuestCatalog
	}
	return nil
}

func (h *Handler) sendQuestErr(player *ZonePlayer, errCode uint16) {
	if player == nil || player.Client == nil {
		return
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_ERR)
	pw.WriteUint16(errCode)
	_ = player.Client.SendPacket(pw)
}

// handleQuestReadReq maneja PROTO_NC_QUEST_READ_REQ (0x4416).
func (h *Handler) handleQuestReadReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	questID, _ := pkt.Reader.ReadUint16()
	var errCode uint16
	if qCat := h.getQuestCatalog(); qCat != nil && qCat.GetQuest(questID) == nil {
		errCode = 0x0B59 // Zone.exe sp_NC_QUEST_READ_REQ: quest not found / DB read fail
	} else if h.charRepo != nil {
		quests, _ := h.charRepo.GetCharacterQuests(player.Char.ID)
		found := false
		for _, q := range quests {
			if uint16(q.QuestNo) == questID {
				found = true
				buf := make([]byte, 30)
				copy(buf, q.Data)
				buf[qdataStatus] = byte(q.Status)
				buf[qdataProgress] |= 1 // flag de lectura en ProgressStep
				_ = h.charRepo.SetCharacterQuest(player.Char.ID, int(questID), q.Status, buf)
				break
			}
		}
		if !found {
			errCode = 0x0B59
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_READ_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint16(questID)
	_ = client.SendPacket(ack)
}

// handleQuestStartReq maneja PROTO_NC_QUEST_START_REQ (0x4414).
func (h *Handler) handleQuestStartReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	questID, _ := pkt.Reader.ReadUint16()
	errCode := h.startQuestForPlayer(player, questID, 0)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_START_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

// handleQuestSelectStartReq maneja PROTO_NC_QUEST_SELECT_START_REQ (0x440F).
func (h *Handler) handleQuestSelectStartReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 4 {
		return
	}
	npcID, _ := pkt.Reader.ReadUint16()
	questID, _ := pkt.Reader.ReadUint16()
	errCode := h.startQuestForPlayer(player, questID, npcID)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_SELECT_START_ACK)
	ack.WriteUint16(npcID)
	ack.WriteUint16(questID)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

func (h *Handler) startQuestForPlayer(player *ZonePlayer, questID, npcID uint16) uint16 {
	qCat := h.getQuestCatalog()
	if qCat != nil {
		qDef := qCat.GetQuest(questID)
		if qDef == nil {
			if npcID != 0 {
				return 0x0B47 // Recv_NC_QUEST_SELECT_START_REQ: QuestData nil
			}
			return 0x0B51 // sp_NC_QUEST_START_REQ: invalid quest/NPC
		}
		if qDef.CheckLevel && (player.Char.Level < qDef.MinLevel || (qDef.MaxLevel > 0 && player.Char.Level > qDef.MaxLevel)) {
			return 0x0B50 // QuestPlayer_NotCheckNPC / requisitos incumplidos
		}
		if npcID != 0 && qDef.StartNPCID != 0 && qDef.StartNPCID != npcID {
			return 0x0B47 // Recv_NC_QUEST_SELECT_START_REQ: GetQuestStatusWithNPC == 0
		}
		if qDef.CheckPrevQuest && qDef.PrevQuestID > 0 && h.charRepo != nil {
			doneList, _ := h.charRepo.GetCharacterDoneQuests(player.Char.ID)
			foundPrev := false
			for _, dq := range doneList {
				if uint16(dq) == qDef.PrevQuestID {
					foundPrev = true
					break
				}
			}
			if !foundPrev {
				return 0x0B50
			}
		}
	}
	if h.charRepo != nil {
		existing, _ := h.charRepo.GetCharacterQuests(player.Char.ID)
		for _, q := range existing {
			if uint16(q.QuestNo) == questID && q.Status == 1 {
				return 0x0B41 // ya activa
			}
		}
		data := make([]byte, 30)
		data[qdataStatus] = 1
		binary.LittleEndian.PutUint64(data[qdataStartTime:qdataStartTime+8], uint64(time.Now().Unix()))
		data[qdataProgress] = 1
		if err := h.charRepo.AddCharacterQuest(player.Char.ID, int(questID), 1, data); err != nil {
			return 6
		}
		h.syncQuestDoingPacket(player)
	}
	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType9(1, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int64(npcID), int(questID), 1)
	}
	return 0
}

func (h *Handler) syncQuestDoingPacket(player *ZonePlayer) {
	if player == nil || player.Client == nil || h.charRepo == nil || player.Char == nil {
		return
	}
	allQuests, _ := h.charRepo.GetCharacterQuests(player.Char.ID)
	var doing []*database.QuestData
	for _, q := range allQuests {
		if q.Status == 1 {
			doing = append(doing, q)
		}
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD)
	pw.WriteUint32(uint32(player.Char.ID))
	pw.WriteByte(1)
	pw.WriteByte(byte(len(doing)))
	for _, q := range doing {
		pw.WriteUint16(uint16(q.QuestNo))
		if len(q.Data) >= 30 {
			pw.WriteBytes(q.Data[:30])
		} else {
			pw.WriteBytes(q.Data)
			pw.Fill(30-len(q.Data), 0)
		}
	}
	_ = player.Client.SendPacket(pw)
}

// handleQuestGiveUpReq maneja PROTO_NC_QUEST_GIVE_UP_REQ (0x4407).
func (h *Handler) handleQuestGiveUpReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	questID, _ := pkt.Reader.ReadUint16()
	now := time.Now()
	var errCode uint16
	if !player.LastQuestGiveUpAt.IsZero() && now.Sub(player.LastQuestGiveUpAt) < 100*time.Millisecond {
		errCode = 0x0B49 // Zone.exe sp_NC_QUEST_GIVE_UP_REQ 100ms cooldown
	} else {
		if qCat := h.getQuestCatalog(); qCat != nil && qCat.GetQuest(questID) == nil {
			errCode = 0x0B45
		}
		if errCode == 0 && h.charRepo != nil {
			quests, _ := h.charRepo.GetCharacterQuests(player.Char.ID)
			doing := false
			for _, q := range quests {
				if uint16(q.QuestNo) == questID && q.Status == 1 {
					doing = true
					break
				}
			}
			if !doing {
				errCode = 0x0B45 // !IsDoingQuest
			} else if err := h.charRepo.DeleteCharacterQuest(player.Char.ID, int(questID)); err != nil {
				errCode = 0x0B45
			} else {
				player.LastQuestGiveUpAt = now
			}
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_GIVE_UP_ACK)
	ack.WriteUint16(questID)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
	if errCode == 0 && h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType9(2, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, 0, int(questID), 0)
	}
}

// handleQuestClientScenarioDoneReq maneja PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ (0x440B).
func (h *Handler) handleQuestClientScenarioDoneReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	scenarioID, _ := pkt.Reader.ReadUint16()
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK)
	ack.WriteUint16(scenarioID)
	_ = client.SendPacket(ack)
}

// handleQuestScriptCmdAck maneja PROTO_NC_QUEST_SCRIPT_CMD_ACK (0x4402).
func (h *Handler) handleQuestScriptCmdAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 7 {
		return
	}
	questID, _ := pkt.Reader.ReadUint16()
	qscStage, _ := pkt.Reader.ReadByte()
	result, _ := pkt.Reader.ReadUint32()

	if result == 0 {
		h.sendQuestErr(player, 0x0C0C) // Zone.exe Recv_NC_QUEST_SCRIPT_CMD_ACK (3084)
		return
	}
	qCat := h.getQuestCatalog()
	var qDef *QuestDef
	if qCat != nil {
		qDef = qCat.GetQuest(questID)
	}
	if qDef == nil {
		if h.charRepo != nil {
			buf := make([]byte, 30)
			buf[qdataStatus] = 2
			binary.LittleEndian.PutUint64(buf[qdataEndTime:qdataEndTime+8], uint64(time.Now().Unix()))
			buf[qdataProgress] = qscStage
			_ = h.charRepo.SetCharacterQuest(player.Char.ID, int(questID), 2, buf)
		}
		return
	}
	// Verifica si hay múltiples ítems seleccionables (Mode == 2) sin selección previa
	var selectableCount int
	for _, r := range qDef.Rewards {
		if r.Mode == 2 {
			selectableCount++
		}
	}
	if selectableCount > 1 {
		player.PendingQuestRewardSelect = questID
		need := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_REWARD_NEED_SELECT_ITEM_CMD)
		need.WriteUint16(questID)
		_ = client.SendPacket(need)
		return
	}
	h.completeQuestWithRewards(player, qDef, 0)
}

// handleQuestRewardSelectItemIndexCmd maneja PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD (0x4411).
func (h *Handler) handleQuestRewardSelectItemIndexCmd(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 6 {
		return
	}
	questID, _ := pkt.Reader.ReadUint16()
	selectedIdx, _ := pkt.Reader.ReadUint32()
	qCat := h.getQuestCatalog()
	if qCat == nil {
		return
	}
	qDef := qCat.GetQuest(questID)
	if qDef == nil {
		h.sendQuestErr(player, 1)
		return
	}
	player.PendingQuestRewardSelect = 0
	h.completeQuestWithRewards(player, qDef, int(selectedIdx))
}

func (h *Handler) completeQuestWithRewards(player *ZonePlayer, qDef *QuestDef, selectedIdx int) {
	if player == nil || player.Char == nil || qDef == nil {
		return
	}
	// Consumir ítems requeridos si están definidos en QuestData.shn
	if h.charRepo != nil {
		for _, reqItem := range qDef.ItemCollects {
			inv, _ := h.charRepo.GetCharacterItems(player.Char.ID)
			rem := uint32(reqItem.Count)
			for _, it := range inv {
				if it.StorageType == 9 && it.ItemID == reqItem.ItemID && rem > 0 {
					take := uint32(it.Flags)
					if take == 0 {
						take = 1
					}
					if take > rem {
						take = rem
					}
					if h.consumeItemUnits(player, it, take) {
						rem -= take
					}
				}
			}
		}
	}
	// Entregar EXP, Dinero e Ítems
	if qDef.ExpEnabled && qDef.RewardExp > 0 {
		player.Char.TCharacterState.Exp += int64(qDef.RewardExp)
	}
	if qDef.MoneyEnabled && qDef.RewardMoney > 0 {
		player.Char.TCharacterState.Money += int64(qDef.RewardMoney)
	}
	selCounter := 0
	for _, rw := range qDef.Rewards {
		if rw.Mode == 1 {
			h.addItemToInventory(player, rw.ItemID, uint32(rw.Count), nil)
		} else if rw.Mode == 2 {
			if selCounter == selectedIdx {
				h.addItemToInventory(player, rw.ItemID, uint32(rw.Count), nil)
			}
			selCounter++
		}
	}
	if h.charRepo != nil {
		doneBlob := make([]byte, 30)
		doneBlob[qdataStatus] = 2
		binary.LittleEndian.PutUint64(doneBlob[qdataEndTime:qdataEndTime+8], uint64(time.Now().Unix()))
		_ = h.charRepo.SetCharacterQuest(player.Char.ID, int(qDef.QuestID), 2, doneBlob)
		_ = h.charRepo.UpdateCharacter(player.Char)
	}
	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueType9(3, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int64(qDef.EndNPCID), int(qDef.QuestID), int(qDef.RewardExp))
	}
	if player.Client != nil {
		// PROTO_NC_QUEST_SCRIPT_CMD_REQ (0x4401, sizeof = 103 -> wire = 105: nQuestID + STRUCT_QSC)
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_SCRIPT_CMD_REQ)
		pw.WriteUint16(qDef.QuestID)
		var qscData [2]byte
		binary.LittleEndian.PutUint16(qscData[0:2], qDef.QuestID)
		writeStructQSC(pw, qscDone, 0, qscData[:])
		_ = player.Client.SendPacket(pw)
	}
}

// OnQuestMobKilled actualiza misiones Doing cuando un jugador mata un monstruo.
func (h *Handler) OnQuestMobKilled(player *ZonePlayer, mobID uint16) {
	if h == nil || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	qCat := h.getQuestCatalog()
	if qCat == nil {
		return
	}
	quests, _ := h.charRepo.GetCharacterQuests(player.Char.ID)
	type mobQuestHit struct {
		mobID   uint16
		questID uint16
	}
	var hits []mobQuestHit
	speciesList := qCat.MobToSpecies[mobID]
	for _, q := range quests {
		if q.Status != 1 {
			continue
		}
		qDef := qCat.GetQuest(uint16(q.QuestNo))
		if qDef == nil {
			continue
		}
		buf := make([]byte, 30)
		copy(buf, q.Data)
		updated := false
		for slotIdx, req := range qDef.MobKills {
			if slotIdx >= 5 || !req.Enabled {
				continue
			}
			matches := !req.IsSpecies && req.MobOrSpeciesID == mobID
			if !matches && req.IsSpecies {
				for _, spID := range speciesList {
					if spID == req.MobOrSpeciesID {
						matches = true
						break
					}
				}
			}
			if matches && buf[qdataMobCount+slotIdx] < req.Count {
				buf[qdataMobCount+slotIdx]++
				updated = true
				hits = append(hits, mobQuestHit{mobID: mobID, questID: qDef.QuestID})
			}
		}
		if updated {
			_ = h.charRepo.SetCharacterQuest(player.Char.ID, q.QuestNo, q.Status, buf)
		}
		for _, drp := range qDef.Drops {
			if drp.MobID == mobID && drp.ItemID > 0 {
				h.addItemToInventory(player, drp.ItemID, 1, nil)
			}
		}
	}
	if len(hits) > 0 && player.Client != nil {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD)
		pw.WriteByte(byte(len(hits)))
		for _, hit := range hits {
			pw.WriteUint16(hit.mobID)
			pw.WriteUint16(hit.questID)
		}
		_ = player.Client.SendPacket(pw)
	}
}

func (h *Handler) advanceSpecificQuestMobKill(player *ZonePlayer, questID, mobID uint16) {
	if h == nil || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	quests, _ := h.charRepo.GetCharacterQuests(player.Char.ID)
	for _, q := range quests {
		if uint16(q.QuestNo) == questID && q.Status == 1 {
			buf := make([]byte, 30)
			copy(buf, q.Data)
			buf[qdataMobCount]++
			_ = h.charRepo.SetCharacterQuest(player.Char.ID, q.QuestNo, q.Status, buf)
			if player.Client != nil {
				pw := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_NOTIFY_MOB_KILL_CMD)
				pw.WriteByte(1)
				pw.WriteUint16(mobID)
				pw.WriteUint16(questID)
				_ = player.Client.SendPacket(pw)
			}
			return
		}
	}
}

// ============================================================================
// Dept 15, Dept 8 (NPC Click/Menu), Dept 27 (Scenario) y Dept 4 (Tutorial)
// ============================================================================

// handleMenuServerMenuAck maneja PROTO_NC_MENU_SERVERMENU_ACK (0x3C02).
func (h *Handler) handleMenuServerMenuAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	reply, _ := pkt.Reader.ReadByte()
	idx := int(reply)
	if idx >= 10 || (player.OpenNPCHandle != 0 && !h.withinNPCRange(player, player.OpenNPCX, player.OpenNPCY)) {
		// Zone.exe sp_NC_MENU_SERVERMENU_ACK / sma_Reply (0x50FD40): distancia > 1000^2 o reply >= 10 -> sma_Reset
		h.resetPlayerNPCState(player)
		return
	}
	var cbName string
	if idx >= 0 && idx < len(player.PendingMenuCallbacks) {
		cbName = player.PendingMenuCallbacks[idx]
	}
	npcH := player.PendingMenuNPC
	player.PendingMenuCallbacks = nil
	if cbName != "" && player.CurrentMap != nil {
		player.CurrentMap.mu.RLock()
		eng := player.CurrentMap.luaEngine
		player.CurrentMap.mu.RUnlock()
		if eng != nil {
			_, _ = eng.CallFunction(cbName, glua.LNumber(npcH), glua.LNumber(player.MapObjectID), glua.LNumber(reply))
		}
	}
}

func (h *Handler) resetPlayerNPCState(player *ZonePlayer) {
	if player == nil {
		return
	}
	player.OpenNPCMobID = 0
	player.OpenNPCHandle = 0
	player.OpenNPCX = 0
	player.OpenNPCY = 0
	player.OpenMultyLinkNPC = 0
	player.OpenMultyLinkOptions = nil
	player.PendingMenuCallbacks = nil
}

func (h *Handler) withinNPCRange(player *ZonePlayer, x, y int32) bool {
	if player == nil || player.Char == nil {
		return false
	}
	dx := int64(x) - int64(player.Char.XPos)
	dy := int64(y) - int64(player.Char.YPos)
	return dx*dx+dy*dy <= 0xF4240
}

func (h *Handler) sendShopOpenTable(player *ZonePlayer, npc *ZoneNPC) {
	if player == nil || player.Client == nil || npc == nil || h.dataStore == nil || h.dataStore.NPCItem == nil || h.itemStore == nil {
		return
	}
	merchantKey := npc.Name
	if merchantKey == "" {
		merchantKey = npc.RoleArg
	}
	merchant, ok := h.dataStore.NPCItem.GetMerchant(merchantKey)
	if (!ok || merchant == nil) && npc.RoleArg != "" && npc.RoleArg != merchantKey {
		merchant, ok = h.dataStore.NPCItem.GetMerchant(npc.RoleArg)
	}
	if !ok || merchant == nil {
		return
	}

	type shopEntry struct {
		slot uint8
		id   uint16
	}
	var entries []shopEntry
	for _, it := range merchant.GetAllItems() {
		info, found := h.itemStore.GetItemInfoByName(it.Item)
		if !found || info == nil {
			continue
		}
		entries = append(entries, shopEntry{slot: it.Slot, id: info.ID})
		if len(entries) >= 80 {
			break
		}
	}

	opCode := protocol.PROTO_NC_MENU_SHOPOPENTABLE_ITEM_CMD
	switch strings.ToLower(npc.RoleArg) {
	case "weapon":
		opCode = protocol.PROTO_NC_MENU_SHOPOPENTABLE_WEAPON_CMD
	case "skill":
		opCode = protocol.PROTO_NC_MENU_SHOPOPENTABLE_SKILL_CMD
	}

	pw := protocol.NewPacketWriter(opCode)
	pw.WriteUint16(uint16(len(entries)))
	pw.WriteUint16(npc.Handle)
	for _, e := range entries {
		pw.WriteByte(e.slot)
		pw.WriteUint16(e.id)
	}
	_ = player.Client.SendPacket(pw)
}

func (h *Handler) openNPCRole(player *ZonePlayer, npc *ZoneNPC) {
	if player == nil || npc == nil {
		return
	}
	switch npc.Role {
	case "Merchant":
		h.sendShopOpenTable(player, npc)
	case "StoreManager":
		h.sendStoragePage(player, 10, 0, 1)
	}
}

// handleActNPCClickCmd maneja PROTO_NC_ACT_NPCCLICK_CMD (0x200A -> sp_NC_ACT_NPCCLICK_CMD @ 0x004C4B70).
func (h *Handler) handleActNPCClickCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	npcHandle, _ := pkt.Reader.ReadUint16()
	player.CurrentMap.mu.RLock()
	npc := player.CurrentMap.npcs[npcHandle]
	cbName := player.CurrentMap.NPCInteractCB[uint32(npcHandle)]
	eng := player.CurrentMap.luaEngine
	player.CurrentMap.mu.RUnlock()

	// Zone.exe 0x004C4660: NPC inexistente -> sp_EndOfTrade
	if npc == nil {
		h.resetPlayerNPCState(player)
		return
	}
	// Zone.exe 0x004C4660 / cm_MenuAction: distancia cuadrada <= 0xF4240 (1000^2)
	if !h.withinNPCRange(player, npc.X, npc.Y) {
		h.resetPlayerNPCState(player)
		return
	}

	player.OpenNPCMobID = npc.MobID
	player.OpenNPCHandle = npc.Handle
	player.OpenNPCX = npc.X
	player.OpenNPCY = npc.Y
	player.OpenMultyLinkNPC = 0
	player.OpenMultyLinkOptions = nil

	if npc.IsGate && npc.GateToWhere != "" {
		if strings.Contains(npc.GateToWhere, ",") {
			// PROTO_NC_MAP_MULTY_LINK_CMD (0x181F, sizeof = 73 -> wire = 75)
			parts := strings.Split(npc.GateToWhere, ",")
			if len(parts) > 5 {
				parts = parts[:5]
			}
			opts := make([]string, 0, len(parts))
			ml := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_MULTY_LINK_CMD)
			ml.WriteUint16(npcHandle)
			ml.WriteUint32(uint32(npc.X))
			ml.WriteUint32(uint32(npc.Y))
			ml.WriteUint16(350)
			ml.WriteByte(byte(len(parts)))
			for i := 0; i < 5; i++ {
				name := ""
				if i < len(parts) {
					name = strings.TrimSpace(parts[i])
					if name != "" {
						opts = append(opts, name)
					}
				}
				writeFixedStr(ml, name, 12)
			}
			player.OpenMultyLinkNPC = npc.Handle
			player.OpenMultyLinkOptions = opts
			_ = client.SendPacket(ml)
			return
		}
		h.triggerGateWarp(player, npc.GateToWhere)
		return
	}

	if cbName != "" && eng != nil {
		_, _ = eng.CallFunction(cbName, glua.LNumber(npcHandle), glua.LNumber(player.MapObjectID))
	}
	if h.npcStore != nil {
		if act, okAct := h.npcStore.GetActionByID(int(npc.MobID)); okAct {
			actPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPC_ACTION_CMD)
			actPkt.WriteByte(1)
			actPkt.WriteUint16(npcHandle)
			actPkt.WriteUint32(uint32(act.ID))
			_ = client.SendPacket(actPkt)
		}
	}
	if npc.Menu != 0 {
		openReq := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCMENUOPEN_REQ)
		openReq.WriteUint16(npc.MobID)
		_ = client.SendPacket(openReq)
		return
	}
	h.openNPCRole(player, npc)
}

// handleActNPCMenuOpenAck maneja PROTO_NC_ACT_NPCMENUOPEN_ACK (0x201D -> sp_NC_ACT_NPCMENUOPEN_ACK @ 0x004C4910).
func (h *Handler) handleActNPCMenuOpenAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.CurrentMap == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	reply, _ := pkt.Reader.ReadByte()
	if player.OpenNPCHandle == 0 {
		return
	}
	npc, found := player.CurrentMap.GetNPC(player.OpenNPCHandle)
	if !found || npc == nil || !h.withinNPCRange(player, npc.X, npc.Y) {
		h.resetPlayerNPCState(player)
		return
	}
	// sp_NC_ACT_NPCMENUOPEN_ACK @ 0x004C4A07: ejecuta NPCRoleBase::nrb_Role solo cuando reply == 1
	if reply == 1 {
		h.openNPCRole(player, npc)
	}
}

// handleActNPCMenuCmd maneja PROTO_NC_ACT_NPC_MENU_CMD (0x2073).
func (h *Handler) handleActNPCMenuCmd(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	menuAck, _ := pkt.Reader.ReadByte()
	if menuAck == 0 {
		player.OpenNPCMobID = 0
	}
}

// handleScenarioAreaEntryAck maneja PROTO_NC_SCENARIO_AREAENTRY_ACK (0x6C06).
func (h *Handler) handleScenarioAreaEntryAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || pkt.Reader.Remaining() < 32 {
		return
	}
	areaBytes, _ := pkt.Reader.ReadBytes(32)
	areaName := strings.TrimRight(string(areaBytes), "\x00")
	if player.CurrentMap != nil {
		player.CurrentMap.mu.RLock()
		eng := player.CurrentMap.luaEngine
		player.CurrentMap.mu.RUnlock()
		if eng != nil {
			_, _ = eng.CallFunction("AreaEntry", glua.LString(areaName), glua.LNumber(player.MapObjectID))
		}
	}
}

// handleScenarioAreaLeaveAck maneja PROTO_NC_SCENARIO_AREALEAVE_ACK (0x6C08).
func (h *Handler) handleScenarioAreaLeaveAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || pkt.Reader.Remaining() < 32 {
		return
	}
	areaBytes, _ := pkt.Reader.ReadBytes(32)
	areaName := strings.TrimRight(string(areaBytes), "\x00")
	if player.CurrentMap != nil {
		player.CurrentMap.mu.RLock()
		eng := player.CurrentMap.luaEngine
		player.CurrentMap.mu.RUnlock()
		if eng != nil {
			_, _ = eng.CallFunction("AreaLeave", glua.LString(areaName), glua.LNumber(player.MapObjectID))
		}
	}
}

// handleCharTutorialStepAck maneja PROTO_NC_CHAR_TUTORIAL_STEP_ACK (0x1117).
func (h *Handler) handleCharTutorialStepAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	step, _ := pkt.Reader.ReadByte()
	player.Char.TutorialStep = step
	if h.charRepo != nil {
		_ = h.charRepo.UpdateCharacter(player.Char)
	}
	if player.CurrentMap != nil {
		player.CurrentMap.mu.RLock()
		eng := player.CurrentMap.luaEngine
		player.CurrentMap.mu.RUnlock()
		if eng != nil {
			if _, err := eng.CallFunction("Tutorial_Step", glua.LNumber(player.MapObjectID), glua.LNumber(step)); err != nil {
				log.Printf("[Zone][Tutorial] Lua step %d callback info: %v", step, err)
			}
		}
	}
	saveCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_TUTORIAL_STEP_SAVE_ACK)
	saveCmd.WriteUint32(uint32(player.Char.ID))
	saveCmd.WriteByte(step)
	saveCmd.Fill(4, 0)
	_ = client.SendPacket(saveCmd)
}
