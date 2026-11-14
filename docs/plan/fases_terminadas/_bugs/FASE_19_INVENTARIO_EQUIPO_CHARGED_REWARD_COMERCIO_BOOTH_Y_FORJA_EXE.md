# FASE 19: Saneamiento Binario 1:1 — Inventario (`NC_ITEM`), Equipamiento, Charged/Reward Inventory, Tiendas NPC, SafeTrade, Booth (`NC_BOOTH`) y Forja

> **Nivel de Modelo Recomendado:** **`High`**
> **Tiempo / Alcance:** 1 sesión de chat (ejecutable con `/goal`).
> **Fuente de la Verdad Obligatoria:** `../AstraServer/Server/Zone00/Zone.exe` + `pdo/Zone.pdb` (vía `scripts/disasm_pdb_handler.py`), `docs/Servicios/Zone.md` y `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`.

---

## 1. Prompt Listo para Copiar y Pegar

```text
/goal Ejecutá la FASE 19 completa (`docs/plan/FASE_19_INVENTARIO_EQUIPO_CHARGED_REWARD_COMERCIO_BOOTH_Y_FORJA_EXE.md`) siguiendo estrictamente `docs/PLAN_SANEAMIENTO_BINARIO_EXE_PDB_100.md`, `docs/Servicios/Zone.md` y `AGENTS.md`.
Regla Mandatoria: Antes de modificar cada handler o paquete S2C, ejecutá `python3 scripts/disasm_pdb_handler.py zone <símbolo>` sobre `../AstraServer/Server/Zone00/Zone.exe` y verificá cada offset y `sizeof` contra `docs/Servicios/Zone.md`. Prohibido inventar bytes, omitir campos o dejar stubs.
1) Desensamblá en `Zone.exe` cómo se serializa `STRUCT_ITEM` / `SHINEDATA_ITEM` dentro de `PROTO_NC_CHAR_CLIENT_ITEM_CMD` (`0x1047`), `PROTO_NC_ITEM_CELLCHANGE_CMD` (`0x3001`) y `PROTO_NC_ITEM_EQUIPCHANGE_CMD` (`0x3002`), y cómo operan `sp_NC_ITEM_RELOC_REQ`, `sp_NC_ITEM_SPLIT_REQ`, `sp_NC_ITEM_EQUIP_REQ`, `sp_NC_ITEM_UNEQUIP_REQ`, `sp_NC_ITEM_DROP_REQ`, `sp_NC_ITEM_PICK_REQ`, `sp_NC_ITEM_DISMANTLE_REQ` y `sp_NC_ITEM_AUTO_ARRANGE_INVEN_REQ`. Corregí cualquier diferencia de tamaño/offset y recalculá stats de personaje al equipar/desequipar.
2) Desensamblá en `Zone.exe` el flujo completo de Charged Inventory y Reward Inventory: `sp_NC_ITEM_CHARGEDINVENOPEN_REQ` (`0x3020`), `sp_NC_ITEM_CHARGED_WITHDRAW_REQ` (`0x3022`), `sp_NC_ITEM_REWARDINVENOPEN_REQ` (`0x302C`) y `sp_NC_ITEM_REWARDINVENGET_REQ` (`0x302F` / `iac_NC_ITEM_REWARDINVENOPEN_ACK`). Corregí `handleChargedInvenOpen`, `handleChargedWithdraw`, `handleRewardInvenOpen` y el retiro de ítems para que al sacar un ítem del Charged/Reward Inventory se inserte en una celda libre real de `tItem`, se envíe `PROTO_NC_ITEM_CELLCHANGE_CMD` al cliente con el struct exacto y el ítem quede 100% usable/movible en el inventario.
3) Desensamblá y alineá 1:1 contra `Zone.exe` los restantes handlers de `NC_ITEM` (48/48 C2S), `NC_ITEMDB` (13), `NC_BOOTH` (16) y `NC_CHARSAVE` (31): compra/venta NPC (`BUY_REQ`, `SELL_REQ`, `REPURCHASE_REQ`), uso de consumibles/pergaminos (`USE_REQ`, `USECOMPLETE_CMD`), almacén (`STORAGEOPEN_REQ`, `DEPOSIT`, `WITHDRAW`), SafeTrade (`TRADEREQ`..`TRADECOMPLETE`), Upgrade/Enchant/RandomOption/Mix/Produce y Kioscos `BOOTH`.
4) Creá `pkg/audit/fase19_binary_items_charged_reward_economy_test.go`, verificá que `go build ./... && go test -p 1 -count=1 ./...` pase en `exit 0`, recompilá `bin/LoginService`, `bin/WorldService` y `bin/ZoneService`, marcá Fase 19 como "Completado" en `docs/plan/README.md` y hacé `git commit` en `main`.
```

---

## 2. Checklist de Verificación Obligatoria (Fase 19)

1. **Formato Binario de Ítems (`PROTO_NC_CHAR_CLIENT_ITEM_CMD`, `CELLCHANGE_CMD`, `EQUIPCHANGE_CMD`):**
   - Verificar con `Zone.md` y `Zone.exe` el tamaño exacto de un slot ocupado vs. un slot vacío en `PROTO_NC_ITEM_CELLCHANGE_CMD` y `PROTO_NC_ITEM_EQUIPCHANGE_CMD`, así como los códigos de resultado de `PROTO_NC_ITEM_EQUIP_ACK` / `UNEQUIP` / `AUTO_ARRANGE_INVEN_ACK`.
2. **Charged Inventory y Reward Inventory 100% Funcionales:**
   - Verificar en `Zone.exe` la estructura exacta de `PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK`, `PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK`, `PROTO_NC_ITEM_REWARDINVENOPEN_ACK` y `PROTO_NC_ITEM_REWARDINVENGET_ACK`.
   - Confirmar que al retirar un ítem (`PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ` o `REWARD`), el servidor crea el ítem en el inventario normal (`nStorageType = 2`) en un slot libre, envía `PROTO_NC_ITEM_CELLCHANGE_CMD` + `WITHDRAW_ACK`, y permite moverlo, usarlo o equiparlo sin errores.
3. **Cobertura 100% de `NC_ITEM` (48 C2S), `NC_ITEMDB`, `NC_BOOTH` (8 C2S) y `NC_CHARSAVE`:**
   - Todos los handlers de estos 4 departamentos auditados contra `Zone.exe` y `Zone.md`.

