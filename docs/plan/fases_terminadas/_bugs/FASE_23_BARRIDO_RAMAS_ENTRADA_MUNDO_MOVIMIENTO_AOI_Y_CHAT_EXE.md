# FASE 23: Barrido Rama por Rama (`.exe`/`.pdb`) — Entrada, Mundo, Movimiento `.shbd`, AOI y Chat

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/WorldManager/WorldManager.exe` + `pdo/WorldManager.pdb`, `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/WorldManager.md`, `docs/Servicios/Zone.md` y `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 23 (`docs/plan/FASE_23_BARRIDO_RAMAS_ENTRADA_MUNDO_MOVIMIENTO_AOI_Y_CHAT_EXE.md`) siguiendo estrictamente `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`, `docs/Servicios/WorldManager.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Desensamblá cada handler y subfunción con `python3 scripts/disasm_pdb_handler.py world <símbolo>` y `python3 scripts/disasm_pdb_handler.py zone <símbolo>` e inspeccioná todos los saltos condicionales (`je`/`jne`/`jb`/`ja`/`switch`), códigos de error oficiales y secuencias de paquetes S2C emitidos en cada rama.
1) Auditá rama por rama en `WorldManager.exe` y `Zone.exe` los departamentos `NC_MISC` (2), `NC_USER` (3), `NC_CHAR` (4), `NC_AVATAR` (5), `NC_MAP` (6), `NC_BRIEFINFO` (7), `NC_ACT` (8) y `NC_MOVER` (36): login a World, creación/borrado/renombrado/selección de avatar, transferencia a Zone (`MAP_LOGIN_REQ`, `CAN_WARP_ACK`, portales `LINKSAME`/`LINKOTHER`), colisiones `.shbd`/`.sbi` (`MOVEWALK`, `MOVERUN`, `STOP`, `JUMP`, `MOVEFAIL_CMD`), interacción con NPCs (`sp_NC_ACT_NPCCLICK_CMD @ 0x004C4B70` -> `0x004C4660` -> `QuestPlayer_NPCClickCheck` / `PROTO_NC_ACT_NPCMENUOPEN_REQ` -> `sp_NC_ACT_NPCMENUOPEN_ACK @ 0x004C4910` -> `cm_MenuAction` -> `NPCRoleBase` / `PROTO_NC_MENU_SERVERMENU_CMD` / tienda / almacén), monturas `MOVER`, visibilidad `BRIEFINFO` y chat/GM (`AdminLvSet.shn`).
2) Alineá en `pkg/server/world/` y `pkg/server/zone/` todas las ramas de validación, códigos de error y orden de paquetes S2C para que coincidan 1:1 con el ensamblador de `WorldManager.exe` y `Zone.exe`.
3) Creá `pkg/audit/fase23_branch_entry_move_aoi_test.go` probando tanto el camino de éxito como todas las ramas de rechazo/error (token inválido, nombre duplicado, slot ocupado, mapa inválido, movimiento bloqueado en `.shbd`, nivel GM insuficiente, link bloqueado).
4) Verificá que `go test -p 1 -count=1 ./...` y los 4 scripts `verify_*_100.sh` pasen en `exit 0`, marcá Fase 23 como "Completado" en `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 23)

1. **Ramas de éxito y error de `NC_USER`, `NC_CHAR` y `NC_AVATAR` en `WorldService` alineadas 1:1 con `WorldManager.exe`.**
2. **Ramas de `NC_MAP`, `NC_ACT`, `NC_MOVER` y `NC_BRIEFINFO` en `ZoneService` alineadas 1:1 con `Zone.exe`.**
3. **Suite `pkg/audit/fase23_branch_entry_move_aoi_test.go` en `PASS` (`exit 0`).**

