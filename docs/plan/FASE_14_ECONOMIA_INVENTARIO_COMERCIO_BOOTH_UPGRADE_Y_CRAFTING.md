# Fase 14 — Economía In-Game Completa: `48/48` C2S `NC_ITEM`, `62` Métodos `iac_*`, Mercaderes NPC, `SafeTrade`, `Booth`, `Upgrade`, `Enchant` y Crafting

- **Nivel recomendado:** `High`
- **Documento maestro de referencia obligatoria:** [PLAN_ACCION_ZONESERVICE_GAMELOG.md](../PLAN_ACCION_ZONESERVICE_GAMELOG.md) (**Secciones 1-Bug #13, 4.1-Depts 2/4/8/12/19/26/38, 4.2 y 6.4**), [`docs/Servicios/Zone.md`](../Servicios/Zone.md) y [`AGENTS.md`](../../AGENTS.md).
- **Dependencias previas:** Fases 08 a 13 completadas sobre `main`.

---

## Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá únicamente la FASE 14 (`docs/plan/FASE_14_ECONOMIA_INVENTARIO_COMERCIO_BOOTH_UPGRADE_Y_CRAFTING.md`) leyendo obligatoriamente las Secciones 1 (Bug #13), 4.1 (Depts 2-`iac_FoodEat`, 4-`MINIMON`, 8-`PRODUCE`/`ACTIONBYITEM`/`ENDOFTRADE`, 12, 19, 26 y 38), 4.2 y 6.4 de `docs/PLAN_ACCION_ZONESERVICE_GAMELOG.md`, junto con `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: 100% de precios, tasas de éxito/rotura, recetas, sockets y cofres provienen de `data/NPCItemList/`, `data/World/Item/` y las tablas `.shn` oficiales, persistiendo de forma atómica con `CALL <SP>` en MySQL y auditando cada operación en `GameLog`.
1) Completá el 100% de los 48 handlers C2S del Departamento 12 (`NC_ITEM`) + los 62 métodos `ItemAttributeClass` (`iac_*`, incluyendo `0x0804`/`0x0805` `iac_FoodEat` y los 4 handlers `CHAR_USEITEM_MINIMON_*` en Dept 4): compra/venta/recompra en Mercaderes NPC (`NPCItemList/*.json`, `ItemShop.shn`), almacén personal/cuenta/gremio/academia (`Dept 38` `GUILD_ACADEMY_REWARD_STORAGE_OPEN/CLOSE`), pergaminos `MAPLINK_SCROLL_REQ`, `CLASS_CHANGE_REQ`, `CHAT_COLOR_CHANGE_REQ`, `TERMEXTEND_REQ`, `EQUIP_BELONGED_CANCEL_USE_REQ` y reparación de durabilidad (`WEAPONENDURE`, `SHIELDENDURE`, `MH_FURNITURE_ENDURE_KIT`).
2) Implementá el comercio seguro entre jugadores (`Dept 19` `NC_TRADE`: los 10 handlers C2S + `ACT_ENDOFTRADE_CMD` con transferencia atómica libre de duplicación) y el sistema completo de Tiendas Personales / Kioscos (`Dept 26` `NC_BOOTH`: los 8 handlers C2S incluyendo búsqueda categorizada `BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ` y `SEARCH_BOOTH_POSITION_REQ` contra `BoothSearchItemList.shn`).
3) Implementá Mejora (`ITEM_UPGRADE_REQ`, `ITEM_NEW_UPGRADE_REQ`, `BRACELET_UPGRADE_REQ`, `MOVER_UPGRADE_REQ` con `UpgradeInfo.shn`, `AccUpGradeInfo.shn`, `BraceletUpgrade.shn`), Gemas y Sockets (`ENCHANT_ADD_GEM`, `REMOVE_GEM`, `ADD_NEW_SOCKET` con `ItemSocket.shn`), `RANDOMOPTION_*` (`RandomOptionTable.json`), `DISMANTLE_REQ` (`ItemDismantle.shn`), `MIX_ITEM_REQ` (`ItemMix.shn`), Producción (`ACT_PRODUCE_CAST_REQ`/`ABORT` con `Produce.shn` y `RawInfo.shn`) y Cofres (`TreasureReward.json`), registrando `GameLog` Tipos `3, 4, 5, 6, 7, B, C, D`.
4) Creá tests transaccionales exhaustivos contra MySQL verificando cero duplicación de ítems/Cen y trazabilidad en `tGameLog`. Al terminar, marcá Fase 14 como "Completado" en `docs/plan/README.md` y hacé `git commit` limpio en `main`. No avances a la Fase 15.
```

---

## Alcance Exacto (Cero Omisiones)

1. **Los `48` Handlers C2S de `NC_ITEM`, `62` Métodos `iac_*`, Mercaderes y Almacenes (`pkg/server/zone/items_*.go`):**
   - **Mercaderes NPC:** `ITEM_BUY_REQ` (`0x3003`), `ITEM_SELL_REQ` (`0x3006`) e `ITEM_REPURCHASE_REQ` (`0x3059`) validando proximidad al NPC, catálogo real en `NPCItemList/<NPC>.json` / `ItemShop.shn`, moneda (Cen/Fama/Tokens) y capacidad de inventario (`GameLog` Tipo `3`).
   - **Almacenes:** `ITEM_DEPOSIT_REQ` (`0x301C`), `ITEM_WITHDRAW_REQ` (`0x301E`), `ITEM_OPENSTORAGEPAGE_REQ` (`0x3028`), `ITEM_ACCOUNT_STORAGE_OPEN_CMD` (`0x304C`), `ITEM_ACCOUNT_STORAGE_CLOSE_CMD` (`0x304D`), `ITEM_GUILD_STORAGE_WITHDRAW_REQ` (`0x3026`), `ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ` (`0x302F`), `DEPOSIT_REQ` (`0x3031`), `GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ` (`0x9831`) y `CLOSE_CMD` (`0x9833`) (`GameLog` Tipos `5` y `B`).
   - **Consumibles, Atributos `iac_*`, Minimon y Utilidades:** `iac_FoodEat` (`0x0804`, `0x0805`), `ITEM_USEABORT_CMD` (`0x3019`), `ACT_ACTIONBYITEM_REQ` (`0x2067`), `ITEM_MINIMON_EQUIP_REQ` (`0x3050`), `CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON/OFF_REQ` (`0x1124`/`0x1126`), `CHARGED_ITEM_ON/OFF_REQ` (`0x1128`/`0x112A`), `ITEM_WEAPONENDURE_CHARGE_REQ` (`0x3035`), `SHIELDENDURE_CHARGE_REQ` (`0x3064`), `MH_FURNITURE_ENDURE_KIT_USE_REQ` (`0x3033`), `EQUIP_BELONGED_CANCEL_USE_REQ` (`0x3048`), `CHAT_COLOR_CHANGE_REQ` (`0x3055`), `TERMEXTEND_REQ` (`0x3057`), `CLASS_CHANGE_REQ` (`0x305E`), `MAPLINK_SCROLL_REQ` (`0x3067`) y apertura de cofres `TreasureReward.json`.

2. **Comercio Seguro (`NC_TRADE` — `10` C2S) y Kioscos (`NC_BOOTH` — `8` C2S):**
   - **`SafeTrade` (`Dept 19`):** `TRADE_PROPOSE_REQ`, `PROPOSE_ASKNO_ACK`, `PROPOSEYES_ACK`, `PROPOSE_CANCEL_CMD`, `CANCEL_REQ`, `UPBOARD_REQ`, `DOWNBOARD_REQ`, `CENBOARDING_REQ`, `BOARDLOCK_REQ`, `DECIDE_REQ` + `ACT_ENDOFTRADE_CMD` con bloqueo mutuo de slots/Cen, validación anti-race y ejecución atómica en MySQL (`GameLog` Tipo `4`).
   - **`Booth` (`Dept 26`):** `BOOTH_OPEN_REQ`, `CLOSE_REQ`, `ENTRY_REQ`, `REFRESH_REQ`, `ITEMTRADE_REQ`, `INTERIORSTART_REQ`, `SEARCH_ITEM_LIST_CATEGORIZED_REQ` y `SEARCH_BOOTH_POSITION_REQ` contra `BoothSearchItemList.shn` (`GameLog` Tipo `C`).

3. **Mejora (`Upgrade`), Gemas (`Enchant`), Opciones Aleatorias y Artesanía (`Produce` / `Mix` / `Dismantle`):**
   - **Mejora y Gemas:** `ITEM_UPGRADE_REQ` (`0x3017`), `ITEM_NEW_UPGRADE_REQ` (`0x306D`), `ITEM_BRACELET_UPGRADE_REQ` (`0x306F`), `ITEM_MOVER_UPGRADE_REQ` (`0x3051`), `ITEM_ENCHANT_ADD_GEM_REQ` (`0x3040`), `REMOVE_GEM_REQ` (`0x3042`), `ADD_NEW_SOCKET_REQ` (`0x3044`), `RANDOMOPTION_CHANGE_REQ` (`0x3053`), `CHANGE_ACCEPT_REQ` (`0x3061`) y `RECOVER_COUNT_LIMIT_REQ` (`0x306B`) leyendo `UpgradeInfo.shn`, `AccUpGradeInfo.shn`, `BraceletUpgrade.shn`, `ItemSocket.shn` y `RandomOptionTable.json` (`GameLog` Tipo `7` + `PROTO_NC_LOG_GAME_ITEM_UPGRADE` / `CW_BREAKATZERO`).
   - **Producción, Mezcla y Desarme:** `ACT_PRODUCE_CAST_REQ` (`0x2035`), `ACT_PRODUCE_CASTABORT_CMD` (`0x2039`), `ITEM_DISMANTLE_REQ` (`0x302A`) e `ITEM_MIX_ITEM_REQ` (`0x3069`) contra `Produce.shn`, `RawInfo.shn`, `ItemDismantle.shn` e `ItemMix.shn` (`GameLog` Tipo `D`).

---

## Criterio de Cierre de la Fase 14

1. `go build ./...` -> `OK`.
2. `go test -p 1 -count=1 -v ./pkg/server/zone/...` -> `PASS` (incluyendo tests de compra/venta NPC, `SafeTrade`, `Booth`, `Upgrade`, `Enchant`, `Mix`, `Dismantle` y `Produce` + verificación en `tGameLog`).
3. `go test -p 1 -count=1 ./...` -> `PASS`.
4. Marcar Fase 14 como `Completado` en `docs/plan/README.md` y ejecutar:
   `git add -A && git commit -m "feat(zone): Fase 14 - 48/48 C2S NC_ITEM, 62 iac_*, SafeTrade, Booth, Upgrade, Enchant y Crafting con GameLog"`
