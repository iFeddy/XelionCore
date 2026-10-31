#!/usr/bin/env bash
# ==============================================================================
# verify_login_account_100.sh
# Suite de Auditoría Automatizada de Completitud 100% (Sección 9 del Plan Maestro)
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
NC='\033[0m' # No Color

echo -e "${BLUE}${BOLD}======================================================================${NC}"
echo -e "${BLUE}${BOLD}   XELIONCORE: AUDITORÍA AUTOMATIZADA DE COMPLETITUD 100%             ${NC}"
echo -e "${BLUE}${BOLD}   (LoginService, Account, AccountLog, StatisticsData, PDB 1:1)      ${NC}"
echo -e "${BLUE}${BOLD}======================================================================${NC}"
echo ""

FAIL_COUNT=0

run_control() {
    local num="$1"
    local name="$2"
    local test_pattern="$3"

    echo -e "${YELLOW}${BOLD}>> Ejecutando Control ${num}: ${name}...${NC}"
    if go test -v -count=1 ./pkg/audit -run "${test_pattern}"; then
        echo -e "${GREEN}${BOLD}>> [PASS] Control ${num}: ${name}${NC}\n"
    else
        echo -e "${RED}${BOLD}>> [FAIL] Control ${num}: ${name}${NC}\n"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
}

# 9.1 Control 1: Tablas, Vistas y Columnas contra information_schema
run_control "1" "Esquema DB 1:1 (Tablas, Vistas y Columnas contra information_schema)" "^TestFase01_Control1"

# 9.2 Control 2: 100% Stored Procedures contra information_schema.ROUTINES (134 SPs)
run_control "2" "Catálogo Exhaustivo de Stored Procedures (Account, AccountLog, Statistics, Session)" "^(TestFase02_Control2|TestFase03_Control2)"

# 9.3 Control 3: Escaneo AST de Go (Cero SQL Inline, Cero JSON/Memoria, CALL <SP> estricto)
run_control "3" "Escaneo AST de Go (Cero SQL Inline, 100% CALL Stored Procedures)" "^TestFase05_Control3"

# 9.4 Control 4: Paridad de OpCodes y sizeof contra Login.md, Account.md y AccountLog.md
run_control "4" "Paridad Oficial de Protocolo (OpCodes y sizeof al byte contra PDBs)" "^TestFase04_Control4"

# 9.5 Control 5: Cobertura Total de Handlers en el Dispatcher de LoginService
run_control "5" "Cobertura Total de Handlers en Dispatchers (fc_NC_* de Login.pdb)" "^TestFase06_Control5"

# 9.6 Control 6: Checksums de los 49 .shn de ChecksumEngine y persistencia en tHackScanLog
run_control "6" "Verificación de los 49 SHNs Oficiales y Persistencia en tHackScanLog" "^TestFase07_Control6"

# 9.7 Control 7: Prueba E2E Integral con Verificación Forense en MySQL
run_control "7" "E2E Integral con Verificación Forense en MySQL (Handshake -> Login -> World -> OTP -> Re-Login -> Cash -> Logout)" "^TestFase07_Control7"

# 9.8 Control 8: Cero Hardcode y Fidelidad Binaria PDB (Cero números mágicos, 100% SPs oficiales)
run_control "8" "Cero Hardcode y Fidelidad Binaria PDB (Zero Magic Literals, 100% SP Remediation)" "^TestFase08_Control8"

echo -e "${BLUE}${BOLD}======================================================================${NC}"
echo -e "${BLUE}${BOLD}   RESUMEN FINAL DE LA AUDITORÍA DE COMPLETITUD                       ${NC}"
echo -e "${BLUE}${BOLD}======================================================================${NC}"

if [ ${FAIL_COUNT} -eq 0 ]; then
    echo -e "${GREEN}${BOLD}>>> TODOS LOS 8 CONTROLES OBLIGATORIOS PASARON EN VERDE (100% PASS) <<<${NC}"
    echo -e "${GREEN}${BOLD}>>> XelionCore cumple con el 100% de fidelidad binaria, PDB y DB.    <<<${NC}"
    echo -e "${BLUE}${BOLD}======================================================================${NC}"
    exit 0
else
    echo -e "${RED}${BOLD}>>> ATENCIÓN: Hubo ${FAIL_COUNT} control(es) fallido(s).               <<<${NC}"
    echo -e "${BLUE}${BOLD}======================================================================${NC}"
    exit 1
fi
