package zone

import (
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ItemSlot identifiers matching NextGen.FiestaLib/ItemSlot.cs
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

// BaseStats holds primary attributes and vitals for a character.
type BaseStats struct {
	Level       byte
	Str         byte
	End         byte
	Dex         byte
	Int         byte
	Spr         byte
	MaxHPStones int16
	MaxSPStones int16
	MaxHP       int
	MaxSP       int
}

// GetJobBaseStats returns the level 1 stats matching NextGen DataProvider.
func GetJobBaseStats(job byte) BaseStats {
	switch job {
	case byte(protocol.Job_Archer):
		return BaseStats{Level: 1, Str: 4, End: 4, Dex: 6, Int: 1, Spr: 3, MaxHPStones: 13, MaxSPStones: 11, MaxHP: 46, MaxSP: 24}
	case byte(protocol.Job_Cleric):
		return BaseStats{Level: 1, Str: 5, End: 4, Dex: 3, Int: 1, Spr: 4, MaxHPStones: 15, MaxSPStones: 11, MaxHP: 46, MaxSP: 32}
	case byte(protocol.Job_Fighter):
		return BaseStats{Level: 1, Str: 6, End: 5, Dex: 3, Int: 1, Spr: 1, MaxHPStones: 15, MaxSPStones: 7, MaxHP: 52, MaxSP: 10}
	case byte(protocol.Job_Mage):
		return BaseStats{Level: 1, Str: 1, End: 3, Dex: 3, Int: 10, Spr: 5, MaxHPStones: 12, MaxSPStones: 15, MaxHP: 42, MaxSP: 46}
	case byte(protocol.Job_Trickster):
		return BaseStats{Level: 1, Str: 5, End: 5, Dex: 4, Int: 1, Spr: 3, MaxHPStones: 15, MaxSPStones: 11, MaxHP: 48, MaxSP: 21}
	default:
		return BaseStats{Level: 1, Str: 6, End: 5, Dex: 3, Int: 1, Spr: 1, MaxHPStones: 15, MaxSPStones: 7, MaxHP: 52, MaxSP: 10}
	}
}

// ZonePlayer represents an active, authenticated player in the Zone.
type ZonePlayer struct {
	Client      *network.Client
	Char        *database.Character
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
}

// NewZonePlayer initializes a ZonePlayer with loaded stats.
func NewZonePlayer(client *network.Client, ch *database.Character) *ZonePlayer {
	stats := GetJobBaseStats(ch.Job)

	curHP := ch.CurHP
	if curHP <= 0 {
		curHP = stats.MaxHP
		ch.CurHP = curHP
	}
	curSP := ch.CurSP
	if curSP <= 0 {
		curSP = stats.MaxSP
		ch.CurSP = curSP
	}

	return &ZonePlayer{
		Client:      client,
		Char:        ch,
		State:       1, // PlayerState.Normal
		Rotation:    0,
		StonesHP:    stats.MaxHPStones,
		StonesSP:    stats.MaxSPStones,
		MaxStonesHP: stats.MaxHPStones,
		MaxStonesSP: stats.MaxSPStones,
		MaxHP:       stats.MaxHP,
		MaxSP:       stats.MaxSP,
		Str:         stats.Str,
		End:         stats.End,
		Dex:         stats.Dex,
		Int:         stats.Int,
		Spr:         stats.Spr,
	}
}

// WriteLook serializes avatar appearance (NextGen ZoneCharacter.WriteLook).
func (p *ZonePlayer) WriteLook(pw *protocol.PacketWriter) {
	var genderBit byte = 0
	if p.Char.Male {
		genderBit = 1
	}
	val := byte(0x01 | (p.Char.Job << 2) | (genderBit << 7))
	pw.WriteByte(val)
	pw.WriteByte(p.Char.Hair)
	pw.WriteByte(p.Char.HairColor)
	pw.WriteByte(p.Char.Face)
}

// WriteEquipment serializes worn gear slots (NextGen ZoneCharacter.WriteEquipment).
func (p *ZonePlayer) WriteEquipment(pw *protocol.PacketWriter) {
	pw.WriteUint16(p.Char.GetEquip(Slot_Helm))
	pw.WriteUint16(p.Char.GetEquip(Slot_Weapon))
	pw.WriteUint16(p.Char.GetEquip(Slot_Armor))
	pw.WriteUint16(p.Char.GetEquip(Slot_Weapon2))
	pw.WriteUint16(p.Char.GetEquip(Slot_Pants))
	pw.WriteUint16(p.Char.GetEquip(Slot_Boots))
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumeBoots))
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumePants))
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumeArmor))
	pw.Fill(6, 0xFF) // UNK
	pw.WriteUint16(p.Char.GetEquip(Slot_Glasses))
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumeHelm))
	pw.Fill(2, 0xFF) // UNK
	pw.WriteUint16(p.Char.GetEquip(Slot_CostumeWeapon))
	pw.WriteUint16(p.Char.GetEquip(Slot_Wing))
	pw.Fill(2, 0xFF) // UNK
	pw.WriteUint16(p.Char.GetEquip(Slot_Tail))
	pw.WriteUint16(p.Char.GetEquip(Slot_Pet))
}

// WriteRefinement serializes upgrade glows (NextGen ZoneCharacter.WriteRefinement).
func (p *ZonePlayer) WriteRefinement(pw *protocol.PacketWriter) {
	pw.WriteByte(0)
	pw.WriteByte(0)
	pw.WriteByte(0)
}

// WriteDetailedInfo serializes primary chunk info (NextGen ZoneCharacter.WriteDetailedInfo).
func (p *ZonePlayer) WriteDetailedInfo(pw *protocol.PacketWriter) {
	pw.WriteInt32(int32(p.Char.ID))
	pw.WriteString(p.Char.Name, 16)
	pw.WriteByte(p.Char.Slot)
	pw.WriteByte(p.Char.Level)
	pw.WriteInt64(0)        // Exp
	pw.WriteInt32(12345678) // UNK
	pw.WriteInt16(p.StonesHP)
	pw.WriteInt16(p.StonesSP)
	pw.WriteUint32(uint32(p.Char.CurHP))
	pw.WriteUint32(uint32(p.Char.CurSP))
	pw.WriteInt32(0) // Fame
	pw.WriteInt64(0) // Money
	shortName := "Rou"
	if p.Char.Map != 0 {
		shortName = "Rou"
	}
	pw.WriteString(shortName, 12)
	pw.WriteInt32(int32(p.Char.XPos))
	pw.WriteInt32(int32(p.Char.YPos))
	pw.WriteByte(p.Rotation)
	pw.WriteByte(p.Str)
	pw.WriteByte(p.End)
	pw.WriteByte(p.Dex)
	pw.WriteByte(p.Int)
	pw.WriteByte(p.Spr)
	pw.WriteInt16(0) // UNK
	pw.WriteUint32(0) // Killpoints
	pw.Fill(7, 0)     // UNK
}

// WriteDetailedInfoExtra serializes stats chunk (NextGen ZoneCharacter.WriteDetailedInfoExtra).
func (p *ZonePlayer) WriteDetailedInfoExtra(pw *protocol.PacketWriter, levelUp bool) {
	if !levelUp {
		pw.WriteUint16(p.MapObjectID)
	}

	pw.WriteInt64(0)    // Exp
	pw.WriteUint64(280) // MaxExp for level 1

	pw.WriteInt32(int32(p.Str))
	pw.WriteInt32(int32(p.Str))
	pw.WriteInt32(int32(p.End))
	pw.WriteInt32(int32(p.End))
	pw.WriteInt32(int32(p.Dex))
	pw.WriteInt32(int32(p.Dex))
	pw.WriteInt32(int32(p.Int))
	pw.WriteInt32(int32(p.Int))
	pw.WriteInt32(0) // Wisdom
	pw.WriteInt32(0) // UNK
	pw.WriteInt32(int32(p.Spr))
	pw.WriteInt32(int32(p.Spr))

	weaponDmg := int32(15)
	magicDmg := int32(5)
	weaponDef := int32(10)
	aim := int32(10)
	evasion := int32(5)
	magicDef := int32(5)

	pw.WriteInt32(weaponDmg)
	pw.WriteInt32(weaponDmg)
	pw.WriteInt32(magicDmg)
	pw.WriteInt32(magicDmg)
	pw.WriteInt32(weaponDef)
	pw.WriteInt32(weaponDef)
	pw.WriteInt32(aim)
	pw.WriteInt32(aim)
	pw.WriteInt32(evasion)
	pw.WriteInt32(evasion)
	pw.WriteInt32(weaponDmg)
	pw.WriteInt32(weaponDmg)
	pw.WriteInt32(magicDmg)
	pw.WriteInt32(magicDmg)
	pw.WriteInt32(magicDef)
	pw.WriteInt32(magicDef)

	pw.WriteInt32(1)
	pw.WriteInt32(20)
	pw.WriteInt32(2)
	pw.WriteInt32(40)

	pw.WriteUint32(uint32(p.MaxHP))
	pw.WriteUint32(uint32(p.MaxSP))

	pw.WriteInt32(0) // UNK
	pw.WriteInt32(int32(p.MaxStonesHP))
	pw.WriteInt32(int32(p.MaxStonesSP))
	pw.Fill(64, 0)

	if !levelUp {
		pw.WriteInt32(int32(p.Char.XPos))
		pw.WriteInt32(int32(p.Char.YPos))
	}
}

// WriteCharacterDisplay serializes in-world avatar entity (NextGen ZoneCharacter.WriteCharacterDisplay).
func (p *ZonePlayer) WriteCharacterDisplay(pw *protocol.PacketWriter) {
	pw.WriteUint16(p.MapObjectID)
	pw.WriteString(p.Char.Name, 16)
	pw.WriteInt32(int32(p.Char.XPos))
	pw.WriteInt32(int32(p.Char.YPos))
	pw.WriteByte(p.Rotation)
	pw.WriteByte(p.State)
	pw.WriteByte(p.Char.Job)
	p.WriteLook(pw)
	p.WriteEquipment(pw)
	p.WriteRefinement(pw)
	pw.WriteUint16(0xFFFF) // Mount handle
	pw.WriteUint16(0xFFFF)
	pw.WriteByte(0xFF) // Emote
	pw.WriteUint16(0xFFFF)
	pw.WriteInt16(0)
	pw.WriteUint16(0) // Mob ID
	pw.Fill(55, 0)    // Buff bits
	pw.WriteInt32(0)  // Guild ID
	pw.WriteByte(0x02) // UNK
	pw.WriteBool(false) // In Academy
	pw.WriteBool(true)  // Pet AutoPickup
	pw.WriteByte(p.Char.Level)
}
