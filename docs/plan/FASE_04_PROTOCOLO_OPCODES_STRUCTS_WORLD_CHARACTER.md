# FASE 04: Protocolo Oficial de Red: OpCodes y Estructuras PDB 1:1

## 1. Objetivo Técnico
Definir con precisión binaria de 100% y sin ninguna omisión (`...` prohibido) todas las estructuras de red C++ y OpCodes oficiales de `Character.exe` y `WorldManager.exe` documentadas en [docs/Servicios/Character.md](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/Servicios/Character.md) y [docs/Servicios/WorldManager.md](file:///Users/fedemac/Desktop/Xelion%20Online/XelionCore/docs/Servicios/WorldManager.md).

---

## 2. Archivos a Modificar / Crear
- `pkg/protocol/opcodes.go`: Constantes oficiales `PROTO_NC_CHAR_*`, `PROTO_NC_WORLD_*`, `PROTO_NC_GUILD_*`, `PROTO_NC_FRIEND_*`, `PROTO_NC_PARTY_*`, `PROTO_NC_KQ_*`, `PROTO_NC_AVATAR_*`, calculadas con la fórmula oficial `(department << 10) | command`.
- `pkg/protocol/pdb/proto_nc_char.go`: Estructuras de ciclo de vida de avatares, stats, apariencia y login de personaje.
- `pkg/protocol/pdb/proto_nc_item.go`: Estructuras de inventario, equipamiento, opciones mágicas y comercio.
- `pkg/protocol/pdb/proto_nc_guild.go`: Estructuras de gremios, academia, rangos y torneos.
- `pkg/protocol/pdb/proto_nc_social.go`: Estructuras de amigos, party, susurros y chat global.
- `pkg/audit/fase04_world_character_sizeof_test.go`: Test automatizado del **Control 4** verificando la paridad al byte (`binary.Size`) contra las especificaciones de los PDBs.

---

## 3. Estructuras Oficiales Críticas a Definir
1. **Selección y Carga de Avatar:**
   - `PROTO_NC_CHAR_CHARDATA_REQ` (OpCode `0x1004`, 2B handle + 4B userNo)
   - `PROTO_NC_CHAR_CHARDATA_ACK` (OpCode `0x1005`, cabecera + N slots de personajes)
   - `PROTO_NC_CHAR_LOGIN_REQ` (OpCode `0x1001`, selección de avatar)
   - `PROTO_NC_CHAR_LOGIN_ACK` (OpCode `0x1002`, token de transferencia a ZoneServer)
2. **Creación y Borrado:**
   - `PROTO_NC_CHAR_CREATE_REQ` (OpCode `0x1007`, datos cosméticos, clase, nombre)
   - `PROTO_NC_CHAR_CREATE_ACK` (OpCode `0x1008`, resultado y slot asignado)
   - `PROTO_NC_CHAR_DELETE_REQ` (OpCode `0x1009`, confirmación de baja)
   - `PROTO_NC_CHAR_DELETE_ACK` (OpCode `0x100A`, respuesta oficial)
3. **Persistencia S2S (CharServer -> ZoneServer / WorldServer):**
   - `PROTO_NC_CHAR_LOGIN_DB` (OpCode `0x1036`)
   - `PROTO_NC_CHAR_LOGOUT_DB` (OpCode `0x1037`)
   - `PROTO_NC_CHARSAVE_ALL_REQ`
   - `PROTO_NC_CHAR_STAT_INCPOINT_DB_REQ` (OpCode `0x105D`)

---

## 4. Criterios de Aceptación y Cierre
1. El 100% de los OpCodes coincide con la fórmula bitfield oficial `(dept << 10) | cmd`.
2. `pkg/audit/fase04_world_character_sizeof_test.go` valida que el 100% de los `sizeof` binarios en Go coincide con los offsets de C++ documentados.
3. Cero alias deprecados o nombres inventados.
4. `go build ./...` compila limpiamente.
5. `git commit` en `main`:
   `git add -A && git commit -m "feat(protocol): Fase 04 - 100% estructuras y opcodes oficiales de Character y WorldManager PDB"`
