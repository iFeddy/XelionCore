package zone

// Config holds runtime options for the Zone Server.
type Config struct {
	BindAddr string `json:"bind_addr"`
	Port     int    `json:"port"`
	IP       string `json:"ip"`
	MapID    uint16 `json:"map_id"`
}
