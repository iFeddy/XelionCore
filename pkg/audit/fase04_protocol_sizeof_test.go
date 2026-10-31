package audit_test

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/protocol/pdb"
)

// TestFase04_Control4_OpCodes verifica la paridad bit a bit de todos los OpCodes
// de Login.md, Account.md y AccountLog.md contra la fórmula (dept << 10) | cmd,
// y asegura que PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK sea 0x0C02.
func TestFase04_Control4_OpCodes(t *testing.T) {
	// Verificación explícita de Bug Crítico #1
	if protocol.PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK != 0x0C02 {
		t.Fatalf("CRITICAL: PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK must be 0x0C02, got 0x%04X",
			protocol.PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK)
	}

	docs := []string{
		"../../docs/Servicios/Login.md",
		"../../docs/Servicios/Account.md",
		"../../docs/Servicios/AccountLog.md",
	}

	// Regex para tablas de OpCodes: | `0x0C02` | 3074 | 3 | 2 | ...
	tableRowRe := regexp.MustCompile(`\|\s*` + "`" + `(0x[0-9A-Fa-f]{4})` + "`" + `\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*([^|]+)\|`)

	totalVerified := 0

	for _, docPath := range docs {
		content, err := os.ReadFile(docPath)
		if err != nil {
			t.Fatalf("Failed to read %s: %v", docPath, err)
		}

		matches := tableRowRe.FindAllStringSubmatch(string(content), -1)
		for _, m := range matches {
			hexStr := m[1]
			decVal, _ := strconv.Atoi(m[2])
			dept, _ := strconv.Atoi(m[3])
			cmd, _ := strconv.Atoi(m[4])
			handler := strings.TrimSpace(m[5])
			handler = strings.ReplaceAll(handler, "`", "")

			// Excluir OpTool según Sección 7
			if strings.Contains(handler, "OPTOOL") {
				continue
			}

			// Validar fórmula matemática oficial: protocol = (dept << 10) | cmd
			expectedProtocol := uint16((dept << 10) | cmd)

			parsedHex, err := strconv.ParseUint(strings.TrimPrefix(hexStr, "0x"), 16, 16)
			if err != nil {
				t.Errorf("%s: invalid hex %s: %v", docPath, hexStr, err)
				continue
			}

			if uint16(parsedHex) != expectedProtocol {
				t.Errorf("%s: OpCode %s formula mismatch: expected 0x%04X, got 0x%04X",
					docPath, hexStr, expectedProtocol, parsedHex)
			}

			if uint16(decVal) != expectedProtocol {
				t.Errorf("%s: decimal value mismatch for %s: %d vs formula %d",
					docPath, hexStr, decVal, expectedProtocol)
			}

			totalVerified++
		}
	}

	if totalVerified < 60 {
		t.Fatalf("Insufficient opcodes verified from docs: got %d, expected >= 60", totalVerified)
	}
	t.Logf("Control 4: Verified %d OpCodes from PDB markdown documentation against (dept << 10) | cmd", totalVerified)
}

// TestFase04_Control4_StructSizes verifica el sizeof exacto en bytes de cada struct
// registrada en pkg/protocol/pdb contra los valores documentados en el PDB.
func TestFase04_Control4_StructSizes(t *testing.T) {
	expectedSizes := map[string]int{
		// Critical User Handshakes & Requests
		"PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ":      64,
		"PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK": 1,
		"PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK": 1,
		"PROTO_NC_USER_XTRAP_REQ":                     1,
		"PROTO_NC_USER_XTRAP_ACK":                     1,
		"PROTO_NC_USER_LOGIN_REQ":                     272,
		"PROTO_NC_USER_PASSWORD_CHECK_REQ":            274,
		"PROTO_NC_USER_PASSWORD_CHECK_ACK":            287,
		"PROTO_NC_USER_PASSWORD_CHECK_ACK_SEND":       290,
		"PROTO_NC_USER_LOGINFAIL_ACK":                 2,
		"PROTO_NC_USER_LOGIN_ACK":                     1,
		"PROTO_NC_USER_LOGIN_ACK_WorldInfo":           18,
		"PROTO_NC_USER_WORLDSELECT_REQ":               1,
		"PROTO_NC_USER_WORLDSELECT_ACK":               83,
		"PROTO_NC_USER_WILLLOGIN_REQ":                 347,
		"PROTO_NC_USER_WILLLOGIN_ACK":                 348,
		"PROTO_NC_USER_WORLD_LOGIN_REQ":               83,
		"PROTO_NC_USER_WORLD_LOGIN_ACK":               26,
		"PROTO_NC_USER_LOGIN_DB":                      9,
		"PROTO_NC_USER_LOGOUT_DB":                     9,
		"PROTO_NC_USER_KICKOFFFROMWORLD_CMD":          4,
		"PROTO_NC_USER_CONNECTCUT_CMD":                2,
		"PROTO_NC_USER_WORLD_STATUS_REQ":              1,
		"PROTO_NC_USER_WORLD_STATUS_ACK":              1,
		"PROTO_NC_USER_LOGIN_NETMARBLE_REQ":           2,
		"PROTO_NC_USER_LOGIN_NETMARBLE_DB_REQ":        278,
		"PROTO_NC_USER_LOGIN_OUTSPARK_REQ":            85,
		"PROTO_NC_USER_LOGIN_WITH_OTP_REQ":            32,
		"PROTO_NC_USER_CREATE_OTP_REQ":                266,
		"PROTO_NC_USER_CREATE_OTP_ACK":                36,
		"PROTO_NC_USER_WILL_WORLD_SELECT_REQ":         6,
		"PROTO_NC_USER_WILL_WORLD_SELECT_ACK":         34,
		"PROTO_NC_USER_TEENAGER_REQ":                  258,
		"PROTO_NC_USER_TEENAGER_ACK":                  313,
		"PROTO_NC_USER_TEENAGER_CMD":                  55,
		"PROTO_NC_USER_TEENAGER_SET_CMD":              54,
		"PROTO_NC_USER_IS_IP_BLOCK_REQ":               298,
		"PROTO_NC_USER_IS_IP_BLOCK_ACK":               275,
		"PROTO_NC_USER_GER_PASSWORD_CHECK_REQ":        274,
		"PROTO_NC_USER_GER_PASSWORD_CHECK_ACK":        287,
		"PROTO_NC_USER_GER_PASSWORD_CHECK_ACK_SEND":   290,
		"PROTO_NC_USER_GER_IS_IP_BLOCK_REQ":           313,
		"PROTO_NC_USER_GER_IS_IP_BLOCK_ACK":           290,
		"PROTO_NC_USER_TW_PASSWORD_CHECK_REQ":         102,
		"PROTO_NC_USER_TW_PASSWORD_CHECK_ACK":         8,
		"PROTO_NC_USER_TW_PASSWORD_CHECK_ACK_SEND":    11,
		"PROTO_NC_USER_TW_IS_IP_BLOCK_REQ":            34,
		"PROTO_NC_USER_TW_IS_IP_BLOCK_ACK":            11,
		"PROTO_NC_USER_JP_PASSWORD_CHECK_REQ":         336,
		"PROTO_NC_USER_JP_PASSWORD_CHECK_ACK":         306,
		"PROTO_NC_USER_JP_IS_IP_BLOCK_REQ":            332,
		"PROTO_NC_USER_JP_IS_IP_BLOCK_ACK":            309,
		"PROTO_NC_USER_CH_PASSWORD_CHECK_REQ":         274,
		"PROTO_NC_USER_CH_PASSWORD_CHECK_ACK":         287,
		"PROTO_NC_USER_CH_PASSWORD_CHECK_ACK_SEND":    290,
		"PROTO_NC_USER_US_LOGIN_REQ":                  316,
		"PROTO_NC_USER_US_PASSWORD_CHECK_REQ":         338,
		"PROTO_NC_USER_US_PASSWORD_CHECK_ACK":         353,
		"PROTO_NC_USER_US_IS_IP_BLOCK_REQ":            379,
		"PROTO_NC_USER_US_IS_IP_BLOCK_ACK":            356,
		"PROTO_NC_USER_TUTORIAL_CAN_SKIP_CMD":         1,
		"PROTO_NC_USER_RETURN_CHECK_REQ":              6,
		"PROTO_NC_USER_RETURN_CHECK_ACK":              11,
		"PROTO_NC_USER_LOGOUT_LAST_DAY_REQ":           6,
		"PROTO_NC_USER_LOGOUT_LAST_DAY_ACK":           10,
		"PROTO_NC_USER_SET_RETURN_CMD":                4,
		"PROTO_NC_USP_USER_CHARACTER_DELETE":          9,
		"PROTO_NC_USP_USER_CHARACTER_INSERT":          50,

		// LOG
		"PROTO_NC_LOG_USER_LOGIN":     9,
		"PROTO_NC_LOG_USER_LOGINFAIL": 278,
		"PROTO_NC_LOG_USER_LOGOUT":    9,

		// SYSLOG
		"PROTO_NC_SYSLOG_ACCOUNT_LOGIN_SUCCESS": 40,
		"PROTO_NC_SYSLOG_ACCOUNT_LOGIN_FAILURE": 296,
		"PROTO_NC_SYSLOG_ACCOUNT_LOGOUT":        8,
		"PROTO_NC_SYSLOG_SERVER_CCU":            5,

		// MISC
		"PROTO_NC_MISC_SEED_ACK":                 2,
		"PROTO_NC_MISC_S2SCONNECTION_RDY_SEND":   3,
		"PROTO_NC_MISC_S2SCONNECTION_REQ":        7,
		"PROTO_NC_MISC_S2SCONNECTION_ACK":        3,
		"PROTO_NC_MISC_HEARTBEAT_ACK_SEND":       3,
		"PROTO_NC_MISC_USER_COUNT_CMD":           4,
		"PROTO_NC_MISC_CONNECTER_REQ":            1,
		"PROTO_NC_MISC_CONNECTER_ACK":            131,
		"PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD": 2,
		"PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD": 2,
		"PROTO_NC_MISC_WEB_DB_KEY_REQ":           6,
		"PROTO_NC_MISC_WEB_DB_KEY_ACK":           40,
		"PROTO_NC_MISC_CHAR_LOGOFF_STATISTICS":   21,
		"PROTO_NC_MISC_EVENT_HIT3_ADD_CASH":     28,
		"PROTO_NC_MISC_CONNECTFROMWHERE_REQ":     10,
		"PROTO_NC_MISC_CONNECTFROMWHERE_ACK":     7,
		"PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ":  12,
		"PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK":  9,
		"PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC": 13,
		"PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM": 15,
		"PROTO_NC_MISC_EVENT_L20_DB_REQ":         272,
		"PROTO_NC_MISC_EVENT_L20_DB_ACK":         13,
		"PROTO_NC_MISC_HACK_SCAN_STORE_CMD":      2,
		"PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD":   51,
		"PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD": 6,
		"PROTO_NC_MISC_ITEMSHOP_URL_DB_REQ":      8,
		"PROTO_NC_MISC_ITEMSHOP_URL_DB_ACK":      8,

		// REPORT
		"PROTO_NC_DATA_REPORT_ADD_REQ":                           2369,
		"PROTO_NC_DATA_REPORT_ADD_ACK":                           28,
		"PROTO_NC_DATA_REPORT_CANCEL_REQ":                        6,
		"PROTO_NC_DATA_REPORT_CANCEL_ACK":                        8,
		"PROTO_NC_DATA_REPORT_GET_REQ":                           7,
		"PROTO_NC_DATA_REPORT_GET_ACK":                           137,
		"PROTO_NC_DATA_SEND1_PROTO_NC_DATA_REPORT_ADD_ACK":       31,
		"PROTO_NC_DATA_SEND1_PROTO_NC_DATA_REPORT_CANCEL_ACK":    11,
		"PROTO_NC_DATA_SEND1_PROTO_NC_DATA_REPORT_GET_ACK":       140,
		"PROTO_NC_DATA_SEND1_PROTO_NC_MISC_EVENT_L20_DB_ACK":     16,
		"PROTO_NC_DATA_SEND1_PROTO_NC_USER_GER_IS_IP_BLOCK_ACK":  293,
		"PROTO_NC_DATA_SEND1_PROTO_NC_USER_IS_IP_BLOCK_ACK":      278,
		"PROTO_NC_DATA_SEND1_PROTO_NC_USER_TW_IS_IP_BLOCK_ACK":   14,
		"PROTO_NC_DATA_SEND1_PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK": 12,
		"PROTO_NC_DATA_SEND1_PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM": 18,

		// SHARED
		"PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD": 3,
		"PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD":      5,
		"PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD":       6,
		"PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD":        9,
		"PROTO_NC_BRIEFINFO_UNEQUIP_CMD":             3,
		"PROTO_NC_CHARSAVE_ITEMACTIONCOOLTIME_ACK":   2,
		"PROTO_NC_CHAR_ABSTATE_CMD":                  6,
		"PROTO_NC_CHAR_ARENA_CMD":                    1,
		"PROTO_NC_CHAR_BASE_CMD":                     105,
		"PROTO_NC_CHAR_CHARGEDBUFF_CMD":              2,
		"PROTO_NC_CHAR_CHESTINFO_CMD":                2,
		"PROTO_NC_CHAR_EMBLEM_CMD":                   1,
		"PROTO_NC_CHAR_FRIEND_CMD":                   2,
		"PROTO_NC_CHAR_FRIEND_CMD_INFO":              25,
		"PROTO_NC_CHAR_GUILD_CMD":                    4,
		"PROTO_NC_CHAR_HOUSE_CMD":                    1,
		"PROTO_NC_CHAR_MASPUP_CMD":                   1,
		"PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_LIST_CMD":   1,
		"PROTO_NC_CHAR_OPTION_GAME":                  64,
		"PROTO_NC_CHAR_OPTION_KEYMAPPING":            308,
		"PROTO_NC_CHAR_OPTION_SHORTCUTDATA":          1024,
		"PROTO_NC_CHAR_OPTION_SHORTCUTSIZE":          24,
		"PROTO_NC_CHAR_OPTION_SOUND":                 1,
		"PROTO_NC_CHAR_OPTION_VIDEO":                 60,
		"PROTO_NC_CHAR_OPTION_WINDOWPOS":             392,
		"PROTO_NC_CHAR_PET_CMD":                      1,
		"PROTO_NC_CHAR_QUEST_DOING_CMD":              6,
		"PROTO_NC_CHAR_QUEST_READ_CMD":               6,
		"PROTO_NC_CHAR_QUEST_REPEAT_CMD":             6,
		"PROTO_NC_CHAR_SEAWAR_CMD":                   1,
		"PROTO_NC_CHAR_SINGLE_OPTION_CMD":            15,
		"PROTO_NC_CHAR_SKILL_CMD":                    9,
		"PROTO_NC_CHAR_SKILL_PASSIVE_CMD":            2,
		"PROTO_NC_CHAR_WEDDINGDATA_ACK":              43,
		"PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK":      14,
		"PROTO_NC_CHAT_RESTRICT_DB_LIST_CMD":         2,
		"PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_REQ": 12,
		"PROTO_NC_ITEMDB_CHARGED_LIST_ACK":           11,
		"PROTO_NC_ITEM_BUY_REQ":                      6,
		"PROTO_NC_ITEM_DROP_REQ":                     14,
		"PROTO_NC_ITEM_SELL_REQ":                     5,
	}

	for name, expectedSize := range expectedSizes {
		instance, exists := pdb.AllPDBStructs[name]
		if !exists {
			t.Errorf("Missing struct in pdb.AllPDBStructs: %s", name)
			continue
		}

		actualSize := binary.Size(instance)
		if actualSize != expectedSize {
			t.Errorf("Struct %s sizeof mismatch: expected %d bytes, got %d bytes",
				name, expectedSize, actualSize)
		}
	}

	t.Logf("Control 4: Successfully verified %d struct sizes against binary.Size", len(expectedSizes))
}

// TestFase04_Control4_NoDeprecatedAliases escanea el repositorio en busca de
// identificadores deprecados no oficiales (PROTO_NC_USER_FILE_CHECK_*, etc.).
func TestFase04_Control4_NoDeprecatedAliases(t *testing.T) {
	banned := []string{
		"PROTO_NC_USER_FILE_CHECK_REQ",
		"PROTO_NC_USER_FILE_CHECK_ACK",
		"PROTO_NC_USER_WORLDSTATUS_REQ",
		"PROTO_NC_USER_WORLDSTATUS_ACK",
		"PROTO_NC_MISC_SEED_CMD",
	}

	rootDir := "../../pkg"
	var violations []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Solo escanear archivos .go (excluyendo este test)
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "fase04_protocol_sizeof_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		strContent := string(content)
		for _, b := range banned {
			if strings.Contains(strContent, b) {
				violations = append(violations, fmt.Sprintf("%s in %s", b, path))
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Walk error: %v", err)
	}

	if len(violations) > 0 {
		t.Fatalf("Found banned deprecated aliases in codebase:\n  %s",
			strings.Join(violations, "\n  "))
	}
	t.Log("Control 4: Zero deprecated aliases found in codebase")
}
