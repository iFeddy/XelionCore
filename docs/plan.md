# Plan de Corrección 1:1 del Sistema de Skills (`Zone.exe` -> `XelionCore`)

**Rama:** `fix/skills`
**Fuente de la verdad:** `re/XelZone.exe` + `re/pdo/Zone.pdb`, `data/ActiveSkill.shn`, `data/ActiveSkillInfoServer.shn`, `data/World/Item/ItemUseFunction.json`.
**Archivos objetivo:** `pkg/server/zone/skill_catalog.go` (nuevo), `pkg/server/zone/combat.go`, `pkg/server/zone/skills.go`, `pkg/server/zone/character.go`, `pkg/server/zone/data_store.go`, `pkg/server/zone/gameloop.go`.

---

## Tabla Maestra de Diferencias (`Zone.exe` vs Go)

### 1. `PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ` (`0x2440`) & `PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ` (`0x2441`)
- `0x0057FA50` (`sp_NC_BAT_SKILLBASH_OBJ_CAST_REQ`), `0x0057F400` (`sp_NC_BAT_SKILLBASH_FLD_CAST_REQ`), `0x0055B920` (`csl_SPCheck`), `0x0046A420` (`csl_CoolTimeCheck`), `0x0055B6A0` (`smo_ply_SkillCoolTime`), `0x0057BEF0` (`smo_SkillBash_Casting`).
- **Estado del jugador:** `0x0057FAEA`, `0x0057F4A2` -> bitmask `(1 << state) & 0x0004` (solo estado `2` en EXE; en Go `State == 1` es Normal/InGame -> mapear consistentemente), error `0x0FC0`.
- **Mapa `canUseSkill`:** `0x0057FB0E`, `0x0057F515` -> `[field+0x9F] != 0`, error `0x0FC3`.
- **Bloqueos / Silence / AbStates:**
  - Ya casteando `[esi+0x1E78] & 0x04` -> `0x0FCC` (`0x0057FB5C`).
  - Silence `smo_SkillBlastOption() == 2` -> `0x0FC7` (`0x0058001E`, `0x0057F7BF`) — Go emitía `0x0FCC`.
  - AbStates `0x123` (291) / `0x1F3` (499) en caster -> `0x0FC1` (`0x0057FC6A`, `0x0057F4D8`); en target `0x1F3` -> abort silencioso (`0x0057FCC1`), `0x123` -> `0x0FC2` (`0x0057FDB6`).
- **LoS (`so_FindBlockLocate`):** `0x0057FD52` (OBJ) y `0x0057F605` (FLD) -> `0x0FD1` con `0x0049DF00` (attack block); FLD con teleport (`[edi+0x178] != 0`) también chequea `0x0049DF70` (move block, `0x0057F6F8`).
- **Existencia y aprendizaje:** `0x0057FD92`, `0x0057F648` -> skill desconocida o pasiva (`[+0xC7] == 2`) -> `0x0FC6`; no aprendida en `CharaterSkillList` -> `0x0FC8` (`0x0046A420`).
- **Costos SP/HP/LP (`csl_SPCheck` `0x0055B920`):**
  - `HPCost = baseHP + (MaxHP * hpRate) / 1000` -> `CurHP <= HPCost` => `0x0FD5`.
  - `LPCost = baseLP` -> `CurLP <= LPCost` (OBJ) / `< LPCost` (FLD) => `0x0FDC`.
  - `SPCost = baseSP + (MaxSP * spRate) / 1000 - empowerSP[nibble1]` -> `CurSP < SPCost` => `0x0FC9`.
- **Cooldown (`csl_CoolTimeCheck` `0x0046A420`):** `DlyTime` de `ActiveSkill.shn` menos empower cooldown (nibble 3), error `0x0FC8`.
- **Rango (`0x0058016E`, `0x0057F8EE`):** `baseRange = objSize + skill.Range`, `maxRange = min((baseRange * 10) >> 3, baseRange + 60)`, error `0x0FCA`.
- **`castIndex` y AOI (`0x0057BEF0`, `0x00437110`):** `castIndex` dinámico (`sbec_GetIndex`), broadcast por AOI (`so_AllOfRange`) para `0x244F`, `0x2451`, `0x2452`.

### 2. Daño, Heal, AoE, AbStates y `CASTABORT` (`0x2444`)
- `0x00506B40` (`RulesOfEngagementPhisycalSkill::roe_AttackPower`): `minAtk = minWC + (minWC * MinWCRate) / 1000 + MinWCAdd + empDmg[nibble0]` (y análogo `maxAtk`), roll uniforme `[minAtk, maxAtk]`.
- `0x00506ED0` (`RulesOfEngagementMagicalSkill::roe_AttackPower`): análogo con `MinMA`/`MaxMA`, `MinMARate`/`MaxMARate`, `MinMAAdd`/`MaxMAAdd`, `empDmg[nibble0]`.
- `0x00505C00` (`roe_CalcDamage`): AbState `0x16D` -> daño `1`, crítico con `well512_GetRandom(1000) < critRate`, elemental `*(1024 + diff) / 1024`.
- `0x005817CE` (`smo_SkillBlast` heal): `heal = baseHeal + empDmg[nibble0]`, marca `isHeal=1` en `0x2452`.
- `0x00583BC0` / `0x00584160` (`sbe_BlastObject` / `sbe_BlastLocate`): radio **`Area / 2`**, sector cónico si `SkillDegree != 360`, ordenamiento por distancia y tope `TargetNumber`.
- `0x005824A7` (`smo_SkillBlast` AbStates): hasta 4 AbStates (`StaNameA..D`), roll `rand(1000) < StaSucRate`, duración + empower nibble 2 (`nT2`).
- `0x0057C480` / `0x00579F10` (`CASTABORT_REQ` / `CASTABORT_CMD`): estado `6` -> `0x20E9`, estado inválido -> `0x20EA`, éxito -> `0x20E8` + `0x2438` por **AOI** si casteaba (o solo a sí mismo si no).

### 3. Mover Skills (`0xCC0E` / `0xCC0F`)
- `0x00580800` (`so_smv_SkillBash_Object`) y `0x00580CF0` (`so_smv_SkillBash_Location`): bitmask `0x0006`, sin SP check, cooldown propio en `MoverSkillList`, rango `min(1.25x, +60)`, LoS en OBJ y FLD, target inválido -> `0x0FD6`.

### 4. `EMPOWALLOC` (`0x4811`), `PRODUCTFIELD` (`0x4822`), `UNLEARN` (`0x4824`)
- `0x0044F150` / `0x00448570`: bitmask `0x0076` (`0x0882`), `minusMask` requiere ítem reset (`0x0886`), tope global de puntos (`0x0883`), tope `0..5` por track + underflow + skill aprendida (`0x0884`), `0x4810` con puntos restantes reales.
- `0x0046C820`: bitmask `0x0016` (`0x0B04`), whitelist de 5 profesiones de producción (`0x0B12`), duplicado (`0x0B13`), máx. 2 profesiones (`0x0B14`), persistencia real.
- `0x0057AD80`: `csl_IsLearn` (`0x1503`), validación de `slot` (`0x0702`), lock (`0x070A`), ítem Skill Erase Scroll (`0x0702`), `lot > 0` (`0x0703`), consumo del ítem.

### 5. `REPLYREVIVE` (`0x482A`) y `ITEM_USE_ACTIVESKILL` (`0x304E`)
- `0x00418A30` / `0x0057D290`: mapa + flag `PendingRevive` (`0x1506`), HP = `(MaxHP * 1000) / 1000` (100%), `0x481F` por **AOI**.
- `0x00540B90`: ejecuta skill activa de ítem consumible (NO aprende skill), errores `0x3129..0x312F`, éxito `0x3128`.

---

## Fases de Implementación

### Fase 1 — Catálogo oficial de skills (`SkillCatalog`) y estado de skills en `ZonePlayer`
- **Estado:** **Completada**
- **Alcance:**
  - Crear `pkg/server/zone/skill_catalog.go` con `SkillDef` y `BuildSkillCatalog(ds *DataStore)` leyendo `ActiveSkill.shn` + `ActiveSkillInfoServer.shn` (reutilizando `ds.ActiveSkillSHN` y `ds.GetSHN("ActiveSkillInfoServer")`).
  - Integrar `SkillCatalog` en `DataStore`.
  - Agregar estado de skills en `ZonePlayer`: `LearnedSkills map[uint16]*database.CharacterSkill`, `NextCastIndex uint16`, `CastingTargetHandle uint16`, `CastingFieldX/Y int32`, `CastingEndTime time.Time`, `PendingRevive bool`, `PendingReviveRate uint16`, `ToggleSkills map[uint16]bool`, `MoverSkillCooldowns map[uint16]time.Time`, más helpers `HasSkill`, `GetSkillEmpower`, `TotalEmpowerUsed`, `MaxEmpowerPoints`, `NextSkillCastIndex`.
  - Poblar `LearnedSkills` al entrar al mapa (`ZonePlayer` login/map join).
  - Tests unitarios para `SkillCatalog` y helpers de `ZonePlayer`.

### Fase 2 — Validaciones de casteo (`0x2440` / `0x2441`), costos SP/HP/LP, cooldown y rango oficial
- **Estado:** **Completada**
- **Alcance:**
  - Actualizar `handleSkillBashObjCastReq` (`0x0057FA50`) y `handleSkillBashFldCastReq` (`0x0057F400`):
    - Estado válido (`0x0FC0`), mapa permite skills (`0x0FC3`), ya casteando (`0x0FCC`), Silence (`0x0FC7`), AbStates `0x123`/`0x1F3` (`0x0FC1`/`0x0FC2`).
    - Lookup en `SkillCatalog` (`0x0FC6`) y `HasSkill` (`0x0FC8`), LoS en OBJ y FLD (`0x0FD1`) + move-block en FLD teleport.
    - Fórmula oficial de costos HP (`0x0FD5`), LP (`0x0FDC`), SP con empower track 1 (`0x0FC9`).
    - Cooldown oficial `DlyTime` menos empower track 3 (`0x0FC8`).
    - Rango oficial `min((baseRange*10)>>3, baseRange+60)` (`0x0FCA`).
    - `castIndex` dinámico y broadcast por AOI (`grid.BroadcastAOI`).
  - Tests unitarios de todas las ramas de rechazo y fórmulas.

### Fase 3 — Motor de resolución (`smo_SkillBlast`, `RulesOfEngagement`, Heal, AoE, AbStates, CastAbort)
- **Estado:** **Completada**
- **Alcance:**
  - Implementar cálculo de daño de skill físico (`0x00506B40`) y mágico (`0x00506ED0`) con empower track 0 (`nT0`), hit/crit rolls y AbState `0x16D` / `0xC9` / `0xCA`.
  - Implementar skills de curación (`0x005817CE`) con flag `isHeal=1` en `0x2452`.
  - Implementar selección AoE oficial (`Area / 2`, sector cónico por `SkillDegree`, ordenamiento por distancia, tope `TargetNumber`) tanto en OBJ AoE como en FLD AoE.
  - Aplicar AbStates `StaNameA..D` con `StaSucRateA..D`, `StaStrengthA..D` y duración + empower track 2 (`nT2`).
  - Actualizar `handleSkillBashCastAbortReq` (`0x0057C480` + `0x00579F10`) con broadcast AOI de `0x2438` y respuesta propia si no casteaba.
  - Tests unitarios de daño físico/mágico, heal, AoE cono/radio/cap, AbStates y abort.

### Fase 4 — Mover Skills (`0xCC0E` / `0xCC0F`)
- **Estado:** **Completada**
- **Alcance:**
  - Implementar `handleMoverSkillBashObjCastReq` (`0x00580800`) y `handleMoverSkillBashFldCastReq` (`0x00580CF0`) con bitmask `0x0006`, cooldown en `MoverSkillCooldowns`, sin costo SP del jugador, LoS en OBJ/FLD, rango oficial y error `0x0FD6`.
  - Tests unitarios dedicados.

### Fase 5 — `EMPOWALLOC` (`0x4811`), `PRODUCTFIELD` (`0x4822`) y `UNLEARN` (`0x4824`)
- **Estado:** Pendiente
- **Alcance:**
  - `handleSkillEmpowAllocReq`: estado (`0x0882`), `minusMask` + ítem reset (`0x0886`), tope global (`0x0883`), tope `0..5` por track + underflow + skill aprendida (`0x0884`), `0x4810` con puntos restantes reales.
  - `handleSkillProductFieldReq`: estado (`0x0B04`), whitelist de 5 skills de producción (`0x0B12`), duplicado (`0x0B13`), máx. 2 (`0x0B14`), aprendizaje y persistencia.
  - `handleSkillUnlearnReq`: `HasSkill` (`0x1503`), validación de slot (`0x0702`), lock (`0x070A`), ítem Skill Erase Scroll (`0x0702`), `lot > 0` (`0x0703`), consumo del ítem y baja en memoria + DB.
  - Tests unitarios cubriendo todos los códigos de error y éxito.

### Fase 6 — `REPLYREVIVE` (`0x482A`) y `ITEM_USE_ACTIVESKILL` (`0x304E`)
- **Estado:** Pendiente
- **Alcance:**
  - `handleSkillReplyReviveCmd`: mapa + `PendingRevive` (`0x1506`), HP = `(MaxHP * rate) / 1000` (100%), `0x481F` por AOI.
  - `handleItemUseActiveSkillReq`: ejecutar skill activa asociada al ítem consumible con validaciones `0x3129..0x312F` y éxito `0x3128`.
  - Tests unitarios y verificación completa de toda la suite `go test ./pkg/server/zone/...`.
