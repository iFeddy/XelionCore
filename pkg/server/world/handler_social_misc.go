package world

import (
	"log"
	"strings"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ============================================================================
// DEPARTAMENTO 2: MISC, SEGURIDAD, URLS, APEX, SPAM
// ============================================================================

func (h *Handler) handleClientDebug(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_CLIENT_DEBUG_MSG_CMD)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleHeartbeat(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_HEARTBEAT_ACK)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handlePing(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB)
	resp.WriteUint32(uint32(time.Now().Unix()))
	_ = client.SendPacket(resp)
}

func (h *Handler) handleWebKey(client *network.Client, pkt *protocol.Packet) {
	var key string
	if h.accountRepo != nil {
		k, err := h.accountRepo.GetOrIssueWebKey(client.AccountID, "")
		if err == nil {
			key = k
		}
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_WEB_KEY_ACK)
	resp.WriteString(key, 64)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleItemShopURL(client *network.Client, pkt *protocol.Packet) {
	shopURL := "https://xelion.online/itemshop"
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_ITEMSHOP_URL_ACK)
	resp.WriteByte(byte(len(shopURL)))
	resp.WriteBytes([]byte(shopURL))
	_ = client.SendPacket(resp)
}

func (h *Handler) handleApexServerData(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[World][%s] Received Apex anticheat data len=%d", client.RemoteAddr(), pkt.Reader.Remaining())
}

func (h *Handler) handleCSReq(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_CS_ACK)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleCSAck(client *network.Client, pkt *protocol.Packet) {
	// CS Ack processed
}

func (h *Handler) handleUseCenBilling(client *network.Client, pkt *protocol.Packet) {
	// Use cen billing processed
}

func (h *Handler) handleXtrapData(client *network.Client, pkt *protocol.Packet) {
	// X-Trap data telemetry processed
}

func (h *Handler) handleHackScan(client *network.Client, pkt *protocol.Packet) {
	hackType, _ := pkt.Reader.ReadInt16()
	desc, _ := pkt.Reader.ReadString(pkt.Reader.Remaining())
	if h.logRepo != nil {
		_ = h.logRepo.AddHackScanLogFull(client.AccountID, h.worldID, client.CharacterID, hackType, client.RemoteIP(), client.Username, desc)
	}
	log.Printf("[World][%s] HackScan logged: type=%d, desc=%s", client.RemoteAddr(), hackType, desc)
}

func (h *Handler) handleClientLoadingBug(client *network.Client, pkt *protocol.Packet) {
	mapName, _ := pkt.Reader.ReadString(16)
	failType, _ := pkt.Reader.ReadInt16()
	failCnt, _ := pkt.Reader.ReadInt16()
	if h.charRepo != nil {
		_ = h.charRepo.SetClientLoadingBug(client.CharacterID, strings.TrimRight(mapName, "\x00"), failType, failCnt)
	}
	log.Printf("[World][%s] Client loading bug detected on map %s (failType=%d, failCnt=%d)", client.RemoteAddr(), mapName, failType, failCnt)
}

func (h *Handler) handleSpammerReport(client *network.Client, pkt *protocol.Packet) {
	targetCharID, _ := pkt.Reader.ReadUint32()
	blockTime, _ := pkt.Reader.ReadInt32()
	standbyTime, _ := pkt.Reader.ReadInt32()
	if h.charRepo != nil {
		_ = h.charRepo.SetSpamBlockTime(int(targetCharID), int(blockTime), int(standbyTime))
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_SPAMMER_REPORT_ACK)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleSpamFilterSet(client *network.Client, pkt *protocol.Packet) {
	spread, _ := pkt.Reader.ReadByte()
	block, _ := pkt.Reader.ReadByte()
	rate, _ := pkt.Reader.ReadByte()
	if h.charRepo != nil {
		_ = h.charRepo.SetSpamFilter(client.AccountID, spread, block, rate)
	}
}

func (h *Handler) handleChatBan(client *network.Client, pkt *protocol.Packet) {
	targetCharID, _ := pkt.Reader.ReadUint32()
	duration, _ := pkt.Reader.ReadInt32()
	if h.charRepo != nil {
		_ = h.charRepo.SetSpamBlockTime(int(targetCharID), int(duration), 0)
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_SPAMMER_CHAT_BAN_ACK)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

// ============================================================================
// DEPARTAMENTO 8: CHAT, WHISPERS, NOTICES
// ============================================================================

func (h *Handler) handleWhisper(client *network.Client, pkt *protocol.Packet) {
	_, _ = pkt.Reader.ReadByte() // itemLinkDataCount
	targetNameBytes, _ := pkt.Reader.ReadBytes(20)
	targetName := strings.TrimRight(string(targetNameBytes), "\x00")
	msgLen, _ := pkt.Reader.ReadByte()
	msgBytes, _ := pkt.Reader.ReadBytes(int(msgLen))
	msg := string(msgBytes)

	targetKey := strings.ToLower(strings.TrimSpace(targetName))
	h.mu.Lock()
	var targetClient *network.Client
	for _, c := range h.activeClients {
		if strings.ToLower(c.Username) == targetKey {
			targetClient = c
			break
		}
	}
	h.mu.Unlock()

	if targetClient == nil || targetClient.CharacterID == 0 {
		failResp := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WHISPERFAIL_ACK)
		failResp.WriteUint16(1) // Not online
		failResp.WriteString(targetName, 20)
		_ = client.SendPacket(failResp)
		return
	}

	// Sender success ACK
	succResp := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WHISPERSUCCESS_ACK)
	succResp.WriteByte(0) // itemLinkDataCount
	succResp.WriteString(targetName, 20)
	succResp.WriteByte(byte(len(msg)))
	succResp.WriteBytes([]byte(msg))
	_ = client.SendPacket(succResp)

	// Delivery to receiver: PROTO_NC_ACT_WHISPER_CMD (OpCode 0x200D)
	cmd := protocol.NewPacketWriter((8 << 10) | 13)
	cmd.WriteByte(0) // itemLinkDataCount
	cmd.WriteString(client.Username, 20)
	cmd.WriteByte(byte(len(msg)))
	cmd.WriteBytes([]byte(msg))
	_ = targetClient.SendPacket(cmd)
}

func (h *Handler) handleNotice(client *network.Client, pkt *protocol.Packet) {
	if client.Admin == 0 {
		return
	}
	msgLen, _ := pkt.Reader.ReadByte()
	msgBytes, _ := pkt.Reader.ReadBytes(int(msgLen))
	noticeMsg := string(msgBytes)

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NOTICE_CMD)
	cmd.WriteByte(byte(len(noticeMsg)))
	cmd.WriteBytes([]byte(noticeMsg))

	h.mu.Lock()
	defer h.mu.Unlock()
	for _, c := range h.activeClients {
		_ = c.SendPacket(cmd)
	}
}

func (h *Handler) handlePartyChat(client *network.Client, pkt *protocol.Packet) {
	_, _ = pkt.Reader.ReadByte() // itemLinkDataCount
	msgLen, _ := pkt.Reader.ReadByte()
	msgBytes, _ := pkt.Reader.ReadBytes(int(msgLen))
	msg := string(msgBytes)

	h.partyMgr.mu.RLock()
	defer h.partyMgr.mu.RUnlock()
	pID, inParty := h.partyMgr.charToParty[client.CharacterID]
	if !inParty {
		return
	}
	party, ok := h.partyMgr.parties[pID]
	if !ok {
		return
	}

	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PARTYCHAT_CMD)
	cmd.WriteString(client.Username, 20)
	cmd.WriteUint16(uint16(len(msg)))
	cmd.WriteBytes([]byte(msg))

	for _, m := range party.Members {
		if m.Client != nil {
			_ = m.Client.SendPacket(cmd)
		}
	}
}

// ============================================================================
// DEPARTAMENTO 31: PRISON (tPrison)
// ============================================================================

func (h *Handler) handlePrisonSendIn(client *network.Client, pkt *protocol.Packet) {
	targetName, _ := pkt.Reader.ReadString(20)
	minutes, _ := pkt.Reader.ReadInt16()
	reason, _ := pkt.Reader.ReadString(pkt.Reader.Remaining())

	if h.charRepo != nil {
		_ = h.charRepo.AddPrison(client.Username, targetName, minutes, reason, "")
	}

	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_ADD_GM_ACK)
	resp.WriteString(targetName, 20)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handlePrisonTakeOut(client *network.Client, pkt *protocol.Packet) {
	targetCharID, _ := pkt.Reader.ReadUint32()
	if h.charRepo != nil {
		_ = h.charRepo.EndPrison(int(targetCharID))
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_END_ACK)
	resp.WriteUint32(targetCharID)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handlePrisonAlterTime(client *network.Client, pkt *protocol.Packet) {
	targetName, _ := pkt.Reader.ReadString(20)
	minutes, _ := pkt.Reader.ReadInt16()
	reason, _ := pkt.Reader.ReadString(pkt.Reader.Remaining())

	if h.charRepo != nil {
		_ = h.charRepo.AlterPrison(client.Username, targetName, minutes, reason, "")
	}

	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_ALTER_GM_ACK)
	resp.WriteString(targetName, 20)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handlePrisonConfirm(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_GET_ACK)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

// ============================================================================
// DEPARTAMENTO 32: REPORT (tReport en AccountLog)
// ============================================================================

func (h *Handler) handleReportAdd(client *network.Client, pkt *protocol.Packet) {
	targetName, _ := pkt.Reader.ReadString(20)
	reason, _ := pkt.Reader.ReadString(64)
	chatLog, _ := pkt.Reader.ReadString(pkt.Reader.Remaining())

	var regNo int
	if h.logRepo != nil {
		r := &database.Report{
			Reporter:       client.AccountID,
			ReporterCharID: client.Username,
			Harmer:         targetName,
			WorldNo:        h.worldID,
			ReportType:     "IN_GAME",
			Reason:         reason,
			ChatLog:        chatLog,
		}
		regNo, _ = h.logRepo.AddReport(r)
	}

	resp := protocol.NewPacketWriter(protocol.PROTO_NC_REPORT_ADD_ACK)
	resp.WriteUint32(uint32(regNo))
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleReportCancel(client *network.Client, pkt *protocol.Packet) {
	regNo, _ := pkt.Reader.ReadUint32()
	var retCode int16
	if h.logRepo != nil {
		retCode, _ = h.logRepo.DeleteReport(int(regNo))
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_REPORT_CANCEL_ACK)
	resp.WriteUint32(regNo)
	resp.WriteInt16(retCode)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleReportGet(client *network.Client, pkt *protocol.Packet) {
	var reports []*database.ReportSummary
	if h.logRepo != nil {
		reports, _ = h.logRepo.GetReports(client.AccountID, h.worldID)
	}
	if len(reports) == 0 {
		resp := protocol.NewPacketWriter(protocol.PROTO_NC_REPORT_GET_ACK)
		resp.WriteUint16(0) // err
		resp.WriteUint32(0) // nRegNumber
		resp.WriteByte(0)   // byReportResult
		var emptyRemark [128]byte
		resp.WriteBytes(emptyRemark[:])
		_ = client.SendPacket(resp)
		return
	}
	for _, r := range reports {
		resp := protocol.NewPacketWriter(protocol.PROTO_NC_REPORT_GET_ACK)
		resp.WriteUint16(0) // err
		resp.WriteUint32(uint32(r.RegNumber))
		resp.WriteByte(r.Result)
		var remarkBuf [128]byte
		copy(remarkBuf[:], r.Remark)
		resp.WriteBytes(remarkBuf[:])
		_ = client.SendPacket(resp)
	}
}

// ============================================================================
// DEPARTAMENTO 37: HOLY PROMISE (tHolyPromise)
// ============================================================================

func (h *Handler) handleHolyPromiseSetUp(client *network.Client, pkt *protocol.Packet) {
	charID, _ := pkt.Reader.ReadString(20)
	upCharID, _ := pkt.Reader.ReadString(20)

	targetName := upCharID
	if strings.EqualFold(client.Username, upCharID) {
		targetName = charID
	}

	var targetClient *network.Client
	var targetID int
	targetKey := strings.ToLower(strings.TrimSpace(targetName))
	h.mu.Lock()
	for _, c := range h.activeClients {
		if strings.ToLower(c.Username) == targetKey {
			targetClient = c
			targetID = c.CharacterID
			break
		}
	}
	h.mu.Unlock()

	h.holyPromiseMgr.HandleSetUpReq(client, charID, upCharID, targetClient, targetID)
}

func (h *Handler) handleHolyPromiseSetUpConfirm(client *network.Client, pkt *protocol.Packet) {
	charID, _ := pkt.Reader.ReadString(20)
	upCharID, _ := pkt.Reader.ReadString(20)
	accept, _ := pkt.Reader.ReadByte()

	requesterName := charID
	if strings.EqualFold(client.Username, charID) {
		requesterName = upCharID
	}

	var reqClient *network.Client
	targetKey := strings.ToLower(strings.TrimSpace(requesterName))
	h.mu.Lock()
	for _, c := range h.activeClients {
		if strings.ToLower(c.Username) == targetKey {
			reqClient = c
			break
		}
	}
	h.mu.Unlock()

	h.holyPromiseMgr.HandleSetUpConfirmAck(client, charID, upCharID, accept, reqClient)
}

func (h *Handler) handleHolyPromiseDelUp(client *network.Client, pkt *protocol.Packet) {
	var masterClient *network.Client
	data, _ := h.charRepo.GetHolyPromise(client.CharacterID)
	if data != nil && data.UpCharNo > 0 {
		h.mu.Lock()
		for _, c := range h.activeClients {
			if c.CharacterID == data.UpCharNo {
				masterClient = c
				break
			}
		}
		h.mu.Unlock()
	}
	h.holyPromiseMgr.HandleDelUpReq(client, masterClient)
}

func (h *Handler) handleHolyPromiseDelDown(client *network.Client, pkt *protocol.Packet) {
	downCharName, _ := pkt.Reader.ReadString(20)
	var pupilClient *network.Client
	var pupilID int
	targetKey := strings.ToLower(strings.TrimSpace(downCharName))
	h.mu.Lock()
	for _, c := range h.activeClients {
		if strings.ToLower(c.Username) == targetKey {
			pupilClient = c
			pupilID = c.CharacterID
			break
		}
	}
	h.mu.Unlock()
	if pupilID == 0 {
		ch, err := h.charRepo.GetCharacterByName(downCharName)
		if err == nil && ch != nil {
			pupilID = ch.ID
		}
	}
	h.holyPromiseMgr.HandleDelDownReq(client, downCharName, pupilID, pupilClient)
}

// ============================================================================
// DEPARTAMENTO 41: BATTLE RANK / INSTANCE DUNGEONS
// ============================================================================

func (h *Handler) handleBattleRankClearTime(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGEDEALT_ACK)
	resp.WriteUint16(0)
	resp.WriteUint16(0) // Count
	_ = client.SendPacket(resp)
}

func (h *Handler) handleBattleRankScore(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGETAKEN_ACK)
	resp.WriteUint16(0)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleBattleRankPoint(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_HEALINGDONE_ACK)
	resp.WriteUint16(0)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleBattleRankMyScore(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_MYRANK_ACK)
	resp.WriteUint16(0)
	resp.WriteUint32(0)
	_ = client.SendPacket(resp)
}

// ============================================================================
// DEPARTAMENTO 42: CHAT RESTRICT (tChatBlock)
// ============================================================================

func (h *Handler) handleChatRestrictAdd(client *network.Client, pkt *protocol.Packet) {
	targetCharID, _ := pkt.Reader.ReadUint32()
	if h.charRepo != nil {
		_ = h.charRepo.AddChatRestrict(client.AccountID, int(targetCharID))
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_CHAT_RESTRICT_ADD_ACK)
	resp.WriteUint32(targetCharID)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleChatRestrictDel(client *network.Client, pkt *protocol.Packet) {
	targetCharID, _ := pkt.Reader.ReadUint32()
	if h.charRepo != nil {
		_ = h.charRepo.DelChatRestrict(client.AccountID, int(targetCharID))
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_CHAT_RESTRICT_DEL_ACK)
	resp.WriteUint32(targetCharID)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleChatRestrictDelAll(client *network.Client, pkt *protocol.Packet) {
	if h.charRepo != nil {
		_ = h.charRepo.DelAllChatRestrict(client.AccountID)
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_CHAT_RESTRICT_DEL_ALL_ACK)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}
