package zone

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/database"
)

func TestFase1_SkillCatalogAndPlayerState(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	activeSHN, err := shn.ReadFile(filepath.Join(dataDir, "ActiveSkill.shn"))
	if err != nil {
		t.Fatalf("error leyendo ActiveSkill.shn: %v", err)
	}
	srvSHN, err := shn.ReadFile(filepath.Join(dataDir, "ActiveSkillInfoServer.shn"))
	if err != nil {
		t.Fatalf("error leyendo ActiveSkillInfoServer.shn: %v", err)
	}

	ds := &DataStore{
		ActiveSkillSHN: activeSHN,
		SHNTables: map[string]*shn.File{
			"ActiveSkill":           activeSHN,
			"ActiveSkillInfoServer": srvSHN,
		},
	}

	cat := BuildSkillCatalog(ds)
	if len(cat.ByID) == 0 {
		t.Fatalf("SkillCatalog.ByID vacío")
	}

	// Verificar que al menos las skills con empower y costos/rangos reales se cargaron
	var foundWithCost, foundWithEmpower bool
	for _, sk := range cat.ByID {
		if sk.SP > 0 && sk.Range > 0 && sk.DlyTimeMs > 0 {
			foundWithCost = true
		}
		if sk.EmpowerDamage[0] != 0 || sk.EmpowerSP[0] != 0 || sk.EmpowerCoolTime[0] != 0 {
			foundWithEmpower = true
		}
		if foundWithCost && foundWithEmpower {
			break
		}
	}
	if !foundWithCost {
		t.Errorf("no se encontró ninguna skill con SP, Range y DlyTimeMs > 0 en ActiveSkill.shn")
	}
	if !foundWithEmpower {
		t.Errorf("no se encontró ninguna skill con tablas de empower nT0..nT3 pobladas en ActiveSkill.shn")
	}

	// Verificar helpers de ZonePlayer
	ch := &database.Character{}
	ch.Level = 15
	p := NewZonePlayer(nil, ch, nil, nil)
	if p.MaxEmpowerPoints() != 14 {
		t.Errorf("MaxEmpowerPoints esperado 14, obtenido %d", p.MaxEmpowerPoints())
	}

	p.SyncLearnedSkills([]*database.CharacterSkill{
		{SkillNo: 101, SkillPowerDemage: 2, SkillPowerSP: 1, SkillPowerKeepTime: 0, SkillPowerCoolTime: 3},
	})
	if !p.HasSkill(101) || p.HasSkill(999) {
		t.Errorf("HasSkill(101)=%v HasSkill(999)=%v", p.HasSkill(101), p.HasSkill(999))
	}
	if p.TotalEmpowerUsed() != 6 {
		t.Errorf("TotalEmpowerUsed esperado 6, obtenido %d", p.TotalEmpowerUsed())
	}
	d, s, k, c := p.GetSkillEmpower(101)
	if d != 2 || s != 1 || k != 0 || c != 3 {
		t.Errorf("GetSkillEmpower inesperado: %d %d %d %d", d, s, k, c)
	}

	c1 := p.NextSkillCastIndex()
	c2 := p.NextSkillCastIndex()
	if c1 != 1 || c2 != 2 {
		t.Errorf("NextSkillCastIndex esperado 1, 2; obtenido %d, %d", c1, c2)
	}
}
