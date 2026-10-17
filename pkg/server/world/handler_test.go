package world

import (
	"encoding/binary"
	"io"
	"net"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestWriteBasicCharInfo(t *testing.T) {
	ch := &database.Character{
		ID:        1,
		AccountID: 10,
		Name:      "PlayerOne",
		Slot:      0,
		Job:       byte(protocol.Job_Fighter),
		Level:     1,
		Male:      true,
		Hair:      1,
		HairColor: 2,
		Face:      3,
		TownName:  "Roumen",
		Equips:    map[byte]uint16{SlotWeapon: 250},
	}

	pw := protocol.New_PROTO_NC_CHAR_CREATE_ACK(func(pw *protocol.PacketWriter) {
		WriteBasicCharInfo(ch, pw)
	})

	wire := pw.BuildWirePacket()
	if len(wire) < 50 {
		t.Fatalf("expected packet to be at least 50 bytes, got %d", len(wire))
	}
}

func TestCharacterRepositoryOperations(t *testing.T) {
	repo := database.NewMemoryCharacterRepositoryWithFile("")

	ch1 := &database.Character{
		AccountID: 1,
		Name:      "Warrior",
		Slot:      0,
		Job:       byte(protocol.Job_Fighter),
		Level:     1,
		Male:      true,
	}

	err := repo.CreateCharacter(ch1)
	if err != nil {
		t.Fatalf("unexpected error creating char: %v", err)
	}

	// Duplicate name error
	chDup := &database.Character{
		AccountID: 2,
		Name:      "Warrior",
		Slot:      0,
		Job:       byte(protocol.Job_Archer),
	}
	if err := repo.CreateCharacter(chDup); err != database.ErrCharacterExists {
		t.Fatalf("expected ErrCharacterExists, got: %v", err)
	}

	// Slot collision error
	chSlot := &database.Character{
		AccountID: 1,
		Name:      "MageHero",
		Slot:      0,
		Job:       byte(protocol.Job_Mage),
	}
	if err := repo.CreateCharacter(chSlot); err != database.ErrSlotInUse {
		t.Fatalf("expected ErrSlotInUse, got: %v", err)
	}

	// Retrieve
	chars, err := repo.GetCharactersByAccount(1)
	if err != nil || len(chars) != 1 {
		t.Fatalf("expected 1 character, got %d (err: %v)", len(chars), err)
	}

	// Delete
	err = repo.DeleteCharacter(1, 0)
	if err != nil {
		t.Fatalf("unexpected error deleting character: %v", err)
	}

	charsAfter, _ := repo.GetCharactersByAccount(1)
	if len(charsAfter) != 0 {
		t.Fatalf("expected 0 characters after deletion, got %d", len(charsAfter))
	}
}

// readPacketFromConn lee y decodifica un paquete oficial de red desde el canal de pruebas.
func readPacketFromConn(r io.Reader) (uint16, []byte, error) {
	lenBuf := make([]byte, 1)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return 0, nil, err
	}

	packetLen := int(lenBuf[0])
	if packetLen == 0xFF {
		extLen := make([]byte, 2)
		if _, err := io.ReadFull(r, extLen); err != nil {
			return 0, nil, err
		}
		packetLen = int(binary.LittleEndian.Uint16(extLen))
	}

	raw := make([]byte, packetLen)
	if _, err := io.ReadFull(r, raw); err != nil {
		return 0, nil, err
	}

	opCode := binary.LittleEndian.Uint16(raw[0:2])
	payload := raw[2:]
	return opCode, payload, nil
}

func setupTestHandler() (*Handler, *database.MemoryCharacterRepository, *database.TransferManager) {
	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	badNames := &BadNameFilter{
		exactBadNames: map[string]bool{"forbidden": true},
		subStrings:    []string{"admin", "gm"},
	}

	defaultData := &DefaultCharacterDataStore{
		Classes: map[string]DefaultClassData{
			"1": {
				ClassID:   1,
				ClassName: "Fighter",
				StartMap:  "RouN",
				StartX:    6888,
				StartY:    7045,
				HP:        52,
				SP:        10,
				Money:     100,
			},
		},
	}

	dataStore := &WorldDataStore{
		BadNames:    badNames,
		DefaultData: defaultData,
	}

	h := NewHandler(charRepo, transfers, dataStore, "127.0.0.1", 9210)
	return h, charRepo, transfers
}

func TestHandler_GameTime(t *testing.T) {
	h, _, _ := setupTestHandler()
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client, err := network.NewClient(serverConn, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	pkt := &protocol.Packet{
		OpCode: protocol.PROTO_NC_MISC_GAMETIME_REQ,
		Header: protocol.ExtractHeader(protocol.PROTO_NC_MISC_GAMETIME_REQ),
		Type:   protocol.ExtractType(protocol.PROTO_NC_MISC_GAMETIME_REQ),
	}

	go h.HandlePacket(client, pkt)

	opCode, payload, err := readPacketFromConn(clientConn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if opCode != protocol.PROTO_NC_MISC_GAMETIME_ACK {
		t.Fatalf("expected OpCode 0x%04X, got 0x%04X", protocol.PROTO_NC_MISC_GAMETIME_ACK, opCode)
	}

	if len(payload) != 3 {
		t.Fatalf("expected 3 bytes payload (hour, minute, second), got %d", len(payload))
	}
}

func TestHandler_CreateCharacter_BadNameRejected(t *testing.T) {
	h, _, _ := setupTestHandler()
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client, err := network.NewClient(serverConn, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.AccountID = 100
	client.Username = "TestAccount"
	client.IsAuthenticated = true

	// Construir payload con nombre prohibido "AdminGuy"
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CREATE_REQ)
	pw.WriteByte(0) // Slot 0
	pw.WriteString("AdminGuy", 20)
	pw.WriteByte(byte(1 << 2)) // Job 1 (Fighter)
	pw.WriteByte(1)            // hair
	pw.WriteByte(1)            // color
	pw.WriteByte(1)            // face

	wire := pw.BuildWirePacket()
	pkt := protocol.NewPacket(wire[1:])
	if pkt == nil {
		t.Fatalf("failed to parse packet")
	}

	go h.HandlePacket(client, pkt)

	opCode, payload, err := readPacketFromConn(clientConn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if opCode != protocol.PROTO_NC_CHAR_CREATE_FAIL_ACK {
		t.Fatalf("expected PROTO_NC_CHAR_CREATE_FAIL_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_CHAR_CREATE_FAIL_ACK, opCode)
	}

	errCode := binary.LittleEndian.Uint16(payload[0:2])
	if errCode != uint16(protocol.CreateCharError_FailedToCreate) {
		t.Fatalf("expected CreateCharError_FailedToCreate, got %d", errCode)
	}
}

func TestHandler_CreateCharacter_Success(t *testing.T) {
	h, repo, _ := setupTestHandler()
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client, err := network.NewClient(serverConn, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.AccountID = 100
	client.Username = "TestAccount"
	client.IsAuthenticated = true

	// Construir payload válido con "Conan"
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CREATE_REQ)
	pw.WriteByte(0) // Slot 0
	pw.WriteString("Conan", 20)
	pw.WriteByte(byte((1 << 2) | (1 << 7))) // Job 1 (Fighter), Male
	pw.WriteByte(2)                         // hair
	pw.WriteByte(1)                         // color
	pw.WriteByte(3)                         // face

	wire := pw.BuildWirePacket()
	pkt := protocol.NewPacket(wire[1:])
	if pkt == nil {
		t.Fatalf("failed to parse packet")
	}

	go h.HandlePacket(client, pkt)

	opCode, _, err := readPacketFromConn(clientConn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if opCode != protocol.PROTO_NC_CHAR_CREATE_ACK {
		t.Fatalf("expected PROTO_NC_CHAR_CREATE_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_CHAR_CREATE_ACK, opCode)
	}

	// Verificar que el personaje fue persistido correctamente
	created, err := repo.GetCharacterByName("Conan")
	if err != nil || created == nil {
		t.Fatalf("expected character Conan in repo, got err: %v", err)
	}

	if created.Slot != 0 || created.Job != 1 || !created.Male || created.CurHP != 52 {
		t.Fatalf("character attributes mismatch: %+v", created)
	}

	// Verificar arma inicial de guerrero
	if created.GetEquip(SlotWeapon) != 250 {
		t.Fatalf("expected sword 250 in SlotWeapon, got %d", created.GetEquip(SlotWeapon))
	}
}

func TestHandler_CharSelect_ZoneHandover(t *testing.T) {
	h, repo, transfers := setupTestHandler()
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client, err := network.NewClient(serverConn, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.AccountID = 100
	client.Username = "TestAccount"
	client.IsAuthenticated = true

	// Crear personaje en slot 0
	ch := &database.Character{
		AccountID: 100,
		Name:      "Legolas",
		Slot:      0,
		Job:       byte(protocol.Job_Archer),
		Level:     1,
		Male:      true,
	}
	_ = repo.CreateCharacter(ch)

	// Enviar selección de personaje en slot 0 (PROTO_NC_CHAR_LOGIN_REQ)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	pw.WriteByte(0)

	wire := pw.BuildWirePacket()
	pkt := protocol.NewPacket(wire[1:])
	if pkt == nil {
		t.Fatalf("failed to parse packet")
	}

	go h.HandlePacket(client, pkt)

	opCode, payload, err := readPacketFromConn(clientConn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if opCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("expected PROTO_NC_CHAR_LOGIN_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_CHAR_LOGIN_ACK, opCode)
	}

	// PROTO_NC_CHAR_LOGIN_ACK oficial: 18 bytes (16 bytes IP + 2 bytes Port)
	if len(payload) != 18 {
		t.Fatalf("expected 18 bytes payload for PROTO_NC_CHAR_LOGIN_ACK, got %d", len(payload))
	}

	zonePort := binary.LittleEndian.Uint16(payload[16:18])
	if zonePort != 9210 {
		t.Fatalf("expected ZonePort 9210, got %d", zonePort)
	}

	// Verificar que el token de transferencia a Zone fue generado y es consumible
	transfer, ok := transfers.ConsumeZoneTransfer("Legolas", 0, nil)
	if !ok || transfer == nil {
		t.Fatalf("expected active zone transfer in TransferManager")
	}
	if transfer.CharName != "Legolas" || transfer.AccountID != 100 {
		t.Fatalf("zone transfer data mismatch: %+v", transfer)
	}
}

func TestHandler_DeleteCharacter_Flow(t *testing.T) {
	h, repo, _ := setupTestHandler()
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client, err := network.NewClient(serverConn, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.AccountID = 100
	client.Username = "TestAccount"
	client.IsAuthenticated = true

	// Crear personaje en slot 1
	ch := &database.Character{
		AccountID: 100,
		Name:      "Gandalf",
		Slot:      1,
		Job:       byte(protocol.Job_Mage),
		Level:     1,
		Male:      true,
	}
	_ = repo.CreateCharacter(ch)

	// Enviar borrado de slot 1 (PROTO_NC_CHAR_DELETE_REQ)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_DELETE_REQ)
	pw.WriteByte(1)

	wire := pw.BuildWirePacket()
	pkt := protocol.NewPacket(wire[1:])
	if pkt == nil {
		t.Fatalf("failed to parse packet")
	}

	go h.HandlePacket(client, pkt)

	opCode, payload, err := readPacketFromConn(clientConn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if opCode != protocol.PROTO_NC_CHAR_DELETE_ACK {
		t.Fatalf("expected PROTO_NC_CHAR_DELETE_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_CHAR_DELETE_ACK, opCode)
	}

	if payload[0] != 1 {
		t.Fatalf("expected slot 1 in delete response, got %d", payload[0])
	}

	// Verificar que ya no está activo
	chars, _ := repo.GetCharactersByAccount(100)
	if len(chars) != 0 {
		t.Fatalf("expected 0 characters after deletion, got %d", len(chars))
	}
}

func TestHandler_BackToWorldSelect(t *testing.T) {
	h, _, _ := setupTestHandler()
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client, err := network.NewClient(serverConn, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	pkt := &protocol.Packet{
		OpCode: protocol.PROTO_NC_USER_BACKTOWORLDLIST_REQ,
		Header: protocol.ExtractHeader(protocol.PROTO_NC_USER_BACKTOWORLDLIST_REQ),
		Type:   protocol.ExtractType(protocol.PROTO_NC_USER_BACKTOWORLDLIST_REQ),
	}

	go h.HandlePacket(client, pkt)

	opCode, _, err := readPacketFromConn(clientConn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if opCode != protocol.PROTO_NC_USER_BACKTOWORLDLIST_ACK {
		t.Fatalf("expected PROTO_NC_USER_BACKTOWORLDLIST_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_USER_BACKTOWORLDLIST_ACK, opCode)
	}
}
