# FASE 20: Saneamiento Binario 1:1 — Quests (`NC_QUEST`), Diálogos NPC (`NC_MENU`/`NC_SCRIPT`), Lua, Portales/Instancias, `MiniHouse`, `Pet`, Títulos, Cartas y Casino

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/Zone.md` y `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 20 completa (`docs/plan/FASE_20_QUESTS_DIALOGOS_LUA_INSTANCIAS_MINIHOUSE_PETS_TITULOS_Y_CASINO_EXE.md`) siguiendo estrictamente `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Antes de modificar cada handler o paquete S2C, ejecutá `python3 scripts/disasm_pdb_handler.py zone <símbolo>` sobre `../AstraServer/Server/Zone00/Zone.exe` y verificá cada offset y `sizeof` contra `docs/Servicios/Zone.md`. Prohibido inventar bytes, omitir campos o dejar stubs.
1) Desensamblá y alineá 1:1 contra `Zone.exe` y `Zone.md` todos los handlers y respuestas S2C de `NC_MENU` (11), `NC_QUEST` (17) y `NC_SCRIPT` (28): `sp_NC_ACT_NPCCLICK_CMD`, `sp_NC_MENU_SERVERMENU_ACK`, `PROTO_NC_MENU_SERVERMENU_CMD`, `sp_NC_QUEST_SCRIPT_CMD_ACK`, `sp_NC_QUEST_START_REQ`, `sp_NC_QUEST_GIVEUP_REQ`, `sp_NC_QUEST_READ_REQ`, `sp_NC_QUEST_REWARD_REQ`, `sp_NC_QUEST_JOBDUNGEON_FIND_RNG_REQ`, `sp_NC_QUEST_SCENARIO_DONE_REQ` y `sp_NC_SCRIPT_END_CMD`.
2) Desensamblá en `Zone.exe` la transición de mapas/portales (`Gate` / `LinkTable` / `PROTO_NC_MAP_LINKOTHER_CMD` / `PROTO_NC_MAP_LINKSAME_CMD`) para que cruzar un portal dentro del mismo nodo o hacia otro nodo funcione al 100% con el cliente oficial.
3) Desensamblá y alineá 1:1 contra `Zone.exe` y `Zone.md` todos los handlers de `NC_CT` (27 - Títulos de personaje), `NC_MINIHOUSE` (33 - 23 handlers C2S), `NC_COLLECT` (35 - Cartas de colección), `NC_PET` (37 - Mascotas), `NC_DICE` (40 - 9 handlers C2S) y `NC_GAMBLE` (42 - 18 handlers C2S).
4) Creá `pkg/audit/fase20_binary_quests_menus_portals_subsystems_test.go`, verificá que `go build ./... && go test -p 1 -count=1 ./...` pase en `exit 0`, recompilá `bin/LoginService`, `bin/WorldService` y `bin/ZoneService`, marcá Fase 20 como "Completado" en `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 20)

1. **Diálogos NPC y Menús (`NC_MENU` / `NC_SCRIPT`):**
   - Verificar en `Zone.exe` cómo `sp_NC_ACT_NPCCLICK_CMD` y `sp_NC_MENU_SERVERMENU_ACK` arman `PROTO_NC_MENU_SERVERMENU_CMD` y cómo abren tienda, almacén o diálogo de misión en el cliente.
2. **Ciclo Completo de Quests (`NC_QUEST`):**
   - Verificar el `sizeof` y offsets exactos de `PROTO_NC_QUEST_START_ACK`, `PROTO_NC_QUEST_STEP_INFO_CMD`, `PROTO_NC_QUEST_REWARD_ACK` y los paquetes de misiones del In-Game Chunk (`103A`, `103B`, `10CE`, `10D7`).
3. **Portales (`LINKSAME` / `LINKOTHER`) y Subsistemas (`NC_MINIHOUSE`, `NC_PET`, `NC_CT`, `NC_COLLECT`, `NC_DICE`, `NC_GAMBLE`):**
   - Verificar en `Zone.exe` los bytes exactos de `PROTO_NC_MAP_LINKSAME_CMD` y `PROTO_NC_MAP_LINKOTHER_CMD` al pisar un portal (`Gate*`) y al usar `&linkto` / pergaminos de retorno.

