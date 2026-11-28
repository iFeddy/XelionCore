package zone

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

func TestGMCommand_LearnSkill_ByIDAndAll(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	activeSHN, err := shn.ReadFile(filepath.Join(dataDir, "ActiveSkill.shn"))
	if err != nil {
		t.Fatalf("error leyendo ActiveSkill.shn: %v", err)
	}
	viewSHN, err := shn.ReadFile(filepath.Join(dataDir, "View", "ActiveSkillView.shn"))
	if err != nil {
		t.Fatalf("error leyendo ActiveSkillView.shn: %v", err)
	}
	ucSHN, err := shn.ReadFile(filepath.Join(dataDir, "UseClassTypeInfo.shn"))
	if err != nil {
		t.Fatalf("error leyendo UseClassTypeInfo.shn: %v", err)
	}
	adminSHN, err := shn.ReadFile(filepath.Join(dataDir, "AdminLvSet.shn"))
	if err != nil {
		t.Fatalf("error leyendo AdminLvSet.shn: %v", err)
	}

	h, _, player, _ := setupFase12TestEnvironment(t)
	ds := &DataStore{
		ActiveSkillSHN: activeSHN,
		SHNTables: map[string]*shn.File{
			"ActiveSkill":      activeSHN,
			"ActiveSkillView":  viewSHN,
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

	// 2. Buscar una skill exclusiva de Cleric (UseClass=8) y una skill de monstruo (UseClass=0) -> DEBE RECHAZARLAS
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
	if clericSkillID == 0 || mobSkillID == 0 {
		t.Fatalf("no se encontraron skills de prueba cleric=%d mob=%d", clericSkillID, mobSkillID)
	}

	h.handleGMCommand(player, "&learnskill "+itoa(int(clericSkillID)))
	if player.HasSkill(clericSkillID) {
		t.Errorf("Fighter NO debería poder aprender skill de Cleric ID=%d", clericSkillID)
	}
	h.handleGMCommand(player, "&learnskill "+itoa(int(mobSkillID)))
	if player.HasSkill(mobSkillID) {
		t.Errorf("Fighter NO debería poder aprender skill de monstruo ID=%d", mobSkillID)
	}

	// 3. &learnskill all: en Fighter Lv.10 debe aprender TripleHit01 (ID=0, DemandLv=3),
	// y NO aprender ninguna skill con DemandLv > 10 ni de otra clase.
	player.LearnedSkills = nil
	if !h.handleGMCommand(player, "&learnskill all") {
		t.Fatalf("&learnskill all debería haber sido manejado como comando GM")
	}
	if !player.HasSkill(0) {
		t.Errorf("&learnskill all debería haber aprendido TripleHit01 (ID=0, DemandLv=3)")
	}
	if player.HasSkill(1) {
		t.Errorf("&learnskill all NO debería haber aprendido TripleHit02 (ID=1, DemandLv=13 > Lv.10)")
	}
	for sid := range player.LearnedSkills {
		sk := ds.SkillCatalog.Get(sid)
		if sk == nil {
			t.Errorf("skill aprendida %d no existe en SkillCatalog", sid)
			continue
		}
		if !ds.SkillCatalog.CanClassUseSkill(sk.UseClass, player.Char.Job) {
			t.Errorf("&learnskill all aprendió skill %d (%s) de otra clase (UseClass=%d)", sid, sk.Name, sk.UseClass)
		}
		if sk.DemandLv <= 0 || sk.DemandLv > 10 {
			t.Errorf("&learnskill all aprendió skill %d (%s) con DemandLv=%d fuera de [1..10]", sid, sk.Name, sk.DemandLv)
		}
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
