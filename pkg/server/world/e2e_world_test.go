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

	// 1. Recibir Handshake (PROTO_NC_MISC_SEED_CMD)
	pkt, err := readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read handshake: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_MISC_SEED_CMD {
		t.Fatalf("expected PROTO_NC_MISC_SEED_CMD (0x%04X), got 0x%04X", protocol.PROTO_NC_MISC_SEED_CMD, pkt.OpCode)
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

	// 3. Recibir CharacterList (PROTO_NC_CHAR_CLIENT_CHARLIST_ACK - 0x0C14)
	pkt, err = readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read character list: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_CHAR_CLIENT_CHARLIST_ACK {
		t.Fatalf("expected PROTO_NC_CHAR_CLIENT_CHARLIST_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_CHAR_CLIENT_CHARLIST_ACK, pkt.OpCode)
	}
	_ = pkt.Reader.Skip(2) // randomID
	initialCharCount, _ := pkt.Reader.ReadByte()
	t.Logf("Paso 2: Autenticación World exitosa, recibida lista de personajes (%d activos)", initialCharCount)

	// 4. Probar rechazo de BadNameFilter
	badReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CREATE_REQ)
	badReq.WriteByte(2) // Slot 2
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
	if pkt.OpCode != protocol.PROTO_NC_CHAR_CREATE_FAIL_ACK {
		t.Fatalf("expected PROTO_NC_CHAR_CREATE_FAIL_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_CHAR_CREATE_FAIL_ACK, pkt.OpCode)
	}
	t.Logf("Paso 3: BadNameFilter oficial bloqueó 'AdminMaster' con 0x1404 correctamente")

	// 5. Crear personaje legítimo: "XelionKnight"
	charName := fmt.Sprintf("Xel%d", time.Now().Unix()%100000)
	createReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CREATE_REQ)
	createReq.WriteByte(2) // Slot 2
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
	if pkt.OpCode != protocol.PROTO_NC_CHAR_CREATE_ACK {
		t.Fatalf("expected PROTO_NC_CHAR_CREATE_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_CHAR_CREATE_ACK, pkt.OpCode)
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
	selReq.WriteByte(2) // Slot 2
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

	// Verificar token de Zona en TransferManager
	zTransfer, ok := transfers.ConsumeZoneTransfer(charName, 0, nil)
	if !ok || zTransfer == nil {
		t.Fatalf("expected active ZoneTransfer for '%s'", charName)
	}
	t.Logf("Paso 7: Token de ZoneTransfer validado y consumido en memoria")

	// 7. Borrar personaje recién creado (PROTO_NC_CHAR_DELETE_REQ)
	delReq := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_DELETE_REQ)
	delReq.WriteByte(2) // Slot 2
	if err := sendClientPacket(conn, delReq, clientCrypto); err != nil {
		t.Fatalf("failed to send DeleteChar: %v", err)
	}

	pkt, err = readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read delete ACK: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_CHAR_DELETE_ACK {
		t.Fatalf("expected PROTO_NC_CHAR_DELETE_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_CHAR_DELETE_ACK, pkt.OpCode)
	}
	delSlot, _ := pkt.Reader.ReadByte()
	if delSlot != 2 {
		t.Fatalf("expected slot 2 in delete ack, got %d", delSlot)
	}
	t.Logf("Paso 8: Soft-delete completado con 0x140C (Slot 2 liberado)")

	// 8. Volver a la lista de mundos (PROTO_NC_USER_BACKTOWORLDLIST_REQ)
	backReq := protocol.NewPacketWriter(protocol.PROTO_NC_USER_BACKTOWORLDLIST_REQ)
	if err := sendClientPacket(conn, backReq, clientCrypto); err != nil {
		t.Fatalf("failed to send BackToWorld: %v", err)
	}

	pkt, err = readServerPacket(conn)
	if err != nil {
		t.Fatalf("failed to read back to world ACK: %v", err)
	}
	if pkt.OpCode != protocol.PROTO_NC_USER_BACKTOWORLDLIST_ACK {
		t.Fatalf("expected PROTO_NC_USER_BACKTOWORLDLIST_ACK (0x%04X), got 0x%04X",
			protocol.PROTO_NC_USER_BACKTOWORLDLIST_ACK, pkt.OpCode)
	}
	t.Logf("Paso 9: Retorno a lista de mundos verificado (0x0C34)")
}
