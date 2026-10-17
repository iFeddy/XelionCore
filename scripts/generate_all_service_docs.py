import struct
import os
import sys
from scripts.generate_service_docs import PDBParser, format_struct_markdown

OPCODE_MAP = {
    # MISC (Dept: 2)
    "PROTO_NC_MISC_S2SCONNECTION_RDY": "0x0801 (Dept: 2, Cmd: 1)",
    "PROTO_NC_MISC_S2SCONNECTION_RDY_SEND": "0x0801 (Dept: 2, Cmd: 1)",
    "PROTO_NC_MISC_S2SCONNECTION_REQ": "0x0802 (Dept: 2, Cmd: 2)",
    "PROTO_NC_MISC_S2SCONNECTION_REQ_SEND": "0x0802 (Dept: 2, Cmd: 2)",
    "PROTO_NC_MISC_S2SCONNECTION_ACK": "0x0803 (Dept: 2, Cmd: 3)",
    "PROTO_NC_MISC_HEARTBEAT_REQ": "0x0804 (Dept: 2, Cmd: 4)",
    "PROTO_NC_MISC_HEARTBEAT_REQ_SEND": "0x0804 (Dept: 2, Cmd: 4)",
    "PROTO_NC_MISC_HEARTBEAT_ACK": "0x0805 (Dept: 2, Cmd: 5)",
    "PROTO_NC_MISC_HEARTBEAT_ACK_SEND": "0x0805 (Dept: 2, Cmd: 5)",
    "PROTO_NC_MISC_USER_COUNT_CMD": "0x0808 (Dept: 2, Cmd: 8)",
    "PROTO_NC_MISC_CONNECTER_REQ": "0x0809 (Dept: 2, Cmd: 9)",
    "PROTO_NC_MISC_CONNECTER_ACK": "0x080A (Dept: 2, Cmd: 10)",
    "PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD": "0x080B (Dept: 2, Cmd: 11)",
    "PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD": "0x080C (Dept: 2, Cmd: 12)",
    "PROTO_NC_MISC_WEB_DB_KEY_REQ": "0x081D (Dept: 2, Cmd: 29)",
    "PROTO_NC_MISC_CHAR_LOGOFF_STATISTICS": "0x081F (Dept: 2, Cmd: 31)",
    "PROTO_NC_MISC_EVENT_HIT3_ADD_CASH": "0x0820 (Dept: 2, Cmd: 32)",
    "PROTO_NC_MISC_START_THE_BOOM_CMD": "0x0823 (Dept: 2, Cmd: 35)",
    "PROTO_NC_MISC_EVENT_L20_DB_REQ": "0x082C (Dept: 2, Cmd: 44)",
    "PROTO_NC_MISC_EVENT_L20_DB_ACK": "0x082D (Dept: 2, Cmd: 45)",
    "PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ": "0x082F (Dept: 2, Cmd: 47)",
    "PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD": "0x0835 (Dept: 2, Cmd: 53)",
    "PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC": "0x0838 (Dept: 2, Cmd: 56)",
    "PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD": "0x0843 (Dept: 2, Cmd: 67)",
    "PROTO_NC_MISC_SEED_ACK": "0x0807 (Dept: 2, Cmd: 7)",

    # USER (Dept: 3)
    "PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK": "0x0C02 (Dept: 3, Cmd: 2)",
    "PROTO_NC_USER_XTRAP_REQ": "0x0C04 (Dept: 3, Cmd: 4)",
    "PROTO_NC_USER_XTRAP_ACK": "0x0C05 (Dept: 3, Cmd: 5)",
    "PROTO_NC_USER_LOGIN_REQ": "0x0C06 (Dept: 3, Cmd: 6)",
    "PROTO_NC_USER_PASSWORD_CHECK_REQ": "0x0C07 (Dept: 3, Cmd: 7)",
    "PROTO_NC_USER_PASSWORD_CHECK_REQ_SEND": "0x0C07 (Dept: 3, Cmd: 7)",
    "PROTO_NC_USER_PASSWORD_CHECK_ACK": "0x0C08 (Dept: 3, Cmd: 8)",
    "PROTO_NC_USER_PASSWORD_CHECK_ACK_SEND": "0x0C08 (Dept: 3, Cmd: 8)",
    "PROTO_NC_USER_LOGINFAIL_ACK": "0x0C09 (Dept: 3, Cmd: 9)",
    "PROTO_NC_USER_LOGINFAIL_ACK_SEND": "0x0C09 (Dept: 3, Cmd: 9)",
    "PROTO_NC_USER_LOGIN_ACK": "0x0C0A (Dept: 3, Cmd: 10)",
    "PROTO_NC_USER_WORLDSELECT_REQ": "0x0C0B (Dept: 3, Cmd: 11)",
    "PROTO_NC_USER_WORLDSELECT_ACK": "0x0C0C (Dept: 3, Cmd: 12)",
    "PROTO_NC_USER_WORLDSELECT_ACK_SEND": "0x0C0C (Dept: 3, Cmd: 12)",
    "PROTO_NC_USER_WILLLOGIN_REQ": "0x0C0D (Dept: 3, Cmd: 13)",
    "PROTO_NC_USER_WILLLOGIN_ACK": "0x0C0E (Dept: 3, Cmd: 14)",
    "PROTO_NC_USER_NORMALLOGOUT_CMD": "0x0C18 (Dept: 3, Cmd: 24)",
    "PROTO_NC_USER_WORLD_STATUS_REQ": "0x0C1B (Dept: 3, Cmd: 27)",
    "PROTO_NC_USER_WORLD_STATUS_ACK": "0x0C1C (Dept: 3, Cmd: 28)",
    "PROTO_NC_USER_LOGIN_NETMARBLE_REQ": "0x0C1D (Dept: 3, Cmd: 29)",
    "PROTO_NC_USER_LOGIN_NETMARBLE_DB_REQ": "0x0C1E (Dept: 3, Cmd: 30)",
    "PROTO_NC_USER_LOGIN_OUTSPARK_REQ": "0x0C20 (Dept: 3, Cmd: 32)",
    "PROTO_NC_USER_TEENAGER_REQ": "0x0C26 (Dept: 3, Cmd: 38)",
    "PROTO_NC_USER_TEENAGER_ACK": "0x0C27 (Dept: 3, Cmd: 39)",
    "PROTO_NC_USER_TEENAGER_SET_CMD": "0x0C28 (Dept: 3, Cmd: 40)",
    "PROTO_NC_USER_IS_IP_BLOCK_REQ": "0x0C2A (Dept: 3, Cmd: 42)",
    "PROTO_NC_USER_IS_IP_BLOCK_ACK": "0x0C2B (Dept: 3, Cmd: 43)",
    "PROTO_NC_USER_CREATE_OTP_REQ": "0x0C35 (Dept: 3, Cmd: 53)",
    "PROTO_NC_USER_CREATE_OTP_ACK": "0x0C36 (Dept: 3, Cmd: 54)",
    "PROTO_NC_USER_LOGIN_WITH_OTP_REQ": "0x0C37 (Dept: 3, Cmd: 55)",
    "PROTO_NC_USER_GER_PASSWORD_CHECK_REQ": "0x0C39 (Dept: 3, Cmd: 57)",
    "PROTO_NC_USER_GER_PASSWORD_CHECK_ACK": "0x0C3A (Dept: 3, Cmd: 58)",
    "PROTO_NC_USER_GER_PASSWORD_CHECK_ACK_SEND": "0x0C3A (Dept: 3, Cmd: 58)",
    "PROTO_NC_USER_GER_IS_IP_BLOCK_REQ": "0x0C3B (Dept: 3, Cmd: 59)",
    "PROTO_NC_USER_GER_IS_IP_BLOCK_ACK": "0x0C3C (Dept: 3, Cmd: 60)",
    "PROTO_NC_USER_TW_PASSWORD_CHECK_REQ": "0x0C3E (Dept: 3, Cmd: 62)",
    "PROTO_NC_USER_TW_PASSWORD_CHECK_ACK": "0x0C3F (Dept: 3, Cmd: 63)",
    "PROTO_NC_USER_TW_PASSWORD_CHECK_ACK_SEND": "0x0C3F (Dept: 3, Cmd: 63)",
    "PROTO_NC_USER_TW_IS_IP_BLOCK_REQ": "0x0C40 (Dept: 3, Cmd: 64)",
    "PROTO_NC_USER_TW_IS_IP_BLOCK_ACK": "0x0C41 (Dept: 3, Cmd: 65)",
    "PROTO_NC_USER_JP_PASSWORD_CHECK_REQ": "0x0C4B (Dept: 3, Cmd: 75)",
    "PROTO_NC_USER_JP_PASSWORD_CHECK_ACK": "0x0C4C (Dept: 3, Cmd: 76)",
    "PROTO_NC_USER_JP_IS_IP_BLOCK_REQ": "0x0C4D (Dept: 3, Cmd: 77)",
    "PROTO_NC_USER_JP_IS_IP_BLOCK_ACK": "0x0C4E (Dept: 3, Cmd: 78)",
    "PROTO_NC_USER_CH_PASSWORD_CHECK_REQ": "0x0C53 (Dept: 3, Cmd: 83)",
    "PROTO_NC_USER_CH_PASSWORD_CHECK_ACK": "0x0C54 (Dept: 3, Cmd: 84)",
    "PROTO_NC_USER_CH_PASSWORD_CHECK_ACK_SEND": "0x0C54 (Dept: 3, Cmd: 84)",
    "PROTO_NC_USER_US_LOGIN_REQ": "0x0C5A (Dept: 3, Cmd: 90)",
    "PROTO_NC_USER_US_PASSWORD_CHECK_REQ": "0x0C5B (Dept: 3, Cmd: 91)",
    "PROTO_NC_USER_US_PASSWORD_CHECK_ACK": "0x0C5C (Dept: 3, Cmd: 92)",
    "PROTO_NC_USER_US_IS_IP_BLOCK_REQ": "0x0C5D (Dept: 3, Cmd: 93)",
    "PROTO_NC_USER_US_IS_IP_BLOCK_ACK": "0x0C5E (Dept: 3, Cmd: 94)",
    "PROTO_NC_USER_RETURN_CHECK_REQ": "0x0C60 (Dept: 3, Cmd: 96)",
    "PROTO_NC_USER_RETURN_CHECK_ACK": "0x0C61 (Dept: 3, Cmd: 97)",
    "PROTO_NC_USER_LOGOUT_LAST_DAY_REQ": "0x0C62 (Dept: 3, Cmd: 98)",
    "PROTO_NC_USER_LOGOUT_LAST_DAY_ACK": "0x0C63 (Dept: 3, Cmd: 99)",
    "PROTO_NC_USER_SET_RETURN_CMD": "0x0C64 (Dept: 3, Cmd: 100)",
    "PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ": "0x0C65 (Dept: 3, Cmd: 101)",
    "PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK": "0x0C67 (Dept: 3, Cmd: 103)",

    # LOG (Dept: 1)
    "PROTO_NC_LOG_USER_LOGINFAIL": "0x0720 (Dept: 1, Cmd: 800)",
    "PROTO_NC_LOG_USER_LOGINFAIL_SEND": "0x0720 (Dept: 1, Cmd: 800)",
    "PROTO_NC_LOG_USER_LOGIN": "0x0721 (Dept: 1, Cmd: 801)",
    "PROTO_NC_LOG_USER_LOGIN_SEND": "0x0721 (Dept: 1, Cmd: 801)",
    "PROTO_NC_LOG_USER_LOGOUT": "0x0722 (Dept: 1, Cmd: 802)",

    # OPTOOL (Dept: 10)
    "PROTO_NC_OPTOOL_S2SCONNECT_LIST_REQ": "0x2801 (Dept: 10, Cmd: 1)",
    "PROTO_NC_OPTOOL_S2SCONNECT_LIST_ACK": "0x2802 (Dept: 10, Cmd: 2)",
    "PROTO_NC_OPTOOL_CLOSE_SERVER_REQ": "0x2803 (Dept: 10, Cmd: 3)",
    "PROTO_NC_OPTOOL_CLOSE_SERVER_ACK": "0x2804 (Dept: 10, Cmd: 4)",
    "PROTO_NC_OPTOOL_MAP_USER_LIST_REQ": "0x2805 (Dept: 10, Cmd: 5)",
    "PROTO_NC_OPTOOL_MAP_USER_LIST_ACK": "0x2806 (Dept: 10, Cmd: 6)",
    "PROTO_NC_OPTOOL_CONNECT_BRIF_REQ": "0x2807 (Dept: 10, Cmd: 7)",
    "PROTO_NC_OPTOOL_CONNECT_BRIF_ACK": "0x2808 (Dept: 10, Cmd: 8)",
    "PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_REQ": "0x2812 (Dept: 10, Cmd: 18)",
    "PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_ACK": "0x2813 (Dept: 10, Cmd: 19)",
    "PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_SET_CMD": "0x2814 (Dept: 10, Cmd: 20)",
    "PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_CLR": "0x2817 (Dept: 10, Cmd: 23)",
    "PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_REQ": "0x2818 (Dept: 10, Cmd: 24)",
    "PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_ACK": "0x2819 (Dept: 10, Cmd: 25)",

    # REPORT (Dept: 32)
    "PROTO_NC_DATA_REPORT_ADD_REQ": "0x8007 (Dept: 32, Cmd: 7)",
    "PROTO_NC_DATA_REPORT_ADD_ACK": "0x8008 (Dept: 32, Cmd: 8)",
    "PROTO_NC_DATA_REPORT_CANCEL_REQ": "0x8009 (Dept: 32, Cmd: 9)",
    "PROTO_NC_DATA_REPORT_CANCEL_ACK": "0x800A (Dept: 32, Cmd: 10)",
    "PROTO_NC_DATA_REPORT_GET_REQ": "0x800B (Dept: 32, Cmd: 11)",
    "PROTO_NC_DATA_REPORT_GET_ACK": "0x800C (Dept: 32, Cmd: 12)",

    # SYSLOG (Dept: 50)
    "PROTO_NC_SYSLOG_ACCOUNT_LOGIN_SUCCESS": "0xC801 (Dept: 50, Cmd: 1)",
    "PROTO_NC_SYSLOG_ACCOUNT_LOGIN_FAILURE": "0xC802 (Dept: 50, Cmd: 2)",
    "PROTO_NC_SYSLOG_ACCOUNT_LOGOUT": "0xC803 (Dept: 50, Cmd: 3)",
    "PROTO_NC_SYSLOG_SERVER_CCU": "0xC811 (Dept: 50, Cmd: 17)",
}

def generate_doc(pdb_path, output_md, service_title, service_header):
    print(f"Generating {output_md} from {pdb_path}...")
    parser = PDBParser(pdb_path)
    proto_structs = [s for s in parser.structs if "tagNC_" in s or "PROTO_NC_" in s]
    
    lines = []
    lines.append(f"# Especificación Oficial: `{service_title}`\n")
    lines.append(service_header)
    lines.append("\n## Catálogo Exhaustivo de Estructuras de Paquetes\n")
    lines.append(f"> **Total de Estructuras Catalogadas:** `{len(proto_structs)}`\n")
    lines.append("> Todas las estructuras han sido extraídas directamente del flujo TPI del PDB oficial con sus tamaños, tipos primitivos, offsets y variables exactas.\n")

    for sname in sorted(proto_structs):
        sinfo = parser.structs[sname]
        members = parser.get_members(sname)
        # Find opcode
        opcode_info = ""
        # Match plain or strip template/subclass
        clean_name = sname.split("<")[0].split("::")[-1]
        if sname in OPCODE_MAP: opcode_info = OPCODE_MAP[sname]
        elif clean_name in OPCODE_MAP: opcode_info = OPCODE_MAP[clean_name]
        
        md_text = format_struct_markdown(sname, sinfo, members, opcode_info)
        lines.append(md_text)

    with open(output_md, "w", encoding="utf-8") as f:
        f.write("\n".join(lines))
    print(f"Successfully generated {output_md} ({len(proto_structs)} structs).")

LOGIN_HEADER = """
## Ficha Técnica del Binario
- **Nombre de Binario:** `Login.exe`
- **Símbolos:** `Login.pdb` (Microsoft C/C++ MSF 7.00, TPI v4.10)
- **Ruta del Proyecto Original:** `E:\\projectf2\\server\\3loginserver2\\loginserver\\`
- **Base Address en Memoria:** `0x00400000`
- **Conexión SQL / ODBC:** **NINGUNA (0 conexiones directas a SQL)**.
  - `Login.exe` **no** importa `ODBC32.dll`.
  - Todas las operaciones de autenticación de cuenta y registro de logs se delegan a través de sockets S2S (`Account.exe` y `AccountLog.exe`).
- **Archivos Locales Leídos:**
  - `ClientVersionKeyInfo.txt`: Archivo de texto de 14 bytes (ej. `"10022024000000"`) leído en tiempo de ejecución para validar versiones en `fc_NC_USER_CLIENT_VERSION_CHECK_REQ`.
  - **Archivos SHN:** **0 archivos**. `Login.exe` no carga ni verifica SHNs de gameplay.
- **Puertos de Red:**
  - **9010 (TCP):** Puerto C2S (Client to Server) donde se conectan los clientes del juego.
  - **9011 (TCP):** Puerto S2S (Server to Server) donde conectan `WorldManager` y herramientas administrativas (`OPTool`).
  - **9031 (TCP Cliente):** Conexión saliente hacia `Account.exe` (`CLoginAccountDBSession`).
  - **9032 (TCP Cliente):** Conexión saliente hacia `AccountLog.exe` (`CLoginAccountLogDBSession`).

## Dispatcher de Protocolos (OpCodes y Handlers)
El protocolo Fiesta utiliza un entero de 16 bits (`NETCOMMAND`) compuesto por un bitfield:
- `department`: **6 bits** (Head / Categoría del sistema, valores 0-63).
- `command`: **10 bits** (Tail / Comando específico, valores 0-1023).
- Fórmula de cálculo: `protocol = (department << 10) | command`.

| OpCode (Hex) | OpCode (Dec) | Dept (Head) | Cmd (Tail) | Función de Callback (Handler) | Dirección | Propósito |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `0x0801` | 2049 | 2 | 1 | `fc_NC_MISC_S2SCONNECTION_RDY` | S2S | Handshake inicial S2S listo |
| `0x0802` | 2050 | 2 | 2 | `fc_NC_MISC_S2SCONNECTION_REQ` | S2S | Solicitud de enlace S2S |
| `0x0803` | 2051 | 2 | 3 | `fc_NC_MISC_S2SCONNECTION_ACK` | S2S | Confirmación de enlace S2S |
| `0x0804` | 2052 | 2 | 4 | `fc_NC_MISC_HEARTBEAT_REQ` | C2S / S2S | Keep-alive heartbeat request |
| `0x0805` | 2053 | 2 | 5 | `fc_NC_MISC_HEARTBEAT_ACK` | S2C / S2S | Keep-alive heartbeat response |
| `0x0808` | 2056 | 2 | 8 | `fc_NC_MISC_USER_COUNT_CMD` | S2S | Actualización de CCU (usuarios conectados) |
| `0x080B` | 2059 | 2 | 11 | `fc_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD` | S2S | Entrega de mensaje WM -> Login -> Account |
| `0x080C` | 2060 | 2 | 12 | `fc_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD` | S2S | Entrega de mensaje WM -> Login -> AccountLog |
| `0x0823` | 2083 | 2 | 35 | `fc_NC_MISC_START_THE_BOOM_CMD` | S2S | Notificación de reinicio/apagado inmediato |
| `0x0C02` | 3074 | 3 | 2 | *(S2C)* | S2C | Notificación de versión incorrecta al cliente |
| `0x0C04` | 3076 | 3 | 4 | `fc_NC_USER_XTRAP_REQ` | C2S | Verificación de seguridad XTrap |
| `0x0C06` | 3078 | 3 | 6 | `fc_NC_USER_LOGIN_REQ` | C2S | Solicitud de Login estándar (User + Pass) |
| `0x0C08` | 3080 | 3 | 8 | `fc_NC_USER_PASSWORD_CHECK_ACK` | S2S (In) | Respuesta de AccountDB validando credenciales |
| `0x0C09` | 3081 | 3 | 9 | *(S2C)* | S2C | Notificación de fallo de inicio de sesión |
| `0x0C0A` | 3082 | 3 | 10 | *(S2C)* | S2C | Login exitoso + Lista de mundos disponibles |
| `0x0C0B` | 3083 | 3 | 11 | `fc_NC_USER_WORLDSELECT_REQ` | C2S | Selección de servidor / mundo por el cliente |
| `0x0C0C` | 3084 | 3 | 12 | *(S2C)* | S2C | Credenciales y Token de traspaso a WorldServer |
| `0x0C0E` | 3086 | 3 | 14 | `fc_NC_USER_WILLLOGIN_ACK` | S2S (In) | Confirmación de WorldManager para ingreso de jugador |
| `0x0C18` | 3096 | 3 | 24 | `fc_NC_USER_NORMALLOGOUT_CMD` | C2S | Desconexión normal solicitada por cliente |
| `0x0C1B` | 3099 | 3 | 27 | `fc_NC_USER_WORLD_STATUS_REQ` | C2S | Consulta de estado de mundos |
| `0x0C1D` | 3101 | 3 | 29 | `fc_NC_USER_LOGIN_NETMARBLE_REQ` | C2S | Login para plataforma Netmarble (Corea) |
| `0x0C20` | 3104 | 3 | 32 | `fc_NC_USER_LOGIN_OUTSPARK_REQ` | C2S | Login con Token de Outspark |
| `0x0C27` | 3111 | 3 | 39 | `fc_NC_USER_TEENAGER_ACK` | S2S (In) | Restricción horaria para menores (Korea/China) |
| `0x0C2B` | 3115 | 3 | 43 | `fc_NC_USER_IS_IP_BLOCK_ACK` | S2S (In) | Respuesta de AccountLog con estado de bloqueo IP |
| `0x0C35` | 3125 | 3 | 53 | `fc_NC_USER_CREATE_OTP_REQ` | S2S | Creación de token One-Time Password |
| `0x0C37` | 3127 | 3 | 55 | `fc_NC_USER_LOGIN_WITH_OTP_REQ` | C2S | Login autenticado vía OTP |
| `0x0C3A` | 3130 | 3 | 58 | `fc_NC_USER_GER_PASSWORD_CHECK_ACK` | S2S (In) | Respuesta de validación región Alemania |
| `0x0C3C` | 3132 | 3 | 60 | `fc_NC_USER_GER_IS_IP_BLOCK_ACK` | S2S (In) | Respuesta bloqueo IP región Alemania |
| `0x0C3F` | 3135 | 3 | 63 | `fc_NC_USER_TW_PASSWORD_CHECK_ACK` | S2S (In) | Respuesta validación región Taiwán |
| `0x0C41` | 3137 | 3 | 65 | `fc_NC_USER_TW_IS_IP_BLOCK_ACK` | S2S (In) | Respuesta bloqueo IP región Taiwán |
| `0x0C4C` | 3148 | 3 | 76 | `fc_NC_USER_JP_PASSWORD_CHECK_ACK` | S2S (In) | Respuesta validación región Japón |
| `0x0C4E` | 3150 | 3 | 78 | `fc_NC_USER_JP_IS_IP_BLOCK_ACK` | S2S (In) | Respuesta bloqueo IP región Japón |
| `0x0C5A` | 3162 | 3 | 90 | `fc_NC_USER_US_LOGIN_REQ` | C2S | Login oficial NA/US (316 bytes payload) |
| `0x0C5C` | 3164 | 3 | 92 | `fc_NC_USER_US_PASSWORD_CHECK_ACK` | S2S (In) | Respuesta autenticación NA/US desde Account |
| `0x0C5E` | 3166 | 3 | 94 | `fc_NC_USER_US_IS_IP_BLOCK_ACK` | S2S (In) | Respuesta verificación IP NA/US desde AccountLog |
| `0x0C65` | 3173 | 3 | 101 | `fc_NC_USER_CLIENT_VERSION_CHECK_REQ` | C2S | Verificación de versión del cliente |
| `0x0C67` | 3175 | 3 | 103 | *(S2C)* | S2C | Versión correcta aprobada (con semilla XTrap) |
| `0x2801` | 10241 | 10 | 1 | `fc_NC_OPTOOL_S2SCONNECT_LIST_REQ` | S2S | OPTool: Lista de conexiones activas |
| `0x2803` | 10243 | 10 | 3 | `fc_NC_OPTOOL_CLOSE_SERVER_REQ` | S2S | OPTool: Comando de apagado de servidor |
| `0x2805` | 10245 | 10 | 5 | `fc_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_CLR` | S2S | OPTool: Limpieza de métricas de logon time |
| `0x2807` | 10247 | 10 | 7 | `fc_NC_OPTOOL_CONNECT_BRIF_REQ` | S2S | OPTool: Resumen de conexiones |
| `0x2812` | 10258 | 10 | 18 | `fc_NC_OPTOOL_LOGIN_USER_RATABLE_GET_REQ` | S2S | OPTool: Consulta de ratio de usuarios permitidos |
| `0x2814` | 10260 | 10 | 20 | `fc_NC_OPTOOL_LOGIN_USER_RATABLE_SET_CMD` | S2S | OPTool: Modificación de ratio de usuarios |
| `0x2818` | 10264 | 10 | 24 | `fc_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_REQ` | S2S | OPTool: Consulta de métricas de logon |
"""

ACCOUNT_HEADER = """
## Ficha Técnica del Binario
- **Nombre de Binario:** `Account.exe`
- **Símbolos:** `Account.pdb` (Microsoft C/C++ MSF 7.00, TPI v4.10)
- **Ruta del Proyecto Original:** `E:\\projectf2\\server\\3loginserver2\\accountserver\\`
- **Base Address en Memoria:** `0x00400000`
- **Conexión SQL / ODBC:** **ACTIVA (`ODBC32.dll`)**.
  - Conecta a la base de datos `Account`.
  - Origen de datos ODBC configurado para ejecutar Stored Procedures de autenticación, alta de usuarios y validación de bloques.
- **Archivos Locales y SHN Leídos:**
  - Compila `shndatafilechecksum.cpp`: Calcula y verifica en memoria los checksums de 49 a 53 archivos `.shn` de cliente (ej. `ItemInfo.shn`, `ActionViewInfo.shn`, `UpgradeInfo.shn`, `MapInfo.shn`) para detectar modificaciones no autorizadas en el cliente.
- **Puertos de Red:**
  - **9031 (TCP):** Puerto de escucha S2S (`PG_AccDB`). Recibe conexiones de `Login.exe` y de los distintos servidores del cluster.

## Procedimientos Almacenados SQL Invocados
Todas las consultas ejecutadas por `Account.exe` están resueltas en los scripts SQL originales recuperados en `docs/sql/procedures/Account_all_procedures.sql`:

| Stored Procedure | Parámetros (ODBC CALL) | Propósito Oficial |
| :--- | :--- | :--- |
| `usp_User_USLoginGame` | `{CALL usp_User_USLoginGame(?,?,?,?,?,?)}` | Login para la versión US: verifica usuario, hash de contraseña, IP y devuelve UserNo, Block, AuthID |
| `usp_User_loginGame` | `{CALL usp_User_loginGame(?,?,?,?,?,?)}` | Login estándar / internacional con verificación de clave y estado de cuenta |
| `usp_User_LoginGame` | `{CALL usp_User_LoginGame(?,?,?,?,?)}` | Variante histórica de autenticación de usuario |
| `usp_User_LoginNetmarble` | `{CALL usp_User_LoginNetmarble(?,?,?,?,?)}` | Autenticación delegada para plataforma Netmarble |
| `usp_Pcbang_chkIP` | `{CALL usp_Pcbang_chkIP(?,?)}` | Verificación de IP de cibercafé / PC Bang |
| `usp_ReturnUser_GetDatediffLastlog` | `{CALL usp_ReturnUser_GetDatediffLastlog(?,?)}` | Cálculo de días desde el último login para recompensas de retorno |
| `usp_ReturnUser_SetUserItem` | `{CALL usp_ReturnUser_SetUserItem(?,?)}` | Asignación de paquete de ítems para jugadores que retornan |
| `usp_User_GetDatediffCreated` | `{CALL usp_User_GetDatediffCreated(?,?)}` | Días transcurridos desde la creación de la cuenta |
| `usp_Event_addCash` | `{CALL usp_Event_addCash(?,?,?)}` | Acreditación de moneda premium (Cash/Spark) por evento |
| `usp_Item_setSerialKey` | `{CALL usp_Item_setSerialKey(?,?,?)}` | Canje de códigos de serie y cupones promocionales |
| `p_UserKey_Get_FromAccount` | `{CALL p_UserKey_Get_FromAccount(?,?,?,?)}` | Obtención de claves y tokens de cuenta |
| `p_UserPlayRemainTime_Get_FromAccount`| `{CALL p_UserPlayRemainTime_Get_FromAccount(?,?,?,?)}` | Control parental / tiempo restante de juego disponible |
| `p_UserPlayRemainTime_Set_FromAccount`| `{CALL p_UserPlayRemainTime_Set_FromAccount(?,?,?,?)}` | Descuento del tiempo de juego consumido |

## Dispatcher de Protocolos (OpCodes y Handlers)

| OpCode (Hex) | OpCode (Dec) | Dept (Head) | Cmd (Tail) | Handler | Propósito |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `0x0801` | 2049 | 2 | 1 | `fc_NC_MISC_HEARTBEAT_ACK` | S2S Heartbeat Ack |
| `0x0802` | 2050 | 2 | 2 | `fc_NC_MISC_S2SCONNECTION_REQ` | S2S Solicitud de Enlace |
| `0x0803` | 2051 | 2 | 3 | `fc_NC_MISC_S2SCONNECTION_ACK` | S2S Aceptación de Enlace |
| `0x0804` | 2052 | 2 | 4 | `fc_NC_MISC_HEARTBEAT_REQ` | S2S Heartbeat Request |
| `0x081D` | 2077 | 2 | 29 | `fc_NC_MISC_WEB_DB_KEY_REQ` | Consulta de token generado vía Web |
| `0x0820` | 2080 | 2 | 32 | `fc_NC_MISC_EVENT_HIT3_ADD_CASH` | Adición de Cash por evento Hit3 |
| `0x082F` | 2095 | 2 | 47 | `fc_NC_MISC_CONNECTFROMWHERE_DB_REQ` | Consulta de región de origen de conexión |
| `0x0838` | 2104 | 2 | 56 | `fc_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC`| Notificación de evento de bono cultural Munsang |
| `0x0C07` | 3079 | 3 | 7 | `fc_NC_USER_PASSWORD_CHECK_REQ` | Verificación de credenciales (Login -> Account) |
| `0x0C1E` | 3102 | 3 | 30 | `fc_NC_USER_LOGIN_NETMARBLE_DB_REQ` | Verificación de credenciales Netmarble |
| `0x0C26` | 3110 | 3 | 38 | `fc_NC_USER_TEENAGER_REQ` | Consulta de política de protección de menores |
| `0x0C28` | 3112 | 3 | 40 | `fc_NC_USER_TEENAGER_SET_CMD` | Aplicación de restricción de tiempo a menores |
| `0x0C39` | 3129 | 3 | 57 | `fc_NC_USER_GER_PASSWORD_CHECK_REQ` | Verificación credenciales Alemania |
| `0x0C3E` | 3134 | 3 | 62 | `fc_NC_USER_TW_PASSWORD_CHECK_REQ` | Verificación credenciales Taiwán |
| `0x0C4B` | 3147 | 3 | 75 | `fc_NC_USER_JP_PASSWORD_CHECK_REQ` | Verificación credenciales Japón |
| `0x0C53` | 3155 | 3 | 83 | `fc_NC_USER_CH_PASSWORD_CHECK_REQ` | Verificación credenciales China |
| `0x0C5B` | 3163 | 3 | 91 | `fc_NC_USER_US_PASSWORD_CHECK_REQ` | Verificación credenciales Norteamérica (US) |
| `0x0C60` | 3168 | 3 | 96 | `fc_NC_USER_RETURN_CHECK_REQ` | Consulta de condición de usuario retornante |
| `0x0C64` | 3172 | 3 | 100 | `fc_NC_USER_SET_RETURN_CMD` | Marca cuenta como usuario retornante |
| `0x2801` | 10241 | 10 | 1 | `fc_NC_OPTOOL_S2SCONNECT_LIST_REQ` | OPTool: Lista conexiones |
| `0x2803` | 10243 | 10 | 3 | `fc_NC_OPTOOL_CLOSE_SERVER_REQ` | OPTool: Cierre de servidor |
| `0x2805` | 10245 | 10 | 5 | `fc_NC_OPTOOL_MAP_USER_LIST_REQ` | OPTool: Lista de usuarios por mapa |
| `0x2807` | 10247 | 10 | 7 | `fc_NC_OPTOOL_CONNECT_BRIF_REQ` | OPTool: Resumen de estado |
"""

ACCOUNTLOG_HEADER = """
## Ficha Técnica del Binario
- **Nombre de Binario:** `AccountLog.exe`
- **Símbolos:** `AccountLog.pdb` (Microsoft C/C++ MSF 7.00, TPI v4.10)
- **Ruta del Proyecto Original:** `E:\\projectf2\\server\\3loginserver2\\accountlogserver\\`
- **Base Address en Memoria:** `0x00400000`
- **Conexión SQL / ODBC:** **ACTIVA (`ODBC32.dll`)**.
  - Conecta a la base de datos `AccountLog`.
  - Ejecuta Stored Procedures para registrar inicios de sesión, fallos de autenticación, cierres de sesión, métricas de CCU y validación de listas de IPs bloqueadas.
- **Archivos Locales y SHN Leídos:**
  - **0 archivos SHN**. `AccountLog.exe` no interactúa con datos estáticos de juego.
- **Puertos de Red:**
  - **9032 (TCP):** Puerto de escucha S2S (`PG_ALDB`). Recibe conexiones de `Login.exe` y servidores del cluster para auditoría y logs.

## Procedimientos Almacenados SQL Invocados
Todas las consultas ejecutadas por `AccountLog.exe` están catalogadas y recuperadas en `docs/sql/procedures/AccountLog_all_procedures.sql`:

| Stored Procedure | Parámetros (ODBC CALL) | Propósito Oficial |
| :--- | :--- | :--- |
| `p_AccountLog_Add` | `{CALL p_AccountLog_Add(?,?,?,?,?,?,?,?,?)}` | Inserción de registro de login exitoso (UserNo, IP, MAC, World, Timestamp) |
| `p_IPBlock_IsBlock` | `{CALL p_IPBlock_IsBlock(?,?,?)}` | Comprobación de si la IP del cliente se encuentra en la tabla de baneo/bloqueo |
| `usp_Log_setLogout` | `{CALL usp_Log_setLogout(?,?,?,?,?,?,?,?)}` | Registro de desconexión / logout con cálculo de tiempo jugado |
| `p_Report_Add` | `{CALL p_Report_Add(?,?,?,?,?,?,?,?)}` | Alta de reporte de usuario / denuncia de jugador |
| `p_Report_Del` | `{CALL p_Report_Del(?,?)}` | Eliminación de reporte procesado |
| `p_Report_Get` | `{CALL p_Report_Get(?,?)}` | Consulta de reportes pendientes |
| `p_SameTimeChar_Add` | `{CALL p_SameTimeChar_Add(?,?,?)}` | Registro de métrica de concurrentes totales (CCU) en el servidor |
| `p_SameTimeCharClass_Add` | `{CALL p_SameTimeCharClass_Add(?,?,?,?)}` | Desglose de concurrentes por clase de personaje |
| `p_Event_L20_Get` | `{CALL p_Event_L20_Get(?,?,?,?,?,?)}` | Auditoría de evento especial nivel 20 |
| `p_Event_L20_Set` | `{CALL p_Event_L20_Set(?,?,?,?,?,?)}` | Registro de recompensa entregada de evento nivel 20 |
| `usp_AccountLog_GetDateDiffLastlogout` | `{CALL usp_AccountLog_GetDateDiffLastlogout(?,?)}` | Consulta de tiempo transcurrido desde el último deslogueo |
| `usp_EventiPad_setTicket` | `{CALL usp_EventiPad_setTicket(?,?,?)}` | Registro de boletos de sorteo de evento promocional |

## Dispatcher de Protocolos (OpCodes y Handlers)

| OpCode (Hex) | OpCode (Dec) | Dept (Head) | Cmd (Tail) | Handler | Propósito |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `0x0720` | 1824 | 1 | 800 | `fc_NC_LOG_USER_LOGINFAIL` | Registro en base de datos de intento de login fallido |
| `0x0721` | 1825 | 1 | 801 | `fc_NC_LOG_USER_LOGIN` | Registro de login exitoso |
| `0x0722` | 1826 | 1 | 802 | `fc_NC_LOG_USER_LOGOUT` | Registro de deslogueo / fin de sesión |
| `0x0801` | 2049 | 2 | 1 | `fc_NC_MISC_HEARTBEAT_ACK` | S2S Heartbeat Ack |
| `0x0802` | 2050 | 2 | 2 | `fc_NC_MISC_S2SCONNECTION_REQ` | S2S Solicitud de Enlace |
| `0x0803` | 2051 | 2 | 3 | `fc_NC_MISC_S2SCONNECTION_ACK` | S2S Aceptación de Enlace |
| `0x0804` | 2052 | 2 | 4 | `fc_NC_MISC_HEARTBEAT_REQ` | S2S Heartbeat Request |
| `0x080A` | 2058 | 2 | 10 | `fc_NC_MISC_CONNECTER_ACK` | Confirmación de conector de red |
| `0x081F` | 2079 | 2 | 31 | `fc_NC_MISC_CHAR_LOGOFF_STATISTICS` | Estadísticas de logoff de personaje |
| `0x0823` | 2083 | 2 | 35 | `fc_NC_MISC_START_THE_BOOM_CMD` | Notificación de reinicio forzado |
| `0x082C` | 2092 | 2 | 44 | `fc_NC_MISC_EVENT_L20_DB_REQ` | Auditoría de evento nivel 20 |
| `0x0835` | 2101 | 2 | 53 | `fc_NC_MISC_HACK_SCAN_STORE_DB_CMD` | Almacenamiento de logs de escaneo anti-hack |
| `0x0843` | 2115 | 2 | 67 | `fc_NC_MISC_EVENT_LONG_TIME_PLAY_CMD` | Log de jugador con sesión prolongada |
| `0x0C2A` | 3114 | 3 | 42 | `fc_NC_USER_IS_IP_BLOCK_REQ` | Solicitud de comprobación de baneo de IP |
| `0x0C3B` | 3131 | 3 | 59 | `fc_NC_USER_GER_IS_IP_BLOCK_REQ` | Comprobación IP Alemania |
| `0x0C40` | 3136 | 3 | 64 | `fc_NC_USER_TW_IS_IP_BLOCK_REQ` | Comprobación IP Taiwán |
| `0x0C4D` | 3149 | 3 | 77 | `fc_NC_USER_JP_IS_IP_BLOCK_REQ` | Comprobación IP Japón |
| `0x0C5D` | 3165 | 3 | 93 | `fc_NC_USER_US_IS_IP_BLOCK_REQ` | Comprobación IP Estados Unidos (US) |
| `0x0C62` | 3170 | 3 | 98 | `fc_NC_USER_LOGOUT_LAST_DAY_REQ` | Consulta de último día de login |
| `0x2801` | 10241 | 10 | 1 | `fc_NC_OPTOOL_S2SCONNECT_LIST_REQ` | OPTool: Lista conexiones |
| `0x2803` | 10243 | 10 | 3 | `fc_NC_OPTOOL_CLOSE_SERVER_REQ` | OPTool: Cierre de servidor |
| `0x2805` | 10245 | 10 | 5 | `fc_NC_OPTOOL_MAP_USER_LIST_REQ` | OPTool: Lista de usuarios por mapa |
| `0x2807` | 10247 | 10 | 7 | `fc_NC_OPTOOL_CONNECT_BRIF_REQ` | OPTool: Resumen de estado |
| `0x8007` | 32775 | 32 | 7 | `fc_NC_DATA_REPORT_ADD_REQ` | Reporte: Alta de denuncia |
| `0x8009` | 32777 | 32 | 9 | `fc_NC_DATA_REPORT_CANCEL_REQ` | Reporte: Cancelación de denuncia |
| `0x800B` | 32779 | 32 | 11 | `fc_NC_DATA_REPORT_GET_REQ` | Reporte: Consulta de reportes |
| `0xC801` | 51201 | 50 | 1 | `fc_NC_SYSLOG_ACCOUNT_LOGIN_SUCCESS`| Syslog: Éxito de login |
| `0xC802` | 51202 | 50 | 2 | `fc_NC_SYSLOG_ACCOUNT_LOGIN_FAILURE`| Syslog: Fallo de login |
| `0xC803` | 51203 | 50 | 3 | `fc_NC_SYSLOG_ACCOUNT_LOGOUT` | Syslog: Desconexión |
| `0xC811` | 51217 | 50 | 17 | `fc_NC_SYSLOG_SERVER_CCU` | Syslog: Métrica de concurrentes |
"""

if __name__ == "__main__":
    generate_doc("/Users/fedemac/Desktop/Xelion Online/AstraServer/Server/Login/Login.pdb",
                 "docs/Servicios/Login.md",
                 "Login (LoginServer)",
                 LOGIN_HEADER)

    generate_doc("/Users/fedemac/Desktop/Xelion Online/AstraServer/Server/Account/Account.pdb",
                 "docs/Servicios/Account.md",
                 "Account (AccountServer)",
                 ACCOUNT_HEADER)

    generate_doc("/Users/fedemac/Desktop/Xelion Online/AstraServer/Server/AccountLog/AccountLog.pdb",
                 "docs/Servicios/AccountLog.md",
                 "AccountLog (AccountLogServer)",
                 ACCOUNTLOG_HEADER)
