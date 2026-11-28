package zone

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestRuntime1_AbStateEffectsAndExpiration(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	loadSHN := func(rel ...string) *shn.File {
		t.Helper()
		f, err := shn.ReadFile(filepath.Join(append([]string{dataDir}, rel...)...))
		if err != nil {
			t.Fatalf("error leyendo %v: %v", rel, err)
		}
		return f
	}

	activeSHN := loadSHN("ActiveSkill.shn")
	abSHN := loadSHN("Abstate.shn")
	subSHN := loadSHN("SubAbstate.shn")

	h, zMap, player, _ := setupFase12TestEnvironment(t)
	ds := &DataStore{
		ActiveSkillSHN: activeSHN,
		SHNTables: map[string]*shn.File{
			"ActiveSkill": activeSHN,
			"Abstate":     abSHN,
			"SubAbstate":  subSHN,
		},
	}
	ds.SkillCatalog = BuildSkillCatalog(ds)
	h.dataStore = ds

	// MentalShout01 (Resist 1, ID=1760) -> StaMentalShoutMain (AbStataIndex=26) + PartyState1 StaMentalShout (AbStataIndex=27, +11 MA, +11 MR)
	resist1 := ds.SkillCatalog.Get(1760)
	if resist1 == nil {
		t.Fatalf("MentalShout01 (1760) no encontrado en ActiveSkill.shn")
	}
	if !resist1.IsPureBuff() {
		t.Fatalf("MentalShout01 debe clasificarse como IsPureBuff()")
	}

	res, dead := h.resolveSkillBlastOnEntity(player, player, resist1, 0, 0)
	if dead || res.Damage != 0 || !res.IsAbState {
		t.Fatalf("Resist 1 debe aplicar AbState sin daño: got damage=%d dead=%v isAbState=%v", res.Damage, dead, res.IsAbState)
	}
	if !player.HasAbState(26) || !player.HasAbState(27) {
		t.Fatalf("Player debe tener AbStataIndex 26 (StaMentalShoutMain) y 27 (StaMentalShout)")
	}
	if player.AbMA != 11 || player.AbMR != 11 {
		t.Fatalf("Resist 1 debe sumar +11 MA y +11 MR (got AbMA=%d AbMR=%d)", player.AbMA, player.AbMR)
	}

	// Verificar expiración en tick1s() de MapGameLoop
	player.ActiveStates[27].ExpiresAt = time.Now().Add(-time.Second)
	gl := NewMapGameLoop(zMap)
	gl.tick1s()
	if player.HasAbState(27) {
		t.Errorf("AbState 27 debería haber expirado tras tick1s()")
	}
	if player.AbMA != 0 || player.AbMR != 0 {
		t.Errorf("Tras expirar AbState 27, AbMA/AbMR deben volver a 0 (got %d/%d)", player.AbMA, player.AbMR)
	}
}

func TestRuntime2_PassiveSkillBonuses(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	loadSHN := func(rel ...string) *shn.File {
		t.Helper()
		f, err := shn.ReadFile(filepath.Join(append([]string{dataDir}, rel...)...))
		if err != nil {
			t.Fatalf("error leyendo %v: %v", rel, err)
		}
		return f
	}

	activeSHN := loadSHN("ActiveSkill.shn")
	passiveSHN := loadSHN("PassiveSkill.shn")
	itemSHN := loadSHN("ItemInfo.shn")

	h, _, player, _ := setupFase12TestEnvironment(t)
	ds := &DataStore{
		ActiveSkillSHN:  activeSHN,
		PassiveSkillSHN: passiveSHN,
		SHNTables: map[string]*shn.File{
			"ActiveSkill":  activeSHN,
			"PassiveSkill": passiveSHN,
			"ItemInfo":     itemSHN,
		},
	}
	ds.SkillCatalog = BuildSkillCatalog(ds)
	h.dataStore = ds

	// Buscar un ítem 1H Sword (WeaponType == 1) en ItemInfo.shn
	var swordID uint16
	for id, wt := range ds.SkillCatalog.ItemWeaponType {
		if wt == 1 {
			swordID = id
			break
		}
	}
	if swordID == 0 {
		t.Fatalf("no se encontró espada 1H (WeaponType=1) en ItemInfo.shn")
	}

	// Aprender BraveMastery01 (ID=0 -> MstRtSword1=100) + OHSwdMastery01 (ID=9 -> +3) + OHSwdMastery02 (ID=10 -> +7)
	player.SyncLearnedSkills([]*database.CharacterSkill{
		{SkillNo: PassiveSkillIDOffset + 0, SkillLevel: 1},
		{SkillNo: PassiveSkillIDOffset + 9, SkillLevel: 1},
		{SkillNo: PassiveSkillIDOffset + 10, SkillLevel: 1},
	})

	// Sin arma equipada -> maestría de espada no aplica
	player.Char.Equips = map[byte]uint16{}
	unarmed := ds.SkillCatalog.ComputePassiveBonuses(player)
	if unarmed.WeaponFlatPlus != 0 || unarmed.WeaponRatePermille != 0 {
		t.Errorf("sin arma no debería aplicar Sword1 mastery (got flat=%d rate=%d)", unarmed.WeaponFlatPlus, unarmed.WeaponRatePermille)
	}

	// Con espada 1H equipada -> max(3,7)=7 plano + 100 permille (+10%)
	player.Char.Equips = map[byte]uint16{10: swordID}
	armed := ds.SkillCatalog.ComputePassiveBonuses(player)
	if armed.WeaponFlatPlus != 7 || armed.WeaponRatePermille != 100 {
		t.Errorf("con espada 1H esperado flat=7 rate=100, got flat=%d rate=%d", armed.WeaponFlatPlus, armed.WeaponRatePermille)
	}
}

func TestRuntime3_CastTimeDelayAndAbort(t *testing.T) {
	h, zMap, player, client := setupFase12TestEnvironment(t)
	h.dataStore = &DataStore{
		SkillCatalog: &SkillCatalog{
			ByID: map[uint16]*SkillDef{
				501: {ID: 501, SP: 10, Range: 300, CastTimeMs: 1200, DlyTimeMs: 3000, MinMA: 50, MaxMA: 50, IsMovingSkill: false},
			},
		},
	}
	mob := NewZoneMonster(60, 102, "CastTarget", "Rou", int32(player.Char.XPos+40), int32(player.Char.YPos), 0, 5, 300, 50, 10, 5, 10, 20, 30, 10, 50)
	zMap.AddMonster(mob)
	player.SyncLearnedSkills([]*database.CharacterSkill{{SkillNo: 501, SkillLevel: 1}})
	player.Char.CurSP = 100

	// 1. Iniciar casteo con CastTime=1200ms -> NO debe dañar inmediatamente
	pkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	pkt.WriteUint16(501)
	pkt.WriteUint16(mob.Handle)
	h.HandlePacket(client, protocol.NewPacket(pkt.Data()))
	if player.PendingCast == nil || player.CastingSkillID != 501 {
		t.Fatalf("Skill con CastTimeMs=1200 debe quedar en PendingCast")
	}
	if mob.CurHP != 300 {
		t.Fatalf("Skill con CastTimeMs=1200 no debe dañar antes de terminar el cast (mob HP=%d)", mob.CurHP)
	}

	// 2. Moverse con !IsMovingSkill cancela el casteo
	h.abortNonMovingCastOnMove(player)
	if player.PendingCast != nil || player.CastingSkillID != 0 {
		t.Fatalf("Moverse con !IsMovingSkill debe cancelar PendingCast")
	}

	// 3. Volver a castear y completar en FinishPendingSkillCast al vencer ExecuteAt
	h.HandlePacket(client, protocol.NewPacket(pkt.Data()))
	if player.PendingCast == nil {
		t.Fatalf("Segundo casteo debería crear PendingCast")
	}
	h.FinishPendingSkillCast(player, player.PendingCast.ExecuteAt.Add(10*time.Millisecond))
	if player.PendingCast != nil || player.CastingSkillID != 0 {
		t.Errorf("Tras FinishPendingSkillCast, PendingCast y CastingSkillID deben limpiarse")
	}
	if mob.CurHP >= 300 {
		t.Errorf("Tras completar CastTime, el mob debería haber recibido daño (HP=%d)", mob.CurHP)
	}
}
