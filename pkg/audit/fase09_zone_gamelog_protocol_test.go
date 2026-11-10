package audit_test

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// TestFase09_Control4_Zone_C2S_OpCodes verifica la paridad 100% de los 245 OpCodes C2S
// oficiales de Zone.exe (protocolstore(ShinePlayer) en 0x004D3510), validando:
// 1. La fórmula matemática oficial: protocol = (dept << 10) | cmd.
// 2. Que el 100% de los 245 OpCodes estén registrados en pkg/protocol/opcodes.go.
func TestFase09_Control4_Zone_C2S_OpCodes(t *testing.T) {
	rootPath := findRepoRoot(t)
	planPath := filepath.Join(rootPath, "docs", "PLAN_ACCION_ZONESERVICE_GAMELOG.md")

	planContent, err := os.ReadFile(planPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", planPath, err)
	}

	// Extraer tabla de la Sección 4.1
	sec41Re := regexp.MustCompile(`(?s)### 4\.1 Las 245 Entradas Oficiales.*?(\|[^\n]+\|\n\|(?:\s*:\s*---\s*:?\s*\|)+\n(?:\|[^\n]+\|\n)+)`)
	m := sec41Re.FindStringSubmatch(string(planContent))
	if len(m) < 2 {
		t.Fatalf("Failed to find Section 4.1 table in %s", planPath)
	}

	tableLines := strings.Split(strings.TrimSpace(m[1]), "\n")[2:]
	if len(tableLines) != 245 {
		t.Fatalf("Expected exactly 245 C2S OpCodes in Section 4.1, found %d", len(tableLines))
	}

	// Cargar todos los OpCodes definidos en pkg/protocol/opcodes.go
	opcodesPath := filepath.Join(rootPath, "pkg", "protocol", "opcodes.go")
	opcodesContent, err := os.ReadFile(opcodesPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", opcodesPath, err)
	}

	definedValues := make(map[uint16]bool)
	constRe := regexp.MustCompile(`(\bPROTO_NC_[A-Za-z0-9_]+)\s*(?:uint16)?\s*=\s*(?:\(\s*(\d+)\s*<<\s*10\s*\)\s*\|\s*(\d+)|(0x[0-9A-Fa-f]{4}))`)
	for _, match := range constRe.FindAllStringSubmatch(string(opcodesContent), -1) {
		if match[2] != "" {
			dept, _ := strconv.Atoi(match[2])
			cmd, _ := strconv.Atoi(match[3])
			definedValues[uint16((dept<<10)|cmd)] = true
		} else {
			val, _ := strconv.ParseUint(strings.TrimPrefix(match[4], "0x"), 16, 16)
			definedValues[uint16(val)] = true
		}
	}

	verifiedCount := 0
	for _, line := range tableLines {
		rawCols := strings.Split(line, "|")
		if len(rawCols) < 8 {
			continue
		}
		cols := make([]string, len(rawCols)-2)
		for i := 1; i < len(rawCols)-1; i++ {
			cols[i-1] = strings.Trim(strings.TrimSpace(rawCols[i]), "` ")
		}

		idx := cols[0]
		hexStr := cols[1]
		decVal, err := strconv.Atoi(cols[2])
		if err != nil {
			t.Fatalf("Row %s: invalid decimal value %s: %v", idx, cols[2], err)
		}
		deptMatch := regexp.MustCompile(`\d+`).FindString(cols[3])
		dept, err := strconv.Atoi(deptMatch)
		if err != nil {
			t.Fatalf("Row %s: invalid department %s: %v", idx, cols[3], err)
		}
		cmd, err := strconv.Atoi(cols[4])
		if err != nil {
			t.Fatalf("Row %s: invalid command %s: %v", idx, cols[4], err)
		}

		expectedProtocol := uint16((dept << 10) | (cmd & 0x3FF))
		parsedHex, err := strconv.ParseUint(strings.TrimPrefix(hexStr, "0x"), 16, 16)
		if err != nil {
			t.Fatalf("Row %s: invalid hex %s: %v", idx, hexStr, err)
		}

		if uint16(parsedHex) != expectedProtocol {
			t.Errorf("Row %s: hex %s does not match formula (dept %d, cmd %d): expected 0x%04X",
				idx, hexStr, dept, cmd, expectedProtocol)
		}
		if uint16(decVal) != expectedProtocol {
			t.Errorf("Row %s: dec %d does not match formula (dept %d, cmd %d): expected %d",
				idx, decVal, dept, cmd, expectedProtocol)
		}

		if !definedValues[expectedProtocol] {
			t.Errorf("Row %s: OpCode 0x%04X (%s) not found in pkg/protocol/opcodes.go",
				idx, expectedProtocol, cols[5])
		} else {
			verifiedCount++
		}
	}

	if verifiedCount != 245 {
		t.Fatalf("Control 4 (C2S OpCodes): Expected 245/245 verified, got %d", verifiedCount)
	}

	t.Logf("Control 4 (C2S OpCodes): 100%% Parity verified (%d/245 C2S OpCodes match official PDB)", verifiedCount)
}

// TestFase09_Control4_GameLog_StructSizes verifica que el 100% de las 44 estructuras
// oficiales de docs/Servicios/GameLog.md estén registradas en pkg/protocol y produzcan
// exactamente el sizeof en bytes del PDB TPI.
func TestFase09_Control4_GameLog_StructSizes(t *testing.T) {
	rootPath := findRepoRoot(t)
	docPath := filepath.Join(rootPath, "docs", "Servicios", "GameLog.md")

	content, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", docPath, err)
	}

	structSectionRe := regexp.MustCompile(`(?s)###\s+\x60([A-Za-z0-9_<>]+)\x60\s*\n.*?\*\*Tamaño Total en Memoria:\*\*\s*\x60(\d+)\s*bytes\x60`)
	matches := structSectionRe.FindAllStringSubmatch(string(content), -1)

	if len(matches) != 44 {
		t.Fatalf("Expected exactly 44 structs in GameLog.md, found %d", len(matches))
	}

	verifiedCount := 0
	for _, m := range matches {
		structName := m[1]
		expectedSize, _ := strconv.Atoi(m[2])

		instance, exists := protocol.ZoneGameLogPDBStructs[structName]
		if !exists {
			t.Errorf("GameLog struct %s (expected size %d) is missing in protocol.ZoneGameLogPDBStructs",
				structName, expectedSize)
			continue
		}

		actualSize := binary.Size(instance)
		if actualSize == 0 && expectedSize == 1 {
			actualSize = 1
		}

		if actualSize != expectedSize {
			t.Errorf("GameLog struct %s sizeof mismatch: expected %d bytes, got %d bytes",
				structName, expectedSize, actualSize)
			continue
		}

		verifiedCount++
	}

	if verifiedCount != 44 {
		t.Fatalf("Control 4 (GameLog StructSizes): Expected 44/44 verified, got %d", verifiedCount)
	}

	t.Logf("Control 4 (GameLog StructSizes): 100%% Parity verified (%d/44 structs match sizeof)", verifiedCount)
}

// TestFase09_Control4_Zone_StructSizes verifica que el 100% de las 1.293 estructuras PROTO_NC_*
// y las 9 plantillas ShineSendType<PROTO_NC_*> de docs/Servicios/Zone.md (excluyendo únicamente
// OpTool y los 40 tipos internos STL) existan en pkg/protocol con sizeof idéntico al byte.
func TestFase09_Control4_Zone_StructSizes(t *testing.T) {
	rootPath := findRepoRoot(t)
	docPath := filepath.Join(rootPath, "docs", "Servicios", "Zone.md")

	content, err := os.ReadFile(docPath)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", docPath, err)
	}

	sections := strings.Split(string(content), "\n### `")[1:]
	expectedSizes := make(map[string]int)

	sizeRe := regexp.MustCompile(`Tamaño Total en Memoria:\*\*\s*\x60(\d+)\s*bytes\x60`)

	for _, s := range sections {
		nameEnd := strings.Index(s, "`")
		if nameEnd == -1 {
			continue
		}
		name := s[:nameEnd]

		// Exclusiones oficiales según Sección 7 de PLAN_ACCION_ZONESERVICE_GAMELOG.md:
		// 1. OpTool (PROTO_NC_OPTOOL_CONNECT_BRIF_ACK)
		// 2. 40 artefactos internos C++ STL (std::vector, std::allocator, reverse_iterator)
		if strings.Contains(name, "OPTOOL") || strings.HasPrefix(name, "std::") ||
			strings.Contains(name, "reverse_iterator") || strings.Contains(name, "allocator") {
			continue
		}

		sizeMatch := sizeRe.FindStringSubmatch(s)
		if len(sizeMatch) < 2 {
			continue
		}
		size, err := strconv.Atoi(sizeMatch[1])
		if err != nil {
			continue
		}

		expectedSizes[name] = size
	}

	// 1.293 structs PROTO_NC_* + 9 ShineSendType = 1.302 structs
	if len(expectedSizes) != 1302 {
		t.Fatalf("Expected exactly 1.302 non-OpTool, non-STL structs in Zone.md, found %d", len(expectedSizes))
	}

	verifiedCount := 0
	sizeMismatches := 0
	missingInRegistry := 0

	for structName, expectedSize := range expectedSizes {
		instance, exists := protocol.ZoneGameLogPDBStructs[structName]
		if !exists {
			t.Errorf("Zone struct %s (expected size %d) is missing in protocol.ZoneGameLogPDBStructs",
				structName, expectedSize)
			missingInRegistry++
			continue
		}

		actualSize := binary.Size(instance)
		// En ISO C++, una struct vacía mide 1 byte (sizeof == 1), mientras que en Go struct{} tiene 0 bytes.
		if actualSize == 0 && expectedSize == 1 {
			actualSize = 1
		}

		if actualSize != expectedSize {
			t.Errorf("Zone struct %s sizeof mismatch: expected %d bytes, got %d bytes",
				structName, expectedSize, actualSize)
			sizeMismatches++
			continue
		}

		verifiedCount++
	}

	if missingInRegistry > 0 || sizeMismatches > 0 {
		t.Fatalf("Control 4 (Zone StructSizes) failed: %d missing in registry, %d size mismatches",
			missingInRegistry, sizeMismatches)
	}

	t.Logf("Control 4 (Zone StructSizes): 100%% Parity verified (%d/1.302 structs match sizeof)", verifiedCount)
}

// TestFase09_Control4_ForbiddenAliases_Scan escanea todo pkg/ y falla si reaparece
// cualquiera de los alias prohibidos de emuladores privados.
func TestFase09_Control4_ForbiddenAliases_Scan(t *testing.T) {
	bannedAliases := []string{
		"PROTO_NC_ACT_REST_START_REQ",
		"PROTO_NC_ACT_REST_START_ACK",
		"PROTO_NC_ACT_REST_STOP_REQ",
		"PROTO_NC_ACT_REST_STOP_ACK",
		"PROTO_NC_CHAR_INGAME_REQ",
		"PROTO_NC_CHAR_INGAME_ACK",
		"PROTO_NC_MISC_WORLD_CLIENT_KEY_REQ",
		"PROTO_NC_CHAR_CLIENT_CHARLIST_ACK",
		"PROTO_NC_CHAR_CREATE_REQ",
		"PROTO_NC_CHAR_CREATE_ACK",
		"PROTO_NC_CHAR_CREATE_FAIL_ACK",
		"PROTO_NC_CHAR_DELETE_REQ",
		"PROTO_NC_CHAR_DELETE_ACK",
		"PROTO_NC_CHAR_DELETE_FAIL_ACK",
		"PROTO_NC_CHAR_GET_QUICKBAR_REQ",
		"PROTO_NC_CHAR_GET_QUICKBAR_STATE_REQ",
		"PROTO_NC_CHAR_GET_GAME_SETTINGS_REQ",
		"PROTO_NC_CHAR_GET_CLIENT_SETTINGS_REQ",
		"PROTO_NC_CHAR_GET_SHORTCUTS_REQ",
		"PROTO_NC_CHAR_SAVE_QUICKBAR_REQ",
		"PROTO_NC_CHAR_SAVE_QUICKBAR_STATE_REQ",
		"PROTO_NC_CHAR_SAVE_GAME_SETTINGS_REQ",
		"PROTO_NC_CHAR_SAVE_CLIENT_SETTINGS_REQ",
		"PROTO_NC_CHAR_SAVE_SHORTCUTS_REQ",
		"PROTO_NC_USER_BACKTOWORLDLIST_REQ",
		"PROTO_NC_USER_BACKTOWORLDLIST_ACK",
	}

	rootPath := findRepoRoot(t)
	pkgDir := filepath.Join(rootPath, "pkg")

	var violations []string

	err := filepath.Walk(pkgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		// Excluir archivos de test de auditoría que contienen la lista de alias prohibidos
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content := string(data)
		for _, alias := range bannedAliases {
			re := regexp.MustCompile(`\b` + regexp.QuoteMeta(alias) + `\b`)
			if re.MatchString(content) {
				relPath, _ := filepath.Rel(rootPath, path)
				violations = append(violations, relPath+": contains banned alias "+alias)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to scan pkg/ directory: %v", err)
	}

	if len(violations) > 0 {
		t.Fatalf("CRITICAL: Found banned private emulator aliases in codebase:\n  %s",
			strings.Join(violations, "\n  "))
	}

	t.Logf("Control 4 (ForbiddenAliases): Zero banned private emulator aliases found in codebase (%d scanned)",
		len(bannedAliases))
}

// TestFase09_Control4_NoPdbZoneGameLogSubpackage verifica que NO existan
// archivos de Zone ni GameLog en el subpaquete pkg/protocol/pdb.
func TestFase09_Control4_NoPdbZoneGameLogSubpackage(t *testing.T) {
	rootPath := findRepoRoot(t)
	pdbDir := filepath.Join(rootPath, "pkg", "protocol", "pdb")

	if _, err := os.Stat(pdbDir); os.IsNotExist(err) {
		t.Log("Control 4 (NoPdb): pkg/protocol/pdb does not exist")
		return
	}

	zoneForbiddenFiles := []string{
		"zone", "bat", "briefinfo", "casino", "act",
		"mover", "gamelog", "shine", "subsystem",
	}

	entries, err := os.ReadDir(pdbDir)
	if err != nil {
		t.Fatalf("Failed to read pkg/protocol/pdb: %v", err)
	}

	for _, entry := range entries {
		name := strings.ToLower(entry.Name())
		for _, forbidden := range zoneForbiddenFiles {
			if strings.Contains(name, forbidden) {
				t.Fatalf("CRITICAL: Forbidden Zone/GameLog file found in pkg/protocol/pdb: %s (must be in pkg/protocol)",
					entry.Name())
			}
		}
	}

	t.Log("Control 4 (NoPdb): No Zone or GameLog protocol files found in pkg/protocol/pdb (all in pkg/protocol)")
}
