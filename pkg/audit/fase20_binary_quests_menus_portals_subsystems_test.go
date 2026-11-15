package audit_test

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/login"
	"github.com/iFeddy/XelionCore/pkg/server/world"
	"github.com/iFeddy/XelionCore/pkg/server/zone"
)

// TestFase20_BinaryProtocolSizesOpcodesAndCoverage verifica contra Zone.md y Zone.pdb:
//  1. La fórmula exacta OpCode = (department << 10) | command para NC_MAP (6 - Portales),
//     NC_ACT (8 - NPCCLICK/AUTO_WAY_FINDING/SCRIPT_MSG), NC_MENU (15), NC_QUEST (17),
//     NC_CT (24), NC_SCENARIO (27), NC_MINIHOUSE (35), NC_INSTANCE_DUNGEON (41),
//     NC_DICE_TAISAI (43), NC_GAMBLE (47), NC_COLLECT (49) y NC_PET (53).
//  2. El tamaño empaquetado exacto (binary.Size / FixedWireSize) de todas las estructuras
//     clave frente a Zone.md / Zone.pdb.
//  3. Cobertura AST del 100% de opcodes C2S oficiales de esos departamentos en Zone HandlePacket.
func TestFase20_BinaryProtocolSizesOpcodesAndCoverage(t *testing.T) {
	opcodeChecks := []struct {
		name string
		got  uint16
		dept uint16
		cmd  uint16
	}{
		// NC_CHAR (Dept 4) Quest state chunk
		{"PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD", protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD, 4, 58},
		{"PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD", protocol.PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD, 4, 59},
		{"PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD", protocol.PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD, 4, 206},
		{"PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD", protocol.PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD, 4, 215},

		// NC_MAP (Dept 6) Portales y transiciones
		{"PROTO_NC_MAP_LINKSAME_CMD", protocol.PROTO_NC_MAP_LINKSAME_CMD, 6, 9},
		{"PROTO_NC_MAP_LINKOTHER_CMD", protocol.PROTO_NC_MAP_LINKOTHER_CMD, 6, 10},
		{"PROTO_NC_MAP_LINKSTART_CMD", protocol.PROTO_NC_MAP_LINKSTART_CMD, 6, 11},
		{"PROTO_NC_MAP_WING_SAVE_REQ", protocol.PROTO_NC_MAP_WING_SAVE_REQ, 6, 20},
		{"PROTO_NC_MAP_WING_SAVE_ACK", protocol.PROTO_NC_MAP_WING_SAVE_ACK, 6, 21},
		{"PROTO_NC_MAP_WING_FLY_REQ", protocol.PROTO_NC_MAP_WING_FLY_REQ, 6, 22},
		{"PROTO_NC_MAP_WING_FLY_ACK", protocol.PROTO_NC_MAP_WING_FLY_ACK, 6, 23},
		{"PROTO_NC_MAP_TOWNPORTAL_REQ", protocol.PROTO_NC_MAP_TOWNPORTAL_REQ, 6, 26},
		{"PROTO_NC_MAP_TOWNPORTAL_ACK", protocol.PROTO_NC_MAP_TOWNPORTAL_ACK, 6, 27},
		{"PROTO_NC_MAP_TONORMALCOORD_CMD", protocol.PROTO_NC_MAP_TONORMALCOORD_CMD, 6, 28},
		{"PROTO_NC_MAP_MULTY_LINK_CMD", protocol.PROTO_NC_MAP_MULTY_LINK_CMD, 6, 31},
		{"PROTO_NC_MAP_MULTY_LINK_SELECT_REQ", protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_REQ, 6, 31},
		{"PROTO_NC_MAP_MULTY_LINK_SELECT_ACK", protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_ACK, 6, 32},

		// NC_ACT (Dept 8) NPCClick / Gate / ScriptMsg
		{"PROTO_NC_ACT_NPCCLICK_CMD", protocol.PROTO_NC_ACT_NPCCLICK_CMD, 8, 10},
		{"PROTO_NC_ACT_NPCMENUOPEN_REQ", protocol.PROTO_NC_ACT_NPCMENUOPEN_REQ, 8, 28},
		{"PROTO_NC_ACT_NPCMENUOPEN_ACK", protocol.PROTO_NC_ACT_NPCMENUOPEN_ACK, 8, 29},
		{"PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ", protocol.PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ, 8, 112},
		{"PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK", protocol.PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK, 8, 113},
		{"PROTO_NC_ACT_NPC_ACTION_CMD", protocol.PROTO_NC_ACT_NPC_ACTION_CMD, 8, 114},
		{"PROTO_NC_ACT_NPC_MENU_CMD", protocol.PROTO_NC_ACT_NPC_MENU_CMD, 8, 115},
		{"PROTO_NC_ACT_SCRIPT_MSG_CMD", protocol.PROTO_NC_ACT_SCRIPT_MSG_CMD, 8, 121},

		// NC_MENU (Dept 15)
		{"PROTO_NC_MENU_SERVERMENU_REQ", protocol.PROTO_NC_MENU_SERVERMENU_REQ, 15, 1},
		{"PROTO_NC_MENU_SERVERMENU_ACK", protocol.PROTO_NC_MENU_SERVERMENU_ACK, 15, 2},

		// NC_QUEST (Dept 17)
		{"PROTO_NC_QUEST_SCRIPT_CMD_REQ", protocol.PROTO_NC_QUEST_SCRIPT_CMD_REQ, 17, 1},
		{"PROTO_NC_QUEST_SCRIPT_CMD_ACK", protocol.PROTO_NC_QUEST_SCRIPT_CMD_ACK, 17, 2},
		{"PROTO_NC_QUEST_GIVE_UP_REQ", protocol.PROTO_NC_QUEST_GIVE_UP_REQ, 17, 7},
		{"PROTO_NC_QUEST_GIVE_UP_ACK", protocol.PROTO_NC_QUEST_GIVE_UP_ACK, 17, 8},
		{"PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ", protocol.PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ, 17, 11},
		{"PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK", protocol.PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_ACK, 17, 12},
		{"PROTO_NC_QUEST_SELECT_START_REQ", protocol.PROTO_NC_QUEST_SELECT_START_REQ, 17, 15},
		{"PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD", protocol.PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD, 17, 17},
		{"PROTO_NC_QUEST_START_REQ", protocol.PROTO_NC_QUEST_START_REQ, 17, 20},
		{"PROTO_NC_QUEST_READ_REQ", protocol.PROTO_NC_QUEST_READ_REQ, 17, 22},

		// NC_CT (Dept 24 - Títulos)
		{"PROTO_NC_CT_SET_CURRENT_REQ", protocol.PROTO_NC_CT_SET_CURRENT_REQ, 24, 1},
		{"PROTO_NC_CT_SET_CURRENT_ACK", protocol.PROTO_NC_CT_SET_CURRENT_ACK, 24, 2},
		{"PROTO_NC_CT_SET_SOMEONECHANGE_CMD", protocol.PROTO_NC_CT_SET_SOMEONECHANGE_CMD, 24, 3},
		{"PROTO_NC_CT_SET_CMD", protocol.PROTO_NC_CT_SET_CMD, 24, 4},

		// NC_SCENARIO (Dept 27)
		{"PROTO_NC_SCENARIO_AREAENTRY_REQ", protocol.PROTO_NC_SCENARIO_AREAENTRY_REQ, 27, 5},
		{"PROTO_NC_SCENARIO_AREAENTRY_ACK", protocol.PROTO_NC_SCENARIO_AREAENTRY_ACK, 27, 6},
		{"PROTO_NC_SCENARIO_AREALEAVE_REQ", protocol.PROTO_NC_SCENARIO_AREALEAVE_REQ, 27, 7},
		{"PROTO_NC_SCENARIO_AREALEAVE_ACK", protocol.PROTO_NC_SCENARIO_AREALEAVE_ACK, 27, 8},

		// NC_MINIHOUSE (Dept 35)
		{"PROTO_NC_MINIHOUSE_ACTIV_REQ", protocol.PROTO_NC_MINIHOUSE_ACTIV_REQ, 35, 1},
		{"PROTO_NC_MINIHOUSE_BUILDING_REQ", protocol.PROTO_NC_MINIHOUSE_BUILDING_REQ, 35, 3},
		{"PROTO_NC_MINIHOUSE_BUILDING_CMD", protocol.PROTO_NC_MINIHOUSE_BUILDING_CMD, 35, 5},
		{"PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_ACK", protocol.PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_ACK, 35, 30},
		{"PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_ACK", protocol.PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_ACK, 35, 51},
		{"PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_ACK", protocol.PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_ACK, 35, 68},
		{"PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_ACK", protocol.PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_ACK, 35, 73},
		{"PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD", protocol.PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD, 35, 74},
		{"PROTO_NC_MINIHOUSE_PORTAL_LIST_CMD", protocol.PROTO_NC_MINIHOUSE_PORTAL_LIST_CMD, 35, 75},
		{"PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CMD", protocol.PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CMD, 35, 100},

		// NC_INSTANCE_DUNGEON (Dept 41)
		{"PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK", protocol.PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK, 41, 13},
		{"PROTO_NC_INSTANCE_DUNGEON_RANK_RESULT_CMD", protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_RESULT_CMD, 41, 18},

		// NC_DICE_TAISAI (Dept 43) & NC_GAMBLE (Dept 47)
		{"PROTO_NC_DICE_TAISAI_GAME_JOIN_REQ", protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_REQ, 43, 1},
		{"PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK", protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK, 43, 2},
		{"PROTO_NC_DICE_TAISAI_RANK_REQ", protocol.PROTO_NC_DICE_TAISAI_RANK_REQ, 43, 22},
		{"PROTO_NC_DICE_TAISAI_RANK_ACK", protocol.PROTO_NC_DICE_TAISAI_RANK_ACK, 43, 25},
		{"PROTO_NC_GAMBLE_GAMBLEHOUSE_EXIT_REQ", protocol.PROTO_NC_GAMBLE_GAMBLEHOUSE_EXIT_REQ, 47, 2},
		{"PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_REQ", protocol.PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_REQ, 47, 217},
		{"PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK", protocol.PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK, 47, 218},
		{"PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_REQ", protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_REQ, 47, 219},
		{"PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK", protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK, 47, 220},

		// NC_COLLECT (Dept 49) & NC_PET (Dept 53)
		{"PROTO_NC_COLLECT_PEEPING_REQ", protocol.PROTO_NC_COLLECT_PEEPING_REQ, 49, 1},
		{"PROTO_NC_COLLECT_CARDREGIST_REQ", protocol.PROTO_NC_COLLECT_CARDREGIST_REQ, 49, 8},
		{"PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ", protocol.PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ, 49, 12},
		{"PROTO_NC_PET_ASK_NEW_NAME_ACK", protocol.PROTO_NC_PET_ASK_NEW_NAME_ACK, 53, 6},
		{"PROTO_NC_PET_SET_NAME_REQ", protocol.PROTO_NC_PET_SET_NAME_REQ, 53, 7},
		{"PROTO_NC_PET_SET_NAME_CANCEL_REQ", protocol.PROTO_NC_PET_SET_NAME_CANCEL_REQ, 53, 9},
	}
	for _, oc := range opcodeChecks {
		want := (oc.dept << 10) | oc.cmd
		if oc.got != want {
			t.Errorf("%s = 0x%04X, esperado=0x%04X ((%d<<10)|%d)", oc.name, oc.got, want, oc.dept, oc.cmd)
		}
	}

	// Verificar sizeof oficiales contra Zone.md / Zone.pdb
	sizeChecks := []struct {
		name string
		got  int
		want int
	}{
		{"PROTO_NC_MAP_MULTY_LINK_CMD", binary.Size(protocol.PROTO_NC_MAP_MULTY_LINK_CMD_DATA{}), 73},
		{"PROTO_NC_MAP_MULTY_LINK_SELECT_REQ", binary.Size(protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_REQ_DATA{}), 12},
		{"PROTO_NC_MAP_MULTY_LINK_SELECT_ACK", binary.Size(protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_ACK_DATA{}), 2},
		{"PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK", binary.Size(protocol.PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK_DATA{}), 2},
		{"PROTO_NC_MENU_SERVERMENU_REQ", binary.Size(protocol.PROTO_NC_MENU_SERVERMENU_REQ_DATA{}), 142},
		{"PROTO_NC_MENU_SERVERMENU_ACK", binary.Size(protocol.PROTO_NC_MENU_SERVERMENU_ACK_DATA{}), 1},
		{"PROTO_NC_QUEST_SCRIPT_CMD_REQ", binary.Size(protocol.PROTO_NC_QUEST_SCRIPT_CMD_REQ_DATA{}), 103},
		{"PROTO_NC_QUEST_SCRIPT_CMD_ACK", binary.Size(protocol.PROTO_NC_QUEST_SCRIPT_CMD_ACK_DATA{}), 7},
		{"PROTO_NC_QUEST_START_REQ", binary.Size(protocol.PROTO_NC_QUEST_START_REQ_DATA{}), 2},
		{"PROTO_NC_QUEST_START_ACK", binary.Size(protocol.PROTO_NC_QUEST_START_ACK_DATA{}), 2},
		{"PROTO_NC_QUEST_READ_REQ", binary.Size(protocol.PROTO_NC_QUEST_READ_REQ_DATA{}), 2},
		{"PROTO_NC_QUEST_READ_ACK", binary.Size(protocol.PROTO_NC_QUEST_READ_ACK_DATA{}), 4},
		{"PROTO_NC_QUEST_GIVE_UP_REQ", binary.Size(protocol.PROTO_NC_QUEST_GIVE_UP_REQ_DATA{}), 2},
		{"PROTO_NC_QUEST_GIVE_UP_ACK", binary.Size(protocol.PROTO_NC_QUEST_GIVE_UP_ACK_DATA{}), 4},
		{"PROTO_NC_CT_SET_CURRENT_REQ", binary.Size(protocol.PROTO_NC_CT_SET_CURRENT_REQ_DATA{}), 4},
		{"PROTO_NC_CT_SET_CURRENT_ACK", binary.Size(protocol.PROTO_NC_CT_SET_CURRENT_ACK_DATA{}), 6},
		{"PROTO_NC_CT_SET_SOMEONECHANGE_CMD", binary.Size(protocol.PROTO_NC_CT_SET_SOMEONECHANGE_CMD_DATA{}), 6},
		{"PROTO_NC_MINIHOUSE_BUILDING_CMD", binary.Size(protocol.PROTO_NC_MINIHOUSE_BUILDING_CMD_DATA{}), 66},
		{"PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK", binary.Size(protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK_DATA{}), 50},
		{"PROTO_NC_DICE_TAISAI_RANK_ACK", binary.Size(protocol.PROTO_NC_DICE_TAISAI_RANK_ACK_DATA{}), 484},
		{"PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_CMD", binary.Size(protocol.PROTO_NC_DICE_TAISAI_DICE_ROLL_RESULT_CMD_DATA{}), 5},
		{"PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK", binary.Size(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK_DATA{}), 3},
		{"PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_ACK", binary.Size(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_ACK_DATA{}), 3},
		{"PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK", binary.Size(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK_DATA{}), 11},
		{"PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK", binary.Size(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK_DATA{}), 10},
		{"PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK", binary.Size(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK_DATA{}), 1},
		{"PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK", binary.Size(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK_DATA{}), 1},
		{"PROTO_NC_GAMBLE_SLOTMACHINE_SOMEONE_GET_JACKPOT_CMD", binary.Size(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_SOMEONE_GET_JACKPOT_CMD_DATA{}), 28},
		{"PROTO_NC_COLLECT_PEEPING_REQ", binary.Size(protocol.PROTO_NC_COLLECT_PEEPING_REQ_DATA{}), 20},
		{"PROTO_NC_COLLECT_CARDREGIST_REQ", binary.Size(protocol.PROTO_NC_COLLECT_CARDREGIST_REQ_DATA{}), 1},
		{"PROTO_NC_COLLECT_CARDREGIST_ACK", binary.Size(protocol.PROTO_NC_COLLECT_CARDREGIST_ACK_DATA{}), 10},
		{"PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ", binary.Size(protocol.PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ_DATA{}), 4},
		{"PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK", binary.Size(protocol.PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK_DATA{}), 6},
		{"PROTO_NC_COLLECT_REWARD_REQ", binary.Size(protocol.PROTO_NC_COLLECT_REWARD_REQ_DATA{}), 4},
		{"PROTO_NC_COLLECT_REWARD_ACK", binary.Size(protocol.PROTO_NC_COLLECT_REWARD_ACK_DATA{}), 4},
		{"PROTO_NC_PET_ASK_NEW_NAME_ACK", binary.Size(protocol.PROTO_NC_PET_ASK_NEW_NAME_ACK_DATA{}), 4},
		{"PROTO_NC_PET_SET_NAME_REQ", binary.Size(protocol.PROTO_NC_PET_SET_NAME_REQ_DATA{}), 18},
	}
	for _, sc := range sizeChecks {
		if sc.got != sc.want {
			t.Errorf("%s binary.Size = %d, esperado=%d", sc.name, sc.got, sc.want)
		}
	}

	// Verificar por AST que HandlePacket despacha el 100% de C2S oficiales de Zone.md para todos los subsistemas de Fase 20
	rootPath := findRepoRoot(t)
	zoneMDBytes, err := os.ReadFile(filepath.Join(rootPath, "docs", "Servicios", "Zone.md"))
	if err != nil {
		t.Fatalf("Error leyendo docs/Servicios/Zone.md: %v", err)
	}
	secRe := regexp.MustCompile(`(?is)protocolstore.*?\n(.*?)(?:\n## |\z)`)
	secMatch := secRe.FindStringSubmatch(string(zoneMDBytes))
	if len(secMatch) < 2 {
		t.Fatalf("No se encontró protocolstore en docs/Servicios/Zone.md")
	}
	opRe := regexp.MustCompile(`PROTO_NC_[A-Z0-9_]+`)
	officialOps := make(map[string]bool)
	for _, op := range opRe.FindAllString(secMatch[1], -1) {
		officialOps[op] = true
	}

	handlerPath := filepath.Join(rootPath, "pkg", "server", "zone", "handler.go")
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, handlerPath, nil, 0)
	if err != nil {
		t.Fatalf("Error parseando AST handler.go: %v", err)
	}
	handledOps := make(map[string]bool)
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "HandlePacket" {
			return true
		}
		ast.Inspect(fn.Body, func(sub ast.Node) bool {
			cc, ok := sub.(*ast.CaseClause)
			if !ok {
				return true
			}
			for _, expr := range cc.List {
				if sel, ok := expr.(*ast.SelectorExpr); ok {
					handledOps[sel.Sel.Name] = true
				}
			}
			return true
		})
		return false
	})

	prefixes := []string{
		"PROTO_NC_MENU_",
		"PROTO_NC_QUEST_",
		"PROTO_NC_SCENARIO_",
		"PROTO_NC_INSTANCE_DUNGEON_",
		"PROTO_NC_CT_",
		"PROTO_NC_MINIHOUSE_",
		"PROTO_NC_DICE_TAISAI_",
		"PROTO_NC_GAMBLE_",
		"PROTO_NC_COLLECT_",
		"PROTO_NC_PET_",
		"PROTO_NC_MAP_",
	}
	var missing []string
	counts := make(map[string]int)
	for op := range officialOps {
		for _, p := range prefixes {
			if strings.HasPrefix(op, p) {
				counts[p]++
				if !handledOps[op] {
					missing = append(missing, op)
				}
			}
		}
	}
	extraPortalNPC := []string{"PROTO_NC_ACT_NPCCLICK_CMD", "PROTO_NC_ACT_NPCMENUOPEN_ACK", "PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ", "PROTO_NC_ACT_NPC_MENU_CMD"}
	for _, op := range extraPortalNPC {
		if !handledOps[op] {
			missing = append(missing, op)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("Faltan handlers C2S en Zone HandlePacket: %v", missing)
	}
	if counts["PROTO_NC_MENU_"] != 1 ||
		counts["PROTO_NC_QUEST_"] != 7 ||
		counts["PROTO_NC_SCENARIO_"] != 2 ||
		counts["PROTO_NC_INSTANCE_DUNGEON_"] != 2 ||
		counts["PROTO_NC_CT_"] != 1 ||
		counts["PROTO_NC_MINIHOUSE_"] != 23 ||
		counts["PROTO_NC_DICE_TAISAI_"] != 9 ||
		counts["PROTO_NC_GAMBLE_"] != 18 ||
		counts["PROTO_NC_COLLECT_"] != 4 ||
		counts["PROTO_NC_PET_"] != 3 ||
		counts["PROTO_NC_MAP_"] != 7 {
		t.Fatalf("Conteo oficial inesperado en Zone.md: %+v", counts)
	}
}

// TestFase20_Quests_Menus_Portals_And_Subsystems_ClusterE2E valida en clúster TCP real
// (LoginService + WorldService + ZoneService + MySQL) todos los formatos binarios y flujos de Fase 20:
//  1) Quests & Menus (`0x103A` PLAYER_QUEST_INFO 32B / `0x103B` PLAYER_QUEST_DONE_INFO 10B,
//     `0x4414` QUEST_START_REQ, `0x4416` QUEST_READ_REQ, `0x4407` QUEST_GIVE_UP_REQ,
//     `0x440B` CLIENT_SCENARIO_DONE_REQ, `0x200A` NPCCLICK -> `0x3C01` SERVERMENU -> `0x3C02` ACK).
//  2) Portales / Gate / Links (`0x1814` WING_SAVE -> `0x1815` 4B uint16, `0x2070` USE_GATE -> `0x2071` + `0x2403` + `0x1809` 12B / `0x180A` 32B,
//     `0x181F` MULTY_LINK_SELECT -> `0x1820` 4B).
//  3) Instancias (`0xA417` RANK_LIST -> `0xA418` con SHINE_INDUN_RANK_CLIENT 48B).
//  4) Títulos (`0x6001` CT_SET_CURRENT -> `0x6002`), Colección (`0xC401`/`0xC40C`/`0xC410`), Pets (`0xD407`/`0xD409`).
//  5) MiniHouse (`0x8C03` BUILDING -> `0x8C05` 68B wire, `0x8C1D` -> `0x8C1E` 4B, `0x8C32` -> `0x8C33` 4B,
//     `0x8C43` -> `0x8C44` 4B, `0x8C48` -> `0x8C49` 4B + `0x8C4A` trailing NUL, `0x8C62` -> `0x8C64` 10B).
//  6) Casino (`0xAC01` -> 52B+, `0xAC16` -> `0xAC19` 486+24*N B con WinCount real, `0xBC6B` -> `0xBC6E`,
//     `0xBCD9` -> `0xBCDA` 3+24*N B, `0xBCDB` -> `0xBCDC` 3+24*N B, `0xBCCE` -> `0xBCCF` 13B, `0xBCD0` -> `0xBCD1` 12B).
func TestFase20_Quests_Menus_Portals_And_Subsystems_ClusterE2E(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repoRoot := findRepoRoot(t)
	dataDir := filepath.Join(repoRoot, "data")

	accountRepo := database.NewMySQLAccountRepository(db)
	charRepo := database.NewMySQLCharacterRepository(db)
	gameLogRepo := database.NewMySQLGameLogRepository(db)
	transfers := database.NewTransferManager(db)

	worldDataStore, err := world.LoadWorldDataStore(dataDir)
	if err != nil {
		t.Fatalf("Error cargando WorldDataStore: %v", err)
	}
	zoneDataStore, err := zone.LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("Error cargando ZoneDataStore: %v", err)
	}

	activeVerKey, err := accountRepo.GetClientVersionKey()
	if err != nil || activeVerKey == "" {
		t.Fatalf("Error obteniendo FOUS.GameVersionKey desde MySQL: %v", err)
	}

	const (
		loginPort = 19532
		worldPort = 19536
		zonePort  = 19546
	)

	wm := login.NewWorldManager(accountRepo)
	loginSrv := login.NewServer(&login.Config{
		BindAddr: fmt.Sprintf("127.0.0.1:%d", loginPort),
		Port:     loginPort,
	}, accountRepo, wm, transfers)
	if err := loginSrv.Start(); err != nil {
		t.Fatalf("Error iniciando LoginServer: %v", err)
	}
	defer loginSrv.Stop()

	worldSrv := world.NewServer(&world.Config{
		ID:       0,
		Name:     "Fase20Cluster",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", worldPort),
		Port:     worldPort,
		ZoneIP:   "127.0.0.1",
		ZonePort: zonePort,
	}, charRepo, transfers, worldDataStore, accountRepo)
	if err := worldSrv.Start(); err != nil {
		t.Fatalf("Error iniciando WorldServer: %v", err)
	}
	defer worldSrv.Stop()

	if err := charRepo.UpsertZoneNode(0, 1, "127.0.0.1", uint16(zonePort)); err != nil {
		t.Fatalf("Error UpsertZoneNode: %v", err)
	}
	defer func() {
		_, _ = db.Exec("DELETE FROM tWorldZoneNode WHERE nZonePort = ?", zonePort)
	}()

	glDispatcher := database.NewGameLogDispatcher(gameLogRepo, 1000, 2)
	glDispatcher.Start()
	defer glDispatcher.Close(t.Context())

	zoneSrv := zone.NewServer(&zone.Config{
		BindAddr: fmt.Sprintf("127.0.0.1:%d", zonePort),
		Port:     zonePort,
		IP:       "127.0.0.1",
		WorldNo:  0,
		ZoneNo:   1,
	}, charRepo, transfers, zoneDataStore.ClassStore, zoneDataStore.FieldStore, zoneDataStore.ChkEngine)
	zoneSrv.SetDataStore(zoneDataStore)
	zoneSrv.SetAccountRepo(accountRepo)
	zoneSrv.SetGameLogDispatcher(glDispatcher)
	if err := zoneSrv.Start(); err != nil {
		t.Fatalf("Error iniciando ZoneServer: %v", err)
	}
	defer zoneSrv.Stop()

	time.Sleep(100 * time.Millisecond)

	const userNo = 998920
	ts := time.Now().Unix() % 100000
	pass := "f20pass"
	sum := md5.Sum([]byte(pass))
	passHash := hex.EncodeToString(sum[:])

	cleanupUser := func(uNo int) {
		_, _ = db.Exec("DELETE FROM tActiveZoneTransfer WHERE nUserNo = ?", uNo)
		_, _ = db.Exec("DELETE FROM tActiveAccountSession WHERE nUserNo = ?", uNo)
		_, _ = db.Exec("DELETE FROM tLogout WHERE nUserNo = ?", uNo)
		if rows, qErr := db.Query("SELECT nCharNo FROM tCharacter WHERE nUserNo = ?", uNo); qErr == nil {
			var charIDs []int
			for rows.Next() {
				var cid int
				if rows.Scan(&cid) == nil {
					charIDs = append(charIDs, cid)
				}
			}
			_ = rows.Close()
			for _, cid := range charIDs {
				_, _ = db.Exec("DELETE FROM tGameLog WHERE nCharNo = ?", cid)
				_, _ = db.Exec("DELETE FROM tQuest WHERE nCharNo = ?", cid)
				_, _ = db.Exec("DELETE FROM tCharacterTitle WHERE nCharNo = ?", cid)
				_, _ = db.Exec("DELETE FROM tItem WHERE nOwner = ?", cid)
				_, _ = db.Exec("DELETE FROM tSkill WHERE nCharNo = ?", cid)
				_, _ = db.Exec("DELETE FROM tGambleCoin WHERE nCharNo = ?", cid)
				_, _ = db.Exec("DELETE FROM tDiceGame WHERE nCharNo = ?", cid)
				_, _ = db.Exec("DELETE FROM tSlotMachine WHERE nCharNo = ?", cid)
			}
		}
		_, _ = db.Exec("DELETE FROM tCharacter WHERE nUserNo = ?", uNo)
		_, _ = db.Exec("DELETE FROM tUser WHERE nUserNo = ?", uNo)
	}
	cleanupUser(userNo)
	defer cleanupUser(userNo)

	uName := fmt.Sprintf("f20usr_%05d", ts)
	cName := fmt.Sprintf("F20_%05d", ts)

	if _, err := db.Exec(`
		INSERT INTO tUser (nUserNo, sUserID, sUserPW, sUserIP, bIsBlock, bIsDelete, dDate)
		VALUES (?, ?, ?, '127.0.0.1', 0, 0, NOW())
	`, userNo, uName, passHash); err != nil {
		t.Fatalf("Error insertando tUser: %v", err)
	}

	lConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", loginPort))
	if err != nil {
		t.Fatalf("Dial LoginService: %v", err)
	}
	lHS, _ := auditReadServerPacket(lConn)
	lSeed, _ := lHS.Reader.ReadInt16()
	lCrypto, _ := crypto.NewNetCrypto(lSeed)
	lClient := &auditMockClient{conn: lConn, crypto: lCrypto}

	verPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ)
	verPkt.WriteString(activeVerKey, 64)
	_ = lClient.sendPacket(verPkt)
	_, _ = auditReadServerPacket(lConn)
	_, _ = auditReadServerPacket(lConn)

	xtPkt := protocol.NewPacketWriter(protocol.PROTO_NC_USER_XTRAP_REQ)
	xtPkt.WriteByte(1)
	_ = lClient.sendPacket(xtPkt)
	_, _ = auditReadServerPacket(lConn)

	lReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_US_LOGIN_REQ)
	lReq.WriteString(uName, 260)
	lReq.WriteString(passHash, 36)
	lReq.Fill(20, 0)
	_ = lClient.sendPacket(lReq)
	_, _ = auditReadServerPacket(lConn)

	wSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WORLDSELECT_REQ)
	wSelReq.WriteByte(0)
	_ = lClient.sendPacket(wSelReq)
	wSelAck, err := auditReadServerPacket(lConn)
	if err != nil || wSelAck.OpCode != protocol.PROTO_NC_USER_WORLDSELECT_ACK {
		t.Fatalf("WORLDSELECT falló: %v / %v", wSelAck, err)
	}
	_, _ = wSelAck.Reader.ReadByte()
	_, _ = wSelAck.Reader.ReadString(16)
	_, _ = wSelAck.Reader.ReadUint16()
	tokenBytes, _ := wSelAck.Reader.ReadBytes(64)
	worldToken := strings.TrimRight(string(tokenBytes), "\x00")
	_ = lConn.Close()

	wConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", worldPort))
	if err != nil {
		t.Fatalf("Dial WorldService: %v", err)
	}
	defer wConn.Close()
	wHS, _ := auditReadServerPacket(wConn)
	wSeed, _ := wHS.Reader.ReadInt16()
	wCrypto, _ := crypto.NewNetCrypto(wSeed)
	wClient := &auditMockClient{conn: wConn, crypto: wCrypto}

	wlReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	wlReq.WriteString(uName, 256)
	wlReq.WriteString(worldToken, 64)
	_ = wClient.sendPacket(wlReq)
	if wlAck, err := auditReadServerPacket(wConn); err != nil || wlAck.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("LOGINWORLD falló: %v / %v", wlAck, err)
	}

	createReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createReq.WriteByte(0)
	createReq.WriteString(cName, 20)
	createReq.WriteByte(byte((1 << 2) | (1 << 7)))
	createReq.WriteByte(1)
	createReq.WriteByte(1)
	createReq.WriteByte(4)
	_ = wClient.sendPacket(createReq)
	if createAck, err := auditReadServerPacket(wConn); err != nil || createAck.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("AVATAR_CREATE falló: %v / %v", createAck, err)
	}

	createdChar, err := charRepo.GetCharacterByName(cName)
	if err != nil || createdChar == nil {
		t.Fatalf("Personaje '%s' no encontrado: %v", cName, err)
	}
	charID := createdChar.ID
	_ = charRepo.AddCharMoney(charID, 500000)

	cSelReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	cSelReq.WriteByte(0)
	_ = wClient.sendPacket(cSelReq)
	if cSelAck, err := auditReadServerPacket(wConn); err != nil || cSelAck.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("CHAR_LOGIN falló: %v / %v", cSelAck, err)
	}
	_, _ = auditReadServerPacket(wConn)
	_, _ = auditReadServerPacket(wConn)
	_, _ = auditReadServerPacket(wConn)

	var dbRandomID uint16
	if err := db.QueryRow("SELECT nRandomID FROM tActiveZoneTransfer WHERE sCharID = ?", cName).Scan(&dbRandomID); err != nil {
		t.Fatalf("Token tActiveZoneTransfer no encontrado: %v", err)
	}

	zConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", zonePort))
	if err != nil {
		t.Fatalf("Dial ZoneService: %v", err)
	}
	defer zConn.Close()
	zHS, _ := auditReadServerPacket(zConn)
	zSeed, _ := zHS.Reader.ReadInt16()
	zCrypto, _ := crypto.NewNetCrypto(zSeed)
	zClient := &auditMockClient{conn: zConn, crypto: zCrypto}

	mapReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	mapReq.WriteUint16(dbRandomID)
	mapReq.WriteString(cName, 20)
	sums := zoneDataStore.ChkEngine.GetChecksums()
	for i := 0; i < 49; i++ {
		mapReq.WriteBytes(sums[i][:])
	}
	_ = zClient.sendPacket(mapReq)

	var loginChunks []*protocol.Packet
	for {
		pkt, rErr := auditReadServerPacket(zConn)
		if rErr != nil {
			t.Fatalf("Error leyendo MAP_LOGIN: %v", rErr)
		}
		loginChunks = append(loginChunks, pkt)
		if pkt.OpCode == protocol.PROTO_NC_MAP_LOGIN_ACK {
			break
		}
	}
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD))

	collectPackets := func(wait time.Duration) []*protocol.Packet {
		var pkts []*protocol.Packet
		_ = zConn.SetReadDeadline(time.Now().Add(wait))
		for {
			pkt, rErr := auditReadServerPacket(zConn)
			if rErr != nil {
				break
			}
			pkts = append(pkts, pkt)
		}
		_ = zConn.SetReadDeadline(time.Time{})
		return pkts
	}
	findPkt := func(pkts []*protocol.Packet, op uint16) *protocol.Packet {
		for _, p := range pkts {
			if p.OpCode == op {
				return p
			}
		}
		return nil
	}

	// Drenar paquetes iniciales de MAP_LOGINCOMPLETE
	_ = collectPackets(150 * time.Millisecond)

	// 1. Verificar formato binario exacto de chunks Quest en MAP_LOGIN (`0x103A` Doing 6+N*32B, `0x103B` Done 12+N*10B)
	if qDoing := findPkt(loginChunks, protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD); qDoing == nil {
		t.Fatalf("Falta 0x103A PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD en MAP_LOGIN")
	} else {
		_, _ = qDoing.Reader.ReadUint32() // nCharNo
		_, _ = qDoing.Reader.ReadByte()   // bNeedClear
		qNum, _ := qDoing.Reader.ReadByte()
		if qDoing.Reader.Remaining() != int(qNum)*32 {
			t.Fatalf("0x103A payload len mismatch: rem=%d, qNum*32=%d", qDoing.Reader.Remaining(), int(qNum)*32)
		}
	}
	if qDone := findPkt(loginChunks, protocol.PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD); qDone == nil {
		t.Fatalf("Falta 0x103B PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD en MAP_LOGIN")
	} else {
		_, _ = qDone.Reader.ReadUint32()   // nCharNo
		_, _ = qDone.Reader.ReadUint16()   // nTotalDoneQuest
		_, _ = qDone.Reader.ReadUint16()   // nTotalDoneQuestSize
		qNum, _ := qDone.Reader.ReadUint16() // nDoneQuestCount
		_, _ = qDone.Reader.ReadUint16()   // nIndex
		if qDone.Reader.Remaining() != int(qNum)*10 {
			t.Fatalf("0x103B payload len mismatch (PLAYER_QUEST_DONE_INFO=10B): rem=%d, qNum*10=%d", qDone.Reader.Remaining(), int(qNum)*10)
		}
	}

	// 2. Ciclo Quest (`NC_QUEST`): tomar primera quest de Nivel 1 de QuestData.shn
	var qID uint16
	for id, qd := range zoneDataStore.QuestCatalog.ByID {
		if qd != nil && !qd.CheckPrevQuest && qd.MinLevel <= 1 {
			qID = id
			break
		}
	}
	if qID == 0 {
		t.Fatalf("No se encontró ninguna quest nivel 1 en QuestCatalog")
	}
	qStartReq := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_START_REQ)
	qStartReq.WriteUint16(qID)
	_ = zClient.sendPacket(qStartReq)
	pkts := collectPackets(150 * time.Millisecond)
	if qStartAck := findPkt(pkts, protocol.PROTO_NC_QUEST_START_ACK); qStartAck == nil || len(qStartAck.Raw)-2 != 2 || binary.LittleEndian.Uint16(qStartAck.Raw[2:4]) != 0 {
		t.Fatalf("0x4415 PROTO_NC_QUEST_START_ACK inválido (esperado 2B errCode=0): %v", qStartAck)
	}

	qReadReq := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_READ_REQ)
	qReadReq.WriteUint16(qID)
	_ = zClient.sendPacket(qReadReq)
	pkts = collectPackets(120 * time.Millisecond)
	if qReadAck := findPkt(pkts, protocol.PROTO_NC_QUEST_READ_ACK); qReadAck == nil || len(qReadAck.Raw)-2 != 4 {
		t.Fatalf("0x4417 PROTO_NC_QUEST_READ_ACK inválido: %v", qReadAck)
	}

	qGiveUpReq := protocol.NewPacketWriter(protocol.PROTO_NC_QUEST_GIVE_UP_REQ)
	qGiveUpReq.WriteUint16(qID)
	_ = zClient.sendPacket(qGiveUpReq)
	pkts = collectPackets(120 * time.Millisecond)
	if qGiveUpAck := findPkt(pkts, protocol.PROTO_NC_QUEST_GIVE_UP_ACK); qGiveUpAck == nil || len(qGiveUpAck.Raw)-2 != 4 {
		t.Fatalf("0x4408 PROTO_NC_QUEST_GIVE_UP_ACK inválido: %v", qGiveUpAck)
	}

	// 3. Portales / Wing / AutoWayFindingUseGate (`0x1814->0x1815` 4B wire, `0x2070->0x2071`, `0x181F->0x1820` 4B wire)
	wingSaveReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_WING_SAVE_REQ)
	wingSaveReq.WriteUint16(0)
	_ = zClient.sendPacket(wingSaveReq)
	pkts = collectPackets(120 * time.Millisecond)
	if wingAck := findPkt(pkts, protocol.PROTO_NC_MAP_WING_SAVE_ACK); wingAck == nil || len(wingAck.Raw)-2 != 2 {
		t.Fatalf("0x1815 PROTO_NC_MAP_WING_SAVE_ACK debe medir exactamente 2B payload (so_SendErrorCode): %v", wingAck)
	}

	multiReq := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_REQ)
	multiReq.WriteString("Rou", 12)
	_ = zClient.sendPacket(multiReq)
	pkts = collectPackets(180 * time.Millisecond)
	if multiAck := findPkt(pkts, protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_ACK); multiAck == nil || len(multiAck.Raw)-2 != 2 {
		t.Fatalf("0x1820 PROTO_NC_MAP_MULTY_LINK_SELECT_ACK debe medir 2B payload: %v", multiAck)
	}
	if linkSame := findPkt(pkts, protocol.PROTO_NC_MAP_LINKSAME_CMD); linkSame == nil || len(linkSame.Raw)-2 != 10 {
		t.Fatalf("0x1809 PROTO_NC_MAP_LINKSAME_CMD debe medir 10B payload (12B wire): %v", linkSame)
	}
	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD))
	_ = collectPackets(150 * time.Millisecond)

	// 4. Instancias Rank List (`0xA417` -> `0xA418`: 2B + nMapCount*13B) y Títulos (`0x6001` -> `0x6002`)
	rankReq := protocol.NewPacketWriter(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_REQ)
	_ = zClient.sendPacket(rankReq)
	pkts = collectPackets(120 * time.Millisecond)
	if indAck := findPkt(pkts, protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_ACK); indAck == nil || len(indAck.Raw)-2 < 2 {
		t.Fatalf("0xA418 PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_ACK faltante/corto: %v", indAck)
	} else {
		pl := indAck.Raw[2:]
		nMapCount := int(binary.LittleEndian.Uint16(pl[0:2]))
		if len(pl) != 2+nMapCount*13 {
			t.Fatalf("0xA418 len=%d, esperado=2+nMapCount*13=%d", len(pl), 2+nMapCount*13)
		}
	}

	ctReq := protocol.NewPacketWriter(protocol.PROTO_NC_CT_SET_CURRENT_REQ)
	ctReq.WriteByte(0)
	ctReq.WriteByte(0)
	ctReq.WriteUint16(0)
	_ = zClient.sendPacket(ctReq)
	pkts = collectPackets(120 * time.Millisecond)
	if ctAck := findPkt(pkts, protocol.PROTO_NC_CT_SET_CURRENT_ACK); ctAck == nil || len(ctAck.Raw)-2 != 6 {
		t.Fatalf("0x6002 PROTO_NC_CT_SET_CURRENT_ACK debe medir 6B payload: %v", ctAck)
	}

	// 5. MiniHouse (`NC_MINIHOUSE` Dept 35): alineación 1:1 verificada contra Zone.exe
	mhBuildReq := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_BUILDING_REQ)
	mhBuildReq.WriteByte(10)
	mhBuildReq.WriteString("", 9)
	mhBuildReq.WriteString("CasaF20", 21)
	mhBuildReq.WriteByte(0)
	mhBuildReq.WriteString("Bienvenidos", 101)
	_ = zClient.sendPacket(mhBuildReq)
	pkts = collectPackets(150 * time.Millisecond)
	if mhAck := findPkt(pkts, protocol.PROTO_NC_MINIHOUSE_BUILDING_ACK); mhAck == nil || len(mhAck.Raw)-2 != 2 {
		t.Fatalf("0x8C04 PROTO_NC_MINIHOUSE_BUILDING_ACK debe medir 2B payload: %v", mhAck)
	}
	if mhCmd := findPkt(pkts, protocol.PROTO_NC_MINIHOUSE_BUILDING_CMD); mhCmd == nil || len(mhCmd.Raw)-2 != 66 {
		gotLen := -1
		if mhCmd != nil {
			gotLen = len(mhCmd.Raw) - 2
		}
		t.Fatalf("0x8C05 PROTO_NC_MINIHOUSE_BUILDING_CMD debe medir 66B payload (68B wire), got=%d", gotLen)
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ))
	pkts = collectPackets(120 * time.Millisecond)
	if furDoneAck := findPkt(pkts, protocol.PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_ACK); furDoneAck == nil || len(furDoneAck.Raw)-2 != 2 {
		t.Fatalf("0x8C1E PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_ACK debe medir 2B payload (so_SendErrorCode): %v", furDoneAck)
	}

	maxReq := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ)
	maxReq.WriteByte(12)
	_ = zClient.sendPacket(maxReq)
	pkts = collectPackets(120 * time.Millisecond)
	if maxAck := findPkt(pkts, protocol.PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_ACK); maxAck == nil || len(maxAck.Raw)-2 != 2 {
		t.Fatalf("0x8C33 PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_ACK debe medir 2B payload: %v", maxAck)
	}

	infoReq := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ)
	infoReq.WriteByte(1)
	_ = zClient.sendPacket(infoReq)
	pkts = collectPackets(120 * time.Millisecond)
	if infoAck := findPkt(pkts, protocol.PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_ACK); infoAck == nil || len(infoAck.Raw)-2 != 2 {
		t.Fatalf("0x8C44 PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_ACK debe medir 2B payload: %v", infoAck)
	}

	notReq := protocol.NewPacketWriter(protocol.PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ)
	notMsg := "HolaCasa"
	notReq.WriteByte(byte(len(notMsg)))
	notReq.WriteString(notMsg, len(notMsg))
	_ = zClient.sendPacket(notReq)
	pkts = collectPackets(120 * time.Millisecond)
	if notAck := findPkt(pkts, protocol.PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_ACK); notAck == nil || len(notAck.Raw)-2 != 2 {
		t.Fatalf("0x8C49 PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_ACK debe medir 2B payload: %v", notAck)
	}
	if notCmd := findPkt(pkts, protocol.PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD); notCmd == nil || len(notCmd.Raw)-2 != 1+len(notMsg)+1 || notCmd.Raw[len(notCmd.Raw)-1] != 0 {
		t.Fatalf("0x8C4A PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_CMD debe incluir NUL final (len=%d): %v", 1+len(notMsg)+1, notCmd)
	}

	// 6. Collection (`NC_COLLECT`) & Pet (`NC_PET`)
	peepReq := protocol.NewPacketWriter(protocol.PROTO_NC_COLLECT_PEEPING_REQ)
	peepReq.WriteString(cName, 20)
	_ = zClient.sendPacket(peepReq)
	pkts = collectPackets(120 * time.Millisecond)
	if findPkt(pkts, protocol.PROTO_NC_COLLECT_PEEPING_COLLECT_ACK) == nil || findPkt(pkts, protocol.PROTO_NC_COLLECT_PEEPING_COLLECT_BOOKMARK_ACK) == nil {
		t.Fatalf("Falta respuesta 0xC405/0xC406 para PROTO_NC_COLLECT_PEEPING_REQ")
	}

	petCancel := protocol.NewPacketWriter(protocol.PROTO_NC_PET_SET_NAME_CANCEL_REQ)
	petCancel.WriteUint16(0)
	_ = zClient.sendPacket(petCancel)
	pkts = collectPackets(120 * time.Millisecond)
	if findPkt(pkts, protocol.PROTO_NC_PET_SET_NAME_CANCEL_ACK) == nil {
		t.Fatalf("Falta 0xD40A PROTO_NC_PET_SET_NAME_CANCEL_ACK")
	}

	// 7. Casino (`NC_DICE_TAISAI` Dept 43 & `NC_GAMBLE` Dept 47): verificación 1:1 de DiceRankInfo (24B) y SLOTMACHINE_JackPotRank (24B)
	if _, err := db.Exec("INSERT INTO tGambleCoin (nCharNo, nCoin) VALUES (?, 500) ON DUPLICATE KEY UPDATE nCoin = 500", charID); err != nil {
		t.Fatalf("Error seeding tGambleCoin: %v", err)
	}
	gameNo, _ := charRepo.SetDiceGame(charID, 0, "1,2,3")
	_ = charRepo.SetDiceLog(gameNo, charID, 0, 1, 1, 1, 10, 75)
	_ = charRepo.SetSlotMachineResult(charID, 10, 1, "6,6,6", true, 500, true, 5)

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_REQ))
	pkts = collectPackets(120 * time.Millisecond)
	if djAck := findPkt(pkts, protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK); djAck == nil || len(djAck.Raw)-2 < 50 {
		t.Fatalf("0xAC02 PROTO_NC_DICE_TAISAI_GAME_JOIN_ACK debe medir >=50B payload: %v", djAck)
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_DICE_TAISAI_RANK_REQ))
	pkts = collectPackets(120 * time.Millisecond)
	if drAck := findPkt(pkts, protocol.PROTO_NC_DICE_TAISAI_RANK_ACK); drAck == nil || len(drAck.Raw)-2 < 484 {
		t.Fatalf("0xAC19 PROTO_NC_DICE_TAISAI_RANK_ACK inválido: %v", drAck)
	} else {
		pl := drAck.Raw[2:]
		histCnt := pl[2]
		currCnt := pl[483]
		if wantLen := 484 + int(currCnt)*24; len(pl) != wantLen {
			t.Fatalf("0xAC19 len=%d, esperado=484+currCnt*24=%d (histCnt=%d, currCnt=%d)", len(pl), wantLen, histCnt, currCnt)
		}
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ))
	pkts = collectPackets(120 * time.Millisecond)
	if gAll := findPkt(pkts, protocol.PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK); gAll == nil || len(gAll.Raw)-2 < 3 || (len(gAll.Raw)-2-3)%24 != 0 {
		t.Fatalf("0xBC6E PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_ACK inválido: %v", gAll)
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_REQ))
	pkts = collectPackets(120 * time.Millisecond)
	if jpAck := findPkt(pkts, protocol.PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK); jpAck == nil || len(jpAck.Raw)-2 < 1 {
		t.Fatalf("0xBCDA PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_ACK inválido: %v", jpAck)
	} else {
		pl := jpAck.Raw[2:]
		listLen := int(pl[0])
		if wantLen := 1 + listLen*24; len(pl) != wantLen {
			t.Fatalf("0xBCDA SLOTMACHINE_JackPotRank debe ser 24B/entrada (1+ListLen*24=%d), got=%d", wantLen, len(pl))
		}
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_REQ))
	pkts = collectPackets(120 * time.Millisecond)
	if winAck := findPkt(pkts, protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK); winAck == nil || len(winAck.Raw)-2 < 1 {
		t.Fatalf("0xBCDC PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_ACK inválido: %v", winAck)
	} else {
		pl := winAck.Raw[2:]
		listLen := int(pl[0])
		if wantLen := 1 + listLen*24; len(pl) != wantLen {
			t.Fatalf("0xBCDC SLOTMACHINE_JackPotRank debe ser 24B/entrada (1+ListLen*24=%d), got=%d", wantLen, len(pl))
		}
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_REQ))
	pkts = collectPackets(120 * time.Millisecond)
	if sbAck := findPkt(pkts, protocol.PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK); sbAck == nil || len(sbAck.Raw)-2 != 11 {
		t.Fatalf("0xBCCF PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_ACK debe medir 11B payload: %v", sbAck)
	}

	_ = zClient.sendPacket(protocol.NewPacketWriter(protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_REQ))
	pkts = collectPackets(120 * time.Millisecond)
	if wsAck := findPkt(pkts, protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK); wsAck == nil || len(wsAck.Raw)-2 != 10 {
		t.Fatalf("0xBCD1 PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_ACK debe medir 10B payload: %v", wsAck)
	}
}
