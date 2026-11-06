# Fase 06 — World Coordination, Topología y Subsistemas Sociales

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_WORLDSERVICE_CHARACTER.md](./PLAN_ACCION_WORLDSERVICE_CHARACTER.md) (**Sección 4 y 5**) y [`docs/Servicios/WorldManager.md`](../Servicios/WorldManager.md).
- **Dependencias previas:** Fase 05 completada sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 06 (`docs/plan/FASE_06_WORLD_COORDINATION_TOPOLOGY_SUBSYSTEMS.md`) leyendo como especificación obligatoria las Secciones de Coordinación de `docs/plan/PLAN_ACCION_WORLDSERVICE_CHARACTER.md` y `docs/Servicios/WorldManager.md`.
Regla Mandatoria: CERO datos hardcodeados y CERO SQL inline. Los gremios, academia, amigos, parties y filtros de chat persisten y consultan exclusivamente a través de los SPs oficiales en MySQL contra las tablas `tGuild*`, `tFriend*`, `tParty*` y `tSpam*`.
Implementá en `pkg/server/world/` los subsistemas de coordinación global: chat y mensajería (susurros, canales, gritos, spam filter con tSpamBlockTime), administración de gremios y academia (PROTO_NC_GUILD_*), lista de amigos (PROTO_NC_FRIEND_*), grupos y party (PROTO_NC_PARTY_*), y sincronización de topología S2S con ZoneService. Creá `pkg/audit/fase06_world_social_coverage_test.go` (Control 6) y verificá que compile y pase al 100%.
No avances a la Fase 07. Al terminar, actualizá el estado de la Fase 6 a "Completado" en `docs/plan/README.md` y hacé un git commit limpio en main.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Implementar Subsistemas Sociales en [`pkg/server/world/social_handler.go`](../../pkg/server/world/):**
   - **Chat y Mensajería Global:**
     * Susurros (`fc_NC_CHAT_WHISPER_REQ`), chat de mapa, canales y gritos.
     * Verificación de bloqueo de spammers en `tSpamBlockTime` y `tSpamFilter`.
   - **Gremios y Academia (`PROTO_NC_GUILD_*`):**
     * Creación de gremio (`p_Guild_Create`), disolución, ascensos/descensos de rango (`p_Guild_ChangeGrade`).
     * Invitación, ingreso y graduación de alumnos en la academia (`p_Guild_Academy_*`).
   - **Lista de Amigos (`PROTO_NC_FRIEND_*`):**
     * Solicitud de amistad, aceptación, borrado (`p_Friend_Delete`) y notificación de estado online/offline a contactos.
   - **Party / Grupo (`PROTO_NC_PARTY_*`):**
     * Creación de grupo (`p_Party_Create`), invitación de miembros, cambio de líder y reglas de distribución de botín.

2. **Topología Dinámica y Balanceo de Zonas en [`pkg/server/world/zone_manager.go`](../../pkg/server/world/):**
   - Registro de zonas activas (`ZoneService`), heartbeat periódico y control de carga.
   - Coordinación de transferencias inter-mapas (`PROTO_NC_CHAR_MAPLOGIN_REQ`).

3. **Crear [`pkg/audit/fase06_world_social_coverage_test.go`](../../pkg/audit/fase06_world_social_coverage_test.go) (Control 6):**
   - Test automatizado que audita el despacho e interacción de los subsistemas sociales y de topología.

---

## Criterio de Cierre de la Fase 06

1. `go build ./...` -> OK.
2. `go test -v ./pkg/audit -run "^TestFase06_WorldSocial_Control6"` -> `PASS`.
3. `go test -v ./pkg/server/world/...` -> `PASS`.
4. Marcar Fase 6 como `Completado` en `docs/plan/README.md`.
5. Ejecutar commit limpio:
   `git add -A && git commit -m "feat(world): Fase 06 - coordinacion global de mundo, chat, gremios, amigos y party 100% MySQL"`
