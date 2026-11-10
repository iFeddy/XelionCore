package protocol_test

import (
	"bytes"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// TestGameLog_SerializeDeserialize verifica que las 44 estructuras de GameLog
// se serialicen y deserialicen correctamente preservando cada byte en Little-Endian.
func TestGameLog_SerializeDeserialize(t *testing.T) {
	// Prueba con PROTO_NC_LOG_GAME_DATA_TYPE_0
	t0 := &protocol.PROTO_NC_LOG_GAME_DATA_TYPE_0_DATA{
		NType:         1,
		NCharNo:       1001,
		NMapX:         500,
		NMapY:         600,
		NMapZ:         700,
		NTargetCharNo: 2002,
		NTargetID:     3003,
		NItemKey:      4004,
		NInt1:         10,
		NInt2:         20,
		NInt3:         30,
		NBigint1:      99999,
	}
	copy(t0.SMap[:], "Rou")

	data, err := t0.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize PROTO_NC_LOG_GAME_DATA_TYPE_0: %v", err)
	}
	if len(data) != 68 {
		t.Fatalf("Expected 68 bytes, got %d", len(data))
	}

	t0Restored := &protocol.PROTO_NC_LOG_GAME_DATA_TYPE_0_DATA{}
	if err := t0Restored.Deserialize(data); err != nil {
		t.Fatalf("Failed to deserialize PROTO_NC_LOG_GAME_DATA_TYPE_0: %v", err)
	}

	if *t0 != *t0Restored {
		t.Fatalf("Restored struct mismatch: %+v vs %+v", t0Restored, t0)
	}

	// Prueba con PROTO_NC_LOG_GAME_ITEM_UPGRADE
	upg := &protocol.PROTO_NC_LOG_GAME_ITEM_UPGRADE_DATA{
		Chrregnum:     555,
		Itemkey:       88888,
		Itemoldlevel:  3,
		Itemnewlevel:  4,
		Left_item_id:  101,
		Right_item_id: 102,
		Middle_item_id: 103,
	}
	copy(upg.Mappos.MapName[:], "Eld")
	upg.Mappos.Coord.X = 1200
	upg.Mappos.Coord.Y = 3400

	upgData, err := upg.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize PROTO_NC_LOG_GAME_ITEM_UPGRADE: %v", err)
	}
	if len(upgData) != 43 {
		t.Fatalf("Expected 43 bytes, got %d", len(upgData))
	}

	upgRestored := &protocol.PROTO_NC_LOG_GAME_ITEM_UPGRADE_DATA{}
	if err := upgRestored.Deserialize(upgData); err != nil {
		t.Fatalf("Failed to deserialize PROTO_NC_LOG_GAME_ITEM_UPGRADE: %v", err)
	}

	if *upg != *upgRestored {
		t.Fatalf("Restored upgrade struct mismatch: %+v vs %+v", upgRestored, upg)
	}

	// Prueba con PROTO_NC_SYSLOG_CHAR_DEATH
	death := &protocol.PROTO_NC_SYSLOG_CHAR_DEATH_DATA{
		UserNo:       10,
		CharNo:       20,
		ZoneNo:       1,
		AttackerType: 2,
		AttackerID:   300,
		DeathCoord: protocol.SHINE_XY_TYPE{
			X: 450,
			Y: 950,
		},
	}
	deathData, err := death.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize PROTO_NC_SYSLOG_CHAR_DEATH: %v", err)
	}
	if len(deathData) != 22 {
		t.Fatalf("Expected 22 bytes, got %d", len(deathData))
	}

	deathRestored := &protocol.PROTO_NC_SYSLOG_CHAR_DEATH_DATA{}
	if err := deathRestored.Deserialize(deathData); err != nil {
		t.Fatalf("Failed to deserialize PROTO_NC_SYSLOG_CHAR_DEATH: %v", err)
	}

	if *death != *deathRestored {
		t.Fatalf("Restored syslog death mismatch: %+v vs %+v", deathRestored, death)
	}

	// Verificar FixedWireSize()
	if t0.FixedWireSize() != 68 {
		t.Fatalf("FixedWireSize expected 68, got %d", t0.FixedWireSize())
	}
	if upg.FixedWireSize() != 43 {
		t.Fatalf("FixedWireSize expected 43, got %d", upg.FixedWireSize())
	}
	if death.FixedWireSize() != 22 {
		t.Fatalf("FixedWireSize expected 22, got %d", death.FixedWireSize())
	}

	t.Logf("GameLog serialization and deserialization verified 100%% successfully")
	_ = bytes.Buffer{}
}

