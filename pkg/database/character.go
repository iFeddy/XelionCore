package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrCharacterNotFound = errors.New("character not found")
	ErrCharacterExists   = errors.New("character name already exists")
	ErrSlotInUse         = errors.New("character slot already in use")
	ErrSlotOutOfRange    = errors.New("character slot out of range (0-5)")
)

// TCharacter represents the canonical tCharacter table from World00_Character.bak.
type TCharacter struct {
	CharNo   int    `json:"nCharNo"`
	UserNo   int    `json:"nUserNo"`
	CharName string `json:"sCharName"`
	Slot     byte   `json:"nSlot"`
	Class    byte   `json:"nClass"`
	Level    byte   `json:"nLevel"`
	IsDelete int    `json:"bIsDelete"` // 0 = Active, 1 = Soft Deleted
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
	CharNo         int    `json:"nCharNo"`
	QuickBar       []byte `json:"sQuickBar"`
	QuickBarState  []byte `json:"sQuickBarState"`
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
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Name      string `json:"name"`
	Slot      byte   `json:"slot"`
	Job       byte   `json:"job"`
	Level     byte   `json:"level"`
	Male      bool   `json:"male"`
	Hair      byte   `json:"hair"`
	HairColor byte   `json:"hair_color"`
	Face      byte   `json:"face"`
	Map       uint16 `json:"map"`
	XPos      int    `json:"x_pos"`
	YPos      int    `json:"y_pos"`
	CurHP     int    `json:"cur_hp"`
	CurSP     int    `json:"cur_sp"`
	TownName  string `json:"town_name"`
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

// CharacterOptions representa las opciones binarias del cliente persistidas en tCharacterOptions
// (p_Char_GetOptShortCutData, p_Char_SetOptShortCutData, p_Char_GetOptGame, etc.).
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

// CharacterRepository defines persistence operations for player characters.
type CharacterRepository interface {
	GetCharactersByAccount(accountID int) ([]*Character, error)
	GetCharacterByName(name string) (*Character, error)
	GetCharacterByID(id int) (*Character, error)
	CreateCharacter(ch *Character) error
	UpdateCharacter(ch *Character) error
	DeleteCharacter(accountID int, slot byte) error

	// Official inventory, shortcuts and skill routines
	GetCharacterItems(charNo int) ([]*Item, error)
	GetCharacterShortcuts(charNo int) ([]*Shortcut, error)
	GetCharacterSkills(charNo int) ([]*CharacterSkill, error)
	SaveCharacterItem(item *Item) error
	SaveCharacterShortcuts(charNo int, shortcuts []*Shortcut) error
	SaveCharacterSkills(charNo int, skills []*CharacterSkill) error

	// Official client options / settings routines (tCharacterOptions)
	GetCharacterOptions(charNo int) (*CharacterOptions, error)
	SaveCharacterOption(charNo int, optType string, data []byte) error

	// Official friend routines (tFriend - p_Friend_Get, p_Friend_Set, p_Friend_Del)
	GetFriends(charNo int) ([]*FriendInfo, error)
	AddFriend(charNo, friendCharNo int) error
	RemoveFriend(charNo, friendCharNo int) error
	GetFriendPoints(charNo int) (int, error)
	SetFriendPoints(charNo int, points int) error

	// Official guild routines (tGuild, tGuildMember - p_Guild_*)
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

	// Official guild academy routines (tGuildAcademy, tGuildAcademyMember - usp_GuildAcademy_*)
	CreateGuildAcademy(guildID, masterCharID int) error
	GetGuildAcademy(guildID int) (*GuildAcademyData, error)
	GetAcademyByCharID(charID int) (*GuildAcademyData, *GuildAcademyMemberData, error)
	GetAcademyMembers(guildID int) ([]*GuildAcademyMemberData, error)
	AddAcademyMember(guildID, charID int, joinLevel byte) error
	RemoveAcademyMember(guildID, charID int) error
	SetAcademyNotify(guildID int, charName, notify string) error
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

// GuildMemberData modela tGuildMember enriquecido con los datos del personaje.
type GuildMemberData struct {
	GuildID   int       `json:"guild_id"`
	CharID    int       `json:"char_id"`
	CharName  string    `json:"char_name"`
	JoinDate  time.Time `json:"join_date"`
	Grade     byte      `json:"grade"` // 0=Master, 1=SubMaster, 2=Officer, 3=Member, 4=Newbie
	Flags     uint32    `json:"flags"`
	ExpRatio  uint16    `json:"exp_ratio"`
	Intro     string    `json:"intro"`
	Class     byte      `json:"class"`
	Level     byte      `json:"level"`
	Map       string    `json:"map"`
	LastLogin time.Time `json:"last_login"`
	IsOnline  bool      `json:"is_online"`
}

// GuildAcademyData modela tGuildAcademy según el esquema oficial.
type GuildAcademyData struct {
	GuildID        int       `json:"guild_id"`
	MasterCharID   int       `json:"master_char_id"`
	MasterCharName string    `json:"master_char_name"`
	MaxMember      uint16    `json:"max_member"`
	Money          uint64    `json:"money"`
	Point          uint32    `json:"point"`
	BuffEndTime    time.Time `json:"buff_end_time"`
	NotifyDate     time.Time `json:"notify_date"`
	NotifyCharID   string    `json:"notify_char_id"`
	Notify         string    `json:"notify"`
	Intro          string    `json:"intro"`
}

// GuildAcademyMemberData modela tGuildAcademyMember con datos del aprendiz.
type GuildAcademyMemberData struct {
	GuildID   int       `json:"guild_id"`
	CharID    int       `json:"char_id"`
	CharName  string    `json:"char_name"`
	JoinDate  time.Time `json:"join_date"`
	JoinLevel byte      `json:"join_level"`
	Grade     byte      `json:"grade"`
	Class     byte      `json:"class"`
	Level     byte      `json:"level"`
	Map       string    `json:"map"`
	LastLogin time.Time `json:"last_login"`
	IsOnline  bool      `json:"is_online"`
}

// GuildListItem representa una fila resumen para la lista pública de gremios.
type GuildListItem struct {
	GuildID        uint32 `json:"guild_id"`
	GuildName      string `json:"guild_name"`
	MasterCharName string `json:"master_char_name"`
	Type           byte   `json:"type"`
	Grade          byte   `json:"grade"`
	MemberCount    uint16 `json:"member_count"`
	MaxMembers     uint16 `json:"max_members"`
}

// FriendInfo representa la información de un amigo persistido en tFriend.
type FriendInfo struct {
	CharID   int    `json:"char_id"`
	CharName string `json:"char_name"`
	Class    byte   `json:"class"`
	Level    byte   `json:"level"`
	Map      string `json:"map"`
	Flags    byte   `json:"flags"`
}

type charStore struct {
	NextID     int          `json:"next_id"`
	Characters []*Character `json:"characters"`
}

// MemoryCharacterRepository provides thread-safe character storage with JSON persistence.
type MemoryCharacterRepository struct {
	mu         sync.RWMutex
	filePath   string
	characters map[int]*Character // Key: Character.ID
	options    map[int]*CharacterOptions // Key: Character.ID
	friends    map[int][]int             // Key: Character.ID -> Friend CharIDs
	points     map[int]int               // Key: Character.ID -> Friend Points
	guilds     map[int]*GuildData        // Key: GuildID
	members    map[int][]*GuildMemberData // Key: GuildID -> Members
	academies  map[int]*GuildAcademyData // Key: GuildID
	acadMembers map[int][]*GuildAcademyMemberData // Key: GuildID -> Academy Members
	nextGuildID int
	nextID     int
}

// NewMemoryCharacterRepositoryWithFile initializes a character repository with a custom file path (or empty for in-memory).
func NewMemoryCharacterRepositoryWithFile(filePath string) *MemoryCharacterRepository {
	repo := &MemoryCharacterRepository{
		filePath:    filePath,
		characters:  make(map[int]*Character),
		options:     make(map[int]*CharacterOptions),
		friends:     make(map[int][]int),
		points:      make(map[int]int),
		guilds:      make(map[int]*GuildData),
		members:     make(map[int][]*GuildMemberData),
		academies:   make(map[int]*GuildAcademyData),
		acadMembers: make(map[int][]*GuildAcademyMemberData),
		nextGuildID: 1,
		nextID:      1,
	}
	repo.load()
	return repo
}

// NewMemoryCharacterRepository initializes a character repository persisted at data/characters.json.
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
	var store charStore
	if err := json.Unmarshal(data, &store); err == nil {
		if store.NextID > 0 {
			r.nextID = store.NextID
		}
		for _, ch := range store.Characters {
			ch.Sync()
			r.characters[ch.ID] = ch
			if ch.ID >= r.nextID {
				r.nextID = ch.ID + 1
			}
		}
	}
}

func (r *MemoryCharacterRepository) save() {
	if r.filePath == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(r.filePath), 0755)
	var list []*Character
	for _, ch := range r.characters {
		ch.Sync()
		list = append(list, ch)
	}
	store := charStore{
		NextID:     r.nextID,
		Characters: list,
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err == nil {
		_ = os.WriteFile(r.filePath, data, 0644)
	}
}

func (r *MemoryCharacterRepository) GetCharactersByAccount(accountID int) ([]*Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	var list []*Character
	for _, ch := range r.characters {
		if ch.AccountID == accountID && ch.TCharacter.IsDelete == 0 {
			copied := *ch
			list = append(list, &copied)
		}
	}
	return list, nil
}

func (r *MemoryCharacterRepository) GetCharacterByName(name string) (*Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	for _, ch := range r.characters {
		if strings.EqualFold(ch.Name, name) && ch.TCharacter.IsDelete == 0 {
			copied := *ch
			return &copied, nil
		}
	}
	return nil, ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) GetCharacterByID(id int) (*Character, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	ch, exists := r.characters[id]
	if !exists || ch.TCharacter.IsDelete != 0 {
		return nil, ErrCharacterNotFound
	}
	copied := *ch
	return &copied, nil
}

func (r *MemoryCharacterRepository) CreateCharacter(ch *Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	if ch.Slot > 5 {
		return ErrSlotOutOfRange
	}

	for _, existing := range r.characters {
		if existing.TCharacter.IsDelete == 0 {
			if strings.EqualFold(existing.Name, ch.Name) {
				return ErrCharacterExists
			}
			if existing.AccountID == ch.AccountID && existing.Slot == ch.Slot {
				return ErrSlotInUse
			}
		}
	}

	ch.ID = r.nextID
	r.nextID++
	ch.Sync()
	r.characters[ch.ID] = ch
	r.save()
	return nil
}

func (r *MemoryCharacterRepository) UpdateCharacter(ch *Character) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	existing, exists := r.characters[ch.ID]
	if !exists {
		return ErrCharacterNotFound
	}

	ch.Sync()
	*existing = *ch
	r.save()
	return nil
}

func (r *MemoryCharacterRepository) DeleteCharacter(accountID int, slot byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.load()

	for id, ch := range r.characters {
		if ch.AccountID == accountID && ch.Slot == slot {
			ch.TCharacter.IsDelete = 1
			r.characters[id] = ch
			r.save()
			return nil
		}
	}
	return ErrCharacterNotFound
}

func (r *MemoryCharacterRepository) GetCharacterItems(charNo int) ([]*Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, exists := r.characters[charNo]
	if !exists {
		return nil, ErrCharacterNotFound
	}
	res := make([]*Item, len(ch.Items))
	copy(res, ch.Items)
	return res, nil
}

func (r *MemoryCharacterRepository) GetCharacterShortcuts(charNo int) ([]*Shortcut, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, exists := r.characters[charNo]
	if !exists {
		return nil, ErrCharacterNotFound
	}
	res := make([]*Shortcut, len(ch.Shortcuts))
	copy(res, ch.Shortcuts)
	return res, nil
}

func (r *MemoryCharacterRepository) GetCharacterSkills(charNo int) ([]*CharacterSkill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, exists := r.characters[charNo]
	if !exists {
		return nil, ErrCharacterNotFound
	}
	res := make([]*CharacterSkill, len(ch.Skills))
	copy(res, ch.Skills)
	return res, nil
}

func (r *MemoryCharacterRepository) SaveCharacterItem(item *Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, exists := r.characters[item.Owner]
	if !exists {
		return ErrCharacterNotFound
	}
	ch.Items = append(ch.Items, item)
	r.save()
	return nil
}

func (r *MemoryCharacterRepository) SaveCharacterShortcuts(charNo int, shortcuts []*Shortcut) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, exists := r.characters[charNo]
	if !exists {
		return ErrCharacterNotFound
	}
	ch.Shortcuts = append(ch.Shortcuts, shortcuts...)
	r.save()
	return nil
}

func (r *MemoryCharacterRepository) SaveCharacterSkills(charNo int, skills []*CharacterSkill) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch, exists := r.characters[charNo]
	if !exists {
		return ErrCharacterNotFound
	}
	ch.Skills = append(ch.Skills, skills...)
	r.save()
	return nil
}

func (r *MemoryCharacterRepository) GetCharacterOptions(charNo int) (*CharacterOptions, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	opt, exists := r.options[charNo]
	if !exists {
		return &CharacterOptions{CharNo: charNo}, nil
	}
	copied := *opt
	return &copied, nil
}

func (r *MemoryCharacterRepository) SaveCharacterOption(charNo int, optType string, data []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	opt, exists := r.options[charNo]
	if !exists {
		opt = &CharacterOptions{CharNo: charNo}
		r.options[charNo] = opt
	}
	switch optType {
	case "shortcut_data":
		opt.ShortCutData = data
	case "shortcut_size":
		opt.ShortCutSize = data
	case "game":
		opt.Game = data
	case "sound":
		opt.Sound = data
	case "window_pos":
		opt.WindowsPos = data
	case "key_mapping":
		opt.KeyMapping = data
	case "video":
		opt.Video = data
	}
	return nil
}

func (r *MemoryCharacterRepository) GetFriends(charNo int) ([]*FriendInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	friendIDs := r.friends[charNo]
	var list []*FriendInfo
	for _, fID := range friendIDs {
		if friendChar, exists := r.characters[fID]; exists && friendChar.TCharacter.IsDelete == 0 {
			list = append(list, &FriendInfo{
				CharID:   friendChar.ID,
				CharName: friendChar.Name,
				Class:    friendChar.Job,
				Level:    friendChar.Level,
				Map:      friendChar.TownName,
			})
		}
	}
	return list, nil
}

func (r *MemoryCharacterRepository) AddFriend(charNo, friendCharNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.friends[charNo] = append(r.friends[charNo], friendCharNo)
	r.friends[friendCharNo] = append(r.friends[friendCharNo], charNo)
	return nil
}

func (r *MemoryCharacterRepository) RemoveFriend(charNo, friendCharNo int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	remove := func(list []int, target int) []int {
		res := make([]int, 0, len(list))
		for _, v := range list {
			if v != target {
				res = append(res, v)
			}
		}
		return res
	}

	r.friends[charNo] = remove(r.friends[charNo], friendCharNo)
	r.friends[friendCharNo] = remove(r.friends[friendCharNo], charNo)
	return nil
}

func (r *MemoryCharacterRepository) GetFriendPoints(charNo int) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.points[charNo], nil
}

func (r *MemoryCharacterRepository) SetFriendPoints(charNo int, points int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.points[charNo] = points
	return nil
}

func (r *MemoryCharacterRepository) CreateGuild(masterCharID int, name, password string, guildType byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Chequear si existe nombre
	for _, g := range r.guilds {
		if g.Name == name {
			return -1, fmt.Errorf("guild name already exists")
		}
	}

	// Chequear si el char ya pertenece a un gremio
	for _, mList := range r.members {
		for _, m := range mList {
			if m.CharID == masterCharID {
				return -3, fmt.Errorf("character already belongs to a guild")
			}
		}
	}

	guildID := r.nextGuildID
	r.nextGuildID++

	now := time.Now()
	g := &GuildData{
		ID:          guildID,
		Name:        name,
		Password:    password,
		MaxMembers:  30,
		Money:       0,
		Type:        guildType,
		Grade:       1,
		CreateDate:  now,
		NotifyDate:  now,
		DismissDate: now,
	}
	r.guilds[guildID] = g

	// Agregar al creador como Guild Master (Grade 0)
	ch := r.characters[masterCharID]
	charName := ""
	var classID, level byte
	mapName := "Rou"
	if ch != nil {
		charName = ch.Name
		classID = ch.Job
		level = ch.Level
		mapName = strconv.Itoa(int(ch.Map))
	}

	masterMember := &GuildMemberData{
		GuildID:   guildID,
		CharID:    masterCharID,
		CharName:  charName,
		JoinDate:  now,
		Grade:     0,
		Class:     classID,
		Level:     level,
		Map:       mapName,
		LastLogin: now,
		IsOnline:  true,
	}
	r.members[guildID] = []*GuildMemberData{masterMember}

	return guildID, nil
}

func (r *MemoryCharacterRepository) GetGuild(guildID int) (*GuildData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.guilds[guildID]
	if !ok {
		return nil, fmt.Errorf("guild not found")
	}
	return g, nil
}

func (r *MemoryCharacterRepository) GetGuildByCharID(charID int) (*GuildData, *GuildMemberData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for gID, mList := range r.members {
		for _, m := range mList {
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
	mList, ok := r.members[guildID]
	if !ok {
		return []*GuildMemberData{}, nil
	}
	res := make([]*GuildMemberData, len(mList))
	copy(res, mList)
	return res, nil
}

func (r *MemoryCharacterRepository) GetAllGuilds() ([]*GuildListItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*GuildListItem, 0, len(r.guilds))
	for _, g := range r.guilds {
		if g.DismissStatus != 0 {
			continue
		}
		masterName := ""
		mList := r.members[g.ID]
		for _, m := range mList {
			if m.Grade == 0 {
				masterName = m.CharName
				break
			}
		}
		res = append(res, &GuildListItem{
			GuildID:        uint32(g.ID),
			GuildName:      g.Name,
			MasterCharName: masterName,
			Type:           g.Type,
			Grade:          g.Grade,
			MemberCount:    uint16(len(mList)),
			MaxMembers:     g.MaxMembers,
		})
	}
	return res, nil
}

func (r *MemoryCharacterRepository) AddGuildMember(guildID, charID int, grade byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch := r.characters[charID]
	charName := ""
	var classID, level byte
	mapName := "Rou"
	if ch != nil {
		charName = ch.Name
		classID = ch.Job
		level = ch.Level
		mapName = strconv.Itoa(int(ch.Map))
	}

	m := &GuildMemberData{
		GuildID:   guildID,
		CharID:    charID,
		CharName:  charName,
		JoinDate:  time.Now(),
		Grade:     grade,
		Class:     classID,
		Level:     level,
		Map:       mapName,
		LastLogin: time.Now(),
		IsOnline:  true,
	}
	r.members[guildID] = append(r.members[guildID], m)
	return nil
}

func (r *MemoryCharacterRepository) RemoveGuildMember(guildID, charID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	mList := r.members[guildID]
	res := make([]*GuildMemberData, 0, len(mList))
	for _, m := range mList {
		if m.CharID != charID {
			res = append(res, m)
		}
	}
	r.members[guildID] = res
	return nil
}

func (r *MemoryCharacterRepository) SetGuildMemberGrade(guildID, charID int, grade byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.members[guildID] {
		if m.CharID == charID {
			m.Grade = grade
			return nil
		}
	}
	return fmt.Errorf("member not found in guild")
}

func (r *MemoryCharacterRepository) SetGuildNotify(guildID int, charName, notify string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.guilds[guildID]
	if !ok {
		return fmt.Errorf("guild not found")
	}
	g.NotifyDate = time.Now()
	g.NotifyCharID = charName
	g.Notify = notify
	return nil
}

func (r *MemoryCharacterRepository) DismissGuild(guildID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	g, ok := r.guilds[guildID]
	if !ok {
		return fmt.Errorf("guild not found")
	}
	g.DismissStatus = 1
	g.DismissDate = time.Now()
	delete(r.members, guildID)
	return nil
}

func (r *MemoryCharacterRepository) CreateGuildAcademy(guildID, masterCharID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch := r.characters[masterCharID]
	masterName := ""
	if ch != nil {
		masterName = ch.Name
	}

	now := time.Now()
	r.academies[guildID] = &GuildAcademyData{
		GuildID:        guildID,
		MasterCharID:   masterCharID,
		MasterCharName: masterName,
		MaxMember:      50,
		NotifyDate:     now,
		BuffEndTime:    now,
	}
	return nil
}

func (r *MemoryCharacterRepository) GetGuildAcademy(guildID int) (*GuildAcademyData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.academies[guildID]
	if !ok {
		return nil, fmt.Errorf("academy not found")
	}
	return a, nil
}

func (r *MemoryCharacterRepository) GetAcademyByCharID(charID int) (*GuildAcademyData, *GuildAcademyMemberData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for gID, mList := range r.acadMembers {
		for _, m := range mList {
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
	mList, ok := r.acadMembers[guildID]
	if !ok {
		return []*GuildAcademyMemberData{}, nil
	}
	res := make([]*GuildAcademyMemberData, len(mList))
	copy(res, mList)
	return res, nil
}

func (r *MemoryCharacterRepository) AddAcademyMember(guildID, charID int, joinLevel byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ch := r.characters[charID]
	charName := ""
	var classID, level byte
	mapName := "Rou"
	if ch != nil {
		charName = ch.Name
		classID = ch.Job
		level = ch.Level
		mapName = strconv.Itoa(int(ch.Map))
	}

	m := &GuildAcademyMemberData{
		GuildID:   guildID,
		CharID:    charID,
		CharName:  charName,
		JoinDate:  time.Now(),
		JoinLevel: joinLevel,
		Grade:     0,
		Class:     classID,
		Level:     level,
		Map:       mapName,
		LastLogin: time.Now(),
		IsOnline:  true,
	}
	r.acadMembers[guildID] = append(r.acadMembers[guildID], m)
	return nil
}

func (r *MemoryCharacterRepository) RemoveAcademyMember(guildID, charID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	mList := r.acadMembers[guildID]
	res := make([]*GuildAcademyMemberData, 0, len(mList))
	for _, m := range mList {
		if m.CharID != charID {
			res = append(res, m)
		}
	}
	r.acadMembers[guildID] = res
	return nil
}

func (r *MemoryCharacterRepository) SetAcademyNotify(guildID int, charName, notify string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	a, ok := r.academies[guildID]
	if !ok {
		return fmt.Errorf("academy not found")
	}
	a.NotifyDate = time.Now()
	a.NotifyCharID = charName
	a.Notify = notify
	return nil
}
