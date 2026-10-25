package login

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler procesa los paquetes entrantes para LoginService respetando estrictamente
// la nomenclatura y estructuras oficiales de red (docs/Servicios/Login.md).
type Handler struct {
	repo          database.AccountRepository
	worldManager  *WorldManager
	transfers     *database.TransferManager
	clientVersion string
}

// NewHandler inicializa el despachador de protocolos oficiales de LoginService.
func NewHandler(repo database.AccountRepository, wm *WorldManager, tm *database.TransferManager) *Handler {
	if tm == nil {
		tm = database.NewTransferManager()
	}
	return &Handler{
		repo:         repo,
		worldManager: wm,
		transfers:    tm,
	}
}

// SetClientVersion configura la versión oficial esperada del cliente (ClientVersionKeyInfo.txt).
func (h *Handler) SetClientVersion(ver string) {
	h.clientVersion = strings.TrimSpace(ver)
}

// HandlePacket enruta cada paquete entrante según su OpCode oficial PROTO_NC_*.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	case protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ, protocol.PROTO_NC_USER_NA2026_VERSION_CHECK_REQ:
		h.handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_US_LOGIN_REQ, protocol.PROTO_NC_USER_LOGIN_REQ, protocol.PROTO_NC_USER_LOGIN_WITH_OTP_REQ, protocol.PROTO_NC_USER_LOGIN_WITH_OTP_LEGACY, protocol.PROTO_NC_USER_NA2026_LOGIN_REQ:
		h.handle_PROTO_NC_USER_US_LOGIN_REQ(client, pkt)
	case protocol.PROTO_NC_USER_FILE_CHECK_REQ:
		h.handle_PROTO_NC_USER_FILE_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_WORLDSTATUS_REQ:
		h.handle_PROTO_NC_USER_WORLDSTATUS_REQ(client, pkt)
	case protocol.PROTO_NC_USER_WORLDSELECT_REQ, protocol.PROTO_NC_USER_NA2026_WORLDSELECT_REQ, protocol.PROTO_NC_USER_LOGINWORLDFAIL_ACK:
		h.handle_PROTO_NC_USER_WORLDSELECT_REQ(client, pkt)
	case protocol.PROTO_NC_USER_NA2026_SESSION_PING:
		h.handle_PROTO_NC_USER_NA2026_SESSION_PING(client, pkt)
	case protocol.PROTO_NC_USER_NA2026_TOS_ACCEPT_HASH_REQ:
		h.handle_PROTO_NC_USER_NA2026_TOS_ACCEPT_HASH_REQ(client, pkt)
	case protocol.PROTO_NC_USER_NA2026_TOS_CONFIRM_REQ:
		h.handle_PROTO_NC_USER_NA2026_TOS_CONFIRM_REQ(client, pkt)
	case protocol.PROTO_NC_MISC_HEARTBEAT_REQ:
		h.handle_PROTO_NC_MISC_HEARTBEAT_REQ(client, pkt)
	case protocol.PROTO_NC_USER_NORMALLOGOUT_CMD:
		// En NA2026, si el cliente ya está autenticado, 0x0C18 se usa para solicitar/refrescar la lista de mundos tras TOS
		if client.IsAuthenticated {
			h.handle_PROTO_NC_USER_WORLDSTATUS_REQ(client, pkt)
		} else {
			h.handle_PROTO_NC_USER_NORMALLOGOUT_CMD(client, pkt)
		}
	default:
		log.Printf("[%s] OpCode no manejado en LoginService: 0x%04X (Dept=%d, Cmd=%d) Len=%d\n%s",
			client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type, len(pkt.Raw),
			hex.Dump(pkt.Raw))
	}
}

// handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ procesa el handshake de versión del cliente (0x0C65 o 0x0C2C en NA2026).
func (h *Handler) handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var versionKey string

	// En NA2026 el opcode es 0x0C2C y envía 32 bytes de hash MD5 ASCII
	if pkt.OpCode == protocol.PROTO_NC_USER_NA2026_VERSION_CHECK_REQ {
		client.ClientYear = 2026
		client.ClientVersion = 2026
		if rem >= 32 {
			s, _ := pkt.Reader.ReadString(32)
			versionKey = strings.TrimSpace(s)
		} else if rem > 0 {
			s, _ := pkt.Reader.ReadString(rem)
			versionKey = strings.TrimSpace(s)
		}
		log.Printf("[%s] PROTO_NC_USER_NA2026_VERSION_CHECK_REQ (0x0C2C) recibido (Hash='%s')", client.RemoteAddr(), versionKey)

		// Validar versión contra CLIENT_VERSION en .ENV
		if h.clientVersion != "" && !matchClientVersion(h.clientVersion, versionKey) {
			log.Printf("[%s] Versión NA2026 no coincide: '%s' (Esperada en .ENV: '%s'). Rechazando conexión con ACK de error.",
				client.RemoteAddr(), versionKey, h.clientVersion)
			resp := protocol.New_PROTO_NC_USER_NA2026_VERSION_CHECK_ACK(1) // Status != 0 indica versión incorrecta
			_ = client.SendPacket(resp)
			client.Close()
			return
		}

		// Versión correcta: responder con 0x0C2D (3 bytes wire: 2D 0C 00)
		resp := protocol.New_PROTO_NC_USER_NA2026_VERSION_CHECK_ACK(0)
		_ = client.SendPacket(resp)
		return
	}

	if rem >= 64 {
		// Estructura oficial C++: char sVersionKey[64]
		s, err := pkt.Reader.ReadString(64)
		if err != nil {
			log.Printf("[%s] Error leyendo sVersionKey: %v", client.RemoteAddr(), err)
			_ = client.SendPacket(protocol.New_PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK())
			client.Close()
			return
		}
		versionKey = strings.TrimSpace(s)
	} else if rem == 4 {
		// Compatibilidad legacy (2 uint16: Year, Version)
		year, _ := pkt.Reader.ReadUint16()
		ver, _ := pkt.Reader.ReadUint16()
		client.ClientYear = year
		client.ClientVersion = ver
		versionKey = fmt.Sprintf("%d%02d", year, ver)
	} else if rem > 0 {
		versionKey, _ = pkt.Reader.ReadString(rem)
		versionKey = strings.TrimSpace(versionKey)
	}

	log.Printf("[%s] PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ recibido (Key='%s')", client.RemoteAddr(), versionKey)

	// Validar versión si está configurada en .ENV (CLIENT_VERSION)
	if h.clientVersion != "" && !matchClientVersion(h.clientVersion, versionKey) {
		log.Printf("[%s] Versión de cliente no coincide: '%s' (Esperada en .ENV: '%s'). Rechazando conexión.",
			client.RemoteAddr(), versionKey, h.clientVersion)
		_ = client.SendPacket(protocol.New_PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK())
		client.Close()
		return
	}

	// Responder con éxito de verificación de versión: PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK (0x0C67)
	resp := protocol.New_PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK()
	if err := client.SendPacket(resp); err != nil {
		log.Printf("[%s] Error enviando PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK: %v", client.RemoteAddr(), err)
	}
}

// matchClientVersion compara la versión configurada (que puede contener múltiples versiones separadas por coma)
// con la versión reportada por el cliente.
func matchClientVersion(configured, received string) bool {
	if configured == "" {
		return true
	}
	parts := strings.Split(configured, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) == received {
			return true
		}
	}
	return false
}

// handle_PROTO_NC_USER_US_LOGIN_REQ procesa la solicitud de autenticación (0x0C5A / 0x0C03 / 0x0C37 / 0x0C38).
// Soporta la estructura oficial (sUserName[260], sPassword[36], spawnapps[20] = 316 bytes) y formatos legacy.
func (h *Handler) handle_PROTO_NC_USER_US_LOGIN_REQ(client *network.Client, pkt *protocol.Packet) {
	var (
		username  string
		password  string
		spawnApps string
		err       error
	)

	rem := pkt.Reader.Remaining()

	// Si el paquete tiene ~351 bytes (PROTO_NC_USER_LOGIN_REQ 0x0C01) NA2026
	if pkt.OpCode == protocol.PROTO_NC_USER_NA2026_LOGIN_REQ || rem >= 340 {
		client.ClientYear = 2026
		client.ClientVersion = 2026
		// 1. Saltar padding inicial de sesión (32 bytes)
		if _, err = pkt.Reader.ReadBytes(32); err != nil {
			log.Printf("[%s] Error leyendo padding inicial: %v", client.RemoteAddr(), err)
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrException)
			return
		}
		// 2. Leer usuario (260 bytes)
		username, err = pkt.Reader.ReadString(260)
		if err != nil {
			log.Printf("[%s] Error leyendo sUserName: %v", client.RemoteAddr(), err)
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrException)
			return
		}
		// 3. Leer password hash MD5 (32 bytes)
		password, err = pkt.Reader.ReadString(32)
		if err != nil {
			log.Printf("[%s] Error leyendo sPassword: %v", client.RemoteAddr(), err)
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrException)
			return
		}
		// 4. Leer client string / spawnApps restante si existe (ej. "Original")
		if pkt.Reader.Remaining() > 0 {
			spawnApps, _ = pkt.Reader.ReadString(pkt.Reader.Remaining())
		}
	} else if rem >= 292 {
		// Estructura US_LOGIN clásica sin padding inicial
		username, err = pkt.Reader.ReadString(260)
		if err != nil {
			log.Printf("[%s] Error leyendo sUserName: %v", client.RemoteAddr(), err)
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrException)
			return
		}
		if pkt.Reader.Remaining() >= 36 {
			password, err = pkt.Reader.ReadString(36)
		} else {
			password, err = pkt.Reader.ReadString(32)
		}
		if err != nil {
			log.Printf("[%s] Error leyendo sPassword: %v", client.RemoteAddr(), err)
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrException)
			return
		}
		if pkt.Reader.Remaining() >= 20 {
			spawnApps, _ = pkt.Reader.ReadString(20)
		} else if pkt.Reader.Remaining() >= 8 {
			spawnApps, _ = pkt.Reader.ReadString(8)
		}
	} else if rem >= 32 {
		username, _ = pkt.Reader.ReadString(16)
		password, _ = pkt.Reader.ReadString(rem - 16)
	} else {
		username, _ = pkt.Reader.ReadString(rem)
	}

	username = strings.TrimSpace(strings.Trim(username, "\x00"))
	password = strings.TrimSpace(strings.Trim(password, "\x00"))
	spawnApps = strings.TrimSpace(strings.Trim(spawnApps, "\x00"))
	clientIP := client.RemoteIP()

	log.Printf("[%s] Intento de login para usuario '%s' (spawnapps: '%s')", client.RemoteAddr(), username, spawnApps)

	if len(username) == 0 || len(username) > 60 {
		log.Printf("[%s] Longitud de usuario inválida", client.RemoteAddr())
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrBlocked)
		return
	}

	// 1. Verificar si la IP está en el firewall tIPBlock
	if mysqlRepo, ok := h.repo.(*database.MySQLAccountRepository); ok {
		blocked, err := mysqlRepo.IsIPBlocked(clientIP)
		if err == nil && blocked {
			log.Printf("[%s] IP bloqueada por tIPBlock: %s", client.RemoteAddr(), clientIP)
			_ = mysqlRepo.LogAccount(username, 2, 0, 0, uint16(protocol.ErrBlocked), clientIP)
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrBlocked)
			return
		}
	}

	// 2. Consulta en base de datos oficial tUser (usp_User_USloginGame)
	acc, err := h.repo.GetAccountByUsername(username)
	if err != nil {
		if errors.Is(err, database.ErrAccountNotFound) {
			log.Printf("[%s] Usuario '%s' no encontrado", client.RemoteAddr(), username)
			if mysqlRepo, ok := h.repo.(*database.MySQLAccountRepository); ok {
				_ = mysqlRepo.LogAccount(username, 2, 0, 0, uint16(protocol.ErrInvalidCredentials), clientIP)
			}
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrInvalidCredentials)
			return
		}

		log.Printf("[%s] Error en base de datos al autenticar: %v", client.RemoteAddr(), err)
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrDatabase)
		return
	}

	// 3. Validar credenciales según el estándar estricto oficial de Fiesta Online:
	// El cliente envía el hash MD5 (32 caracteres hexadecimales) y tUser.sUserPW almacena char(32) MD5.
	{
		clientPassHash := strings.ToLower(password)
		if len(clientPassHash) != 32 {
			clientPassHash = strings.ToLower(md5Hash(password))
		}
		accPassHash := strings.ToLower(acc.Password)
		if len(accPassHash) != 32 {
			accPassHash = strings.ToLower(md5Hash(acc.Password))
		}

		if clientPassHash != accPassHash {
			log.Printf("[%s] Contraseña incorrecta para usuario '%s'", client.RemoteAddr(), username)
			if mysqlRepo, ok := h.repo.(*database.MySQLAccountRepository); ok {
				_ = mysqlRepo.LogAccount(username, 2, acc.ID, 0, uint16(protocol.ErrInvalidCredentials), clientIP)
			}
			h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrInvalidCredentials)
			return
		}
	}

	// 4. Validar estado de la cuenta (bloqueada, eliminada o sin permisos)
	if acc.Banned {
		log.Printf("[%s] Cuenta '%s' bloqueada o inhabilitada", client.RemoteAddr(), username)
		if mysqlRepo, ok := h.repo.(*database.MySQLAccountRepository); ok {
			_ = mysqlRepo.LogAccount(username, 2, acc.ID, 0, uint16(protocol.ErrBlocked), clientIP)
		}
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrBlocked)
		return
	}

	// 5. Autenticación exitosa: registrar auditoría y actualizar sesión en tUser
	if mysqlRepo, ok := h.repo.(*database.MySQLAccountRepository); ok {
		_ = mysqlRepo.UpdateLoginSession(acc.ID, clientIP)
		_ = mysqlRepo.LogAccount(username, 1, acc.ID, 0, 0, clientIP)

		// 5a. Detección y premio para jugador retornante (usp_ReturnUser_GetDatediffLastlog & usp_ReturnUser_SetUserItem)
		diffLastLog, err := mysqlRepo.GetReturnUserDateDiffLastLog(acc.ID)
		if err == nil {
			// Si nunca ha recibido el pack (-1) o pasaron más de 30 días (2592000 seg) desde el último retorno
			if diffLastLog == -1 || diffLastLog >= 30*24*3600 {
				diffCreated, errCreated := mysqlRepo.GetUserDateDiffCreated(acc.ID)
				if errCreated == nil && diffCreated >= 30*24*3600 {
					log.Printf("[%s] Jugador retornante '%s' elegible (Inactivo %d s). Asignando paquete de bienvenida.",
						client.RemoteAddr(), username, diffLastLog)
					_ = mysqlRepo.SetReturnUserItem(acc.ID, 30001)
				}
			}
		}
	}

	client.AccountID = acc.ID
	client.Username = acc.Username
	client.Admin = acc.Admin
	client.IsAuthenticated = true

	log.Printf("[%s] Autenticación oficial exitosa para '%s' (UserNo: %d, AuthID: %d)",
		client.RemoteAddr(), client.Username, client.AccountID, client.Admin)

	// 6. Enviar secuencia oficial de autenticación:
	// Si el cliente conectó con NA2026 (0x0C01), enviar la secuencia oficial capturada:
	// 0x0C40 (OK) -> 0x0C47 -> 0x0C06 (Mundos) -> 0x0C3E (Tokens)
	if pkt.OpCode == protocol.PROTO_NC_USER_NA2026_LOGIN_REQ {
		_ = client.SendPacket(protocol.New_PROTO_NC_USER_NA2026_AUTH_SUCCESS_ACK())
		_ = client.SendPacket(protocol.New_PROTO_NC_USER_NA2026_PRE_WORLDS_ACK())
		h.send_PROTO_NC_USER_LOGIN_ACK(client, false)
		_ = client.SendPacket(protocol.New_PROTO_NC_USER_NA2026_ACCOUNT_TOKEN_ACK(acc.ID))
		return
	}

	// Flujo legacy estándar (0x0C0A)
	h.send_PROTO_NC_USER_LOGIN_ACK(client, false)
}

// handle_PROTO_NC_USER_NA2026_SESSION_PING responde al ping de sesión en ventana de TOS (0x0C34 -> 0x0C35).
func (h *Handler) handle_PROTO_NC_USER_NA2026_SESSION_PING(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	data, _ := pkt.Reader.ReadBytes(rem)
	resp := protocol.New_PROTO_NC_USER_NA2026_SESSION_PONG(data)
	_ = client.SendPacket(resp)
}

// handle_PROTO_NC_USER_NA2026_TOS_ACCEPT_HASH_REQ procesa la aceptación de términos (0x0C3F).
func (h *Handler) handle_PROTO_NC_USER_NA2026_TOS_ACCEPT_HASH_REQ(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[%s] PROTO_NC_USER_NA2026_TOS_ACCEPT_HASH_REQ (0x0C3F) recibido y registrado", client.RemoteAddr())
}

// handle_PROTO_NC_USER_NA2026_TOS_CONFIRM_REQ procesa la confirmación de lectura de TOS (0x0C41).
func (h *Handler) handle_PROTO_NC_USER_NA2026_TOS_CONFIRM_REQ(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[%s] PROTO_NC_USER_NA2026_TOS_CONFIRM_REQ (0x0C41) recibido. Enviando listado de mundos activo.", client.RemoteAddr())
	h.send_PROTO_NC_USER_LOGIN_ACK(client, true)
}

// handle_PROTO_NC_MISC_HEARTBEAT_REQ procesa el keep-alive heartbeat del cliente (0x0804).
func (h *Handler) handle_PROTO_NC_MISC_HEARTBEAT_REQ(client *network.Client, pkt *protocol.Packet) {
	_ = client.SendPacket(protocol.New_PROTO_NC_MISC_HEARTBEAT_ACK())
}

// handle_PROTO_NC_USER_NORMALLOGOUT_CMD procesa la desconexión normal solicitada por el cliente (0x0C18).
func (h *Handler) handle_PROTO_NC_USER_NORMALLOGOUT_CMD(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[%s] PROTO_NC_USER_NORMALLOGOUT_CMD recibido. Cerrando conexión limpiamente.", client.RemoteAddr())
	client.Close()
}

// handle_PROTO_NC_USER_FILE_CHECK_REQ procesa la verificación de archivos del cliente (0x0C04).
func (h *Handler) handle_PROTO_NC_USER_FILE_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.New_PROTO_NC_USER_FILE_CHECK_ACK(true)
	_ = client.SendPacket(resp)
}

// handle_PROTO_NC_USER_WORLDSTATUS_REQ procesa la solicitud de refresco del listado de mundos (0x0C1B).
func (h *Handler) handle_PROTO_NC_USER_WORLDSTATUS_REQ(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		log.Printf("[%s] Refresco de mundos no autenticado rechazado", client.RemoteAddr())
		return
	}
	h.send_PROTO_NC_USER_LOGIN_ACK(client, true)
}

// handle_PROTO_NC_USER_WORLDSELECT_REQ procesa la selección de mundo (0x0C0B, 1 byte worldno).
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
		// Fallback automático al mundo principal (World 8 para NA2026 o World 0) si el ID solicitado no está explícitamente registrado
		if w8, ok := h.worldManager.GetWorld(8); ok {
			world = w8
			exists = true
		} else if w0, ok := h.worldManager.GetWorld(0); ok {
			world = w0
			exists = true
		} else if len(h.worldManager.WorldsList()) > 0 {
			world = h.worldManager.WorldsList()[0]
			exists = true
		} else {
			log.Printf("[%s] Usuario '%s' seleccionó un mundo inexistente: %d", client.RemoteAddr(), client.Username, worldID)
			return
		}
	}

	if world.Status == protocol.WorldStatus_Maintenance || world.Status == protocol.WorldStatus_Offline {
		log.Printf("[%s] Mundo %d fuera de servicio o en mantenimiento", client.RemoteAddr(), worldID)
		h.send_PROTO_NC_USER_LOGINFAIL_ACK(client, protocol.ErrServerMaintenance)
		return
	}

	token := h.transfers.CreateTransfer(client.AccountID, client.Username, client.Admin, client.RemoteAddr())
	client.IsTransferring = true

	log.Printf("[%s] Transfiriendo '%s' a World '%s' (%s:%d)",
		client.RemoteAddr(), client.Username, world.Name, world.IP, world.Port)

	// Si el cliente usa NA2026 (0x0C0A o 0x0C15), responder con 0x0C0B (86 bytes)
	if pkt.OpCode == protocol.PROTO_NC_USER_NA2026_WORLDSELECT_REQ || pkt.OpCode == protocol.PROTO_NC_USER_LOGINWORLDFAIL_ACK || client.ClientYear == 2026 {
		resp := protocol.New_PROTO_NC_USER_NA2026_WORLDSELECT_ACK(byte(world.Status), world.IP, world.Port, token, worldID)
		_ = client.SendPacket(resp)
		return
	}

	// Respuesta oficial legacy: PROTO_NC_USER_WORLDSELECT_ACK (0x0C0C, 83 bytes)
	resp := protocol.New_PROTO_NC_USER_WORLDSELECT_ACK(byte(world.Status), world.IP, world.Port, token)
	_ = client.SendPacket(resp)
}

// send_PROTO_NC_USER_LOGIN_ACK emite el listado de mundos activos (0x0C0A / 0x0C1C en legacy, o 0x0C06 / 0x0C19 en NA2026).
func (h *Handler) send_PROTO_NC_USER_LOGIN_ACK(client *network.Client, isResend bool) {
	// Determinar si el cliente es NA2026 según el año/versión negociado o la estructura de login
	isNA2026 := client.ClientYear == 2026 || client.ClientVersion == 2026

	var opCode uint16
	if isNA2026 {
		opCode = protocol.PROTO_NC_USER_NA2026_LOGIN_ACK // 0x0C06
		if isResend {
			opCode = protocol.PROTO_NC_USER_NA2026_WORLDSTATUS_ACK // 0x0C19
		}
	} else {
		opCode = protocol.PROTO_NC_USER_LOGIN_ACK // 0x0C0A
		if isResend {
			opCode = protocol.PROTO_NC_USER_WORLDSTATUS_ACK // 0x0C1C
		}
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

	if isNA2026 {
		// Para el cliente NA2026, aseguramos la estructura oficial de 16 mundos con World 8 (Enid ES / Xelion) activo
		worldInfos = protocol.BuildNA2026OfficialWorldList(worldInfos)
	}

	resp := protocol.New_PROTO_NC_USER_LOGIN_ACK(opCode, worldInfos)
	_ = client.SendPacket(resp)
}

// send_PROTO_NC_USER_LOGINFAIL_ACK emite el paquete oficial de rechazo de login (0x0C09 o 0x0C07 en NA2026).
func (h *Handler) send_PROTO_NC_USER_LOGINFAIL_ACK(client *network.Client, errCode protocol.ServerError) {
	if client.ClientYear == 2026 || client.ClientVersion == 2026 {
		// En NA2026 se envía 0x0C07 con uint16 (0x0045)
		respNA := protocol.NewPacketWriter(protocol.PROTO_NC_USER_NA2026_LOGINFAIL_ACK)
		respNA.WriteUint16(0x0045) // Código oficial capturado en wrong_pass
		_ = client.SendPacket(respNA)
		return
	}

	// Estándar oficial legacy 0x0C09
	resp := protocol.New_PROTO_NC_USER_LOGINFAIL_ACK(uint16(errCode))
	_ = client.SendPacket(resp)
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
