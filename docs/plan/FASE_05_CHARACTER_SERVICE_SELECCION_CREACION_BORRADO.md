# FASE 05: Character Service: Selección, Creación, Borrado y Carga de Personajes

## 1. Objetivo Técnico
Completar la lógica de negocio oficial de `Character.exe` integrada en `WorldService` (`pkg/server/world/`), implementando la deserialización estricta de paquetes C2S/S2S, invocación de Stored Procedures en MySQL y emisión de respuestas binarias sin valores inventados ni datos hardcodeados.

---

## 2. Archivos a Modificar / Crear
- `pkg/server/world/character_handler.go`: Despachador de protocolos oficiales de gestión de personajes:
  - `handle_PROTO_NC_CHAR_CHARDATA_REQ`: Consulta a MySQL vía `p_Char_GetListOfUserChar` y serialización de avatares con apariencia y equipamiento visible.
  - `handle_PROTO_NC_CHAR_CREATE_REQ`: Creación de avatar validando nombre único (`p_Char_CheckName`), inicialización de stats según clase (`DefaultCharacterData.txt`), equipamiento base e inserción atómica (`p_Char_Create`).
  - `handle_PROTO_NC_CHAR_DELETE_REQ`: Soft-delete oficial (`p_Char_Delete`), validación de no pertenencia como líder de gremio y persistencia de fecha de baja.
  - `handle_PROTO_NC_CHAR_LOGIN_REQ`: Selección de personaje, carga de datos completos (`p_Char_GetLoginData`), emisión de token de zona de un solo uso y respuesta con coordenadas iniciales (`sLoginZone`, `X`, `Y`).
- `pkg/audit/fase05_world_character_dispatcher_test.go`: Test automatizado del **Control 5** verificando la cobertura total de handlers de gestión de personajes.
- `pkg/server/world/character_e2e_test.go`: Pruebas de integración E2E de ciclo de vida de avatares.

---

## 3. Puntos Críticos de Fidelidad
1. **Unicidad de Nombre:** Validación case-insensitive de nombres de avatar sin depender de la collation del cliente.
2. **Cero Hardcode de Equipamiento Base:** El equipamiento inicial de cada clase (arma de madera, ropa de aprendiz, pociones básicas) se obtiene de la base de datos o de los datos oficiales cargados en memoria desde los archivos del servidor, nunca de IDs fijos en código.
3. **Control de Slots:** Asignación secuencial de slots disponibles (0 al 4) y rechazo si la cuenta alcanza el cupo máximo permitido.

---

## 4. Criterios de Aceptación y Cierre
1. El despachador procesa el flujo completo: Listado -> Creación -> Borrado -> Selección.
2. `pkg/audit/fase05_world_character_dispatcher_test.go` pasa 100% en verde.
3. `pkg/server/world/character_e2e_test.go` verifica que los avatares se crean, modifican y borran correctamente en MySQL `xelion_online`.
4. `go build ./...` compila limpiamente.
5. `git commit` en `main`:
   `git add -A && git commit -m "feat(world): Fase 05 - gestion oficial de avatares (lista, creacion, borrado y seleccion) 100% MySQL"`
