#!/usr/bin/env bash
# ==============================================================================
# verify_world_character_100.sh
# Suite de Auditoría Automatizada de Completitud 100% (Sección 9 del Plan Maestro)
# WorldService, Character, Topología de Zona y Gameplay In-Game Chunk
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
echo -e "${BLUE}${BOLD}   (WorldService, Character, ZoneTopology, PDB 1:1, MySQL 100%)       ${NC}"
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

# 9.1 Control 1: Tablas de World00_Character.bak y Tablas Operativas S2S
run_control "1" "Esquema DB 1:1 (Tablas de World00_Character y Tablas Operativas S2S)" "^TestFase01_WorldCharacter_Control1"

# 9.2 Control 2: 100% Stored Procedures contra information_schema.ROUTINES (425 SPs)
run_control "2" "Catálogo Exhaustivo de Stored Procedures (425 SPs de World/Character)" "^TestFase02_WorldCharacter_Control2"

# 9.3 Control 3: Escaneo AST de Go (Cero SQL Inline, Cero Mapas en Memoria de Transferencia)
run_control "3" "Escaneo AST de Go (Cero SQL Inline, 100% CALL Stored Procedures)" "^TestFase03_Control3_AST_WorldCharacter"

# 9.4 Control 4: Paridad de OpCodes y sizeof contra WorldManager.md y Character.md
run_control "4" "Paridad Oficial de Protocolo (OpCodes y sizeof al byte contra PDBs)" "^(TestFase04_Control4_WorldCharacter_)"

# 9.5 Control 5: Cobertura Total de Handlers de WorldManager.pdb (167 C2S + S2S/Persistencia)
run_control "5" "Cobertura Total de Handlers en Dispatchers (167 C2S de WorldManager.pdb)" "^(TestFase05_Control5A_|TestFase06_Control5B_)"

# 9.6 Control 6: Carga y Checksums de los 54 + 71 Archivos .shn (74 únicos) y ChecksumEngine
run_control "6" "Carga y Checksums de los 74 SHNs Únicos y ChecksumEngine" "^TestFase07_Control6_SHNFilesAndChecksums"

# 9.7 Control 7: Prueba E2E Integral Forense Multi-Servicio WorldService <-> ZoneService en MySQL
run_control "7" "E2E Integral Forense Multi-Servicio WorldService <-> ZoneService en MySQL" "^TestFase07_Control7_E2E_Forensic_MultiService"

# 9.8 Control 8: Escaneo AST Cero Datos Hardcodeados (Cero 'onson', 'gamigo', 'Mara Pirates', 'RouN')
run_control "8" "Cero Hardcode y Fidelidad Binaria PDB (Zero Magic Literals, 100% SP Remediation)" "^TestFase07_Control8_AST_ZeroHardcode"

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
