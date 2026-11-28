package zone

import (
	"testing"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestFase6_ReplyReviveAndItemUseActiveSkill(t *testing.T) {
	h, zMap, player, client := setupFase12TestEnvironment(t)
	h.dataStore = &DataStore{
		SkillCatalog: &SkillCatalog{
			ByID: map[uint16]*SkillDef{
				7001: {ID: 7001, MinWC: 60, MaxWC: 60, Range: 300},
			},
		},
	}

	// 1. REPLYREVIVE sin PendingRevive -> debe rechazar (0x1506) y mantener muerto
	player.MaxHP = 500
	player.State = 3
	player.Char.CurHP = 0
	player.PendingRevive = false
	revPkt := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_REPLYREVIVE_CMD)
	h.HandlePacket(client, protocol.NewPacket(revPkt.Data()))
	if player.State != 3 || player.Char.CurHP != 0 {
		t.Errorf("Sin PendingRevive el jugador no debería revivir, State=%d CurHP=%d", player.State, player.Char.CurHP)
	}

	// 2. REPLYREVIVE con PendingRevive -> revive con 100% MaxHP (rate=1000 en Zone.exe 0x00418ACD)
	player.PendingRevive = true
	h.HandlePacket(client, protocol.NewPacket(revPkt.Data()))
	if player.State != 1 || player.Char.CurHP != 500 {
		t.Errorf("ReplyRevive esperado State=1 CurHP=500 (100%%), obtenido State=%d CurHP=%d", player.State, player.Char.CurHP)
	}
	if player.PendingRevive {
		t.Errorf("PendingRevive debería limpiarse tras revivir")
	}

	// 3. ITEM_USE_ACTIVESKILL: ejecuta skill activa del ítem consumible (NO aprende la skill!)
	mob := NewZoneMonster(60, 102, "ScrollTarget", "Rou", 40, 0, 0, 5, 200, 50, 0, 0, 0, 0, 0, 0, 50)
	mob.IsKillableByUser = true
	zMap.AddMonster(mob)

	_ = h.charRepo.SaveCharacterItem(&database.Item{
		ItemKey:     9100,
		Owner:       player.Char.ID,
		StorageType: 9,
		Storage:     2,
		ItemID:      7001,
		Flags:       1,
	})

	usePkt := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_USE_ACTIVESKILL_REQ)
	usePkt.WriteUint8(2)           // slot
	usePkt.WriteUint16(mob.Handle) // target
	usePkt.WriteInt32(40)
	usePkt.WriteInt32(0)
	h.HandlePacket(client, protocol.NewPacket(usePkt.Data()))

	if mob.CurHP >= 200 {
		t.Errorf("ItemUseActiveSkill debería haber dañado al mob, CurHP=%d", mob.CurHP)
	}
	if player.HasSkill(7001) {
		t.Errorf("ItemUseActiveSkill NO debe aprender la skill 7001")
	}
	if it := h.findItemAt(player.Char.ID, 9, 2); it != nil {
		t.Errorf("ItemUseActiveSkill debería haber consumido el ítem en slot 2")
	}
}
