# Fase 03 — Repositorios Go 100% `CALL` Stored Procedures (Eliminación de SQL Inline)

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Sección 3 y 4**) y [`pkg/database/mysql_character.go`](../../pkg/database/mysql_character.go).
- **Dependencias previas:** Fase 02 completada sobre `main` (311 Stored Procedures aplicados en `xelion_online`).

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 03 (`docs/plan/FASE_03_REPOSITORIOS_GO_WORLD_CHARACTER_MYSQL.md`) leyendo como especificación obligatoria las Secciones 3 y 4 de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md` y respetando `AGENTS.md`.
Regla Mandatoria: CERO datos hardcodeados y CERO SQL inline. Los repositorios de Go interactúan con la base de datos exclusivamente mediante `CALL <StoredProcedure>` contra las 61 tablas oficiales en `xelion_online`. La parametrización estática proviene de `data/`.
Eliminá el 100% del SQL inline (SELECT, INSERT, UPDATE, DELETE) en `pkg/database/mysql_character.go` y subsistemas asociados, reemplazándolo exclusivamente por llamadas `CALL <SP>` contra MySQL xelion_online. Definí interfaces fuertemente tipadas en `pkg/database/character.go`, adaptá todos los llamadores en el servidor de mundo, creá `pkg/audit/fase03_world_character_ast_test.go` (Control 3 AST) y verificá que compile y pase al 100%.
No avances a la Fase 04. Al terminar, actualizá el estado de la Fase 3 a "Completado" en `docs/plan/README.md` y hacé un git commit limpio en main.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Definir Interfaces Oficiales en [`pkg/database/character.go`](../../pkg/database/character.go):**
   - `CharacterRepository`: `GetCharactersByAccount`, `GetCharacterByID`, `GetCharacterLoginData`, `CreateCharacter`, `DeleteCharacter`, `CheckNameExists`, `SaveLogoutData`, `IncreaseStatPoint`.
   - `ItemRepository`: `GetItemList`, `InsertItem`, `UpdateItemSlot`, `UpdateItemCount`, `DeleteItem`, `SaveItemOptions`.
   - `SkillRepository`: `GetSkillList`, `InsertSkill`, `UpdateSkill`, `SaveSkillTree`.
   - `GuildRepository`: `CreateGuild`, `GetGuildInfo`, `AddMember`, `RemoveMember`, `ChangeMemberGrade`, `AddAcademyMember`.
   - `SocialRepository`: `GetFriendList`, `AddFriend`, `DeleteFriend`, `CreateParty`, `AddPartyMember`, `RemovePartyMember`.

2. **Refactorizar [`pkg/database/mysql_character.go`](../../pkg/database/mysql_character.go):**
   - Eliminar todas las sentencias SQL embebidas (`SELECT ... FROM tCharacter`, `INSERT INTO tItem ...`, etc.).
   - Migrar cada operación a invocaciones oficiales `CALL <SP>(...)`.
   - Asegurar manejo correcto de transacciones, escaneo de filas y nulos mediante `sql.NullString` y `sql.NullTime`.

3. **Crear [`pkg/audit/fase03_world_character_ast_test.go`](../../pkg/audit/fase03_world_character_ast_test.go) (Control 3 AST):**
   - Analizador de árbol sintáctico (AST) que recorre los archivos `.go` de repositorios de Character verificando que no existan queries SQL inline (`SELECT`, `INSERT`, `UPDATE`, `DELETE`).
   - Verifica que el 100% de las consultas comiencen con `"CALL "`.

---

## Criterio de Cierre de la Fase 03

1. `go build ./...` -> OK.
2. `go test -v ./pkg/audit -run "^TestFase03_WorldCharacter_AST"` -> `PASS` (Cero SQL inline).
3. `go test -v ./pkg/database/...` -> `PASS`.
4. Marcar Fase 3 como `Completado` en `docs/plan/README.md`.
5. Ejecutar commit limpio:
   `git add -A && git commit -m "refactor(database): Fase 03 - migracion 100% a CALL Stored Procedures en repositorios de Character"`
