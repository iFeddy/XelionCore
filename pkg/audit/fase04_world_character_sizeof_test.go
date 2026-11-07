package audit_test

import (
	"encoding/binary"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// TestFase04_Control4_WorldCharacter_OpCodes_Completeness verifica:
// 1. Que el 100% de los OpCodes no-OpTool de WorldManager.md y Character.md cumplan la fórmula (dept << 10) | cmd.
// 2. Que el 100% de dichos OpCodes estén definidos en pkg/protocol/opcodes.go con su valor oficial exacto.
func TestFase04_Control4_WorldCharacter_OpCodes_Completeness(t *testing.T) {
	rootPath := findRepoRoot(t)

	docs := []string{
		filepath.Join(rootPath, "docs", "Servicios", "WorldManager.md"),
		filepath.Join(rootPath, "docs", "Servicios", "Character.md"),
	}

	// Extraer todas las constantes de OpCodes definidas en pkg/protocol/opcodes.go
	opcodesGoPath := filepath.Join(rootPath, "pkg", "protocol", "opcodes.go")
	opcodesContent, err := os.ReadFile(opcodesGoPath)
	if err != nil {
		t.Fatalf("Failed to read opcodes.go: %v", err)
	}

	// Regex para constantes en Go: PROTO_NC_* uint16 = (dept << 10) | cmd // 0xHEX
	constFormulaRe := regexp.MustCompile(`(\bPROTO_NC_[A-Z0-9_]+)\s*(?:uint16)?\s*=\s*\(\s*(\d+)\s*<<\s*10\s*\)\s*\|\s*(\d+)`)
	definedOpCodes := make(map[string]uint16)
	for _, m := range constFormulaRe.FindAllStringSubmatch(string(opcodesContent), -1) {
		name := m[1]
		dept, _ := strconv.Atoi(m[2])
		cmd, _ := strconv.Atoi(m[3])
		definedOpCodes[name] = uint16((dept << 10) | (cmd & 0x3FF))
	}

	// También parsear si alguna constante fue definida directamente en hexadecimal
	constHexRe := regexp.MustCompile(`(\bPROTO_NC_[A-Z0-9_]+)\s*(?:uint16)?\s*=\s*(0x[0-9A-Fa-f]{4})`)
	for _, m := range constHexRe.FindAllStringSubmatch(string(opcodesContent), -1) {
		name := m[1]
		val, _ := strconv.ParseUint(strings.TrimPrefix(m[2], "0x"), 16, 16)
		definedOpCodes[name] = uint16(val)
	}

	tableRowRe := regexp.MustCompile(`(?m)\|\s*\x60(0x[0-9A-Fa-f]{4})\x60\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*\x60([^\x60]+)\x60\s*\|\s*\x60([^\x60]+)\x60\s*\|`)

	totalVerified := 0
	docOpcodes := make(map[string]map[uint16]bool)

	for _, docPath := range docs {
		content, err := os.ReadFile(docPath)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", docPath, err)
		}

		matches := tableRowRe.FindAllStringSubmatch(string(content), -1)
		for _, m := range matches {
			hexStr := m[1]
			dept, _ := strconv.Atoi(m[3])
			cmd, _ := strconv.Atoi(m[4])
			fnName := strings.TrimSpace(m[5])
			protoName := strings.TrimSpace(m[6])

			// Excluir OpTool según Sección 7 de PLAN_ACCION_WORLDSERVICE_CHARACTER.md
			if dept == 10 || strings.Contains(protoName, "OPTOOL") || strings.Contains(fnName, "OPTOOL") {
				continue
			}

			// Validar fórmula matemática oficial: protocol = (dept << 10) | (cmd & 0x3FF)
			expectedProtocol := uint16((dept << 10) | (cmd & 0x3FF))

			parsedHex, err := strconv.ParseUint(strings.TrimPrefix(hexStr, "0x"), 16, 16)
			if err != nil {
				t.Errorf("%s: invalid hex %s: %v", docPath, hexStr, err)
				continue
			}

			if uint16(parsedHex) != expectedProtocol {
				t.Errorf("%s: OpCode %s (%s) formula mismatch: expected 0x%04X, got 0x%04X",
					docPath, protoName, hexStr, expectedProtocol, parsedHex)
			}

			if docOpcodes[protoName] == nil {
				docOpcodes[protoName] = make(map[uint16]bool)
			}
			docOpcodes[protoName][expectedProtocol] = true
			totalVerified++
		}
	}

	if len(docOpcodes) < 700 {
		t.Fatalf("Insufficient unique opcodes parsed from WorldManager and Character: got %d, expected >= 700",
			len(docOpcodes))
	}

	// Verificar que cada OpCode no-OpTool esté presente en pkg/protocol/opcodes.go con un valor oficial válido
	missingCount := 0
	valueMismatches := 0
	for protoName, validVals := range docOpcodes {
		val, exists := definedOpCodes[protoName]
		if !exists {
			t.Errorf("OpCode %s is missing in pkg/protocol/opcodes.go", protoName)
			missingCount++
			continue
		}
		if !validVals[val] {
			t.Errorf("OpCode %s value mismatch in opcodes.go: got 0x%04X, expected one of %v",
				protoName, val, validVals)
			valueMismatches++
		}
	}

	if missingCount > 0 || valueMismatches > 0 {
		t.Fatalf("OpCodes completeness check failed: %d missing, %d value mismatches",
			missingCount, valueMismatches)
	}

	t.Logf("Control 4 (OpCodes): 100%% Parity verified (%d OpCodes non-OpTool, 0 missing, 0 value mismatches)",
		len(docOpcodes))
}

// TestFase04_Control4_WorldCharacter_StructSizes verifica el sizeof exacto al byte
// de todas las estructuras oficiales C++ no-OpTool de WorldManager.md y Character.md
// registradas en pkg/protocol/world_character_registry.go.
func TestFase04_Control4_WorldCharacter_StructSizes(t *testing.T) {
	rootPath := findRepoRoot(t)

	docs := []string{
		filepath.Join(rootPath, "docs", "Servicios", "WorldManager.md"),
		filepath.Join(rootPath, "docs", "Servicios", "Character.md"),
	}

	expectedSizes := make(map[string]int)
	structSectionRe := regexp.MustCompile(`(?m)###\s+\x60?(PROTO_[A-Za-z0-9_]+)\x60?\s*\n\s*-\s*\*\*Tamaño Total en Memoria:\*\*\s*\x60?(\d+)\s*bytes\x60?`)

	for _, docPath := range docs {
		content, err := os.ReadFile(docPath)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", docPath, err)
		}

		matches := structSectionRe.FindAllStringSubmatch(string(content), -1)
		for _, m := range matches {
			name := m[1]
			size, _ := strconv.Atoi(m[2])

			// Excluir OpTool según Sección 7
			if strings.Contains(name, "OPTOOL") {
				continue
			}

			expectedSizes[name] = size
		}
	}

	if len(expectedSizes) < 1400 {
		t.Fatalf("Insufficient struct sizes extracted from docs: got %d, expected >= 1400", len(expectedSizes))
	}

	verifiedCount := 0
	sizeMismatches := 0
	missingInRegistry := 0

	for structName, expectedSize := range expectedSizes {
		instance, exists := protocol.WorldCharacterPDBStructs[structName]
		if !exists {
			t.Errorf("Struct %s (expected size %d) is missing in protocol.WorldCharacterPDBStructs",
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
			t.Errorf("Struct %s sizeof mismatch: expected %d bytes, got %d bytes",
				structName, expectedSize, actualSize)
			sizeMismatches++
			continue
		}

		verifiedCount++
	}

	if missingInRegistry > 0 || sizeMismatches > 0 {
		t.Fatalf("Control 4 (StructSizes) failed: %d missing in registry, %d size mismatches",
			missingInRegistry, sizeMismatches)
	}

	t.Logf("Control 4 (StructSizes): 100%% Parity verified (%d structs checked, 0 mismatches)", verifiedCount)
}

// TestFase04_Control4_WorldCharacter_ForbiddenAliases_Scan escanea todo pkg/
// y falla si reaparece cualquiera de los 19 alias prohibidos de emuladores privados.
func TestFase04_Control4_WorldCharacter_ForbiddenAliases_Scan(t *testing.T) {
	bannedAliases := []string{
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
		"PROTO_NC_CHAR_INGAME_REQ",
		"PROTO_NC_CHAR_INGAME_ACK",
		"PROTO_NC_ACT_REST_START_REQ",
		"PROTO_NC_ACT_REST_START_ACK",
		"PROTO_NC_ACT_REST_STOP_REQ",
		"PROTO_NC_ACT_REST_STOP_ACK",
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
		// Solo escanear archivos de código .go (excluyendo este archivo de test)
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "fase04_world_character_sizeof_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			// Si falla parseo de comentarios, leer directamente texto
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			str := string(content)
			for _, banned := range bannedAliases {
				if strings.Contains(str, banned) {
					violations = append(violations, fmt.Sprintf("%s in %s", banned, path))
				}
			}
			return nil
		}

		// Inspeccionar identificadores en AST
		for _, ident := range node.Unresolved {
			for _, banned := range bannedAliases {
				if ident.Name == banned {
					relPath, _ := filepath.Rel(rootPath, path)
					violations = append(violations, fmt.Sprintf("%s in %s (line %d)",
						banned, relPath, fset.Position(ident.Pos()).Line))
				}
			}
		}

		// También buscar en contenido de tokens y declaraciones
		content, _ := os.ReadFile(path)
		for _, banned := range bannedAliases {
			re := regexp.MustCompile(`\b` + banned + `\b`)
			if re.Match(content) {
				relPath, _ := filepath.Rel(rootPath, path)
				violations = append(violations, fmt.Sprintf("%s found in %s", banned, relPath))
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Walk error in pkg: %v", err)
	}

	if len(violations) > 0 {
		t.Fatalf("CRITICAL: Found banned private emulator aliases in codebase:\n  %s",
			strings.Join(violations, "\n  "))
	}

	t.Logf("Control 4 (Aliases): Zero banned private emulator aliases found in codebase (%d scanned)",
		len(bannedAliases))
}

// TestFase04_Control4_WorldCharacter_NoPdbWorldSubpackage verifica que NO existan
// archivos de WorldManager ni Character en el subpaquete pkg/protocol/pdb.
func TestFase04_Control4_WorldCharacter_NoPdbWorldSubpackage(t *testing.T) {
	rootPath := findRepoRoot(t)
	pdbDir := filepath.Join(rootPath, "pkg", "protocol", "pdb")

	if _, err := os.Stat(pdbDir); os.IsNotExist(err) {
		// Si no existe, cumple perfectamente la regla
		t.Log("Control 4 (NoPdb): pkg/protocol/pdb does not exist")
		return
	}

	// Si existe, verificar que NO contenga archivos de World o Character
	worldCharForbiddenFiles := []string{
		"avatar", "world", "guild", "friend", "party", "raid",
		"kq", "itemdb", "option", "prison", "rank", "indun",
		"holy_promise", "minihouse", "pet", "char",
	}

	entries, err := os.ReadDir(pdbDir)
	if err != nil {
		t.Fatalf("Failed to read pkg/protocol/pdb: %v", err)
	}

	for _, entry := range entries {
		name := strings.ToLower(entry.Name())
		for _, forbidden := range worldCharForbiddenFiles {
			if strings.Contains(name, forbidden) {
				t.Fatalf("CRITICAL: Forbidden World/Character file found in pkg/protocol/pdb: %s (must be in pkg/protocol)",
					entry.Name())
			}
		}
	}

	t.Log("Control 4 (NoPdb): No World or Character protocol files found in pkg/protocol/pdb (all in pkg/protocol)")
}
