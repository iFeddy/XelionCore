package zone

import (
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/classparam"
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

// JobToClassName mapea el ID numérico de profesión a su nombre de clase en data/World/ClassParam/.
func JobToClassName(job byte) string {
	switch job {
	case byte(protocol.Job_Archer):
		return "Archer"
	case byte(protocol.Job_Cleric):
		return "Cleric"
	case byte(protocol.Job_Fighter):
		return "Fighter"
	case byte(protocol.Job_Mage):
		return "Mage"
	case byte(protocol.Job_Trickster):
		return "Trickster"
	case byte(protocol.Job_Crusader):
		return "Crusader"
	default:
		return "Fighter"
	}
}

// ZonePlayer representa a un jugador autenticado y activo en una instancia de ZoneServer.
type ZonePlayer struct {
	Client      *network.Client
	Char        *database.Character
	CurrentMap  *ZoneMap
	MapObjectID uint16
	State       byte // 1=Normal, 2=Player, 3=Dead, 4=Resting, 5=Vendor, 6=Mount
	Rotation    byte
	StonesHP    int16
	StonesSP    int16
	MaxStonesHP int16
	MaxStonesSP int16
	MaxHP       int
	MaxSP       int
	Str         byte
	End         byte
	Dex         byte
	Int         byte
	Spr         byte
	LevelParam  classparam.LevelParam
	LoginTime   time.Time
}

// NewZonePlayer inicializa un ZonePlayer cargando sus stats oficiales desde ClassParam.
func NewZonePlayer(client *network.Client, ch *database.Character, classStore *classparam.Store) *ZonePlayer {
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

	maxHP := lp.MaxHP
	if maxHP <= 0 {
		maxHP = 52
	}
	maxSP := lp.MaxSP
	if maxSP <= 0 {
		maxSP = 10
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

	return &ZonePlayer{
		Client:      client,
		Char:        ch,
		State:       1,    // Normal
		Rotation:    0x55, // 85 grados
		LoginTime:   time.Now(),
		StonesHP:    maxStonesHP,
		StonesSP:    maxStonesSP,
		MaxStonesHP: maxStonesHP,
		MaxStonesSP: maxStonesSP,
		MaxHP:       maxHP,
		MaxSP:       maxSP,
		Str:         0,
		End:         0,
		Dex:         0,
		Int:         0,
		Spr:         0,
		LevelParam:  lp,
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
	pw.WriteBytes([]byte{0, 0, 0})                      // upgrade (3 bytes oficiales: lefthand, righthand, body, leg, shoes)
}

// WriteDetailedInfo serializa PROTO_NC_CHAR_BASE_CMD (105 bytes oficiales en PDB).
func (p *ZonePlayer) WriteDetailedInfo(pw *protocol.PacketWriter) {
	pw.WriteUint32(uint32(p.Char.ID))                          // chrregnum
	pw.WriteString(p.Char.Name, 20)                            // charid (Name5: 20 bytes)
	pw.WriteByte(p.Char.Slot)                                  // slotno
	pw.WriteByte(p.Char.Level)                                 // Level
	pw.WriteUint64(uint64(p.Char.TCharacterState.Exp))         // Experience
	pw.WriteUint16(0)                                          // CurPwrStone
	pw.WriteUint16(0)                                          // CurGrdStone
	pw.WriteUint16(uint16(p.StonesHP))                         // CurHPStone
	pw.WriteUint16(uint16(p.StonesSP))                         // CurSPStone
	pw.WriteUint32(uint32(p.Char.CurHP))                       // CurHP
	pw.WriteUint32(uint32(p.Char.CurSP))                       // CurSP
	pw.WriteUint32(0)                                          // CurLP
	pw.WriteUint32(uint32(p.Char.TCharacterState.Fame))        // fame
	pw.WriteUint64(uint64(p.Char.TCharacterState.Money))       // Cen (money)

	// logininfo: PROTO_NC_CHAR_BASE_CMD::LoginLocation (21 bytes)
	shortName := p.Char.TownName
	if p.CurrentMap != nil && p.CurrentMap.MapID != "" {
		shortName = p.CurrentMap.MapID
	}
	if shortName == "" {
		shortName = "RouN"
	}
	pw.WriteString(shortName, 12)       // currentmap (Name3: 12 bytes)
	pw.WriteUint32(uint32(p.Char.XPos)) // currentcoord.xy.x
	pw.WriteUint32(uint32(p.Char.YPos)) // currentcoord.xy.y
	pw.WriteByte(p.Rotation)            // currentcoord.dir

	// statdistribute: CHARSTATDISTSTR (6 bytes)
	pw.WriteByte(p.Str) // Strength
	pw.WriteByte(p.End) // Constitute
	pw.WriteByte(p.Dex) // Dexterity
	pw.WriteByte(p.Int) // Intelligence
	pw.WriteByte(p.Spr) // MentalPower
	pw.WriteByte(0)     // RedistributePoint

	pw.WriteByte(0)       // pkyellowtime
	pw.WriteUint32(0)     // pkcount
	pw.WriteUint16(0)     // prisonmin
	pw.WriteByte(0)       // adminlevel
	pw.Fill(4, 0)         // flags (4 bytes)
}

// WriteDetailedInfoExtra serializa PROTO_NC_CHAR_MAPLOGIN_ACK (242 bytes oficiales en PDB).
func (p *ZonePlayer) WriteDetailedInfoExtra(pw *protocol.PacketWriter, levelUp bool) {
	if !levelUp {
		pw.WriteUint16(p.MapObjectID) // charhandle (2 bytes)
	}

	// CHAR_PARAMETER_DATA (232 bytes)
	pw.WriteUint64(0)                               // PrevExp
	pw.WriteUint64(uint64(280))                     // NextExp
	pw.WriteUint32(uint32(p.LevelParam.Strength))   // Strength.base
	pw.WriteUint32(uint32(p.LevelParam.Strength + int(p.Str))) // Strength.change
	pw.WriteUint32(uint32(p.LevelParam.Constitution)) // Constitute.base
	pw.WriteUint32(uint32(p.LevelParam.Constitution + int(p.End))) // Constitute.change
	pw.WriteUint32(uint32(p.LevelParam.Dexterity))  // Dexterity.base
	pw.WriteUint32(uint32(p.LevelParam.Dexterity + int(p.Dex))) // Dexterity.change
	pw.WriteUint32(uint32(p.LevelParam.Intelligence)) // Intelligence.base
	pw.WriteUint32(uint32(p.LevelParam.Intelligence + int(p.Int))) // Intelligence.change
	pw.WriteUint32(uint32(p.LevelParam.Wizdom))       // Wizdom.base
	pw.WriteUint32(0)                               // Wizdom.change
	pw.WriteUint32(uint32(p.LevelParam.MentalPower))  // MentalPower.base
	pw.WriteUint32(uint32(p.LevelParam.MentalPower + int(p.Spr))) // MentalPower.change

	weaponDmg := uint32(15)
	magicDmg := uint32(5)
	weaponDef := uint32(10)
	aim := uint32(10)
	evasion := uint32(5)
	magicDef := uint32(5)

	pw.WriteUint32(weaponDmg) // WClow.base
	pw.WriteUint32(weaponDmg) // WClow.change
	pw.WriteUint32(magicDmg)  // WChigh.base
	pw.WriteUint32(magicDmg)  // WChigh.change
	pw.WriteUint32(weaponDef) // AC.base
	pw.WriteUint32(weaponDef) // AC.change
	pw.WriteUint32(aim)       // TH.base
	pw.WriteUint32(aim)       // TH.change
	pw.WriteUint32(evasion)   // TB.base
	pw.WriteUint32(evasion)   // TB.change
	pw.WriteUint32(weaponDmg) // MAlow.base
	pw.WriteUint32(weaponDmg) // MAlow.change
	pw.WriteUint32(magicDmg)  // MAhigh.base
	pw.WriteUint32(magicDmg)  // MAhigh.change
	pw.WriteUint32(magicDef)  // MR.base
	pw.WriteUint32(magicDef)  // MR.change
	pw.WriteUint32(aim)       // MH.base
	pw.WriteUint32(aim)       // MH.change
	pw.WriteUint32(evasion)   // MB.base
	pw.WriteUint32(evasion)   // MB.change

	pw.WriteUint32(uint32(p.MaxHP)) // MaxHp
	pw.WriteUint32(uint32(p.MaxSP)) // MaxSp
	pw.WriteUint32(0)               // MaxLp
	pw.WriteUint32(0)               // MaxAp
	pw.WriteUint32(uint32(p.MaxStonesHP)) // MaxHPStone
	pw.WriteUint32(uint32(p.MaxStonesSP)) // MaxSPStone

	// PwrStone (16 bytes)
	pw.WriteUint32(0)
	pw.WriteUint32(uint32(p.LevelParam.PwrStoneWC))
	pw.WriteUint32(uint32(p.LevelParam.PwrStoneMA))
	pw.WriteUint32(uint32(p.LevelParam.MaxPwrStone))

	// GrdStone (16 bytes)
	pw.WriteUint32(0)
	pw.WriteUint32(uint32(p.LevelParam.GrdStoneAC))
	pw.WriteUint32(uint32(p.LevelParam.GrdStoneMR))
	pw.WriteUint32(uint32(p.LevelParam.MaxGrdStone))

	// Resistances (32 bytes)
	pw.WriteUint32(uint32(p.LevelParam.PainRes))
	pw.WriteUint32(0)
	pw.WriteUint32(uint32(p.LevelParam.RestraintRes))
	pw.WriteUint32(0)
	pw.WriteUint32(uint32(p.LevelParam.CurseRes))
	pw.WriteUint32(0)
	pw.WriteUint32(uint32(p.LevelParam.ShockRes))
	pw.WriteUint32(0)

	// logincoord: SHINE_XY_TYPE (8 bytes)
	if !levelUp {
		pw.WriteUint32(uint32(p.Char.XPos))
		pw.WriteUint32(uint32(p.Char.YPos))
	}
}

// WriteCharacterDisplay serializes PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD (235 bytes exactos en PDB).
func (p *ZonePlayer) WriteCharacterDisplay(pw *protocol.PacketWriter) {
	pw.WriteUint16(p.MapObjectID)     // 0x000 (0): unsigned short handle (2)
	pw.WriteString(p.Char.Name, 20)   // 0x002 (2): Name5 charid (20)
	pw.WriteInt32(int32(p.Char.XPos)) // 0x016 (22): coord.x (4)
	pw.WriteInt32(int32(p.Char.YPos)) // coord.y (4)
	pw.WriteByte(p.Rotation)          // coord.dir (1)
	pw.WriteByte(p.State)             // 0x01f (31): mode (1)
	pw.WriteByte(p.Char.Job)          // 0x020 (32): chrclass (1)
	p.WriteLook(pw)                   // 0x021 (33): PROTO_AVATAR_SHAPE_INFO (4)
	// 0x025 (37): shapedata (45 bytes = equip 43 bytes + 2 bytes padding)
	p.WriteEquipment(pw)              // 43 bytes
	pw.WriteUint16(0)                 // 2 bytes pad
	pw.WriteUint16(0)                 // 0x052 (82): polymorph (2)
	pw.Fill(3, 0)                     // 0x054 (84): STOPEMOTICON_DESCRIPT emoticon (3)
	pw.Fill(4, 0)                     // 0x057 (87): CHARTITLE_BRIEFINFO chartitle (4)
	pw.Fill(99, 0)                    // 0x05b (91): ABNORMAL_STATE_BIT abstatebit (99)
	pw.WriteUint32(0)                 // 0x0be (190): myguild (4)
	pw.WriteByte(0)                   // 0x0c2 (194): type (1)
	pw.WriteByte(0)                   // 0x0c3 (195): isGuildAcademyMember (1)
	pw.WriteByte(1)                   // 0x0c4 (196): IsAutoPick (1)
	pw.WriteByte(p.Char.Level)        // 0x0c5 (197): Level (1)
	pw.WriteString("", 32)            // 0x0c6 (198): sAnimation (32)
	pw.WriteUint16(0)                 // 0x0e6 (230): nMoverHnd (2)
	pw.WriteByte(0)                   // 0x0e8 (232): nMoverSlot (1)
	pw.WriteByte(0)                   // 0x0e9 (233): nKQTeamType (1)
	pw.WriteByte(0)                   // 0x0ea (234): IsUseItemMinimon (1)
}
