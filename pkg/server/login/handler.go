package login

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler procesa el 100% de los paquetes entrantes para LoginService respetando estrictamente
// la arquitectura unificada y las estructuras oficiales de red (docs/Servicios/Login.md, Account.md, AccountLog.md).
type Handler struct {
	repo         database.AccountRepository
	logRepo      database.AccountLogRepository
	sessionRepo  database.LoginSessionRepository
	worldManager *WorldManager
	transfers    *database.TransferManager
}

// NewHandler inicializa el despachador de protocolos oficiales de LoginService.
func NewHandler(repo database.AccountRepository, wm *WorldManager, tm *database.TransferManager, sessionRepo ...database.LoginSessionRepository) *Handler {
	if tm == nil {
		tm = database.NewTransferManager()
	}
	var sRepo database.LoginSessionRepository
	if len(sessionRepo) > 0 && sessionRepo[0] != nil {
		sRepo = sessionRepo[0]
	} else if r, ok := repo.(database.LoginSessionRepository); ok {
		sRepo = r
	}
	var lRepo database.AccountLogRepository
	if r, ok := repo.(database.AccountLogRepository); ok {
		lRepo = r
	}
	return &Handler{
		repo:         repo,
		logRepo:      lRepo,
		sessionRepo:  sRepo,
		worldManager: wm,
		transfers:    tm,
	}
}

// HandlePacket enruta cada paquete entrante según su OpCode oficial PROTO_NC_*.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	// 1. Handshake y versión C2S
	case protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ: // 0x0C65
		h.handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_XTRAP_REQ: // 0x0C04
		h.handle_PROTO_NC_USER_XTRAP_REQ(client, pkt)

	// 2. Variantes de autenticación C2S
	case protocol.PROTO_NC_USER_US_LOGIN_REQ: // 0x0C5A (US)
		h.handle_PROTO_NC_USER_US_LOGIN_REQ(client, pkt)
	case protocol.PROTO_NC_USER_LOGIN_REQ: // 0x0C06 (Estándar internacional)
		h.handle_PROTO_NC_USER_LOGIN_REQ(client, pkt)
	case protocol.PROTO_NC_USER_LOGIN_WITH_OTP_REQ: // 0x0C37 (OTP 32B)
		h.handle_PROTO_NC_USER_LOGIN_WITH_OTP_REQ(client, pkt)
	case protocol.PROTO_NC_USER_LOGIN_NETMARBLE_REQ: // 0x0C1D (Netmarble)
		h.handle_PROTO_NC_USER_LOGIN_NETMARBLE_REQ(client, pkt)
	case protocol.PROTO_NC_USER_LOGIN_OUTSPARK_REQ: // 0x0C20 (Outspark 85B)
		h.handle_PROTO_NC_USER_LOGIN_OUTSPARK_REQ(client, pkt)

	// 3. Handlers regionales de contraseña y autenticación (GER, TW, JP, CH, US)
	case protocol.PROTO_NC_USER_GER_PASSWORD_CHECK_REQ: // 0x0C39
		h.handle_PROTO_NC_USER_GER_PASSWORD_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_TW_PASSWORD_CHECK_REQ: // 0x0C3E
		h.handle_PROTO_NC_USER_TW_PASSWORD_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_JP_PASSWORD_CHECK_REQ: // 0x0C4B
		h.handle_PROTO_NC_USER_JP_PASSWORD_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_CH_PASSWORD_CHECK_REQ: // 0x0C53
		h.handle_PROTO_NC_USER_CH_PASSWORD_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_US_PASSWORD_CHECK_REQ: // 0x0C5B
		h.handle_PROTO_NC_USER_US_PASSWORD_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_PASSWORD_CHECK_REQ: // 0x0C07
		h.handle_PROTO_NC_USER_PASSWORD_CHECK_REQ(client, pkt)

	// 4. Chequeos de IP Block (C2S/S2S)
	case protocol.PROTO_NC_USER_IS_IP_BLOCK_REQ: // 0x0C2A
		h.handle_PROTO_NC_USER_IS_IP_BLOCK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_GER_IS_IP_BLOCK_REQ: // 0x0C3B
		h.handle_PROTO_NC_USER_GER_IS_IP_BLOCK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_TW_IS_IP_BLOCK_REQ: // 0x0C40
		h.handle_PROTO_NC_USER_TW_IS_IP_BLOCK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_JP_IS_IP_BLOCK_REQ: // 0x0C4D
		h.handle_PROTO_NC_USER_JP_IS_IP_BLOCK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_US_IS_IP_BLOCK_REQ: // 0x0C5D
		h.handle_PROTO_NC_USER_US_IS_IP_BLOCK_REQ(client, pkt)

	// 5. Restricción Teenager (Control parental de minutos jugados)
	case protocol.PROTO_NC_USER_TEENAGER_REQ: // 0x0C26
		h.handle_PROTO_NC_USER_TEENAGER_REQ(client, pkt)
	case protocol.PROTO_NC_USER_TEENAGER_SET_CMD: // 0x0C28
		h.handle_PROTO_NC_USER_TEENAGER_SET_CMD(client, pkt)

	// 6. Selección y estado de mundos
	case protocol.PROTO_NC_USER_WORLD_STATUS_REQ: // 0x0C1B
		h.handle_PROTO_NC_USER_WORLD_STATUS_REQ(client, pkt)
	case protocol.PROTO_NC_USER_WORLDSELECT_REQ: // 0x0C0B
		h.handle_PROTO_NC_USER_WORLDSELECT_REQ(client, pkt)

	// 7. Desconexión y Heartbeat
	case protocol.PROTO_NC_MISC_HEARTBEAT_REQ: // 0x0804
		h.handle_PROTO_NC_MISC_HEARTBEAT_REQ(client, pkt)
	case protocol.PROTO_NC_MISC_HEARTBEAT_ACK: // 0x0805
		h.handle_PROTO_NC_MISC_HEARTBEAT_ACK(client, pkt)
	case protocol.PROTO_NC_USER_NORMALLOGOUT_CMD: // 0x0C18
		h.handle_PROTO_NC_USER_NORMALLOGOUT_CMD(client, pkt)

	// 8. Handlers S2S de topología y control
	case protocol.PROTO_NC_MISC_S2SCONNECTION_RDY: // 0x0801
		h.handle_PROTO_NC_MISC_S2SCONNECTION_RDY(client, pkt)
	case protocol.PROTO_NC_MISC_S2SCONNECTION_REQ: // 0x0802
		h.handle_PROTO_NC_MISC_S2SCONNECTION_REQ(client, pkt)
	case protocol.PROTO_NC_MISC_USER_COUNT_CMD: // 0x0808
		h.handle_PROTO_NC_MISC_USER_COUNT_CMD(client, pkt)
	case protocol.PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD: // 0x080B
		h.handle_PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD(client, pkt)
	case protocol.PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD: // 0x080C
		h.handle_PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD(client, pkt)
	case protocol.PROTO_NC_MISC_START_THE_BOOM_CMD: // 0x0823
		h.handle_PROTO_NC_MISC_START_THE_BOOM_CMD(client, pkt)

	// 9. OTP generation (S2S desde WorldManager)
	case protocol.PROTO_NC_USER_CREATE_OTP_REQ: // 0x0C35
		h.handle_PROTO_NC_USER_CREATE_OTP_REQ(client, pkt)

	// 10. Retorno y último logout
	case protocol.PROTO_NC_USER_RETURN_CHECK_REQ: // 0x0C60
		h.handle_PROTO_NC_USER_RETURN_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_SET_RETURN_CMD: // 0x0C64
		h.handle_PROTO_NC_USER_SET_RETURN_CMD(client, pkt)
	case protocol.PROTO_NC_USER_LOGOUT_LAST_DAY_REQ: // 0x0C62
		h.handle_PROTO_NC_USER_LOGOUT_LAST_DAY_REQ(client, pkt)

	// 11. Respuestas S2S que pueden ser recibidas internamente
	case protocol.PROTO_NC_MISC_S2SCONNECTION_ACK: // 0x0803
	case protocol.PROTO_NC_USER_PASSWORD_CHECK_ACK: // 0x0C08
	case protocol.PROTO_NC_USER_WILLLOGIN_ACK: // 0x0C0E
	case protocol.PROTO_NC_USER_TEENAGER_ACK: // 0x0C27
	case protocol.PROTO_NC_USER_IS_IP_BLOCK_ACK: // 0x0C2B
	case protocol.PROTO_NC_USER_GER_PASSWORD_CHECK_ACK: // 0x0C3A
	case protocol.PROTO_NC_USER_GER_IS_IP_BLOCK_ACK: // 0x0C3C
	case protocol.PROTO_NC_USER_TW_PASSWORD_CHECK_ACK: // 0x0C3F
	case protocol.PROTO_NC_USER_TW_IS_IP_BLOCK_ACK: // 0x0C41
	case protocol.PROTO_NC_USER_JP_PASSWORD_CHECK_ACK: // 0x0C4C
	case protocol.PROTO_NC_USER_JP_IS_IP_BLOCK_ACK: // 0x0C4E
	case protocol.PROTO_NC_USER_CH_PASSWORD_CHECK_ACK: // 0x0C54
	case protocol.PROTO_NC_USER_US_PASSWORD_CHECK_ACK: // 0x0C5C
	case protocol.PROTO_NC_USER_US_IS_IP_BLOCK_ACK: // 0x0C5E
		log.Printf("[%s] Paquete S2S recibido: 0x%04X", client.RemoteAddr(), pkt.OpCode)

	default:
		log.Printf("[%s] OpCode no manejado en LoginService: 0x%04X (Dept=%d, Cmd=%d)",
			client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type)
	}
}

// ============================================================================
// 1. Handshake y Versión
// ============================================================================

// handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ procesa el handshake de versión del cliente (0x0C65).
func (h *Handler) handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var versionKey string

	if rem >= 64 {
		s, err := pkt.Reader.ReadString(64)
		if err != nil {
			log.Printf("[%s] Error leyendo sVersionKey: %v", client.RemoteAddr(), err)
			_ = client.SendPacket(protocol.New_PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK())
			client.Close()
			return
		}
		versionKey = strings.TrimSpace(s)
	} else if rem == 4 {
		year, _ := pkt.Reader.ReadUint16()
		ver, _ := pkt.Reader.ReadUint16()
		client.ClientYear = year
		client.ClientVersion = ver
	} else if rem > 0 {
		versionKey, _ = pkt.Reader.ReadString(rem)
		versionKey = strings.TrimSpace(versionKey)
	}

	log.Printf("[%s] PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ recibido (Key='%s')", client.RemoteAddr(), versionKey)

	// Resolver versión esperada obligatoriamente desde MySQL (FOUS.GameVersionKey)
	if h.sessionRepo == nil {
		log.Printf("[%s] Repositorio de sesión no disponible para leer FOUS.GameVersionKey. Enviando 0x0C02 y rechazando.",
			client.RemoteAddr())
		_ = client.SendPacket(protocol.New_PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK())
		client.Close()
		return
	}

	dbKey, err := h.sessionRepo.GetClientVersionKey()
	expectedVer := strings.TrimSpace(dbKey)
	if err != nil || expectedVer == "" {
		log.Printf("[%s] Error o clave vacía en FOUS.GameVersionKey (err=%v). Enviando 0x0C02 y rechazando.",
			client.RemoteAddr(), err)
		_ = client.SendPacket(protocol.New_PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK())
		client.Close()
		return
	}

	if versionKey == "" || versionKey != expectedVer {
		log.Printf("[%s] Versión de cliente no coincide: '%s' (Esperada en FOUS: '%s'). Enviando 0x0C02 y rechazando.",
			client.RemoteAddr(), versionKey, expectedVer)
		_ = client.SendPacket(protocol.New_PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK())
		client.Close()
		return
	}

	// 1. Responder con éxito de versión: PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK (0x0C67, 1B)
	_ = client.SendPacket(protocol.New_PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK())
}

// handle_PROTO_NC_USER_XTRAP_REQ procesa la verificación de archivos/XTrap (0x0C04).
func (h *Handler) handle_PROTO_NC_USER_XTRAP_REQ(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.New_PROTO_NC_USER_XTRAP_ACK(true)
	_ = client.SendPacket(resp)
}

// ============================================================================
// 2. Variantes de Autenticación C2S
// ============================================================================

// handle_PROTO_NC_USER_US_LOGIN_REQ procesa el login oficial para la versión US (0x0C5A, 316 bytes).
func (h *Handler) handle_PROTO_NC_USER_US_LOGIN_REQ(client *network.Client, pkt *protocol.Packet) {
	var username, password string
	rem := pkt.Reader.Remaining()
	if rem >= 316 {
		username, _ = pkt.Reader.ReadString(260)
		password, _ = pkt.Reader.ReadString(36)
		_, _ = pkt.Reader.ReadString(20) // spawnapps
	} else if rem >= 292 {
		username, _ = pkt.Reader.ReadString(260)
		password, _ = pkt.Reader.ReadString(32)
	} else if rem >= 32 {
		username, _ = pkt.Reader.ReadString(16)
		password, _ = pkt.Reader.ReadString(rem - 16)
	} else {
		username, _ = pkt.Reader.ReadString(rem)
	}

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	clientIP := client.RemoteIP()

	log.Printf("[%s] PROTO_NC_USER_US_LOGIN_REQ para usuario '%s'", client.RemoteAddr(), username)

	var acc *database.Account
	var err error
	if h.repo != nil {
		acc, err = h.repo.LoginGameUS(username, password, clientIP)
	} else {
		err = database.ErrAccountNotFound
	}

	h.executeLoginPipeline(client, acc, username, password, clientIP, err)
}

// handle_PROTO_NC_USER_LOGIN_REQ procesa el login internacional oficial (0x0C06, 272 bytes: User[256], Password[16]).
func (h *Handler) handle_PROTO_NC_USER_LOGIN_REQ(client *network.Client, pkt *protocol.Packet) {
	var username, password string
	rem := pkt.Reader.Remaining()
	if rem >= 272 {
		username, _ = pkt.Reader.ReadString(256)
		password, _ = pkt.Reader.ReadString(16)
	} else if rem >= 32 {
		username, _ = pkt.Reader.ReadString(16)
		password, _ = pkt.Reader.ReadString(rem - 16)
	} else {
		username, _ = pkt.Reader.ReadString(rem)
	}

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	clientIP := client.RemoteIP()

	log.Printf("[%s] PROTO_NC_USER_LOGIN_REQ para usuario '%s'", client.RemoteAddr(), username)

	var acc *database.Account
	var err error
	if h.repo != nil {
		acc, err = h.repo.LoginGame(username, password, clientIP)
	} else {
		err = database.ErrAccountNotFound
	}

	h.executeLoginPipeline(client, acc, username, password, clientIP, err)
}

// handle_PROTO_NC_USER_LOGIN_WITH_OTP_REQ procesa la reconexión autenticada vía OTP (0x0C37, 32 bytes).
func (h *Handler) handle_PROTO_NC_USER_LOGIN_WITH_OTP_REQ(client *network.Client, pkt *protocol.Packet) {
	var otp [32]byte
	rem := pkt.Reader.Remaining()
	if rem >= 32 {
		b, err := pkt.Reader.ReadBytes(32)
		if err == nil {
			copy(otp[:], b)
		}
	} else {
		b, _ := pkt.Reader.ReadBytes(rem)
		copy(otp[:], b)
	}

	clientIP := client.RemoteIP()
	log.Printf("[%s] PROTO_NC_USER_LOGIN_WITH_OTP_REQ recibido (OTP: %s)", client.RemoteAddr(), string(otp[:8]))

	if h.sessionRepo == nil || h.repo == nil {
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrInvalidCredentials)
		return
	}

	userNo, userID, err := h.sessionRepo.ConsumeOTP(otp, clientIP)
	if err != nil || userNo <= 0 {
		log.Printf("[%s] OTP inválido o expirado: %v", client.RemoteAddr(), err)
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrInvalidCredentials)
		return
	}

	acc, err := h.repo.GetAccountByUsername(userID)
	if err != nil {
		log.Printf("[%s] Error cargando cuenta por OTP (UserNo: %d, ID: '%s'): %v", client.RemoteAddr(), userNo, userID, err)
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrDatabase)
		return
	}

	h.executeLoginPipeline(client, acc, userID, "", clientIP, nil)
}

// handle_PROTO_NC_USER_LOGIN_NETMARBLE_REQ procesa el login de plataforma Netmarble (0x0C1D).
func (h *Handler) handle_PROTO_NC_USER_LOGIN_NETMARBLE_REQ(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var uniID, userID string
	if rem >= 278 {
		_, _ = pkt.Reader.ReadUint16() // Netpacketheader
		uniID, _ = pkt.Reader.ReadString(20)
		userID, _ = pkt.Reader.ReadString(256)
	} else if rem >= 60 {
		uniID, _ = pkt.Reader.ReadString(30)
		userID, _ = pkt.Reader.ReadString(30)
	} else if rem >= 2 {
		cpLen, _ := pkt.Reader.ReadUint16()
		if int(cpLen) <= pkt.Reader.Remaining() {
			userID, _ = pkt.Reader.ReadString(int(cpLen))
		}
	} else {
		userID, _ = pkt.Reader.ReadString(rem)
	}

	uniID = strings.TrimSpace(uniID)
	userID = strings.TrimSpace(userID)
	clientIP := client.RemoteIP()

	log.Printf("[%s] PROTO_NC_USER_LOGIN_NETMARBLE_REQ para usuario '%s' (UniID: '%s')", client.RemoteAddr(), userID, uniID)

	var acc *database.Account
	var err error
	if h.repo != nil {
		acc, err = h.repo.LoginNetmarble(uniID, userID)
	} else {
		err = database.ErrAccountNotFound
	}

	h.executeLoginPipeline(client, acc, userID, "", clientIP, err)
}

// handle_PROTO_NC_USER_LOGIN_OUTSPARK_REQ procesa el login de plataforma Outspark (0x0C20, 85 bytes: LogonToken[65], Spawnapps[20]).
func (h *Handler) handle_PROTO_NC_USER_LOGIN_OUTSPARK_REQ(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var token string
	if rem >= 85 {
		token, _ = pkt.Reader.ReadString(65)
		_, _ = pkt.Reader.ReadString(20) // spawnapps
	} else if rem >= 65 {
		token, _ = pkt.Reader.ReadString(65)
	} else {
		token, _ = pkt.Reader.ReadString(rem)
	}

	token = strings.TrimSpace(token)
	clientIP := client.RemoteIP()

	log.Printf("[%s] PROTO_NC_USER_LOGIN_OUTSPARK_REQ (Token length: %d)", client.RemoteAddr(), len(token))

	if h.repo == nil {
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrInvalidCredentials)
		return
	}

	acc, err := h.repo.LoginOutspark(token, clientIP)
	h.executeLoginPipeline(client, acc, token, "", clientIP, err)
}

// ============================================================================
// 3. Handlers Regionales S2S y Verificación de Contraseña
// ============================================================================

func (h *Handler) handle_PROTO_NC_USER_GER_PASSWORD_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.handleGenericPasswordCheck(client, pkt, protocol.PROTO_NC_USER_GER_PASSWORD_CHECK_ACK)
}

func (h *Handler) handle_PROTO_NC_USER_TW_PASSWORD_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.handleGenericPasswordCheck(client, pkt, protocol.PROTO_NC_USER_TW_PASSWORD_CHECK_ACK)
}

func (h *Handler) handle_PROTO_NC_USER_JP_PASSWORD_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.handleGenericPasswordCheck(client, pkt, protocol.PROTO_NC_USER_JP_PASSWORD_CHECK_ACK)
}

func (h *Handler) handle_PROTO_NC_USER_CH_PASSWORD_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.handleGenericPasswordCheck(client, pkt, protocol.PROTO_NC_USER_CH_PASSWORD_CHECK_ACK)
}

func (h *Handler) handle_PROTO_NC_USER_US_PASSWORD_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.handleGenericPasswordCheck(client, pkt, protocol.PROTO_NC_USER_US_PASSWORD_CHECK_ACK)
}

func (h *Handler) handle_PROTO_NC_USER_PASSWORD_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.handleGenericPasswordCheck(client, pkt, protocol.PROTO_NC_USER_PASSWORD_CHECK_ACK)
}

func (h *Handler) handleGenericPasswordCheck(client *network.Client, pkt *protocol.Packet, ackOpCode uint16) {
	rem := pkt.Reader.Remaining()
	var username, password string

	switch pkt.OpCode {
	case protocol.PROTO_NC_USER_TW_PASSWORD_CHECK_REQ: // 0x0C3E (102B)
		if rem >= 102 {
			_, _ = pkt.Reader.ReadUint16() // Netpacketheader
			username, _ = pkt.Reader.ReadString(64)
			password, _ = pkt.Reader.ReadString(36)
		} else {
			username, _ = pkt.Reader.ReadString(rem)
		}
	case protocol.PROTO_NC_USER_JP_PASSWORD_CHECK_REQ: // 0x0C4B (336B)
		if rem >= 336 {
			_, _ = pkt.Reader.ReadUint16() // Netpacketheader
			username, _ = pkt.Reader.ReadString(24)
			password, _ = pkt.Reader.ReadString(36)
			_, _ = pkt.Reader.ReadString(254) // channel
			_, _ = pkt.Reader.ReadString(20)  // ip
		} else {
			username, _ = pkt.Reader.ReadString(rem)
		}
	case protocol.PROTO_NC_USER_US_PASSWORD_CHECK_REQ: // 0x0C5B (338B)
		if rem >= 338 {
			_, _ = pkt.Reader.ReadUint16() // Netpacketheader
			username, _ = pkt.Reader.ReadString(260)
			password, _ = pkt.Reader.ReadString(36)
			_, _ = pkt.Reader.ReadString(20) // ip
			_, _ = pkt.Reader.ReadString(20) // spawnapps
		} else {
			username, _ = pkt.Reader.ReadString(rem)
		}
	default: // 0x0C07, 0x0C39 (GER), 0x0C53 (CH) - 274B
		if rem >= 274 {
			_, _ = pkt.Reader.ReadUint16() // Netpacketheader
			username, _ = pkt.Reader.ReadString(256)
			password, _ = pkt.Reader.ReadString(16)
		} else if rem >= 20 {
			username, _ = pkt.Reader.ReadString(20)
			if pkt.Reader.Remaining() >= 32 {
				password, _ = pkt.Reader.ReadString(32)
			}
		} else {
			username, _ = pkt.Reader.ReadString(rem)
		}
	}

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	var (
		userNo    int
		nRet      int
		outUserID string
		admin     byte
		active    byte = 1
		block     byte
	)

	if h.repo != nil {
		uNo, ret, uID, adm, act, blk, err := h.repo.CheckPasswordS2S(username, password, true)
		if err == nil {
			userNo = uNo
			nRet = ret
			outUserID = uID
			admin = adm
			active = act
			block = blk
		} else {
			nRet = -1
		}
	}

	if outUserID == "" {
		outUserID = username
	}

	pw := protocol.NewPacketWriter(ackOpCode)
	switch ackOpCode {
	case protocol.PROTO_NC_USER_TW_PASSWORD_CHECK_ACK: // 8B
		pw.WriteUint16(0) // Netpacketheader
		if nRet == 0 {
			pw.WriteUint16(0)
		} else {
			pw.WriteUint16(1)
		}
		pw.WriteUint32(uint32(userNo))

	case protocol.PROTO_NC_USER_JP_PASSWORD_CHECK_ACK: // 306B
		pw.WriteUint16(0) // Netpacketheader
		if nRet == 0 {
			pw.WriteUint16(0)
		} else {
			pw.WriteUint16(1)
		}
		pw.WriteUint32(uint32(userNo))
		pw.WriteString(outUserID, 24)
		pw.Fill(254, 0) // channel
		pw.WriteString(client.RemoteIP(), 20)

	case protocol.PROTO_NC_USER_US_PASSWORD_CHECK_ACK: // 353B
		pw.WriteUint16(0) // Netpacketheader
		if nRet == 0 {
			pw.WriteUint16(0)
		} else {
			pw.WriteUint16(1)
		}
		pw.WriteUint32(uint32(userNo))
		pw.WriteString(outUserID, 260)
		pw.WriteString(password, 36)
		pw.WriteString(client.RemoteIP(), 20)
		pw.WriteByte(admin)
		pw.WriteUint32(uint32(block))
		pw.WriteUint32(uint32(active))
		pw.Fill(20, 0) // spawnapps

	default: // 0x0C08, 0x0C3A (GER), 0x0C54 (CH) - 287B
		pw.WriteUint16(0) // Netpacketheader
		pw.WriteString(outUserID, 256)
		pw.WriteString(password, 16)
		pw.WriteUint32(uint32(userNo))
		pw.WriteByte(admin)
		pw.WriteUint32(uint32(block))
		pw.WriteUint32(uint32(active))
	}

	_ = client.SendPacket(pw)
}

// ============================================================================
// 4. Chequeos de IP Block
// ============================================================================

func (h *Handler) handle_PROTO_NC_USER_IS_IP_BLOCK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.respondIPBlock(client, pkt, protocol.PROTO_NC_USER_IS_IP_BLOCK_ACK)
}

func (h *Handler) handle_PROTO_NC_USER_GER_IS_IP_BLOCK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.respondIPBlock(client, pkt, protocol.PROTO_NC_USER_GER_IS_IP_BLOCK_ACK)
}

func (h *Handler) handle_PROTO_NC_USER_TW_IS_IP_BLOCK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.respondIPBlock(client, pkt, protocol.PROTO_NC_USER_TW_IS_IP_BLOCK_ACK)
}

func (h *Handler) handle_PROTO_NC_USER_JP_IS_IP_BLOCK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.respondIPBlock(client, pkt, protocol.PROTO_NC_USER_JP_IS_IP_BLOCK_ACK)
}

func (h *Handler) handle_PROTO_NC_USER_US_IS_IP_BLOCK_REQ(client *network.Client, pkt *protocol.Packet) {
	h.respondIPBlock(client, pkt, protocol.PROTO_NC_USER_US_IS_IP_BLOCK_ACK)
}

func (h *Handler) respondIPBlock(client *network.Client, pkt *protocol.Packet, ackOpCode uint16) {
	rem := pkt.Reader.Remaining()
	var userNo uint32
	var clientHandle uint16
	ip := client.RemoteIP()

	switch pkt.OpCode {
	case protocol.PROTO_NC_USER_IS_IP_BLOCK_REQ: // 298B (LoginData 272B + UserNo 4B + Handle 2B + IP 20B)
		if rem >= 298 {
			_, _ = pkt.Reader.ReadBytes(272)
			userNo, _ = pkt.Reader.ReadUint32()
			clientHandle, _ = pkt.Reader.ReadUint16()
			sIP, _ := pkt.Reader.ReadString(20)
			if s := strings.TrimSpace(sIP); s != "" {
				ip = s
			}
		}
	case protocol.PROTO_NC_USER_GER_IS_IP_BLOCK_REQ: // 313B (LoginData 287B + UserNo 4B + Handle 2B + IP 20B)
		if rem >= 313 {
			_, _ = pkt.Reader.ReadBytes(287)
			userNo, _ = pkt.Reader.ReadUint32()
			clientHandle, _ = pkt.Reader.ReadUint16()
			sIP, _ := pkt.Reader.ReadString(20)
			if s := strings.TrimSpace(sIP); s != "" {
				ip = s
			}
		}
	case protocol.PROTO_NC_USER_TW_IS_IP_BLOCK_REQ: // 34B (LoginData 8B + UserNo 4B + Handle 2B + IP 20B)
		if rem >= 34 {
			_, _ = pkt.Reader.ReadBytes(8)
			userNo, _ = pkt.Reader.ReadUint32()
			clientHandle, _ = pkt.Reader.ReadUint16()
			sIP, _ := pkt.Reader.ReadString(20)
			if s := strings.TrimSpace(sIP); s != "" {
				ip = s
			}
		}
	case protocol.PROTO_NC_USER_JP_IS_IP_BLOCK_REQ: // 332B (LoginData 306B + UserNo 4B + Handle 2B + IP 20B)
		if rem >= 332 {
			_, _ = pkt.Reader.ReadBytes(306)
			userNo, _ = pkt.Reader.ReadUint32()
			clientHandle, _ = pkt.Reader.ReadUint16()
			sIP, _ := pkt.Reader.ReadString(20)
			if s := strings.TrimSpace(sIP); s != "" {
				ip = s
			}
		}
	case protocol.PROTO_NC_USER_US_IS_IP_BLOCK_REQ: // 379B (LoginData 353B + UserNo 4B + Handle 2B + IP 20B)
		if rem >= 379 {
			_, _ = pkt.Reader.ReadBytes(353)
			userNo, _ = pkt.Reader.ReadUint32()
			clientHandle, _ = pkt.Reader.ReadUint16()
			sIP, _ := pkt.Reader.ReadString(20)
			if s := strings.TrimSpace(sIP); s != "" {
				ip = s
			}
		}
	}

	var isBlocked uint8
	if h.repo != nil {
		b, err := h.repo.IsIPBlocked(int(userNo), ip)
		if err == nil && b {
			isBlocked = 1
		}
	}

	pw := protocol.NewPacketWriter(ackOpCode)
	switch ackOpCode {
	case protocol.PROTO_NC_USER_IS_IP_BLOCK_ACK: // 275B (LoginData 272B + Handle 2B + Block 1B)
		pw.Fill(272, 0)
		pw.WriteUint16(clientHandle)
		pw.WriteByte(isBlocked)
	case protocol.PROTO_NC_USER_GER_IS_IP_BLOCK_ACK: // 290B (LoginData 287B + Handle 2B + Block 1B)
		pw.Fill(287, 0)
		pw.WriteUint16(clientHandle)
		pw.WriteByte(isBlocked)
	case protocol.PROTO_NC_USER_TW_IS_IP_BLOCK_ACK: // 11B (LoginData 8B + Handle 2B + Block 1B)
		pw.Fill(8, 0)
		pw.WriteUint16(clientHandle)
		pw.WriteByte(isBlocked)
	case protocol.PROTO_NC_USER_JP_IS_IP_BLOCK_ACK: // 309B (LoginData 306B + Handle 2B + Block 1B)
		pw.Fill(306, 0)
		pw.WriteUint16(clientHandle)
		pw.WriteByte(isBlocked)
	case protocol.PROTO_NC_USER_US_IS_IP_BLOCK_ACK: // 356B (LoginData 353B + Handle 2B + Block 1B)
		pw.Fill(353, 0)
		pw.WriteUint16(clientHandle)
		pw.WriteByte(isBlocked)
	default:
		pw.WriteByte(isBlocked)
	}

	_ = client.SendPacket(pw)
}

// ============================================================================
// 5. Restricción Teenager
// ============================================================================

func (h *Handler) handle_PROTO_NC_USER_TEENAGER_REQ(client *network.Client, pkt *protocol.Packet) {
	userNo, _ := pkt.Reader.ReadInt32()
	var remainSec int
	if h.repo != nil {
		remainSec, _ = h.repo.GetPlayRemainTime(int(userNo))
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_USER_TEENAGER_ACK)
	pw.WriteInt32(userNo)
	pw.WriteInt32(int32(remainSec))
	_ = client.SendPacket(pw)
}

func (h *Handler) handle_PROTO_NC_USER_TEENAGER_SET_CMD(client *network.Client, pkt *protocol.Packet) {
	userNo, _ := pkt.Reader.ReadInt32()
	remainSec, _ := pkt.Reader.ReadInt32()
	if h.repo != nil {
		_ = h.repo.SetPlayRemainTime(int(userNo), int(remainSec))
	}
}

// ============================================================================
// 6. Selección y Estado de Mundos
// ============================================================================

// handle_PROTO_NC_USER_WORLD_STATUS_REQ refresca el listado dinámico de mundos (0x0C1B).
func (h *Handler) handle_PROTO_NC_USER_WORLD_STATUS_REQ(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		log.Printf("[%s] Refresco de mundos no autenticado rechazado", client.RemoteAddr())
		return
	}
	h.send_PROTO_NC_USER_LOGIN_ACK(client, true)
}

// handle_PROTO_NC_USER_WORLDSELECT_REQ procesa la selección de mundo (0x0C0B, 1B worldno).
func (h *Handler) handle_PROTO_NC_USER_WORLDSELECT_REQ(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.IsTransferring {
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrException)
		return
	}

	worldID, err := pkt.Reader.ReadByte()
	if err != nil {
		log.Printf("[%s] Paquete PROTO_NC_USER_WORLDSELECT_REQ inválido: %v", client.RemoteAddr(), err)
		return
	}

	world, exists := h.worldManager.GetWorld(worldID)
	if !exists {
		log.Printf("[%s] Usuario '%s' seleccionó mundo inexistente: %d", client.RemoteAddr(), client.Username, worldID)
		return
	}

	if world.Status == protocol.WorldStatus_Maintenance || world.Status == protocol.WorldStatus_Offline {
		log.Printf("[%s] Mundo %d fuera de servicio o en mantenimiento", client.RemoteAddr(), worldID)
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrServerMaintenance)
		return
	}
	if world.Status == protocol.WorldStatus_Full {
		log.Printf("[%s] Mundo %d completo", client.RemoteAddr(), worldID)
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrServerFull)
		return
	}

	// Generar token criptográfico oficial de 64 bytes
	var tokenBytes [64]byte
	_, _ = rand.Read(tokenBytes[:])

	if h.sessionRepo != nil {
		_ = h.sessionRepo.IssueWorldToken(client.AccountID, worldID, tokenBytes)
	}

	if h.transfers != nil {
		h.transfers.CreateTransfer(client.AccountID, client.Username, client.Admin, client.RemoteAddr(), tokenBytes[:], worldID)
	}
	client.IsTransferring = true

	log.Printf("[%s] Transfiriendo '%s' a World '%s' (%s:%d)",
		client.RemoteAddr(), client.Username, world.Name, world.IP, world.Port)

	resp := protocol.New_PROTO_NC_USER_WORLDSELECT_ACK_64(byte(world.Status), world.IP, world.Port, tokenBytes)
	_ = client.SendPacket(resp)
}

// ============================================================================
// 7. Desconexión y Heartbeat
// ============================================================================

func (h *Handler) handle_PROTO_NC_MISC_HEARTBEAT_REQ(client *network.Client, pkt *protocol.Packet) {
	_ = client.SendPacket(protocol.New_PROTO_NC_MISC_HEARTBEAT_ACK())
}

func (h *Handler) handle_PROTO_NC_MISC_HEARTBEAT_ACK(client *network.Client, pkt *protocol.Packet) {
	// Acknowledge received from peer
}

func (h *Handler) handle_PROTO_NC_USER_NORMALLOGOUT_CMD(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[%s] PROTO_NC_USER_NORMALLOGOUT_CMD recibido. Cerrando conexión limpiamente.", client.RemoteAddr())
	h.OnClientDisconnect(client)
	client.Close()
}

// OnClientDisconnect realiza el cierre ordenado de sesión y registro de auditoría en MySQL.
func (h *Handler) OnClientDisconnect(client *network.Client) {
	if client.IsAuthenticated {
		if client.IsTransferring {
			// El cliente está traspasándose a WorldServer, no cerramos la sesión en tActiveAccountSession
			client.IsAuthenticated = false
			return
		}
		var playMin uint16
		if h.sessionRepo != nil {
			m, err := h.sessionRepo.CloseSessionAndGetPlayMin(client.AccountID)
			if err == nil {
				playMin = m
			}
		}
		if h.repo != nil {
			_ = h.repo.LogAccount(client.Username, 2, client.AccountID, playMin, "", 0, 0, client.RemoteIP())
		}
		if h.logRepo != nil {
			_ = h.logRepo.AddSyslog(0xC803, client.AccountID, client.Username, client.RemoteIP(), 0, 0, 0)
		}
		client.IsAuthenticated = false
	}
}

// ============================================================================
// 8. Handlers S2S de Topología y Control
// ============================================================================

func (h *Handler) handle_PROTO_NC_MISC_S2SCONNECTION_RDY(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[%s] S2S Connection Ready recibido", client.RemoteAddr())
}

func (h *Handler) handle_PROTO_NC_MISC_S2SCONNECTION_REQ(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var echoData uint8 = 1
	var serverFromWorldNum uint8 = 0
	if rem >= 7 {
		echoData, _ = pkt.Reader.ReadByte()
		_, _ = pkt.Reader.ReadByte() // ServerToID
		_, _ = pkt.Reader.ReadByte() // ServerFromID
		serverFromWorldNum, _ = pkt.Reader.ReadByte()
		_, _ = pkt.Reader.ReadByte() // ServerFromZoneNum
		_, _ = pkt.Reader.ReadUint16() // Key
	} else if rem >= 1 {
		echoData, _ = pkt.Reader.ReadByte()
	}

	if h.sessionRepo != nil && serverFromWorldNum > 0 {
		_ = h.sessionRepo.WorldServerHeartbeat(serverFromWorldNum)
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_MISC_S2SCONNECTION_ACK(echoData, 0))
}

func (h *Handler) handle_PROTO_NC_MISC_USER_COUNT_CMD(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	if rem >= 5 {
		worldID, _ := pkt.Reader.ReadByte()
		userCount, _ := pkt.Reader.ReadUint16()
		maxUser, _ := pkt.Reader.ReadUint16()
		if h.sessionRepo != nil {
			_ = h.sessionRepo.UpdateWorldServerCCUWithMax(worldID, int(userCount), int(maxUser))
		}
	} else if rem >= 4 {
		userCount, _ := pkt.Reader.ReadUint16()
		maxUser, _ := pkt.Reader.ReadUint16()
		if h.sessionRepo != nil {
			_ = h.sessionRepo.UpdateWorldServerCCUWithMax(1, int(userCount), int(maxUser))
		}
	} else if rem >= 3 {
		worldID, _ := pkt.Reader.ReadByte()
		userCount, _ := pkt.Reader.ReadUint16()
		if h.sessionRepo != nil {
			_ = h.sessionRepo.UpdateWorldServerCCU(worldID, int(userCount))
		}
	}
}

func (h *Handler) handle_PROTO_NC_MISC_DELIVER_WM_LOGIN_ACDB_CMD(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[%s] Bridge S2S Account DB command recibido", client.RemoteAddr())
}

func (h *Handler) handle_PROTO_NC_MISC_DELIVER_WM_LOGIN_ALDB_CMD(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[%s] Bridge S2S AccountLog DB command recibido", client.RemoteAddr())
}

func (h *Handler) handle_PROTO_NC_MISC_START_THE_BOOM_CMD(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[%s] START_THE_BOOM_CMD recibido. Activando mantenimiento general en tWorldServer.", client.RemoteAddr())
	if h.sessionRepo != nil {
		_ = h.sessionRepo.SetMaintenanceAll(true)
	}
}

// ============================================================================
// 9. OTP Generation (S2S)
// ============================================================================

func (h *Handler) handle_PROTO_NC_USER_CREATE_OTP_REQ(client *network.Client, pkt *protocol.Packet) {
	clientHandle, _ := pkt.Reader.ReadUint16()
	userNo, _ := pkt.Reader.ReadInt32()

	var otp [32]byte
	rawBytes := make([]byte, 16)
	_, _ = rand.Read(rawBytes)
	hexStr := hex.EncodeToString(rawBytes)
	copy(otp[:], hexStr)

	var errCode uint16
	if h.sessionRepo != nil {
		err := h.sessionRepo.CreateOTP(int(userNo), otp)
		if err != nil {
			errCode = 1
		}
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_USER_CREATE_OTP_ACK(clientHandle, errCode, otp))
}

// ============================================================================
// 10. Retorno y Último Logout
// ============================================================================

func (h *Handler) handle_PROTO_NC_USER_RETURN_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	handle, _ := pkt.Reader.ReadUint16()
	userNo, _ := pkt.Reader.ReadUint32()

	var isEligible byte
	var lastSec int64 = -1

	if h.repo != nil {
		diffLast, errLast := h.repo.GetReturnUserDateDiffLastLog(int(userNo))
		if errLast == nil {
			lastSec = diffLast
			if diffLast == -1 || diffLast >= 30*24*3600 {
				diffCreated, errCreated := h.repo.GetUserDateDiffCreated(int(userNo))
				if errCreated == nil && diffCreated >= 30*24*3600 {
					isEligible = 1
				}
			}
		}
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_USER_RETURN_CHECK_ACK(handle, userNo, isEligible, int32(lastSec)))
}

func (h *Handler) handle_PROTO_NC_USER_SET_RETURN_CMD(client *network.Client, pkt *protocol.Packet) {
	userNo, _ := pkt.Reader.ReadUint32()
	if h.repo != nil {
		_ = h.repo.SetReturnUserItem(int(userNo))
	}
}

func (h *Handler) handle_PROTO_NC_USER_LOGOUT_LAST_DAY_REQ(client *network.Client, pkt *protocol.Packet) {
	handle, _ := pkt.Reader.ReadUint16()
	userNo, _ := pkt.Reader.ReadUint32()

	var lastSec int64 = -1
	if h.logRepo != nil {
		diff, err := h.logRepo.GetDateDiffLastLogout(int(userNo))
		if err == nil {
			lastSec = diff
		}
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_USER_LOGOUT_LAST_DAY_ACK(handle, userNo, int32(lastSec)))
}

// ============================================================================
// Pipeline Unificado de Autenticación Oficial
// ============================================================================

func (h *Handler) executeLoginPipeline(client *network.Client, acc *database.Account, username, password, clientIP string, err error) {
	// 1. Verificar error de credenciales
	if err != nil {
		if errors.Is(err, database.ErrAccountNotFound) || errors.Is(err, database.ErrInvalidPassword) {
			if h.repo != nil {
				_ = h.repo.LogAccount(username, 3, 0, 0, password, 0, uint16(protocol.ErrInvalidCredentials), clientIP)
			}
			if h.logRepo != nil {
				_ = h.logRepo.AddSyslog(0xC802, 0, username, clientIP, 0, 0, uint16(protocol.ErrInvalidCredentials))
			}
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrInvalidCredentials)
			return
		}
		if errors.Is(err, database.ErrAccountBlocked) {
			if h.repo != nil {
				_ = h.repo.LogAccount(username, 3, 0, 0, password, 0, uint16(protocol.ErrBlocked), clientIP)
			}
			if h.logRepo != nil {
				_ = h.logRepo.AddSyslog(0xC802, 0, username, clientIP, 0, 0, uint16(protocol.ErrBlocked))
			}
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrBlocked)
			return
		}
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrDatabase)
		return
	}

	if acc == nil {
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrInvalidCredentials)
		return
	}

	// 2. Verificar IP Block (p_IPBlock_IsBlock) pasando acc.ID y clientIP
	if h.repo != nil {
		blocked, errBlock := h.repo.IsIPBlocked(acc.ID, clientIP)
		if errBlock == nil && blocked {
			log.Printf("[%s] IP bloqueada por tIPBlock: %s para usuario '%s'", client.RemoteAddr(), clientIP, username)
			_ = h.repo.LogAccount(username, 3, acc.ID, 0, "", 0, 73, clientIP)
			if h.logRepo != nil {
				_ = h.logRepo.AddSyslog(0xC802, acc.ID, username, clientIP, 0, 0, 73)
			}
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrBlocked)
			return
		}
	}

	// 3. Evaluar PC Bang y Teenager (restricción horaria)
	if h.repo != nil {
		_, _ = h.repo.CheckPcbangIP(clientIP)
		remainSec, errRemain := h.repo.GetPlayRemainTime(acc.ID)
		if errRemain == nil && remainSec == 0 {
			log.Printf("[%s] Cuenta teenager sin saldo de juego para '%s'", client.RemoteAddr(), username)
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrServerFull)
			return
		}
	}

	// 4. Registrar sesión activa en tActiveAccountSession y auditar éxito en tAccountLog + tSyslog
	if h.sessionRepo != nil {
		wasActive, _ := h.sessionRepo.RegisterLoginSession(acc.ID, acc.Username, clientIP)
		if wasActive {
			log.Printf("[%s] Sesión previa activa detectada para UserNo %d, reemplazada en tActiveAccountSession",
				client.RemoteAddr(), acc.ID)
		}
	}
	if h.repo != nil {
		_ = h.repo.LogAccount(username, 1, acc.ID, 0, "", 0, 0, clientIP)
	}
	if h.logRepo != nil {
		_ = h.logRepo.AddSyslog(0xC801, acc.ID, username, clientIP, 0, 0, 0)
	}

	// 5. Configurar cliente autenticado
	client.AccountID = acc.ID
	client.Username = acc.Username
	client.Admin = acc.Admin
	client.IsAuthenticated = true

	log.Printf("[%s] Autenticación oficial exitosa para '%s' (UserNo: %d, AuthID: %d)",
		client.RemoteAddr(), client.Username, client.AccountID, client.Admin)

	// 6. Consultar si puede omitir el tutorial y enviar PROTO_NC_USER_TUTORIAL_CAN_SKIP_CMD si aplica
	if h.repo != nil {
		canSkip, errSkip := h.repo.GetTutorialIsSkip(acc.ID)
		if errSkip == nil && canSkip {
			_ = client.SendPacket(protocol.New_PROTO_NC_USER_TUTORIAL_CAN_SKIP_CMD(true))
		}
	}

	// 7. Enviar listado de mundos oficiales (PROTO_NC_USER_LOGIN_ACK 0x0C0A)
	h.send_PROTO_NC_USER_LOGIN_ACK(client, false)
}

// send_PROTO_NC_USER_LOGIN_ACK emite el listado de mundos oficiales (0x0C0A / 0x0C1C).
func (h *Handler) send_PROTO_NC_USER_LOGIN_ACK(client *network.Client, isResend bool) {
	opCode := uint16(protocol.PROTO_NC_USER_LOGIN_ACK)
	if isResend {
		opCode = protocol.PROTO_NC_USER_WORLD_STATUS_ACK
	}

	worlds := h.worldManager.WorldsList()
	worldInfos := make([]protocol.PROTO_NC_USER_LOGIN_ACK_WorldInfo, 0, len(worlds))
	for _, w := range worlds {
		worldInfos = append(worldInfos, protocol.PROTO_NC_USER_LOGIN_ACK_WorldInfo{
			WorldNo:     w.ID,
			WorldName:   w.Name,
			WorldStatus: byte(w.Status),
		})
	}

	resp := protocol.New_PROTO_NC_USER_LOGIN_ACK(opCode, worldInfos)
	_ = client.SendPacket(resp)
}

// send_PROTO_NC_USER_LOGINFAIL_ACK emite el paquete oficial de rechazo de login (0x0C09).
func (h *Handler) send_PROTO_NC_USER_LOGINFAIL_ACK(client *network.Client, errCode protocol.ServerError) {
	resp := protocol.New_PROTO_NC_USER_LOGINFAIL_ACK(uint16(errCode))
	_ = client.SendPacket(resp)
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
