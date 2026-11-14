# Fase 13 — Motor de Monstruos, IA (`MobBehavior`), Patrullaje (`MobRoam`), Tablas Oficiales de `Drops`, `AutoPick`, Títulos (`CT`) y Cartas (`COLLECT`)

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (**Secciones 1-Bugs #8 y #15, 4.1-Depts 4/24/49 y 6.3**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fases 08 a 12 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 13 (`docs/plan/FASE_13_MOTOR_MONSTRUOS_IA_MOBBEHAVIOR_DROPS_Y_COLECCION.md`) leyendo obligatoriamente las Secciones 1 (Bugs #8 y #15), 4.1 (Depts 4, 24 y 49) y 6.3 de `docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md`, junto con `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: 100% de los spawns, rutas, IA, secuencias de ataque, drops, títulos y cartas provienen de `data/MobRegen/`, `data/MobRoam/`, `data/MobSetting/`, `data/MobBehaviorDescript.json`, `data/MobAttackSequence/`, `data/World/Drop/` y tablas `.shn`; CERO drops o monstruos hardcodeados.
1) Incorporá `ZoneMob` vivo en `ZoneMap` alimentado por `data/MobRegen/*.json`, `MobInfo.shn`, `MobInfoServer.shn` y `MobSpecies.shn`, ejecutando en el Game Loop su patrullaje (`MobRoam`), detección/sigilo (`MobSetting` + `StealthMatrix.json`), tabla de aggro, persecución validada con `.shbd`, retorno (leash), llamadas de auxilio, diálogos (`MobChat.json`), secuencias de ataque (`MobAttackSequence`) y máquina `MobBehaviorDescript.json` + scripts Lua de `data/LuaScript/AIScript/`.
2) Al morir un monstruo, calculá y distribuí EXP/Fama individual, en `Party` y en `Raid` (`ExpTable.json`, `RestExp`), generá Cen e ítems con opciones aleatorias desde `ItemDropTable.json`, `ItemDropGroup.json`, `WorldDrop.shn` y `RandomOptionTable.json`, implementá propiedad temporal, limpieza de suelo, `PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ` (`0x10BC` -> `CALL usp_Char_SetAutoPick` / `tAutoPick`), anuncios `MobKillAnnounce.shn` y logs `GameLog` Tipo `6` + `PROTO_NC_LOG_GAME_MK*`.
3) Implementá el sistema completo de Títulos (`Dept 24`: `PROTO_NC_CT_SET_CURRENT_REQ` `0x6001` + progreso automático de `CharacterTitleData.shn` vía `CALL p_Char_*Title*`) y Colección de Cartas (`Dept 49`: `COLLECT_PEEPING_REQ`, `COLLECT_CARDREGIST_REQ`, `COLLECT_BOOKMARK_REGIST_REQ`, `COLLECT_REWARD_REQ` con `CollectCard*.shn` y `CALL p_Collection_*`).
4) Creá tests completos de ciclo de vida de monstruos, aggro/roam, drops probabilísticos, `AutoPick`, títulos y cartas. Al terminar, marcá Fase 13 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 14.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Spawns, Reaparición e IA Completa de Monstruos (`pkg/server/zone/mob.go`, `mob_ai.go`):**
   - Instanciar al cargar cada mapa todos sus grupos y áreas de `data/MobRegen/<Map>.json` resolviendo stats desde `MobInfo.shn` y `MobInfoServer.shn`, difundiendo apariciones por AOI (`PROTO_NC_BRIEFINFO_MOB_CMD` / `REGENMOB_CMD`).
   - Implementar la máquina de estados oficial (`Idle` -> `Roam` -> `Aggro/Chase` -> `Attack/Skill` -> `Flee/CallHelp` -> `Return/Reset` -> `Dead/RespawnTimer`) guiada por `data/MobRoam/*.json`, `data/MobSetting/*.json`, `data/MobBehaviorDescript.json`, `data/MobAttackSequence/*.json`, `MobChat.json`, `StealthMatrix.json`, colisión `.shbd` y scripts Lua `AIScript`.

2. **Reparto de Experiencia, Tablas Oficiales de `Drops` y `AutoPick` (`pkg/server/zone/drops.go`):**
   - Calcular recompensa de EXP y Fama según contribución de daño, penalización/bono de nivel, bonificación de `Party` / `Raid` y `RestExp` (`tRestExp`), disparando subida de nivel (`PROTO_NC_CHAR_LEVELUP_*` + `CALL p_Char_SetLevelAndExp` + `GameLog` Tipo `2`).
   - Evaluar `ItemDropTable.json`, `ItemDropGroup.json` y `WorldDrop.shn` generando drops de ítems (con `RandomOptionTable.json` cuando corresponda) y Cen en el suelo con reserva temporal de dueño, expiración automática y soporte completo de `PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ` (`0x10BC`), más anuncios `MobKillAnnounce.shn` y auditoría `GameLog` Tipo `6` (`p_GameLog_Add_Type6`, `PROTO_NC_LOG_GAME_MKED`, `PROTO_NC_LOG_GAME_MK_DROP_ITEM`).

3. **Títulos de Personaje (`NC_CT` — `Dept 24`) y Colección de Cartas (`NC_COLLECT` — `Dept 49`):**
   - **Títulos:** `PROTO_NC_CT_SET_CURRENT_REQ` (`0x6001`), incremento automático de contadores por muerte de monstruos, nivel, fama, misiones o producción según `CharacterTitleData.shn`, desbloqueo de títulos y persistencia vía `CALL p_Char_AddTitle` / `CALL p_Char_SetTitle`.
   - **Cartas de Colección (`4` C2S):** `PROTO_NC_COLLECT_PEEPING_REQ` (`0xC401`), `PROTO_NC_COLLECT_CARDREGIST_REQ` (`0xC405`), `PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ` (`0xC409`) y `PROTO_NC_COLLECT_REWARD_REQ` (`0xC40D`) contra `CollectCard*.shn` y `CALL p_Collection_*` (`tCollectionCard`, `tCollectionFavorites`, `tCollectionReward`, `tCollectionSerialCount`).

---

## Criterio de Cierre de la Fase 13

1. `go build ./...` -> `OK`.
2. `go test -p 1 -count=1 -v ./pkg/server/zone/...` -> `PASS` (incluyendo tests de `MobRegen`, IA `MobBehavior`, `Drops`, `AutoPick`, `CT` y `COLLECT`).
3. `go test -p 1 -count=1 ./...` -> `PASS`.
4. Marcar Fase 13 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(zone): Fase 13 - motor completo de monstruos, IA MobBehavior, Drops oficiales, AutoPick, Titulos y Cartas"`
