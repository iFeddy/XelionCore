# Fase 12 — Motor de Combate (`NC_BAT`), Habilidades Activas/Pasivas (`NC_SKILL`), `AbState`, `SoulStone` y Reglas de PvP/PK

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (**Secciones 0-Regla #5, 1-Bugs #10 y #12, 4.1-Depts 4/9/18/20 y 6.2**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fases 08 a 11 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 12 (`docs/plan/FASE_12_MOTOR_COMBATE_SKILLS_ABSTATE_SOULSTONE_Y_PVP.md`) leyendo obligatoriamente las Secciones 0 (Regla #5), 1 (Bugs #10 y #12), 4.1 (Depts 4, 9, 18 y 20) y 6.2 de `docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md`, junto con `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: 100% de las fórmulas de daño, crítico, ángulo, brecha de nivel, enfriamientos y estados provienen de `data/World/Combat/`, `data/AbState/`, `data/AreaBMP/` y las tablas `.shn` oficiales; CERO números mágicos inventados.
1) Implementá el motor de combate (`Dept 9`: los 9 handlers C2S `TARGETTING_REQ`, `UNTARGET_REQ`, `BASHSTART_CMD`, `BASHSTOP_CMD`, `ASSIST_REQ`, `SKILLBASH_OBJ_CAST_REQ`, `SKILLBASH_FLD_CAST_REQ`, `SKILLBASH_CASTABORT_REQ`, `ABSTATE_ERASE_REQ` + auxiliares desde montura) usando `Common.json`, `DamageByAngle.json`, `DamageLvGap*.json`, `SwingDamageFactor.shn`, `ActionRangeFactor.shn`, `ArmorMaterialFactor.shn`, LoS `.shbd` (`IsAttackBlocked`) y máscaras `pkg/data/areabmp`.
2) Implementá el subsistema completo de Habilidades (`Dept 18`: `SKILL_EMPOWALLOC_REQ`, `SKILL_PRODUCTFIELD_REQ`, `SKILL_UNLEARN_REQ`, `SKILL_REPLYREVIVE_CMD` + `ITEM_USE_ACTIVESKILL_REQ`, `ActiveSkill.shn`, `PassiveSkill.shn`, `SkillMechanics.json` para tótems/trampas/escudos y persistencia con `CALL p_Skill_*`).
3) Implementá la máquina de estados `AbState` (`AbState.shn`, `SubAbState.shn`, `data/AbState/*.json`, ticks DoT/HoT, stuns/miedos/escudos, `CHAR_DEPOLYMORPH_CMD`, `CHAR_CHARGEDBUFF_ERASE_REQ` y guardado en `tCharacterState`), `SoulStone` (`Dept 20` 4 C2S + bono de daño/curación), asignación/reducción de stats (`CHAR_STAT_INCPOINT_REQ` / `DECPOINT_REQ`), muerte/resurrección (`CHAR_REVIVE_REQ`, `CHAR_REBIRTH_CMD`/`REJECT_CMD`, `ITEM_REVIVEITEMUSE_CMD`), PvP/PK (`FieldLvCondition.shn`, arenas libres, `PKCount`/`PKYellowTime`) y auditoría `GameLog` Tipos `2`, `8` y `E`.
4) Creá tests exhaustivos de fórmulas de combate, casteo AoE, ticks `AbState`, muerte/revivir y PvP. Al terminar, marcá Fase 12 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 13.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Motor Oficial de Combate Físico y Mágico (`pkg/server/zone/combat.go` — `Dept 9` `9` C2S):**
   - Implementar `PROTO_NC_BAT_TARGETTING_REQ` (`0x2401`), `UNTARGET_REQ` (`0x2408`), `BASHSTART_CMD` (`0x242B`), `BASHSTOP_CMD` (`0x2432`), `ASSIST_REQ` (`0x243E`), `SKILLBASH_OBJ_CAST_REQ` (`0x2440`), `SKILLBASH_FLD_CAST_REQ` (`0x2441`), `SKILLBASH_CASTABORT_REQ` (`0x2444`) y `ABSTATE_ERASE_REQ` (`0x2454`), más los casteos desde montura (`MOVER_SKILLBASH_OBJ_CAST_REQ` / `FLD_CAST_REQ`).
   - Aplicar las fórmulas exactas de `data/World/Combat/`: `WC`/`MA` vs `AC`/`MR`, `TH` vs `TB` (hit/evasión), tasa y multiplicador crítico, ángulo atacante-objetivo (`DamageByAngle.json`), diferencia de nivel (`DamageLvGap*.json`), `SwingDamageFactor.shn`, `ArmorMaterialFactor.shn`, verificación de línea de visión en `.shbd` (`IsAttackBlocked`) y difusión AOI con `PROTO_NC_BAT_SWING_CMD` / `SWING_DAMAGE_CMD` / `SKILLBASH_HIT_OBJ_START_CMD` / `SKILLBASH_HIT_DAMAGE_CMD` / `HPCHANGE_CMD` / `SPCHANGE_CMD` / `LPCHANGE_CMD`.

2. **Habilidades Activas, Pasivas, Empoderamiento y Campos Mágicos (`pkg/server/zone/skills.go` — `Dept 18` `4` C2S):**
   - Implementar `PROTO_NC_SKILL_EMPOWALLOC_REQ` (`0x4811`), `SKILL_PRODUCTFIELD_REQ` (`0x4822`), `SKILL_UNLEARN_REQ` (`0x4824`), `SKILL_REPLYREVIVE_CMD` (`0x482A`), `ITEM_USE_ACTIVESKILL_REQ` (`0x304E`) y reinicios pagos (`SKILL_EMPOW_RESET_REQ` / `CHARGED_SKILLEMPOW_INITIALIZE_REQ`).
   - Validar nivel requerido, arma requerida, SP/LP, `CoolTime` por grupo (`ActiveSkillGroup.shn`), formas de área (`ActionRangeFactor.shn` y `pkg/data/areabmp`), efectos pasivos (`PassiveSkill.shn`) y entidades de campo (`SkillMechanics.json`: tótems, trampas, auras), persistiendo en MySQL mediante `CALL p_Skill_*` y auditando en `GameLog` Tipo `E` (`p_GameLog_Add_TypeE`).

3. **Estados Anormales (`AbState`), `SoulStone`, Distribución de Stats, Muerte/Resurrección y PvP/PK:**
   - **`AbState` (`pkg/server/zone/abstate_engine.go`):** motor de estados con prioridad/apilamiento, resistencias (`SubAbState.shn`), ticks periódicos DoT/HoT, crowd-control (stun, fear, root, sleep, silence), polimorfismo (`PROTO_NC_CHAR_DEPOLYMORPH_CMD` `0x10B0`), buffs pagos (`PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ` `0x112E`) y persistencia en `tCharacterState` (`CALL p_Char_SetStateWithRestTime`).
   - **`SoulStone` (`Dept 20` `4` C2S) y Stats Libres (`Dept 4`):** validar y completar `SOULSTONE_HP_BUY_REQ`, `SP_BUY_REQ`, `HP_USE_REQ`, `SP_USE_REQ` aplicando multiplicadores reales de `DamageBySoul` / `SingleData.shn`; implementar `CHAR_STAT_INCPOINT_REQ` (`0x105C`), `CHAR_STAT_DECPOINT_REQ` (`0x1062`) y `CHARGED_STAT_INITIALIZE_REQ` contra `FreeStatTable.json` y `CALL p_Char_Set*` + `GameLog` Tipo `2`.
   - **Muerte, Resurrección y PvP/PK:** penalización/protección de EXP al morir, `CHAR_REVIVE_REQ` (`0x104E`), `CHAR_REBIRTH_CMD` (`0x1102`), `CHAR_REBIRTH_REJECT_CMD` (`0x1103`), `ITEM_REVIVEITEMUSE_CMD` (`0x3039`), validación de zonas PvP/PK (`FieldList.json`, `FieldLvCondition.shn`), conteo `PKCount`/`PKYellowTime` y registro en `GameLog` Tipo `8` (`p_GameLog_Add_Type8` / `PROTO_NC_LOG_GAME_PKED`).

---

## Criterio de Cierre de la Fase 12

1. `go build ./...` -> `OK`.
2. `go test -p 1 -count=1 -v ./pkg/server/zone/...` -> `PASS` (incluyendo tests de combate, ángulos, AoE `AreaBMP`, `AbState`, `SoulStone`, stats y resurrección).
3. `go test -p 1 -count=1 ./...` -> `PASS`.
4. Marcar Fase 12 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(zone): Fase 12 - motor oficial de combate, ActiveSkill/PassiveSkill, AbState, SoulStone y PvP/PK"`
