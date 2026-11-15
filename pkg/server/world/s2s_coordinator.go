package world

import (
	"fmt"
	"log"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/s2s"
)

// S2SSink define el canal vivo para publicar eventos wms_NC_* hacia nodos ZoneService conectados.
type S2SSink interface {
	Broadcast(ev s2s.Event)
	BroadcastExceptZone(ev s2s.Event, exceptZoneNo uint16)
	SendToZone(zoneNo uint16, ev s2s.Event) bool
}

// SetS2SSink inyecta el coordinador S2S vivo en el despachador de WorldService.
func (h *Handler) SetS2SSink(sink S2SSink) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.s2sSink = sink
}

// BroadcastToZones publica un evento oficial wms_NC_* hacia todos los nodos ZoneService conectados.
func (h *Handler) BroadcastToZones(symbol string, payload map[string]interface{}) {
	h.mu.Lock()
	sink := h.s2sSink
	worldID := h.worldID
	h.mu.Unlock()
	if sink == nil {
		return
	}
	sink.Broadcast(s2s.Event{
		Symbol:  symbol,
		WorldID: worldID,
		Payload: payload,
	})
}

// HandleZoneEvent procesa eventos wms_NC_* recibidos en vivo desde un proceso ZoneService.
func (h *Handler) HandleZoneEvent(zoneNo uint16, ev s2s.Event) {
	switch ev.Symbol {
	case "wms_NC_MISC_S2SCONNECTION_RDY":
		log.Printf("[World][S2S] ZoneService node #%d connected for World %d", zoneNo, h.worldID)

	case "wms_NC_USER_CONNECTCUT2ZONE_CMD":
		accID, ok := s2sInt(ev.Payload, "account_id")
		if !ok {
			accID, _ = s2sInt(ev.Payload, "userNo")
		}
		charID, ok := s2sInt(ev.Payload, "char_id")
		if !ok {
			charID, _ = s2sInt(ev.Payload, "charID")
		}
		logoutType, ok := s2sInt(ev.Payload, "logout_type")
		if !ok {
			logoutType, ok = s2sInt(ev.Payload, "logoutType")
			if !ok {
				logoutType, _ = s2sInt(ev.Payload, "type")
			}
		}
		h.handleZoneLogoutSync(accID, charID, byte(logoutType))

	case "wms_NC_CT_CHARTTING_CMD", "wms_NC_MISC_START_THE_BOOM_CMD":
		msg := s2sString(ev.Payload, "message")
		if msg == "" {
			return
		}
		typ := s2sString(ev.Payload, "type")
		sender := s2sString(ev.Payload, "sender")

		// Reenviar a otros nodos ZoneService distintos al origen para coherencia multi-zona
		h.mu.Lock()
		sink := h.s2sSink
		worldID := h.worldID
		h.mu.Unlock()
		if sink != nil {
			fwdPayload := make(map[string]interface{}, len(ev.Payload)+1)
			for k, v := range ev.Payload {
				fwdPayload[k] = v
			}
			fwdPayload["inbound"] = true
			sink.BroadcastExceptZone(s2s.Event{
				Symbol:  ev.Symbol,
				WorldID: worldID,
				ZoneNo:  ev.ZoneNo,
				Payload: fwdPayload,
			}, zoneNo)
		}

		if typ == "roar" || typ == "shout" {
			if len(msg) > 250 {
				msg = msg[:250]
			}
			cmd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONESHOUT_CMD)
			cmd.WriteByte(0)
			cmd.WriteString(sender, 20)
			cmd.WriteByte(1)
			cmd.WriteByte(byte(len(msg)))
			cmd.WriteBytes([]byte(msg))
			h.broadcastNonTransferring(cmd)
		}

	case "wms_NC_USER_CONNECTION_ZONE_CHANGE_CHAR_ID_CMD":
		charID, _ := s2sInt(ev.Payload, "char_id")
		newName := s2sString(ev.Payload, "new_name")
		if charID > 0 && newName != "" {
			h.syncCharRenameInWorld(charID, newName)
		}

	case "wms_NC_GUILD_WORLD_RENAME_ACK":
		guildID, _ := s2sInt(ev.Payload, "guild_id")
		newName := s2sString(ev.Payload, "new_name")
		if guildID > 0 && newName != "" && h.guildMgr != nil {
			h.guildMgr.SyncRenameCache(guildID, newName)
		}

	case "wms_NC_GUILD_WORLD_RETYPE_ACK":
		guildID, _ := s2sInt(ev.Payload, "guild_id")
		newType, _ := s2sInt(ev.Payload, "new_type")
		if guildID > 0 && h.guildMgr != nil {
			h.guildMgr.SyncRetypeCache(guildID, byte(newType))
		}

	case "wms_NC_PARTY_INFORM_REQ":
		if partyNo, ok := s2sInt(ev.Payload, "party_no"); ok && partyNo > 0 {
			h.BroadcastToZones("wms_NC_PARTY_INFORM_ACK", map[string]interface{}{
				"party_no": partyNo,
				"inbound":  true,
			})
		}
	}
}

func (h *Handler) handleZoneLogoutSync(accID, charID int, logoutType byte) {
	var charName string
	var target *network.Client
	h.mu.Lock()
	if accID > 0 {
		target = h.activeClients[accID]
	}
	if target == nil && charID > 0 {
		for _, c := range h.activeClients {
			if c != nil && c.CharacterID == charID {
				target = c
				break
			}
		}
	}
	if target != nil {
		charName = target.CharName
		if logoutType == 1 {
			// Retorno a Selección de Personaje: conservar socket en WorldService pero limpiar personaje activo
			target.CharacterID = 0
			target.CharName = ""
			target.IsTransferring = false
		} else {
			if target.AccountID > 0 && h.activeClients[target.AccountID] == target {
				delete(h.activeClients, target.AccountID)
			}
			target.CharacterID = 0
			target.CharName = ""
		}
	}
	h.mu.Unlock()

	if charID > 0 {
		if h.partyMgr != nil {
			h.partyMgr.LeavePartyByCharID(charID)
		}
		if h.guildMgr != nil {
			h.guildMgr.UnregisterOnline(charID)
		}
		if h.kqMgr != nil {
			h.kqMgr.UnregisterOnline(charID)
			h.kqMgr.RemovePlayer(charID)
		}
		if h.friendMgr != nil && charName != "" {
			h.friendMgr.BroadcastLogout(charID, charName)
		}
	}

	if target != nil && logoutType != 1 {
		target.Close()
	}
	log.Printf("[World][S2S] Zone logout sync processed (Account=%d Char=%d LogoutType=%d)", accID, charID, logoutType)
}

func (h *Handler) syncCharRenameInWorld(charID int, newName string) {
	h.mu.Lock()
	for _, c := range h.activeClients {
		if c != nil && c.CharacterID == charID {
			c.CharName = newName
		}
	}
	h.mu.Unlock()

	if h.partyMgr != nil {
		h.partyMgr.RenameOnlineChar(charID, newName)
	}
}

func (h *Handler) broadcastNonTransferring(pkt *protocol.PacketWriter) {
	h.mu.Lock()
	clients := make([]*network.Client, 0, len(h.activeClients))
	for _, c := range h.activeClients {
		if c != nil && !c.IsTransferring && c.CharacterID > 0 {
			clients = append(clients, c)
		}
	}
	h.mu.Unlock()

	for _, c := range clients {
		_ = c.SendPacket(pkt)
	}
}

func s2sString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func s2sInt(m map[string]interface{}, key string) (int, bool) {
	if m == nil {
		return 0, false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case uint16:
		return int(n), true
	case uint32:
		return int(n), true
	case byte:
		return int(n), true
	case float64:
		return int(n), true
	default:
		if s, ok := v.(string); ok {
			s = strings.TrimSpace(s)
			var val int
			if _, err := fmt.Sscanf(s, "%d", &val); err == nil {
				return val, true
			}
		}
		return 0, false
	}
}
