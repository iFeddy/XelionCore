# Plan Maestro Etapa 4: Saneamiento Binario 1:1 de Handlers contra `.exe` + `.pdb` (`ZoneService` & `WorldService`)

## 0. Causa Raíz y Reglas Inquebrantables de la Etapa 4

Durante las Etapas 1–3 se completaron:
- Las **176 tablas** y **573 Stored Procedures** en MySQL (`xelion_online`).
- La carga en RAM de los **180 `.shn`**, **301 `BlockInfo` (`.shbd`/`.sbi`/`.aid`/`.bdt`/`.shab`)**, **52 `AreaBMP`** y **274 APIs Lua**.
- El catálogo de OpCodes (`pkg/protocol/opcodes.go`) y structs PDB (`pkg/protocol/pdb/`).

Sin embargo, al conectar el **cliente oficial `.exe` de Fiesta Online** con `DEBUG=true`, quedó expuesta una brecha crítica en la capa de **handlers de `pkg/server/zone/` y `pkg/server/world/`**:
1. **Serialización/Deserialización manual desalineada del PDB y del `.exe`:** Muchos handlers leían o escribían bytes a mano mediante `PacketReader`/`PacketWriter` omitiendo campos o usando tamaños distintos a los de `docs/Servicios/Zone.md` y `docs/Servicios/WorldManager.md` (por ejemplo, `handleChatNormal` no leía `itemLinkDataCount` de `PROTO_NC_ACT_CHAT_REQ`).
2. **Flujos de negocio simplificados sin desensamblar el `.exe`:** Handlers como `sp_NC_USER_NORMALLOGOUT_CMD` cerraban la conexión en lugar de ejecutar el flujo real de `Zone.exe` / `WorldManager.exe` (regreso a Selección de Personaje o Selección de Mundo), y `sp_NC_ITEM_CHARGED_WITHDRAW_REQ` no enviaba la actualización de celda al inventario del cliente.
3. **Geometría `.shbd` (`BlockInfo`) rechazando movimiento válido (`MOVEFAIL`):** `IsWorldMoveBlocked` / `handleMoveRun` interpretaba las coordenadas/bits de `.shbd` de forma incompatible con `Zone.exe`, devolviendo `PROTO_NC_ACT_MOVEFAIL_CMD` en pleno terreno transitable.
4. **Tests de auditoría acoplados al código Go y no al `.exe`:** Los tests anteriores enviaban paquetes armados según lo que el handler Go esperaba leer, dando `PASS` aunque el handler difiriera del binario oficial `Zone.exe`.

---

### Reglas de Oro Mandatorias para toda la Etapa 4 (Fases 17 a 21)

1. **Desensamblado Obligatorio con `scripts/disasm_pdb_handler.py` antes de codificar:**
   - Está **terminantemente prohibido** adivinar qué lee, valida o responde un handler.
   - Para cada handler `sp_NC_*` (`Zone.exe`) o `fc_NC_*` (`WorldManager.exe`), el agente **DEBE** ejecutar:
     ```bash
     python3 scripts/disasm_pdb_handler.py zone <sp_NC_NOMBRE_HANDLER> 600
     ```
     (y desensamblar las subfunciones `SendToClient*` / helpers que invoque) junto con la consulta de la struct exacta en `docs/Servicios/Zone.md` y `docs/Servicios/WorldManager.md`.
2. **Paridad de Bytes 1:1 en `PacketReader` y `PacketWriter` contra `Zone.md` / `WorldManager.md`:**
   - Todo paquete C2S leído por un handler debe consumir **exactamente** los campos y tipos declarados en `struct PROTO_NC_*` de `docs/Servicios/Zone.md` / `WorldManager.md` en sus offsets exactos.
   - Todo paquete S2C emitido por un handler debe producir **exactamente** `2 + sizeof(PROTO_NC_*)` bytes (2 bytes de OpCode + el tamaño exacto documentado en `Zone.md` / `WorldManager.md`, o `2 + sizeof(header) + N * sizeof(entry)` para paquetes con array variable `content[0]` / `list[0]`).
3. **Cero Valores Inventados, Cero Hardcodes y Cero Stubs Silenciosos:**
   - Ningún handler puede devolver constantes mágicas inventadas ni ignorar silenciosamente acciones que en `Zone.exe` actualizan estado y envían `ACK`/`CMD` al cliente.
4. **Tests de Paridad Binaria por Handler (Verificación de Wire Size y Offsets):**
   - Cada fase incluye un test en `pkg/audit/` que invoca los handlers reales y verifica que el buffer S2C emitido (`pw.Data()`) tenga **exactamente el `sizeof` oficial de `docs/Servicios/Zone.md`** y que el parser C2S consuma un payload construido estrictamente desde la struct oficial de `pkg/protocol/`.

---

## 1. Partición Exhaustiva del 100% de los Handlers (Fases 17 a 21)

Para cubrir **los 245 handlers C2S de `Zone.exe` y los 167 handlers C2S de `WorldManager.exe` (412 handlers C2S en total + todos sus paquetes S2C)** sin omitir ni uno solo y sin agotar la ventana de contexto, la Etapa 4 se divide en 5 fases por departamentos completos del protocolo:

| Fase | Departamentos Cubiertos (100% de sus OpCodes C2S y S2C) | Alcance Técnico Principal |
| :---: | :--- | :--- |
| **Fase 17** | `NC_MISC` (2), `NC_USER` (3), `NC_CHAR` (4), `NC_MAP` (6), `NC_ACT` (8), `NC_MOVER` (36) + Geometría `.shbd` | **Entrada al Juego y Núcleo Vital:** In-Game Chunk (`CHAR_CLIENT_BASE`, `SHAPE`, `ITEM`, `QUEST_*`, `SKILL`, `PASSIVE`, `GAME`, `CHARTITLE`, `MAP_LOGIN_ACK`), colisiones `.shbd` (`BlockInfo` / `ShineBlockData` en `Zone.exe`), movimiento (`MOVEWALK`, `MOVERUN`, `STOP`, `JUMP`, `MOVER`), Chat (`itemLinkDataCount`), comandos GM (`AdminLvSet.shn`) y flujo real de Logout / Volver a Selección de Personaje / Selección de Mundo. |
| **Fase 18** | `NC_BRIEFINFO` (7), `NC_BAT` (9), `NC_SKILL` (18), `NC_SOULSTONE` (22) + Spawns/IA de Monstruos | **Visibilidad AOI y Combate 1:1:** Structs S2C exactos de `BRIEFINFO_LOGINCHARACTER`, `CHARACTER`, `REGENMOB`, `MOB`, `DROPEDITEM`, `ITEMONFIELD`, `MAGICFIELDINFO`, `MOVER`, `DOOR`; target (`BAT_TARGETTING_REQ` -> `BAT_TARGETINFO_CMD`), auto-ataque, skills (`SKILL_CAST_REQ`, `SKILL_STOP_REQ`), `AbState`, `SoulStone`, HP/SP change, muerte/revivir y drops. |
| **Fase 19** | `NC_ITEM` (12), `NC_ITEMDB` (13), `NC_BOOTH` (16), `NC_CHARSAVE` (31) | **Economía, Inventario y Forja 1:1 (48 C2S `NC_ITEM` + `BOOTH` + `CHARSAVE`):** Formato binario exacto de `ITEM_CELLCHANGE_CMD`, `ITEM_EQUIPCHANGE_CMD`, `EQUIP`/`UNEQUIP`, `RELOC`, `SPLIT`, `DROP`/`PICK`, `BUY`/`SELL`, `STORAGE`, **`CHARGEDINVENOPEN` + `CHARGED_WITHDRAW` (con entrega real en celda de inventario)**, **`REWARDINVENOPEN` + `REWARDINVENGET`**, `USE`, `UPGRADE`, `ENCHANT`, `RANDOMOPTION`, `MIX`, `DISMANTLE`, `PRODUCE`, SafeTrade y Kioscos `BOOTH`. |
| **Fase 20** | `NC_MENU` (11), `NC_QUEST` (17), `NC_CT` (27), `NC_SCRIPT` (28), `NC_MINIHOUSE` (33), `NC_COLLECT` (35), `NC_PET` (37), `NC_DICE` (40), `NC_GAMBLE` (42) | **Misiones, NPCs, Lua y Subsistemas de Contenido 1:1:** Diálogos/menús NPC (`MENU_SERVERMENU_CMD`, `NPCCLICK`), ciclo completo `QUEST_START`/`GIVEUP`/`REWARD`/`READ`/`SCRIPT_CMD`, triggers Lua, portales e instancias (`InstanceDungeon`, `KQ`), `MiniHouse` (23 C2S), `Pet`, `CharacterTitle` (`NC_CT`), `CollectCard` y Casino (`DICE` + `GAMBLE`). |
| **Fase 21** | `NC_PARTY` (14), `NC_GUILD` (21), `NC_PRISON` (23), `NC_FRIEND` (25), `NC_KQ` (29), `NC_WT` (30), `NC_PROMISE` (34), `NC_WEDDING` (38), `NC_GUILD_ACADEMY` (39), `NC_HOLY_PROMISE` (45) + `WorldService` (167 C2S) | **Sistemas Sociales, Puente S2S Vivo `World <-> Zone` y Auditoría Binaria Total:** Paridad 1:1 de Party, Raid, Guild, Friends, KQ, Academy, Master/Pupil, Wedding, HolyPromise, WeaponTitle y Prison entre `WorldManager.exe` y `Zone.exe`, enlace S2S en vivo entre procesos, y script final `scripts/verify_binary_exe_parity_100.sh` (`8/8 PASS`). |

