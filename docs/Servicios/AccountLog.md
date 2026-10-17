# Especificación Oficial: `AccountLog (AccountLogServer)`


## Ficha Técnica del Binario
- **Nombre de Binario:** `AccountLog.exe`
- **Símbolos:** `AccountLog.pdb` (Microsoft C/C++ MSF 7.00, TPI v4.10)
- **Ruta del Proyecto Original:** `E:\projectf2\server\3loginserver2\accountlogserver\`
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


## Catálogo Exhaustivo de Estructuras de Paquetes

> **Total de Estructuras Catalogadas:** `106`

> Todas las estructuras han sido extraídas directamente del flujo TPI del PDB oficial con sus tamaños, tipos primitivos, offsets y variables exactas.

### `PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD`

- **Tamaño Total en Memoria:** `3 bytes` (`0x3`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD {
    unsigned short handle;
    unsigned char count;
    ABSTATE_INFORMATION infoList[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `handle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned char` | `count` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `ABSTATE_INFORMATION[0]` | `infoList` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD`

- **Tamaño Total en Memoria:** `5 bytes` (`0x5`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD {
    unsigned short handle;
    unsigned short item;
    unsigned char nSlotNum;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `handle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short` | `item` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `nSlotNum` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD`

- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD {
    unsigned short handle;
    unsigned short item;
    unsigned char upgrade;
    unsigned char nSlotNum;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `handle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short` | `item` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `upgrade` | Campo de paquete oficial extraído de PDB TPI |
| `0x005` | `5` | `unsigned char` | `nSlotNum` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD`

- **Tamaño Total en Memoria:** `9 bytes` (`0x9`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD {
    PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD upgradeinfo;
    unsigned short currentmobid;
    unsigned char currentkilllevel;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD` | `upgradeinfo` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned short` | `currentmobid` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char` | `currentkilllevel` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_BRIEFINFO_UNEQUIP_CMD`

- **Tamaño Total en Memoria:** `3 bytes` (`0x3`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_BRIEFINFO_UNEQUIP_CMD {
    unsigned short handle;
    unsigned char slot;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `handle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned char` | `slot` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHARSAVE_ITEMACTIONCOOLTIME_ACK`

- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHARSAVE_ITEMACTIONCOOLTIME_ACK {
    unsigned short num;
    ActionCooltime group[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `num` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `ActionCooltime[0]` | `group` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_ABSTATE_CMD`

- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_ABSTATE_CMD {
    unsigned int chrregnum;
    unsigned short number;
    ABSTATEREADBLOCK abstate[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned short` | `number` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `ABSTATEREADBLOCK[0]` | `abstate` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_ARENA_CMD`

- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_ARENA_CMD {
    // Estructura vacía / Payload de 0 bytes
};
#pragma pack(pop)
```

> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*

### `PROTO_NC_CHAR_BASE_CMD`

- **Tamaño Total en Memoria:** `105 bytes` (`0x69`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_BASE_CMD {
    unsigned int chrregnum;
    Name5 charid;
    unsigned char slotno;
    unsigned char Level;
    unsigned __int64 Experience;
    unsigned short CurPwrStone;
    unsigned short CurGrdStone;
    unsigned short CurHPStone;
    unsigned short CurSPStone;
    unsigned int CurHP;
    unsigned int CurSP;
    unsigned int CurLP;
    unsigned int fame;
    unsigned __int64 Cen;
    PROTO_NC_CHAR_BASE_CMD::LoginLocation logininfo;
    CHARSTATDISTSTR statdistribute;
    unsigned char pkyellowtime;
    unsigned int pkcount;
    unsigned short prisonmin;
    unsigned char adminlevel;
    PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags> flags;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `Name5` | `charid` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned char` | `slotno` | Campo de paquete oficial extraído de PDB TPI |
| `0x019` | `25` | `unsigned char` | `Level` | Campo de paquete oficial extraído de PDB TPI |
| `0x01A` | `26` | `unsigned __int64` | `Experience` | Campo de paquete oficial extraído de PDB TPI |
| `0x022` | `34` | `unsigned short` | `CurPwrStone` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned short` | `CurGrdStone` | Campo de paquete oficial extraído de PDB TPI |
| `0x026` | `38` | `unsigned short` | `CurHPStone` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `unsigned short` | `CurSPStone` | Campo de paquete oficial extraído de PDB TPI |
| `0x02A` | `42` | `unsigned int` | `CurHP` | Campo de paquete oficial extraído de PDB TPI |
| `0x02E` | `46` | `unsigned int` | `CurSP` | Campo de paquete oficial extraído de PDB TPI |
| `0x032` | `50` | `unsigned int` | `CurLP` | Campo de paquete oficial extraído de PDB TPI |
| `0x036` | `54` | `unsigned int` | `fame` | Campo de paquete oficial extraído de PDB TPI |
| `0x03A` | `58` | `unsigned __int64` | `Cen` | Campo de paquete oficial extraído de PDB TPI |
| `0x042` | `66` | `PROTO_NC_CHAR_BASE_CMD::LoginLocation` | `logininfo` | Campo de paquete oficial extraído de PDB TPI |
| `0x057` | `87` | `CHARSTATDISTSTR` | `statdistribute` | Campo de paquete oficial extraído de PDB TPI |
| `0x05D` | `93` | `unsigned char` | `pkyellowtime` | Campo de paquete oficial extraído de PDB TPI |
| `0x05E` | `94` | `unsigned int` | `pkcount` | Campo de paquete oficial extraído de PDB TPI |
| `0x062` | `98` | `unsigned short` | `prisonmin` | Campo de paquete oficial extraído de PDB TPI |
| `0x064` | `100` | `unsigned char` | `adminlevel` | Campo de paquete oficial extraído de PDB TPI |
| `0x065` | `101` | `PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>` | `flags` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>`

- **Tamaño Total en Memoria:** `4 bytes` (`0x4`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags> {
    unsigned int bin;
    PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>::<unnamed-type-str> str;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `bin` | Campo de paquete oficial extraído de PDB TPI |
| `0x000` | `0` | `PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>::<unnamed-type-str>` | `str` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>::<unnamed-type-str>`

- **Tamaño Total en Memoria:** `4 bytes` (`0x4`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>::<unnamed-type-str> {
    type_0x1205 skillempower_can_reset;
    type_0x1205 restunion;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `type_0x1205` | `skillempower_can_reset` | Campo de paquete oficial extraído de PDB TPI |
| `0x000` | `0` | `type_0x1205` | `restunion` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_BASE_CMD::LoginLocation`

- **Tamaño Total en Memoria:** `21 bytes` (`0x15`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_BASE_CMD::LoginLocation {
    Name3 currentmap;
    SHINE_COORD_TYPE currentcoord;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `Name3` | `currentmap` | Campo de paquete oficial extraído de PDB TPI |
| `0x00C` | `12` | `SHINE_COORD_TYPE` | `currentcoord` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_CHARGEDBUFF_CMD`

- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_CHARGEDBUFF_CMD {
    unsigned short NumOfChargedBuff;
    PROTO_CHARGEDBUFF_INFO ChargedBuff[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `NumOfChargedBuff` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `PROTO_CHARGEDBUFF_INFO[0]` | `ChargedBuff` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_CHESTINFO_CMD`

- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_CHESTINFO_CMD {
    unsigned short itemnum;
    unsigned short location[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `itemnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short[0]` | `location` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_EMBLEM_CMD`

- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_EMBLEM_CMD {
    // Estructura vacía / Payload de 0 bytes
};
#pragma pack(pop)
```

> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*

### `PROTO_NC_CHAR_FRIEND_CMD`

- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_FRIEND_CMD {
    unsigned short numberoffriend;
    PROTO_NC_CHAR_FRIEND_CMD_INFO friend_array[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `numberoffriend` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `PROTO_NC_CHAR_FRIEND_CMD_INFO[0]` | `friend_array` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_FRIEND_CMD_INFO`

- **Tamaño Total en Memoria:** `25 bytes` (`0x19`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_FRIEND_CMD_INFO {
    unsigned int charno;
    Name5 charid;
    unsigned char flag;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `charno` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `Name5` | `charid` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned char` | `flag` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_GUILD_CMD`

- **Tamaño Total en Memoria:** `4 bytes` (`0x4`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_GUILD_CMD {
    unsigned int nGuildNo;
    GUILD_CLIENT Guild[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nGuildNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `GUILD_CLIENT[0]` | `Guild` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_HOUSE_CMD`

- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_HOUSE_CMD {
    // Estructura vacía / Payload de 0 bytes
};
#pragma pack(pop)
```

> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*

### `PROTO_NC_CHAR_MASPUP_CMD`

- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_MASPUP_CMD {
    // Estructura vacía / Payload de 0 bytes
};
#pragma pack(pop)
```

> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*

### `PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_LIST_CMD`

- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_LIST_CMD {
    unsigned char nGuideViewNum;
    unsigned char nGuideViewList[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `nGuideViewNum` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned char[0]` | `nGuideViewList` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_OPTION_GAME`

- **Tamaño Total en Memoria:** `64 bytes` (`0x40`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_OPTION_GAME {
    unsigned char Data[64];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[64]` | `Data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_OPTION_KEYMAPPING`

- **Tamaño Total en Memoria:** `308 bytes` (`0x134`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_OPTION_KEYMAPPING {
    unsigned char Data[308];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[308]` | `Data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_OPTION_SHORTCUTDATA`

- **Tamaño Total en Memoria:** `1024 bytes` (`0x400`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_OPTION_SHORTCUTDATA {
    unsigned char Data[1024];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[1024]` | `Data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_OPTION_SHORTCUTSIZE`

- **Tamaño Total en Memoria:** `24 bytes` (`0x18`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_OPTION_SHORTCUTSIZE {
    unsigned char Data[24];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[24]` | `Data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_OPTION_SOUND`

- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_OPTION_SOUND {
    unsigned char Data[1];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[1]` | `Data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_OPTION_VIDEO`

- **Tamaño Total en Memoria:** `60 bytes` (`0x3C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_OPTION_VIDEO {
    unsigned char Data[60];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[60]` | `Data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_OPTION_WINDOWPOS`

- **Tamaño Total en Memoria:** `392 bytes` (`0x188`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_OPTION_WINDOWPOS {
    unsigned char Data[392];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[392]` | `Data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_PET_CMD`

- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_PET_CMD {
    // Estructura vacía / Payload de 0 bytes
};
#pragma pack(pop)
```

> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*

### `PROTO_NC_CHAR_QUEST_DOING_CMD`

- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_QUEST_DOING_CMD {
    unsigned int chrregnum;
    unsigned char bNeedClear;
    unsigned char nNumOfDoingQuest;
    PLAYER_QUEST_INFO QuestDoingArray[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `bNeedClear` | Campo de paquete oficial extraído de PDB TPI |
| `0x005` | `5` | `unsigned char` | `nNumOfDoingQuest` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `PLAYER_QUEST_INFO[0]` | `QuestDoingArray` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_QUEST_READ_CMD`

- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_QUEST_READ_CMD {
    unsigned int chrregnum;
    unsigned short nNumOfReadQuest;
    unsigned short QuestReadIDArray[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned short` | `nNumOfReadQuest` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned short[0]` | `QuestReadIDArray` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_QUEST_REPEAT_CMD`

- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_QUEST_REPEAT_CMD {
    unsigned int chrregnum;
    unsigned short nNumOfRepeatQuest;
    PLAYER_QUEST_INFO QuestRepeatArray[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned short` | `nNumOfRepeatQuest` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `PLAYER_QUEST_INFO[0]` | `QuestRepeatArray` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_SEAWAR_CMD`

- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_SEAWAR_CMD {
    // Estructura vacía / Payload de 0 bytes
};
#pragma pack(pop)
```

> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*

### `PROTO_NC_CHAR_SINGLE_OPTION_CMD`

- **Tamaño Total en Memoria:** `15 bytes` (`0xF`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_SINGLE_OPTION_CMD {
    unsigned int chrregnum;
    int resttime;
    unsigned char autopick;
    unsigned char mystery_vault;
    unsigned int nChatBan;
    unsigned char bFirstJoin;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `int` | `resttime` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char` | `autopick` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `unsigned char` | `mystery_vault` | Campo de paquete oficial extraído de PDB TPI |
| `0x00A` | `10` | `unsigned int` | `nChatBan` | Campo de paquete oficial extraído de PDB TPI |
| `0x00E` | `14` | `unsigned char` | `bFirstJoin` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_SKILL_CMD`

- **Tamaño Total en Memoria:** `9 bytes` (`0x9`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_SKILL_CMD {
    unsigned int chrregnum;
    PARTMARK PartMark;
    unsigned short nMaxNum;
    unsigned short number;
    PROTO_SKILLREADBLOCK skill[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `PARTMARK` | `PartMark` | Campo de paquete oficial extraído de PDB TPI |
| `0x005` | `5` | `unsigned short` | `nMaxNum` | Campo de paquete oficial extraído de PDB TPI |
| `0x007` | `7` | `unsigned short` | `number` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `PROTO_SKILLREADBLOCK[0]` | `skill` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_SKILL_PASSIVE_CMD`

- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_SKILL_PASSIVE_CMD {
    unsigned short number;
    unsigned short passive[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `number` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short[0]` | `passive` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_WEDDINGDATA_ACK`

- **Tamaño Total en Memoria:** `43 bytes` (`0x2B`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_WEDDINGDATA_ACK {
    NETPACKETZONEHEADER header;
    unsigned short Error;
    PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Status> Status;
    ShineDateTime When;
    ShineDateTime SubWhen;
    ShineDateTime ApplicationLimit;
    ShineDateTime DelayLimit;
    ShineDateTime DivorceLimit;
    ShineDateTime NewlyLimit;
    PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Partner> Partner;
    unsigned char WeddingCount;
    unsigned char RefusalCount;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETZONEHEADER` | `header` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned short` | `Error` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Status>` | `Status` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `ShineDateTime` | `When` | Campo de paquete oficial extraído de PDB TPI |
| `0x00D` | `13` | `ShineDateTime` | `SubWhen` | Campo de paquete oficial extraído de PDB TPI |
| `0x00D` | `13` | `ShineDateTime` | `ApplicationLimit` | Campo de paquete oficial extraído de PDB TPI |
| `0x00D` | `13` | `ShineDateTime` | `DelayLimit` | Campo de paquete oficial extraído de PDB TPI |
| `0x00D` | `13` | `ShineDateTime` | `DivorceLimit` | Campo de paquete oficial extraído de PDB TPI |
| `0x00D` | `13` | `ShineDateTime` | `NewlyLimit` | Campo de paquete oficial extraído de PDB TPI |
| `0x011` | `17` | `PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Partner>` | `Partner` | Campo de paquete oficial extraído de PDB TPI |
| `0x029` | `41` | `unsigned char` | `WeddingCount` | Campo de paquete oficial extraído de PDB TPI |
| `0x02A` | `42` | `unsigned char` | `RefusalCount` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Partner>`

- **Tamaño Total en Memoria:** `24 bytes` (`0x18`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Partner> {
    Name5 Name;
    unsigned int RegisterNumber;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `Name5` | `Name` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `RegisterNumber` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Status>`

- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_WEDDINGDATA_ACK::<unnamed-type-Status> {
    type_0x1205 Main;
    type_0x1205 Sub;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `type_0x1205` | `Main` | Campo de paquete oficial extraído de PDB TPI |
| `0x000` | `0` | `type_0x1205` | `Sub` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK`

- **Tamaño Total en Memoria:** `14 bytes` (`0xE`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK {
    NETPACKETZONEHEADER header;
    PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK::<unnamed-type-PartnerInfo> PartnerInfo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETZONEHEADER` | `header` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK::<unnamed-type-PartnerInfo>` | `PartnerInfo` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK::<unnamed-type-PartnerInfo>`

- **Tamaño Total en Memoria:** `8 bytes` (`0x8`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAR_WEDDING_PARTNER_INFO_ACK::<unnamed-type-PartnerInfo> {
    unsigned char Level;
    unsigned char Connected;
    ShineDateTime LastConnect;
    unsigned short CurrentMap;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `Level` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned char` | `Connected` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `ShineDateTime` | `LastConnect` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned short` | `CurrentMap` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_CHAT_RESTRICT_DB_LIST_CMD`

- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_CHAT_RESTRICT_DB_LIST_CMD {
    unsigned short nChatRestrictListCnt;
    PROTO_CHAT_RESTRICT_DB_INFO ChatRestrictData[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nChatRestrictListCnt` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `PROTO_CHAT_RESTRICT_DB_INFO[0]` | `ChatRestrictData` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_REPORT_ADD_ACK`

- **OpCode Asociado:** `0x8008 (Dept: 32, Cmd: 8)`
- **Tamaño Total en Memoria:** `28 bytes` (`0x1C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_REPORT_ADD_ACK {
    unsigned short err;
    unsigned short hReporter;
    Name5 sHarmer;
    int nRegNumber;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `err` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short` | `hReporter` | Puerto TCP de escucha |
| `0x004` | `4` | `Name5` | `sHarmer` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `int` | `nRegNumber` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_REPORT_ADD_REQ`

- **OpCode Asociado:** `0x8007 (Dept: 32, Cmd: 7)`
- **Tamaño Total en Memoria:** `2369 bytes` (`0x941`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_REPORT_ADD_REQ {
    unsigned int dwReporter;
    unsigned short hReporter;
    unsigned char nWorldNo;
    Name5 sReporterCharID;
    Name5 sHarmer;
    unsigned char byReportType[16];
    unsigned char byReason[256];
    unsigned short nChatLogSize;
    unsigned char byChatLog[2048];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `dwReporter` | Puerto TCP de escucha |
| `0x004` | `4` | `unsigned short` | `hReporter` | Puerto TCP de escucha |
| `0x006` | `6` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |
| `0x007` | `7` | `Name5` | `sReporterCharID` | Puerto TCP de escucha |
| `0x01B` | `27` | `Name5` | `sHarmer` | Campo de paquete oficial extraído de PDB TPI |
| `0x02F` | `47` | `unsigned char[16]` | `byReportType` | Puerto TCP de escucha |
| `0x03F` | `63` | `unsigned char[256]` | `byReason` | Campo de paquete oficial extraído de PDB TPI |
| `0x13F` | `319` | `unsigned short` | `nChatLogSize` | Campo de paquete oficial extraído de PDB TPI |
| `0x141` | `321` | `unsigned char[2048]` | `byChatLog` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_REPORT_CANCEL_ACK`

- **OpCode Asociado:** `0x800A (Dept: 32, Cmd: 10)`
- **Tamaño Total en Memoria:** `8 bytes` (`0x8`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_REPORT_CANCEL_ACK {
    int nRegNumber;
    unsigned short hReporter;
    unsigned short err;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `int` | `nRegNumber` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned short` | `hReporter` | Puerto TCP de escucha |
| `0x006` | `6` | `unsigned short` | `err` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_REPORT_CANCEL_REQ`

- **OpCode Asociado:** `0x8009 (Dept: 32, Cmd: 9)`
- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_REPORT_CANCEL_REQ {
    int nRegNumber;
    unsigned short hReporter;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `int` | `nRegNumber` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned short` | `hReporter` | Puerto TCP de escucha |

### `PROTO_NC_DATA_REPORT_GET_ACK`

- **OpCode Asociado:** `0x800C (Dept: 32, Cmd: 12)`
- **Tamaño Total en Memoria:** `137 bytes` (`0x89`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_REPORT_GET_ACK {
    unsigned short hReporter;
    unsigned short err;
    int nRegNumber;
    unsigned char byReportResult;
    unsigned char byRemark[128];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `hReporter` | Puerto TCP de escucha |
| `0x002` | `2` | `unsigned short` | `err` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `int` | `nRegNumber` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char` | `byReportResult` | Puerto TCP de escucha |
| `0x009` | `9` | `unsigned char[128]` | `byRemark` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_REPORT_GET_REQ`

- **OpCode Asociado:** `0x800B (Dept: 32, Cmd: 11)`
- **Tamaño Total en Memoria:** `7 bytes` (`0x7`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_REPORT_GET_REQ {
    unsigned short hReporter;
    unsigned int dwReporter;
    unsigned char nWorldNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `hReporter` | Puerto TCP de escucha |
| `0x002` | `2` | `unsigned int` | `dwReporter` | Puerto TCP de escucha |
| `0x006` | `6` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |

### `PROTO_NC_DATA_SEND1<PROTO_NC_DATA_REPORT_ADD_ACK>`

- **Tamaño Total en Memoria:** `31 bytes` (`0x1F`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_DATA_REPORT_ADD_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_DATA_REPORT_ADD_ACK /* base */ PROTO_NC_DATA_REPORT_ADD_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_DATA_REPORT_ADD_ACK` | `/* base */ PROTO_NC_DATA_REPORT_ADD_ACK` | Puerto TCP de escucha |

### `PROTO_NC_DATA_SEND1<PROTO_NC_DATA_REPORT_CANCEL_ACK>`

- **Tamaño Total en Memoria:** `11 bytes` (`0xB`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_DATA_REPORT_CANCEL_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_DATA_REPORT_CANCEL_ACK /* base */ PROTO_NC_DATA_REPORT_CANCEL_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_DATA_REPORT_CANCEL_ACK` | `/* base */ PROTO_NC_DATA_REPORT_CANCEL_ACK` | Puerto TCP de escucha |

### `PROTO_NC_DATA_SEND1<PROTO_NC_DATA_REPORT_GET_ACK>`

- **Tamaño Total en Memoria:** `140 bytes` (`0x8C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_DATA_REPORT_GET_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_DATA_REPORT_GET_ACK /* base */ PROTO_NC_DATA_REPORT_GET_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_DATA_REPORT_GET_ACK` | `/* base */ PROTO_NC_DATA_REPORT_GET_ACK` | Puerto TCP de escucha |

### `PROTO_NC_DATA_SEND1<PROTO_NC_MISC_EVENT_L20_DB_ACK>`

- **Tamaño Total en Memoria:** `16 bytes` (`0x10`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_MISC_EVENT_L20_DB_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_MISC_EVENT_L20_DB_ACK /* base */ PROTO_NC_MISC_EVENT_L20_DB_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_MISC_EVENT_L20_DB_ACK` | `/* base */ PROTO_NC_MISC_EVENT_L20_DB_ACK` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_SEND1<PROTO_NC_USER_GER_IS_IP_BLOCK_ACK>`

- **Tamaño Total en Memoria:** `293 bytes` (`0x125`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_USER_GER_IS_IP_BLOCK_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_USER_GER_IS_IP_BLOCK_ACK /* base */ PROTO_NC_USER_GER_IS_IP_BLOCK_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_USER_GER_IS_IP_BLOCK_ACK` | `/* base */ PROTO_NC_USER_GER_IS_IP_BLOCK_ACK` | Dirección IP del cliente o del servidor |

### `PROTO_NC_DATA_SEND1<PROTO_NC_USER_IS_IP_BLOCK_ACK>`

- **Tamaño Total en Memoria:** `278 bytes` (`0x116`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_USER_IS_IP_BLOCK_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_USER_IS_IP_BLOCK_ACK /* base */ PROTO_NC_USER_IS_IP_BLOCK_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_USER_IS_IP_BLOCK_ACK` | `/* base */ PROTO_NC_USER_IS_IP_BLOCK_ACK` | Dirección IP del cliente o del servidor |

### `PROTO_NC_DATA_SEND1<PROTO_NC_USER_TW_IS_IP_BLOCK_ACK>`

- **Tamaño Total en Memoria:** `14 bytes` (`0xE`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_USER_TW_IS_IP_BLOCK_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_USER_TW_IS_IP_BLOCK_ACK /* base */ PROTO_NC_USER_TW_IS_IP_BLOCK_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_USER_TW_IS_IP_BLOCK_ACK` | `/* base */ PROTO_NC_USER_TW_IS_IP_BLOCK_ACK` | Dirección IP del cliente o del servidor |

### `PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_REQ`

- **Tamaño Total en Memoria:** `12 bytes` (`0xC`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_REQ {
    unsigned int nGTNo;
    unsigned int OwnerGuildNum;
    int nItem;
    PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_REQ::ITEM_OPT ItemOpt[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nGTNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `OwnerGuildNum` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `int` | `nItem` | Campo de paquete oficial extraído de PDB TPI |
| `0x00C` | `12` | `PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_REQ::ITEM_OPT[0]` | `ItemOpt` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_REQ::ITEM_OPT`

- **Tamaño Total en Memoria:** `33 bytes` (`0x21`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_GUILD_TOURNAMENT_DB_REWARD_CREATE_REQ::ITEM_OPT {
    unsigned char RewardType;
    unsigned int Value1;
    unsigned int Value2;
    unsigned int Value3;
    unsigned int IO_Str;
    unsigned int IO_Con;
    unsigned int IO_Dex;
    unsigned int IO_Int;
    unsigned int IO_Men;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `RewardType` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned int` | `Value1` | Campo de paquete oficial extraído de PDB TPI |
| `0x005` | `5` | `unsigned int` | `Value2` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `unsigned int` | `Value3` | Campo de paquete oficial extraído de PDB TPI |
| `0x00D` | `13` | `unsigned int` | `IO_Str` | Campo de paquete oficial extraído de PDB TPI |
| `0x011` | `17` | `unsigned int` | `IO_Con` | Campo de paquete oficial extraído de PDB TPI |
| `0x015` | `21` | `unsigned int` | `IO_Dex` | Campo de paquete oficial extraído de PDB TPI |
| `0x019` | `25` | `unsigned int` | `IO_Int` | Campo de paquete oficial extraído de PDB TPI |
| `0x01D` | `29` | `unsigned int` | `IO_Men` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_ITEMDB_CHARGED_LIST_ACK`

- **Tamaño Total en Memoria:** `11 bytes` (`0xB`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_ITEMDB_CHARGED_LIST_ACK {
    NETPACKETZONEHEADER packheader;
    unsigned short ErrorCode;
    unsigned char nPartMark;
    unsigned short NumOfChargedItem;
    PROTO_CHARGED_ITEM_INFO ChargedItemInfoList[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETZONEHEADER` | `packheader` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned short` | `ErrorCode` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char` | `nPartMark` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `unsigned short` | `NumOfChargedItem` | Campo de paquete oficial extraído de PDB TPI |
| `0x00B` | `11` | `PROTO_CHARGED_ITEM_INFO[0]` | `ChargedItemInfoList` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_ITEM_BUY_REQ`

- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_ITEM_BUY_REQ {
    unsigned short itemid;
    unsigned int lot;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `itemid` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned int` | `lot` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_ITEM_DROP_REQ`

- **Tamaño Total en Memoria:** `14 bytes` (`0xE`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_ITEM_DROP_REQ {
    ITEM_INVEN slot;
    unsigned int lot;
    SHINE_XY_TYPE loc;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `ITEM_INVEN` | `slot` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned int` | `lot` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `SHINE_XY_TYPE` | `loc` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_ITEM_SELL_REQ`

- **Tamaño Total en Memoria:** `5 bytes` (`0x5`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_ITEM_SELL_REQ {
    unsigned char slot;
    unsigned int lot;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `slot` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned int` | `lot` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_USER_LOGIN`

- **OpCode Asociado:** `0x0721 (Dept: 1, Cmd: 801)`
- **Tamaño Total en Memoria:** `9 bytes` (`0x9`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_USER_LOGIN {
    unsigned int usrregnum;
    unsigned char worldnum;
    unsigned char ip[4];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `usrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `worldnum` | Identificador o información del WorldServer |
| `0x005` | `5` | `unsigned char[4]` | `ip` | Dirección IP del cliente o del servidor |

### `PROTO_NC_LOG_USER_LOGINFAIL`

- **OpCode Asociado:** `0x0720 (Dept: 1, Cmd: 800)`
- **Tamaño Total en Memoria:** `278 bytes` (`0x116`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_USER_LOGINFAIL {
    char userid[256]; /* Name256Byte */
    char userpassword[16]; /* Name4 */
    unsigned char userip[4];
    unsigned short error;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `char[256] /* Name256Byte */` | `userid` | Campo de paquete oficial extraído de PDB TPI |
| `0x100` | `256` | `char[16] /* Name4 */` | `userpassword` | Contraseña en texto plano o hash del cliente |
| `0x110` | `272` | `unsigned char[4]` | `userip` | Dirección IP del cliente o del servidor |
| `0x114` | `276` | `unsigned short` | `error` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_USER_LOGOUT`

- **OpCode Asociado:** `0x0722 (Dept: 1, Cmd: 802)`
- **Tamaño Total en Memoria:** `9 bytes` (`0x9`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_USER_LOGOUT {
    unsigned int usrregnum;
    unsigned char worldnum;
    unsigned short playmin;
    unsigned short error;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `usrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `worldnum` | Identificador o información del WorldServer |
| `0x005` | `5` | `unsigned short` | `playmin` | Campo de paquete oficial extraído de PDB TPI |
| `0x007` | `7` | `unsigned short` | `error` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_CHAR_LOGOFF_STATISTICS`

- **OpCode Asociado:** `0x081F (Dept: 2, Cmd: 31)`
- **Tamaño Total en Memoria:** `21 bytes` (`0x15`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_CHAR_LOGOFF_STATISTICS {
    unsigned char WorldNo;
    unsigned int UserNo;
    unsigned int CharNo;
    unsigned char CharClass;
    unsigned char CharLevel;
    unsigned __int64 CharGetExp;
    unsigned short CharPlayMin;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `WorldNo` | Identificador o información del WorldServer |
| `0x001` | `1` | `unsigned int` | `UserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x005` | `5` | `unsigned int` | `CharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `unsigned char` | `CharClass` | Campo de paquete oficial extraído de PDB TPI |
| `0x00A` | `10` | `unsigned char` | `CharLevel` | Campo de paquete oficial extraído de PDB TPI |
| `0x00B` | `11` | `unsigned __int64` | `CharGetExp` | Campo de paquete oficial extraído de PDB TPI |
| `0x013` | `19` | `unsigned short` | `CharPlayMin` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_CONNECTER_ACK`

- **OpCode Asociado:** `0x080A (Dept: 2, Cmd: 10)`
- **Tamaño Total en Memoria:** `131 bytes` (`0x83`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_CONNECTER_ACK {
    unsigned char worldno;
    unsigned short connectercount;
    unsigned short nEachClass[64];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `worldno` | Identificador o información del WorldServer |
| `0x001` | `1` | `unsigned short` | `connectercount` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `unsigned short[64]` | `nEachClass` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_CONNECTER_REQ`

- **OpCode Asociado:** `0x0809 (Dept: 2, Cmd: 9)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_CONNECTER_REQ {
    // Estructura vacía / Payload de 0 bytes
};
#pragma pack(pop)
```

> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*

### `PROTO_NC_MISC_EVENT_L20_DB_ACK`

- **OpCode Asociado:** `0x082D (Dept: 2, Cmd: 45)`
- **Tamaño Total en Memoria:** `13 bytes` (`0xD`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_EVENT_L20_DB_ACK {
    unsigned short nWMHandle;
    unsigned int nUserNo;
    unsigned char nEventNo;
    unsigned int nCharNo;
    unsigned short nError;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nWMHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned char` | `nEventNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x007` | `7` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x00B` | `11` | `unsigned short` | `nError` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_EVENT_L20_DB_REQ`

- **OpCode Asociado:** `0x082C (Dept: 2, Cmd: 44)`
- **Tamaño Total en Memoria:** `272 bytes` (`0x110`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_EVENT_L20_DB_REQ {
    unsigned short nWMHandle;
    char sUserID[256]; /* Name256Byte */
    unsigned int nUserNo;
    unsigned char nEventNo;
    unsigned char nWorldNo;
    unsigned int nCharNo;
    unsigned int nEventPoint;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nWMHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `char[256] /* Name256Byte */` | `sUserID` | Campo de paquete oficial extraído de PDB TPI |
| `0x102` | `258` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x106` | `262` | `unsigned char` | `nEventNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x107` | `263` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |
| `0x108` | `264` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x10C` | `268` | `unsigned int` | `nEventPoint` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD`

- **OpCode Asociado:** `0x0843 (Dept: 2, Cmd: 67)`
- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD {
    unsigned int nUserNo;
    unsigned short nPlayTime;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned short` | `nPlayTime` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_HACK_SCAN_STORE_CMD`

- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_HACK_SCAN_STORE_CMD {
    unsigned short nDescSize;
    unsigned char sDesc[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nDescSize` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned char[0]` | `sDesc` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD`

- **OpCode Asociado:** `0x0835 (Dept: 2, Cmd: 53)`
- **Tamaño Total en Memoria:** `51 bytes` (`0x33`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_HACK_SCAN_STORE_DB_CMD {
    unsigned int nUserNo;
    unsigned char nWorldNo;
    unsigned int nCharNo;
    unsigned char sIP[20];
    Name5 sCharID;
    PROTO_NC_MISC_HACK_SCAN_STORE_CMD ScanInfo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |
| `0x005` | `5` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `unsigned char[20]` | `sIP` | Dirección IP del cliente o del servidor |
| `0x01D` | `29` | `Name5` | `sCharID` | Campo de paquete oficial extraído de PDB TPI |
| `0x031` | `49` | `PROTO_NC_MISC_HACK_SCAN_STORE_CMD` | `ScanInfo` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_HEARTBEAT_ACK_SEND`

- **OpCode Asociado:** `0x0805 (Dept: 2, Cmd: 5)`
- **Tamaño Total en Memoria:** `3 bytes` (`0x3`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_HEARTBEAT_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |

### `PROTO_NC_MISC_ITEMSHOP_URL_DB_ACK`

- **Tamaño Total en Memoria:** `8 bytes` (`0x8`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_ITEMSHOP_URL_DB_ACK {
    unsigned int nCharNo;
    unsigned short nError;
    unsigned short nLen;
    char sURL[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned short` | `nError` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned short` | `nLen` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `char[0]` | `sURL` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_ITEMSHOP_URL_DB_REQ`

- **Tamaño Total en Memoria:** `8 bytes` (`0x8`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_ITEMSHOP_URL_DB_REQ {
    unsigned int nCharNo;
    unsigned int nUserNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_S2SCONNECTION_ACK`

- **OpCode Asociado:** `0x0803 (Dept: 2, Cmd: 3)`
- **Tamaño Total en Memoria:** `3 bytes` (`0x3`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_S2SCONNECTION_ACK {
    unsigned char echo_data;
    unsigned short error;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `echo_data` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned short` | `error` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_S2SCONNECTION_RDY_SEND`

- **OpCode Asociado:** `0x0801 (Dept: 2, Cmd: 1)`
- **Tamaño Total en Memoria:** `3 bytes` (`0x3`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_S2SCONNECTION_RDY_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |

### `PROTO_NC_MISC_S2SCONNECTION_REQ`

- **OpCode Asociado:** `0x0802 (Dept: 2, Cmd: 2)`
- **Tamaño Total en Memoria:** `7 bytes` (`0x7`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_S2SCONNECTION_REQ {
    unsigned char echo_data;
    unsigned char server_to_id;
    unsigned char server_from_id;
    unsigned char server_from_world_num;
    unsigned char server_from_zone_num;
    unsigned short key;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `echo_data` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned char` | `server_to_id` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned char` | `server_from_id` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `unsigned char` | `server_from_world_num` | Identificador o información del WorldServer |
| `0x004` | `4` | `unsigned char` | `server_from_zone_num` | Campo de paquete oficial extraído de PDB TPI |
| `0x005` | `5` | `unsigned short` | `key` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_OPTOOL_CLOSE_SERVER_ACK`

- **OpCode Asociado:** `0x2804 (Dept: 10, Cmd: 4)`
- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_CLOSE_SERVER_ACK {
    unsigned short error;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `error` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_OPTOOL_CONNECT_BRIF_ACK`

- **OpCode Asociado:** `0x2808 (Dept: 10, Cmd: 8)`
- **Tamaño Total en Memoria:** `170 bytes` (`0xAA`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_CONNECT_BRIF_ACK {
    unsigned char echo_data;
    unsigned char my_server_id;
    unsigned short count[2][8][168];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `echo_data` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned char` | `my_server_id` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short[2][8][168]` | `count` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_OPTOOL_CONNECT_BRIF_REQ`

- **OpCode Asociado:** `0x2807 (Dept: 10, Cmd: 7)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_CONNECT_BRIF_REQ {
    unsigned char echo_data;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `echo_data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_OPTOOL_MAP_USER_LIST_ACK`

- **OpCode Asociado:** `0x2806 (Dept: 10, Cmd: 6)`
- **Tamaño Total en Memoria:** `4 bytes` (`0x4`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_MAP_USER_LIST_ACK {
    unsigned char echo_data;
    unsigned char my_server_id;
    unsigned short num_of_list;
    NC_OPTOOL_MAP_USER_LIST_INFO user_info[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `echo_data` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned char` | `my_server_id` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short` | `num_of_list` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `NC_OPTOOL_MAP_USER_LIST_INFO[0]` | `user_info` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_OPTOOL_MAP_USER_LIST_REQ`

- **OpCode Asociado:** `0x2805 (Dept: 10, Cmd: 5)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_MAP_USER_LIST_REQ {
    unsigned char echo_data;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `echo_data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_OPTOOL_S2SCONNECT_LIST_ACK`

- **OpCode Asociado:** `0x2802 (Dept: 10, Cmd: 2)`
- **Tamaño Total en Memoria:** `4 bytes` (`0x4`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_S2SCONNECT_LIST_ACK {
    unsigned char echo_data;
    unsigned char my_server_id;
    unsigned short num_of_connections;
    SERVER_CONN_INFO connection_info[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `echo_data` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned char` | `my_server_id` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short` | `num_of_connections` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `SERVER_CONN_INFO[0]` | `connection_info` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_OPTOOL_S2SCONNECT_LIST_REQ`

- **OpCode Asociado:** `0x2801 (Dept: 10, Cmd: 1)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_S2SCONNECT_LIST_REQ {
    unsigned char echo_data;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `echo_data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_SYSLOG_ACCOUNT_LOGIN_FAILURE`

- **OpCode Asociado:** `0xC802 (Dept: 50, Cmd: 2)`
- **Tamaño Total en Memoria:** `296 bytes` (`0x128`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_SYSLOG_ACCOUNT_LOGIN_FAILURE {
    char UserName[260];
    char UserIP[16]; /* Name4 */
    Name5 spawnapps;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `char[260]` | `UserName` | Nombre de cuenta o usuario (string) |
| `0x104` | `260` | `char[16] /* Name4 */` | `UserIP` | Dirección IP del cliente o del servidor |
| `0x114` | `276` | `Name5` | `spawnapps` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_SYSLOG_ACCOUNT_LOGIN_SUCCESS`

- **OpCode Asociado:** `0xC801 (Dept: 50, Cmd: 1)`
- **Tamaño Total en Memoria:** `40 bytes` (`0x28`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_SYSLOG_ACCOUNT_LOGIN_SUCCESS {
    unsigned int UserNo;
    char UserIP[16]; /* Name4 */
    Name5 spawnapps;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `UserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `char[16] /* Name4 */` | `UserIP` | Dirección IP del cliente o del servidor |
| `0x014` | `20` | `Name5` | `spawnapps` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_SYSLOG_ACCOUNT_LOGOUT`

- **OpCode Asociado:** `0xC803 (Dept: 50, Cmd: 3)`
- **Tamaño Total en Memoria:** `8 bytes` (`0x8`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_SYSLOG_ACCOUNT_LOGOUT {
    unsigned int UserNo;
    unsigned int LoginKeepTime;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `UserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `LoginKeepTime` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_SYSLOG_SERVER_CCU`

- **OpCode Asociado:** `0xC811 (Dept: 50, Cmd: 17)`
- **Tamaño Total en Memoria:** `5 bytes` (`0x5`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_SYSLOG_SERVER_CCU {
    unsigned int UserCount;
    unsigned char WorldNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `UserCount` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `WorldNo` | Identificador o información del WorldServer |

### `PROTO_NC_USER_GER_IS_IP_BLOCK_ACK`

- **OpCode Asociado:** `0x0C3C (Dept: 3, Cmd: 60)`
- **Tamaño Total en Memoria:** `290 bytes` (`0x122`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_GER_IS_IP_BLOCK_ACK {
    PROTO_NC_USER_GER_PASSWORD_CHECK_ACK LoginData;
    unsigned short nClientHandle;
    unsigned char bUserIPBlock;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_GER_PASSWORD_CHECK_ACK` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x11F` | `287` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x121` | `289` | `unsigned char` | `bUserIPBlock` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_GER_IS_IP_BLOCK_REQ`

- **OpCode Asociado:** `0x0C3B (Dept: 3, Cmd: 59)`
- **Tamaño Total en Memoria:** `313 bytes` (`0x139`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_GER_IS_IP_BLOCK_REQ {
    PROTO_NC_USER_GER_PASSWORD_CHECK_ACK LoginData;
    unsigned int nUserNo;
    unsigned short nClientHandle;
    char sUserIP[20];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_GER_PASSWORD_CHECK_ACK` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x11F` | `287` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x123` | `291` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x125` | `293` | `char[20]` | `sUserIP` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_GER_PASSWORD_CHECK_ACK`

- **OpCode Asociado:** `0x0C3A (Dept: 3, Cmd: 58)`
- **Tamaño Total en Memoria:** `287 bytes` (`0x11F`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_GER_PASSWORD_CHECK_ACK {
    NETPACKETHEADER netpacketheader;
    char id[256]; /* Name256Byte */
    char password[16]; /* Name4 */
    unsigned int usrregnum;
    unsigned char authid;
    unsigned int block;
    unsigned int loginable;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `char[256] /* Name256Byte */` | `id` | Campo de paquete oficial extraído de PDB TPI |
| `0x102` | `258` | `char[16] /* Name4 */` | `password` | Contraseña en texto plano o hash del cliente |
| `0x112` | `274` | `unsigned int` | `usrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x116` | `278` | `unsigned char` | `authid` | Nivel de autorización (Player, GM, Admin) |
| `0x117` | `279` | `unsigned int` | `block` | Estado de bloqueo / baneo de cuenta o IP |
| `0x11B` | `283` | `unsigned int` | `loginable` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_IS_IP_BLOCK_ACK`

- **OpCode Asociado:** `0x0C2B (Dept: 3, Cmd: 43)`
- **Tamaño Total en Memoria:** `275 bytes` (`0x113`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_IS_IP_BLOCK_ACK {
    PROTO_NC_USER_LOGIN_REQ LoginData;
    unsigned short nClientHandle;
    unsigned char bUserIPBlock;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_LOGIN_REQ` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x110` | `272` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x112` | `274` | `unsigned char` | `bUserIPBlock` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_IS_IP_BLOCK_REQ`

- **OpCode Asociado:** `0x0C2A (Dept: 3, Cmd: 42)`
- **Tamaño Total en Memoria:** `298 bytes` (`0x12A`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_IS_IP_BLOCK_REQ {
    PROTO_NC_USER_LOGIN_REQ LoginData;
    unsigned int nUserNo;
    unsigned short nClientHandle;
    char sUserIP[20];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_LOGIN_REQ` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x110` | `272` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x114` | `276` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x116` | `278` | `char[20]` | `sUserIP` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_JP_IS_IP_BLOCK_ACK`

- **OpCode Asociado:** `0x0C4E (Dept: 3, Cmd: 78)`
- **Tamaño Total en Memoria:** `309 bytes` (`0x135`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_JP_IS_IP_BLOCK_ACK {
    PROTO_NC_USER_JP_PASSWORD_CHECK_ACK LoginData;
    unsigned short nClientHandle;
    unsigned char bUserIPBlock;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_JP_PASSWORD_CHECK_ACK` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x132` | `306` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x134` | `308` | `unsigned char` | `bUserIPBlock` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_JP_IS_IP_BLOCK_REQ`

- **OpCode Asociado:** `0x0C4D (Dept: 3, Cmd: 77)`
- **Tamaño Total en Memoria:** `332 bytes` (`0x14C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_JP_IS_IP_BLOCK_REQ {
    PROTO_NC_USER_JP_PASSWORD_CHECK_ACK LoginData;
    unsigned int nUserNo;
    unsigned short nClientHandle;
    char sUserIP[20];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_JP_PASSWORD_CHECK_ACK` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x132` | `306` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x136` | `310` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x138` | `312` | `char[20]` | `sUserIP` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_JP_PASSWORD_CHECK_ACK`

- **OpCode Asociado:** `0x0C4C (Dept: 3, Cmd: 76)`
- **Tamaño Total en Memoria:** `306 bytes` (`0x132`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_JP_PASSWORD_CHECK_ACK {
    NETPACKETHEADER netpacketheader;
    unsigned short nError;
    unsigned int nUserNo;
    char sUserName[24];
    char sChannel[254];
    char sUserIP[20];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `unsigned short` | `nError` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `char[24]` | `sUserName` | Nombre de cuenta o usuario (string) |
| `0x020` | `32` | `char[254]` | `sChannel` | Campo de paquete oficial extraído de PDB TPI |
| `0x11E` | `286` | `char[20]` | `sUserIP` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_LOGIN_REQ`

- **OpCode Asociado:** `0x0C06 (Dept: 3, Cmd: 6)`
- **Tamaño Total en Memoria:** `272 bytes` (`0x110`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGIN_REQ {
    char user[256]; /* Name256Byte */
    char password[16]; /* Name4 */
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `char[256] /* Name256Byte */` | `user` | Campo de paquete oficial extraído de PDB TPI |
| `0x100` | `256` | `char[16] /* Name4 */` | `password` | Contraseña en texto plano o hash del cliente |

### `PROTO_NC_USER_LOGOUT_LAST_DAY_ACK`

- **OpCode Asociado:** `0x0C63 (Dept: 3, Cmd: 99)`
- **Tamaño Total en Memoria:** `10 bytes` (`0xA`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGOUT_LAST_DAY_ACK {
    unsigned short nHandle;
    unsigned int nUserNo;
    int tLastSec;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `int` | `tLastSec` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_LOGOUT_LAST_DAY_REQ`

- **OpCode Asociado:** `0x0C62 (Dept: 3, Cmd: 98)`
- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGOUT_LAST_DAY_REQ {
    unsigned short nHandle;
    unsigned int nUserNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_TW_IS_IP_BLOCK_ACK`

- **OpCode Asociado:** `0x0C41 (Dept: 3, Cmd: 65)`
- **Tamaño Total en Memoria:** `11 bytes` (`0xB`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_TW_IS_IP_BLOCK_ACK {
    PROTO_NC_USER_TW_PASSWORD_CHECK_ACK LoginData;
    unsigned short nClientHandle;
    unsigned char bUserIPBlock;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_TW_PASSWORD_CHECK_ACK` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x00A` | `10` | `unsigned char` | `bUserIPBlock` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_TW_IS_IP_BLOCK_REQ`

- **OpCode Asociado:** `0x0C40 (Dept: 3, Cmd: 64)`
- **Tamaño Total en Memoria:** `34 bytes` (`0x22`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_TW_IS_IP_BLOCK_REQ {
    PROTO_NC_USER_TW_PASSWORD_CHECK_ACK LoginData;
    unsigned int nUserNo;
    unsigned short nClientHandle;
    char sUserIP[20];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_TW_PASSWORD_CHECK_ACK` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x00C` | `12` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x00E` | `14` | `char[20]` | `sUserIP` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_TW_PASSWORD_CHECK_ACK`

- **OpCode Asociado:** `0x0C3F (Dept: 3, Cmd: 63)`
- **Tamaño Total en Memoria:** `8 bytes` (`0x8`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_TW_PASSWORD_CHECK_ACK {
    NETPACKETHEADER netpacketheader;
    unsigned short nError;
    unsigned int nUserNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `unsigned short` | `nError` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_US_IS_IP_BLOCK_ACK`

- **OpCode Asociado:** `0x0C5E (Dept: 3, Cmd: 94)`
- **Tamaño Total en Memoria:** `356 bytes` (`0x164`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_US_IS_IP_BLOCK_ACK {
    PROTO_NC_USER_US_PASSWORD_CHECK_ACK LoginData;
    unsigned short nClientHandle;
    unsigned char bUserIPBlock;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_US_PASSWORD_CHECK_ACK` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x161` | `353` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x163` | `355` | `unsigned char` | `bUserIPBlock` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_US_IS_IP_BLOCK_REQ`

- **OpCode Asociado:** `0x0C5D (Dept: 3, Cmd: 93)`
- **Tamaño Total en Memoria:** `379 bytes` (`0x17B`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_US_IS_IP_BLOCK_REQ {
    PROTO_NC_USER_US_PASSWORD_CHECK_ACK LoginData;
    unsigned int nUserNo;
    unsigned short nClientHandle;
    char sUserIP[20];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_US_PASSWORD_CHECK_ACK` | `LoginData` | Campo de paquete oficial extraído de PDB TPI |
| `0x161` | `353` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x165` | `357` | `unsigned short` | `nClientHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x167` | `359` | `char[20]` | `sUserIP` | Dirección IP del cliente o del servidor |

### `PROTO_NC_USER_US_PASSWORD_CHECK_ACK`

- **OpCode Asociado:** `0x0C5C (Dept: 3, Cmd: 92)`
- **Tamaño Total en Memoria:** `353 bytes` (`0x161`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_US_PASSWORD_CHECK_ACK {
    NETPACKETHEADER netpacketheader;
    unsigned short nError;
    unsigned int nUserNo;
    char sUserName[260];
    char sPassword[36];
    char sUserIP[20];
    unsigned char nAuthID;
    unsigned int nBlock;
    unsigned int bLoginAble;
    Name5 spawnapps;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `unsigned short` | `nError` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `char[260]` | `sUserName` | Nombre de cuenta o usuario (string) |
| `0x10C` | `268` | `char[36]` | `sPassword` | Contraseña en texto plano o hash del cliente |
| `0x130` | `304` | `char[20]` | `sUserIP` | Dirección IP del cliente o del servidor |
| `0x144` | `324` | `unsigned char` | `nAuthID` | Nivel de autorización (Player, GM, Admin) |
| `0x145` | `325` | `unsigned int` | `nBlock` | Estado de bloqueo / baneo de cuenta o IP |
| `0x149` | `329` | `unsigned int` | `bLoginAble` | Campo de paquete oficial extraído de PDB TPI |
| `0x14D` | `333` | `Name5` | `spawnapps` | Campo de paquete oficial extraído de PDB TPI |
