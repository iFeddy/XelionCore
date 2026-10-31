# FASE 07: Integración Transversal ZoneService y Suite de Auditoría Final

## 1. Objetivo Técnico
Completar la integración transaccional entre `WorldService` y `ZoneService`, garantizando la transferencia fluida del personaje hacia el mapa de destino (`Zone`), la carga de sus estadísticas, inventario y habilidades desde MySQL, la verificación de los 54 SHNs de `Character.exe` y 71 SHNs de `WorldManager.exe`, y la ejecución automatizada de los 8 Controles de Auditoría mediante `scripts/verify_world_character_100.sh`.

---

## 2. Archivos a Modificar / Crear
- `pkg/server/zone/transfer_handler.go`: Recepción del token de transferencia desde WorldService, autenticación del avatar y despacho de datos iniciales al cliente.
- `pkg/server/world/zone_bridge.go`: Protocolo S2S de notificación de cambio de zona y persistencia de ubicación (`sLoginZone`, `X`, `Y`).
- `pkg/audit/world_character_completeness_test.go`: Suite de pruebas forenses de ciclo completo:
  - Creación de cuenta -> Login -> Selección de Mundo -> Creación de Personaje -> Selección -> Ingreso a Zona -> Modificación de Inventario/EXP -> Logout -> Validación de filas en MySQL.
- `scripts/verify_world_character_100.sh`: Script bash ejecutable que corre los 8 Controles Obligatorios de la Etapa 2.

---

## 3. Controles Obligatorios de la Suite (`verify_world_character_100.sh`)
1. **Control 1:** Esquema relacional de las 61 tablas de `World00_Character` en `xelion_online`.
2. **Control 2:** Catálogo exhaustivo de los 311 Stored Procedures oficiales en `information_schema.ROUTINES`.
3. **Control 3:** Análisis AST en Go confirmando cero SQL inline y 100% `CALL` en repositorios.
4. **Control 4:** Paridad de OpCodes y `sizeof` al byte contra `Character.md` y `WorldManager.md`.
5. **Control 5:** Cobertura de handlers de avatares en el Dispatcher de `WorldService`.
6. **Control 6:** Verificación de consistencia de los 54+71 archivos `.shn`.
7. **Control 7:** Prueba E2E integral forense de punta a punta.
8. **Control 8:** Cero constantes mágicas y cero hardcode de datos de juego.

---

## 4. Criterios de Aceptación y Cierre
1. `go build ./...` compila limpiamente sin warnings ni errores.
2. `go test -count=1 ./...` pasa 100% en verde en todos los paquetes del repositorio.
3. `./scripts/verify_world_character_100.sh` finaliza con código de salida `0` y todos los controles en `[PASS]`.
4. `git commit` en `main`:
   `git add -A && git commit -m "feat(core): Fase 07 - integracion completa WorldService y Character con Zone y auditoria 100% PASS"`
