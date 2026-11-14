package zone

import (
	"encoding/binary"
	"strings"
	"sync"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

type memWMSBridge struct {
	mu     sync.Mutex
	events []WMSEvent
}

func (b *memWMSBridge) PublishWMS(ev WMSEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events = append(b.events, ev)
}

func (b *memWMSBridge) countBySymbol(sym string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	n := 0
	for _, ev := range b.events {
		if ev.Symbol == sym {
			n++
		}
	}
	return n
}

func TestFase16_All22NewC2SHandlersAndNoUnhandled(t *testing.T) {
	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	bridge := &memWMSBridge{}
	h.SetWMSBridge(bridge)

	p, conn := newTestClientAndPlayer(t, repo, ds, zMap, h, 9601, "C2SHero", 45, 2000, 2000)
	p2, conn2 := newTestClientAndPlayer(t, repo, ds, zMap, h, 9602, "PartyPal", 45, 2100, 2100)
	_ = conn2

	// 1. PINGTEST_CLIENT_ZONE_DB & TOOL_WM_CLIENT_ZONE_DB
	ping90 := make([]byte, 90)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB, ping90))
	if len(conn.findByOpcode(protocol.PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB)) == 0 {
		t.Fatalf("expected PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB response")
	}

	ping130 := make([]byte, 130)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB, ping130))
	if len(conn.findByOpcode(protocol.PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB)) == 0 {
		t.Fatalf("expected PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB response")
	}
	if bridge.countBySymbol("wms_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB") == 0 {
		t.Fatalf("expected wms_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB bridge event")
	}

	// 2. SERVERPARAMETER_REQ
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_MISC_SERVERPARAMETER_REQ, nil))
	if len(conn.findByOpcode(protocol.PROTO_NC_MISC_SERVERPARAMETER_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_MISC_SERVERPARAMETER_ACK")
	}

	// 3. USER_REGISENUMBER_REQ & USER_USE_BEAUTY_SHOP_CMD
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_USER_REGISENUMBER_REQ, nil))
	if len(conn.findByOpcode(protocol.PROTO_NC_USER_REGISENUMBER_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_USER_REGISENUMBER_ACK")
	}
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_USER_USE_BEAUTY_SHOP_CMD, []byte{4, 3, 2}))
	if !p.BeautyShopOpen || p.Char.Hair != 4 || p.Char.HairColor != 3 || p.Char.Face != 2 {
		t.Fatalf("expected BeautyShop state updated, got open=%v hair=%d color=%d face=%d",
			p.BeautyShopOpen, p.Char.Hair, p.Char.HairColor, p.Char.Face)
	}

	// 4. CHAR_LOGOUTREADY_CMD & CHAR_LOGOUTCANCEL_CMD
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_CHAR_LOGOUTREADY_CMD, nil))
	if !p.LogoutPending {
		t.Fatalf("expected LogoutPending=true after LOGOUTREADY_CMD")
	}
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_CHAR_LOGOUTCANCEL_CMD, nil))
	if p.LogoutPending {
		t.Fatalf("expected LogoutPending=false after LOGOUTCANCEL_CMD")
	}

	// 5. CHAR_SAVE_LINK_REQ (usp_LinkScroll_SetDestMap)
	linkBuf := make([]byte, 12)
	copy(linkBuf, "Eld")
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_CHAR_SAVE_LINK_REQ, linkBuf))
	if p.Char.KQMapName != "Eld" || bridge.countBySymbol("wms_NC_CHAR_ZONE_LINK_FROM_CMD") == 0 {
		t.Fatalf("expected CHAR_SAVE_LINK_REQ to set KQMapName='Eld' and emit wms_NC_CHAR_ZONE_LINK_FROM_CMD, got %q", p.Char.KQMapName)
	}

	// 6. ACT_CHANGEMODE_REQ, ACT_SHOUT_CMD & ACT_ROAR_REQ
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ACT_CHANGEMODE_REQ, []byte{2}))
	if len(conn.findByOpcode(protocol.PROTO_NC_ACT_SOMEONECHANGEMODE_CMD)) == 0 || p.ActionMode != 2 {
		t.Fatalf("expected ActionMode=2 and PROTO_NC_ACT_SOMEONECHANGEMODE_CMD")
	}

	shoutMsg := "Hello Zone World!"
	shoutBuf := append([]byte{0, byte(len(shoutMsg))}, []byte(shoutMsg)...)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ACT_SHOUT_CMD, shoutBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_ACT_SOMEONESHOUT_CMD)) == 0 {
		t.Fatalf("expected PROTO_NC_ACT_SOMEONESHOUT_CMD broadcast")
	}

	_ = repo.CreateItem(&database.Item{ItemKey: 777001, Owner: p.Char.ID, ItemID: 8001, StorageType: 9, Storage: 10, Flags: 2})
	roarMsg := "World Roar!"
	roarBuf := append([]byte{10, byte(len(roarMsg))}, []byte(roarMsg)...)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ACT_ROAR_REQ, roarBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_ACT_ROAR_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_ACT_ROAR_ACK")
	}

	// 7. PARTY_MEMBERINFOREQ_CMD & CHARSAVE_UI_STATE_SAVE_REQ
	partyNo, err := repo.CreateParty(p.Char.ID, p2.Char.ID)
	if err != nil || partyNo <= 0 {
		t.Fatalf("CreateParty failed: %v", err)
	}
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_PARTY_MEMBERINFOREQ_CMD, nil))
	if len(conn.findByOpcode(protocol.PROTO_NC_PARTY_MEMBERINFORM_CMD)) == 0 {
		t.Fatalf("expected PROTO_NC_PARTY_MEMBERINFORM_CMD response")
	}

	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ, []byte{7}))
	if len(conn.findByOpcode(protocol.PROTO_NC_CHARSAVE_UI_STATE_SAVE_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_CHARSAVE_UI_STATE_SAVE_ACK")
	}

	// 8. FRIEND_UES_FRIEND_POINT_REQ
	_ = repo.SetFriendPoints(p.Char.ID, 50)
	fpBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(fpBuf, 15)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_FRIEND_UES_FRIEND_POINT_REQ, fpBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_FRIEND_UES_FRIEND_POINT_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_FRIEND_UES_FRIEND_POINT_ACK")
	}
	remPts, _ := repo.GetFriendPoints(p.Char.ID)
	if remPts != 35 {
		t.Fatalf("expected remaining friend points 35, got %d", remPts)
	}

	// 9. GUILD_STORAGEOPEN_REQ, ITEMDB_CLOSE_GUILD_STORAGE_CMD, GUILD_RENAME_REQ & GUILD_RETYPE_REQ
	gID, err := repo.CreateGuild(p.Char.ID, "OrderOne", "pass", 1)
	if err != nil {
		t.Fatalf("CreateGuild failed: %v", err)
	}
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GUILD_STORAGEOPEN_REQ, nil))
	if len(conn.findByOpcode(protocol.PROTO_NC_GUILD_STORAGEOPEN_ACK)) == 0 || !p.GuildStorageOpen {
		t.Fatalf("expected GuildStorageOpen=true and ACK")
	}
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_ITEMDB_CLOSE_GUILD_STORAGE_CMD, nil))
	if p.GuildStorageOpen {
		t.Fatalf("expected GuildStorageOpen=false after close cmd")
	}

	_ = repo.CreateItem(&database.Item{ItemKey: 777002, Owner: p.Char.ID, ItemID: 8002, StorageType: 9, Storage: 11, Flags: 1})
	renameBuf := make([]byte, 17)
	renameBuf[0] = 11
	copy(renameBuf[1:17], "OrderTwo")
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GUILD_RENAME_REQ, renameBuf))
	if len(conn.findByOpcode(protocol.PROTO_NC_GUILD_RENAME_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_GUILD_RENAME_ACK")
	}
	gAfter, _, _ := repo.GetGuildByCharID(p.Char.ID)
	if gAfter == nil || gAfter.ID != gID || gAfter.Name != "OrderTwo" {
		t.Fatalf("expected guild renamed to OrderTwo, got %+v", gAfter)
	}

	_ = repo.CreateItem(&database.Item{ItemKey: 777003, Owner: p.Char.ID, ItemID: 8003, StorageType: 9, Storage: 12, Flags: 1})
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_GUILD_RETYPE_REQ, []byte{12, 2}))
	if len(conn.findByOpcode(protocol.PROTO_NC_GUILD_RETYPE_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_GUILD_RETYPE_ACK")
	}

	// 10. USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ & PET_ASK_NEW_NAME_ACK
	_ = repo.CreateItem(&database.Item{ItemKey: 777004, Owner: p.Char.ID, ItemID: 8005, StorageType: 9, Storage: 13, Flags: 1})
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ, []byte{13}))
	if len(conn.findByOpcode(protocol.PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_ACK)) == 0 {
		t.Fatalf("expected PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_ACK")
	}

	petAckBuf := make([]byte, 4)
	binary.LittleEndian.PutUint16(petAckBuf[0:2], 0)
	binary.LittleEndian.PutUint16(petAckBuf[2:4], 4321)
	h.HandlePacket(p.Client, buildTestPacket(protocol.PROTO_NC_PET_ASK_NEW_NAME_ACK, petAckBuf))
	if p.ActivePet == nil || p.ActivePet.Handle != 4321 {
		t.Fatalf("expected ActivePet handle 4321 after PROTO_NC_PET_ASK_NEW_NAME_ACK")
	}

	// 11. USER_NORMALLOGOUT_CMD
	h.HandlePacket(p2.Client, buildTestPacket(protocol.PROTO_NC_USER_NORMALLOGOUT_CMD, []byte{1}))
	if bridge.countBySymbol("wms_NC_USER_CONNECTCUT2ZONE_CMD") == 0 {
		t.Fatalf("expected wms_NC_USER_CONNECTCUT2ZONE_CMD after USER_NORMALLOGOUT_CMD")
	}

	if glEntry := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool { return e.CharNo == p.Char.ID && (e.Type == 4 || e.Type == 6) }); glEntry == nil {
		t.Fatalf("expected GameLog entries dispatched from Phase 16 C2S flows")
	}
}

func TestFase16_S2S_138WMS_And_196GDS_Coordination(t *testing.T) {
	h, ds, repo, zMap, _ := setupFase14Handler(t)
	bridge := &memWMSBridge{}
	h.SetWMSBridge(bridge)

	if h.WMSHandlerCount() != 138 {
		t.Fatalf("expected exactly 138 registered wms_NC_* handlers, got %d", h.WMSHandlerCount())
	}
	if h.GDSMappingCount() != 196 {
		t.Fatalf("expected exactly 196 official gds_NC_* mappings, got %d", h.GDSMappingCount())
	}

	for _, sym := range OfficialWMSSymbols138 {
		if !h.DispatchWMS(sym, nil) {
			t.Fatalf("missing wms handler for %s", sym)
		}
	}
	for _, sym := range OfficialGDSSymbols196 {
		if sp, ok := h.DispatchGDS(sym); !ok || sp == "" {
			t.Fatalf("missing gds->SP mapping for %s", sym)
		}
	}

	p, conn := newTestClientAndPlayer(t, repo, ds, zMap, h, 9610, "S2SHero", 50, 2500, 2500)

	// Inbound wms_NC_CT_CHARTTING_CMD notice
	if !h.DispatchWMS("wms_NC_CT_CHARTTING_CMD", map[string]interface{}{"message": "Cluster notice test"}) {
		t.Fatalf("DispatchWMS wms_NC_CT_CHARTTING_CMD failed")
	}
	if len(conn.findByOpcode(protocol.PROTO_NC_ACT_NOTICE_CMD)) == 0 {
		t.Fatalf("expected cluster notice broadcast to player")
	}

	// Inbound wms_NC_CHAR_ZONE_CHAT_COLOR_CMD
	if !h.DispatchWMS("wms_NC_CHAR_ZONE_CHAT_COLOR_CMD", map[string]interface{}{"charID": p.Char.ID}) {
		t.Fatalf("DispatchWMS wms_NC_CHAR_ZONE_CHAT_COLOR_CMD failed")
	}
	if len(conn.findByOpcode(protocol.PROTO_NC_CHAR_CHAT_COLOR_CMD)) == 0 {
		t.Fatalf("expected PROTO_NC_CHAR_CHAT_COLOR_CMD packet")
	}

	// Outbound GDS SP verification
	if sp, ok := h.DispatchGDS("gds_NC_ITEMDB_ADMINCREATESUC_ACK"); !ok || (!strings.Contains(sp, "p_") && !strings.Contains(sp, "usp_")) {
		t.Fatalf("expected GDS mapping to Character SP, got ok=%v sp=%q", ok, sp)
	}
}

func TestFase16_MultiNode_Zone00ToZone04_And_OfficialGM_AdminLvSet(t *testing.T) {
	// 1. Validar topología multi-nodo oficial Zone00..Zone04
	expectedPorts := map[string]struct {
		zoneNo uint16
		port   int
	}{
		"Zone00": {0, 9016},
		"Zone01": {1, 9019},
		"Zone02": {2, 9022},
		"Zone03": {3, 9025},
		"Zone04": {4, 9028},
	}
	for id, exp := range expectedPorts {
		zNo, port, ok := ResolveZoneNodeConfig(id)
		if !ok || zNo != exp.zoneNo || port != exp.port {
			t.Fatalf("ResolveZoneNodeConfig(%s) = (%d, %d, %v), want (%d, %d, true)", id, zNo, port, ok, exp.zoneNo, exp.port)
		}
	}

	h, ds, repo, zMap, glRepo := setupFase14Handler(t)
	if ds.AdminLvCatalog == nil || len(ds.AdminLvCatalog.Levels) == 0 {
		t.Fatalf("expected AdminLvSet.shn parsed into DataStore.AdminLvCatalog")
	}
	reqMakeItemLv, ok := ds.AdminLvCatalog.RequiredLevel("makeitem")
	if !ok {
		t.Fatalf("expected 'makeitem' command present in AdminLvSet.shn")
	}

	// Test MapManager zone filtering
	h.MapManager().SetZoneFilter(0)
	if !h.MapManager().OwnsMap("Rou") {
		t.Fatalf("expected Zone00 to own 'Rou' (Fiesta=0)")
	}
	h.MapManager().SetZoneFilter(1)
	if h.MapManager().OwnsMap("Rou") {
		t.Fatalf("expected Zone01 NOT to own 'Rou'")
	}
	h.MapManager().SetZoneFilter(-1)

	// 2. Test authorization denial for non-GM player (AdminLevel = 0)
	regular, regConn := newTestClientAndPlayer(t, repo, ds, zMap, h, 9620, "NormalUser", 10, 1500, 1500)
	regular.Char.AdminLevel = 0
	regCmd := "&makeitem 100 1"
	regPayload := append([]byte{0, byte(len(regCmd))}, []byte(regCmd)...)
	h.HandlePacket(regular.Client, buildTestPacket(protocol.PROTO_NC_ACT_CHAT_REQ, regPayload))
	if len(regConn.findByOpcode(protocol.PROTO_NC_ACT_NOTICE_CMD)) == 0 {
		t.Fatalf("expected denial notice for unauthorized GM command")
	}

	// 3. Test GM execution strictly via AdminLvSet.shn + tGMItemReward* + GameLog Type L
	gm, gmConn := newTestClientAndPlayer(t, repo, ds, zMap, h, 9621, "GameMaster", 99, 1600, 1600)
	gm.Char.AdminLevel = reqMakeItemLv

	// Spawn a monster so &killallmob has work to do
	_ = zMap.SpawnMobFromTemplate(&MobTemplate{ID: 1, Name: "TestSlime", Level: 1, MaxHP: 100}, 1620, 1620, 0, "", 0, true)

	// Run &makeitem
	mkCmd := "&makeitem 101 1"
	mkPayload := append([]byte{0, byte(len(mkCmd))}, []byte(mkCmd)...)
	h.HandlePacket(gm.Client, buildTestPacket(protocol.PROTO_NC_ACT_CHAT_REQ, mkPayload))

	if len(gmConn.findByOpcode(protocol.PROTO_NC_ITEM_CELLCHANGE_CMD)) == 0 {
		t.Fatalf("expected PROTO_NC_ITEM_CELLCHANGE_CMD after GM &makeitem")
	}
	invItems, err := repo.GetItemsByStorage(gm.Char.ID, 9)
	if err != nil || len(invItems) == 0 || invItems[0].ItemID != 101 {
		t.Fatalf("expected GM reward item 101 persisted in storage 9, got err=%v items=%+v", err, invItems)
	}

	// Test &adminlevel updating target player admin level
	if admLvReq, exists := ds.AdminLvCatalog.RequiredLevel("adminlevel"); exists && admLvReq > gm.Char.AdminLevel {
		gm.Char.AdminLevel = admLvReq
	}
	admCmd := "&adminlevel NormalUser 50"
	h.HandlePacket(gm.Client, buildTestPacket(protocol.PROTO_NC_ACT_CHAT_REQ, append([]byte{0, byte(len(admCmd))}, []byte(admCmd)...)))
	if regular.Char.AdminLevel != 50 {
		t.Fatalf("expected NormalUser AdminLevel updated to 50 via &adminlevel, got %d", regular.Char.AdminLevel)
	}

	// Run &hide, &immortal, &killallmob
	for _, cmdStr := range []string{"&hide", "&immortal", "&killallmob"} {
		reqLv, exists := ds.AdminLvCatalog.RequiredLevel(strings.TrimPrefix(cmdStr, "&"))
		if exists && reqLv > gm.Char.AdminLevel {
			gm.Char.AdminLevel = reqLv
		}
		payload := append([]byte{0, byte(len(cmdStr))}, []byte(cmdStr)...)
		h.HandlePacket(gm.Client, buildTestPacket(protocol.PROTO_NC_ACT_CHAT_REQ, payload))
	}
	if !gm.IsGMHidden || !gm.IsGMImmortal {
		t.Fatalf("expected GM hidden & immortal flags enabled, got hide=%v immortal=%v", gm.IsGMHidden, gm.IsGMImmortal)
	}

	glL := waitForGameLogType(t, glRepo, func(e *database.GameLogEntry) bool { return e.CharNo == gm.Char.ID })
	if glL == nil {
		t.Fatalf("expected GameLog Type L entry persisted for GM action")
	}
}
