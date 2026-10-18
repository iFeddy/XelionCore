# Directivas y Arquitectura de XelionCore

## 1. Misión del Proyecto
Desarrollar un emulador del servidor oficial de Fiesta Online en Go (`XelionCore`) con **100% de fidelidad binaria, de protocolo y de base de datos** respecto a los binarios y símbolos originales (`AstraServer/Server/`).

---

## 2. Reglas de Oro (Mandatorias para la IA)
1. **Cero Omisiones, Cero Resúmenes y Nomenclatura Oficial 100% Idéntica:**
   - Está terminantemente prohibido usar "...", simplificar campos, o agrupar paquetes en listas genéricas.
   - **Nombres de Paquetes Idénticos a la Documentación:** En el código Go (structs, constantes, funciones y tests) se debe usar **estrictamente el nombre oficial exacto del paquete** documentado (`PROTO_NC_*`, por ejemplo `PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ`, `PROTO_NC_USER_US_LOGIN_REQ`, `PROTO_NC_USER_LOGIN_ACK`, etc.). Prohibido renombrar a nombres genéricos o usar nomenclaturas de emuladores privados (ej. `SH3Type_*`, `LoginReq`, etc.).
   - Toda estructura de paquete debe coincidir exactamente con los tipos C++, tamaños en bytes (`sizeof`) y offsets documentados en `docs/Servicios/`.
   - Nunca inventar estructuras ni usar esquemas de emuladores no oficiales / NextGen. Toda verdad técnica proviene de los PDBs y ejecutables oficiales.
2. **Fuente de la Verdad Oficial (Ground Truth):**
   - **Estructuras y Paquetes de Red:** [docs/Servicios/](docs/Servicios/) (generados directamente desde los flujos TPI de los PDBs oficiales).
   - **Bases de Datos y Tablas:** [docs/TablasLegacy.md](docs/TablasLegacy.md) (176 tablas oficiales con sus tipos exactos, PKs y FKs).
   - **Procedimientos Almacenados:** [docs/ProcedimientosLegacy.md](docs/ProcedimientosLegacy.md) y [docs/sql/procedures/](docs/sql/procedures/) (573 stored procedures recuperados byte a byte de los `.bak`).
   - **Archivos de Datos del Juego:** `../AstraServer/Server/9Data/Shine/` (archivos `.shn` y scripts `.txt`).
3. **Codificación de Protocolo (OpCodes):**
   - La cabecera `NETCOMMAND` oficial es una unión bitfield de 16 bits:
     * `command`: **10 bits** (Tail / Comando: valores 0 a 1023, máscara `0x3FF`).
     * `department`: **6 bits** (Head / Categoría: valores 0 a 63, shift de 10 bits).
   - Fórmula matemática estricta: `protocol = (department << 10) | command`.

---

## 3. Mapa de la Documentación Técnica Oficial

| Archivo | Contenido y Alcance Técnico |
| :--- | :--- |
| [docs/Servicios.md](docs/Servicios.md) | Visión general de la arquitectura, tabla comparativa Legacy vs. Go y enlaces a todas las especificaciones. |
| [docs/Servicios/Login.md](docs/Servicios/Login.md) | 85 estructuras de red completas, sockets C2S (`9010`) y S2S (`9011`), verificación de versión en `ClientVersionKeyInfo.txt`, zero-SQL. |
| [docs/Servicios/Account.md](docs/Servicios/Account.md) | 101 estructuras, autenticación en base de datos `Account`, 13 Stored Procedures con firmas ODBC exactas y checksums de 53 SHNs. |
| [docs/Servicios/AccountLog.md](docs/Servicios/AccountLog.md) | 106 estructuras, auditoría de sesiones y firewall de IPs en base de datos `AccountLog`, 12 Stored Procedures. |
| [docs/Servicios/WorldManager.md](docs/Servicios/WorldManager.md) | 1.379 estructuras, 463 handlers S2S de coordinación de mundo, gremios, bodas, selección de personajes y 71 SHNs leídos. |
| [docs/Servicios/Zone.md](docs/Servicios/Zone.md) | 1.343 estructuras, game loop, combate, inventario, movimiento, 189 SHNs y 31 scripts `.txt` de configuración. |
| [docs/Servicios/Character.md](docs/Servicios/Character.md) | 1.369 estructuras S2S, 196 handlers, 311 Stored Procedures contra `World00_Character` y 54 SHNs verificados. |
| [docs/Servicios/GameLog.md](docs/Servicios/GameLog.md) | 44 estructuras `PROTO_NC_LOG_*`, 23 Stored Procedures contra `World00_GameLog` y 49 SHNs verificados. |
| [docs/TablasLegacy.md](docs/TablasLegacy.md) | Diccionario de datos de 176 tablas del clúster oficial. |
| [docs/ProcedimientosLegacy.md](docs/ProcedimientosLegacy.md) | Catálogo de 573 procedimientos almacenados legacy. |
| [docs/sql/procedures/](docs/sql/procedures/) | Scripts `.sql` originales listos para migración o consulta. |

---

## 4. Arquitectura de Implementación en Go (`XelionCore`)

En el legado (`AstraServer`), los servidores se comunicaban mediante complejas cadenas de sockets S2S (`Zone` -> `Character` -> SQL, `Login` -> `Account` -> SQL). En `XelionCore` en Go, unificamos la lógica aprovechando la concurrencia nativa (goroutines, canales y pools de MySQL):

1. **`LoginService`:**
   - Unifica la recepción de clientes de `Login.exe` con la validación de cuenta de `Account.exe` y el registro de auditoría de `AccountLog.exe`.
   - Lee `ClientVersionKeyInfo.txt` para el handshake inicial (`0x0C65` / `0x0C67`).
2. **`WorldService`:**
   - Unifica la coordinación global de `WorldManager.exe` con la persistencia de personajes de `Character.exe`.
3. **`ZoneService`:**
   - Maneja directamente el bucle de juego, NPCs, mobs y combate, cargando los `.shn` de mapas y scripts de datos.
4. **`GamigoZR`:**
   - **No se emula ni se requiere.** Fue únicamente un mock HTTP local para saltear un chequeo DRM que Gamigo introdujo en `Zone.exe` en 2015.

---

## 5. Mapa de Recursos Modernizados para `ZoneService` (`data/` y `pkg/data/`)

Todos los archivos de configuración y datos del mundo de `Zone.exe` han sido modernizados a JSON modular estructurado con cero pérdida de datos y cuentan con sus cargadores oficiales en Go (`pkg/data/*`):

| Subsistema | Ubicación de Datos (`data/`) | Cargador Go (`pkg/data/`) | Responsabilidad en `ZoneService` |
| :--- | :--- | :--- | :--- |
| **Field / Maps** | `data/World/Map/` (`FieldList.json`, `InstanceDungeon.json`, `GuildTournament.json`, `TutorialMap.json`) | `pkg/data/field` | Configuración de campos, atributos de mapa, coordenadas de revivir, flags de PvP y binding de zonas (`zoneNumber`). |
| **NPCs & Links** | `data/World/NPC/` (`Maps/<Map>.json`, `LinkTable.json`, `NPCAction.json`) | `pkg/data/npc` | Spawns de NPCs por mapa (activos e inactivos), portales/warps entre mapas y máquinas de estado/acciones de NPCs. |
| **Mob Regen** | `data/MobRegen/<Map>.json` | `pkg/data/mobregen` | Generación y respawn de monstruos en cada mapa del mundo. |
| **Mob Setting** | `data/MobSetting/<Mob>.json` | `pkg/data/mobsetting` | Atributos individuales de monstruos, rangos de visión, persecución y triggers de acción unificados. |
| **Mob Behavior** | `data/MobBehaviorDescript.json` | `pkg/data/mobsetting` | Máquinas de estados y comportamientos de IA de monstruos. |
| **Mob Roam** | `data/MobRoam/<Map>.json` | `pkg/data/mobroam` | Rutas de patrullaje y movimiento de monstruos por mapa. |
| **Mob Attacks** | `data/MobAttackSequence/<Mob>.json` | `pkg/data/mobattack` | Secuencias de animación y patrones de ataque de monstruos. |
| **Drops** | `data/World/Drop/` (`ItemDropTable.json`, `ItemDropGroup.json`) | `pkg/data/drop` | Tablas de drop por monstruo (arrays dinámicos sin slots nulos) y grupos de drop probabilísticos. |
| **Combat** | `data/World/Combat/` (`Common.json`, `ExpTable.json`, `FreeStatTable.json`, `DamageByAngle.json`, etc.) | `pkg/data/combat` | Fórmulas de combate, tablas de experiencia, distribución de stats y multiplicadores posicionales. |
| **Class Param** | `data/World/ClassParam/<Class>.json` (27 clases) | `pkg/data/classparam` | Tablas de stats base, crecimiento y fórmulas por nivel (1 al 150) para cada clase del juego. |
| **Items** | `data/World/Item/` (`ItemOptions.json`, `ItemUseFunction.json`, `RandomOptionTable.json`, `TreasureReward.json`) | `pkg/data/item` | Opciones aleatorias de equipamiento, cofres del tesoro y scripts de uso de ítems. |
| **NPC Merchant** | `data/NPCItemList/<NPC>.json` | `pkg/data/npcitem` | Inventario y venta de ítems de mercaderes por NPC individual. |
| **AbState** | `data/AbState/<AbState>.json` | `pkg/data/abstate` | Efectos de estados anormales, buffs, debuffs y ticks periódicos. |
| **Misc** | `data/World/Misc/` (`RecallCoord.json`, `StealthMatrix.json`, `SkillMechanics.json`, `MobChat.json`) | `pkg/data/misc` | Pergaminos de retorno, capas de invisibilidad, totems/trampas y diálogos de monstruos. |
| **Quest Parser** | `data/World/Quest/QuestParser.json` | `pkg/data/questparser` | Mapeo de IDs de clases de misiones (`gQuestParserData`). |
| **Scripts & Lua** | `data/Script/*.json` y `data/LuaScript/` | `pkg/data/script`, `pkg/script/lua` | Lógica de eventos, diálogos y scripts de misiones ejecutados en runtime Lua. |
| **SHN Tables** | `../AstraServer/Server/9Data/Shine/` | `pkg/data/shn` | Universo binario de datos oficiales (189 SHNs). |
