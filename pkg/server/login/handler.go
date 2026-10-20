package login

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
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
	case protocol.PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ:
		h.handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_US_LOGIN_REQ, protocol.PROTO_NC_USER_LOGIN_REQ, protocol.PROTO_NC_USER_LOGIN_WITH_OTP_REQ, protocol.PROTO_NC_USER_LOGIN_WITH_OTP_LEGACY:
		h.handle_PROTO_NC_USER_US_LOGIN_REQ(client, pkt)
	case protocol.PROTO_NC_USER_FILE_CHECK_REQ:
		h.handle_PROTO_NC_USER_FILE_CHECK_REQ(client, pkt)
	case protocol.PROTO_NC_USER_WORLDSTATUS_REQ:
		h.handle_PROTO_NC_USER_WORLDSTATUS_REQ(client, pkt)
	case protocol.PROTO_NC_USER_WORLDSELECT_REQ:
		h.handle_PROTO_NC_USER_WORLDSELECT_REQ(client, pkt)
	case protocol.PROTO_NC_MISC_HEARTBEAT_REQ:
		h.handle_PROTO_NC_MISC_HEARTBEAT_REQ(client, pkt)
	case protocol.PROTO_NC_USER_NORMALLOGOUT_CMD:
		h.handle_PROTO_NC_USER_NORMALLOGOUT_CMD(client, pkt)
	default:
		log.Printf("[%s] OpCode no manejado en LoginService: 0x%04X (Dept=%d, Cmd=%d)",
			client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type)
	}
}

// handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ procesa el handshake de versión del cliente (0x0C65).
func (h *Handler) handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var versionKey string

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
	} else if rem > 0 {
		versionKey, _ = pkt.Reader.ReadString(rem)
		versionKey = strings.TrimSpace(versionKey)
	}

	log.Printf("[%s] PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ recibido (Key='%s')", client.RemoteAddr(), versionKey)

	// Validar versión si está configurada (ClientVersionKeyInfo.txt)
	if h.clientVersion != "" && versionKey != "" && versionKey != h.clientVersion {
		log.Printf("[%s] Versión de cliente no coincide: '%s' (Esperada: '%s'). Rechazando conexión.",
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
	if rem >= 292 {
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

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
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

	// 6. Enviar listado de mundos oficiales (PROTO_NC_USER_LOGIN_ACK 0x0C0A)
	h.send_PROTO_NC_USER_LOGIN_ACK(client, false)
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
		log.Printf("[%s] Usuario '%s' seleccionó un mundo inexistente: %d", client.RemoteAddr(), client.Username, worldID)
		return
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

	// Respuesta oficial: PROTO_NC_USER_WORLDSELECT_ACK (0x0C0C, 83 bytes)
	resp := protocol.New_PROTO_NC_USER_WORLDSELECT_ACK(byte(world.Status), world.IP, world.Port, token)
	_ = client.SendPacket(resp)
}

// send_PROTO_NC_USER_LOGIN_ACK emite el listado de mundos activos (0x0C0A / 0x0C1C).
func (h *Handler) send_PROTO_NC_USER_LOGIN_ACK(client *network.Client, isResend bool) {
	opCode := uint16(protocol.PROTO_NC_USER_LOGIN_ACK)
	if isResend {
		opCode = protocol.PROTO_NC_USER_WORLDSTATUS_ACK
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
