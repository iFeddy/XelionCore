package zone

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestCanEquipItemChecksOfficialClassAndLevel(t *testing.T) {
	store := loadOfficialItemStore(t)

	handler := &Handler{itemStore: store}
	player := &ZonePlayer{Char: &database.Character{Job: byte(protocol.Job_Fighter), Level: 2}}
	leatherBoots := &database.Item{ItemID: 0}
	if !handler.canEquipItem(player, leatherBoots) {
		t.Fatal("Fighter at the required level should equip LeatherBoots")
	}

	player.Char.Level = 1
	if handler.canEquipItem(player, leatherBoots) {
		t.Fatal("character below DemandLv should not equip LeatherBoots")
	}

	player.Char.Level = 2
	player.Char.Job = byte(protocol.Job_Cleric)
	if handler.canEquipItem(player, leatherBoots) {
		t.Fatal("class excluded by UseClassTypeInfo should not equip LeatherBoots")
	}
}

func TestSortInventoryItemsUsesOfficialCategoryOrder(t *testing.T) {
	store := loadOfficialItemStore(t)
	boots, ok := store.GetItemInfo(0)
	if !ok || boots.SortOrder < 0 {
		t.Fatal("LeatherBoots is missing its official ItemSort order")
	}
	helmet, ok := store.GetItemInfo(1)
	if !ok || helmet.SortOrder < 0 {
		t.Fatal("LeatherHelmet is missing its official ItemSort order")
	}
	if boots.SortOrder == helmet.SortOrder {
		t.Fatal("test items should belong to different official sort categories")
	}

	items := []*database.Item{
		{ItemKey: 1, ItemID: 1, Storage: 0},
		{ItemKey: 2, ItemID: 0, Storage: 1},
	}
	sortInventoryItems(items, store)
	if boots.SortOrder < helmet.SortOrder && items[0].ItemID != boots.ID {
		t.Fatal("inventory did not put the earlier official category first")
	}
	if helmet.SortOrder < boots.SortOrder && items[0].ItemID != helmet.ID {
		t.Fatal("inventory did not put the earlier official category first")
	}
}

func loadOfficialItemStore(t *testing.T) *item.Store {
	t.Helper()
	dataDir := filepath.Join("..", "..", "..", "data")
	store := item.NewStore()
	if err := store.LoadItemInfoSHN(filepath.Join(dataDir, "ItemInfo.shn")); err != nil {
		t.Fatal(err)
	}
	if err := store.LoadInventoryDataSHN(
		filepath.Join(dataDir, "ItemInfoServer.shn"),
		filepath.Join(dataDir, "ItemSort.shn"),
		filepath.Join(dataDir, "UseClassTypeInfo.shn"),
		filepath.Join(dataDir, "ClassName.shn"),
	); err != nil {
		t.Fatal(err)
	}
	return store
}
