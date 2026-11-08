package world

import (
	"log"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// PendingHolyPromiseRequest representa una solicitud pendiente de promesa sagrada.
type PendingHolyPromiseRequest struct {
	RequesterID   int
	RequesterName string
	TargetID      int
	TargetName    string
	IsRequesterUp bool
}

// HolyPromiseManager gestiona los vínculos de Maestro-Alumno y Promesa Sagrada (Dept 37).
type HolyPromiseManager struct {
	mu             sync.RWMutex
	charRepo       database.CharacterRepository
	pendingInvites map[int]*PendingHolyPromiseRequest // TargetID -> Request
}

// NewHolyPromiseManager inicializa el administrador de promesas sagradas con persistencia en Base de Datos.
func NewHolyPromiseManager(charRepo database.CharacterRepository) *HolyPromiseManager {
	return &HolyPromiseManager{
		charRepo:       charRepo,
		pendingInvites: make(map[int]*PendingHolyPromiseRequest),
	}
}

// HandleSetUpReq procesa la solicitud de promesa sagrada (0x9401 PROTO_NC_HOLY_PROMISE_SET_UP_REQ).
func (hpm *HolyPromiseManager) HandleSetUpReq(client *network.Client, charID, upCharID string, targetClient *network.Client, targetID int) {
	hpm.mu.Lock()
	defer hpm.mu.Unlock()

	reqName := client.Username
	var isUp bool
	if strings.EqualFold(reqName, upCharID) {
		isUp = true
	} else if strings.EqualFold(reqName, charID) {
		isUp = false
	} else {
		_ = client.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_ACK(1))
		return
	}

	if targetClient == nil || targetID == 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_ACK(2)) // Target not online
		return
	}

	hpm.pendingInvites[targetID] = &PendingHolyPromiseRequest{
		RequesterID:   client.CharacterID,
		RequesterName: reqName,
		TargetID:      targetID,
		TargetName:    targetClient.Username,
		IsRequesterUp: isUp,
	}

	_ = targetClient.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_CONFIRM_REQ(charID, upCharID))
	log.Printf("[HolyPromise] '%s' requested holy promise with '%s'", reqName, targetClient.Username)
}

// HandleSetUpConfirmAck procesa la respuesta del invitado (0x9405 PROTO_NC_HOLY_PROMISE_SET_UP_CONFIRM_ACK).
func (hpm *HolyPromiseManager) HandleSetUpConfirmAck(client *network.Client, charID, upCharID string, accept byte, requesterClient *network.Client) {
	hpm.mu.Lock()
	defer hpm.mu.Unlock()

	targetID := client.CharacterID
	req, exists := hpm.pendingInvites[targetID]
	if !exists {
		_ = client.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_ACK(3))
		return
	}
	delete(hpm.pendingInvites, targetID)

	if accept == 0 {
		_ = client.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_ACK(4)) // Rejected
		if requesterClient != nil {
			_ = requesterClient.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_ACK(4))
		}
		return
	}

	var downCharNo, upCharNo int
	if req.IsRequesterUp {
		upCharNo = req.RequesterID
		downCharNo = req.TargetID
	} else {
		upCharNo = req.TargetID
		downCharNo = req.RequesterID
	}

	if hpm.charRepo != nil {
		if err := hpm.charRepo.SetHolyPromise(downCharNo, upCharNo, 0); err != nil {
			log.Printf("[HolyPromise] Error persisting holy promise (%d, %d): %v", downCharNo, upCharNo, err)
			_ = client.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_ACK(5))
			return
		}
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_ACK(0))
	if requesterClient != nil {
		_ = requesterClient.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_SET_UP_ACK(0))
	}
	log.Printf("[HolyPromise] Relationship established between Master %d and Pupil %d", upCharNo, downCharNo)
}

// HandleDelUpReq procesa la disolución solicitada por el alumno (0x9406 PROTO_NC_HOLY_PROMISE_DEL_UP_REQ).
func (hpm *HolyPromiseManager) HandleDelUpReq(client *network.Client, masterClient *network.Client) {
	charNo := client.CharacterID
	if hpm.charRepo != nil {
		data, err := hpm.charRepo.GetHolyPromise(charNo)
		if err == nil && data != nil && data.UpCharNo > 0 {
			_ = hpm.charRepo.DelHolyPromise(charNo, data.UpCharNo)
		} else {
			_ = hpm.charRepo.DelHolyPromiseChar(charNo)
		}
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_DEL_UP_ACK(0, 0))
	if masterClient != nil {
		_ = masterClient.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_DEL_CMD(client.Username, 0))
	}
	log.Printf("[HolyPromise] Pupil '%s' severed relationship with Master", client.Username)
}

// HandleDelDownReq procesa la disolución solicitada por el maestro (0x940A PROTO_NC_HOLY_PROMISE_DEL_DOWN_REQ).
func (hpm *HolyPromiseManager) HandleDelDownReq(client *network.Client, downCharName string, pupilID int, pupilClient *network.Client) {
	upCharNo := client.CharacterID
	if hpm.charRepo != nil && pupilID > 0 {
		_ = hpm.charRepo.DelHolyPromise(pupilID, upCharNo)
	}

	_ = client.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_DEL_DOWN_ACK(downCharName, 0))
	if pupilClient != nil {
		_ = pupilClient.SendPacket(protocol.New_PROTO_NC_HOLY_PROMISE_DEL_CMD(client.Username, 1))
	}
	log.Printf("[HolyPromise] Master '%s' severed relationship with Pupil '%s'", client.Username, downCharName)
}
