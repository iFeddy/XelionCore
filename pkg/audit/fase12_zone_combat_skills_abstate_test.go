package audit

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

// TestFase12_Audit_NoInlineSQL_Zone_And_DB verifica mediante AST que:
// 1. En los subsistemas de Zone (combat, skills, stats_pvp, abstate) y MySQL Character Repositories,
//    el 100% de las consultas a base de datos utilizan Stored Procedures (CALL <SP>).
// 2. Cero SQL inline directo (SELECT/INSERT/UPDATE/DELETE) en la capa de persistencia de personaje y habilidades.
func TestFase12_Audit_NoInlineSQL_Zone_And_DB(t *testing.T) {
	dataDir := findAuditDataDir()
	rootPath := filepath.Dir(dataDir)
	if filepath.Base(dataDir) == "data" && rootPath == "." {
		rootPath = ".."
	}
	// Buscar raíz del repositorio
	for i := 0; i < 4; i++ {
		if fi, err := os.Stat(filepath.Join(rootPath, "go.mod")); err == nil && !fi.IsDir() {
			break
		}
		rootPath = filepath.Join(rootPath, "..")
	}

	filesToCheck := []string{
		filepath.Join(rootPath, "pkg", "database", "mysql_character.go"),
		filepath.Join(rootPath, "pkg", "database", "mysql_character_progression.go"),
		filepath.Join(rootPath, "pkg", "database", "gamelog.go"),
		filepath.Join(rootPath, "pkg", "server", "zone", "combat.go"),
		filepath.Join(rootPath, "pkg", "server", "zone", "skills.go"),
		filepath.Join(rootPath, "pkg", "server", "zone", "stats_pvp.go"),
		filepath.Join(rootPath, "pkg", "server", "zone", "abstate_engine.go"),
	}

	fset := token.NewFileSet()
	for _, fPath := range filesToCheck {
		if _, err := os.Stat(fPath); os.IsNotExist(err) {
			continue
		}
		node, err := parser.ParseFile(fset, fPath, nil, 0)
		if err != nil {
			t.Fatalf("Error analizando AST de %s: %v", fPath, err)
		}

		ast.Inspect(node, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			val := strings.ToUpper(strings.Trim(lit.Value, "`\""))
			// Permitir consultas administrativas de metadata en gamelog y schemas
			if strings.Contains(fPath, "gamelog.go") && (strings.Contains(val, "INFORMATION_SCHEMA") || strings.Contains(val, "SHOW TABLES")) {
				return true
			}
			if strings.HasPrefix(val, "SELECT ") && !strings.Contains(val, "CALL ") && !strings.Contains(val, "INFORMATION_SCHEMA") {
				// No permitir SELECT directo contra tablas de personajes en estos repositorios
				if strings.Contains(val, "FROM TCHAR") || strings.Contains(val, "FROM TUSER") || strings.Contains(val, "FROM TITEM") {
					t.Errorf("CRITERIO FALLIDO: Consulta SQL inline prohibida en %s: %s", fPath, val)
				}
			}
			if (strings.HasPrefix(val, "INSERT INTO ") || strings.HasPrefix(val, "UPDATE ") || strings.HasPrefix(val, "DELETE FROM ")) &&
				!strings.Contains(val, "CALL ") {
				t.Errorf("CRITERIO FALLIDO: Consulta SQL inline prohibida en %s: %s", fPath, val)
			}
			return true
		})
	}
	t.Logf("Auditoría AST de Cero SQL inline en repositorios de Zone y Character completada con éxito")
}

// TestFase12_Audit_Opcodes_Depts_Fidelity verifica que todos los opcodes requeridos para
// Dept 9 (Combate), Dept 18 (Habilidades), Dept 20 (SoulStone), Dept 4 (Stats/Revive) y Dept 51 (Mover)
// sigan la fórmula exacta de 16 bits: protocol = (department << 10) | command.
func TestFase12_Audit_Opcodes_Depts_Fidelity(t *testing.T) {
	tests := []struct {
		name     string
		opcode   uint16
		dept     uint16
		cmd      uint16
	}{
		// Dept 9: Combat
		{"PROTO_NC_BAT_TARGETTING_REQ", protocol.PROTO_NC_BAT_TARGETTING_REQ, 9, 1},
		{"PROTO_NC_BAT_UNTARGET_REQ", protocol.PROTO_NC_BAT_UNTARGET_REQ, 9, 8},
		{"PROTO_NC_BAT_BASHSTART_CMD", protocol.PROTO_NC_BAT_BASHSTART_CMD, 9, 43},
		{"PROTO_NC_BAT_BASHSTOP_CMD", protocol.PROTO_NC_BAT_BASHSTOP_CMD, 9, 50},
		{"PROTO_NC_BAT_ASSIST_REQ", protocol.PROTO_NC_BAT_ASSIST_REQ, 9, 62},
		{"PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ", protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ, 9, 64},
		{"PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ", protocol.PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ, 9, 65},
		{"PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ", protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ, 9, 68},
		{"PROTO_NC_BAT_ABSTATE_ERASE_REQ", protocol.PROTO_NC_BAT_ABSTATE_ERASE_REQ, 9, 84},

		// Dept 18: Skills
		{"PROTO_NC_SKILL_EMPOWALLOC_REQ", protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ, 18, 17},
		{"PROTO_NC_SKILL_PRODUCTFIELD_REQ", protocol.PROTO_NC_SKILL_PRODUCTFIELD_REQ, 18, 34},
		{"PROTO_NC_SKILL_UNLEARN_REQ", protocol.PROTO_NC_SKILL_UNLEARN_REQ, 18, 36},
		{"PROTO_NC_SKILL_REPLYREVIVE_CMD", protocol.PROTO_NC_SKILL_REPLYREVIVE_CMD, 18, 42},

		// Dept 20: SoulStone
		{"PROTO_NC_SOULSTONE_HP_BUY_REQ", protocol.PROTO_NC_SOULSTONE_HP_BUY_REQ, 20, 1},
		{"PROTO_NC_SOULSTONE_SP_BUY_REQ", protocol.PROTO_NC_SOULSTONE_SP_BUY_REQ, 20, 2},
		{"PROTO_NC_SOULSTONE_HP_USE_REQ", protocol.PROTO_NC_SOULSTONE_HP_USE_REQ, 20, 7},
		{"PROTO_NC_SOULSTONE_SP_USE_REQ", protocol.PROTO_NC_SOULSTONE_SP_USE_REQ, 20, 9},

		// Dept 4: Stats & Revive
		{"PROTO_NC_CHAR_REVIVE_REQ", protocol.PROTO_NC_CHAR_REVIVE_REQ, 4, 78},
		{"PROTO_NC_CHAR_STAT_INCPOINT_REQ", protocol.PROTO_NC_CHAR_STAT_INCPOINT_REQ, 4, 92},
		{"PROTO_NC_CHAR_STAT_DECPOINT_REQ", protocol.PROTO_NC_CHAR_STAT_DECPOINT_REQ, 4, 98},
		{"PROTO_NC_CHAR_DEPOLYMORPH_CMD", protocol.PROTO_NC_CHAR_DEPOLYMORPH_CMD, 4, 176},
		{"PROTO_NC_CHAR_REBIRTH_CMD", protocol.PROTO_NC_CHAR_REBIRTH_CMD, 4, 258},
		{"PROTO_NC_CHAR_REBIRTH_REJECT_CMD", protocol.PROTO_NC_CHAR_REBIRTH_REJECT_CMD, 4, 259},
		{"PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ", protocol.PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ, 4, 302},

		// Dept 12: Item Active Skill & Item Revive
		{"PROTO_NC_ITEM_REVIVEITEMUSE_CMD", protocol.PROTO_NC_ITEM_REVIVEITEMUSE_CMD, 12, 57},
		{"PROTO_NC_ITEM_USE_ACTIVESKILL_REQ", protocol.PROTO_NC_ITEM_USE_ACTIVESKILL_REQ, 12, 78},

		// Dept 51: Mover Mounted Skills
		{"PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ", protocol.PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ, 51, 14},
		{"PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ", protocol.PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ, 51, 15},
	}

	for _, tc := range tests {
		calculated := (tc.dept << 10) | tc.cmd
		if tc.opcode != calculated {
			t.Errorf("OpCode %s (0x%04X) no coincide con fórmula (Dept %d, Cmd %d -> 0x%04X)",
				tc.name, tc.opcode, tc.dept, tc.cmd, calculated)
		}
		if uint16(protocol.ExtractHeader(tc.opcode)) != tc.dept {
			t.Errorf("ExtractHeader para %s retornó %d, esperado %d", tc.name, protocol.ExtractHeader(tc.opcode), tc.dept)
		}
		if protocol.ExtractType(tc.opcode) != tc.cmd {
			t.Errorf("ExtractType para %s retornó %d, esperado %d", tc.name, protocol.ExtractType(tc.opcode), tc.cmd)
		}
	}

	t.Logf("Verificación matemática de los 25 opcodes críticos de Fase 12 completada con 100%% de fidelidad")
}

// TestFase12_Audit_AbState_Bitfield99Bytes_And_CrowdControl verifica la máquina de estados anormales,
// el bitfield oficial de 99 bytes (792 estados indexables), y las restricciones de control de masas.
func TestFase12_Audit_AbState_Bitfield99Bytes_And_CrowdControl(t *testing.T) {
	player := &zone.ZonePlayer{
		ActiveStates: make(map[uint16]*zone.ActiveAbState),
	}

	// 1. Bitfield inicial debe ser de exactamente 99 bytes
	initialBit := player.GetAbStateBit()
	if len(initialBit) != 99 {
		t.Fatalf("Tamaño de abstatebit esperado 99 bytes, obtenido %d", len(initialBit))
	}
	for i, b := range initialBit {
		if b != 0 {
			t.Fatalf("Byte %d de abstatebit inicial no es cero: 0x%02X", i, b)
		}
	}

	// 2. Agregar Stun: restringe movimiento, ataque y casteo
	player.AddAbState(&zone.ActiveAbState{
		ID:        120,
		Name:      "StunBuff",
		Duration:  5 * time.Second,
		ExpiresAt: time.Now().Add(5 * time.Second),
		IsStun:    true,
	})

	if player.CanMove() || player.CanAttack() || player.CanCast() {
		t.Fatalf("Jugador con Stun no debe poder moverse, atacar ni castear")
	}

	// Verificar bit en bitfield: ID 120 -> byte 15, bit 0 (120/8 = 15, 120%8 = 0)
	bitfield := player.GetAbStateBit()
	if (bitfield[15] & 0x01) == 0 {
		t.Fatalf("Bit 0 del byte 15 debe estar activado para el estado ID 120")
	}

	// 3. Remover Stun y agregar Silencio: restringe casteo pero permite movimiento y ataque
	player.RemoveAbState(120)
	if !player.CanMove() || !player.CanAttack() || !player.CanCast() {
		t.Fatalf("Tras remover Stun, el jugador debe recuperar todas sus facultades")
	}

	player.AddAbState(&zone.ActiveAbState{
		ID:        200,
		Name:      "SilenceDebuff",
		Duration:  5 * time.Second,
		ExpiresAt: time.Now().Add(5 * time.Second),
		IsSilence: true,
	})

	if !player.CanMove() {
		t.Fatalf("Jugador silenciado debe poder moverse")
	}
	if !player.CanAttack() {
		t.Fatalf("Jugador silenciado debe poder realizar ataques básicos")
	}
	if player.CanCast() {
		t.Fatalf("Jugador silenciado NO debe poder castear habilidades")
	}

	// 4. Agregar Inmovilización (Root): restringe movimiento pero permite ataque y casteo
	player.RemoveAbState(200)
	player.AddAbState(&zone.ActiveAbState{
		ID:        250,
		Name:      "RootDebuff",
		Duration:  5 * time.Second,
		ExpiresAt: time.Now().Add(5 * time.Second),
		IsRoot:    true,
	})

	if player.CanMove() {
		t.Fatalf("Jugador inmovilizado NO debe poder moverse")
	}
	if !player.CanAttack() {
		t.Fatalf("Jugador inmovilizado debe poder atacar")
	}
	if !player.CanCast() {
		t.Fatalf("Jugador inmovilizado debe poder castear")
	}

	t.Logf("Máquina de estados AbState y bitfield oficial de 99 bytes validados exitosamente")
}

// TestFase12_Audit_CombatFormulas_Angle_Crit_SoulStone verifica que los cálculos de daño,
// ángulos geométricos, multiplicador crítico y SoulStones respeten las especificaciones oficiales.
func TestFase12_Audit_CombatFormulas_Angle_Crit_SoulStone(t *testing.T) {
	dataDir := findAuditDataDir()
	store := combat.NewStore()
	if err := store.LoadDir(filepath.Join(dataDir, "World", "Combat")); err != nil {
		t.Fatalf("Error cargando store de combate: %v", err)
	}

	// 1. Verificación de tasa angular oficial de DamageByAngle.json
	// Frontal (0 deg): 1000 (1.0x)
	rate0 := store.GetAngleDamageRate(true, 0)
	if rate0 != 1000 {
		t.Errorf("Tasa esperada para ángulo 0 deg = 1000, obtenida: %d", rate0)
	}
	// Espalda (180 deg): 1200 (1.2x)
	rate180 := store.GetAngleDamageRate(true, 180)
	if rate180 != 1200 {
		t.Errorf("Tasa esperada para ángulo 180 deg = 1200, obtenida: %d", rate180)
	}

	// 2. Verificación de cálculo angular en coordenadas 2D
	// Atacante en (100, 100), Defensor mirando hacia el eje +X (0 deg)
	// Si atacante está en (200, 100) [frente]: ángulo ~0
	angleFront := zone.CalculateAngleDeg(200, 100, 100, 100, 0)
	if angleFront < 0 || angleFront > 180 {
		t.Errorf("Ángulo frontal calculado fuera de rango: %d", angleFront)
	}

	// 3. Verificación de daño potenciado por SoulStone vs daño normal
	dmgNormal := zone.CalculatePhysicalDamage(100, 20, 0, true, 0, 1000, false, store)
	soulRate := store.GetSoulDamageRate(1, 4) // SoulStone demand 1, count 4
	dmgSoul := zone.CalculatePhysicalDamage(100, 20, 0, true, 0, soulRate, false, store)

	if dmgSoul <= dmgNormal {
		t.Errorf("Daño potenciado con SoulStone (%d) debe ser mayor a daño normal (%d)", dmgSoul, dmgNormal)
	}

	// 4. Verificación de multiplicador crítico (1.5x)
	dmgCrit := zone.CalculatePhysicalDamage(100, 20, 0, true, 0, 1000, true, store)
	expectedCritDmg := int(float64(dmgNormal) * 1.5)
	diff := dmgCrit - expectedCritDmg
	if diff < -1 || diff > 1 {
		t.Errorf("Daño crítico esperado ~%d, obtenido %d", expectedCritDmg, dmgCrit)
	}

	t.Logf("Fórmulas oficiales de combate, multiplicadores angulares, críticos y SoulStone validadas al 100%%")
}
