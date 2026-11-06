#!/bin/bash
# ==============================================================================
# XelionCore - Script de Aplicación de Esquema y Stored Procedures (MySQL 8.4)
# ==============================================================================
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-lerd}"
MYSQL_DATABASE="${MYSQL_DATABASE:-xelion_online}"
PODMAN_CONTAINER="${PODMAN_CONTAINER:-lerd-mysql}"

run_sql_file() {
    local file_path="$1"
    local desc="$2"
    echo "==> Aplicando $desc: $file_path"

    if command -v mysql >/dev/null 2>&1; then
        mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$file_path"
    elif command -v podman >/dev/null 2>&1 && podman ps --format "{{.Names}}" | grep -q "^${PODMAN_CONTAINER}$"; then
        podman exec -i "$PODMAN_CONTAINER" mysql -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$file_path"
    elif command -v docker >/dev/null 2>&1 && docker ps --format "{{.Names}}" | grep -q "^${PODMAN_CONTAINER}$"; then
        docker exec -i "$PODMAN_CONTAINER" mysql -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$file_path"
    else
        echo "ERROR: No se encontró el cliente mysql local ni el contenedor $PODMAN_CONTAINER en podman/docker." >&2
        exit 1
    fi
}

echo "=== Iniciando aplicación de base de datos en '$MYSQL_DATABASE' ==="

# 1. Esquema base de tablas y vistas
run_sql_file "$DIR/database/schema.sql" "Esquema Base (schema.sql)"

# 2. Stored Procedures Oficiales de Account
run_sql_file "$DIR/database/procedures_account.sql" "Stored Procedures de Account (75 SPs)"

# 3. Stored Procedures Oficiales de AccountLog
run_sql_file "$DIR/database/procedures_accountlog.sql" "Stored Procedures de AccountLog (28 SPs)"

# 4. Stored Procedures Oficiales de StatisticsData
run_sql_file "$DIR/database/procedures_statistics.sql" "Stored Procedures de StatisticsData (18 SPs)"

# 5. Stored Procedures Operacionales de LoginSession
run_sql_file "$DIR/database/procedures_login_session.sql" "Stored Procedures de LoginSession (17 SPs)"

# 6. Semillas Oficiales Idempotentes (FOUS, tUserAuth, tWorldServer, tAccountEventReward)
run_sql_file "$DIR/database/seed_login_world.sql" "Semillas Oficiales (seed_login_world.sql)"

# 7. Esquema Oficial World00_Character y Semillas de Personajes
run_sql_file "$DIR/database/schema_character.sql" "Esquema World00_Character (schema_character.sql)"
run_sql_file "$DIR/database/seed_character_defaults.sql" "Semillas de Personajes (seed_character_defaults.sql)"

echo "=== Aplicación completada exitosamente. Total 143 Stored Procedures cargados. ==="
