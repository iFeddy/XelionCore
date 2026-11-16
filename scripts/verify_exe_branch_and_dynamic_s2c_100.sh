#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

GREEN='\033[0;32m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

PASS=0
TOTAL=8

run_control() {
    local idx="$1"
    local title="$2"
    local pattern="$3"
    echo -e "${CYAN}[${idx}/${TOTAL}]${NC} ${title}..."
    if go test -v -count=1 ./pkg/audit -run "^${pattern}$"; then
        echo -e "${GREEN}[PASS]${NC} ${title}"
        PASS=$((PASS + 1))
    else
        echo -e "${RED}[FAIL]${NC} ${title}"
        exit 1
    fi
}

echo "==========================================================================="
echo "  XelionCore - Auditoría Completa Rama-por-Rama, Lua & Dynamic S2C (FASE 26)"
echo "==========================================================================="

run_control 1 "Control 1: Dynamic S2C 567/567 NewPacketWriter & Formulas" "TestFase26_Control1_DynamicS2C_567NewPacketWriter"
run_control 2 "Control 2: Fase 23 Entry/World/Move/AOI/Chat Branches" "TestFase26_Control2_Fase23_EntryWorldMoveAOIChat"
run_control 3 "Control 3: Fase 24 Combat Formulas, Timers & Mob AI" "TestFase26_Control3_Fase24_CombatFormulasTimersMobsDrops"
run_control 4 "Control 4: Fase 25 48/48 NC_ITEM, Trade, Booth & Forge" "TestFase26_Control4_Fase25_Economy48ItemTradeBoothForge"
run_control 5 "Control 5: Quests, Lua 274 APIs + Real Scripts, MiniHouse, Pet, Dice, Gamble" "TestFase26_Control5_QuestsLua274ScriptsMiniHousePetDiceGamble"
run_control 6 "Control 6: Social Departments, Autonomous World Timers & S2S" "TestFase26_Control6_SocialDepartments_Timers_S2S"
run_control 7 "Control 7: Live E2E Cluster Multi-Branch + KQ Lua + Timers" "TestFase26_Control7_E2E_Cluster_MultiBranch_Lua_Timers"
run_control 8 "Control 8: 412 Handlers C2S, AST 0 SQL/Hardcodes & 5 Suites Check" "TestFase26_Control8_All5Suites_And_412Handlers_BranchCompleteness"

echo "==========================================================================="
echo -e "  RESUMEN FASE 26: ${GREEN}${PASS}/${TOTAL} [PASS]${NC}"
echo "==========================================================================="
