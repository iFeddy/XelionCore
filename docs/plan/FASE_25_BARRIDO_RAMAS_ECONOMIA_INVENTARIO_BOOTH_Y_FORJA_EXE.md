# FASE 25: Barrido Rama por Rama (`.exe`/`.pdb`) — Economía, Inventario (48 `NC_ITEM`), Comercio, Booth y Forja

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/Zone.md` y `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 25 (`docs/plan/FASE_25_BARRIDO_RAMAS_ECONOMIA_INVENTARIO_BOOTH_Y_FORJA_EXE.md`) siguiendo estrictamente `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Desensamblá cada handler de los 48 `NC_ITEM`, `NC_ITEMDB`, `NC_BOOTH` y `NC_CHARSAVE` en `Zone.exe` con `python3 scripts/disasm_pdb_handler.py zone <símbolo>` e inspeccioná todas las ramas condicionales (`if/else`/`switch`), códigos de error y paquetes S2C emitidos.
1) Auditá rama por rama en `Zone.exe` los 48 handlers C2S de `NC_ITEM` (12), `NC_ITEMDB` (13), `NC_BOOTH` (16), `NC_CHARSAVE` (31) y `SafeTrade`: `EQUIP`/`UNEQUIP`, `RELOC`, `SPLIT`, `DROP`/`PICK`, `BUY`/`SELL`, `STORAGE`, `CHARGEDINVENOPEN`/`CHARGED_WITHDRAW`, `REWARDINVENOPEN`/`REWARDINVENGET`, `USE`, `UPGRADE` (éxito, fallo sin cambio, bajada de nivel, rotura/destrucción, piedras protectoras), `ENCHANT`, `RANDOMOPTION`, `MIX`, `DISMANTLE`, `PRODUCE`, comercio entre jugadores y kioscos `BOOTH`.
2) Alineá en `pkg/server/zone/` todos los códigos de error oficiales, actualizaciones de celdas (`ITEM_CELLCHANGE_CMD`, `ITEM_EQUIPCHANGE_CMD`), descuentos/sumas de Cen y registros `GameLog` en cada rama.
3) Creá `pkg/audit/fase25_branch_item_economy_forge_test.go` probando tanto el camino de éxito como todas las ramas de rechazo (inventario lleno, Cen insuficiente, restricción de clase/nivel, piedra incorrecta, rotura/degradación en upgrade, trade cancelado, booth cerrado).
4) Verificá que `go test -p 1 -count=1 ./...` y los 4 scripts `verify_*_100.sh` pasen en `exit 0`, marcá Fase 25 como "Completado" en `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 25)

1. **Los 48 handlers de `NC_ITEM` verificados rama por rama contra `Zone.exe` (incluyendo todos los resultados de `UPGRADE`/`ENCHANT`/`PRODUCE`).**
2. **Ramas de `SafeTrade`, `NC_BOOTH`, `CHARGED` y `REWARD` verificadas 1:1.**
3. **Suite `pkg/audit/fase25_branch_item_economy_forge_test.go` en `PASS` (`exit 0`).**

