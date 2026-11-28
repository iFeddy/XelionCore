package zone

import (
	"testing"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestFase5_EmpowerProductFieldAndUnlearn(t *testing.T) {
	h, _, player, client := setupFase12TestEnvironment(t)
	h.dataStore = &DataStore{
		SpecialSkills: &SpecialSkillConfig{
			ProductionSkills: [5]uint16{3001, 3002, 3003, 3004, 3005},
			EmpowerResetItem: 5001,
			SkillEraseItem:   5002,
		},
	}

	player.Char.Level = 10 // MaxEmpowerPoints = 9
	player.SyncLearnedSkills([]*database.CharacterSkill{
		{CharNo: player.Char.ID, SkillNo: 101, SkillLevel: 1},
	})

	// 1. EMPOWALLOC: asignar +3 en track 0
	empowPkt := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ)
	empowPkt.WriteUint16(101)
	empowPkt.WriteUint16(0x0003) // +3 dmg
	empowPkt.WriteUint16(0)
	h.HandlePacket(client, protocol.NewPacket(empowPkt.Data()))
	if d, _, _, _ := player.GetSkillEmpower(101); d != 3 {
		t.Errorf("Empower dmg esperado 3, obtenido %d", d)
	}

	// 2. EMPOWALLOC: intentar pasar de 5 (+3 adicional -> 6 > 5) debe fallar (0x0884)
	empowPkt2 := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ)
	empowPkt2.WriteUint16(101)
	empowPkt2.WriteUint16(0x0003)
	empowPkt2.WriteUint16(0)
	h.HandlePacket(client, protocol.NewPacket(empowPkt2.Data()))
	if d, _, _, _ := player.GetSkillEmpower(101); d != 3 {
		t.Errorf("Empower dmg no debería superar 5, obtenido %d", d)
	}

	// 3. EMPOWALLOC: minusMask sin ítem reset -> falla (0x0886)
	empowMinusFail := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ)
	empowMinusFail.WriteUint16(101)
	empowMinusFail.WriteUint16(0)
	empowMinusFail.WriteUint16(0x0002)
	h.HandlePacket(client, protocol.NewPacket(empowMinusFail.Data()))
	if d, _, _, _ := player.GetSkillEmpower(101); d != 3 {
		t.Errorf("Sin ítem reset no debería restar empower, obtenido %d", d)
	}

	// 4. EMPOWALLOC: dar 2 ítems reset (5001) en box 9 y restar 2 puntos -> éxito y consume ítems
	_ = h.charRepo.SaveCharacterItem(&database.Item{
		ItemKey:     9001,
		Owner:       player.Char.ID,
		StorageType: 9,
		Storage:     0,
		ItemID:      5001,
		Flags:       2,
	})
	empowMinusOK := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ)
	empowMinusOK.WriteUint16(101)
	empowMinusOK.WriteUint16(0)
	empowMinusOK.WriteUint16(0x0002)
	h.HandlePacket(client, protocol.NewPacket(empowMinusOK.Data()))
	if d, _, _, _ := player.GetSkillEmpower(101); d != 1 {
		t.Errorf("Tras restar 2 puntos con ítem reset, empower dmg esperado 1, obtenido %d", d)
	}

	// 5. PRODUCTFIELD: skill fuera de whitelist (9999) -> 0x0B12 (no se aprende)
	prodBad := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_PRODUCTFIELD_REQ)
	prodBad.WriteUint16(9999)
	h.HandlePacket(client, protocol.NewPacket(prodBad.Data()))
	if player.HasSkill(9999) {
		t.Errorf("Skill 9999 fuera de whitelist no debería aprenderse")
	}

	// 6. PRODUCTFIELD: aprender 2 profesiones válidas (3001 y 3002); la 3ra (3003) debe rechazarse (0x0B14)
	for _, pid := range []uint16{3001, 3002, 3003} {
		p := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_PRODUCTFIELD_REQ)
		p.WriteUint16(pid)
		h.HandlePacket(client, protocol.NewPacket(p.Data()))
	}
	if !player.HasSkill(3001) || !player.HasSkill(3002) || player.HasSkill(3003) {
		t.Errorf("ProductField máx 2 esperado: 3001=%v 3002=%v 3003=%v", player.HasSkill(3001), player.HasSkill(3002), player.HasSkill(3003))
	}

	// 7. UNLEARN: sin ítem SkillEraseItem (5002) en slot 1 -> falla (0x0702)
	unlearnFail := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_UNLEARN_REQ)
	unlearnFail.WriteUint8(1)
	unlearnFail.WriteUint16(101)
	h.HandlePacket(client, protocol.NewPacket(unlearnFail.Data()))
	if !player.HasSkill(101) {
		t.Errorf("Unlearn sin pergamino no debería eliminar la skill 101")
	}

	// 8. UNLEARN: con ítem SkillEraseItem (5002) en slot 1 -> elimina skill 101 y consume pergamino
	_ = h.charRepo.SaveCharacterItem(&database.Item{
		ItemKey:     9002,
		Owner:       player.Char.ID,
		StorageType: 9,
		Storage:     1,
		ItemID:      5002,
		Flags:       1,
	})
	unlearnOK := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_UNLEARN_REQ)
	unlearnOK.WriteUint8(1)
	unlearnOK.WriteUint16(101)
	h.HandlePacket(client, protocol.NewPacket(unlearnOK.Data()))
	if player.HasSkill(101) {
		t.Errorf("Unlearn con pergamino debería haber eliminado la skill 101")
	}
}
