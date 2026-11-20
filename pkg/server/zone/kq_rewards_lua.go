package zone

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/data/kq"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/script/lua"
	glua "github.com/yuin/gopher-lua"
)

// ============================================================================
// Sub-fase 27D: Kingdom Quest Lua, Recompensas (Zone.exe + Character.exe),
//               Guardado/Restauración de Posición KQ y GameLog KQ_ENTER/LEAVE
// ============================================================================

const (
	// Constantes de sp_KQReward (Zone.exe @ 0x0052DA70) y TreasureChestMaker (0x00595D00)
	kqRewardRollMax       = 0x3E8 // exe: 0x0052DDD0 (well512_GetRandom(0x3E8) = 1000)
	kqRewardSlots         = 0x0F  // exe: 0x0052E031 (15 slots en KINGDOM_QUEST_REW)
	kqRewardChestMaxItems = 8     // exe: 0x00595D5F (cmp dword ptr [esi+0x3e8], 8)
	kqRewardChestType     = 7     // exe: 0x0052DD90 (push 7 -> TreasureChestMaker)

	// Códigos de error oficiales de recompensa KQ (Zone.exe & Character.exe)
	kqErrRewardItemCreate = uint16(0x1101) // exe: 0x00432224 (Character.exe fc_NC_KQ_REWARD_REQ)
	kqErrRewardMoneyAdd   = uint16(0x1102) // exe: 0x0043225D (Character.exe fc_NC_KQ_REWARD_REQ)
	kqErrRewardFameAdd    = uint16(0x1103) // exe: 0x00432290 (Character.exe fc_NC_KQ_REWARD_REQ)
	kqErrRewardNoMobKill  = uint16(0x1104) // exe: 0x004284D5 (Zone.exe AxialListKQRewardLua::ali_Work -> 0x5823)

	// Máscaras de limpieza de mapa y comandos de fin KQ
	kqClearObjectMaskEndLua      = 0xB0  // exe: 0x005DF2EE (cEndOfKingdomQuest) / 0x004F60B2 (ShineEndOfKingdomQuest::sa_Step)
	kqClearObjectMaskQuestResult = 0x1B0 // exe: 0x005E323D (cQuestResult)
	kqEndDept                    = 0x16  // exe: 0x005E31C0 / 0x005E320F (Dept 22 = NC_KQ)
	kqEndCmdComplete             = 0x12  // exe: 0x005E31BE (0x5812 PROTO_NC_KQ_COMPLETE_CMD)
	kqEndCmdFail                 = 0x13  // exe: 0x005E320D (0x5813 PROTO_NC_KQ_FAIL_CMD)

	// Constantes de so_SaveLocation (Zone.exe @ 0x004542F0)
	kqSaveCoordMinValid = 0x0A // exe: 0x004546CE (cmp coord, 0xa -> fallback RouN 0x192D, 0x21B6)

	// Constantes de cScoreBoard_AllInMap (Zone.exe @ 0x005EDAD0) y cScoreInfo_AllInMap (0x005E3970)
	kqScoreBoardMaxWin   = 5    // exe: 0x005EDC55 / 0x005EDC6F
	kqScoreBoardMaxScore = 0x63 // exe: 0x005EDC65 / 0x005EDC79 (99)
	kqScoreInfoMaxCount  = 2    // exe: 0x005E3A31 (cmp ebx, 2)

	// Tipos de título oficiales (CharacterTitleData.shn) para CT_KQSuccess / CT_KQFail
	kqTitleTypeKQJoin    = uint8(23) // exe: 0x005CBF90 (CCharacterTitleZone::CT_KQJoin)
	kqTitleTypeKQSuccess = uint8(24) // exe: 0x005CBFB0 (CCharacterTitleZone::CT_KQSuccess)
	kqTitleTypeKQFail    = uint8(25) // exe: 0x005CBFD0 (CCharacterTitleZone::CT_KQFail)
)

// ============================================================================
// KQContributeList: kqcl_FindKQ / kqcl_GetMobKill / kqcl_SendMobKill
// ============================================================================

// kqcl_FindKQ ports KQContributeList::kqcl_FindKQ (exe: 0x00499CA0).
func (c *KQContributeList) kqcl_FindKQ(handle uint32) (map[uint32]int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m, ok := c.active[handle]
	return m, ok
}

// kqcl_GetMobKill ports KQContributeList::kqcl_GetMobKill (exe: 0x00499D80).
func (c *KQContributeList) kqcl_GetMobKill(handle uint32, charRegNum uint32) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m, ok := c.active[handle]
	if !ok {
		return 0
	}
	return m[charRegNum]
}

// kqcl_SendMobKill ports KQContributeList::kqcl_SendMobKill (exe: 0x00499FF0).
// Registra contribución de muertes de monstruos dentro de la instancia KQ activa.
func (c *KQContributeList) kqcl_SendMobKill(handle uint32, charRegNum uint32, mobKillDelta int) {
	if handle == kqInvalidHandle || charRegNum == 0 || mobKillDelta <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active == nil {
		c.active = make(map[uint32]map[uint32]int)
	}
	m, ok := c.active[handle]
	if !ok {
		m = make(map[uint32]int)
		c.active[handle] = m
	}
	m[charRegNum] += mobKillDelta
}

// RecordKQMobKill actualiza KQContributeList cuando un jugador mata un mob en un mapa KQ (exe: 0x00499FF0).
func (h *Handler) RecordKQMobKill(zMap *ZoneMap, killer *ZonePlayer) {
	if h == nil || zMap == nil || killer == nil || killer.Char == nil {
		return
	}
	kqHandle := zMap.fm_GetKQhandle()
	if kqHandle == kqInvalidHandle {
		return
	}
	h.kqState().contribute.kqcl_SendMobKill(kqHandle, uint32(killer.Char.ID), 1)
}

// ============================================================================
// Pending Reward Orders (gds_NC_KQ_REWARDSUC_ACK / gds_NC_KQ_REWARDFAIL_ACK)
// ============================================================================

// KQRewardItemSpec describe un ítem generado por TreasureChestMaker::tcm_ItemMake (exe: 0x00595D00).
type KQRewardItemSpec struct {
	ItemID   uint16
	ItemInx  string
	Quantity uint32
	ItemKey  int64 // clave asignada al persistir (para rollback si falla Money/Fame)
}

// KQRewardOrder representa una orden de recompensa KQ enviada vía PROTO_NC_KQ_REWARD_REQ (0x5815).
type KQRewardOrder struct {
	Handle     uint16
	CharNo     uint32
	OrderIndex uint16
	RewardID   uint32
	RewardName string
	BoxItemID  uint16
	Exp        uint32
	Money      uint64
	Fame       uint32
	Items      []KQRewardItemSpec
	Committed  bool
	Failed     bool
	FailCode   uint16
}

// KQRewardTracker almacena órdenes pendientes y resultados recientes para verificación 1:1.
type KQRewardTracker struct {
	mu       sync.Mutex
	pending  map[uint64]*KQRewardOrder // key = (uint64(charNo)<<16) | uint64(orderIndex)
	history  []*KQRewardOrder
	randFunc func(n int) int // opcional para tests determinísticos
}

func newKQRewardTracker() *KQRewardTracker {
	return &KQRewardTracker{
		pending: make(map[uint64]*KQRewardOrder),
	}
}

func kqRewardOrderKey(charNo uint32, orderIdx uint16) uint64 {
	return (uint64(charNo) << 16) | uint64(orderIdx)
}

func (rt *KQRewardTracker) roll(n int) int {
	rt.mu.Lock()
	fn := rt.randFunc
	rt.mu.Unlock()
	if fn != nil {
		return fn(n)
	}
	return rand.Intn(n)
}

// SetKQRewardRand configura un generador determinístico para tests de sp_KQReward (exe: 0x0052DDD0).
func (h *Handler) SetKQRewardRand(fn func(n int) int) {
	st := h.kqState()
	st.mu.Lock()
	if st.rewardTracker == nil {
		st.rewardTracker = newKQRewardTracker()
	}
	st.rewardTracker.mu.Lock()
	st.rewardTracker.randFunc = fn
	st.rewardTracker.mu.Unlock()
	st.mu.Unlock()
}

func (h *Handler) kqRewardTracker() *KQRewardTracker {
	st := h.kqState()
	st.mu.Lock()
	defer st.mu.Unlock()
	if st.rewardTracker == nil {
		st.rewardTracker = newKQRewardTracker()
	}
	return st.rewardTracker
}

// LastKQRewardOrders retorna una copia de las órdenes de recompensa procesadas (diagnóstico/tests).
func (h *Handler) LastKQRewardOrders() []*KQRewardOrder {
	rt := h.kqRewardTracker()
	rt.mu.Lock()
	defer rt.mu.Unlock()
	out := make([]*KQRewardOrder, len(rt.history))
	copy(out, rt.history)
	return out
}

// ============================================================================
// AxialListKQReward & AxialListKQRewardLua (Zone.exe 0x00428410 / 0x004288C0 / 0x00428490)
// ============================================================================

// AxialListKQReward ports AxialListKQReward (exe: 0x00428410) & AxialListKQRewardLua (exe: 0x004288C0).
// En Zone.exe ambos comparten ali_Work @ 0x00428490 (MSVC ICF).
type AxialListKQReward struct {
	Handler       *Handler
	Elem          *KQElement // +0x04 (kqc_SearchHandle)
	DemandMobKill int        // +0x08 (elem->Info.DemandMobKill @ +0x99)
}

// NewAxialListKQReward ports AxialListKQReward::AxialListKQReward (exe: 0x00428410)
// and AxialListKQRewardLua::AxialListKQRewardLua (exe: 0x004288C0).
func (h *Handler) NewAxialListKQReward(kqHandle uint32) *AxialListKQReward {
	elem, _, ok := h.kqState().container.kqc_SearchHandle(kqHandle) // exe: 0x00428453 / 0x00428903
	a := &AxialListKQReward{
		Handler: h,
		Elem:    elem,
	}
	if ok && elem != nil {
		a.DemandMobKill = int(elem.Info.DemandMobKill) // exe: 0x0042845F / 0x0042890F ([eax+0x99])
	}
	return a
}

// ali_Work ports AxialListKQRewardLua::ali_Work (exe: 0x00428490).
func (a *AxialListKQReward) ali_Work(player *ZonePlayer) bool {
	if a == nil || a.Handler == nil || a.Elem == nil || player == nil || player.Char == nil {
		return true // exe: 0x0042849B
	}
	kills := a.Handler.kqState().contribute.kqcl_GetMobKill(a.Elem.Info.Handle, uint32(player.Char.ID)) // exe: 0x004284AB
	if kills >= a.DemandMobKill {                                                                        // exe: 0x004284B2
		a.Handler.so_ply_KQRewardStruct(player, a.Elem) // exe: 0x004284C3 -> 0x0052E9C0
		return true
	}
	// exe: 0x004284CD..0x004284DE: push 0x1104; push 0x23; push 0x16 -> PROTO_NC_KQ_NOREWARD_CMD (0x5823)
	if player.Client != nil {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_KQ_NOREWARD_CMD)
		pw.WriteUint16(kqErrRewardNoMobKill)
		_ = player.Client.SendPacket(pw)
	}
	return true
}

// RunInMap aplica AxialListKQReward(Lua) sobre todos los jugadores del mapa (so_AllInMap @ 0x005DF19B).
func (a *AxialListKQReward) RunInMap(zMap *ZoneMap) {
	if a == nil || zMap == nil {
		return
	}
	for _, p := range zMap.GetPlayers() {
		a.ali_Work(p)
	}
}

// ============================================================================
// ShinePlayer::so_ply_KQRewardStruct (0x0052E9C0), so_ply_KQRewardIndex (0x0052EA80),
// ShinePlayer::sp_KQReward (0x0052DA70),
// Character.exe CPFsCharacter::fc_NC_KQ_REWARD_REQ (0x004321D0),
// GameDBSession::gds_NC_KQ_REWARDSUC_ACK (0x0052E570) & gds_NC_KQ_REWARDFAIL_ACK (0x0052E680)
// ============================================================================

// so_ply_KQRewardStruct ports ShinePlayer::so_ply_KQRewardStruct (exe: 0x0052E9C0).
func (h *Handler) so_ply_KQRewardStruct(player *ZonePlayer, elem *KQElement) {
	if h == nil || player == nil || elem == nil {
		return
	}
	if h.dataStore == nil || h.dataStore.KQStore == nil {
		return
	}
	// exe: 0x0052EA26..0x0052EA3A: g_KQRewardDataBox[uint16(elem->Info.RewardIndex)] (0x0049A950)
	rew, ok := h.dataStore.KQStore.GetReward(uint32(elem.Info.RewardIndex))
	if !ok || rew == nil {
		return // exe: 0x0052EA41
	}
	if idx := rew.IndexName(); idx != "" {
		player.KQRewardIdx = idx
	} else {
		player.KQRewardIdx = strconv.Itoa(int(rew.ID))
	}
	h.sp_KQReward(player, rew) // exe: 0x0052EA46 -> 0x0052DA70
}

// so_ply_KQRewardIndex ports ShinePlayer::so_ply_KQRewardIndex (exe: 0x0052EA80).
func (h *Handler) so_ply_KQRewardIndex(player *ZonePlayer, rewardIndex string) {
	if h == nil || player == nil || rewardIndex == "" {
		return
	}
	if h.dataStore == nil || h.dataStore.KQStore == nil {
		return
	}
	// exe: 0x0052EAF3: g_KQRewardDataBox[rewardIndex] (0x0049A9A0, exact IndexString match)
	rew, ok := h.dataStore.KQStore.GetRewardByName(rewardIndex)
	if !ok || rew == nil {
		// Fallback numérico solo si rewardIndex es un ID entero explícito (p. ej. "1")
		if id, err := strconv.ParseUint(strings.TrimSpace(rewardIndex), 10, 32); err == nil {
			rew, ok = h.dataStore.KQStore.GetReward(uint32(id))
		}
	}
	if !ok || rew == nil {
		return // exe: 0x0052EAFA
	}
	player.KQRewardIdx = rew.IndexName()
	h.sp_KQReward(player, rew) // exe: 0x0052EAFF -> 0x0052DA70
}

// igc_Getitem ports ItemGroupClassifier::igc_Getitem (Zone.exe @ 0x004903E0).
// Primero busca el ítem directo en ItemInfo.shn; si no existe, busca en ItemDropGroup.json
// filtrando por clase equipable si corresponde.
func (h *Handler) igc_Getitem(argName string, playerClass byte) (uint16, string, bool) {
	name := strings.TrimSpace(argName)
	if name == "" || h == nil || h.dataStore == nil || h.dataStore.ItemStore == nil {
		return 0xFFFF, "", false
	}
	// exe: 0x00490450: direct lookup in ItemInfo BTree
	if info, ok := h.dataStore.ItemStore.GetItemInfoByName(name); ok && info != nil {
		return info.ID, info.InxName, true
	}
	// exe: 0x0049049C..0x004904DF: lookup in ItemGroupClassifier (ItemDropGroup)
	if h.dataStore.DropStore != nil {
		if entries, ok := h.dataStore.DropStore.GetGroup(name); ok && len(entries) > 0 {
			rt := h.kqRewardTracker()
			// Filtrar por clase equipable si ItemInfo provee WhoEquip
			var classMatched []string
			var allValid []string
			for _, ent := range entries {
				if info, found := h.dataStore.ItemStore.GetItemInfoByName(ent.ItemID); found && info != nil {
					allValid = append(allValid, info.InxName)
					if info.UseClass == 0 || classMatchesWhoEquip(playerClass, info.UseClass) {
						classMatched = append(classMatched, info.InxName)
					}
				}
			}
			pool := classMatched
			if len(pool) == 0 {
				pool = allValid
			}
			if len(pool) > 0 {
				picked := pool[rt.roll(len(pool))]
				if info, found := h.dataStore.ItemStore.GetItemInfoByName(picked); found && info != nil {
					return info.ID, info.InxName, true
				}
			}
		}
	}
	return 0xFFFF, "", false
}

// sp_GetItemWhoEquip_ClassGroup mirrors ShinePlayer::sp_GetItemWhoEquip_ClassGroup (exe: 0x005630F0).
func classMatchesWhoEquip(job byte, whoEquip uint32) bool {
	if whoEquip == 0 || job == 0 {
		return true
	}
	// Base class bitmask: Fighter(1..5), Cleric(6..10), Archer(11..15), Mage(16..20), Trickster(21..25), Crusader(26..27)
	bit := uint32(1) << (job - 1)
	return (whoEquip & bit) != 0
}

// sp_KQReward ports ShinePlayer::sp_KQReward (Zone.exe @ 0x0052DA70).
func (h *Handler) sp_KQReward(player *ZonePlayer, rew *kq.KINGDOM_QUEST_REW) *KQRewardOrder {
	if h == nil || player == nil || player.Char == nil || rew == nil {
		return nil
	}
	// exe: 0x0052DCC8..0x0052DCEA: skip if state == 0 or state == 3 (Dead)
	if player.State == 0 || player.State == 3 {
		return nil
	}

	// exe: 0x0052DCF0..0x0052DD89: resolve KQBoxItemIDX in ItemInfo
	var boxItemID uint16 = 0xFFFF
	if h.dataStore != nil && h.dataStore.ItemStore != nil {
		if boxInfo, ok := h.dataStore.ItemStore.GetItemInfoByName(rew.BoxItemIndex()); ok && boxInfo != nil {
			boxItemID = boxInfo.ID
		}
	}

	rt := h.kqRewardTracker()
	var (
		totalExp   uint32
		totalMoney uint64
		totalFame  uint32
		items      []KQRewardItemSpec
	)

	// exe: 0x0052DDB5..0x0052E034: loop i = 0 .. 14 (0x0F slots)
	for i := 0; i < kqRewardSlots; i++ {
		rate := rew.RewardRate[i]
		roll := rt.roll(kqRewardRollMax) // exe: 0x0052DDD0 (well512_GetRandom(0x3E8))
		if roll >= int(rate) {           // exe: 0x0052DDE2 (cmp ax, word ptr [edi+0x1e]; jb -> take)
			continue
		}
		if h.dataStore == nil || h.dataStore.KQStore == nil {
			continue
		}
		sr, ok := h.dataStore.KQStore.GetShineReward(rew.Reward[i]) // exe: 0x0052DED8 (rd_FindHandle @ 0x0064C850)
		if !ok || sr == nil {
			continue // exe: 0x0052DEE1 (Skiped by Null info)
		}

		// exe: 0x0052DF3C..0x0052DFDB (jump table @ 0x0052E348)
		switch sr.RewardType {
		case kq.ShineRewardTypeNone: // exe: 0x0052DF50
			continue
		case kq.ShineRewardTypeItem: // exe: 0x0052DF7D (tcm_ItemMake @ 0x00595D00)
			if len(items) >= kqRewardChestMaxItems { // exe: 0x00595D5F (Too many content)
				continue
			}
			itemID, itemInx, found := h.igc_Getitem(sr.ArgumentString(), player.Char.Job) // exe: 0x00595D9A
			if !found || itemID == 0xFFFF {
				continue
			}
			qty := sr.Quantity
			if qty == 0 {
				qty = 1
			}
			items = append(items, KQRewardItemSpec{
				ItemID:   itemID,
				ItemInx:  itemInx,
				Quantity: qty,
			})
		case kq.ShineRewardTypeExp: // exe: 0x0052DFA3
			totalExp += sr.Quantity
		case kq.ShineRewardTypeCen: // exe: 0x0052DFBF
			totalMoney += uint64(sr.Quantity)
		case kq.ShineRewardTypeFame: // exe: 0x0052DFDB
			totalFame += sr.Quantity
		}
	}

	// exe: 0x0052E0D5..0x0052E239: asignar orderIndex y enviar PROTO_NC_KQ_REWARD_REQ (0x5815) a Character.exe
	orderIdx := player.KQRewardOrderSeq
	player.KQRewardOrderSeq++ // exe: 0x0052E235 (add word ptr [eax+4], 1)

	order := &KQRewardOrder{
		Handle:     player.MapObjectID,
		CharNo:     uint32(player.Char.ID),
		OrderIndex: orderIdx,
		RewardID:   rew.ID,
		RewardName: rew.IndexName(),
		BoxItemID:  boxItemID,
		Exp:        totalExp,
		Money:      totalMoney,
		Fame:       totalFame,
		Items:      items,
	}

	rt.mu.Lock()
	rt.pending[kqRewardOrderKey(order.CharNo, order.OrderIndex)] = order
	rt.history = append(rt.history, order)
	rt.mu.Unlock()

	// exe: 0x0052E1BD -> Character.exe CPFsCharacter::fc_NC_KQ_REWARD_REQ (0x004321D0)
	h.fc_NC_KQ_REWARD_REQ(player, order)

	// exe: 0x0052E239..0x0052E24C: sp_GainExp(totalExp, 0xFFFF, 0xFFFF)
	if totalExp > 0 {
		mapID := ""
		if player.CurrentMap != nil {
			mapID = player.CurrentMap.MapID
		}
		h.AddPlayerExpAndCheckLevelUp(player, uint64(totalExp), 0xFFFF, 0xFFFF, int32(player.Char.XPos), int32(player.Char.YPos), mapID)
	}
	return order
}

// fc_NC_KQ_REWARD_REQ ports CPFsCharacter::fc_NC_KQ_REWARD_REQ (Character.exe @ 0x004321D0).
// Ejecuta la transacción de recompensa en orden exacto con rollback en cada etapa:
//  1. Crear ítems (fallo -> 0x1101 -> gds_NC_KQ_REWARDFAIL_ACK 0x5817)
//  2. Si money > 0: p_Char_MoneyAdd (fallo -> rollback ítems -> 0x1102 -> 0x5817)
//  3. Si fame > 0: p_Char_FameAdd (fallo -> rollback money via p_Char_MoneySub + ítems -> 0x1103 -> 0x5817)
//  4. Emitir GameLog money/fame y confirmar vía gds_NC_KQ_REWARDSUC_ACK (0x5816).
func (h *Handler) fc_NC_KQ_REWARD_REQ(player *ZonePlayer, order *KQRewardOrder) {
	if h == nil || player == nil || player.Char == nil || order == nil {
		return
	}
	charNo := int(order.CharNo)

	// Paso 1: Crear ítems de recompensa (exe: 0x00432215..0x00432232)
	createdKeys := make([]int64, 0, len(order.Items))
	for i := range order.Items {
		spec := &order.Items[i]
		added, okAdd := h.addItemToInventory(player, spec.ItemID, spec.Quantity, nil)
		if !okAdd || added == nil {
			// Rollback de ítems ya creados (exe: 0x004322B5..0x00432337)
			h.rollbackKQRewardItems(player, createdKeys)
			h.gds_NC_KQ_REWARDFAIL_ACK(order.Handle, order.CharNo, order.OrderIndex, kqErrRewardItemCreate)
			return
		}
		spec.ItemKey = added.ItemKey
		createdKeys = append(createdKeys, added.ItemKey)
	}

	// Paso 2: Si money > 0 -> p_Char_MoneyAdd (exe: 0x00432233..0x00432269)
	if order.Money > 0 && h.charRepo != nil {
		if err := h.charRepo.AddCharMoney(charNo, int64(order.Money)); err != nil {
			h.rollbackKQRewardItems(player, createdKeys)
			h.gds_NC_KQ_REWARDFAIL_ACK(order.Handle, order.CharNo, order.OrderIndex, kqErrRewardMoneyAdd)
			return
		}
	}

	// Paso 3: Si fame > 0 -> p_Char_FameAdd (exe: 0x0043226B..0x004322B4)
	if order.Fame > 0 && h.charRepo != nil {
		if err := h.charRepo.AddCharFame(charNo, int(order.Fame)); err != nil {
			// Rollback de money vía p_Char_MoneySub (exe: 0x004322A2..0x004322B0 -> 0x00402C60)
			if order.Money > 0 {
				_, _ = h.charRepo.SubCharMoney(charNo, int64(order.Money))
			}
			// Rollback de ítems (exe: 0x004322B5..0x00432337)
			h.rollbackKQRewardItems(player, createdKeys)
			h.gds_NC_KQ_REWARDFAIL_ACK(order.Handle, order.CharNo, order.OrderIndex, kqErrRewardFameAdd)
			return
		}
	}

	// Paso 4: GameLog para Money (0x0519) y Fame (0x051C) (exe: 0x00432386..0x00432480)
	mapID := ""
	if player.CurrentMap != nil {
		mapID = player.CurrentMap.MapID
	}
	if order.Money > 0 && h.gameLogDispatcher != nil {
		h.gameLogDispatcher.Enqueue(&database.GameLogEntry{
			Type:    0x0519, // exe: Character.exe @ 0x00432394 (NC_LOG_GAME_MONEY_ADD)
			CharNo:  charNo,
			Map:     mapID,
			MapX:    player.Char.XPos,
			MapY:    player.Char.YPos,
			Int1:    int(order.Money),
			Bigint1: int64(player.Char.TCharacterState.Money) + int64(order.Money),
		})
	}
	if order.Fame > 0 && h.gameLogDispatcher != nil {
		h.gameLogDispatcher.Enqueue(&database.GameLogEntry{
			Type:   0x051C, // exe: Character.exe @ 0x00432418 (NC_LOG_GAME_FAME_ADD)
			CharNo: charNo,
			Map:    mapID,
			MapX:   player.Char.XPos,
			MapY:   player.Char.YPos,
			Int1:   int(order.Fame),
			Int2:   int(player.Char.TCharacterState.Fame) + int(order.Fame),
		})
	}

	// Paso 5: Confirmar vía 0x5816 (exe: 0x00432480..0x004324C3 -> Zone 0x0052E570)
	h.gds_NC_KQ_REWARDSUC_ACK(order.Handle, order.CharNo, order.OrderIndex)
}

func (h *Handler) rollbackKQRewardItems(player *ZonePlayer, keys []int64) {
	if h == nil || player == nil || player.Char == nil || len(keys) == 0 {
		return
	}
	keySet := make(map[int64]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
		if h.charRepo != nil {
			_ = h.charRepo.DeleteItem(k)
		}
	}
	filtered := player.Char.Items[:0]
	for _, it := range player.Char.Items {
		if it != nil && !keySet[it.ItemKey] {
			filtered = append(filtered, it)
		}
	}
	player.Char.Items = filtered
}

// gds_NC_KQ_REWARDSUC_ACK ports GameDBSession::gds_NC_KQ_REWARDSUC_ACK (Zone.exe @ 0x0052E570).
func (h *Handler) gds_NC_KQ_REWARDSUC_ACK(handle uint16, charNo uint32, orderIdx uint16) {
	if h == nil {
		return
	}
	rt := h.kqRewardTracker()
	key := kqRewardOrderKey(charNo, orderIdx)
	rt.mu.Lock()
	order, ok := rt.pending[key]
	if ok {
		delete(rt.pending, key)
		order.Committed = true
	}
	rt.mu.Unlock()
	if !ok || order == nil {
		return
	}

	// exe: 0x0052E5DF..0x0052E609: validar som_GetObject(handle) y CharNo
	p := h.findPlayerByHandle(handle)
	if p == nil || p.Char == nil || uint32(p.Char.ID) != charNo {
		return
	}

	// Aplicar money y fame en memoria del jugador y notificar al cliente (exe: 0x0052E626 commit order)
	if order.Money > 0 {
		p.Char.TCharacterState.Money += int64(order.Money)
		if p.Client != nil {
			pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CENCHANGE_CMD)
			pw.WriteUint64(uint64(p.Char.TCharacterState.Money))
			_ = p.Client.SendPacket(pw)
		}
	}
	if order.Fame > 0 {
		p.Char.TCharacterState.Fame += int(order.Fame)
		if p.Client != nil {
			pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_FAMECHANGE_CMD)
			pw.WriteUint32(uint32(p.Char.TCharacterState.Fame))
			_ = p.Client.SendPacket(pw)
		}
	}
}

// gds_NC_KQ_REWARDFAIL_ACK ports GameDBSession::gds_NC_KQ_REWARDFAIL_ACK (Zone.exe @ 0x0052E680).
func (h *Handler) gds_NC_KQ_REWARDFAIL_ACK(handle uint16, charNo uint32, orderIdx uint16, errCode uint16) {
	if h == nil {
		return
	}
	rt := h.kqRewardTracker()
	key := kqRewardOrderKey(charNo, orderIdx)
	rt.mu.Lock()
	order, ok := rt.pending[key]
	if ok {
		delete(rt.pending, key)
		order.Failed = true
		order.FailCode = errCode
	}
	rt.mu.Unlock()
}

// ============================================================================
// AxialListKQEnd & AxialListKQEndLua (Zone.exe 0x004284F0 / 0x004289D0 / 0x00429E00)
// + so_ply_KQSuccess (0x0055AE10) / so_ply_KQFail (0x0055AE30)
// + GameLog PROTO_NC_LOG_GAME_KQ_LEAVE (0x0442 @ 0x00429E52)
// ============================================================================

// so_ply_KQSuccess ports ShinePlayer::so_ply_KQSuccess (Zone.exe @ 0x0055AE10 -> CT_KQSuccess @ 0x005CBFB0).
func (h *Handler) so_ply_KQSuccess(player *ZonePlayer) {
	if h == nil || player == nil || !h.isPlayerOnKQMap(player) { // exe: 0x0055AE1E (KingdomMap == 1)
		return
	}
	h.UnlockOrAdvanceTitle(player, kqTitleTypeKQSuccess, 1, false) // exe: 0x0055AE2A -> 0x005CBFB0
}

// so_ply_KQFail ports ShinePlayer::so_ply_KQFail (Zone.exe @ 0x0055AE30 -> CT_KQFail @ 0x005CBFD0).
func (h *Handler) so_ply_KQFail(player *ZonePlayer) {
	if h == nil || player == nil || !h.isPlayerOnKQMap(player) { // exe: 0x0055AE3E (KingdomMap == 1)
		return
	}
	h.UnlockOrAdvanceTitle(player, kqTitleTypeKQFail, 1, false) // exe: 0x0055AE4A -> 0x005CBFD0
}

// AxialListKQEnd ports AxialListKQEnd (exe: 0x004284F0) & AxialListKQEndLua (exe: 0x004289D0).
// En Zone.exe ambos comparten ali_Work @ 0x00429E00 (MSVC ICF).
type AxialListKQEnd struct {
	Handler *Handler
	Dept    uint16 // +0x04 (0x16 = NC_KQ)
	Cmd     uint16 // +0x08 (0x12 = COMPLETE_CMD, 0x13 = FAIL_CMD)
}

// NewAxialListKQEndLua ports AxialListKQEndLua::AxialListKQEndLua (exe: 0x004289D0).
func (h *Handler) NewAxialListKQEndLua(dept, cmd uint16) *AxialListKQEnd {
	return &AxialListKQEnd{
		Handler: h,
		Dept:    dept,
		Cmd:     cmd,
	}
}

// ali_Work ports AxialListKQEndLua::ali_Work (Zone.exe @ 0x00429E00).
func (a *AxialListKQEnd) ali_Work(player *ZonePlayer) bool {
	if a == nil || a.Handler == nil || player == nil || player.Char == nil {
		return true // exe: 0x00429E0C
	}
	isSuccess := a.Cmd == kqEndCmdComplete // exe: 0x00429E26 (cmp dword ptr [ebx+8], 0x12)
	if isSuccess {
		a.Handler.so_ply_KQSuccess(player) // exe: 0x00429E30 -> 0x0055AE10
	} else {
		a.Handler.so_ply_KQFail(player) // exe: 0x00429E38 -> 0x0055AE30
	}

	// exe: 0x00429E40..0x00429EE2: emitir PROTO_NC_LOG_GAME_KQ_LEAVE (0x0442, 31 bytes payload)
	if player.CurrentMap != nil && a.Handler.gameLogDispatcher != nil {
		var leavePkt protocol.PROTO_NC_LOG_GAME_KQ_LEAVE_DATA
		leavePkt.Chrregnum = uint32(player.Char.ID)
		copy(leavePkt.Mappos.MapName[:], player.CurrentMap.MapID)
		leavePkt.Mappos.Coord.X = uint32(player.Char.XPos)
		leavePkt.Mappos.Coord.Y = uint32(player.Char.YPos)
		player.CurrentMap.mu.RLock()
		leavePkt.KQ_ID = player.CurrentMap.KQID
		leavePkt.KQ_Handle = player.CurrentMap.KQHandle
		player.CurrentMap.mu.RUnlock()
		if isSuccess {
			leavePkt.Leave_Status = 1 // exe: 0x00429EA7 (sete al)
		}
		a.Handler.gameLogDispatcher.EnqueueFromKQLeave(&leavePkt)
	}

	// exe: 0x00429EE7..0x00429EF9: enviar (dept=0x16, cmd=0x12|0x13) sin payload (0x5812 / 0x5813)
	if player.Client != nil {
		opcode := (a.Dept << 10) | (a.Cmd & 0x3FF)
		pw := protocol.NewPacketWriter(opcode)
		_ = player.Client.SendPacket(pw)
	}
	return true
}

// RunInMap aplica AxialListKQEnd(Lua) sobre todos los jugadores del mapa (so_AllInMap @ 0x005E322C).
func (a *AxialListKQEnd) RunInMap(zMap *ZoneMap) {
	if a == nil || zMap == nil {
		return
	}
	for _, p := range zMap.GetPlayers() {
		a.ali_Work(p)
	}
}

// ============================================================================
// PineEventScriptNode::ShineEndOfKingdomQuest::sa_Step (Zone.exe @ 0x004F5FC0)
// ============================================================================

// ShineEndOfKingdomQuest ports PineEventScriptNode::ShineEndOfKingdomQuest::sa_Step (exe: 0x004F5FC0).
// Idéntico en comportamiento a cEndOfKingdomQuest (exe: 0x005DF1F0):
// envía wms_EndOfKQPacket(kqHandle) (0x5810) y limpia objetos del mapa con máscara 0xB0.
func (h *Handler) ShineEndOfKingdomQuest(zMap *ZoneMap) bool {
	if h == nil || zMap == nil {
		return false
	}
	kqHandle := zMap.fm_GetKQhandle() // exe: 0x004F6078
	if kqHandle == kqInvalidHandle {  // exe: 0x004F607D (cmp edi, -1)
		return false
	}
	h.wms_EndOfKQPacket(kqHandle)               // exe: 0x004F609D -> 0x005A4800
	zMap.fm_ClearObject(kqClearObjectMaskEndLua) // exe: 0x004F60B7 -> fm_ClearObject(0xB0)
	return true
}

// ============================================================================
// GameLog PROTO_NC_LOG_GAME_KQ_ENTER (Zone.exe @ 0x00496740)
// ============================================================================

// EmitKQEnterGameLog emite PROTO_NC_LOG_GAME_KQ_ENTER (0x0441, 30 bytes payload)
// conforme a KingdomQuestGateKeeper::erb_happen (Zone.exe @ 0x00496740..0x004967DA).
func (h *Handler) EmitKQEnterGameLog(player *ZonePlayer, prevMap string, prevX, prevY int, kqID uint16, kqHandle uint32) {
	if h == nil || h.gameLogDispatcher == nil || player == nil || player.Char == nil || prevMap == "" {
		return
	}
	var enterPkt protocol.PROTO_NC_LOG_GAME_KQ_ENTER_DATA
	enterPkt.Chrregnum = uint32(player.Char.ID)
	copy(enterPkt.Mappos.MapName[:], prevMap)
	enterPkt.Mappos.Coord.X = uint32(prevX)
	enterPkt.Mappos.Coord.Y = uint32(prevY)
	enterPkt.KQ_ID = kqID
	enterPkt.KQ_Handle = kqHandle
	h.gameLogDispatcher.EnqueueFromKQEnter(&enterPkt)
}

// ============================================================================
// so_SaveLocation (Zone.exe @ 0x004542F0) — Guardado de Posición Normal + KQ
// ============================================================================

// SavePreKQReturnPos guarda el mapa y coordenadas previas al entrar a una KQ para rollback seguro.
func (h *Handler) SavePreKQReturnPos(player *ZonePlayer, prevMap string, prevX, prevY int) {
	if player == nil || prevMap == "" {
		return
	}
	// Solo sobreescribir si el mapa previo NO es ya una KQ
	if h != nil && h.dataStore != nil {
		if mi, ok := h.dataStore.GetKQMapInfo(prevMap); ok && mi != nil && mi.KingdomMap == 1 {
			return
		}
	}
	player.KQReturnMap = prevMap
	player.KQReturnX = prevX
	player.KQReturnY = prevY
}

// so_SaveLocation ports ShinePlayer::so_SaveLocation (Zone.exe @ 0x004542F0..0x00454773)
// y persiste los 8 parámetros oficiales vía p_Char_SaveLocation (Character.exe @ 0x00459E90).
func (h *Handler) so_SaveLocation(player *ZonePlayer) (loginMap string, loginX, loginY int, kqHandle uint32, kqMap string, kqX, kqY uint32, ok bool) {
	if h == nil || player == nil || player.Char == nil || player.CurrentMap == nil {
		return "", 0, 0, kqInvalidHandle, "", 0, 0, false // exe: 0x0045435A
	}

	curMap := player.CurrentMap
	curMapName := curMap.MapID
	curX := player.Char.XPos
	curY := player.Char.YPos

	// 1. Determinar posición de respawn segura (sLoginZone, nLoginZoneX, nLoginZoneY)
	// exe: 0x004544DA..0x0045461F: si está muerto o en mapa sin IsSavePos (KQ), usar mapa de retorno/RegenCity
	onKQ := h.isPlayerOnKQMap(player)
	if onKQ || player.State == 3 {
		if player.KQReturnMap != "" && player.KQReturnX > kqSaveCoordMinValid && player.KQReturnY > kqSaveCoordMinValid {
			loginMap = player.KQReturnMap
			loginX = player.KQReturnX
			loginY = player.KQReturnY
		} else {
			loginMap, loginX, loginY = h.resolveRegenCityLocation(curMap) // exe: 0x004545C0..0x0045461C
		}
	} else {
		loginMap = curMapName
		loginX = curX
		loginY = curY
	}

	// 2. Campos KQ actuales (exe: 0x0045461F..0x00454649)
	kqHandle = curMap.fm_GetKQhandle() // [edi+0x18] = CurrentMap+0x3C
	kqMap = curMapName                 // [edi+0x1C..0x27] = CurrentMap name
	kqX = uint32(curX)                 // [edi+0x28] = pos.X
	kqY = uint32(curY)                 // [edi+0x2C] = pos.Y

	// 3. Validación de coordenadas mínimas <= 0x0A -> fallback RouN (0x192D, 0x21B6) (exe: 0x004546CE..0x00454750)
	if loginX <= kqSaveCoordMinValid || loginY <= kqSaveCoordMinValid || kqX <= kqSaveCoordMinValid || kqY <= kqSaveCoordMinValid {
		loginMap = kqRoumenFallbackMap
		loginX = kqRoumenFallbackX
		loginY = kqRoumenFallbackY
	}

	// 4. Persistir en CharacterRepository vía p_Char_SaveLocation (Character.exe @ 0x00459E90)
	if h.charRepo != nil {
		_ = h.charRepo.SaveCharLocationWithKQ(player.Char.ID, loginMap, loginX, loginY, kqHandle, kqMap, kqX, kqY)
	}
	player.Char.KQHandle = kqHandle
	player.Char.KQMapName = kqMap
	player.Char.KQX = kqX
	player.Char.KQY = kqY
	return loginMap, loginX, loginY, kqHandle, kqMap, kqX, kqY, true
}

func (h *Handler) resolveRegenCityLocation(zMap *ZoneMap) (string, int, int) {
	destMap := kqRoumenFallbackMap
	destX := kqRoumenFallbackX
	destY := kqRoumenFallbackY
	if h != nil && h.fieldStore != nil && zMap != nil {
		fa, ok := h.fieldStore.GetField(zMap.MapID)
		if !ok && zMap.FieldEntry != nil {
			fa = zMap.FieldEntry
			ok = true
		}
		if ok && fa != nil {
			revCity := strings.TrimSpace(fa.Respawn.City)
			targetFA := fa
			if revCity != "" && revCity != "-" {
				if cityFA, okCity := h.fieldStore.GetField(revCity); okCity && cityFA != nil {
					targetFA = cityFA
				}
			}
			if targetFA.Respawn.City != "" && targetFA.Respawn.City != "-" {
				destMap = targetFA.Respawn.City
			}
			if len(targetFA.Respawn.Spots) > 0 && (targetFA.Respawn.Spots[0].X != 0 || targetFA.Respawn.Spots[0].Y != 0) {
				destX = targetFA.Respawn.Spots[0].X
				destY = targetFA.Respawn.Spots[0].Y
			}
		}
	}
	return destMap, destX, destY
}

// ============================================================================
// Ejecución Real de Scripts Lua de KQ (23 directorios de data/LuaScript/KQ/)
// ============================================================================

// kqTableScriptEntry mapea scripts PineEventScript-migrados que exponen <Tabla>.Main(Field[, InitFlag]).
var kqTableScriptEntry = map[string]string{
	"conditionofhero": "KQConditionOfHero",
	"gordonmaster":    "KQGordonMaster",
	"honeying":        "KQHoneying",
	"kqhbat":          "KQHBat",
	"underhall":       "KQUnderHall",
}

// ResolveKQScriptPath resuelve ScriptLanguage de KingdomQuest.shn (o nombre de directorio KQ)
// al archivo .lua ejecutable dentro de data/LuaScript/.
func ResolveKQScriptPath(dataDir, scriptLanguage string) (scriptPath string, tableEntry string, ok bool) {
	lang := strings.TrimSpace(strings.ReplaceAll(scriptLanguage, "\\", "/"))
	if lang == "" {
		return "", "", false
	}
	if dataDir == "" {
		dataDir = resolveZoneDataDir()
	}
	luaRoot := filepath.Join(dataDir, "LuaScript")

	parts := strings.Split(lang, "/")
	base := parts[len(parts)-1]
	baseNoDigits := strings.TrimRight(base, "0123456789")
	parent := ""
	if len(parts) > 1 {
		parent = strings.Join(parts[:len(parts)-1], "/")
	}

	candidates := []string{
		filepath.Join(luaRoot, filepath.FromSlash(lang)+".lua"),
		filepath.Join(luaRoot, filepath.FromSlash(lang), base+".lua"),
	}
	if baseNoDigits != "" && baseNoDigits != base {
		candidates = append(candidates,
			filepath.Join(luaRoot, filepath.FromSlash(parent), baseNoDigits, baseNoDigits+".lua"),
			filepath.Join(luaRoot, "KQ", baseNoDigits, baseNoDigits+".lua"),
		)
	}
	if !strings.HasPrefix(strings.ToUpper(lang), "KQ/") {
		candidates = append(candidates,
			filepath.Join(luaRoot, "KQ", filepath.FromSlash(lang)+".lua"),
			filepath.Join(luaRoot, "KQ", filepath.FromSlash(lang), base+".lua"),
		)
	}

	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			dirKey := strings.ToLower(strings.TrimSuffix(filepath.Base(c), filepath.Ext(c)))
			return c, kqTableScriptEntry[dirKey], true
		}
	}
	return "", "", false
}

// StartKQScript carga y ejecuta el script Lua oficial de la KQ en un motor aislado sobre ZoneMap.
// Invocado por cc_PlayFilm (Zone.exe @ 0x005089A0).
func (m *ZoneMap) StartKQScript(scriptName, initValue string) error {
	if m == nil {
		return nil
	}
	dataDir := ""
	if m.dataStore != nil && m.dataStore.DataDir != "" {
		dataDir = m.dataStore.DataDir
	} else {
		dataDir = resolveZoneDataDir()
	}
	scriptPath, tableEntry, found := ResolveKQScriptPath(dataDir, scriptName)
	if !found {
		return nil
	}

	luaRoot := filepath.Join(dataDir, "LuaScript")
	bridge := NewZoneWorldBridge(m)
	eng, err := lua.NewEngine(lua.Config{
		BasePath: luaRoot,
		Bridge:   bridge,
	})
	if err != nil {
		return err
	}

	// Exponer InitValue como global y como argumento numérico/string si corresponde
	eng.SetGlobal("InitValue", glua.LString(initValue))
	if err := eng.DoFile(scriptPath); err != nil {
		eng.Close()
		return err
	}

	m.mu.Lock()
	if m.luaEngine != nil {
		m.luaEngine.Close()
	}
	m.luaEngine = eng
	m.worldBridge = bridge
	m.kqScriptTable = tableEntry
	m.mu.Unlock()

	// Invocar punto de entrada inicial: Main(Field) o <KQTable>.Main(Field, InitValue)
	if tableEntry != "" {
		if initInt, parseErr := strconv.Atoi(strings.TrimSpace(initValue)); parseErr == nil {
			_, err = eng.CallTableFunction(tableEntry, "Main", glua.LString(m.MapID), glua.LNumber(initInt))
		} else {
			_, err = eng.CallTableFunction(tableEntry, "Main", glua.LString(m.MapID), glua.LString(initValue))
		}
		if err != nil {
			log.Printf("[Zone][KQ-Lua] Error ejecutando %s.Main(%s, %q): %v", tableEntry, m.MapID, initValue, err)
		}
		return err
	}

	if eng.GetGlobal("Main").Type() == glua.LTFunction {
		if _, callErr := eng.CallFunction("Main", glua.LString(m.MapID)); callErr != nil {
			log.Printf("[Zone][KQ-Lua] Error ejecutando Main(%s) de %s: %v", m.MapID, scriptName, callErr)
			return callErr
		}
	}
	return nil
}

// TickKQScript ejecuta un paso del bucle Main(Field) del script Lua activo de KQ.
func (m *ZoneMap) TickKQScript() error {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	active := m.KQScriptActive
	eng := m.luaEngine
	tableEntry := m.kqScriptTable
	mapID := m.MapID
	m.mu.RUnlock()

	if !active || eng == nil || tableEntry != "" {
		return nil
	}
	if eng.GetGlobal("Main").Type() != glua.LTFunction {
		return nil
	}
	_, err := eng.CallFunction("Main", glua.LString(mapID))
	return err
}

// StopKQScript detiene y libera el motor Lua de la instancia KQ (invocado por cc_DropFilm @ 0x00509150).
func (m *ZoneMap) StopKQScript() {
	if m == nil {
		return
	}
	m.mu.Lock()
	eng := m.luaEngine
	m.luaEngine = nil
	m.kqScriptTable = ""
	m.mu.Unlock()
	if eng != nil {
		eng.Close()
	}
}
