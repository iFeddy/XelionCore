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
	ErrCharacterNotFound = errors.New("character not found")
	ErrCharacterExists   = errors.New("character name already exists")
	ErrSlotInUse         = errors.New("character slot already in use")
	ErrSlotOutOfRange    = errors.New("character slot out of range (0-5)")
	ErrItemNotFound      = errors.New("item not found")
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

	// Equips: Slot index -> item ID
	Equips map[byte]uint16 `json:"equips"`

	// Flat convenience fields for rapid gameplay access
	ID        int               `json:"id"`
	AccountID int               `json:"account_id"`
	Name      string            `json:"name"`
	Slot      byte              `json:"slot"`
	Job       byte              `json:"job"`
	Level     byte              `json:"level"`
	Male      bool              `json:"male"`
	Hair      byte              `json:"hair"`
	HairColor byte              `json:"hair_color"`
	Face      byte              `json:"face"`
	Map       uint16            `json:"map"`
	XPos      int               `json:"x_pos"`
	YPos      int               `json:"y_pos"`
	CurHP     int               `json:"cur_hp"`
	CurSP     int               `json:"cur_sp"`
	TownName  string            `json:"town_name"`
	Items     []*Item           `json:"items"`
	Shortcuts []*Shortcut       `json:"shortcuts"`
	Skills    []*CharacterSkill `json:"skills"`
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
	CharNo     int    `json:"char_no"`
	SkillNo    uint16 `json:"skill_no"`
	SkillLevel byte   `json:"skill_level"`
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
	CharNo    int    `json:"char_no"`
	TitleType int16  `json:"title_type"`
	Title     string `json:"title"`
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

// CollectionCardData modela cartas de colección en tCollection.
type CollectionCardData struct {
	CardID int16 `json:"card_id"`
	Count  int16 `json:"count"`
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
	GTNo int `json:"gt_no"`
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
		c.TCharacterState.HPStones = 15
		c.TCharacterState.SPStones = 7

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
		if c.TownName == "" {
			c.TownName = "Roumen"
		}
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
	AddCharMoney(charNo int, amount int64) error
	SubCharMoney(charNo int, amount int64) (bool, error)
	SaveCharFreeStats(charNo int, str, con, dex, intel, men, freePoints byte) error
	GetCharacterStates(charNo int) ([]*CharacterStateData, error)
	SetCharacterState(charNo int, stateNo int16, level int16, stateTime int, restTime int, count int16, pd, ps, pk, pc byte) error
	GetCharacterTitles(charNo int) ([]*CharacterTitleData, error)
	SetCharacterTitle(charNo int, titleType int16, title string) error
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
	GetItemActionCoolTime(charNo int) ([]*ItemActionCoolTimeData, error)
	SetItemActionCoolTime(charNo int, actionGroup int16, endTime time.Time) error
}

// CharacterProgressionRepository gobierna habilidades, misiones, atajos, opciones de cliente, asistencia y tutorial.
type CharacterProgressionRepository interface {
	GetCharacterSkills(charNo int) ([]*CharacterSkill, error)
	SaveCharacterSkill(charNo int, skill *CharacterSkill) error
	SaveCharacterSkills(charNo int, skills []*CharacterSkill) error
	SetSkillCoolTime(charNo int, skillNo int16, writeTime int64, coolTime int) error
	SetSkillPower(charNo int, skillNo int16, pd, ps, keepTime, coolTime byte) error
	GetCharacterShortcuts(charNo int) ([]*Shortcut, error)
	SaveCharacterShortcuts(charNo int, shortcuts []*Shortcut) error
	GetCharacterOptions(charNo int) (*CharacterOptions, error)
	SaveCharacterOption(charNo int, optType string, data []byte) error
	GetCharacterQuests(charNo int) ([]*QuestData, error)
	GetCharacterDoneQuests(charNo int) ([]int, error)
	AddCharacterQuest(charNo int, questNo int, status byte, data []byte) error
	SetCharacterQuest(charNo int, questNo int, status byte, data []byte) error
	DeleteCharacterQuest(charNo int, questNo int) error
	GetAttendance(charNo int) ([]time.Time, error)
	CheckAttendance(charNo int, dDate time.Time) error
	GetAttendanceRewards(charNo int) ([]int, error)
	SetAttendanceReward(charNo int, rewardID int) error
	GetCollectionCards(charNo int) ([]*CollectionCardData, error)
	SetCollectionCard(charNo int, cardID int16, count int16) error
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

	items       map[int][]*Item
	shortcuts   map[int][]*Shortcut
	skills      map[int][]*CharacterSkill
	options     map[int]*CharacterOptions
	friends     map[int][]*FriendInfo
	friendPts   map[int]int
	guilds      map[int]*GuildData
	guildMbrs   map[int][]*GuildMemberData
	academies   map[int]*GuildAcademyData
	academyMbrs map[int][]*GuildAcademyMemberData
	zoneNodes   map[uint16]*ZoneNodeInfo
}

func NewMemoryCharacterRepositoryWithFile(filePath string) *MemoryCharacterRepository {
	repo := &MemoryCharacterRepository{
		filePath:     filePath,
		characters:   make(map[int]*Character),
		charNameToID: make(map[string]int),
		nextCharID:   1,
		nextItemKey:  1,
		items:        make(map[int][]*Item),
		shortcuts:    make(map[int][]*Shortcut),
		skills:       make(map[int][]*CharacterSkill),
		options:      make(map[int]*CharacterOptions),
		friends:      make(map[int][]*FriendInfo),
		friendPts:    make(map[int]int),
		guilds:       make(map[int]*GuildData),
		guildMbrs:    make(map[int][]*GuildMemberData),
		academies:    make(map[int]*GuildAcademyData),
		academyMbrs:  make(map[int][]*GuildAcademyMemberData),
		zoneNodes:    make(map[uint16]*ZoneNodeInfo),
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
	delete(r.charNameToID, strings.ToLower(ch.Name))
	ch.Name = newName
	r.charNameToID[strings.ToLower(newName)] = charNo
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

func (r *MemoryCharacterRepository) GetCharacterStates(charNo int) ([]*CharacterStateData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetCharacterState(charNo int, stateNo int16, level int16, stateTime int, restTime int, count int16, pd, ps, pk, pc byte) error {
	return nil
}

func (r *MemoryCharacterRepository) GetCharacterTitles(charNo int) ([]*CharacterTitleData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetCharacterTitle(charNo int, titleType int16, title string) error {
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
	return time.Now(), nil
}

func (r *MemoryCharacterRepository) SetRestExpExecTime(charNo int, t time.Time) error {
	return nil
}

func (r *MemoryCharacterRepository) GetAutoPick(charNo int) (bool, error) {
	return false, nil
}

func (r *MemoryCharacterRepository) SetAutoPick(charNo int, autoPick bool) error {
	return nil
}

func (r *MemoryCharacterRepository) SaveCharLink(charNo int, linkZone string, linkX, linkY int) error {
	return nil
}

func (r *MemoryCharacterRepository) AddCharExp(charNo int, exp int64) error {
	return nil
}

func (r *MemoryCharacterRepository) AddCharFame(charNo int, fame int) error {
	return nil
}

func (r *MemoryCharacterRepository) SubCharFame(charNo int, fame int) error {
	return nil
}

func (r *MemoryCharacterRepository) GetCharFame(charNo int) (int, error) {
	return 0, nil
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
	return nil
}

func (r *MemoryCharacterRepository) GetItemOptions(itemKey int64) (map[int16]int64, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetItemOption(itemKey int64, optionType int16, optionData int64) error {
	return nil
}

func (r *MemoryCharacterRepository) GetItemEnchant(itemKey int64) ([]*ItemEnchantData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetItemData(itemKey int64, owner int, storageType byte, storage int16, itemID int, flags int, dDate time.Time) error {
	return nil
}

func (r *MemoryCharacterRepository) MergeItems(ownerFrom, ownerTo int, keyFrom, keyTo int64, optType int16, optData int64) error {
	return nil
}

func (r *MemoryCharacterRepository) SplitAndMergeItems(keyFrom int64, ownerFrom int, optTypeFrom int16, optDataFrom int64, keyTo int64, ownerTo int, optTypeTo int16, optDataTo int64) error {
	return nil
}

func (r *MemoryCharacterRepository) UpgradeItem(owner int, itemKey int64, itemDestroy int, upOptType int16, upOptValue int64, upFailCountType int16, upFailCountValue int64) error {
	return nil
}

func (r *MemoryCharacterRepository) GetItemActionCoolTime(charNo int) ([]*ItemActionCoolTimeData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetItemActionCoolTime(charNo int, actionGroup int16, endTime time.Time) error {
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
	case "shortcut", "shortcutdata":
		opt.ShortCutData = data
	case "shortcutsize":
		opt.ShortCutSize = data
	case "game":
		opt.Game = data
	case "video":
		opt.Video = data
	case "sound":
		opt.Sound = data
	case "windowspos":
		opt.WindowsPos = data
	case "keymapping":
		opt.KeyMapping = data
	}
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
	return nil
}

func (r *MemoryCharacterRepository) GetCollectionCards(charNo int) ([]*CollectionCardData, error) {
	return nil, nil
}

func (r *MemoryCharacterRepository) SetCollectionCard(charNo int, cardID int16, count int16) error {
	return nil
}

func (r *MemoryCharacterRepository) GetTutorialCharInfo(charNo int) (byte, byte, error) {
	return 0, 0, nil
}

func (r *MemoryCharacterRepository) SetTutorialCharInfo(userNo, charNo int, state, step byte) error {
	return nil
}

func (r *MemoryCharacterRepository) GetTutorialIsSkip(userNo int) (bool, error) {
	return false, nil
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
