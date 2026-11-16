package zone

import "strings"

// Config holds runtime options for the Zone Server.
type Config struct {
	BindAddr     string `json:"bind_addr"`
	Port         int    `json:"port"`
	IP           string `json:"ip"`
	MapID        uint16 `json:"map_id"`
	WorldNo      byte   `json:"world_no"`
	ZoneNo       uint16 `json:"zone_no"`
	ZoneID       string `json:"zone_id"`
	MultiNode    bool   `json:"multi_node"`
	PreloadAll   bool   `json:"preload_all"`
	WorldS2SAddr string `json:"world_s2s_addr"`
}

// ZoneNodePorts contiene los puertos oficiales S2S/C2S por nodo Zone00..Zone04 (Zone.md / FASE 16).
var ZoneNodePorts = map[string]struct {
	ZoneNo uint16
	Port   int
}{
	"zone00": {ZoneNo: 0, Port: 9016},
	"zone01": {ZoneNo: 1, Port: 9019},
	"zone02": {ZoneNo: 2, Port: 9022},
	"zone03": {ZoneNo: 3, Port: 9025},
	"zone04": {ZoneNo: 4, Port: 9028},
}

// ResolveZoneNodeConfig resuelve el número de zona y puerto oficial para Zone00..Zone04.
func ResolveZoneNodeConfig(zoneID string) (zoneNo uint16, port int, ok bool) {
	norm := strings.ToLower(strings.TrimSpace(zoneID))
	if norm == "" {
		return 0, 0, false
	}
	entry, found := ZoneNodePorts[norm]
	if !found {
		return 0, 0, false
	}
	return entry.ZoneNo, entry.Port, true
}
