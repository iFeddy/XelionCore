package world

import (
	"log"
	"strings"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func (h *Handler) getAvatarRepo() database.CharacterAvatarRepository {
	if ar, ok := h.charRepo.(database.CharacterAvatarRepository); ok {
		return ar
	}
	return nil
}

func (h *Handler) getProgRepo() database.CharacterProgressionRepository {
	if pr, ok := h.charRepo.(database.CharacterProgressionRepository); ok {
		return pr
	}
	return nil
}

func (h *Handler) getTopoRepo() database.ZoneTopologyRepository {
	if tr, ok := h.charRepo.(database.ZoneTopologyRepository); ok {
		return tr
	}
	return nil
}

func (h *Handler) getSessionRepo() database.LoginSessionRepository {
	if h.sessionRepo != nil {
		return h.sessionRepo
	}
	if sr, ok := h.charRepo.(database.LoginSessionRepository); ok {
		return sr
	}
	return nil
}

func (h *Handler) getLogRepo() database.AccountLogRepository {
	if h.logRepo != nil {
		return h.logRepo
	}
	if lr, ok := h.charRepo.(database.AccountLogRepository); ok {
		return lr
	}
	return nil
}

// handleNormalLogout procesa el logout oficial de cliente (PROTO_NC_USER_NORMALLOGOUT_CMD - 0x0C18).
func (h *Handler) handleNormalLogout(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}

	logoutType, _ := pkt.Reader.ReadByte()
	log.Printf("[World][%s] Normal logout requested (type: %d) for user %d (%s)",
		client.RemoteAddr(), logoutType, client.AccountID, client.Username)

	if sessionRepo := h.getSessionRepo(); sessionRepo != nil {
		playMin, err := sessionRepo.CloseSessionAndGetPlayMin(client.AccountID)
		if err == nil {
			if logRepo := h.getLogRepo(); logRepo != nil {
				_ = logRepo.LogAccount(client.Username, 2, client.AccountID, uint16(playMin), "", h.worldID, 0, client.RemoteIP())
			}
		}
	}

	if client.CharacterID > 0 {
		_ = h.charRepo.CharacterLogout(client.CharacterID, 0)
	}

	client.Close()
}

// handleAvatarList despacha la lista actualizada de avatares (PROTO_NC_USER_AVATAR_LIST_REQ - 0x0C1F).
func (h *Handler) handleAvatarList(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}
	h.sendCharacterList(client)
}

// handleRenameCharacter procesa el renombramiento de personaje (PROTO_NC_AVATAR_RENAME_REQ - 0x140F).
func (h *Handler) handleRenameCharacter(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil || slot > 5 {
		return
	}

	newName, err := pkt.Reader.ReadString(20)
	if err != nil {
		return
	}
	newName = strings.TrimSpace(newName)

	if len(newName) < 3 || len(newName) > 16 {
		resp := protocol.New_PROTO_NC_AVATAR_RENAME_ACK(slot, newName, 1)
		_ = client.SendPacket(resp)
		return
	}

	if h.dataStore != nil && h.dataStore.BadNames != nil && h.dataStore.BadNames.IsBadName(newName) {
		resp := protocol.New_PROTO_NC_AVATAR_RENAME_ACK(slot, newName, 2)
		_ = client.SendPacket(resp)
		return
	}

	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err != nil {
		resp := protocol.New_PROTO_NC_AVATAR_RENAME_ACK(slot, newName, 1)
		_ = client.SendPacket(resp)
		return
	}

	var targetChar *database.Character
	for _, ch := range chars {
		if ch.Slot == slot {
			targetChar = ch
			break
		}
	}
	if targetChar == nil {
		resp := protocol.New_PROTO_NC_AVATAR_RENAME_ACK(slot, newName, 1)
		_ = client.SendPacket(resp)
		return
	}

	avatarRepo := h.getAvatarRepo()
	if avatarRepo == nil {
		resp := protocol.New_PROTO_NC_AVATAR_RENAME_ACK(slot, newName, 1)
		_ = client.SendPacket(resp)
		return
	}

	if err := avatarRepo.ChangeCharacterName(targetChar.ID, newName); err != nil {
		log.Printf("[World][%s] Error renaming character in slot %d: %v", client.RemoteAddr(), slot, err)
		resp := protocol.New_PROTO_NC_AVATAR_RENAME_ACK(slot, newName, 3)
		_ = client.SendPacket(resp)
		return
	}

	targetChar.Name = newName
	resp := protocol.New_PROTO_NC_AVATAR_RENAME_ACK(slot, newName, 0)
	_ = client.SendPacket(resp)
}

// handleChangeCharID ejecuta el cambio de ID persistido (PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ - 0xB409).
func (h *Handler) handleChangeCharID(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil || slot > 5 {
		return
	}

	newName, err := pkt.Reader.ReadString(20)
	if err != nil {
		return
	}
	newName = strings.TrimSpace(newName)

	if len(newName) < 3 || len(newName) > 16 {
		resp := protocol.New_PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK(1, slot, newName)
		_ = client.SendPacket(resp)
		return
	}

	if h.dataStore != nil && h.dataStore.BadNames != nil && h.dataStore.BadNames.IsBadName(newName) {
		resp := protocol.New_PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK(2, slot, newName)
		_ = client.SendPacket(resp)
		return
	}

	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err != nil {
		resp := protocol.New_PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK(1, slot, newName)
		_ = client.SendPacket(resp)
		return
	}

	var targetChar *database.Character
	for _, ch := range chars {
		if ch.Slot == slot {
			targetChar = ch
			break
		}
	}
	if targetChar == nil {
		resp := protocol.New_PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK(1, slot, newName)
		_ = client.SendPacket(resp)
		return
	}

	avatarRepo := h.getAvatarRepo()
	if avatarRepo == nil {
		resp := protocol.New_PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK(1, slot, newName)
		_ = client.SendPacket(resp)
		return
	}

	rowNo := 0
	if targetChar.ChangeID != nil {
		rowNo = int(targetChar.ChangeID.RowNo)
	}
	if rowNo == 0 {
		list, _ := avatarRepo.GetCharacterChangeIDList(client.AccountID)
		if len(list) > 0 {
			rowNo = int(list[0].RowNo)
		}
	}

	ret, err := avatarRepo.ExecuteCharacterChangeID(rowNo, targetChar.ID, newName, 0, client.AccountID)
	if err != nil || ret != 0 {
		log.Printf("[World][%s] Error executing ChangeCharID for char %d (ret=%d): %v", client.RemoteAddr(), targetChar.ID, ret, err)
		resp := protocol.New_PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK(3, slot, newName)
		_ = client.SendPacket(resp)
		return
	}

	targetChar.Name = newName
	if targetChar.ChangeID != nil {
		targetChar.ChangeID.NeedChangeID = 0
	}

	resp := protocol.New_PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ACK(0, slot, newName)
	_ = client.SendPacket(resp)
}

// handleSetStyleGetInfo consulta los cupones de belleza disponibles (PROTO_NC_CHAR_SET_STYLE_GET_INFO_REQ - 0x107C).
func (h *Handler) handleSetStyleGetInfo(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}

	var items [6][2]uint16
	if h.dataStore != nil {
		for grade := 1; grade <= 5; grade++ {
			couponID := h.dataStore.GetStyleCouponID(byte(grade))
			items[grade][0] = couponID
			count := uint16(0)
			if client.CharacterID > 0 {
				if charItems, err := h.charRepo.GetCharacterItems(client.CharacterID); err == nil {
					for _, item := range charItems {
						if item != nil && item.ItemID == couponID {
							count += uint16(item.Flags)
							if item.Flags <= 0 {
								count++
							}
						}
					}
				}
			}
			items[grade][1] = count
		}
	}

	resp := protocol.New_PROTO_NC_CHAR_SET_STYLE_GET_INFO_ACK(0, items)
	_ = client.SendPacket(resp)
}

// handleSetStyle actualiza el corte, color o rostro (PROTO_NC_CHAR_SET_STYLE_REQ - 0x1078).
func (h *Handler) handleSetStyle(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		resp := protocol.New_PROTO_NC_CHAR_SET_STYLE_ACK(1)
		_ = client.SendPacket(resp)
		return
	}

	// PROTO_AVATAR_SHAPE_INFO: 4 bytes (jobGender, hair, color, face)
	jobGender, _ := pkt.Reader.ReadByte()
	hair, _ := pkt.Reader.ReadByte()
	color, _ := pkt.Reader.ReadByte()
	face, _ := pkt.Reader.ReadByte()
	_ = jobGender

	var ok bool
	if h.dataStore != nil {
		_, ok = h.dataStore.ValidateStyleChange(hair, color, face)
	} else {
		ok = true
	}

	if !ok {
		resp := protocol.New_PROTO_NC_CHAR_SET_STYLE_ACK(2)
		_ = client.SendPacket(resp)
		return
	}

	avatarRepo := h.getAvatarRepo()
	if avatarRepo == nil {
		resp := protocol.New_PROTO_NC_CHAR_SET_STYLE_ACK(1)
		_ = client.SendPacket(resp)
		return
	}

	charInfo, err := h.charRepo.GetCharacterByID(client.CharacterID)
	if err != nil || charInfo == nil {
		resp := protocol.New_PROTO_NC_CHAR_SET_STYLE_ACK(1)
		_ = client.SendPacket(resp)
		return
	}

	genderByte := byte(0)
	if charInfo.Male {
		genderByte = 1
	}

	if err := avatarRepo.SetCharacterShape(client.CharacterID, charInfo.Race, byte(charInfo.Job), genderByte, hair, color, face); err != nil {
		log.Printf("[World][%s] SetCharacterShape failed for char %d: %v", client.RemoteAddr(), client.CharacterID, err)
		resp := protocol.New_PROTO_NC_CHAR_SET_STYLE_ACK(1)
		_ = client.SendPacket(resp)
		return
	}

	charInfo.Hair = hair
	charInfo.HairColor = color
	charInfo.Face = face

	resp := protocol.New_PROTO_NC_CHAR_SET_STYLE_ACK(0)
	_ = client.SendPacket(resp)
}

// handleTutorialPopup registra la elección de omitir o realizar el tutorial (PROTO_NC_CHAR_TUTORIAL_POPUP_ACK - 0x1111).
func (h *Handler) handleTutorialPopup(client *network.Client, pkt *protocol.Packet) {
	bIsSkip, _ := pkt.Reader.ReadByte()
	client.TutorialChosen = true
	client.TutorialSkip = (bIsSkip != 0)
}

// handleAttendanceReward reclama la recompensa de asistencia diaria (PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_REQ - 0x10F1).
func (h *Handler) handleAttendanceReward(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		return
	}

	index, _ := pkt.Reader.ReadByte()
	rewardID, _ := pkt.Reader.ReadByte()

	progRepo := h.getProgRepo()
	if progRepo == nil {
		resp := protocol.New_PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_ACK(index, rewardID, 1)
		_ = client.SendPacket(resp)
		return
	}

	if _, err := progRepo.SetAttendanceRewardAt(client.CharacterID, time.Now(), int(rewardID)); err != nil {
		log.Printf("[World][%s] SetAttendanceRewardAt failed for char %d: %v", client.RemoteAddr(), client.CharacterID, err)
		resp := protocol.New_PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_ACK(index, rewardID, 1)
		_ = client.SendPacket(resp)
		return
	}

	resp := protocol.New_PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_ACK(index, rewardID, 0)
	_ = client.SendPacket(resp)
}

// handleNewbieGuideViewSet marca una guía de novato como vista (PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD - 0x10DA).
func (h *Handler) handleNewbieGuideViewSet(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		return
	}

	guideNo, _ := pkt.Reader.ReadByte()
	if progRepo := h.getProgRepo(); progRepo != nil {
		_ = progRepo.SetInGameGuideView(client.CharacterID, guideNo)
	}
}

// handleGetAllSettings responde con todas las opciones clásicas (PROTO_NC_CHAR_OPTION_GET_ALL_REQ - 0x7001).
func (h *Handler) handleGetAllSettings(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		return
	}
	h.handleGetQuickBar(client, pkt)
	h.handleGetQuickBarState(client, pkt)
	h.handleGetVideoSettings(client, pkt)
	h.handleGetSoundSettings(client, pkt)
	h.handleGetGameSettings(client, pkt)
	h.handleGetClientSettings(client, pkt)
	h.handleGetShortCuts(client, pkt)
}

// handleGetVideoSettings responde opciones de video (PROTO_NC_CHAR_OPTION_GET_VIDEO_REQ - 0x7006).
func (h *Handler) handleGetVideoSettings(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		return
	}
	opts, _ := h.charRepo.GetCharacterOptions(client.CharacterID)
	var data []byte
	if opts != nil && len(opts.Video) > 0 {
		data = opts.Video
	}
	resp := protocol.New_PROTO_NC_CHAR_OPTION_GET_VIDEO_ACK(data)
	_ = client.SendPacket(resp)
}

// handleGetSoundSettings responde opciones de sonido (PROTO_NC_CHAR_OPTION_GET_SOUND_REQ - 0x7008).
func (h *Handler) handleGetSoundSettings(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		return
	}
	opts, _ := h.charRepo.GetCharacterOptions(client.CharacterID)
	var data []byte
	if opts != nil && len(opts.Sound) > 0 {
		data = opts.Sound
	}
	resp := protocol.New_PROTO_NC_CHAR_OPTION_GET_SOUND_ACK(data)
	_ = client.SendPacket(resp)
}

// handleImproveSetShortcut guarda accesos directos mejorados (PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ - 0x7037).
func (h *Handler) handleImproveSetShortcut(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK(1)
		_ = client.SendPacket(resp)
		return
	}
	cnt, _ := pkt.Reader.ReadByte()
	var entries []database.ImproveShortcutEntry
	for i := 0; i < int(cnt) && pkt.Reader.Remaining() >= 7; i++ {
		slot, _ := pkt.Reader.ReadByte()
		code, _ := pkt.Reader.ReadUint16()
		val, _ := pkt.Reader.ReadInt32()
		entries = append(entries, database.ImproveShortcutEntry{
			SlotNo: slot,
			CodeNo: code,
			Value:  val,
		})
	}
	if progRepo := h.getProgRepo(); progRepo != nil {
		for _, e := range entries {
			_ = progRepo.SetImproveShortCut(client.CharacterID, e.SlotNo, e.CodeNo, e.Value)
		}
	}
	resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK(0)
	_ = client.SendPacket(resp)
}

// handleImproveSetKeymap guarda mapeo de teclas mejorado (PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_REQ - 0x7039).
func (h *Handler) handleImproveSetKeymap(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_ACK(1)
		_ = client.SendPacket(resp)
		return
	}
	cnt, _ := pkt.Reader.ReadUint16()
	var entries []database.ImproveKeyMapEntry
	for i := 0; i < int(cnt) && pkt.Reader.Remaining() >= 4; i++ {
		fn, _ := pkt.Reader.ReadUint16()
		ext, _ := pkt.Reader.ReadByte()
		asc, _ := pkt.Reader.ReadByte()
		entries = append(entries, database.ImproveKeyMapEntry{
			FunctionNo: fn,
			ExtendKey:  ext,
			ASCIICode:  asc,
		})
	}
	if progRepo := h.getProgRepo(); progRepo != nil {
		for _, e := range entries {
			_ = progRepo.SetImproveKeyMap(client.CharacterID, e.FunctionNo, e.ExtendKey, e.ASCIICode)
		}
	}
	resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_ACK(0)
	_ = client.SendPacket(resp)
}

// handleImproveSetGameOption guarda opciones de juego mejoradas (PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_REQ - 0x703B).
func (h *Handler) handleImproveSetGameOption(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_ACK(1)
		_ = client.SendPacket(resp)
		return
	}
	cnt, _ := pkt.Reader.ReadUint16()
	var entries []database.ImproveGameOptionEntry
	for i := 0; i < int(cnt) && pkt.Reader.Remaining() >= 3; i++ {
		optNo, _ := pkt.Reader.ReadUint16()
		val, _ := pkt.Reader.ReadByte()
		entries = append(entries, database.ImproveGameOptionEntry{
			OptionNo: optNo,
			Value:    val,
		})
	}
	if progRepo := h.getProgRepo(); progRepo != nil {
		for _, e := range entries {
			_ = progRepo.SetImproveGameOption(client.CharacterID, e.OptionNo, e.Value)
		}
	}
	resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_ACK(0)
	_ = client.SendPacket(resp)
}

// handleImproveInitKeymap restaura mapeo de teclas a valores de fábrica (PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_REQ - 0x7043).
func (h *Handler) handleImproveInitKeymap(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_ACK(1, nil)
		_ = client.SendPacket(resp)
		return
	}
	var entries []protocol.KEY_MAP_DATA
	if progRepo := h.getProgRepo(); progRepo != nil {
		if err := progRepo.InitImproveKeyMap(client.CharacterID); err == nil {
			if dbEntries, err := progRepo.GetImproveKeyMap(client.CharacterID); err == nil {
				for _, e := range dbEntries {
					entries = append(entries, protocol.KEY_MAP_DATA{
						FunctionNo: e.FunctionNo,
						ExtendKey:  e.ExtendKey,
						ASCIICode:  e.ASCIICode,
					})
				}
			}
		}
	}
	resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_ACK(0, entries)
	_ = client.SendPacket(resp)
}

// handleImproveInitGameOption restaura opciones de juego a valores de fábrica (PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_REQ - 0x7045).
func (h *Handler) handleImproveInitGameOption(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.CharacterID <= 0 {
		resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_ACK(1, nil)
		_ = client.SendPacket(resp)
		return
	}
	var entries []protocol.GAME_OPTION_DATA
	if progRepo := h.getProgRepo(); progRepo != nil {
		if err := progRepo.InitImproveGameOptions(client.CharacterID); err == nil {
			if dbEntries, err := progRepo.GetImproveGameOptions(client.CharacterID); err == nil {
				for _, e := range dbEntries {
					entries = append(entries, protocol.GAME_OPTION_DATA{
						OptionNo: e.OptionNo,
						Value:    e.Value,
					})
				}
			}
		}
	}
	resp := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_ACK(0, entries)
	_ = client.SendPacket(resp)
}

// handleLoginBlockDataErr registra error de datos bloqueados (PROTO_NC_CHAR_OPTION_LOGIN_BLOCKDATA_ERR - 0x705F).
func (h *Handler) handleLoginBlockDataErr(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[World][%s] Received PROTO_NC_CHAR_OPTION_LOGIN_BLOCKDATA_ERR", client.RemoteAddr())
}
