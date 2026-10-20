# Comparativa Real del Chunk de Entrada a Zona (`Zone.pdb` vs `XelionCore` vs `NextGen`)

Este documento detalla la comparativa técnica, paquete por paquete y byte por byte, de lo que ocurre desde que el jugador hace clic en su personaje (`PROTO_NC_CHAR_LOGIN_REQ` / `0x1001`) hasta que el cliente debe validar el chunk recibido en Zona y responder con `PROTO_NC_MAP_LOGINCOMPLETE_CMD` (`0x1803`), contrastando:

1. **Binarios Oficiales (`WorldManager.exe` / `WorldManager.pdb` + `Zone.exe` / `Zone.pdb`)**
2. **`XelionCore` (Go)** ([`pkg/server/world/handler.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/world/handler.go) y [`pkg/server/zone/handler.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/zone/handler.go))
3. **`NextGen-Emulator` (C#)** ([`NextGen.World/Handlers/Handler4.cs`](file:///C:/Users/fperez/Desktop/Xelion/NextGen-Emulator/NextGen.World/Handlers/Handler4.cs) y [`NextGen.Zone/Game/ZoneCharacter.cs`](file:///C:/Users/fperez/Desktop/Xelion/NextGen-Emulator/NextGen.Zone/Game/ZoneCharacter.cs))

---

## 1. Resumen Ejecutivo: ¿Por qué el cliente NO valida el chunk en `XelionCore` y nunca envía `0x1803`?

Al descompilar las funciones exactas de `Zone.exe` y `WorldManager.exe` y compararlas contra `XelionCore`, aparecen **3 errores críticos de protocolo en `XelionCore`** que hacen que el cliente se quede esperando o aborte el parseo del chunk:

### Causa Crítica #1: `PartMark = 1` en lugar de `PartMark = 3` en `PROTO_NC_CHAR_CLIENT_SKILL_CMD` (`0x103D`)
- **En `Zone.exe` (`so_charskilldata` @ `0x0044BF20`):**
  El byte `PartMark` (offset `+1` del payload de `0x103D`) es un bitfield:
  ```cpp
  struct PartMark {
      unsigned char first : 1; // bit 0 (0x01) = Primer fragmento
      unsigned char last  : 1; // bit 1 (0x02) = Último fragmento
      unsigned char seq   : 6; // bits 2..7    = Secuencia
  };
  ```
  En `0x0044C039`, `Zone.exe` hace `mov byte ptr [esi + 3], 1` (`first = 1`), y luego en `0x0044C06D`, cuando termina de meter las skills en el paquete (`eax >= nMaxNum`), ejecuta:
  ```asm
  0x0044C06D: or byte ptr [esi + 3], 2   ; PartMark |= 0x02 (last = 1) -> PartMark = 0x03
  ```
  Por lo tanto, cuando todas las skills entran en un único paquete, `Zone.exe` envía **`PartMark = 3` (`0x03`, `first=1 | last=1`)**.
- **En `XelionCore` ([`pkg/server/zone/handler.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/zone/handler.go#L292-L311)):**
  En la línea 295 escribe fijo:
  ```go
  skillPkt.WriteByte(1) // PartMark (1 byte) -> first=1, last=0 !!
  ```
  Como todo personaje nuevo tiene 25 `starter_skills` cargadas desde [`default_character_data.json`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/data/default_character_data.json#L55-L81), `XelionCore` envía `0x103D` con **`PartMark = 1` (`last = 0`)**. El cliente interpreta que **la lista de skills está incompleta** y se queda esperando el paquete con el bit `last = 1` (`0x02`).

### Causa Crítica #2: `XelionCore` omite por completo `PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD` (`0x104A`)
- **En `Zone.exe` (`so_CharChargedBuffData` @ `0x0044E1E0`, llamado desde `wms_NC_CHAR_ZONE_CHARGEDBUFF_CMD` @ `0x0044BCF0`):**
  Incluso cuando el personaje tiene `0` buffs cargados (`NumOfChargedBuff == 0`), el salto en `0x0044E20B` (`jae 0x44e461`) va a `0x0044E461`, arma el paquete `0x104A` con `NumOfChargedBuff = 0` (`uint16`, 2 bytes de payload / 4 bytes en cable) y **SIEMPRE lo envía al cliente**:
  ```asm
  0x0044E46F: mov edx, 0x104a
  0x0044E474: mov word ptr [ebx], dx
  0x0044E477: mov word ptr [ebx + 2], ax ; NumOfChargedBuff = 0
  0x0044E47B: push 4
  0x0044E482: call pp_SetPacketLen
  ```
- **En `XelionCore` ([`pkg/server/zone/handler.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/zone/handler.go#L284-L315)):**
  No existe envío de `0x104A` (`PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD`). Salta directamente de `0x103E` (`PASSIVE`) a `0x103D` (`SKILL`).

### Causa Crítica #3: `XelionCore` envía siempre `0x1047` (`PROTO_NC_CHAR_CLIENT_ITEM_CMD`) para las 4 cajas aunque estén vacías y con tamaño de atributo fijo (`65` / `37` bytes)
1. **Cajas vacías (`8, 9, 12` vs `15`):**
   - En `Zone.exe` (`so_StoreInventoryFromServer` @ `0x0044DDF0`), cuando `box != 15` (`0x0044DF5E: cmp al, 0xf; jne 0x44e006`), salta a `0x0044E006` y llama a `ci_FillBufferInventoryItem` (`0x00643D40`).
   - Si la caja (`8` Equip, `9` Inven o `12` House) tiene **0 ítems**, `ci_FillBufferInventoryItem` retorna `0` bytes (`eax = 0`), y en `0x0044E02F`:
     ```asm
     0x0044E02F: test eax, eax
     0x0044E031: jle  0x44e0ad ; <-- NO ENVÍA EL PAQUETE 0x1047 SI LA CAJA 8, 9 o 12 ESTÁ VACÍA
     ```
   - Únicamente para **Box `15` (`0x0F`, ActionItem / Emociones)** el código entra por `0x0044DF62..0x0044DFFF` y envía siempre el paquete `0x1047` aunque tenga `0` ítems (`numofitem = 0, box = 15, flag = 1`).
   - En cambio, `XelionCore` ([`pkg/server/zone/handler.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/zone/handler.go#L346-L377)) fuerza el envío de `0x1047` con `numofitem = 0` para las cajas `8, 9, 12, 15`.
2. **Tamaño de atributo dinámico de ítems (`iac_fill2buffer`) en `0x1047`:**
   - En `Zone.exe` (`ItemAttrCls_Weapon::iac_fill2buffer` @ `0x00640540`), el tamaño de un arma (`Class 5`, como la `ShortSword` ID `250` equipada al crear un Fighter) se calcula dinámicamente según las opciones aleatorias:
     ```asm
     0x00640558: movzx eax, byte ptr [edx + 0x40] ; byte en offset 64 (0x40): bitfield (num_options << 1) | Identified
     0x0064055C: shr   eax, 1                     ; eax = num_options (0 si está limpio)
     0x0064055E: lea   ecx, [edx + eax*2 + 0x40]
     0x00640563: lea   eax, [eax + ecx + 1]       ; total = 0x40 (64) + 1 + (num_options * 3) = 65 bytes + 3*N
     ```
     Para Armas (`Class 5`) con 0 opciones es `65` bytes (`datasize = 69`), pero para **Armaduras/Escudos/Botas/Brazaletes (`Class 6, 7, 8, 38` @ `0x00640570`)** con 0 opciones:
     ```asm
     0x00640588: movzx eax, byte ptr [edx + 0x0c] ; byte en offset 12 (0x0C): (num_options << 1) | Identified
     0x0064058C: shr   eax, 1                     ; eax = num_options (0)
     0x0064058E: lea   ecx, [edx + eax*2 + 0x0c]
     0x00640593: lea   eax, [eax + ecx + 1]       ; total = 0x0C (12) + 1 + (num_options * 3) = 13 bytes (¡NO 37 bytes!)
     ```
   - Para **Amuletos/Anillos (`Class 4` @ `0x00640510`)** con 0 opciones:
     ```asm
     0x00640529: movzx eax, byte ptr [edx + 0x25] ; byte en offset 37 (0x25)
     0x0064052D: shr   eax, 1
     0x0064052F: lea   ecx, [edx + eax*2 + 0x25]
     0x00640534: lea   eax, [eax + ecx + 1]       ; total = 0x25 (37) + 1 + (num_options * 3) = 38 bytes
     ```
   - Para **Cofres (`Class 15` @ `0x0063E8E0`)** con 0 ítems dentro: `1 + (count * 8) = 1 byte` (¡no 65!).
   - Para **BindItem (`Class 13` @ `0x0063E8B0`)** con 0 entradas: `1 + (count * 10) = 1 byte` (¡no 101!).
   - En el ítem inicial `46003` (`MVault_01`, `Class = 10` Decoration), `Zone.exe` (`0x0063E890`) envía `8` bytes (`datasize = 12`). Si `XelionCore` rellena con `0` un `Class 6/7/8` (Armor/Shield/Boot) enviando `37` ceros cuando el byte en offset `12` (`num_options`) es `0`, el cliente lee `num_options = 0` (13 bytes de atributo) y los otros 24 bytes desincronizan la lectura del paquete `0x1047`.

### Causa Adicional en `WorldService` (Justo al hacer click en el personaje, antes de conectar a Zone)
- **En `WorldManager.exe` (`fc_NC_CHAR_CHARDATA_ACK` @ `0x00410465` -> `SendToClient_CharOptionData` @ `0x00431680`):**
  Inmediatamente después de enviar `PROTO_NC_CHAR_LOGIN_ACK` (`0x1003`) al cliente en el socket de World, `WorldManager.exe` llama a `SendToClient_CharOptionData()` (`0x00431682`) y le envía automáticamente al cliente por el socket de World sus opciones guardadas:
  - Si `DataType == 2` (sistema mejorado NA2016):
    1. `0x7032` (`PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD`)
    2. `0x7033` (`PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD`)
    3. `0x7034` (`PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD`)
  - Si `DataType == 1` (sistema clásico):
    `0x7003`, `0x7005`, `0x7007`, `0x7009`, `0x700B`, `0x700D`, `0x700F`.
- **En `XelionCore` ([`pkg/server/world/handler.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/world/handler.go#L578-L638)):**
  `handleCharSelect` **solo envía `0x1003` (`PROTO_NC_CHAR_LOGIN_ACK`)** y no envía ningún paquete `PROTO_NC_CHAR_OPTION_*`. Además, en [`pkg/protocol/proto_nc_char.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/protocol/proto_nc_char.go#L127-L161), si el cliente llegara a pedir opciones clásicas, los tamaños de `0x700D` (`WINDOWPOS_ACK`) y `0x700F` (`KEYMAPPING_ACK`) en `XelionCore` son `137B` y `53B`, mientras que en `WorldManager.exe` (`0x00431BBB` y `0x00431C4A`) son **`0x189` (`393` bytes)** y **`0x135` (`309` bytes)** respectivamente.

---

## 2. Tabla Comparativa Global: Secuencia tras Click en Personaje

### Fase A: En el Servidor `World` (al recibir `PROTO_NC_CHAR_LOGIN_REQ` `0x1001`)

| Orden | `WorldManager.pdb` / `WorldManager.exe` (Oficial NA2016) | `XelionCore` (Go) [`world/handler.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/world/handler.go#L578-L638) | `NextGen` (C#) [`Handler4.cs`](file:///C:/Users/fperez/Desktop/Xelion/NextGen-Emulator/NextGen.World/Handlers/Handler4.cs#L15-L50) |
| :---: | :--- | :--- | :--- |
| **1** | Recibe `0x1001` (`PROTO_NC_CHAR_LOGIN_REQ`, `1B` slot). Pide datos a `Character.exe` (`0x1004`). Al recibir `0x101A` (`NC_CHAR_CHARDATA_ACK`), verifica si es Tutorial (`Slot` nuevo -> mapa `"Rou"`, `X=6478, Y=10558`) o Prisión (`PrisonMin > 0` -> `"Pri"`, `X=704, Y=962`). | Recibe `0x1001` (`1B` slot). Busca el personaje en MySQL y genera token en `TransferManager`. | Recibe `SH4Type.CharSelect` (`0x1001`, `1B` slot). Genera `TransferClient` hacia Zone. |
| **2** | Envía a `Character.exe` `0x1036` (`PROTO_NC_CHAR_REGISTNUMBER_CMD`) y a `Account.exe` `0xC806` (`PROTO_NC_USER_CHAR_BORING_CMD`). | *(Unificado en memoria/DB)* | *(Unificado en memoria)* |
| **3** | Envía al cliente **`0x1003` `PROTO_NC_CHAR_LOGIN_ACK`** (`18B` payload / `20B` wire: `char zoneip[16]`, `uint16 zoneport`). | Envía al cliente **`0x1003` `PROTO_NC_CHAR_LOGIN_ACK`** (`18B` payload / `20B` wire: `zoneip[16]`, `zoneport`). | Envía al cliente **`0x1003` `SH4Type.ServerIP`** (`18B` payload: `ip[16]`, `ushort port`). |
| **4** | Llama a `SendToClient_CharOptionData()` (`0x00431680`) y envía al cliente en World:<br>• **`0x7032` `PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD`** (`2 + N*9` bytes)<br>• **`0x7033` `PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD`** (`2 + N*4` bytes)<br>• **`0x7034` `PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD`** (`2 + N*3` bytes) | **❌ NO ENVÍA NADA.** Solo responde si el cliente envía explícitamente `0x7002..0x700E`. | No envía nada en `CharSelect`; responde a `SH28Type` (`0x7002..0x700E`) y `SH31Type` (`0x7C02`) cuando el cliente los pide. |

---

### Fase B: En el Servidor `Zone` (al recibir `PROTO_NC_MAP_LOGIN_REQ` `0x1801` -> Envío del Chunk)

Cuando el cliente se conecta al puerto TCP de `Zone`:
1. `Zone` envía `0x0807` (`PROTO_NC_MISC_SEED_ACK`, `2B` xor seed).
2. El cliente envía `0x1801` (`PROTO_NC_MAP_LOGIN_REQ`, `1590B` payload: `uint16 wldmanhandle`, `char charid[20]`, `Name8 checksum[49]`).
3. En `Zone.exe` (`sp_NC_MAP_LOGIN_REQ` @ `0x00452590`), valida los 49 hashes MD5 contra `DataFileCheck`, envía `0x104B` (`PROTO_NC_CHAR_ZONE_CHARDATA_REQ`) a `WorldManager.exe`, y a medida que `WorldManager.exe` (`fc_NC_CHAR_ZONE_CHARDATA_REQ` @ `0x00428BA0`) le inyecta los paquetes `NC_CHAR_ZONE_*`, **`Zone.exe` envía al cliente la siguiente secuencia exacta**:

| Orden Oficial | OpCode | Nombre Oficial (`Zone.pdb`) | `Zone.exe` Oficial (`Zone.pdb` NA2016) | `XelionCore` (Go) [`zone/handler.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/zone/handler.go#L214-L397) | `NextGen` (C#) [`ZoneCharacter.cs`](file:///C:/Users/fperez/Desktop/Xelion/NextGen-Emulator/NextGen.Zone/Game/ZoneCharacter.cs#L411-L434) |
| :---: | :---: | :--- | :--- | :--- | :--- |
| **1** | `0x1038` | `PROTO_NC_CHAR_CLIENT_BASE_CMD` | **Siempre se envía.**<br>Payload: **`105 B`** (Wire: `107 B`)<br>Handler: `so_charbasedata` (`0x0044D890`) | **Siempre se envía (#1).**<br>Payload: **`105 B`** (Wire: `107 B`)<br>✅ Idéntico a `Zone.pdb`. | **Siempre se envía (#1).**<br>Payload: **`93 B`** (versión vieja: `Name` 16B, sin `CurLP`, cola de 7B). |
| **2** | `0x1039` | `PROTO_NC_CHAR_CLIENT_SHAPE_CMD` | **Siempre se envía.**<br>Payload: **`4 B`** (Wire: `6 B`)<br>Handler: `so_charsharpdata` (`0x0044AE80`) | **Siempre se envía (#2).**<br>Payload: **`4 B`** (Wire: `6 B`)<br>✅ Idéntico a `Zone.pdb`. | **Siempre se envía (#2).**<br>Payload: **`4 B`**. |
| **3** | `0x103A` | `PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD` | **Siempre se envía.**<br>Payload: **`6 B + 32*N`** (Wire: `8 B` vacío)<br>`chrregnum`(u32), `bNeedClear=1`(u8), `nNum=0`(u8).<br>Handler: `so_StoreQuestDoing` (`0x0044B3C0`) | **Siempre se envía (#3).**<br>Payload: **`6 B`** (`chrregnum`, `1`, `0`).<br>✅ Idéntico a `Zone.pdb`. | **Siempre se envía (#3).**<br>Payload: **`6 B`** (`ID`, `1`, `0`). |
| **4** | `0x103B` | `PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD` | **Siempre se envía.**<br>Payload: **`12 B + 10*N`** (Wire: `14 B` vacío)<br>`chrregnum`(u32), `nTotalDoneQuest`(u16), `nTotalDoneQuestSize`(u16), `nDoneQuestCount`(u16), `nIndex`(u16).<br>Handler: `so_StoreQuestDone` (`0x0044B500`) | **Siempre se envía (#4).**<br>Payload: **`12 B`** (`chrregnum`, `0, 0, 0, 0`).<br>✅ Idéntico a `Zone.pdb`. | **Siempre se envía (#4).**<br>Payload: **`6 B`** (`int ID`, `short 0` — struct de cliente antiguo). |
| **5** | `0x10CE` | `PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD` | **Siempre se envía.**<br>Payload: **`6 B + 2*N`** (Wire: `8 B` vacío)<br>`chrregnum`(u32), `nNumOfReadQuest`(u16).<br>Handler: `so_StoreQuestRead` (`0x0044B640`) | **Siempre se envía (#5).**<br>Payload: **`6 B`** (`chrregnum`, `0`).<br>✅ Idéntico a `Zone.pdb`. | ❌ No existe en cliente antiguo. |
| **6** | `0x10D7` | `PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD` | **Siempre se envía.**<br>Payload: **`6 B + 32*N`** (Wire: `8 B` vacío)<br>`chrregnum`(u32), `nNumOfRepeatQuest`(u16).<br>Handler: `so_StoreQuestRepeat` (`0x0044B770`) | **Siempre se envía (#6).**<br>Payload: **`6 B`** (`chrregnum`, `0`).<br>✅ Idéntico a `Zone.pdb`. | ❌ No existe en cliente antiguo. |
| **7** | `0x103E` | `PROTO_NC_CHAR_CLIENT_PASSIVE_CMD` | **Siempre se envía.**<br>Payload: **`2 B + 2*N`** (Wire: `4 B` vacío)<br>`number`(u16), `skillid[N]`(u16).<br>Handler: `so_charpassivedata` (`0x0044E5A0`) | **Siempre se envía (#7).**<br>Payload: **`2 B`** (`passiveCount = 0`).<br>✅ Idéntico si tiene 0 pasivas. | **Siempre se envía (#6 en NextGen, después de ActiveSkills).**<br>Payload: **`2 B + 2*N`**. |
| **8** | `0x104A` | `PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD` | **⚠️ SIEMPRE SE ENVÍA** (incluso con 0 buffs).<br>Payload: **`2 B + 14*N`** (Wire: `4 B` vacío)<br>`NumOfChargedBuff = 0` (`uint16`).<br>Handler: `so_CharChargedBuffData` (`0x0044E1E0`) | **❌ NO SE ENVÍA NUNCA (Omitido en `XelionCore`).** | ❌ No lo envía en `SendGetIngameChunk`. |
| **9** | `0x103D` | `PROTO_NC_CHAR_CLIENT_SKILL_CMD` | Se envía si `nMaxNum > 0`.<br>Payload: **`10 B + 12*N`** (Wire: `12 B + 12*N`)<br>• `restempow` (`u8`)<br>• **`PartMark` (`u8`) = `0x03`** (`first=1 \| last=1`)<br>• `nMaxNum` (`u16`)<br>• `chrregnum` (`u32`)<br>• `number` (`u16`)<br>• `block[N]` (`12B` c/u)<br>Handler: `so_charskilldata` (`0x0044BF20`) | **⚠️ BUG CRÍTICO (#8):**<br>Envía **`PartMark = 1` (`0x01`)** en vez de **`3` (`0x03`)** en [`handler.go:295`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/zone/handler.go#L295).<br>El bit `last` (`0x02`) queda en `0` y el cliente espera más paquetes `0x103D`. | **Siempre se envía (#5 en NextGen).**<br>Payload antiguo: **`7 B + 12*N`** (`byte 0`, `int ID`, `ushort count`, sin `PartMark` ni `nMaxNum`). |
| **10** | `0x1048` | `PROTO_NC_CHAR_CLIENT_GAME_CMD` | **Siempre se envía** (antes de los ítems).<br>Payload: **`4 B`** (Wire: `6 B`)<br>`partyno` (`u16`, `0xFFFF`), `nRaidNo` (`u16`, `0xFFFF`).<br>Handler: `so_Storegamedata` (`0x0044C220`) | **Siempre se envía (#9).**<br>Payload: **`4 B`** (`0xFFFFFFFF`).<br>✅ Idéntico a `Zone.pdb`. | **Se envía en #13 (después de ítems y títulos).**<br>Payload: **`2 B`** (`ushort 0xFFFF`, sin `nRaidNo`). |
| **11** | `0x1047` | `PROTO_NC_CHAR_CLIENT_ITEM_CMD` *(Box `8`: Equip)* | **Solo se envía si la caja `8` tiene `> 0` ítems** (`0x0044E02F: test eax, eax; jle skip`).<br>Payload: `3 B + items` (`numofitem` u8, `box=8` u8, `flag=1` u8). | **Se envía siempre (#10)** incluso si tiene `0` ítems.<br>⚠️ Si equipa `Class 6/7/8/38`, rellena `37B` en vez de `13B`. | Se envía siempre (#7).<br>*(Nota: NextGen tiene un bug en `Handler4.cs:103` donde escribe `InventoryItems.Count` en vez de `EquippedItems.Count`).* |
| **12** | `0x1047` | `PROTO_NC_CHAR_CLIENT_ITEM_CMD` *(Box `9`: Inven)* | **Solo se envía si la caja `9` tiene `> 0` ítems.**<br>Payload: `3 B + items` (`numofitem` u8, `box=9` u8, `flag=1` u8). | **Se envía siempre (#11)** incluso si tiene `0` ítems. | Se envía siempre (#8). |
| **13** | `0x1047` | `PROTO_NC_CHAR_CLIENT_ITEM_CMD` *(Box `12`: House)* | **Solo se envía si la caja `12` tiene `> 0` ítems.**<br>Para un personaje sin muebles en MiniHouse, **`Zone.exe` NO envía `0x1047` para Box `12`**. | **⚠️ Se envía siempre (#12)** con `numofitem = 0, box = 12, flag = 1`. | Se envía siempre (#9) con `count = 0, box = 0x0C, unk = 1`. |
| **14** | `0x1047` | `PROTO_NC_CHAR_CLIENT_ITEM_CMD` *(Box `15`: ActionItem)* | **Siempre se envía** (`0x0044DF5E: cmp al, 0xf`), incluso con `0` ítems.<br>Payload vacío: **`3 B`** (Wire: `5 B`): `numofitem=0, box=15, flag=1`. | **Siempre se envía (#13).**<br>Payload vacío: **`3 B`** (`0, 15, 1`).<br>✅ Idéntico a `Zone.pdb`. | Se envía siempre (#10) (`0, 0x0F, 1`), más un 5to paquete `0x1047` corrupto de 2 bytes (`SendPremiumItemList`). |
| **15** | `0x1049` | `PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD` | **Siempre se envía.**<br>Payload: **`6 B + 2*N`** (Wire: `8 B + 2*N`)<br>`CurrentTitle`(u8), `CurrentTitleElement`(u8), `CurrentTitleMobID`(u16), `NumOfTitle`(u16), `TitleArray[N]`(2B c/u).<br>Handler: `wms_NC_CHAR_ZONE_CHARTITLE_CMD` (`0x005A39F0`) | **Siempre se envía (#14).**<br>Payload: **`6 B`** (`0, 0, 0, 0` — Nota: para personaje sin título en `Zone.exe`, `CurrentTitle = 0xFF` / `255` en `InitNewPlayer`, aunque `0` también es leído). | Se envía en #12.<br>Payload: **`6 B`** (`int 0, short 0`). |
| **16** | `0x3801`<br>`0x3802`<br>`0x3803` | `PROTO_NC_COLLECT_CARDREGIST_CMD` *(si aplica)* | En `wms_NC_CHAR_ZONE_CARDCOLLECT_CMD` (`0x0044CEF0`), envía datos de colección de cartas si el personaje tiene cartas/bookmarks/rewards registrados. | No se envía (0 cartas). | No existe en cliente antiguo. |
| **17** | `0x1802` | `PROTO_NC_CHAR_MAPLOGIN_ACK` *(alias `PROTO_NC_MAP_LOGIN_ACK`)* | **Siempre se envía al final** (cuando `Zone.exe` recibe `0x101D` `NC_CHAR_ZONE_CHARDATA_ACK` en `0x005A3590`).<br>Marca al jugador en el mapa (`so_MapMarking`), calcula todos los stats (`so_ply_ChangeAbility`) y envía **`0x1802`**.<br>Payload: **`242 B`** (Wire: **`244 B` / `0xF4`**). | **Siempre se envía al final (#15).**<br>Payload: **`242 B`** (Wire: **`244 B`**).<br>✅ Tamaño y estructura coinciden con `Zone.pdb`. | Se envía en #14 (`SH6Type.DetailedCharacterInfo`).<br>Payload: **`238 B`** (sin `MaxAp` u32). Luego envía `0x3021` y `0x302D`. |

---

## 3. Comparativa Detallada Byte a Byte de las Estructuras Clave

### 3.1. `0x1038` — `PROTO_NC_CHAR_CLIENT_BASE_CMD`
- **`Zone.pdb` / `Zone.exe` (`so_charbasedata` @ `0x0044D890`):** **105 bytes** de payload (`107` bytes en cable, `0x6B`).
- **`XelionCore` ([`WriteDetailedInfo` en `pkg/server/zone/character.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/zone/character.go#L118-L152)):** **105 bytes** de payload (`107` bytes en cable).
- **`NextGen` (`WriteDetailedInfo` en [`ZoneCharacter.cs`](file:///C:/Users/fperez/Desktop/Xelion/NextGen-Emulator/NextGen.Zone/Game/ZoneCharacter.cs#L436-L469)):** **93 bytes** de payload.

| Offset Payload | Campo Oficial (`Zone.pdb`) | Tipo / Tamaño | `Zone.exe` (`0x0044D890`) | `XelionCore` (`character.go`) | `NextGen` (`ZoneCharacter.cs`) |
| :---: | :--- | :---: | :--- | :--- | :--- |
| `0x00` (`0`) | `chrregnum` | `uint32` (4B) | `[edi + 0x564]` (Char RegNum) | `uint32(p.Char.ID)` ✅ | `int ID` (4B) |
| `0x04` (`4`) | `charid` | `char[20]` (`Name5`) | `[edi + 0x568]` (20B) | `WriteString(p.Char.Name, 20)` ✅ | `WriteString(Name, 16)` ❌ (16B) |
| `0x18` (`24`) | `slot` | `uint8` (1B) | `[edi + 0x57C]` | `p.Char.Slot` ✅ | `byte Slot` (1B) |
| `0x19` (`25`) | `level` | `uint8` (1B) | `[edi + 0x57D]` | `p.Char.Level` ✅ | `byte Level` (1B) |
| `0x1A` (`26`) | `exp` | `uint64` (8B) | `[edi + 0x580]` | `uint64(p.Char.Exp)` ✅ | `ulong Exp` (8B) |
| `0x22` (`34`) | `Rudra / Unk` | `uint32` (4B) | `[edi + 0x588]` | `uint32(0)` ✅ | `int 0` (4B) |
| `0x26` (`38`) | `hpstone` | `uint16` (2B) | `[edi + 0x58C]` | `uint16(p.StonesHP)` ✅ | `short StonesHP` (2B) |
| `0x28` (`40`) | `spstone` | `uint16` (2B) | `[edi + 0x58E]` | `uint16(p.StonesSP)` ✅ | `short StonesSP` (2B) |
| `0x2A` (`42`) | `curhp` | `uint32` (4B) | `[edi + 0x590]` | `uint32(p.Char.CurHP)` ✅ | `uint CurHP` (4B) |
| `0x2E` (`46`) | `cursp` | `uint32` (4B) | `[edi + 0x594]` | `uint32(p.Char.CurSP)` ✅ | `uint CurSP` (4B) |
| `0x32` (`50`) | `curlp` | `uint32` (4B) | `[edi + 0x598]` | `uint32(p.Char.CurLP)` ✅ | *(No existe en NextGen, escribe `int 0` antes de Stones)* |
| `0x36` (`54`) | `fame` | `uint32` (4B) | `[edi + 0x59C]` | `uint32(p.Fame)` ✅ | `int Fame` (4B) |
| `0x3A` (`58`) | `cen` | `uint64` (8B) | `[edi + 0x5A0]` | `uint64(p.Money)` ✅ | `long Money` (8B) |
| `0x42` (`66`) | `mapname` | `char[12]` (`Name3`) | `[edi + 0x5A8]` (12B) | `WriteString(mapName, 12)` ✅ | `WriteString(MapName, 12)` |
| `0x4E` (`78`) | `loc.x` | `int32` (4B) | `[edi + 0x5B4]` | `int32(p.Char.XPos)` ✅ | `int XPos` (4B) |
| `0x52` (`82`) | `loc.y` | `int32` (4B) | `[edi + 0x5B8]` | `int32(p.Char.YPos)` ✅ | `int YPos` (4B) |
| `0x56` (`86`) | `dir` | `uint8` (1B) | `[edi + 0x5BC]` | `p.Rotation` ✅ | `byte Rotation` (1B) |
| `0x57` (`87`) | `str, con, dex, int, men` | `5 x uint8` (5B) | `[edi + 0x5BD..0x5C1]` | `Str, Con, Dex, Int, Men` ✅ | `5 x byte` (5B) |
| `0x5C` (`92`) | `point` (`UsablePoints`) | `uint8` (1B) | `[edi + 0x5C2]` | `p.UsablePoints` ✅ | `byte UsablePoints` (1B) |
| `0x5D` (`93`) | `killpoint` | `uint32` (4B) | `[edi + 0x5C4]` | `uint32(p.KillPoints)` ✅ | Cola de 7 bytes (`ushort`, `int`, `byte`) |
| `0x61` (`97`) | `prisontime` | `uint16` (2B) | `[edi + 0x5C8]` | `uint16(p.PrisonMin)` ✅ | — |
| `0x63` (`99`) | `AdminLevel` / `Tail` | `6 bytes` | `[edi + 0x5CA..0x5CF]` (`u32 + u16`) | `Fill(5, 0)` + `WriteByte(Admin)` (6B) ✅ | — |

---

### 3.2. `0x103D` — `PROTO_NC_CHAR_CLIENT_SKILL_CMD`
Aquí está la diferencia exacta en el encabezado de `0x103D`:

```cpp
// ESTRUCTURA OFICIAL EN Zone.pdb (10 bytes de cabecera + 12 bytes * number):
struct PROTO_NC_CHAR_CLIENT_SKILL_CMD {
    unsigned char restempow;   // +0x00 (1B)
    struct {
        unsigned char first : 1; // bit 0 (0x01)
        unsigned char last  : 1; // bit 1 (0x02)
        unsigned char seq   : 6; // bits 2..7
    } PartMark;                // +0x01 (1B) -> DEBE SER 0x03 (first=1 | last=1) cuando entra en 1 paquete
    unsigned short nMaxNum;    // +0x02 (2B) -> Total de skills del personaje
    unsigned long  chrregnum;  // +0x04 (4B) -> ID de registro del personaje
    unsigned short number;     // +0x08 (2B) -> Cantidad de bloques en este paquete
    PROTO_SKILLREADBLOCKCLIENT block[number]; // +0x0A (12B cada uno)
};
```

| Campo | `Zone.exe` Oficial (`0x0044BF20`) | `XelionCore` (`zone/handler.go:293-310`) | `NextGen` (`Handler4.cs:79-92`) |
| :--- | :--- | :--- | :--- |
| `restempow` (`+0`) | `uint8` (`[edi + 0x7BA1]`) | `byte(0)` | `byte(0)` |
| `PartMark` (`+1`) | **`0x03`** (`first=1, last=1`, `or byte ptr [esi+3], 2` en `0x0044C06D`) | **`0x01` ❌ (`first=1, last=0`)** | *(No existe en NextGen)* |
| `nMaxNum` (`+2`) | `uint16` (Total skills) | `uint16(len(activeSkills))` | *(No existe en NextGen)* |
| `chrregnum` (`+4`) | `uint32` (Char ID) | `uint32(p.Char.ID)` | `int Character.ID` (`+1`) |
| `number` (`+8`) | `uint16` (Skills en este paquete) | `uint16(len(activeSkills))` | `ushort ActiveSkills.Count` (`+5`) |
| `block[i]` (`+10`) | `12 bytes`:<br>• `skillid` (`uint16`)<br>• `cooltime` (`uint32`)<br>• `empow` (`uint16`, bitfield 4x4 bits)<br>• `mastery` (`uint32`) | `12 bytes`:<br>• `sk.SkillNo` (`uint16`)<br>• `0` (`uint32`)<br>• `0` (`uint16`)<br>• `0` (`uint32`) | `12 bytes`:<br>• `ID` (`ushort`)<br>• `9000` (`int`)<br>• `Upgrades` (`ushort`)<br>• `int 0` (`4B`) |

---

### 3.3. `0x1047` — `PROTO_NC_CHAR_CLIENT_ITEM_CMD` y `PROTO_ITEMPACKET_INFORM`

```cpp
// ESTRUCTURA OFICIAL EN Zone.pdb:
struct PROTO_NC_CHAR_CLIENT_ITEM_CMD {
    unsigned char numofitem; // +0x00 (1B)
    unsigned char box;       // +0x01 (1B) -> 8 (Equip), 9 (Inven), 12 (House), 15 (ActionItem)
    struct {
        unsigned char clear : 1; // bit 0 (1 en el primer paquete de esa box, 0 en los siguientes)
        unsigned char dummy : 7;
    } flag;                  // +0x02 (1B) -> 0x01
    PROTO_ITEMPACKET_INFORM ItemArray[numofitem]; // +0x03
};

struct PROTO_ITEMPACKET_INFORM {
    unsigned char datasize;  // +0x00 (1B) -> Tamaño de `info` en bytes: 2 (loc) + 2 (id) + attrSize
    struct {
        ITEM_INVEN location;        // +0x01 (2B): uint16 = (box << 10) | (slot & 0x03FF)
        unsigned short id;          // +0x03 (2B): ItemID en ItemInfo.shn
        SHINE_ITEM_ATTRIBUTE attr;  // +0x05 (variable según ItemAttrCls_*::iac_fill2buffer)
    } info;
};
```

En `Zone.exe` (`ItemAttributeClassContainer` @ `0x00640170` + `iac_fill2buffer`), el tamaño exacto de `attr` cuando un ítem **no tiene opciones aleatorias ni entradas dinámicas** (`byte` de conteo = `0`) es:

| `ItemInfo.shn` `Class` | Nombre de Clase en `Zone.pdb` | Método `iac_fill2buffer` en `Zone.exe` | Tamaño Real en `Zone.exe` (sin opciones) | Tamaño que envía `XelionCore` ([`pkg/data/item/item.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/data/item/item.go#L15-L55)) | ¿Coincide? |
| :---: | :--- | :---: | :---: | :---: | :---: |
| `0, 11, 12, 14, 16, 19, 20, 22, 25, 28, 33, 34` | `ByteLot`, `SkillScroll`, `RecallScroll`, `UpSource`, `WeaponTitleLicence`, `UpgradeGems`, `Feed`, `Enchant`, `GBCoin` | `0x0063E830` / `0x0063E970` | **`1 B`** (`datasize = 5`) | `1 B` (`datasize = 5`) | ✅ Sí |
| `1, 3, 21, 35` | `WordLot`, `QuestLot`, `KQStep`, `ActiveSkill` | `0x0063E850` | **`2 B`** (`datasize = 6`) | `2 B` (`datasize = 6`) | ✅ Sí |
| `2, 18, 24, 26, 27, 37` | `DwrdLot`, `MiniHouseSkin`, `Amount`, `CostumWeapon`, `ActionItem`, `CostumShield` | `0x0063E950` | **`4 B`** (`datasize = 8`) | `4 B` (`datasize = 8`) | ✅ Sí |
| `4` | `ItemAttrCls_Amulet` | `0x00640510` | **`38 B` + `3 * numOpts`** (`datasize = 42` con 0 opts) | `38 B` (`datasize = 42`) | ✅ Sí (con 0 opts) |
| `5` | `ItemAttrCls_Weapon` | `0x00640540` | **`65 B` + `3 * numOpts`** (`datasize = 69` con 0 opts) | `65 B` (`datasize = 69`) | ✅ Sí (con 0 opts) |
| **`6, 7, 8, 38`** | **`ItemAttrCls_Armor` / `Shield` / `Boot` / `Bracelet`** | **`0x00640570`** | **`13 B` + `3 * numOpts`** (`datasize = 17` con 0 opts; `37 B` solo si tiene el máximo de 8 opciones!) | **`37 B` fijo ❌ (`datasize = 41`)** | **❌ NO (Desincroniza 24 bytes por cada armadura/escudo/bota/brazalete limpio)** |
| `9` | `ItemAttrCls_Furniture` | `0x0063E870` | **`36 B`** (`datasize = 40`) | `36 B` (`datasize = 40`) | ✅ Sí |
| `10` | `ItemAttrCls_Decoration` (ej. `46003` `MVault_01`) | `0x0063E890` | **`8 B`** (`datasize = 12`) | `8 B` (`datasize = 12`) | ✅ Sí |
| **`13`** | **`ItemAttrCls_BindItem`** | **`0x0063E8B0`** | **`1 B` + `10 * count`** (`datasize = 5` cuando `count=0`) | **`101 B` fijo ❌** | **❌ NO (Envía 100 bytes de más si `count=0`)** |
| **`15`** | **`ItemAttrCls_ItemChest`** | **`0x0063E8E0`** | **`1 B` + `8 * count`** (`datasize = 5` cuando `count=0`) | **`65 B` fijo ❌** | **❌ NO (Envía 64 bytes de más si `count=0`)** |
| **`17, 32`** | **`KingdomQuestItem` / `NoEffect`** | **`0x004021C0`** | **`0 B`** (`datasize = 4`) | **`1 B` ❌** | **❌ NO** |
| `23` | `ItemAttrCls_Riding` | `0x0063E910` | **`19 B`** (`datasize = 23`) | `19 B` | ✅ Sí |
| `29` | `ItemAttrCls_Capsule` | `0x0063E980` | **`12 B`** (`datasize = 16`) | `12 B` | ✅ Sí |
| `30` | `ItemAttrCls_CardCollect_Close` | `0x0063E9B0` | **`9 B`** (`datasize = 13`) | `9 B` | ✅ Sí |
| `31` | `ItemAttrCls_CardCollect` | `0x0063E9E0` | **`5 B`** (`datasize = 9`) | `5 B` | ✅ Sí |
| `36` | `ItemAttrCls_Pet` | `0x0063EA00` | **`26 B`** (`datasize = 30`) | `26 B` | ✅ Sí |

---

### 3.4. `0x1802` — `PROTO_NC_CHAR_MAPLOGIN_ACK` (`PROTO_NC_MAP_LOGIN_ACK`)
- **En `Zone.exe` (`wms_NC_CHAR_ZONE_CHARDATA_ACK` @ `0x005A37C1` -> `so_fill_PROTO_NC_CHAR_MAPLOGIN_ACK` @ `0x005A3140`):**
  Envía exactamente **`0xF4` (`244` bytes en cable / `242` bytes de payload)**:
  - `+0x00` (`2B`): `uint16 handle` (Object Handle en el mapa de `Zone`, asignado por `ShineObjectManager`).
  - `+0x02` (`116B`): `CHAR_PARAMETER_DATA param` (`29 x uint32`: MaxHP, MaxSP, MaxLP, Direct5Stats, IndirectStats, Min/Max WC, AC, WalkSpeed, RunSpeed, etc.).
  - `+0x76` (`104B`): `PROTO_EQUIPMENT shape` (`26` slots de equipo visibles + `upgrade` bitfields + `belong` + `CosEff`).
  - `+0xDE` (`20B`): Coordenadas de casa/comercio, `PKCount` (`u32`), `AutoPick` (`u8`), `Bitfields` (`u8`), `nTemp` (`u16`), `MoverHandle` (`u16`), `MoverSlot` (`u16`).
- **En `XelionCore` ([`WriteDetailedInfoExtra` en `pkg/server/zone/character.go`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/zone/character.go#L223-L299)):**
  Envía exactamente **`242` bytes de payload (`244` bytes en cable)** con la misma estructura y offsets que `so_fill_PROTO_NC_CHAR_MAPLOGIN_ACK`.
- **En `NextGen` ([`WriteDetailedInfoExtra` en `ZoneCharacter.cs`](file:///C:/Users/fperez/Desktop/Xelion/NextGen-Emulator/NextGen.Zone/Game/ZoneCharacter.cs#L483-L526)):**
  Envía **`238` bytes de payload** porque pertenece a una versión anterior del cliente que no tenía el campo `MaxLP` (`uint32`, 4 bytes) dentro de `CHAR_PARAMETER_DATA`.

---

## 4. Conclusión Puntual

El cliente oficial (`Fiesta.bin`) no envía `PROTO_NC_MAP_LOGINCOMPLETE_CMD` (`0x1803`) después de recibir el chunk de `XelionCore` principalmente por:
1. **`0x103D` (`PROTO_NC_CHAR_CLIENT_SKILL_CMD`) lleva `PartMark = 1` en lugar de `PartMark = 3`** ([`pkg/server/zone/handler.go#L295`](file:///C:/Users/fperez/Desktop/Xelion/XelionCore/pkg/server/zone/handler.go#L295)), dejando abierto el stream de skills (`last == 0`).
2. **Falta el paquete obligatorio `0x104A` (`PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD`, 2 bytes de payload = `0x0000`)** entre `0x103E` (`PASSIVE`) y `0x103D` (`SKILL`).
3. **En `0x1047` (`PROTO_NC_CHAR_CLIENT_ITEM_CMD`)**, `Zone.exe` **no envía paquetes para las cajas `8`, `9` ni `12` si están vacías** (solo envía vacía la caja `15`), y si el personaje tiene equipada o en inventario alguna armadura/escudo/bota/brazalete (`Class 6, 7, 8, 38`), `Zone.exe` envía **`13` bytes de atributo** (cuando tiene 0 opciones) mientras que `XelionCore` envía **`37` bytes**, rompiendo el parser de `0x1047`.
