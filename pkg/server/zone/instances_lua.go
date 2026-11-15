package zone

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/script/lua"
	glua "github.com/yuin/gopher-lua"
)

// GuildTournamentMatch mantiene el estado vivo de una arena GuildTournament en ZoneMap.
type GuildTournamentMatch struct {
	Seed       uint8
	GuildA     string
	GuildB     string
	GuildAID   uint32
	GuildBID   uint32
	ScoreA     uint32
	ScoreB     uint32
	Flags      map[uint32]uint32 // flagHandle -> ownerGuildID
	Portals    map[uint32]uint32 // portalHandle -> linkedFlagHandle
	Gates      map[uint32]string // gateHandle -> guildName
	Observers  map[int]bool
	StartedAt  time.Time
}

// WeddingCeremonyState mantiene el estado en vivo de una sala de bodas Lua.
type WeddingCeremonyState struct {
	GroomCharNo  int
	BrideCharNo  int
	GroomName    string
	BrideName    string
	ReservedUnix int64
	Completed    bool
}

// GetPlayers retorna una instantánea segura de los jugadores activos en el mapa.
func (m *ZoneMap) GetPlayers() []*ZonePlayer {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*ZonePlayer, 0, len(m.players))
	for _, p := range m.players {
		out = append(out, p)
	}
	return out
}

// GetMap retorna un mapa ya inicializado sin crearlo.
func (mm *MapManager) GetMap(mapID string) *ZoneMap {
	if mm == nil {
		return nil
	}
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	return mm.maps[strings.ToLower(strings.TrimSpace(mapID))]
}

// CreateInstanceMap crea una copia aislada de un ZoneMap y carga su script Lua (ID/KQ/GT/Wedding).
func (mm *MapManager) CreateInstanceMap(baseMapID, instanceKey, kind string, defID uint16, limitSec int, scriptRelPath string) *ZoneMap {
	if mm == nil {
		return nil
	}
	mm.mu.Lock()
	defer mm.mu.Unlock()
	k := strings.ToLower(strings.TrimSpace(instanceKey))
	if existing, ok := mm.maps[k]; ok {
		return existing
	}

	var serial uint16 = defID
	name := baseMapID
	width, height := 20000, 20000
	if mm.fieldStore != nil {
		if entry, ok := mm.fieldStore.GetField(baseMapID); ok && entry != nil {
			serial = uint16(entry.Serial)
			name = entry.Name
			width = entry.Width
			height = entry.Height
		}
	}

	zMap := NewZoneMap(serial, instanceKey, name)
	zMap.Width = width
	zMap.Height = height
	zMap.IsInstance = true
	zMap.InstanceKind = kind
	zMap.InstanceDefID = defID
	zMap.InstanceLimitSec = limitSec
	zMap.InstanceStartedAt = time.Now()
	zMap.dataStore = mm.dataStore
	zMap.handler = mm.handler
	if mm.dataStore != nil {
		zMap.mobCatalog = mm.dataStore.MobCatalog
	}
	if mm.fieldStore != nil {
		if entry, ok := mm.fieldStore.GetField(baseMapID); ok {
			zMap.FieldEntry = entry
		}
	}
	if mm.blockStore != nil {
		mbi := mm.blockStore.GetSHBD(baseMapID)
		doors := mm.blockStore.NewInstanceDoors(baseMapID)
		zMap.SetCollisionData(mbi, doors)
	}

	mm.maps[k] = zMap
	mm.loadMapNPCs(zMap)

	// Inicializar Lua Engine para la instancia si existe el script
	bridge := NewZoneWorldBridge(zMap)
	zMap.worldBridge = bridge
	dataDir := resolveZoneDataDir()
	luaRoot := filepath.Join(dataDir, "LuaScript")
	eng, err := lua.NewEngine(lua.Config{
		BasePath: luaRoot,
		Bridge:   bridge,
	})
	if err == nil {
		zMap.luaEngine = eng
		if scriptRelPath != "" {
			fullScript := filepath.Join(luaRoot, scriptRelPath)
			if _, stErr := os.Stat(fullScript); stErr == nil {
				if loadErr := eng.DoFile(fullScript); loadErr != nil {
					log.Printf("[Zone][Instance] Error cargando script Lua %s: %v", fullScript, loadErr)
				} else {
					_, _ = eng.CallFunction("Main", glua.LString(instanceKey))
				}
			}
		}
	}
	return zMap
}

// ============================================================================
// Implementaciones ZoneMap Bridge para KQ, ID, GT, Bodas e Items Lua
// ============================================================================

func (m *ZoneMap) bridgeShowKQTimerWithLife(seconds int) {
	for _, p := range m.GetPlayers() {
		m.bridgeShowKQTimerWithLifeObj(uint32(p.MapObjectID), seconds)
	}
}

func (m *ZoneMap) bridgeShowKQTimerWithLifeObj(playerHandle uint32, seconds int) {
	p := m.GetPlayer(uint16(playerHandle))
	if p == nil && m.handler != nil {
		if cand := m.handler.findPlayerByHandle(uint16(playerHandle)); cand != nil && cand.CurrentMap == m {
			p = cand
		}
	}
	if p == nil || p.Client == nil {
		return
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_SCENARIO_TIMER_CMD)
	pw.WriteUint16(uint16(seconds))
	_ = p.Client.SendPacket(pw)
}

func (m *ZoneMap) bridgeGetKQLimitSecond() int {
	if m.InstanceLimitSec > 0 {
		rem := m.InstanceLimitSec - int(time.Since(m.InstanceStartedAt).Seconds())
		if rem < 0 {
			return 0
		}
		return rem
	}
	return 1800
}

func (m *ZoneMap) bridgeGetKQTeamType(handle uint32) int {
	if p := m.GetPlayer(uint16(handle)); p != nil {
		return int(p.KQTeam)
	}
	return 0
}

func (m *ZoneMap) bridgeKQRewardIndex(playerHandle uint32, rewardIndex string) {
	p := m.GetPlayer(uint16(playerHandle))
	if p != nil {
		p.KQRewardIdx = rewardIndex
	}
}

func (m *ZoneMap) bridgeScoreBoardAllInMap(show bool, round, redWin, redScore, blueWin, blueScore int) {
	m.mu.Lock()
	m.KQScoreRed = redScore
	m.KQScoreBlue = blueScore
	m.KQWinRed = redWin
	m.KQWinBlue = blueWin
	m.mu.Unlock()

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_SCORE_BOARD_INFO_CMD)
	if show {
		pw.WriteByte(1)
	} else {
		pw.WriteByte(0)
	}
	pw.WriteByte(byte(round))
	pw.WriteUint16(uint16(redScore))
	pw.WriteUint16(uint16(blueScore))
	for _, p := range m.GetPlayers() {
		if p.Client != nil {
			_ = p.Client.SendPacket(pw)
		}
	}
}

func (m *ZoneMap) bridgeIndunRankResult(args []int) {
	// PROTO_NC_INSTANCE_DUNGEON_RANK_RESULT_CMD (0xA41B, sizeof = 5 + N*48 SHINE_INDUN_RANK_CLIENT)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_RESULT_CMD)
	elapsed := 0
	if !m.InstanceStartedAt.IsZero() {
		elapsed = int(time.Since(m.InstanceStartedAt).Seconds())
		if elapsed < 0 {
			elapsed = 0
		}
	}
	hour := byte(elapsed / 3600)
	min := byte((elapsed % 3600) / 60)
	sec := byte(elapsed % 60)
	score := uint64(0)
	if len(args) > 0 {
		score = uint64(args[0])
	}
	players := m.GetPlayers()
	pw.WriteByte(hour)
	pw.WriteByte(min)
	pw.WriteByte(sec)
	pw.WriteUint16(uint16(len(players)))
	for _, pl := range players {
		charName := ""
		var charClass byte
		if pl.Char != nil {
			charName = pl.Char.Name
			charClass = pl.Char.Job
		}
		writeFixedStr(pw, charName, 20) // nCharID (Name5)
		pw.WriteByte(charClass)         // nCharClass
		pw.Fill(3, 0)                   // padding +21..+23 (offset +24 nDamageDealt)
		pw.WriteUint64(score)           // nDamageDealt
		pw.WriteUint64(0)               // nDamageTaken
		pw.WriteUint64(0)               // nHealingDone
	}
	for _, p := range players {
		if p.Client != nil {
			_ = p.Client.SendPacket(pw)
		}
	}
}

func (m *ZoneMap) bridgeRewardAll(rewardType string) {
	if m.handler == nil {
		return
	}
	var kqRew *KQRewardGroup
	if m.dataStore != nil && m.dataStore.KQCatalog != nil {
		if def := m.dataStore.KQCatalog.ByID[m.InstanceDefID]; def != nil {
			kqRew = m.dataStore.KQCatalog.RewardByKey[strings.ToLower(def.RewardIndex)]
		}
	}
	for _, p := range m.GetPlayers() {
		if p.Char == nil {
			continue
		}
		rewGrp := kqRew
		if p.KQRewardIdx != "" && m.dataStore != nil && m.dataStore.KQCatalog != nil {
			if custom := m.dataStore.KQCatalog.RewardByKey[strings.ToLower(p.KQRewardIdx)]; custom != nil {
				rewGrp = custom
			}
		}
		var firstItemID uint16
		if rewGrp != nil && len(rewGrp.Items) > 0 && m.dataStore.ItemStore != nil {
			if info, ok := m.dataStore.ItemStore.GetItemInfoByName(rewGrp.Items[0].ItemInx); ok && info != nil {
				firstItemID = info.ID
				m.handler.addItemToInventory(p, info.ID, 1, nil)
			}
		}
		p.Char.TCharacterState.Fame += 50
		p.Char.TCharacterState.Money += 1000
		if m.handler.charRepo != nil {
			_ = m.handler.charRepo.UpdateCharacter(p.Char)
		}
		if p.Client != nil {
			pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_REWARD_REQ)
			pw.WriteUint32(50)   // Fame
			pw.WriteUint64(1000) // Cen
			var itmBuf [23]byte
			itmBuf[0] = byte(firstItemID & 0xFF)
			itmBuf[1] = byte(firstItemID >> 8)
			pw.WriteBytes(itmBuf[:])
			_ = p.Client.SendPacket(pw)
		}
		if m.handler.gameLogDispatcher != nil {
			m.handler.gameLogDispatcher.EnqueueTypeH(1, p.Char.ID, m.MapID, p.Char.XPos, p.Char.YPos, 0, int(m.InstanceDefID))
		}
	}
}

func (m *ZoneMap) bridgeEndOfKingdomQuest(success bool) {
	op := protocol.PROTO_NC_KQ_COMPLETE_CMD
	if !success {
		op = protocol.PROTO_NC_KQ_FAIL_CMD
	}
	pw := protocol.NewPacketWriter(op)
	pw.WriteUint16(m.InstanceDefID)
	for _, p := range m.GetPlayers() {
		if p.Client != nil {
			_ = p.Client.SendPacket(pw)
		}
	}
}

func (m *ZoneMap) bridgeRewardItem(playerHandle uint32, itemInx string, count int) {
	if m.handler == nil || m.dataStore == nil || m.dataStore.ItemStore == nil {
		return
	}
	info, ok := m.dataStore.ItemStore.GetItemInfoByName(itemInx)
	if !ok || info == nil {
		return
	}
	if count <= 0 {
		count = 1
	}
	if playerHandle == 0 {
		for _, p := range m.GetPlayers() {
			m.handler.addItemToInventory(p, info.ID, uint32(count), nil)
		}
		return
	}
	if p := m.GetPlayer(uint16(playerHandle)); p != nil {
		m.handler.addItemToInventory(p, info.ID, uint32(count), nil)
	}
}

func (m *ZoneMap) bridgeInvenItemDestroy(playerHandle uint32, itemInx string, count int) int {
	if m.handler == nil || m.handler.charRepo == nil || m.dataStore == nil || m.dataStore.ItemStore == nil {
		return 0
	}
	p := m.GetPlayer(uint16(playerHandle))
	info, ok := m.dataStore.ItemStore.GetItemInfoByName(itemInx)
	if p == nil || p.Char == nil || !ok || info == nil {
		return 0
	}
	inv, _ := m.handler.charRepo.GetCharacterItems(p.Char.ID)
	for _, it := range inv {
		if it.StorageType == 9 && it.ItemID == info.ID {
			lot := uint32(it.Flags)
			if lot == 0 {
				lot = 1
			}
			q := uint32(count)
			if count < 0 || q > lot {
				q = lot
			}
			if m.handler.consumeItemUnits(p, it, q) {
				return 1
			}
		}
	}
	return 0
}

func (m *ZoneMap) bridgeGetItemLot(playerHandle uint32, itemInx string) (int, bool) {
	if m.handler == nil || m.handler.charRepo == nil || m.dataStore == nil || m.dataStore.ItemStore == nil {
		return 0, false
	}
	p := m.GetPlayer(uint16(playerHandle))
	info, ok := m.dataStore.ItemStore.GetItemInfoByName(itemInx)
	if p == nil || p.Char == nil || !ok || info == nil {
		return 0, false
	}
	inv, _ := m.handler.charRepo.GetCharacterItems(p.Char.ID)
	total := 0
	for _, it := range inv {
		if it.StorageType == 9 && it.ItemID == info.ID {
			lot := int(it.Flags)
			if lot <= 0 {
				lot = 1
			}
			total += lot
		}
	}
	return total, false
}

func (m *ZoneMap) bridgeDropItem(itemInx string, mobHandle, playerHandle uint32, rate int) {
	if m.handler == nil || m.dataStore == nil || m.dataStore.ItemStore == nil {
		return
	}
	if rate > 0 && rand.Intn(1000) >= rate {
		return
	}
	info, ok := m.dataStore.ItemStore.GetItemInfoByName(itemInx)
	if !ok || info == nil {
		return
	}
	if p := m.GetPlayer(uint16(playerHandle)); p != nil {
		m.handler.addItemToInventory(p, info.ID, 1, nil)
	}
}

func (m *ZoneMap) bridgeScenarioEnd() {
	m.mu.Lock()
	m.IsInstance = false
	m.mu.Unlock()
}

func (m *ZoneMap) bridgeCompleteWedding(groomHandle, brideHandle uint32) {
	if m.handler == nil || m.handler.charRepo == nil {
		return
	}
	g := m.GetPlayer(uint16(groomHandle))
	b := m.GetPlayer(uint16(brideHandle))
	if g != nil && b != nil && g.Char != nil && b.Char != nil {
		_ = m.handler.charRepo.CompleteWedding(g.Char.ID, b.Char.ID)
		if m.handler.gameLogDispatcher != nil {
			m.handler.gameLogDispatcher.EnqueueTypeJ(3, g.Char.ID, m.MapID, g.Char.XPos, g.Char.YPos, 0, b.Char.ID, 0, 1, 0)
		}
	}
}

func (m *ZoneMap) bridgeApplyAbstateByName(handle uint32, staName string, level, durationMs int) {
	if m.dataStore == nil || m.dataStore.AbState == nil {
		return
	}
	if st, ok := m.dataStore.AbState.GetByName(staName); ok && st != nil {
		sec := uint32(durationMs / 1000)
		if sec == 0 && durationMs > 0 {
			sec = 1
		}
		if m.worldBridge != nil {
			m.worldBridge.SetAbstate(handle, st.ID, uint16(level), sec)
		}
	}
}

func (m *ZoneMap) bridgeRegisterGTFlag(npcHandle uint32, redEffect, blueEffect string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.GTState == nil {
		m.GTState = &GuildTournamentMatch{
			Flags:     make(map[uint32]uint32),
			Portals:   make(map[uint32]uint32),
			Gates:     make(map[uint32]string),
			Observers: make(map[int]bool),
			StartedAt: time.Now(),
		}
	}
	m.GTState.Flags[npcHandle] = 0
}

func (m *ZoneMap) bridgeRegisterGTPortal(npcHandle, flagHandle uint32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.GTState == nil {
		m.GTState = &GuildTournamentMatch{
			Flags:     make(map[uint32]uint32),
			Portals:   make(map[uint32]uint32),
			Gates:     make(map[uint32]string),
			Observers: make(map[int]bool),
		}
	}
	m.GTState.Portals[npcHandle] = flagHandle
}

func (m *ZoneMap) bridgeRegisterGTGate(npcHandle uint32, guildName string, flagHandle uint32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.GTState == nil {
		m.GTState = &GuildTournamentMatch{
			Flags:     make(map[uint32]uint32),
			Portals:   make(map[uint32]uint32),
			Gates:     make(map[uint32]string),
			Observers: make(map[int]bool),
		}
	}
	m.GTState.Gates[npcHandle] = guildName
}

func (m *ZoneMap) bridgeEvaluateGTMatchScore() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.GTState == nil {
		return ""
	}
	if m.GTState.ScoreA >= m.GTState.ScoreB {
		return m.GTState.GuildA
	}
	return m.GTState.GuildB
}

func (m *ZoneMap) bridgeRewardGuildTournament(winningGuild string) {
	if m.handler == nil {
		return
	}
	for _, p := range m.GetPlayers() {
		if p.Char != nil && p.Client != nil {
			res := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_RESULT_CMD)
			writeFixedStr(res, winningGuild, 20)
			_ = p.Client.SendPacket(res)
		}
	}
}

// ============================================================================
// Dept 41: Instance Dungeons (2 C2S) & Dept 22: Kingdom Quests en Mapa
// ============================================================================

// handleInstanceDungeonLevelSelectMenuAck maneja PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK (0xA40D).
func (h *Handler) handleInstanceDungeonLevelSelectMenuAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 6 {
		return
	}
	_, _ = pkt.Reader.ReadUint16() // npcHandle
	levelType, _ := pkt.Reader.ReadUint32()

	if h.fieldStore == nil || h.mapManager == nil {
		return
	}
	dungeons := h.fieldStore.GetAllInstanceDungeons()
	if len(dungeons) == 0 {
		return
	}
	chosen := dungeons[0]
	for _, d := range dungeons {
		if d.ModeIDLv == int(levelType) || d.IDNo == int(levelType) {
			chosen = d
			break
		}
	}
	instKey := fmt.Sprintf("%s_id_%d_%d", chosen.MapID, chosen.IDNo, player.Char.ID)
	scriptPath := ""
	if chosen.ScriptName != "" {
		scriptPath = filepath.Join("ID", chosen.ScriptName+".lua")
	}
	instMap := h.mapManager.CreateInstanceMap(chosen.MapID, instKey, "ID", uint16(chosen.IDNo), chosen.LimitTime, scriptPath)
	if instMap != nil {
		h.WarpPlayer(player, instMap.MapID, 1000, 1000, 0)
	}
	// PROTO_NC_INSTANCE_DUNGEON_MAP_REGIST_CMD (0xA413, sizeof = 2 + N*13 SHINE_INDUN_INFO)
	reg := protocol.NewPacketWriter(protocol.PROTO_NC_INSTANCE_DUNGEON_MAP_REGIST_CMD)
	reg.WriteUint16(1) // nMapCount
	writeFixedStr(reg, chosen.MapID, 12)
	reg.WriteByte(byte(chosen.ModeIDLv))
	_ = client.SendPacket(reg)
}

// handleInstanceDungeonRankListReq maneja PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_REQ (0xA417).
func (h *Handler) handleInstanceDungeonRankListReq(client *network.Client, pkt *protocol.Packet) {
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_ACK)
	ack.WriteUint16(0) // nMapCount
	_ = client.SendPacket(ack)
}

// handleKQEntryResponseAck maneja PROTO_NC_KQ_ENTRYRESPONCE_ACK (0x581A).
func (h *Handler) handleKQEntryResponseAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 5 {
		return
	}
	reply, _ := pkt.Reader.ReadByte()
	encHandle, _ := pkt.Reader.ReadUint32()
	if reply == 0 || h.mapManager == nil {
		return
	}

	kqID := uint16(encHandle)
	if kqID == 0 {
		kqID = 1
	}
	baseMap := "KQ_HMiniDragon"
	limitSec := 1800
	scriptPath := ""
	if h.dataStore != nil && h.dataStore.KQCatalog != nil {
		if def := h.dataStore.KQCatalog.ByID[kqID]; def != nil {
			if def.BaseMap != "" {
				baseMap = def.BaseMap
			}
			if def.LimitTimeSec > 0 {
				limitSec = def.LimitTimeSec
			}
			if def.ScriptName != "" {
				scriptPath = filepath.Join("KQ", def.ScriptName, def.ScriptName+".lua")
			}
		}
	}
	instKey := fmt.Sprintf("kq_%d_%d", kqID, time.Now().Unix()/300)
	kqMap := h.mapManager.CreateInstanceMap(baseMap, instKey, "KQ", kqID, limitSec, scriptPath)
	if kqMap != nil {
		if player.KQTeam == 0 {
			player.KQTeam = 1
		}
		h.WarpPlayer(player, kqMap.MapID, 1000, 1000, 0)
		kqMap.AddPlayer(player)
		kqMap.bridgeShowKQTimerWithLifeObj(uint32(player.MapObjectID), limitSec)
	}
	if h.gameLogDispatcher != nil {
		h.gameLogDispatcher.EnqueueTypeH(2, player.Char.ID, instKey, player.Char.XPos, player.Char.YPos, 0, int(kqID))
	}
}

// ============================================================================
// Dept 29: Guild Tournament en Mapa (6 C2S)
// ============================================================================

func (h *Handler) ensureGTMatch(zMap *ZoneMap) *GuildTournamentMatch {
	if zMap == nil {
		return nil
	}
	zMap.mu.Lock()
	defer zMap.mu.Unlock()
	if zMap.GTState == nil {
		zMap.GTState = &GuildTournamentMatch{
			Seed:      1,
			Flags:     make(map[uint32]uint32),
			Portals:   make(map[uint32]uint32),
			Gates:     make(map[uint32]string),
			Observers: make(map[int]bool),
			StartedAt: time.Now(),
		}
	}
	return zMap.GTState
}

// handleGuildTournamentRecallEnterReq maneja PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_REQ (0x749C).
func (h *Handler) handleGuildTournamentRecallEnterReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	accept, _ := pkt.Reader.ReadByte()
	var errCode uint16
	if accept != 0 && h.mapManager != nil {
		gtMapID := "GT_Map01"
		if h.fieldStore != nil {
			for k := range map[string]struct{}{"GT_Map01": {}, "GuildTournament": {}} {
				if _, ok := h.fieldStore.GetGuildTournament(k); ok {
					gtMapID = k
					break
				}
			}
		}
		instMap := h.mapManager.CreateInstanceMap(gtMapID, "gt_active_1", "GT", 1, 1800, filepath.Join("Guild", "GuildTournament.lua"))
		if instMap != nil {
			h.ensureGTMatch(instMap)
			h.WarpPlayer(player, instMap.MapID, 1964, 2124, 0)
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

// handleGuildTournamentFlagCaptureReq maneja PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_REQ (0x74CA).
func (h *Handler) handleGuildTournamentFlagCaptureReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || player.CurrentMap == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	targetFlag, _ := pkt.Reader.ReadUint16()
	gt := h.ensureGTMatch(player.CurrentMap)

	var guildNo uint32 = 1
	if h.charRepo != nil {
		if gd, _, err := h.charRepo.GetGuildByCharID(player.Char.ID); err == nil && gd != nil && gd.ID > 0 {
			guildNo = uint32(gd.ID)
		}
	}
	scoreInc := uint32(100)
	if h.dataStore != nil && h.dataStore.GTCatalog != nil && h.dataStore.GTCatalog.OccupyScore > 0 {
		scoreInc = h.dataStore.GTCatalog.OccupyScore
	}

	player.CurrentMap.mu.Lock()
	gt.Flags[uint32(targetFlag)] = guildNo
	if gt.GuildAID == 0 || gt.GuildAID == guildNo {
		gt.GuildAID = guildNo
		gt.ScoreA += scoreInc
	} else {
		gt.GuildBID = guildNo
		gt.ScoreB += scoreInc
	}
	player.CurrentMap.mu.Unlock()

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	res := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_RESULT_CMD)
	writeFixedStr(res, player.Char.Name, 20)
	writeFixedStr(res, fmt.Sprintf("Flag_%d", targetFlag), 20)
	res.WriteUint32(guildNo)
	res.WriteUint32(guildNo)
	for _, pl := range player.CurrentMap.GetPlayers() {
		if pl.Client != nil {
			_ = pl.Client.SendPacket(res)
		}
	}
}

// handleGuildTournamentUseSkillReq maneja PROTO_NC_GUILD_TOURNAMENT_USESKILL_REQ (0x74CF).
func (h *Handler) handleGuildTournamentUseSkillReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || player.CurrentMap == nil || pkt.Reader.Remaining() < 2 {
		return
	}
	skillNo, _ := pkt.Reader.ReadUint16()
	var errCode uint16
	if h.dataStore != nil && h.dataStore.GTCatalog != nil {
		if sk, ok := h.dataStore.GTCatalog.Skills[uint32(skillNo)]; ok && sk.StaName != "" {
			player.CurrentMap.bridgeApplyAbstateByName(uint32(player.MapObjectID), sk.StaName, 1, 30000)
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_USESKILL_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint16(skillNo)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_USESKILL_CMD)
	cmd.WriteUint32(1)
	cmd.WriteUint16(skillNo)
	for _, pl := range player.CurrentMap.GetPlayers() {
		if pl.Client != nil {
			_ = pl.Client.SendPacket(cmd)
		}
	}
}

// handleGuildTournamentDiceGameThrowReq maneja PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_REQ (0x74D3).
func (h *Handler) handleGuildTournamentDiceGameThrowReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.CurrentMap == nil {
		return
	}
	roll := uint32(rand.Intn(100) + 1)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_CMD)
	cmd.WriteUint32(1)
	cmd.WriteUint32(roll)
	for _, pl := range player.CurrentMap.GetPlayers() {
		if pl.Client != nil {
			_ = pl.Client.SendPacket(cmd)
		}
	}
}

// handleGuildTournamentObserverEnterReq maneja PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_REQ (0x74D7).
func (h *Handler) handleGuildTournamentObserverEnterReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	seed, _ := pkt.Reader.ReadByte()
	player.GTObserving = true
	if player.CurrentMap != nil {
		gt := h.ensureGTMatch(player.CurrentMap)
		player.CurrentMap.mu.Lock()
		gt.Observers[player.Char.ID] = true
		player.CurrentMap.mu.Unlock()
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_ACK)
	ack.WriteUint16(0)
	ack.WriteByte(seed)
	ack.WriteUint32(1)
	ack.WriteUint32(2)
	_ = client.SendPacket(ack)
}

// handleGuildTournamentObserverOutReq maneja PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_REQ (0x74DB).
func (h *Handler) handleGuildTournamentObserverOutReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || pkt.Reader.Remaining() < 1 {
		return
	}
	seed, _ := pkt.Reader.ReadByte()
	player.GTObserving = false
	if player.CurrentMap != nil && player.CurrentMap.GTState != nil {
		player.CurrentMap.mu.Lock()
		delete(player.CurrentMap.GTState.Observers, player.Char.ID)
		player.CurrentMap.mu.Unlock()
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_ACK)
	ack.WriteUint16(0)
	ack.WriteByte(seed)
	_ = client.SendPacket(ack)
	h.WarpPlayer(player, "Eld", 17214, 13445, 0)
}

// ============================================================================
// Dept 8 & Dept 4: Bodas (8 C2S + Partner Info) y Dept 37: HolyPromise Cen
// ============================================================================

var (
	pendingWeddingMu       sync.Mutex
	pendingWeddingProposals = make(map[int]int) // targetCharID -> proposerCharID
	pendingDivorces         = make(map[int]int) // targetCharID -> requesterCharID
)

// handleActWeddingProposeReq maneja PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ (0x204F).
func (h *Handler) handleActWeddingProposeReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || player.CurrentMap == nil || pkt.Reader.Remaining() < 39 {
		return
	}
	toHandle, _ := pkt.Reader.ReadUint16()
	wordBytes, _ := pkt.Reader.ReadBytes(37)
	target := player.CurrentMap.GetPlayer(toHandle)
	if target == nil || target.Char == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK)
		ack.WriteUint16(toHandle)
		ack.WriteUint16(1)
		ack.WriteBytes(wordBytes)
		_ = client.SendPacket(ack)
		return
	}
	pendingWeddingMu.Lock()
	pendingWeddingProposals[target.Char.ID] = player.Char.ID
	pendingWeddingMu.Unlock()

	if target.Client != nil {
		fwd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WEDDING_PROPOSEACK_REQ)
		fwd.WriteUint16(player.MapObjectID)
		fwd.WriteBytes(wordBytes)
		_ = target.Client.SendPacket(fwd)
	}
}

// handleActWeddingProposeAckAck maneja PROTO_NC_ACT_WEDDING_PROPOSEACK_ACK (0x2051).
func (h *Handler) handleActWeddingProposeAckAck(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil || player.CurrentMap == nil || pkt.Reader.Remaining() < 40 {
		return
	}
	fromHandle, _ := pkt.Reader.ReadUint16()
	reply, _ := pkt.Reader.ReadByte()
	respBytes, _ := pkt.Reader.ReadBytes(37)

	pendingWeddingMu.Lock()
	proposerID := pendingWeddingProposals[player.Char.ID]
	delete(pendingWeddingProposals, player.Char.ID)
	pendingWeddingMu.Unlock()

	if reply != 0 && h.charRepo != nil && proposerID > 0 {
		_ = h.charRepo.ProposeWedding(proposerID, player.Char.ID)
		if h.gameLogDispatcher != nil {
			h.gameLogDispatcher.EnqueueTypeJ(1, proposerID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, player.Char.ID, 0, 1, 0)
		}
	}
	if proposer := player.CurrentMap.GetPlayer(fromHandle); proposer != nil && proposer.Client != nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_ACK)
		ack.WriteUint16(player.MapObjectID)
		if reply != 0 {
			ack.WriteUint16(0)
		} else {
			ack.WriteUint16(1)
		}
		ack.WriteBytes(respBytes)
		_ = proposer.Client.SendPacket(ack)
	}
}

// handleActWeddingHallReservReq maneja PROTO_NC_ACT_WEDDING_HALL_RESERV_REQ (0x2053).
func (h *Handler) handleActWeddingHallReservReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil {
		return
	}
	nowUnix := time.Now().Unix()
	if h.charRepo != nil {
		if st, err := h.charRepo.GetWeddingStatus(player.Char.ID); err == nil && st != nil {
			_ = h.charRepo.SetWeddingStatus(player.Char.ID, st.PartnerCharNo, st.Status, 1, time.Unix(nowUnix, 0), time.Unix(nowUnix, 0), st.MarriageCount, st.DivorceCancel)
		}
	}
	if h.mapManager != nil {
		baseWedMap := ""
		if player.CurrentMap != nil {
			baseWedMap = player.CurrentMap.MapID
		}
		if h.fieldStore != nil {
			if fw, ok := h.fieldStore.GetField("Wedding"); ok && fw != nil {
				baseWedMap = fw.MapID
			} else if baseWedMap == "" {
				if f0, ok := h.fieldStore.GetFieldBySerial(0); ok && f0 != nil {
					baseWedMap = f0.MapID
				}
			}
		}
		wedMap := h.mapManager.CreateInstanceMap(baseWedMap, fmt.Sprintf("wedding_%d", player.Char.ID), "Wedding", 1, 900, filepath.Join("Wedding", "Wedding.lua"))
		if wedMap != nil {
			wedMap.mu.Lock()
			wedMap.WeddingState = &WeddingCeremonyState{
				GroomCharNo:  player.Char.ID,
				GroomName:    player.Char.Name,
				ReservedUnix: nowUnix,
			}
			wedMap.mu.Unlock()
		}
	}
	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueTypeJ(2, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, 0, nowUnix, 1, 0)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WEDDING_HALL_RESERV_ACK)
	ack.WriteUint16(0)
	ack.WriteInt64(nowUnix)
	writeFixedStr(ack, time.Unix(nowUnix, 0).Format(time.RFC3339), 36)
	_ = client.SendPacket(ack)
}

// handleActWeddingHallGuestEnterReadyReq maneja PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_REQ (0x2056).
func (h *Handler) handleActWeddingHallGuestEnterReadyReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil {
		return
	}
	now := time.Now().Unix()
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_ACK)
	ack.WriteUint16(0)
	writeFixedStr(ack, player.Char.Name, 20)
	writeFixedStr(ack, "", 20)
	ack.WriteInt64(now)
	ack.WriteInt64(now + 60)
	ack.WriteInt64(now + 900)
	ack.Fill(174-66, 0)
	_ = client.SendPacket(ack)
}

// handleActWeddingHallGuestEnterReq maneja PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_REQ (0x2058).
func (h *Handler) handleActWeddingHallGuestEnterReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// handleActWeddingAgreementDivorceReq maneja PROTO_NC_ACT_WEDDING_AGREEMENT_DIVORCE_REQ (0x205B).
func (h *Handler) handleActWeddingAgreementDivorceReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil {
		return
	}
	if h.charRepo != nil {
		if st, err := h.charRepo.GetWeddingStatus(player.Char.ID); err == nil && st != nil && st.PartnerCharNo > 0 {
			_ = h.charRepo.RequestAgreeDivorce(player.Char.ID, st.PartnerCharNo)
		}
	}
	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueTypeJ(4, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, 0, 0, 1, 0)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WEDDING_DIVORCE_REQ_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// handleActWeddingCompulsoryDivorceReq maneja PROTO_NC_ACT_WEDDING_COMPULSORY_DIVORCE_REQ (0x205C).
func (h *Handler) handleActWeddingCompulsoryDivorceReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil {
		return
	}
	if h.charRepo != nil {
		if st, err := h.charRepo.GetWeddingStatus(player.Char.ID); err == nil && st != nil && st.PartnerCharNo > 0 {
			_ = h.charRepo.RequestForceDivorce(player.Char.ID, st.PartnerCharNo)
		}
	}
	if h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueTypeJ(5, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, 0, 0, 2, 0)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WEDDING_DIVORCE_REQ_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

// handleActWeddingDivorceRefuseCmd maneja PROTO_NC_ACT_WEDDING_DIVORCE_REFUSE_CMD (0x205E).
func (h *Handler) handleActWeddingDivorceRefuseCmd(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil {
		return
	}
	pendingWeddingMu.Lock()
	delete(pendingDivorces, player.Char.ID)
	pendingWeddingMu.Unlock()
}

// handleCharWeddingPartnerInfoReq maneja PROTO_NC_CHAR_WEDDING_PARTNER_INFO_REQ (0x1093).
func (h *Handler) handleCharWeddingPartnerInfoReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil {
		return
	}
	var partnerName string
	if h.charRepo != nil {
		if st, err := h.charRepo.GetWeddingStatus(player.Char.ID); err == nil && st != nil && st.PartnerCharNo > 0 {
			if pCh, errCh := h.charRepo.GetCharacterByID(st.PartnerCharNo); errCh == nil && pCh != nil {
				partnerName = pCh.Name
			}
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK)
	ack.WriteUint16(0)
	writeFixedStr(ack, partnerName, 20)
	_ = client.SendPacket(ack)
}

// handleHolyPromiseGetCenRewardReq maneja PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_REQ (0x943C).
func (h *Handler) handleHolyPromiseGetCenRewardReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil {
		return
	}
	var money uint64
	var errCode uint16
	if h.charRepo != nil {
		m, err := h.charRepo.GetHolyPromiseRemainMoney(player.Char.ID)
		if err != nil {
			errCode = 1
		} else if m > 0 {
			money = uint64(m)
		}
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint64(money)
	_ = client.SendPacket(ack)
}

// handleHolyPromiseWithdrawCenRewardReq maneja PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_REQ (0x9440).
func (h *Handler) handleHolyPromiseWithdrawCenRewardReq(client *network.Client, pkt *protocol.Packet) {
	player := h.GetPlayer(client)
	if player == nil || player.Char == nil {
		return
	}
	var money uint64
	var errCode uint16
	if h.charRepo != nil {
		m, err := h.charRepo.WithdrawHolyPromiseRemainMoney(player.Char.ID)
		if err != nil {
			errCode = 1
		} else if m > 0 {
			money = uint64(m)
			player.Char.TCharacterState.Money += m
			_ = h.charRepo.SaveCharMoney(player.Char.ID, player.Char.TCharacterState.Money)
		}
	}
	if errCode == 0 && h.gameLogDispatcher != nil && player.CurrentMap != nil {
		h.gameLogDispatcher.EnqueueTypeK(1, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(money), 0, player.Char.TCharacterState.Money)
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_ACK)
	ack.WriteUint16(errCode)
	ack.WriteUint64(money)
	_ = client.SendPacket(ack)
}
