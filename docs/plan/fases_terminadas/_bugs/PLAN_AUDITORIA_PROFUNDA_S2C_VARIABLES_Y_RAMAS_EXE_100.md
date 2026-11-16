# Plan Maestro Etapa 5: Auditoría 1:1 de Paquetes S2C Variables y Barrido Rama por Rama de los 412 Handlers C2S contra `.exe` + `.pdb`

## 0. Objetivo y Alcance Técnico

Tras completar la **Etapa 4 (Fases 17 a 21)**, el proyecto cuenta con:
- **100% de OpCodes C2S registrados:** `49/49` en `LoginService`, `167/167` en `WorldService` y `245/245` en `ZoneService` (**412 handlers C2S** entre World y Zone).
- **100% de structs de tamaño fijo verificados en `pkg/protocol/`** y escaneo AST estático de paquetes S2C de tamaño fijo en `Control 8` de `pkg/audit/binary_exe_parity_completeness_test.go`.

Sin embargo, para alcanzar **paridad pre-lanzamiento total sin depender de descubrir errores manualmente con el cliente oficial** y sin que quede ningún subsistema interno en el aire, debemos cerrar **tres brechas técnicas** que un chequeo estático de tamaño fijo no cubre:

1. **Punto 1 — Paquetes S2C de tamaño variable (`[0];` / `[];`) y bloques `NewPacketWriter` con bucles o condicionales:**
   - En `docs/Servicios/WorldManager.md` existen **148 structs dinámicos** (sobre 1.054 con tamaño) y en `docs/Servicios/Zone.md` existen **155 structs dinámicos** (sobre 1.196 con tamaño).
   - En `pkg/server/world/` y `pkg/server/zone/` hay **567 llamadas a `protocol.NewPacketWriter(...)`** que emiten **436 OpCodes `PROTO_NC_*` únicos**, de los cuales **55 corresponden a structs con arreglos flexibles** (`list[0]`, `data[0]`, strings con longitud dinámica, etc.) o bloques con `for` / `if-else` que `auditFixedS2CInSource` hoy saltea cuando detecta escritura dinámica.
2. **Punto 2 — Barrido de lógica rama por rama (`if/else`, validaciones, códigos de error y secuencia de paquetes S2C) de los 412 handlers C2S contra el ensamblador de `WorldManager.exe` y `Zone.exe`:**
   - Un handler puede recibir el struct C2S correcto y emitir un paquete S2C con `sizeof` válido en el caso feliz, pero diferir del `.exe` en:
     - Qué condiciones de guardia (`if`) evalúa y en qué orden.
     - Qué código de error oficial (`uint16`) devuelve en el paquete `*_ACK` / `*_FAIL` cuando una validación falla.
     - Qué secuencia exacta de múltiples paquetes S2C (a sí mismo, a la lista AOI, al Party/Guild o por S2S a World/Zone) dispara en cada rama.
3. **Punto 3 — Subfunciones de Fórmulas Internas, Game Loop / Timers Autónomos (`OnTimer`) y Ejecución Real de `data/LuaScript/` sobre los 274 Bindings C-API:**
   - **Subfunciones matemáticas de combate (`Zone.exe`):** Fórmulas internas de ataque físico/mágico, precisión/evasión (`hit/evade`), golpe crítico, bloqueo, defensa, `ClassParam` y multiplicadores posicionales (`DamageByAngle`).
   - **Timers autónomos (`OnTimer`) en `Zone.exe` y `WorldManager.exe`:** Ticks periódicos que se ejecutan sin paquete C2S del cliente: ticks DoT/HoT y expiración de `AbState`, regeneración natural de HP/SP/LP, máquinas de estado `MobBehavior` / `MobRegen` / `MobRoam`, propiedad y desaparición de `Drops` en el suelo, auto-guardado periódico (`CHARSAVE`), y máquinas de tiempo de `KingdomQuest` (`KQ`), `GuildTournament`, `GuildWar`, `Wedding` e `InstanceDungeon`.
   - **Ejecución real de `data/LuaScript/` con los 274 bindings C-API (`cMobRegen_XY`, `cAIScriptFunc`, `cTimer`, `cLinkTo`, `cQuest*`, etc.):** Verificación contra `Zone.exe` de que cada binding C-API llamado por los scripts `.lua` oficiales de instancias, KQs, misiones y eventos produzca exactamente los mismos efectos sobre el mapa y los mismos paquetes S2C.

---

## 1. Reglas Inquebrantables de la Etapa 5 (Fases 22 a 26)

1. **Desensamblado Completo de Todas las Ramas, Subfunciones y Timers (`scripts/disasm_pdb_handler.py`):**
   - Para cada handler `sp_NC_*` (`Zone.exe`) y `fc_NC_*` (`WorldManager.exe`), además de las subfunciones matemáticas y timers `OnTimer` que invoca el motor, se deben inspeccionar **todos los saltos condicionales (`je`, `jne`, `jb`, `ja`, `jg`, `jl`, `switch`)** y las llamadas a `SendToClient` / `Broadcast` / `SendToWorld` en cada camino de ejecución.
2. **Fórmula Matemática Exacta para Paquetes S2C Variables:**
   - Todo paquete S2C con arreglo flexible `Entry list[0]` documentado en `Zone.md` o `WorldManager.md` debe cumplir en el wire:
     $$\text{WirePayloadSize} = \text{sizeof}(\text{HeaderBase}) + N \times \text{sizeof}(\text{Elemento})$$
     donde $N$ es exactamente el campo contador escrito en la cabecera del paquete.
   - Para paquetes con longitud de string variable (`msgLen` + `msg`), el número de bytes escritos tras la cabecera debe ser exactamente igual al campo `len` prefijado en el struct.
3. **Tests de Caso Feliz + Casos de Rechazo + Timers + Scripts Lua:**
   - Cada fase debe agregar pruebas automatizadas en `pkg/audit/` que ejerciten tanto la **rama de éxito** como las **ramas de error/rechazo**, además de los **ticks de timers autónomos**, **fórmulas de combate** y **scripts `.lua` oficiales**.

---

## 2. Partición en 5 Fases Modulares (Fases 22 a 26)

| Fase | Punto Cubierto | Alcance Técnico y Departamentos | Entregable de Auditoría |
| :---: | :---: | :--- | :--- |
| **Fase 22** | **Punto 1 (100% S2C Variables)** | Auditoría de los **148 structs dinámicos de World** + **155 de Zone**, revisión de las **567 llamadas a `NewPacketWriter`** (`436` OpCodes S2C únicos, incluyendo los `55` con `[0]`/`[]` y bloques con loops/ramas) e interceptor automático de tamaños dinámicos. | `pkg/audit/fase22_dynamic_s2c_parity_test.go` |
| **Fase 23** | **Punto 2 — Bloque A (Entrada, Mundo, AOI y Movimiento)** | Desensamblado rama por rama de `WorldService` (`NC_USER`, `NC_CHAR`, `NC_AVATAR`) y `ZoneService` (`NC_MISC`, `NC_USER`, `NC_CHAR`, `NC_MAP`, `NC_ACT`, `NC_MOVER`, `NC_BRIEFINFO`): validaciones, errores y secuencias S2C. | `pkg/audit/fase23_branch_entry_move_aoi_test.go` |
| **Fase 24** | **Puntos 2 y 3 — Bloque B (Combate, Fórmulas, Skills, AbState, Mobs, Drops y Timers de Zona)** | Desensamblado rama por rama en `Zone.exe` de `NC_BAT` (9), `NC_SKILL` (18), `NC_SOULSTONE` (22), `NC_CT` (27), `NC_COLLECT` (35) + **subfunciones matemáticas de combate (`hit/evade/crit/block/damage/angle`)** + **timers autónomos (`AbState` ticks, regen HP/SP/LP, IA/aggro/regen de Mobs, expiración/propiedad de Drops y `CHARSAVE` periódico)**. | `pkg/audit/fase24_branch_combat_skill_mob_test.go` |
| **Fase 25** | **Punto 2 — Bloque C (Economía, 48 `NC_ITEM`, Trade, Booth y Forja)** | Desensamblado rama por rama en `Zone.exe` de `NC_ITEM` (12), `NC_ITEMDB` (13), `NC_BOOTH` (16), `NC_CHARSAVE` (31) y `SafeTrade` (fallos de celda, oro, nivel, clase, upgrade/enchant/mix/dismantle/produce). | `pkg/audit/fase25_branch_item_economy_forge_test.go` |
| **Fase 26** | **Puntos 2 y 3 — Bloque D (Quests, 274 Bindings Lua + `data/LuaScript/`, Timers de Mundo/Instancias, Subsistemas y Sociales) + Suite Final** | Desensamblado rama por rama de `NC_MENU`, `NC_QUEST`, `NC_SCRIPT`, `NC_MINIHOUSE`, `NC_PET`, `NC_DICE`, `NC_GAMBLE` y sociales (`NC_PARTY`, `NC_GUILD`, `NC_FRIEND`, `NC_KQ`, `NC_WEDDING`, `NC_GUILD_ACADEMY`, `NC_PROMISE`, `NC_HOLY_PROMISE`, `NC_WT`, `NC_PRISON`) + **ejecución real de `data/LuaScript/` sobre los 274 bindings C-API** + **timers autónomos de `KQ`/`GuildTournament`/`GuildWar`/`Wedding`/`InstanceDungeon`** + script `scripts/verify_exe_branch_and_dynamic_s2c_100.sh` (`8/8 PASS`). | `pkg/audit/exe_branch_and_dynamic_s2c_completeness_test.go` + `scripts/verify_exe_branch_and_dynamic_s2c_100.sh` |

