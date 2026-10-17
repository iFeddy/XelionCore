# Catálogo Completo de Procedimientos Almacenados Oficiales - Xelion Legacy (573 Stored Procedures)

Este documento contiene el inventario exhaustivo de los **573 procedimientos almacenados (Stored Procedures)** programados originalmente por OnsOn Soft para Fiesta Online / AstraServer, extraídos directamente de las páginas internas de los backups oficiales (`.bak`), **excluyendo la totalidad de procedimientos internos de Microsoft SQL Server** (`sp_MS...`, `sp_help...`, `xp_...`, etc.).

En la arquitectura clásica de Fiesta Online, estos procedimientos constituyen la **API de comunicación interna** entre los servidores C++ (`Account.exe`, `Character.exe`, `World.exe`, `Zone.exe`, `GameLog.exe`, `OperatorTool.exe`) y la base de datos a través de ODBC.

> [!NOTE]
> Para el esquema físico de tablas, columnas, tipos de datos y relaciones de claves foráneas (FK), consulta [docs/TablasLegacy.md](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/TablasLegacy.md) (176 tablas oficiales).

> [!IMPORTANT]
> **Respaldo Completo de Código Fuente SQL en el Repositorio**:
> Para garantizar que el código SQL nunca se pierda si se borran o dañan los archivos `.bak`, el código fuente completo de los **573 procedimientos almacenados** ha sido extraído y versionado en texto plano en la carpeta [docs/sql/procedures/](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/):
> * **Scripts Consolidados por Base de Datos**:
>   * [`Account_all_procedures.sql`](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/Account_all_procedures.sql) (77 procedures)
>   * [`World00_Character_all_procedures.sql`](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/World00_Character_all_procedures.sql) (408 procedures)
>   * [`AccountLog_all_procedures.sql`](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/AccountLog_all_procedures.sql) (27 procedures)
>   * [`World00_GameLog_all_procedures.sql`](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/World00_GameLog_all_procedures.sql) (26 procedures)
>   * [`StatisticsData_all_procedures.sql`](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/StatisticsData_all_procedures.sql) (18 procedures)
>   * [`OperatorTool_all_procedures.sql`](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/sql/procedures/OperatorTool_all_procedures.sql) (17 procedures)
> * **Archivos Individuales**: Cada procedimiento cuenta además con su archivo `.sql` individual dentro de su subcarpeta correspondiente (ej. `docs/sql/procedures/World00_Character/p_Char_Login.sql`).

---

## Resumen General por Base de Datos

| Base de Datos | Archivo Backup | Cantidad de Procedures | Módulos Clave Cubiertos |
| :--- | :--- | :---: | :--- |
| **`World00_Character`** | `World00_Character.bak` | **408** | Ciclo de vida de personaje, inventario, habilidades, gremios, academia, misiones, bodas, mascotas, casino |
| **`Account`** | `Account.bak` | **77** | Autenticación, control de accesos, facturación, tienda web, saldo de Cash |
| **`AccountLog`** | `AccountLog.bak` | **27** | Auditoría de login/logout, denuncias a GM, bloqueos de IP, eventos de cuenta |
| **`World00_GameLog`** | `World00_GameLog.bak` | **26** | Logs de eventos de juego (drops, muertes, compras, trades, mejoras, comandos GM) |
| **`StatisticsData`** | `StatisticsData.bak` | **18** | Concurrencia (CCU), usuarios activos diarios (DAU), demografía y progresión |
| **`OperatorTool`** | `OperatorTool.bak` | **17** | Gestión de GameMasters, registro de auditoría de comandos y permisos |
| **TOTAL** | | **573** | |

---

## Account (`Account.bak`) - Autenticación, Facturación y Cuentas (77 Procedures)

| # | Nombre del Stored Procedure | Servicio Invocador | Propósito / Operación en el Servidor |
| :---: | :--- | :--- | :--- |
| 1 | `p_UserGetAuthList` | AccountServer | Autenticación de credenciales de usuario y verificación de servidor |
| 2 | `p_UserSetAuthLoginable` | AccountServer | Autenticación de credenciales de usuario y verificación de servidor |
| 3 | `usp_Cash_getCash` | AccountServer | Consulta de saldo de Cash / Moneda premium |
| 4 | `usp_Cash_setCash` | AccountServer | Acreditación de saldo de Cash por recarga o evento |
| 5 | `usp_Cash_setCash_EVENT_L20` | AccountServer | Acreditación de saldo de Cash por recarga o evento |
| 6 | `usp_Cash_useCash` | AccountServer | Deducción de saldo de Cash por compra en Item Shop |
| 7 | `usp_Charge_ItemDel` | AccountServer | Gestión de ítems comprados en tienda web pendientes de entrega |
| 8 | `usp_Charge_ItemDraw` | AccountServer | Gestión de ítems comprados en tienda web pendientes de entrega |
| 9 | `usp_Charge_ItemInsert` | AccountServer | Gestión de ítems comprados en tienda web pendientes de entrega |
| 10 | `usp_Charge_ItemSelectList` | AccountServer | Gestión de ítems comprados en tienda web pendientes de entrega |
| 11 | `usp_Item_getBuyCount` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 12 | `usp_Item_getBuyList` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 13 | `usp_Item_getCategoryName` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 14 | `usp_Item_getEventList` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 15 | `usp_Item_getGroup` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 16 | `usp_Item_getGroupList` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 17 | `usp_Item_getGroupPrice` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 18 | `usp_Item_getItem` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 19 | `usp_Item_getOrder` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 20 | `usp_Item_getOrderList` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 21 | `usp_Item_getSubCategory` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 22 | `usp_Item_getTotalCount` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 23 | `usp_Item_getUser` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 24 | `usp_Item_setBuyOrder` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 25 | `usp_Item_setCategory` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 26 | `usp_Item_setItem` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 27 | `usp_Item_setOrder` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 28 | `usp_Item_setSerialKey` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 29 | `usp_OPT_getBuyList` | AccountServer | Búsqueda y opciones de usuario para herramientas administrativas |
| 30 | `usp_OPT_getPayAgreement` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 31 | `usp_OPT_getUserCount` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 32 | `usp_OPT_getUserID` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 33 | `usp_OPT_getUserNo` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 34 | `usp_OPT_searchUser` | AccountServer | Gestión de información de usuario |
| 35 | `usp_OPT_setPayAgreement` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 36 | `usp_OPT_setUserAuthID` | AccountServer | Autenticación de credenciales de usuario y verificación de servidor |
| 37 | `usp_OPT_setUserBlock` | AccountServer | Actualización de información o estado de la cuenta |
| 38 | `usp_Pay_chkLimit` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 39 | `usp_Pay_getAgreement` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 40 | `usp_Pay_getOrder` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 41 | `usp_Pay_setAgreement` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 42 | `usp_Pay_setOrder` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 43 | `usp_Pay_setResult` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 44 | `usp_ReturnUser_GetDatediffLastlog` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 45 | `usp_ReturnUser_SetUserItem` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 46 | `usp_Stat_sellitem` | AccountServer | Consulta de catálogo de ítems de la tienda web y precios grupales |
| 47 | `usp_Tutorial_getCharInfo` | AccountServer | Operación de base de datos |
| 48 | `usp_Tutorial_getIsSkip` | AccountServer | Operación de base de datos |
| 49 | `usp_Tutorial_setCharInfo` | AccountServer | Operación de base de datos |
| 50 | `usp_User_GetDatediffCreated` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 51 | `usp_User_USloginGame` | AccountServer | Autenticación de credenciales de usuario y verificación de servidor |
| 52 | `usp_User_USloginGame_O` | AccountServer | Autenticación de credenciales de usuario y verificación de servidor |
| 53 | `usp_User_chkBlock` | AccountServer | Verificación de estado de baneo o bloqueo de cuenta |
| 54 | `usp_User_chkRegister` | AccountServer | Comprobación de registro y duplicados de ID de cuenta |
| 55 | `usp_User_chkUserID` | AccountServer | Gestión de información de usuario |
| 56 | `usp_User_delete` | AccountServer | Marcado o eliminación de cuenta de usuario |
| 57 | `usp_User_getAge` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 58 | `usp_User_getAuthList` | AccountServer | Autenticación de credenciales de usuario y verificación de servidor |
| 59 | `usp_User_getBuyCount` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 60 | `usp_User_getBuyList` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 61 | `usp_User_getLogin` | AccountServer | Autenticación de credenciales de usuario y verificación de servidor |
| 62 | `usp_User_getOption` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 63 | `usp_User_getPayCount` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 64 | `usp_User_getPayList` | AccountServer | Registro y procesamiento de órdenes de pago / facturación |
| 65 | `usp_User_getUser` | AccountServer | Obtención de perfil, datos de cuenta y permisos de usuario |
| 66 | `usp_User_insert` | AccountServer | Gestión de información de usuario |
| 67 | `usp_User_insertOption` | AccountServer | Gestión de información de usuario |
| 68 | `usp_User_insertOption_Old` | AccountServer | Gestión de información de usuario |
| 69 | `usp_User_insert_Old` | AccountServer | Gestión de información de usuario |
| 70 | `usp_User_loginGame` | AccountServer | Autenticación de credenciales de usuario y verificación de servidor |
| 71 | `usp_User_search` | AccountServer | Gestión de información de usuario |
| 72 | `usp_User_searchZipcode` | AccountServer | Gestión de información de usuario |
| 73 | `usp_User_setChannelLog` | AccountServer | Actualización de información o estado de la cuenta |
| 74 | `usp_User_setDrawComment` | AccountServer | Actualización de información o estado de la cuenta |
| 75 | `usp_User_setLoginable` | AccountServer | Autenticación de credenciales de usuario y verificación de servidor |
| 76 | `usp_User_update` | AccountServer | Actualización de información o estado de la cuenta |
| 77 | `usp_User_updatePassword` | AccountServer | Actualización de información o estado de la cuenta |

---

## World00_Character (`World00_Character.bak`) - Mundo, Personajes, Items y Sistemas Sociales (408 Procedures)

| # | Nombre del Stored Procedure | Servicio Invocador | Propósito / Operación en el Servidor |
| :---: | :--- | :--- | :--- |
| 1 | `Adam_ListProc` | CharServer / ZoneServer | Operación de base de datos |
| 2 | `CommandExecute` | CharServer / ZoneServer | Operación de base de datos |
| 3 | `IndexOptimize` | CharServer / ZoneServer | Operación de base de datos |
| 4 | `p_Char_AddExp` | CharServer / ZoneServer | Adición o pérdida de experiencia y cálculo de nivel |
| 5 | `p_Char_ChangeClass` | CharServer / ZoneServer | Promoción o cambio de clase (Job Advancement) |
| 6 | `p_Char_ChangeID` | CharServer / ZoneServer | Servicio de cambio de nombre de personaje |
| 7 | `p_Char_CheckOptionsExistance` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 8 | `p_Char_Create` | CharServer / ZoneServer | Creación de nuevo personaje con apariencia y atributos base |
| 9 | `p_Char_CreateSetDefaultData` | CharServer / ZoneServer | Creación de nuevo personaje con apariencia y atributos base |
| 10 | `p_Char_DelCharTitle` | CharServer / ZoneServer | Gestión de títulos obtenidos y equipados por el personaje |
| 11 | `p_Char_DelChargedBuff` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 12 | `p_Char_DelOptions` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 13 | `p_Char_DelState` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 14 | `p_Char_DelStateAll` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 15 | `p_Char_Delete` | CharServer / ZoneServer | Eliminación suave o permanente de personaje |
| 16 | `p_Char_FameAdd` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 17 | `p_Char_FameGet` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 18 | `p_Char_FameSet` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 19 | `p_Char_FameSub` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 20 | `p_Char_Find` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 21 | `p_Char_GetAllData` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 22 | `p_Char_GetCharID` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 23 | `p_Char_GetCharNo` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 24 | `p_Char_GetChargedBufferAll` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 25 | `p_Char_GetChargedBufferCount` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 26 | `p_Char_GetIsDeleted` | CharServer / ZoneServer | Eliminación suave o permanente de personaje |
| 27 | `p_Char_GetKQMap` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 28 | `p_Char_GetListOfUserChar` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 29 | `p_Char_GetListOfUserCharAll` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 30 | `p_Char_GetLoginData` | CharServer / ZoneServer | Carga inicial de personaje al iniciar sesión en el mundo |
| 31 | `p_Char_GetNumOfUserChar` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 32 | `p_Char_GetNumOfUserCharAll` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 33 | `p_Char_GetOptGame` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 34 | `p_Char_GetOptKeyMapping` | CharServer / ZoneServer | Guardado de la barra de habilidades y configuración de teclas |
| 35 | `p_Char_GetOptShortCutData` | CharServer / ZoneServer | Guardado de la barra de habilidades y configuración de teclas |
| 36 | `p_Char_GetOptShortCutSize` | CharServer / ZoneServer | Guardado de la barra de habilidades y configuración de teclas |
| 37 | `p_Char_GetOptSound` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 38 | `p_Char_GetOptVideo` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 39 | `p_Char_GetOptWindowsPos` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 40 | `p_Char_GetShape` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 41 | `p_Char_GetState` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 42 | `p_Char_GetStateAll` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 43 | `p_Char_GetStateAllWithRestTime` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 44 | `p_Char_GetStateAllWithSetItem` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 45 | `p_Char_GetStateAllWithSetItemAndRestTime` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 46 | `p_Char_GetStateWithRestTime` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 47 | `p_Char_GetStateWithSetItem` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 48 | `p_Char_GetStateWithSetItemAndRestTime` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 49 | `p_Char_GetTitleAll` | CharServer / ZoneServer | Gestión de títulos obtenidos y equipados por el personaje |
| 50 | `p_Char_GetTitleCount` | CharServer / ZoneServer | Gestión de títulos obtenidos y equipados por el personaje |
| 51 | `p_Char_GetUserNo` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 52 | `p_Char_Login` | CharServer / ZoneServer | Carga inicial de personaje al iniciar sesión en el mundo |
| 53 | `p_Char_Logout` | CharServer / ZoneServer | Guardado completo del estado del personaje al desconectarse |
| 54 | `p_Char_MoneyAdd` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 55 | `p_Char_MoneyGet` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 56 | `p_Char_MoneyGetCharsByUserNo` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 57 | `p_Char_MoneyGetGreaterThen` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 58 | `p_Char_MoneyMove` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 59 | `p_Char_MoneySet` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 60 | `p_Char_MoneySub` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 61 | `p_Char_Rename` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 62 | `p_Char_SaveLevelExpFrame` | CharServer / ZoneServer | Adición o pérdida de experiencia y cálculo de nivel |
| 63 | `p_Char_SaveLocation` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 64 | `p_Char_SavePKCount` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 65 | `p_Char_SaveStat` | CharServer / ZoneServer | Asignación o reseteo de puntos de estadísticas (STR, DEX, etc.) |
| 66 | `p_Char_SetAdminLevel` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 67 | `p_Char_SetAllData` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 68 | `p_Char_SetChargedBuffer` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 69 | `p_Char_SetLevelExp` | CharServer / ZoneServer | Adición o pérdida de experiencia y cálculo de nivel |
| 70 | `p_Char_SetOptGame` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 71 | `p_Char_SetOptKeyMapping` | CharServer / ZoneServer | Guardado de la barra de habilidades y configuración de teclas |
| 72 | `p_Char_SetOptShortCutData` | CharServer / ZoneServer | Guardado de la barra de habilidades y configuración de teclas |
| 73 | `p_Char_SetOptShortCutSize` | CharServer / ZoneServer | Guardado de la barra de habilidades y configuración de teclas |
| 74 | `p_Char_SetOptSound` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 75 | `p_Char_SetOptVideo` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 76 | `p_Char_SetOptWindowsPos` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 77 | `p_Char_SetShape` | CharServer / ZoneServer | Gestión de ciclo de vida y propiedades del personaje |
| 78 | `p_Char_SetState` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 79 | `p_Char_SetStateWithRestTime` | CharServer / ZoneServer | Guardado de buffs activos, pociones y efectos temporales |
| 80 | `p_Char_SetStateWithSetItem` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 81 | `p_Char_SetStateWithSetItemAndRestTime` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 82 | `p_Char_SetTitle` | CharServer / ZoneServer | Gestión de títulos obtenidos y equipados por el personaje |
| 83 | `p_Char_StatPoint_Dec` | CharServer / ZoneServer | Asignación o reseteo de puntos de estadísticas (STR, DEX, etc.) |
| 84 | `p_Char_StatPoint_Inc` | CharServer / ZoneServer | Asignación o reseteo de puntos de estadísticas (STR, DEX, etc.) |
| 85 | `p_Char_SubExp` | CharServer / ZoneServer | Adición o pérdida de experiencia y cálculo de nivel |
| 86 | `p_Char_WeddingStatus_Get` | WorldServer | Consulta o actualización del estado conyugal |
| 87 | `p_Char_WeddingStatus_Set` | WorldServer | Consulta o actualización del estado conyugal |
| 88 | `p_Char_Wedding_Cancel_Propose` | WorldServer | Envío o cancelación de propuesta de matrimonio |
| 89 | `p_Char_Wedding_Divorce_Cancel` | WorldServer | Proceso de divorcio (mutuo acuerdo o forzado) |
| 90 | `p_Char_Wedding_Divorce_Do` | WorldServer | Proceso de divorcio (mutuo acuerdo o forzado) |
| 91 | `p_Char_Wedding_Divorce_Do_by_Item` | WorldServer | Proceso de divorcio (mutuo acuerdo o forzado) |
| 92 | `p_Char_Wedding_Divorce_by_Agree` | WorldServer | Proceso de divorcio (mutuo acuerdo o forzado) |
| 93 | `p_Char_Wedding_Divorce_by_Force` | WorldServer | Proceso de divorcio (mutuo acuerdo o forzado) |
| 94 | `p_Char_Wedding_Do` | WorldServer | Consumación de matrimonio y registro de anillo |
| 95 | `p_Char_Wedding_Escape_Divorce_State` | WorldServer | Proceso de divorcio (mutuo acuerdo o forzado) |
| 96 | `p_Char_Wedding_Get_Partner_Info` | WorldServer | Gestión del sistema de bodas y matrimonios |
| 97 | `p_Char_Wedding_Propose` | WorldServer | Envío o cancelación de propuesta de matrimonio |
| 98 | `p_Char_Wipe` | CharServer / ZoneServer | Eliminación suave o permanente de personaje |
| 99 | `p_Event_Friend_Reward` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 100 | `p_Event_WorldCup_Init` | CharServer / ZoneServer | Operación de base de datos |
| 101 | `p_Friend_Del` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 102 | `p_Friend_Del_All` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 103 | `p_Friend_Get` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 104 | `p_Friend_GetAll` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 105 | `p_Friend_GetAll2` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 106 | `p_Friend_GetNumOfFriend` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 107 | `p_Friend_Set` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 108 | `p_GuildTournament_Get` | WorldServer | Inscripción y gestión de torneos de gremios |
| 109 | `p_GuildTournament_GetAllTournamentList` | WorldServer | Inscripción y gestión de torneos de gremios |
| 110 | `p_GuildTournament_GetLastMatch` | WorldServer | Inscripción y gestión de torneos de gremios |
| 111 | `p_GuildTournament_Set` | WorldServer | Inscripción y gestión de torneos de gremios |
| 112 | `p_Guild_Delete` | WorldServer | Disolución o eliminación de gremio |
| 113 | `p_Guild_Get` | WorldServer | Gestión de información y configuración de gremios |
| 114 | `p_Guild_GetAll` | WorldServer | Gestión de información y configuración de gremios |
| 115 | `p_Guild_GetCount` | WorldServer | Gestión de información y configuración de gremios |
| 116 | `p_Guild_GetIntro` | WorldServer | Gestión de información y configuración de gremios |
| 117 | `p_Guild_GetList` | WorldServer | Gestión de información y configuración de gremios |
| 118 | `p_Guild_GetListDetail` | WorldServer | Gestión de información y configuración de gremios |
| 119 | `p_Guild_GetMyNo` | WorldServer | Gestión de información y configuración de gremios |
| 120 | `p_Guild_GetName` | WorldServer | Gestión de información y configuración de gremios |
| 121 | `p_Guild_GetNotify` | WorldServer | Gestión de información y configuración de gremios |
| 122 | `p_Guild_Make` | WorldServer | Creación de un nuevo gremio |
| 123 | `p_Guild_Make2` | WorldServer | Creación de un nuevo gremio |
| 124 | `p_Guild_Member_Get` | WorldServer | Gestión de miembros, rangos, permisos e intros de gremio |
| 125 | `p_Guild_Member_GetAll` | WorldServer | Gestión de miembros, rangos, permisos e intros de gremio |
| 126 | `p_Guild_Member_GetCount` | WorldServer | Gestión de miembros, rangos, permisos e intros de gremio |
| 127 | `p_Guild_Member_Join` | WorldServer | Gestión de miembros, rangos, permisos e intros de gremio |
| 128 | `p_Guild_Member_Leave` | WorldServer | Gestión de miembros, rangos, permisos e intros de gremio |
| 129 | `p_Guild_Member_SetExpRatio` | WorldServer | Gestión de miembros, rangos, permisos e intros de gremio |
| 130 | `p_Guild_Member_SetFlags` | WorldServer | Gestión de miembros, rangos, permisos e intros de gremio |
| 131 | `p_Guild_Member_SetGrade` | WorldServer | Gestión de miembros, rangos, permisos e intros de gremio |
| 132 | `p_Guild_Member_SetIntro` | WorldServer | Gestión de miembros, rangos, permisos e intros de gremio |
| 133 | `p_Guild_MoneyAdd` | WorldServer | Depósito y extracción del fondo monetario del gremio |
| 134 | `p_Guild_MoneySub` | WorldServer | Depósito y extracción del fondo monetario del gremio |
| 135 | `p_Guild_Set` | WorldServer | Gestión de información y configuración de gremios |
| 136 | `p_Guild_SetDismiss` | WorldServer | Disolución o eliminación de gremio |
| 137 | `p_Guild_SetExp` | WorldServer | Gestión de información y configuración de gremios |
| 138 | `p_Guild_SetFame` | WorldServer | Gestión de información y configuración de gremios |
| 139 | `p_Guild_SetGrade` | WorldServer | Gestión de información y configuración de gremios |
| 140 | `p_Guild_SetIntro` | WorldServer | Gestión de información y configuración de gremios |
| 141 | `p_Guild_SetNotify` | WorldServer | Gestión de información y configuración de gremios |
| 142 | `p_Guild_SetStoneLevel` | WorldServer | Gestión de información y configuración de gremios |
| 143 | `p_Guild_SetType` | WorldServer | Gestión de información y configuración de gremios |
| 144 | `p_Guild_SetWar` | WorldServer | Declaración, registro y resultado de guerras de gremios |
| 145 | `p_Guild_SetWarResult` | WorldServer | Declaración, registro y resultado de guerras de gremios |
| 146 | `p_Guild_setName` | WorldServer | Gestión de información y configuración de gremios |
| 147 | `p_HolyPromise_CheckRing` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 148 | `p_HolyPromise_Del` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 149 | `p_HolyPromise_DelChar` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 150 | `p_HolyPromise_Get` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 151 | `p_HolyPromise_GetMember` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 152 | `p_HolyPromise_GetMemberCount` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 153 | `p_HolyPromise_Set` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 154 | `p_HolyPromise_SetDate` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 155 | `p_Item_BuyAll` | ZoneServer | Compra de ítems en tiendas NPC o lotes |
| 156 | `p_Item_BuyLot` | ZoneServer | Compra de ítems en tiendas NPC o lotes |
| 157 | `p_Item_ClrWeaponTitle` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 158 | `p_Item_CostumWeapon_Endure_Set` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 159 | `p_Item_CostumWeapon_Endure_Update` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 160 | `p_Item_Create` | ZoneServer | Creación o spawn de ítem en slot vacío de inventario |
| 161 | `p_Item_Create_UseEmptyStorage` | ZoneServer | Creación o spawn de ítem en slot vacío de inventario |
| 162 | `p_Item_Create_UseEmptyStorage_tmp` | ZoneServer | Creación o spawn de ítem en slot vacío de inventario |
| 163 | `p_Item_Delete` | ZoneServer | Destrucción o consumo de ítem del inventario |
| 164 | `p_Item_ExchangeStorage` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 165 | `p_Item_ExchangeStorageAndOwner` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 166 | `p_Item_Furniture_Endure_Update` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 167 | `p_Item_GetData` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 168 | `p_Item_GetList` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 169 | `p_Item_GetListType` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 170 | `p_Item_GetMobList` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 171 | `p_Item_GetMobListAll` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 172 | `p_Item_GetNumOfItem` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 173 | `p_Item_GetOption` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 174 | `p_Item_GetOptionAll` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 175 | `p_Item_GetOptionCount` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 176 | `p_Item_GetStorage` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 177 | `p_Item_GetStorageAndOwner` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 178 | `p_Item_GetStorageEmpty` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 179 | `p_Item_GetStorageEmptyList` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 180 | `p_Item_GetStorageListType` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 181 | `p_Item_GetStorageTypeCount` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 182 | `p_Item_GetTitle` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 183 | `p_Item_IsValidOwner` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 184 | `p_Item_Merge` | ZoneServer | División y apilamiento de consumibles/materiales |
| 185 | `p_Item_Merge_New` | ZoneServer | División y apilamiento de consumibles/materiales |
| 186 | `p_Item_SellAll` | ZoneServer | Venta de ítems a comerciantes NPC |
| 187 | `p_Item_SellAll_LotCheck` | ZoneServer | Venta de ítems a comerciantes NPC |
| 188 | `p_Item_SellLot` | ZoneServer | Venta de ítems a comerciantes NPC |
| 189 | `p_Item_SellLot_LotCheck` | ZoneServer | Venta de ítems a comerciantes NPC |
| 190 | `p_Item_SetData` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 191 | `p_Item_SetMobList` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 192 | `p_Item_SetOption` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 193 | `p_Item_SetStorage` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 194 | `p_Item_SetStorageAndOwner` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 195 | `p_Item_SetStorageAndOwner_UseEmptyStorage` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 196 | `p_Item_SetTitle` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 197 | `p_Item_SetWeaponTitle` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 198 | `p_Item_Split_N_Merge` | ZoneServer | División y apilamiento de consumibles/materiales |
| 199 | `p_Item_Split_N_Merge_New` | ZoneServer | División y apilamiento de consumibles/materiales |
| 200 | `p_Item_Upgrade` | ZoneServer | Refinamiento de equipamiento con piedras celestiales |
| 201 | `p_Item_Upgrade_CelllInAll` | ZoneServer | Refinamiento de equipamiento con piedras celestiales |
| 202 | `p_Item_Upgrade_GOLD_NINE` | ZoneServer | Refinamiento de equipamiento con piedras celestiales |
| 203 | `p_Item_Upgrade_OLD` | ZoneServer | Refinamiento de equipamiento con piedras celestiales |
| 204 | `p_KingdomQuest_Del` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 205 | `p_KingdomQuest_Get` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 206 | `p_KingdomQuest_GetAll` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 207 | `p_KingdomQuest_GetNumOfKQ` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 208 | `p_KingdomQuest_Set` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 209 | `p_MasPup_DelAll` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 210 | `p_MasPup_DelMaster` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 211 | `p_MasPup_DelPupils` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 212 | `p_MasPup_GetAll` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 213 | `p_MasPup_SetMaster` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 214 | `p_MasPup_SetPupils` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 215 | `p_Prison_Add` | WorldServer | Envío a prisión por sanciones GM y cuenta regresiva de penalización |
| 216 | `p_Prison_Add2` | WorldServer | Envío a prisión por sanciones GM y cuenta regresiva de penalización |
| 217 | `p_Prison_Alter` | WorldServer | Envío a prisión por sanciones GM y cuenta regresiva de penalización |
| 218 | `p_Prison_Alter2` | WorldServer | Envío a prisión por sanciones GM y cuenta regresiva de penalización |
| 219 | `p_Prison_End` | WorldServer | Envío a prisión por sanciones GM y cuenta regresiva de penalización |
| 220 | `p_Prison_Get` | WorldServer | Envío a prisión por sanciones GM y cuenta regresiva de penalización |
| 221 | `p_Prison_UpdateCharPrisonMin` | WorldServer | Envío a prisión por sanciones GM y cuenta regresiva de penalización |
| 222 | `p_Quest_Add` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 223 | `p_Quest_Del` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 224 | `p_Quest_Get` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 225 | `p_Quest_GetAll` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 226 | `p_Quest_GetAllDoing` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 227 | `p_Quest_GetAllDone` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 228 | `p_Quest_Set` | ZoneServer | Aceptación, progreso, entrega y reseteo de misiones |
| 229 | `p_Skill_Del` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 230 | `p_Skill_Get` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 231 | `p_Skill_GetAll` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 232 | `p_Skill_GetNumOfSkill` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 233 | `p_Skill_ResetPower` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 234 | `p_Skill_Set` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 235 | `p_Skill_SetCoolTime` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 236 | `p_Skill_SetPassive` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 237 | `p_Skill_SetPower` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 238 | `p_Skill_SetPowerAll` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 239 | `p_User_Deposit` | CharServer / ZoneServer | Operación de base de datos |
| 240 | `p_User_MoneyAdd` | CharServer / ZoneServer | Operación de base de datos |
| 241 | `p_User_MoneyGet` | CharServer / ZoneServer | Operación de base de datos |
| 242 | `p_User_MoneyGetGreaterThen` | CharServer / ZoneServer | Operación de base de datos |
| 243 | `p_User_MoneySet` | CharServer / ZoneServer | Operación de base de datos |
| 244 | `p_User_MoneySub` | CharServer / ZoneServer | Operación de base de datos |
| 245 | `p_User_Withdraw` | CharServer / ZoneServer | Operación de base de datos |
| 246 | `p_View_GetData` | CharServer / ZoneServer | Operación de base de datos |
| 247 | `usp_Attendance_getCallDateList` | WorldServer | Calendario de asistencia diaria y entrega de premios acumulados |
| 248 | `usp_Attendance_getCount` | WorldServer | Calendario de asistencia diaria y entrega de premios acumulados |
| 249 | `usp_Attendance_getRewardList` | WorldServer | Calendario de asistencia diaria y entrega de premios acumulados |
| 250 | `usp_Attendance_setCall` | WorldServer | Calendario de asistencia diaria y entrega de premios acumulados |
| 251 | `usp_Attendance_setReward` | WorldServer | Calendario de asistencia diaria y entrega de premios acumulados |
| 252 | `usp_Char_GetAutoPick` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 253 | `usp_Char_GetRestExpExecTime` | CharServer / ZoneServer | Adición o pérdida de experiencia y cálculo de nivel |
| 254 | `usp_Char_SetAutoPick` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 255 | `usp_Char_SetRestExpExecTime` | CharServer / ZoneServer | Adición o pérdida de experiencia y cálculo de nivel |
| 256 | `usp_Char_chkCharNo` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 257 | `usp_CharacterLoading` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 258 | `usp_Character_GetDiffFriendPoint` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 259 | `usp_Character_GetFriendpoint` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 260 | `usp_Character_SetFriendpoint` | WorldServer | Lista de amigos, puntos de amistad y sistema de Maestro/Pupilo |
| 261 | `usp_Character_chkKeyMap` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 262 | `usp_Character_getAutoScroll` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 263 | `usp_Character_getChangeIDList` | CharServer / ZoneServer | Servicio de cambio de nombre de personaje |
| 264 | `usp_Character_getGameOption` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 265 | `usp_Character_getItemActionCoolTime` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 266 | `usp_Character_getKeyMap` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 267 | `usp_Character_getShortCutList` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 268 | `usp_Character_getSingleOption` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 269 | `usp_Character_initKeyMap` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 270 | `usp_Character_setAutoPick` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 271 | `usp_Character_setBuffScroll` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 272 | `usp_Character_setChangeID` | CharServer / ZoneServer | Servicio de cambio de nombre de personaje |
| 273 | `usp_Character_setChangeIDItemUsed` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 274 | `usp_Character_setChargedbuffScroll` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 275 | `usp_Character_setChatBan` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 276 | `usp_Character_setGameOption` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 277 | `usp_Character_setGameOptionInit` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 278 | `usp_Character_setItemActionCoolTime` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 279 | `usp_Character_setKeyMap` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 280 | `usp_Character_setMysteryBolt` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 281 | `usp_Character_setRestExpExecTime` | CharServer / ZoneServer | Adición o pérdida de experiencia y cálculo de nivel |
| 282 | `usp_Character_setShortCut` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 283 | `usp_Character_setStatPoint` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 284 | `usp_ChargeItem_Draw` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 285 | `usp_ChargeItem_GetList` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 286 | `usp_ChargeItem_UseSkillPowerInitDel` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 287 | `usp_ChargeItem_UseSkillPowerInitLotChange` | ZoneServer | Árbol de habilidades, niveles de skill, cooldowns y potenciadores |
| 288 | `usp_ChargeItem_UseStatPointDecDel` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 289 | `usp_ChargeItem_UseStatPointDecLotChange` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 290 | `usp_ChargeItem_UseStatPointInitDel` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 291 | `usp_ChargeItem_UseStatPointInitLotChange` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 292 | `usp_ChatBlock_delBlock` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 293 | `usp_ChatBlock_delBlockALL` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 294 | `usp_ChatBlock_getBlockList` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 295 | `usp_ChatBlock_setBlock` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 296 | `usp_ChatColor_get` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 297 | `usp_ChatColor_setBalloon` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 298 | `usp_ChatColor_setFont` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 299 | `usp_Collection_getCardList` | CharServer / ZoneServer | Álbum de colección de cartas de monstruos y recompensas de estrellas |
| 300 | `usp_Collection_getFavoritesList` | CharServer / ZoneServer | Álbum de colección de cartas de monstruos y recompensas de estrellas |
| 301 | `usp_Collection_getRewardList` | CharServer / ZoneServer | Álbum de colección de cartas de monstruos y recompensas de estrellas |
| 302 | `usp_Collection_getSerialList` | CharServer / ZoneServer | Álbum de colección de cartas de monstruos y recompensas de estrellas |
| 303 | `usp_Collection_setCard` | CharServer / ZoneServer | Álbum de colección de cartas de monstruos y recompensas de estrellas |
| 304 | `usp_Collection_setFavorites` | CharServer / ZoneServer | Álbum de colección de cartas de monstruos y recompensas de estrellas |
| 305 | `usp_Collection_setReward` | CharServer / ZoneServer | Álbum de colección de cartas de monstruos y recompensas de estrellas |
| 306 | `usp_Collection_setSerial` | CharServer / ZoneServer | Álbum de colección de cartas de monstruos y recompensas de estrellas |
| 307 | `usp_Dice_chgMachineToMoneycard` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 308 | `usp_Dice_getDiceExpect` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 309 | `usp_Dice_getDiceRank` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 310 | `usp_Dice_setDice` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 311 | `usp_Dice_setDiceLog` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 312 | `usp_Dice_setDicePay` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 313 | `usp_Dice_setDiceResult` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 314 | `usp_GM_getItemRewardGroupNo` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 315 | `usp_GM_getItemRewardList` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 316 | `usp_GM_setItemReward` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 317 | `usp_Gamble_addCoin` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 318 | `usp_Gamble_getCoin` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 319 | `usp_Gamble_getDiceRank` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 320 | `usp_Gamble_getExchangedCoin` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 321 | `usp_Gamble_getJackpotCoin` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 322 | `usp_Gamble_getJackpotRanking` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 323 | `usp_Gamble_setCoinToMoney` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 324 | `usp_Gamble_setDice` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 325 | `usp_Gamble_setDiceLog` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 326 | `usp_Gamble_setDiceResult` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 327 | `usp_Gamble_setMoneyToCoin` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 328 | `usp_Gamble_setSlotMachineResult` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 329 | `usp_Gamble_subCoin` | ZoneServer | Minijuegos de azar: tragamonedas (slot machine), dados y jackpot de fichas |
| 330 | `usp_GuildAcademy_delMember` | WorldServer | Gestión de aprendices en la academia del gremio |
| 331 | `usp_GuildAcademy_delReward` | WorldServer | Reclamación y distribución de recompensas de academia |
| 332 | `usp_GuildAcademy_delRewardAll` | WorldServer | Reclamación y distribución de recompensas de academia |
| 333 | `usp_GuildAcademy_depositMoney` | WorldServer | Administración de la academia del gremio |
| 334 | `usp_GuildAcademy_getGraduateCount` | WorldServer | Graduación de aprendices al alcanzar el nivel tope |
| 335 | `usp_GuildAcademy_getGraduateList` | WorldServer | Graduación de aprendices al alcanzar el nivel tope |
| 336 | `usp_GuildAcademy_getGuildAcademyList` | WorldServer | Administración de la academia del gremio |
| 337 | `usp_GuildAcademy_getGuildNo` | WorldServer | Administración de la academia del gremio |
| 338 | `usp_GuildAcademy_getHistoryList` | WorldServer | Administración de la academia del gremio |
| 339 | `usp_GuildAcademy_getMemberList` | WorldServer | Gestión de aprendices en la academia del gremio |
| 340 | `usp_GuildAcademy_getMoney` | WorldServer | Administración de la academia del gremio |
| 341 | `usp_GuildAcademy_getRanking` | WorldServer | Cálculo del ranking de academias |
| 342 | `usp_GuildAcademy_getReward` | WorldServer | Reclamación y distribución de recompensas de academia |
| 343 | `usp_GuildAcademy_getRewardAll` | WorldServer | Reclamación y distribución de recompensas de academia |
| 344 | `usp_GuildAcademy_raiseMember` | WorldServer | Gestión de aprendices en la academia del gremio |
| 345 | `usp_GuildAcademy_setAcademy` | WorldServer | Administración de la academia del gremio |
| 346 | `usp_GuildAcademy_setChatBanTime` | WorldServer | Administración de la academia del gremio |
| 347 | `usp_GuildAcademy_setGraduate` | WorldServer | Graduación de aprendices al alcanzar el nivel tope |
| 348 | `usp_GuildAcademy_setHistory` | WorldServer | Administración de la academia del gremio |
| 349 | `usp_GuildAcademy_setMaster` | WorldServer | Administración de la academia del gremio |
| 350 | `usp_GuildAcademy_setMember` | WorldServer | Gestión de aprendices en la academia del gremio |
| 351 | `usp_GuildAcademy_setMemberLevelup` | WorldServer | Gestión de aprendices en la academia del gremio |
| 352 | `usp_GuildAcademy_setNotify` | WorldServer | Administración de la academia del gremio |
| 353 | `usp_GuildAcademy_setRanking` | WorldServer | Cálculo del ranking de academias |
| 354 | `usp_GuildAcademy_setReward` | WorldServer | Reclamación y distribución de recompensas de academia |
| 355 | `usp_GuildAcademy_withdrawMoney` | WorldServer | Administración de la academia del gremio |
| 356 | `usp_GuildEmblem_GetInfo` | WorldServer | Carga y guardado del emblema/escudo del gremio |
| 357 | `usp_GuildEmblem_GetInfoAll` | WorldServer | Carga y guardado del emblema/escudo del gremio |
| 358 | `usp_GuildEmblem_SetInfo` | WorldServer | Carga y guardado del emblema/escudo del gremio |
| 359 | `usp_GuildTournament_getChoiceList` | WorldServer | Inscripción y gestión de torneos de gremios |
| 360 | `usp_GuildTournament_getMoney` | WorldServer | Inscripción y gestión de torneos de gremios |
| 361 | `usp_GuildTournament_getRegistList` | WorldServer | Inscripción y gestión de torneos de gremios |
| 362 | `usp_GuildTournament_setChoice` | WorldServer | Inscripción y gestión de torneos de gremios |
| 363 | `usp_GuildTournament_setRegist` | WorldServer | Inscripción y gestión de torneos de gremios |
| 364 | `usp_GuildTournament_setResult` | WorldServer | Inscripción y gestión de torneos de gremios |
| 365 | `usp_Guild_GetToken` | WorldServer | Gestión de fichas y monedas especiales de gremio |
| 366 | `usp_Guild_GetTokenAll` | WorldServer | Gestión de fichas y monedas especiales de gremio |
| 367 | `usp_Guild_SetToken` | WorldServer | Gestión de fichas y monedas especiales de gremio |
| 368 | `usp_Guild_TokenSub` | WorldServer | Gestión de fichas y monedas especiales de gremio |
| 369 | `usp_HolyPromise_getRemainMoney` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 370 | `usp_HolyPromise_setMoney` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 371 | `usp_HolyPromise_withDrawRemainMoney` | WorldServer | Sistema de compromiso Holy Promise (anillos de promesa y bendiciones) |
| 372 | `usp_InGameGuide_getViewList` | CharServer / ZoneServer | Operación de base de datos |
| 373 | `usp_InGameGuide_setGuideNo` | CharServer / ZoneServer | Operación de base de datos |
| 374 | `usp_Item_Add_EnchantGem` | ZoneServer | Engarce de gemas, sockets y mejoras de armas/armaduras |
| 375 | `usp_Item_Add_EnchantSocket` | ZoneServer | Engarce de gemas, sockets y mejoras de armas/armaduras |
| 376 | `usp_Item_Del_EnchantGem` | ZoneServer | Engarce de gemas, sockets y mejoras de armas/armaduras |
| 377 | `usp_Item_Get_EnchantAll` | ZoneServer | Engarce de gemas, sockets y mejoras de armas/armaduras |
| 378 | `usp_Item_GuildToken_BuyAll` | WorldServer | Gestión de fichas y monedas especiales de gremio |
| 379 | `usp_Item_GuildToken_BuyLot` | WorldServer | Gestión de fichas y monedas especiales de gremio |
| 380 | `usp_Item_Set_EnchantGemLot` | ZoneServer | Engarce de gemas, sockets y mejoras de armas/armaduras |
| 381 | `usp_Item_getAutoMonList` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 382 | `usp_Item_getEmptyStorage` | ZoneServer | Movimiento entre inventario, almacén de cuenta y equipo |
| 383 | `usp_Item_setAutoMon` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 384 | `usp_Item_setPutOnBelonged` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 385 | `usp_Item_setPutOnBelongedCancel` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 386 | `usp_Item_setPutOnBelongedExchange` | ZoneServer | Operaciones de inventario y manipulación de equipamiento |
| 387 | `usp_Minihouse_delGate` | ZoneServer | Sistema de casa personal Minihouse (portales, visitas y muebles) |
| 388 | `usp_Minihouse_getGateList` | ZoneServer | Sistema de casa personal Minihouse (portales, visitas y muebles) |
| 389 | `usp_Minihouse_getVisitorCount` | ZoneServer | Sistema de casa personal Minihouse (portales, visitas y muebles) |
| 390 | `usp_Minihouse_setGate` | ZoneServer | Sistema de casa personal Minihouse (portales, visitas y muebles) |
| 391 | `usp_Minihouse_setVisitorCount` | ZoneServer | Sistema de casa personal Minihouse (portales, visitas y muebles) |
| 392 | `usp_Pet_getInfo` | ZoneServer | Creación, renombramiento, invocación y habilidades de mascotas |
| 393 | `usp_Pet_setCreate` | ZoneServer | Creación, renombramiento, invocación y habilidades de mascotas |
| 394 | `usp_Pet_setName` | ZoneServer | Creación, renombramiento, invocación y habilidades de mascotas |
| 395 | `usp_SafeTrade_setAll` | ZoneServer | Comercio seguro entre dos personajes sin duplicaciones |
| 396 | `usp_SafeTrade_setLot` | ZoneServer | Comercio seguro entre dos personajes sin duplicaciones |
| 397 | `usp_SafeTrade_setMerge` | ZoneServer | División y apilamiento de consumibles/materiales |
| 398 | `usp_SpamBlock_getTime` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 399 | `usp_SpamBlock_setTime` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 400 | `usp_SpamFilter_get` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 401 | `usp_SpamFilter_set` | WorldServer | Filtro anti-spam, silenciamiento temporal (mute) y bloqueo de chat |
| 402 | `usp_Statistics_getCharacterLevel` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 403 | `usp_Statistics_setCharacterLevel` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 404 | `usp_Statistics_setGameMoney` | CharServer / ZoneServer | Operación de base de datos |
| 405 | `usp_Statistics_setTotalMoney` | CharServer / ZoneServer | Operación de base de datos |
| 406 | `usp_Tutorial_getCharInfo` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |
| 407 | `usp_Tutorial_getIsSkip` | CharServer / ZoneServer | Operación de base de datos |
| 408 | `usp_Tutorial_setCharInfo` | CharServer / ZoneServer | Persistencia de puntos de vida y maná |

---

## AccountLog (`AccountLog.bak`) - Auditoría de Cuentas y Seguridad (27 Procedures)

| # | Nombre del Stored Procedure | Servicio Invocador | Propósito / Operación en el Servidor |
| :---: | :--- | :--- | :--- |
| 1 | `p_AccountLog_Add` | AccountLogServer | Inserción y consulta del historial de logins, IPs y desconexiones |
| 2 | `p_AccountLog_Get` | AccountLogServer | Inserción y consulta del historial de logins, IPs y desconexiones |
| 3 | `p_AccountLog_GetCount` | AccountLogServer | Inserción y consulta del historial de logins, IPs y desconexiones |
| 4 | `p_AccountLog_GetLog` | AccountLogServer | Inserción y consulta del historial de logins, IPs y desconexiones |
| 5 | `p_AccountLog_GetLogCount` | AccountLogServer | Inserción y consulta del historial de logins, IPs y desconexiones |
| 6 | `p_Event_Friend_DelRequest` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 7 | `p_Event_Friend_EventComplete` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 8 | `p_Event_Friend_Get` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 9 | `p_Event_Friend_GetCount` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 10 | `p_Event_Friend_Giveup` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 11 | `p_Event_Friend_SetHideMsg` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 12 | `p_Event_WorldCup_Enter` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 13 | `p_IPBlock_Add` | AccountLogServer | Gestión de lista negra de IPs bloqueadas por el firewall / servidor |
| 14 | `p_IPBlock_Delete` | AccountLogServer | Gestión de lista negra de IPs bloqueadas por el firewall / servidor |
| 15 | `p_IPBlock_Get` | AccountLogServer | Gestión de lista negra de IPs bloqueadas por el firewall / servidor |
| 16 | `p_IPBlock_GetList` | AccountLogServer | Gestión de lista negra de IPs bloqueadas por el firewall / servidor |
| 17 | `p_IPBlock_IsBlock` | AccountLogServer | Gestión de lista negra de IPs bloqueadas por el firewall / servidor |
| 18 | `p_Report_Add` | AccountLogServer | Registro y seguimiento de denuncias de jugadores a operadores GM |
| 19 | `p_Report_Deal_OPT` | AccountLogServer | Registro y seguimiento de denuncias de jugadores a operadores GM |
| 20 | `p_Report_Del` | AccountLogServer | Registro y seguimiento de denuncias de jugadores a operadores GM |
| 21 | `p_Report_Get` | AccountLogServer | Registro y seguimiento de denuncias de jugadores a operadores GM |
| 22 | `p_Report_Get_Count_OPT` | AccountLogServer | Registro y seguimiento de denuncias de jugadores a operadores GM |
| 23 | `p_Report_Get_OPT` | AccountLogServer | Registro y seguimiento de denuncias de jugadores a operadores GM |
| 24 | `usp_EventiPad_getList` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 25 | `usp_EventiPad_getTicket` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 26 | `usp_EventiPad_setTicket` | AccountLogServer | Registro de participación y recompensas en eventos de cuenta |
| 27 | `usp_Log_setLogout` | AccountLogServer | Auditoría de estadísticas y tiempo de juego al desconectarse |

---

## World00_GameLog (`World00_GameLog.bak`) - Logs de Acciones de Juego (26 Procedures)

| # | Nombre del Stored Procedure | Servicio Invocador | Propósito / Operación en el Servidor |
| :---: | :--- | :--- | :--- |
| 1 | `p_GameLog_Add` | GameLogServer | Registro general de eventos del mundo en el log |
| 2 | `p_GameLog_Add_Type1` | GameLogServer | Log de conexión y entrada al mundo |
| 3 | `p_GameLog_Add_Type2` | GameLogServer | Log de desconexión y salida del mundo |
| 4 | `p_GameLog_Add_Type3` | GameLogServer | Log de creación de personaje |
| 5 | `p_GameLog_Add_Type4` | GameLogServer | Log de eliminación de personaje |
| 6 | `p_GameLog_Add_Type5` | GameLogServer | Log de cambio de mapa / teleport |
| 7 | `p_GameLog_Add_Type6` | GameLogServer | Log de subida de nivel (Level Up) |
| 8 | `p_GameLog_Add_Type7` | GameLogServer | Log de muerte y resurrección de personaje |
| 9 | `p_GameLog_Add_Type8` | GameLogServer | Log de compra / venta de ítems en NPC |
| 10 | `p_GameLog_Add_Type9` | GameLogServer | Log de intercambio de ítems entre jugadores (Trade) |
| 11 | `p_GameLog_Add_TypeA` | GameLogServer | Log de drop de ítems al suelo o destrucción |
| 12 | `p_GameLog_Add_TypeB` | GameLogServer | Log de pickup / recolección de botín |
| 13 | `p_GameLog_Add_TypeC` | GameLogServer | Log de mejora y encantamiento de equipamiento |
| 14 | `p_GameLog_Add_TypeD` | GameLogServer | Log de transacciones de dinero / banco |
| 15 | `p_GameLog_Add_TypeE` | GameLogServer | Log de comandos administrativos ejecutados por GameMasters |
| 16 | `p_GameLog_Add_TypeF` | GameLogServer | Log de misiones aceptadas, completadas y canceladas |
| 17 | `p_GameLog_Add_TypeG` | GameLogServer | Log de gremios (creación, uniones, abandonos, expulsiones) |
| 18 | `p_GameLog_Add_TypeH` | GameLogServer | Log de PvP y guerras de gremios |
| 19 | `p_GameLog_Add_TypeI` | GameLogServer | Log de producción y crafteo de ítems |
| 20 | `p_GameLog_Add_TypeJ` | GameLogServer | Log de matrimonio y divorcio |
| 21 | `p_GameLog_Add_TypeK` | GameLogServer | Log de academia de gremio |
| 22 | `p_GameLog_Add_TypeL` | GameLogServer | Log de minijuegos y casino |
| 23 | `p_GameLog_Get` | GameLogServer | Consulta y filtrado de logs de auditoría para soporte técnico |
| 24 | `p_GameLog_GetCount` | GameLogServer | Consulta y filtrado de logs de auditoría para soporte técnico |
| 25 | `p_GameLog_GetLog` | GameLogServer | Consulta y filtrado de logs de auditoría para soporte técnico |
| 26 | `p_GameLog_GetLogCount` | GameLogServer | Consulta y filtrado de logs de auditoría para soporte técnico |

---

## StatisticsData (`StatisticsData.bak`) - Métricas y Estadísticas del Servidor (18 Procedures)

| # | Nombre del Stored Procedure | Servicio Invocador | Propósito / Operación en el Servidor |
| :---: | :--- | :--- | :--- |
| 1 | `pFindCCU` | WorldServer / Cron | Registro y cálculo de usuarios conectados en simultáneo (CCU) |
| 2 | `p_DayClassLevelPlay_Add` | WorldServer / Cron | Distribución de jugadores por clase, nivel y tiempo jugado |
| 3 | `p_DayClassLevelPlay_Query` | WorldServer / Cron | Distribución de jugadores por clase, nivel y tiempo jugado |
| 4 | `p_DayClassLevel_Add` | WorldServer / Cron | Distribución de jugadores por clase, nivel y tiempo jugado |
| 5 | `p_DayClassLevel_Query` | WorldServer / Cron | Distribución de jugadores por clase, nivel y tiempo jugado |
| 6 | `p_DayConnectionTime` | WorldServer / Cron | Cálculo periódico de estadísticas del servidor |
| 7 | `p_DayUser` | WorldServer / Cron | Estadísticas de usuarios diarios únicos (DAU) |
| 8 | `p_DayUser_Add` | WorldServer / Cron | Estadísticas de usuarios diarios únicos (DAU) |
| 9 | `p_DayUser_Query` | WorldServer / Cron | Estadísticas de usuarios diarios únicos (DAU) |
| 10 | `p_SameTimeChar_Add` | WorldServer / Cron | Registro y cálculo de usuarios conectados en simultáneo (CCU) |
| 11 | `p_SameTimeChar_Query` | WorldServer / Cron | Registro y cálculo de usuarios conectados en simultáneo (CCU) |
| 12 | `p_SameTimeWorld` | WorldServer / Cron | Registro y cálculo de usuarios conectados en simultáneo (CCU) |
| 13 | `p_UserAgeReport` | WorldServer / Cron | Métricas demográficas de jugadores (edad, género, región) |
| 14 | `p_UserAge_Query` | WorldServer / Cron | Métricas demográficas de jugadores (edad, género, región) |
| 15 | `p_UserGenderTest` | WorldServer / Cron | Métricas demográficas de jugadores (edad, género, región) |
| 16 | `p_UserGender_Query` | WorldServer / Cron | Métricas demográficas de jugadores (edad, género, región) |
| 17 | `p_UserLocation_Query` | WorldServer / Cron | Métricas demográficas de jugadores (edad, género, región) |
| 18 | `p_getQuestTimes` | WorldServer / Cron | Métricas de misiones completadas en el servidor |

---

## OperatorTool (`OperatorTool.bak`) - Herramientas de GameMasters y Operadores (17 Procedures)

| # | Nombre del Stored Procedure | Servicio Invocador | Propósito / Operación en el Servidor |
| :---: | :--- | :--- | :--- |
| 1 | `p_OPLog_Add` | OperatorTool (GM) | Registro de auditoría de cada acción ejecutada por GameMasters |
| 2 | `p_OPLog_Add_All` | OperatorTool (GM) | Registro de auditoría de cada acción ejecutada por GameMasters |
| 3 | `p_OPLog_Add_ToIP` | OperatorTool (GM) | Registro de auditoría de cada acción ejecutada por GameMasters |
| 4 | `p_OPLog_Add_ToUserNo` | OperatorTool (GM) | Registro de auditoría de cada acción ejecutada por GameMasters |
| 5 | `p_OPLog_Add_ToWorldNo` | OperatorTool (GM) | Registro de auditoría de cada acción ejecutada por GameMasters |
| 6 | `p_OPLog_Get` | OperatorTool (GM) | Registro de auditoría de cada acción ejecutada por GameMasters |
| 7 | `p_OPLog_GetCount` | OperatorTool (GM) | Registro de auditoría de cada acción ejecutada por GameMasters |
| 8 | `p_OPTMenu_GetAuth` | OperatorTool (GM) | Control de accesos y permisos por pestañas de la herramienta GM |
| 9 | `p_OPTMenu_GetAuthAll` | OperatorTool (GM) | Control de accesos y permisos por pestañas de la herramienta GM |
| 10 | `p_OPTMenu_GetCount` | OperatorTool (GM) | Control de accesos y permisos por pestañas de la herramienta GM |
| 11 | `p_OPTMenu_SetAuth` | OperatorTool (GM) | Control de accesos y permisos por pestañas de la herramienta GM |
| 12 | `p_Operator_Add` | OperatorTool (GM) | Alta, baja y modificación de privilegios de operadores |
| 13 | `p_Operator_Del` | OperatorTool (GM) | Alta, baja y modificación de privilegios de operadores |
| 14 | `p_Operator_GetAll` | OperatorTool (GM) | Consulta de operadores del sistema |
| 15 | `p_Operator_Logon` | OperatorTool (GM) | Autenticación de cuenta de GameMaster en la herramienta GM |
| 16 | `p_Operator_SetLevel` | OperatorTool (GM) | Alta, baja y modificación de privilegios de operadores |
| 17 | `p_Operator_SetPW` | OperatorTool (GM) | Alta, baja y modificación de privilegios de operadores |

---
