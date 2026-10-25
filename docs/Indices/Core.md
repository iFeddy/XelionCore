# Índice: config/crypto/logger/network

> Generado automáticamente por scripts/gen_indices. Solo firmas con n° de línea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. 

## pkg/config/config.go
- L12: `type WorldConfig struct`
- L22: `type Config struct`
- L41: `func Load(envPath ...string) (*Config, error)`

## pkg/config/config_test.go
- L8: `func TestLoadConfig(t *testing.T)`

## pkg/crypto/net_crypto.go
- L47: `type NetCrypto struct`
- L52: `func NewNetCrypto(offset int16) (*NetCrypto, error)`
- L60: `func RandomXorOffset() int16`
- L69: `func (nc *NetCrypto) XorPos() int16`
- L74: `func (nc *NetCrypto) Crypt(buf []byte)`

## pkg/crypto/net_crypto_test.go
- L8: `func TestNetCryptoSymmetry(t *testing.T)`
- L37: `func TestNetCryptoWrapAround(t *testing.T)`

## pkg/logger/logger.go
- L15: `func InitSessionLogger(serviceName string) (*os.File, error)`

## pkg/network/client.go
- L18: `type Client struct`
- L41: `func NewClient(conn net.Conn, onClose func(*Client)) (*Client, error)`
- L59: `func (c *Client) RemoteAddr() string`
- L64: `func (c *Client) RemoteIP() string`
- L73: `func (c *Client) XorOffset() int16`
- L81: `func (c *Client) Crypto() *crypto.NetCrypto`
- L86: `func (c *Client) SendPacket(pw *protocol.PacketWriter) error`
- L107: `func (c *Client) Close()`
- L117: `func (c *Client) StartReadLoop(packetHandler func(*Client, *protocol.Packet))`

## pkg/network/server.go
- L13: `type TCPServer struct`
- L25: `func NewTCPServer(addr string, onConnect func(*Client), packetHandler func(*Client, *protocol.Packet)) *TCPServer`
- L35: `func (s *TCPServer) SetOnDisconnect(fn func(*Client))`
- L40: `func (s *TCPServer) Start() error`
- L54: `func (s *TCPServer) acceptLoop()`
- L103: `func (s *TCPServer) Stop()`
