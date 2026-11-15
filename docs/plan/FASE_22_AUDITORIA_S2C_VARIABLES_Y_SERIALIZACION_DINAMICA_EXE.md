# FASE 22: Auditoría 1:1 de Paquetes S2C Variables (`[0]`/`[]`), Loops Dinámicos y las 567 Llamadas a `NewPacketWriter`

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/WorldManager/WorldManager.exe` + `pdo/WorldManager.pdb`, `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/WorldManager.md`, `docs/Servicios/Zone.md` y `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 22 (`docs/plan/FASE_22_AUDITORIA_S2C_VARIABLES_Y_SERIALIZACION_DINAMICA_EXE.md`) siguiendo estrictamente `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`, `docs/Servicios/WorldManager.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Antes de modificar cualquier emisor S2C variable o helper, ejecutá `python3 scripts/disasm_pdb_handler.py world <símbolo>` y `python3 scripts/disasm_pdb_handler.py zone <símbolo>` sobre los `.exe` originales y verificá la cabecera base (`sizeof(Header)`), el offset del contador `N` y el tamaño exacto de cada elemento repetido (`sizeof(Entry)`).
1) Auditá los 148 structs dinámicos (`[0];`/`[];`) de `WorldManager.md` y los 155 structs dinámicos de `Zone.md`, junto con las 567 llamadas a `protocol.NewPacketWriter` en `pkg/server/world/` y `pkg/server/zone/` (incluyendo los 55 OpCodes S2C variables y cualquier bloque con `for` o `if/else` que `auditFixedS2CInSource` omitía).
2) Alineá byte a byte contra `Zone.exe` y `WorldManager.exe` todos los emisores S2C variables: In-Game Chunk (`PROTO_NC_CHAR_ITEM_CMD`, `SKILL_CMD`, `PASSIVE_CMD`, `QUEST_DOING_CMD`, `QUEST_DONE_CMD`, `QUEST_READ_CMD`, `ABSTATE_CMD`, `CHARTITLE_CMD`), listas AOI (`PROTO_NC_BRIEFINFO_*`), listas de tiendas/inventarios/almacenes/booth (`NC_ITEM`, `NC_MENU`, `NC_BOOTH`), listas sociales (`NC_PARTY`, `NC_GUILD`, `NC_FRIEND`, `NC_KQ`, `NC_GUILD_ACADEMY`, `NC_HOLY_PROMISE`) y subsistemas (`NC_MINIHOUSE`, `NC_COLLECT`, `NC_GAMBLE`, `NC_DICE`).
3) Creá `pkg/audit/fase22_dynamic_s2c_parity_test.go` con un validador automático que parsee las fórmulas `sizeof(Header) + N * sizeof(Entry)` (y longitudes prefijadas de strings) de `WorldManager.md` y `Zone.md` y verifique en runtime y vía AST que el 100% de los paquetes S2C variables emitidos por `WorldService` y `ZoneService` cumplan la ecuación exacta al byte.
4) Verificá que `go test -p 1 -count=1 ./...` y los 4 scripts `verify_*_100.sh` pasen en `exit 0`, marcá Fase 22 como "Completado" en `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 22)

1. **Catálogo de Fórmulas de Tamaño Variable (`WorldManager.md` + `Zone.md`):**
   - Extracción automática de `HeaderSize`, offset/tipo del campo `count`/`len` y `ElemSize` para los structs con `[0];` / `[];`.
2. **Inspección de las 567 llamadas a `NewPacketWriter` en `pkg/server/world/` y `pkg/server/zone/`:**
   - Cero bloques `NewPacketWriter` ignorados sin validación de tamaño (fijo o variable).
3. **Suite `pkg/audit/fase22_dynamic_s2c_parity_test.go` en `PASS` (`exit 0`).**

