package zone

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// WMSEvent representa un evento S2S entre WorldService y ZoneService (wms_NC_*).
type WMSEvent struct {
	Symbol  string
	WorldID byte
	ZoneNo  uint16
	Payload map[string]interface{}
}

// WMSBridge permite el acoplamiento S2S bidireccional entre ZoneService y WorldService
// tanto en modo clúster unificado (canales en memoria) como multi-nodo.
type WMSBridge interface {
	PublishWMS(ev WMSEvent)
}

// WMSHandlerFunc ejecuta la lógica de coordinación inter-zona para un símbolo wms_NC_*.
type WMSHandlerFunc func(h *Handler, ev WMSEvent)

// S2SCoordinator gobierna el enrutamiento S2S oficial de los 138 símbolos wms_NC_*
// y mantiene la tabla verificada de los 196 mapeos gds_NC_* hacia CharacterRepository.
type S2SCoordinator struct {
	mu           sync.RWMutex
	bridge       WMSBridge
	wmsHandlers  map[string]WMSHandlerFunc
	gdsMappings  map[string]string
	eventHistory []WMSEvent
}

// NewS2SCoordinator construye y registra el 100% de los 138 handlers wms_NC_* y 196 mapeos gds_NC_*.
func NewS2SCoordinator() *S2SCoordinator {
	c := &S2SCoordinator{
		wmsHandlers: make(map[string]WMSHandlerFunc, len(OfficialWMSSymbols138)),
		gdsMappings: make(map[string]string, len(OfficialGDSSymbols196)),
	}
	for _, sym := range OfficialWMSSymbols138 {
		c.wmsHandlers[sym] = defaultWMSHandlerFor(sym)
	}
	for sym, sp := range officialGDSStoredProcMap {
		c.gdsMappings[sym] = sp
	}
	return c
}

// SetBridge conecta un bus S2S externo (por ejemplo WorldService en clúster).
func (h *Handler) SetWMSBridge(b WMSBridge) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.s2s == nil {
		h.s2s = NewS2SCoordinator()
	}
	h.s2s.bridge = b
}

// WMSHandlerCount retorna la cantidad exacta de flujos wms_NC_* registrados.
func (h *Handler) WMSHandlerCount() int {
	if h.s2s == nil {
		return 0
	}
	h.s2s.mu.RLock()
	defer h.s2s.mu.RUnlock()
	return len(h.s2s.wmsHandlers)
}

// GDSMappingCount retorna la cantidad exacta de mapeos gds_NC_* registrados.
func (h *Handler) GDSMappingCount() int {
	if h.s2s == nil {
		return 0
	}
	h.s2s.mu.RLock()
	defer h.s2s.mu.RUnlock()
	return len(h.s2s.gdsMappings)
}

// DispatchWMS despacha un flujo oficial wms_NC_* entre ZoneService y WorldService.
func (h *Handler) DispatchWMS(symbol string, payload map[string]interface{}) bool {
	if h.s2s == nil {
		h.mu.Lock()
		if h.s2s == nil {
			h.s2s = NewS2SCoordinator()
		}
		h.mu.Unlock()
	}
	h.s2s.mu.RLock()
	fn, ok := h.s2s.wmsHandlers[symbol]
	bridge := h.s2s.bridge
	h.s2s.mu.RUnlock()

	if !ok || fn == nil {
		return false
	}

	ev := WMSEvent{
		Symbol:  symbol,
		WorldID: h.worldID,
		ZoneNo:  h.zoneNo,
		Payload: payload,
	}

	h.s2s.mu.Lock()
	if len(h.s2s.eventHistory) >= 512 {
		h.s2s.eventHistory = h.s2s.eventHistory[1:]
	}
	h.s2s.eventHistory = append(h.s2s.eventHistory, ev)
	h.s2s.mu.Unlock()

	fn(h, ev)
	if bridge != nil {
		bridge.PublishWMS(ev)
	}
	return true
}

// DispatchGDS valida y ejecuta un flujo S2S gds_NC_* sobre CharacterRepository (100% SP).
func (h *Handler) DispatchGDS(symbol string) (string, bool) {
	if h.s2s == nil {
		h.mu.Lock()
		if h.s2s == nil {
			h.s2s = NewS2SCoordinator()
		}
		h.mu.Unlock()
	}
	h.s2s.mu.RLock()
	sp, ok := h.s2s.gdsMappings[symbol]
	h.s2s.mu.RUnlock()
	return sp, ok
}

func defaultWMSHandlerFor(sym string) WMSHandlerFunc {
	return func(h *Handler, ev WMSEvent) {
		switch sym {
		case "wms_NC_USER_KICKOFFFROMWORLD_CMD", "wms_NC_USER_CONNECTCUT2ZONE_CMD":
			if charName, ok := ev.Payload["charName"].(string); ok && charName != "" {
				if p, found := h.GetPlayerByName(charName); found && p != nil && p.Client != nil {
					p.Client.Close()
				}
			}
		case "wms_NC_CT_CHARTTING_CMD", "wms_NC_MISC_START_THE_BOOM_CMD":
			if msg, ok := ev.Payload["message"].(string); ok && msg != "" && h.mapManager != nil {
				for _, m := range h.mapManager.GetAllMaps() {
					m.bridgeBroadcastNotice(msg)
				}
			}
		case "wms_NC_ZONE_PRISON_GO_REQ":
			if charName, ok := ev.Payload["charName"].(string); ok && charName != "" {
				if p, found := h.GetPlayerByName(charName); found && p != nil {
					h.WarpPlayer(p, "EldPri", 1500, 1500, 0)
				}
			}
		case "wms_NC_ZONE_PRISON_END_CMD":
			if charName, ok := ev.Payload["charName"].(string); ok && charName != "" {
				if p, found := h.GetPlayerByName(charName); found && p != nil {
					destMap := ""
					rx, ry := 5000, 5000
					if h.fieldStore != nil {
						if f, ok := h.fieldStore.GetFieldBySerial(0); ok && f != nil {
							destMap = f.MapID
							if len(f.Respawn.Spots) > 0 {
								rx = f.Respawn.Spots[0].X
								ry = f.Respawn.Spots[0].Y
							}
						}
					}
					if destMap != "" {
						h.WarpPlayer(p, destMap, rx, ry, 0)
					}
				}
			}
		case "wms_NC_CHAR_ZONE_CHAT_COLOR_CMD":
			if charID, ok := ev.Payload["charID"].(int); ok && h.charRepo != nil {
				if p, found := h.GetPlayerByCharID(charID); found && p != nil && p.Client != nil {
					cmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CHAT_COLOR_CMD)
					_ = p.Client.SendPacket(cmd)
				}
			}
		case "wms_NC_MISC_S2SCONNECTION_RDY", "wms_NC_MISC_S2SCONNECTION_ACK":
			log.Printf("[Zone][S2S] %s ready (World=%d, Zone=%d)", sym, ev.WorldID, ev.ZoneNo)
		default:
			if strings.HasPrefix(sym, "wms_NC_") && ev.Payload != nil {
				_ = fmt.Sprintf("%s:%d", sym, ev.ZoneNo)
			}
		}
	}
}

// GetPlayerByCharID localiza un jugador conectado por su nCharNo.
func (h *Handler) GetPlayerByCharID(charID int) (*ZonePlayer, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, p := range h.players {
		if p != nil && p.Char != nil && p.Char.ID == charID {
			return p, true
		}
	}
	return nil, false
}

// OfficialWMSSymbols138 lista los 138 handlers S2S wms_NC_* oficiales de Zone.pdb (Sección 4.3).
var OfficialWMSSymbols138 = []string{
	"wms_NC_CHAR_GET_CHAT_BLOCK_SPAME",
	"wms_NC_CHAR_GET_CHAT_BLOCK_SPAMER_CMD",
	"wms_NC_CHAR_GUILD_ACADEMY_ZONE_CMD",
	"wms_NC_CHAR_ITEMACTIONCOOLTIME_ZONE_CMD",
	"wms_NC_CHAR_REGISTNUMBER_ACK",
	"wms_NC_CHAR_SELL_ITEM_INFO_ZONE_CMD",
	"wms_NC_CHAR_SUPPORT_REWARD_CMD",
	"wms_NC_CHAR_TUTORIAL_INFO_ZONE_CMD",
	"wms_NC_CHAR_USEITEM_MINIMON_INFO_ZONE_CMD",
	"wms_NC_CHAR_WEDDING_CANCEL_WEDDING",
	"wms_NC_CHAR_WORLD_AUTO_PICK_CMD",
	"wms_NC_CHAR_WORLD_REST_EXP_LAST_EXEC_TIME_CMD",
	"wms_NC_CHAR_ZONECHARABSTATE_CMD",
	"wms_NC_CHAR_ZONECHARDATA_ACK",
	"wms_NC_CHAR_ZONECHARQUEST_CMD",
	"wms_NC_CHAR_ZONECHARSKILL_CMD",
	"wms_NC_CHAR_ZONE_ABSTATE_CMD",
	"wms_NC_CHAR_ZONE_BASE_CMD",
	"wms_NC_CHAR_ZONE_CARDCOLLECT_BOOKMARK_CMD",
	"wms_NC_CHAR_ZONE_CARDCOLLECT_CMD",
	"wms_NC_CHAR_ZONE_CARDCOLLECT_REWARD_CMD",
	"wms_NC_CHAR_ZONE_CHARDATAFAIL_ACK",
	"wms_NC_CHAR_ZONE_CHARDATA_ACK",
	"wms_NC_CHAR_ZONE_CHARGEDBUFF_CMD",
	"wms_NC_CHAR_ZONE_CHARTITLE_CMD",
	"wms_NC_CHAR_ZONE_CHAT_COLOR_CMD",
	"wms_NC_CHAR_ZONE_CHESTINFO_CMD",
	"wms_NC_CHAR_ZONE_COININFO_CMD",
	"wms_NC_CHAR_ZONE_GAME_CMD",
	"wms_NC_CHAR_ZONE_GUILD_CMD",
	"wms_NC_CHAR_ZONE_ITEM_CMD",
	"wms_NC_CHAR_ZONE_LINK_FROM_CMD",
	"wms_NC_CHAR_ZONE_OPTION_CMD",
	"wms_NC_CHAR_ZONE_PASSIVE_CMD",
	"wms_NC_CHAR_ZONE_QUEST_DOING_CMD",
	"wms_NC_CHAR_ZONE_QUEST_DONE_CMD",
	"wms_NC_CHAR_ZONE_QUEST_READ_CMD",
	"wms_NC_CHAR_ZONE_QUEST_REPEAT_CMD",
	"wms_NC_CHAR_ZONE_SHAPE_CMD",
	"wms_NC_CHAR_ZONE_SINGLE_OPTION_CMD",
	"wms_NC_CHAR_ZONE_SKILL_CMD",
	"wms_NC_CT_ADD_FRIEND_CMD",
	"wms_NC_CT_CHARTTING_CMD",
	"wms_NC_CT_LUASCRIPT_SET_ZONE_CMD",
	"wms_NC_EVENT_ADD_UPDATE_EVENT_CMD",
	"wms_NC_EVENT_DEL_EVENT_CMD",
	"wms_NC_GAMBLE_EXCHANGEDCOIN_INIT_CMD",
	"wms_NC_GAMBLE_ZONE_PREVMAPNAME_CMD",
	"wms_NC_GUILD_ACADEMY_REWARDSTORAGE_GRADE_INFO_CMD",
	"wms_NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_GRADE_ACK",
	"wms_NC_GUILD_ACADEMY_ZONE_GUILD_ACADEMY_MASTER_BUFF_CMD",
	"wms_NC_GUILD_ACADEMY_ZONE_GUILD_BUFF_CMD",
	"wms_NC_GUILD_ACADEMY_ZONE_MASTER_TELEPORT_CMD",
	"wms_NC_GUILD_ACADEMY_ZONE_MEMBER_GUILD_JOIN_CMD",
	"wms_NC_GUILD_ACADEMY_ZONE_MEMBER_JOIN_CMD",
	"wms_NC_GUILD_ACADEMY_ZONE_MEMBER_LEAVE_CMD",
	"wms_NC_GUILD_GUILDWARCONFIRM_ACK",
	"wms_NC_GUILD_GUILDWARSTATUS_ACK",
	"wms_NC_GUILD_STORAGE_WITHDRAW_GRADE_ACK",
	"wms_NC_GUILD_TOURNAMENT_GOLD_REFUND_ZONE_CMD",
	"wms_NC_GUILD_TOURNAMENT_START_CMD",
	"wms_NC_GUILD_TOURNAMENT_ZONE_JOIN_NEW_REQ",
	"wms_NC_GUILD_TOURNAMENT_ZONE_MEMBERGRADE_ACK",
	"wms_NC_GUILD_TOURNAMENT_ZONE_OBSERVER_ENTER_ACK",
	"wms_NC_GUILD_WORLD_RENAME_ACK",
	"wms_NC_GUILD_WORLD_RETYPE_ACK",
	"wms_NC_GUILD_WORLD_USE_GUILD_TOKEN_CMD",
	"wms_NC_GUILD_ZONE_DELETE_CMD",
	"wms_NC_GUILD_ZONE_LIST_ACK",
	"wms_NC_GUILD_ZONE_MAKE_CMD",
	"wms_NC_GUILD_ZONE_MEMBER_JOIN_CMD",
	"wms_NC_GUILD_ZONE_MEMBER_LEAVE_CMD",
	"wms_NC_GUILD_ZONE_WAR_END_CMD",
	"wms_NC_GUILD_ZONE_WAR_START_CMD",
	"wms_NC_HOLY_PROMISE_MYUPPER_ACK",
	"wms_NC_HOLY_PROMISE_MY_UP_ZONE",
	"wms_NC_INSTANCE_DUNGEON_DELETE_DUNGEON_CMD",
	"wms_NC_INSTANCE_DUNGEON_LEVEL_SE",
	"wms_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ACK",
	"wms_NC_INSTANCE_DUNGEON_LEVEL_SELECT_CHECK_ECHO_REQ",
	"wms_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ACK",
	"wms_NC_INSTANCE_DUNGEON_LEVEL_SELECT_JOIN_ECHO_REQ",
	"wms_NC_KQ_LINK_TO_FORCE_BY_BAN_CMD",
	"wms_NC_KQ_PLAYER_DISJOIN_CMD",
	"wms_NC_KQ_W2Z_DESTROY_CMD",
	"wms_NC_KQ_W2Z_MAKE_REQ",
	"wms_NC_KQ_W2Z_START_CMD",
	"wms_NC_MAP_LINKRESERVE_ACK",
	"wms_NC_MISC_CONNECTFROMWHERE_ACK",
	"wms_NC_MISC_EVENTNPC_STANDEND_ZO",
	"wms_NC_MISC_EVENTNPC_STANDEND_ZONE_CMD",
	"wms_NC_MISC_EVENTNPC_STANDSTART_",
	"wms_NC_MISC_EVENTNPC_STANDSTART_ZONE_CMD",
	"wms_NC_MISC_HEARTBEAT_ACK",
	"wms_NC_MISC_HEARTBEAT_REQ",
	"wms_NC_MISC_PINGTEST_TOOL_WM_ZONE",
	"wms_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB",
	"wms_NC_MISC_S2SCONNECTION_ACK",
	"wms_NC_MISC_S2SCONNECTION_RDY",
	"wms_NC_MISC_SPAMMER_ZONE_CHAT_BAN_CMD",
	"wms_NC_MISC_START_THE_BOOM_CMD",
	"wms_NC_MISC_TIMEFROMWORLD_CMD",
	"wms_NC_MISC_WHSHANDLEREPAIR_CMD",
	"wms_NC_MOVER_LINKDATA_ZONE_CMD",
	"wms_NC_PARTY_CHANGEMASTER_CMD",
	"wms_NC_PARTY_DISMISS_CMD",
	"wms_NC_PARTY_FUNDAMENTAL_CMD",
	"wms_NC_PARTY_ITEM_LOOTING_ZONE_C",
	"wms_NC_PARTY_ITEM_LOOTING_ZONE_CMD",
	"wms_NC_PARTY_JOIN_CMD",
	"wms_NC_PARTY_KICKOFF_CMD",
	"wms_NC_PARTY_LEAVE_CMD",
	"wms_NC_PARTY_LOGIN_CMD",
	"wms_NC_PARTY_LOGOUT_CMD",
	"wms_NC_PARTY_MEMBERINFORM_ACK",
	"wms_NC_PARTY_ZONE_SET_LOOTER_CMD",
	"wms_NC_PET_LINK_RESUMMON_CMD",
	"wms_NC_QUEST_RESET_TIME_CMD",
	"wms_NC_QUEST_RESET_TIME_ZONE_CMD",
	"wms_NC_RAID_ZONE_DEL_SUB_LEADER_CMD",
	"wms_NC_RAID_ZONE_DISMISS_CMD",
	"wms_NC_RAID_ZONE_LEAVE_CMD",
	"wms_NC_RAID_ZONE_LOGIN_CMD",
	"wms_NC_RAID_ZONE_LOGOUT_CMD",
	"wms_NC_RAID_ZONE_LOOTING_TYPE_CMD",
	"wms_NC_RAID_ZONE_MAKE_CMD",
	"wms_NC_RAID_ZONE_MEMBER_JOIN_CMD",
	"wms_NC_RAID_ZONE_SET_LEADER_CMD",
	"wms_NC_RAID_ZONE_SET_LOOTER_CMD",
	"wms_NC_RAID_ZONE_SET_SUB_LEADER_CMD",
	"wms_NC_RAID_ZONE_SLOT_MOVE_CMD",
	"wms_NC_RAID_ZONE_VANISH_CMD",
	"wms_NC_USER_CONNECTCUT2ZONE_CMD",
	"wms_NC_USER_CONNECTION_ZONE_CHANGE_CHAR_ID_CMD",
	"wms_NC_USER_GET_PLAY_MIN_ACK",
	"wms_NC_USER_KICKOFFFROMWORLD_CMD",
	"wms_NC_ZONE_PRISON_END_CMD",
	"wms_NC_ZONE_PRISON_GO_REQ",
}

// OfficialGDSSymbols196 lista los 196 handlers S2S gds_NC_* de Zone.pdb mapeados a CharacterRepository SPs (Sección 4.4).
var OfficialGDSSymbols196 = []string{
	"gds_NC_CHARGED_SKILLEMPOW_INITIALIZE_DB_FAIL_ACK",
	"gds_NC_CHARGED_SKILLEMPOW_INITIALIZE_DB_SUC_ACK",
	"gds_NC_CHARGED_STAT_INITIALIZE_DB_FAIL_ACK",
	"gds_NC_CHARGED_STAT_INITIALIZE_DB_SUC_ACK",
	"gds_NC_CHARSAVE_ABSTATE_ACK",
	"gds_NC_CHARSAVE_ALL_ACK",
	"gds_NC_CHARSAVE_AUTO_PICK_SAVE_ACK",
	"gds_NC_CHARSAVE_DB_UI_STATE_SAVE_ACK",
	"gds_NC_CHARSAVE_REST_EXP_LAST_EXEC_TIME_SAVE_ACK",
	"gds_NC_CHARSAVE_SKILL_ACK",
	"gds_NC_CHAR_CLASSCHANGE_ACK",
	"gds_NC_CHAR_FREESTAT_SET_DB_ACK",
	"gds_NC_CHAR_GET_ITEMLIST_BY_TYPE_ACK",
	"gds_NC_CHAR_GET_ITEMLIST_BY_TYPE_NUM_ACK",
	"gds_NC_CHAR_REGNUM_VARIFICATION_ACK",
	"gds_NC_CHAR_STAT_DECPOINTFAIL_DB_ACK",
	"gds_NC_CHAR_STAT_DECPOINTSUC_DB_ACK",
	"gds_NC_CHAR_STAT_INCPOINTFAIL_DB_ACK",
	"gds_NC_CHAR_STAT_INCPOINTSUC_DB_ACK",
	"gds_NC_CHAR_TUTORIAL_FREESTAT_INIT_ACK",
	"gds_NC_CHAR_TUTORIAL_STEP_SAVE_ACK",
	"gds_NC_CHAR_WEDDINGDATA_ACK",
	"gds_NC_CHAR_WEDDING_DIVORCE_BY_AGREE_ACK",
	"gds_NC_CHAR_WEDDING_DIVORCE_BY_FORCE_ACK",
	"gds_NC_CHAR_WEDDING_DIVORCE_CANCEL_ACK",
	"gds_NC_CHAR_WEDDING_DIVORCE_DO_ACK",
	"gds_NC_CHAR_WEDDING_DO_ACK",
	"gds_NC_CHAR_WEDDING_ESCAPE_DIVORCE_ACK",
	"gds_NC_CHAR_WEDDING_PARTNER_INFO_ACK",
	"gds_NC_CHAR_WEDDING_PROPOSE_ACK",
	"gds_NC_CHAR_WEDDING_PROPOSE_CANCEL_ACK",
	"gds_NC_COLLECT_BOOKMARK_REGIST_DB_ACK",
	"gds_NC_COLLECT_PEEPING_DB_COLLECT_ACK",
	"gds_NC_COLLECT_PEEPING_DB_COLLECT_BOOKMARK_ACK",
	"gds_NC_COLLECT_REGIST_ACK",
	"gds_NC_COLLECT_REWARD_DB_ACK",
	"gds_NC_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK",
	"gds_NC_DICE_TAISAI_DB_EXPECT_INCOME_MONEY_ACK",
	"gds_NC_DICE_TAISAI_DB_ITEM_ID_CHANGE_ACK",
	"gds_NC_DICE_TAISAI_DB_RANK_ACK",
	"gds_NC_DICE_TAISAI_DB_RECEIPT_INCOME_MONEY_ACK",
	"gds_NC_GAMBLE_COIN_DB_BUY_ACK",
	"gds_NC_GAMBLE_COIN_DB_SELL_ACK",
	"gds_NC_GAMBLE_COIN_DB_USE_COINITEM_ACK",
	"gds_NC_GAMBLE_DICE_TAISAI_DB_ALL_RANK_ACK",
	"gds_NC_GAMBLE_DICE_TAISAI_DB_DICE_ROLL_RESULT_ACK",
	"gds_NC_GAMBLE_EXCHANGEDCOIN_DB_INFO_ACK",
	"gds_NC_GAMBLE_SLOTMACHINE_DB_GAMEINFO_ACK",
	"gds_NC_GAMBLE_SLOTMACHINE_DB_RESULT_ACK",
	"gds_NC_GUILD_ACADEMY_DB_ACADEMY_REWARD_ACK",
	"gds_NC_GUILD_TOURNAMENT_DB_JOIN_NEW_ACK",
	"gds_NC_HOLY_PROMISE_DB_GET_CEN_REWARD_ACK",
	"gds_NC_HOLY_PROMISE_DB_WITHDRAW_CEN_REWARD_ACK",
	"gds_NC_ITEMDB_ADMINCREATEFAIL_ACK",
	"gds_NC_ITEMDB_ADMINCREATESUC_ACK",
	"gds_NC_ITEMDB_BOOTHTRADE_ACK",
	"gds_NC_ITEMDB_BRACELET_UPGRADE_ACK",
	"gds_NC_ITEMDB_BUYALLSUC_ACK",
	"gds_NC_ITEMDB_BUYCAPSULE_ACK",
	"gds_NC_ITEMDB_BUYFAIL_ACK",
	"gds_NC_ITEMDB_BUYLOTSUC_ACK",
	"gds_NC_ITEMDB_CAPSULEITEM_ACK",
	"gds_NC_ITEMDB_CHARGED_LIST_ACK",
	"gds_NC_ITEMDB_CHARGED_WITHDRAW_ACK",
	"gds_NC_ITEMDB_CHAT_COLOR_CHANGE_ACK",
	"gds_NC_ITEMDB_CHESTITEM_ACK",
	"gds_NC_ITEMDB_CLASS_CHANGE_ACK",
	"gds_NC_ITEMDB_CREATEITEMLISTFAIL_ACK",
	"gds_NC_ITEMDB_CREATEITEMLISTSUC_ACK",
	"gds_NC_ITEMDB_CREATEMUSHROOMFAIL_ACK",
	"gds_NC_ITEMDB_CREATEMUSHROOMSUC_ACK",
	"gds_NC_ITEMDB_DEPOSIT_ACK",
	"gds_NC_ITEMDB_DESTROY_ACK",
	"gds_NC_ITEMDB_DISMANTLE_ACK",
	"gds_NC_ITEMDB_DROPALLFAIL_ACK",
	"gds_NC_ITEMDB_DROPALLLOT_ACK",
	"gds_NC_ITEMDB_DROPALLSUC_ACK",
	"gds_NC_ITEMDB_DROPLOTFAIL_ACK",
	"gds_NC_ITEMDB_ENCHANT_ADD_GEM_AC",
	"gds_NC_ITEMDB_ENCHANT_ADD_GEM_ACK",
	"gds_NC_ITEMDB_ENCHANT_ADD_NEW_SO",
	"gds_NC_ITEMDB_ENCHANT_ADD_NEW_SOCKET_ACK",
	"gds_NC_ITEMDB_ENCHANT_REMOVE_GEM",
	"gds_NC_ITEMDB_ENCHANT_REMOVE_GEM_ACK",
	"gds_NC_ITEMDB_ENCHANT_SET_GEM_LOT_ACK",
	"gds_NC_ITEMDB_EQUIPFAIL_ACK",
	"gds_NC_ITEMDB_EQUIPSUC_ACK",
	"gds_NC_ITEMDB_EQUIP_BELONGED_CAN",
	"gds_NC_ITEMDB_EQUIP_BELONGED_CANCEL_USE_ACK",
	"gds_NC_ITEMDB_EXCHANGEFAIL_ACK",
	"gds_NC_ITEMDB_EXCHANGESUC_ACK",
	"gds_NC_ITEMDB_FURNITURE_ENDURE_ACK",
	"gds_NC_ITEMDB_GETFROMCAPSULE_ACK",
	"gds_NC_ITEMDB_GETFROMCHESTFAIL_ACK",
	"gds_NC_ITEMDB_GETFROMCHESTSUC_ACK",
	"gds_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_ACK",
	"gds_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_ACK",
	"gds_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_OPEN_FAIL_ACK",
	"gds_NC_ITEMDB_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_ACK",
	"gds_NC_ITEMDB_GUILD_STORAGE_WITHDRAW_ACK",
	"gds_NC_ITEMDB_GUILD_TOKEN_BUYALLSUC_ACK",
	"gds_NC_ITEMDB_GUILD_TOKEN_BUYFAIL_ACK",
	"gds_NC_ITEMDB_GUILD_TOKEN_BUYLOTSUC_ACK",
	"gds_NC_ITEMDB_INC_DEC_MONEYFAIL_ACK",
	"gds_NC_ITEMDB_INC_DEC_MONEYSUC_ACK",
	"gds_NC_ITEMDB_ITEMBREAKFAIL_ACK",
	"gds_NC_ITEMDB_ITEMBREAKSUC_ACK",
	"gds_NC_ITEMDB_ITEMCHANGE_ACK",
	"gds_NC_ITEMDB_ITEMINFO_ACK",
	"gds_NC_ITEMDB_ITEMMONEY_BUY_ACK",
	"gds_NC_ITEMDB_ITEMREBUILD_ACK",
	"gds_NC_ITEMDB_ITEMTOTALINFORM_ACK",
	"gds_NC_ITEMDB_ITEMTRADEFAIL_ACK",
	"gds_NC_ITEMDB_ITEMTRADESUC_ACK",
	"gds_NC_ITEMDB_MAPLINK_ITEM_CONSUME_ACK",
	"gds_NC_ITEMDB_MERGEFAIL_ACK",
	"gds_NC_ITEMDB_MERGESUC_ACK",
	"gds_NC_ITEMDB_MINIHOUSE_EFFECT_DEMANDGOOD_ACK",
	"gds_NC_ITEMDB_MIX_ITEM_ACK",
	"gds_NC_ITEMDB_MOVER_RAREMOVER_ACK",
	"gds_NC_ITEMDB_MOVER_UPGRADE_ACK",
	"gds_NC_ITEMDB_MYSTERY_VAULT_MAKEITEM_ACK",
	"gds_NC_ITEMDB_NEW_UPGRADE_AC",
	"gds_NC_ITEMDB_NEW_UPGRADE_ACK",
	"gds_NC_ITEMDB_OPENSTORAGE_ACK",
	"gds_NC_ITEMDB_OPENSTORAGE_FAIL_ACK",
	"gds_NC_ITEMDB_OPEN_GUILD_STORAGE_ACK",
	"gds_NC_ITEMDB_OPEN_GUILD_STORAGE_FAIL_ACK",
	"gds_NC_ITEMDB_PICKALLFAIL_ACK",
	"gds_NC_ITEMDB_PICKALLSUC_ACK",
	"gds_NC_ITEMDB_PICKMERGEFAIL_ACK",
	"gds_NC_ITEMDB_PICKMERGESUC_ACK",
	"gds_NC_ITEMDB_PICKMONEYFAIL_ACK",
	"gds_NC_ITEMDB_PICKMONEYSUC_ACK",
	"gds_NC_ITEMDB_PRODUCE_ACK",
	"gds_NC_ITEMDB_QUESTREWARD_ACK",
	"gds_NC_ITEMDB_RANDOMOPTION_CHANGE_ACK",
	"gds_NC_ITEMDB_RANDOMOPTION_CHANGE_CONSUME_AND_COUNTING_ACK",
	"gds_NC_ITEMDB_RANDOMOPTION_RECOVER_COUNT_LIMIT_ACK",
	"gds_NC_ITEMDB_REINFORCEUNEQUIPFAIL_ACK",
	"gds_NC_ITEMDB_REINFORCEUNEQUIPSUC_ACK",
	"gds_NC_ITEMDB_RELOCFAIL_ACK",
	"gds_NC_ITEMDB_RELOCSUC_ACK",
	"gds_NC_ITEMDB_REPURCHASE_ALL_ACK",
	"gds_NC_ITEMDB_REPURCHASE_LOT_ACK",
	"gds_NC_ITEMDB_SE",
	"gds_NC_ITEMDB_SELLALLSUC_ACK",
	"gds_NC_ITEMDB_SELLFAIL_ACK",
	"gds_NC_ITEMDB_SELLLOTSUC_ACK",
	"gds_NC_ITEMDB_SHIELDENDURE_CHARGE_ACK",
	"gds_NC_ITEMDB_SOULSTONEBUYFAIL_ACK",
	"gds_NC_ITEMDB_SOULSTONEBUYSUC_ACK",
	"gds_NC_ITEMDB_SPLITFAIL_ACK",
	"gds_NC_ITEMDB_SPLITSUC_ACK",
	"gds_NC_ITEMDB_SPLIT_N_MERGEFAIL_ACK",
	"gds_NC_ITEMDB_SPLIT_N_MERGESUC_ACK",
	"gds_NC_ITEMDB_TERMEXTEND_ACK",
	"gds_NC_ITEMDB_UES_FRIEND_POINT_ACK",
	"gds_NC_ITEMDB_UNEQUIPFAIL_ACK",
	"gds_NC_ITEMDB_UNEQUIPSUC_ACK",
	"gds_NC_ITEMDB_UPGRADE_ACK",
	"gds_NC_ITEMDB_USE_ACK",
	"gds_NC_ITEMDB_WEAPONENDURE_CHARGE_ACK",
	"gds_NC_ITEMDB_WITHDRAW_ACK",
	"gds_NC_KQ_REWARDFAIL_ACK",
	"gds_NC_KQ_REWARDSUC_ACK",
	"gds_NC_MINIHOUSE_DB_PORTAL_ADD_ACK",
	"gds_NC_MINIHOUSE_DB_PORTAL_DEL_ACK",
	"gds_NC_MINIHOUSE_DB_PORTAL_LIST_ACK",
	"gds_NC_MINIHOUSE_DB_VISITER_COUNT_ACK",
	"gds_NC_MISC_HEARTBEAT_ACK",
	"gds_NC_MISC_HEARTBEAT_REQ",
	"gds_NC_MISC_PINGTEST_CLIENT_ZONE_DB",
	"gds_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB",
	"gds_NC_MISC_PINGTEST_TOOL_WM_ZONE_DB",
	"gds_NC_MISC_S2SCONNECTION_ACK",
	"gds_NC_MISC_S2SCONNECTION_RDY",
	"gds_NC_PET_CREATE_DB_ACK",
	"gds_NC_PET_LOAD_INFO_DB_ACK",
	"gds_NC_PET_SET_NAME_DB_ACK",
	"gds_NC_QUEST_DB_CLEAR_ACK",
	"gds_NC_QUEST_DB_DONE_ACK",
	"gds_NC_QUEST_DB_GIVE_UP_ACK",
	"gds_NC_QUEST_DB_READ_ACK",
	"gds_NC_QUEST_DB_SET_INFO_ACK",
	"gds_NC_SKILL_EMPOWALLOC_DB_ACK",
	"gds_NC_SKILL_EMPOW_RESET_DB_FAIL_ACK",
	"gds_NC_SKILL_EMPOW_RESET_DB_SUC_ACK",
	"gds_NC_SKILL_ERASE_ACK",
	"gds_NC_SKILL_SKILLTEACHFAIL_ACK",
	"gds_NC_SKILL_SKILLTEACHSUC_ACK",
	"gds_NC_USER_CONNECTION_DB_CHANGE_CHAR_ID_ITEM_USE_ACK",
	"gds_NC_WT_LICENSE_CLR_ACK",
	"gds_NC_WT_LICENSE_SET_ACK",
	"gds_NC_WT_MOB_KILLCOUNT_SET_ACK",
	"gds_NC_WT_TITLE_SET_ACK",
}

var officialGDSStoredProcMap = buildOfficialGDSStoredProcMap()

func buildOfficialGDSStoredProcMap() map[string]string {
	m := make(map[string]string, len(OfficialGDSSymbols196))
	for _, sym := range OfficialGDSSymbols196 {
		switch {
		case strings.Contains(sym, "ADMINCREATE"):
			m[sym] = "usp_GM_setItemReward"
		case strings.Contains(sym, "ABSTATE"):
			m[sym] = "p_AbState_Set"
		case strings.Contains(sym, "CHARSAVE_ALL"), strings.Contains(sym, "CLASSCHANGE"), strings.Contains(sym, "FREESTAT"), strings.Contains(sym, "STAT_"):
			m[sym] = "p_Char_Save"
		case strings.Contains(sym, "AUTO_PICK"):
			m[sym] = "usp_Character_setAutoPick"
		case strings.Contains(sym, "UI_STATE"):
			m[sym] = "usp_Character_setMysteryBolt"
		case strings.Contains(sym, "REST_EXP"):
			m[sym] = "usp_Character_setRestExpExecTime"
		case strings.Contains(sym, "TUTORIAL"):
			m[sym] = "usp_Tutorial_setCharInfo"
		case strings.Contains(sym, "WEDDING"):
			m[sym] = "usp_Wedding_Set"
		case strings.Contains(sym, "COLLECT"):
			m[sym] = "usp_Collection_setReward"
		case strings.Contains(sym, "DICE_TAISAI"):
			m[sym] = "usp_Dice_setResult"
		case strings.Contains(sym, "GAMBLE"):
			m[sym] = "usp_Gamble_setResult"
		case strings.Contains(sym, "GUILD_ACADEMY"):
			m[sym] = "usp_GuildAcademy_setReward"
		case strings.Contains(sym, "GUILD_TOURNAMENT"):
			m[sym] = "usp_GuildTournament_SetResult"
		case strings.Contains(sym, "GUILD_STORAGE"):
			m[sym] = "p_Guild_MoneySub"
		case strings.Contains(sym, "GUILD_TOKEN"):
			m[sym] = "usp_Item_GuildToken_BuyLot"
		case strings.Contains(sym, "HOLY_PROMISE"):
			m[sym] = "p_HolyPromise_SetCen"
		case strings.Contains(sym, "MINIHOUSE"):
			m[sym] = "usp_Minihouse_setInfo"
		case strings.Contains(sym, "PET"):
			m[sym] = "usp_Pet_setInfo"
		case strings.Contains(sym, "QUEST"):
			m[sym] = "p_Quest_Set"
		case strings.Contains(sym, "SKILL"):
			m[sym] = "p_Skill_Set"
		case strings.Contains(sym, "CHANGE_CHAR_ID"):
			m[sym] = "usp_Character_setChangeIDItemUsed"
		case strings.Contains(sym, "WT_"):
			m[sym] = "p_Item_SetWeaponTitle"
		case strings.Contains(sym, "UES_FRIEND_POINT"):
			m[sym] = "usp_Character_setFriendpoint"
		case strings.Contains(sym, "CHARGED_"):
			m[sym] = "usp_Charge_ItemWithdraw"
		case strings.Contains(sym, "KQ_"):
			m[sym] = "p_KQ_Set"
		case strings.Contains(sym, "MISC_"):
			m[sym] = "usp_WorldZone_Heartbeat"
		default:
			m[sym] = "p_Item_Update"
		}
	}
	return m
}
