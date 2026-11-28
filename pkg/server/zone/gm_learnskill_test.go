package zone

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/database"
)

func TestGMCommand_LearnSkill_ByID_All_And_Reset(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	loadSHN := func(rel ...string) *shn.File {
		t.Helper()
		p := filepath.Join(append([]string{dataDir}, rel...)...)
		f, err := shn.ReadFile(p)
		if err != nil {
			t.Fatalf("error leyendo %s: %v", p, err)
		}
		return f
	}

	activeSHN := loadSHN("ActiveSkill.shn")
	viewSHN := loadSHN("View", "ActiveSkillView.shn")
	passiveSHN := loadSHN("PassiveSkill.shn")
	passiveViewSHN := loadSHN("View", "PassiveSkillView.shn")
	itemSHN := loadSHN("ItemInfo.shn")
	ucSHN := loadSHN("UseClassTypeInfo.shn")
	adminSHN := loadSHN("AdminLvSet.shn")

	passiveMap := make(map[uint16]bool, len(passiveSHN.Rows))
	for i := range passiveSHN.Rows {
		passiveMap[passiveSHN.Rows[i].GetUint16("ID")] = true
	}

	h, _, player, _ := setupFase12TestEnvironment(t)
	ds := &DataStore{
		ActiveSkillSHN:  activeSHN,
		PassiveSkillSHN: passiveSHN,
		PassiveSkillMap: passiveMap,
		SHNTables: map[string]*shn.File{
			"ActiveSkill":      activeSHN,
			"ActiveSkillView":  viewSHN,
			"PassiveSkill":     passiveSHN,
			"PassiveSkillView": passiveViewSHN,
			"ItemInfo":         itemSHN,
			"UseClassTypeInfo": ucSHN,
			"AdminLvSet":       adminSHN,
			"AdminLvSet.shn":   adminSHN,
		},
	}
	ds.AdminLvCatalog = BuildAdminLvCatalog(ds)
	if reqLv, ok := ds.AdminLvCatalog.RequiredLevel("learnskill"); !ok || reqLv != 50 {
		t.Fatalf("AdminLvSet.shn debe definir 'learnskill' con ALS_Lv=50 (got ok=%v lv=%d)", ok, reqLv)
	}
	ds.SkillCatalog = BuildSkillCatalog(ds)
	h.dataStore = ds

	// 0. Verificar fix crítico de colisión ActiveSkill/PassiveSkill (offset 0x7530 = 30000)
	if ds.IsPassiveSkill(0) || ds.IsPassiveSkill(1) || ds.IsPassiveSkill(20) || ds.IsPassiveSkill(60) {
		t.Fatalf("ActiveSkill IDs bajos (0, 1, 20, 60) NUNCA deben clasificarse como pasivas")
	}
	if !ds.IsPassiveSkill(PassiveSkillIDOffset + 0) {
		t.Fatalf("PassiveSkill DB ID 30000 (BraveMastery01) debe clasificarse como pasiva")
	}

	// Configurar GM como Fighter (Job=1) nivel 10 con AdminLevel=60
	player.Char.Job = 1
	player.Char.Level = 10
	player.Char.AdminLevel = 60

	// 1. &learnskill 1 (TripleHit02: Fighter UseClass=2, DemandLv=13 > player Lv.10) -> DEBE PERMITIRLO
	if !h.handleGMCommand(player, "&learnskill 1") {
		t.Fatalf("&learnskill 1 debería haber sido manejado como comando GM")
	}
	if !player.HasSkill(1) {
		t.Errorf("Fighter Lv.10 debería poder aprender TripleHit02 (ID=1, DemandLv=13) con &learnskill <id>")
	}

	// 2. Rechazar skills de otra clase (Cleric UseClass=8), monstruos (UseClass=0) y globales (UseClass=1 SummonMate01)
	var clericSkillID, mobSkillID uint16
	for _, sk := range ds.SkillCatalog.ByID {
		if clericSkillID == 0 && sk.UseClass == 8 {
			clericSkillID = sk.ID
		}
		if mobSkillID == 0 && sk.UseClass == 0 && sk.ID > 0 {
			mobSkillID = sk.ID
		}
		if clericSkillID != 0 && mobSkillID != 0 {
			break
		}
	}
	h.handleGMCommand(player, "&learnskill "+itoa(int(clericSkillID)))
	if player.HasSkill(clericSkillID) {
		t.Errorf("Fighter NO debería poder aprender skill de Cleric ID=%d", clericSkillID)
	}
	h.handleGMCommand(player, "&learnskill "+itoa(int(mobSkillID)))
	if player.HasSkill(mobSkillID) {
		t.Errorf("Fighter NO debería poder aprender skill de monstruo ID=%d", mobSkillID)
	}
	h.handleGMCommand(player, "&learnskill 15000") // SummonMate01 (UseClass=1)
	if player.HasSkill(15000) {
		t.Errorf("&learnskill <id> NO debería permitir aprender SummonMate01 (UseClass=1)")
	}

	// 3. &learnskill all en Fighter Lv.20 con starter skills cargadas:
	// Debe aprender exactamente 8 activas reales + 3 pasivas reales, sin SummonMate01, DT_* ni PowerofLove01.
	player.Char.Level = 20
	player.SyncLearnedSkills([]*database.CharacterSkill{
		{CharNo: player.Char.ID, SkillNo: 29200, SkillLevel: 1}, // starter emoji/action
		{CharNo: player.Char.ID, SkillNo: 29468, SkillLevel: 1}, // starter BallSnow01
	})
	if !h.handleGMCommand(player, "&learnskill all") {
		t.Fatalf("&learnskill all debería haber sido manejado como comando GM")
	}

	expectedActive := []uint16{0, 1, 20, 21, 40, 41, 60, 80}
	for _, sid := range expectedActive {
		if !player.HasSkill(sid) {
			t.Errorf("&learnskill all (Fighter Lv.20) debería haber aprendido activa ID=%d", sid)
		}
	}
	expectedPassive := []uint16{PassiveSkillIDOffset + 0, PassiveSkillIDOffset + 9, PassiveSkillIDOffset + 10}
	for _, pid := range expectedPassive {
		if !player.HasSkill(pid) {
			t.Errorf("&learnskill all (Fighter Lv.20) debería haber aprendido pasiva DB ID=%d", pid)
		}
	}
	for _, forbidden := range []uint16{15000, 29460, PassiveSkillIDOffset + 400} {
		if player.HasSkill(forbidden) {
			t.Errorf("&learnskill all NO debería haber aprendido skill prohibida ID=%d", forbidden)
		}
	}

	// 4. &learnskill reset: borra todas las activas y pasivas reales, mantiene starter skills (29200, 29468)
	if !h.handleGMCommand(player, "&learnskill reset") {
		t.Fatalf("&learnskill reset debería haber sido manejado como comando GM")
	}
	if !player.HasSkill(29200) || !player.HasSkill(29468) {
		t.Errorf("&learnskill reset NO debe borrar starter skills (29200, 29468)")
	}
	for _, sid := range expectedActive {
		if player.HasSkill(sid) {
			t.Errorf("&learnskill reset debería haber eliminado activa ID=%d", sid)
		}
	}
	for _, pid := range expectedPassive {
		if player.HasSkill(pid) {
			t.Errorf("&learnskill reset debería haber eliminado pasiva DB ID=%d", pid)
		}
	}
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	var buf [12]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
