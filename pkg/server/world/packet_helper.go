package world

import (
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ItemSlot and ItemEquip identifiers matching official Fiesta client equip slots (ItemEquipEnum).
const (
	EquipHead         byte = 1
	EquipAccHeadB     byte = 2
	EquipBody         byte = 7
	EquipAccBody      byte = 8
	EquipAccBack      byte = 9
	EquipLeftHand     byte = 10
	EquipAccLeftHand  byte = 11
	EquipRightHand    byte = 12
	EquipAccRightHand byte = 13
	EquipCosEff       byte = 17
	EquipAccHip       byte = 18
	EquipPant         byte = 19
	EquipAccPant      byte = 20
	EquipBoot         byte = 21
	EquipAccBoot      byte = 22
	EquipMouth        byte = 24
	EquipMinimon      byte = 25
	EquipEye          byte = 26
	EquipAccHeadA     byte = 27
	EquipMiniMonR     byte = 28
	EquipAccShield    byte = 29

	// Aliases for compatibility
	SlotHelm          byte = EquipHead
	SlotCostumeHelm   byte = EquipAccHeadB
	SlotArmor         byte = EquipBody
	SlotCostumeArmor  byte = EquipAccBody
	SlotWing          byte = EquipAccBack
	SlotWeapon2       byte = EquipLeftHand
	SlotCostumeShield byte = EquipAccLeftHand
	SlotWeapon        byte = EquipRightHand
	SlotCostumeWeapon byte = EquipAccRightHand
	SlotTail          byte = EquipAccHip
	SlotPants         byte = EquipPant
	SlotCostumePants  byte = EquipAccPant
	SlotBoots         byte = EquipBoot
	SlotCostumeBoots  byte = EquipAccBoot
	SlotGlasses       byte = EquipEye
	SlotPet           byte = EquipMiniMonR
)

// WriteBasicCharInfo serializa la información del personaje según el struct oficial
// PROTO_AVATARINFORMATION (130 bytes exactos) extraído de WorldManager.pdb / Character.pdb.
func WriteBasicCharInfo(ch *database.Character, pw *protocol.PacketWriter) {
	// 0x000 (4 bytes): unsigned int chrregnum
	pw.WriteUint32(uint32(ch.ID))

	// 0x004 (20 bytes): Name5 name
	pw.WriteString(ch.Name, 20)

	// 0x018 (2 bytes): unsigned short level
	pw.WriteUint16(uint16(ch.Level))

	// 0x01a (1 byte): unsigned char slot
	pw.WriteByte(ch.Slot)

	// 0x01b (12 bytes): Name3 loginmap
	mapShortName := ch.TownName
	pw.WriteString(mapShortName, 12)

	// 0x027 (5 bytes): PROTO_AVATAR_DELETE_INFO delinfo (year, month, day, hour, min)
	delInfo := [5]byte{0, 0, 0, 0, 0}
	if ch.DeletedAt != nil && !ch.DeletedAt.IsZero() {
		delInfo[0] = byte(ch.DeletedAt.Year() - 2000)
		delInfo[1] = byte(ch.DeletedAt.Month())
		delInfo[2] = byte(ch.DeletedAt.Day())
		delInfo[3] = byte(ch.DeletedAt.Hour())
		delInfo[4] = byte(ch.DeletedAt.Minute())
	}
	pw.WriteBytes(delInfo[:])

	// 0x02c (4 bytes): PROTO_AVATAR_SHAPE_INFO shape (race:2, chrclass:5, gender:1, hairtype, haircolor, faceshape)
	race := ch.Race
	if race == 0 {
		if ch.Job == 21 || ch.Job == 26 {
			race = 2
		} else {
			race = 1
		}
	}
	genderBit := byte(0)
	if ch.Male {
		genderBit = 1
	}
	shapeByte0 := byte((race & 0x03) | ((ch.Job & 0x1F) << 2) | ((genderBit & 0x01) << 7))
	pw.WriteByte(shapeByte0)
	pw.WriteByte(ch.Hair)
	pw.WriteByte(ch.HairColor)
	pw.WriteByte(ch.Face)

	// 0x030 (43 bytes): PROTO_EQUIPMENT equip (20 slots * 2 bytes + 3 bytes upgrade)
	// Orden oficial Fiesta PDB:
	// 1. Equ_Head (Slot 1)
	pw.WriteUint16(ch.GetEquip(EquipHead))
	// 2. Equ_Mouth (Slot 24)
	pw.WriteUint16(ch.GetEquip(EquipMouth))
	// 3. Equ_RightHand (Slot 12)
	pw.WriteUint16(ch.GetEquip(EquipRightHand))
	// 4. Equ_Body (Slot 7)
	pw.WriteUint16(ch.GetEquip(EquipBody))
	// 5. Equ_LeftHand (Slot 10)
	pw.WriteUint16(ch.GetEquip(EquipLeftHand))
	// 6. Equ_Pant (Slot 19)
	pw.WriteUint16(ch.GetEquip(EquipPant))
	// 7. Equ_Boot (Slot 21)
	pw.WriteUint16(ch.GetEquip(EquipBoot))
	// 8. Equ_AccBoot (Slot 22)
	pw.WriteUint16(ch.GetEquip(EquipAccBoot))
	// 9. Equ_AccPant (Slot 20)
	pw.WriteUint16(ch.GetEquip(EquipAccPant))
	// 10. Equ_AccBody (Slot 8)
	pw.WriteUint16(ch.GetEquip(EquipAccBody))
	// 11. Equ_AccHeadA (Slot 27 / Slot 2)
	accHead := ch.GetEquip(EquipAccHeadA)
	if accHead == 0 {
		accHead = ch.GetEquip(EquipAccHeadB)
	}
	pw.WriteUint16(accHead)
	// 12. Equ_MiniMon_R (Slot 28)
	pw.WriteUint16(ch.GetEquip(EquipMiniMonR))
	// 13. Equ_Eye (Slot 26)
	pw.WriteUint16(ch.GetEquip(EquipEye))
	// 14. Equ_AccLeftHand (Slot 11)
	pw.WriteUint16(ch.GetEquip(EquipAccLeftHand))
	// 15. Equ_AccRightHand (Slot 13)
	pw.WriteUint16(ch.GetEquip(EquipAccRightHand))
	// 16. Equ_AccBack (Slot 9)
	pw.WriteUint16(ch.GetEquip(EquipAccBack))
	// 17. Equ_CosEff (Slot 17)
	pw.WriteUint16(ch.GetEquip(EquipCosEff))
	// 18. Equ_AccHip (Slot 18)
	pw.WriteUint16(ch.GetEquip(EquipAccHip))
	// 19. Equ_Minimon (Slot 25)
	pw.WriteUint16(ch.GetEquip(EquipMinimon))
	// 20. Equ_AccShield (Slot 29)
	pw.WriteUint16(ch.GetEquip(EquipAccShield))

	// 3 bytes upgrade bitfield:
	// b0: (u10 & 0x0F) | ((u12 & 0x0F) << 4) (LeftHand | RightHand)
	// b1: (u7 & 0x0F)  | ((u19 & 0x0F) << 4) (Body | Leg)
	// b2: (u21 & 0x0F) (Shoes)
	u10 := ch.GetEquipUpgrade(EquipLeftHand)
	u12 := ch.GetEquipUpgrade(EquipRightHand)
	u7 := ch.GetEquipUpgrade(EquipBody)
	u19 := ch.GetEquipUpgrade(EquipPant)
	u21 := ch.GetEquipUpgrade(EquipBoot)
	b0 := byte((u10 & 0x0F) | ((u12 & 0x0F) << 4))
	b1 := byte((u7 & 0x0F) | ((u19 & 0x0F) << 4))
	b2 := byte(u21 & 0x0F)
	pw.WriteBytes([]byte{b0, b1, b2})

	// 0x05b (4 bytes): unsigned int nKQHandle
	pw.WriteUint32(ch.KQHandle)

	// 0x05f (12 bytes): Name3 sKQMapName
	pw.WriteString(ch.KQMapName, 12)

	// 0x06b (8 bytes): SHINE_XY_TYPE nKQCoord (x, y)
	pw.WriteUint32(ch.KQX)
	pw.WriteUint32(ch.KQY)

	// 0x073 (4 bytes): SHINE_DATETIME dKQDate
	pw.WriteUint32(ch.KQDate)

	// 0x077 (6 bytes): CHAR_ID_CHANGE_DATA CharIDChangeData (bNeedChangeID, bInit, nRowNo)
	needChange := byte(0)
	initFlag := byte(0)
	var rowNo uint32 = 0
	if ch.ChangeID != nil {
		if ch.ChangeID.NeedChangeID != 0 {
			needChange = 1
		}
		if ch.ChangeID.Init != 0 {
			initFlag = 1
		}
		rowNo = ch.ChangeID.RowNo
	}
	pw.WriteByte(needChange)
	pw.WriteByte(initFlag)
	pw.WriteUint32(rowNo)

	// 0x07d (5 bytes): PROTO_TUTORIAL_INFO TutorialInfo (nTutorialState, nTutorialStep)
	pw.WriteUint32(uint32(ch.TutorialState))
	pw.WriteByte(ch.TutorialStep)
}
