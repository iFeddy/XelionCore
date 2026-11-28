package zone

import (
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ItemSlot identifiers matching official server equipment slots.
const (
	Slot_None          byte = 0
	Slot_Helm          byte = 1
	Slot_CostumeHelm   byte = 2
	Slot_Necklace      byte = 6
	Slot_Armor         byte = 7
	Slot_CostumeArmor  byte = 8
	Slot_Wing          byte = 9
	Slot_Weapon2       byte = 10
	Slot_CostumeShield byte = 11
	Slot_Weapon        byte = 12
	Slot_CostumeWeapon byte = 13
	Slot_Ring          byte = 15
	Slot_Ring2         byte = 16
	Slot_Tail          byte = 18
	Slot_Pants         byte = 19
	Slot_CostumePants  byte = 20
	Slot_Boots         byte = 21
	Slot_CostumeBoots  byte = 22
	Slot_Earings       byte = 23
	Slot_Glasses       byte = 26
	Slot_Pet           byte = 28
)

// JobToClassName mapea el ID oficial de las 27 clases del juego a su definición en data/World/ClassParam/.
func JobToClassName(job byte) string {
	switch job {
	case 1:
		return "Fighter"
	case 2:
		return "CleverFighter"
	case 3:
		return "Warrior"
	case 4:
		return "Gladiator"
	case 5:
		return "Knight"
	case 6:
		return "Cleric"
	case 7:
		return "HighCleric"
	case 8:
		return "Paladin"
	case 9:
		return "HolyKnight"
	case 10:
		return "Guardian"
	case 11:
		return "Archer"
	case 12:
		return "HawkArcher"
	case 13:
		return "Scout"
	case 14:
		return "SharpShooter"
	case 15:
		return "Ranger"
	case 16:
		return "Mage"
	case 17:
		return "WizMage"
	case 18:
		return "Enchanter"
	case 19:
		return "Warrock"
	case 20:
		return "Wizard"
	case 21:
		return "Joker"
	case 22:
		return "Chaser"
	case 23:
		return "Cruel"
	case 24:
		return "Closer"
	case 25:
		return "Assassin"
	case 26:
		return "Sentinel"
	case 27:
		return "Savior"
	default:
		return "Fighter"
	}
}

// PendingSkillCast representa una habilidad con CastTime > 0 en proceso de canalización en el tick de 100ms.
type PendingSkillCast struct {
	SkillID   uint16
	IsField   bool
	TargetObj uint16
	FldX      int32
	FldY      int32
	ExecuteAt time.Time
	IsMoving  bool
	HPCost    int
	SPCost    int
	EmpDmg    uint8
	EmpKeep   uint8
	EmpCool   uint8
	SkillDef  *SkillDef
}

// ZonePlayer representa a un jugador autenticado y activo en una instancia de ZoneServer.
type ZonePlayer struct {
	Client        *network.Client
	Char          *database.Character
	AccountID     int
	Username      string
	Admin         uint8
	CurrentMap    *ZoneMap
	MapObjectID   uint16
	WorldNo       byte
	State         byte // 1=Normal, 2=Player, 3=Dead, 4=Resting, 5=Vendor, 6=Mount
	Rotation      byte
	StonesHP      int16
	StonesSP      int16
	MaxStonesHP   int16
	MaxStonesSP   int16
	HPChangeOrder uint16
	MaxHP         int
	MaxSP         int
	Str           byte
	End           byte
	Dex           byte
	Int           byte
	Spr           byte
	LevelParam    classparam.LevelParam
	EquipMinWC    uint32
	EquipMaxWC    uint32
	EquipAC       uint32
	EquipMinMA    uint32
	EquipMaxMA    uint32
	EquipMR       uint32
	EquipTH       uint32
	EquipTB       uint32
	// Bonificadores dinámicos acumulados por AbStates activos (SubAbstate.shn SAA_* oficial de Zone.pdb)
	AbWC           int
	AbWCRate       int
	AbAC           int
	AbACRate       int
	AbMA           int
	AbMARate       int
	AbMR           int
	AbMRRate       int
	AbTH           int
	AbTHRate       int
	AbTB           int
	AbTBRate       int
	AbMaxHP        int
	AbMaxHPRate    int
	AbMaxSP        int
	AbMaxSPRate    int
	AbCritRate     int
	AbDamageRate   int
	AbDmgDownRate  int
	AbReflectRate  int
	AbHealRate     int
	AbCastTimeMod  int
	AbMinHPFloor   int
	LoginTime     time.Time
	LastWarpTime  time.Time

	// Campos de subsistemas de la Fase 11
	HideExceptMe bool
	TentID       uint16
	MoverHandle  uint16
	MoverSlot    uint8
	MoverSpeed   uint16
	WalkSpeed    uint16
	RunSpeed     uint16
	PrevExp      uint64
	NextExp      uint64
	GatheringObj uint16

	// Campos de subsistemas de la Fase 12 (Combate, Skills, AbState, PK)
	TargetHandle          uint16
	SecondaryTargetHandle uint16
	IsAutoAttacking       bool
	LastAttackTime        time.Time
	CastingSkillID        uint16
	CastingTargetHandle   uint16
	PendingCast           *PendingSkillCast
	CastingFieldX         int32
	CastingFieldY         int32
	CastingEndTime        time.Time
	NextCastIndex         uint16
	LearnedSkills         map[uint16]*database.CharacterSkill
	ToggleSkills          map[uint16]bool
	SkillCooldowns        map[uint16]time.Time
	SkillGroupCooldowns   map[int]time.Time // Cooldowns compartidos por DlyGroupNum (ActiveSkill.shn)
	MoverSkillCooldowns   map[uint16]time.Time
	PendingRevive         bool
	PendingReviveRate     uint16
	ActiveStates          map[uint16]*ActiveAbState
	PKCount               int
	PKYellowTime          int
	LastSaveAt            time.Time
	RestTickCount         int
	HouseHPTickMs         uint32
	HouseSPTickMs         uint32

	// Campos de subsistemas de la Fase 13 (AutoPick, Títulos, RestExp)
	IsAutoPick          bool
	CurrentTitleType    byte
	CurrentTitleElement byte
	CurrentTitleMobID   uint16
	RestExpBonusRate    int

	// Campos de subsistemas de la Fase 14 (Economía, Comercio, Booth, Forja & Crafting)
	OpenNPCMobID         uint16
	OpenNPCHandle        uint16
	OpenNPCX             int32
	OpenNPCY             int32
	OpenMultyLinkNPC     uint16
	OpenMultyLinkOptions []string
	ChargedInvenOpened   bool
	AccountStorageOpen   bool
	AcademyStorageOpen   bool
	RepurchaseItems      []*RepurchaseEntry
	PendingRandomOpt     *PendingRandomOption
	ActiveTrade          *SafeTradeSession
	PendingTradeFrom     uint16
	ActiveBooth          *PlayerBooth
	MinimonNormalSlots   [12]uint16
	MinimonChargedSlots  [12]uint16
	ActiveProduce        *ProduceSession

	// Campos de subsistemas de la Fase 15 (Quests, Instancias/KQ/GT/Wedding, MiniHouse, Pets y Casino)
	PendingQuestRewardSelect uint16
	LastQuestGiveUpAt        time.Time
	PendingMenuNPC           uint16
	PendingMenuCallbacks     []string
	ActivePet                *ActivePetState
	InsideMiniHouseOwner     int
	MiniHouseOpen            bool
	MiniHouseTitle           string
	MiniHousePassword        string
	MiniHouseMaxGuests       uint8
	CasinoCoins              uint64
	CasinoDiceJoined         bool
	KQTeam                   byte
	KQRewardIdx              string
	KQReturnMap              string
	KQReturnX                int
	KQReturnY                int
	KQRewardOrderSeq         uint16
	GTObserving              bool

	// Campos de subsistemas de la Fase 16 (GM, Sesión, Beauty, GuildStorage, ActionMode)
	IsGMHidden       bool
	IsGMImmortal     bool
	ActionMode       byte
	BeautyShopOpen   bool
	GuildStorageOpen bool
	LogoutPending    bool
	LogoutType       byte
}

// ActivePetState mantiene el estado vivo en mapa de la mascota invocada por un jugador.
type ActivePetState struct {
	Handle     uint16
	PetID      int
	ItemKey    int64
	PetType    uint16
	Mind       uint8
	Hunger     uint8
	Stress     uint8
	Tendency   int16
	SummonedAt time.Time
}

// RepurchaseEntry guarda un ítem vendido recientemente para PROTO_NC_ITEM_REPURCHASE_REQ (0x3059).
type RepurchaseEntry struct {
	Handle  uint16
	ItemID  uint16
	Lot     uint32
	Price   uint64
	Options map[int16]int64
}

// PendingRandomOption almacena la tirada de opciones aleatorias pendiente de aceptar (0x3053 -> 0x3061).
type PendingRandomOption struct {
	ItemKey   int64
	InvenSlot uint8
	Options   []ItemOptionPair
}

// ItemOptionPair representa un par tipo/valor de opción aleatoria en ItemOptionStorage (3B wire).
type ItemOptionPair struct {
	Type  uint8
	Value uint16
}

// ProduceSession representa una sesión activa de crafteo (`0x2035` PROTO_NC_ACT_PRODUCE_CAST_REQ).
type ProduceSession struct {
	ProductID  uint16
	StartedAt  time.Time
	CastMillis int
}

// GetHandle implementa GridEntity para ZonePlayer.
func (p *ZonePlayer) GetHandle() uint16 {
	return p.MapObjectID
}

// GetPosition implementa GridEntity para ZonePlayer.
func (p *ZonePlayer) GetPosition() (int, int) {
	if p.Char == nil {
		return 0, 0
	}
	return p.Char.XPos, p.Char.YPos
}

// GetEntityType implementa GridEntity para ZonePlayer.
func (p *ZonePlayer) GetEntityType() GridEntityType {
	return GridEntityPlayer
}

// GetClient implementa GridEntity para ZonePlayer.
func (p *ZonePlayer) GetClient() *network.Client {
	return p.Client
}

// NewZonePlayer inicializa un ZonePlayer cargando sus stats oficiales desde ClassParam y FreeStatTable.
func NewZonePlayer(client *network.Client, ch *database.Character, classStore *classparam.Store, combatStore *combat.Store) *ZonePlayer {
	className := JobToClassName(ch.Job)
	lvl := int(ch.Level)
	if lvl <= 0 {
		lvl = 1
	}

	var lp classparam.LevelParam
	if classStore != nil {
		if p, ok := classStore.GetParam(className, lvl); ok && p != nil {
			lp = *p
		}
	}

	// Obtener bonos oficiales de FreeStatTable.json
	var freeBonus combat.FreeStatBonusResult
	if combatStore != nil {
		freeBonus = combatStore.GetFreeStatBonus(
			ch.TCharacter.Strength,
			ch.TCharacter.Intelligence,
			ch.TCharacter.Dexterity,
			ch.TCharacter.Constitute,
			ch.TCharacter.MentalPower,
		)
	}

	// Calcular el MaxHP real: Vida base de la clase + bonus por Constitución distribuida de FreeStatTable
	maxHP := lp.MaxHP + freeBonus.MaxHP
	if maxHP <= 0 {
		maxHP = 52
	}

	// Calcular el MaxSP real: Maná base de la clase + bonus por MentalPower distribuido de FreeStatTable
	maxSP := lp.MaxSP + freeBonus.MaxSP
	if maxSP <= 0 {
		maxSP = 10
	}

	// Obtener umbrales de experiencia real de ExpTable.json
	var prevExp, nextExp int64
	if combatStore != nil {
		prevExp, nextExp = combatStore.GetExpProgression(lvl)
	} else {
		prevExp = 0
		nextExp = 10
	}

	curHP := ch.CurHP
	if curHP <= 0 {
		curHP = maxHP
		ch.CurHP = curHP
	}
	curSP := ch.CurSP
	if curSP <= 0 {
		curSP = maxSP
		ch.CurSP = curSP
	}

	maxStonesHP := int16(lp.MAXSoulHP)
	if maxStonesHP <= 0 {
		maxStonesHP = 15
	}
	maxStonesSP := int16(lp.MAXSoulSP)
	if maxStonesSP <= 0 {
		maxStonesSP = 7
	}

	stonesHP := ch.TCharacterState.HPStones
	if stonesHP <= 0 {
		stonesHP = maxStonesHP
		ch.TCharacterState.HPStones = stonesHP
	}
	stonesSP := ch.TCharacterState.SPStones
	if stonesSP <= 0 {
		stonesSP = maxStonesSP
		ch.TCharacterState.SPStones = stonesSP
	}

	walkSpeed := uint16(33)
	runSpeed := uint16(127)
	if combatStore != nil {
		common := combatStore.Common()
		if common.WalkSpeed > 0 {
			walkSpeed = uint16(common.WalkSpeed)
		}
		if common.RunSpeed > 0 {
			runSpeed = uint16(common.RunSpeed)
		}
	}

	return &ZonePlayer{
		Client:       client,
		Char:         ch,
		State:        1,    // Normal
		Rotation:     0x55, // 85 grados
		LoginTime:    time.Now(),
		StonesHP:     stonesHP,
		StonesSP:     stonesSP,
		MaxStonesHP:  maxStonesHP,
		MaxStonesSP:  maxStonesSP,
		MaxHP:        maxHP,
		MaxSP:        maxSP,
		WalkSpeed:    walkSpeed,
		RunSpeed:     runSpeed,
		PrevExp:      uint64(prevExp),
		NextExp:      uint64(nextExp),
		Str:          byte(ch.TCharacter.Strength),
		End:          byte(ch.TCharacter.Constitute),
		Dex:          byte(ch.TCharacter.Dexterity),
		Int:          byte(ch.TCharacter.Intelligence),
		Spr:          byte(ch.TCharacter.MentalPower),
		LevelParam:   lp,
		ActiveStates: make(map[uint16]*ActiveAbState),
	}
}

// GetMoveSpeed devuelve la velocidad actual de caminata o corrida según Zone.exe (so_WalkSpeed / so_RunSpeed).
func (p *ZonePlayer) GetMoveSpeed(walk bool) uint16 {
	if p != nil && p.MoverHandle > 0 && p.MoverSpeed > 0 {
		return p.MoverSpeed
	}
	if walk {
		if p != nil && p.WalkSpeed > 0 {
			return p.WalkSpeed
		}
		return 33
	}
	if p != nil && p.RunSpeed > 0 {
		return p.RunSpeed
	}
	return 127
}

// RecalculateEquipmentStats suma los atributos oficiales de los ítems equipados desde ItemInfo.shn (so_RecalcEquipParam).
func (p *ZonePlayer) RecalculateEquipmentStats(itemStore *item.Store) {
	p.EquipMinWC = 0
	p.EquipMaxWC = 0
	p.EquipAC = 0
	p.EquipMinMA = 0
	p.EquipMaxMA = 0
	p.EquipMR = 0
	p.EquipTH = 0
	p.EquipTB = 0

	if p.Char == nil || p.Char.Equips == nil || itemStore == nil {
		return
	}

	for _, itemID := range p.Char.Equips {
		if itemID == 0 || itemID == 0xFFFF {
			continue
		}
		if info, ok := itemStore.GetItemInfo(itemID); ok && info != nil {
			p.EquipMinWC += info.MinWC
			p.EquipMaxWC += info.MaxWC
			p.EquipAC += info.AC
			p.EquipMinMA += info.MinMA
			p.EquipMaxMA += info.MaxMA
			p.EquipMR += info.MR
			p.EquipTH += info.TH
			p.EquipTB += info.TB
		}
	}
}

// WriteLook serializa la apariencia visual del avatar según PROTO_AVATAR_SHAPE_INFO (4 bytes).
func (p *ZonePlayer) WriteLook(pw *protocol.PacketWriter) {
	var genderBit byte = 0
	if p.Char.Male {
		genderBit = 1
	}
	val := byte((1 & 0x03) | ((p.Char.Job & 0x1F) << 2) | ((genderBit & 0x01) << 7))
	pw.WriteByte(val)
	pw.WriteByte(p.Char.Hair)
	pw.WriteByte(p.Char.HairColor)
	pw.WriteByte(p.Char.Face)
}

// WriteEquipment serializa el equipamiento según PROTO_EQUIPMENT (43 bytes oficiales en PDB).
// Empaqueta los niveles de mejora (0..15) de las 5 partes visuales en 3 bytes contiguos.
func (p *ZonePlayer) WriteEquipment(pw *protocol.PacketWriter) {
	pw.WriteUint16(p.Char.GetEquip(Slot_Helm))          // Equ_Head
	pw.WriteUint16(0)                                   // Equ_Mouth
	pw.WriteUint16(p.Char.GetEquip(Slot_Weapon))        // Equ_RightHand
	pw.WriteUint16(p.Char.GetEquip(Slot_Armor))         // Equ_Body
	pw.WriteUint16(p.Char.GetEquip(Slot_Weapon2))       // Equ_LeftHand (Shield/Dual)
	pw.WriteUint16(p.Char.GetEquip(Slot_Pants))         // Equ_Pant
	pw.WriteUint16(p.Char.GetEquip(Slot_Boots))         // Equ_Boot
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumeBoots))  // Equ_AccBoot
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumePants))  // Equ_AccPant
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumeArmor))  // Equ_AccBody
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumeHelm))   // Equ_AccHeadA
	pw.WriteUint16(0)                                   // Equ_MiniMon_R
	pw.WriteUint16(p.Char.GetEquip(Slot_Glasses))       // Equ_Eye
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumeShield)) // Equ_AccLeftHand
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumeWeapon)) // Equ_AccRightHand
	pw.WriteUint16(p.Char.GetEquip(Slot_Wing))          // Equ_AccBack
	pw.WriteUint16(0)                                   // Equ_CosEff
	pw.WriteUint16(p.Char.GetEquip(Slot_Tail))          // Equ_AccHip
	pw.WriteUint16(p.Char.GetEquip(Slot_Pet))           // Equ_Minimon
	pw.WriteUint16(0)                                   // Equ_AccShield

	// Serialización exacta de los 3 bytes de upgrade visual de Zone.pdb
	var rh, lh, body, pants, boots byte
	if p.Char != nil && p.Char.EquipUpgrades != nil {
		rh = p.Char.EquipUpgrades[Slot_Weapon] & 0x0F
		lh = p.Char.EquipUpgrades[Slot_Weapon2] & 0x0F
		body = p.Char.EquipUpgrades[Slot_Armor] & 0x0F
		pants = p.Char.EquipUpgrades[Slot_Pants] & 0x0F
		boots = p.Char.EquipUpgrades[Slot_Boots] & 0x0F
	}
	pw.WriteBytes([]byte{
		rh | (lh << 4),
		body | (pants << 4),
		boots,
	})
}

// WriteDetailedInfo serializa PROTO_NC_CHAR_BASE_CMD (105 bytes oficiales en PDB).
func (p *ZonePlayer) WriteDetailedInfo(pw *protocol.PacketWriter) {
	pw.WriteUint32(uint32(p.Char.ID))                    // chrregnum
	pw.WriteString(p.Char.Name, 20)                      // charid (Name5: 20 bytes)
	pw.WriteByte(p.Char.Slot)                            // slotno
	pw.WriteByte(p.Char.Level)                           // Level
	pw.WriteUint64(uint64(p.Char.TCharacterState.Exp))   // Experience
	pw.WriteUint16(0)                                    // CurPwrStone
	pw.WriteUint16(0)                                    // CurGrdStone
	pw.WriteUint16(uint16(p.StonesHP))                   // CurHPStone
	pw.WriteUint16(uint16(p.StonesSP))                   // CurSPStone
	pw.WriteUint32(uint32(p.Char.CurHP))                 // CurHP
	pw.WriteUint32(uint32(p.Char.CurSP))                 // CurSP
	pw.WriteUint32(0)                                    // CurLP
	pw.WriteUint32(uint32(p.Char.TCharacterState.Fame))  // fame
	pw.WriteUint64(uint64(p.Char.TCharacterState.Money)) // Cen (money)

	// logininfo: PROTO_NC_CHAR_BASE_CMD::LoginLocation (21 bytes)
	shortName := p.Char.TownName
	if p.CurrentMap != nil && p.CurrentMap.MapID != "" {
		shortName = p.CurrentMap.MapID
	}
	pw.WriteString(shortName, 12)       // currentmap (Name3: 12 bytes)
	pw.WriteUint32(uint32(p.Char.XPos)) // currentcoord.xy.x
	pw.WriteUint32(uint32(p.Char.YPos)) // currentcoord.xy.y
	pw.WriteByte(p.Rotation)            // currentcoord.dir

	// statdistribute: CHARSTATDISTSTR (6 bytes)
	pw.WriteByte(p.Str)                                     // Strength
	pw.WriteByte(p.End)                                     // Constitute
	pw.WriteByte(p.Dex)                                     // Dexterity
	pw.WriteByte(p.Int)                                     // Intelligence
	pw.WriteByte(p.Spr)                                     // MentalPower
	pw.WriteByte(byte(p.Char.TCharacter.RedistributePoint)) // RedistributePoint real

	pw.WriteByte(0)                 // pkyellowtime
	pw.WriteUint32(0)               // pkcount
	pw.WriteUint16(0)               // prisonmin
	pw.WriteByte(p.Char.AdminLevel) // adminlevel real
	pw.Fill(4, 0)                   // flags (4 bytes)
}

// WriteDetailedInfoExtra serializa PROTO_NC_MAP_LOGIN_ACK (242 bytes oficiales en PDB).
func (p *ZonePlayer) WriteDetailedInfoExtra(pw *protocol.PacketWriter, levelUp bool) {
	if !levelUp {
		pw.WriteUint16(p.MapObjectID) // charhandle (2 bytes)
	}

	// CHAR_PARAMETER_DATA (232 bytes)
	pw.WriteUint64(p.PrevExp)                                      // PrevExp real
	pw.WriteUint64(p.NextExp)                                      // NextExp real
	pw.WriteUint32(uint32(p.LevelParam.Strength))                  // Strength.base
	pw.WriteUint32(uint32(p.LevelParam.Strength + int(p.Str)))     // Strength.change
	pw.WriteUint32(uint32(p.LevelParam.Constitution))              // Constitute.base
	pw.WriteUint32(uint32(p.LevelParam.Constitution + int(p.End))) // Constitute.change
	pw.WriteUint32(uint32(p.LevelParam.Dexterity))                 // Dexterity.base
	pw.WriteUint32(uint32(p.LevelParam.Dexterity + int(p.Dex)))    // Dexterity.change
	pw.WriteUint32(uint32(p.LevelParam.Intelligence))              // Intelligence.base
	pw.WriteUint32(uint32(p.LevelParam.Intelligence + int(p.Int))) // Intelligence.change
	pw.WriteUint32(0)                                              // Wizdom.base (no usado en so_fill_CHAR_PARAMETER_DATA)
	pw.WriteUint32(0)                                              // Wizdom.change
	pw.WriteUint32(uint32(p.LevelParam.MentalPower))               // MentalPower.base
	pw.WriteUint32(uint32(p.LevelParam.MentalPower + int(p.Spr)))  // MentalPower.change

	// WClow..MB: c_Storepure escribe 0 en todos los base (slots vtable CharClass +0x00..+0x1C = 0x449600);
	// total refleja so_RecalcEquipParam/Abstate/Passive/Last (so_fill_CHAR_PARAMETER_DATA @ 0x005A2E00).
	pw.WriteUint32(0)            // WClow.base
	pw.WriteUint32(p.EquipMinWC) // WClow.change
	pw.WriteUint32(0)            // WChigh.base
	pw.WriteUint32(p.EquipMaxWC) // WChigh.change
	pw.WriteUint32(0)            // AC.base
	pw.WriteUint32(p.EquipAC)    // AC.change
	pw.WriteUint32(0)            // TH.base
	pw.WriteUint32(p.EquipTH)    // TH.change
	pw.WriteUint32(0)            // TB.base
	pw.WriteUint32(p.EquipTB)    // TB.change
	pw.WriteUint32(0)            // MAlow.base
	pw.WriteUint32(p.EquipMinMA) // MAlow.change
	pw.WriteUint32(0)            // MAhigh.base
	pw.WriteUint32(p.EquipMaxMA) // MAhigh.change
	pw.WriteUint32(0)            // MR.base
	pw.WriteUint32(p.EquipMR)    // MR.change
	pw.WriteUint32(0)            // MH.base
	pw.WriteUint32(0)            // MH.change
	pw.WriteUint32(0)            // MB.base
	pw.WriteUint32(0)            // MB.change

	pw.WriteUint32(uint32(p.MaxHP))       // MaxHp
	pw.WriteUint32(uint32(p.MaxSP))       // MaxSp
	pw.WriteUint32(0)                     // MaxLp
	pw.WriteUint32(0)                     // MaxAp
	pw.WriteUint32(uint32(p.MaxStonesHP)) // MaxHPStone
	pw.WriteUint32(uint32(p.MaxStonesSP)) // MaxSPStone

	// PwrStone (16 bytes), GrdStone (16 bytes), Resistances (32 bytes): puestos a 0 en so_fill_CHAR_PARAMETER_DATA (+0xA8..+0xE7)
	pw.WriteBytes(make([]byte, 64))

	// logincoord: SHINE_XY_TYPE (8 bytes)
	if !levelUp {
		pw.WriteUint32(uint32(p.Char.XPos))
		pw.WriteUint32(uint32(p.Char.YPos))
	}
}

// Funcion para incrementar el contador de cambio de HP que lleva el paquete/cliente de hpchange
func (p *ZonePlayer) NextHPChangeOrder() uint16 {
	p.HPChangeOrder++
	if p.HPChangeOrder == 0 {
		p.HPChangeOrder = 1
	}
	return p.HPChangeOrder
}

// WriteCharacterDisplay serializa PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD (235 bytes exactos en PDB).
// Alineado con Zone.exe 0BriefInformationCharacter (0x00549090), sp_MakeTentBriefInfo (0x0056E720) y bia_Fill.
func (p *ZonePlayer) WriteCharacterDisplay(pw *protocol.PacketWriter) {
	var polymorphID uint16
	for _, st := range p.ActiveStates {
		if st != nil && st.PolymorphID > 0 {
			polymorphID = st.PolymorphID
			break
		}
	}
	abBit := p.GetAbStateBit()

	pw.WriteUint16(p.MapObjectID)     // 0x000 (0): unsigned short handle (2)
	pw.WriteString(p.Char.Name, 20)   // 0x002 (2): Name5 charid (20)
	pw.WriteInt32(int32(p.Char.XPos)) // 0x016 (22): coord.x (4)
	pw.WriteInt32(int32(p.Char.YPos)) // coord.y (4)
	pw.WriteByte(p.Rotation)          // coord.dir (1)
	pw.WriteByte(p.State)             // 0x01f (31): mode (1)
	pw.WriteByte(p.Char.Job)          // 0x020 (32): chrclass (1)
	p.WriteLook(pw)                   // 0x021 (33): PROTO_AVATAR_SHAPE_INFO (4)
	// 0x025 (37): shapedata union (45 bytes: CHARBRIEFINFO_NOTCAMP / CHARBRIEFINFO_CAMP / CHARBRIEFINFO_RIDE)
	if p.State == 4 && p.TentID > 0 {
		pw.WriteUint16(p.TentID)   // tentID (2)
		pw.WriteUint32(0xFFFFFFFF) // default camp placeholders (sp_MakeTentBriefInfo 0x0056E78C)
		pw.WriteUint32(0xFFFFFFFF)
		pw.WriteUint16(0xFFFF)
		pw.Fill(33, 0)
	} else {
		p.WriteEquipment(pw) // 43 bytes
		pw.WriteUint16(0)    // 2 bytes pad
	}
	pw.WriteUint16(polymorphID)         // 0x052 (82): polymorph (2)
	pw.Fill(3, 0)                       // 0x054 (84): STOPEMOTICON_DESCRIPT emoticon (3)
	pw.WriteByte(p.CurrentTitleType)    // 0x057 (87): CHARTITLE_BRIEFINFO.Type (1)
	pw.WriteByte(p.CurrentTitleElement) // 0x058 (88): CHARTITLE_BRIEFINFO.ElementNo (1)
	pw.WriteUint16(p.CurrentTitleMobID) // 0x059 (89): CHARTITLE_BRIEFINFO.MobID (2)
	pw.WriteBytes(abBit[:])             // 0x05b (91): ABNORMAL_STATE_BIT abstatebit (99)
	pw.WriteUint32(0)                   // 0x0be (190): myguild (4)
	pw.WriteByte(2)                     // 0x0c2 (194): type (1) = 2 en 0BriefInformationCharacter (0x005490F5)
	pw.WriteByte(0)                     // 0x0c3 (195): isGuildAcademyMember (1)
	if p.IsAutoPick {                   // 0x0c4 (196): IsAutoPick (1)
		pw.WriteByte(1)
	} else {
		pw.WriteByte(0)
	}
	pw.WriteByte(p.Char.Level)    // 0x0c5 (197): Level (1)
	pw.WriteString("", 32)        // 0x0c6 (198): sAnimation (32)
	pw.WriteUint16(p.MoverHandle) // 0x0e6 (230): nMoverHnd (2)
	pw.WriteByte(p.MoverSlot)     // 0x0e8 (232): nMoverSlot (1)
	pw.WriteByte(p.KQTeam)        // 0x0e9 (233): nKQTeamType (1)
	pw.WriteByte(0)               // 0x0ea (234): IsUseItemMinimon (1)
}

// SyncLearnedSkills sincroniza el caché en memoria de habilidades aprendidas del jugador.
func (p *ZonePlayer) SyncLearnedSkills(skills []*database.CharacterSkill) {
	if p == nil {
		return
	}
	p.LearnedSkills = make(map[uint16]*database.CharacterSkill, len(skills))
	for _, sk := range skills {
		if sk != nil {
			cp := *sk
			p.LearnedSkills[sk.SkillNo] = &cp
		}
	}
}

// HasSkill verifica si el jugador tiene aprendida la skill dada (csl_IsLearn).
func (p *ZonePlayer) HasSkill(skillID uint16) bool {
	if p == nil || p.LearnedSkills == nil {
		return false
	}
	_, ok := p.LearnedSkills[skillID]
	return ok
}

// GetLearnedSkill retorna la entrada aprendida de la skill o nil.
func (p *ZonePlayer) GetLearnedSkill(skillID uint16) *database.CharacterSkill {
	if p == nil || p.LearnedSkills == nil {
		return nil
	}
	return p.LearnedSkills[skillID]
}

// SetLearnedSkill registra o actualiza una skill aprendida en el caché en memoria.
func (p *ZonePlayer) SetLearnedSkill(sk *database.CharacterSkill) {
	if p == nil || sk == nil {
		return
	}
	if p.LearnedSkills == nil {
		p.LearnedSkills = make(map[uint16]*database.CharacterSkill)
	}
	cp := *sk
	p.LearnedSkills[sk.SkillNo] = &cp
}

// RemoveLearnedSkill elimina una skill aprendida del caché en memoria.
func (p *ZonePlayer) RemoveLearnedSkill(skillID uint16) {
	if p == nil || p.LearnedSkills == nil {
		return
	}
	delete(p.LearnedSkills, skillID)
}

// GetSkillEmpower retorna los 4 niveles de empoderamiento (0..5) de una skill aprendida:
// track 0 = damage/heal, track 1 = SP, track 2 = keepTime, track 3 = coolTime.
func (p *ZonePlayer) GetSkillEmpower(skillID uint16) (dmg, sp, keepTime, coolTime uint8) {
	sk := p.GetLearnedSkill(skillID)
	if sk == nil {
		return 0, 0, 0, 0
	}
	return sk.SkillPowerDemage, sk.SkillPowerSP, sk.SkillPowerKeepTime, sk.SkillPowerCoolTime
}

// TotalEmpowerUsed suma todos los puntos de empower invertidos en todas las skills aprendidas (csl_EmpowerSummary).
func (p *ZonePlayer) TotalEmpowerUsed() int {
	if p == nil || p.LearnedSkills == nil {
		return 0
	}
	total := 0
	for _, sk := range p.LearnedSkills {
		if sk != nil {
			total += int(sk.SkillPowerDemage) + int(sk.SkillPowerSP) + int(sk.SkillPowerKeepTime) + int(sk.SkillPowerCoolTime)
		}
	}
	return total
}

// MaxEmpowerPoints retorna el máximo de puntos de empower disponibles según nivel del personaje.
// En Fiesta Online oficial se gana 1 punto de empower por cada nivel a partir de nivel 2 (level - 1).
func (p *ZonePlayer) MaxEmpowerPoints() int {
	if p == nil || p.Char == nil || p.Char.Level <= 1 {
		return 0
	}
	return int(p.Char.Level) - 1
}

// NextSkillCastIndex genera un castIndex dinámico no nulo (1..65535), equivalente a sbec_GetIndex en Zone.exe.
func (p *ZonePlayer) NextSkillCastIndex() uint16 {
	if p == nil {
		return 1
	}
	p.NextCastIndex++
	if p.NextCastIndex == 0 {
		p.NextCastIndex = 1
	}
	return p.NextCastIndex
}
