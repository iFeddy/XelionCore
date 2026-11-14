# FASE 17: Saneamiento Binario 1:1 — In-Game Chunk, Geometría `.shbd`, Movimiento, Chat/GM y Retorno a Selección de Personaje/World

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/Zone.md`, `docs/Servicios/WorldManager.md` y `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 17 completa (`docs/plan/FASE_17_NUCLEO_INGAME_CHUNK_MOVIMIENTO_SHBD_CHAT_GM_Y_LOGOUT_EXE.md`) siguiendo estrictamente `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Antes de modificar cada handler o subsistema, ejecutá `python3 scripts/disasm_pdb_handler.py zone <símbolo>` sobre `../AstraServer/Server/Zone00/Zone.exe` (y `world` sobre `WorldManager.exe`) y verificá cada offset y `sizeof` contra `docs/Servicios/Zone.md`. Prohibido inventar bytes, omitir campos o dejar stubs.
1) Desensamblá en `Zone.exe` la carga y consulta de `.shbd` (`ShineBlockData` / funciones de bloqueo de mapa) y `sp_NC_ACT_MOVEWALK_CMD` / `sp_NC_ACT_MOVERUN_CMD`. Corregí `pkg/data/blockinfo` y `handleMoveWalk`/`handleMoveRun` para que la conversión de coordenadas de mundo a celda/bit `.shbd` sea 100% idéntica a `Zone.exe` y elimine los `PROTO_NC_ACT_MOVEFAIL_CMD` falsos en terreno caminable (ej. `RouCos02` `9647,9002 -> 9588,9200` y `RouN`).
2) Desensamblá y alineá byte a byte contra `Zone.exe` y `Zone.md` todos los paquetes del In-Game Chunk y de los departamentos `NC_MISC` (2), `NC_USER` (3), `NC_CHAR` (4), `NC_MAP` (6), `NC_ACT` (8) y `NC_MOVER` (36): `PROTO_NC_CHAR_CLIENT_BASE_CMD`, `PROTO_NC_CHAR_CLIENT_SHAPE_CMD`, `PROTO_NC_CHAR_CLIENT_QUEST_*`, `PROTO_NC_CHAR_CLIENT_SKILL_CMD`, `PROTO_NC_CHAR_CLIENT_PASSIVE_CMD`, `PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD`, `PROTO_NC_CHAR_CLIENT_GAME_CMD`, `PROTO_NC_CHAR_CLIENT_ITEM_CMD`, `PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD` y `PROTO_NC_MAP_LOGIN_ACK` (verificando HP/SP/MaxHP/MaxSP, stats en `CharParam`, velocidad de movimiento y coordenadas).
3) Desensamblá `sp_NC_ACT_CHAT_REQ` (`0x0045A400`) en `Zone.exe` y corregí `handleChatNormal` para leer exactamente `struct PROTO_NC_ACT_CHAT_REQ { uint8 itemLinkDataCount; uint8 len; char content[]; }` (incluyendo links de ítems si `itemLinkDataCount > 0`), emitir `PROTO_NC_ACT_SOMEONECHAT_CMD` con su estructura oficial exacta y despachar todos los comandos GM (`&` / `/`) gobernados por `AdminLvSet.shn`.
4) Desensamblá `sp_NC_CHAR_LOGOUTREADY_CMD` (`0x1071`), `sp_NC_CHAR_LOGOUTCANCEL_CMD` (`0x1072`) y `sp_NC_USER_NORMALLOGOUT_CMD` (`0x0C18`) en `Zone.exe` + el flujo en `WorldManager.exe`. Implementá el retorno real a Selección de Personaje (`WorldService`) y Selección de Mundo (`LoginService`) emitiendo los paquetes S2C oficiales que espera el cliente y preservando la sesión activa en MySQL para que el cliente reconecte a `WorldService` / `LoginService` sin desconectarse.
5) Creá `pkg/audit/fase17_binary_core_move_chat_logout_test.go`, verificá que `go build ./... && go test -p 1 -count=1 ./...` pase en `exit 0`, recompilá `bin/LoginService`, `bin/WorldService` y `bin/ZoneService`, marcá Fase 17 como "Completado" en `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 17)

1. **Colisiones `.shbd` y Movimiento (`NC_ACT` / `NC_MOVER`):**
   - Desensamblar con `scripts/disasm_pdb_handler.py zone --search Block` / `sp_NC_ACT_MOVERUN_CMD` cómo `Zone.exe` calcula el índice de byte y bit en el archivo `.shbd` a partir de `(x, y)` y del encabezado del `.shbd`.
   - Verificar con un test sobre `data/BlockInfo/RouCos02.shbd` y `RouN.shbd` que las coordenadas de spawn (`9647, 9002` en `RouCos02` y `6888, 7045` en `RouN`) y los desplazamientos adyacentes son transitables (`IsWorldMoveBlocked == false`) y que `handleMoveRun` / `handleMoveWalk` responden `PROTO_NC_ACT_SOMEONEMOVERUN_CMD` / `SOMEONEMOVEWALK_CMD` con el `sizeof` exacto de `Zone.md`.
2. **In-Game Chunk y Stats (`NC_CHAR` / `NC_MAP`):**
   - Comparar campo por campo el buffer generado por cada paquete del In-Game Chunk y `PROTO_NC_MAP_LOGIN_ACK` (`sizeof = 242` bytes payload + `2` bytes OpCode = `244` bytes) contra la definición de `Zone.md` y el llenado de `CharLoginAck` en `Zone.exe`.
   - Asegurar que `MaxHP`, `MaxSP`, `WalkSpeed`, `RunSpeed`, `AC`, `MR`, `WC`, `MA`, `TH`, `TB` en `PROTO_NC_MAP_LOGIN_ACK` se calculen e inyecten en los offsets exactos que lee el cliente `.exe`.
3. **Chat y Comandos GM (`sp_NC_ACT_CHAT_REQ`):**
   - Leer `itemLinkDataCount` (1B), `len` (1B), `itemLinks` (`itemLinkDataCount * sizeof(SHINELINKDATA)`) y `content` (`len` bytes).
   - Verificar que el paquete `0120 00 0b 2661646d696e6c6576656c` (`&adminlevel`) capturado del cliente real sea parseado como `&adminlevel` y ejecutado sin descartarse.
4. **Retorno a Selección de Personaje y Selección de Mundo (`NC_CHAR_LOGOUTREADY_CMD` / `NC_USER_NORMALLOGOUT_CMD`):**
   - Desensamblar `sp_NC_CHAR_LOGOUTREADY_CMD` y `sp_NC_USER_NORMALLOGOUT_CMD` en `Zone.exe` y verificar qué paquetes S2C envía `Zone.exe` / `WorldManager.exe` (ej. `PROTO_NC_CHAR_LOGOUTREADY_ACK` / `PROTO_NC_USER_WORLD_SELECT_ACK` / transferencia de regreso a `WorldService`) y asegurar que `WorldService` re-acepte la sesión del usuario.
5. **Departamentos C2S Cubiertos al 100% en esta Fase:**
   - Todos los handlers de `NC_MISC` (`2`), `NC_USER` (`3`), `NC_CHAR` (`4`), `NC_MAP` (`6`), `NC_ACT` (`8`) y `NC_MOVER` (`36`) auditados contra `Zone.exe` y `Zone.md`.

