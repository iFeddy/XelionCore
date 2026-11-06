# Fase 07 — Integración Transversal ZoneService y Suite de Auditoría Final

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Sección 5 y Controles 1 a 8**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fase 06 completada sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 07 (`docs/plan/FASE_07_INTEGRACION_ZONE_GAMEPLAY_Y_AUDITORIA_FINAL.md`) leyendo como especificación obligatoria la Sección 5 y Controles de Auditoría de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: CERO datos hardcodeados y CERO tablas inventadas. Toda la información estática del juego proviene exclusivamente de los JSON y SHN de `data/` (ClassParam, Map, Item, etc.), y toda la persistencia y datos de personajes provienen 100% de los Stored Procedures oficiales y tablas de MySQL.
Completá la integración transversal WorldService -> ZoneService: traspaso de avatar autenticado mediante token, carga inicial de stats e inventario desde MySQL y data/, verificación de los 54 SHNs de Character y 71 SHNs de WorldManager. Creá `pkg/audit/world_character_completeness_test.go` (Control 7 E2E y Control 8 Cero Hardcode) y el script ejecutable `scripts/verify_world_character_100.sh`. No te detengas hasta que la suite completa pase al 100% (8 controles en [PASS] con Exit 0).
Al terminar, actualizá el estado de la Fase 7 a "Completado" en `docs/plan/README.md` y hacé un git commit limpio en main.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Integración Transversal en [`pkg/server/zone/`](../../pkg/server/zone/):**
   - Recepción y consumo de un solo uso del token de transferencia emitido por WorldService.
   - Carga atómica de datos del personaje:
     * Stats base calculados desde [`data/World/ClassParam/<Class>.json`](../../data/World/ClassParam/) + stats asignados en `tCharacter`.
     * Carga de inventario y equipamiento desde `tItem` y `tItemOptions` mapeados contra `ItemInfo.shn`.
     * Carga de habilidades desde `tSkill`.
     * Carga de buffs y estados anormales activos desde `tCharacterState`.
   - Persistencia periódica y de desconexión en MySQL mediante `p_Char_Logout` (HP, SP, EXP, Dinero, Zona, X, Y, D).

2. **Verificación de Consistencia de Archivos SHN:**
   - Validar que el motor cargue y verifique los 54 archivos `.shn` de `Character` y los 71 archivos `.shn` de `WorldManager` en memoria con sus checksums oficiales.

3. **Crear [`pkg/audit/world_character_completeness_test.go`](../../pkg/audit/world_character_completeness_test.go):**
   - **Control 7 (E2E Integral Forense):** Flujo completo de conexión: Login -> WorldSelect -> CharacterList -> CharacterCreate -> CharacterSelect -> ZoneTransfer -> ItemMove -> LevelUp -> Logout -> Verificación de columnas exactas en MySQL `xelion_online`.
   - **Control 8 (Cero Datos Hardcodeados):** Escaneo sintáctico AST verificando ausencia de constantes mágicas en stats, equipamiento y fórmulas.

4. **Crear [`scripts/verify_world_character_100.sh`](../../scripts/verify_world_character_100.sh):**
   - Script bash ejecutable que ejecuta secuencialmente los 8 Controles de Auditoría:
     * Control 1: Esquema de 61 tablas en `information_schema.TABLES`.
     * Control 2: Catálogo de 311 Stored Procedures en `information_schema.ROUTINES`.
     * Control 3: Escaneo AST (Cero SQL inline en repositorios de Character).
     * Control 4: Paridad de OpCodes y `sizeof` al byte contra PDBs de Zone y WorldManager.
     * Control 5: Cobertura de handlers de avatares en Dispatcher de WorldService.
     * Control 6: Verificación de consistencia de los 54+71 archivos `.shn`.
     * Control 7: Prueba E2E integral forense de punta a punta.
     * Control 8: Cero constantes mágicas y cero datos hardcodeados.

---

## Criterio de Cierre de la Fase 07

1. `go build ./...` -> OK sin errores ni warnings.
2. `go test -count=1 ./...` -> `PASS` en todos los paquetes del repositorio.
3. `./scripts/verify_world_character_100.sh` -> **100% PASS (8/8 controles en verde, Exit 0)**.
4. Marcar Fase 7 como `Completado` en `docs/plan/README.md`.
5. Ejecutar commit limpio:
   `git add -A && git commit -m "feat(core): Fase 07 - integracion completa WorldService y Character con Zone y auditoria 100% PASS"`
