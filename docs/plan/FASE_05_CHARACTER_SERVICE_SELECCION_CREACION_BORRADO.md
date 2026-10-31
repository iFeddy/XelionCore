# Fase 05 — Character Service: Selección, Creación, Borrado y Carga de Personajes

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Sección 4 y 5**) y [`docs/Servicios/Character.md`](../Servicios/Character.md).
- **Dependencias previas:** Fase 04 completada sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 05 (`docs/plan/FASE_05_CHARACTER_SERVICE_SELECCION_CREACION_BORRADO.md`) leyendo como especificación obligatoria las Secciones 4 y 5 de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md` y `docs/Servicios/Character.md`.
Implementá en `pkg/server/world/` el 100% de la lógica oficial de gestión de avatares: listado de personajes (PROTO_NC_CHAR_CHARDATA_REQ/ACK), creación atómica con validación de nombre y equipamiento base desde data/ (PROTO_NC_CHAR_CREATE_REQ/ACK), borrado oficial con soft-delete (PROTO_NC_CHAR_DELETE_REQ/ACK) y selección de avatar (PROTO_NC_CHAR_LOGIN_REQ/ACK) emitiendo el token hacia ZoneServer. Creá `pkg/audit/fase05_world_character_dispatcher_test.go` (Control 5) y pruebas E2E en `pkg/server/world/character_e2e_test.go`.
No avances a la Fase 06. Al terminar, actualizá el estado de la Fase 5 a "Completado" en `docs/plan/README.md` y hacé un git commit limpio en main.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Implementar Handlers Oficiales en [`pkg/server/world/character_handler.go`](../../pkg/server/world/):**
   - `handle_PROTO_NC_CHAR_CHARDATA_REQ` (`0x1004`):
     * Consulta avatares del usuario en el mundo mediante `p_Char_GetListOfUserChar`.
     * Carga equipamiento visible e ítems cosméticos de cada slot.
     * Emite `PROTO_NC_CHAR_CHARDATA_ACK` (`0x1005`) con el total de slots y personajes serializados.
   - `handle_PROTO_NC_CHAR_CREATE_REQ` (`0x1007`):
     * Valida unicidad de nombre de avatar mediante `p_Char_CheckName`.
     * Carga stats base y equipamiento inicial desde `data/World/ClassParam/<Class>.json` y `DefaultCharacterData.txt` (cero hardcode).
     * Ejecuta alta atómica en MySQL mediante `p_Char_Create`.
     * Emite `PROTO_NC_CHAR_CREATE_ACK` (`0x1008`).
   - `handle_PROTO_NC_CHAR_DELETE_REQ` (`0x1009`):
     * Valida pertenencia de la cuenta y comprueba que no sea líder de gremio activo.
     * Ejecuta baja lógica en MySQL mediante `p_Char_Delete` (`bDeleted=1, dDeletedDate=NOW()`).
     * Emite `PROTO_NC_CHAR_DELETE_ACK` (`0x100A`).
   - `handle_PROTO_NC_CHAR_LOGIN_REQ` (`0x1001`):
     * Carga datos completos de sesión del avatar mediante `p_Char_GetLoginData`.
     * Genera token criptográfico de transferencia hacia `ZoneService`.
     * Emite `PROTO_NC_CHAR_LOGIN_ACK` (`0x1002`) con coordenadas oficiales (`sLoginZone`, `X`, `Y`).

2. **Crear [`pkg/audit/fase05_world_character_dispatcher_test.go`](../../pkg/audit/fase05_world_character_dispatcher_test.go) (Control 5):**
   - Test automatizado que audita el registro y cobertura de los handlers `fc_NC_*` de gestión de personajes en el Dispatcher.

3. **Crear / Ampliar [`pkg/server/world/character_e2e_test.go`](../../pkg/server/world/):**
   - Pruebas E2E de conexión de cliente, listado de avatares, creación de nuevo personaje, selección y verificación de persistencia en MySQL.

---

## Criterio de Cierre de la Fase 05

1. `go build ./...` -> OK.
2. `go test -v ./pkg/audit -run "^TestFase05_WorldCharacter_Control5"` -> `PASS`.
3. `go test -v ./pkg/server/world/...` -> `PASS`.
4. Marcar Fase 5 como `Completado` en `docs/plan/README.md`.
5. Ejecutar commit limpio:
   `git add -A && git commit -m "feat(world): Fase 05 - gestion oficial de avatares (lista, creacion, borrado y seleccion) 100% MySQL"`
