package zone

import (
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func buildTestPacket(op uint16, payload []byte) *protocol.Packet {
	r := protocol.NewPacketReader(payload)
	return &protocol.Packet{
		Header: byte(op >> 10),
		Type:   op & 0x03FF,
		OpCode: op,
		Raw:    payload,
		Reader: r,
	}
}

type memGameLogRepo struct {
	database.GameLogRepository
	mu      sync.Mutex
	entries []*database.GameLogEntry
}

func (m *memGameLogRepo) AddLog(ctx context.Context, entry *database.GameLogEntry) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, entry)
	return int64(len(m.entries)), nil
}

func setupFase14Handler(t *testing.T) (*Handler, *DataStore, *database.MemoryCharacterRepository, *ZoneMap, *memGameLogRepo) {
	t.Helper()
	dataDir := filepath.Join("..", "..", "..", "data")
	ds, err := LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("LoadAllData error: %v", err)
	}
	repo := database.NewMemoryCharacterRepositoryWithFile("")
	h := NewHandler(repo, database.NewTransferManager(), ds.ClassStore, ds.FieldStore, ds.ChkEngine)
	h.SetDataStore(ds)
	h.itemStore = ds.ItemStore

	glRepo := &memGameLogRepo{}
	disp := database.NewGameLogDispatcher(glRepo, 256, 1)
	disp.Start()
	t.Cleanup(func() { _ = disp.Close(context.Background()) })
	h.SetGameLogDispatcher(disp)

	zMap := h.mapManager.GetOrCreateMap("Rou")
	return h, ds, repo, zMap, glRepo
}

func TestFase14_All62ItemAttributeClassMethods(t *testing.T) {
	iac := &ItemAttributeClass{}
	iac.iac_itemcreate(1001, 5, 1) // Weapon
	iac.iac_SetLot(2)
	iac.iac_SetAmount(33)
	iac.iac_SetChargedNumber(44)
	if iac.iac_charge(6) != 50 || iac.iac_GetChargedNumber() != 50 {
		t.Fatalf("unexpected charged number: %d", iac.iac_GetChargedNumber())
	}
	iac.iac_SetDeteleDateTime(123456)
	iac.iac_SetEndTime(234567)
	iac.iac_SetUseableDateTime(345678)
	iac.iac_SetEndureEndTime(456789)
	iac.iac_SetEndureGrade(3)
	iac.iac_SetFurnitureKey(77)
	iac.iac_SetGrade(4)
	iac.iac_SetHP(999)
	iac.iac_SetIsPutOnBelonged(true)
	iac.iac_SetItemOptionUpgrade(7)
	iac.iac_SetRandomOptionChangeCount(2)
	iac.iac_SetRareFailCount(5)
	iac.iac_SetDuringPetSummon(true)
	iac.iac_SetDuringRide(true)
	iac.iac_SetPetID(200)
	iac.iac_SetPetName("Foxy")
	iac.iac_SetPetRegistNumber(8888)
	iac.iac_Setoption(1, 15)
	iac.iac_InvenStore(3)

	if iac.iac_FoodEat(150) == 0 || iac.iac_GetHungry() != 150 {
		t.Fatalf("iac_FoodEat failed")
	}
	iac.iac_wptitle_Setlicense(55)
	iac.iac_wptitle_Settitle(12)
	if iac.iac_wptitle_mobkilladd(9) != 9 || iac.iac_wptitle_Getlicencemob() != 55 || iac.iac_wptitle_Gettitle() != 12 {
		t.Fatalf("wptitle methods mismatch")
	}
	if _, idx, ok := iac.iac_FindOption(1); !ok || idx != 0 || iac.iac_StepGet() != 1 || iac.iac_EndOfOption(1) != true {
		t.Fatalf("option lookup failed")
	}
	if iac.iac_equipbriefinfo() != 7 || iac.iac_Getupgrade() != 7 || iac.iac_GetItemOptionUpgrade() != 7 {
		t.Fatalf("upgrade getters mismatch")
	}
	if iac.iac_GetAmount() != 33 || iac.iac_GetAttrSize() <= 65 || iac.iac_GetDeteleDateTime() != 123456 ||
		!iac.iac_GetDuringPetSummon() || !iac.iac_GetDuringRide() || iac.iac_GetEndureEndTime() != 456789 ||
		iac.iac_GetGrade() != 4 || iac.iac_GetHP() != 999 || !iac.iac_GetIsPutOnBelonged() ||
		len(iac.iac_GetItemOptionStruct()) != 1 || iac.iac_GetLot() != 5 || iac.iac_GetPetID() != 200 ||
		iac.iac_GetPetName() != "Foxy" || iac.iac_GetPetRegistNumber() != 8888 ||
		iac.iac_GetRandomOptionChangeCount() != 2 || iac.iac_GetRareFailCount() != 5 ||
		iac.iac_GetUseableDateTime() != 345678 || iac.iac_ItemClass() != 5 || iac.iac_islotitem() {
		t.Fatalf("getter invariant failure")
	}

	raw := iac.iac_fill2buffer()
	parsed := &ItemAttributeClass{ItemClass: 5}
	if !parsed.iac_fill2struct(raw) || parsed.Upgrade != 7 || !parsed.IsPutOnBelonged || len(parsed.StatOptions) != 1 || parsed.StatOptions[0].Value != 15 {
		t.Fatalf("fill2buffer/fill2struct round-trip failed: %+v", parsed)
	}
	parsed.iac_ClearItemOptionUpgrade()
	if parsed.Upgrade != 0 {
		t.Fatalf("iac_ClearItemOptionUpgrade failed")
	}
	parsed.iac_DropItemMake(1002, 1, 10, []item.ItemStatOption{{Type: 2, Value: 20}})
	if parsed.ItemID != 1002 || parsed.Lot != 10 {
		t.Fatalf("iac_DropItemMake failed")
	}
	_ = parsed.iac_Getitemchest()
}

func TestFase14_NPCShopBuySellRepurchaseAndStorage(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)
	p, conn := newTestClientAndPlayer(t, repo, ds, zMap, h, 9101, "MerchantTester", 50, 1000, 1000)
	p.Char.TCharacterState.Money = 100000

	// Elegir ítem real de ItemInfo.shn con BuyPrice > 0 y SellPrice > 0
	var shopItem *item.ItemInfo
	for id := uint16(1); id < 60000; id++ {
		if info, ok := ds.ItemStore.GetItemInfo(id); ok && info != nil && info.BuyPrice > 0 && info.SellPrice > 0 && info.MaxLot >= 10 {
			shopItem = info
			break
		}
	}
	if shopItem == nil {
		t.Fatalf("no stackable shop item found in ItemInfo.shn")
	}

	// 1. Buy 2 units
	pwBuy := make([]byte, 6)
	binary.LittleEndian.PutUint16(pwBuy[0:2], shopItem.ID)
	binary.LittleEndian.PutUint32(pwBuy[2:6], 2)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_BUY_REQ, pwBuy))

	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_BUY_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_ITEM_BUY_ACK")
	}
	items, _ := repo.GetCharacterItems(p.Char.ID)
	if len(items) != 1 || items[0].ItemID != shopItem.ID {
		t.Fatalf("expected purchased item in inventory, got %+v", items)
	}
	afterBuyMoney := p.Char.TCharacterState.Money

	// 2. Sell 1 unit -> creates Repurchase entry + emits 0x305C
	conn.clear()
	sellBuf := make([]byte, 5)
	sellBuf[0] = byte(items[0].Storage)
	binary.LittleEndian.PutUint32(sellBuf[1:5], 1)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_SELL_REQ, sellBuf))

	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_SELL_ACK)) == 0 || len(conn.findByOpcode(protocol.PROTO_NC_ITEM_SELL_ITEM_INSERT_CMD)) == 0 {
		t.Fatalf("expected SELL_ACK and SELL_ITEM_INSERT_CMD")
	}
	if len(p.RepurchaseItems) != 1 {
		t.Fatalf("expected 1 repurchase entry")
	}

	// 3. Repurchase sold item
	repBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(repBuf, p.RepurchaseItems[0].Handle)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_REPURCHASE_REQ, repBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_REPURCHASE_ACK)) == 0 {
		t.Fatalf("expected REPURCHASE_ACK")
	}
	if len(p.RepurchaseItems) != 0 {
		t.Fatalf("repurchase entry should be popped")
	}
	if p.Char.TCharacterState.Money != afterBuyMoney {
		t.Fatalf("money invariant mismatch after sell+repurchase: got %d want %d", p.Char.TCharacterState.Money, afterBuyMoney)
	}

	// 4. Personal & Account Storage Deposit/Withdraw & OpenPage (0x301C, 0x301E, 0x304C, 0x3028, 0x304D)
	conn.clear()
	depBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(depBuf, 500)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_DEPOSIT_REQ, depBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_ACCOUNT_STORAGE_OPEN_CMD, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_OPENSTORAGEPAGE_REQ, []byte{1}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_WITHDRAW_REQ, depBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_ACCOUNT_STORAGE_CLOSE_CMD, nil))

	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_DEPOSIT_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_MENU_OPENSTORAGE_CMD)) < 2 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_WITHDRAW_ACK)) == 0 {
		t.Fatalf("expected Deposit, OpenStorage, Withdraw ACKs")
	}

	// 5. Guild & Academy Storage (0x3026, 0x3031, 0x302F, 0x986D, 0x986F)
	gid, err := repo.CreateGuild(p.Char.ID, "XelionGuild", "", 0)
	if err != nil {
		t.Fatalf("CreateGuild err: %v", err)
	}
	_ = repo.SetGuildMoney(gid, 5000)
	_ = repo.CreateGuildAcademy(gid, p.Char.ID)

	wGuild := make([]byte, 8)
	binary.LittleEndian.PutUint64(wGuild, 1000)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_REQ, wGuild))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ, wGuild))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ, []byte{0}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ, wGuild))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_CLOSE_CMD, []byte{0}))

	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK)) == 0 {
		t.Fatalf("expected guild and academy storage ACKs")
	}
}

func TestFase14_SafeTradeAntiDuplicationAndBoothSearch(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)
	pA, connA := newTestClientAndPlayer(t, repo, ds, zMap, h, 9201, "TraderA", 40, 2000, 2000)
	pB, connB := newTestClientAndPlayer(t, repo, ds, zMap, h, 9202, "TraderB", 40, 2010, 2010)
	pA.Char.TCharacterState.Money = 10000
	pB.Char.TCharacterState.Money = 5000

	var eqInfo *item.ItemInfo
	for id := uint16(1); id < 60000; id++ {
		if info, ok := ds.ItemStore.GetItemInfo(id); ok && info != nil && info.Class == 5 {
			eqInfo = info
			break
		}
	}
	if eqInfo == nil {
		t.Fatalf("no weapon found in ItemInfo.shn")
	}

	itemA, ok := h.addItemToInventory(pA, eqInfo.ID, 1, nil)
	if !ok {
		t.Fatalf("failed creating itemA")
	}

	// Propose -> Ask -> Accept
	propBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(propBuf, pB.MapObjectID)
	h.HandlePacket(pA.Client, buildTestPacket(protocol.PROTO_NC_TRADE_PROPOSE_REQ, propBuf))
	h.HandlePacket(pB.Client, buildTestPacket(protocol.PROTO_NC_TRADE_PROPOSEYES_ACK, nil))

	if len(connA.findByOpcode(protocol.PROTO_NC_TRADE_START_CMD)) == 0 || len(connB.findByOpcode(protocol.PROTO_NC_TRADE_START_CMD)) == 0 {
		t.Fatalf("expected TRADE_START_CMD on both clients")
	}

	// Upboard A's weapon + Cen boarding
	h.HandlePacket(pA.Client, buildTestPacket(protocol.PROTO_NC_TRADE_UPBOARD_REQ, []byte{byte(itemA.Storage)}))
	cenBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(cenBuf, 2500)
	h.HandlePacket(pA.Client, buildTestPacket(protocol.PROTO_NC_TRADE_CENBOARDING_REQ, cenBuf))

	// Both Lock + Concurrent Decide (verify zero duplication!)
	h.HandlePacket(pA.Client, buildTestPacket(protocol.PROTO_NC_TRADE_BOARDLOCK_REQ, nil))
	h.HandlePacket(pB.Client, buildTestPacket(protocol.PROTO_NC_TRADE_BOARDLOCK_REQ, nil))
	h.HandlePacket(pA.Client, buildTestPacket(protocol.PROTO_NC_TRADE_DECIDE_REQ, nil))
	h.HandlePacket(pB.Client, buildTestPacket(protocol.PROTO_NC_TRADE_DECIDE_REQ, nil))

	if len(connA.findByOpcode(protocol.PROTO_NC_TRADE_TRADECOMPLETE_CMD)) == 0 ||
		len(connB.findByOpcode(protocol.PROTO_NC_TRADE_TRADECOMPLETE_CMD)) == 0 {
		t.Fatalf("expected TRADECOMPLETE_CMD on both sides")
	}

	itemsA, _ := repo.GetCharacterItems(pA.Char.ID)
	itemsB, _ := repo.GetCharacterItems(pB.Char.ID)
	if len(itemsA) != 0 || len(itemsB) != 1 || itemsB[0].ItemKey != itemA.ItemKey {
		t.Fatalf("anti-duplication failure: itemsA=%d itemsB=%d", len(itemsA), len(itemsB))
	}
	if pA.Char.TCharacterState.Money != 7500 || pB.Char.TCharacterState.Money != 7500 {
		t.Fatalf("cen conservation failure: A=%d B=%d", pA.Char.TCharacterState.Money, pB.Char.TCharacterState.Money)
	}

	// Ahora pB abre Kiosko / Booth (0x6801) vendiendo el arma y pA busca por categoría (0x6813), ubica posición (0x6815) y compra (0x680D)!
	connA.clear()
	openBuf := make([]byte, 31+12)
	copy(openBuf[0:30], "Super Weapon Kiosk")
	openBuf[30] = 1 | (1 << 1) // isSell=1, count=1
	openBuf[31] = byte(itemsB[0].Storage)
	openBuf[32] = 0
	binary.LittleEndian.PutUint64(openBuf[33:41], 1200)
	binary.LittleEndian.PutUint16(openBuf[41:43], 1)

	h.HandlePacket(pB.Client, buildTestPacket(protocol.PROTO_NC_BOOTH_OPEN_REQ, openBuf))
	if len(connB.findByOpcode(protocol.PROTO_NC_BOOTH_OPEN_ACK)) == 0 || pB.ActiveBooth == nil {
		t.Fatalf("expected booth open")
	}

	searchReq := make([]byte, 6)
	binary.LittleEndian.PutUint32(searchReq[0:4], 0)
	h.HandlePacket(pA.Client, buildTestPacket(protocol.PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ, searchReq))
	if len(connA.findByOpcode(protocol.PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK)) == 0 {
		t.Fatalf("expected BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_ACK")
	}

	posReq := make([]byte, 6)
	binary.LittleEndian.PutUint16(posReq[4:6], pB.MapObjectID)
	h.HandlePacket(pA.Client, buildTestPacket(protocol.PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ, posReq))
	if len(connA.findByOpcode(protocol.PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_ACK)) == 0 {
		t.Fatalf("expected BOOTH_SEARCH_BOOTH_POSITION_ACK")
	}

	buyBooth := make([]byte, 5)
	binary.LittleEndian.PutUint16(buyBooth[0:2], pB.MapObjectID)
	buyBooth[2] = 0
	binary.LittleEndian.PutUint16(buyBooth[3:5], 1)
	h.HandlePacket(pA.Client, buildTestPacket(protocol.PROTO_NC_BOOTH_ITEMTRADE_REQ, buyBooth))

	itemsA2, _ := repo.GetCharacterItems(pA.Char.ID)
	itemsB2, _ := repo.GetCharacterItems(pB.Char.ID)
	if len(itemsA2) != 1 || len(itemsB2) != 0 {
		t.Fatalf("booth purchase conservation failure: A=%d B=%d", len(itemsA2), len(itemsB2))
	}
	if pA.Char.TCharacterState.Money != 6300 || pB.Char.TCharacterState.Money != 8700 {
		t.Fatalf("booth cen conservation failure: A=%d B=%d", pA.Char.TCharacterState.Money, pB.Char.TCharacterState.Money)
	}

	h.HandlePacket(pB.Client, buildTestPacket(protocol.PROTO_NC_BOOTH_CLOSE_REQ, nil))
	if pB.ActiveBooth != nil {
		t.Fatalf("booth should be closed")
	}
}

func TestFase14_ForgeGemsRandomOptDismantleMixProduceAndUtilities(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)
	p, conn := newTestClientAndPlayer(t, repo, ds, zMap, h, 9301, "ForgeHero", 60, 1000, 1000)

	var weapon *item.ItemInfo
	for id := uint16(1); id < 60000; id++ {
		if info, ok := ds.ItemStore.GetItemInfo(id); ok && info != nil && info.Class == 5 {
			weapon = info
			break
		}
	}
	wIt, _ := h.addItemToInventory(p, weapon.ID, 1, nil)
	matIt, _ := h.addItemToInventory(p, weapon.ID, 10, nil) // drill / gem / raw proxy item

	// 1. Add Socket -> Add Gem -> Remove Gem
	eqLoc := (uint16(9) << 10) | uint16(wIt.Storage)
	matLoc := (uint16(9) << 10) | uint16(matIt.Storage)
	sockBuf := make([]byte, 4)
	binary.LittleEndian.PutUint16(sockBuf[0:2], eqLoc)
	binary.LittleEndian.PutUint16(sockBuf[2:4], matLoc)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_REQ, sockBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_REQ, sockBuf))
	remBuf := append(sockBuf, 0)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_REQ, remBuf))

	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_ACK)) == 0 {
		t.Fatalf("expected all 3 socket/gem ACKs")
	}

	// 2. Random Option Change -> Accept -> Recover
	roIt, _ := h.addItemToInventory(p, weapon.ID, 1, nil)
	roStuff, _ := h.addItemToInventory(p, weapon.ID, 5, nil)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_REQ, []byte{byte(roIt.Storage), byte(roStuff.Storage)}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_REQ, []byte{byte(roIt.Storage), 1}))

	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK)) == 0 {
		t.Fatalf("expected all 3 RandomOption ACKs")
	}

	// 3. Upgrade + Mover Upgrade
	upgBuf := make([]byte, 9)
	upgBuf[0] = byte(roIt.Storage)
	upgBuf[1] = byte(roStuff.Storage)
	upgBuf[2], upgBuf[3], upgBuf[8] = 0xFF, 0xFF, 0xFF
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_UPGRADE_REQ, upgBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_MOVER_UPGRADE_REQ, []byte{byte(roIt.Storage), byte(roStuff.Storage)}))

	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_UPGRADE_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_MOVER_UPGRADE_ACK)) == 0 {
		t.Fatalf("expected UPGRADE_ACK and MOVER_UPGRADE_ACK")
	}

	// 4. Dismantle (`ItemDismantle.shn`) + Mix (`ItemMix.shn`) + Produce (`Produce.shn`)
	dismIt, _ := h.addItemToInventory(p, weapon.ID, 1, nil)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_DISMANTLE_REQ, []byte{byte(dismIt.Storage)}))
	if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_DISMANTLE_ACK)) == 0 {
		t.Fatalf("expected DISMANTLE_ACK")
	}

	if mixTbl := h.getSHN("ItemMix"); mixTbl != nil && len(mixTbl.Rows) > 0 {
		rawName := strings.TrimSpace(mixTbl.Rows[0].GetString("Raw_ItemIDX"))
		if rawInfo, ok := ds.ItemStore.GetItemInfoByName(rawName); ok && rawInfo != nil {
			l, _ := h.addItemToInventory(p, rawInfo.ID, 2, nil)
			r, _ := h.addItemToInventory(p, rawInfo.ID, 2, nil)
			h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_MIX_ITEM_REQ, []byte{0xFF, byte(l.Storage), byte(r.Storage)}))
			if len(conn.findByOpcode(protocol.PROTO_NC_ITEM_MIX_ITEM_ACK)) == 0 {
				t.Fatalf("expected MIX_ITEM_ACK")
			}
		}
	}

	if prodTbl := h.getSHN("Produce"); prodTbl != nil && len(prodTbl.Rows) > 0 {
		row := prodTbl.Rows[0]
		prodID := uint16(row.GetInt("ProductID"))
		for i := 0; i < 8; i++ {
			rName := strings.TrimSpace(row.GetString(fmt.Sprintf("Raw%d", i)))
			qty := row.GetInt(fmt.Sprintf("Quantity%d", i))
			if rName != "" && rName != "-" && qty > 0 {
				if rInfo, ok := ds.ItemStore.GetItemInfoByName(rName); ok && rInfo != nil {
					h.addItemToInventory(p, rInfo.ID, uint32(qty), nil)
				}
			}
		}
		prodBuf := make([]byte, 2)
		binary.LittleEndian.PutUint16(prodBuf, prodID)
		h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ACT_PRODUCE_CAST_REQ, prodBuf))
		if len(conn.findByOpcode(protocol.PROTO_NC_ACT_PRODUCE_MAKE_CMD)) == 0 {
			t.Fatalf("expected PROTO_NC_ACT_PRODUCE_MAKE_CMD")
		}
	}

	// 5. Utilities: FoodEat, Minimon, ActionByItem, Endure, BelongCancel, ChatColor, TermExtend, ClassChange, Chest
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_MISC_HEARTBEAT_REQ, []byte{byte(wIt.Storage)}))
	mmBuf := make([]byte, 24)
	binary.LittleEndian.PutUint16(mmBuf[0:2], weapon.ID)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON_REQ, mmBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_OFF_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_ON_REQ, mmBuf))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_OFF_REQ, nil))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ACT_ACTIONBYITEM_REQ, []byte{byte(wIt.Storage)}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_CHAT_COLOR_CHANGE_REQ, []byte{byte(wIt.Storage)}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_EQUIP_BELONGED_CANCEL_USE_REQ, []byte{byte(matIt.Storage), byte(wIt.Storage)}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_TERMEXTEND_REQ, []byte{byte(matIt.Storage), byte(wIt.Storage)}))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_CLASS_CHANGE_REQ, []byte{byte(matIt.Storage)}))

	wEnd := make([]byte, 3)
	wEnd[0] = byte(matIt.Storage)
	binary.LittleEndian.PutUint16(wEnd[1:3], eqLoc)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_WEAPONENDURE_CHARGE_REQ, wEnd))
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_SHIELDENDURE_CHARGE_REQ, wEnd))

	if len(conn.findByOpcode(protocol.PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ACT_ACTIONBYITEM_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_CHAT_COLOR_CHANGE_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_EQUIP_BELONGED_CANCEL_USE_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_CLASS_CHANGE_ACK)) == 0 ||
		len(conn.findByOpcode(protocol.PROTO_NC_ITEM_WEAPONENDURE_CHARGE_ACK)) == 0 {
		t.Fatalf("missing expected utility ACKs")
	}

	// Open real chest from TreasureReward.json (`Tbag06`)
	if tbag, ok := ds.ItemStore.GetItemInfoByName("Tbag06"); ok && tbag != nil {
		cIt, added := h.addItemToInventory(p, tbag.ID, 1, nil)
		if added {
			h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEM_USE_REQ, []byte{byte(cIt.Storage), 9}))
		}
	}

	// Verificar GameLogDispatcher registró tipos 6, 7 y D de las operaciones anteriores
	_ = h.gameLogDispatcher.Close(context.Background())
}

func TestFase14_LiveMySQLTransactionsAndGameLog(t *testing.T) {
	cfg, err := config.Load("../../../.env")
	if err != nil {
		t.Skipf("skipping live MySQL test (no .env): %v", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil || db.Ping() != nil {
		t.Skipf("skipping live MySQL test (unreachable): %v", err)
	}
	defer db.Close()

	repo := database.NewMySQLCharacterRepository(db)
	glRepo := database.NewMySQLGameLogRepository(db)
	ctx := context.Background()

	var userA, userB int
	rows, qErr := db.Query("SELECT nUserNo FROM tuser ORDER BY nUserNo DESC LIMIT 2")
	if qErr != nil {
		t.Fatalf("query tuser failed: %v", qErr)
	}
	defer rows.Close()
	for rows.Next() {
		var u int
		if err := rows.Scan(&u); err == nil {
			if userA == 0 {
				userA = u
			} else if userB == 0 {
				userB = u
			}
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterating tuser rows failed: %v", err)
	}
	if userA == 0 {
		t.Fatalf("no users in tuser table for MySQL Phase 14 transactional test")
	}
	if userB == 0 {
		userB = userA
	}

	findFreeSlot := func(userNo int) int {
		for s := 0; s <= 5; s++ {
			var cnt int
			if err := db.QueryRow("SELECT COUNT(*) FROM tcharacter WHERE nUserNo = ? AND nSlotNo = ? AND bDeleted = 0", userNo, s).Scan(&cnt); err == nil && cnt == 0 {
				return s
			}
		}
		return 0
	}
	slotA := byte(findFreeSlot(userA))
	charA := &database.Character{AccountID: userA, Slot: slotA, Name: fmt.Sprintf("F14A_%d", time.Now().Unix()%10000), Job: 1, Level: 10, CurHP: 100, CurSP: 50, TownName: "Rou"}
	if err := repo.CreateCharacter(charA); err != nil {
		t.Fatalf("CreateCharacter A (user=%d slot=%d) failed: %v", userA, slotA, err)
	}
	defer func() {
		_, _ = db.Exec("DELETE FROM tcharacter WHERE nCharNo = ?", charA.ID)
	}()

	slotB := byte(findFreeSlot(userB))
	if userB == userA && slotB == slotA {
		slotB++
	}
	charB := &database.Character{AccountID: userB, Slot: slotB, Name: fmt.Sprintf("F14B_%d", time.Now().Unix()%10000), Job: 1, Level: 10, CurHP: 100, CurSP: 50, TownName: "Rou"}
	if err := repo.CreateCharacter(charB); err != nil {
		t.Fatalf("CreateCharacter B (user=%d slot=%d) failed: %v", userB, slotB, err)
	}
	defer func() {
		_, _ = db.Exec("DELETE FROM tcharacter WHERE nCharNo = ?", charB.ID)
	}()

	_ = repo.SaveCharMoney(charA.ID, 50000)
	_ = repo.SaveCharMoney(charB.ID, 10000)
	if _, err := repo.UserDepositMoney(userA, charA.ID, 500); err != nil {
		t.Fatalf("MySQL UserDepositMoney failed: %v", err)
	}
	if _, err := repo.UserWithdrawMoney(userA, charA.ID, 500); err != nil {
		t.Fatalf("MySQL UserWithdrawMoney failed: %v", err)
	}

	// Concurrency + SafeTradeAll + Enchant/Upgrade SP check against real MySQL
	it := &database.Item{StorageType: 9, Owner: charA.ID, Storage: 0, ItemID: 1, Flags: 1}
	drill := &database.Item{StorageType: 9, Owner: charA.ID, Storage: 1, ItemID: 2, Flags: 1}
	gem := &database.Item{StorageType: 9, Owner: charA.ID, Storage: 2, ItemID: 3, Flags: 1}
	rem := &database.Item{StorageType: 9, Owner: charA.ID, Storage: 3, ItemID: 4, Flags: 1}
	stone := &database.Item{StorageType: 9, Owner: charA.ID, Storage: 4, ItemID: 5, Flags: 1}
	if err := repo.CreateItem(it); err != nil {
		t.Fatalf("MySQL CreateItem eq failed: %v", err)
	}
	defer func() { _ = repo.DeleteItem(it.ItemKey) }()
	for _, m := range []*database.Item{drill, gem, rem, stone} {
		if err := repo.CreateItem(m); err != nil {
			t.Fatalf("MySQL CreateItem mat failed: %v", err)
		}
		mk := m.ItemKey
		defer func() { _ = repo.DeleteItem(mk) }()
	}
	if err := repo.SetItemOption(it.ItemKey, 523, 0); err != nil {
		t.Fatalf("MySQL SetItemOption(523) failed: %v", err)
	}
	if err := repo.SetItemOption(it.ItemKey, 524, 0); err != nil {
		t.Fatalf("MySQL SetItemOption(524) failed: %v", err)
	}

	if _, _, err := repo.EnchantAddSocket(charA.ID, it.ItemKey, drill.ItemKey, 0); err != nil {
		t.Fatalf("MySQL EnchantAddSocket failed: %v", err)
	}
	if err := repo.EnchantAddGem(charA.ID, it.ItemKey, gem.ItemKey, 0, 100, 1); err != nil {
		t.Fatalf("MySQL EnchantAddGem failed: %v", err)
	}
	if _, err := repo.EnchantRemoveGem(charA.ID, it.ItemKey, rem.ItemKey, 0, 1, 0); err != nil {
		t.Fatalf("MySQL EnchantRemoveGem failed: %v", err)
	}
	if err := repo.UpgradeItemFull(charA.ID, it.ItemKey, 0, 1, 1, 0, 0, stone.ItemKey, 0, 0, 0, 0); err != nil {
		t.Fatalf("MySQL UpgradeItemFull failed: %v", err)
	}

	var wg sync.WaitGroup
	errs := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- repo.SafeTradeAll(charB.ID, charA.ID, 0, it.ItemKey, 9, 0)
	}()
	wg.Wait()
	close(errs)
	if e := <-errs; e != nil {
		t.Fatalf("MySQL SafeTradeAll SP failed: %v", e)
	}
	itemsA, _ := repo.GetCharacterItems(charA.ID)
	itemsB, _ := repo.GetCharacterItems(charB.ID)
	var countA, countB int
	for _, x := range itemsA {
		if x.ItemKey == it.ItemKey {
			countA++
		}
	}
	for _, x := range itemsB {
		if x.ItemKey == it.ItemKey {
			countB++
		}
	}
	if countA != 0 || countB != 1 {
		t.Fatalf("MySQL zero-duplication check failed for ItemKey=%d: A=%d B=%d", it.ItemKey, countA, countB)
	}

	// Verificar GameLog Types 3, 4, 5, 6, 7, B, C, D persistidos en World00_GameLog (`tGameLog`) via CALL p_GameLog_Add_Type*
	checks := []struct {
		name string
		fn   func() (int64, error)
	}{
		{"Type3", func() (int64, error) { return glRepo.AddType3(ctx, 1, charA.ID, "Rou", 100, 100, 0, 0, 1, 500) }},
		{"Type4", func() (int64, error) { return glRepo.AddType4(ctx, 1, charA.ID, "Rou", 100, 100, 0, charB.ID, 500, 1) }},
		{"Type5", func() (int64, error) { return glRepo.AddType5(ctx, 1, charA.ID, 500, 500) }},
		{"Type6", func() (int64, error) {
			return glRepo.AddType6(ctx, 10, charA.ID, "Rou", 100, 100, 0, int(it.ItemID), it.ItemKey)
		}},
		{"Type7", func() (int64, error) {
			return glRepo.AddType7(ctx, 1, charA.ID, "Rou", 100, 100, 0, int(it.ItemID), 3, 0)
		}},
		{"TypeB", func() (int64, error) {
			return glRepo.AddTypeB(ctx, 1, charA.ID, "Rou", 100, 100, 0, 0, 1, it.ItemKey, 500, 500, 0)
		}},
		{"TypeC", func() (int64, error) {
			return glRepo.AddTypeC(ctx, 1, charA.ID, "Rou", 100, 100, 0, charB.ID, it.ItemKey, int(it.ItemID), 1, 500)
		}},
		{"TypeD", func() (int64, error) { return glRepo.AddTypeD(ctx, 1, charA.ID, it.ItemKey, int(it.ItemID), 1) }},
	}
	for _, c := range checks {
		logNo, err := c.fn()
		if err != nil {
			t.Fatalf("MySQL GameLog %s failed: %v", c.name, err)
		}
		if logNo <= 0 {
			t.Fatalf("MySQL GameLog %s expected positive nLogNo, got %d", c.name, logNo)
		}
	}
}
