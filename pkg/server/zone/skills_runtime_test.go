package zone

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
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
