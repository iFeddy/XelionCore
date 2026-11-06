#!/bin/bash
# ==============================================================================
# XelionCore - Script de Aplicación del Esquema y Semillas de World00_Character
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

echo "=== Aplicando Esquema y Semillas de World00_Character en '$MYSQL_DATABASE' ==="

# 1. Esquema oficial 1:1 de World00_Character (91 tablas oficiales)
run_sql_file "$DIR/database/schema_character.sql" "Esquema World00_Character (schema_character.sql)"

# 2. Semillas oficiales de personajes, opciones y exp
run_sql_file "$DIR/database/seed_character_defaults.sql" "Semillas de Personajes (seed_character_defaults.sql)"

echo "=== Esquema y semillas de World00_Character aplicados exitosamente. ==="
