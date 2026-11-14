package zone

import (
	"encoding/binary"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Constantes de tipos en tItemOptions (sin colisionar con RandomOption 1..14 ni con p_Item_Upgrade IOT_NUMBER=0).
const (
	iotUpgradeGrade        int16 = 0
	iotFailCount           int16 = 510
	iotEndureEndTime       int16 = 511
	iotEndureGrade         int16 = 512
	iotHP                  int16 = 513
	iotIsPutOnBelonged     int16 = 514
	iotRandomChangeCount   int16 = 515
	iotRareFailCount       int16 = 516
	iotUseableDateTime     int16 = 517
	iotDeleteDateTime      int16 = 518
	iotAmount              int16 = 519
	iotChargedNumber       int16 = 520
	iotFurnitureKey        int16 = 521
	iotMoverGrade          int16 = 522
	iotMaxSocket           int16 = 523
	iotCreatedSocket       int16 = 524
	iotSocketGem0          int16 = 525
	iotSocketGem1          int16 = 526
	iotSocketGem2          int16 = 527
	iotSocketGemLot0       int16 = 528
	iotSocketGemLot1       int16 = 529
	iotSocketGemLot2       int16 = 530
	iotWPTitleID           int16 = 531
	iotWPLicenseMob        int16 = 532
	iotWPKillCount         int16 = 533
	iotPetID               int16 = 534
	iotPetRegNo            int16 = 535
	iotPetHungry           int16 = 536
	iotDuringRide          int16 = 537
	iotDuringPetSummon     int16 = 538
	iotRandomOptSlotBase   int16 = 600 // 600..607 -> (Type << 16) | Value
)

// ItemAttributeClass implementa los 62 métodos oficiales iac_* de Zone.pdb para
// codificación/decodificación binaria de atributos, mejora, durabilidad, montura, mascota, cofres y wptitle.
type ItemAttributeClass struct {
	ItemKey           int64
	ItemID            uint16
	ItemClass         uint8
	Lot               uint32
	Amount            uint32
	ChargedNumber     uint32
	Upgrade           uint8
	FailCount         uint16
	RareFailCount     uint16
	Grade             uint8
	EndureGrade       uint8
	EndureEndTime     uint32
	HP                uint32
	Hungry            uint16
	UseableDateTime   uint32
	DeleteDateTime    uint32
	IsPutOnBelonged   bool
	RandomChangeCount uint16
	MaxSocket         uint8
	CreatedSocket     uint8
	SocketGems        [3]uint16
	SocketGemLots     [3]uint8
	DuringPetSummon   bool
	DuringRide        bool
	PetID             uint16
	PetName           string
	PetRegistNumber   uint32
	FurnitureKey      uint32
	StatOptions       []item.ItemStatOption
	ChestItems        []uint64
	WPTitleID         uint16
	WPLicenseMob      uint16
	WPKillCount       uint32
}

// NewItemAttributeFromDB construye un ItemAttributeClass poblado desde un database.Item + tItemOptions.
func NewItemAttributeFromDB(it *database.Item, opts map[int16]int64, itemStore *item.Store) *ItemAttributeClass {
	iac := &ItemAttributeClass{}
	if it == nil {
		return iac
	}
	iac.ItemKey = it.ItemKey
	iac.ItemID = it.ItemID
	iac.Lot = uint32(it.Flags)
	if iac.Lot == 0 {
		iac.Lot = 1
	}
	if itemStore != nil {
		if cls, ok := itemStore.GetItemClass(it.ItemID); ok {
			iac.ItemClass = cls
		}
	}
	if opts == nil {
		return iac
	}
	if v, ok := opts[iotUpgradeGrade]; ok && v >= 0 {
		iac.Upgrade = uint8(v)
	}
	if v, ok := opts[iotFailCount]; ok {
		iac.FailCount = uint16(v)
	}
	if v, ok := opts[iotRareFailCount]; ok {
		iac.RareFailCount = uint16(v)
	}
	if v, ok := opts[iotEndureEndTime]; ok {
		iac.EndureEndTime = uint32(v)
	}
	if v, ok := opts[iotEndureGrade]; ok {
		iac.EndureGrade = uint8(v)
	}
	if v, ok := opts[iotHP]; ok {
		iac.HP = uint32(v)
	}
	if v, ok := opts[iotIsPutOnBelonged]; ok && v != 0 {
		iac.IsPutOnBelonged = true
	}
	if v, ok := opts[iotRandomChangeCount]; ok {
		iac.RandomChangeCount = uint16(v)
	}
	if v, ok := opts[iotUseableDateTime]; ok {
		iac.UseableDateTime = uint32(v)
	}
	if v, ok := opts[iotDeleteDateTime]; ok {
		iac.DeleteDateTime = uint32(v)
	}
	if v, ok := opts[iotAmount]; ok {
		iac.Amount = uint32(v)
	}
	if v, ok := opts[iotChargedNumber]; ok {
		iac.ChargedNumber = uint32(v)
	}
	if v, ok := opts[iotFurnitureKey]; ok {
		iac.FurnitureKey = uint32(v)
	}
	if v, ok := opts[iotMoverGrade]; ok {
		iac.Grade = uint8(v)
	}
	if v, ok := opts[iotMaxSocket]; ok {
		iac.MaxSocket = uint8(v)
	}
	if v, ok := opts[iotCreatedSocket]; ok {
		iac.CreatedSocket = uint8(v)
	}
	for i := 0; i < 3; i++ {
		if v, ok := opts[iotSocketGem0+int16(i)]; ok {
			iac.SocketGems[i] = uint16(v)
		}
		if v, ok := opts[iotSocketGemLot0+int16(i)]; ok {
			iac.SocketGemLots[i] = uint8(v)
		}
	}
	if v, ok := opts[iotWPTitleID]; ok {
		iac.WPTitleID = uint16(v)
	}
	if v, ok := opts[iotWPLicenseMob]; ok {
		iac.WPLicenseMob = uint16(v)
	}
	if v, ok := opts[iotWPKillCount]; ok {
		iac.WPKillCount = uint32(v)
	}
	if v, ok := opts[iotPetID]; ok {
		iac.PetID = uint16(v)
	}
	if v, ok := opts[iotPetRegNo]; ok {
		iac.PetRegistNumber = uint32(v)
	}
	if v, ok := opts[iotPetHungry]; ok {
		iac.Hungry = uint16(v)
	}
	if v, ok := opts[iotDuringRide]; ok && v != 0 {
		iac.DuringRide = true
	}
	if v, ok := opts[iotDuringPetSummon]; ok && v != 0 {
		iac.DuringPetSummon = true
	}

	// Leer RandomOptions guardadas en slots 600..607 primero
	for s := int16(0); s < 8; s++ {
		if encoded, ok := opts[iotRandomOptSlotBase+s]; ok && encoded > 0 {
			optType := uint8(encoded >> 16)
			optVal := uint16(encoded & 0xFFFF)
			iac.StatOptions = append(iac.StatOptions, item.ItemStatOption{Type: optType, Value: optVal})
		}
	}
	// Compatibilidad con opciones guardadas directamente por tipo 1..14 en drops.go
	if len(iac.StatOptions) == 0 {
		for k := int16(1); k <= 14; k++ {
			if val, ok := opts[k]; ok && val > 0 {
				iac.StatOptions = append(iac.StatOptions, item.ItemStatOption{Type: uint8(k), Value: uint16(val)})
			}
		}
	}
	return iac
}

// === Los 62 Métodos Oficiales ItemAttributeClass (`iac_*` de Zone.pdb) ===

func (c *ItemAttributeClass) iac_ClearItemOptionUpgrade() {
	c.Upgrade = 0
}

func (c *ItemAttributeClass) iac_DropItemMake(itemID uint16, itemClass uint8, lot uint32, opts []item.ItemStatOption) {
	c.ItemID = itemID
	c.ItemClass = itemClass
	if lot == 0 {
		lot = 1
	}
	c.Lot = lot
	c.StatOptions = append([]item.ItemStatOption(nil), opts...)
}

func (c *ItemAttributeClass) iac_EndOfOption(idx int) bool {
	return idx < 0 || idx >= len(c.StatOptions)
}

func (c *ItemAttributeClass) iac_FindOption(optType uint8) (item.ItemStatOption, int, bool) {
	for i, o := range c.StatOptions {
		if o.Type == optType {
			return o, i, true
		}
	}
	return item.ItemStatOption{}, -1, false
}

// iac_FoodEat aplica restauración de comida/poción o alimento de mascota/montura sobre el ItemAttributeClass y devuelve energía/HP/SP ganados.
func (c *ItemAttributeClass) iac_FoodEat(feedAmount uint16) uint16 {
	const maxHungry = uint16(10000)
	if c.Hungry+feedAmount > maxHungry {
		c.Hungry = maxHungry
	} else {
		c.Hungry += feedAmount
	}
	return c.Hungry
}

func (c *ItemAttributeClass) iac_GetAmount() uint32            { return c.Amount }
func (c *ItemAttributeClass) iac_GetAttrSize() int             { return item.CalcAttrSizeForClass(c.ItemClass, c.attrInnerCount()) }
func (c *ItemAttributeClass) iac_GetChargedNumber() uint32     { return c.ChargedNumber }
func (c *ItemAttributeClass) iac_GetDeteleDateTime() uint32    { return c.DeleteDateTime }
func (c *ItemAttributeClass) iac_GetDuringPetSummon() bool     { return c.DuringPetSummon }
func (c *ItemAttributeClass) iac_GetDuringRide() bool          { return c.DuringRide }
func (c *ItemAttributeClass) iac_GetEndureEndTime() uint32     { return c.EndureEndTime }
func (c *ItemAttributeClass) iac_GetGrade() uint8              { return c.Grade }
func (c *ItemAttributeClass) iac_GetHP() uint32                { return c.HP }
func (c *ItemAttributeClass) iac_GetHungry() uint16            { return c.Hungry }
func (c *ItemAttributeClass) iac_GetIsPutOnBelonged() bool     { return c.IsPutOnBelonged }
func (c *ItemAttributeClass) iac_GetItemOptionStruct() []item.ItemStatOption {
	out := make([]item.ItemStatOption, len(c.StatOptions))
	copy(out, c.StatOptions)
	return out
}
func (c *ItemAttributeClass) iac_GetItemOptionUpgrade() uint8        { return c.Upgrade }
func (c *ItemAttributeClass) iac_GetLot() uint32                     { return c.Lot }
func (c *ItemAttributeClass) iac_GetPetID() uint16                   { return c.PetID }
func (c *ItemAttributeClass) iac_GetPetName() string                 { return c.PetName }
func (c *ItemAttributeClass) iac_GetPetRegistNumber() uint32         { return c.PetRegistNumber }
func (c *ItemAttributeClass) iac_GetRandomOptionChangeCount() uint16 { return c.RandomChangeCount }
func (c *ItemAttributeClass) iac_GetRareFailCount() uint16           { return c.RareFailCount }
func (c *ItemAttributeClass) iac_GetUseableDateTime() uint32         { return c.UseableDateTime }
func (c *ItemAttributeClass) iac_Getitemchest() []uint64             { return append([]uint64(nil), c.ChestItems...) }
func (c *ItemAttributeClass) iac_Getupgrade() uint8                  { return c.Upgrade }
func (c *ItemAttributeClass) iac_InvenStore(lot uint32)              { c.Lot += lot }
func (c *ItemAttributeClass) iac_ItemClass() uint8                   { return c.ItemClass }
func (c *ItemAttributeClass) iac_SetAmount(v uint32)                 { c.Amount = v }
func (c *ItemAttributeClass) iac_SetChargedNumber(v uint32)          { c.ChargedNumber = v }
func (c *ItemAttributeClass) iac_SetDeteleDateTime(v uint32)         { c.DeleteDateTime = v }
func (c *ItemAttributeClass) iac_SetDuringPetSummon(v bool)          { c.DuringPetSummon = v }
func (c *ItemAttributeClass) iac_SetDuringRide(v bool)               { c.DuringRide = v }
func (c *ItemAttributeClass) iac_SetEndTime(v uint32)                { c.UseableDateTime = v }
func (c *ItemAttributeClass) iac_SetEndureEndTime(v uint32)          { c.EndureEndTime = v }
func (c *ItemAttributeClass) iac_SetEndureGrade(v uint8)             { c.EndureGrade = v }
func (c *ItemAttributeClass) iac_SetFurnitureKey(v uint32)           { c.FurnitureKey = v }
func (c *ItemAttributeClass) iac_SetGrade(v uint8)                   { c.Grade = v }
func (c *ItemAttributeClass) iac_SetHP(v uint32)                     { c.HP = v }
func (c *ItemAttributeClass) iac_SetIsPutOnBelonged(v bool)          { c.IsPutOnBelonged = v }
func (c *ItemAttributeClass) iac_SetItemOptionUpgrade(v uint8)       { c.Upgrade = v }
func (c *ItemAttributeClass) iac_SetLot(v uint32)                    { c.Lot = v }
func (c *ItemAttributeClass) iac_SetPetID(v uint16)                  { c.PetID = v }
func (c *ItemAttributeClass) iac_SetPetName(v string)                { c.PetName = v }
func (c *ItemAttributeClass) iac_SetPetRegistNumber(v uint32)        { c.PetRegistNumber = v }
func (c *ItemAttributeClass) iac_SetRandomOptionChangeCount(v uint16) {
	c.RandomChangeCount = v
}
func (c *ItemAttributeClass) iac_SetRareFailCount(v uint16)  { c.RareFailCount = v }
func (c *ItemAttributeClass) iac_SetUseableDateTime(v uint32) { c.UseableDateTime = v }

func (c *ItemAttributeClass) iac_Setoption(optType uint8, val uint16) {
	for i := range c.StatOptions {
		if c.StatOptions[i].Type == optType {
			c.StatOptions[i].Value = val
			return
		}
	}
	if len(c.StatOptions) < 8 {
		c.StatOptions = append(c.StatOptions, item.ItemStatOption{Type: optType, Value: val})
	}
}

func (c *ItemAttributeClass) iac_StepGet() int {
	return len(c.StatOptions)
}

func (c *ItemAttributeClass) iac_charge(delta uint32) uint32 {
	c.ChargedNumber += delta
	return c.ChargedNumber
}

// iac_equipbriefinfo retorna el campo compacto de 1 byte de apariencia de equipo: (Upgrade & 0x0F).
func (c *ItemAttributeClass) iac_equipbriefinfo() uint8 {
	return c.Upgrade & 0x0F
}

func (c *ItemAttributeClass) attrInnerCount() int {
	if c.ItemClass == 15 {
		return len(c.ChestItems)
	}
	if c.ItemClass == 13 {
		return int(c.MaxSocket)
	}
	return len(c.StatOptions)
}

// iac_fill2buffer serializa ItemAttributeClass al formato binario exacto de Fiesta Online.
func (c *ItemAttributeClass) iac_fill2buffer() []byte {
	buf := item.BuildItemAttrBuffer(c.ItemClass, c.Lot, c.StatOptions, c.attrInnerCount())
	if len(buf) == 0 {
		return buf
	}
	switch c.ItemClass {
	case 4: // Amulet / Ring (38 + 3*n)
		buf[0] = c.Upgrade
		if c.IsPutOnBelonged {
			buf[1] |= 0x01
		}
	case 5: // Weapon (65 + 3*n)
		buf[0] = c.Upgrade
		if c.IsPutOnBelonged {
			buf[1] |= 0x01
		}
		buf[2] = c.CreatedSocket
		for i := 0; i < 3; i++ {
			if 4+i*2+2 <= len(buf) {
				binary.LittleEndian.PutUint16(buf[4+i*2:6+i*2], c.SocketGems[i])
			}
		}
		if c.WPTitleID != 0 {
			binary.LittleEndian.PutUint16(buf[10:12], c.WPTitleID)
			binary.LittleEndian.PutUint16(buf[12:14], c.WPLicenseMob)
			binary.LittleEndian.PutUint32(buf[14:18], c.WPKillCount)
		}
	case 6, 7, 8, 38: // Armor, Shield, Boot, Bracelet (13 + 3*n)
		buf[0] = c.Upgrade
		if c.IsPutOnBelonged {
			buf[1] |= 0x01
		}
		buf[2] = c.CreatedSocket
		for i := 0; i < 3 && 4+i*2+2 <= 12; i++ {
			binary.LittleEndian.PutUint16(buf[4+i*2:6+i*2], c.SocketGems[i])
		}
	case 23: // Riding / Mover (19B)
		buf[0] = c.Grade
		buf[1] = c.Upgrade
		binary.LittleEndian.PutUint16(buf[2:4], c.Hungry)
		binary.LittleEndian.PutUint32(buf[4:8], c.UseableDateTime)
		if c.DuringRide {
			buf[8] = 1
		}
	case 36: // Pet / Minimon (26B)
		binary.LittleEndian.PutUint16(buf[0:2], c.PetID)
		binary.LittleEndian.PutUint32(buf[2:6], c.PetRegistNumber)
		binary.LittleEndian.PutUint16(buf[6:8], c.Hungry)
		if c.DuringPetSummon {
			buf[8] = 1
		}
		copy(buf[9:25], c.PetName)
	case 2, 18, 24, 26, 27, 37: // 4-byte classes including CostumWeapon(26), CostumShield(37), ActionItem(27)
		val := c.Lot
		if c.EndureEndTime > 0 && (c.ItemClass == 26 || c.ItemClass == 37) {
			val = c.EndureEndTime
		} else if c.ChargedNumber > 0 && c.ItemClass == 24 {
			val = c.ChargedNumber
		}
		binary.LittleEndian.PutUint32(buf[0:4], val)
	case 9: // Furniture (36B)
		binary.LittleEndian.PutUint32(buf[0:4], c.FurnitureKey)
		binary.LittleEndian.PutUint32(buf[4:8], c.EndureEndTime)
		buf[8] = c.EndureGrade
	}
	return buf
}

// iac_fill2struct deserializa un buffer binario de atributo en los campos de ItemAttributeClass.
func (c *ItemAttributeClass) iac_fill2struct(buf []byte) bool {
	if len(buf) == 0 {
		return c.ItemClass == 17 || c.ItemClass == 32
	}
	switch c.ItemClass {
	case 4:
		if len(buf) < 38 {
			return false
		}
		c.Upgrade = buf[0]
		c.IsPutOnBelonged = (buf[1] & 0x01) != 0
		return c.parseTailOptions(buf[37], buf[38:])
	case 5:
		if len(buf) < 65 {
			return false
		}
		c.Upgrade = buf[0]
		c.IsPutOnBelonged = (buf[1] & 0x01) != 0
		c.CreatedSocket = buf[2]
		for i := 0; i < 3; i++ {
			c.SocketGems[i] = binary.LittleEndian.Uint16(buf[4+i*2 : 6+i*2])
		}
		c.WPTitleID = binary.LittleEndian.Uint16(buf[10:12])
		c.WPLicenseMob = binary.LittleEndian.Uint16(buf[12:14])
		c.WPKillCount = binary.LittleEndian.Uint32(buf[14:18])
		return c.parseTailOptions(buf[64], buf[65:])
	case 6, 7, 8, 38:
		if len(buf) < 13 {
			return false
		}
		c.Upgrade = buf[0]
		c.IsPutOnBelonged = (buf[1] & 0x01) != 0
		c.CreatedSocket = buf[2]
		for i := 0; i < 3 && 6+i*2 <= 12; i++ {
			c.SocketGems[i] = binary.LittleEndian.Uint16(buf[4+i*2 : 6+i*2])
		}
		return c.parseTailOptions(buf[12], buf[13:])
	case 0, 11, 12, 14, 16, 19, 20, 22, 25, 28, 33, 34:
		c.Lot = uint32(buf[0])
		return true
	case 1, 3, 21, 35:
		if len(buf) < 2 {
			return false
		}
		c.Lot = uint32(binary.LittleEndian.Uint16(buf[0:2]))
		return true
	case 2, 18, 24, 26, 27, 37:
		if len(buf) < 4 {
			return false
		}
		v := binary.LittleEndian.Uint32(buf[0:4])
		c.Lot = v
		if c.ItemClass == 26 || c.ItemClass == 37 {
			c.EndureEndTime = v
		}
		return true
	case 23:
		if len(buf) < 19 {
			return false
		}
		c.Grade = buf[0]
		c.Upgrade = buf[1]
		c.Hungry = binary.LittleEndian.Uint16(buf[2:4])
		c.UseableDateTime = binary.LittleEndian.Uint32(buf[4:8])
		c.DuringRide = buf[8] != 0
		return true
	case 36:
		if len(buf) < 26 {
			return false
		}
		c.PetID = binary.LittleEndian.Uint16(buf[0:2])
		c.PetRegistNumber = binary.LittleEndian.Uint32(buf[2:6])
		c.Hungry = binary.LittleEndian.Uint16(buf[6:8])
		c.DuringPetSummon = buf[8] != 0
		c.PetName = strings.TrimRight(string(buf[9:25]), "\x00")
		return true
	default:
		return true
	}
}

func (c *ItemAttributeClass) parseTailOptions(flag uint8, tail []byte) bool {
	n := int(flag >> 1)
	if len(tail) < n*3 {
		return false
	}
	c.StatOptions = make([]item.ItemStatOption, 0, n)
	for i := 0; i < n; i++ {
		off := i * 3
		c.StatOptions = append(c.StatOptions, item.ItemStatOption{
			Type:  tail[off],
			Value: binary.LittleEndian.Uint16(tail[off+1 : off+3]),
		})
	}
	return true
}

func (c *ItemAttributeClass) iac_islotitem() bool {
	return c.ItemClass <= 3 || c.ItemClass == 11 || c.ItemClass == 12 || c.ItemClass == 14 || c.ItemClass == 16
}

func (c *ItemAttributeClass) iac_itemcreate(itemID uint16, itemClass uint8, lot uint32) {
	c.ItemID = itemID
	c.ItemClass = itemClass
	if lot == 0 {
		lot = 1
	}
	c.Lot = lot
}

func (c *ItemAttributeClass) iac_wptitle_Getlicencemob() uint16 { return c.WPLicenseMob }
func (c *ItemAttributeClass) iac_wptitle_Gettitle() uint16      { return c.WPTitleID }
func (c *ItemAttributeClass) iac_wptitle_Setlicense(mobID uint16) {
	c.WPLicenseMob = mobID
}
func (c *ItemAttributeClass) iac_wptitle_Settitle(titleID uint16) {
	c.WPTitleID = titleID
}
func (c *ItemAttributeClass) iac_wptitle_mobkilladd(delta uint32) uint32 {
	c.WPKillCount += delta
	return c.WPKillCount
}

// buildItemAttrFromDB serializa un ítem con todas sus opciones DB en []byte para red.
func (h *Handler) buildItemAttrFromDB(it *database.Item) []byte {
	if it == nil {
		return nil
	}
	var opts map[int16]int64
	if h.charRepo != nil && it.ItemKey > 0 {
		opts, _ = h.charRepo.GetItemOptions(it.ItemKey)
	}
	iac := NewItemAttributeFromDB(it, opts, h.itemStore)
	return iac.iac_fill2buffer()
}

// buildItemPacketInform construye un bloque [datasize(1B = 2+len(attr)) | itemid(2B) | attr].
func (h *Handler) buildItemPacketInform(it *database.Item) []byte {
	if it == nil || it.ItemID == 0 || it.ItemID == 0xFFFF {
		return nil
	}
	attr := h.buildItemAttrFromDB(it)
	out := make([]byte, 3+len(attr))
	out[0] = byte(2 + len(attr))
	binary.LittleEndian.PutUint16(out[1:3], it.ItemID)
	copy(out[3:], attr)
	return out
}

// buildItemInlineAttr serializa iti_fill2buffer_union (`itemid(2B) | attr`) sin byte datasize inicial,
// usado por PROTO_NC_TRADE_OPPOSITUPBOARD_CMD (0x4C10), BoothItemList SELL y BoothItemInfo.
func (h *Handler) buildItemInlineAttr(it *database.Item) []byte {
	if it == nil || it.ItemID == 0 || it.ItemID == 0xFFFF {
		return nil
	}
	attr := h.buildItemAttrFromDB(it)
	out := make([]byte, 2+len(attr))
	binary.LittleEndian.PutUint16(out[0:2], it.ItemID)
	copy(out[2:], attr)
	return out
}

// buildItemPacketInformWithLoc construye un PROTO_ITEMPACKET_INFORM oficial de ias_FillBufferItem:
// [0] datasize (= 4 + len(attr)) | [1..2] location ((box<<10)|slot) | [3..4] itemid | [5..] attr.
func (h *Handler) buildItemPacketInformWithLoc(box uint8, slot int16, it *database.Item) []byte {
	if it == nil || it.ItemID == 0 || it.ItemID == 0xFFFF {
		return nil
	}
	attr := h.buildItemAttrFromDB(it)
	loc := (uint16(box) << 10) | (uint16(slot) & 0x03FF)
	out := make([]byte, 5+len(attr))
	out[0] = byte(4 + len(attr))
	binary.LittleEndian.PutUint16(out[1:3], loc)
	binary.LittleEndian.PutUint16(out[3:5], it.ItemID)
	copy(out[5:], attr)
	return out
}

// buildFixed103ItemInfo rellena el buffer fijo de 103 bytes usado en Trade/Booth/Repurchase
// (`ItemInfo [103]byte` = 1B datasize + 2B itemID + up to 100B attr).
func (h *Handler) buildFixed103ItemInfo(it *database.Item) [103]byte {
	var out [103]byte
	raw := h.buildItemPacketInform(it)
	if len(raw) > 103 {
		raw = raw[:103]
		raw[0] = 102
	}
	copy(out[:], raw)
	return out
}

// syncPlayerMoney persiste y notifica PROTO_NC_CHAR_CENCHANGE_CMD (0x1033) al cliente.
func (h *Handler) syncPlayerMoney(player *ZonePlayer) {
	if player == nil || player.Char == nil {
		return
	}
	if h.charRepo != nil {
		_ = h.charRepo.SaveCharMoney(player.Char.ID, player.Char.TCharacterState.Money)
	}
	if player.Client != nil {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CENCHANGE_CMD)
		pw.WriteUint64(uint64(player.Char.TCharacterState.Money))
		_ = player.Client.SendPacket(pw)
	}
}

// addItemToInventory crea o apila un ítem en el inventario (box=9) y notifica PROTO_NC_ITEM_CELLCHANGE_CMD.
func (h *Handler) addItemToInventory(player *ZonePlayer, itemID uint16, lot uint32, opts []item.ItemStatOption) (*database.Item, bool) {
	if player == nil || player.Char == nil || h.charRepo == nil || itemID == 0 || itemID == 0xFFFF {
		return nil, false
	}
	if lot == 0 {
		lot = 1
	}
	maxLot := uint32(1)
	if h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(itemID); ok && info != nil && info.MaxLot > 0 {
			maxLot = info.MaxLot
		}
	}
	// Intentar apilar si MaxLot > 1 y no tiene opciones personalizadas
	if maxLot > 1 && len(opts) == 0 {
		items, err := h.charRepo.GetCharacterItems(player.Char.ID)
		if err == nil {
			for _, existing := range items {
				if existing != nil && existing.StorageType == 9 && existing.ItemID == itemID {
					curLot := uint32(existing.Flags)
					if curLot == 0 {
						curLot = 1
					}
					if curLot+lot <= maxLot {
						existing.Flags = int32(curLot + lot)
						_ = h.charRepo.UpdateItemFlags(existing.ItemKey, int(existing.Flags))
						if player.Client != nil {
							h.sendCellChange(player.Client, 9, existing.Storage, existing, 0xFFFF)
						}
						return existing, true
					}
				}
			}
		}
	}
	slot, ok := h.findFreeInvenSlot(player.Char.ID)
	if !ok {
		return nil, false
	}
	newIt := &database.Item{
		StorageType: 9,
		Owner:       player.Char.ID,
		Storage:     slot,
		ItemID:      itemID,
		Flags:       int32(lot),
	}
	if err := h.charRepo.CreateItem(newIt); err != nil {
		return nil, false
	}
	for i, opt := range opts {
		if i >= 8 {
			break
		}
		encoded := (int64(opt.Type) << 16) | int64(opt.Value)
		_ = h.charRepo.SetItemOption(newIt.ItemKey, iotRandomOptSlotBase+int16(i), encoded)
	}
	if player.Client != nil {
		h.sendCellChangeFull(player.Client, 9, slot, newIt, 0xFFFF)
	}
	return newIt, true
}

// sendCellChangeFull envía PROTO_NC_ITEM_CELLCHANGE_CMD / EQUIPCHANGE_CMD usando buildItemAttrFromDB completo.
func (h *Handler) sendCellChangeFull(client interface{ SendPacket(*protocol.PacketWriter) error }, box uint8, slot int16, it *database.Item, exchange uint16) {
	if client == nil {
		return
	}
	if box == 8 {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIPCHANGE_CMD)
		pw.WriteUint16(exchange)
		pw.WriteByte(byte(slot))
		if it == nil || it.ItemID == 0xFFFF {
			pw.WriteUint16(0xFFFF)
		} else {
			pw.WriteUint16(it.ItemID)
			pw.WriteBytes(h.buildItemAttrFromDB(it))
		}
		_ = client.SendPacket(pw)
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
		pw.WriteBytes(h.buildItemAttrFromDB(it))
	}
	_ = client.SendPacket(pw)
}

// consumeItemUnits consume `qty` unidades de un ítem (decrementa Flags o elimina de DB y envía CellChange).
func (h *Handler) consumeItemUnits(player *ZonePlayer, it *database.Item, qty uint32) bool {
	if player == nil || it == nil || h.charRepo == nil || qty == 0 {
		return false
	}
	curLot := uint32(it.Flags)
	if curLot == 0 {
		curLot = 1
	}
	if curLot < qty {
		return false
	}
	if curLot == qty {
		if err := h.charRepo.DeleteItem(it.ItemKey); err != nil {
			return false
		}
		h.sendCellChangeFull(player.Client, it.StorageType, it.Storage, nil, 0xFFFF)
		return true
	}
	it.Flags = int32(curLot - qty)
	if err := h.charRepo.UpdateItemFlags(it.ItemKey, int(it.Flags)); err != nil {
		return false
	}
	h.sendCellChangeFull(player.Client, it.StorageType, it.Storage, it, 0xFFFF)
	return true
}

// === Caches y Consultas SHN de la Fase 14 ===

type BelongRule struct {
	PutOnBelonged bool
	NoDrop        bool
	NoSell        bool
	NoStorage     bool
	NoTrade       bool
	NoDelete      bool
}

type phase14SHNCache struct {
	once          sync.Once
	belongRules   map[string]BelongRule
	itemShopLots  map[string]uint32
	itemMoneyCost map[string]struct {
		MoneyIdx string
		MoneyLot uint16
	}
	marketSubTypes map[string]uint32
}

var p14Cache phase14SHNCache

func (h *Handler) getSHN(name string) *shn.File {
	if h == nil || h.dataStore == nil {
		return nil
	}
	return h.dataStore.GetSHN(name)
}

func (h *Handler) getBelongRule(itemID uint16) BelongRule {
	if h == nil || h.dataStore == nil || h.itemStore == nil {
		return BelongRule{}
	}
	info, ok := h.itemStore.GetItemInfo(itemID)
	if !ok || info == nil {
		return BelongRule{}
	}
	// Leer BelongType de ItemInfo.shn si existe columna, y mapear en BelongTypeInfo.shn
	itemInfoSHN := h.getSHN("ItemInfo")
	belongTbl := h.getSHN("BelongTypeInfo")
	if itemInfoSHN == nil || belongTbl == nil {
		return BelongRule{}
	}
	row := itemInfoSHN.FindRow("InxName", info.InxName)
	if row == nil {
		return BelongRule{}
	}
	btInx := strings.TrimSpace(row.GetString("BelongType"))
	if btInx == "" || btInx == "-" {
		return BelongRule{}
	}
	p14Cache.once.Do(func() {
		p14Cache.belongRules = make(map[string]BelongRule)
		for _, r := range belongTbl.Rows {
			k := strings.TrimSpace(r.GetString("BT_Inx"))
			if k == "" {
				continue
			}
			p14Cache.belongRules[strings.ToLower(k)] = BelongRule{
				PutOnBelonged: r.GetInt("PutOnBelonged") != 0,
				NoDrop:        r.GetInt("NoDrop") != 0,
				NoSell:        r.GetInt("NoSell") != 0,
				NoStorage:     r.GetInt("NoStorage") != 0,
				NoTrade:       r.GetInt("NoTrade") != 0,
				NoDelete:      r.GetInt("NoDelete") != 0,
			}
		}
	})
	if rule, ok := p14Cache.belongRules[strings.ToLower(btInx)]; ok {
		return rule
	}
	return BelongRule{}
}

func (h *Handler) isItemBound(it *database.Item) bool {
	if it == nil {
		return false
	}
	if h.charRepo != nil && it.ItemKey > 0 {
		opts, err := h.charRepo.GetItemOptions(it.ItemKey)
		if err == nil && opts[iotIsPutOnBelonged] != 0 {
			return true
		}
	}
	return false
}

func (h *Handler) getShopBundleLot(inxName string) uint32 {
	tbl := h.getSHN("ItemShop")
	if tbl == nil || inxName == "" {
		return 1
	}
	if row := tbl.FindRow("InxName", inxName); row != nil {
		if lot := row.GetUint32("Lot"); lot > 0 {
			return lot
		}
	}
	return 1
}

func (h *Handler) getMarketSearchType(inxName string) uint32 {
	if h == nil || inxName == "" {
		return 0
	}
	srvTbl := h.getSHN("ItemInfoServer")
	mktTbl := h.getSHN("MarketSearchInfo")
	if srvTbl == nil || mktTbl == nil {
		return 0
	}
	row := srvTbl.FindRow("InxName", inxName)
	if row == nil {
		return 0
	}
	subType := row.GetInt("AuctionSubType")
	if subType <= 0 {
		subType = row.GetInt("MarketSubType")
	}
	if subType <= 0 {
		return 0
	}
	if mRow := mktTbl.FindRowInt("AuctionSubType", int64(subType)); mRow != nil {
		return uint32(mRow.GetInt("MarketSearchType"))
	}
	return 0
}

// encodeSHNDateTimestamp codifica la fecha actual + días en formato compacto uint32 compatible con ShineDateTime.
func encodeShineEndTimeDays(baseTS uint32, addDays uint32) uint32 {
	now := uint32(time.Now().Unix())
	if baseTS > now {
		return baseTS + addDays*86400
	}
	return now + addDays*86400
}

// encodeSHNDateTimeDummy helper exported for tests
func EncodeItemOptionsMap(opts []item.ItemStatOption) [25]byte {
	var out [25]byte
	n := len(opts)
	if n > 8 {
		n = 8
	}
	out[0] = byte((n << 1) | 1)
	for i := 0; i < n; i++ {
		off := 1 + i*3
		out[off] = opts[i].Type
		binary.LittleEndian.PutUint16(out[off+1:off+3], opts[i].Value)
	}
	return out
}
