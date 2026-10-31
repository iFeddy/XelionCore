# Fase 04 — Protocolo Oficial de Red: OpCodes y Estructuras PDB 1:1

- **Nivel recomendado:** `Medium`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md), [`docs/Servicios/Character.md`](../Servicios/Character.md) y [`docs/Servicios/WorldManager.md`](../Servicios/WorldManager.md).
- **Dependencias previas:** Fase 03 completada sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 04 (`docs/plan/FASE_04_PROTOCOLO_OPCODES_STRUCTS_WORLD_CHARACTER.md`) leyendo como especificación obligatoria las Secciones de Protocolo de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md`, `docs/Servicios/Character.md` y `docs/Servicios/WorldManager.md`.
Definí el 100% de los OpCodes oficiales de Character y WorldManager con la fórmula (dept << 10) | cmd en `pkg/protocol/opcodes.go`. Creá las estructuras C++ de red con nombres oficiales exactos (PROTO_NC_CHAR_*, PROTO_NC_ITEMDB_*, PROTO_NC_GUILD_*, PROTO_NC_FRIEND_*, PROTO_NC_PARTY_*, PROTO_NC_KQ_*) y sizeof idéntico al byte en `pkg/protocol/pdb/`. Creá `pkg/audit/fase04_world_character_sizeof_test.go` (Control 4) verificando sizeof con binary.Size y ausencia de alias inventados.
No avances a la Fase 05. Al terminar, actualizá el estado de la Fase 4 a "Completado" en `docs/plan/README.md` y hacé un git commit limpio en main.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Definir Constantes Oficiales en [`pkg/protocol/opcodes.go`](../../pkg/protocol/opcodes.go):**
   - Calcular cada OpCode mediante `(department << 10) | command`:
     * Departamento 04 (`NC_CHAR`): `0x1001` (`PROTO_NC_CHAR_LOGIN_REQ`), `0x1002`, `0x1004` (`PROTO_NC_CHAR_CHARDATA_REQ`), `0x1005`, `0x1007` (`PROTO_NC_CHAR_CREATE_REQ`), `0x1008`, `0x1009` (`PROTO_NC_CHAR_DELETE_REQ`), `0x100A`, `0x1036` (`PROTO_NC_CHAR_LOGIN_DB`), `0x1037` (`PROTO_NC_CHAR_LOGOUT_DB`), `0x105D` (`PROTO_NC_CHAR_STAT_INCPOINT_DB_REQ`), etc.
     * Departamento 05 (`NC_AVATAR`): `0x1408` (`PROTO_NC_AVATAR_ERASEDATA_REQ`), `0x1415` (`PROTO_NC_AVATAR_GUILD_MEMBER_DATA_REQ`).
     * Departamento 13 (`NC_GUILD`): `0x3401`..`0x3446`.
     * Departamento 14 (`NC_PARTY`): `0x3801`..`0x3820`.
     * Departamento 16 (`NC_KQ`): `0x4001`..`0x4030`.
     * Departamento 20 (`NC_FRIEND`): `0x5001`..`0x5020`.

2. **Crear / Ampliar Archivos PDB en [`pkg/protocol/pdb/`](../../pkg/protocol/pdb/):**
   - [`proto_nc_char.go`](../../pkg/protocol/pdb/proto_nc_char.go): structs oficiales de selección, datos de avatares, stats, apariencia y persistencia.
   - [`proto_nc_item.go`](../../pkg/protocol/pdb/proto_nc_item.go): structs de manipulación de inventario, equipamiento y opciones mágicas.
   - [`proto_nc_guild.go`](../../pkg/protocol/pdb/proto_nc_guild.go): structs de gremios y academia.
   - [`proto_nc_social.go`](../../pkg/protocol/pdb/proto_nc_social.go): structs de amigos, party y chat.
   - Asegurar alineación `#pragma pack(push, 1)` y paridad exacta en bytes.

3. **Crear [`pkg/audit/fase04_world_character_sizeof_test.go`](../../pkg/audit/fase04_world_character_sizeof_test.go) (Control 4):**
   - Comprueba cada OpCode contra la documentación oficial de `Character.md` y `WorldManager.md`.
   - Comprueba con `binary.Size` que cada struct Go mide exactamente los mismos bytes documentados en C++.
   - Audita que no existan alias no oficiales ni estructuras deprecadas.

---

## Criterio de Cierre de la Fase 04

1. `go build ./...` -> OK.
2. `go test -v ./pkg/audit -run "^TestFase04_WorldCharacter_Control4"` -> `PASS` (100% paridad de OpCodes y `sizeof`).
3. Marcar Fase 4 como `Completado` en `docs/plan/README.md`.
4. Ejecutar commit limpio:
   `git add -A && git commit -m "feat(protocol): Fase 04 - 100% estructuras y opcodes oficiales de Character y WorldManager PDB"`
