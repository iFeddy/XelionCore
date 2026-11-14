package zone

import (
	"log"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/data/drop"
	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

const (
	// DefaultDropOwnershipDuration define la ventana estándar de propiedad temporal tras morir un monstruo.
	DefaultDropOwnershipDuration = 30 * time.Second
	// DefaultDropDespawnAge define el tiempo máximo que un ítem permanece en el suelo antes de limpiarse.
	DefaultDropDespawnAge = 120 * time.Second
)

// RandomOptRule representa una regla de opción aleatoria de RandomOption.shn.
type RandomOptRule struct {
	OptionType uint8
	Min        uint32
	Max        uint32
	DropRate   uint32
}

// RandomOptCountRule representa una regla de cantidad máxima de opciones de RandomOptionCount.shn.
type RandomOptCountRule struct {
	LimitCount uint16
	DropRate   uint16
}

// RandomOptionCatalog indexa RandomOption.shn y RandomOptionCount.shn.
type RandomOptionCatalog struct {
	RulesByItem  map[string][]RandomOptRule
	CountsByItem map[string][]RandomOptCountRule
}

// PartyBonusTables indexa PartyBonusByMember.shn, PartyBonusByLvDiff.shn y PartyBonusLimit.shn.
type PartyBonusTables struct {
	ByMember map[int]int   // count -> ratio (1000 = 100%)
	ByLvDiff map[int]int   // abs(maxLv - memberLv) -> ratio (1000 = 100%)
	MaxLvLim map[uint8]uint64 // level -> max bonus exp cap
	MaxMemK  int
	MaxDiffK int
}

var (
	randOptCache sync.Map // *DataStore -> *RandomOptionCatalog
	partyCache   sync.Map // *DataStore -> *PartyBonusTables
)

func getRandomOptCatalog(ds *DataStore) *RandomOptionCatalog {
	if ds == nil {
		return nil
	}
	if v, ok := randOptCache.Load(ds); ok {
		return v.(*RandomOptionCatalog)
	}
	cat := &RandomOptionCatalog{
		RulesByItem:  make(map[string][]RandomOptRule),
		CountsByItem: make(map[string][]RandomOptCountRule),
	}
	if f := ds.GetSHN("RandomOption"); f != nil {
		for _, row := range f.Rows {
			inx := strings.ToLower(strings.TrimSpace(row.GetString("DropItemIndex")))
			if inx == "" {
				continue
			}
			cat.RulesByItem[inx] = append(cat.RulesByItem[inx], RandomOptRule{
				OptionType: uint8(row.GetInt("RandomOptionType")),
				Min:        uint32(row.GetInt("Min")),
				Max:        uint32(row.GetInt("Max")),
				DropRate:   uint32(row.GetInt("TypeDropRate")),
			})
		}
	}
	if f := ds.GetSHN("RandomOptionCount"); f != nil {
		for _, row := range f.Rows {
			inx := strings.ToLower(strings.TrimSpace(row.GetString("InxName")))
			if inx == "" {
				continue
			}
			cat.CountsByItem[inx] = append(cat.CountsByItem[inx], RandomOptCountRule{
				LimitCount: row.GetUint16("LimitCount"),
				DropRate:   row.GetUint16("LimitDropRate"),
			})
		}
	}
	actual, _ := randOptCache.LoadOrStore(ds, cat)
	return actual.(*RandomOptionCatalog)
}

func getPartyBonusTables(ds *DataStore) *PartyBonusTables {
	if ds == nil {
		return nil
	}
	if v, ok := partyCache.Load(ds); ok {
		return v.(*PartyBonusTables)
	}
	t := &PartyBonusTables{
		ByMember: make(map[int]int),
		ByLvDiff: make(map[int]int),
		MaxLvLim: make(map[uint8]uint64),
	}
	if f := ds.GetSHN("PartyBonusByMember"); f != nil {
		for _, row := range f.Rows {
			cnt := int(row.GetByte("PB_PartyMember"))
			if cnt == 0 {
				cnt = row.GetInt("PB_PartyMember")
			}
			ratio := int(row.GetUint16("PB_BonusRatio"))
			if ratio == 0 {
				ratio = row.GetInt("PB_BonusRatio")
			}
			t.ByMember[cnt] = ratio
			if cnt > t.MaxMemK {
				t.MaxMemK = cnt
			}
		}
	}
	if f := ds.GetSHN("PartyBonusByLvDiff"); f != nil {
		for _, row := range f.Rows {
			diff := int(int16(row.GetUint16("PB_LvDiff")))
			if diff < 0 {
				diff = -diff
			}
			ratio := int(row.GetUint16("PB_BonusRatio"))
			t.ByLvDiff[diff] = ratio
			if diff > t.MaxDiffK {
				t.MaxDiffK = diff
			}
		}
	}
	if f := ds.GetSHN("PartyBonusLimit"); f != nil {
		for _, row := range f.Rows {
			lv := row.GetByte("PSE_ChrLv")
			lim := uint64(row.GetInt("PSE_ExpLimit"))
			if lv > 0 && lim > 0 {
				t.MaxLvLim[lv] = lim
			}
		}
	}
	actual, _ := partyCache.LoadOrStore(ds, t)
	return actual.(*PartyBonusTables)
}

// SpawnMonsterDrop coloca un ítem o Cen en el suelo con propiedad temporal y opciones aleatorias y difunde 0x1C0A.
func (m *ZoneMap) SpawnMonsterDrop(itemID uint16, lot uint32, flags int, x, y int32, dropMobHandle uint16, ownerCharID int, ownerDur time.Duration, allCanLoot bool, isCen bool, statOpts []item.ItemStatOption) *ZoneDroppedItem {
	if m == nil {
		return nil
	}
	var expiresAt time.Time
	if !allCanLoot && ownerCharID > 0 && ownerDur > 0 {
		expiresAt = time.Now().Add(ownerDur)
	}

	m.mu.Lock()
	h := m.nextObjID
	m.nextObjID++
	di := &ZoneDroppedItem{
		Handle:         h,
		ItemID:         itemID,
		Lot:            lot,
		Flags:          flags,
		X:              x,
		Y:              y,
		DropMobHandle:  dropMobHandle,
		Attr:           0x01,
		OwnerCharID:    ownerCharID,
		OwnerExpiresAt: expiresAt,
		AllCanLoot:     allCanLoot,
		IsCenDrop:      isCen,
		StatOptions:    statOpts,
		CreatedAt:      time.Now(),
	}
	if m.droppedItems == nil {
		m.droppedItems = make(map[uint16]*ZoneDroppedItem)
	}
	m.droppedItems[h] = di
	m.mu.Unlock()

	if m.grid != nil {
		m.grid.AddEntity(di)
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_DROPEDITEM_CMD)
	pw.WriteUint16(di.Handle)
	pw.WriteUint16(di.ItemID)
	pw.WriteInt32(di.X)
	pw.WriteInt32(di.Y)
	pw.WriteUint16(di.DropMobHandle)
	if di.OwnerCharID > 0 {
		pw.WriteByte(0)
	} else {
		pw.WriteByte(1)
	}
	if m.grid != nil {
		m.grid.BroadcastAOI(pw, int(di.X), int(di.Y), 0)
	}
	return di
}

// HandleMonsterDeath ejecuta el pipeline oficial completo cuando muere un monstruo:
// 1. EXP + Party/Raid sharing + RestExp + LevelUp + Fame
// 2. Drops probabilísticos (Cen, ItemDropTable, ItemDropGroup, CollectCard, RandomOption)
// 3. Propiedad temporal, AutoPick, MobKillAnnounce y GameLog Tipo 6
// 4. Progreso automático de títulos (CharacterTitleData.shn)
func (h *Handler) HandleMonsterDeath(zMap *ZoneMap, mob *ZoneMonster, killer *ZonePlayer) {
	if h == nil || zMap == nil || mob == nil {
		return
	}

	// Resolver killer si vino nil (ej. DoT/AoE) a partir del mayor daño registrado
	if killer == nil {
		topCharID := mob.TopDamageChar()
		if topCharID > 0 {
			killer = h.findPlayerInMapByCharID(zMap, topCharID)
		}
	}

	mob.mu.RLock()
	mobID := mob.MobID
	mobInx := mob.InxName
	mobHandle := mob.Handle
	mobX, mobY := mob.X, mob.Y
	monEXP := mob.Exp
	expRange := mob.EXPRange
	allCanLoot := mob.AllCanLoot
	mobRank := mob.Rank
	mob.mu.RUnlock()

	if expRange <= 0 {
		expRange = 1500
	}

	// 1. Distribución de EXP (Solo / Party / Raid) y Fama
	if killer != nil && killer.Char != nil {
		h.distributeMonsterExpAndFame(zMap, mobID, mobX, mobY, monEXP, expRange, mobRank, killer)
	}

	// 2. Generación de Drops (Cen, Items, Cartas) + AutoPick + GameLog Tipo 6
	ownerCharID := 0
	if killer != nil && killer.Char != nil {
		ownerCharID = killer.Char.ID
	} else {
		ownerCharID = mob.TopDamageChar()
	}
	h.generateMonsterDrops(zMap, mobID, mobInx, mobHandle, mobX, mobY, allCanLoot, ownerCharID, killer)

	// 3. Anuncio oficial MobKillAnnounce.shn (0x2459 / 0x2458)
	if h.dataStore != nil && h.dataStore.MobCatalog != nil {
		if textIdx, hasAnn := h.dataStore.MobCatalog.KillAnnounce[mobID]; hasAnn {
			h.broadcastMobKillAnnounce(zMap, textIdx)
		}
	}

	// 4. Progreso automático de títulos (Dept 24)
	if killer != nil && killer.Char != nil {
		h.OnMonsterKilledTitles(killer, mob)
	}
}

func (h *Handler) findPlayerInMapByCharID(zMap *ZoneMap, charID int) *ZonePlayer {
	if zMap == nil || charID <= 0 {
		return nil
	}
	zMap.mu.RLock()
	defer zMap.mu.RUnlock()
	for _, p := range zMap.players {
		if p != nil && p.Char != nil && p.Char.ID == charID {
			return p
		}
	}
	return nil
}

// distributeMonsterExpAndFame distribuye experiencia individual, en Party o en Raid aplicando RestExp y Fame.
func (h *Handler) distributeMonsterExpAndFame(zMap *ZoneMap, mobID uint16, mobX, mobY int32, monEXP uint32, expRange int, mobRank uint8, killer *ZonePlayer) {
	if killer == nil || killer.Char == nil {
		return
	}

	recipients := []*ZonePlayer{killer}
	if h.charRepo != nil {
		if partyNo, _, _, err := h.charRepo.GetPartyByChar(killer.Char.ID); err == nil && partyNo > 0 {
			if members, err := h.charRepo.GetPartyMembers(partyNo); err == nil && len(members) > 0 {
				var inRange []*ZonePlayer
				seen := make(map[int]bool)
				for _, m := range members {
					if m == nil || seen[m.CharNo] {
						continue
					}
					p := h.findPlayerInMapByCharID(zMap, m.CharNo)
					if p == nil || p.Char == nil || p.State == 3 {
						continue
					}
					px, py := p.GetPosition()
					if math.Hypot(float64(px-int(mobX)), float64(py-int(mobY))) <= float64(expRange) {
						inRange = append(inRange, p)
						seen[m.CharNo] = true
					}
				}
				if len(inRange) > 0 {
					recipients = inRange
				}
			}
		}
	}

	if monEXP == 0 {
		return
	}

	var maxLv uint8
	for _, p := range recipients {
		if p.Char.Level > maxLv {
			maxLv = p.Char.Level
		}
	}

	nMembers := len(recipients)
	partyTables := getPartyBonusTables(h.dataStore)
	memberRatio := 1000
	if nMembers > 1 {
		if partyTables != nil {
			if r, ok := partyTables.ByMember[nMembers]; ok {
				memberRatio = r
			} else if nMembers > partyTables.MaxMemK && partyTables.MaxMemK > 0 {
				// Soporta grupos tipo Raid (>5 miembros) escalando desde la última entrada SHN
				memberRatio = partyTables.ByMember[partyTables.MaxMemK] + (nMembers-partyTables.MaxMemK)*100
			}
		} else if h.dataStore != nil && h.dataStore.CombatStore != nil {
			memberRatio = h.dataStore.CombatStore.GetPartyExpBonus(nMembers)
		}
	}

	totalPartyExp := (uint64(monEXP) * uint64(memberRatio)) / 1000
	if totalPartyExp == 0 {
		totalPartyExp = uint64(monEXP)
	}

	for _, p := range recipients {
		share := totalPartyExp / uint64(nMembers)
		if nMembers > 1 && partyTables != nil {
			diff := int(maxLv) - int(p.Char.Level)
			if diff < 0 {
				diff = -diff
			}
			diffRatio := 1000
			if r, ok := partyTables.ByLvDiff[diff]; ok {
				diffRatio = r
			} else if diff > partyTables.MaxDiffK && partyTables.MaxDiffK > 0 {
				diffRatio = partyTables.ByLvDiff[partyTables.MaxDiffK]
			}
			share = (share * uint64(diffRatio)) / 1000
			if lim, ok := partyTables.MaxLvLim[p.Char.Level]; ok && lim > 0 && share > lim {
				share = lim
			}
		}
		if share == 0 {
			share = 1
		}

		// Aplicar RestExp desde Common.json + tRestExp
		share = h.applyRestExpBonus(p, share)

		h.AddPlayerExpAndCheckLevelUp(p, share, mobID, mobX, mobY, zMap.MapID)

		// Otorgar fama al matar jefes/héroes (Rank >= 4)
		if mobRank >= 4 && h.charRepo != nil {
			fameGain := int(mobRank) * 2
			beforeFame := p.Char.TCharacterState.Fame
			if err := h.charRepo.AddCharFame(p.Char.ID, fameGain); err == nil {
				if p.Char.TCharacterState.Fame == beforeFame {
					p.Char.TCharacterState.Fame += fameGain
				}
				if p.Client != nil {
					famePkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_FAMECHANGE_CMD)
					famePkt.WriteUint32(uint32(p.Char.TCharacterState.Fame))
					_ = p.Client.SendPacket(famePkt)
				}
			}
		}
	}
}

func (h *Handler) applyRestExpBonus(p *ZonePlayer, baseExp uint64) uint64 {
	if h.dataStore == nil || h.dataStore.CombatStore == nil || p == nil || p.Char == nil {
		return baseExp
	}
	common := h.dataStore.CombatStore.Common()
	if common.RestExpRate <= 1000 {
		return baseExp
	}

	hasRestBonus := p.RestExpBonusRate > 0
	if !hasRestBonus && h.charRepo != nil {
		if t, err := h.charRepo.GetRestExpExecTime(p.Char.ID); err == nil && !t.IsZero() {
			minDur := time.Duration(common.MinNeedTime) * time.Second
			if minDur <= 0 {
				minDur = 6 * time.Hour
			}
			if time.Since(t) >= minDur {
				hasRestBonus = true
			}
		}
	}
	if !hasRestBonus {
		return baseExp
	}
	boosted := (baseExp * uint64(common.RestExpRate)) / 1000
	if boosted <= baseExp {
		return baseExp
	}
	return boosted
}

// AddPlayerExpAndCheckLevelUp añade EXP al jugador, gestiona múltiples subidas de nivel y emite paquetes + GameLog.
func (h *Handler) AddPlayerExpAndCheckLevelUp(p *ZonePlayer, expGain uint64, mobID uint16, mobX, mobY int32, mapID string) {
	if p == nil || p.Char == nil || expGain == 0 {
		return
	}

	p.Char.TCharacterState.Exp += int64(expGain)

	var cs *combat.Store
	if h.dataStore != nil {
		cs = h.dataStore.CombatStore
	}
	maxLv := 125
	if cs != nil {
		if m := cs.MaxLevel(); m > 0 {
			maxLv = m
		}
		if c := cs.Common(); c.LevelLimit > 0 && c.LevelLimit < maxLv {
			maxLv = c.LevelLimit
		}
	}

	leveledUp := false
	for int(p.Char.Level) < maxLv {
		if cs == nil {
			break
		}
		nextExp, ok := cs.GetNextExp(int(p.Char.Level))
		if !ok || nextExp <= 0 || p.Char.TCharacterState.Exp < int64(nextExp) {
			break
		}

		p.Char.Level++
		p.Char.TCharacter.Level = p.Char.Level
		leveledUp = true

		// Actualizar ClassParam + recarga total de HP/SP al subir nivel
		if h.classStore != nil {
			if lp, ok := h.classStore.GetParam(JobToClassName(p.Char.Job), int(p.Char.Level)); ok && lp != nil {
				p.LevelParam = *lp
				free := cs.GetFreeStatBonus(
					p.Char.TCharacter.Strength,
					p.Char.TCharacter.Intelligence,
					p.Char.TCharacter.Dexterity,
					p.Char.TCharacter.Constitute,
					p.Char.TCharacter.MentalPower,
				)
				if lp.MaxHP+free.MaxHP > 0 {
					p.MaxHP = lp.MaxHP + free.MaxHP
				}
				if lp.MaxSP+free.MaxSP > 0 {
					p.MaxSP = lp.MaxSP + free.MaxSP
				}
			}
		}
		p.Char.CurHP = p.MaxHP
		p.Char.CurSP = p.MaxSP

		prevExp, newNext := cs.GetExpProgression(int(p.Char.Level))
		p.PrevExp = uint64(prevExp)
		p.NextExp = uint64(newNext)

		if p.Client != nil {
			// PROTO_NC_CHAR_LEVEL_CHANGED_CMD (0x1074, 7B): wmhandle(2), CharNo(4), CharNewLevel(1)
			lvCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LEVEL_CHANGED_CMD)
			lvCmd.WriteUint16(p.MapObjectID)
			lvCmd.WriteUint32(uint32(p.Char.ID))
			lvCmd.WriteByte(p.Char.Level)
			_ = p.Client.SendPacket(lvCmd)

			hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
			hpCmd.WriteUint32(uint32(p.Char.CurHP))
			hpCmd.WriteUint16(0)
			_ = p.Client.SendPacket(hpCmd)

			spCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
			spCmd.WriteUint32(uint32(p.Char.CurSP))
			_ = p.Client.SendPacket(spCmd)
		}

		if h.gameLogDispatcher != nil {
			h.gameLogDispatcher.EnqueueType2(1, p.Char.ID, mapID, int(mobX), int(mobY), 0, int(p.Char.Level), 0)
		}

		h.OnPlayerLevelUpTitles(p)
	}

	if h.charRepo != nil {
		if err := h.charRepo.SaveCharLevelExp(p.Char.ID, int16(p.Char.Level), p.Char.TCharacterState.Exp); err != nil {
			log.Printf("[Zone][Drops] Error SaveCharLevelExp char %d: %v", p.Char.ID, err)
		}
	}

	if p.Client != nil {
		// PROTO_NC_CHAR_EXP_CHANGED_CMD (0x1073, 14B): wmhandle(2), CharNo(4), CurrentExp(8)
		expCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_EXP_CHANGED_CMD)
		expCmd.WriteUint16(p.MapObjectID)
		expCmd.WriteUint32(uint32(p.Char.ID))
		expCmd.WriteUint64(uint64(p.Char.TCharacterState.Exp))
		_ = p.Client.SendPacket(expCmd)
	}

	if h.gameLogDispatcher != nil {
		// Tipo 6 nType=2: ganancia de experiencia por muerte de mob (PROTO_NC_LOG_GAME_MKED)
		h.gameLogDispatcher.EnqueueType6(2, p.Char.ID, mapID, int(mobX), int(mobY), 0, int(mobID), int64(expGain))
	}

	if leveledUp {
		p.RecalculateEquipmentStats(h.itemStore)
	}
}

func (h *Handler) generateMonsterDrops(zMap *ZoneMap, mobID uint16, mobInx string, mobHandle uint16, mobX, mobY int32, allCanLoot bool, ownerCharID int, killer *ZonePlayer) {
	if h.dataStore == nil || h.itemStore == nil {
		return
	}

	spawnCount := 0
	offsetPos := func() (int32, int32) {
		ang := float64(spawnCount) * 1.15
		dist := 18.0 + float64(spawnCount%4)*6.0
		spawnCount++
		return mobX + int32(math.Cos(ang)*dist), mobY + int32(math.Sin(ang)*dist)
	}

	// A. Drops de tabla oficial (ItemDropTable.json + ItemDropGroup.json)
	if h.dataStore.DropStore != nil {
		if tbl, ok := h.dataStore.DropStore.GetMobDrop(mobInx); ok && tbl != nil {
			excluded := make(map[string]bool, len(tbl.ExcludedItems))
			for _, ex := range tbl.ExcludedItems {
				excluded[strings.ToLower(strings.TrimSpace(ex))] = true
			}

			// 1. Money / Cen
			if tbl.Money.Rate > 0 && tbl.Money.Max >= tbl.Money.Min && tbl.Money.Max > 0 {
				if rand.Intn(1_000_000) < tbl.Money.Rate {
					cenAmt := tbl.Money.Min
					if tbl.Money.Max > tbl.Money.Min {
						cenAmt += rand.Intn(tbl.Money.Max - tbl.Money.Min + 1)
					}
					if cenAmt > 0 {
						dx, dy := offsetPos()
						moneyItemID := h.resolveMoneyItemID()
						di := zMap.SpawnMonsterDrop(moneyItemID, uint32(cenAmt), cenAmt, dx, dy, mobHandle, ownerCharID, DefaultDropOwnershipDuration, allCanLoot, true, nil)
						if h.gameLogDispatcher != nil {
							h.gameLogDispatcher.EnqueueType6(1, ownerCharID, zMap.MapID, int(dx), int(dy), 0, int(moneyItemID), int64(cenAmt))
						}
						if di != nil && killer != nil && killer.IsAutoPick && (allCanLoot || di.OwnerCharID == 0 || di.OwnerCharID == killer.Char.ID) {
							h.autoPickDroppedItem(zMap, killer, di)
						}
					}
				}
			}

			// 2. TradeBoxes
			for _, tb := range tbl.TradeBoxes {
				if tb.Rate <= 0 {
					continue
				}
				if rand.Intn(1_000_000) < tb.Rate {
					h.rollAndSpawnDropEntry(zMap, tb.Box, 1, 0, 0, excluded, offsetPos, mobHandle, ownerCharID, allCanLoot, killer)
				}
			}

			// 3. Drops regulares o grupos
			for _, de := range tbl.Drops {
				if de.Rate <= 0 {
					continue
				}
				if rand.Intn(1_000_000) < de.Rate {
					qty := de.Num
					if qty <= 0 {
						qty = 1
					}
					h.rollAndSpawnDropEntry(zMap, de.Item, qty, de.MinUpgrade, de.MaxUpgrade, excluded, offsetPos, mobHandle, ownerCharID, allCanLoot, killer)
				}
			}
		}
	}

	// B. Drops de Cartas de Colección (CollectCardMobGroup.shn + CollectCardDropRate.shn)
	if h.dataStore.CollectCatalog != nil && mobInx != "" {
		cards := h.dataStore.CollectCatalog.CardsForMob(mobInx)
		for _, c := range cards {
			if c.DropRate <= 0 {
				continue
			}
			// CC_CardGetRate en Fiesta está en escala por 1,000,000
			if rand.Intn(1_000_000) < int(c.DropRate) {
				if info, ok := h.itemStore.GetItemInfoByName(c.ItemInx); ok && info != nil {
					dx, dy := offsetPos()
					di := zMap.SpawnMonsterDrop(info.ID, 1, 1, dx, dy, mobHandle, ownerCharID, DefaultDropOwnershipDuration, allCanLoot, false, nil)
					if h.gameLogDispatcher != nil {
						h.gameLogDispatcher.EnqueueType6(1, ownerCharID, zMap.MapID, int(dx), int(dy), 0, int(info.ID), 1)
					}
					if di != nil && killer != nil && killer.IsAutoPick {
						h.autoPickDroppedItem(zMap, killer, di)
					}
				}
			}
		}
	}
}

func (h *Handler) rollAndSpawnDropEntry(
	zMap *ZoneMap,
	targetInx string,
	qty int,
	minUp int,
	maxUp int,
	excluded map[string]bool,
	offsetPos func() (int32, int32),
	mobHandle uint16,
	ownerCharID int,
	allCanLoot bool,
	killer *ZonePlayer,
) {
	itemInx := targetInx
	if entries, isGrp := h.dataStore.DropStore.GetGroup(targetInx); isGrp && len(entries) > 0 {
		chosen := entries[rand.Intn(len(entries))]
		itemInx = chosen.ItemID
		if chosen.MaxQuantity >= chosen.MinQuantity && chosen.MaxQuantity > 0 {
			qty = chosen.MinQuantity
			if chosen.MaxQuantity > chosen.MinQuantity {
				qty += rand.Intn(chosen.MaxQuantity - chosen.MinQuantity + 1)
			}
		}
	}
	if itemInx == "" || excluded[strings.ToLower(strings.TrimSpace(itemInx))] {
		return
	}
	info, ok := h.itemStore.GetItemInfoByName(itemInx)
	if !ok || info == nil {
		return
	}
	if qty <= 0 {
		qty = 1
	}

	flags := qty
	if minUp > 0 || maxUp > 0 {
		_ = minUp
	}

	opts := h.RollItemRandomOptions(info.InxName)
	dx, dy := offsetPos()
	di := zMap.SpawnMonsterDrop(info.ID, uint32(qty), flags, dx, dy, mobHandle, ownerCharID, DefaultDropOwnershipDuration, allCanLoot, false, opts)

	if h.gameLogDispatcher != nil {
		h.gameLogDispatcher.EnqueueType6(1, ownerCharID, zMap.MapID, int(dx), int(dy), 0, int(info.ID), int64(qty))
	}

	if di != nil && killer != nil && killer.IsAutoPick && (allCanLoot || di.OwnerCharID == 0 || di.OwnerCharID == killer.Char.ID) {
		h.autoPickDroppedItem(zMap, killer, di)
	}
}

// RollItemRandomOptions genera opciones aleatorias oficiales desde RandomOptionCount.shn + RandomOption.shn + RandomOptionTable.json.
func (h *Handler) RollItemRandomOptions(itemInx string) []item.ItemStatOption {
	if h == nil || h.dataStore == nil || itemInx == "" {
		return nil
	}
	key := strings.ToLower(strings.TrimSpace(itemInx))
	roc := getRandomOptCatalog(h.dataStore)

	numOpts := 0
	if roc != nil {
		if cntRules := roc.CountsByItem[key]; len(cntRules) > 0 {
			var totalWeight int
			for _, r := range cntRules {
				totalWeight += int(r.DropRate)
			}
			if totalWeight > 0 {
				roll := rand.Intn(totalWeight)
				acc := 0
				for _, r := range cntRules {
					acc += int(r.DropRate)
					if roll < acc {
						numOpts = int(r.LimitCount)
						break
					}
				}
			}
		}
		if numOpts > 0 {
			rules := roc.RulesByItem[key]
			if len(rules) > 0 {
				if numOpts > len(rules) {
					numOpts = len(rules)
				}
				usedType := make(map[uint8]bool)
				var totalWeight int
				for _, r := range rules {
					totalWeight += int(r.DropRate)
				}
				var out []item.ItemStatOption
				for len(out) < numOpts && totalWeight > 0 {
					roll := rand.Intn(totalWeight)
					acc := 0
					var picked *RandomOptRule
					for i := range rules {
						acc += int(rules[i].DropRate)
						if roll < acc {
							picked = &rules[i]
							break
						}
					}
					if picked == nil || usedType[picked.OptionType] {
						// Seleccionar secuencialmente si colisiona tipo repetido
						found := false
						for i := range rules {
							if !usedType[rules[i].OptionType] {
								picked = &rules[i]
								found = true
								break
							}
						}
						if !found {
							break
						}
					}
					usedType[picked.OptionType] = true
					val := picked.Min
					if picked.Max > picked.Min {
						val += uint32(rand.Intn(int(picked.Max - picked.Min + 1)))
					}
					if val > 0 {
						out = append(out, item.ItemStatOption{Type: picked.OptionType, Value: uint16(val)})
					}
				}
				if len(out) > 0 {
					return out
				}
			}
		}
	}

	// Fallback a RandomOptionTable.json en ItemStore
	if h.itemStore != nil {
		if entry, ok := h.itemStore.GetRandomOption(key); ok && entry != nil && entry.MaxOptions > 0 {
			n := entry.MinOptions
			if entry.MaxOptions > entry.MinOptions {
				n += rand.Intn(entry.MaxOptions - entry.MinOptions + 1)
			}
			type statDef struct {
				optType uint8
				rng     item.StatRange
			}
			pool := []statDef{
				{0, entry.STR},
				{1, entry.CON},
				{2, entry.DEX},
				{3, entry.INT},
				{4, entry.MEN},
			}
			var avail []statDef
			for _, s := range pool {
				if s.rng.Max > 0 {
					avail = append(avail, s)
				}
			}
			rand.Shuffle(len(avail), func(i, j int) { avail[i], avail[j] = avail[j], avail[i] })
			if n > len(avail) {
				n = len(avail)
			}
			out := make([]item.ItemStatOption, 0, n)
			for i := 0; i < n; i++ {
				val := avail[i].rng.Min
				if avail[i].rng.Max > avail[i].rng.Min {
					val += rand.Intn(avail[i].rng.Max - avail[i].rng.Min + 1)
				}
				if val > 0 {
					out = append(out, item.ItemStatOption{Type: avail[i].optType, Value: uint16(val)})
				}
			}
			return out
		}
	}
	return nil
}

func (h *Handler) resolveMoneyItemID() uint16 {
	if h.dataStore != nil {
		if f := h.dataStore.GetSHN("ItemMoney"); f != nil && len(f.Rows) > 0 {
			for _, row := range f.Rows {
				if id := row.GetUint16("ItemID"); id > 0 {
					return id
				}
				if id := row.GetUint16("ID"); id > 0 {
					return id
				}
				if inx := strings.TrimSpace(row.GetString("InxName")); inx != "" {
					if info, ok := h.itemStore.GetItemInfoByName(inx); ok && info != nil {
						return info.ID
					}
				}
			}
		}
	}
	if h.itemStore != nil {
		for _, cand := range []string{"Money", "Money01", "Money00", "Gold"} {
			if info, ok := h.itemStore.GetItemInfoByName(cand); ok && info != nil {
				return info.ID
			}
		}
	}
	return 20 // Money01 fallback ID estándar cuando ItemMoney.shn sólo define rangos gráficos
}

// autoPickDroppedItem recoge automáticamente un ítem o Cen recién tirado para jugadores con AutoPick activo.
func (h *Handler) autoPickDroppedItem(zMap *ZoneMap, player *ZonePlayer, di *ZoneDroppedItem) {
	if zMap == nil || player == nil || player.Char == nil || di == nil {
		return
	}
	if di.IsCenDrop {
		consumed, ok := zMap.ConsumeDroppedItem(di.Handle)
		if !ok || consumed == nil {
			return
		}
		amt := int64(consumed.Lot)
		if amt <= 0 {
			amt = 1
		}
		beforeMoney := player.Char.TCharacterState.Money
		if h.charRepo != nil {
			_ = h.charRepo.AddCharMoney(player.Char.ID, amt)
		}
		if player.Char.TCharacterState.Money == beforeMoney {
			player.Char.TCharacterState.Money += amt
		}
		if player.Client != nil {
			cenCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CENCHANGE_CMD)
			cenCmd.WriteUint64(uint64(player.Char.TCharacterState.Money))
			_ = player.Client.SendPacket(cenCmd)
		}
		return
	}

	if h.charRepo == nil {
		return
	}
	freeSlot, hasSlot := h.findFreeInvenSlot(player.Char.ID)
	if !hasSlot {
		return
	}
	consumed, ok := zMap.ConsumeDroppedItem(di.Handle)
	if !ok || consumed == nil {
		return
	}
	lot := int32(consumed.Lot)
	if lot <= 0 {
		lot = 1
	}
	newItem := &database.Item{
		StorageType: 9,
		Owner:       player.Char.ID,
		Storage:     freeSlot,
		ItemID:      consumed.ItemID,
		Flags:       lot,
	}
	if err := h.charRepo.CreateItem(newItem); err != nil {
		return
	}
	for _, opt := range consumed.StatOptions {
		_ = h.charRepo.SetItemOption(newItem.ItemKey, int16(opt.Type), int64(opt.Value))
	}
	if player.Client != nil {
		h.sendCellChangeWithOptions(player.Client, 9, freeSlot, newItem, consumed.StatOptions, 0xFFFF)
	}
}

func (h *Handler) sendCellChangeWithOptions(client *network.Client, box uint8, slot int16, it *database.Item, opts []item.ItemStatOption, exchange uint16) {
	if client == nil {
		return
	}
	if len(opts) == 0 {
		h.sendCellChange(client, box, slot, it, exchange)
		return
	}
	loc := (uint16(box) << 10) | (uint16(slot) & 0x03FF)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CELLCHANGE_CMD)
	pw.WriteUint16(exchange)
	pw.WriteUint16(loc)
	if it == nil || it.ItemID == 0xFFFF {
		pw.WriteUint16(0xFFFF)
	} else {
		pw.WriteUint16(it.ItemID)
		var attrBuf []byte
		if h.itemStore != nil {
			attrBuf = h.itemStore.BuildItemAttr(it.ItemID, uint32(it.Flags), opts, 0)
		} else {
			attrBuf = []byte{byte(it.Flags)}
		}
		pw.WriteBytes(attrBuf)
	}
	_ = client.SendPacket(pw)
}

func (h *Handler) broadcastMobKillAnnounce(zMap *ZoneMap, textIndex uint32) {
	cliCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_CLIENT_MOB_KILL_ANNOUNCE_CMD)
	cliCmd.WriteUint32(textIndex)
	if zMap != nil {
		zMap.Broadcast(cliCmd, 0)
	}

	wldCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_WORLD_MOB_KILL_ANNOUNCE_CMD)
	wldCmd.WriteUint32(textIndex)
	h.mu.RLock()
	clients := make([]*network.Client, 0, len(h.players))
	for c := range h.players {
		clients = append(clients, c)
	}
	h.mu.RUnlock()
	for _, c := range clients {
		_ = c.SendPacket(wldCmd)
	}
}

// handleAutoPickReq implementa PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ (0x10BC -> CALL usp_Char_SetAutoPick).
func (h *Handler) handleAutoPickReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}
	bEnable, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	enabled := bEnable != 0

	var errCode uint16
	if h.charRepo != nil {
		if err := h.charRepo.SetAutoPick(player.Char.ID, enabled); err != nil {
			log.Printf("[Zone][AutoPick] Error SetAutoPick char %d: %v", player.Char.ID, err)
			errCode = 1
		}
	}
	if errCode == 0 {
		player.IsAutoPick = enabled
	}

	// PROTO_NC_CHAR_CLIENT_AUTO_PICK_ACK (0x10BD, 3B): err(2), bEnable(1)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_AUTO_PICK_ACK)
	ack.WriteUint16(errCode)
	if enabled {
		ack.WriteByte(1)
	} else {
		ack.WriteByte(0)
	}
	_ = client.SendPacket(ack)

	// PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD (0x10BE, 3B): player(2), bEnable(1)
	if errCode == 0 {
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD)
		cmd.WriteUint16(player.MapObjectID)
		if enabled {
			cmd.WriteByte(1)
		} else {
			cmd.WriteByte(0)
		}
		if player.CurrentMap != nil && player.CurrentMap.grid != nil {
			px, py := player.GetPosition()
			player.CurrentMap.grid.BroadcastAOI(cmd, px, py, 0)
		} else {
			_ = client.SendPacket(cmd)
		}
	}
}

// Unused import safeguard cleanly referencing drop package type if referenced indirectly
var _ = (*drop.MobDropTable)(nil)
