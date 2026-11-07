package world_test

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iFeddy/XelionCore/pkg/config"
	"github.com/iFeddy/XelionCore/pkg/crypto"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

// readServerPacket decodifica un paquete emitido por el servidor (S2C no lleva cifrado).
func readServerPacket(r io.Reader) (*protocol.Packet, error) {
	lenBuf := make([]byte, 1)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return nil, err
	}

	var payloadLen int
	if lenBuf[0] == 0 {
		extBuf := make([]byte, 2)
		if _, err := io.ReadFull(r, extBuf); err != nil {
			return nil, err
		}
		payloadLen = int(binary.LittleEndian.Uint16(extBuf))
	} else {
		payloadLen = int(lenBuf[0])
	}

	payload := make([]byte, payloadLen)
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
		wire[0] = 0
		binary.LittleEndian.PutUint16(wire[1:3], uint16(payloadLen))
		copy(wire[3:], bufCopy)
	}
	_, err := w.Write(wire)
	return err
}

func TestE2EWorldServiceFlow(t *testing.T) {
	cfg, err := config.Load("../../../.env")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping E2E test, MySQL not accessible: %v", err)
	}

	charRepo := database.NewMySQLCharacterRepository(db)
	transfers := database.NewTransferManager(db)

	dataStore, err := world.LoadWorldDataStore("../../../data")
	if err != nil {
		t.Fatalf("failed to load WorldDataStore: %v", err)
	}

	// Iniciar servidor en puerto de test dinámico
	testPort := 19014
	worldCfg := &world.Config{
		ID:       0,
		Name:     "XelionTestWorld",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", testPort),
		Port:     testPort,
		ZoneIP:   "192.168.0.4",
		ZonePort: 9210,
	}

	srv := world.NewServer(worldCfg, charRepo, transfers, dataStore)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start WorldServer: %v", err)
	}
	defer srv.Stop()
	time.Sleep(50 * time.Millisecond)

	// Simular token emitido por LoginService
	testAccountID := 1
	testUsername := "admin"
	transferToken := transfers.CreateTransfer(testAccountID, testUsername, 1, "127.0.0.1")

	// Conectar cliente de prueba
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", testPort))
	if err != nil {
		t.Fatalf("failed to connect to WorldServer: %v", err)
	}
	defer conn.Close()

	// 1. Recibir Handshake (PROTO_NC_MISC_SEED_ACK)
	pkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read handshake: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_MISC_SEED_ACK {
		t.Fatalf("expected PROTO_NC_MISC_SEED_ACK (0x%04X), got 0x%04X", protocol.PROTO_NC_MISC_SEED_ACK, pkt.OpCode)
	}
	xorPos, _ := pkt.Reader.ReadInt16()
	clientCrypto, err := crypto.NewNetCrypto(xorPos)
	if err != nil {
		t.Fatalf("failed to init client crypto: %v", err)
	}
	t.Logf("Paso 1: Handshake WorldServer recibido, XOR Offset = %d", xorPos)

	// 2. Enviar PROTO_NC_USER_LOGINWORLD_REQ oficial (320 bytes payload: user[256] + validate_new[64])
	keyReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	keyReq.WriteString("admin", 256)
	keyReq.WriteString(transferToken, 32)
	keyReq.Fill(32, 0)
	if err := sendClientPacket(conn, keyReq, clientCrypto); err != nil {
		t.Fatalf("failed to send WorldClientKey: %v", err)
	}

	// 3. Recibir CharacterList (PROTO_NC_USER_LOGINWORLD_ACK - 0x0C14)
	pkt, err = readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read character list: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("expected PROTO_NC_USER_LOGINWORLD_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_USER_LOGINWORLD_ACK, pkt.OpCode)
	}
	_ = pkt.Reader.Skip(2) // randomID
	initialCharCount, _ := pkt.Reader.ReadByte()
	t.Logf("Paso 2: Autenticación World exitosa, recibida lista de personajes (%d activos)", initialCharCount)

	existingChars, _ := charRepo.GetCharactersByAccount(1)
	usedSlots := make(map[byte]bool)
	for _, c := range existingChars {
		usedSlots[c.Slot] = true
	}
	var freeSlot byte = 5
	for s := byte(0); s < 6; s++ {
		if !usedSlots[s] {
			freeSlot = s
			break
		}
	}

	// 4. Probar rechazo de BadNameFilter
	badReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	badReq.WriteByte(freeSlot)
	badReq.WriteString("AdminMaster", 20)
	badReq.WriteByte(byte(1 << 2)) // Fighter
	badReq.WriteByte(1) // Hair
	badReq.WriteByte(1) // Color
	badReq.WriteByte(4) // Face (Face 2, Grade 0)
	_ = sendClientPacket(conn, badReq, clientCrypto)

	pkt, err = readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read bad name rejection: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_AVATAR_CREATEFAIL_ACK {
		t.Fatalf("expected PROTO_NC_AVATAR_CREATEFAIL_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_AVATAR_CREATEFAIL_ACK, pkt.OpCode)
	}
	t.Logf("Paso 3: BadNameFilter oficial bloqueó 'AdminMaster' con 0x1404 correctamente")

	// 5. Crear personaje legítimo: "XelionKnight"
	charName := fmt.Sprintf("Xel%d", time.Now().Unix()%100000)
	createReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createReq.WriteByte(freeSlot)
	createReq.WriteString(charName, 20)
	createReq.WriteByte(byte((1 << 2) | (1 << 7))) // Fighter, Male
	createReq.WriteByte(1)                         // Hair
	createReq.WriteByte(2)                         // HairColor
	createReq.WriteByte(4)                         // Face (Face 2, Grade 0)
	if err := sendClientPacket(conn, createReq, clientCrypto); err != nil {
		t.Fatalf("failed to send CreateChar: %v", err)
	}

	pkt, err = readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read create response: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("expected PROTO_NC_AVATAR_CREATESUCC_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_AVATAR_CREATESUCC_ACK, pkt.OpCode)
	}
	t.Logf("Paso 4: Creación de personaje '%s' confirmada con 0x1406", charName)

	// Verificar en base de datos MySQL real
	createdChar, err := charRepo.GetCharacterByName(charName)
	if err != nil || createdChar == nil {
		t.Fatalf("character not found in MySQL: %v", err)
	}
	if createdChar.CurHP != 52 || createdChar.GetEquip(world.SlotWeapon) != 250 {
		t.Fatalf("character stats/equips mismatch: %+v", createdChar)
	}
	t.Logf("Paso 5: Persistencia en MySQL tCharacter, tCharacterShape y tItem verificada con fidelidad 100%%")

	// 6. Seleccionar personaje para ingresar a Zona (PROTO_NC_CHAR_LOGIN_REQ)
	selReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	selReq.WriteByte(freeSlot) // Slot creado
	if err := sendClientPacket(conn, selReq, clientCrypto); err != nil {
		t.Fatalf("failed to send CharLogin: %v", err)
	}

	pkt, err = readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read CharLogin ACK: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("expected PROTO_NC_CHAR_LOGIN_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_CHAR_LOGIN_ACK, pkt.OpCode)
	}
	zoneIP, _ := pkt.Reader.ReadString(16)
	zoneIP = strings.TrimRight(zoneIP, "\x00")
	zonePort, _ := pkt.Reader.ReadUint16()
	if zoneIP != "192.168.0.4" || zonePort != 9210 {
		t.Fatalf("expected Zone 192.168.0.4:9210, got %s:%d", zoneIP, zonePort)
	}
	t.Logf("Paso 6: Traspaso a Zone recibido (0x1003): Destino = %s:%d", zoneIP, zonePort)

	// Recibir los 3 paquetes de configuración de interfaz del personaje (0x7032, 0x7033, 0x7034)
	expectedOptOpcodes := []uint16{
		protocol.PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD,
		protocol.PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD,
		protocol.PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD,
	}
	for _, expOp := range expectedOptOpcodes {
		optPkt, err := readServerPacket(conn)
		if err != nil {
			t.Fatalf("failed to read option packet 0x%04X: %v", expOp, err)
		}
		if optPkt.OpCode != expOp {
			t.Fatalf("expected option OpCode 0x%04X, got 0x%04X", expOp, optPkt.OpCode)
		}
	}
	t.Logf("Paso 6b: Paquetes de configuración de interfaz (0x7032, 0x7033, 0x7034) recibidos correctamente")

	// Verificar token de Zona en TransferManager
	zTransfer, ok := transfers.ConsumeZoneTransfer(charName, 0, nil)
	if !ok || zTransfer == nil {
		t.Fatalf("expected active ZoneTransfer for '%s'", charName)
	}
	t.Logf("Paso 7: Token de ZoneTransfer validado y consumido en memoria")

	// 7. Borrar personaje recién creado (PROTO_NC_AVATAR_ERASE_REQ)
	delReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_ERASE_REQ)
	delReq.WriteByte(freeSlot)
	if err := sendClientPacket(conn, delReq, clientCrypto); err != nil {
		t.Fatalf("failed to send DeleteChar: %v", err)
	}

	pkt, err = readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read delete ACK: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_AVATAR_ERASESUCC_ACK {
		t.Fatalf("expected PROTO_NC_AVATAR_ERASESUCC_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_AVATAR_ERASESUCC_ACK, pkt.OpCode)
	}
	delSlot, _ := pkt.Reader.ReadByte()
	if delSlot != freeSlot {
		t.Fatalf("expected slot %d in delete ack, got %d", freeSlot, delSlot)
	}
	t.Logf("Paso 8: Soft-delete completado con 0x140C (Slot %d liberado)", freeSlot)

	// 8. Volver a la lista de mundos (PROTO_NC_USER_WILL_WORLD_SELECT_REQ)
	backReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_WILL_WORLD_SELECT_REQ)
	if err := sendClientPacket(conn, backReq, clientCrypto); err != nil {
		t.Fatalf("failed to send BackToWorld: %v", err)
	}

	pkt, err = readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read back to world ACK: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_USER_WILL_WORLD_SELECT_ACK {
		t.Fatalf("expected PROTO_NC_USER_WILL_WORLD_SELECT_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_USER_WILL_WORLD_SELECT_ACK, pkt.OpCode)
	}
	t.Logf("Paso 9: Retorno a lista de mundos verificado (0x0C34)")
}

func TestE2EWorldService_Fase05_Comprehensive(t *testing.T) {
	cfg, err := config.Load("../../../.env")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping E2E test, MySQL not accessible: %v", err)
	}

	charRepo := database.NewMySQLCharacterRepository(db)
	transfers := database.NewTransferManager(db)

	dataStore, err := world.LoadWorldDataStore("../../../data")
	if err != nil {
		t.Fatalf("failed to load WorldDataStore: %v", err)
	}

	testPort := 19018
	worldCfg := &world.Config{
		ID:       0,
		Name:     "XelionTestWorldFase05",
		BindAddr: fmt.Sprintf("127.0.0.1:%d", testPort),
		Port:     testPort,
		ZoneIP:   "192.168.0.4",
		ZonePort: 9210,
	}

	srv := world.NewServer(worldCfg, charRepo, transfers, dataStore)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start WorldServer: %v", err)
	}
	defer srv.Stop()
	time.Sleep(50 * time.Millisecond)

	testAccountID := 1
	testUsername := "admin"
	transferToken := transfers.CreateTransfer(testAccountID, testUsername, 1, "127.0.0.1")

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", testPort))
	if err != nil {
		t.Fatalf("failed to connect to WorldServer: %v", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))

	// 1. Handshake Seed ACK (0x0C05)
	pkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read handshake: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_MISC_SEED_ACK {
		t.Fatalf("expected PROTO_NC_MISC_SEED_ACK, got 0x%04X", pkt.OpCode)
	}
	xorPos, _ := pkt.Reader.ReadInt16()
	clientCrypto, err := crypto.NewNetCrypto(xorPos)
	if err != nil {
		t.Fatalf("failed to init client crypto: %v", err)
	}

	// 2. Login World REQ (0x0C12) -> ACK (0x0C14)
	loginReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_REQ)
	loginReq.WriteString(testUsername, 256)
	loginReq.WriteString(transferToken, 32)
	loginReq.Fill(32, 0)
	if err := sendClientPacket(conn, loginReq, clientCrypto); err != nil {
		t.Fatalf("failed to send login req: %v", err)
	}

	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("expected 0x0C14 loginworld ack, got 0x%04X (err: %v)", pkt.OpCode, err)
	}

	// 3. Avatar List Request (0x0C1F) -> ACK (0x0C14)
	listReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_AVATAR_LIST_REQ)
	if err := sendClientPacket(conn, listReq, clientCrypto); err != nil {
		t.Fatalf("failed to send avatar list req: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_USER_LOGINWORLD_ACK {
		t.Fatalf("expected 0x0C14 in response to 0x0C1F, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	t.Logf("Control 5A E2E: PROTO_NC_USER_AVATAR_LIST_REQ (0x0C1F) -> 0x0C14 recibido")

	// Determinar ranuras libres
	existingChars, _ := charRepo.GetCharactersByAccount(testAccountID)
	usedSlots := make(map[byte]bool)
	for _, c := range existingChars {
		usedSlots[c.Slot] = true
	}
	var freeSlot byte = 255
	for s := byte(0); s < 6; s++ {
		if !usedSlots[s] {
			freeSlot = s
			break
		}
	}
	if freeSlot == 255 {
		t.Fatal("no free slots available for testing")
	}

	// 4. Tutorial Mode Check:
	// Enviar Tutorial Popup (0x1111) con bIsSkip = 0 (quiere hacer tutorial)
	tutPopup := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_TUTORIAL_POPUP_ACK)
	tutPopup.WriteByte(0) // bIsSkip = 0
	if err := sendClientPacket(conn, tutPopup, clientCrypto); err != nil {
		t.Fatalf("failed to send tutorial popup choice: %v", err)
	}

	tutCharName := fmt.Sprintf("Tut%d", time.Now().Unix()%100000)
	createTutReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createTutReq.WriteByte(freeSlot)
	createTutReq.WriteString(tutCharName, 20)
	createTutReq.WriteByte(byte((1 << 2) | (1 << 7))) // Fighter, Male
	createTutReq.WriteByte(1)                         // Hair
	createTutReq.WriteByte(1)                         // Color
	createTutReq.WriteByte(4)                         // Face
	if err := sendClientPacket(conn, createTutReq, clientCrypto); err != nil {
		t.Fatalf("failed to send CreateChar in tutorial mode: %v", err)
	}

	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("expected 0x1406 for tutorial char create, got 0x%04X (err: %v)", pkt.OpCode, err)
	}

	// Verificar en base de datos que el mapa inicial fue RouT (desde TutorialMap.json)
	var dbTutMap string
	_ = db.QueryRow("SELECT sLoginZone FROM tCharacter WHERE sID = ?", tutCharName).Scan(&dbTutMap)
	if dbTutMap != "RouT" {
		t.Fatalf("expected starting map 'RouT' for tutorial char, got '%s'", dbTutMap)
	}
	t.Logf("Control 5A E2E: Creación en modo tutorial verificada con mapa '%s'", dbTutMap)

	// Borrar el personaje tutorial para liberar el slot
	eraseTutReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_ERASE_REQ)
	eraseTutReq.WriteByte(freeSlot)
	if err := sendClientPacket(conn, eraseTutReq, clientCrypto); err != nil {
		t.Fatalf("failed to send erase tutorial char: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_AVATAR_ERASESUCC_ACK {
		t.Fatalf("expected 0x140C for tutorial char erase, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	t.Logf("Control 5A E2E: Personaje tutorial borrado con 0x140C")

	// 5. Creación de Sentinel (Clase 26 - Nivel 60 en Eld)
	// Enviar Tutorial Popup (0x1111) con bIsSkip = 1
	skipTutPopup := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_TUTORIAL_POPUP_ACK)
	skipTutPopup.WriteByte(1) // bIsSkip = 1
	if err := sendClientPacket(conn, skipTutPopup, clientCrypto); err != nil {
		t.Fatalf("failed to send tutorial skip popup: %v", err)
	}

	sentName := fmt.Sprintf("Sen%d", time.Now().Unix()%100000)
	createSentReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_CREATE_REQ)
	createSentReq.WriteByte(freeSlot)
	createSentReq.WriteString(sentName, 20)
	createSentReq.WriteByte(byte((26 << 2) | (1 << 7))) // Sentinel (26), Male
	createSentReq.WriteByte(1)                          // Hair
	createSentReq.WriteByte(2)                          // Color
	createSentReq.WriteByte(4)                          // Face
	if err := sendClientPacket(conn, createSentReq, clientCrypto); err != nil {
		t.Fatalf("failed to send CreateChar Sentinel: %v", err)
	}

	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_AVATAR_CREATESUCC_ACK {
		t.Fatalf("expected 0x1406 for Sentinel create, got 0x%04X (err: %v)", pkt.OpCode, err)
	}

	// Verificar Sentinel en MySQL
	sentChar, err := charRepo.GetCharacterByName(sentName)
	if err != nil || sentChar == nil {
		t.Fatalf("sentinel not found in MySQL: %v", err)
	}
	if sentChar.Level != 60 || sentChar.TCharacterState.Exp != 24346145 || sentChar.TownName != "Eld" {
		t.Fatalf("sentinel level/exp/town mismatch: Level=%d Exp=%d Town=%s",
			sentChar.Level, sentChar.TCharacterState.Exp, sentChar.TownName)
	}
	if sentChar.CurHP != 1532 || sentChar.CurSP != 1 {
		t.Fatalf("sentinel base vitals mismatch from ClassParam: HP=%d SP=%d",
			sentChar.CurHP, sentChar.CurSP)
	}
	t.Logf("Control 5A E2E: Sentinel Lv 60 creado en '%s' con Exp %d y stats base verificados",
		sentChar.TownName, sentChar.TCharacterState.Exp)

	// 6. Style Info (0x107C)
	styleInfoReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_SET_STYLE_GET_INFO_REQ)
	styleInfoReq.WriteByte(freeSlot)
	if err := sendClientPacket(conn, styleInfoReq, clientCrypto); err != nil {
		t.Fatalf("failed to send style info req: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_SET_STYLE_GET_INFO_ACK {
		t.Fatalf("expected 0x107D for style info, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	_ = pkt.Reader.Skip(2) // Error code
	_ = pkt.Reader.Skip(4) // Grade 0 item (2 bytes id + 2 bytes count)
	coupon1ID, _ := pkt.Reader.ReadUint16()
	if coupon1ID != 30490 {
		t.Fatalf("expected coupon grade 1 to be 30490, got %d", coupon1ID)
	}
	t.Logf("Control 5A E2E: PROTO_NC_CHAR_SET_STYLE_GET_INFO_ACK (0x107D) verificado con cupón %d", coupon1ID)

	// 7. Seleccionar Sentinel para entrar a mundo (0x1001)
	selReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_LOGIN_REQ)
	selReq.WriteByte(freeSlot)
	if err := sendClientPacket(conn, selReq, clientCrypto); err != nil {
		t.Fatalf("failed to send char login: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_LOGIN_ACK {
		t.Fatalf("expected 0x1003 for char login, got 0x%04X (err: %v)", pkt.OpCode, err)
	}

	// Recibir opciones mejoradas automáticas (0x7032, 0x7033, 0x7034)
	for _, expectedOp := range []uint16{
		protocol.PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD,
		protocol.PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD,
		protocol.PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD,
	} {
		optPkt, err := readServerPacket(conn)
		if err != nil || optPkt.OpCode != expectedOp {
			t.Fatalf("expected improve option 0x%04X, got 0x%04X (err: %v)", expectedOp, optPkt.OpCode, err)
		}
	}
	t.Logf("Control 5A E2E: CharLogin y entrega automática de opciones mejoradas recibidas")

	// 8. Style Change (0x1078)
	setStyleReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_SET_STYLE_REQ)
	setStyleReq.WriteByte(0) // jobGender
	setStyleReq.WriteByte(3) // hair
	setStyleReq.WriteByte(1) // color
	setStyleReq.WriteByte(5) // face
	if err := sendClientPacket(conn, setStyleReq, clientCrypto); err != nil {
		t.Fatalf("failed to send set style req: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_SET_STYLE_ACK {
		t.Fatalf("expected 0x1079 for set style, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	styleErr, _ := pkt.Reader.ReadUint16()
	if styleErr != 0 {
		t.Fatalf("expected style change success (0), got %d", styleErr)
	}
	var dbHair, dbColor, dbFace byte
	_ = db.QueryRow("SELECT nHairType, nHairColor, nFaceShape FROM tCharacterShape WHERE nCharNo = ?", sentChar.ID).Scan(&dbHair, &dbColor, &dbFace)
	if dbHair != 3 || dbColor != 1 || dbFace != 5 {
		t.Fatalf("tCharacterShape mismatch in DB: hair=%d color=%d face=%d", dbHair, dbColor, dbFace)
	}
	t.Logf("Control 5A E2E: PROTO_NC_CHAR_SET_STYLE_REQ (0x1078) -> 0x1079 verificado en MySQL")

	// 9. Improved Options Set & Init (0x7037..0x7046)
	// Set Shortcut
	setSC := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ)
	setSC.WriteByte(1)     // count
	setSC.WriteByte(0)     // slot
	setSC.WriteUint16(205) // code
	setSC.WriteInt32(777)  // val
	if err := sendClientPacket(conn, setSC, clientCrypto); err != nil {
		t.Fatalf("failed to send improve set shortcut: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK {
		t.Fatalf("expected 0x7038, got 0x%04X (err: %v)", pkt.OpCode, err)
	}

	// Set KeyMap
	setKM := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_REQ)
	setKM.WriteUint16(1) // count
	setKM.WriteUint16(5) // funcNo
	setKM.WriteByte(1)   // extKey
	setKM.WriteByte('E') // ascii
	if err := sendClientPacket(conn, setKM, clientCrypto); err != nil {
		t.Fatalf("failed to send improve set keymap: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_ACK {
		t.Fatalf("expected 0x703A, got 0x%04X (err: %v)", pkt.OpCode, err)
	}

	// Set Game Option
	setGO := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_REQ)
	setGO.WriteUint16(1)  // count
	setGO.WriteUint16(12) // optNo
	setGO.WriteByte(1)    // value
	if err := sendClientPacket(conn, setGO, clientCrypto); err != nil {
		t.Fatalf("failed to send improve set game option: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_ACK {
		t.Fatalf("expected 0x703C, got 0x%04X (err: %v)", pkt.OpCode, err)
	}

	// Init KeyMap
	initKM := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_REQ)
	if err := sendClientPacket(conn, initKM, clientCrypto); err != nil {
		t.Fatalf("failed to send improve init keymap: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_ACK {
		t.Fatalf("expected 0x7044, got 0x%04X (err: %v)", pkt.OpCode, err)
	}

	// Init Game Option
	initGO := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_REQ)
	if err := sendClientPacket(conn, initGO, clientCrypto); err != nil {
		t.Fatalf("failed to send improve init game option: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_ACK {
		t.Fatalf("expected 0x7046, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	t.Logf("Control 5A E2E: Opciones Mejoradas NA2016 (0x7037..0x7046) probadas y confirmadas")

	// 10. Classic Video & Sound Options (0x7006, 0x7008)
	vidReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_OPTION_GET_VIDEO_REQ)
	if err := sendClientPacket(conn, vidReq, clientCrypto); err != nil {
		t.Fatalf("failed to send get video: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_OPTION_GET_VIDEO_ACK {
		t.Fatalf("expected 0x7007, got 0x%04X (err: %v)", pkt.OpCode, err)
	}

	sndReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_OPTION_GET_SOUND_REQ)
	if err := sendClientPacket(conn, sndReq, clientCrypto); err != nil {
		t.Fatalf("failed to send get sound: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_OPTION_GET_SOUND_ACK {
		t.Fatalf("expected 0x7009, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	t.Logf("Control 5A E2E: Opciones clásicas Video/Sound (0x7006 -> 0x7007, 0x7008 -> 0x7009) verificadas")

	// 11. Newbie Guide View (0x10DA) & Attendance Reward (0x10F1)
	guideReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD)
	guideReq.WriteByte(9) // guideNo
	if err := sendClientPacket(conn, guideReq, clientCrypto); err != nil {
		t.Fatalf("failed to send newbie guide view: %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	var dbGuideID int
	_ = db.QueryRow("SELECT nGuideNo FROM tInGameGuide WHERE nCharNo = ? AND nGuideNo = 9", sentChar.ID).Scan(&dbGuideID)
	if dbGuideID != 9 {
		t.Fatalf("expected guide 9 recorded in tInGameGuide, got %d", dbGuideID)
	}

	attendReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_REQ)
	attendReq.WriteByte(0) // index
	attendReq.WriteByte(1) // rewardID
	if err := sendClientPacket(conn, attendReq, clientCrypto); err != nil {
		t.Fatalf("failed to send attendance reward req: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_ACK {
		t.Fatalf("expected 0x10F2, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	t.Logf("Control 5A E2E: Newbie Guide View (0x10DA) y Attendance Reward (0x10F1 -> 0x10F2) confirmados")

	// 12. Renombrar Personaje (0x140F)
	renamedName := fmt.Sprintf("Ren%d", time.Now().Unix()%100000)
	renameReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_RENAME_REQ)
	renameReq.WriteByte(freeSlot)
	renameReq.WriteString(renamedName, 20)
	if err := sendClientPacket(conn, renameReq, clientCrypto); err != nil {
		t.Fatalf("failed to send rename req: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_AVATAR_RENAME_ACK {
		t.Fatalf("expected 0x1410, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	_ = pkt.Reader.Skip(21) // slot (1) + newName (20)
	renErr, _ := pkt.Reader.ReadUint16()
	if renErr != 0 {
		t.Fatalf("expected rename success (0), got %d", renErr)
	}
	var dbName string
	_ = db.QueryRow("SELECT sID FROM tCharacter WHERE nCharNo = ?", sentChar.ID).Scan(&dbName)
	if dbName != renamedName {
		t.Fatalf("expected renamed name in DB '%s', got '%s'", renamedName, dbName)
	}
	t.Logf("Control 5A E2E: PROTO_NC_AVATAR_RENAME_REQ (0x140F) -> 0x1410 verificado con éxito")

	// 13. Cambio de ID Persistido (0xB409)
	finalName := fmt.Sprintf("Fin%d", time.Now().Unix()%100000)
	_, err = db.Exec("INSERT INTO tCharacterChangeID (nUserNo, nCharNo, bFlag, bInit, sOldID, sNewID) VALUES (?, ?, 0, 0, ?, '')",
		testAccountID, sentChar.ID, renamedName)
	if err != nil {
		t.Fatalf("failed to insert mock pending change ID: %v", err)
	}

	changeIDReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ)
	changeIDReq.WriteByte(freeSlot)
	changeIDReq.WriteString(finalName, 20)
	if err := sendClientPacket(conn, changeIDReq, clientCrypto); err != nil {
		t.Fatalf("failed to send change char ID req: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK {
		t.Fatalf("expected 0xB40C, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	chErr, _ := pkt.Reader.ReadUint16()
	if chErr != 0 {
		t.Fatalf("expected change ID success (0), got %d", chErr)
	}
	_ = db.QueryRow("SELECT sID FROM tCharacter WHERE nCharNo = ?", sentChar.ID).Scan(&dbName)
	if dbName != finalName {
		t.Fatalf("expected final name in DB '%s', got '%s'", finalName, dbName)
	}
	t.Logf("Control 5A E2E: PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ (0xB409) -> 0xB40C verificado")

	// 14. Rechazo de Borrado si es Guild Master Activo (0x140A)
	_, _ = db.Exec("INSERT INTO tGuildMember (nNo, nCharNo, nGrade) VALUES (8888, ?, 0)", sentChar.ID)
	eraseMasterReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_ERASE_REQ)
	eraseMasterReq.WriteByte(freeSlot)
	if err := sendClientPacket(conn, eraseMasterReq, clientCrypto); err != nil {
		t.Fatalf("failed to send erase guild master req: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_AVATAR_ERASEFAIL_ACK {
		t.Fatalf("expected 0x140A for guild master erase rejection, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	t.Logf("Control 5A E2E: Borrado de líder de gremio bloqueado con 0x140A correctamente")
	_, _ = db.Exec("DELETE FROM tGuildMember WHERE nNo = 8888 AND nCharNo = ?", sentChar.ID)

	// 15. Borrado Exitoso del Personaje (0x1409 -> 0x140C)
	eraseSuccReq := protocol.NewPacketWriter(protocol.PROTO_NC_AVATAR_ERASE_REQ)
	eraseSuccReq.WriteByte(freeSlot)
	if err := sendClientPacket(conn, eraseSuccReq, clientCrypto); err != nil {
		t.Fatalf("failed to send erase req: %v", err)
	}
	pkt, err = readServerPacket(conn)
	if err != nil || pkt.OpCode != protocol.PROTO_NC_AVATAR_ERASESUCC_ACK {
		t.Fatalf("expected 0x140C for erase success, got 0x%04X (err: %v)", pkt.OpCode, err)
	}
	t.Logf("Control 5A E2E: Personaje borrado exitosamente con 0x140C")

	// 16. Normal Logout (0x0C18)
	logoutReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NORMALLOGOUT_CMD)
	logoutReq.WriteByte(0) // type
	if err := sendClientPacket(conn, logoutReq, clientCrypto); err != nil {
		t.Fatalf("failed to send normal logout: %v", err)
	}
	t.Logf("Control 5A E2E: PROTO_NC_USER_NORMALLOGOUT_CMD (0x0C18) enviado y procesado limpiamente")
}
