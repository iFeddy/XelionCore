# Directivas y Arquitectura de XelionCore

## 1. Misión del Proyecto
Desarrollar un emulador del servidor oficial de Fiesta Online en Go (`XelionCore`) con **100% de fidelidad binaria, de protocolo y de base de datos** respecto a los binarios y símbolos originales (`AstraServer/Server/`).

---

## 2. Reglas de Oro (Mandatorias para la IA)
1. **Cero Omisiones y Cero Resúmenes:**
   - Está terminantemente prohibido usar "...", simplificar campos, o agrupar paquetes en listas genéricas.
   - Toda estructura de paquete debe coincidir exactamente con los tipos C++, tamaños en bytes (`sizeof`) y offsets documentados.
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
