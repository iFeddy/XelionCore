# FASE 24: Barrido Rama por Rama (`.exe`/`.pdb`) — Combate, Fórmulas Internas, Skills, AbState, Mobs, Drops y Timers de Zona

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/Zone.md` y `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 24 (`docs/plan/FASE_24_BARRIDO_RAMAS_COMBATE_SKILLS_ABSTATE_MOBS_Y_DROPS_EXE.md`) siguiendo estrictamente `docs/plan/PLAN_AUDITORIA_PROFUNDA_S2C_VARIABLES_Y_RAMAS_EXE_100.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Desensamblá cada handler, subfunción matemática de combate y timer periódico (`OnTimer`) en `Zone.exe` con `python3 scripts/disasm_pdb_handler.py zone <símbolo>` e inspeccioná todas las ramas condicionales (`if/else`/`switch`), fórmulas, códigos de error y secuencias de paquetes S2C.
1) Auditá rama por rama en `Zone.exe` todos los handlers C2S y subfunciones internas de `NC_BAT` (9), `NC_SKILL` (18), `NC_SOULSTONE` (22), `NC_CT` (27) y `NC_COLLECT` (35): targeteo, auto-ataque, rango/ángulo, consumo de SP/HP/LP, cooldowns, máscaras AoE (`AreaBMP`), `AbState` (buffs/debuffs/inmunidades), `SoulStone`, muerte, `sp_ReviveReqProcess` (`REVIVESAME`/`REVIVEOTHER`), títulos y colección de cartas.
2) Desensamblá y alineá 1:1 contra `Zone.exe` las subfunciones matemáticas de combate (`hit/evade`, `crit`, `block`, daño físico/mágico, `ClassParam`, `DamageByAngle`) y los timers autónomos (`OnTimer`) del Game Loop de Zona: ticks periódicos DoT/HoT y expiración de `AbState`, regeneración pasiva de HP/SP/LP, máquina de estados `MobBehavior` / `MobRegen` / `MobRoam` / `MobAttackSequence`, propiedad/expiración de `Drops` en el piso y auto-guardado periódico (`CHARSAVE`).
3) Alineá en `pkg/server/zone/` todas las condiciones de rechazo, códigos de error (`SKILL_CAST_FAIL_ACK`, etc.), fórmulas de daño y orden exacto de broadcasts S2C (`SWING_CMD`, `DAMAGE_CMD`, `HPCHANGE_CMD`, `SPCHANGE_CMD`, `ABSTATE_SET_CMD`/`RESET_CMD`, `REALLYKILL_CMD`, `EXP_GAIN_CMD`).
4) Creá `pkg/audit/fase24_branch_combat_skill_mob_test.go` ejercitando las ramas de éxito y de fallo (sin SP, en cooldown, fuera de rango, target muerto, stun/abstate bloqueante, soulstone insuficiente, revivir en ciudad vs en el lugar), las fórmulas matemáticas de combate y los ticks de timers autónomos (`AbState`, regen, mobs y drops).
5) Verificá que `go test -p 1 -count=1 ./...` y los 4 scripts `verify_*_100.sh` pasen en `exit 0`, marcá Fase 24 como "Completado" en `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 24)

1. **Ramas de validación, daño, muerte y resurrección de `NC_BAT` y `NC_SKILL` + subfunciones matemáticas (`hit/evade/crit/block/damage/angle`) alineadas 1:1 con `Zone.exe`.**
2. **Timers autónomos (`OnTimer`) de `AbState` (ticks/expiración), regeneración HP/SP/LP, IA/aggro/regen de Mobs, expiración de Drops y `CHARSAVE` verificados contra `Zone.exe`.**
3. **Suite `pkg/audit/fase24_branch_combat_skill_mob_test.go` en `PASS` (`exit 0`).**

