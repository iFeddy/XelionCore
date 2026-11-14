package zone

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// CharacterTitleStage representa uno de los 4 niveles (0..3) de un título en CharacterTitleData.shn.
type CharacterTitleStage struct {
	ElementNo byte
	Title     string
	Value     uint32
	Fame      uint32
}

// CharacterTitleStateDef representa un buff pasivo asociado a un título en CharacterTitleStateServer.shn.
type CharacterTitleStateDef struct {
	Type      uint8
	TitleLV   uint8
	StateName string
	Strength  uint8
}

// CharacterTitleDef almacena la definición completa de un título oficial desde CharacterTitleData.shn.
type CharacterTitleDef struct {
	Type    uint8
	Permit  int
	Refresh int
	Stages  [4]CharacterTitleStage
	States  map[uint8]CharacterTitleStateDef // TitleLV -> estado
}

// TitleCatalog indexa CharacterTitleData.shn, CharacterTitleStateServer.shn y CollectCardTitle.shn.
type TitleCatalog struct {
	ByType             map[uint8]*CharacterTitleDef
	KillTitleTypes     []uint8          // Tipos que avanzan por matar monstruos o jefes
	LevelTitleTypes    []uint8          // Tipos que avanzan por nivel alcanzado
	CardTitleByItemInx map[string]uint8 // CC_ItemInx -> CharacterTitle Type (CollectCardTitle.shn)
}

// BuildTitleCatalog construye el índice oficial de títulos desde DataStore.
func BuildTitleCatalog(ds *DataStore) *TitleCatalog {
	cat := &TitleCatalog{
		ByType:             make(map[uint8]*CharacterTitleDef),
		CardTitleByItemInx: make(map[string]uint8),
	}
	if ds == nil {
		return cat
	}

	// 1. CharacterTitleData.shn
	if f := ds.GetSHN("CharacterTitleData"); f != nil {
		for _, row := range f.Rows {
			tVal := row.GetInt("Type")
			if tVal < 0 || tVal > 255 {
				continue
			}
			t := uint8(tVal)
			def := &CharacterTitleDef{
				Type:    t,
				Permit:  row.GetInt("Permit"),
				Refresh: row.GetInt("Refresh"),
				States:  make(map[uint8]CharacterTitleStateDef),
			}
			def.Stages[0] = CharacterTitleStage{
				ElementNo: 0,
				Title:     strings.TrimSpace(row.GetString("Title0")),
				Value:     uint32(row.GetInt("Value0")),
				Fame:      uint32(row.GetInt("Fame0")),
			}
			def.Stages[1] = CharacterTitleStage{
				ElementNo: 1,
				Title:     strings.TrimSpace(row.GetString("Title1")),
				Value:     uint32(row.GetInt("Value1")),
				Fame:      uint32(row.GetInt("Fame1")),
			}
			def.Stages[2] = CharacterTitleStage{
				ElementNo: 2,
				Title:     strings.TrimSpace(row.GetString("Title2")),
				Value:     uint32(row.GetInt("Value2")),
				Fame:      uint32(row.GetInt("Fame2")),
			}
			def.Stages[3] = CharacterTitleStage{
				ElementNo: 3,
				Title:     strings.TrimSpace(row.GetString("Title3")),
				Value:     uint32(row.GetInt("Value3")),
				Fame:      uint32(row.GetInt("Fame3")),
			}
			cat.ByType[t] = def

			// Clasificar heurísticamente por prefijo/contenido de título oficial si hay metas > 0
			lower0 := strings.ToLower(def.Stages[0].Title)
			if strings.Contains(lower0, "lv") || strings.Contains(lower0, "level") {
				cat.LevelTitleTypes = append(cat.LevelTitleTypes, t)
			} else if strings.Contains(lower0, "kill") || strings.Contains(lower0, "hunt") || strings.Contains(lower0, "slime") || strings.Contains(lower0, "mob") || strings.Contains(lower0, "boss") {
				cat.KillTitleTypes = append(cat.KillTitleTypes, t)
			}
		}
	}

	// 2. CharacterTitleStateServer.shn
	if f := ds.GetSHN("CharacterTitleStateServer"); f != nil {
		for _, row := range f.Rows {
			tVal := row.GetInt("Type")
			if tVal < 0 || tVal > 255 {
				continue
			}
			t := uint8(tVal)
			lv := row.GetByte("TitleLV")
			if def := cat.ByType[t]; def != nil {
				def.States[lv] = CharacterTitleStateDef{
					Type:      t,
					TitleLV:   lv,
					StateName: strings.TrimSpace(row.GetString("StateName")),
					Strength:  row.GetByte("Strength"),
				}
			}
		}
	}

	// 3. CollectCardTitle.shn
	if f := ds.GetSHN("CollectCardTitle"); f != nil {
		for _, row := range f.Rows {
			tVal := row.GetInt("Type")
			itemInx := strings.ToLower(strings.TrimSpace(row.GetString("CC_ItemInx")))
			if itemInx != "" && tVal >= 0 && tVal <= 255 {
				cat.CardTitleByItemInx[itemInx] = uint8(tVal)
			}
		}
	}

	return cat
}

// CollectCardDef representa una carta de colección desde CollectCard.shn + CollectCardDropRate.shn.
type CollectCardDef struct {
	CardID        uint16
	ItemInx       string
	CardGradeType int
	CardMobGroup  int
	DropRate      uint16
}

// CollectRewardDef representa una recompensa de álbum desde CollectCardReward.shn.
type CollectRewardDef struct {
	RewardID       uint16
	CardRewardType int
	CardLot        uint16
	RewardItemInx  string
	RewardLot      uint16
}

// CollectStarRule representa la probabilidad de estrella según el grado en CollectCardStarRate.shn.
type CollectStarRule struct {
	CardGradeType int
	StarLot       uint8
	StarRate      uint16
}

// CollectCatalog indexa todas las tablas CollectCard*.shn oficiales.
type CollectCatalog struct {
	CardsByID      map[uint16]*CollectCardDef
	CardsByItemInx map[string]*CollectCardDef
	MobGroupByMob  map[string][]int
	CardsByMobGrp  map[int][]*CollectCardDef
	RewardsByID    map[uint16]*CollectRewardDef
	StarRatesByGrd map[int][]CollectStarRule
}

// BuildCollectCatalog construye el índice completo de cartas de colección desde DataStore.
func BuildCollectCatalog(ds *DataStore) *CollectCatalog {
	cat := &CollectCatalog{
		CardsByID:      make(map[uint16]*CollectCardDef),
		CardsByItemInx: make(map[string]*CollectCardDef),
		MobGroupByMob:  make(map[string][]int),
		CardsByMobGrp:  make(map[int][]*CollectCardDef),
		RewardsByID:    make(map[uint16]*CollectRewardDef),
		StarRatesByGrd: make(map[int][]CollectStarRule),
	}
	if ds == nil {
		return cat
	}

	dropRates := make(map[uint16]uint16)
	if f := ds.GetSHN("CollectCardDropRate"); f != nil {
		for _, row := range f.Rows {
			id := row.GetUint16("CC_CardID")
			rate := row.GetUint16("CC_CardGetRate")
			if id > 0 {
				dropRates[id] = rate
			}
		}
	}

	if f := ds.GetSHN("CollectCard"); f != nil {
		for _, row := range f.Rows {
			cardID := row.GetUint16("CC_CardID")
			itemInx := strings.TrimSpace(row.GetString("CC_ItemInx"))
			if cardID == 0 || itemInx == "" {
				continue
			}
			grp := row.GetInt("CC_CardMobGroup")
			c := &CollectCardDef{
				CardID:        cardID,
				ItemInx:       itemInx,
				CardGradeType: row.GetInt("CC_CardGradeType"),
				CardMobGroup:  grp,
				DropRate:      dropRates[cardID],
			}
			cat.CardsByID[cardID] = c
			cat.CardsByItemInx[strings.ToLower(itemInx)] = c
			cat.CardsByMobGrp[grp] = append(cat.CardsByMobGrp[grp], c)
		}
	}

	if f := ds.GetSHN("CollectCardMobGroup"); f != nil {
		for _, row := range f.Rows {
			mobInx := strings.ToLower(strings.TrimSpace(row.GetString("CC_MobInx")))
			grp := row.GetInt("CC_CardMobGroup")
			if mobInx != "" && grp != 0 {
				cat.MobGroupByMob[mobInx] = append(cat.MobGroupByMob[mobInx], grp)
			}
		}
	}

	if f := ds.GetSHN("CollectCardReward"); f != nil {
		for _, row := range f.Rows {
			rwID := row.GetUint16("CC_RewardID")
			if rwID == 0 {
				continue
			}
			cat.RewardsByID[rwID] = &CollectRewardDef{
				RewardID:       rwID,
				CardRewardType: row.GetInt("CC_CardRewardType"),
				CardLot:        row.GetUint16("CC_CardLot"),
				RewardItemInx:  strings.TrimSpace(row.GetString("CC_RewardItemInx")),
				RewardLot:      row.GetUint16("CC_RewardLot"),
			}
		}
	}

	if f := ds.GetSHN("CollectCardStarRate"); f != nil {
		for _, row := range f.Rows {
			grd := row.GetInt("CC_CardGradeType")
			cat.StarRatesByGrd[grd] = append(cat.StarRatesByGrd[grd], CollectStarRule{
				CardGradeType: grd,
				StarLot:       row.GetByte("CC_StarLot"),
				StarRate:      row.GetUint16("CC_StarRate"),
			})
		}
	}

	return cat
}

// CardsForMob devuelve todas las cartas asociables a un MobInx oficial.
func (c *CollectCatalog) CardsForMob(mobInx string) []*CollectCardDef {
	if c == nil || mobInx == "" {
		return nil
	}
	groups := c.MobGroupByMob[strings.ToLower(strings.TrimSpace(mobInx))]
	if len(groups) == 0 {
		return nil
	}
	var out []*CollectCardDef
	for _, g := range groups {
		out = append(out, c.CardsByMobGrp[g]...)
	}
	return out
}

// RollStar calcula las estrellas (1..N) de una carta registrada usando CollectCardStarRate.shn.
func (c *CollectCatalog) RollStar(gradeType int) byte {
	if c == nil {
		return 1
	}
	rules := c.StarRatesByGrd[gradeType]
	if len(rules) == 0 {
		return 1
	}
	var total int
	for _, r := range rules {
		total += int(r.StarRate)
	}
	if total <= 0 {
		return 1
	}
	roll := rand.Intn(total)
	acc := 0
	for _, r := range rules {
		acc += int(r.StarRate)
		if roll < acc {
			if r.StarLot == 0 {
				return 1
			}
			return r.StarLot
		}
	}
	return 1
}

// ---------- Dept 24: Títulos de Personaje (PROTO_NC_CT_*) ----------

// handleCTSetCurrentReq implementa PROTO_NC_CT_SET_CURRENT_REQ (0x6001):
// Cambia el título equipado activo del jugador y difunde PROTO_NC_CT_SET_SOMEONECHANGE_CMD (0x6003).
func (h *Handler) handleCTSetCurrentReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	titleType, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	elementNo, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	mobID, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	var errCode uint16
	if titleType != 0 || elementNo != 0 {
		// Validar título oficial y posesión en DB
		if h.dataStore != nil && h.dataStore.TitleCatalog != nil {
			if def := h.dataStore.TitleCatalog.ByType[titleType]; def == nil {
				errCode = 1
			}
		}
		if errCode == 0 && h.charRepo != nil {
			titles, err := h.charRepo.GetCharacterTitles(player.Char.ID)
			if err != nil {
				errCode = 1
			} else {
				found := false
				for _, t := range titles {
					if t != nil && t.TitleType == int16(titleType) {
						if t.ElementNo >= elementNo || (t.ElementNo == 0 && elementNo == 0) {
							found = true
							break
						}
					}
				}
				if !found {
					errCode = 1
				}
			}
		}
	}

	if errCode == 0 {
		player.CurrentTitleType = titleType
		player.CurrentTitleElement = elementNo
		player.CurrentTitleMobID = mobID
		if h.charRepo != nil {
			_ = h.charRepo.SetCharacterTitle(player.Char.ID, int16(titleType), fmt.Sprintf("%d:%d", elementNo, mobID))
		}
	}

	// PROTO_NC_CT_SET_CURRENT_ACK (0x6002, 6B): Type(1), ElementNo(1), MobID(2), Error(2)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_CT_SET_CURRENT_ACK)
	ack.WriteByte(titleType)
	ack.WriteByte(elementNo)
	ack.WriteUint16(mobID)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)

	// PROTO_NC_CT_SET_SOMEONECHANGE_CMD (0x6003, 6B): hnd(2), Type(1), ElementNo(1), MobID(2)
	if errCode == 0 {
		bcast := protocol.NewPacketWriter(protocol.PROTO_NC_CT_SET_SOMEONECHANGE_CMD)
		bcast.WriteUint16(player.MapObjectID)
		bcast.WriteByte(titleType)
		bcast.WriteByte(elementNo)
		bcast.WriteUint16(mobID)
		if player.CurrentMap != nil && player.CurrentMap.grid != nil {
			px, py := player.GetPosition()
			player.CurrentMap.grid.BroadcastAOI(bcast, px, py, player.MapObjectID)
		}
	}
}

// UnlockOrAdvanceTitle registra o incrementa progreso hacia un título en CharacterTitleData.shn
// y emite PROTO_NC_CT_SET_CMD (0x6004) cuando se alcanza un nuevo nivel del título.
func (h *Handler) UnlockOrAdvanceTitle(p *ZonePlayer, titleType uint8, deltaOrAbsolute int64, isAbsolute bool) {
	if h == nil || p == nil || p.Char == nil {
		return
	}
	var def *CharacterTitleDef
	if h.dataStore != nil && h.dataStore.TitleCatalog != nil {
		def = h.dataStore.TitleCatalog.ByType[titleType]
	}

	var curElement byte
	var curVal int64
	hasExisting := false
	if h.charRepo != nil {
		if existing, err := h.charRepo.GetCharacterTitles(p.Char.ID); err == nil {
			for _, t := range existing {
				if t != nil && t.TitleType == int16(titleType) {
					curElement = t.ElementNo
					curVal = t.ElementValue
					hasExisting = true
					break
				}
			}
		}
	}

	newVal := curVal
	if isAbsolute {
		if deltaOrAbsolute > newVal {
			newVal = deltaOrAbsolute
		}
	} else {
		newVal += deltaOrAbsolute
	}

	newElem := curElement
	var fameGain uint32
	unlockedNewStage := false
	if def != nil {
		for idx, st := range def.Stages {
			if st.Title == "" && st.Value == 0 {
				continue
			}
			if newVal >= int64(st.Value) {
				stageElem := byte(idx)
				if !hasExisting || stageElem > newElem {
					newElem = stageElem
					fameGain += st.Fame
					unlockedNewStage = true
				}
			}
		}
	} else if !hasExisting {
		unlockedNewStage = true
	}

	if !hasExisting && !unlockedNewStage && newVal <= curVal {
		return
	}

	if h.charRepo != nil {
		if err := h.charRepo.SetCharacterTitleProgress(p.Char.ID, titleType, newElem, newVal); err != nil {
			log.Printf("[Zone][Title] Error SetCharacterTitleProgress char %d type %d: %v", p.Char.ID, titleType, err)
			return
		}
	}

	if fameGain > 0 && h.charRepo != nil {
		beforeFame := p.Char.TCharacterState.Fame
		_ = h.charRepo.AddCharFame(p.Char.ID, int(fameGain))
		if p.Char.TCharacterState.Fame == beforeFame {
			p.Char.TCharacterState.Fame += int(fameGain)
		}
		if p.Client != nil {
			famePkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_FAMECHANGE_CMD)
			famePkt.WriteUint32(uint32(p.Char.TCharacterState.Fame))
			_ = p.Client.SendPacket(famePkt)
		}
	}

	if unlockedNewStage && p.Client != nil {
		// PROTO_NC_CT_SET_CMD (0x6004, 2B): Type(1), ElementNo(1)
		cmd := protocol.NewPacketWriter(protocol.PROTO_NC_CT_SET_CMD)
		cmd.WriteByte(titleType)
		cmd.WriteByte(newElem)
		_ = p.Client.SendPacket(cmd)
	}
}

// OnMonsterKilledTitles avanza los títulos asociados a cacería de monstruos según CharacterTitleData.shn.
func (h *Handler) OnMonsterKilledTitles(p *ZonePlayer, mob *ZoneMonster) {
	if h == nil || p == nil || mob == nil || h.dataStore == nil || h.dataStore.TitleCatalog == nil {
		return
	}
	tc := h.dataStore.TitleCatalog
	for _, tType := range tc.KillTitleTypes {
		def := tc.ByType[tType]
		if def == nil {
			continue
		}
		// Coincidencia de familia de jefe/slime o contador general de muerte de monstruos
		lowerMob := strings.ToLower(mob.InxName)
		lowerT := strings.ToLower(def.Stages[0].Title)
		if strings.Contains(lowerT, "slime") && !strings.Contains(lowerMob, "slime") {
			continue
		}
		if strings.Contains(lowerT, "boss") && mob.Rank < 4 {
			continue
		}
		h.UnlockOrAdvanceTitle(p, tType, 1, false)
	}
}

// OnPlayerLevelUpTitles avanza los títulos de nivel al alcanzar hitos de CharacterTitleData.shn.
func (h *Handler) OnPlayerLevelUpTitles(p *ZonePlayer) {
	if h == nil || p == nil || p.Char == nil || h.dataStore == nil || h.dataStore.TitleCatalog == nil {
		return
	}
	for _, tType := range h.dataStore.TitleCatalog.LevelTitleTypes {
		h.UnlockOrAdvanceTitle(p, tType, int64(p.Char.Level), true)
	}
}

// ---------- Dept 49: Colección de Cartas (PROTO_NC_COLLECT_*) ----------

// writeCardCollectBlock serializa PROTO_NC_CHAR_CARDCOLLECT_CMD: chrregnum(4), number(2), []CardInform(11B c/u).
func writeCardCollectBlock(pw *protocol.PacketWriter, charID uint32, cards []*database.CollectionCardData) {
	pw.WriteUint32(charID)
	pw.WriteUint16(uint16(len(cards)))
	for _, c := range cards {
		var rt uint32
		if !c.RegistTime.IsZero() {
			rt = database.PackShineDateTime(c.RegistTime)
		}
		pw.WriteUint32(rt)
		pw.WriteUint32(uint32(c.SerialNo))
		pw.WriteUint16(uint16(c.CardID))
		pw.WriteByte(c.StarNumber)
	}
}

// writeCardBookmarkBlock serializa PROTO_NC_CHAR_CARDCOLLECT_BOOKMARK_CMD: chrregnum(4), number(2), []Bookmark(4B c/u).
func writeCardBookmarkBlock(pw *protocol.PacketWriter, charID uint32, bookmarks []*database.CollectionBookmarkData) {
	pw.WriteUint32(charID)
	pw.WriteUint16(uint16(len(bookmarks)))
	for _, b := range bookmarks {
		pw.WriteUint16(uint16(b.Slot))
		pw.WriteUint16(uint16(b.CardID))
	}
}

// sendPlayerCollectionOnLogin envía las 3 tablas de colección del personaje al entrar al mapa.
func (h *Handler) sendPlayerCollectionOnLogin(p *ZonePlayer) {
	if h == nil || p == nil || p.Client == nil || p.Char == nil || h.charRepo == nil {
		return
	}
	cards, _ := h.charRepo.GetCollectionCards(p.Char.ID)
	if len(cards) > 0 {
		pkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_CARDCOLLECT_CMD)
		writeCardCollectBlock(pkt, uint32(p.Char.ID), cards)
		_ = p.Client.SendPacket(pkt)
	}

	bookmarks, _ := h.charRepo.GetCollectionFavorites(p.Char.ID)
	if len(bookmarks) > 0 {
		bmPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_CARDCOLLECT_BOOKMARK_CMD)
		writeCardBookmarkBlock(bmPkt, uint32(p.Char.ID), bookmarks)
		_ = p.Client.SendPacket(bmPkt)
	}

	rewards, _ := h.charRepo.GetCollectionRewards(p.Char.ID)
	if len(rewards) > 0 {
		rwPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_CARDCOLLECT_REWARD_CMD)
		rwPkt.WriteUint32(uint32(p.Char.ID))
		rwPkt.WriteUint16(uint16(len(rewards)))
		for _, rwID := range rewards {
			rwPkt.WriteUint16(uint16(rwID))
		}
		_ = p.Client.SendPacket(rwPkt)
	}
}

// handleCollectPeepingReq implementa PROTO_NC_COLLECT_PEEPING_REQ (0xC401):
// Inspecciona la colección de cartas y favoritos de otro jugador y responde 0xC405 y 0xC406.
func (h *Handler) handleCollectPeepingReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	targetName, err := pkt.Reader.ReadString(20)
	if err != nil {
		return
	}

	var targetCharID int
	var errCode uint16
	if tp, found := h.GetPlayerByName(targetName); found && tp != nil && tp.Char != nil {
		targetCharID = tp.Char.ID
	} else if h.charRepo != nil {
		if ch, err := h.charRepo.GetCharacterByName(targetName); err == nil && ch != nil {
			targetCharID = ch.ID
		} else {
			errCode = 1
		}
	} else {
		errCode = 1
	}

	var cards []*database.CollectionCardData
	var bookmarks []*database.CollectionBookmarkData
	if errCode == 0 && h.charRepo != nil {
		cards, _ = h.charRepo.GetCollectionCards(targetCharID)
		bookmarks, _ = h.charRepo.GetCollectionFavorites(targetCharID)
	}

	// 1. PROTO_NC_COLLECT_PEEPING_COLLECT_ACK (0xC405)
	ack1 := protocol.NewPacketWriter(protocol.PROTO_NC_COLLECT_PEEPING_COLLECT_ACK)
	ack1.WriteString(targetName, 20)
	ack1.WriteUint16(errCode)
	writeCardCollectBlock(ack1, uint32(targetCharID), cards)
	_ = client.SendPacket(ack1)

	// 2. PROTO_NC_COLLECT_PEEPING_COLLECT_BOOKMARK_ACK (0xC406)
	ack2 := protocol.NewPacketWriter(protocol.PROTO_NC_COLLECT_PEEPING_COLLECT_BOOKMARK_ACK)
	ack2.WriteString(targetName, 20)
	ack2.WriteUint16(errCode)
	writeCardBookmarkBlock(ack2, uint32(targetCharID), bookmarks)
	_ = client.SendPacket(ack2)
}

// handleCollectCardRegistReq implementa PROTO_NC_COLLECT_CARDREGIST_REQ (0xC408):
// Registra una carta desde el inventario (Box 9), calcula su estrella desde CollectCardStarRate.shn,
// persiste serial/registro vía usp_Collection_* y responde PROTO_NC_COLLECT_CARDREGIST_ACK (0xC40B).
func (h *Handler) handleCollectCardRegistReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	invenSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	sendAck := func(cardID uint16, star byte, serial uint32, errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_COLLECT_CARDREGIST_ACK)
		ack.WriteByte(invenSlot)
		ack.WriteUint16(cardID)
		ack.WriteByte(star)
		ack.WriteUint32(serial)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	srcItem := h.findItemAt(player.Char.ID, 9, int16(invenSlot))
	if srcItem == nil {
		sendAck(0, 0, 0, 1)
		return
	}

	if h.dataStore == nil || h.dataStore.CollectCatalog == nil || h.itemStore == nil {
		sendAck(0, 0, 0, 1)
		return
	}
	info, ok := h.itemStore.GetItemInfo(srcItem.ItemID)
	if !ok || info == nil {
		sendAck(0, 0, 0, 1)
		return
	}
	cardDef := h.dataStore.CollectCatalog.CardsByItemInx[strings.ToLower(info.InxName)]
	if cardDef == nil {
		sendAck(0, 0, 0, 1)
		return
	}

	star := h.dataStore.CollectCatalog.RollStar(cardDef.CardGradeType)

	// Generar número de serie global del servidor para la carta
	var nextSerial int64 = 1
	if serials, err := h.charRepo.GetCollectionSerialList(); err == nil {
		nextSerial = serials[int(cardDef.CardID)] + 1
	}
	if err := h.charRepo.SetCollectionSerialCount(int(cardDef.CardID), nextSerial); err != nil {
		log.Printf("[Zone][Collect] SetCollectionSerialCount err: %v", err)
	}

	if err := h.charRepo.RegisterCollectionCard(player.Char.ID, int(cardDef.CardID), nextSerial, star, srcItem.ItemKey); err != nil {
		log.Printf("[Zone][Collect] RegisterCollectionCard err: %v", err)
		sendAck(cardDef.CardID, star, uint32(nextSerial), 1)
		return
	}

	// Consumir 1 carta del inventario
	if srcItem.Flags > 1 {
		srcItem.Flags--
		_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(srcItem.Flags))
		h.sendCellChange(client, 9, int16(invenSlot), srcItem, 0xFFFF)
	} else {
		_ = h.charRepo.DeleteItem(srcItem.ItemKey)
		h.sendCellChange(client, 9, int16(invenSlot), nil, 0xFFFF)
	}

	// Desbloquear título asociado si la carta figura en CollectCardTitle.shn
	if h.dataStore.TitleCatalog != nil {
		if titleType, hasTitle := h.dataStore.TitleCatalog.CardTitleByItemInx[strings.ToLower(cardDef.ItemInx)]; hasTitle {
			h.UnlockOrAdvanceTitle(player, titleType, 1, true)
		}
	}

	sendAck(cardDef.CardID, star, uint32(nextSerial), 0)
}

// handleCollectBookmarkRegistReq implementa PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ (0xC40C).
func (h *Handler) handleCollectBookmarkRegistReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	slot, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	cardID, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	var errCode uint16
	regist := cardID != 0
	if regist {
		cards, _ := h.charRepo.GetCollectionCards(player.Char.ID)
		owned := false
		for _, c := range cards {
			if c != nil && uint16(c.CardID) == cardID {
				owned = true
				break
			}
		}
		if !owned {
			errCode = 1
		}
	}

	if errCode == 0 {
		if err := h.charRepo.SetCollectionFavorite(player.Char.ID, byte(slot), int(cardID), regist); err != nil {
			errCode = 1
		}
	}

	// PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK (0xC40F, 6B): BookMarkSlot(2), CardID(2), err(2)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK)
	ack.WriteUint16(slot)
	ack.WriteUint16(cardID)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

// handleCollectRewardReq implementa PROTO_NC_COLLECT_REWARD_REQ (0xC410).
func (h *Handler) handleCollectRewardReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	rewardID, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	_, _ = pkt.Reader.ReadUint16() // RewardItem solicitado por cliente (validado contra CollectCardReward.shn)

	sendAck := func(errCode uint16) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_COLLECT_REWARD_ACK)
		ack.WriteUint16(rewardID)
		ack.WriteUint16(errCode)
		_ = client.SendPacket(ack)
	}

	if h.dataStore == nil || h.dataStore.CollectCatalog == nil || h.itemStore == nil {
		sendAck(1)
		return
	}
	rwDef := h.dataStore.CollectCatalog.RewardsByID[rewardID]
	if rwDef == nil {
		sendAck(1)
		return
	}

	// Verificar que no haya sido reclamada antes
	claimed, _ := h.charRepo.GetCollectionRewards(player.Char.ID)
	for _, cid := range claimed {
		if uint16(cid) == rewardID {
			sendAck(1)
			return
		}
	}

	// Verificar requisito mínimo de cartas registradas
	cards, _ := h.charRepo.GetCollectionCards(player.Char.ID)
	if len(cards) < int(rwDef.CardLot) {
		sendAck(1)
		return
	}

	rewardItemInfo, ok := h.itemStore.GetItemInfoByName(rwDef.RewardItemInx)
	if !ok || rewardItemInfo == nil {
		sendAck(1)
		return
	}

	freeSlot, hasSlot := h.findFreeInvenSlot(player.Char.ID)
	if !hasSlot {
		sendAck(1)
		return
	}

	lot := int32(rwDef.RewardLot)
	if lot <= 0 {
		lot = 1
	}

	if err := h.charRepo.SetCollectionReward(player.Char.ID, int(rewardID)); err != nil {
		sendAck(1)
		return
	}

	newItem := &database.Item{
		StorageType: 9,
		Owner:       player.Char.ID,
		Storage:     freeSlot,
		ItemID:      rewardItemInfo.ID,
		Flags:       lot,
	}
	if err := h.charRepo.CreateItem(newItem); err != nil {
		sendAck(1)
		return
	}
	h.sendCellChange(client, 9, freeSlot, newItem, 0xFFFF)

	sendAck(0)
	_ = time.Now()
}
