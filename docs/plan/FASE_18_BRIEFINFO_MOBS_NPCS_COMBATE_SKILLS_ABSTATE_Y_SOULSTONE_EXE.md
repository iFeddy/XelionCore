# FASE 18: Saneamiento Binario 1:1 — `NC_BRIEFINFO` (Mobs, NPCs, Players, Gates, Drops), Combate (`NC_BAT`), Skills (`NC_SKILL`), `AbState` y `SoulStone`

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/Zone.md` y `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 18 completa (`docs/plan/FASE_18_BRIEFINFO_MOBS_NPCS_COMBATE_SKILLS_ABSTATE_Y_SOULSTONE_EXE.md`) siguiendo estrictamente `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Antes de modificar cada handler o paquete S2C, ejecutá `python3 scripts/disasm_pdb_handler.py zone <símbolo>` sobre `../AstraServer/Server/Zone00/Zone.exe` y verificá cada offset y `sizeof` contra `docs/Servicios/Zone.md`. Prohibido inventar bytes, omitir campos o dejar stubs.
1) Desensamblá en `Zone.exe` la construcción de todos los paquetes del departamento `NC_BRIEFINFO` (Dept 7): `PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD`, `PROTO_NC_BRIEFINFO_CHARACTER_CMD`, `PROTO_NC_BRIEFINFO_REGENMOB_CMD`, `PROTO_NC_BRIEFINFO_MOB_CMD`, `PROTO_NC_BRIEFINFO_DROPEDITEM_CMD`, `PROTO_NC_BRIEFINFO_ITEMONFIELD_CMD`, `PROTO_NC_BRIEFINFO_MAGICFIELDINFO_CMD`, `PROTO_NC_BRIEFINFO_MOVER_CMD`, `PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD`, `PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD`, `PROTO_NC_BRIEFINFO_UNEQUIP_CMD`, `PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD` y `PROTO_NC_BRIEFINFO_DOOR_CMD`. Corregí su serialización en `pkg/server/zone/` para que cada entry y paquete coincida al byte con `Zone.exe` y `Zone.md` (modo/estado de NPCs, portales/gates, monstruos vivos, velocidades, coordenadas y animaciones).
2) Desensamblá y alineá 1:1 contra `Zone.exe` todos los handlers y respuestas S2C de `NC_BAT` (Dept 9): `sp_NC_BAT_TARGETTING_REQ` (`PROTO_NC_BAT_TARGETINFO_CMD` con HP/MaxHP/SP/MaxSP/LP/Level reales del mob/jugador), `sp_NC_BAT_HIT_REQ`, `sp_NC_BAT_UNTARGET_REQ`, `sp_NC_BAT_BASHSTART_CMD`, `sp_NC_BAT_BASHSTOP_CMD`, `PROTO_NC_BAT_SWING_CMD`, `PROTO_NC_BAT_SWINGDAMAGE_CMD`, `PROTO_NC_BAT_HPCHANGE_CMD`, `PROTO_NC_BAT_SPCHANGE_CMD`, `PROTO_NC_BAT_LPCHANGE_CMD`, `PROTO_NC_BAT_EXPGAIN_CMD`, `PROTO_NC_BAT_LEVELUP_CMD`, `PROTO_NC_BAT_ REALLYKILL_CMD` y `PROTO_NC_BAT_ABSTATESET_CMD`/`RESET_CMD`.
3) Desensamblá y alineá 1:1 contra `Zone.exe` todos los handlers y respuestas S2C de `NC_SKILL` (Dept 18) y `NC_SOULSTONE` (Dept 22): `sp_NC_SKILL_SKILLTEACH_REQ`, `sp_NC_SKILL_CAST_REQ`, `sp_NC_SKILL_STOP_REQ`, `sp_NC_SKILL_ITEMACTIONCOOLTIME_REQ`, `sp_NC_SKILL_REVIVE_SAME_SC_REQ`, `sp_NC_SKILL_WARP_CANCLE_REQ`, `sp_NC_SOULSTONE_HP_USE_REQ`, `sp_NC_SOULSTONE_SP_USE_REQ`, `sp_NC_SOULSTONE_BUY_REQ` y `sp_NC_SOULSTONE_USE_REQ`.
4) Creá `pkg/audit/fase18_binary_briefinfo_combat_skills_test.go`, verificá que `go build ./... && go test -p 1 -count=1 ./...` pase en `exit 0`, recompilá `bin/LoginService`, `bin/WorldService` y `bin/ZoneService`, marcá Fase 18 como "Completado" en `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 18)

1. **Paridad Exacta de `NC_BRIEFINFO` (`Zone.md` vs `Zone.exe`):**
   - Verificar el `sizeof` exacto de cada elemento dentro de `PROTO_NC_BRIEFINFO_MOB_CMD` (`PROTO_NC_BRIEFINFO_REGENMOB_CMD`) y `PROTO_NC_BRIEFINFO_CHARACTER_CMD`.
   - Confirmar cómo `Zone.exe` diferencia gates/portales (`Gate*`), NPCs y monstruos dentro de `BRIEFINFO_MOB_CMD` (flags, `mode`, velocidades, equipo/aspecto y `abstate`).
2. **Paridad Exacta de `NC_BAT` (`Target`, `Swing`, `Damage`, `HP/SP Change`):**
   - Desensamblar `sp_NC_BAT_TARGETTING_REQ` y verificar los `30` bytes de payload (`32` con OpCode) de `PROTO_NC_BAT_TARGETINFO_CMD` y el orden exacto de `HP`, `MaxHP`, `SP`, `MaxSP`, `LP`, `MaxLP`, `Level`.
   - Verificar `PROTO_NC_BAT_HPCHANGE_CMD` y `PROTO_NC_BAT_SPCHANGE_CMD` para que la barra de HP/SP del jugador nunca se corrompa ni muestre valores desfasados.
3. **Paridad Exacta de `NC_SKILL` y `NC_SOULSTONE`:**
   - Desensamblar `sp_NC_SKILL_CAST_REQ` y `sp_NC_SOULSTONE_HP_USE_REQ` / `SP_USE_REQ` / `BUY_REQ` para confirmar que las piedras de HP/SP y las skills lean y respondan los structs oficiales exactos.

