import struct
import os
import sys
import re
from scripts.generate_service_docs import PDBParser, format_struct_markdown
from scripts.generate_all_service_docs import OPCODE_MAP

# Load WM handlers from WorldManager.pdb
def extract_wm_handlers():
    from scripts.parse_pdb_types import MSFReader
    msf = MSFReader("../AstraServer/Server/WorldManager/WorldManager.pdb")
    dbi = msf.get_stream(3)
    symrec_stream, = struct.unpack_from("<H", dbi, 20)
    syms = msf.get_stream(symrec_stream)
    pos = 0
    handlers = {}
    while pos < len(syms):
        rec_len, = struct.unpack_from("<H", syms, pos)
        if rec_len == 0: break
        rectype, = struct.unpack_from("<H", syms, pos + 2)
        rdata = syms[pos+4 : pos+2+rec_len]
        if rectype == 0x110e and len(rdata) >= 10:
            flags, off, seg = struct.unpack_from("<IIH", rdata, 0)
            name = rdata[10:].split(b"\x00")[0].decode("latin1", errors="ignore")
            if "fc_NC_" in name:
                clean = name.split("@")[0].replace("?fc_", "fc_")
                handlers[clean] = (seg, off)
        pos += 2 + rec_len
    return handlers

WM_HEADER = """
## Ficha Técnica del Binario
- **Nombre de Binario:** `WorldManager.exe`
- **Símbolos:** `WorldManager.pdb` (Microsoft C/C++ MSF 7.00, TPI v4.10)
- **Ruta del Proyecto Original:** `E:\\projectf2\\server\\3loginserver2\\worldmanager\\`
- **Base Address en Memoria:** `0x00400000`
- **Conexión SQL / ODBC:** **NINGUNA (0 conexiones directas a SQL)**.
  - `WorldManager.exe` **no** importa `ODBC32.dll`.
  - Todas las operaciones de base de datos de juego (creación/borrado de avatares, gremios, bodas, rankings de torneos) se delegan vía socket S2S a `Character.exe` (`PG_W00_CDB`, puerto `9033`).
- **Archivos Locales Leídos:**
  - `WMServerInfo.txt`: Configuración de enlace de mundos y parámetros del WorldManager.
  - `CallStack.txt`, `ExecUpdateInfo.txt`: Diagnóstico y versionado.
- **Archivos SHN Leídos (71 Archivos):**
  - Carga tablas globales para sincronizar el estado entre todas las zonas: `MapInfo.shn`, `MapLinkPoint.shn`, `MobInfo.shn`, `NPCViewInfo.shn`, `UpgradeInfo.shn`, `WeaponAttrib.shn`, `DiceDividind.shn`, `MiniHouse.shn`, etc.
- **Puertos de Red:**
  - **9014 (TCP):** Puerto principal S2S donde se conectan `Login.exe` y todas las instancias de `Zone.exe` (`PG_W00_WM`).
  - **9015 (TCP):** Puerto secundario S2S para sincronización inter-servidor y herramientas administrativas.
  - **9033 (TCP Cliente):** Conexión saliente hacia `Character.exe` (`PG_W00_CDB`).

## Arquitectura de Red y Coordinación de Mundo
`WorldManager` es el nodo central de coordinación del mundo (WorldServer):
1. **Handover de Login a Zone:** Cuando un jugador elige un servidor en `Login.exe`, Login solicita permiso a `WorldManager` (`PROTO_NC_USER_WILLLOGIN_REQ` / `ACK`). Si WM aprueba el cupo, Login entrega al cliente la IP y puerto de la `Zone` correspondiente (`PROTO_NC_USER_WORLDSELECT_ACK`).
2. **Sincronización Inter-Zona:** Cuando un jugador cruza un portal o teletransporte entre dos zonas distintas (ej. Zona 0 a Zona 1), `WorldManager` coordina la desconexión segura en la zona origen y la admisión en la zona destino (`PROTO_NC_CHAR_MAPLOGIN_REQ`).
3. **Gestión Global:** Administra el estado global de Gremios (`GuildAcademy`, torneos), Bodas (`Wedding`), Canales de Chat global, y el recuento total de jugadores en línea.
"""

ZONE_HEADER = """
## Ficha Técnica del Binario
- **Nombre de Binario:** `Zone.exe` (ejecutado en instancias `Zone00`, `Zone01`, `Zone02`, `Zone03`, `Zone04`)
- **Símbolos:** `Zone.pdb` (Microsoft C/C++ MSF 7.00, TPI v4.10)
- **Ruta del Proyecto Original:** `E:\\projectf2\\server\\3loginserver2\\zoneserver\\`
- **Base Address en Memoria:** `0x00400000`
- **Conexión SQL / ODBC:** **NINGUNA (0 conexiones directas a SQL)**.
  - `Zone.exe` **no** importa `ODBC32.dll`.
  - Toda la persistencia de inventario, ítems, experiencia, stats, misiones y habilidades se delega por socket S2S a `Character.exe` (`PG_W00_CDB`, puerto `9033`).
  - Todos los eventos de auditoría y logs de juego se envían por socket S2S a `GameLog.exe` (`PG_W00_GLDB`, puerto `9034`).
- **Archivos Locales Leídos (31 Archivos TXT):**
  - Scripts de inicialización de mapa y contenido: `MapLoad.txt`, `NPC.txt`, `QuestParser.txt`, `TreasureReward.txt`, `ServerInfo_Zone*.txt`, `DefaultCharacterData.txt`, etc.
- **Archivos SHN Leídos (189 Archivos):**
  - Carga el universo completo de datos del juego en memoria RAM: `ItemInfo.shn`, `ActiveSkill.shn`, `PassiveSkill.shn`, `MobInfo.shn`, `UpgradeInfo.shn`, `MapWayPoint.shn`, `QuestInfo.shn`, `ProduceView.shn`, `SpamerPenalty.shn`, etc.
- **Puertos de Red:**
  - **Puertos C2S (Clientes):**
    - `9016`: Instancia `Zone00`
    - `9019`: Instancia `Zone01`
    - `9022`: Instancia `Zone02`
    - `9025`: Instancia `Zone03`
    - `9028`: Instancia `Zone04`
  - **Puertos S2S (Inter-Servidor):**
    - `9017`, `9018`: Enlaces S2S `Zone00` hacia `WorldManager`
    - `9020`, `9021`: Enlaces S2S `Zone01` hacia `WorldManager`
    - `9023`, `9024`: Enlaces S2S `Zone02` hacia `WorldManager`
    - `9026`, `9027`: Enlaces S2S `Zone03` hacia `WorldManager`
    - `9029`, `9030`: Enlaces S2S `Zone04` hacia `WorldManager`
  - **Conexiones Salientes:**
    - Puerto `9033`: Hacia `Character.exe` (`GameDBSession`).
    - Puerto `9034`: Hacia `GameLog.exe` (`GameLogSession`).

## Arquitectura del Game Loop y Sesiones
1. **`ClientSession`:** Maneja el socket directo con cada cliente conectado al mapa. Valida semillas de encriptación (`Send_NC_MISC_SEED_ACK`) y procesa acciones de movimiento, combate y casteo de habilidades.
2. **`GameDBSession`:** Enlace directo con `Character.exe`. Despacha guardados asíncronos periódicos (`gds_NC_CHARSAVE_ALL_REQ`), uso de ítems (`gds_NC_ITEMDB_USE`), mejoras (`gds_NC_ITEMDB_UPGRADE`), y recompensas de Kingdom Quest (`gds_NC_KQ_REWARD`).
3. **`GameLogSession`:** Enlace directo con `GameLog.exe`. Envía logs de drop, compra/venta, muerte y eventos con las estructuras `PROTO_NC_LOG_GAME_DATA_TYPE_0` a `TYPE_J`.
"""

CHARACTER_HEADER = """
## Ficha Técnica del Binario
- **Nombre de Binario:** `Character.exe`
- **Ruta del Proyecto Original:** `D:\\Jenkins\\workspace\\Server - US - Release\\Server\\2DataServer\\DataServer\\Character Release US\\`
- **Base Address en Memoria:** `0x00400000`
- **Conexión SQL / ODBC:** **ACTIVA (`ODBC32.dll`)**.
  - Conecta a la base de datos `World00_Character`.
  - Ejecuta **311 Stored Procedures** para todas las transacciones críticas del juego.
- **Archivos Locales Leídos:**
  - `DataServerInfo_Character.txt`: Configuración de puertos y conexión ODBC.
  - `DefaultCharacterData.txt`: Estadísticas, equipo inicial, dinero y ubicación de inicio al crear un personaje por clase.
  - `TutorialCharacterData.txt`: Parámetros iniciales para la zona de tutorial.
- **Archivos SHN Verificados (54 Archivos):**
  - Compila `shndatafilechecksum.cpp` que valida los hashes MD5 de 54 archivos `.shn` para asegurar consistencia con el cliente y las zonas.
- **Puertos de Red:**
  - **9033 (TCP):** Puerto de escucha S2S (`PG_W00_CDB`). Recibe conexiones de `WorldManager.exe` y todas las instancias de `Zone.exe`.

## Procedimientos Almacenados SQL Invocados (World00_Character)
Todos los Stored Procedures están completamente recuperados y documentados en [docs/sql/procedures/World00_Character_all_procedures.sql](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/World00_Character_all_procedures.sql).
"""

GAMELOG_HEADER = """
## Ficha Técnica del Binario
- **Nombre de Binario:** `GameLog.exe`
- **Base Address en Memoria:** `0x00400000`
- **Conexión SQL / ODBC:** **ACTIVA (`ODBC32.dll`)**.
  - Conecta a la base de datos `World00_GameLog`.
  - Ejecuta Stored Procedures para inserción de logs estructurados por tipo de evento.
- **Archivos Locales Leídos:**
  - `DataServerInfo_GameLog.txt`: Configuración de puerto y conexión ODBC.
- **Archivos SHN Verificados (49 Archivos):**
  - Compila `shndatafilechecksum.cpp` para validar integridad de datos.
- **Puertos de Red:**
  - **9034 (TCP):** Puerto de escucha S2S (`PG_W00_GLDB`). Recibe conexiones de todas las instancias de `Zone.exe` y `WorldManager.exe`.

## Procedimientos Almacenados SQL Invocados (World00_GameLog)
Todos los Stored Procedures están completamente recuperados en [docs/sql/procedures/World00_GameLog_all_procedures.sql](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/World00_GameLog_all_procedures.sql):

| Stored Procedure | Parámetros (ODBC CALL) | Tipo de Evento / Propósito |
| :--- | :--- | :--- |
| `p_GameLog_Add` | `{CALL p_GameLog_Add( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Inserción genérica de log de juego |
| `p_GameLog_Add_Type1` | `{CALL p_GameLog_Add_Type1( ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo 1: Login / Logout de personaje en mapa |
| `p_GameLog_Add_Type2` | `{CALL p_GameLog_Add_Type2( ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo 2: Subida de nivel y estadísticas |
| `p_GameLog_Add_Type3` | `{CALL p_GameLog_Add_Type3( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo 3: Movimientos de dinero (compra/venta NPC) |
| `p_GameLog_Add_Type4` | `{CALL p_GameLog_Add_Type4( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo 4: Comercio directo entre jugadores (Trade) |
| `p_GameLog_Add_Type5` | `{CALL p_GameLog_Add_Type5( ?, ?, ?, ?, ? )}` | Log Tipo 5: Interacción con almacén / depósito |
| `p_GameLog_Add_Type6` | `{CALL p_GameLog_Add_Type6( ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo 6: Obtención y descarte de ítems (Drop / Pick) |
| `p_GameLog_Add_Type7` | `{CALL p_GameLog_Add_Type7( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo 7: Mejora y encantamiento de armas/armaduras |
| `p_GameLog_Add_Type8` | `{CALL p_GameLog_Add_Type8( ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo 8: Muerte de personaje (PvP / PvE) |
| `p_GameLog_Add_Type9` | `{CALL p_GameLog_Add_Type9( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo 9: Progreso y finalización de misiones (Quests) |
| `p_GameLog_Add_TypeA` | `{CALL p_GameLog_Add_TypeA( ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo A: Eventos de gremio (creación, rangos) |
| `p_GameLog_Add_TypeB` | `{CALL p_GameLog_Add_TypeB( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo B: Almacén de gremio |
| `p_GameLog_Add_TypeC` | `{CALL p_GameLog_Add_TypeC( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo C: Tiendas personales (Kiosk / Booth) |
| `p_GameLog_Add_TypeD` | `{CALL p_GameLog_Add_TypeD( ?, ?, ?, ?, ?, ? )}` | Log Tipo D: Producción / Crafteo de ítems |
| `p_GameLog_Add_TypeE` | `{CALL p_GameLog_Add_TypeE( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo E: Aprendizaje y reinicio de habilidades |
| `p_GameLog_Add_TypeF` | `{CALL p_GameLog_Add_TypeF( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo F: Casa personal (MiniHouse) |
| `p_GameLog_Add_TypeG` | `{CALL p_GameLog_Add_TypeG( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo G: Mascotas y monturas |
| `p_GameLog_Add_TypeH` | `{CALL p_GameLog_Add_TypeH( ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo H: Kingdom Quest y eventos comunitarios |
| `p_GameLog_Add_TypeI` | `{CALL p_GameLog_Add_TypeI( ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo I: Casino / Tragamonedas / Apuestas |
| `p_GameLog_Add_TypeJ` | `{CALL p_GameLog_Add_TypeJ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo J: Bodas y divorcios |
| `p_GameLog_Add_TypeK` | `{CALL p_GameLog_Add_TypeK( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo K: Academia de gremio |
| `p_GameLog_Add_TypeL` | `{CALL p_GameLog_Add_TypeL( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )}` | Log Tipo L: Sistema de reportes y sanciones |
| `usp_Login_setChannel` | `{CALL usp_Login_setChannel( ?, ?, ?, ? )}` | Registro de canal de conexión |
"""

def generate_doc_from_pdb(pdb_path, output_md, title, header):
    print(f"Generating {output_md} from {pdb_path}...")
    parser = PDBParser(pdb_path)
    proto_structs = [s for s in parser.structs if "tagNC_" in s or "PROTO_NC_" in s]
    
    lines = []
    lines.append(f"# Especificación Oficial: `{title}`\n")
    lines.append(header)
    lines.append("\n## Catálogo Exhaustivo de Estructuras de Paquetes\n")
    lines.append(f"> **Total de Estructuras Catalogadas:** `{len(proto_structs)}`\n")
    lines.append("> Todas las estructuras han sido extraídas directamente del flujo TPI del PDB oficial con sus tamaños, tipos primitivos, offsets y variables exactas.\n")

    for sname in sorted(proto_structs):
        sinfo = parser.structs[sname]
        members = parser.get_members(sname)
        opcode_info = OPCODE_MAP.get(sname, "")
        if not opcode_info:
            clean_name = sname.split("<")[0].split("::")[-1]
            opcode_info = OPCODE_MAP.get(clean_name, "")
        
        md_text = format_struct_markdown(sname, sinfo, members, opcode_info)
        lines.append(md_text)

    with open(output_md, "w", encoding="utf-8") as f:
        f.write("\n".join(lines))
    print(f"Successfully generated {output_md} ({len(proto_structs)} structs).")

def generate_character_doc():
    output_md = "docs/Servicios/Character.md"
    print(f"Generating {output_md}...")
    # Character structures are in Zone and WorldManager PDBs
    zone_parser = PDBParser("../AstraServer/Server/Zone00/Zone.pdb")
    wm_parser = PDBParser("../AstraServer/Server/WorldManager/WorldManager.pdb")
    
    with open("../AstraServer/Server/Character/Character.exe", "rb") as f: data = f.read()
    handlers = sorted(set(re.findall(rb"fc_NC_[a-zA-Z0-9_]+", data)))
    
    # Extract unique character-related structures
    char_structs = set()
    for s in list(zone_parser.structs.keys()) + list(wm_parser.structs.keys()):
        if any(k in s for k in ["PROTO_NC_CHAR", "PROTO_NC_ITEMDB", "PROTO_NC_GUILD", "PROTO_NC_PET", "PROTO_NC_QUEST", "PROTO_NC_SKILL", "PROTO_NC_COLLECT", "PROTO_NC_GAMBLE", "PROTO_NC_HOLY", "PROTO_NC_MINIHOUSE", "PROTO_NC_USP_USER_CHARACTER", "PROTO_NC_AVATAR"]):
            char_structs.add(s)

    lines = []
    lines.append("# Especificación Oficial: `Character (CharacterDataServer)`\n")
    lines.append(CHARACTER_HEADER)
    lines.append("\n## Catálogo de Handlers de Red (`fc_NC_*` en `Character.exe`)")
    lines.append(f"> **Total de Handlers Extraídos del Binario:** `{len(handlers)}`\n")
    lines.append("| N° | Función Handler | Propósito Técnico en Base de Datos |")
    lines.append("| :--- | :--- | :--- |")
    for i, h_b in enumerate(handlers, 1):
        h = h_b.decode("latin1")
        lines.append(f"| {i} | `{h}` | Despacho de persistencia / lectura hacia `World00_Character` |")

    lines.append(f"\n## Catálogo Exhaustivo de Estructuras de Paquetes S2S de `Character`\n")
    lines.append(f"> **Total de Estructuras Catalogadas:** `{len(char_structs)}`\n")
    lines.append("> Estructuras resueltas con paridad binaria exacta a partir de `Zone.pdb` y `WorldManager.pdb`.\n")

    for sname in sorted(char_structs):
        parser = zone_parser if sname in zone_parser.structs else wm_parser
        sinfo = parser.structs[sname]
        members = parser.get_members(sname)
        opcode_info = OPCODE_MAP.get(sname, "")
        md_text = format_struct_markdown(sname, sinfo, members, opcode_info)
        lines.append(md_text)

    with open(output_md, "w", encoding="utf-8") as f:
        f.write("\n".join(lines))
    print(f"Successfully generated {output_md} ({len(char_structs)} structs).")

def generate_gamelog_doc():
    output_md = "docs/Servicios/GameLog.md"
    print(f"Generating {output_md}...")
    zone_parser = PDBParser("../AstraServer/Server/Zone00/Zone.pdb")
    
    log_structs = [s for s in zone_parser.structs if "PROTO_NC_LOG" in s or "SYSLOG" in s]

    lines = []
    lines.append("# Especificación Oficial: `GameLog (GameLogDataServer)`\n")
    lines.append(GAMELOG_HEADER)
    lines.append(f"\n## Catálogo Exhaustivo de Estructuras de Logs S2S (`PROTO_NC_LOG_*`)\n")
    lines.append(f"> **Total de Estructuras Catalogadas:** `{len(log_structs)}`\n")
    lines.append("> Estructuras resueltas con paridad binaria exacta a partir de `Zone.pdb`.\n")

    for sname in sorted(log_structs):
        sinfo = zone_parser.structs[sname]
        members = zone_parser.get_members(sname)
        opcode_info = OPCODE_MAP.get(sname, "")
        md_text = format_struct_markdown(sname, sinfo, members, opcode_info)
        lines.append(md_text)

    with open(output_md, "w", encoding="utf-8") as f:
        f.write("\n".join(lines))
    print(f"Successfully generated {output_md} ({len(log_structs)} structs).")

if __name__ == "__main__":
    generate_doc_from_pdb("../AstraServer/Server/WorldManager/WorldManager.pdb",
                          "docs/Servicios/WorldManager.md",
                          "WorldManager (WorldServer)",
                          WM_HEADER)

    generate_doc_from_pdb("../AstraServer/Server/Zone00/Zone.pdb",
                          "docs/Servicios/Zone.md",
                          "Zone (ZoneServer)",
                          ZONE_HEADER)

    generate_character_doc()
    generate_gamelog_doc()
