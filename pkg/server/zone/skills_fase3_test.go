package zone

import (
	"testing"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestFase3_SkillDamageAndHealAndAoE(t *testing.T) {
	ch := &database.Character{}
	ch.Level = 10
	p := NewZonePlayer(nil, ch, nil, nil)
	p.Str = 100
	p.EquipMinWC = 50
	p.EquipMaxWC = 50
	p.Int = 80
	p.EquipMinMA = 20
	p.EquipMaxMA = 20

	// 1. Física: base=150, rate=200 (20% -> +30), add=20, empDmg[2]=15 => 215
	physSkill := &SkillDef{
		ID:            10,
		MinWC:         20,
		MaxWC:         20,
		MinWCRate:     200,
		MaxWCRate:     200,
		EmpowerDamage: [5]int{5, 10, 15, 20, 25},
	}
	gotPhys := calcOfficialSkillAttackPower(p, physSkill, 3)
	if gotPhys != 215 {
		t.Errorf("phys attack esperado 215, obtenido %d", gotPhys)
	}

	// 2. Mágica: base=100, rate=500 (50% -> +50), add=40, empDmg[1]=12 => 202
	magSkill := &SkillDef{
		ID:            11,
		MinMA:         40,
		MaxMA:         40,
		MinMARate:     500,
		MaxMARate:     500,
		EmpowerDamage: [5]int{6, 12, 18, 24, 30},
	}
	gotMag := calcOfficialSkillAttackPower(p, magSkill, 2)
	if gotMag != 202 {
		t.Errorf("mag attack esperado 202, obtenido %d", gotMag)
	}

	// 3. Heal: base=120, empDmg[0]=25 => 145
	healSkill := &SkillDef{
		ID:            12,
		Specials:      [5]SkillSpecialEntry{{Index: 1, Value: 120}},
		EmpowerDamage: [5]int{25, 50, 75, 100, 125},
	}
	if gotHeal := calcOfficialSkillHeal(healSkill, 1); gotHeal != 145 {
		t.Errorf("heal esperado 145, obtenido %d", gotHeal)
	}

	// 4. AoE: Area=200 -> radio=100, TargetNumber=2, ordenamiento por distancia
	zMap := NewZoneMap(0, "Rou", "Roumen")
	mFar := NewZoneMonster(1, 101, "Far", "Rou", 150, 0, 0, 1, 100, 10, 5, 5, 5, 5, 5, 5, 10)
	mMid := NewZoneMonster(2, 101, "Mid", "Rou", 80, 0, 0, 1, 100, 10, 5, 5, 5, 5, 5, 5, 10)
	mNear := NewZoneMonster(3, 101, "Near", "Rou", 30, 0, 0, 1, 100, 10, 5, 5, 5, 5, 5, 5, 10)
	mMid2 := NewZoneMonster(4, 101, "Mid2", "Rou", 50, 0, 0, 1, 100, 10, 5, 5, 5, 5, 5, 5, 10)
	zMap.AddMonster(mFar)
	zMap.AddMonster(mMid)
	zMap.AddMonster(mNear)
	zMap.AddMonster(mMid2)

	aoeSkill := &SkillDef{ID: 13, Area: 200, TargetNumber: 2, SkillDegree: 360}
	picked := selectAoEMonsters(zMap, 0, 0, 0, 0, aoeSkill)
	if len(picked) != 2 {
		t.Fatalf("AoE TargetNumber=2 esperado 2 mobs, obtenido %d", len(picked))
	}
	if picked[0].Handle != mNear.Handle || picked[1].Handle != mMid2.Handle {
		t.Errorf("AoE orden por distancia esperado [Near(3), Mid2(4)], obtenido [%d, %d]", picked[0].Handle, picked[1].Handle)
	}

	// 5. Sector cónico (SkillDegree=90 hacia +X) excluye mob detrás (-X)
	mBehind := NewZoneMonster(5, 101, "Behind", "Rou", -40, 0, 0, 1, 100, 10, 5, 5, 5, 5, 5, 5, 10)
	zMap.AddMonster(mBehind)
	coneSkill := &SkillDef{ID: 14, Area: 200, TargetNumber: 10, SkillDegree: 90}
	pickedCone := selectAoEMonsters(zMap, 0, 0, 50, 0, coneSkill)
	for _, m := range pickedCone {
		if m.Handle == mBehind.Handle {
			t.Errorf("Mob detrás del cono de 90 grados no debería ser seleccionado")
		}
	}
}

func TestFase3_CastAbortAndOneDmgShield(t *testing.T) {
	h, zMap, player, client := setupFase12TestEnvironment(t)

	// AbState 0x16D en mob fuerza daño = 1 (Zone.exe 0x00505CCD)
	mob := NewZoneMonster(60, 102, "ShieldMob", "Rou", 30, 0, 0, 5, 200, 50, 0, 0, 0, 0, 0, 0, 50)
	mob.ActiveStates = map[uint16]*ActiveAbState{0x16D: {ID: 0x16D}}
	zMap.AddMonster(mob)

	sk := &SkillDef{ID: 200, MinWC: 500, MaxWC: 500}
	res, dead := h.resolveSkillBlastOnEntity(player, mob, sk, 0, 0)
	if dead || res.Damage != 1 || mob.CurHP != 199 {
		t.Errorf("AbState 0x16D esperado dmg=1 CurHP=199, obtenido dmg=%d CurHP=%d", res.Damage, mob.CurHP)
	}

	// CastAbort limpia CastingSkillID en estado normal
	player.CastingSkillID = 200
	abortPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ)
	h.HandlePacket(client, protocol.NewPacket(abortPkt.Data()))
	if player.CastingSkillID != 0 {
		t.Errorf("CastingSkillID debería limpiarse a 0 tras CASTABORT_REQ")
	}
}
