# Especificación Oficial: `Login (LoginServer)`


## Ficha Técnica del Binario
- **Nombre de Binario:** `Login.exe`
- **Símbolos:** `Login.pdb` (Microsoft C/C++ MSF 7.00, TPI v4.10)
- **Ruta del Proyecto Original:** `E:\projectf2\server\3loginserver2\loginserver\`
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


## Catálogo Exhaustivo de Estructuras de Paquetes

> **Total de Estructuras Catalogadas:** `85`

> Todas las estructuras han sido extraídas directamente del flujo TPI del PDB oficial con sus tamaños, tipos primitivos, offsets y variables exactas.

### `CLoginClientSession::Send_NC_MISC_SEED_ACK::__l2::PROTO_NC_MISC_SEED_ACK_SEND`

- **Tamaño Total en Memoria:** `5 bytes` (`0x5`)

```cpp
#pragma pack(push, 1)
struct CLoginClientSession::Send_NC_MISC_SEED_ACK::__l2::PROTO_NC_MISC_SEED_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_MISC_SEED_ACK ack;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_MISC_SEED_ACK` | `ack` | Campo de paquete oficial extraído de PDB TPI |

### `CParserClient::fc_NC_USER_XTRAP_REQ::__l4::PROTO_NC_USER_XTRAP_ACK_SEND`

- **Tamaño Total en Memoria:** `4 bytes` (`0x4`)

```cpp
#pragma pack(push, 1)
struct CParserClient::fc_NC_USER_XTRAP_REQ::__l4::PROTO_NC_USER_XTRAP_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_XTRAP_ACK ack;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_XTRAP_ACK` | `ack` | Campo de paquete oficial extraído de PDB TPI |

### `CParserWorldManager::fc_NC_USER_CREATE_OTP_REQ::__l14::PROTO_NC_USER_CREATE_OTP_ACK_SEND`

- **Tamaño Total en Memoria:** `39 bytes` (`0x27`)

```cpp
#pragma pack(push, 1)
struct CParserWorldManager::fc_NC_USER_CREATE_OTP_REQ::__l14::PROTO_NC_USER_CREATE_OTP_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_CREATE_OTP_ACK ack;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_CREATE_OTP_ACK` | `ack` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_DATA_SEND1<PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_ACK>`

- **Tamaño Total en Memoria:** `16 bytes` (`0x10`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_ACK /* base */ PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_ACK` | `/* base */ PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_ACK` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_SEND1<PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_ACK>`

- **Tamaño Total en Memoria:** `27 bytes` (`0x1B`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_ACK /* base */ PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_ACK` | `/* base */ PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_ACK` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_SEND1<PROTO_NC_USER_GER_IS_IP_BLOCK_REQ>`

- **Tamaño Total en Memoria:** `316 bytes` (`0x13C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_USER_GER_IS_IP_BLOCK_REQ> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_USER_GER_IS_IP_BLOCK_REQ /* base */ PROTO_NC_USER_GER_IS_IP_BLOCK_REQ;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_USER_GER_IS_IP_BLOCK_REQ` | `/* base */ PROTO_NC_USER_GER_IS_IP_BLOCK_REQ` | Dirección IP del cliente o del servidor |

### `PROTO_NC_DATA_SEND1<PROTO_NC_USER_TEENAGER_CMD>`

- **Tamaño Total en Memoria:** `58 bytes` (`0x3A`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_USER_TEENAGER_CMD> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_USER_TEENAGER_CMD /* base */ PROTO_NC_USER_TEENAGER_CMD;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_USER_TEENAGER_CMD` | `/* base */ PROTO_NC_USER_TEENAGER_CMD` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_SEND1<PROTO_NC_USER_TW_IS_IP_BLOCK_REQ>`

- **Tamaño Total en Memoria:** `37 bytes` (`0x25`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_USER_TW_IS_IP_BLOCK_REQ> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_USER_TW_IS_IP_BLOCK_REQ /* base */ PROTO_NC_USER_TW_IS_IP_BLOCK_REQ;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_USER_TW_IS_IP_BLOCK_REQ` | `/* base */ PROTO_NC_USER_TW_IS_IP_BLOCK_REQ` | Dirección IP del cliente o del servidor |

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

### `PROTO_NC_LOG_USER_LOGINFAIL_SEND`

- **OpCode Asociado:** `0x0720 (Dept: 1, Cmd: 800)`
- **Tamaño Total en Memoria:** `281 bytes` (`0x119`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_USER_LOGINFAIL_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_LOG_USER_LOGINFAIL log;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_LOG_USER_LOGINFAIL` | `log` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_USER_LOGIN_SEND`

- **OpCode Asociado:** `0x0721 (Dept: 1, Cmd: 801)`
- **Tamaño Total en Memoria:** `12 bytes` (`0xC`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_USER_LOGIN_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_LOG_USER_LOGIN log;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_LOG_USER_LOGIN` | `log` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD`

- **OpCode Asociado:** `0x080B (Dept: 2, Cmd: 11)`
- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD {
    unsigned short packet_size;
    unsigned char packet[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `packet_size` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned char[0]` | `packet` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD`

- **OpCode Asociado:** `0x080C (Dept: 2, Cmd: 12)`
- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD {
    unsigned short packet_size;
    unsigned char packet[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `packet_size` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned char[0]` | `packet` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_MISC_HEARTBEAT_REQ_SEND`

- **OpCode Asociado:** `0x0804 (Dept: 2, Cmd: 4)`
- **Tamaño Total en Memoria:** `3 bytes` (`0x3`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_HEARTBEAT_REQ_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |

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

### `PROTO_NC_MISC_S2SCONNECTION_REQ_SEND`

- **OpCode Asociado:** `0x0802 (Dept: 2, Cmd: 2)`
- **Tamaño Total en Memoria:** `10 bytes` (`0xA`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_S2SCONNECTION_REQ_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_MISC_S2SCONNECTION_REQ req;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_MISC_S2SCONNECTION_REQ` | `req` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_SEED_ACK`

- **OpCode Asociado:** `0x0807 (Dept: 2, Cmd: 7)`
- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_SEED_ACK {
    unsigned short seed;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `seed` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_USER_COUNT_CMD`

- **OpCode Asociado:** `0x0808 (Dept: 2, Cmd: 8)`
- **Tamaño Total en Memoria:** `4 bytes` (`0x4`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_USER_COUNT_CMD {
    unsigned short NumOfUser;
    unsigned short MaxOfUser;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `NumOfUser` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short` | `MaxOfUser` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_ACK`

- **OpCode Asociado:** `0x2813 (Dept: 10, Cmd: 19)`
- **Tamaño Total en Memoria:** `13 bytes` (`0xD`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_ACK {
    unsigned char nWorldNo;
    LOGIN_USER_RATABLE RateTable;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |
| `0x001` | `1` | `LOGIN_USER_RATABLE` | `RateTable` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_REQ`

- **OpCode Asociado:** `0x2812 (Dept: 10, Cmd: 18)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_GET_REQ {
    unsigned char nWorldNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |

### `PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_SET_CMD`

- **OpCode Asociado:** `0x2814 (Dept: 10, Cmd: 20)`
- **Tamaño Total en Memoria:** `13 bytes` (`0xD`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_LOGIN_USER_RATABLE_SET_CMD {
    unsigned char nWorldNo;
    LOGIN_USER_RATABLE RateTable;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |
| `0x001` | `1` | `LOGIN_USER_RATABLE` | `RateTable` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_ACK`

- **OpCode Asociado:** `0x2819 (Dept: 10, Cmd: 25)`
- **Tamaño Total en Memoria:** `24 bytes` (`0x18`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_OPTOOL_LOGON_PROCESS_TIME_VIEW_ACK {
    unsigned int Connect_Count;
    unsigned int Connect_Time;
    unsigned int Login_Count;
    unsigned int Login_Time;
    unsigned int IPBlock_Count;
    unsigned int IPBlock_Time;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `Connect_Count` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `Connect_Time` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned int` | `Login_Count` | Campo de paquete oficial extraído de PDB TPI |
| `0x00C` | `12` | `unsigned int` | `Login_Time` | Campo de paquete oficial extraído de PDB TPI |
| `0x010` | `16` | `unsigned int` | `IPBlock_Count` | Dirección IP del cliente o del servidor |
| `0x014` | `20` | `unsigned int` | `IPBlock_Time` | Dirección IP del cliente o del servidor |

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

### `PROTO_NC_USER_CH_PASSWORD_CHECK_ACK`

- **OpCode Asociado:** `0x0C54 (Dept: 3, Cmd: 84)`
- **Tamaño Total en Memoria:** `287 bytes` (`0x11F`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CH_PASSWORD_CHECK_ACK {
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

### `PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK`

- **OpCode Asociado:** `0x0C67 (Dept: 3, Cmd: 103)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK {
    unsigned char XTrapServerKeyLength;
    unsigned char XTrapServerKey[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `XTrapServerKeyLength` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned char[0]` | `XTrapServerKey` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ`

- **OpCode Asociado:** `0x0C65 (Dept: 3, Cmd: 101)`
- **Tamaño Total en Memoria:** `64 bytes` (`0x40`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ {
    char sVersionKey[64];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `char[64]` | `sVersionKey` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK`

- **OpCode Asociado:** `0x0C02 (Dept: 3, Cmd: 2)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK {
    // Estructura vacía / Payload de 0 bytes
};
#pragma pack(pop)
```

> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*

### `PROTO_NC_USER_CONNECTCUT_CMD`

- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CONNECTCUT_CMD {
    unsigned short reason;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `reason` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_CONNECTCUT_CMD_SEND`

- **Tamaño Total en Memoria:** `5 bytes` (`0x5`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CONNECTCUT_CMD_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_CONNECTCUT_CMD cmd;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_CONNECTCUT_CMD` | `cmd` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_CREATE_OTP_ACK`

- **OpCode Asociado:** `0x0C36 (Dept: 3, Cmd: 54)`
- **Tamaño Total en Memoria:** `36 bytes` (`0x24`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CREATE_OTP_ACK {
    unsigned short nError;
    unsigned short hWM;
    Name8 sOTP;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nError` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short` | `hWM` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `Name8` | `sOTP` | One-Time Password token |

### `PROTO_NC_USER_CREATE_OTP_REQ`

- **OpCode Asociado:** `0x0C35 (Dept: 3, Cmd: 53)`
- **Tamaño Total en Memoria:** `266 bytes` (`0x10A`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CREATE_OTP_REQ {
    unsigned short hWM;
    unsigned int nUserNo;
    char sUserID[256]; /* Name256Byte */
    in_addr nClientIP;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `hWM` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `char[256] /* Name256Byte */` | `sUserID` | Campo de paquete oficial extraído de PDB TPI |
| `0x106` | `262` | `in_addr` | `nClientIP` | Dirección IP del cliente o del servidor |

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

### `PROTO_NC_USER_KICKOFFFROMWORLD_CMD`

- **Tamaño Total en Memoria:** `4 bytes` (`0x4`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_KICKOFFFROMWORLD_CMD {
    unsigned int usrregnum;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `usrregnum` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_KICKOFFFROMWORLD_CMD_SEND`

- **Tamaño Total en Memoria:** `13 bytes` (`0xD`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_KICKOFFFROMWORLD_CMD_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    NETPACKETZONEHEADER netpacketzoneheader;
    PROTO_NC_USER_KICKOFFFROMWORLD_CMD cmd;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `NETPACKETZONEHEADER` | `netpacketzoneheader` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `PROTO_NC_USER_KICKOFFFROMWORLD_CMD` | `cmd` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_LOGINFAIL_ACK`

- **OpCode Asociado:** `0x0C09 (Dept: 3, Cmd: 9)`
- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGINFAIL_ACK {
    unsigned short error;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `error` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_LOGINFAIL_ACK_SEND`

- **OpCode Asociado:** `0x0C09 (Dept: 3, Cmd: 9)`
- **Tamaño Total en Memoria:** `5 bytes` (`0x5`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGINFAIL_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_LOGINFAIL_ACK ack;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_LOGINFAIL_ACK` | `ack` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_LOGIN_ACK`

- **OpCode Asociado:** `0x0C0A (Dept: 3, Cmd: 10)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGIN_ACK {
    unsigned char numofworld;
    PROTO_NC_USER_LOGIN_ACK::WorldInfo worldinfo[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `numofworld` | Identificador o información del WorldServer |
| `0x001` | `1` | `PROTO_NC_USER_LOGIN_ACK::WorldInfo[0]` | `worldinfo` | Identificador o información del WorldServer |

### `PROTO_NC_USER_LOGIN_ACK::WorldInfo`

- **Tamaño Total en Memoria:** `18 bytes` (`0x12`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGIN_ACK::WorldInfo {
    unsigned char worldno;
    char worldname[16]; /* Name4 */
    unsigned char worldstatus;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `worldno` | Identificador o información del WorldServer |
| `0x001` | `1` | `char[16] /* Name4 */` | `worldname` | Identificador o información del WorldServer |
| `0x011` | `17` | `unsigned char` | `worldstatus` | Identificador o información del WorldServer |

### `PROTO_NC_USER_LOGIN_DB`

- **Tamaño Total en Memoria:** `9 bytes` (`0x9`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGIN_DB {
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

### `PROTO_NC_USER_LOGIN_DB_SEND`

- **Tamaño Total en Memoria:** `12 bytes` (`0xC`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGIN_DB_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_LOGIN_DB db;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_LOGIN_DB` | `db` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_LOGIN_NETMARBLE_REQ`

- **OpCode Asociado:** `0x0C1D (Dept: 3, Cmd: 29)`
- **Tamaño Total en Memoria:** `2 bytes` (`0x2`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGIN_NETMARBLE_REQ {
    unsigned short CPCookieLength;
    unsigned char CPCookie[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `CPCookieLength` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned char[0]` | `CPCookie` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_LOGIN_OUTSPARK_REQ`

- **OpCode Asociado:** `0x0C20 (Dept: 3, Cmd: 32)`
- **Tamaño Total en Memoria:** `85 bytes` (`0x55`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGIN_OUTSPARK_REQ {
    unsigned char LogonToken[65];
    Name5 spawnapps;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[65]` | `LogonToken` | Token de autenticación externa |
| `0x041` | `65` | `Name5` | `spawnapps` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_USER_LOGIN_WITH_OTP_REQ`

- **OpCode Asociado:** `0x0C37 (Dept: 3, Cmd: 55)`
- **Tamaño Total en Memoria:** `32 bytes` (`0x20`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGIN_WITH_OTP_REQ {
    Name8 sOTP;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `Name8` | `sOTP` | One-Time Password token |

### `PROTO_NC_USER_PASSWORD_CHECK_ACK`

- **OpCode Asociado:** `0x0C08 (Dept: 3, Cmd: 8)`
- **Tamaño Total en Memoria:** `287 bytes` (`0x11F`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_PASSWORD_CHECK_ACK {
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

### `PROTO_NC_USER_PASSWORD_CHECK_REQ`

- **OpCode Asociado:** `0x0C07 (Dept: 3, Cmd: 7)`
- **Tamaño Total en Memoria:** `274 bytes` (`0x112`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_PASSWORD_CHECK_REQ {
    NETPACKETHEADER netpacketheader;
    char user[256]; /* Name256Byte */
    char password[16]; /* Name4 */
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `char[256] /* Name256Byte */` | `user` | Campo de paquete oficial extraído de PDB TPI |
| `0x102` | `258` | `char[16] /* Name4 */` | `password` | Contraseña en texto plano o hash del cliente |

### `PROTO_NC_USER_PASSWORD_CHECK_REQ_SEND`

- **OpCode Asociado:** `0x0C07 (Dept: 3, Cmd: 7)`
- **Tamaño Total en Memoria:** `277 bytes` (`0x115`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_PASSWORD_CHECK_REQ_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_PASSWORD_CHECK_REQ req;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_PASSWORD_CHECK_REQ` | `req` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_TEENAGER_ACK`

- **OpCode Asociado:** `0x0C27 (Dept: 3, Cmd: 39)`
- **Tamaño Total en Memoria:** `313 bytes` (`0x139`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_TEENAGER_ACK {
    char sUserID[256]; /* Name256Byte */
    unsigned short nHandle;
    PROTO_NC_USER_TEENAGER_CMD Data;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `char[256] /* Name256Byte */` | `sUserID` | Campo de paquete oficial extraído de PDB TPI |
| `0x100` | `256` | `unsigned short` | `nHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x102` | `258` | `PROTO_NC_USER_TEENAGER_CMD` | `Data` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_TEENAGER_CMD`

- **Tamaño Total en Memoria:** `55 bytes` (`0x37`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_TEENAGER_CMD {
    unsigned char sUserKey[50];
    unsigned char nType;
    unsigned int nMinute;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[50]` | `sUserKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x032` | `50` | `unsigned char` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x033` | `51` | `unsigned int` | `nMinute` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_USER_US_LOGIN_REQ`

- **OpCode Asociado:** `0x0C5A (Dept: 3, Cmd: 90)`
- **Tamaño Total en Memoria:** `316 bytes` (`0x13C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_US_LOGIN_REQ {
    char sUserName[260];
    char sPassword[36];
    Name5 spawnapps;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `char[260]` | `sUserName` | Nombre de cuenta o usuario (string) |
| `0x104` | `260` | `char[36]` | `sPassword` | Contraseña en texto plano o hash del cliente |
| `0x128` | `296` | `Name5` | `spawnapps` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_USER_US_PASSWORD_CHECK_REQ`

- **OpCode Asociado:** `0x0C5B (Dept: 3, Cmd: 91)`
- **Tamaño Total en Memoria:** `338 bytes` (`0x152`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_US_PASSWORD_CHECK_REQ {
    NETPACKETHEADER netpacketheader;
    char sUserName[260];
    char sPassword[36];
    char sUserIP[20];
    Name5 spawnapps;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `char[260]` | `sUserName` | Nombre de cuenta o usuario (string) |
| `0x106` | `262` | `char[36]` | `sPassword` | Contraseña en texto plano o hash del cliente |
| `0x12A` | `298` | `char[20]` | `sUserIP` | Dirección IP del cliente o del servidor |
| `0x13E` | `318` | `Name5` | `spawnapps` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_WILLLOGIN_ACK`

- **OpCode Asociado:** `0x0C0E (Dept: 3, Cmd: 14)`
- **Tamaño Total en Memoria:** `348 bytes` (`0x15C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_WILLLOGIN_ACK {
    PROTO_NC_USER_WILLLOGIN_REQ will_login_req;
    unsigned char ok;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `PROTO_NC_USER_WILLLOGIN_REQ` | `will_login_req` | Campo de paquete oficial extraído de PDB TPI |
| `0x15B` | `347` | `unsigned char` | `ok` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_WILLLOGIN_REQ`

- **OpCode Asociado:** `0x0C0D (Dept: 3, Cmd: 13)`
- **Tamaño Total en Memoria:** `347 bytes` (`0x15B`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_WILLLOGIN_REQ {
    NETPACKETHEADER netpacketheader;
    unsigned char willlogin_worldno;
    unsigned int usrregnum;
    char userid[256]; /* Name256Byte */
    unsigned short validate_new[32];
    Name5 spawnapps;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `unsigned char` | `willlogin_worldno` | Identificador o información del WorldServer |
| `0x003` | `3` | `unsigned int` | `usrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x007` | `7` | `char[256] /* Name256Byte */` | `userid` | Campo de paquete oficial extraído de PDB TPI |
| `0x107` | `263` | `unsigned short[32]` | `validate_new` | Campo de paquete oficial extraído de PDB TPI |
| `0x147` | `327` | `Name5` | `spawnapps` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_WORLDSELECT_ACK`

- **OpCode Asociado:** `0x0C0C (Dept: 3, Cmd: 12)`
- **Tamaño Total en Memoria:** `83 bytes` (`0x53`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_WORLDSELECT_ACK {
    unsigned char worldstatus;
    char ip[16]; /* Name4 */
    unsigned short port;
    unsigned short validate_new[32];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `worldstatus` | Identificador o información del WorldServer |
| `0x001` | `1` | `char[16] /* Name4 */` | `ip` | Dirección IP del cliente o del servidor |
| `0x011` | `17` | `unsigned short` | `port` | Puerto TCP de escucha |
| `0x013` | `19` | `unsigned short[32]` | `validate_new` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_WORLDSELECT_ACK_SEND`

- **OpCode Asociado:** `0x0C0C (Dept: 3, Cmd: 12)`
- **Tamaño Total en Memoria:** `86 bytes` (`0x56`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_WORLDSELECT_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_WORLDSELECT_ACK ack;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_WORLDSELECT_ACK` | `ack` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_WORLDSELECT_REQ`

- **OpCode Asociado:** `0x0C0B (Dept: 3, Cmd: 11)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_WORLDSELECT_REQ {
    unsigned char worldno;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `worldno` | Identificador o información del WorldServer |

### `PROTO_NC_USER_WORLD_STATUS_REQ`

- **OpCode Asociado:** `0x0C1B (Dept: 3, Cmd: 27)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_WORLD_STATUS_REQ {
    // Estructura vacía / Payload de 0 bytes
};
#pragma pack(pop)
```

> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*

### `PROTO_NC_USER_XTRAP_ACK`

- **OpCode Asociado:** `0x0C05 (Dept: 3, Cmd: 5)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_XTRAP_ACK {
    unsigned char bSuccess;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `bSuccess` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_XTRAP_REQ`

- **OpCode Asociado:** `0x0C04 (Dept: 3, Cmd: 4)`
- **Tamaño Total en Memoria:** `1 bytes` (`0x1`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_XTRAP_REQ {
    unsigned char XTrapClientKeyLength;
    unsigned char XTrapClientKey[0];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `XTrapClientKeyLength` | Campo de paquete oficial extraído de PDB TPI |
| `0x001` | `1` | `unsigned char[0]` | `XTrapClientKey` | Campo de paquete oficial extraído de PDB TPI |
