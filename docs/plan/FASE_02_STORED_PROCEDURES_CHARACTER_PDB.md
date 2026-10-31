# FASE 02: Stored Procedures Oficiales de `World00_Character` (311 SPs)

## 1. Objetivo Técnico
Recrear e implementar en MySQL (`xelion_online`) el catálogo completo de los **311 Stored Procedures** oficiales de `Character.exe` y `WorldManager.exe` documentados en [docs/sql/procedures/World00_Character_all_procedures.sql](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/World00_Character_all_procedures.sql) y [docs/Servicios/Character.md](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/Servicios/Character.md), garantizando firmas de parámetros idénticas a las llamadas ODBC oficiales.

---

## 2. Archivos a Modificar / Crear
- `database/procedures_character.sql`: Código SQL de los 311 Stored Procedures oficiales adaptados a la sintaxis y tipos de MySQL 8.0 preservando semántica transaccional.
- `scripts/apply_mysql_procedures_character.sh`: Script bash para compilación y carga de procedimientos en MySQL.
- `pkg/audit/fase02_world_character_sp_test.go`: Test automatizado del **Control 2** verificando la existencia de los 311 procedures en `information_schema.ROUTINES`.

---

## 3. Familias de Stored Procedures Críticos
1. **Ciclo de Vida de Avatar:** `p_Char_Create`, `p_Char_Delete`, `p_Char_GetListOfUserChar`, `p_Char_GetLoginData`, `p_Char_Login`, `p_Char_Logout`, `p_Char_CheckName`.
2. **Inventario Transaccional:** `p_Item_GetList`, `p_Item_Insert`, `p_Item_UpdateSlot`, `p_Item_UpdateCount`, `p_Item_Delete`, `p_Item_SaveOptions`.
3. **Estadísticas y Nivel:** `p_Char_IncStatPoint`, `p_Char_SetExp`, `p_Char_SetMoney`, `p_Char_SetLevel`, `p_Char_UpdateLocation`.
4. **Habilidades y Talentos:** `p_Skill_GetList`, `p_Skill_Insert`, `p_Skill_Update`, `p_SkillTree_Set`.
5. **Gremios y Red Social:** `p_Guild_Create`, `p_Guild_AddMember`, `p_Guild_ChangeGrade`, `p_Guild_Delete`, `p_Friend_GetList`, `p_Friend_Add`, `p_Friend_Delete`.
6. **Misiones y Recompensas:** `p_Quest_GetList`, `p_Quest_Insert`, `p_Quest_Update`, `p_Quest_Delete`.

---

## 4. Criterios de Aceptación y Cierre
1. `bash scripts/apply_mysql_procedures_character.sh` compila los 311 procedures sin warnings ni errores de sintaxis en `lerd-mysql`.
2. `go test -v ./pkg/audit -run "^TestFase02_WorldCharacter_Control2"` pasa 100% en verde validando que `COUNT(*) == 311` en `information_schema.ROUTINES`.
3. Tests funcionales en Go validan la ejecución atómica de creación, carga y borrado de avatar mediante `CALL`.
4. `git commit` en `main`:
   `git add -A && git commit -m "feat(database): Fase 02 - implementacion de los 311 Stored Procedures de World00_Character"`
