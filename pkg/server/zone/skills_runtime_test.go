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
	if player.HasAbState(26) {
		t.Fatalf("Resist 1 NO debe duplicar el portador oculto 26 (StaMentalShoutMain)")
	}
	if !player.HasAbState(27) {
		t.Fatalf("Player debe tener AbStataIndex 27 (StaMentalShout)")
	}
	if player.AbMA != 11 || player.AbMR != 11 {
		t.Fatalf("Resist 1 debe sumar +11 MA y +11 MR (got AbMA=%d AbMR=%d)", player.AbMA, player.AbMR)
	}

	// Invincible (DivineShield01, ID=1800) -> ShieldHP=1040 absorbe 500 de daño
	invincible := ds.SkillCatalog.Get(1800)
	if invincible == nil {
		t.Fatalf("DivineShield01 (1800) no encontrado en ActiveSkill.shn")
	}
	h.resolveSkillBlastOnEntity(player, player, invincible, 0, 0)
	if rem := player.AbsorbDamageWithShield(zMap, 500); rem != 0 {
		t.Fatalf("Invincible (ShieldHP=1040) debe absorber 500 de daño completo, got rem=%d", rem)
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

func TestRuntime4_SharedCooldownGroupAndClericRevive(t *testing.T) {
	h, zMap, cleric, clericClient := setupFase12TestEnvironment(t)
	h.dataStore = &DataStore{
		SkillCatalog: &SkillCatalog{
			ByID: map[uint16]*SkillDef{
				601:  {ID: 601, SP: 5, Range: 200, DlyTimeMs: 5000, DlyGroupNum: 77, DlyTimeGroupMs: 4000, MinMA: 30, MaxMA: 30},
				602:  {ID: 602, SP: 5, Range: 200, DlyTimeMs: 5000, DlyGroupNum: 77, DlyTimeGroupMs: 4000, MinMA: 30, MaxMA: 30},
				1620: {ID: 1620, SP: 10, Range: 400, DlyTimeMs: 5000, SkillTargetState: 1, Specials: [5]SkillSpecialEntry{{Index: 12}}},
			},
		},
	}
	mob := NewZoneMonster(60, 102, "GroupCDMob", "Rou", int32(cleric.Char.XPos+40), int32(cleric.Char.YPos), 0, 5, 500, 50, 10, 5, 10, 20, 30, 10, 50)
	zMap.AddMonster(mob)
	cleric.SyncLearnedSkills([]*database.CharacterSkill{
		{SkillNo: 601, SkillLevel: 1},
		{SkillNo: 602, SkillLevel: 1},
		{SkillNo: 1620, SkillLevel: 1},
	})
	cleric.Char.CurSP = 100

	// 1. Castear skill 601 (DlyGroupNum=77) debe bloquear skill 602 del mismo grupo
	p1 := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	p1.WriteUint16(601)
	p1.WriteUint16(mob.Handle)
	h.HandlePacket(clericClient, protocol.NewPacket(p1.Data()))
	hpAfter601 := mob.CurHP
	if hpAfter601 >= 500 {
		t.Fatalf("Skill 601 debería haber dañado al mob")
	}

	p2 := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	p2.WriteUint16(602)
	p2.WriteUint16(mob.Handle)
	h.HandlePacket(clericClient, protocol.NewPacket(p2.Data()))
	if mob.CurHP != hpAfter601 {
		t.Errorf("Skill 602 del mismo DlyGroupNum=77 debería estar bloqueada por cooldown compartido")
	}

	// 2. Cleric lanza Revival01 (1620, SkillTargetState=1) sobre un jugador muerto -> activa PendingRevive y envía 0x4829
	deadChar := &database.Character{ID: 777, Name: "DeadAlly", Level: 10, CurHP: 0, CurSP: 0}
	deadAlly := NewZonePlayer(nil, deadChar, nil, nil)
	deadAlly.CurrentMap = zMap
	deadAlly.MapObjectID = 10999
	deadAlly.MaxHP = 200
	deadAlly.MaxSP = 100
	deadAlly.State = 3 // Dead
	deadAlly.Char.XPos = cleric.Char.XPos + 30
	deadAlly.Char.YPos = cleric.Char.YPos
	zMap.AddPlayer(deadAlly)

	revPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	revPkt.WriteUint16(1620)
	revPkt.WriteUint16(deadAlly.MapObjectID)
	h.HandlePacket(clericClient, protocol.NewPacket(revPkt.Data()))

	if !deadAlly.PendingRevive || deadAlly.PendingReviveRate != 1000 {
		t.Fatalf("Revival01 sobre jugador muerto debe marcar PendingRevive=true (rate=1000), got pending=%v rate=%d",
			deadAlly.PendingRevive, deadAlly.PendingReviveRate)
	}
}
