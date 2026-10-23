package zone

import (
	"encoding/binary"
	"io"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/npc"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func TestHandler_NormalLogoutReturnsToWorldSelect(t *testing.T) {
	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	char := &database.Character{
		ID:        23,
		AccountID: 1,
		Name:      "LogoutHero",
		TownName:  "RouN",
		CurHP:     52,
		CurSP:     10,
	}
	if err := charRepo.CreateCharacter(char); err != nil {
		t.Fatalf("failed to create character: %v", err)
	}

	h := NewHandler(charRepo, nil, nil, nil, nil)
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	client, err := network.NewClient(serverConn, nil)
	if err != nil {
		t.Fatalf("failed to create network client: %v", err)
	}
	h.players[client] = &ZonePlayer{Client: client, Char: char}

	req := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NORMALLOGOUT_CMD)
	req.WriteByte(0)
	pkt := protocol.NewPacket(req.Data())

	done := make(chan struct{})
	go func() {
		h.HandlePacket(client, pkt)
		close(done)
	}()

	resp, err := readServerPacket(clientConn)
	if err != nil {
		t.Fatalf("failed to read logout response: %v", err)
	}
	if resp.OpCode != protocol.PROTO_NC_USER_WILL_WORLD_SELECT_ACK {
		t.Fatalf("expected PROTO_NC_USER_WILL_WORLD_SELECT_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_USER_WILL_WORLD_SELECT_ACK, resp.OpCode)
	}
	if resp.Reader.Remaining() != 34 {
		t.Fatalf("expected 34-byte logout response payload, got %d bytes", resp.Reader.Remaining())
	}
	errCode, err := resp.Reader.ReadUint16()
	if err != nil {
		t.Fatalf("failed to read world-select response error code: %v", err)
	}
	if errCode != 0 {
		t.Fatalf("expected successful world-select response, got error code %d", errCode)
	}
	otp, err := resp.Reader.ReadBytes(32)
	if err != nil {
		t.Fatalf("failed to read world-select OTP: %v", err)
	}
	for i, b := range otp {
		if b != 0 {
			t.Fatalf("expected empty OTP, byte %d was 0x%02X", i, b)
		}
	}

	<-done
	if _, ok := h.getPlayer(client); ok {
		t.Fatal("expected logged-out character to be removed from active players")
	}
}

func readServerPacket(r io.Reader) (*protocol.Packet, error) {
	lenBuf := make([]byte, 1)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return nil, err
	}

	length := int(lenBuf[0])
	if length == 0 {
		extBuf := make([]byte, 2)
		if _, err := io.ReadFull(r, extBuf); err != nil {
			return nil, err
		}
		length = int(binary.LittleEndian.Uint16(extBuf))
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	return protocol.NewPacket(payload), nil
}

func sendClientPacket(w io.Writer, pw *protocol.PacketWriter, nc *crypto.NetCrypto) error {
	bufCopy := make([]byte, len(pw.Data()))
	copy(bufCopy, pw.Data())
	if nc != nil {
		nc.Crypt(bufCopy)
	}

	payloadLen := len(bufCopy)
	var wire []byte
	if payloadLen <= 255 {
		wire = make([]byte, 1+payloadLen)
		wire[0] = byte(payloadLen)
		copy(wire[1:], bufCopy)
	} else {
		wire = make([]byte, 3+payloadLen)
		wire[0] = 0x00
		binary.LittleEndian.PutUint16(wire[1:3], uint16(payloadLen))
		copy(wire[3:], bufCopy)
	}

	_, err := w.Write(wire)
	return err
}

func TestZoneServer_Phase1_LifecycleAndHandshake(t *testing.T) {
	classDir := filepath.Join("..", "..", "..", "data", "World", "ClassParam")
	mapDir := filepath.Join("..", "..", "..", "data", "World", "Map")

	classStore, err := classparam.Load(classDir)
	if err != nil {
		t.Fatalf("failed to load classparam: %v", err)
	}

	fieldStore, err := field.Load(mapDir)
	if err != nil {
		t.Fatalf("failed to load field: %v", err)
	}

	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	cfg := &Config{
		BindAddr: "127.0.0.1:19210",
		Port:     19210,
		IP:       "127.0.0.1",
		MapID:    0,
	}

	srv := NewServer(cfg, charRepo, transfers, classStore, fieldStore, nil)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start zone server: %v", err)
	}
	defer srv.Stop()

	// 1. Conectar cliente TCP
	conn, err := net.DialTimeout("tcp", "127.0.0.1:19210", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect to zone server: %v", err)
	}
	defer conn.Close()

	// 2. Leer Handshake de semilla XOR (PROTO_NC_MISC_SEED_CMD, OpCode 0x0807)
	pkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read handshake packet: %v", err)
	}

	if pkt.OpCode != protocol.PROTO_NC_MISC_SEED_CMD {
		t.Fatalf("expected PROTO_NC_MISC_SEED_CMD (0x0807), got 0x%04X", pkt.OpCode)
	}

	xorPos, err := pkt.Reader.ReadInt16()
	if err != nil {
		t.Fatalf("failed to read xor offset: %v", err)
	}
	if xorPos < 0 {
		t.Fatalf("invalid negative xor offset: %d", xorPos)
	}
}

func TestZoneServer_Phase1_TransferKeyVerification(t *testing.T) {
	classDir := filepath.Join("..", "..", "..", "data", "World", "ClassParam")
	mapDir := filepath.Join("..", "..", "..", "data", "World", "Map")

	classStore, _ := classparam.Load(classDir)
	fieldStore, _ := field.Load(mapDir)
	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	cfg := &Config{
		BindAddr: "127.0.0.1:19211",
		Port:     19211,
		IP:       "127.0.0.1",
		MapID:    0,
	}

	srv := NewServer(cfg, charRepo, transfers, classStore, fieldStore, nil)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start zone server: %v", err)
	}
	defer srv.Stop()

	// Crear personaje en el repo
	char := &database.Character{
		ID:        1,
		AccountID: 1,
		Name:      "ZoneHero",
		Slot:      0,
		Job:       byte(protocol.Job_Fighter),
		Level:     1,
		Male:      true,
		XPos:      5600,
		YPos:      8800,
		CurHP:     52,
		CurSP:     10,
	}
	if err := charRepo.CreateCharacter(char); err != nil {
		t.Fatalf("failed to create test char: %v", err)
	}

	// Registrar transferencia válida desde WorldService
	transfers.CreateZoneTransfer(1, "TestUser", char.Name, char.ID, 12345, 0, "127.0.0.1")

	// Conectar cliente
	conn, err := net.DialTimeout("tcp", "127.0.0.1:19211", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Leer handshake inicial y configurar NetCrypto
	hsPkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read handshake: %v", err)
	}
	xorPos, _ := hsPkt.Reader.ReadInt16()
	clientCrypto, err := crypto.NewNetCrypto(xorPos)
	if err != nil {
		t.Fatalf("failed to create client crypto: %v", err)
	}

	// Enviar PROTO_NC_MAP_LOGIN_REQ con credenciales válidas
	req := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	req.WriteUint16(12345) // randomID
	req.WriteString(char.Name, 16)

	if err := sendClientPacket(conn, req, clientCrypto); err != nil {
		t.Fatalf("failed to send PROTO_NC_MAP_LOGIN_REQ: %v", err)
	}

	// Leer secuencia completa de paquetes del InGame chunk para personaje limpio (sin skills ni ítems en cajas 8/9/12):
	expectedOpCodes := []uint16{
		protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD,          // 1. 0x1038
		protocol.PROTO_NC_CHAR_CLIENT_SHAPE_CMD,         // 2. 0x1039
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD,   // 3. 0x103A
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD,    // 4. 0x103B
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD,    // 5. 0x10CE
		protocol.PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD,  // 6. 0x10D7
		protocol.PROTO_NC_CHAR_CLIENT_PASSIVE_CMD,       // 7. 0x103E
		protocol.PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD,   // 8. 0x104A
		protocol.PROTO_NC_CHAR_CLIENT_GAME_CMD,          // 10. 0x1048 (0x103D omitido si len(skills)==0)
		protocol.PROTO_NC_CHAR_CLIENT_ITEM_CMD,          // 14. 0x1047 (Solo Box 15 vacía; 8, 9 y 12 omitidas si vacías)
		protocol.PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD,     // 15. 0x1049
		protocol.PROTO_NC_CHAR_MAPLOGIN_ACK,             // 17. 0x1802 (Último paquete del bloque)
	}

	for idx, expOp := range expectedOpCodes {
		respPkt, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("step %d: failed to read packet (expected 0x%04X): %v", idx+1, expOp, err)
		}
		if respPkt.OpCode != expOp {
			t.Fatalf("step %d: expected OpCode 0x%04X, got 0x%04X", idx+1, expOp, respPkt.OpCode)
		}
	}
}

func TestZoneServer_Phase1_InvalidTransferRejection(t *testing.T) {
	classDir := filepath.Join("..", "..", "..", "data", "World", "ClassParam")
	mapDir := filepath.Join("..", "..", "..", "data", "World", "Map")

	classStore, _ := classparam.Load(classDir)
	fieldStore, _ := field.Load(mapDir)
	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	cfg := &Config{
		BindAddr: "127.0.0.1:19212",
		Port:     19212,
		IP:       "127.0.0.1",
		MapID:    0,
	}

	srv := NewServer(cfg, charRepo, transfers, classStore, fieldStore, nil)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start zone server: %v", err)
	}
	defer srv.Stop()

	conn, err := net.DialTimeout("tcp", "127.0.0.1:19212", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Descartar handshake y configurar NetCrypto
	hsPkt, _ := readServerPacket(conn)
	xorPos, _ := hsPkt.Reader.ReadInt16()
	clientCrypto, _ := crypto.NewNetCrypto(xorPos)

	// Enviar transferencia inexistente
	req := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	req.WriteUint16(9999) // RandomID inválido
	req.WriteString("GhostChar", 16)

	if err := sendClientPacket(conn, req, clientCrypto); err != nil {
		t.Fatalf("failed to send packet: %v", err)
	}

	respPkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read rejection: %v", err)
	}

	if respPkt.OpCode != protocol.PROTO_NC_MAP_LOGINFAIL_ACK {
		t.Errorf("expected PROTO_NC_MAP_LOGINFAIL_ACK (0x%04X), got 0x%04X", protocol.PROTO_NC_MAP_LOGINFAIL_ACK, respPkt.OpCode)
	}
}

func TestZoneServer_Phase1_ChecksumVerification(t *testing.T) {
	classDir := filepath.Join("..", "..", "..", "data", "World", "ClassParam")
	mapDir := filepath.Join("..", "..", "..", "data", "World", "Map")
	dataDir := filepath.Join("..", "..", "..", "data")
	viewDir := filepath.Join(dataDir, "View")

	classStore, _ := classparam.Load(classDir)
	fieldStore, _ := field.Load(mapDir)
	chkEngine, err := shn.NewChecksumEngine(dataDir, viewDir)
	if err != nil {
		t.Fatalf("failed to init ChecksumEngine: %v", err)
	}

	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	cfg := &Config{
		BindAddr: "127.0.0.1:19213",
		Port:     19213,
		IP:       "127.0.0.1",
		MapID:    0,
	}

	srv := NewServer(cfg, charRepo, transfers, classStore, fieldStore, chkEngine)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start zone server: %v", err)
	}
	defer srv.Stop()

	char := &database.Character{
		ID:        100,
		AccountID: 1,
		Name:      "ChecksumTester",
		Job:       1,
		Level:     1,
		Male:      true,
		XPos:      5600,
		YPos:      8800,
		CurHP:     52,
		CurSP:     10,
	}
	_ = charRepo.CreateCharacter(char)
	transfers.CreateZoneTransfer(1, "TestUser", char.Name, char.ID, 54321, 0, "127.0.0.1")

	// Conectar cliente y enviar checksum corrupto en el índice 2 (CharacterTitleData.shn)
	conn, err := net.DialTimeout("tcp", "127.0.0.1:19213", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	hsPkt, _ := readServerPacket(conn)
	xorPos, _ := hsPkt.Reader.ReadInt16()
	clientCrypto, _ := crypto.NewNetCrypto(xorPos)

	req := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	req.WriteUint16(54321) // RandomID
	req.WriteString(char.Name, 20)
	checksums := chkEngine.GetChecksums()
	checksums[2][0] ^= 0xFF // Corromper índice 2 (CharacterTitleData.shn)
	for i := 0; i < 49; i++ {
		req.WriteBytes(checksums[i][:])
	}

	if err := sendClientPacket(conn, req, clientCrypto); err != nil {
		t.Fatalf("failed to send packet: %v", err)
	}

	respPkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	if respPkt.OpCode != protocol.PROTO_NC_MAP_LOGINFAIL_ACK {
		t.Fatalf("expected PROTO_NC_MAP_LOGINFAIL_ACK (0x%04X), got 0x%04X", protocol.PROTO_NC_MAP_LOGINFAIL_ACK, respPkt.OpCode)
	}

	errCode, _ := respPkt.Reader.ReadUint16()
	wrongIdx, _ := respPkt.Reader.ReadByte()
	if errCode != 0x0147 {
		t.Errorf("expected errCode 0x0147, got 0x%04X", errCode)
	}
	if wrongIdx != 2 {
		t.Errorf("expected wrongIdx 2, got %d", wrongIdx)
	}
}

func TestZoneServer_ReportsAndLogoutLog(t *testing.T) {
	cfg := &Config{
		BindAddr: "127.0.0.1:19213",
		Port:     19213,
		IP:       "127.0.0.1",
		MapID:    0,
	}

	charRepo := database.NewMemoryCharacterRepository()
	transfers := database.NewTransferManager(nil)
	classStore, _ := classparam.Load(filepath.Join("..", "..", "data", "World", "ClassParam"))
	fieldStore, _ := field.Load(filepath.Join("..", "..", "data", "World", "Map"))

	srv := NewServer(cfg, charRepo, transfers, classStore, fieldStore, nil)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start zone server: %v", err)
	}
	defer srv.Stop()

	char := &database.Character{
		ID:        201,
		AccountID: 10,
		Name:      "ReporterHero",
		Slot:      0,
		Job:       byte(protocol.Job_Fighter),
		Level:     5,
		Male:      true,
		XPos:      5600,
		YPos:      8800,
		CurHP:     100,
		CurSP:     50,
	}
	_ = charRepo.CreateCharacter(char)
	transfers.CreateZoneTransfer(10, "User10", char.Name, char.ID, 33445, 0, "127.0.0.1")

	conn, err := net.DialTimeout("tcp", "127.0.0.1:19213", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	hsPkt, _ := readServerPacket(conn)
	xorPos, _ := hsPkt.Reader.ReadInt16()
	clientCrypto, _ := crypto.NewNetCrypto(xorPos)

	// Login to Zone
	req := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	req.WriteUint16(33445)
	req.WriteString(char.Name, 16)
	_ = sendClientPacket(conn, req, clientCrypto)

	// Consumir paquetes iniciales (LoadMap, LoginComplete)
	_, _ = readServerPacket(conn) // LoadMap
	_, _ = readServerPacket(conn) // LoginComplete

	// 1. Probar PROTO_NC_DATA_REPORT_ADD_REQ (0x8007)
	repPkt := protocol.NewPacketWriter(protocol.PROTO_NC_DATA_REPORT_ADD_REQ)
	repPkt.WriteUint32(10) // dwReporter
	repPkt.WriteUint16(1)  // hReporter
	repPkt.WriteByte(0)    // nWorldNo
	repPkt.WriteString(char.Name, 20)
	repPkt.WriteString("Offender", 20)
	repPkt.WriteString("Abuse", 16)
	repPkt.WriteString("Harassment in town", 256)
	repPkt.WriteUint16(12)
	repPkt.WriteString("chat message", 12)

	if err := sendClientPacket(conn, repPkt, clientCrypto); err != nil {
		t.Fatalf("Failed to send report add: %v", err)
	}

	var ackPkt *protocol.Packet
	for i := 0; i < 20; i++ {
		p, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("Failed to read server packet: %v", err)
		}
		if p.OpCode == protocol.PROTO_NC_DATA_REPORT_ADD_ACK {
			ackPkt = p
			break
		}
	}
	if ackPkt == nil {
		t.Fatalf("Did not receive PROTO_NC_DATA_REPORT_ADD_ACK")
	}
	t.Logf("PROTO_NC_DATA_REPORT_ADD_ACK verificado exitosamente (0x%04X)", ackPkt.OpCode)

	// 2. Probar PROTO_NC_DATA_REPORT_CANCEL_REQ (0x8009)
	cancelPkt := protocol.NewPacketWriter(protocol.PROTO_NC_DATA_REPORT_CANCEL_REQ)
	cancelPkt.WriteInt32(100) // RegNumber
	cancelPkt.WriteUint16(1)  // hReporter
	if err := sendClientPacket(conn, cancelPkt, clientCrypto); err != nil {
		t.Fatalf("Failed to send report cancel: %v", err)
	}

	var cancelAck *protocol.Packet
	for i := 0; i < 20; i++ {
		p, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("Failed to read server packet: %v", err)
		}
		if p.OpCode == protocol.PROTO_NC_DATA_REPORT_CANCEL_ACK {
			cancelAck = p
			break
		}
	}
	if cancelAck == nil {
		t.Fatalf("Did not receive PROTO_NC_DATA_REPORT_CANCEL_ACK")
	}
	t.Logf("PROTO_NC_DATA_REPORT_CANCEL_ACK verificado exitosamente (0x%04X)", cancelAck.OpCode)

	// 3. Probar desconexión y conteo de jugadores activos
	activeBefore := srv.Handler().ActivePlayersCount()
	if activeBefore == 0 {
		t.Errorf("Expected at least 1 active player before disconnect")
	}

	conn.Close()
	time.Sleep(50 * time.Millisecond)

	activeAfter := srv.Handler().ActivePlayersCount()
	if activeAfter != 0 {
		t.Errorf("Expected 0 active players after disconnect, got %d", activeAfter)
	}
	t.Logf("Desconexión limpia y actualización de CCU verificada")
}

func TestZoneServer_Phase1_NPCSpawnsAndGateWarps(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	classDir := filepath.Join(dataDir, "World", "ClassParam")
	mapDir := filepath.Join(dataDir, "World", "Map")
	npcDir := filepath.Join(dataDir, "World", "NPC")

	classStore, err := classparam.Load(classDir)
	if err != nil {
		t.Fatalf("failed to load classparam: %v", err)
	}

	fieldStore, err := field.Load(mapDir)
	if err != nil {
		t.Fatalf("failed to load field: %v", err)
	}

	npcStore, err := npc.Load(npcDir)
	if err != nil {
		t.Fatalf("failed to load npc: %v", err)
	}

	mobInfoSHN, err := shn.ReadFile(filepath.Join(dataDir, "MobInfo.shn"))
	if err != nil {
		t.Fatalf("failed to load MobInfo.shn: %v", err)
	}

	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	transfers := database.NewTransferManager()

	testChar := &database.Character{
		ID:        99,
		AccountID: 1,
		Name:      "WarpHero",
		Slot:      0,
		Job:       1,
		Level:     1,
		TownName:  "RouN",
		XPos:      6208,
		YPos:      6405,
		CurHP:     100,
		CurSP:     50,
	}
	if err := charRepo.CreateCharacter(testChar); err != nil {
		t.Fatalf("failed to create test char: %v", err)
	}

	cfg := &Config{
		BindAddr: "127.0.0.1:19214",
		Port:     19214,
		IP:       "127.0.0.1",
		MapID:    0,
	}

	srv := NewServer(cfg, charRepo, transfers, classStore, fieldStore, nil)
	ds := &DataStore{
		ClassStore: classStore,
		FieldStore: fieldStore,
		NPCStore:   npcStore,
		MobInfoSHN: mobInfoSHN,
	}
	srv.SetDataStore(ds)

	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start zone server: %v", err)
	}
	defer srv.Stop()

	// Preparar transfer key
	randomID := uint16(44556)
	transfers.CreateZoneTransfer(1, "TestUser", testChar.Name, testChar.ID, randomID, 0, "127.0.0.1")

	// Conectar cliente TCP
	conn, err := net.DialTimeout("tcp", "127.0.0.1:19214", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to connect to zone: %v", err)
	}
	defer conn.Close()

	// Leer XOR key
	xorPkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read XOR key: %v", err)
	}
	xorPos, _ := xorPkt.Reader.ReadInt16()
	clientCrypto, err := crypto.NewNetCrypto(xorPos)
	if err != nil {
		t.Fatalf("failed to init crypto: %v", err)
	}

	// Enviar Login Req
	transKeyPkt := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_REQ)
	transKeyPkt.WriteUint16(randomID)
	transKeyPkt.WriteString(testChar.Name, 16)
	if err := sendClientPacket(conn, transKeyPkt, clientCrypto); err != nil {
		t.Fatalf("failed to send transfer key: %v", err)
	}

	// Consumir chunk inicial (0x1038, 0x1039, etc.)
	var mapLoginAckReceived bool
	for i := 0; i < 20; i++ {
		p, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("failed to read chunk packet: %v", err)
		}
		if p.OpCode == protocol.PROTO_NC_CHAR_MAPLOGIN_ACK {
			mapLoginAckReceived = true
			break
		}
	}
	if !mapLoginAckReceived {
		t.Fatalf("Did not receive PROTO_NC_CHAR_MAPLOGIN_ACK in chunk")
	}

	// Enviar PROTO_NC_MAP_LOGINCOMPLETE_CMD (0x1803)
	readyPkt := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD)
	if err := sendClientPacket(conn, readyPkt, clientCrypto); err != nil {
		t.Fatalf("failed to send ready pkt: %v", err)
	}

	// Leer paquetes tras LoginComplete: debe llegar PROTO_NC_BRIEFINFO_MOB_CMD (0x1C09)
	var mobPkt *protocol.Packet
	for i := 0; i < 10; i++ {
		p, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("failed to read post-login packet: %v", err)
		}
		if p.OpCode == protocol.PROTO_NC_BRIEFINFO_MOB_CMD {
			mobPkt = p
			break
		}
	}

	if mobPkt == nil {
		t.Fatalf("Expected PROTO_NC_BRIEFINFO_MOB_CMD (0x1C09) containing map NPCs")
	}

	mobCount, err := mobPkt.Reader.ReadByte()
	if err != nil || mobCount == 0 {
		t.Fatalf("Expected at least 1 NPC in RouN, got %d (err=%v)", mobCount, err)
	}
	t.Logf("PROTO_NC_BRIEFINFO_MOB_CMD recibido con %d NPCs/Gates en RouN", mobCount)

	// Validar que se enviaron los campos oficiales de 149 bytes de PROTO_NC_BRIEFINFO_REGENMOB_CMD
	// y encontrar un NPC interactivo y un Gate
	var (
		foundSmithHandle uint16
		foundSmithMobID  uint16
		foundGateWhere   string
	)

	rouMap := srv.Handler().mapManager.GetOrCreateMap("RouN")
	for _, n := range rouMap.GetAllNPCs() {
		if n.Name == "RouSmithJames" {
			foundSmithHandle = n.Handle
			foundSmithMobID = n.MobID
		}
		if n.IsGate && n.GateToWhere == "GateRou1" {
			foundGateWhere = n.GateToWhere
		}
	}

	if foundSmithHandle == 0 || foundSmithMobID != 28 {
		t.Fatalf("Expected RouSmithJames with MobID=28, got handle=%d, mobid=%d", foundSmithHandle, foundSmithMobID)
	}
	if foundGateWhere != "GateRou1" {
		t.Fatalf("Expected GateRou1 on RouN")
	}
	t.Logf("NPC RouSmithJames verificado (Handle: %d, MobID: %d), GateRou1 verificado", foundSmithHandle, foundSmithMobID)

	// 1. Probar clic en NPC (PROTO_NC_ACT_NPCCLICK_CMD, 0x200A)
	clickPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCCLICK_CMD)
	clickPkt.WriteUint16(foundSmithHandle)
	if err := sendClientPacket(conn, clickPkt, clientCrypto); err != nil {
		t.Fatalf("failed to send NPC click: %v", err)
	}

	// Debe responder con PROTO_NC_ACT_NPCMENUOPEN_REQ (0x201C) con MobID=28
	var openReqPkt *protocol.Packet
	for i := 0; i < 5; i++ {
		p, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("failed to read NPC click response: %v", err)
		}
		if p.OpCode == protocol.PROTO_NC_ACT_NPCMENUOPEN_REQ {
			openReqPkt = p
			break
		}
	}
	if openReqPkt == nil {
		t.Fatalf("Expected PROTO_NC_ACT_NPCMENUOPEN_REQ (0x201C)")
	}
	receivedMobID, _ := openReqPkt.Reader.ReadUint16()
	if receivedMobID != 28 {
		t.Fatalf("Expected NPCMENUOPEN_REQ with mobid 28, got %d", receivedMobID)
	}
	t.Logf("PROTO_NC_ACT_NPCMENUOPEN_REQ (0x201C) recibido exitosamente con MobID %d", receivedMobID)

	// 2. Responder con PROTO_NC_ACT_NPCMENUOPEN_ACK (0x201D)
	menuAckPkt := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCMENUOPEN_ACK)
	menuAckPkt.WriteByte(0) // OK
	if err := sendClientPacket(conn, menuAckPkt, clientCrypto); err != nil {
		t.Fatalf("failed to send NPCMENUOPEN_ACK: %v", err)
	}

	// 3. Probar Warp mediante Gate (PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ, 0x2070)
	// Posicionar al jugador cerca de GateRou1 (6477, 10579)
	player, ok := srv.Handler().getPlayer(srv.Handler().ActivePlayersSnapshot()[0].Client)
	if !ok || player == nil {
		t.Fatalf("Failed to retrieve ZonePlayer")
	}
	player.Char.XPos = 6477
	player.Char.YPos = 10575

	gateReq := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ)
	if err := sendClientPacket(conn, gateReq, clientCrypto); err != nil {
		t.Fatalf("failed to send use gate request: %v", err)
	}

	// Debe responder con PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK (0x2071)
	var gateAckPkt *protocol.Packet
	for i := 0; i < 10; i++ {
		p, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("failed to read gate ack: %v", err)
		}
		if p.OpCode == protocol.PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK {
			gateAckPkt = p
			break
		}
	}
	if gateAckPkt == nil {
		t.Fatalf("Expected PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK (0x2071)")
	}
	gateErr, _ := gateAckPkt.Reader.ReadUint16()
	if gateErr != 0 {
		t.Fatalf("Expected gate ack error code 0, got %d", gateErr)
	}

	// Validar que el jugador fue teletransportado a RouCos02 (Sand Hill) según LinkTable.json
	time.Sleep(50 * time.Millisecond)
	if player.CurrentMap.MapID != "RouCos02" {
		t.Fatalf("Expected player in map RouCos02, got %s", player.CurrentMap.MapID)
	}
	if player.Char.XPos != 17418 || player.Char.YPos != 13832 {
		t.Fatalf("Expected position (17418, 13832), got (%d, %d)", player.Char.XPos, player.Char.YPos)
	}
	t.Logf("Warp por Gate exitoso a %s (%d, %d)", player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos)

	// Validar que se persistió la ubicación en base de datos
	persistedChar, err := charRepo.GetCharacterByID(testChar.ID)
	if err != nil || persistedChar == nil {
		t.Fatalf("Failed to fetch persisted char %d: %v", testChar.ID, err)
	}
	if persistedChar.TownName != "RouCos02" || persistedChar.XPos != 17418 {
		t.Fatalf("Expected persisted location RouCos02 (17418), got %s (%d)", persistedChar.TownName, persistedChar.XPos)
	}
	t.Logf("Persistencia de ubicación de Warp en base de datos validada")

	// 4. Probar Multi Link Select (PROTO_NC_MAP_MULTY_LINK_SELECT_REQ, 0x181F)
	// Teletransportar de regreso a RouN
	multiPkt := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_REQ)
	multiPkt.WriteString("RouN", 12)
	if err := sendClientPacket(conn, multiPkt, clientCrypto); err != nil {
		t.Fatalf("failed to send multi link select: %v", err)
	}

	var multiAck *protocol.Packet
	for i := 0; i < 15; i++ {
		p, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("failed to read multi link ack: %v", err)
		}
		if p.OpCode == protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_ACK {
			multiAck = p
			break
		}
	}
	if multiAck == nil {
		t.Fatalf("Expected PROTO_NC_MAP_MULTY_LINK_SELECT_ACK (0x1820)")
	}
	multiErr, _ := multiAck.Reader.ReadUint16()
	if multiErr != 0 {
		t.Fatalf("Expected multi link err=0, got %d", multiErr)
	}

	time.Sleep(50 * time.Millisecond)
	if player.CurrentMap.MapID != "RouN" {
		t.Fatalf("Expected player back in RouN, got %s", player.CurrentMap.MapID)
	}
	t.Logf("Multi link warp de regreso a RouN verificado exitosamente")
}
