package zone

import (
	"testing"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestFase4_MoverSkillBashObjAndFld(t *testing.T) {
	h, zMap, player, client := setupFase12TestEnvironment(t)
	h.dataStore = &DataStore{
		SkillCatalog: &SkillCatalog{
			ByID: map[uint16]*SkillDef{
				401: {ID: 401, SP: 50, Range: 200, DlyTimeMs: 3000, MinWC: 40, MaxWC: 40},
				402: {ID: 402, SP: 60, Range: 250, Area: 200, TargetNumber: 3, DlyTimeMs: 2000, MinWC: 30, MaxWC: 30},
			},
		},
	}

	player.State = 6 // Montado (bitmask 0x0006 válido en Mover)
	player.MoverHandle = 777
	player.Char.CurSP = 5 // Menor que SP de la skill (50/60): Mover NO debe cobrar SP del jugador!

	mob := NewZoneMonster(60, 102, "MoverTarget", "Rou", 50, 0, 0, 5, 200, 50, 0, 0, 0, 0, 0, 0, 50)
	mob.IsKillableByUser = true
	zMap.AddMonster(mob)

	// 1. Mover OBJ cast: debe hacer daño sin consumir SP y aplicar cooldown en MoverSkillCooldowns
	objPkt := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ)
	objPkt.WriteUint16(401)
	objPkt.WriteUint16(mob.Handle)
	h.HandlePacket(client, protocol.NewPacket(objPkt.Data()))

	if player.Char.CurSP != 5 {
		t.Errorf("Mover skill no debería consumir SP del jugador, CurSP=%d", player.Char.CurSP)
	}
	if mob.CurHP >= 200 {
		t.Errorf("Mover skill OBJ debería dañar al mob, CurHP=%d", mob.CurHP)
	}
	if _, hasCD := player.MoverSkillCooldowns[401]; !hasCD {
		t.Errorf("MoverSkillCooldowns[401] debería estar seteado")
	}

	// 2. Mob no matable -> 0x0FD6 (no recibe daño)
	mob2 := NewZoneMonster(61, 103, "InvulnMob", "Rou", 50, 0, 0, 5, 200, 50, 0, 0, 0, 0, 0, 0, 50)
	mob2.IsKillableByUser = false
	zMap.AddMonster(mob2)
	delete(player.MoverSkillCooldowns, 401)

	objPkt2 := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ)
	objPkt2.WriteUint16(401)
	objPkt2.WriteUint16(mob2.Handle)
	h.HandlePacket(client, protocol.NewPacket(objPkt2.Data()))
	if mob2.CurHP != 200 {
		t.Errorf("Mob no matable no debería recibir daño de Mover skill")
	}

	// 3. Mover FLD cast: AoE sin consumir SP
	hpBefore := mob.CurHP
	fldPkt := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ)
	fldPkt.WriteUint16(402)
	fldPkt.WriteInt32(50)
	fldPkt.WriteInt32(0)
	h.HandlePacket(client, protocol.NewPacket(fldPkt.Data()))
	if player.Char.CurSP != 5 {
		t.Errorf("Mover skill FLD no debería consumir SP, CurSP=%d", player.Char.CurSP)
	}
	if mob.CurHP >= hpBefore {
		t.Errorf("Mover skill FLD debería dañar al mob en radio AoE")
	}
}
