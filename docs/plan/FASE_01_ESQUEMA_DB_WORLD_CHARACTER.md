# Fase 01 — Esquema Oficial 1:1 `World00_Character` en MySQL: Tablas, Columnas, PKs y Semillas

- **Nivel recomendado:** `Medium`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Secciones 0, 1 y 2 completas**) y [docs/TablasLegacy.md](../TablasLegacy.md).
- **Dependencias previas:** Etapa 1 completada sobre `main` (base única `xelion_online`).

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 01 (`docs/plan/FASE_01_ESQUEMA_DB_WORLD_CHARACTER.md`) leyendo como especificación obligatoria las Secciones 0, 1 y 2 completas de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md`, `docs/TablasLegacy.md` y respetando `AGENTS.md`.
Regla Mandatoria: CERO datos hardcodeados y CERO tablas inventadas. Toda la información estática del juego proviene de los JSON y SHN de `data/`, y la persistencia relacional se limita estrictamente a las 61 tablas oficiales de `World00_Character` en `xelion_online`.
Garantizá que las 61 tablas oficiales queden 100% creadas en xelion_online con sus tipos exactos y PKs oficiales, agregá las semillas de DefaultCharacterData.txt y TutorialCharacterData.txt en database/seed_character_defaults.sql, aplicá en Podman lerd-mysql, creá pkg/audit/fase01_world_character_schema_test.go y verificá que compile y pase el Control 1 al 100%.
No avances a la Fase 02. Al terminar, actualizá el estado de la Fase 1 a "Completado" en `docs/plan/README.md` y hacé un git commit limpio en main.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Crear / Actualizar [`database/schema_character.sql`](../../database/schema_character.sql) en `xelion_online`:**
   - Desplegar las **61 tablas oficiales de `World00_Character.bak`** catalogadas en [`docs/TablasLegacy.md`](../TablasLegacy.md) con sus tipos C++ exactos, claves primarias y foráneas:
     * **Dominio de Personaje y Apariencia:** `tCharacter` (con sus 39 columnas oficiales), `tCharacterShape` (1:1 cosmética), `tCharacterState` (buffs persistentes), `tCharacterTitle`, `tRestExp`, `tTutorial`.
     * **Dominio de Inventario y Objetos:** `tItem` (con `nItemKey BIGINT PK`), `tItemOptions` (modificadores mágicos), `tItemCoolTime`, `tChargeItemBuff`.
     * **Dominio de Habilidades:** `tSkill`, `tSkillTree`.
     * **Dominio Social (Gremios y Academia):** `tGuild`, `tGuildMember`, `tGuildAcademy`, `tGuildAcademyMember`.
     * **Dominio de Amigos y Grupos:** `tFriend`, `tFriendGroup`, `tParty`, `tPartyMember`.
     * **Dominio de Misiones y Progresión:** `tQuest`, `tQuestTimes`.
     * **Dominio de Mascotas y Monturas:** `tPet`, `tPetSkill`, `tPetNew`.
     * **Dominio de Bodas, MiniHouses y Tiendas:** `tWedding`, `tMiniHouse`, `tStore`.
     * **Tablas de Soporte y Estadísticas Oficiales:** `tPrison`, `tSpamBlockTime`, `tSpamFilter`, `tStatCharLevel`, `tStatGameMoney`, `tStatisticsGameMoney`, `tUserMoney`, etc.
   - Preservar compatibilidad con el esquema previo de `Account`, `AccountLog` y `StatisticsData` en la base única `xelion_online`.

2. **Crear [`database/seed_character_defaults.sql`](../../database/seed_character_defaults.sql) (idempotente con `INSERT ... ON DUPLICATE KEY UPDATE`):**
   - Extraer y registrar los parámetros iniciales oficiales de `DefaultCharacterData.txt` y `TutorialCharacterData.txt` para las clases jugables:
     * Fighter, Cleric, Archer, Mage, Trickster, Gladiator, etc.
     * Stats base oficiales (STR, CON, DEX, INT, MEN), HP y SP iniciales, dinero base (`nMoney`).
     * Coordenadas de spawn inicial en Roumen / Tutorial Map (`sLoginZone`, `nLoginZoneX`, `nLoginZoneY`, `nLoginZoneD`).

3. **Crear [`scripts/apply_mysql_schema_character.sh`](../../scripts/apply_mysql_schema_character.sh):**
   - Script de aprovisionamiento que ejecuta `schema_character.sql` y `seed_character_defaults.sql` dentro del contenedor Podman `lerd-mysql`.

4. **Crear [`pkg/audit/fase01_world_character_schema_test.go`](../../pkg/audit/fase01_world_character_schema_test.go) (Control 1 de World/Character):**
   - Test automatizado que audita contra `information_schema.TABLES` y `information_schema.COLUMNS` la presencia y tipos de las 61 tablas de `World00_Character` en `xelion_online`.

---

## Criterio de Cierre de la Fase 01

1. `bash scripts/apply_mysql_schema_character.sh` -> Ejecutado con éxito (`Exit 0`).
2. `go build ./...` -> OK.
3. `go test -v ./pkg/audit -run "^TestFase01_WorldCharacter_Control1"` -> `PASS`.
4. Marcar Fase 1 como `Completado` en `docs/plan/README.md`.
5. Ejecutar commit limpio:
   `git add -A && git commit -m "feat(database): Fase 01 - esquema oficial 1:1 World00_Character y semillas de personajes en MySQL"`
