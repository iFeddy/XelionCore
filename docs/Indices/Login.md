# Índice: LoginService

> Generado automáticamente por scripts/gen_indices. Solo firmas con n° de línea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. 

## cmd/LoginService/main.go
- L20: `func main()`

## pkg/server/login/e2e_login_test.go
- L20: `type mockClient struct`
- L25: `func readServerPacket(conn net.Conn) (*protocol.Packet, error)`
- L50: `func (c *mockClient) sendPacket(pw *protocol.PacketWriter) error`
- L72: `func TestE2ELoginServiceFlow(t *testing.T)`
- L253: `func TestLoginService_HeartbeatAndLogout(t *testing.T)`

## pkg/server/login/handler.go
- L18: `type Handler struct`
- L26: `func NewHandler(repo database.AccountRepository, wm *WorldManager, tm *database.TransferManager) *Handler`
- L38: `func (h *Handler) SetClientVersion(ver string)`
- L43: `func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet)`
- L78: `func (h *Handler) handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ(client *network.Client, pkt *protocol.Packet)`
- L153: `func matchClientVersion(configured, received string) bool`
- L168: `func (h *Handler) handle_PROTO_NC_USER_US_LOGIN_REQ(client *network.Client, pkt *protocol.Packet)`
- L353: `func (h *Handler) handle_PROTO_NC_USER_NA2026_SESSION_PING(client *network.Client, pkt *protocol.Packet)`
- L361: `func (h *Handler) handle_PROTO_NC_USER_NA2026_TOS_ACCEPT_HASH_REQ(client *network.Client, pkt *protocol.Packet)`
- L366: `func (h *Handler) handle_PROTO_NC_USER_NA2026_TOS_CONFIRM_REQ(client *network.Client, pkt *protocol.Packet)`
- L372: `func (h *Handler) handle_PROTO_NC_MISC_HEARTBEAT_REQ(client *network.Client, pkt *protocol.Packet)`
- L377: `func (h *Handler) handle_PROTO_NC_USER_NORMALLOGOUT_CMD(client *network.Client, pkt *protocol.Packet)`
- L383: `func (h *Handler) handle_PROTO_NC_USER_FILE_CHECK_REQ(client *network.Client, pkt *protocol.Packet)`
- L389: `func (h *Handler) handle_PROTO_NC_USER_WORLDSTATUS_REQ(client *network.Client, pkt *protocol.Packet)`
- L398: `func (h *Handler) handle_PROTO_NC_USER_WORLDSELECT_REQ(client *network.Client, pkt *protocol.Packet)`
- L453: `func (h *Handler) send_PROTO_NC_USER_LOGIN_ACK(client *network.Client, isResend bool)`
- L490: `func (h *Handler) send_PROTO_NC_USER_LOGINFAIL_ACK(client *network.Client, errCode protocol.ServerError)`
- L504: `func md5Hash(text string) string`

## pkg/server/login/na2026_test.go
- L13: `func TestNA2026_LoginFlow(t *testing.T)`
- L128: `func TestLogin_VersionValidationMismatch(t *testing.T)`
- L213: `func readServerPacket(conn net.Conn) (*protocol.Packet, error)`
- L241: `func sendClientPacket(conn net.Conn, pw *protocol.PacketWriter, nc *crypto.NetCrypto) error`

## pkg/server/login/server.go
- L12: `type Config struct`
- L19: `type Server struct`
- L28: `func NewServer(cfg *Config, repo database.AccountRepository, wm *WorldManager, tm *database.TransferManager) *Server`
- L51: `func (s *Server) Start() error`
- L57: `func (s *Server) Stop()`

## pkg/server/login/world_manager.go
- L12: `type World struct`
- L21: `type WorldManager struct`
- L27: `func NewWorldManager() *WorldManager`
- L45: `func (wm *WorldManager) AddWorld(w *World)`
- L52: `func (wm *WorldManager) GetWorld(id byte) (*World, bool)`
- L60: `func (wm *WorldManager) WorldsList() []*World`
- L71: `func GenerateTransferToken() string`
