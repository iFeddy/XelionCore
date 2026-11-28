package zone

import (
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestFase2_OfficialSkillFormulas(t *testing.T) {
	// 1. Rango oficial: min((base * 10) >> 3, base + 60)
	// base = 100 -> scaled = 125 (diff 25 <= 60) => 125
	if got := calcOfficialSkillMaxRange(0, 100); got != 125 {
		t.Errorf("range(0,100) esperado 125, obtenido %d", got)
	}
	// base = 400 -> scaled = 500 (diff 100 > 60) => 460
	if got := calcOfficialSkillMaxRange(0, 400); got != 460 {
		t.Errorf("range(0,400) esperado 460, obtenido %d", got)
	}
	// objSize = 30, skillRange = 200 -> base = 230, scaled = 287 (diff 57 <= 60) => 287
	if got := calcOfficialSkillMaxRange(30, 200); got != 287 {
		t.Errorf("range(30,200) esperado 287, obtenido %d", got)
	}

	// 2. Costos HP/SP/LP con empower track 1
	sk := &SkillDef{
		ID:        500,
		SP:        30,
		SPRate:    100, // 10% de MaxSP
		HP:        10,
		HPRate:    50, // 5% de MaxHP
		LP:        5,
		DlyTimeMs: 4000,
		EmpowerSP: [5]int{3, 6, 9, 12, 15},
		EmpowerCoolTime: [5]int{500, 1000, 1500, 2000, 2500},
	}
	hpCost, spCost, lpCost := calcOfficialSkillCosts(sk, 1000, 500, 2)
	// HP = 10 + (1000*50)/1000 = 60
	// SP = 30 + (500*100)/1000 - 6 = 74
	// LP = 5
	if hpCost != 60 || spCost != 74 || lpCost != 5 {
		t.Errorf("costs inesperados: hp=%d sp=%d lp=%d (esperado 60, 74, 5)", hpCost, spCost, lpCost)
	}

	// 3. Cooldown oficial con empower track 3
	cd := calcOfficialSkillCooldown(sk, 3)
	if cd != 2500*time.Millisecond {
		t.Errorf("cooldown esperado 2500ms, obtenido %v", cd)
	}
}

func TestFase2_CastValidationsAndErrorCodes(t *testing.T) {
	h, zMap, player, client := setupFase12TestEnvironment(t)
	h.dataStore = &DataStore{
		SkillCatalog: &SkillCatalog{
			ByID: map[uint16]*SkillDef{
				101: {ID: 101, SP: 20, Range: 150, DlyTimeMs: 2000},
				102: {ID: 102, HP: 500, Range: 150, DlyTimeMs: 1000},
				103: {ID: 103, LP: 10, Range: 150, DlyTimeMs: 1000},
			},
		},
		PassiveSkillMap: map[uint16]bool{900: true},
	}

	mob := NewZoneMonster(60, 102, "TargetMob", "Rou", int32(player.Char.XPos+40), int32(player.Char.YPos), 0, 5, 200, 50, 10, 5, 10, 20, 30, 10, 50)
	zMap.AddMonster(mob)

	player.SyncLearnedSkills([]*database.CharacterSkill{
		{SkillNo: 101, SkillLevel: 1},
		{SkillNo: 102, SkillLevel: 1},
		{SkillNo: 103, SkillLevel: 1},
	})
	player.Char.CurHP = 100
	player.Char.CurSP = 100

	// 1. Silence -> 0x0FC7 (no 0x0FCC)
	player.AddAbState(&ActiveAbState{ID: 51, IsSilence: true, ExpiresAt: time.Now().Add(time.Minute)})
	pkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	pkt.WriteUint16(101)
	pkt.WriteUint16(mob.Handle)
	h.HandlePacket(client, protocol.NewPacket(pkt.Data()))
	if mob.CurHP != 200 {
		t.Errorf("Silenced player no debería dañar al mob")
	}
	player.RemoveAbState(51)

	// 2. Skill pasiva -> 0x0FC6
	pkt = protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	pkt.WriteUint16(900)
	pkt.WriteUint16(mob.Handle)
	h.HandlePacket(client, protocol.NewPacket(pkt.Data()))
	if mob.CurHP != 200 {
		t.Errorf("Skill pasiva no debería ejecutarse")
	}

	// 3. HP insuficiente (CurHP 100 <= HPCost 500) -> 0x0FD5
	pkt = protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	pkt.WriteUint16(102)
	pkt.WriteUint16(mob.Handle)
	h.HandlePacket(client, protocol.NewPacket(pkt.Data()))
	if mob.CurHP != 200 {
		t.Errorf("Skill con HPCost > CurHP no debería ejecutarse")
	}

	// 4. LP insuficiente -> 0x0FDC
	pkt = protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	pkt.WriteUint16(103)
	pkt.WriteUint16(mob.Handle)
	h.HandlePacket(client, protocol.NewPacket(pkt.Data()))
	if mob.CurHP != 200 {
		t.Errorf("Skill con LPCost > 0 no debería ejecutarse")
	}

	// 5. Skill válida 101 consume 20 SP y aplica cooldown de 2000ms
	pkt = protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	pkt.WriteUint16(101)
	pkt.WriteUint16(mob.Handle)
	h.HandlePacket(client, protocol.NewPacket(pkt.Data()))
	if player.Char.CurSP != 80 {
		t.Errorf("CurSP esperado 80 (100-20), obtenido %d", player.Char.CurSP)
	}
	if mob.CurHP >= 200 {
		t.Errorf("Mob debería haber recibido daño de skill 101")
	}
	if player.NextCastIndex != 1 {
		t.Errorf("NextCastIndex esperado 1, obtenido %d", player.NextCastIndex)
	}

	// 6. Segundo casteo inmediato debe ser rechazado por cooldown (0x0FC8)
	hpAfterFirst := mob.CurHP
	pkt = protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	pkt.WriteUint16(101)
	pkt.WriteUint16(mob.Handle)
	h.HandlePacket(client, protocol.NewPacket(pkt.Data()))
	if mob.CurHP != hpAfterFirst {
		t.Errorf("Segundo casteo en cooldown no debería haber hecho daño")
	}
}
