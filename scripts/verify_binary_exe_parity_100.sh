#!/usr/bin/env bash
# ==============================================================================
# verify_binary_exe_parity_100.sh
# Suite Final de Auditoría Binaria 1:1 contra EXE/PDB (Fase 17 a Fase 21)
# WorldService (167 C2S), ZoneService (245 C2S), Departamentos Sociales & S2S Vivo
# ==============================================================================
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${ROOT_DIR}"

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m'

echo -e "${BLUE}${BOLD}======================================================================${NC}"
echo -e "${BLUE}${BOLD}   XELIONCORE: AUDITORÍA BINARIA 1:1 EXE/PDB FINAL (FASE 21 - 100%)   ${NC}"
echo -e "${BLUE}${BOLD}   (World 167 C2S, Zone 245 C2S, Social, SHBD, S2S Vivo y AST S2C)   ${NC}"
echo -e "${BLUE}${BOLD}======================================================================${NC}"
echo ""

FAIL_COUNT=0
PASS_COUNT=0

run_control() {
    local num="$1"
    local name="$2"
    local test_pattern="$3"

    echo -e "${YELLOW}${BOLD}>> Ejecutando Control ${num}/8: ${name}...${NC}"
    if go test -v -count=1 ./pkg/audit -run "${test_pattern}"; then
        echo -e "${GREEN}${BOLD}>> [PASS] Control ${num}/8: ${name}${NC}\n"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}${BOLD}>> [FAIL] Control ${num}/8: ${name}${NC}\n"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
}

run_control "1" "Herramienta disasm_pdb_handler.py & Fuentes Oficiales PDB/MD" "^TestFase21_Control1_PDB_DisassemblyToolAndGroundTruth$"
run_control "2" "Paridad In-Game Chunk, SHBD (RouCos02/RouN), Movimiento & Chat" "^TestFase21_Control2_InGameChunk_SHBD_Movement_Chat_Parity$"
run_control "3" "Paridad NC_BRIEFINFO, NC_BAT, NC_SKILL y NC_SOULSTONE" "^TestFase21_Control3_BriefInfo_Bat_Skill_Soulstone_Parity$"
run_control "4" "Cobertura 48/48 NC_ITEM, CELLCHANGE, EQUIP, CHARGED y REWARD" "^TestFase21_Control4_Item48_CellChange_Equip_Charged_Reward_Parity$"
run_control "5" "Paridad NC_MENU, NC_QUEST, Portales, MINIHOUSE, PET, CT, COLLECT, DICE y GAMBLE" "^TestFase21_Control5_Menu_Quest_Portals_Subsystems_Parity$"
run_control "6" "Cobertura 167/167 World C2S, 245/245 Zone C2S y Departamentos Sociales 1:1" "^TestFase21_Control6_World167C2S_Zone245C2S_And_SocialDepartments_Parity$"
run_control "7" "E2E Clúster en Vivo (Login -> World -> Zone -> S2S Return-to-CharSelect)" "^TestFase21_Control7_E2E_LiveCluster_S2S_ReturnToCharSelect$"
run_control "8" "Escaneo AST Cero SQL Inline, Cero Hardcodes y Paridad S2C Zone.md/WorldManager.md" "^TestFase21_Control8_AST_S2C_WireSizes_ZeroSQL_ZeroHardcode$"

echo -e "${BLUE}${BOLD}======================================================================${NC}"
echo -e "${BLUE}${BOLD}   RESUMEN FINAL AUDITORÍA BINARIA EXE/PDB (FASE 21): ${PASS_COUNT}/8 [PASS]  ${NC}"
echo -e "${BLUE}${BOLD}======================================================================${NC}"

if [ ${FAIL_COUNT} -eq 0 ] && [ ${PASS_COUNT} -eq 8 ]; then
    echo -e "${GREEN}${BOLD}>>> TODOS LOS 8 CONTROLES OBLIGATORIOS PASARON AL 100% (8/8 [PASS]) <<<${NC}"
    exit 0
else
    echo -e "${RED}${BOLD}>>> ATENCIÓN: ${FAIL_COUNT} control(es) fallido(s) (${PASS_COUNT}/8 [PASS]). <<<${NC}"
    exit 1
fi

