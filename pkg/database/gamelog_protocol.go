package database

import (
	"strings"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// EnqueueFromLogType0 encola un evento PROTO_NC_LOG_GAME_DATA_TYPE_0 en el despachador de logs.
func (d *GameLogDispatcher) EnqueueFromLogType0(p *protocol.PROTO_NC_LOG_GAME_DATA_TYPE_0_DATA) bool {
	if p == nil {
		return false
	}
	sMap := strings.TrimRight(string(p.SMap[:]), "\x00")
	return d.Enqueue(&GameLogEntry{
		Type:         int(p.NType),
		CharNo:       int(p.NCharNo),
		Map:          sMap,
		MapX:         int(p.NMapX),
		MapY:         int(p.NMapY),
		MapZ:         int(p.NMapZ),
		TargetCharNo: int(p.NTargetCharNo),
		TargetID:     int(p.NTargetID),
		ItemKey:      int64(p.NItemKey),
		Int1:         int(p.NInt1),
		Int2:         int(p.NInt2),
		Int3:         int(p.NInt3),
		Bigint1:      int64(p.NBigint1),
	})
}

// EnqueueFromLogType1 encola un evento PROTO_NC_LOG_GAME_DATA_TYPE_1 en el despachador de logs.
func (d *GameLogDispatcher) EnqueueFromLogType1(p *protocol.PROTO_NC_LOG_GAME_DATA_TYPE_1_DATA) bool {
	if p == nil {
		return false
	}
	sMap := strings.TrimRight(string(p.SMap[:]), "\x00")
	return d.EnqueueType1(int(p.NType), int(p.NCharNo), sMap, int(p.NMapX), int(p.NMapY), int(p.NMapZ))
}

// EnqueueFromLogType2 encola un evento PROTO_NC_LOG_GAME_DATA_TYPE_2 en el despachador de logs.
func (d *GameLogDispatcher) EnqueueFromLogType2(p *protocol.PROTO_NC_LOG_GAME_DATA_TYPE_2_DATA) bool {
	if p == nil {
		return false
	}
	sMap := strings.TrimRight(string(p.SMap[:]), "\x00")
	return d.EnqueueType2(int(p.NType), int(p.NCharNo), sMap, int(p.NMapX), int(p.NMapY), int(p.NMapZ), int(p.NInt1), int(p.NInt2))
}

// EnqueueFromLogType7 encola un evento PROTO_NC_LOG_GAME_DATA_TYPE_7 en el despachador de logs.
func (d *GameLogDispatcher) EnqueueFromLogType7(p *protocol.PROTO_NC_LOG_GAME_DATA_TYPE_7_DATA) bool {
	if p == nil {
		return false
	}
	sMap := strings.TrimRight(string(p.SMap[:]), "\x00")
	return d.EnqueueType7(int(p.NType), int(p.NCharNo), sMap, int(p.NMapX), int(p.NMapY), int(p.NMapZ), int(p.NTargetCharNo), int(p.NInt1), int(p.NInt2))
}

// EnqueueFromLogType8 encola un evento PROTO_NC_LOG_GAME_DATA_TYPE_8 en el despachador de logs.
func (d *GameLogDispatcher) EnqueueFromLogType8(p *protocol.PROTO_NC_LOG_GAME_DATA_TYPE_8_DATA) bool {
	if p == nil {
		return false
	}
	sMap := strings.TrimRight(string(p.SMap[:]), "\x00")
	return d.EnqueueType8(int(p.NType), int(p.NCharNo), sMap, int(p.NMapX), int(p.NMapY), int(p.NMapZ), int(p.NInt1))
}

