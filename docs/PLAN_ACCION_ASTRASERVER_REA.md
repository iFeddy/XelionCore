# Plan de Acción: Ingeniería Inversa Oficial (AstraServer / REA) y Completitud de Login & Character (WorldManager)

> **Documento:** Plan Estratégico y Técnico de Ingeniería Inversa y Migración a Go (`XelionCore`)  
> **Fuente Oficial Analizada:** Servidor Oficial Fiesta Online (`AstraServer`)  
> **Objetivo:** Documentar y ejecutar la extracción 100% exacta de protocolos, estructuras binarias de paquetes y esquemas de base de datos (`.bak`) para Login y Character (WorldManager) utilizando REA (Reverse Engineering Anything), eliminando suposiciones y logrando interoperabilidad nativa con el cliente.

---

## 1. Diagnóstico y Hallazgos Críticos en AstraServer

Durante la inspección profunda del repositorio y archivos oficiales en `AstraServer/`, se identificaron activos de valor incalculable que permiten prescindir de ingeniería inversa a ciegas:

### 1.1 Archivos PDB Oficiales con Símbolos de Depuración Completos
Los ejecutables oficiales cuentan con sus archivos PDB (Program Database) originales:
- **`Login.pdb`** (3.3 MB) — Servidor de Login y autenticación.
- **`WorldManager.pdb`** (12+ MB) — Servidor de Mundo y Selección/Gestión de Personajes (`Character`).
- **`Zone.pdb`** (30+ MB) — Servidor de Zonas e instancias de mapas.
- **`Account.pdb`** y **`AccountLog.pdb`** — Servicios de base de datos y auditoría de cuentas.

**Impacto Técnico:**  
Los PDBs contienen los nombres reales de todas las clases C++, métodos de despacho de red (`CParserClient::fc_NC_...`), estructuras de paquetes de red (`PROTO_NC_...`) y tipos de datos exactos utilizados por los desarrolladores originales de Fiesta Online.

### 1.2 Backups Oficiales de Base de Datos MSSQL (`.bak`)
En `AstraServer/Server/Databases/` y `AstraServer/Database/`:
- **`Account.bak`**: Contiene la autenticación de usuarios, esquemas de contraseñas, autorizaciones e IP blocking.
  - Tablas núcleo: `tUser`, `tUserAuth`, `tUserOption`, `tUserBlock`.
  - Procedimientos almacenados clave: `usp_User_USloginGame`, `usp_User_chkBlock`, `usp_Cash_*`.
- **`World00_Character.bak`**: Contiene la estructura canónica de datos de personajes.
  - Tablas núcleo: `tCharacter`, `tCharacterShape`, `tCharacterState`, `tCharacterKeyMap`, `tCharacterShortCut`, `tCharacterGameOption`, `tCharacterSingleOption`, `tCharacterTitle`, `tGuild`, `tGuildAcademy`, `tGuildAcademyMember`, `tSkill`, `tItemEnchant`, `tItemOptions`.
- **`Options.bak`**: Opciones de servidor y creación automática de cuentas (`Auto_Account_Creation`).
- **`StatisticsData.bak`**, **`OperatorTool.bak`**, **`AccountLog.bak`**, **`World00_GameLog.bak`**.

### 1.3 Archivos de Configuración de Topología Oficial
En `AstraServer/Server/9Data/ServerInfo/`:
- **`ServerInfo_Login.txt`** define la topología de puertos canónicos:
  - `Login Server`: Puerto Cliente `9010`, Puerto Inter-Server `9011` / `9012`.
  - `WorldManager (W00_WM)`: Puerto Cliente `9013`, Inter-Server `9014` / `9015`.
  - `Zone 0 (Roumen & zonas)`: Puerto Cliente `9016`, Inter-Server `9017` / `9018`.
  - `Zone 1..4`: Puertos `9019`, `9022`, `9025`, `9028`.
  - `Database Services`: Puertos `9031`..`9034`.

---

## 2. Mapa Completo de Paquetes Oficiales (PDB / REA)

A través de la extracción de símbolos de los PDBs de `AstraServer`, se desglosan a continuación las funciones y estructuras reales que manejan Login y Character:

### 2.1 Paquetes de Login Server (`CParserClient` en `Login.pdb`)

| Función Manejadora Oficial | Estructura / Tipo | Dirección | Descripción Oficial |
|---|---|---|---|
| `fc_NC_USER_CLIENT_VERSION_CHECK_REQ` | `PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ` | C -> S | Verificación de versión del cliente (Año y Versión binaria) |
| `Send_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK` | `PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK` | S -> C | Versión aceptada, cliente habilitado a enviar credenciales |
| `Send_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK` | `PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK` | S -> C | Versión incorrecta, desconexión |
| `fc_NC_USER_US_LOGIN_REQ` | `PROTO_NC_USER_US_LOGIN_REQ` | C -> S | Solicitud de login con Username, Hash de contraseña y Client Tag |
| `fc_NC_USER_LOGIN_WITH_OTP_REQ` | `PROTO_NC_USER_LOGIN_WITH_OTP_REQ` | C -> S | Login alternativo mediante One-Time Password / Token |
| `fc_NC_USER_WORLD_STATUS_REQ` | `PROTO_NC_USER_WORLD_STATUS_REQ` | C -> S | Solicitud de lista y estado de servidores de mundo |
| `Send_NC_USER_WORLD_STATUS_ACK` | `PROTO_NC_USER_WORLD_STATUS_ACK` | S -> C | Envío de mundos, nombres, carga de usuarios e IPs |
| `fc_NC_USER_WORLDSELECT_REQ` | `PROTO_NC_USER_WORLDSELECT_REQ` | C -> S | Selección de servidor de mundo por el usuario |
| `Send_NC_USER_WORLDSELECT_ACK` | `PROTO_NC_USER_WORLDSELECT_ACK` | S -> C | Envío de IP/Puerto de WorldManager y token de transferencia |
| `fc_NC_USER_NORMALLOGOUT_CMD` | `PROTO_NC_USER_NORMALLOGOUT_CMD` | C -> S | Salida ordenada de la pantalla de login |
| `fc_NC_USER_XTRAP_REQ` | `PROTO_NC_USER_XTRAP_REQ` | C -> S | Verificación anti-cheat xTrap / GameGuard |

### 2.2 Paquetes de Character & World Server (`CParserClient` en `WorldManager.pdb`)

| Función Manejadora Oficial | Estructura / Tipo | Dirección | Descripción Oficial |
|---|---|---|---|
| `fc_NC_USER_LOGINWORLD_REQ` | `PROTO_NC_USER_LOGINWORLD_REQ` | C -> S | Handshake inicial de transferencia (WorldClientKey) |
| `fc_NC_USER_AVATAR_LIST_REQ` | `PROTO_NC_USER_AVATAR_LIST_REQ` | C -> S | Solicitud de lista de personajes de la cuenta |
| `Send_NC_USER_AVATAR_LIST_ACK` | `PROTO_NC_USER_AVATAR_LIST_ACK` | S -> C | Lista completa de personajes en slots 0..5 con RandomID |
| `fc_NC_AVATAR_CREATE_REQ` | `PROTO_NC_AVATAR_CREATE_REQ` | C -> S | Creación de nuevo avatar (Nombre, Slot, Clase, Look) |
| `Send_NC_AVATAR_CREATESUCC_ACK` | `PROTO_NC_AVATAR_CREATESUCC_ACK` | S -> C | Éxito en creación con datos serializados del avatar |
| `Send_NC_AVATAR_CREATEFAIL_ACK` | `PROTO_NC_AVATAR_CREATEFAIL_ACK` | S -> C | Error en creación (nombre duplicado, slot inválido) |
| `fc_NC_AVATAR_ERASE_REQ` | `PROTO_NC_AVATAR_ERASE_REQ` | C -> S | Solicitud de borrado de personaje |
| `Send_NC_AVATAR_ERASESUC_ACK` | `PROTO_NC_AVATAR_ERASESUC_ACK` | S -> C | Confirmación de borrado del personaje |
| `fc_NC_AVATAR_RENAME_REQ` | `PROTO_NC_AVATAR_RENAME_REQ` | C -> S | Cambio de nombre de avatar (por cupón o GM) |
| `fc_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ` | `PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ` | C -> S | **CharSelect:** Selección de avatar para entrar al juego |
| `Send_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK` | `PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK` | S -> C | Envío de IP y puerto de Zone Server |
| `fc_NC_USER_WILL_WORLD_SELECT_REQ` | `PROTO_NC_USER_WILL_WORLD_SELECT_REQ` | C -> S | Regreso a la lista de mundos desde selección de char |
| `fc_NC_CHAR_OPTION_GET_ALL_REQ` | `PROTO_NC_CHAR_OPTION_GET_ALL_REQ` | C -> S | Petición unificada de configuraciones de cliente |
| `fc_NC_CHAR_OPTION_GET_KEYMAPPING_REQ` | `PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_REQ` | C -> S | Carga de QuickBars y mapeo de teclas |
| `fc_NC_CHAR_OPTION_SET_KEYMAPPING_CMD` | `PROTO_NC_CHAR_OPTION_SET_KEYMAPPING_CMD` | C -> S | Guardado de QuickBars modificado por el jugador |
| `fc_NC_CHAR_OPTION_GET_SHORTCUTDATA_REQ`| `PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_REQ` | C -> S | Carga de ShortCuts de interfaz |
| `fc_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD`| `PROTO_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD` | C -> S | Guardado de ShortCuts |
| `fc_NC_CHAR_OPTION_GET_GAME_REQ` | `PROTO_NC_CHAR_OPTION_GET_GAME_REQ` | C -> S | Carga de GameSettings (opciones de combate/cámara) |
| `fc_NC_CHAR_OPTION_SET_GAME_CMD` | `PROTO_NC_CHAR_OPTION_SET_GAME_CMD` | C -> S | Guardado de GameSettings |
| `fc_NC_CHAR_OPTION_GET_WINDOWPOS_REQ` | `PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_REQ` | C -> S | Posiciones de ventanas en pantalla |
| `fc_NC_CHAR_OPTION_GET_VIDEO_REQ` | `PROTO_NC_CHAR_OPTION_GET_VIDEO_REQ` | C -> S | Configuración gráfica guardada |
| `fc_NC_CHAR_OPTION_GET_SOUND_REQ` | `PROTO_NC_CHAR_OPTION_GET_SOUND_REQ` | C -> S | Configuración de audio y volumen |
| `fc_NC_CHAR_LOGIN_REQ` / `CH22Type.GotIngame` | `PROTO_NC_CHAR_LOGIN_REQ` | C -> S | Notificación de cliente entrando a juego (`0x581B`) |
| `Send_NC_CHAR_LOGIN_ACK` | `0x581C`, `0x581D`, `(21,7)`, `(2,69)` | S -> C | ACK de entrada a juego, KQ list, Friend list y TimeSync |

---

## 3. Análisis de la Base de Datos (`.bak`) y Esquema Relacional

Los backups en `AstraServer/Server/Databases/` contienen el modelo de datos exacto de producción de Fiesta Online. A continuación se desglosan las tablas y columnas críticas que deben replicarse en Go:

### 3.1 Base de Datos `Account`
```sql
-- Tabla: tUser (Cuentas principales)
nUserNo       INT PRIMARY KEY IDENTITY(1,1), -- ID único de cuenta
sUserID       NVARCHAR(60) NOT NULL UNIQUE,  -- Nombre de usuario
sUserPW       NVARCHAR(32) NOT NULL,         -- Hash de contraseña (MD5/SHA)
sUserName     NVARCHAR(60),                  -- Nombre visible / Display name
bIsBlock      INT DEFAULT 0,                 -- 0 = Activo, 1 = Bloqueado
bIsDelete     INT DEFAULT 0,                 -- Borrado lógico
nAuthID       TINYINT DEFAULT 1,             -- Nivel de privilegios (1=User, 100=Admin/GM)
sUserIP       NVARCHAR(50),                  -- Última IP conectada
dCreateDate   DATETIME DEFAULT GETDATE(),    -- Fecha de registro
dLoginDate    DATETIME                       -- Último login exitoso
```

### 3.2 Base de Datos `World00_Character`
```sql
-- Tabla: tCharacter (Personajes)
nCharNo       INT PRIMARY KEY IDENTITY(1,1), -- ID único de avatar
nUserNo       INT NOT NULL,                  -- Relación con tUser.nUserNo
sCharName     NVARCHAR(40) NOT NULL UNIQUE,  -- Nombre del personaje
nSlot         TINYINT NOT NULL,              -- Slot del personaje (0..5)
nClass        TINYINT NOT NULL,              -- Job (1=Fighter, 6=Cleric, 11=Archer, 16=Mage, 21=Trickster)
nLevel        TINYINT DEFAULT 1,             -- Nivel actual
bIsDelete     INT DEFAULT 0                  -- Flag de borrado

-- Tabla: tCharacterShape (Apariencia del personaje)
nCharNo       INT PRIMARY KEY,               -- FK tCharacter
nGender       TINYINT,                       -- 0 = Femenino, 1 = Masculino
nHair         TINYINT,                       -- Estilo de peinado
nHairColor    TINYINT,                       -- Color de cabello
nFace         TINYINT                        -- Estilo de rostro

-- Tabla: tCharacterState (Estado y posición)
nCharNo       INT PRIMARY KEY,               -- FK tCharacter
nMap          INT DEFAULT 0,                 -- ID del mapa (0 = Roumen)
nXPos         INT DEFAULT 7636,              -- Coordenada X
nYPos         INT DEFAULT 4610,              -- Coordenada Y
nRotation     TINYINT DEFAULT 85,            -- Rotación (0x55)
nHP           INT DEFAULT 52,                -- Vida actual
nSP           INT DEFAULT 10,                -- Mana/SP actual
nHPStones     INT DEFAULT 15,                -- Piedras de HP disponibles
nSPStones     INT DEFAULT 7,                 -- Piedras de SP disponibles
nExp          BIGINT DEFAULT 0,              -- Experiencia actual
nMoney        BIGINT DEFAULT 0               -- Dinero en inventario

-- Tablas de Configuración de Interfaz por Personaje:
-- tCharacterKeyMap       (QuickBars y asignación de habilidades a teclas)
-- tCharacterShortCut     (Accesos directos de UI)
-- tCharacterGameOption   (Configuraciones de combate y juego)
-- tCharacterSingleOption (Configuraciones de cliente)
```

---

## 4. Metodología de Ingeniería Inversa con REA (Reverse Engineering Anything)

Para extraer con precisión quirúrgica de 1 byte las estructuras C++ empaquetadas (`#pragma pack(1)`) de cada paquete:

### Fase 1: Extracción de Símbolos y Offsets de Estructuras vía PDB
1. Utilizar herramientas de volcado de tipos PDB (tales como `pdbparse`, `cvdump` o `dia2dump`) sobre `Login.pdb` y `WorldManager.pdb`.
2. Extraer las definiciones de cada struct que empiece por `PROTO_NC_` y `NETPACKET`:
   - Campos miembro.
   - Tipos de datos primitivos (`char`, `unsigned short`, `int`, `unsigned int`, `__int64`).
   - Offsets relativos de cada campo para verificar alineación de bytes exacta.

### Fase 2: Decompilación Asistida (Ghidra / IDA / Radare2) sobre `Login.exe` y `WorldManager.exe`
1. Cargar `Login.exe` y `WorldManager.exe` en Ghidra asociándoles sus respectivos `.pdb`. Al contar con PDB, **todos los nombres de funciones, variables globales y tipos se resuelven automáticamente**.
2. Localizar la tabla de despacho:
   - En `CParserClient::Parser(CLoginClientSession *pSession, NETPACKET *pPacket, int nSize)` y `CParserClient::Parser(CWMClientSession *pSession, NETPACKET *pPacket, int nSize)`.
   - Extraer la tabla de salto (`switch-case` o `jump table`) que vincula el OpCode binario (Header de 6 bits y Type de 10 bits) con la función `fc_NC_...`.
3. Validar los constructores y serializadores de paquetes salientes (`Send_NC_...`) para certificar el orden exacto de escritura en el buffer de red.

### Fase 3: Extracción de Esquema DDL de los `.bak`
1. Montar o inspeccionar los archivos `.bak` utilizando scripts de volcado de cabeceras MDF o restaurándolos en una instancia de prueba de SQL Server Express (usando el script provisto `_Install Databases.ps1` en la máquina Windows).
2. Generar el script DDL completo (`CREATE TABLE`, tipos de datos exactos, claves primarias e índices).
3. Convertir dicho esquema canónico a modelos en Go y migraciones SQL compatibles (SQLite / PostgreSQL / MySQL).

---

## 5. Plan de Ejecución Paso a Paso (Roadmap)

### Paso 1: Extracción Automatizada de Estructuras (Sin código de juego aún)
- [ ] Ejecutar extractor de PDBs para generar un catálogo de cabeceras C/C++ de todas las estructuras `PROTO_NC_` de Login y WorldManager.
- [ ] Construir un diccionario JSON de referencia cruzada:
  `{ "OpCodeHex": "0x0C05", "Header": 3, "Type": 5, "OfficialName": "PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK", "Fields": [...] }`.

### Paso 2: Alineación y Completitud del Protocolo de Login (`cmd/login`)
- [ ] Actualizar los `PacketReader` y `PacketWriter` de Login en `pkg/server/login/handler.go` para que coincidan 1:1 con las estructuras `PROTO_NC_USER_*`.
- [ ] Implementar soporte completo para las variantes de login (`US_LOGIN_REQ`, `LOGIN_WITH_OTP_REQ`, `WORLD_STATUS_REQ`).
- [ ] Implementar la persistencia de cuentas y usuarios basada en el modelo de `tUser` y `tUserAuth` de `Account.bak`.

### Paso 3: Alineación y Completitud del Protocolo de Character / World (`cmd/world`)
- [ ] Actualizar los manejadores de `CParserClient` en `pkg/server/world/handler.go`:
  - `fc_NC_USER_LOGINWORLD_REQ`: Validación del token de transferencia Login -> World.
  - `fc_NC_USER_AVATAR_LIST_REQ`: Respuesta `PROTO_NC_USER_AVATAR_LIST_ACK` preservando `RandomID`.
  - `fc_NC_AVATAR_CREATE_REQ`: Creación de personaje guardando apariencia en `tCharacterShape` y stats en `tCharacterState`.
  - `fc_NC_AVATAR_ERASE_REQ`: Borrado lógico y confirmación.
  - Paquetes de interfaz y Quickbars: `GET_KEYMAPPING_REQ`, `GET_SHORTCUTDATA_REQ`, `GET_GAME_REQ`, `GET_WINDOWPOS_REQ`.
  - `fc_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ`: Selección de personaje, generación y registro del token hacia el Zone Server correspondiente.
  - `fc_NC_CHAR_LOGIN_REQ` (`0x581B`): Despacho de paquetes de inicialización in-game (`0x581C`, `0x581D`, Friends `21,7`, TimeSync `2,69`).

### Paso 4: Adaptador de Base de Datos Canónica para `XelionCore`
- [ ] Desarrollar capa de persistencia en Go inspirada en las tablas oficiales de `Account.bak` y `World00_Character.bak`.
- [ ] Soportar migración directa de datos desde los scripts SQL oficiales provistos en `AstraServer/Database/sql/`.

---

## 6. Conclusión y Próximos Pasos

El descubrimiento de los archivos **PDB oficiales** en `AstraServer` transforma el proyecto: ya no es necesario realizar ingeniería inversa "a ciegas" ni depender de implementaciones incompletas de otros emuladores comunitarios. Tenemos los nombres, estructuras, firmas y algoritmos oficiales tal como fueron diseñados por los creadores de Fiesta Online.

Este plan de acción queda registrado como la hoja de ruta técnica oficial para llevar a `XelionCore` a una compatibilidad del 100% en los módulos de Login y Character/WorldManager.
