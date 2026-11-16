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
	msgBytes, _ := pkt.Reader.ReadBytes(pkt.Reader.Remaining())
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_CLIENT_DEBUG_MSG_CMD)
	resp.WriteUint16(uint16(client.CharacterID)) // wldmanhandle
	resp.WriteUint16(uint16(len(msgBytes)))    // msg_size
	if len(msgBytes) > 0 {
		resp.WriteBytes(msgBytes)
	}
	_ = client.SendPacket(resp)
}

func (h *Handler) handleHeartbeat(client *network.Client, pkt *protocol.Packet) {
	_ = client.SendPacket(protocol.New_PROTO_NC_MISC_HEARTBEAT_ACK())
}

func (h *Handler) handlePing(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB)
	resp.Fill(130, 0)
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
	resp.WriteUint16(0) // Error (2B)
	var keyBuf [32]byte
	copy(keyBuf[:], key)
	resp.WriteBytes(keyBuf[:]) // WebKey (32B)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleItemShopURL(client *network.Client, pkt *protocol.Packet) {
	shopURL := "https://xelion.online/itemshop"
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_ITEMSHOP_URL_ACK)
	resp.WriteUint16(0)                   // nError (2B)
	resp.WriteUint16(uint16(len(shopURL))) // nLen (2B)
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

func (h *Handler) findOnlineClientByCharName(charName string) *network.Client {
	targetKey := strings.ToLower(strings.TrimSpace(strings.TrimRight(charName, "\x00")))
	if targetKey == "" {
		return nil
	}
	h.mu.Lock()
	var found *network.Client
	for _, c := range h.activeClients {
		if strings.ToLower(strings.TrimSpace(c.CharName)) == targetKey {
			found = c
			break
		}
	}
	h.mu.Unlock()
	if found != nil {
		return found
	}
	if h.partyMgr != nil {
		h.partyMgr.mu.RLock()
		if mem, ok := h.partyMgr.onlineChars[targetKey]; ok && mem.Client != nil {
			found = mem.Client
		}
		h.partyMgr.mu.RUnlock()
	}
	return found
}

func (h *Handler) handleChatBan(client *network.Client, pkt *protocol.Packet) {
	var spammerName string
	if pkt.Reader.Remaining() >= 20 {
		spammerName, _ = pkt.Reader.ReadString(20)
		if h.charRepo != nil {
			if ch, err := h.charRepo.GetCharacterByName(spammerName); err == nil && ch != nil {
				_ = h.charRepo.SetSpamBlockTime(ch.ID, 600, 0)
			}
		}
	} else {
		targetCharID, _ := pkt.Reader.ReadUint32()
		duration, _ := pkt.Reader.ReadInt32()
		if h.charRepo != nil {
			_ = h.charRepo.SetSpamBlockTime(int(targetCharID), int(duration), 0)
		}
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_SPAMMER_CHAT_BAN_ACK)
	resp.WriteString(spammerName, 20)
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)
}

// ============================================================================
// DEPARTAMENTO 8: CHAT, WHISPERS, NOTICES
// ============================================================================

func (h *Handler) handleWhisper(client *network.Client, pkt *protocol.Packet) {
	itemLinkCount, _ := pkt.Reader.ReadByte()
	targetNameBytes, _ := pkt.Reader.ReadBytes(20)
	targetName := strings.TrimRight(string(targetNameBytes), "\x00")
	msgLen, _ := pkt.Reader.ReadByte()
	msgBytes, _ := pkt.Reader.ReadBytes(int(msgLen))
	msg := string(msgBytes)
	senderName := h.activeCharName(client)

	targetClient := h.findOnlineClientByCharName(targetName)
	if targetClient == nil || targetClient.CharacterID == 0 {
		failResp := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WHISPERFAIL_ACK)
		failResp.WriteUint16(0x0F69) // fc_NC_ACT_WHISPER_REQ (0x0041C610): target offline -> 3945
		failResp.WriteString(targetName, 20)
		_ = client.SendPacket(failResp)
		return
	}

	// Sender success ACK (22 + len bytes)
	succResp := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_WHISPERSUCCESS_ACK)
	succResp.WriteByte(itemLinkCount)
	succResp.WriteString(targetName, 20)
	succResp.WriteByte(byte(len(msg)))
	succResp.WriteBytes([]byte(msg))
	_ = client.SendPacket(succResp)

	// Delivery to receiver: PROTO_NC_ACT_SOMEONEWHISPER_CMD (OpCode 0x200D)
	// WorldManager.exe fc_NC_ACT_WHISPER_REQ (0x41c80a): itemLinkDataCount[1] + talker[20] + flag[1] + len[1] + content[len]
	var gmFlag byte
	if client.Admin > 0 {
		gmFlag = 1
	}
	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEWHISPER_CMD)
	cmd.WriteByte(itemLinkCount)
	cmd.WriteString(senderName, 20)
	cmd.WriteByte(gmFlag)
	cmd.WriteByte(byte(len(msg)))
	cmd.WriteBytes([]byte(msg))
	_ = targetClient.SendPacket(cmd)

	// Coordinar entrega viva si el destinatario ya transfirió su sesión activa a ZoneService
	h.BroadcastToZones("wms_NC_CT_CHARTTING_CMD", map[string]interface{}{
		"type":          "whisper",
		"sender":        senderName,
		"target":        targetName,
		"message":       msg,
		"itemLinkCount": itemLinkCount,
		"senderFlag":    gmFlag,
	})
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
	clients := make([]*network.Client, 0, len(h.activeClients))
	for _, c := range h.activeClients {
		clients = append(clients, c)
	}
	h.mu.Unlock()
	for _, c := range clients {
		_ = c.SendPacket(cmd)
	}

	h.BroadcastToZones("wms_NC_CT_CHARTTING_CMD", map[string]interface{}{
		"type":    "notice",
		"message": noticeMsg,
	})
}

func (h *Handler) handlePartyChat(client *network.Client, pkt *protocol.Packet) {
	itemLinkCount, _ := pkt.Reader.ReadByte()
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

	// Struct oficial PDB: PROTO_NC_ACT_PARTYCHAT_CMD (talker[20], itemLinkDataCount[1], len[1], content[len])
	cmd := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PARTYCHAT_CMD)
	cmd.WriteString(h.activeCharName(client), 20)
	cmd.WriteByte(itemLinkCount)
	cmd.WriteByte(byte(len(msg)))
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
	var reason, remark string
	if pkt.Reader.Remaining() >= 80 {
		reason, _ = pkt.Reader.ReadString(16)
		remark, _ = pkt.Reader.ReadString(64)
	} else {
		reason, _ = pkt.Reader.ReadString(pkt.Reader.Remaining())
	}

	cleanTarget := strings.TrimSpace(targetName)
	var errCode uint16
	if cleanTarget == "" || minutes <= 0 {
		errCode = 1
	} else {
		if h.charRepo != nil {
			_ = h.charRepo.AddPrison(h.activeCharName(client), cleanTarget, minutes, reason, remark)
		}
		h.prisonMu.Lock()
		if h.activePrisons == nil {
			h.activePrisons = make(map[string]*PrisonSentence)
		}
		h.activePrisons[strings.ToLower(cleanTarget)] = &PrisonSentence{
			CharName:  cleanTarget,
			Reason:    reason,
			Remark:    remark,
			ExpiresAt: time.Now().Add(time.Duration(minutes) * time.Minute),
		}
		h.prisonMu.Unlock()
	}

	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_ADD_GM_ACK)
	resp.WriteString(targetName, 20)
	resp.WriteUint16(errCode)
	_ = client.SendPacket(resp)

	if errCode == 0 {
		h.BroadcastToZones("wms_NC_ZONE_PRISON_GO_REQ", map[string]interface{}{
			"charName": cleanTarget,
		})
	}
}

func (h *Handler) handlePrisonTakeOut(client *network.Client, pkt *protocol.Packet) {
	targetCharID := client.CharacterID
	if pkt.Reader.Remaining() >= 4 {
		v, _ := pkt.Reader.ReadUint32()
		if v > 0 {
			targetCharID = int(v)
		}
	}
	var targetName string
	if h.charRepo != nil && targetCharID > 0 {
		_ = h.charRepo.EndPrison(targetCharID)
		if ch, err := h.charRepo.GetCharacterByID(targetCharID); err == nil && ch != nil {
			targetName = ch.Name
		}
	}
	if targetName != "" {
		h.prisonMu.Lock()
		delete(h.activePrisons, strings.ToLower(strings.TrimSpace(targetName)))
		h.prisonMu.Unlock()
	}
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_END_ACK)
	resp.WriteUint32(uint32(targetCharID))
	resp.WriteUint16(0)
	_ = client.SendPacket(resp)

	h.BroadcastToZones("wms_NC_ZONE_PRISON_END_CMD", map[string]interface{}{
		"charName": targetName,
		"char_id":  targetCharID,
	})
}

func (h *Handler) handlePrisonAlterTime(client *network.Client, pkt *protocol.Packet) {
	targetName, _ := pkt.Reader.ReadString(20)
	minutes, _ := pkt.Reader.ReadInt16()
	var reason, remark string
	if pkt.Reader.Remaining() >= 80 {
		reason, _ = pkt.Reader.ReadString(16)
		remark, _ = pkt.Reader.ReadString(64)
	} else {
		reason, _ = pkt.Reader.ReadString(pkt.Reader.Remaining())
	}

	cleanTarget := strings.TrimSpace(targetName)
	var errCode uint16
	if cleanTarget == "" || minutes <= 0 {
		errCode = 1
	} else {
		if h.charRepo != nil {
			_ = h.charRepo.AlterPrison(h.activeCharName(client), cleanTarget, minutes, reason, remark)
		}
		h.prisonMu.Lock()
		if h.activePrisons == nil {
			h.activePrisons = make(map[string]*PrisonSentence)
		}
		h.activePrisons[strings.ToLower(cleanTarget)] = &PrisonSentence{
			CharName:  cleanTarget,
			Reason:    reason,
			Remark:    remark,
			ExpiresAt: time.Now().Add(time.Duration(minutes) * time.Minute),
		}
		h.prisonMu.Unlock()
	}

	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_ALTER_GM_ACK)
	resp.WriteString(targetName, 20)
	resp.WriteUint16(errCode)
	_ = client.SendPacket(resp)
}

func (h *Handler) handlePrisonConfirm(client *network.Client, pkt *protocol.Packet) {
	// Struct oficial PDB: PROTO_NC_PRISON_GET_ACK (84 bytes: err[2], nMinute[2], sReason[16], sRemark[64])
	charName := strings.ToLower(strings.TrimSpace(h.activeCharName(client)))
	var remMin uint16
	var reason, remark string
	h.prisonMu.Lock()
	if st := h.activePrisons[charName]; st != nil {
		diff := time.Until(st.ExpiresAt)
		if diff > 0 {
			remMin = uint16((diff.Seconds() + 59) / 60)
			reason = st.Reason
			remark = st.Remark
		}
	}
	h.prisonMu.Unlock()
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_GET_ACK)
	resp.WriteUint16(0)      // err
	resp.WriteUint16(remMin) // nMinute
	resp.WriteString(reason, 16)
	resp.WriteString(remark, 64)
	_ = client.SendPacket(resp)
}

// OnPrisonTimer ejecuta CPrisonServer::OnPrisonTimer (WorldManager.exe 0x00458360):
// libera automáticamente personajes cuya condena venció y notifica a ZoneService vía S2S.
func (h *Handler) OnPrisonTimer(now time.Time) int {
	h.prisonMu.Lock()
	var expired []*PrisonSentence
	for k, st := range h.activePrisons {
		if st == nil || !now.Before(st.ExpiresAt) {
			if st != nil {
				expired = append(expired, st)
			}
			delete(h.activePrisons, k)
		}
	}
	h.prisonMu.Unlock()

	for _, st := range expired {
		var charID int
		if h.charRepo != nil {
			if ch, err := h.charRepo.GetCharacterByName(st.CharName); err == nil && ch != nil {
				charID = ch.ID
				_ = h.charRepo.EndPrison(charID)
			}
		}
		h.BroadcastToZones("wms_NC_ZONE_PRISON_END_CMD", map[string]interface{}{
			"charName": st.CharName,
			"char_id":  charID,
		})
	}
	return len(expired)
}

// ForcePrisonExpiryForTest adelanta el vencimiento de prisión de un personaje para tests determinísticos.
func (h *Handler) ForcePrisonExpiryForTest(charName string) {
	h.prisonMu.Lock()
	defer h.prisonMu.Unlock()
	if st := h.activePrisons[strings.ToLower(strings.TrimSpace(charName))]; st != nil {
		st.ExpiresAt = time.Now().Add(-1 * time.Second)
	}
}

// ActivePrisonCount retorna la cantidad de condenas activas registradas en WorldManager.
func (h *Handler) ActivePrisonCount() int {
	h.prisonMu.Lock()
	defer h.prisonMu.Unlock()
	return len(h.activePrisons)
}

// RunWorldTick ejecuta un ciclo autónomo completo de WorldManager (KQ + Prison + GuildWar).
func (h *Handler) RunWorldTick(now time.Time) {
	if h.kqMgr != nil {
		h.kqMgr.DoKingdomQuest(now)
	}
	if h.guildMgr != nil {
		h.guildMgr.UpdateGuildWars(now)
	}
	h.OnPrisonTimer(now)
}

// KQManager retorna el gestor KQ (usado por suites de auditoría).
func (h *Handler) KQManager() *KQManager { return h.kqMgr }

// GuildManager retorna el gestor de gremios (usado por suites de auditoría).
func (h *Handler) GuildManager() *GuildManager { return h.guildMgr }

// PartyManager retorna el gestor de grupos (usado por suites de auditoría).
func (h *Handler) PartyManager() *PartyManager { return h.partyMgr }

// ============================================================================
// DEPARTAMENTO 32: REPORT (tReport en AccountLog)
// ============================================================================

func (h *Handler) handleReportAdd(client *network.Client, pkt *protocol.Packet) {
	targetName, _ := pkt.Reader.ReadString(20)
	var reportType, reason, chatLog string
	if pkt.Reader.Remaining() >= 274 {
		reportType, _ = pkt.Reader.ReadString(16)
		reason, _ = pkt.Reader.ReadString(256)
		chatLogSize, _ := pkt.Reader.ReadUint16()
		if int(chatLogSize) > 0 && pkt.Reader.Remaining() >= int(chatLogSize) {
			chatLog, _ = pkt.Reader.ReadString(int(chatLogSize))
		} else {
			chatLog, _ = pkt.Reader.ReadString(pkt.Reader.Remaining())
		}
	} else {
		reportType = "IN_GAME"
		reason, _ = pkt.Reader.ReadString(64)
		chatLog, _ = pkt.Reader.ReadString(pkt.Reader.Remaining())
	}

	var regNo int
	if h.logRepo != nil {
		r := &database.Report{
			Reporter:       client.AccountID,
			ReporterCharID: h.activeCharName(client),
			Harmer:         targetName,
			WorldNo:        h.worldID,
			ReportType:     reportType,
			Reason:         reason,
			ChatLog:        chatLog,
		}
		regNo, _ = h.logRepo.AddReport(r)
	}

	// Struct oficial PDB: PROTO_NC_REPORT_ADD_ACK (26 bytes: sHarmer[20], nRegNumber[4], err[2])
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_REPORT_ADD_ACK)
	resp.WriteString(targetName, 20)
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
	if strings.EqualFold(h.activeCharName(client), upCharID) {
		targetName = charID
	}

	targetClient := h.findOnlineClientByCharName(targetName)
	var targetID int
	if targetClient != nil {
		targetID = targetClient.CharacterID
	}

	h.holyPromiseMgr.HandleSetUpReq(client, charID, upCharID, targetClient, targetID)
}

func (h *Handler) handleHolyPromiseSetUpConfirm(client *network.Client, pkt *protocol.Packet) {
	charID, _ := pkt.Reader.ReadString(20)
	upCharID, _ := pkt.Reader.ReadString(20)
	accept, _ := pkt.Reader.ReadByte()

	requesterName := charID
	if strings.EqualFold(h.activeCharName(client), charID) {
		requesterName = upCharID
	}

	reqClient := h.findOnlineClientByCharName(requesterName)
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
	pupilClient := h.findOnlineClientByCharName(downCharName)
	var pupilID int
	if pupilClient != nil {
		pupilID = pupilClient.CharacterID
	} else {
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

func writeIndunRankEmptyTab(client *network.Client, op uint16) {
	// Struct oficial PDB: PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_*_ACK (19 bytes header: 3*4B shape, uint16 err, uint32 rank, uint8 count)
	resp := protocol.NewPacketWriter(op)
	resp.Fill(12, 0)    // shapeInfo1st, 2nd, 3th (3 * 4B)
	resp.WriteUint16(0) // nErr
	resp.WriteUint32(0) // nRank
	resp.WriteByte(0)   // nCount
	_ = client.SendPacket(resp)
}

func (h *Handler) handleBattleRankClearTime(client *network.Client, pkt *protocol.Packet) {
	writeIndunRankEmptyTab(client, protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGEDEALT_ACK)
}

func (h *Handler) handleBattleRankScore(client *network.Client, pkt *protocol.Packet) {
	writeIndunRankEmptyTab(client, protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGETAKEN_ACK)
}

func (h *Handler) handleBattleRankPoint(client *network.Client, pkt *protocol.Packet) {
	writeIndunRankEmptyTab(client, protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_HEALINGDONE_ACK)
}

func (h *Handler) handleBattleRankMyScore(client *network.Client, pkt *protocol.Packet) {
	// Struct oficial PDB: PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_MYRANK_ACK (62 bytes)
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_MYRANK_ACK)
	resp.Fill(62, 0)
	_ = client.SendPacket(resp)
}

// ============================================================================
// DEPARTAMENTO 42: CHAT RESTRICT (tChatBlock)
// ============================================================================

func (h *Handler) handleChatRestrictAdd(client *network.Client, pkt *protocol.Packet) {
	var targetName string
	var targetCharID uint32
	if pkt.Reader.Remaining() >= 20 {
		targetName, _ = pkt.Reader.ReadString(20)
		if h.charRepo != nil {
			if ch, err := h.charRepo.GetCharacterByName(targetName); err == nil && ch != nil {
				targetCharID = uint32(ch.ID)
			}
		}
	} else {
		targetCharID, _ = pkt.Reader.ReadUint32()
	}
	if h.charRepo != nil && targetCharID > 0 {
		_ = h.charRepo.AddChatRestrict(client.AccountID, int(targetCharID))
	}
	// Struct oficial PDB: PROTO_NC_CHAT_RESTRICT_ADD_ACK (24 bytes: nError[2], sCharID[20], bIsGM[1], bIsDeleted[1])
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_CHAT_RESTRICT_ADD_ACK)
	resp.WriteUint16(0)
	resp.WriteString(targetName, 20)
	resp.WriteByte(0)
	resp.WriteByte(0)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleChatRestrictDel(client *network.Client, pkt *protocol.Packet) {
	var targetName string
	var targetCharID uint32
	if pkt.Reader.Remaining() >= 20 {
		targetName, _ = pkt.Reader.ReadString(20)
		if h.charRepo != nil {
			if ch, err := h.charRepo.GetCharacterByName(targetName); err == nil && ch != nil {
				targetCharID = uint32(ch.ID)
			}
		}
	} else {
		targetCharID, _ = pkt.Reader.ReadUint32()
	}
	if h.charRepo != nil && targetCharID > 0 {
		_ = h.charRepo.DelChatRestrict(client.AccountID, int(targetCharID))
	}
	// Struct oficial PDB: PROTO_NC_CHAT_RESTRICT_DEL_ACK (22 bytes: nError[2], sChatRestrictCharID[20])
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_CHAT_RESTRICT_DEL_ACK)
	resp.WriteUint16(0)
	resp.WriteString(targetName, 20)
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
