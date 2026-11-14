package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	ErrCharacterNotFound      = errors.New("character not found")
	ErrCharacterExists        = errors.New("character name already exists")
	ErrSlotInUse              = errors.New("character slot already in use")
	ErrSlotOutOfRange         = errors.New("character slot out of range (0-5)")
	ErrCharacterIsGuildMaster = errors.New("character is an active guild master")
	ErrItemNotFound           = errors.New("item not found")
)

// TCharacter represents the canonical tCharacter table from World00_Character.bak.
type TCharacter struct {
	CharNo            int    `json:"nCharNo"`
	UserNo            int    `json:"nUserNo"`
	CharName          string `json:"sCharName"`
	Slot              byte   `json:"nSlot"`
	Class             byte   `json:"nClass"`
	Level             byte   `json:"nLevel"`
	IsDelete          int    `json:"bIsDelete"` // 0 = Active, 1 = Soft Deleted
	Strength          int    `json:"nStrength"`
	Constitute        int    `json:"nConstitute"`
	Dexterity         int    `json:"nDexterity"`
	Intelligence      int    `json:"nIntelligence"`
	MentalPower       int    `json:"nMentalPower"`
	RedistributePoint int    `json:"nRedistributePoint"`
}

// TCharacterShape represents the appearance table tCharacterShape.
type TCharacterShape struct {
	CharNo    int  `json:"nCharNo"`
	Gender    byte `json:"nGender"` // 0 = Female, 1 = Male
	Hair      byte `json:"nHair"`
	HairColor byte `json:"nHairColor"`
	Face      byte `json:"nFace"`
}

// TCharacterState represents the in-world vitals and position from tCharacterState.
type TCharacterState struct {
	CharNo   int    `json:"nCharNo"`
	Map      uint16 `json:"nMap"`
	XPos     int    `json:"nXPos"`
	YPos     int    `json:"nYPos"`
	Rotation byte   `json:"nRotation"`
	HP       int    `json:"nHP"`
	SP       int    `json:"nSP"`
	HPStones int16  `json:"nHPStones"`
	SPStones int16  `json:"nSPStones"`
	Exp      int64  `json:"nExp"`
	Money    int64  `json:"nMoney"`
	Fame     int    `json:"nFame"`
}

// TCharacterKeyMap represents the quickbar settings from tCharacterKeyMap.
type TCharacterKeyMap struct {
	CharNo        int    `json:"nCharNo"`
	QuickBar      []byte `json:"sQuickBar"`
	QuickBarState []byte `json:"sQuickBarState"`
}

// TCharacterShortCut represents shortcut settings from tCharacterShortCut.
type TCharacterShortCut struct {
	CharNo   int    `json:"nCharNo"`
	ShortCut []byte `json:"sShortCut"`
}

// TCharacterGameOption represents game settings from tCharacterGameOption.
type TCharacterGameOption struct {
	CharNo     int    `json:"nCharNo"`
	GameOption []byte `json:"sGameOption"`
}

// TCharacterSingleOption represents client UI settings from tCharacterSingleOption.
type TCharacterSingleOption struct {
	CharNo       int    `json:"nCharNo"`
	SingleOption []byte `json:"sSingleOption"`
}

// Character represents an avatar in XelionCore, embedding the official relational tables.
type Character struct {
	// Canonical relational tables
	TCharacter          TCharacter             `json:"t_character"`
	TCharacterShape     TCharacterShape        `json:"t_character_shape"`
	TCharacterState     TCharacterState        `json:"t_character_state"`
	TCharacterKeyMap    TCharacterKeyMap       `json:"t_character_key_map"`
	TCharacterShortCut  TCharacterShortCut     `json:"t_character_short_cut"`
	TCharacterGameOpt   TCharacterGameOption   `json:"t_character_game_option"`
	TCharacterSingleOpt TCharacterSingleOption `json:"t_character_single_option"`

	// Equips: Slot index (ItemEquipEnum) -> item ID
	Equips map[byte]uint16 `json:"equips"`
	// EquipUpgrades: Slot index (ItemEquipEnum) -> upgrade level (0..15)
	EquipUpgrades map[byte]byte `json:"equip_upgrades,omitempty"`

	// Flat convenience fields for rapid gameplay access
	ID               int                      `json:"id"`
	AccountID        int                      `json:"account_id"`
	Name             string                   `json:"name"`
	Slot             byte                     `json:"slot"`
	Race             byte                     `json:"race"`
	Job              byte                     `json:"job"`
	Level            byte                     `json:"level"`
	AdminLevel       byte                     `json:"admin_level"`
	CreateWorld      byte                     `json:"create_world"`
	Male             bool                     `json:"male"`
	Hair             byte                     `json:"hair"`
	HairColor        byte                     `json:"hair_color"`
	Face             byte                     `json:"face"`
	Map              uint16                   `json:"map"`
	XPos             int                      `json:"x_pos"`
	YPos             int                      `json:"y_pos"`
	CurHP            int                      `json:"cur_hp"`
	CurSP            int                      `json:"cur_sp"`
	TownName         string                   `json:"town_name"`
	DeletedAt        *time.Time               `json:"deleted_at,omitempty"`
	KQHandle         uint32                   `json:"kq_handle,omitempty"`
	KQMapName        string                   `json:"kq_map_name,omitempty"`
	KQX              uint32                   `json:"kq_x,omitempty"`
	KQY              uint32                   `json:"kq_y,omitempty"`
	KQDate           uint32                   `json:"kq_date,omitempty"`
	ChangeID         *CharIDChangeInfo        `json:"change_id,omitempty"`
	TutorialState    int32                    `json:"tutorial_state"`
	TutorialStep     byte                     `json:"tutorial_step"`
	Items            []*Item                  `json:"items"`
	Shortcuts        []*Shortcut              `json:"shortcuts"`
	ImproveShortcuts []*ImproveShortcutEntry  `json:"improve_shortcuts,omitempty"`
	Skills           []*CharacterSkill        `json:"skills"`
	Quests           []*QuestData             `json:"quests,omitempty"`
}

// CharIDChangeInfo modela una fila pendiente de tCharacterChangeID (usp_Character_getChangeIDList).
type CharIDChangeInfo struct {
	RowNo        uint32 `json:"row_no"`
	UserNo       int    `json:"user_no"`
	CharNo       int    `json:"char_no"`
	NeedChangeID byte   `json:"need_change_id"`
	Init         byte   `json:"init"`
}

// ImproveShortcutEntry modela una fila de tCharacterShortCut para el sistema mejorado NA2016 (SHORT_CUT_DATA).
type ImproveShortcutEntry struct {
	SlotNo byte   `json:"slot_no"`
	CodeNo uint16 `json:"code_no"`
	Value  int32  `json:"value"`
}

// ImproveKeyMapEntry modela una fila de tCharacterKeyMap para el sistema mejorado NA2016 (KEY_MAP_DATA).
type ImproveKeyMapEntry struct {
	FunctionNo uint16 `json:"function_no"`
	ExtendKey  byte   `json:"extend_key"`
	ASCIICode  byte   `json:"ascii_code"`
}

// ImproveGameOptionEntry modela una fila de tCharacterGameOption para el sistema mejorado NA2016 (GAME_OPTION_DATA).
type ImproveGameOptionEntry struct {
	OptionNo uint16 `json:"option_no"`
	Value    byte   `json:"value"`
}

// PrisonInfo modela una sanción activa en tPrison (p_Prison_Get).
type PrisonInfo struct {
	Reason string `json:"reason"`
	Remark string `json:"remark"`
}

// Item represents an official tItem database row.
type Item struct {
	ItemKey     int64     `json:"item_key"`
	StorageType byte      `json:"storage_type"` // 0: Equip, 9: Inven, 12: House, 15: Emotions
	Owner       int       `json:"owner"`        // nCharNo
	Storage     int16     `json:"storage"`      // Slot
	ItemID      uint16    `json:"item_id"`
	Flags       int32     `json:"flags"`
	Date        time.Time `json:"date"`
}

// Shortcut represents an entry in tCharacterShortCut.
type Shortcut struct {
	CharNo   int   `json:"char_no"`
	SlotNo   byte  `json:"slot_no"`
	CodeType int16 `json:"code_type"`
	CodeNo   int16 `json:"code_no"`
	Value    int32 `json:"value"`
}

// CharacterSkill represents a learned skill in tSkill.
type CharacterSkill struct {
	CharNo             int    `json:"char_no"`
	SkillNo            uint16 `json:"skill_no"`
	SkillLevel         byte   `json:"skill_level"`
	SkillExp           int    `json:"skill_exp"`
	SkillWriteTime     int64  `json:"skill_write_time"`
	SkillCoolTime      int    `json:"skill_cool_time"`
	SkillPowerDemage   byte   `json:"skill_power_damage"`
	SkillPowerSP       byte   `json:"skill_power_sp"`
	SkillPowerKeepTime byte   `json:"skill_power_keep_time"`
	SkillPowerCoolTime byte   `json:"skill_power_cool_time"`
}

// CharacterStateData modela un estado anormal o buff persistido en tCharacterState.
type CharacterStateData struct {
	CharNo        int   `json:"char_no"`
	StateNo       int16 `json:"state_no"`
	StateLevel    int16 `json:"state_level"`
	StateTime     int   `json:"state_time"`
	RestTime      int   `json:"rest_time"`
	StateCount    int16 `json:"state_count"`
	PowerDamage   byte  `json:"power_damage"`
	PowerSP       byte  `json:"power_sp"`
	PowerKeepTime byte  `json:"power_keep_time"`
	PowerCoolTime byte  `json:"power_cool_time"`
}

// CharacterTitleData modela títulos de personajes en tCharacterTitle.
type CharacterTitleData struct {
	CharNo       int    `json:"char_no"`
	TitleType    int16  `json:"title_type"`
	ElementNo    byte   `json:"element_no"`
	ElementValue int64  `json:"element_value"`
	Title        string `json:"title"`
}

// ChargedBuffData modela buffs pagos de tCharacterChargedBuffer.
type ChargedBuffData struct {
	Key     int       `json:"key"`
	ID      int16     `json:"id"`
	UseDate time.Time `json:"use_date"`
	EndDate time.Time `json:"end_date"`
}

// ItemEnchantData modela ranuras de gemas en tItemEnchant.
type ItemEnchantData struct {
	Slot      byte `json:"slot"`
	GemItemID int  `json:"gem_item_id"`
	RestCount byte `json:"rest_count"`
}

// ItemActionCoolTimeData modela cooltimes de ítems en tItemActionCoolTime.
type ItemActionCoolTimeData struct {
	CharNo      int       `json:"char_no"`
	ActionGroup int16     `json:"action_group"`
	EndTime     time.Time `json:"end_time"`
}

// QuestData modela misiones en tQuest.
type QuestData struct {
	CharNo  int    `json:"char_no"`
	QuestNo int    `json:"quest_no"`
	Status  byte   `json:"status"`
	Data    []byte `json:"data"`
}

// CollectionCardData modela cartas de colección en tCollectionCard.
type CollectionCardData struct {
	CardID     int16     `json:"card_id"`
	Count      int16     `json:"count"`
	SerialNo   int64     `json:"serial_no"`
	StarNumber byte      `json:"star_number"`
	RegistTime time.Time `json:"regist_time"`
}

// CollectionBookmarkData modela cartas favoritas/marcadores en tCollectionFavorites.
type CollectionBookmarkData struct {
	Slot       uint16 `json:"slot"`
	CardID     uint16 `json:"card_id"`
	SerialNo   int64  `json:"serial_no"`
	StarNumber byte   `json:"star_number"`
}

// KingdomQuestData modela datos de Kingdom Quest en tKingdomQuest.
type KingdomQuestData struct {
	QuestNo int    `json:"quest_no"`
	Data    string `json:"data"`
}

// MasterPupilData modela relación maestro-alumno en tMasterPupils.
type MasterPupilData struct {
	CharNo   int  `json:"char_no"`
	IsMaster bool `json:"is_master"`
}

// HolyPromiseData modela relación de promesa sagrada en tHolyPromise.
type HolyPromiseData struct {
	CharNo   int  `json:"char_no"`
	UpCharNo int  `json:"up_char_no"`
	Flags    byte `json:"flags"`
}

// PartyData modela un grupo en tParty.
type PartyData struct {
	PartyNo      int `json:"party_no"`
	MasterCharNo int `json:"master_char_no"`
	Flags        int `json:"flags"`
}

// PartyMemberData modela un miembro de grupo en tPartyMember.
type PartyMemberData struct {
	CharNo    int    `json:"char_no"`
	Flag      byte   `json:"flag"`
	CharID    string `json:"char_id"`
	Class     byte   `json:"class"`
	Level     byte   `json:"level"`
	LoginZone string `json:"login_zone"`
	HP        int16  `json:"hp"`
	SP        int16  `json:"sp"`
}

// PetData modela mascotas en tPet / tPetNew.
type PetData struct {
	PetNo    int    `json:"pet_no"`
	OwnerNo  int    `json:"owner_no"`
	PetName  string `json:"pet_name"`
	Mind     byte   `json:"mind"`
	Stress   byte   `json:"stress"`
	Tendency byte   `json:"tendency"`
}

// GuildHistoryData modela historial de gremio en tGuildHistory.
type GuildHistoryData struct {
	HistoryType byte   `json:"history_type"`
	Param1      int    `json:"param1"`
	Param2      int    `json:"param2"`
	Text        string `json:"text"`
}

// GuildTournamentData modela torneos de gremio en tGuildTournament.
type GuildTournamentData struct {
	GTNo     int  `json:"gt_no"`
	State    byte `json:"state"`
	MinLevel byte `json:"min_level"`
	MaxLevel byte `json:"max_level"`
}

// ZoneNodeInfo representa un nodo ZoneService activo en tWorldZoneNode.
type ZoneNodeInfo struct {
	WorldNo        uint8     `json:"world_no"`
	ZoneNo         uint16    `json:"zone_no"`
	ZoneIP         string    `json:"zone_ip"`
	ZonePort       uint16    `json:"zone_port"`
	CurrentPlayers int       `json:"current_players"`
	IsOnline       bool      `json:"is_online"`
	LastHeartbeat  time.Time `json:"last_heartbeat"`
}

// Sync ensures bidirectional alignment between flat fields and the canonical tables.
func (c *Character) Sync() {
	if c.Equips == nil {
		c.Equips = make(map[byte]uint16)
	}
	if c.EquipUpgrades == nil {
		c.EquipUpgrades = make(map[byte]byte)
	}

	// If flat fields are set but canonical tables are empty:
	if c.TCharacter.CharNo == 0 && c.ID != 0 {
		c.TCharacter.CharNo = c.ID
		c.TCharacter.UserNo = c.AccountID
		c.TCharacter.CharName = c.Name
		c.TCharacter.Slot = c.Slot
		c.TCharacter.Class = c.Job
		c.TCharacter.Level = c.Level
		c.TCharacter.IsDelete = 0

		genderByte := byte(0)
		if c.Male {
			genderByte = 1
		}
		c.TCharacterShape.CharNo = c.ID
		c.TCharacterShape.Gender = genderByte
		c.TCharacterShape.Hair = c.Hair
		c.TCharacterShape.HairColor = c.HairColor
		c.TCharacterShape.Face = c.Face

		c.TCharacterState.CharNo = c.ID
		c.TCharacterState.Map = c.Map
		c.TCharacterState.XPos = c.XPos
		c.TCharacterState.YPos = c.YPos
		c.TCharacterState.Rotation = 0x55 // Standard AstraServer/NextGen default (85 deg)
		c.TCharacterState.HP = c.CurHP
		c.TCharacterState.SP = c.CurSP
		if c.TCharacterState.HPStones == 0 {
			c.TCharacterState.HPStones = 15
		}
		if c.TCharacterState.SPStones == 0 {
			c.TCharacterState.SPStones = 7
		}

		c.TCharacterKeyMap.CharNo = c.ID
		c.TCharacterShortCut.CharNo = c.ID
		c.TCharacterGameOpt.CharNo = c.ID
		c.TCharacterSingleOpt.CharNo = c.ID
	} else if c.TCharacter.CharNo != 0 {
		// Populate flat fields from canonical tables
		c.ID = c.TCharacter.CharNo
		c.AccountID = c.TCharacter.UserNo
		c.Name = c.TCharacter.CharName
		c.Slot = c.TCharacter.Slot
		c.Job = c.TCharacter.Class
		c.Level = c.TCharacter.Level
		c.Male = (c.TCharacterShape.Gender != 0)
		c.Hair = c.TCharacterShape.Hair
		c.HairColor = c.TCharacterShape.HairColor
		c.Face = c.TCharacterShape.Face
		c.Map = c.TCharacterState.Map
		c.XPos = c.TCharacterState.XPos
		c.YPos = c.TCharacterState.YPos
		c.CurHP = c.TCharacterState.HP
		c.CurSP = c.TCharacterState.SP
	}
}

// GetEquip returns the equipped item ID for a slot, or 0xFFFF if empty.
func (c *Character) GetEquip(slot byte) uint16 {
	if c.Equips != nil {
		if val, exists := c.Equips[slot]; exists {
			return val
		}
	}
	return 0xFFFF
}

// GetEquipUpgrade returns the upgrade nibble (0..15) for an equipped slot, or 0 if none.
func (c *Character) GetEquipUpgrade(slot byte) byte {
	if c.EquipUpgrades != nil {
		if val, exists := c.EquipUpgrades[slot]; exists {
			return val & 0x0F
		}
	}
	return 0
}

// CharacterOptions representa las opciones binarias del cliente persistidas en tCharacterOptions.
type CharacterOptions struct {
	CharNo       int    `json:"nCharNo"`
	ShortCutData []byte `json:"sShortCutData"` // 1024 bytes (p_Char_GetOptShortCutData)
	ShortCutSize []byte `json:"sShortCutSize"` // 24 bytes (p_Char_GetOptShortCutSize)
	Video        []byte `json:"sVideo"`        // 60 bytes
	Sound        []byte `json:"sSound"`        // 1 byte (p_Char_GetOptSound)
	Game         []byte `json:"sGame"`         // 64 bytes (p_Char_GetOptGame)
	WindowsPos   []byte `json:"sWindowsPos"`   // 136 bytes (p_Char_GetOptWindowsPos)
	KeyMapping   []byte `json:"sKeyMapping"`   // 52 bytes (p_Char_GetOptKeyMapping)
}

// ============================================================================
// Interfaces Modulares por Dominio (World00_Character)
// ============================================================================

// CharacterAvatarRepository gobierna ciclo de vida, stats, vitales y opciones de avatares.
type CharacterAvatarRepository interface {
	GetCharactersByAccount(accountID int) ([]*Character, error)
	GetCharacterByName(name string) (*Character, error)
	GetCharacterByID(id int) (*Character, error)
	GetCharacterLoginData(charNo int) (*Character, error)
	CreateCharacter(ch *Character) error
	UpdateCharacter(ch *Character) error
	DeleteCharacter(accountID int, slot byte) error
	ChangeCharacterName(charNo int, newName string) error
	ChangeCharacterClass(charNo int, newClass byte) error
	SetCharacterShape(charNo int, race, class, gender, hair, hairColor, face byte) error
	CharacterLogin(charNo int) error
	CharacterLogout(charNo int, playMin int) error
	SaveCharLocation(charNo int, mapName string, x, y int, rotation byte) error
	SaveCharStats(charNo int, hp, sp, hps, sps int16) error
	SaveCharLevelExp(charNo int, level int16, exp int64) error
	SaveCharMoney(charNo int, money int64) error
	AddCharMoney(charNo int, amount int64) error
	SubCharMoney(charNo int, amount int64) (bool, error)
	SaveCharFreeStats(charNo int, str, con, dex, intel, men, freePoints byte) error
	IncStatPoint(charNo int, statType byte) error
	DecStatPoint(charNo int, statType byte) error
	GetCharacterStates(charNo int) ([]*CharacterStateData, error)
	SetCharacterState(charNo int, stateNo int16, level int16, stateTime int, restTime int, count int16, pd, ps, pk, pc byte) error
	GetCharacterTitles(charNo int) ([]*CharacterTitleData, error)
	SetCharacterTitle(charNo int, titleType int16, title string) error
	SetCharacterTitleProgress(charNo int, titleType byte, elementNo byte, elementValue int64) error
	GetChargedBuffs(charNo int) ([]*ChargedBuffData, error)
	SetChargedBuff(charNo int, key int, id int16, useDate, endDate time.Time) error
	DelChargedBuff(charNo int, key int) error
	GetRestExpExecTime(charNo int) (time.Time, error)
	SetRestExpExecTime(charNo int, t time.Time) error
	GetAutoPick(charNo int) (bool, error)
	SetAutoPick(charNo int, autoPick bool) error
	SaveCharLink(charNo int, linkZone string, linkX, linkY int) error
	AddCharExp(charNo int, exp int64) error
	AddCharFame(charNo int, fame int) error
	SubCharFame(charNo int, fame int) error
	GetCharFame(charNo int) (int, error)
	GetCharacterChangeIDList(userNo int) ([]*CharIDChangeInfo, error)
	ExecuteCharacterChangeID(rowNo int, charNo int, newName string, isInit byte, userNo int) (int, error)
	IsActiveGuildMaster(charNo int) (bool, error)
	GetPrisonInfo(charNo int) (*PrisonInfo, error)
}

// CharacterItemRepository gobierna inventario, equipamiento, almacén, opciones y gemas.
type CharacterItemRepository interface {
	GetCharacterItems(charNo int) ([]*Item, error)
	GetItemsByStorage(charNo int, storageType byte) ([]*Item, error)
	CreateItem(item *Item) error
	CreateItemUseEmptyStorage(itemKey int64, owner int, storageType byte, itemID int, flags int) (int16, error)
	SaveCharacterItem(item *Item) error
	DeleteItem(itemKey int64) error
	SetItemStorage(itemKey int64, storageType byte, storageSlot int16) error
	ExchangeItemStorage(itemKey1 int64, storageType1 byte, storageSlot1 int16, itemKey2 int64, storageType2 byte, storageSlot2 int16) error
	ArrangeInventory(charNo int, itemKeys []int64) error
	UpdateItemFlags(itemKey int64, flags int) error
	GetItemOptions(itemKey int64) (map[int16]int64, error)
	SetItemOption(itemKey int64, optionType int16, optionData int64) error
	GetItemEnchant(itemKey int64) ([]*ItemEnchantData, error)
	SetItemData(itemKey int64, owner int, storageType byte, storage int16, itemID int, flags int, dDate time.Time) error
	MergeItems(ownerFrom, ownerTo int, keyFrom, keyTo int64, optType int16, optData int64) error
	SplitAndMergeItems(keyFrom int64, ownerFrom int, optTypeFrom int16, optDataFrom int64, keyTo int64, ownerTo int, optTypeTo int16, optDataTo int64) error
	UpgradeItem(owner int, itemKey int64, itemDestroy int, upOptType int16, upOptValue int64, upFailCountType int16, upFailCountValue int64) error
	UpgradeItemFull(owner int, itemKey int64, itemDestroy int, upOptType int16, upOptValue int64, upFailCountType int16, upFailCountValue int64, rawKey, rawLeftKey, rawRightKey, rawMiddleKey, giftMoney int64) error
	GetItemActionCoolTime(charNo int) ([]*ItemActionCoolTimeData, error)
	SetItemActionCoolTime(charNo int, actionGroup int16, endTime time.Time) error
	BuyItemAll(itemKey int64, owner int, storageType byte, storage int16, itemID int, flags int, price int64, fame int, coin int64) error
	BuyItemLot(itemKey int64, owner int, optType int16, optData int64, price int64, fame int, coin int64) error
	SellItemAll(itemKey int64, owner int, price int64) error
	SellItemLot(itemKey int64, owner int, optType int16, optData int64, price int64) error
	UserDepositMoney(userNo, charNo int, amount int64) (int64, error)
	UserWithdrawMoney(userNo, charNo int, amount int64) (int64, error)
	UserGetMoney(userNo int) (int64, error)
	EnchantAddGem(charNo int, equipKey, gemKey int64, socketIdx byte, gemID int, restGemCount byte) error
	EnchantAddSocket(charNo int, equipKey, drillKey int64, drillLotType int16) (byte, int64, error)
	EnchantRemoveGem(charNo int, equipKey, removerKey int64, socketIdx, needItem byte, removerLotType int16) (int64, error)
	SetEnchantGemLot(charNo int, equipKey int64, slot, restGemCount byte) error
	PutOnBelongedCancel(charNo int, cancelKey int64, cancelOptType int16, cancelOptData int64, targetKey int64, targetOptType int16, targetOptData int64) error
	SetMinimonAutoItem(charNo, monType, slot, itemID int) error
	FurnitureEndureUpdate(owner int, kitKey int64, kitOptType int16, kitLot int16, furnitureKey int64, endureOptType int16, endureOptData int64) error
	CostumeWeaponEndureUpdate(owner int, kitKey int64, kitOptType int16, kitLot int16, costumeKey int64, endureOptType int16, endureOptData int64) error
	SafeTradeAll(buyer, seller int, moveMoney, itemKey int64, storageType byte, storage int16) error
	SafeTradeLot(buyer int, buyerItemKey int64, buyerStorageType byte, buyerStorage int16, buyerItemID, buyerFlags int, buyerOptType int16, seller int, sellerItemKey int64, sellerOptType int16, moveMoney, splitLot int64) error
	SafeTradeMerge(buyer int, buyerItemKey int64, buyerOptType int16, seller int, sellerItemKey int64, sellerOptType int16, moveMoney, splitLot int64, isDelete byte) error
}

// CharacterProgressionRepository gobierna habilidades, misiones, atajos, opciones de cliente, asistencia y tutorial.
type CharacterProgressionRepository interface {
	GetCharacterSkills(charNo int) ([]*CharacterSkill, error)
	SaveCharacterSkill(charNo int, skill *CharacterSkill) error
	SaveCharacterSkills(charNo int, skills []*CharacterSkill) error
	SetSkillCoolTime(charNo int, skillNo int16, writeTime int64, coolTime int) error
	SetSkillPower(charNo int, skillNo int16, pd, ps, keepTime, coolTime byte) error
	SetSingleSkillPower(charNo int, skillNo int16, powerType byte, power byte) error
	DeleteCharacterSkill(charNo int, skillNo int16) error
	GetCharacterShortcuts(charNo int) ([]*Shortcut, error)
	SaveCharacterShortcuts(charNo int, shortcuts []*Shortcut) error
	GetCharacterOptions(charNo int) (*CharacterOptions, error)
	SaveCharacterOption(charNo int, optType string, data []byte) error
	DeleteCharacterOptions(charNo int) error
	GetImproveShortCutList(charNo int) ([]*ImproveShortcutEntry, error)
	SetImproveShortCut(charNo int, slotNo byte, codeNo uint16, value int32) error
	GetImproveKeyMap(charNo int) ([]*ImproveKeyMapEntry, error)
	SetImproveKeyMap(charNo int, funcNo uint16, extendKey, asciiCode byte) error
	InitImproveKeyMap(charNo int) error
	GetImproveGameOptions(charNo int) ([]*ImproveGameOptionEntry, error)
	SetImproveGameOption(charNo int, optionNo uint16, value byte) error
	InitImproveGameOptions(charNo int) error
	GetInGameGuideViewList(charNo int) ([]byte, error)
	SetInGameGuideView(charNo int, guideNo byte) error
	GetCharacterQuests(charNo int) ([]*QuestData, error)
	GetCharacterDoneQuests(charNo int) ([]int, error)
	AddCharacterQuest(charNo int, questNo int, status byte, data []byte) error
	SetCharacterQuest(charNo int, questNo int, status byte, data []byte) error
	DeleteCharacterQuest(charNo int, questNo int) error
	GetAttendance(charNo int) ([]time.Time, error)
	CheckAttendance(charNo int, dDate time.Time) error
	GetAttendanceRewards(charNo int) ([]int, error)
	SetAttendanceReward(charNo int, rewardID int) error
	SetAttendanceRewardAt(charNo int, eventDate time.Time, rewardID int) (int, error)
	GetCollectionCards(charNo int) ([]*CollectionCardData, error)
	SetCollectionCard(charNo int, cardID int16, count int16) error
	RegisterCollectionCard(charNo int, cardID int, serialNo int64, starNumber byte, itemKey int64) error
	GetCollectionFavorites(charNo int) ([]*CollectionBookmarkData, error)
	SetCollectionFavorite(charNo int, slot byte, cardID int, regist bool) error
	GetCollectionRewards(charNo int) ([]int, error)
	SetCollectionReward(charNo int, rewardID int) error
	GetCollectionSerialList() (map[int]int64, error)
	SetCollectionSerialCount(cardID int, count int64) error
	GetTutorialCharInfo(charNo int) (byte, byte, error)
	SetTutorialCharInfo(userNo, charNo int, state, step byte) error
	GetTutorialIsSkip(userNo int) (bool, error)
}

// GuildRepository gobierna creación, membresías, rangos, academia y torneos de gremios.
type GuildRepository interface {
	CreateGuild(masterCharID int, name, password string, guildType byte) (int, error)
	GetGuild(guildID int) (*GuildData, error)
	GetGuildByCharID(charID int) (*GuildData, *GuildMemberData, error)
	GetGuildMembers(guildID int) ([]*GuildMemberData, error)
	GetAllGuilds() ([]*GuildListItem, error)
	AddGuildMember(guildID, charID int, grade byte) error
	RemoveGuildMember(guildID, charID int) error
	SetGuildMemberGrade(guildID, charID int, grade byte) error
	SetGuildNotify(guildID int, charName, notify string) error
	DismissGuild(guildID int) error
	SetGuildMoney(guildID int, money uint64) error
	SetGuildGrade(guildID int, grade byte) error
	GetGuildToken(guildID int) (int, error)
	SetGuildToken(guildID int, token int) error
	GetGuildHistory(guildID int) ([]*GuildHistoryData, error)
	AddGuildHistory(guildID int, historyType byte, param1, param2 int, text string) error
	CreateGuildAcademy(guildID, masterCharID int) error
	GetGuildAcademy(guildID int) (*GuildAcademyData, error)
	GetAcademyByCharID(charID int) (*GuildAcademyData, *GuildAcademyMemberData, error)
	GetAcademyMembers(guildID int) ([]*GuildAcademyMemberData, error)
	AddAcademyMember(guildID, charID int, joinLevel byte) error
	RemoveAcademyMember(guildID, charID int) error
	SetAcademyNotify(guildID int, charName, notify string) error
	SetAcademyChatBan(guildID, charID int, banMinutes int) error
	SetAcademyGraduate(guildID, charID int) error
	GetAcademyGraduates(guildID int) ([]int, error)
	SubGuildMoney(guildID, memberCharNo int, amount int64) (int64, int64, error)
	GetAcademyMoney(guildID int) (int64, error)
	DepositAcademyMoney(guildID, charNo int, amount int64) error
	WithdrawAcademyMoney(guildID, charNo int, amount int64) error
	GetAcademyRewardsAll(guildID, typeID int) ([]*AcademyRewardData, error)
	DelAcademyReward(guildID int, level int16, typeID byte) error
	GetGuildTournament(gtNo int) (*GuildTournamentData, error)
	GetAllGuildTournaments() ([]*GuildTournamentData, error)
}

// WorldSocialRepository gobierna lista de amigos, friendpoints, Kingdom Quests, HolyPromise, mascotas y minihouse.
type WorldSocialRepository interface {
	GetFriends(charNo int) ([]*FriendInfo, error)
	AddFriend(charNo, friendCharNo int) error
	RemoveFriend(charNo, friendCharNo int) error
	GetFriendPoints(charNo int) (int, error)
	SetFriendPoints(charNo int, points int) error
	GetKingdomQuests() ([]*KingdomQuestData, error)
	GetKingdomQuest(questNo int) (*KingdomQuestData, error)
	SetKingdomQuest(questNo int, data string) error
	DelKingdomQuest(questNo int) error
	CheckHolyPromiseRing(charNo, upCharNo int) (bool, error)
	SetHolyPromise(charNo, upCharNo int, flags byte) error
	DelHolyPromise(charNo, upCharNo int) error
	DelHolyPromiseChar(charNo int) error
	GetHolyPromise(charNo int) (*HolyPromiseData, error)
	GetMasterPupils(charNo int) ([]*MasterPupilData, error)
	SetMaster(charNo int, masterCharNo int) error
	SetPupil(charNo int, pupilCharNo int) error
	DelMaster(charNo int) error
	DelPupils(charNo int) error
	DelAllMasterPupils(charNo int) error
	GetPetInfo(ownerNo int) (*PetData, error)
	SetPetName(petNo int, ownerNo int, name string) error
	SetPetTendency(petNo int, tendency byte) error
	GetMiniHouseURL(charNo int) (string, error)
	SetMiniHouseURL(charNo int, url string) error
	GetMiniHouseVisitors(charNo int) (int, error)
	SetMiniHouseVisitors(charNo int, visitors int) error

	// Party (tParty / tPartyMember)
	CreateParty(masterCharNo int, flags int) (int, error)
	AddPartyMember(partyNo int, charNo int, flag byte) error
	RemovePartyMember(partyNo int, charNo int) error
	SetPartyMaster(partyNo int, masterCharNo int) error
	SetPartyFlags(partyNo int, flags int) error
	DeleteParty(partyNo int) error
	GetParty(partyNo int) (*PartyData, error)
	GetPartyMembers(partyNo int) ([]*PartyMemberData, error)
	GetPartyByChar(charNo int) (int, int, int, error)

	// Chat Restrict (tChatBlock)
	AddChatRestrict(userNo int, charNo int) error
	DelChatRestrict(userNo int, charNo int) error
	DelAllChatRestrict(userNo int) error

	// Prison (tPrison)
	AddPrison(gmUserID string, targetCharName string, minutes int16, reason, remark string) error
	EndPrison(charNo int) error
	AlterPrison(gmUserID string, targetCharName string, minutes int16, reason, remark string) error

	// Spammer / LoadingBug (tCharacterLoading / tSpamBlockTime / tSpamFilter)
	SetClientLoadingBug(charNo int, mapName string, failType int16, failCnt int16) error
	SetSpamFilter(userNo int, isSpamSpread, isSpamBlock byte, rate byte) error
	GetSpamFilter(userNo int) (bool, bool, byte, error)
	SetSpamBlockTime(charNo int, blockTime, standbyTime int) error
	GetSpamBlockTime(charNo int) (int, int, error)
}

// ZoneTopologyRepository gobierna transferencias de un solo uso y topología de nodos ZoneService.
type ZoneTopologyRepository interface {
	IssueZoneTransfer(worldNo uint8, zoneNo uint16, userNo int, userID string, charNo int, charID string, randomID uint16, adminLevel uint8, clientIP, targetMap string, targetX, targetY int) error
	ConsumeZoneTransfer(charID string, randomID uint16, clientIP string) (*ZoneTransfer, error)
	UpsertZoneNode(worldNo uint8, zoneNo uint16, zoneIP string, zonePort uint16) error
	HeartbeatZoneNode(worldNo uint8, zoneNo uint16, currentPlayers int) error
	GetZoneNodes(worldNo uint8) ([]*ZoneNodeInfo, error)
	ResolveMapEndpoint(worldNo uint8, zoneNo uint16) (string, uint16, bool, error)
}

// CharacterRepository unifica la totalidad de los repositorios de WorldService y ZoneService.
type CharacterRepository interface {
	CharacterAvatarRepository
	CharacterItemRepository
	CharacterProgressionRepository
	GuildRepository
	WorldSocialRepository
	ZoneTopologyRepository
}

// GuildData modela tGuild según el esquema oficial de base de datos.
type GuildData struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Password      string    `json:"password"`
	MaxMembers    uint16    `json:"max_members"`
	Money         uint64    `json:"money"`
	Type          byte      `json:"type"`
	Grade         byte      `json:"grade"`
	Fame          uint32    `json:"fame"`
	StoneLevel    uint16    `json:"stone_level"`
	Exp           uint64    `json:"exp"`
	CreateDate    time.Time `json:"create_date"`
	DismissStatus byte      `json:"dismiss_status"`
	DismissDate   time.Time `json:"dismiss_date"`
	Intro         string    `json:"intro"`
	NotifyDate    time.Time `json:"notify_date"`
	NotifyCharID  string    `json:"notify_char_id"`
	Notify        string    `json:"notify"`
}

// GuildMemberData modela tGuildMember y sus datos de personaje asociados.
type GuildMemberData struct {
	GuildID   int       `json:"guild_id"`
	CharID    int       `json:"char_id"`
	CharName  string    `json:"char_name"`
	Class     byte      `json:"class"`
	Level     byte      `json:"level"`
	Map       string    `json:"map"`
	Grade     byte      `json:"grade"`
	Flags     int       `json:"flags"`
	ExpRatio  int16     `json:"exp_ratio"`
	Intro     string    `json:"intro"`
	JoinDate  time.Time `json:"join_date"`
	LastLogin time.Time `json:"last_login"`
}

// GuildListItem modela las filas devueltas por p_Guild_GetListDetail.
type GuildListItem struct {
	ID             int    `json:"id"`
	GuildID        uint32 `json:"guild_id"`
	Name           string `json:"name"`
	GuildName      string `json:"guild_name"`
	MasterCharName string `json:"master_char_name"`
	MemberCount    uint16 `json:"member_count"`
	MaxMembers     uint16 `json:"max_members"`
	Type           byte   `json:"type"`
	Grade          byte   `json:"grade"`
	Fame           uint32 `json:"fame"`
	Exp            uint64 `json:"exp"`
}

// GuildAcademyData modela tGuildAcademy.
type GuildAcademyData struct {
	GuildID        int       `json:"guild_id"`
	MasterChar     int       `json:"master_char"`
	MasterCharName string    `json:"master_char_name"`
	MaxMember      uint16    `json:"max_member"`
	Point          uint32    `json:"point"`
	Intro          string    `json:"intro"`
	NotifyCharID   string    `json:"notify_char_id"`
	Notify         string    `json:"notify"`
	NotifyDate     time.Time `json:"notify_date"`
	BuffEndTime    time.Time `json:"buff_end_time"`
}

// GuildAcademyMemberData modela tGuildAcademyMember.
type GuildAcademyMemberData struct {
	GuildID    int       `json:"guild_id"`
	CharID     int       `json:"char_id"`
	CharName   string    `json:"char_name"`
	Class      byte      `json:"class"`
	Level      byte      `json:"level"`
	JoinLevel  byte      `json:"join_level"`
	JoinDate   time.Time `json:"join_date"`
	Grade      byte      `json:"grade"`
	Map        string    `json:"map"`
	LastLogin  time.Time `json:"last_login"`
	Graduate   bool      `json:"graduate"`
	ChatBanEnd time.Time `json:"chat_ban_end"`
}

// AcademyRewardData modela tGuildAcademyReward.
type AcademyRewardData struct {
	Level  int16 `json:"level"`
	TypeID byte  `json:"type_id"`
	Value  int64 `json:"value"`
	ItemID int   `json:"item_id"`
}

// FriendInfo modela relaciones de amistad en tFriend.
type FriendInfo struct {
	CharNo    int       `json:"char_no"`
	CharName  string    `json:"char_name"`
	Class     byte      `json:"class"`
	Level     byte      `json:"level"`
	Map       string    `json:"map"`
	IsOnline  bool      `json:"is_online"`
	CreatedAt time.Time `json:"created_at"`
}

// ============================================================================
// MemoryCharacterRepository (Implementación mock para tests autónomos)
// ============================================================================

type MemoryCharacterRepository struct {
	mu           sync.RWMutex
	filePath     string
	characters   map[int]*Character
	charNameToID map[string]int
	nextCharID   int
	nextItemKey  int64

	items         map[int][]*Item
	shortcuts     map[int][]*Shortcut
	improveSC     map[int][]*ImproveShortcutEntry
	improveKeyMap map[int][]*ImproveKeyMapEntry
	improveGameOp map[int][]*ImproveGameOptionEntry
	guides        map[int][]byte
	changeIDs     map[int][]*CharIDChangeInfo
	attendRewards map[int]map[int]bool
	tutorialSkip  map[int]bool
	skills        map[int][]*CharacterSkill
	options       map[int]*CharacterOptions
	friends       map[int][]*FriendInfo
	friendPts     map[int]int
	guilds        map[int]*GuildData
	guildMbrs     map[int][]*GuildMemberData
	academies     map[int]*GuildAcademyData
	academyMbrs   map[int][]*GuildAcademyMemberData
	zoneNodes     map[uint16]*ZoneNodeInfo

	titles              map[int][]*CharacterTitleData
	autoPick            map[int]bool
	restExpTime         map[int]time.Time
	collectionCards     map[int][]*CollectionCardData
	collectionBookmarks map[int][]*CollectionBookmarkData
	collectionRewards   map[int]map[int]bool
	collectionSerials   map[int]int64
	parties             map[int]*PartyData
	partyMembers        map[int][]*PartyMemberData
	nextPartyNo         int
	itemOptions         map[int64]map[int16]int64
	itemEnchants        map[int64][]*ItemEnchantData
	userMoney           map[int]int64
	academyMoney        map[int]int64
	academyRewards      map[int][]*AcademyRewardData
}

func NewMemoryCharacterRepositoryWithFile(filePath string) *MemoryCharacterRepository {
	repo := &MemoryCharacterRepository{
		filePath:            filePath,
		characters:          make(map[int]*Character),
		charNameToID:        make(map[string]int),
		nextCharID:          1,
		nextItemKey:         1,
		items:               make(map[int][]*Item),
		shortcuts:           make(map[int][]*Shortcut),
		improveSC:           make(map[int][]*ImproveShortcutEntry),
		improveKeyMap:       make(map[int][]*ImproveKeyMapEntry),
		improveGameOp:       make(map[int][]*ImproveGameOptionEntry),
		guides:              make(map[int][]byte),
		changeIDs:           make(map[int][]*CharIDChangeInfo),
		attendRewards:       make(map[int]map[int]bool),
		tutorialSkip:        make(map[int]bool),
		skills:              make(map[int][]*CharacterSkill),
		options:             make(map[int]*CharacterOptions),
		friends:             make(map[int][]*FriendInfo),
		friendPts:           make(map[int]int),
		titles:              make(map[int][]*CharacterTitleData),
		autoPick:            make(map[int]bool),
		restExpTime:         make(map[int]time.Time),
		collectionCards:     make(map[int][]*CollectionCardData),
		collectionBookmarks: make(map[int][]*CollectionBookmarkData),
		collectionRewards:   make(map[int]map[int]bool),
		collectionSerials:   make(map[int]int64),
		parties:             make(map[int]*PartyData),
		partyMembers:        make(map[int][]*PartyMemberData),
		nextPartyNo:         1,
		itemOptions:         make(map[int64]map[int16]int64),
		itemEnchants:        make(map[int64][]*ItemEnchantData),
		userMoney:           make(map[int]int64),
		academyMoney:        make(map[int]int64),
		academyRewards:      make(map[int][]*AcademyRewardData),
		guilds:        make(map[int]*GuildData),
		guildMbrs:     make(map[int][]*GuildMemberData),
		academies:     make(map[int]*GuildAcademyData),
		academyMbrs:   make(map[int][]*GuildAcademyMemberData),
		zoneNodes:     make(map[uint16]*ZoneNodeInfo),
	}
	repo.load()
	return repo
}

func NewMemoryCharacterRepository() *MemoryCharacterRepository {
	return NewMemoryCharacterRepositoryWithFile(filepath.Join("data", "characters.json"))
}

func (r *MemoryCharacterRepository) load() {
	if r.filePath == "" {
		return
	}
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return
	}
	var stored []*Character
	if err := json.Unmarshal(data, &stored); err != nil {
		return
	}
	for _, ch := range stored {
		ch.Sync()
		r.characters[ch.ID] = ch
		r.charNameToID[strings.ToLower(ch.Name)] = ch.ID
		if ch.ID >= r.nextCharID {
			r.nextCharID = ch.ID + 1
		}
	}
}

func (r *MemoryCharacterRepository) save() {
	if r.filePath == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(r.filePath), 0755)
	stored := make([]*Character, 0, len(r.characters))
	for _, ch := range r.characters {
		ch.Sync()
		stored = append(stored, ch)
	}
	data, err := json.MarshalIndent(stored, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(r.filePath, data, 0644)
}

func (r *MemoryCharacterRepository) GetCharactersByAccount(accountID int) ([]*Character, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Character
	for _, ch := range r.characters {
		if ch.AccountID == accountID {
			result = append(result, ch)
		}
	}
	return result, nil
}

func (r *MemoryCharacterRepository) GetCharacterByName(name string) (*Character, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, exists := r.charNameToID[strings.ToLower(name)]
	if !exists {
		return nil, ErrCharacterNotFound
	}
	return r.characters[id], nil
}

func (r *MemoryCharacterRepository) GetCharacterByID(id int) (*Character, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, exists := r.characters[id]
	if !exists {
		return nil, ErrCharacterNotFound
	}
	return ch, nil
}

func (r *MemoryCharacterRepository) GetCharacterLoginData(charNo int) (*Character, error) {
	return r.GetCharacterByID(charNo)
}

func (r *MemoryCharacterRepository) CreateCharacter(ch *Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	lower := strings.ToLower(ch.Name)
	if _, exists := r.charNameToID[lower]; exists {
		return ErrCharacterExists
	}
	for _, existing := range r.characters {
		if existing.AccountID == ch.AccountID && existing.Slot == ch.Slot {
			return ErrSlotInUse
		}
	}

	ch.ID = r.nextCharID
	r.nextCharID++
	ch.Sync()

	r.characters[ch.ID] = ch
	r.charNameToID[lower] = ch.ID
	for slot, itemID := range ch.Equips {
		r.items[ch.ID] = append(r.items[ch.ID], &Item{
			ItemKey:     r.nextItemKey,
			Owner:       ch.ID,
			StorageType: 8,
			Storage:     int16(slot),
			ItemID:      itemID,
			Flags:       0,
			Date:        time.Now(),
		})
		r.nextItemKey++
	}
	for _, it := range ch.Items {
		it.Owner = ch.ID
		if it.ItemKey == 0 {
			it.ItemKey = r.nextItemKey
			r.nextItemKey++
		}
		r.items[ch.ID] = append(r.items[ch.ID], it)
	}
	if len(ch.Shortcuts) > 0 {
		r.shortcuts[ch.ID] = append(r.shortcuts[ch.ID], ch.Shortcuts...)
	}
	if len(ch.ImproveShortcuts) > 0 {
		r.improveSC[ch.ID] = append(r.improveSC[ch.ID], ch.ImproveShortcuts...)
	}
	if len(ch.Skills) > 0 {
		r.skills[ch.ID] = append(r.skills[ch.ID], ch.Skills...)
	}
	r.save()
	return nil
}

func (r *MemoryCharacterRepository) UpdateCharacter(ch *Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.characters[ch.ID]; !exists {
		return ErrCharacterNotFound
	}
	ch.Sync()
	r.characters[ch.ID] = ch
	r.save()
	return nil
}

func (r *MemoryCharacterRepository) DeleteCharacter(accountID int, slot byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, ch := range r.characters {
		if ch.AccountID == accountID && ch.Slot == slot {
			for _, mbrs := range r.guildMbrs {
				for _, m := range mbrs {
					if m.CharID == id && m.Grade == 0 {
						return ErrCharacterIsGuildMaster
					}
				}
			}
			delete(r.charNameToID, strings.ToLower(ch.Name))
			delete(r.characters, id)
			r.save()
			return nil
		}
	}
	return ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) ChangeCharacterName(charNo int, newName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.characters[charNo]
	if !ok {
		return ErrCharacterNotFound
	}
	lowerNew := strings.ToLower(newName)
	if existingID, exists := r.charNameToID[lowerNew]; exists && existingID != charNo {
		return ErrCharacterExists
	}
	delete(r.charNameToID, strings.ToLower(ch.Name))
	ch.Name = newName
	ch.TCharacter.CharName = newName
	r.charNameToID[lowerNew] = charNo
	return nil
}

func (r *MemoryCharacterRepository) ChangeCharacterClass(charNo int, newClass byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.characters[charNo]
	if !ok {
		return ErrCharacterNotFound
	}
	ch.Job = newClass
	return nil
}

func (r *MemoryCharacterRepository) SetCharacterShape(charNo int, race, class, gender, hair, hairColor, face byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.characters[charNo]
	if !ok {
		return ErrCharacterNotFound
	}
	ch.Job = class
	ch.Male = (gender != 0)
	ch.Hair = hair
	ch.HairColor = hairColor
	ch.Face = face
	return nil
}

func (r *MemoryCharacterRepository) CharacterLogin(charNo int) error               { return nil }
func (r *MemoryCharacterRepository) CharacterLogout(charNo int, playMin int) error { return nil }

func (r *MemoryCharacterRepository) SaveCharLocation(charNo int, mapName string, x, y int, rotation byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		ch.TownName = mapName
		ch.XPos = x
		ch.YPos = y
		return nil
	}
	return ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) SaveCharStats(charNo int, hp, sp, hps, sps int16) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		ch.CurHP = int(hp)
		ch.CurSP = int(sp)
		return nil
	}
	return ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) SaveCharLevelExp(charNo int, level int16, exp int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		ch.Level = byte(level)
		ch.TCharacterState.Exp = exp
		return nil
	}
	return ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) SaveCharMoney(charNo int, money int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		ch.TCharacterState.Money = money
		return nil
	}
	return ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) AddCharMoney(charNo int, amount int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		ch.TCharacterState.Money += amount
		return nil
	}
	return ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) SubCharMoney(charNo int, amount int64) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		if ch.TCharacterState.Money < amount {
			return false, nil
		}
		ch.TCharacterState.Money -= amount
		return true, nil
	}
	return false, ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) SaveCharFreeStats(charNo int, str, con, dex, intel, men, freePoints byte) error {
	return nil
}

func (r *MemoryCharacterRepository) IncStatPoint(charNo int, statType byte) error {
	return nil
}

func (r *MemoryCharacterRepository) DecStatPoint(charNo int, statType byte) error {
	return nil
}

func (r *MemoryCharacterRepository) GetCharacterStates(charNo int) ([]*CharacterStateData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetCharacterState(charNo int, stateNo int16, level int16, stateTime int, restTime int, count int16, pd, ps, pk, pc byte) error {
	return nil
}

func (r *MemoryCharacterRepository) GetCharacterTitles(charNo int) ([]*CharacterTitleData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	src := r.titles[charNo]
	out := make([]*CharacterTitleData, 0, len(src))
	for _, t := range src {
		cp := *t
		out = append(out, &cp)
	}
	return out, nil
}

func (r *MemoryCharacterRepository) SetCharacterTitle(charNo int, titleType int16, title string) error {
	return r.SetCharacterTitleProgress(charNo, byte(titleType), 0, 1)
}

func (r *MemoryCharacterRepository) SetCharacterTitleProgress(charNo int, titleType byte, elementNo byte, elementValue int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.titles == nil {
		r.titles = make(map[int][]*CharacterTitleData)
	}
	for _, t := range r.titles[charNo] {
		if byte(t.TitleType) == titleType && t.ElementNo == elementNo {
			t.ElementValue = elementValue
			return nil
		}
	}
	r.titles[charNo] = append(r.titles[charNo], &CharacterTitleData{
		CharNo:       charNo,
		TitleType:    int16(titleType),
		ElementNo:    elementNo,
		ElementValue: elementValue,
	})
	return nil
}

func (r *MemoryCharacterRepository) GetChargedBuffs(charNo int) ([]*ChargedBuffData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetChargedBuff(charNo int, key int, id int16, useDate, endDate time.Time) error {
	return nil
}

func (r *MemoryCharacterRepository) DelChargedBuff(charNo int, key int) error {
	return nil
}

func (r *MemoryCharacterRepository) GetRestExpExecTime(charNo int) (time.Time, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if t, ok := r.restExpTime[charNo]; ok {
		return t, nil
	}
	return time.Now(), nil
}

func (r *MemoryCharacterRepository) SetRestExpExecTime(charNo int, t time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.restExpTime == nil {
		r.restExpTime = make(map[int]time.Time)
	}
	r.restExpTime[charNo] = t
	return nil
}

func (r *MemoryCharacterRepository) GetAutoPick(charNo int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.autoPick[charNo], nil
}

func (r *MemoryCharacterRepository) SetAutoPick(charNo int, autoPick bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.autoPick == nil {
		r.autoPick = make(map[int]bool)
	}
	r.autoPick[charNo] = autoPick
	return nil
}

func (r *MemoryCharacterRepository) SaveCharLink(charNo int, linkZone string, linkX, linkY int) error {
	return nil
}

func (r *MemoryCharacterRepository) AddCharExp(charNo int, exp int64) error {
	return nil
}

func (r *MemoryCharacterRepository) AddCharFame(charNo int, fame int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		ch.TCharacterState.Fame += fame
		return nil
	}
	return ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) SubCharFame(charNo int, fame int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		ch.TCharacterState.Fame -= fame
		if ch.TCharacterState.Fame < 0 {
			ch.TCharacterState.Fame = 0
		}
		return nil
	}
	return ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) GetCharFame(charNo int) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if ch, ok := r.characters[charNo]; ok {
		return ch.TCharacterState.Fame, nil
	}
	return 0, ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) GetCharacterChangeIDList(userNo int) ([]*CharIDChangeInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.changeIDs[userNo], nil
}

func (r *MemoryCharacterRepository) ExecuteCharacterChangeID(rowNo int, charNo int, newName string, isInit byte, userNo int) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, ok := r.characters[charNo]
	if !ok {
		return -3023, nil
	}
	if ch.AccountID != userNo {
		return -3024, nil
	}
	lowerNew := strings.ToLower(newName)
	if _, exists := r.charNameToID[lowerNew]; exists {
		return -3039, nil
	}
	list := r.changeIDs[userNo]
	found := -1
	for i, entry := range list {
		if entry.RowNo == uint32(rowNo) && entry.CharNo == charNo {
			found = i
			break
		}
	}
	if found == -1 {
		return -3024, nil
	}
	delete(r.charNameToID, strings.ToLower(ch.Name))
	ch.Name = newName
	ch.TCharacter.CharName = newName
	ch.ChangeID = nil
	r.charNameToID[lowerNew] = charNo
	r.changeIDs[userNo] = append(list[:found], list[found+1:]...)
	if isInit != 0 {
		delete(r.friends, charNo)
	}
	return 0, nil
}

func (r *MemoryCharacterRepository) IsActiveGuildMaster(charNo int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, mbrs := range r.guildMbrs {
		for _, m := range mbrs {
			if m.CharID == charNo && m.Grade == 0 {
				return true, nil
			}
		}
	}
	return false, nil
}

func (r *MemoryCharacterRepository) GetPrisonInfo(charNo int) (*PrisonInfo, error) {
	return nil, nil
}

// Items
func (r *MemoryCharacterRepository) GetCharacterItems(charNo int) ([]*Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items[charNo], nil
}

func (r *MemoryCharacterRepository) GetItemsByStorage(charNo int, storageType byte) ([]*Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*Item
	for _, it := range r.items[charNo] {
		if it.StorageType == storageType {
			res = append(res, it)
		}
	}
	return res, nil
}

func (r *MemoryCharacterRepository) CreateItem(item *Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if item.ItemKey == 0 {
		item.ItemKey = r.nextItemKey
		r.nextItemKey++
	}
	r.items[item.Owner] = append(r.items[item.Owner], item)
	return nil
}

func (r *MemoryCharacterRepository) CreateItemUseEmptyStorage(itemKey int64, owner int, storageType byte, itemID int, flags int) (int16, error) {
	it := &Item{
		ItemKey:     itemKey,
		Owner:       owner,
		StorageType: storageType,
		Storage:     0,
		ItemID:      uint16(itemID),
		Flags:       int32(flags),
		Date:        time.Now(),
	}
	_ = r.CreateItem(it)
	return 0, nil
}

func (r *MemoryCharacterRepository) SaveCharacterItem(item *Item) error {
	return r.CreateItem(item)
}

func (r *MemoryCharacterRepository) DeleteItem(itemKey int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for owner, list := range r.items {
		for i, it := range list {
			if it.ItemKey == itemKey {
				r.items[owner] = append(list[:i], list[i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) SetItemStorage(itemKey int64, storageType byte, storageSlot int16) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, list := range r.items {
		for _, it := range list {
			if it.ItemKey == itemKey {
				it.StorageType = storageType
				it.Storage = storageSlot
				return nil
			}
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) ExchangeItemStorage(itemKey1 int64, storageType1 byte, storageSlot1 int16, itemKey2 int64, storageType2 byte, storageSlot2 int16) error {
	_ = r.SetItemStorage(itemKey1, storageType1, storageSlot1)
	_ = r.SetItemStorage(itemKey2, storageType2, storageSlot2)
	return nil
}

func (r *MemoryCharacterRepository) ArrangeInventory(charNo int, itemKeys []int64) error {
	return nil
}

func (r *MemoryCharacterRepository) UpdateItemFlags(itemKey int64, flags int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, list := range r.items {
		for _, it := range list {
			if it.ItemKey == itemKey {
				it.Flags = int32(flags)
				return nil
			}
		}
	}
	return errors.New("item not found")
}

func (r *MemoryCharacterRepository) GetItemOptions(itemKey int64) (map[int16]int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	opts, ok := r.itemOptions[itemKey]
	if !ok {
		return nil, nil
	}
	out := make(map[int16]int64, len(opts))
	for k, v := range opts {
		out[k] = v
	}
	return out, nil
}

func (r *MemoryCharacterRepository) SetItemOption(itemKey int64, optionType int16, optionData int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.itemOptions == nil {
		r.itemOptions = make(map[int64]map[int16]int64)
	}
	if r.itemOptions[itemKey] == nil {
		r.itemOptions[itemKey] = make(map[int16]int64)
	}
	r.itemOptions[itemKey][optionType] = optionData
	return nil
}

func (r *MemoryCharacterRepository) GetItemEnchant(itemKey int64) ([]*ItemEnchantData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.itemEnchants[itemKey]
	out := make([]*ItemEnchantData, len(list))
	for i, e := range list {
		cp := *e
		out[i] = &cp
	}
	return out, nil
}

func (r *MemoryCharacterRepository) SetItemData(itemKey int64, owner int, storageType byte, storage int16, itemID int, flags int, dDate time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for oldOwner, list := range r.items {
		for i, it := range list {
			if it.ItemKey == itemKey {
				it.Owner = owner
				it.StorageType = storageType
				it.Storage = storage
				it.ItemID = uint16(itemID)
				it.Flags = int32(flags)
				it.Date = dDate
				if oldOwner != owner {
					r.items[oldOwner] = append(list[:i], list[i+1:]...)
					r.items[owner] = append(r.items[owner], it)
				}
				return nil
			}
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) MergeItems(ownerFrom, ownerTo int, keyFrom, keyTo int64, optType int16, optData int64) error {
	_ = r.UpdateItemFlags(keyTo, int(optData))
	_ = r.DeleteItem(keyFrom)
	return nil
}

func (r *MemoryCharacterRepository) SplitAndMergeItems(keyFrom int64, ownerFrom int, optTypeFrom int16, optDataFrom int64, keyTo int64, ownerTo int, optTypeTo int16, optDataTo int64) error {
	_ = r.UpdateItemFlags(keyFrom, int(optDataFrom))
	_ = r.UpdateItemFlags(keyTo, int(optDataTo))
	return nil
}

func (r *MemoryCharacterRepository) UpgradeItem(owner int, itemKey int64, itemDestroy int, upOptType int16, upOptValue int64, upFailCountType int16, upFailCountValue int64) error {
	return r.UpgradeItemFull(owner, itemKey, itemDestroy, upOptType, upOptValue, upFailCountType, upFailCountValue, 0, 0, 0, 0, 0)
}

func (r *MemoryCharacterRepository) UpgradeItemFull(owner int, itemKey int64, itemDestroy int, upOptType int16, upOptValue int64, upFailCountType int16, upFailCountValue int64, rawKey, rawLeftKey, rawRightKey, rawMiddleKey, giftMoney int64) error {
	if giftMoney > 0 {
		ok, err := r.SubCharMoney(owner, giftMoney)
		if err != nil || !ok {
			return errors.New("insufficient money")
		}
	}
	if itemDestroy == 1 {
		_ = r.DeleteItem(itemKey)
	} else {
		_ = r.SetItemOption(itemKey, upOptType, upOptValue)
		_ = r.SetItemOption(itemKey, upFailCountType, upFailCountValue)
	}
	for _, rk := range []int64{rawKey, rawLeftKey, rawRightKey, rawMiddleKey} {
		if rk > 0 {
			_ = r.DeleteItem(rk)
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) GetItemActionCoolTime(charNo int) ([]*ItemActionCoolTimeData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetItemActionCoolTime(charNo int, actionGroup int16, endTime time.Time) error {
	return nil
}

func (r *MemoryCharacterRepository) BuyItemAll(itemKey int64, owner int, storageType byte, storage int16, itemID int, flags int, price int64, fame int, coin int64) error {
	if price > 0 {
		ok, err := r.SubCharMoney(owner, price)
		if err != nil || !ok {
			return errors.New("insufficient money")
		}
	}
	if fame > 0 {
		if err := r.SubCharFame(owner, fame); err != nil {
			return err
		}
	}
	return r.CreateItem(&Item{
		ItemKey:     itemKey,
		Owner:       owner,
		StorageType: storageType,
		Storage:     storage,
		ItemID:      uint16(itemID),
		Flags:       int32(flags),
		Date:        time.Now(),
	})
}

func (r *MemoryCharacterRepository) BuyItemLot(itemKey int64, owner int, optType int16, optData int64, price int64, fame int, coin int64) error {
	if price > 0 {
		ok, err := r.SubCharMoney(owner, price)
		if err != nil || !ok {
			return errors.New("insufficient money")
		}
	}
	if fame > 0 {
		if err := r.SubCharFame(owner, fame); err != nil {
			return err
		}
	}
	_ = r.UpdateItemFlags(itemKey, int(optData))
	return r.SetItemOption(itemKey, optType, optData)
}

func (r *MemoryCharacterRepository) SellItemAll(itemKey int64, owner int, price int64) error {
	if price > 0 {
		_ = r.AddCharMoney(owner, price)
	}
	return r.DeleteItem(itemKey)
}

func (r *MemoryCharacterRepository) SellItemLot(itemKey int64, owner int, optType int16, optData int64, price int64) error {
	if price > 0 {
		_ = r.AddCharMoney(owner, price)
	}
	_ = r.UpdateItemFlags(itemKey, int(optData))
	return r.SetItemOption(itemKey, optType, optData)
}

func (r *MemoryCharacterRepository) UserDepositMoney(userNo, charNo int, amount int64) (int64, error) {
	if amount <= 0 {
		return 0, errors.New("invalid deposit amount")
	}
	ok, err := r.SubCharMoney(charNo, amount)
	if err != nil || !ok {
		return 0, errors.New("insufficient character money")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.userMoney == nil {
		r.userMoney = make(map[int]int64)
	}
	r.userMoney[userNo] += amount
	return r.userMoney[userNo], nil
}

func (r *MemoryCharacterRepository) UserWithdrawMoney(userNo, charNo int, amount int64) (int64, error) {
	if amount <= 0 {
		return 0, errors.New("invalid withdraw amount")
	}
	r.mu.Lock()
	if r.userMoney == nil || r.userMoney[userNo] < amount {
		r.mu.Unlock()
		return 0, errors.New("insufficient storage money")
	}
	r.userMoney[userNo] -= amount
	rem := r.userMoney[userNo]
	r.mu.Unlock()
	_ = r.AddCharMoney(charNo, amount)
	return rem, nil
}

func (r *MemoryCharacterRepository) UserGetMoney(userNo int) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.userMoney == nil {
		return 0, nil
	}
	return r.userMoney[userNo], nil
}

func (r *MemoryCharacterRepository) EnchantAddGem(charNo int, equipKey, gemKey int64, socketIdx byte, gemID int, restGemCount byte) error {
	if gemKey > 0 {
		_ = r.DeleteItem(gemKey)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.itemEnchants == nil {
		r.itemEnchants = make(map[int64][]*ItemEnchantData)
	}
	list := r.itemEnchants[equipKey]
	replaced := false
	for _, e := range list {
		if e.Slot == socketIdx {
			e.GemItemID = gemID
			e.RestCount = restGemCount
			replaced = true
			break
		}
	}
	if !replaced {
		r.itemEnchants[equipKey] = append(list, &ItemEnchantData{
			Slot:      socketIdx,
			GemItemID: gemID,
			RestCount: restGemCount,
		})
	}
	return nil
}

func (r *MemoryCharacterRepository) EnchantAddSocket(charNo int, equipKey, drillKey int64, drillLotType int16) (byte, int64, error) {
	if drillKey > 0 {
		_ = r.DeleteItem(drillKey)
	}
	opts, _ := r.GetItemOptions(equipKey)
	cur := byte(opts[524])
	cur++
	_ = r.SetItemOption(equipKey, 524, int64(cur))
	return cur, 0, nil
}

func (r *MemoryCharacterRepository) EnchantRemoveGem(charNo int, equipKey, removerKey int64, socketIdx, needItem byte, removerLotType int16) (int64, error) {
	if needItem > 0 && removerKey > 0 {
		_ = r.DeleteItem(removerKey)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	list := r.itemEnchants[equipKey]
	for i, e := range list {
		if e.Slot == socketIdx {
			r.itemEnchants[equipKey] = append(list[:i], list[i+1:]...)
			break
		}
	}
	return 0, nil
}

func (r *MemoryCharacterRepository) SetEnchantGemLot(charNo int, equipKey int64, slot, restGemCount byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.itemEnchants[equipKey] {
		if e.Slot == slot {
			e.RestCount = restGemCount
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) PutOnBelongedCancel(charNo int, cancelKey int64, cancelOptType int16, cancelOptData int64, targetKey int64, targetOptType int16, targetOptData int64) error {
	if cancelOptData <= 0 {
		_ = r.DeleteItem(cancelKey)
	} else {
		_ = r.UpdateItemFlags(cancelKey, int(cancelOptData))
	}
	return r.SetItemOption(targetKey, targetOptType, targetOptData)
}

func (r *MemoryCharacterRepository) SetMinimonAutoItem(charNo, monType, slot, itemID int) error {
	return nil
}

func (r *MemoryCharacterRepository) FurnitureEndureUpdate(owner int, kitKey int64, kitOptType int16, kitLot int16, furnitureKey int64, endureOptType int16, endureOptData int64) error {
	if kitLot <= 0 {
		_ = r.DeleteItem(kitKey)
	} else {
		_ = r.UpdateItemFlags(kitKey, int(kitLot))
	}
	return r.SetItemOption(furnitureKey, endureOptType, endureOptData)
}

func (r *MemoryCharacterRepository) CostumeWeaponEndureUpdate(owner int, kitKey int64, kitOptType int16, kitLot int16, costumeKey int64, endureOptType int16, endureOptData int64) error {
	if kitLot <= 0 {
		_ = r.DeleteItem(kitKey)
	} else {
		_ = r.UpdateItemFlags(kitKey, int(kitLot))
	}
	return r.SetItemOption(costumeKey, endureOptType, endureOptData)
}

func (r *MemoryCharacterRepository) SafeTradeAll(buyer, seller int, moveMoney, itemKey int64, storageType byte, storage int16) error {
	if moveMoney > 0 {
		ok, err := r.SubCharMoney(buyer, moveMoney)
		if err != nil || !ok {
			return errors.New("insufficient buyer money")
		}
		_ = r.AddCharMoney(seller, moveMoney)
	}
	if itemKey > 0 {
		r.mu.Lock()
		defer r.mu.Unlock()
		for oldOwner, list := range r.items {
			for i, it := range list {
				if it.ItemKey == itemKey {
					it.Owner = buyer
					it.StorageType = storageType
					it.Storage = storage
					it.Date = time.Now()
					if oldOwner != buyer {
						r.items[oldOwner] = append(list[:i], list[i+1:]...)
						r.items[buyer] = append(r.items[buyer], it)
					}
					return nil
				}
			}
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) SafeTradeLot(buyer int, buyerItemKey int64, buyerStorageType byte, buyerStorage int16, buyerItemID, buyerFlags int, buyerOptType int16, seller int, sellerItemKey int64, sellerOptType int16, moveMoney, splitLot int64) error {
	if moveMoney > 0 {
		ok, err := r.SubCharMoney(buyer, moveMoney)
		if err != nil || !ok {
			return errors.New("insufficient buyer money")
		}
		_ = r.AddCharMoney(seller, moveMoney)
	}
	if sellerItemKey > 0 && splitLot > 0 {
		for _, it := range r.items[seller] {
			if it.ItemKey == sellerItemKey {
				rem := int(it.Flags) - int(splitLot)
				if rem <= 0 {
					_ = r.DeleteItem(sellerItemKey)
				} else {
					_ = r.UpdateItemFlags(sellerItemKey, rem)
				}
				break
			}
		}
		return r.CreateItem(&Item{
			ItemKey:     buyerItemKey,
			Owner:       buyer,
			StorageType: buyerStorageType,
			Storage:     buyerStorage,
			ItemID:      uint16(buyerItemID),
			Flags:       int32(splitLot),
			Date:        time.Now(),
		})
	}
	return nil
}

func (r *MemoryCharacterRepository) SafeTradeMerge(buyer int, buyerItemKey int64, buyerOptType int16, seller int, sellerItemKey int64, sellerOptType int16, moveMoney, splitLot int64, isDelete byte) error {
	if moveMoney > 0 {
		ok, err := r.SubCharMoney(buyer, moveMoney)
		if err != nil || !ok {
			return errors.New("insufficient buyer money")
		}
		_ = r.AddCharMoney(seller, moveMoney)
	}
	if isDelete == 1 {
		_ = r.DeleteItem(sellerItemKey)
	} else {
		for _, it := range r.items[seller] {
			if it.ItemKey == sellerItemKey {
				_ = r.UpdateItemFlags(sellerItemKey, int(it.Flags)-int(splitLot))
				break
			}
		}
	}
	for _, it := range r.items[buyer] {
		if it.ItemKey == buyerItemKey {
			_ = r.UpdateItemFlags(buyerItemKey, int(it.Flags)+int(splitLot))
			break
		}
	}
	return nil
}

// Progression
func (r *MemoryCharacterRepository) GetCharacterSkills(charNo int) ([]*CharacterSkill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.skills[charNo], nil
}

func (r *MemoryCharacterRepository) SaveCharacterSkill(charNo int, skill *CharacterSkill) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skills[charNo] = append(r.skills[charNo], skill)
	return nil
}

func (r *MemoryCharacterRepository) SaveCharacterSkills(charNo int, skills []*CharacterSkill) error {
	for _, sk := range skills {
		_ = r.SaveCharacterSkill(charNo, sk)
	}
	return nil
}

func (r *MemoryCharacterRepository) SetSkillCoolTime(charNo int, skillNo int16, writeTime int64, coolTime int) error {
	return nil
}

func (r *MemoryCharacterRepository) SetSkillPower(charNo int, skillNo int16, pd, ps, keepTime, coolTime byte) error {
	return nil
}

func (r *MemoryCharacterRepository) SetSingleSkillPower(charNo int, skillNo int16, powerType byte, power byte) error {
	return nil
}

func (r *MemoryCharacterRepository) DeleteCharacterSkill(charNo int, skillNo int16) error {
	return nil
}

func (r *MemoryCharacterRepository) GetCharacterShortcuts(charNo int) ([]*Shortcut, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.shortcuts[charNo], nil
}

func (r *MemoryCharacterRepository) SaveCharacterShortcuts(charNo int, shortcuts []*Shortcut) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shortcuts[charNo] = shortcuts
	return nil
}

func (r *MemoryCharacterRepository) GetCharacterOptions(charNo int) (*CharacterOptions, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	opt, ok := r.options[charNo]
	if !ok {
		return &CharacterOptions{CharNo: charNo}, nil
	}
	return opt, nil
}

func (r *MemoryCharacterRepository) SaveCharacterOption(charNo int, optType string, data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	opt, ok := r.options[charNo]
	if !ok {
		opt = &CharacterOptions{CharNo: charNo}
		r.options[charNo] = opt
	}
	switch strings.ToLower(optType) {
	case "shortcut", "shortcutdata", "shortcut_data":
		opt.ShortCutData = data
	case "shortcutsize", "shortcut_size":
		opt.ShortCutSize = data
	case "game":
		opt.Game = data
	case "video":
		opt.Video = data
	case "sound":
		opt.Sound = data
	case "windowspos", "windows_pos", "window_pos":
		opt.WindowsPos = data
	case "keymapping", "key_mapping", "key_map":
		opt.KeyMapping = data
	default:
		return fmt.Errorf("unknown option type: %s", optType)
	}
	return nil
}

func (r *MemoryCharacterRepository) DeleteCharacterOptions(charNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.options, charNo)
	return nil
}

func (r *MemoryCharacterRepository) GetImproveShortCutList(charNo int) ([]*ImproveShortcutEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*ImproveShortcutEntry
	for _, e := range r.improveSC[charNo] {
		if e.CodeNo > 0 {
			res = append(res, e)
		}
	}
	return res, nil
}

func (r *MemoryCharacterRepository) SetImproveShortCut(charNo int, slotNo byte, codeNo uint16, value int32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.improveSC[charNo] {
		if e.SlotNo == slotNo {
			e.CodeNo = codeNo
			e.Value = value
			return nil
		}
	}
	r.improveSC[charNo] = append(r.improveSC[charNo], &ImproveShortcutEntry{
		SlotNo: slotNo,
		CodeNo: codeNo,
		Value:  value,
	})
	return nil
}

func (r *MemoryCharacterRepository) GetImproveKeyMap(charNo int) ([]*ImproveKeyMapEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.improveKeyMap[charNo], nil
}

func (r *MemoryCharacterRepository) SetImproveKeyMap(charNo int, funcNo uint16, extendKey, asciiCode byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.improveKeyMap[charNo] {
		if e.FunctionNo == funcNo {
			e.ExtendKey = extendKey
			e.ASCIICode = asciiCode
			return nil
		}
	}
	r.improveKeyMap[charNo] = append(r.improveKeyMap[charNo], &ImproveKeyMapEntry{
		FunctionNo: funcNo,
		ExtendKey:  extendKey,
		ASCIICode:  asciiCode,
	})
	return nil
}

func (r *MemoryCharacterRepository) InitImproveKeyMap(charNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.improveKeyMap[charNo] = []*ImproveKeyMapEntry{
		{FunctionNo: 1, ExtendKey: 0, ASCIICode: 'W'},
	}
	return nil
}

func (r *MemoryCharacterRepository) GetImproveGameOptions(charNo int) ([]*ImproveGameOptionEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.improveGameOp[charNo], nil
}

func (r *MemoryCharacterRepository) SetImproveGameOption(charNo int, optionNo uint16, value byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.improveGameOp[charNo] {
		if e.OptionNo == optionNo {
			e.Value = value
			return nil
		}
	}
	r.improveGameOp[charNo] = append(r.improveGameOp[charNo], &ImproveGameOptionEntry{
		OptionNo: optionNo,
		Value:    value,
	})
	return nil
}

func (r *MemoryCharacterRepository) InitImproveGameOptions(charNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.improveGameOp, charNo)
	return nil
}

func (r *MemoryCharacterRepository) GetInGameGuideViewList(charNo int) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.guides[charNo], nil
}

func (r *MemoryCharacterRepository) SetInGameGuideView(charNo int, guideNo byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, g := range r.guides[charNo] {
		if g == guideNo {
			return nil
		}
	}
	r.guides[charNo] = append(r.guides[charNo], guideNo)
	return nil
}

func (r *MemoryCharacterRepository) GetCharacterQuests(charNo int) ([]*QuestData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) GetCharacterDoneQuests(charNo int) ([]int, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) AddCharacterQuest(charNo int, questNo int, status byte, data []byte) error {
	return nil
}

func (r *MemoryCharacterRepository) SetCharacterQuest(charNo int, questNo int, status byte, data []byte) error {
	return nil
}

func (r *MemoryCharacterRepository) DeleteCharacterQuest(charNo int, questNo int) error {
	return nil
}

func (r *MemoryCharacterRepository) GetAttendance(charNo int) ([]time.Time, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) CheckAttendance(charNo int, dDate time.Time) error {
	return nil
}

func (r *MemoryCharacterRepository) GetAttendanceRewards(charNo int) ([]int, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetAttendanceReward(charNo int, rewardID int) error {
	_, err := r.SetAttendanceRewardAt(charNo, time.Now(), rewardID)
	return err
}

func (r *MemoryCharacterRepository) SetAttendanceRewardAt(charNo int, eventDate time.Time, rewardID int) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.attendRewards[charNo] == nil {
		r.attendRewards[charNo] = make(map[int]bool)
	}
	if r.attendRewards[charNo][rewardID] {
		return -1111, nil
	}
	r.attendRewards[charNo][rewardID] = true
	return 0, nil
}

func (r *MemoryCharacterRepository) GetCollectionCards(charNo int) ([]*CollectionCardData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	src := r.collectionCards[charNo]
	out := make([]*CollectionCardData, 0, len(src))
	for _, c := range src {
		cp := *c
		out = append(out, &cp)
	}
	return out, nil
}

func (r *MemoryCharacterRepository) SetCollectionCard(charNo int, cardID int16, count int16) error {
	return r.RegisterCollectionCard(charNo, int(cardID), 1, byte(count), 0)
}

func (r *MemoryCharacterRepository) RegisterCollectionCard(charNo int, cardID int, serialNo int64, starNumber byte, itemKey int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.collectionCards == nil {
		r.collectionCards = make(map[int][]*CollectionCardData)
	}
	if itemKey > 0 {
		items := r.items[charNo]
		for i, it := range items {
			if it.ItemKey == itemKey {
				r.items[charNo] = append(items[:i], items[i+1:]...)
				break
			}
		}
	}
	for _, c := range r.collectionCards[charNo] {
		if int(c.CardID) == cardID {
			c.SerialNo = serialNo
			c.StarNumber = starNumber
			c.Count = int16(starNumber)
			c.RegistTime = time.Now()
			return nil
		}
	}
	r.collectionCards[charNo] = append(r.collectionCards[charNo], &CollectionCardData{
		CardID:     int16(cardID),
		Count:      int16(starNumber),
		SerialNo:   serialNo,
		StarNumber: starNumber,
		RegistTime: time.Now(),
	})
	return nil
}

func (r *MemoryCharacterRepository) GetCollectionFavorites(charNo int) ([]*CollectionBookmarkData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	src := r.collectionBookmarks[charNo]
	out := make([]*CollectionBookmarkData, 0, len(src))
	for _, bm := range src {
		cp := *bm
		for _, c := range r.collectionCards[charNo] {
			if uint16(c.CardID) == cp.CardID {
				cp.SerialNo = c.SerialNo
				cp.StarNumber = c.StarNumber
				break
			}
		}
		out = append(out, &cp)
	}
	return out, nil
}

func (r *MemoryCharacterRepository) SetCollectionFavorite(charNo int, slot byte, cardID int, regist bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.collectionBookmarks == nil {
		r.collectionBookmarks = make(map[int][]*CollectionBookmarkData)
	}
	list := r.collectionBookmarks[charNo]
	if !regist {
		for i, bm := range list {
			if bm.Slot == uint16(slot) {
				r.collectionBookmarks[charNo] = append(list[:i], list[i+1:]...)
				return nil
			}
		}
		return nil
	}
	for _, bm := range list {
		if bm.Slot == uint16(slot) {
			bm.CardID = uint16(cardID)
			return nil
		}
	}
	r.collectionBookmarks[charNo] = append(list, &CollectionBookmarkData{
		Slot:   uint16(slot),
		CardID: uint16(cardID),
	})
	return nil
}

func (r *MemoryCharacterRepository) GetCollectionRewards(charNo int) ([]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rw := r.collectionRewards[charNo]
	out := make([]int, 0, len(rw))
	for id, ok := range rw {
		if ok {
			out = append(out, id)
		}
	}
	return out, nil
}

func (r *MemoryCharacterRepository) SetCollectionReward(charNo int, rewardID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.collectionRewards == nil {
		r.collectionRewards = make(map[int]map[int]bool)
	}
	if r.collectionRewards[charNo] == nil {
		r.collectionRewards[charNo] = make(map[int]bool)
	}
	if r.collectionRewards[charNo][rewardID] {
		return fmt.Errorf("collection reward %d already claimed", rewardID)
	}
	r.collectionRewards[charNo][rewardID] = true
	return nil
}

func (r *MemoryCharacterRepository) GetCollectionSerialList() (map[int]int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[int]int64, len(r.collectionSerials))
	for k, v := range r.collectionSerials {
		out[k] = v
	}
	return out, nil
}

func (r *MemoryCharacterRepository) SetCollectionSerialCount(cardID int, count int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.collectionSerials == nil {
		r.collectionSerials = make(map[int]int64)
	}
	r.collectionSerials[cardID] = count
	return nil
}

func (r *MemoryCharacterRepository) GetTutorialCharInfo(charNo int) (byte, byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if ch, ok := r.characters[charNo]; ok {
		return byte(ch.TutorialState), ch.TutorialStep, nil
	}
	return 0, 0, nil
}

func (r *MemoryCharacterRepository) SetTutorialCharInfo(userNo, charNo int, state, step byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.characters[charNo]; ok {
		ch.TutorialState = int32(state)
		ch.TutorialStep = step
	}
	if state >= 1 {
		r.tutorialSkip[userNo] = true
	}
	return nil
}

func (r *MemoryCharacterRepository) GetTutorialIsSkip(userNo int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.tutorialSkip[userNo], nil
}

// Guild
func (r *MemoryCharacterRepository) CreateGuild(masterCharID int, name, password string, guildType byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := len(r.guilds) + 1
	g := &GuildData{
		ID:         id,
		Name:       name,
		Password:   password,
		MaxMembers: 50,
		Type:       guildType,
		CreateDate: time.Now(),
	}
	r.guilds[id] = g
	masterName := fmt.Sprintf("Char_%d", masterCharID)
	if ch, ok := r.characters[masterCharID]; ok {
		masterName = ch.Name
	}
	r.guildMbrs[id] = append(r.guildMbrs[id], &GuildMemberData{
		GuildID:  id,
		CharID:   masterCharID,
		CharName: masterName,
		Grade:    0,
		JoinDate: time.Now(),
	})
	return id, nil
}

func (r *MemoryCharacterRepository) GetGuild(guildID int) (*GuildData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.guilds[guildID]
	if !ok {
		return nil, errors.New("guild not found")
	}
	return g, nil
}

func (r *MemoryCharacterRepository) GetGuildByCharID(charID int) (*GuildData, *GuildMemberData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for gID, mbrs := range r.guildMbrs {
		for _, m := range mbrs {
			if m.CharID == charID {
				return r.guilds[gID], m, nil
			}
		}
	}
	return nil, nil, nil
}

func (r *MemoryCharacterRepository) GetGuildMembers(guildID int) ([]*GuildMemberData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.guildMbrs[guildID], nil
}

func (r *MemoryCharacterRepository) GetAllGuilds() ([]*GuildListItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*GuildListItem
	for _, g := range r.guilds {
		masterName := ""
		for _, m := range r.guildMbrs[g.ID] {
			if m.Grade == 0 {
				masterName = m.CharName
				break
			}
		}
		res = append(res, &GuildListItem{
			ID:             g.ID,
			GuildID:        uint32(g.ID),
			Name:           g.Name,
			GuildName:      g.Name,
			MasterCharName: masterName,
			MemberCount:    uint16(len(r.guildMbrs[g.ID])),
			MaxMembers:     g.MaxMembers,
			Type:           g.Type,
			Grade:          g.Grade,
			Fame:           g.Fame,
			Exp:            g.Exp,
		})
	}
	return res, nil
}

func (r *MemoryCharacterRepository) AddGuildMember(guildID, charID int, grade byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	mName := fmt.Sprintf("Char_%d", charID)
	if ch, ok := r.characters[charID]; ok {
		mName = ch.Name
	}
	r.guildMbrs[guildID] = append(r.guildMbrs[guildID], &GuildMemberData{
		GuildID:  guildID,
		CharID:   charID,
		CharName: mName,
		Grade:    grade,
		JoinDate: time.Now(),
	})
	return nil
}

func (r *MemoryCharacterRepository) RemoveGuildMember(guildID, charID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	mbrs := r.guildMbrs[guildID]
	for i, m := range mbrs {
		if m.CharID == charID {
			r.guildMbrs[guildID] = append(mbrs[:i], mbrs[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) SetGuildMemberGrade(guildID, charID int, grade byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.guildMbrs[guildID] {
		if m.CharID == charID {
			m.Grade = grade
			return nil
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) SetGuildNotify(guildID int, charName, notify string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g, ok := r.guilds[guildID]; ok {
		g.Notify = notify
		g.NotifyCharID = charName
		g.NotifyDate = time.Now()
		return nil
	}
	return errors.New("guild not found")
}

func (r *MemoryCharacterRepository) DismissGuild(guildID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g, ok := r.guilds[guildID]; ok {
		g.DismissStatus = 1
		g.DismissDate = time.Now()
		return nil
	}
	return errors.New("guild not found")
}

func (r *MemoryCharacterRepository) SetGuildMoney(guildID int, money uint64) error {
	return nil
}

func (r *MemoryCharacterRepository) SetGuildGrade(guildID int, grade byte) error {
	return nil
}

func (r *MemoryCharacterRepository) GetGuildToken(guildID int) (int, error) {
	return 0, nil
}

func (r *MemoryCharacterRepository) SetGuildToken(guildID int, token int) error {
	return nil
}

func (r *MemoryCharacterRepository) GetGuildHistory(guildID int) ([]*GuildHistoryData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) AddGuildHistory(guildID int, historyType byte, param1, param2 int, text string) error {
	return nil
}

func (r *MemoryCharacterRepository) CreateGuildAcademy(guildID, masterCharID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.academies[guildID] = &GuildAcademyData{
		GuildID:    guildID,
		MasterChar: masterCharID,
		MaxMember:  50,
	}
	return nil
}

func (r *MemoryCharacterRepository) GetGuildAcademy(guildID int) (*GuildAcademyData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.academies[guildID]
	if !ok {
		return nil, errors.New("academy not found")
	}
	return a, nil
}

func (r *MemoryCharacterRepository) GetAcademyByCharID(charID int) (*GuildAcademyData, *GuildAcademyMemberData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for gID, mbrs := range r.academyMbrs {
		for _, m := range mbrs {
			if m.CharID == charID {
				return r.academies[gID], m, nil
			}
		}
	}
	return nil, nil, nil
}

func (r *MemoryCharacterRepository) GetAcademyMembers(guildID int) ([]*GuildAcademyMemberData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.academyMbrs[guildID], nil
}

func (r *MemoryCharacterRepository) AddAcademyMember(guildID, charID int, joinLevel byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	aName := fmt.Sprintf("Char_%d", charID)
	if ch, ok := r.characters[charID]; ok {
		aName = ch.Name
	}
	r.academyMbrs[guildID] = append(r.academyMbrs[guildID], &GuildAcademyMemberData{
		GuildID:   guildID,
		CharID:    charID,
		CharName:  aName,
		JoinLevel: joinLevel,
		JoinDate:  time.Now(),
	})
	return nil
}

func (r *MemoryCharacterRepository) RemoveAcademyMember(guildID, charID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	mbrs := r.academyMbrs[guildID]
	for i, m := range mbrs {
		if m.CharID == charID {
			r.academyMbrs[guildID] = append(mbrs[:i], mbrs[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) SetAcademyNotify(guildID int, charName, notify string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if a, ok := r.academies[guildID]; ok {
		a.Notify = notify
		a.NotifyDate = time.Now()
		return nil
	}
	return errors.New("academy not found")
}

func (r *MemoryCharacterRepository) SetAcademyChatBan(guildID, charID int, banMinutes int) error {
	return nil
}

func (r *MemoryCharacterRepository) SetAcademyGraduate(guildID, charID int) error {
	return nil
}

func (r *MemoryCharacterRepository) GetAcademyGraduates(guildID int) ([]int, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SubGuildMoney(guildID, memberCharNo int, amount int64) (int64, int64, error) {
	r.mu.Lock()
	g, ok := r.guilds[guildID]
	if !ok || g.Money < uint64(amount) {
		r.mu.Unlock()
		return 0, 0, errors.New("insufficient guild money")
	}
	g.Money -= uint64(amount)
	rem := int64(g.Money)
	r.mu.Unlock()
	_ = r.AddCharMoney(memberCharNo, amount)
	ch, _ := r.GetCharacterByID(memberCharNo)
	var mMoney int64
	if ch != nil {
		mMoney = ch.TCharacterState.Money
	}
	return mMoney, rem, nil
}

func (r *MemoryCharacterRepository) GetAcademyMoney(guildID int) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.academyMoney == nil {
		return 0, nil
	}
	return r.academyMoney[guildID], nil
}

func (r *MemoryCharacterRepository) DepositAcademyMoney(guildID, charNo int, amount int64) error {
	ok, err := r.SubCharMoney(charNo, amount)
	if err != nil || !ok {
		return errors.New("insufficient char money")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.academyMoney == nil {
		r.academyMoney = make(map[int]int64)
	}
	r.academyMoney[guildID] += amount
	return nil
}

func (r *MemoryCharacterRepository) WithdrawAcademyMoney(guildID, charNo int, amount int64) error {
	r.mu.Lock()
	if r.academyMoney == nil || r.academyMoney[guildID] < amount {
		r.mu.Unlock()
		return errors.New("insufficient academy money")
	}
	r.academyMoney[guildID] -= amount
	r.mu.Unlock()
	return r.AddCharMoney(charNo, amount)
}

func (r *MemoryCharacterRepository) GetAcademyRewardsAll(guildID, typeID int) ([]*AcademyRewardData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*AcademyRewardData
	for _, rew := range r.academyRewards[guildID] {
		if (typeID == 1 && rew.TypeID == 1) || (typeID != 1 && rew.TypeID > 1) {
			cp := *rew
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (r *MemoryCharacterRepository) DelAcademyReward(guildID int, level int16, typeID byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := r.academyRewards[guildID]
	var kept []*AcademyRewardData
	for _, rew := range list {
		if rew.Level == level && ((typeID == 1 && rew.TypeID == 1) || (typeID != 1 && rew.TypeID > 1)) {
			continue
		}
		kept = append(kept, rew)
	}
	r.academyRewards[guildID] = kept
	return nil
}

func (r *MemoryCharacterRepository) GetGuildTournament(gtNo int) (*GuildTournamentData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) GetAllGuildTournaments() ([]*GuildTournamentData, error) {
	return nil, nil
}

// Social
func (r *MemoryCharacterRepository) GetFriends(charNo int) ([]*FriendInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.friends[charNo], nil
}

func (r *MemoryCharacterRepository) AddFriend(charNo, friendCharNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	fName := fmt.Sprintf("Friend_%d", friendCharNo)
	if fChar, ok := r.characters[friendCharNo]; ok {
		fName = fChar.Name
	}
	r.friends[charNo] = append(r.friends[charNo], &FriendInfo{
		CharNo:    friendCharNo,
		CharName:  fName,
		CreatedAt: time.Now(),
	})
	return nil
}

func (r *MemoryCharacterRepository) RemoveFriend(charNo, friendCharNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := r.friends[charNo]
	for i, f := range list {
		if f.CharNo == friendCharNo {
			r.friends[charNo] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) GetFriendPoints(charNo int) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.friendPts[charNo], nil
}

func (r *MemoryCharacterRepository) SetFriendPoints(charNo int, points int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.friendPts[charNo] = points
	return nil
}

func (r *MemoryCharacterRepository) GetKingdomQuests() ([]*KingdomQuestData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) GetKingdomQuest(questNo int) (*KingdomQuestData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetKingdomQuest(questNo int, data string) error {
	return nil
}

func (r *MemoryCharacterRepository) DelKingdomQuest(questNo int) error {
	return nil
}

func (r *MemoryCharacterRepository) CheckHolyPromiseRing(charNo, upCharNo int) (bool, error) {
	return false, nil
}

func (r *MemoryCharacterRepository) GetMasterPupils(charNo int) ([]*MasterPupilData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetMaster(charNo int, masterCharNo int) error {
	return nil
}

func (r *MemoryCharacterRepository) SetPupil(charNo int, pupilCharNo int) error {
	return nil
}

func (r *MemoryCharacterRepository) DelMaster(charNo int) error {
	return nil
}

func (r *MemoryCharacterRepository) DelPupils(charNo int) error {
	return nil
}

func (r *MemoryCharacterRepository) DelAllMasterPupils(charNo int) error {
	return nil
}

func (r *MemoryCharacterRepository) GetPetInfo(ownerNo int) (*PetData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetPetName(petNo int, ownerNo int, name string) error {
	return nil
}

func (r *MemoryCharacterRepository) SetPetTendency(petNo int, tendency byte) error {
	return nil
}

func (r *MemoryCharacterRepository) GetMiniHouseURL(charNo int) (string, error) {
	return "", nil
}

func (r *MemoryCharacterRepository) SetMiniHouseURL(charNo int, url string) error {
	return nil
}

func (r *MemoryCharacterRepository) GetMiniHouseVisitors(charNo int) (int, error) {
	return 0, nil
}

func (r *MemoryCharacterRepository) SetMiniHouseVisitors(charNo int, visitors int) error {
	return nil
}

func (r *MemoryCharacterRepository) SetHolyPromise(charNo, upCharNo int, flags byte) error { return nil }
func (r *MemoryCharacterRepository) DelHolyPromise(charNo, upCharNo int) error             { return nil }
func (r *MemoryCharacterRepository) DelHolyPromiseChar(charNo int) error                 { return nil }
func (r *MemoryCharacterRepository) GetHolyPromise(charNo int) (*HolyPromiseData, error) {
	return &HolyPromiseData{CharNo: charNo, UpCharNo: 0, Flags: 0}, nil
}

func (r *MemoryCharacterRepository) CreateParty(masterCharNo int, flags int) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.parties == nil {
		r.parties = make(map[int]*PartyData)
	}
	if r.partyMembers == nil {
		r.partyMembers = make(map[int][]*PartyMemberData)
	}
	if r.nextPartyNo <= 0 {
		r.nextPartyNo = 1
	}
	partyNo := r.nextPartyNo
	r.nextPartyNo++
	r.parties[partyNo] = &PartyData{PartyNo: partyNo, MasterCharNo: masterCharNo, Flags: flags}
	chName := ""
	var chClass, chLevel byte
	var chZone string
	var chHP, chSP int
	if ch, ok := r.characters[masterCharNo]; ok {
		chName = ch.Name
		chClass = ch.Job
		chLevel = ch.Level
		chZone = ch.TownName
		chHP = ch.CurHP
		chSP = ch.CurSP
	}
	r.partyMembers[partyNo] = []*PartyMemberData{{
		CharNo: masterCharNo, Flag: 1, CharID: chName, Class: chClass, Level: chLevel, LoginZone: chZone, HP: int16(chHP), SP: int16(chSP),
	}}
	return partyNo, nil
}

func (r *MemoryCharacterRepository) AddPartyMember(partyNo int, charNo int, flag byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.partyMembers == nil {
		r.partyMembers = make(map[int][]*PartyMemberData)
	}
	for _, m := range r.partyMembers[partyNo] {
		if m.CharNo == charNo {
			return nil
		}
	}
	chName := ""
	var chClass, chLevel byte
	var chZone string
	var chHP, chSP int
	if ch, ok := r.characters[charNo]; ok {
		chName = ch.Name
		chClass = ch.Job
		chLevel = ch.Level
		chZone = ch.TownName
		chHP = ch.CurHP
		chSP = ch.CurSP
	}
	r.partyMembers[partyNo] = append(r.partyMembers[partyNo], &PartyMemberData{
		CharNo: charNo, Flag: flag, CharID: chName, Class: chClass, Level: chLevel, LoginZone: chZone, HP: int16(chHP), SP: int16(chSP),
	})
	return nil
}

func (r *MemoryCharacterRepository) RemovePartyMember(partyNo int, charNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := r.partyMembers[partyNo]
	for i, m := range list {
		if m.CharNo == charNo {
			r.partyMembers[partyNo] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return nil
}

func (r *MemoryCharacterRepository) SetPartyMaster(partyNo int, masterCharNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.parties[partyNo]; ok {
		p.MasterCharNo = masterCharNo
	}
	return nil
}

func (r *MemoryCharacterRepository) SetPartyFlags(partyNo int, flags int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, ok := r.parties[partyNo]; ok {
		p.Flags = flags
	}
	return nil
}

func (r *MemoryCharacterRepository) DeleteParty(partyNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.parties, partyNo)
	delete(r.partyMembers, partyNo)
	return nil
}

func (r *MemoryCharacterRepository) GetParty(partyNo int) (*PartyData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, ok := r.parties[partyNo]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, fmt.Errorf("party %d not found", partyNo)
}

func (r *MemoryCharacterRepository) GetPartyMembers(partyNo int) ([]*PartyMemberData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	src := r.partyMembers[partyNo]
	out := make([]*PartyMemberData, 0, len(src))
	for _, m := range src {
		cp := *m
		out = append(out, &cp)
	}
	return out, nil
}

func (r *MemoryCharacterRepository) GetPartyByChar(charNo int) (int, int, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for pNo, members := range r.partyMembers {
		for _, m := range members {
			if m.CharNo == charNo {
				if p, ok := r.parties[pNo]; ok {
					return p.PartyNo, p.MasterCharNo, p.Flags, nil
				}
				return pNo, charNo, 0, nil
			}
		}
	}
	return 0, 0, 0, nil
}

func (r *MemoryCharacterRepository) AddChatRestrict(userNo int, charNo int) error { return nil }
func (r *MemoryCharacterRepository) DelChatRestrict(userNo int, charNo int) error { return nil }
func (r *MemoryCharacterRepository) DelAllChatRestrict(userNo int) error          { return nil }

func (r *MemoryCharacterRepository) AddPrison(gmUserID string, targetCharName string, minutes int16, reason, remark string) error {
	return nil
}
func (r *MemoryCharacterRepository) EndPrison(charNo int) error { return nil }
func (r *MemoryCharacterRepository) AlterPrison(gmUserID string, targetCharName string, minutes int16, reason, remark string) error {
	return nil
}

func (r *MemoryCharacterRepository) SetClientLoadingBug(charNo int, mapName string, failType int16, failCnt int16) error {
	return nil
}
func (r *MemoryCharacterRepository) SetSpamFilter(userNo int, isSpamSpread, isSpamBlock byte, rate byte) error {
	return nil
}
func (r *MemoryCharacterRepository) GetSpamFilter(userNo int) (bool, bool, byte, error) {
	return false, false, 0, nil
}
func (r *MemoryCharacterRepository) SetSpamBlockTime(charNo int, blockTime, standbyTime int) error {
	return nil
}
func (r *MemoryCharacterRepository) GetSpamBlockTime(charNo int) (int, int, error) {
	return 0, 0, nil
}

// Topology
func (r *MemoryCharacterRepository) IssueZoneTransfer(worldNo uint8, zoneNo uint16, userNo int, userID string, charNo int, charID string, randomID uint16, adminLevel uint8, clientIP, targetMap string, targetX, targetY int) error {
	return nil
}

func (r *MemoryCharacterRepository) ConsumeZoneTransfer(charID string, randomID uint16, clientIP string) (*ZoneTransfer, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) UpsertZoneNode(worldNo uint8, zoneNo uint16, zoneIP string, zonePort uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.zoneNodes[zoneNo] = &ZoneNodeInfo{
		WorldNo:       worldNo,
		ZoneNo:        zoneNo,
		ZoneIP:        zoneIP,
		ZonePort:      zonePort,
		IsOnline:      true,
		LastHeartbeat: time.Now(),
	}
	return nil
}

func (r *MemoryCharacterRepository) HeartbeatZoneNode(worldNo uint8, zoneNo uint16, currentPlayers int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if n, ok := r.zoneNodes[zoneNo]; ok {
		n.CurrentPlayers = currentPlayers
		n.LastHeartbeat = time.Now()
	}
	return nil
}

func (r *MemoryCharacterRepository) GetZoneNodes(worldNo uint8) ([]*ZoneNodeInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*ZoneNodeInfo
	for _, n := range r.zoneNodes {
		if n.WorldNo == worldNo {
			res = append(res, n)
		}
	}
	return res, nil
}

func (r *MemoryCharacterRepository) ResolveMapEndpoint(worldNo uint8, zoneNo uint16) (string, uint16, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if n, ok := r.zoneNodes[zoneNo]; ok && n.IsOnline {
		return n.ZoneIP, n.ZonePort, true, nil
	}
	if n, ok := r.zoneNodes[0]; ok && n.IsOnline {
		return n.ZoneIP, n.ZonePort, true, nil
	}
	return "127.0.0.1", 9210, true, nil
}
