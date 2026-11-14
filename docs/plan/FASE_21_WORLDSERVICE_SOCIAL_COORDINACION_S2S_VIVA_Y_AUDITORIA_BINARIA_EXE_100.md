# FASE 21: Saneamiento Binario 1:1 — `WorldService` (167 C2S), Sistemas Sociales (`Party`, `Guild`, `Friends`, `KQ`, `Wedding`), Puente S2S Vivo y Auditoría Binaria Final 100%

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/WorldManager/WorldManager.exe` + `pdo/WorldManager.pdb`, `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/WorldManager.md`, `docs/Servicios/Zone.md` y `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 21 final (`docs/plan/FASE_21_WORLDSERVICE_SOCIAL_COORDINACION_S2S_VIVA_Y_AUDITORIA_BINARIA_EXE_100.md`) siguiendo estrictamente `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`, `docs/Servicios/WorldManager.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Antes de modificar cada handler o paquete S2C, ejecutá `python3 scripts/disasm_pdb_handler.py world <símbolo>` y `python3 scripts/disasm_pdb_handler.py zone <símbolo>` sobre los `.exe` originales y verificá cada offset y `sizeof` contra `docs/Servicios/WorldManager.md` y `docs/Servicios/Zone.md`. Prohibido inventar bytes, omitir campos o dejar stubs.
1) Desensamblá y alineá 1:1 contra `Zone.exe` y `WorldManager.exe` todos los departamentos sociales y de mundo compartido en `ZoneService` y `WorldService`: `NC_PARTY` (14), `NC_GUILD` (21), `NC_PRISON` (23), `NC_FRIEND` (25), `NC_KQ` (29), `NC_WT` (30), `NC_PROMISE` (34), `NC_WEDDING` (38), `NC_GUILD_ACADEMY` (39) y `NC_HOLY_PROMISE` (45), junto con la revisión completa de los 167 handlers C2S de `WorldService` (`NC_USER`, `NC_CHAR`, `NC_AVATAR`).
2) Asegurá la sincronización en vivo entre `WorldService` y `ZoneService` cuando corren como procesos separados (para que acciones iniciadas en Zone o World —como regreso a Selección de Personaje, Party, Guild, Friends, Whisper, Roar y KQ— surtan efecto inmediato en ambos servicios).
3) Creá `pkg/audit/binary_exe_parity_completeness_test.go` y el script ejecutable `scripts/verify_binary_exe_parity_100.sh` que audite automáticamente que los 245 handlers C2S de `ZoneService` y los 167 handlers C2S de `WorldService` consuman y emitan paquetes con paridad exacta de `sizeof` y estructura respecto a `Zone.md`, `WorldManager.md` y los `.exe` oficiales (`8/8 [PASS]`, `exit 0`), manteniendo también `verify_login_account_100.sh`, `verify_world_character_100.sh` y `verify_zone_gamelog_100.sh` en `exit 0`.
4) Recompilá `bin/LoginService`, `bin/WorldService` y `bin/ZoneService`, marcá Fase 21 como "Completado" en `docs/plan/README.md` y hacé el `git commit` final en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 21)

1. **Paridad 1:1 de `WorldService` (167 C2S) y Departamentos Sociales de `ZoneService`:**
   - Verificar contra `WorldManager.exe` y `Zone.exe` los paquetes de `NC_AVATAR`, `NC_CHAR`, `NC_PARTY`, `NC_GUILD`, `NC_FRIEND`, `NC_KQ`, `NC_GUILD_ACADEMY`, `NC_PROMISE`, `NC_WEDDING`, `NC_HOLY_PROMISE`, `NC_WT` y `NC_PRISON`.
2. **Suite Final `scripts/verify_binary_exe_parity_100.sh` (`8/8 [PASS]`):**
   - **Control 1:** Verificación de herramienta `scripts/disasm_pdb_handler.py` contra `Zone.exe`/`Zone.pdb` y `WorldManager.exe`/`WorldManager.pdb`.
   - **Control 2:** Paridad binaria de In-Game Chunk, `.shbd` (`RouCos02`/`RouN`), Movimiento y Chat (`itemLinkDataCount`).
   - **Control 3:** Paridad binaria de `NC_BRIEFINFO`, `NC_BAT`, `NC_SKILL` y `NC_SOULSTONE`.
   - **Control 4:** Paridad binaria de `NC_ITEM` (`48/48`), `CELLCHANGE`, `EQUIPCHANGE`, `CHARGED_WITHDRAW` y `REWARD`.
   - **Control 5:** Paridad binaria de `NC_MENU`, `NC_QUEST`, `NC_SCRIPT`, Portales, `NC_MINIHOUSE`, `NC_PET`, `NC_CT`, `NC_COLLECT`, `NC_DICE` y `NC_GAMBLE`.
   - **Control 6:** Paridad binaria de `WorldService` (`167/167` C2S) y Departamentos Sociales (`Party`, `Guild`, `Friends`, `KQ`, `Wedding`).
   - **Control 7:** Flujo E2E completo de sesión en vivo incluyendo Login -> World -> Zone -> Movimiento `.shbd` -> Chat/GM -> Retiro Charged Item -> Regreso a Selección de Personaje (`WorldService`).
   - **Control 8:** Escaneo AST y de tamaños S2C emitidos por los handlers contra el catálogo oficial de `Zone.md` y `WorldManager.md`.

