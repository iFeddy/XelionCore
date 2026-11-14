#!/usr/bin/env bash
# ==============================================================================
# verify_zone_gamelog_100.sh
# Suite de Auditoría Automatizada de Completitud 100% (Etapa 3 - Sección 9)
# ZoneService, GameLogService, Lua Runtime, S2S, AdminLvSet y Multi-Nodo
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
echo -e "${BLUE}${BOLD}   XELIONCORE: AUDITORÍA ETAPA 3 (ZONESERVICE & GAMELOG 100%)         ${NC}"
echo -e "${BLUE}${BOLD}   (176 Tablas, 27 SPs, 245 C2S, 138 WMS, 196 GDS, 180 SHN, 274 Lua)  ${NC}"
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

# 9.1 Control 1: 176 Tablas Legacy en xelion_online y Eliminación Definitiva de database/zone_tables.sql
run_control "1" "Esquema 176 Tablas + tGameLog + TEMP_QuestComplete + Sin zone_tables.sql" "^TestFase08_Control1_Schema_176_Tables_And_No_ZoneTables$"

# 9.2 Control 2: 27 Stored Procedures GameLog y Pool Seguro
run_control "2" "27 Stored Procedures GameLog contra information_schema y Concurrencia" "^TestFase08_Control2_GameLog_StoredProcedures_PoolSafe$"

# 9.3 Control 3: Escaneo AST Cero SQL Inline en ZoneService, Lua y GameLog
run_control "3" "Escaneo AST Go (Cero SQL Inline, 100% CALL Stored Procedures)" "^(TestFase08_Control3_AST_NoInlineSQL_GameLog|TestFase16_Control3_AST_ZoneAndGameLog_NoInlineSQL)$"

# 9.4 Control 4: Paridad de Protocolo Zone (245 OpCodes C2S, 1293 Structs Zone + 44 GameLog, Cero Alias Prohibidos)
run_control "4" "Paridad Oficial de Protocolo Zone & GameLog (Zone.md / GameLog.md)" "^TestFase09_Control4_"

# 9.5 Control 5: Cobertura 245/245 C2S, Cero DATA_REPORT, 138 wms_NC_*, 196 gds_NC_* y Zone00..Zone04
run_control "5" "Cobertura Total C2S (245/245), S2S (138 WMS / 196 GDS) y Multi-Nodo Zone00..04" "^TestFase16_Control5_Zone_245C2S_138WMS_196GDS_Completeness$"

# 9.6 Control 6: Carga RAM 180 SHN + 301 BlockInfo + 44 AreaBMP + 274 APIs Lua Oficiales
run_control "6" "Recursos en RAM (180 SHN, 301 BlockInfo, 44 AreaBMP, 274 APIs Lua)" "^TestFase10_Control6_"

# 9.7 Control 7: Prueba E2E Forense Multi-Servicio (Login -> World -> Zone -> GameLog -> MySQL)
run_control "7" "E2E Forense de Clúster Completo y Auditoría GameLog / AdminLvSet / tGMItemRewardLog" "^TestFase16_Control7_E2E_Forensic_Cluster_Zone_GameLog$"

# 9.8 Control 8: Escaneo AST Cero Hardcodes y Cero Privilegios GM Hardcodeados
run_control "8" "Escaneo AST Cero Hardcodes (Cero 'onson','gamigo','Mara Pirates','RouN','loadDefaults')" "^(TestFase08_Control8_|TestFase16_Control8_AST_ZeroHardcode_ZoneAndLua)$"

echo -e "${BLUE}${BOLD}======================================================================${NC}"
echo -e "${BLUE}${BOLD}   RESUMEN FINAL ETAPA 3 (ZoneService + GameLog): ${PASS_COUNT}/8 [PASS]  ${NC}"
echo -e "${BLUE}${BOLD}======================================================================${NC}"

if [ ${FAIL_COUNT} -eq 0 ] && [ ${PASS_COUNT} -eq 8 ]; then
    echo -e "${GREEN}${BOLD}>>> TODOS LOS 8 CONTROLES OBLIGATORIOS PASARON AL 100% (8/8 [PASS]) <<<${NC}"
    exit 0
else
    echo -e "${RED}${BOLD}>>> ATENCIÓN: ${FAIL_COUNT} control(es) fallido(s) (${PASS_COUNT}/8 [PASS]). <<<${NC}"
    exit 1
fi

