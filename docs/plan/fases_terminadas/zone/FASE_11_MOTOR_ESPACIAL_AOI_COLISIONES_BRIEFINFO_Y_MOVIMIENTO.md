# Fase 11 — Game Loop Determinístico, Motor Espacial AOI, Colisiones `.shbd`/`.sbi`, `BRIEFINFO`, Movimiento, `Tent`, `Gather` y Monturas (`MOVER`)

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (**Secciones 0-Reglas #4 y #5, 1-Bugs #8, #9, #10 y #11, 4.1-Depts 2/6/7/8/51 y 6.1**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fases 08, 09 y 10 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 11 (`docs/plan/FASE_11_MOTOR_ESPACIAL_AOI_COLISIONES_BRIEFINFO_Y_MOVIMIENTO.md`) leyendo obligatoriamente las Secciones 0 (Reglas #4 y #5), 1 (Bugs #8, #9, #10 y #11), 4.1 (Depts 2, 6, 7, 8 y 51) y 6.1 de `docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md`, junto con `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: CERO fallbacks `"RouN"`, CERO fórmulas de stats inventadas (`* 3`, `NextExp = 280`, `JobToClassName` truncado a 6 clases) y CERO broadcasts globales O(N^2) donde corresponde AOI.
1) Implementá en `pkg/server/zone/` el Game Loop determinístico por mapa (`100ms` / `1s` / `30s`) y la grilla espacial por sectores AOI (`SpatialGrid`), reemplazando la difusión global por suscripción incremental de celdas y cubriendo la familia completa `PROTO_NC_BRIEFINFO_*` (`0x1C01` `BRIEFINFO_INFORM_CMD`, altas/bajas `BRIEFINFODELETE_CMD`, personajes, NPCs/gates, drops, tiendas `Camp`, mascotas, muebles y campos mágicos) más `HIDE_EXCEPT_ME_ON_CMD` / `OFF_CMD` y `MISC_SEED_REQ`.
2) Corregí `pkg/server/zone/character.go` y `map.go` para soportar las 27 clases oficiales (`ClassName.shn` + `data/World/ClassParam/*.json`), calcular `MaxHP`/`MaxSP`/`MaxLP`, `PrevExp`/`NextExp` y bonos de stats desde `FreeStatTable.json`, `Common.json` y `ExpTable.json`, serializar los upgrades reales en `PROTO_EQUIPMENT` y erradicar todo `"RouN"` o tamaño ficticio `256x256`.
3) Validá todos los paquetes de movimiento, salto, modo, portales/alas, descanso en tienda (`PITCHTENT_REQ` / `FOLDTENT_REQ`), recolección en campo (`GATHERSTART_REQ` / `GATHERCANCEL_CMD` / `GATHERCOMPLETE_REQ` con `CollectInfo.shn`) y monturas (`MOVER_RIDE_ON_REQ`, `MOVER_RIDE_OFF_REQ`, `ACT_RIDE_FEEDING_REQ` con `MoverMain.shn` / `MoverSub.shn`) contra `pkg/data/blockinfo` (`.shbd` + `.sbi`) y registra entrada/salida de mapa en `GameLog` Tipo `1`.
4) Agregá tests unitarios e integrales de AOI, colisiones `.shbd`/`.sbi`, 27 clases, movimiento, `Tent`, `Gather` y `Mover`. Al terminar, marcá Fase 11 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 12.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Game Loop Determinístico y Particionamiento Espacial AOI (`pkg/server/zone/aoi.go`, `gameloop.go`, `map.go`):**
   - Incorporar en `ZoneMap` la grilla espacial por sectores (`SpatialGrid`) para jugadores, monstruos, NPCs, portales, drops, tiendas de campaña (`Camp`), mascotas, kioscos (`Booth`) y campos mágicos.
   - Implementar actualización diferencial de visibilidad al entrar al mapa, moverse de celda o teletransportarse: envío de `PROTO_NC_BRIEFINFO_*` para nuevos sectores visibles y `PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD` para sectores que salen de rango, más manejo real de `PROTO_NC_BRIEFINFO_INFORM_CMD` (`0x1C01`), `PROTO_NC_MISC_HIDE_EXCEPT_ME_ON_CMD` (`0x0817`), `PROTO_NC_MISC_HIDE_EXCEPT_ME_OFF_CMD` (`0x0818`) y `PROTO_NC_MISC_SEED_REQ` (`0x0806`).
   - Ejecutar ticks periódicos por mapa (`100 ms` movimiento/interpolación, `1 s` regeneración natural y en descanso desde `Common.json`, `30 s` autoguardado periódico).

2. **Las `27` Clases Oficiales y Cálculo Fiel de Parámetros (`pkg/server/zone/character.go` — `Bugs #9, #10 y #11`):**
   - Reemplazar `JobToClassName` de 6 casos por resolución exacta de las **27 clases** desde `ClassName.shn` y `data/World/ClassParam/<Class>.json`.
   - Calcular `MaxHP`, `MaxSP`, `MaxLP`, `PrevExp`, `NextExp`, `MaxStonesHP`, `MaxStonesSP` y stats derivados estrictamente desde `ClassParam`, `FreeStatTable.json`, `Common.json`, `ExpTable.json` e ítems equipados (incluyendo `upgrade` visual de 3 bytes en `PROTO_EQUIPMENT` y campos reales `RedistributePoint`, `PKCount`, `PKYellowTime`, `PrisonMin`, `AdminLevel` en `PROTO_NC_CHAR_BASE_CMD`).
   - Eliminar los fallbacks `normID = "RouN"`, `shortName = "RouN"` y `256x256` en `map.go` y `character.go`, obteniendo dimensiones exactas de `FieldList.json` y cabecera `.shbd`.

3. **Movimiento Validado contra `.shbd`/`.sbi`, Links de Mapa, `Tent`, `Gather` y `Mover`:**
   - **Movimiento y Acciones Base (`Dept 8` y `Dept 6`):** `WALK_REQ`, `RUN_REQ`, `MOVEWALK_CMD`, `MOVERUN_CMD`, `STOP_REQ`, `CHANGEMODE_REQ`, `JUMP_CMD`, `SHOUT_CMD`, `ROAR_REQ`, `EMOTICON_CMD`, `EMOTICONSTOP_CMD` validando colisión `.shbd` + puertas `.sbi` y cota de velocidad anti-speedhack; `MAP_WING_SAVE_REQ` (`0x1814`), `MAP_WING_FLY_REQ` (`0x1816`), `MAP_TOWNPORTAL_REQ` (`0x181A`), `MAP_TONORMALCOORD_CMD` (`0x181C`), `MAP_MULTY_LINK_SELECT_REQ` (`0x181F`), `CHAR_SAVE_LINK_REQ` (`0x10B9`) y `ACT_AUTO_WAY_FINDING_USE_GATE_REQ` (`0x2070`), registrando `GameLog` Tipo `1` (`p_GameLog_Add_Type1`).
   - **Tiendas de Descanso (`Tent`):** `ACT_PITCHTENT_REQ` (`0x2027`) y `ACT_FOLDTENT_REQ` (`0x202A`) con difusión `PROTO_NC_BRIEFINFO_CAMP_CMD` y tasas oficiales de recuperación HP/SP.
   - **Recolección (`Gather`):** `ACT_GATHERSTART_REQ` (`0x202D`), `ACT_GATHERCANCEL_CMD` (`0x2030`) y `ACT_GATHERCOMPLETE_REQ` (`0x2032`) contra nodos y `CollectInfo.shn` / `RawInfo.shn`.
   - **Monturas (`Dept 51` + `Dept 8`):** `MOVER_RIDE_ON_REQ` (`0xCC01`), `MOVER_RIDE_OFF_REQ` (`0xCC02`) y `ACT_RIDE_FEEDING_REQ` (`0x2044`) leyendo velocidades, hambre y ranura desde `MoverMain.shn` y `MoverSub.shn` + `GameLog` Tipo `G`.

---

## Criterio de Cierre de la Fase 11

1. `go build ./...` -> `OK`.
2. `go test -p 1 -count=1 -v ./pkg/server/zone/...` -> `PASS` (incluyendo tests nuevos de AOI, bloqueo `.shbd`/`.sbi`, 27 clases, `Tent`, `Gather` y `Mover`).
3. `go test -p 1 -count=1 ./...` -> `PASS`.
4. Marcar Fase 11 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(zone): Fase 11 - GameLoop, grilla AOI, colisiones shbd/sbi, 27 clases, movimiento, Tent, Gather y Mover"`
