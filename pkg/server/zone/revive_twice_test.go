package zone

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// TestCharReviveReq_DieAndReviveTwice simula morir y revivir dos veces seguidas (en mapa exterior y en ciudad)
// verificando que PROTO_NC_CHAR_REVIVE_REQ (0x104E) en el mismo mapa envíe PROTO_NC_CHAR_REVIVESAME_CMD (0x104F)
// con la columna ID de MapInfo.shn (p. ej. RouN -> 150, Apoline -> 181, LegelVal -> 254) sin disparar sendInGameChunk,
// y que WarpPlayer envíe ese mismo ID en PROTO_NC_MAP_LINKSAME_CMD (0x1809).
func TestCharReviveReq_DieAndReviveTwice(t *testing.T) {
	ds, h, repo, _ := loadFase13TestEnv(t)
	if h.fieldStore == nil {
		t.Skip("FieldStore no disponible")
	}
	mapInfo := ds.GetSHN("MapInfo")
	if mapInfo == nil {
		t.Fatal("MapInfo.shn no cargado en DataStore")
	}

	lookupMapInfoID := func(mapName string) (uint16, bool) {
		for i := range mapInfo.Rows {
			r := &mapInfo.Rows[i]
			if strings.EqualFold(strings.TrimSpace(r.GetString("MapName")), mapName) {
				return uint16(r.GetInt("ID")), true
			}
		}
		return 0, false
	}

	cases := []struct {
		startMap string
	}{
		{"RouCos01"}, // mapa exterior (Respawn.City = RouCos01, ID=1)
		{"Rou"},      // ciudad (Respawn.City = Rou, ID=0)
		{"RouN"},     // ciudad nueva (Respawn.City = RouN, Serial=148, ID=150)
	}

	for _, tc := range cases {
		startField, ok := h.fieldStore.GetField(tc.startMap)
		if !ok || startField == nil {
			t.Fatalf("FieldList sin %s", tc.startMap)
		}
		wantMap := startField.Respawn.City
		wantField, ok := h.fieldStore.GetField(wantMap)
		if !ok || wantField == nil {
			t.Fatalf("Respawn.City %q de %s no existe en FieldList", wantMap, tc.startMap)
		}
		wantMapID, ok := lookupMapInfoID(wantField.MapID)
		if !ok {
			t.Fatalf("MapInfo.shn sin %s", wantField.MapID)
		}
		wantX, wantY := startField.Respawn.Spots[0].X, startField.Respawn.Spots[0].Y

		zMap := h.mapManager.GetOrCreateMap(tc.startMap)
		p, conn := newTestClientAndPlayer(t, repo, ds, zMap, h, 500+len(tc.startMap), "Reviver"+tc.startMap, 20, 1000, 1000)

		for i := 1; i <= 2; i++ {
			// Muerte
			p.Char.CurHP = 0
			p.State = 3
			conn.clear()

			h.HandlePacket(p.Client, protocol.NewPacket(protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_REVIVE_REQ).Data()))

			if p.State != 1 || p.Char.CurHP <= 0 {
				t.Fatalf("[%s #%d] revive: State=%d HP=%d", tc.startMap, i, p.State, p.Char.CurHP)
			}
			if p.CurrentMap == nil || p.CurrentMap.MapID != wantField.MapID {
				t.Fatalf("[%s #%d] mapa destino esperado %s, obtenido %+v", tc.startMap, i, wantField.MapID, p.CurrentMap)
			}
			if p.CurrentMap.ID != uint16(wantField.Serial) {
				t.Fatalf("[%s #%d] Serial esperado %d, obtenido %d", tc.startMap, i, wantField.Serial, p.CurrentMap.ID)
			}
			if p.Char.XPos != wantX || p.Char.YPos != wantY {
				t.Fatalf("[%s #%d] coords esperadas (%d,%d), obtenidas (%d,%d)", tc.startMap, i, wantX, wantY, p.Char.XPos, p.Char.YPos)
			}
			revives := conn.findByOpcode(protocol.PROTO_NC_CHAR_REVIVESAME_CMD)
			if len(revives) != 1 || len(revives[0].Payload) < 10 {
				t.Fatalf("[%s #%d] se esperaba 1 PROTO_NC_CHAR_REVIVESAME_CMD, obtenidos %d", tc.startMap, i, len(revives))
			}
			if got := binary.LittleEndian.Uint16(revives[0].Payload[0:2]); got != wantMapID {
				t.Fatalf("[%s #%d] REVIVESAME mapid esperado %d (col ID MapInfo.shn), obtenido %d", tc.startMap, i, wantMapID, got)
			}
			if n := len(conn.findByOpcode(protocol.PROTO_NC_MAP_LINKSAME_CMD)); n != 0 {
				t.Fatalf("[%s #%d] revive en mismo mapa no debe enviar LINKSAME (enviados=%d)", tc.startMap, i, n)
			}
			if n := len(conn.findByOpcode(protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD)); n != 0 {
				t.Fatalf("[%s #%d] revive en mismo mapa no debe disparar sendInGameChunk (CLIENT_BASE=%d)", tc.startMap, i, n)
			}

			// REVIVE_REQ duplicado estando vivo: debe ignorarse (idempotente).
			conn.clear()
			h.HandlePacket(p.Client, protocol.NewPacket(protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_REVIVE_REQ).Data()))
			if n := len(conn.findByOpcode(protocol.PROTO_NC_CHAR_REVIVESAME_CMD)); n != 0 {
				t.Fatalf("[%s #%d] REVIVE_REQ estando vivo no debe revivir (REVIVESAME=%d)", tc.startMap, i, n)
			}
		}

		// Verificar también que WarpPlayer hacia RouN y Apoline envía la columna ID de MapInfo.shn en LINKSAME_CMD.
		for _, target := range []string{"RouN", "Apoline"} {
			if strings.EqualFold(tc.startMap, target) {
				continue
			}
			wantID, ok := lookupMapInfoID(target)
			if !ok {
				t.Fatalf("MapInfo.shn sin %s", target)
			}
			conn.clear()
			if !h.WarpPlayer(p, target, 4134, 6488, 0) {
				t.Fatalf("[%s] WarpPlayer hacia %s falló", tc.startMap, target)
			}
			links := conn.findByOpcode(protocol.PROTO_NC_MAP_LINKSAME_CMD)
			if len(links) != 1 || len(links[0].Payload) < 10 {
				t.Fatalf("[%s] se esperaba 1 PROTO_NC_MAP_LINKSAME_CMD al warpear a %s, obtenidos %d", tc.startMap, target, len(links))
			}
			if got := binary.LittleEndian.Uint16(links[0].Payload[0:2]); got != wantID {
				t.Fatalf("[%s] LINKSAME mapid para %s esperado %d (col ID MapInfo.shn), obtenido %d", tc.startMap, target, wantID, got)
			}
			if n := len(conn.findByOpcode(protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD)); n != 0 {
				t.Fatalf("[%s] LINKSAME hacia %s no debe disparar sendInGameChunk (CLIENT_BASE=%d)", tc.startMap, target, n)
			}
			if p.MapObjectID < 10000 {
				t.Fatalf("[%s] player MapObjectID=%d debe estar en rango dedicado (>=10000) para no colisionar con NPCs/mobs", tc.startMap, p.MapObjectID)
			}
		}
	}
}
