# FASE 03: Repositorios Go 100% `CALL` Stored Procedures (Eliminación de SQL Inline)

## 1. Objetivo Técnico
Refactorizar íntegramente `pkg/database/mysql_character.go` y los repositorios relacionados con el mundo (`pkg/database/character.go`), erradicando cualquier consulta SQL inline (`SELECT`, `INSERT`, `UPDATE`, `DELETE`) y reemplazándola exclusivamente por invocaciones `CALL <StoredProcedure>(...)` contra MySQL (`xelion_online`), con interfaces fuertemente tipadas y sin aserciones de tipo inseguras.

---

## 2. Archivos a Modificar / Crear
- `pkg/database/character.go`: Definición de interfaces oficiales:
  - `CharacterRepository`: Listado de avatares, creación, borrado, carga integral y persistencia de logout.
  - `ItemRepository`: Inventario, almacén, equipamiento, opciones mágicas y transferencias atómicas.
  - `SkillRepository`: Repertorio de skills y árboles de pasivas.
  - `GuildRepository`: Gestión completa de gremios, miembros y academia.
  - `SocialRepository`: Lista de amigos, grupos, bloqueo de spammers y bodas.
- `pkg/database/mysql_character.go`: Implementación de los métodos de las interfaces invocando exclusivamente `r.db.QueryRow("CALL ...")` y `r.db.Query("CALL ...")`.
- `pkg/audit/fase03_world_character_ast_test.go`: Test automatizado del **Control 3** mediante análisis de árbol sintáctico (AST) verificando cero SQL inline en `pkg/database/mysql_character*.go`.

---

## 3. Puntos Críticos de Refactorización en Go
1. **`GetCharactersByAccount`:** Reemplazar query multicampo de `tCharacter` y `tCharacterShape` por `CALL p_Char_GetListOfUserChar(?, ?)`.
2. **`CreateCharacter`:** Reemplazar inserciones directas por `CALL p_Char_Create(?, ?, ?, ?, ?, ?, ?, ?, ?, @out_charNo, @out_ret)`.
3. **`DeleteCharacter`:** Reemplazar `UPDATE tCharacter SET bDeleted=1` por `CALL p_Char_Delete(?, ?, @ret)`.
4. **`GetCharacterByID` / `GetCharacterLoginData`:** Reemplazar joins inline por `CALL p_Char_GetLoginData(?, ?)`.
5. **Inventario:** Reemplazar consultas inline de `tItem` por `CALL p_Item_GetList(?, ?)` y `CALL p_Item_Insert(...)`.

---

## 4. Criterios de Aceptación y Cierre
1. El analizador AST en `pkg/audit/fase03_world_character_ast_test.go` confirma que no existe ninguna ocurrencia de SQL inline en los repositorios de personajes.
2. `go test -v ./pkg/database/...` pasa en verde contra MySQL `xelion_online`.
3. `go build ./...` compila limpiamente.
4. `git commit` en `main`:
   `git add -A && git commit -m "refactor(database): Fase 03 - migracion 100% a CALL Stored Procedures en repositorios de Character"`
