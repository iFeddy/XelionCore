package database

import (
	"fmt"
	"time"
)

// IssueZoneTransfer registra un traspaso de sesión de un solo uso en tActiveZoneTransfer
// invocando atómicamente el Stored Procedure oficial usp_Session_IssueZoneTransfer.
func (r *MySQLCharacterRepository) IssueZoneTransfer(
	worldNo uint8, zoneNo uint16, userNo int, userID string,
	charNo int, charID string, randomID uint16, adminLevel uint8,
	clientIP, targetMap string, targetX, targetY int,
) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Session_IssueZoneTransfer(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		worldNo, zoneNo, userNo, userID, charNo, charID, randomID, adminLevel, clientIP, targetMap, targetX, targetY,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Session_IssueZoneTransfer: %w", err)
	}
	if ret != 0 {
		return fmt.Errorf("usp_Session_IssueZoneTransfer failed with code %d", ret)
	}
	return nil
}

// ConsumeZoneTransfer valida, consume y destruye de forma atómica y transaccional un traspaso activo
// invocando usp_Session_ConsumeZoneTransfer. Garantiza protección estricta contra ataques de replay.
func (r *MySQLCharacterRepository) ConsumeZoneTransfer(charID string, randomID uint16, clientIP string) (*ZoneTransfer, error) {
	row := r.db.QueryRow("CALL usp_Session_ConsumeZoneTransfer(?, ?, ?)", charID, randomID, clientIP)

	var (
		worldNo, zoneNo, userNo, charNo, adminLevel, targetX, targetY, ret int
		userID, cID, targetMap                                             string
	)

	err := row.Scan(&worldNo, &zoneNo, &userNo, &userID, &charNo, &cID, &adminLevel, &targetMap, &targetX, &targetY, &ret)
	if err != nil {
		return nil, fmt.Errorf("error scanning usp_Session_ConsumeZoneTransfer: %w", err)
	}

	if ret != 0 {
		return nil, fmt.Errorf("transfer verification failed with code %d", ret)
	}

	return &ZoneTransfer{
		AccountID: userNo,
		Username:  userID,
		CharName:  cID,
		CharID:    charNo,
		RandomID:  randomID,
		Admin:     byte(adminLevel),
		HostIP:    clientIP,
		CreatedAt: time.Now(),
	}, nil
}

// UpsertZoneNode registra o actualiza un nodo ZoneService activo en tWorldZoneNode invocando usp_WorldZone_UpsertNode.
func (r *MySQLCharacterRepository) UpsertZoneNode(worldNo uint8, zoneNo uint16, zoneIP string, zonePort uint16) error {
	var ret int
	err := r.db.QueryRow("CALL usp_WorldZone_UpsertNode(?, ?, ?, ?)", worldNo, zoneNo, zoneIP, zonePort).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_WorldZone_UpsertNode: %w", err)
	}
	return nil
}

// HeartbeatZoneNode actualiza el latido de vida y carga de jugadores de un nodo invocando usp_WorldZone_Heartbeat.
func (r *MySQLCharacterRepository) HeartbeatZoneNode(worldNo uint8, zoneNo uint16, currentPlayers int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_WorldZone_Heartbeat(?, ?, ?)", worldNo, zoneNo, currentPlayers).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_WorldZone_Heartbeat: %w", err)
	}
	return nil
}

// GetZoneNodes recupera todos los nodos ZoneService activos de un mundo invocando usp_WorldZone_GetList.
func (r *MySQLCharacterRepository) GetZoneNodes(worldNo uint8) ([]*ZoneNodeInfo, error) {
	rows, err := r.db.Query("CALL usp_WorldZone_GetList(?)", worldNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_WorldZone_GetList(%d): %w", worldNo, err)
	}
	defer rows.Close()

	var nodes []*ZoneNodeInfo
	for rows.Next() {
		node := &ZoneNodeInfo{}
		var bOnline int
		if err := rows.Scan(&node.WorldNo, &node.ZoneNo, &node.ZoneIP, &node.ZonePort, &node.CurrentPlayers, &bOnline, &node.LastHeartbeat); err == nil {
			node.IsOnline = (bOnline != 0)
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

// ResolveMapEndpoint resuelve dinámicamente el IP y puerto de la zona asignada invocando usp_WorldZone_ResolveMapEndpoint.
func (r *MySQLCharacterRepository) ResolveMapEndpoint(worldNo uint8, zoneNo uint16) (string, uint16, bool, error) {
	var ip string
	var port uint16
	var bOnline int
	err := r.db.QueryRow("CALL usp_WorldZone_ResolveMapEndpoint(?, ?)", worldNo, zoneNo).Scan(&ip, &port, &bOnline)
	if err != nil {
		return "", 0, false, fmt.Errorf("error calling usp_WorldZone_ResolveMapEndpoint: %w", err)
	}
	return ip, port, (bOnline != 0), nil
}

// Aseguramiento estático de implementación completa de CharacterRepository sin omisiones.
var _ CharacterRepository = (*MySQLCharacterRepository)(nil)
