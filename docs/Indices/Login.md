# Ãndice: LoginService

> Generado automÃ¡ticamente por scripts/gen_indices.ps1. Solo firmas con nÂ° de lÃ­nea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. 

## pkg/server/login/e2e_login_test.go
- L20: `type mockClient struct`
- L25: `func readServerPacket(conn net.Conn) (*protocol.Packet, error)`
- L50: `func (c *mockClient) sendPacket(pw *protocol.PacketWriter) error`
- L72: `func TestE2ELoginServiceFlow(t *testing.T)`
- L253: `func TestLoginService_HeartbeatAndLogout(t *testing.T)`

## pkg/server/login/handler.go
- L17: `type Handler struct`
- L25: `func NewHandler(repo database.AccountRepository, wm *WorldManager, tm *database.TransferManager) *Handler`
- L37: `func (h *Handler) SetClientVersion(ver string)`
- L42: `func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet)`
- L66: `func (h *Handler) handle_PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ(client *network.Client, pkt *protocol.Packet)`
- L111: `func (h *Handler) handle_PROTO_NC_USER_US_LOGIN_REQ(client *network.Client, pkt *protocol.Packet)`
- L283: `func (h *Handler) handle_PROTO_NC_MISC_HEARTBEAT_REQ(client *network.Client, pkt *protocol.Packet)`
- L288: `func (h *Handler) handle_PROTO_NC_USER_NORMALLOGOUT_CMD(client *network.Client, pkt *protocol.Packet)`
- L294: `func (h *Handler) handle_PROTO_NC_USER_FILE_CHECK_REQ(client *network.Client, pkt *protocol.Packet)`
- L300: `func (h *Handler) handle_PROTO_NC_USER_WORLDSTATUS_REQ(client *network.Client, pkt *protocol.Packet)`
- L309: `func (h *Handler) handle_PROTO_NC_USER_WORLDSELECT_REQ(client *network.Client, pkt *protocol.Packet)`
- L345: `func (h *Handler) send_PROTO_NC_USER_LOGIN_ACK(client *network.Client, isResend bool)`
- L367: `func (h *Handler) send_PROTO_NC_USER_LOGINFAIL_ACK(client *network.Client, errCode protocol.ServerError)`
- L372: `func md5Hash(text string) string`

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

## cmd/LoginService/main.go
- L20: `func main()`

