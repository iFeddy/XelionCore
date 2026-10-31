# Fase 02 — Stored Procedures Oficiales de `World00_Character` (311 SPs)

- **Nivel recomendado:** `Medium` *(o `High`)*
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Sección 3 completa**) y [docs/sql/procedures/World00_Character_all_procedures.sql](../sql/procedures/World00_Character_all_procedures.sql).
- **Dependencias previas:** Fase 01 completada sobre `main` (esquema de 61 tablas aplicado en `xelion_online`).

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 02 (`docs/plan/FASE_02_STORED_PROCEDURES_CHARACTER_PDB.md`) leyendo como especificación obligatoria la Sección 3 de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md`, `docs/ProcedimientosLegacy.md` y `docs/sql/procedures/World00_Character_all_procedures.sql`.
Recreá los 311 Stored Procedures oficiales de World00_Character en `database/procedures_character.sql`, aplicalos en Podman lerd-mysql, creá pkg/audit/fase02_world_character_sp_test.go (Control 2 de World/Character) validando que COUNT(*) == 311 en information_schema.ROUTINES y agregá tests funcionales Go para los flujos críticos.
No avances a la Fase 03. Al terminar, actualizá el estado de la Fase 2 a "Completado" en `docs/plan/README.md` y hacé un git commit limpio en main.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Crear [`database/procedures_character.sql`](../../database/procedures_character.sql):**
   - Implementar los **311 Stored Procedures oficiales** de `Character.exe` y `WorldManager.exe` adaptados a MySQL 8.0:
     * **Ciclo de Vida de Avatar:** `p_Char_Create`, `p_Char_Delete`, `p_Char_GetListOfUserChar`, `p_Char_GetLoginData`, `p_Char_Login`, `p_Char_Logout`, `p_Char_CheckName`, `p_Char_SetExp`, `p_Char_SetMoney`, `p_Char_SetLevel`, `p_Char_IncStatPoint`, `p_Char_UpdateLocation`.
     * **Inventario y Objetos:** `p_Item_GetList`, `p_Item_Insert`, `p_Item_UpdateSlot`, `p_Item_UpdateCount`, `p_Item_Delete`, `p_Item_SaveOptions`, `p_Item_Equip`, `p_Item_Unequip`, `p_Item_Upgrade`.
     * **Habilidades:** `p_Skill_GetList`, `p_Skill_Insert`, `p_Skill_Update`, `p_Skill_Delete`, `p_SkillTree_Set`.
     * **Gremios y Academia:** `p_Guild_Create`, `p_Guild_Delete`, `p_Guild_AddMember`, `p_Guild_DelMember`, `p_Guild_ChangeGrade`, `p_Guild_Academy_AddMember`, `p_Guild_Academy_Graduate`.
     * **Amigos y Grupos:** `p_Friend_GetList`, `p_Friend_Add`, `p_Friend_Delete`, `p_FriendGroup_Add`, `p_Party_Create`, `p_Party_AddMember`, `p_Party_DelMember`.
     * **Misiones y Mascotas:** `p_Quest_GetList`, `p_Quest_Insert`, `p_Quest_Update`, `p_Quest_Delete`, `p_Pet_GetList`, `p_Pet_Insert`, `p_Pet_Update`.

2. **Crear [`scripts/apply_mysql_procedures_character.sh`](../../scripts/apply_mysql_procedures_character.sh):**
   - Script bash de carga atómica de los procedimientos en Podman `lerd-mysql`.

3. **Crear [`pkg/audit/fase02_world_character_sp_test.go`](../../pkg/audit/fase02_world_character_sp_test.go) (Control 2):**
   - Catálogo exhaustivo que valida la presencia de los 311 procedimientos en `information_schema.ROUTINES`.
   - Tests funcionales unitarios ejecutando `p_Char_Create`, `p_Char_GetListOfUserChar`, `p_Item_Insert` y `p_Char_Delete`.

---

## Criterio de Cierre de la Fase 02

1. `bash scripts/apply_mysql_procedures_character.sh` -> OK.
2. `go build ./...` -> OK.
3. `go test -v ./pkg/audit -run "^TestFase02_WorldCharacter_Control2"` -> `PASS` (311 SPs verificados).
4. Marcar Fase 2 como `Completado` en `docs/plan/README.md`.
5. Ejecutar commit limpio:
   `git add -A && git commit -m "feat(database): Fase 02 - implementacion de los 311 Stored Procedures de World00_Character"`
