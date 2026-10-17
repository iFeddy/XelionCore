# Especificación Oficial: `Account (AccountServer)`


## Ficha Técnica del Binario
- **Nombre de Binario:** `Account.exe`
- **Símbolos:** `Account.pdb` (Microsoft C/C++ MSF 7.00, TPI v4.10)
- **Ruta del Proyecto Original:** `E:\projectf2\server\3loginserver2\accountserver\`
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


## Catálogo Exhaustivo de Estructuras de Paquetes

> **Total de Estructuras Catalogadas:** `101`

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

### `PROTO_NC_DATA_SEND1<PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK>`

- **Tamaño Total en Memoria:** `12 bytes` (`0xC`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK /* base */ PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK` | `/* base */ PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_DATA_SEND1<PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM>`

- **Tamaño Total en Memoria:** `18 bytes` (`0x12`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_DATA_SEND1<PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM> {
    SIZE_NETCOMMAND1 /* base */ SIZE_NETCOMMAND1;
    PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM /* base */ PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `SIZE_NETCOMMAND1` | `/* base */ SIZE_NETCOMMAND1` | Campo de paquete oficial extraído de PDB TPI |
| `0x003` | `3` | `PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM` | `/* base */ PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_MISC_CONNECTFROMWHERE_ACK`

- **Tamaño Total en Memoria:** `7 bytes` (`0x7`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_CONNECTFROMWHERE_ACK {
    NETPACKETZONEHEADER netpacketzoneheader;
    unsigned char connectfrom;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETZONEHEADER` | `netpacketzoneheader` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned char` | `connectfrom` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK`

- **Tamaño Total en Memoria:** `9 bytes` (`0x9`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK {
    unsigned short hZone;
    PROTO_NC_MISC_CONNECTFROMWHERE_ACK ToZoneAck;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `hZone` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `PROTO_NC_MISC_CONNECTFROMWHERE_ACK` | `ToZoneAck` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ`

- **OpCode Asociado:** `0x082F (Dept: 2, Cmd: 47)`
- **Tamaño Total en Memoria:** `12 bytes` (`0xC`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ {
    unsigned short hZone;
    PROTO_NC_MISC_CONNECTFROMWHERE_REQ FromZoneReq;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `hZone` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `PROTO_NC_MISC_CONNECTFROMWHERE_REQ` | `FromZoneReq` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_CONNECTFROMWHERE_REQ`

- **Tamaño Total en Memoria:** `10 bytes` (`0xA`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_CONNECTFROMWHERE_REQ {
    NETPACKETZONEHEADER netpacketzoneheader;
    unsigned char connectip[4];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETZONEHEADER` | `netpacketzoneheader` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned char[4]` | `connectip` | Dirección IP del cliente o del servidor |

### `PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM`

- **Tamaño Total en Memoria:** `15 bytes` (`0xF`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM {
    unsigned int nUserNo;
    unsigned int nCharNo;
    unsigned char nWorldNo;
    unsigned short nError;
    unsigned short hZone;
    unsigned short nQuestID;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |
| `0x009` | `9` | `unsigned short` | `nError` | Campo de paquete oficial extraído de PDB TPI |
| `0x00B` | `11` | `unsigned short` | `hZone` | Campo de paquete oficial extraído de PDB TPI |
| `0x00D` | `13` | `unsigned short` | `nQuestID` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC`

- **OpCode Asociado:** `0x0838 (Dept: 2, Cmd: 56)`
- **Tamaño Total en Memoria:** `13 bytes` (`0xD`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC {
    unsigned int nUserNo;
    unsigned int nCharNo;
    unsigned char nWorldNo;
    unsigned short hZone;
    unsigned short nQuestID;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |
| `0x009` | `9` | `unsigned short` | `hZone` | Campo de paquete oficial extraído de PDB TPI |
| `0x00B` | `11` | `unsigned short` | `nQuestID` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_EVENT_HIT3_ADD_CASH`

- **OpCode Asociado:** `0x0820 (Dept: 2, Cmd: 32)`
- **Tamaño Total en Memoria:** `28 bytes` (`0x1C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_EVENT_HIT3_ADD_CASH {
    unsigned char UserID[20];
    unsigned __int64 AddCash;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[20]` | `UserID` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned __int64` | `AddCash` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_MISC_WEB_DB_KEY_ACK`

- **Tamaño Total en Memoria:** `40 bytes` (`0x28`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_WEB_DB_KEY_ACK {
    unsigned short Error;
    unsigned short nWorldSessionHandle;
    unsigned int nUserNo;
    unsigned char WebKey[32];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `Error` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short` | `nWorldSessionHandle` | Identificador o información del WorldServer |
| `0x004` | `4` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char[32]` | `WebKey` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_MISC_WEB_DB_KEY_REQ`

- **OpCode Asociado:** `0x081D (Dept: 2, Cmd: 29)`
- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_MISC_WEB_DB_KEY_REQ {
    unsigned short nWorldSessionHandle;
    unsigned int nUserNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nWorldSessionHandle` | Identificador o información del WorldServer |
| `0x002` | `2` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_USER_CH_PASSWORD_CHECK_ACK_SEND`

- **OpCode Asociado:** `0x0C54 (Dept: 3, Cmd: 84)`
- **Tamaño Total en Memoria:** `290 bytes` (`0x122`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CH_PASSWORD_CHECK_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_CH_PASSWORD_CHECK_ACK ack;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_CH_PASSWORD_CHECK_ACK` | `ack` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_CH_PASSWORD_CHECK_REQ`

- **OpCode Asociado:** `0x0C53 (Dept: 3, Cmd: 83)`
- **Tamaño Total en Memoria:** `274 bytes` (`0x112`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_CH_PASSWORD_CHECK_REQ {
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

### `PROTO_NC_USER_GER_PASSWORD_CHECK_ACK_SEND`

- **OpCode Asociado:** `0x0C3A (Dept: 3, Cmd: 58)`
- **Tamaño Total en Memoria:** `290 bytes` (`0x122`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_GER_PASSWORD_CHECK_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_GER_PASSWORD_CHECK_ACK ack;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_GER_PASSWORD_CHECK_ACK` | `ack` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_GER_PASSWORD_CHECK_REQ`

- **OpCode Asociado:** `0x0C39 (Dept: 3, Cmd: 57)`
- **Tamaño Total en Memoria:** `274 bytes` (`0x112`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_GER_PASSWORD_CHECK_REQ {
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

### `PROTO_NC_USER_JP_PASSWORD_CHECK_REQ`

- **OpCode Asociado:** `0x0C4B (Dept: 3, Cmd: 75)`
- **Tamaño Total en Memoria:** `336 bytes` (`0x150`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_JP_PASSWORD_CHECK_REQ {
    NETPACKETHEADER netpacketheader;
    char sUserName[24];
    char sPassword[36];
    char sChannel[254];
    char sUserIP[20];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `char[24]` | `sUserName` | Nombre de cuenta o usuario (string) |
| `0x01A` | `26` | `char[36]` | `sPassword` | Contraseña en texto plano o hash del cliente |
| `0x03E` | `62` | `char[254]` | `sChannel` | Campo de paquete oficial extraído de PDB TPI |
| `0x13C` | `316` | `char[20]` | `sUserIP` | Dirección IP del cliente o del servidor |

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

### `PROTO_NC_USER_LOGIN_NETMARBLE_DB_REQ`

- **OpCode Asociado:** `0x0C1E (Dept: 3, Cmd: 30)`
- **Tamaño Total en Memoria:** `278 bytes` (`0x116`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGIN_NETMARBLE_DB_REQ {
    NETPACKETHEADER netpacketheader;
    unsigned char UniID[20];
    unsigned char UserID[256];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `unsigned char[20]` | `UniID` | Campo de paquete oficial extraído de PDB TPI |
| `0x016` | `22` | `unsigned char[256]` | `UserID` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_LOGOUT_DB`

- **Tamaño Total en Memoria:** `9 bytes` (`0x9`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_LOGOUT_DB {
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

### `PROTO_NC_USER_PASSWORD_CHECK_ACK_SEND`

- **OpCode Asociado:** `0x0C08 (Dept: 3, Cmd: 8)`
- **Tamaño Total en Memoria:** `290 bytes` (`0x122`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_PASSWORD_CHECK_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_PASSWORD_CHECK_ACK ack;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_PASSWORD_CHECK_ACK` | `ack` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_USER_RETURN_CHECK_ACK`

- **OpCode Asociado:** `0x0C61 (Dept: 3, Cmd: 97)`
- **Tamaño Total en Memoria:** `11 bytes` (`0xB`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_RETURN_CHECK_ACK {
    unsigned short nHandle;
    unsigned int nUserNo;
    unsigned char bReturnUser;
    int tLastSec;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned char` | `bReturnUser` | Campo de paquete oficial extraído de PDB TPI |
| `0x007` | `7` | `int` | `tLastSec` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_RETURN_CHECK_REQ`

- **OpCode Asociado:** `0x0C60 (Dept: 3, Cmd: 96)`
- **Tamaño Total en Memoria:** `6 bytes` (`0x6`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_RETURN_CHECK_REQ {
    unsigned short nHandle;
    unsigned int nUserNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_SET_RETURN_CMD`

- **OpCode Asociado:** `0x0C64 (Dept: 3, Cmd: 100)`
- **Tamaño Total en Memoria:** `4 bytes` (`0x4`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_SET_RETURN_CMD {
    unsigned int nUserNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_USER_TEENAGER_REQ`

- **OpCode Asociado:** `0x0C26 (Dept: 3, Cmd: 38)`
- **Tamaño Total en Memoria:** `258 bytes` (`0x102`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_TEENAGER_REQ {
    char sUserID[256]; /* Name256Byte */
    unsigned short nHandle;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `char[256] /* Name256Byte */` | `sUserID` | Campo de paquete oficial extraído de PDB TPI |
| `0x100` | `256` | `unsigned short` | `nHandle` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_TEENAGER_SET_CMD`

- **OpCode Asociado:** `0x0C28 (Dept: 3, Cmd: 40)`
- **Tamaño Total en Memoria:** `54 bytes` (`0x36`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_TEENAGER_SET_CMD {
    unsigned char sUserKey[50];
    unsigned int nUserPlayRemainTime;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char[50]` | `sUserKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x032` | `50` | `unsigned int` | `nUserPlayRemainTime` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_USER_TW_PASSWORD_CHECK_ACK_SEND`

- **OpCode Asociado:** `0x0C3F (Dept: 3, Cmd: 63)`
- **Tamaño Total en Memoria:** `11 bytes` (`0xB`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_TW_PASSWORD_CHECK_ACK_SEND {
    unsigned char size;
    NETCOMMAND /* union (command:10, dept:6) */ netcmd;
    PROTO_NC_USER_TW_PASSWORD_CHECK_ACK ack;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned char` | `size` | Longitud del payload del paquete en bytes |
| `0x001` | `1` | `NETCOMMAND /* union (command:10, dept:6) */` | `netcmd` | OpCode del protocolo (Dept:6 bits, Cmd:10 bits) |
| `0x003` | `3` | `PROTO_NC_USER_TW_PASSWORD_CHECK_ACK` | `ack` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USER_TW_PASSWORD_CHECK_REQ`

- **OpCode Asociado:** `0x0C3E (Dept: 3, Cmd: 62)`
- **Tamaño Total en Memoria:** `102 bytes` (`0x66`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USER_TW_PASSWORD_CHECK_REQ {
    NETPACKETHEADER netpacketheader;
    char sUserName[36];
    char sSerial[36];
    char sUID[20];
    char sSID[8];
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `NETPACKETHEADER` | `netpacketheader` | Cabecera interna de sesión (ClientHandle uint16) |
| `0x002` | `2` | `char[36]` | `sUserName` | Nombre de cuenta o usuario (string) |
| `0x026` | `38` | `char[36]` | `sSerial` | Campo de paquete oficial extraído de PDB TPI |
| `0x04A` | `74` | `char[20]` | `sUID` | Campo de paquete oficial extraído de PDB TPI |
| `0x05E` | `94` | `char[8]` | `sSID` | Campo de paquete oficial extraído de PDB TPI |

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

### `PROTO_NC_USP_USER_CHARACTER_DELETE`

- **Tamaño Total en Memoria:** `9 bytes` (`0x9`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USP_USER_CHARACTER_DELETE {
    unsigned int nUserNo;
    unsigned char nWorldNo;
    unsigned int nCharNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |
| `0x005` | `5` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_USP_USER_CHARACTER_INSERT`

- **Tamaño Total en Memoria:** `50 bytes` (`0x32`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_USP_USER_CHARACTER_INSERT {
    unsigned int nUserNo;
    unsigned char nWorldNo;
    unsigned int nCharNo;
    char sCharName[40];
    unsigned char nClass;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nUserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `nWorldNo` | Identificador o información del WorldServer |
| `0x005` | `5` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `char[40]` | `sCharName` | Campo de paquete oficial extraído de PDB TPI |
| `0x031` | `49` | `unsigned char` | `nClass` | Campo de paquete oficial extraído de PDB TPI |
