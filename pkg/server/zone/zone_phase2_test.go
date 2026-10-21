package zone

import (
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func waitOpCode(t *testing.T, conn net.Conn, op uint16) *protocol.Packet {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	for i := 0; i < 40; i++ {
		p, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("waiting 0x%04X: %v", op, err)
		}
		if p.OpCode == op {
			return p
		}
	}
	t.Fatalf("did not receive 0x%04X", op)
	return nil
}

func TestZoneServer_Phase2_InventoryEquipConsumablesAndSoulstones(t *testing.T) {
	classStore, _ := classparam.Load(filepath.Join("..", "..", "..", "data", "World", "ClassParam"))
	fieldStore, _ := field.Load(filepath.Join("..", "..", "..", "data", "World", "Map"))
	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	char := &database.Character{ID: 300, AccountID: 1, Name: "InvHero", Job: 1, Level: 1, TownName: "RouN", XPos: 6208, YPos: 6405, CurHP: 10, CurSP: 5}
	if err := charRepo.CreateCharacter(char); err != nil {
		t.Fatal(err)
	}
	_ = charRepo.CreateItem(&database.Item{StorageType: 9, Owner: char.ID, Storage: 0, ItemID: 1, Flags: 1})
	_ = charRepo.CreateItem(&database.Item{StorageType: 9, Owner: char.ID, Storage: 1, ItemID: 2, Flags: 10})

	srv := NewServer(&Config{BindAddr: "127.0.0.1:19220", Port: 19220, IP: "127.0.0.1"}, charRepo, transfers, classStore, fieldStore, nil)
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
	defer srv.Stop()
	transfers.CreateZoneTransfer(1, "U", char.Name, char.ID, 777, 0, "127.0.0.1")

	conn, err := net.DialTimeout("tcp", "127.0.0.1:19220", 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	hs, _ := readServerPacket(conn)
	xor, _ := hs.Reader.ReadInt16()
	nc, _ := crypto.NewNetCrypto(xor)

	login := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	login.WriteUint16(777)
	login.WriteString(char.Name, 16)
	_ = sendClientPacket(conn, login, nc)
	waitOpCode(t, conn, protocol.PROTO_NC_CHAR_MAPLOGIN_ACK)
	_ = sendClientPacket(conn, protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD), nc)

	// Equip (0x300F) -> 0x3011 err 0x0281
	eq := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_REQ)
	eq.WriteByte(0)
	_ = sendClientPacket(conn, eq, nc)
	if e, _ := waitOpCode(t, conn, protocol.PROTO_NC_ITEM_EQUIP_ACK).Reader.ReadUint16(); e != 0x0281 {
		t.Fatalf("equip err 0x%04X", e)
	}

	// Split (0x300D) 3 de slot 1 a slot 5 -> 0x300E err 0x0441
	sp := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SPLIT_REQ)
	sp.WriteUint16(9<<10 | 1)
	sp.WriteUint16(9<<10 | 5)
	sp.WriteByte(3)
	_ = sendClientPacket(conn, sp, nc)
	if e, _ := waitOpCode(t, conn, protocol.PROTO_NC_ITEM_SPLIT_ACK).Reader.ReadUint16(); e != 0x0441 {
		t.Fatalf("split err 0x%04X", e)
	}

	// Reloc (0x300B) slot 5 -> slot 1 (merge) -> 0x300C err 0x0241
	rl := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_REQ)
	rl.WriteUint16(9<<10 | 5)
	rl.WriteUint16(9<<10 | 7)
	_ = sendClientPacket(conn, rl, nc)
	if e, _ := waitOpCode(t, conn, protocol.PROTO_NC_ITEM_RELOC_ACK).Reader.ReadUint16(); e != 0x0241 {
		t.Fatalf("reloc err 0x%04X", e)
	}

	// Unequip (0x3012) slot 12 -> inven 0 -> 0x3011 err 0x0281
	uq := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_UNEQUIP_REQ)
	uq.WriteByte(12)
	uq.WriteByte(0)
	_ = sendClientPacket(conn, uq, nc)
	if e, _ := waitOpCode(t, conn, protocol.PROTO_NC_ITEM_EQUIP_ACK).Reader.ReadUint16(); e != 0x0281 {
		t.Fatalf("unequip err 0x%04X", e)
	}

	// Drop (0x3007) slot 7 -> 0x3008 err 0x0301, luego Pick (0x3009) -> 0x300A err 0x0341
	dr := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DROP_REQ)
	dr.WriteByte(7)
	dr.WriteByte(0)
	dr.WriteInt32(6208)
	dr.WriteInt32(6405)
	_ = sendClientPacket(conn, dr, nc)
	dropCmd := waitOpCode(t, conn, protocol.PROTO_NC_BRIEFINFO_DROPEDITEM_CMD)
	handle, _ := dropCmd.Reader.ReadUint16()
	if e, _ := waitOpCode(t, conn, protocol.PROTO_NC_ITEM_DROP_ACK).Reader.ReadUint16(); e != 0x0301 {
		t.Fatalf("drop err 0x%04X", e)
	}
	pk := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_REQ)
	pk.WriteUint16(handle)
	_ = sendClientPacket(conn, pk, nc)
	pa := waitOpCode(t, conn, protocol.PROTO_NC_ITEM_PICK_ACK)
	_, _ = pa.Reader.ReadUint16()
	_, _ = pa.Reader.ReadUint32()
	if e, _ := pa.Reader.ReadUint16(); e != 0x0341 {
		t.Fatalf("pick err 0x%04X", e)
	}

	// Soulstone HP use (0x5007) -> 0x5008 + 0x240E
	_ = sendClientPacket(conn, protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_USE_REQ), nc)
	waitOpCode(t, conn, protocol.PROTO_NC_SOULSTONE_HP_USESUC_ACK)
	waitOpCode(t, conn, protocol.PROTO_NC_BAT_HPCHANGE_CMD)

	// Soulstone SP use (0x5009) -> 0x500A + 0x240F
	_ = sendClientPacket(conn, protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_SP_USE_REQ), nc)
	waitOpCode(t, conn, protocol.PROTO_NC_SOULSTONE_SP_USESUC_ACK)
	waitOpCode(t, conn, protocol.PROTO_NC_BAT_SPCHANGE_CMD)

	// Auto arrange (0x304A) -> 0x304B err 0x30A8
	_ = sendClientPacket(conn, protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ), nc)
	if e, _ := waitOpCode(t, conn, protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK).Reader.ReadUint16(); e != 0x30A8 {
		t.Fatalf("arrange err 0x%04X", e)
	}
}
