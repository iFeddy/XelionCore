# FASE 26: Barrido Rama por Rama (`.exe`/`.pdb`) — Quests, 274 Bindings C-API Lua + `data/LuaScript/`, Timers de Mundo/Instancias, Subsistemas, Sociales y Suite Final `verify_exe_branch_and_dynamic_s2c_100.sh`

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/WorldManager/WorldManager.exe` + `pdo/WorldManager.pdb`, `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/WorldManager.md`, `docs/Servicios/Zone.md` y `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 26 final (`docs/plan/FASE_26_BARRIDO_RAMAS_QUESTS_SUBSISTEMAS_SOCIAL_Y_SUITE_FINAL_EXE.md`) siguiendo estrictamente `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`, `docs/Servicios/WorldManager.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Desensamblá cada handler, binding C-API de Lua (`c*`) y timer de mundo/instancia con `python3 scripts/disasm_pdb_handler.py world <símbolo>` y `python3 scripts/disasm_pdb_handler.py zone <símbolo>` e inspeccioná todas las ramas condicionales (`if/else`/`switch`), códigos de error y paquetes S2C emitidos.
1) Auditá rama por rama en `Zone.exe` y `WorldManager.exe` todos los handlers de `NC_MENU` (11), `NC_QUEST` (17), `NC_SCRIPT` (28), `NC_MINIHOUSE` (33), `NC_PET` (37), `NC_DICE` (40), `NC_GAMBLE` (42) y los departamentos sociales (`NC_PARTY` 14, `NC_GUILD` 21, `NC_PRISON` 23, `NC_FRIEND` 25, `NC_KQ` 29, `NC_WT` 30, `NC_PROMISE` 34, `NC_WEDDING` 38, `NC_GUILD_ACADEMY` 39, `NC_HOLY_PROMISE` 45).
2) Desensamblá y verificá contra `Zone.exe` y `WorldManager.exe` la ejecución de los scripts `.lua` reales de `data/LuaScript/` sobre los 274 bindings C-API (`pkg/script/lua/`) y los timers autónomos (`OnTimer`) de `KingdomQuest` (`KQ`), `GuildTournament`, `GuildWar`, `Wedding` e `InstanceDungeon` (transiciones de fase, spawns por script, puertas, recompensas y expulsión al finalizar).
3) Alineá en `pkg/server/world/`, `pkg/server/zone/` y `pkg/script/lua/` todas las condiciones de rechazo, códigos de error oficiales, efectos de bindings Lua y secuencias de paquetes S2C / eventos S2S para cada rama.
4) Creá `pkg/audit/exe_branch_and_dynamic_s2c_completeness_test.go` y el script ejecutable `scripts/verify_exe_branch_and_dynamic_s2c_100.sh` (`8/8 [PASS]`, `exit 0`) que consolide la verificación de paquetes S2C variables (Fase 22), el barrido de ramas de éxito/rechazo de los 412 handlers C2S (`167` World + `245` Zone), las fórmulas/timers de zona (Fase 24) y la ejecución de `data/LuaScript/` + timers de instancias/KQ (Fase 26), manteniendo las 4 suites `verify_*_100.sh` anteriores en `exit 0`.
5) Recompilá `bin/LoginService`, `bin/WorldService` y `bin/ZoneService`, marcá Fase 26 como "Completado" en `docs/plan/README.md` y hacé el `git commit` final en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 26)

1. **Ramas de éxito y rechazo de `NC_MENU`, `NC_QUEST`, `NC_SCRIPT`, `NC_MINIHOUSE`, `NC_PET`, `NC_DICE`, `NC_GAMBLE` y departamentos sociales verificadas 1:1 contra `Zone.exe` y `WorldManager.exe`.**
2. **Ejecución real de `data/LuaScript/` sobre los 274 bindings C-API y timers autónomos de `KQ`, `GuildTournament`, `GuildWar`, `Wedding` e `InstanceDungeon` verificados 1:1.**
3. **Suite Final `scripts/verify_exe_branch_and_dynamic_s2c_100.sh` (`8/8 [PASS]`, `exit 0`):**
   - **Control 1:** Validación 100% de fórmulas S2C de tamaño variable (`sizeof(Header) + N * sizeof(Entry)`) y de las 567 llamadas a `NewPacketWriter`.
   - **Control 2:** Ramas de éxito y rechazo de Entrada, Mundo, Movimiento `.shbd`, AOI y Chat (`Fase 23`).
   - **Control 3:** Ramas de éxito y rechazo de Combate, Fórmulas Matemáticas, Skills, AbState, SoulStone, Mobs, Drops y Timers de Zona (`Fase 24`).
   - **Control 4:** Ramas de éxito y rechazo de Economía, 48 `NC_ITEM`, SafeTrade, Booth y Forja (`Fase 25`).
   - **Control 5:** Ramas de éxito y rechazo de Quests, ejecución de `data/LuaScript/` sobre los 274 bindings C-API, MiniHouse, Pets y Casino (`Fase 26A`).
   - **Control 6:** Ramas de éxito y rechazo de Departamentos Sociales, Timers de `KQ`/`GuildTournament`/`GuildWar`/`Wedding` y Coordinación S2S Viva (`Fase 26B`).
   - **Control 7:** Simulación E2E multi-rama (casos felices + errores forzados + ticks de timers + scripts Lua) en clúster real (`LoginService` + `WorldService` + `ZoneService`).
   - **Control 8:** Cobertura combinada de las 5 suites de verificación (`verify_login_account_100.sh`, `verify_world_character_100.sh`, `verify_zone_gamelog_100.sh`, `verify_binary_exe_parity_100.sh` y `verify_exe_branch_and_dynamic_s2c_100.sh`) en `exit 0`.

