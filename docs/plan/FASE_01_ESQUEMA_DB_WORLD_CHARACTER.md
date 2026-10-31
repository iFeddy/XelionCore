# FASE 01: Recreación del Esquema Oficial 1:1 de `World00_Character` en MySQL

## 1. Objetivo Técnico
Garantizar la presencia de las **61 tablas oficiales** de `World00_Character.bak` en la base de datos MySQL (`xelion_online`), con sus tipos de datos, restricciones PK, índices y claves foráneas exactas según [docs/TablasLegacy.md](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/TablasLegacy.md), junto con la carga de semillas de inicialización de personajes (`DefaultCharacterData.txt` y `TutorialCharacterData.txt`).

---

## 2. Archivos a Modificar / Crear
- `database/schema_character.sql`: Definición DDL idempotente (`CREATE TABLE IF NOT EXISTS`) de las 61 tablas de `World00_Character`.
- `database/seed_character_defaults.sql`: Semillas de creación de personajes por clase (stats base, equipamiento inicial, monedas y coordenadas de spawn) extraídas de `DefaultCharacterData.txt` y `TutorialCharacterData.txt`.
- `scripts/apply_mysql_schema_character.sh`: Script de aprovisionamiento en Podman (`podman exec -i lerd-mysql mysql ...`).
- `pkg/audit/fase01_world_character_schema_test.go`: Test automatizado del **Control 1** consultando `information_schema.TABLES` y `information_schema.COLUMNS`.

---

## 3. Tablas Oficiales Clave a Desplegar
1. **Personajes y Cosméticos:** `tCharacter`, `tCharacterShape`, `tCharacterState`, `tCharacterTitle`, `tRestExp`, `tTutorial`.
2. **Inventario y Objetos:** `tItem`, `tItemOptions`, `tItemCoolTime`, `tChargeItemBuff`.
3. **Habilidades:** `tSkill`, `tSkillTree`.
4. **Gremios y Academia:** `tGuild`, `tGuildMember`, `tGuildAcademy`, `tGuildAcademyMember`.
5. **Amigos y Grupos:** `tFriend`, `tFriendGroup`, `tParty`, `tPartyMember`.
6. **Misiones:** `tQuest`, `tQuestTimes`.
7. **Mascotas:** `tPet`, `tPetSkill`.
8. **Social / Vivienda:** `tWedding`, `tMiniHouse`, `tStore`.

---

## 4. Criterios de Aceptación y Cierre
1. `bash scripts/apply_mysql_schema_character.sh` aplica sin errores sobre el contenedor Podman `lerd-mysql`.
2. `go test -v ./pkg/audit -run "^TestFase01_WorldCharacter_Control1"` pasa 100% en verde verificando las 61 tablas y sus columnas principales.
3. `go build ./...` compila limpiamente.
4. `git commit` en `main`:
   `git add -A && git commit -m "feat(database): Fase 01 - esquema oficial 1:1 World00_Character y semillas de personajes en MySQL"`
