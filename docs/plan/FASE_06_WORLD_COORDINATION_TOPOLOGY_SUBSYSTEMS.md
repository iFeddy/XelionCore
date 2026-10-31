# FASE 06: World Coordination, Topología y Subsistemas Sociales

## 1. Objetivo Técnico
Completar el subsistema de coordinación global de `WorldManager.exe` en `WorldService` (`pkg/server/world/`), incluyendo enrutamiento de mensajería (chat general, susurros, canales, gritos), gremios, academia, lista de amigos, grupos / party, eventos Kingdom Quest (KQ) y coordinación de topología de zonas con `ZoneService`.

---

## 2. Archivos a Modificar / Crear
- `pkg/server/world/social_handler.go`: Handlers oficiales de chat, amigos y gremios:
  - Mensajería privada y verificación de destinatario en línea.
  - Creación y administración de gremios (`PROTO_NC_GUILD_*`), rangos de miembros e invitaciones a academia.
  - Lista de amigos (`PROTO_NC_FRIEND_*`), actualización de estado conectado/desconectado y notas personales.
  - Formación de Party (`PROTO_NC_PARTY_*`), asignación de líder, distribución de drops y expulsiones.
- `pkg/server/world/zone_manager.go`: Mapeo dinámico de zonas de juego activas (`ZoneService`), puertos de escucha y balanceo de carga.
- `pkg/audit/fase06_world_social_coverage_test.go`: Test automatizado del **Control 6** validando el enrutamiento de eventos sociales y topología S2S.

---

## 3. Puntos Críticos de Coordinación
1. **Chat y Filtro de Spammers:** Integración con `tSpamBlockTime` y `tSpamFilter` para silenciar automáticamente mensajes abusivos o repetitivos.
2. **Eventos Globales:** Emisión de anuncios de sistema y mensajes de servidor a todos los avatares activos en el clúster.
3. **Persistencia S2S de Gremios:** Cada cambio de rango, expulsión o ingreso se persiste inmediatamente en MySQL mediante los Stored Procedures correspondientes (`p_Guild_AddMember`, `p_Guild_ChangeGrade`, etc.).

---

## 4. Criterios de Aceptación y Cierre
1. El despachador procesa correctamente paquetes de chat, gremios, amigos y grupos entre múltiples conexiones simultáneas.
2. `pkg/audit/fase06_world_social_coverage_test.go` pasa 100% en verde.
3. `go build ./...` compila limpiamente.
4. `git commit` en `main`:
   `git add -A && git commit -m "feat(world): Fase 06 - coordinacion global de mundo, chat, gremios, amigos y party 100% MySQL"`
