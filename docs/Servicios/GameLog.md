# Especificación Oficial: `GameLog (GameLogDataServer)`


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


## Catálogo Exhaustivo de Estructuras de Logs S2S (`PROTO_NC_LOG_*`)

> **Total de Estructuras Catalogadas:** `44`

> Estructuras resueltas con paridad binaria exacta a partir de `Zone.pdb`.

### `PROTO_NC_LOG_GAME_DATA_TYPE_0`

- **Tamaño Total en Memoria:** `68 bytes` (`0x44`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_0 {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetCharNo;
    unsigned int nTargetID;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
    unsigned int nInt3;
    unsigned __int64 nBigint1;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x030` | `48` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x034` | `52` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |
| `0x038` | `56` | `unsigned int` | `nInt3` | Campo de paquete oficial extraído de PDB TPI |
| `0x03C` | `60` | `unsigned __int64` | `nBigint1` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_1`

- **Tamaño Total en Memoria:** `32 bytes` (`0x20`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_1 {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_2`

- **Tamaño Total en Memoria:** `40 bytes` (`0x28`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_2 {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nInt1;
    unsigned int nInt2;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_3`

- **Tamaño Total en Memoria:** `44 bytes` (`0x2C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_3 {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetCharNo;
    unsigned int nTargetID;
    unsigned int nInt1;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_4`

- **Tamaño Total en Memoria:** `44 bytes` (`0x2C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_4 {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nInt1;
    unsigned int nInt2;
    unsigned int nInt3;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `unsigned int` | `nInt3` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_5`

- **Tamaño Total en Memoria:** `16 bytes` (`0x10`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_5 {
    unsigned int nType;
    unsigned int nCharNo;
    unsigned int nInt1;
    unsigned int nInt2;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x00C` | `12` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_6`

- **Tamaño Total en Memoria:** `44 bytes` (`0x2C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_6 {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nInt1;
    unsigned __int64 nBigint1;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned __int64` | `nBigint1` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_7`

- **Tamaño Total en Memoria:** `44 bytes` (`0x2C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_7 {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetCharNo;
    unsigned int nInt1;
    unsigned int nInt2;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_8`

- **Tamaño Total en Memoria:** `36 bytes` (`0x24`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_8 {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nInt1;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_9`

- **Tamaño Total en Memoria:** `48 bytes` (`0x30`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_9 {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x02C` | `44` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_A`

- **Tamaño Total en Memoria:** `36 bytes` (`0x24`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_A {
    unsigned int nType;
    unsigned int nCharNo;
    unsigned int nTargetCharNo;
    unsigned int nTargetID;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
    unsigned int nInt3;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned int` | `nTargetCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x00C` | `12` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x010` | `16` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nInt3` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_B`

- **Tamaño Total en Memoria:** `60 bytes` (`0x3C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_B {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetCharNo;
    unsigned int nTargetID;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
    unsigned int nInt3;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x030` | `48` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x034` | `52` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |
| `0x038` | `56` | `unsigned int` | `nInt3` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_C`

- **Tamaño Total en Memoria:** `56 bytes` (`0x38`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_C {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetID;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
    unsigned int nInt3;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x02C` | `44` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x030` | `48` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |
| `0x034` | `52` | `unsigned int` | `nInt3` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_D`

- **Tamaño Total en Memoria:** `24 bytes` (`0x18`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_D {
    unsigned int nType;
    unsigned int nCharNo;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x010` | `16` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_E`

- **Tamaño Total en Memoria:** `48 bytes` (`0x30`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_E {
    unsigned int nType;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetID;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x010` | `16` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x02C` | `44` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_F`

- **Tamaño Total en Memoria:** `64 bytes` (`0x40`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_F {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetID;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
    unsigned int nInt3;
    unsigned __int64 nBigint1;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x02C` | `44` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x030` | `48` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |
| `0x034` | `52` | `unsigned int` | `nInt3` | Campo de paquete oficial extraído de PDB TPI |
| `0x038` | `56` | `unsigned __int64` | `nBigint1` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_G`

- **Tamaño Total en Memoria:** `60 bytes` (`0x3C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_G {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetID;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
    unsigned __int64 nBigint1;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x02C` | `44` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x030` | `48` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |
| `0x034` | `52` | `unsigned __int64` | `nBigint1` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_H`

- **Tamaño Total en Memoria:** `36 bytes` (`0x24`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_H {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetCharNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetCharNo` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_I`

- **Tamaño Total en Memoria:** `40 bytes` (`0x28`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_I {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetCharNo;
    unsigned int nTargetID;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_J`

- **Tamaño Total en Memoria:** `52 bytes` (`0x34`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_J {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nTargetID;
    unsigned __int64 nItemKey;
    unsigned int nInt1;
    unsigned int nInt2;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nTargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned __int64` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x02C` | `44` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x030` | `48` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_DATA_TYPE_K`

- **Tamaño Total en Memoria:** `48 bytes` (`0x30`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_DATA_TYPE_K {
    unsigned int nType;
    unsigned int nCharNo;
    Name3 sMap;
    unsigned int nMapX;
    unsigned int nMapY;
    unsigned int nMapZ;
    unsigned int nInt1;
    unsigned int nInt2;
    unsigned __int64 nBigint1;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `nType` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `Name3` | `sMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x014` | `20` | `unsigned int` | `nMapX` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `nMapY` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `nMapZ` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned int` | `nInt1` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned int` | `nInt2` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `unsigned __int64` | `nBigint1` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_ITEM_CW_BREAKATZERO`

- **Tamaño Total en Memoria:** `32 bytes` (`0x20`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_ITEM_CW_BREAKATZERO {
    unsigned int chrregnum;
    SHINE_ITEM_REGISTNUMBER itemkey;
    MAPPOS mappos;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `SHINE_ITEM_REGISTNUMBER` | `itemkey` | Campo de paquete oficial extraído de PDB TPI |
| `0x00C` | `12` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_ITEM_PUT_ON_BELONGED`

- **Tamaño Total en Memoria:** `16 bytes` (`0x10`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_ITEM_PUT_ON_BELONGED {
    unsigned int chrregnum;
    SHINE_ITEM_REGISTNUMBER itemkey;
    ITEM_INVEN iteminvenfrom;
    ITEM_INVEN iteminvento;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `SHINE_ITEM_REGISTNUMBER` | `itemkey` | Campo de paquete oficial extraído de PDB TPI |
| `0x00C` | `12` | `ITEM_INVEN` | `iteminvenfrom` | Campo de paquete oficial extraído de PDB TPI |
| `0x00E` | `14` | `ITEM_INVEN` | `iteminvento` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_ITEM_UPGRADE`

- **Tamaño Total en Memoria:** `43 bytes` (`0x2B`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_ITEM_UPGRADE {
    unsigned int chrregnum;
    MAPPOS mappos;
    SHINE_ITEM_REGISTNUMBER itemkey;
    unsigned char itemoldlevel;
    unsigned char itemnewlevel;
    unsigned short left_item_id;
    unsigned short right_item_id;
    unsigned short middle_item_id;
    ItemOptionStorage::Element randomoption_elemental;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `SHINE_ITEM_REGISTNUMBER` | `itemkey` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned char` | `itemoldlevel` | Campo de paquete oficial extraído de PDB TPI |
| `0x021` | `33` | `unsigned char` | `itemnewlevel` | Campo de paquete oficial extraído de PDB TPI |
| `0x022` | `34` | `unsigned short` | `left_item_id` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned short` | `right_item_id` | Campo de paquete oficial extraído de PDB TPI |
| `0x026` | `38` | `unsigned short` | `middle_item_id` | Campo de paquete oficial extraído de PDB TPI |
| `0x028` | `40` | `ItemOptionStorage::Element` | `randomoption_elemental` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_KQ_ENTER`

- **Tamaño Total en Memoria:** `30 bytes` (`0x1E`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_KQ_ENTER {
    unsigned int chrregnum;
    MAPPOS mappos;
    unsigned short KQ_ID;
    unsigned int KQ_Handle;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned short` | `KQ_ID` | Campo de paquete oficial extraído de PDB TPI |
| `0x01A` | `26` | `unsigned int` | `KQ_Handle` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_KQ_LEAVE`

- **Tamaño Total en Memoria:** `31 bytes` (`0x1F`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_KQ_LEAVE {
    unsigned int chrregnum;
    MAPPOS mappos;
    unsigned short KQ_ID;
    unsigned int KQ_Handle;
    unsigned char Leave_Status;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned short` | `KQ_ID` | Campo de paquete oficial extraído de PDB TPI |
| `0x01A` | `26` | `unsigned int` | `KQ_Handle` | Campo de paquete oficial extraído de PDB TPI |
| `0x01E` | `30` | `unsigned char` | `Leave_Status` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_LEVEL_UP`

- **Tamaño Total en Memoria:** `33 bytes` (`0x21`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_LEVEL_UP {
    unsigned int chrregnum;
    MAPPOS mappos;
    unsigned char level;
    unsigned __int64 currentexp;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned char` | `level` | Campo de paquete oficial extraído de PDB TPI |
| `0x019` | `25` | `unsigned __int64` | `currentexp` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_MINIHOUSE_BUILDING`

- **Tamaño Total en Memoria:** `24 bytes` (`0x18`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_MINIHOUSE_BUILDING {
    unsigned int chrregnum;
    MAPPOS mappos;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_MINIHOUSE_VISIT`

- **Tamaño Total en Memoria:** `28 bytes` (`0x1C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_MINIHOUSE_VISIT {
    unsigned int owner_chrregnum;
    unsigned int visitor_chrregnum;
    MAPPOS mappos;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `owner_chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `visitor_chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_MK2_DEAD`

- **Tamaño Total en Memoria:** `26 bytes` (`0x1A`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_MK2_DEAD {
    unsigned short mobid;
    MAPPOS mappos;
    unsigned int nCharNo;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `mobid` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |
| `0x016` | `22` | `unsigned int` | `nCharNo` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_MK2_START`

- **Tamaño Total en Memoria:** `28 bytes` (`0x1C`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_MK2_START {
    unsigned int chrregnum;
    unsigned char nLevel;
    unsigned char nClass;
    unsigned short mobid;
    MAPPOS mappos;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned char` | `nLevel` | Campo de paquete oficial extraído de PDB TPI |
| `0x005` | `5` | `unsigned char` | `nClass` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned short` | `mobid` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_MKED`

- **Tamaño Total en Memoria:** `34 bytes` (`0x22`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_MKED {
    unsigned int chrregnum;
    MAPPOS mappos;
    unsigned short mobid;
    unsigned __int64 exp;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned short` | `mobid` | Campo de paquete oficial extraído de PDB TPI |
| `0x01A` | `26` | `unsigned __int64` | `exp` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_MK_DROP_ITEM`

- **Tamaño Total en Memoria:** `34 bytes` (`0x22`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_MK_DROP_ITEM {
    unsigned short nItemID;
    unsigned short nLot;
    SHINE_ITEM_REGISTNUMBER nItemKey;
    unsigned short mobid;
    MAPPOS mappos;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned short` | `nItemID` | Campo de paquete oficial extraído de PDB TPI |
| `0x002` | `2` | `unsigned short` | `nLot` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `SHINE_ITEM_REGISTNUMBER` | `nItemKey` | Campo de paquete oficial extraído de PDB TPI |
| `0x00C` | `12` | `unsigned short` | `mobid` | Campo de paquete oficial extraído de PDB TPI |
| `0x00E` | `14` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_MK_GETEXPINFIELD`

- **Tamaño Total en Memoria:** `12 bytes` (`0xC`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_MK_GETEXPINFIELD {
    unsigned int chrregnum;
    unsigned short mobid;
    unsigned short number;
    unsigned int exp;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned short` | `mobid` | Campo de paquete oficial extraído de PDB TPI |
| `0x006` | `6` | `unsigned short` | `number` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned int` | `exp` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GAME_PKED`

- **Tamaño Total en Memoria:** `33 bytes` (`0x21`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GAME_PKED {
    unsigned int chrregnum;
    MAPPOS mappos;
    unsigned int killerdiechrrgnum;
    unsigned int pkcount;
    unsigned char pkyellowtime;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `killerdiechrrgnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `unsigned int` | `pkcount` | Campo de paquete oficial extraído de PDB TPI |
| `0x020` | `32` | `unsigned char` | `pkyellowtime` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_GMAE_LOGOUT_ZONEINFO`

- **Tamaño Total en Memoria:** `41 bytes` (`0x29`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_GMAE_LOGOUT_ZONEINFO {
    unsigned int chrregnum;
    MAPPOS mappos;
    unsigned char StoreType;
    unsigned __int64 CharExp;
    unsigned __int64 CharMoney;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `chrregnum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `MAPPOS` | `mappos` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned char` | `StoreType` | Campo de paquete oficial extraído de PDB TPI |
| `0x019` | `25` | `unsigned __int64` | `CharExp` | Campo de paquete oficial extraído de PDB TPI |
| `0x021` | `33` | `unsigned __int64` | `CharMoney` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_REGENLOCATESAVE_CMD`

- **Tamaño Total en Memoria:** `38 bytes` (`0x26`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_REGENLOCATESAVE_CMD {
    unsigned int ChrRegNum;
    Name3 SaveMap;
    SHINE_XY_TYPE SaveLocate;
    unsigned int KQHandle;
    SHINE_XY_TYPE InvalidLocate;
    unsigned short Flag;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `ChrRegNum` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `Name3` | `SaveMap` | Campo de paquete oficial extraído de PDB TPI |
| `0x010` | `16` | `SHINE_XY_TYPE` | `SaveLocate` | Campo de paquete oficial extraído de PDB TPI |
| `0x018` | `24` | `unsigned int` | `KQHandle` | Campo de paquete oficial extraído de PDB TPI |
| `0x01C` | `28` | `SHINE_XY_TYPE` | `InvalidLocate` | Campo de paquete oficial extraído de PDB TPI |
| `0x024` | `36` | `unsigned short` | `Flag` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_WEDDING_HALL_CANCEL`

- **Tamaño Total en Memoria:** `8 bytes` (`0x8`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_WEDDING_HALL_CANCEL {
    unsigned int CharRegisterNumber;
    unsigned int PartnerRegisterNumber;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `CharRegisterNumber` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `PartnerRegisterNumber` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_WEDDING_HALL_RESERVE`

- **Tamaño Total en Memoria:** `8 bytes` (`0x8`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_WEDDING_HALL_RESERVE {
    unsigned int CharRegisterNumber;
    unsigned int PartnerRegisterNumber;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `CharRegisterNumber` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `PartnerRegisterNumber` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_LOG_WEDDING_HALL_START`

- **Tamaño Total en Memoria:** `8 bytes` (`0x8`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_LOG_WEDDING_HALL_START {
    unsigned int CharRegisterNumber;
    unsigned int PartnerRegisterNumber;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `CharRegisterNumber` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `PartnerRegisterNumber` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_SYSLOG_CHAR_DEATH`

- **Tamaño Total en Memoria:** `22 bytes` (`0x16`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_SYSLOG_CHAR_DEATH {
    unsigned int UserNo;
    unsigned int CharNo;
    unsigned char ZoneNo;
    unsigned char AttackerType;
    unsigned int AttackerID;
    SHINE_XY_TYPE DeathCoord;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `UserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `CharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char` | `ZoneNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `unsigned char` | `AttackerType` | Campo de paquete oficial extraído de PDB TPI |
| `0x00A` | `10` | `unsigned int` | `AttackerID` | Campo de paquete oficial extraído de PDB TPI |
| `0x00E` | `14` | `SHINE_XY_TYPE` | `DeathCoord` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_SYSLOG_CHAR_QUEST_STARTED`

- **Tamaño Total en Memoria:** `13 bytes` (`0xD`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_SYSLOG_CHAR_QUEST_STARTED {
    unsigned int UserNo;
    unsigned int CharNo;
    unsigned short QuestID;
    unsigned char ZoneNo;
    unsigned short QuestStartNPC;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `UserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `CharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned short` | `QuestID` | Campo de paquete oficial extraído de PDB TPI |
| `0x00A` | `10` | `unsigned char` | `ZoneNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x00B` | `11` | `unsigned short` | `QuestStartNPC` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_SYSLOG_CHAR_VICTORY`

- **Tamaño Total en Memoria:** `22 bytes` (`0x16`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_SYSLOG_CHAR_VICTORY {
    unsigned int UserNo;
    unsigned int CharNo;
    unsigned char ZoneNo;
    unsigned char TargetType;
    unsigned int TargetID;
    SHINE_XY_TYPE KillCoord;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `UserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `CharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char` | `ZoneNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `unsigned char` | `TargetType` | Campo de paquete oficial extraído de PDB TPI |
| `0x00A` | `10` | `unsigned int` | `TargetID` | Campo de paquete oficial extraído de PDB TPI |
| `0x00E` | `14` | `SHINE_XY_TYPE` | `KillCoord` | Campo de paquete oficial extraído de PDB TPI |

### `PROTO_NC_SYSLOG_CHAR_ZONE_TRANSITION`

- **Tamaño Total en Memoria:** `11 bytes` (`0xB`)

```cpp
#pragma pack(push, 1)
struct PROTO_NC_SYSLOG_CHAR_ZONE_TRANSITION {
    unsigned int UserNo;
    unsigned int CharNo;
    unsigned char DestZoneNo;
    unsigned char SrcZoneNo;
    unsigned char CharLev;
};
#pragma pack(pop)
```

| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |
| :--- | :--- | :--- | :--- | :--- |
| `0x000` | `0` | `unsigned int` | `UserNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x004` | `4` | `unsigned int` | `CharNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x008` | `8` | `unsigned char` | `DestZoneNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x009` | `9` | `unsigned char` | `SrcZoneNo` | Campo de paquete oficial extraído de PDB TPI |
| `0x00A` | `10` | `unsigned char` | `CharLev` | Campo de paquete oficial extraído de PDB TPI |
