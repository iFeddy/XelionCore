package zone

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/item"
	"github.com/iFeddy/XelionCore/pkg/data/npc"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler processes incoming packets for the Zone Server.
type Handler struct {
	worldID           byte
	zoneNo            uint16
	multiNode         bool
	mu                sync.RWMutex
	charRepo          database.CharacterRepository
	accountRepo       *database.MySQLAccountRepository
	transfers         *database.TransferManager
	dataStore         *DataStore
	classStore        *classparam.Store
	fieldStore        *field.Store
	npcStore          *npc.Store
	mobInfoSHN        *shn.File
	chkEngine         *shn.ChecksumEngine
	itemStore         *item.Store
	players           map[*network.Client]*ZonePlayer
	mapManager        *MapManager
	gameLogDispatcher *database.GameLogDispatcher
	s2s               *S2SCoordinator
}

// SetGameLogDispatcher configura el despachador de logs de auditoría para ZoneService.
func (h *Handler) SetGameLogDispatcher(disp *database.GameLogDispatcher) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.gameLogDispatcher = disp
}

// GameLogDispatcher retorna el despachador de logs configurado.
func (h *Handler) GameLogDispatcher() *database.GameLogDispatcher {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.gameLogDispatcher
}

// SetMultiNode habilita o deshabilita el filtrado y enrutamiento multi-nodo Zone00..Zone04.
func (h *Handler) SetMultiNode(enabled bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.multiNode = enabled
}

// IsMultiNode indica si el Handler opera en topología multi-nodo Zone00..Zone04.
func (h *Handler) IsMultiNode() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.multiNode
}

// S2S retorna el coordinador S2S (138 wms_NC_* / 196 gds_NC_*) del Handler.
func (h *Handler) S2S() *S2SCoordinator {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.s2s
}

// MapManager retorna el administrador de mapas de Zone.
func (h *Handler) MapManager() *MapManager {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.mapManager
}

// NewHandler initializes a new Zone packet handler.
func NewHandler(charRepo database.CharacterRepository, transfers *database.TransferManager, classStore *classparam.Store, fieldStore *field.Store, chkEngine *shn.ChecksumEngine) *Handler {
	h := &Handler{
		charRepo:   charRepo,
		transfers:  transfers,
		classStore: classStore,
		fieldStore: fieldStore,
		chkEngine:  chkEngine,
		players:    make(map[*network.Client]*ZonePlayer),
		mapManager: NewMapManager(fieldStore),
		s2s:        NewS2SCoordinator(),
	}
	h.mapManager.SetHandler(h)
	return h
}

// HandlePacket routes packets to their specific handlers using official PROTO_NC_* opcodes.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	case protocol.PROTO_NC_MAP_LOGIN_REQ:
		h.handleTransferKey(client, pkt)
	case protocol.PROTO_NC_MAP_LOGINCOMPLETE_CMD, protocol.PROTO_NC_KQ_LIST_REFRESH_REQ:
		h.handleClientReady(client, pkt)
	case protocol.PROTO_NC_MISC_SEED_REQ:
		h.handleSeedReq(client, pkt)
	case protocol.PROTO_NC_MISC_HIDE_EXCEPT_ME_ON_CMD:
		h.handleHideExceptMeOn(client, pkt)
	case protocol.PROTO_NC_MISC_HIDE_EXCEPT_ME_OFF_CMD:
		h.handleHideExceptMeOff(client, pkt)
	case protocol.PROTO_NC_BRIEFINFO_INFORM_CMD:
		h.handleBriefInfoInform(client, pkt)
	case protocol.PROTO_NC_ACT_WALK_REQ:
		h.handleWalk(client, pkt)
	case protocol.PROTO_NC_ACT_RUN_REQ:
		h.handleRun(client, pkt)
	case protocol.PROTO_NC_ACT_MOVEWALK_CMD:
		h.handleMoveWalk(client, pkt)
	case protocol.PROTO_NC_ACT_MOVERUN_CMD:
		h.handleMoveRun(client, pkt)
	case protocol.PROTO_NC_ACT_STOP_REQ:
		h.handleStop(client, pkt)
	case protocol.PROTO_NC_ACT_JUMP_CMD:
		h.handleJump(client, pkt)
	case protocol.PROTO_NC_ACT_EMOTICON_CMD:
		h.handleEmoticon(client, pkt)
	case protocol.PROTO_NC_ACT_EMOTICONSTOP_CMD:
		h.handleEmoticonStop(client, pkt)
	case protocol.PROTO_NC_ACT_CHAT_REQ:
		h.handleChatNormal(client, pkt)
	case protocol.PROTO_NC_ACT_PITCHTENT_REQ:
		h.handleBeginRest(client, pkt)
	case protocol.PROTO_NC_ACT_FOLDTENT_REQ:
		h.handleEndRest(client, pkt)
	case protocol.PROTO_NC_ACT_GATHERSTART_REQ:
		h.handleGatherStart(client, pkt)
	case protocol.PROTO_NC_ACT_GATHERCANCEL_CMD:
		h.handleGatherCancel(client, pkt)
	case protocol.PROTO_NC_ACT_GATHERCOMPLETE_REQ:
		h.handleGatherComplete(client, pkt)
	case protocol.PROTO_NC_ACT_RIDE_FEEDING_REQ:
		h.handleRideFeeding(client, pkt)
	case protocol.PROTO_NC_MOVER_RIDE_ON_REQ:
		h.handleMoverRideOn(client, pkt)
	case protocol.PROTO_NC_MOVER_RIDE_OFF_REQ:
		h.handleMoverRideOff(client, pkt)
	case protocol.PROTO_NC_MAP_WING_SAVE_REQ:
		h.handleWingSave(client, pkt)
	case protocol.PROTO_NC_MAP_TONORMALCOORD_CMD:
		h.handleToNormalCoord(client, pkt)
	case protocol.PROTO_NC_ACT_ENDOFTRADE_CMD:
		h.handleEndOfTradeCmd(client, pkt)
	case protocol.PROTO_NC_ACT_PRODUCE_CAST_REQ:
		h.handleActProduceCastReq(client, pkt)
	case protocol.PROTO_NC_ACT_PRODUCE_CASTABORT_CMD:
		h.handleActProduceCastAbortCmd(client, pkt)
	case protocol.PROTO_NC_ACT_ACTIONBYITEM_REQ:
		h.handleActionByItemReq(client, pkt)
	case protocol.PROTO_NC_ACT_NPCCLICK_CMD:
		h.handleActNPCClickCmd(client, pkt)
	case protocol.PROTO_NC_ACT_NPCMENUOPEN_ACK:
		h.handleActNPCMenuOpenAck(client, pkt)
	case protocol.PROTO_NC_ACT_NPC_ACTION_CMD:
		h.handleNPCActionCmd(client, pkt)
	case protocol.PROTO_NC_ACT_NPC_MENU_CMD:
		h.handleActNPCMenuCmd(client, pkt)
	case protocol.PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ:
		h.handleUseGate(client, pkt)
	case protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_REQ:
		h.handleMultyLinkSelect(client, pkt)
	case protocol.PROTO_NC_MAP_TOWNPORTAL_REQ:
		h.handleTownPortal(client, pkt)
	case protocol.PROTO_NC_MAP_WING_FLY_REQ:
		h.handleWingFly(client, pkt)
	case protocol.PROTO_NC_ITEM_BUY_REQ:
		h.handleItemBuyReq(client, pkt)
	case protocol.PROTO_NC_ITEM_SELL_REQ:
		h.handleItemSellReq(client, pkt)
	case protocol.PROTO_NC_ITEM_REPURCHASE_REQ:
		h.handleItemRepurchaseReq(client, pkt)
	case protocol.PROTO_NC_ITEM_UPGRADE_REQ:
		h.handleItemUpgradeReq(client, pkt)
	case protocol.PROTO_NC_ITEM_NEW_UPGRADE_REQ:
		h.handleItemNewUpgradeReq(client, pkt)
	case protocol.PROTO_NC_ITEM_BRACELET_UPGRADE_REQ:
		h.handleItemBraceletUpgradeReq(client, pkt)
	case protocol.PROTO_NC_ITEM_MOVER_UPGRADE_REQ:
		h.handleItemMoverUpgradeReq(client, pkt)
	case protocol.PROTO_NC_ITEM_USEABORT_CMD:
		h.handleItemUseAbortCmd(client, pkt)
	case protocol.PROTO_NC_ITEM_DEPOSIT_REQ:
		h.handleItemDepositReq(client, pkt)
	case protocol.PROTO_NC_ITEM_WITHDRAW_REQ:
		h.handleItemWithdrawReq(client, pkt)
	case protocol.PROTO_NC_ITEM_GUILD_STORAGE_WITHDRAW_REQ:
		h.handleItemGuildStorageWithdrawReq(client, pkt)
	case protocol.PROTO_NC_ITEM_OPENSTORAGEPAGE_REQ:
		h.handleItemOpenStoragePageReq(client, pkt)
	case protocol.PROTO_NC_ITEM_ACCOUNT_STORAGE_OPEN_CMD:
		h.handleAccountStorageOpenCmd(client, pkt)
	case protocol.PROTO_NC_ITEM_ACCOUNT_STORAGE_CLOSE_CMD:
		h.handleAccountStorageCloseCmd(client, pkt)
	case protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_REQ:
		h.handleItemAcademyStorageWithdrawReq(client, pkt)
	case protocol.PROTO_NC_ITEM_GUILD_ACADEMY_REWARD_STORAGE_DEPOSIT_REQ:
		h.handleItemAcademyStorageDepositReq(client, pkt)
	case protocol.PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_OPEN_REQ:
		h.handleGuildAcademyRewardStorageOpenReq(client, pkt)
	case protocol.PROTO_NC_GUILD_ACADEMY_REWARD_STORAGE_CLOSE_CMD:
		h.handleGuildAcademyRewardStorageCloseCmd(client, pkt)
	case protocol.PROTO_NC_ITEM_DISMANTLE_REQ:
		h.handleItemDismantleReq(client, pkt)
	case protocol.PROTO_NC_ITEM_MIX_ITEM_REQ:
		h.handleItemMixItemReq(client, pkt)
	case protocol.PROTO_NC_ITEM_MH_FURNITURE_ENDURE_KIT_USE_REQ:
		h.handleMHFurnitureEndureKitUseReq(client, pkt)
	case protocol.PROTO_NC_ITEM_WEAPONENDURE_CHARGE_REQ:
		h.handleWeaponEndureChargeReq(client, pkt)
	case protocol.PROTO_NC_ITEM_SHIELDENDURE_CHARGE_REQ:
		h.handleShieldEndureChargeReq(client, pkt)
	case protocol.PROTO_NC_ITEM_DICE_GAME_START_REQ:
		h.handleDiceGameStartReq(client, pkt)
	case protocol.PROTO_NC_ITEM_ENCHANT_ADD_GEM_REQ:
		h.handleItemEnchantAddGemReq(client, pkt)
	case protocol.PROTO_NC_ITEM_ENCHANT_REMOVE_GEM_REQ:
		h.handleItemEnchantRemoveGemReq(client, pkt)
	case protocol.PROTO_NC_ITEM_ENCHANT_ADD_NEW_SOCKET_REQ:
		h.handleItemEnchantAddNewSocketReq(client, pkt)
	case protocol.PROTO_NC_ITEM_EQUIP_BELONGED_CANCEL_USE_REQ:
		h.handleEquipBelongedCancelUseReq(client, pkt)
	case protocol.PROTO_NC_ITEM_MINIMON_EQUIP_REQ:
		h.handleMinimonEquipReq(client, pkt)
	case protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_REQ:
		h.handleRandomOptionChangeReq(client, pkt)
	case protocol.PROTO_NC_ITEM_RANDOMOPTION_CHANGE_ACCEPT_REQ:
		h.handleRandomOptionChangeAcceptReq(client, pkt)
	case protocol.PROTO_NC_ITEM_RANDOMOPTION_RECOVER_COUNT_LIMIT_REQ:
		h.handleRandomOptionRecoverCountLimitReq(client, pkt)
	case protocol.PROTO_NC_ITEM_CHAT_COLOR_CHANGE_REQ:
		h.handleChatColorChangeReq(client, pkt)
	case protocol.PROTO_NC_ITEM_TERMEXTEND_REQ:
		h.handleTermExtendReq(client, pkt)
	case protocol.PROTO_NC_ITEM_CLASS_CHANGE_REQ:
		h.handleClassChangeReq(client, pkt)
	case protocol.PROTO_NC_ITEM_MAPLINK_SCROLL_REQ:
		h.handleMapLinkScrollReq(client, pkt)
	case protocol.PROTO_NC_ITEM_REWARDINVENOPEN_REQ:
		h.handleRewardInvenOpen(client, pkt)
	case protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ:
		h.handleChargedInvenOpen(client, pkt)
	case protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_REQ:
		h.handleChargedWithdraw(client, pkt)
	case protocol.PROTO_NC_ITEM_RELOC_REQ:
		h.handleItemReloc(client, pkt)
	case protocol.PROTO_NC_ITEM_SPLIT_REQ:
		h.handleItemSplit(client, pkt)
	case protocol.PROTO_NC_ITEM_EQUIP_REQ:
		h.handleItemEquip(client, pkt)
	case protocol.PROTO_NC_ITEM_RINGEQUIP_REQ:
		h.handleItemRingEquip(client, pkt)
	case protocol.PROTO_NC_ITEM_UNEQUIP_REQ:
		h.handleItemUnequip(client, pkt)
	case protocol.PROTO_NC_ITEM_USE_REQ:
		h.handleItemUse(client, pkt)
	case protocol.PROTO_NC_ITEM_DROP_REQ:
		h.handleItemDrop(client, pkt)
	case protocol.PROTO_NC_ITEM_PICK_REQ:
		h.handleItemPick(client, pkt)
	case protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_REQ:
		h.handleAutoArrangeInven(client, pkt)
	case protocol.PROTO_NC_TRADE_PROPOSE_REQ:
		h.handleTradeProposeReq(client, pkt)
	case protocol.PROTO_NC_TRADE_PROPOSE_ASKNO_ACK:
		h.handleTradeProposeAskNoAck(client, pkt)
	case protocol.PROTO_NC_TRADE_PROPOSEYES_ACK:
		h.handleTradeProposeYesAck(client, pkt)
	case protocol.PROTO_NC_TRADE_PROPOSE_CANCEL_CMD:
		h.handleTradeProposeCancelCmd(client, pkt)
	case protocol.PROTO_NC_TRADE_CANCEL_REQ:
		h.handleTradeCancelReq(client, pkt)
	case protocol.PROTO_NC_TRADE_UPBOARD_REQ:
		h.handleTradeUpboardReq(client, pkt)
	case protocol.PROTO_NC_TRADE_DOWNBOARD_REQ:
		h.handleTradeDownboardReq(client, pkt)
	case protocol.PROTO_NC_TRADE_CENBOARDING_REQ:
		h.handleTradeCenBoardingReq(client, pkt)
	case protocol.PROTO_NC_TRADE_BOARDLOCK_REQ:
		h.handleTradeBoardLockReq(client, pkt)
	case protocol.PROTO_NC_TRADE_DECIDE_REQ:
		h.handleTradeDecideReq(client, pkt)
	case protocol.PROTO_NC_BOOTH_OPEN_REQ:
		h.handleBoothOpenReq(client, pkt)
	case protocol.PROTO_NC_BOOTH_CLOSE_REQ:
		h.handleBoothCloseReq(client, pkt)
	case protocol.PROTO_NC_BOOTH_ENTRY_REQ:
		h.handleBoothEntryReq(client, pkt)
	case protocol.PROTO_NC_BOOTH_REFRESH_REQ:
		h.handleBoothRefreshReq(client, pkt)
	case protocol.PROTO_NC_BOOTH_ITEMTRADE_REQ:
		h.handleBoothItemTradeReq(client, pkt)
	case protocol.PROTO_NC_BOOTH_INTERIORSTART_REQ:
		h.handleBoothInteriorStartReq(client, pkt)
	case protocol.PROTO_NC_BOOTH_SEARCH_ITEM_LIST_CATEGORIZED_REQ:
		h.handleBoothSearchItemListCategorizedReq(client, pkt)
	case protocol.PROTO_NC_BOOTH_SEARCH_BOOTH_POSITION_REQ:
		h.handleBoothSearchBoothPositionReq(client, pkt)
	case protocol.PROTO_NC_SOULSTONE_HP_BUY_REQ:
		h.handleSoulStoneHPBuy(client, pkt)
	case protocol.PROTO_NC_SOULSTONE_SP_BUY_REQ:
		h.handleSoulStoneSPBuy(client, pkt)
	case protocol.PROTO_NC_SOULSTONE_HP_USE_REQ:
		h.handleSoulStoneHPUse(client, pkt)
	case protocol.PROTO_NC_SOULSTONE_SP_USE_REQ:
		h.handleSoulStoneSPUse(client, pkt)
	case protocol.PROTO_NC_BAT_TARGETTING_REQ:
		h.handleTargettingReq(client, pkt)
	case protocol.PROTO_NC_BAT_UNTARGET_REQ:
		h.handleUntargetReq(client, pkt)
	case protocol.PROTO_NC_BAT_BASHSTART_CMD:
		h.handleBashStartCmd(client, pkt)
	case protocol.PROTO_NC_BAT_BASHSTOP_CMD:
		h.handleBashStopCmd(client, pkt)
	case protocol.PROTO_NC_BAT_ASSIST_REQ:
		h.handleAssistReq(client, pkt)
	case protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ:
		h.handleSkillBashObjCastReq(client, pkt)
	case protocol.PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ:
		h.handleSkillBashFldCastReq(client, pkt)
	case protocol.PROTO_NC_BAT_SKILLBASH_CASTABORT_REQ:
		h.handleSkillBashCastAbortReq(client, pkt)
	case protocol.PROTO_NC_BAT_ABSTATE_ERASE_REQ:
		h.handleAbStateEraseReq(client, pkt)
	case protocol.PROTO_NC_MOVER_SKILLBASH_OBJ_CAST_REQ:
		h.handleMoverSkillBashObjCastReq(client, pkt)
	case protocol.PROTO_NC_MOVER_SKILLBASH_FLD_CAST_REQ:
		h.handleMoverSkillBashFldCastReq(client, pkt)
	case protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ:
		h.handleSkillEmpowAllocReq(client, pkt)
	case protocol.PROTO_NC_SKILL_PRODUCTFIELD_REQ:
		h.handleSkillProductFieldReq(client, pkt)
	case protocol.PROTO_NC_SKILL_UNLEARN_REQ:
		h.handleSkillUnlearnReq(client, pkt)
	case protocol.PROTO_NC_SKILL_REPLYREVIVE_CMD:
		h.handleSkillReplyReviveCmd(client, pkt)
	case protocol.PROTO_NC_ITEM_USE_ACTIVESKILL_REQ:
		h.handleItemUseActiveSkillReq(client, pkt)
	case protocol.PROTO_NC_CHAR_STAT_INCPOINT_REQ:
		h.handleCharStatIncPointReq(client, pkt)
	case protocol.PROTO_NC_CHAR_STAT_DECPOINT_REQ:
		h.handleCharStatDecPointReq(client, pkt)
	case protocol.PROTO_NC_CHAR_REVIVE_REQ:
		h.handleCharReviveReq(client, pkt)
	case protocol.PROTO_NC_CHAR_REBIRTH_CMD:
		h.handleCharRebirthCmd(client, pkt)
	case protocol.PROTO_NC_CHAR_REBIRTH_REJECT_CMD:
		h.handleCharRebirthRejectCmd(client, pkt)
	case protocol.PROTO_NC_ITEM_REVIVEITEMUSE_CMD:
		h.handleItemReviveItemUseCmd(client, pkt)
	case protocol.PROTO_NC_CHAR_DEPOLYMORPH_CMD:
		h.handleDepolymorphCmd(client, pkt)
	case protocol.PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ:
		h.handleChargedBuffEraseReq(client, pkt)
	case protocol.PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ:
		h.handleAutoPickReq(client, pkt)
	case protocol.PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_ON_REQ:
		h.handleMinimonNormalOnReq(client, pkt)
	case protocol.PROTO_NC_CHAR_USEITEM_MINIMON_NORMAL_ITEM_OFF_REQ:
		h.handleMinimonNormalOffReq(client, pkt)
	case protocol.PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_ON_REQ:
		h.handleMinimonChargedOnReq(client, pkt)
	case protocol.PROTO_NC_CHAR_USEITEM_MINIMON_CHARGED_ITEM_OFF_REQ:
		h.handleMinimonChargedOffReq(client, pkt)
	case protocol.PROTO_NC_CT_SET_CURRENT_REQ:
		h.handleCTSetCurrentReq(client, pkt)
	case protocol.PROTO_NC_COLLECT_PEEPING_REQ:
		h.handleCollectPeepingReq(client, pkt)
	case protocol.PROTO_NC_COLLECT_CARDREGIST_REQ:
		h.handleCollectCardRegistReq(client, pkt)
	case protocol.PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ:
		h.handleCollectBookmarkRegistReq(client, pkt)
	case protocol.PROTO_NC_COLLECT_REWARD_REQ:
		h.handleCollectRewardReq(client, pkt)

	// === Fase 15: Quests, NPCs/Menús, Escenarios, Tutorial, Instancias Lua (ID/KQ/GT/Wedding/HolyPromise), MiniHouse, Pets, WT y Casino ===
	case protocol.PROTO_NC_QUEST_SCRIPT_CMD_ACK:
		h.handleQuestScriptCmdAck(client, pkt)
	case protocol.PROTO_NC_QUEST_GIVE_UP_REQ:
		h.handleQuestGiveUpReq(client, pkt)
	case protocol.PROTO_NC_QUEST_CLIENT_SCENARIO_DONE_REQ:
		h.handleQuestClientScenarioDoneReq(client, pkt)
	case protocol.PROTO_NC_QUEST_SELECT_START_REQ:
		h.handleQuestSelectStartReq(client, pkt)
	case protocol.PROTO_NC_QUEST_REWARD_SELECT_ITEM_INDEX_CMD:
		h.handleQuestRewardSelectItemIndexCmd(client, pkt)
	case protocol.PROTO_NC_QUEST_START_REQ:
		h.handleQuestStartReq(client, pkt)
	case protocol.PROTO_NC_QUEST_READ_REQ:
		h.handleQuestReadReq(client, pkt)
	case protocol.PROTO_NC_MENU_SERVERMENU_ACK:
		h.handleMenuServerMenuAck(client, pkt)
	case protocol.PROTO_NC_SCENARIO_AREAENTRY_ACK:
		h.handleScenarioAreaEntryAck(client, pkt)
	case protocol.PROTO_NC_SCENARIO_AREALEAVE_ACK:
		h.handleScenarioAreaLeaveAck(client, pkt)
	case protocol.PROTO_NC_CHAR_TUTORIAL_STEP_ACK:
		h.handleCharTutorialStepAck(client, pkt)
	case protocol.PROTO_NC_INSTANCE_DUNGEON_LEVEL_SELECT_MENU_ACK:
		h.handleInstanceDungeonLevelSelectMenuAck(client, pkt)
	case protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_LIST_REQ:
		h.handleInstanceDungeonRankListReq(client, pkt)
	case protocol.PROTO_NC_KQ_ENTRYRESPONCE_ACK:
		h.handleKQEntryResponseAck(client, pkt)
	case protocol.PROTO_NC_GUILD_TOURNAMENT_RECALL_ENTER_REQ:
		h.handleGuildTournamentRecallEnterReq(client, pkt)
	case protocol.PROTO_NC_GUILD_TOURNAMENT_FLAGCAPTURE_REQ:
		h.handleGuildTournamentFlagCaptureReq(client, pkt)
	case protocol.PROTO_NC_GUILD_TOURNAMENT_USESKILL_REQ:
		h.handleGuildTournamentUseSkillReq(client, pkt)
	case protocol.PROTO_NC_GUILD_TOURNAMENT_DICEGAME_THROW_REQ:
		h.handleGuildTournamentDiceGameThrowReq(client, pkt)
	case protocol.PROTO_NC_GUILD_TOURNAMENT_OBSERVER_ENTER_REQ:
		h.handleGuildTournamentObserverEnterReq(client, pkt)
	case protocol.PROTO_NC_GUILD_TOURNAMENT_OBSERVER_OUT_REQ:
		h.handleGuildTournamentObserverOutReq(client, pkt)
	case protocol.PROTO_NC_ACT_WEDDING_PROPOSEREQ_REQ:
		h.handleActWeddingProposeReq(client, pkt)
	case protocol.PROTO_NC_ACT_WEDDING_PROPOSEACK_ACK:
		h.handleActWeddingProposeAckAck(client, pkt)
	case protocol.PROTO_NC_ACT_WEDDING_HALL_RESERV_REQ:
		h.handleActWeddingHallReservReq(client, pkt)
	case protocol.PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_READY_REQ:
		h.handleActWeddingHallGuestEnterReadyReq(client, pkt)
	case protocol.PROTO_NC_ACT_WEDDING_HALL_GUEST_ENTER_REQ:
		h.handleActWeddingHallGuestEnterReq(client, pkt)
	case protocol.PROTO_NC_ACT_WEDDING_AGREEMENT_DIVORCE_REQ:
		h.handleActWeddingAgreementDivorceReq(client, pkt)
	case protocol.PROTO_NC_ACT_WEDDING_COMPULSORY_DIVORCE_REQ:
		h.handleActWeddingCompulsoryDivorceReq(client, pkt)
	case protocol.PROTO_NC_ACT_WEDDING_DIVORCE_REFUSE_CMD:
		h.handleActWeddingDivorceRefuseCmd(client, pkt)
	case protocol.PROTO_NC_CHAR_WEDDING_PARTNER_INFO_REQ:
		h.handleCharWeddingPartnerInfoReq(client, pkt)
	case protocol.PROTO_NC_HOLY_PROMISE_GET_CEN_REWARD_REQ:
		h.handleHolyPromiseGetCenRewardReq(client, pkt)
	case protocol.PROTO_NC_HOLY_PROMISE_WITHDRAW_CEN_REWARD_REQ:
		h.handleHolyPromiseWithdrawCenRewardReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_ACTIV_REQ:
		h.handleMiniHouseActivReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_BUILDING_REQ:
		h.handleMiniHouseBuildingReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_LOGINCOMPLETE_REQ:
		h.handleMiniHouseLoginCompleteReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_VISITREADY_REQ:
		h.handleMiniHouseVisitReadyReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_VISIT_REQ:
		h.handleMiniHouseVisitReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_KICKOUT_REQ:
		h.handleMiniHouseKickoutReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_KICKOUTCANCEL_REQ:
		h.handleMiniHouseKickoutCancelReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_EXIT_REQ:
		h.handleMiniHouseExitReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_ARRANGEMODE_REQ:
		h.handleMiniHouseArrangeModeReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_FUNITUREINFOCOMPLETE_REQ:
		h.handleMiniHouseFurnitureInfoCompleteReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_REARRANGE_REQ:
		h.handleMiniHouseRearrangeReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_CREATE_FURNITURE_REQ:
		h.handleMiniHouseCreateFurnitureReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_DELETE_FURNITURE_REQ:
		h.handleMiniHouseDeleteFurnitureReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_COMPULSIONMOVETO_REQ:
		h.handleMiniHouseCompulsionMoveToReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_MODIFY_PASSWORD_REQ:
		h.handleMiniHouseModifyPasswordReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_MODIFY_TITLE_REQ:
		h.handleMiniHouseModifyTitleReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_MODIFY_MAXENTERNUM_REQ:
		h.handleMiniHouseModifyMaxEnterNumReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_MODIFY_OWNERBLOG_REQ:
		h.handleMiniHouseModifyOwnerBlogReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_OWNERBLOG_REQ:
		h.handleMiniHouseOwnerBlogReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_CHAR_ACTION_REQ:
		h.handleMiniHouseCharActionReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_FURNITURE_EFFECT_REQ:
		h.handleMiniHouseFurnitureEffectReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_MODIFY_ITEM_INFO_OPEN_REQ:
		h.handleMiniHouseModifyItemInfoOpenReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_MODIFY_NOTIFY_REQ:
		h.handleMiniHouseModifyNotifyReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_PORTAL_DEL_REQ:
		h.handleMiniHousePortalDelReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_PORTAL_OPEN_REQ:
		h.handleMiniHousePortalOpenReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_PORTAL_CLOSE_REQ:
		h.handleMiniHousePortalCloseReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_PORTAL_EFFECT_REQ:
		h.handleMiniHousePortalEffectReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_REQ:
		h.handleMiniHouseFurnitureEmotionReq(client, pkt)
	case protocol.PROTO_NC_MINIHOUSE_FURNITURE_EMOTION_CANCEL_REQ:
		h.handleMiniHouseFurnitureEmotionCancelReq(client, pkt)
	case protocol.PROTO_NC_PET_ASK_NEW_NAME_REQ:
		h.handlePetAskNewNameReq(client, pkt)
	case protocol.PROTO_NC_PET_ASK_NEW_NAME_ACK:
		h.handlePetAskNewNameAck(client, pkt)
	case protocol.PROTO_NC_PET_SET_NAME_REQ:
		h.handlePetSetNameReq(client, pkt)
	case protocol.PROTO_NC_PET_SET_NAME_CANCEL_REQ:
		h.handlePetSetNameCancelReq(client, pkt)
	case protocol.PROTO_NC_WT_LICENSE_REQ:
		h.handleWTLicenseReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_GAME_JOIN_REQ:
		h.handleDiceTaisaiGameJoinReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_GAME_LEAVE_REQ:
		h.handleDiceTaisaiGameLeaveReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_EXPECT_INCOME_MONEY_REQ:
		h.handleDiceTaisaiExpectIncomeMoneyReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_RECEIPT_INCOME_MONEY_REQ:
		h.handleDiceTaisaiReceiptIncomeMoneyReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_GAME_MODE_CHANGE_REQ:
		h.handleDiceTaisaiGameModeChangeReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_RANK_REQ:
		h.handleDiceTaisaiRankReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_BETTING_REQ:
		h.handleDiceTaisaiBettingReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_BETTING_CANCEL_REQ:
		h.handleDiceTaisaiBettingCancelReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_DICE_ROLL_REQ:
		h.handleDiceTaisaiDiceRollReq(client, pkt)
	case protocol.PROTO_NC_DICE_TAISAI_LARGE_AMOUNT_REQ:
		h.handleDiceTaisaiLargeAmountReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_GAMBLEHOUSE_EXIT_REQ:
		h.handleGambleHouseExitReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_COIN_VIPCARD_UI_OPEN_REQ:
		h.handleGambleCoinVipCardUIOpenReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_COIN_BUY_REQ:
		h.handleGambleCoinBuyReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_COIN_SELL_REQ:
		h.handleGambleCoinSellReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_TYPE_AND_WHERE_STAND_REQ:
		h.handleGambleTypeAndWhereStandReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_DICE_TAISAI_GAME_JOIN_REQ:
		h.handleGambleDiceTaisaiGameJoinReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_DICE_TAISAI_GAME_LEAVE_REQ:
		h.handleGambleDiceTaisaiGameLeaveReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_DICE_TAISAI_ALL_RANK_REQ:
		h.handleGambleDiceTaisaiAllRankReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_DICE_TAISAI_CURR_RANK_REQ:
		h.handleGambleDiceTaisaiCurrRankReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_REQ:
		h.handleGambleDiceTaisaiBettingReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_DICE_TAISAI_BETTING_CANCEL_REQ:
		h.handleGambleDiceTaisaiBettingCancelReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_DICE_TAISAI_LARGE_AMOUNT_REQ:
		h.handleGambleDiceTaisaiLargeAmountReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_SLOTMACHINE_GAME_JOIN_REQ:
		h.handleGambleSlotMachineJoinReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_SLOTMACHINE_GAME_LEAVE_REQ:
		h.handleGambleSlotMachineLeaveReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_SLOTMACHINE_START_REQ:
		h.handleGambleSlotMachineStartReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_SLOTMACHINE_STOPBUTTON_REQ:
		h.handleGambleSlotMachineStopButtonReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WHEELSTOP_REQ:
		h.handleGambleSlotMachineWheelStopReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_SLOTMACHINE_JACKPOTRANKING_REQ:
		h.handleGambleSlotMachineJackpotRankingReq(client, pkt)
	case protocol.PROTO_NC_GAMBLE_SLOTMACHINE_WINRANKING_REQ:
		h.handleGambleSlotMachineWinRankingReq(client, pkt)
	// Fase 16 — Cierre 245/245 OpCodes C2S oficiales (protocolstore ShinePlayer)
	case protocol.PROTO_NC_MISC_PINGTEST_CLIENT_ZONE_DB:
		h.handleMiscPingTestClientZoneDB(client, pkt)
	case protocol.PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB:
		h.handleMiscPingTestToolWMClientZoneDB(client, pkt)
	case protocol.PROTO_NC_MISC_SERVERPARAMETER_REQ:
		h.handleMiscServerParameterReq(client, pkt)
	case protocol.PROTO_NC_USER_NORMALLOGOUT_CMD:
		h.handleUserNormalLogoutCmd(client, pkt)
	case protocol.PROTO_NC_USER_REGISENUMBER_REQ:
		h.handleUserRegiseNumberReq(client, pkt)
	case protocol.PROTO_NC_USER_USE_BEAUTY_SHOP_CMD:
		h.handleUserUseBeautyShopCmd(client, pkt)
	case protocol.PROTO_NC_CHAR_LOGOUTREADY_CMD:
		h.handleCharLogoutReadyCmd(client, pkt)
	case protocol.PROTO_NC_CHAR_LOGOUTCANCEL_CMD:
		h.handleCharLogoutCancelCmd(client, pkt)
	case protocol.PROTO_NC_CHAR_SAVE_LINK_REQ:
		h.handleCharSaveLinkReq(client, pkt)
	case protocol.PROTO_NC_ACT_CHANGEMODE_REQ:
		h.handleActChangeModeReq(client, pkt)
	case protocol.PROTO_NC_ACT_SHOUT_CMD:
		h.handleActShoutCmd(client, pkt)
	case protocol.PROTO_NC_ACT_ROAR_REQ:
		h.handleActRoarReq(client, pkt)
	case protocol.PROTO_NC_PARTY_MEMBERINFOREQ_CMD:
		h.handlePartyMemberInfoReqCmd(client, pkt)
	case protocol.PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ:
		h.handleCharSaveUIStateSaveReq(client, pkt)
	case protocol.PROTO_NC_FRIEND_UES_FRIEND_POINT_REQ:
		h.handleFriendUseFriendPointReq(client, pkt)
	case protocol.PROTO_NC_GUILD_STORAGEOPEN_REQ:
		h.handleGuildStorageOpenReq(client, pkt)
	case protocol.PROTO_NC_ITEMDB_CLOSE_GUILD_STORAGE_CMD:
		h.handleItemDBCloseGuildStorageCmd(client, pkt)
	case protocol.PROTO_NC_GUILD_RENAME_REQ:
		h.handleGuildRenameReq(client, pkt)
	case protocol.PROTO_NC_GUILD_RETYPE_REQ:
		h.handleGuildRetypeReq(client, pkt)
	case protocol.PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_ITEM_USE_REQ:
		h.handleUserConnectionChangeCharIDItemUseReq(client, pkt)
	case protocol.PROTO_NC_MISC_HEARTBEAT_REQ, protocol.PROTO_NC_MISC_HEARTBEAT_ACK:
		h.handleFoodEatCmd(client, pkt)
	case protocol.PROTO_NC_MISC_GAMETIME_REQ:
		// Heartbeat pong received, ignore
	default:
		log.Printf("[Zone][%s] Unhandled OpCode: 0x%04X (Header=%d, Type=%d, Len=%d)",
			client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type, len(pkt.Raw))
	}
}

// OnClientDisconnect cleans up player resources on TCP disconnect.
func (h *Handler) OnClientDisconnect(client *network.Client) {
	h.mu.Lock()
	player, exists := h.players[client]
	delete(h.players, client)
	h.mu.Unlock()

	if exists && player != nil {
		h.cancelTradeSession(player)
		h.closePlayerBooth(player)
		if player.CurrentMap != nil {
			if h.gameLogDispatcher != nil && player.Char != nil {
				h.gameLogDispatcher.EnqueueType1(2, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0)
			}
			player.CurrentMap.RemovePlayer(player)
		}

		if player.Char != nil && h.charRepo != nil {
			mapID := player.Char.TownName
			if player.CurrentMap != nil && player.CurrentMap.MapID != "" {
				mapID = player.CurrentMap.MapID
			}

			playMin := int(time.Since(player.LoginTime).Minutes())
			if playMin < 1 {
				playMin = 1
			}

			// 1. CALL p_Char_Logout oficial
			_ = h.charRepo.CharacterLogout(player.Char.ID, playMin)

			// 2. Persistencia integral de ubicación y estadísticas
			_ = h.charRepo.SaveCharLocation(player.Char.ID, mapID, player.Char.XPos, player.Char.YPos, player.Rotation)
			_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)
			_ = h.charRepo.SaveCharLevelExp(player.Char.ID, int16(player.Char.Level), player.Char.TCharacterState.Exp)
			_ = h.charRepo.SaveCharMoney(player.Char.ID, player.Char.TCharacterState.Money)

			// 3. CALL usp_Char_SetRestExpExecTime oficial
			_ = h.charRepo.SetRestExpExecTime(player.Char.ID, time.Now())

			_ = h.charRepo.UpdateCharacter(player.Char)
		}

		// Persistir registro oficial de desconexión y tiempo jugado en tLogout (usp_Log_setLogout)
		if h.accountRepo != nil && player.Char != nil {
			playMin := int16(time.Since(player.LoginTime).Minutes())
			if playMin < 1 {
				playMin = 1
			}
			worldNo := player.WorldNo
			if worldNo == 0 && h.worldID != 0 {
				worldNo = h.worldID
			}
			_ = h.accountRepo.LogLogout(
				worldNo,
				player.Char.AccountID,
				player.Char.ID,
				player.Char.Job,
				player.Char.Level,
				player.Char.TCharacterState.Exp,
				playMin,
			)

			// Zone.exe nunca cierra tActiveAccountSession en Account DB; la sesión de cuenta la gobierna
			// WorldService/LoginService para permitir el retorno a Selección de Personaje y Selección de Mundo.
			log.Printf("[Zone][%s] Logout de personaje registrado en tLogout: '%s' (User: %d, Char: %d, World: %d, PlayTime: %d min, LogoutType: %d, Transfer: %v)",
				client.RemoteAddr(), player.Char.Name, player.Char.AccountID, player.Char.ID, worldNo, playMin, player.LogoutType, client.IsTransferring)
		}
	}
}

func (h *Handler) getPlayer(client *network.Client) (*ZonePlayer, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	p, ok := h.players[client]
	return p, ok
}

// GetPlayerByName busca un jugador activo por su nombre.
func (h *Handler) GetPlayerByName(name string) (*ZonePlayer, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, p := range h.players {
		if p != nil && p.Char != nil && p.Char.Name == name {
			return p, true
		}
	}
	return nil, false
}

func (h *Handler) handleTransferKey(client *network.Client, pkt *protocol.Packet) {
	req, err := protocol.Read_PROTO_NC_MAP_LOGIN_REQ(pkt.Reader)
	if err != nil {
		log.Printf("[Zone][%s] Error reading PROTO_NC_MAP_LOGIN_REQ: %v", client.RemoteAddr(), err)
		return
	}

	charName := strings.TrimRight(string(req.Chardata.Charid[:]), "\x00")
	randomID := req.Chardata.Wldmanhandle

	log.Printf("[Zone][%s] PROTO_NC_MAP_LOGIN_REQ received for character '%s' (RandomID: %d, HasChecksums: %v)",
		client.RemoteAddr(), charName, randomID, req.HasChecksums())

	// Verificación oficial binaria de DataFileCheck (Zone.exe 0x0045264b)
	if h.chkEngine != nil && req.HasChecksums() {
		var clientChecksums [49]shn.Name8
		for i := 0; i < 49; i++ {
			clientChecksums[i] = req.Checksum[i]
		}
		valid, wrongIdx := h.chkEngine.VerifyDataFiles(clientChecksums)
		if !valid {
			fileName := h.chkEngine.GetFileName(wrongIdx)
			expectedHash := h.chkEngine.GetRawHash(wrongIdx)
			clientHash := string(req.Checksum[wrongIdx][:])
			log.Printf("[Zone][%s] DataFileCheck FAILED at index %d (%s): client hash='%s', server expected='%s'",
				client.RemoteAddr(), wrongIdx, fileName, clientHash, expectedHash)

			// Persistir el incidente en tHackScanLog mediante CALL p_HackScan_AddFull oficial
			if h.accountRepo != nil {
				_ = h.accountRepo.AddHackScanLogFull(
					0,
					h.worldID,
					0,
					1, // hackType = 1 (SHN Checksum mismatch)
					client.RemoteIP(),
					charName,
					fmt.Sprintf("SHN Checksum mismatch at index %d (%s)", wrongIdx, fileName),
				)
			}

			// Error oficial en Zone.exe (0x00452727): código 0x0147 (327 decimal) + nWrongDataFileIndex (1 byte)
			failPkt := protocol.New_PROTO_NC_MAP_LOGINFAIL_ACK(0x0147, byte(wrongIdx))
			_ = client.SendPacket(failPkt)
			client.Close()
			return
		}
		log.Printf("[Zone][%s] DataFileCheck PASSED for all 49 .shn files", client.RemoteAddr())
	}

	var transfer *database.ZoneTransfer
	if h.charRepo != nil {
		var err error
		transfer, err = h.charRepo.ConsumeZoneTransfer(charName, randomID, client.RemoteIP())
		if err != nil {
			transfer, _ = h.charRepo.ConsumeZoneTransfer(charName, randomID, "")
		}
	}
	if transfer == nil && h.transfers != nil {
		var ok bool
		transfer, ok = h.transfers.ConsumeZoneTransfer(charName, randomID, pkt.Raw)
		if !ok {
			transfer = nil
		}
	}

	if transfer == nil {
		log.Printf("[Zone][%s] Transfer verification failed for '%s' (RandomID: %d)",
			client.RemoteAddr(), charName, randomID)
		errPkt := protocol.New_PROTO_NC_MAP_LOGINFAIL_ACK(4, 0) // RequestedCharacterIDNotMatching
		_ = client.SendPacket(errPkt)
		client.Close()
		return
	}

	char, err := h.charRepo.GetCharacterByName(transfer.CharName)
	if err != nil || char == nil {
		char, err = h.charRepo.GetCharacterByID(transfer.CharID)
	}

	if err != nil || char == nil {
		log.Printf("[Zone][%s] Character '%s' not found in database: %v", client.RemoteAddr(), transfer.CharName, err)
		errPkt := protocol.New_PROTO_NC_MAP_LOGINFAIL_ACK(4, 0)
		_ = client.SendPacket(errPkt)
		client.Close()
		return
	}

	var cs *combat.Store
	if h.dataStore != nil {
		cs = h.dataStore.CombatStore
	}
	player := NewZonePlayer(client, char, h.classStore, cs)
	player.WorldNo = h.worldID
	if h.charRepo != nil {
		if ap, err := h.charRepo.GetAutoPick(char.ID); err == nil {
			player.IsAutoPick = ap
		}
		if titles, err := h.charRepo.GetCharacterTitles(char.ID); err == nil {
			for _, t := range titles {
				if t != nil && t.Title != "" && strings.Contains(t.Title, ":") {
					player.CurrentTitleType = byte(t.TitleType)
					player.CurrentTitleElement = t.ElementNo
					break
				}
			}
		}
	}
	mapID := transfer.TargetMap
	if mapID == "" {
		mapID = char.TownName
	}
	if mapID == "" && h.fieldStore != nil {
		if f, ok := h.fieldStore.GetFieldBySerial(0); ok && f != nil {
			mapID = f.MapID
		}
	}
	if transfer.TargetX > 0 && transfer.TargetY > 0 {
		player.Char.XPos = transfer.TargetX
		player.Char.YPos = transfer.TargetY
	}
	zMap := h.mapManager.GetOrCreateMap(mapID)
	player.CurrentMap = zMap
	zMap.AssignObjectID(player)
	h.mu.Lock()
	h.players[client] = player
	h.mu.Unlock()

	client.Username = transfer.Username
	client.AccountID = transfer.AccountID
	client.Admin = transfer.Admin
	client.IsAuthenticated = true

	log.Printf("[Zone][%s] Character '%s' authenticated successfully. Dispatching In-Game Chunk...",
		client.RemoteAddr(), char.Name)

	h.sendInGameChunk(player)
}

// sendInGameChunk despacha la secuencia estricta de paquetes de entrada a Zona (cliente NA2016):
//  1. 0x1038 (PROTO_NC_CHAR_CLIENT_BASE_CMD)             - Siempre
//  2. 0x1039 (PROTO_NC_CHAR_CLIENT_SHAPE_CMD)            - Siempre
//  3. 0x103A (PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD)      - Siempre
//  4. 0x103B (PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD)       - Siempre
//  5. 0x10CE (PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD)       - Siempre
//  6. 0x10D7 (PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD)     - Siempre
//  7. 0x103E (PROTO_NC_CHAR_CLIENT_PASSIVE_CMD)          - Siempre
//  8. 0x104A (PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD)      - Siempre (incondicional, 0 si está vacío)
//  9. 0x103D (PROTO_NC_CHAR_CLIENT_SKILL_CMD)            - Solo si len(activeSkills) > 0 (PartMark = 3)
// 10. 0x1048 (PROTO_NC_CHAR_CLIENT_GAME_CMD)             - Siempre
// 11. 0x1047 (PROTO_NC_CHAR_CLIENT_ITEM_CMD, Box 8)      - Solo si len > 0 (Equip)
// 12. 0x1047 (PROTO_NC_CHAR_CLIENT_ITEM_CMD, Box 9)      - Solo si len > 0 (Inven)
// 13. 0x1047 (PROTO_NC_CHAR_CLIENT_ITEM_CMD, Box 12)     - Solo si len > 0 (House)
// 14. 0x1047 (PROTO_NC_CHAR_CLIENT_ITEM_CMD, Box 15)     - Siempre (ActionItem, incluso vacío)
// 15. 0x1049 (PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD)        - Siempre
// 16. 0x3801/0x3802/0x3803 (Card Regist)                 - Solo si tiene cartas registradas
// 17. 0x1802 (PROTO_NC_MAP_LOGIN_ACK)                    - Siempre (último paquete del bloque)
func (h *Handler) sendInGameChunk(p *ZonePlayer) {
	// 1. 0x1038 - PROTO_NC_CHAR_CLIENT_BASE_CMD (Siempre, wire len = 107, payload = 105)
	infoPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_BASE_CMD)
	p.WriteDetailedInfo(infoPkt)
	_ = p.Client.SendPacket(infoPkt)

	// 2. 0x1039 - PROTO_NC_CHAR_CLIENT_SHAPE_CMD (Siempre, wire len = 6, payload = 4)
	lookPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_SHAPE_CMD)
	p.WriteLook(lookPkt)
	_ = p.Client.SendPacket(lookPkt)

	// 3-6. Quests (Doing, Done, Read, Repeat) cargadas 100% desde MySQL (tQuest / tQuestTimes)
	var doingList []*database.QuestData
	var repeatList []*database.QuestData
	var doneQuests []int
	if h.charRepo != nil && p.Char != nil {
		allQuests, _ := h.charRepo.GetCharacterQuests(p.Char.ID)
		for _, q := range allQuests {
			if q.Status == 1 {
				doingList = append(doingList, q)
			} else if q.Status == 3 {
				repeatList = append(repeatList, q)
			}
		}
		doneQuests, _ = h.charRepo.GetCharacterDoneQuests(p.Char.ID)
	}

	// 3. 0x103A - PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD
	questsBusy := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD)
	questsBusy.WriteUint32(uint32(p.Char.ID))
	questsBusy.WriteByte(1)                    // bNeedClear = 1
	questsBusy.WriteByte(byte(len(doingList))) // nNumOfDoingQuest
	for _, q := range doingList {
		questsBusy.WriteUint16(uint16(q.QuestNo))
		if len(q.Data) >= 30 {
			questsBusy.WriteBytes(q.Data[:30])
		} else {
			questsBusy.WriteBytes(q.Data)
			questsBusy.Fill(30-len(q.Data), 0)
		}
	}
	_ = p.Client.SendPacket(questsBusy)

	// 4. 0x103B - PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD
	questsDone := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD)
	questsDone.WriteUint32(uint32(p.Char.ID))
	questsDone.WriteUint16(uint16(len(doneQuests))) // nTotalDoneQuest
	questsDone.WriteUint16(uint16(len(doneQuests))) // nTotalDoneQuestSize
	questsDone.WriteUint16(uint16(len(doneQuests))) // nDoneQuestCount
	questsDone.WriteUint16(0)                       // nIndex
	for _, qID := range doneQuests {
		questsDone.WriteUint16(uint16(qID))
	}
	_ = p.Client.SendPacket(questsDone)

	// 5. 0x10CE - PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD (Siempre, wire len = 8, payload = 6)
	questsRead := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_READ_CMD)
	questsRead.WriteUint32(uint32(p.Char.ID))
	questsRead.WriteUint16(0) // nNumOfReadQuest = 0
	_ = p.Client.SendPacket(questsRead)

	// 6. 0x10D7 - PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD (Siempre, wire len = 8, payload = 6)
	questsRepeat := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_QUEST_REPEAT_CMD)
	questsRepeat.WriteUint32(uint32(p.Char.ID))
	questsRepeat.WriteUint16(uint16(len(repeatList))) // nNumOfRepeatQuest
	for _, q := range repeatList {
		questsRepeat.WriteUint16(uint16(q.QuestNo))
		if len(q.Data) >= 30 {
			questsRepeat.WriteBytes(q.Data[:30])
		} else {
			questsRepeat.WriteBytes(q.Data)
			questsRepeat.Fill(30-len(q.Data), 0)
		}
	}
	_ = p.Client.SendPacket(questsRepeat)

	// 7. 0x103E - PROTO_NC_CHAR_CLIENT_PASSIVE_CMD (Siempre, clasificado por PassiveSkill.shn)
	var activeSkills []*database.CharacterSkill
	var passiveSkills []*database.CharacterSkill
	if h.charRepo != nil && p.Char != nil {
		skills, _ := h.charRepo.GetCharacterSkills(p.Char.ID)
		for _, sk := range skills {
			if h.dataStore != nil && h.dataStore.IsPassiveSkill(sk.SkillNo) {
				passiveSkills = append(passiveSkills, sk)
			} else {
				activeSkills = append(activeSkills, sk)
			}
		}
	}

	passiveSkillsPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_PASSIVE_CMD)
	passiveSkillsPkt.WriteUint16(uint16(len(passiveSkills)))
	for _, ps := range passiveSkills {
		passiveSkillsPkt.WriteUint16(ps.SkillNo)
	}
	_ = p.Client.SendPacket(passiveSkillsPkt)

	// 8. 0x104A - PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD (CRÍTICO: Siempre, incondicional)
	var chargedBuffs []*database.ChargedBuffData
	if h.charRepo != nil && p.Char != nil {
		chargedBuffs, _ = h.charRepo.GetChargedBuffs(p.Char.ID)
	}
	buffPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_CHARGEDBUFF_CMD)
	buffPkt.WriteUint16(uint16(len(chargedBuffs))) // 0 si no tiene buffs (2 bytes de payload)
	now := time.Now()
	for _, buff := range chargedBuffs {
		timeLeft := uint32(0)
		keepTime := uint32(0)
		if !buff.EndDate.IsZero() && buff.EndDate.After(now) {
			timeLeft = uint32(buff.EndDate.Sub(now).Seconds())
			keepTime = uint32(buff.EndDate.Sub(buff.UseDate).Seconds())
		}
		buffPkt.WriteUint16(uint16(buff.ID))
		buffPkt.WriteUint32(timeLeft)
		buffPkt.WriteUint32(keepTime)
		buffPkt.WriteUint32(uint32(buff.Key))
	}
	_ = p.Client.SendPacket(buffPkt)

	// 9. 0x103D - PROTO_NC_CHAR_CLIENT_SKILL_CMD (Solo si tiene skills activas, con PartMark = 3)
	if len(activeSkills) > 0 {
		skillPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_SKILL_CMD)
		skillPkt.WriteByte(0)                               // restempow (1 byte)
		skillPkt.WriteByte(3)                               // PartMark = 3 (first=1 | last=1, 1 byte)
		skillPkt.WriteUint16(uint16(len(activeSkills)))     // nMaxNum (2 bytes)
		skillPkt.WriteUint32(uint32(p.Char.ID))             // chrregnum (4 bytes)
		skillPkt.WriteUint16(uint16(len(activeSkills)))     // number (2 bytes)
		for _, sk := range activeSkills {
			skillPkt.WriteUint16(sk.SkillNo)                // skillid (2 bytes)
			skillPkt.WriteUint32(uint32(sk.SkillCoolTime))  // cooltime (4 bytes)
			skillPkt.WriteUint16(0)                         // empow (2 bytes)
			skillPkt.WriteUint32(uint32(sk.SkillExp))       // mastery (4 bytes)
		}
		_ = p.Client.SendPacket(skillPkt)
	}

	// 10. 0x1048 - PROTO_NC_CHAR_CLIENT_GAME_CMD (Siempre, wire len = 6, payload = 4)
	partyNo := uint16(0xFFFF)
	raidNo := uint16(0xFFFF)
	if h.charRepo != nil && p.Char != nil {
		_, _ = h.charRepo.GetCharacterOptions(p.Char.ID)
		_, _ = h.charRepo.GetCharacterStates(p.Char.ID)
		_, _ = h.charRepo.GetItemActionCoolTime(p.Char.ID)
		_, _ = h.charRepo.GetRestExpExecTime(p.Char.ID)
		if ptNo, _, _, err := h.charRepo.GetPartyByChar(p.Char.ID); err == nil && ptNo > 0 {
			partyNo = uint16(ptNo)
		}
	}
	gamePkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_GAME_CMD)
	gamePkt.WriteUint16(partyNo)
	gamePkt.WriteUint16(raidNo)
	_ = p.Client.SendPacket(gamePkt)

	// 11-14. 0x1047 - PROTO_NC_CHAR_CLIENT_ITEM_CMD:
	// - Box 8 (Equip):      Solo si len > 0
	// - Box 9 (Inven):      Solo si len > 0
	// - Box 12 (House):     Solo si len > 0
	// - Box 15 (ActionItem): Siempre (incluso vacío: numofitem=0, box=15, flag=1)
	var charItems []*database.Item
	if h.charRepo != nil && p.Char != nil {
		charItems, _ = h.charRepo.GetCharacterItems(p.Char.ID)
	}
	itemsByBox := make(map[uint8][]*database.Item)
	for _, it := range charItems {
		boxType := it.StorageType
		if boxType == 0 {
			boxType = 8
		}
		itemsByBox[boxType] = append(itemsByBox[boxType], it)
		if boxType == 8 {
			if p.Char.Equips == nil {
				p.Char.Equips = make(map[byte]uint16)
			}
			p.Char.Equips[byte(it.Storage)] = it.ItemID
		}
	}

	sendItemBox := func(box uint8, items []*database.Item, sendIfEmpty bool) {
		if len(items) == 0 && !sendIfEmpty {
			return
		}
		itemPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_ITEM_CMD)
		if len(items) > 0 {
			itemPkt.WriteByte(byte(len(items))) // numofitem
			itemPkt.WriteByte(box)              // box
			itemPkt.WriteByte(1)                // flag = 1
			for _, it := range items {
				loc := (uint16(box) << 10) | (uint16(it.Storage) & 0x03FF)
				var statOpts []item.ItemStatOption
				if h.charRepo != nil && it.ItemKey > 0 {
					options, _ := h.charRepo.GetItemOptions(it.ItemKey)
					for k, v := range options {
						statOpts = append(statOpts, item.ItemStatOption{Type: uint8(k), Value: uint16(v)})
					}
				}
				attrBuf := h.itemStore.BuildItemAttr(it.ItemID, uint32(it.Flags), statOpts, 0)
				dataSize := byte(2 + 2 + len(attrBuf))
				itemPkt.WriteByte(dataSize)    // datasize: 2 (location) + 2 (itemid) + len(attr)
				itemPkt.WriteUint16(loc)       // location: ITEM_INVEN (uint16 bitfield: (box << 10) | slot)
				itemPkt.WriteUint16(it.ItemID) // ItemID
				itemPkt.WriteBytes(attrBuf)    // attr: dynamic payload emulating iac_fill2buffer
			}
		} else {
			itemPkt.WriteByte(0)   // numofitem = 0
			itemPkt.WriteByte(box) // box = 15 (ActionItem)
			itemPkt.WriteByte(1)   // flag = 1
		}
		_ = p.Client.SendPacket(itemPkt)
	}

	// 11. 0x1047 (Box 8: Equip) - Solo si len > 0
	sendItemBox(8, itemsByBox[8], false)
	// 12. 0x1047 (Box 9: Inven) - Solo si len > 0
	sendItemBox(9, itemsByBox[9], false)
	// 13. 0x1047 (Box 12: House) - Solo si len > 0
	sendItemBox(12, itemsByBox[12], false)
	// 14. 0x1047 (Box 15: ActionItem) - Siempre (incluso vacío)
	sendItemBox(15, itemsByBox[15], true)

	// 15. 0x1049 - PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD (Siempre, wire len = 8, payload = 6)
	var titles []*database.CharacterTitleData
	if h.charRepo != nil && p.Char != nil {
		titles, _ = h.charRepo.GetCharacterTitles(p.Char.ID)
	}
	titlesPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD)
	titlesPkt.WriteByte(p.CurrentTitleType)    // CurrentTitle
	titlesPkt.WriteByte(p.CurrentTitleElement) // CurrentTitleElement
	titlesPkt.WriteUint16(p.CurrentTitleMobID) // CurrentTitleMobID
	titlesPkt.WriteUint16(uint16(len(titles))) // NumOfTitle
	for _, t := range titles {
		titlesPkt.WriteByte(byte(t.TitleType))
		titlesPkt.WriteByte(t.ElementNo)
	}
	_ = p.Client.SendPacket(titlesPkt)

	// 16. 0x10E4 / 0x10E7 / 0x10EA (Colección de Cartas, Favoritos y Recompensas)
	h.sendPlayerCollectionOnLogin(p)

	// 17. 0x1802 - PROTO_NC_MAP_LOGIN_ACK (CRÍTICO: Absolutamente el último paquete del bloque)
	p.RecalculateEquipmentStats(h.itemStore)
	extraPkt := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_LOGIN_ACK)
	p.WriteDetailedInfoExtra(extraPkt, false)
	_ = p.Client.SendPacket(extraPkt)
}

func (h *Handler) handleClientReady(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	log.Printf("[Zone][%s] ClientReady received from '%s'. Spawning into '%s' (MapID: %s, Serial: %d)...",
		client.RemoteAddr(), player.Char.Name, player.CurrentMap.Name, player.CurrentMap.MapID, player.CurrentMap.ID)

	// Add player to map and broadcast spawn
	player.CurrentMap.AddPlayer(player)

	// Loguear entrada de mapa en GameLog Tipo 1 (nType = 1 para ingreso)
	if h.gameLogDispatcher != nil && player.Char != nil {
		h.gameLogDispatcher.EnqueueType1(1, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0)
	}
}

func (h *Handler) handleWalk(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	oldX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	oldY, _ := pkt.Reader.ReadInt32()
	newX, _ := pkt.Reader.ReadInt32()
	newY, _ := pkt.Reader.ReadInt32()

	player.CurrentMap.MovePlayer(player, int(oldX), int(oldY), int(newX), int(newY), true)
}

func (h *Handler) handleRun(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	oldX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	oldY, _ := pkt.Reader.ReadInt32()
	newX, _ := pkt.Reader.ReadInt32()
	newY, _ := pkt.Reader.ReadInt32()

	player.CurrentMap.MovePlayer(player, int(oldX), int(oldY), int(newX), int(newY), false)
}

// handleMoveWalk implementa sp_NC_ACT_MOVEWALK_CMD (0x2017) según Zone.pdb (RVA 0x56C60):
// Lee PROTO_NC_ACT_MOVEWALK_CMD (from: SHINE_XY_TYPE 8B, to: SHINE_XY_TYPE 8B) y transmite 0x2018 (PROTO_NC_ACT_SOMEONEMOVEWALK_CMD).
func (h *Handler) handleMoveWalk(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	fromX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	fromY, _ := pkt.Reader.ReadInt32()
	toX, _ := pkt.Reader.ReadInt32()
	toY, _ := pkt.Reader.ReadInt32()

	if fromX == toX && fromY == toY {
		player.CurrentMap.sendMoveFail(player)
		return
	}

	player.CurrentMap.MovePlayerCmd(player, int(fromX), int(fromY), int(toX), int(toY), true)
}

// handleMoveRun implementa sp_NC_ACT_MOVERUN_CMD (0x2019) según Zone.pdb (RVA 0x56FD0):
// Lee PROTO_NC_ACT_MOVERUN_CMD (from: SHINE_XY_TYPE 8B, to: SHINE_XY_TYPE 8B) y transmite 0x201A (PROTO_NC_ACT_SOMEONEMOVERUN_CMD).
func (h *Handler) handleMoveRun(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	fromX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	fromY, _ := pkt.Reader.ReadInt32()
	toX, _ := pkt.Reader.ReadInt32()
	toY, _ := pkt.Reader.ReadInt32()

	if fromX == toX && fromY == toY {
		player.CurrentMap.sendMoveFail(player)
		return
	}

	player.CurrentMap.MovePlayerCmd(player, int(fromX), int(fromY), int(toX), int(toY), false)

	// Verificar activación de Gate si el jugador corrió hacia el portal
	if time.Since(player.LastWarpTime) >= 3*time.Second {
		if gate, found := player.CurrentMap.FindGateNear(int(toX), int(toY), 180.0); found && gate.GateToWhere != "" {
			h.triggerGateWarp(player, gate.GateToWhere)
		}
	}
}

func (h *Handler) handleStop(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	newX, err := pkt.Reader.ReadInt32()
	if err != nil {
		return
	}
	newY, _ := pkt.Reader.ReadInt32()

	player.CurrentMap.StopPlayer(player, int(newX), int(newY))

	// Verificar si el jugador pisó la zona de activación de un portal/gate
	if time.Since(player.LastWarpTime) >= 3*time.Second {
		if gate, found := player.CurrentMap.FindGateNear(int(newX), int(newY), 180.0); found && gate.GateToWhere != "" {
			h.triggerGateWarp(player, gate.GateToWhere)
		}
	}
}

func (h *Handler) handleJump(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	player.CurrentMap.JumpPlayer(player)
}

// handleEmoticon maneja 0x2020 (PROTO_NC_ACT_EMOTICON_CMD -> sp_NC_ACT_EMOTICON_CMD / sp_DoEmotion en Zone.pdb).
// Lee PROTO_NC_ACT_EMOTICON_CMD (emoticonid: uint8, 1 byte) y transmite 0x2021 (PROTO_NC_ACT_SOMEONEEMOTICON_CMD).
func (h *Handler) handleEmoticon(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	emoticonID, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	player.CurrentMap.EmoticonPlayer(player, emoticonID)
}

// handleEmoticonStop maneja 0x2022 (PROTO_NC_ACT_EMOTICONSTOP_CMD -> sp_NC_ACT_EMOTICONSTOP_CMD en Zone.pdb).
// Lee STOPEMOTICON_DESCRIPT (emoticonid: uint8 1B, emoticonframe: uint16 2B) y transmite 0x2023 (PROTO_NC_ACT_SOMEONEEMOTICONSTOP_CMD).
func (h *Handler) handleEmoticonStop(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	emoticonID, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	emoticonFrame, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	player.CurrentMap.EmoticonStopPlayer(player, emoticonID, emoticonFrame)
}

func (h *Handler) handleChatNormal(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	// PROTO_NC_ACT_CHAT_REQ (0x2001): uint8 itemLinkDataCount + uint8 len + char content[len] + itemLinkData[count*103]
	itemLinkCount, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	textLen, err := pkt.Reader.ReadByte()
	if err != nil || (textLen == 0 && itemLinkCount == 0) {
		return
	}

	var text string
	if textLen > 0 {
		rawText, err := pkt.Reader.ReadString(int(textLen))
		if err != nil {
			return
		}
		text = rawText
		if itemLinkCount == 0 {
			text = strings.TrimSpace(rawText)
		}
	}

	var linkData []byte
	if itemLinkCount > 0 {
		const itemLinkEntrySize = 103 // imul edx, edx, 0x67 en Zone.exe (so_Chat @ 0x458d7f)
		expected := int(itemLinkCount) * itemLinkEntrySize
		if pkt.Reader.Remaining() < expected {
			return
		}
		linkData, err = pkt.Reader.ReadBytes(expected)
		if err != nil {
			return
		}
	}

	trimmed := strings.TrimSpace(text)
	if trimmed == "" && itemLinkCount == 0 {
		return
	}

	if itemLinkCount == 0 && h.handleGMCommand(player, trimmed) {
		return
	}

	log.Printf("[Zone][Chat][%s] %s: %s (links=%d)", player.CurrentMap.Name, player.Char.Name, trimmed, itemLinkCount)

	color := byte(0x2A) // Normal chat color
	if (player.Client != nil && player.Client.Admin > 0) || (player.Char != nil && player.Char.AdminLevel > 0) {
		color = 0x03 // GM gold color
	}

	player.CurrentMap.NormalChatWithLinks(player, text, color, itemLinkCount, linkData)
}

func (h *Handler) handleBeginRest(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	var tentID uint16 = 1
	if pkt.Reader.Remaining() >= 2 {
		if id, err := pkt.Reader.ReadUint16(); err == nil {
			tentID = id
		}
	}

	player.CurrentMap.PitchTent(player, tentID)

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_PITCHTENT_ACK)
	pw.WriteUint16(0x0A81) // OK to rest
	_ = client.SendPacket(pw)
}

func (h *Handler) handleEndRest(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	player.CurrentMap.FoldTent(player)

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_FOLDTENT_ACK)
	pw.WriteUint16(0x0A81)
	_ = client.SendPacket(pw)
}

func (h *Handler) handleSeedReq(client *network.Client, pkt *protocol.Packet) {
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_SEED_ACK)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)
}

func (h *Handler) handleHideExceptMeOn(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	player.HideExceptMe = true
	surrounding := player.CurrentMap.grid.GetSurroundingPlayers(player.Char.XPos, player.Char.YPos)
	for _, other := range surrounding {
		if other.MapObjectID == player.MapObjectID {
			continue
		}
		delPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
		delPkt.WriteUint16(other.MapObjectID)
		_ = client.SendPacket(delPkt)
	}
}

func (h *Handler) handleHideExceptMeOff(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	player.HideExceptMe = false
	surrounding := player.CurrentMap.grid.GetSurroundingPlayers(player.Char.XPos, player.Char.YPos)
	for _, other := range surrounding {
		if other.MapObjectID == player.MapObjectID {
			continue
		}
		spawnOther := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
		other.WriteCharacterDisplay(spawnOther)
		_ = client.SendPacket(spawnOther)
	}
}

func (h *Handler) handleBriefInfoInform(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	if pkt.Reader.Remaining() < 6 {
		return
	}
	_, _ = pkt.Reader.ReadUint16() // myHnd
	_, _ = pkt.Reader.ReadUint16() // cmd
	targetHnd, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	zMap := player.CurrentMap
	zMap.mu.RLock()
	otherPlayer, isPlayer := zMap.players[targetHnd]
	mobObj, isMob := zMap.mobs[targetHnd]
	npcObj, isNPC := zMap.npcs[targetHnd]
	itemObj, isItem := zMap.droppedItems[targetHnd]
	tentID, isTent := zMap.tents[targetHnd]
	zMap.mu.RUnlock()

	if isPlayer && otherPlayer != nil {
		if !player.HideExceptMe {
			pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD)
			otherPlayer.WriteCharacterDisplay(pw)
			_ = client.SendPacket(pw)
		}
	} else if isMob && mobObj != nil && !mobObj.IsDead {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_REGENMOB_CMD)
		mobObj.WriteRegenMob(pw, 0)
		_ = client.SendPacket(pw)
	} else if isNPC && npcObj != nil {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_MOB_CMD)
		pw.WriteByte(1)
		npcObj.WriteRegenMob(pw)
		_ = client.SendPacket(pw)
	} else if isItem && itemObj != nil {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_DROPEDITEM_CMD)
		itemObj.WriteDropEntry(pw)
		_ = client.SendPacket(pw)
	} else if isTent {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEPITCHTENT_CMD)
		pw.WriteUint16(targetHnd)
		pw.WriteUint16(tentID)
		pw.WriteUint32(0xFFFFFFFF)
		pw.WriteUint32(0xFFFFFFFF)
		pw.WriteUint16(0xFFFF)
		_ = client.SendPacket(pw)
	}
}

func (h *Handler) handleWingSave(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	if player.Char != nil {
		player.Char.TownName = player.CurrentMap.MapID
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_WING_SAVE_ACK)
	ack.WriteByte(0) // success
	_ = client.SendPacket(ack)
}

func (h *Handler) handleToNormalCoord(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	if player.CurrentMap.blockInfo != nil {
		tx, ty := blockinfo.WorldToTile(float64(player.Char.XPos), float64(player.Char.YPos))
		if player.CurrentMap.blockInfo.IsMoveBlocked(tx, ty) {
			if nx, ny, ok := player.CurrentMap.blockInfo.FindNearestWalkable(tx, ty, 50); ok {
				wx, wy := blockinfo.TileToWorld(nx, ny)
				player.Char.XPos = int(wx)
				player.Char.YPos = int(wy)
			}
		}
	}
}

func (h *Handler) handleGatherStart(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	var objHandle uint16
	if pkt.Reader.Remaining() >= 2 {
		if hnd, err := pkt.Reader.ReadUint16(); err == nil {
			objHandle = hnd
		}
	}
	player.GatheringObj = objHandle

	// PROTO_NC_ACT_GATHERSTART_ACK (6 bytes: picker uint16, object uint16, err uint16)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_GATHERSTART_ACK)
	ack.WriteUint16(player.MapObjectID)
	ack.WriteUint16(objHandle)
	ack.WriteUint16(0) // success
	_ = client.SendPacket(ack)

	// PROTO_NC_ACT_SOMEONEGATHERSTART_CMD (6 bytes: player uint16, object uint16, picker uint16)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEGATHERSTART_CMD)
	pw.WriteUint16(player.MapObjectID)
	pw.WriteUint16(objHandle)
	pw.WriteUint16(player.MapObjectID)
	player.CurrentMap.grid.BroadcastAOI(pw, player.Char.XPos, player.Char.YPos, player.MapObjectID)
}

func (h *Handler) handleGatherCancel(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	objHandle := player.GatheringObj
	player.GatheringObj = 0

	// PROTO_NC_ACT_SOMEONEGATHERCANCEL_CMD (4 bytes: player uint16, object uint16)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEGATHERCANCEL_CMD)
	pw.WriteUint16(player.MapObjectID)
	pw.WriteUint16(objHandle)
	player.CurrentMap.grid.BroadcastAOI(pw, player.Char.XPos, player.Char.YPos, 0)
}

func (h *Handler) handleGatherComplete(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	objHandle := player.GatheringObj
	player.GatheringObj = 0

	// PROTO_NC_ACT_GATHERCOMPLETE_ACK (4 bytes: err uint16, object uint16)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_GATHERCOMPLETE_ACK)
	ack.WriteUint16(0) // err = 0
	ack.WriteUint16(objHandle)
	_ = client.SendPacket(ack)

	// PROTO_NC_ACT_SOMEONEGATHERCOMPLETE_CMD (4 bytes: player uint16, object uint16)
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONEGATHERCOMPLETE_CMD)
	pw.WriteUint16(player.MapObjectID)
	pw.WriteUint16(objHandle)
	player.CurrentMap.grid.BroadcastAOI(pw, player.Char.XPos, player.Char.YPos, player.MapObjectID)
}

func (h *Handler) handleRideFeeding(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_RIDE_FEEDING_ACK)
	ack.WriteUint16(0) // success
	_ = client.SendPacket(ack)
}

func (h *Handler) handleMoverRideOn(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	var moverHnd uint16 = 1
	var slot uint8 = 0
	if pkt.Reader.Remaining() >= 3 {
		if hnd, err := pkt.Reader.ReadUint16(); err == nil {
			moverHnd = hnd
		}
		if s, err := pkt.Reader.ReadByte(); err == nil {
			slot = s
		}
	}

	player.MoverHandle = moverHnd
	player.MoverSlot = slot
	player.MoverSpeed = 220
	player.State = 6 // Mount

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_RIDE_ON_CMD)
	ack.WriteUint16(moverHnd)
	ack.WriteByte(slot)
	ack.WriteByte(0)
	_ = client.SendPacket(ack)

	someonePw := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_SOMEONE_RIDE_ON_CMD)
	someonePw.WriteUint16(player.MapObjectID)
	someonePw.WriteUint16(moverHnd)
	someonePw.WriteByte(slot)
	someonePw.WriteByte(0)
	player.CurrentMap.grid.BroadcastAOI(someonePw, player.Char.XPos, player.Char.YPos, player.MapObjectID)

	spdPw := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_MOVESPEED_CMD)
	spdPw.WriteUint16(player.MapObjectID)
	spdPw.WriteUint16(player.MoverSpeed)
	spdPw.WriteUint16(0)
	player.CurrentMap.grid.BroadcastAOI(spdPw, player.Char.XPos, player.Char.YPos, 0)

	if h.gameLogDispatcher != nil && player.Char != nil {
		h.gameLogDispatcher.EnqueueTypeG(1, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(moverHnd), int64(moverHnd), 1, int(slot), 0)
	}
}

func (h *Handler) handleMoverRideOff(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}
	oldHnd := player.MoverHandle
	player.MoverHandle = 0
	player.MoverSlot = 0
	player.MoverSpeed = 0
	player.State = 1 // Normal

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_RIDE_OFF_CMD)
	ack.WriteUint16(0)
	_ = client.SendPacket(ack)

	someonePw := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_SOMEONE_RIDE_OFF_CMD)
	someonePw.WriteUint16(player.MapObjectID)
	player.CurrentMap.grid.BroadcastAOI(someonePw, player.Char.XPos, player.Char.YPos, player.MapObjectID)

	spdPw := protocol.NewPacketWriter(protocol.PROTO_NC_MOVER_MOVESPEED_CMD)
	spdPw.WriteUint16(player.MapObjectID)
	spdPw.WriteUint16(player.GetMoveSpeed(false))
	spdPw.WriteUint16(0)
	player.CurrentMap.grid.BroadcastAOI(spdPw, player.Char.XPos, player.Char.YPos, 0)

	if h.gameLogDispatcher != nil && player.Char != nil {
		h.gameLogDispatcher.EnqueueTypeG(2, player.Char.ID, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, 0, int(oldHnd), int64(oldHnd), 0, 0, 0)
	}
}

// SetAccountRepo configura el repositorio de cuentas y logs oficial (tLogout, tReport, etc.).
// SetItemStore sets the item data store for equipment and attribute serialization.
func (h *Handler) SetItemStore(store *item.Store) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.itemStore = store
}

func (h *Handler) SetAccountRepo(repo *database.MySQLAccountRepository) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.accountRepo = repo
}

// ActivePlayersCount retorna la cantidad actual de jugadores conectados en la zona.
func (h *Handler) ActivePlayersCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.players)
}

// ActivePlayersSnapshot retorna una instantánea de los jugadores para telemetría de clases.
func (h *Handler) ActivePlayersSnapshot() []*ZonePlayer {
	h.mu.RLock()
	defer h.mu.RUnlock()
	list := make([]*ZonePlayer, 0, len(h.players))
	for _, p := range h.players {
		if p != nil {
			list = append(list, p)
		}
	}
	return list
}

// handleRewardInvenOpen implementa sp_NC_ITEM_REWARDINVENOPEN_REQ (0x302C) y so_ply_FillRewardInven (0x302D)
// según Zone.pdb (RewardInven.cpp, RVA 0xFB060 y 0xFB310):
//   - Lee uint16 page (24 slots por página: from = page * 24, to = page * 24 + 23)
//   - Consulta tItem con nStorageType = 2 (Reward Inventory)
//   - Responde con 0x302D (PROTO_NC_ITEM_REWARDINVENOPEN_ACK): itemcounter (1B) + PROTO_ITEMPACKET_INFORM[itemcounter]
func (h *Handler) handleRewardInvenOpen(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	page, err := pkt.Reader.ReadUint16()
	if err != nil {
		page = 0
	}
	fromSlot := int16(page * 24)
	toSlot := int16(page*24 + 23)

	var pageItems []*database.Item
	if h.charRepo != nil {
		allRewardItems, err := h.charRepo.GetItemsByStorage(player.Char.ID, 2) // nStorageType = 2 (Reward Storage en Zone.exe)
		if err == nil {
			for _, it := range allRewardItems {
				if it != nil && it.Storage >= fromSlot && it.Storage <= toSlot {
					pageItems = append(pageItems, it)
				}
			}
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REWARDINVENOPEN_ACK)
	ack.WriteByte(byte(len(pageItems))) // itemcounter (1 byte)
	for _, it := range pageItems {
		loc := (uint16(2) << 10) | (uint16(it.Storage) & 0x03FF)
		var attrBuf []byte
		if h.itemStore != nil {
			attrBuf = h.itemStore.BuildItemAttr(it.ItemID, uint32(it.Flags), nil, 0)
		} else {
			attrBuf = []byte{1}
		}
		dataSize := byte(2 + 2 + len(attrBuf))
		ack.WriteByte(dataSize)    // datasize
		ack.WriteUint16(loc)       // location: ITEM_INVEN
		ack.WriteUint16(it.ItemID) // ItemID
		ack.WriteBytes(attrBuf)    // attr (iac_fill2buffer)
	}
	_ = client.SendPacket(ack)
}

// handleChargedInvenOpen implementa sp_NC_ITEM_CHARGEDINVENOPEN_REQ (0x3020) y so_ply_FillChargedInven (0x3021)
// según Zone.pdb (RVA 0x913C0 y 0x91F30):
//   - Lee uint16 page (24 slots por página: from = page * 24, to = page * 24 + 23)
//   - Consulta tChargeItem (usp_Charge_ItemSelectList / usp_ChargeItem_GetList) filtrando por userNo e isDraw = 0
//   - Responde con 0x3021 (PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK):
//     ErrorCode (uint16 = 0x1041), nPartMark (uint8 = 3), NumOfChargedItem (uint16),
//     seguido de PROTO_CHARGED_ITEM_INFO[NumOfChargedItem] (16 bytes cada uno: ItemOrderNo, ItemCode, ItemAmount, ItemRegisterDate).
func (h *Handler) handleChargedInvenOpen(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	page, err := pkt.Reader.ReadUint16()
	if err != nil {
		page = 0
	}

	var pageItems []*database.ChargedItem
	if h.accountRepo != nil {
		allCharged, err := h.accountRepo.GetChargedItems(player.Char.AccountID)
		if err == nil && len(allCharged) > 0 {
			startIdx := int(page) * 24
			endIdx := startIdx + 24
			if startIdx < len(allCharged) {
				if endIdx > len(allCharged) {
					endIdx = len(allCharged)
				}
				pageItems = allCharged[startIdx:endIdx]
			}
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK)
	ack.WriteUint16(0x1041)               // ErrorCode = 0x1041 (4161 = éxito en so_ply_FillChargedInven, Zone.exe 0x492F8E)
	ack.WriteByte(3)                      // nPartMark = 3 (first=1 | last=1, Zone.exe 0x493007)
	ack.WriteUint16(uint16(len(pageItems))) // NumOfChargedItem
	for _, ci := range pageItems {
		ack.WriteUint32(ci.RowNo)                                  // 0x00: ItemOrderNo (4 bytes)
		ack.WriteUint32(ci.GoodsNo)                                // 0x04: ItemCode (4 bytes)
		ack.WriteUint32(ci.Amount)                                 // 0x08: ItemAmount (4 bytes)
		ack.WriteUint32(database.PackShineDateTime(ci.RegisterDate)) // 0x0C: ItemRegisterDate (ShineDateTime, 4 bytes)
	}
	_ = client.SendPacket(ack)
}

// handleChargedWithdraw retira un ítem premium hacia el inventario mediante usp_Charge_ItemDraw.
func (h *Handler) handleChargedWithdraw(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	itemOrderNo, err := pkt.Reader.ReadUint32()
	if err != nil {
		return
	}

	var errCode uint16 = 0
	if h.accountRepo != nil {
		err := h.accountRepo.DrawChargedItem(int64(itemOrderNo))
		if err != nil {
			log.Printf("[Zone][%s] DrawChargedItem %d failed: %v", client.RemoteAddr(), itemOrderNo, err)
			errCode = 1
		} else {
			log.Printf("[Zone][%s] DrawChargedItem %d exitoso para usuario %d",
				client.RemoteAddr(), itemOrderNo, player.Char.AccountID)
		}
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CHARGED_WITHDRAW_ACK)
	ack.WriteUint32(itemOrderNo)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)
}

// SetWorldID asigna el worldNo para auditoría oficial y telemetría.
func (h *Handler) SetWorldID(id byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.worldID = id
}

// SetZoneNo asigna el zoneNo de la instancia para enrutamiento y saltos inter-zona.
func (h *Handler) SetZoneNo(zoneNo uint16) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.zoneNo = zoneNo
}

// SetDataStore sets the entire loaded game data store for the zone handler.
func (h *Handler) SetDataStore(ds *DataStore) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.dataStore = ds
	if ds != nil {
		if ds.ClassStore != nil {
			h.classStore = ds.ClassStore
		}
		if ds.FieldStore != nil {
			h.fieldStore = ds.FieldStore
		}
		if ds.ItemStore != nil {
			h.itemStore = ds.ItemStore
		}
		if ds.ChkEngine != nil {
			h.chkEngine = ds.ChkEngine
		}
		if ds.NPCStore != nil {
			h.npcStore = ds.NPCStore
		}
		if ds.MobInfoSHN != nil {
			h.mobInfoSHN = ds.MobInfoSHN
		}
		h.mapManager.SetNPCStore(ds.NPCStore, ds.MobInfoSHN)
		if ds.BlockInfo != nil {
			h.mapManager.SetBlockStore(ds.BlockInfo)
		}
		h.mapManager.SetDataStore(ds)
	}
}

// handleNPCClick procesa 0x200A (PROTO_NC_ACT_NPCCLICK_CMD -> sp_NC_ACT_NPCCLICK_CMD en Zone.pdb).
func (h *Handler) handleNPCClick(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	npcHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	npc, found := player.CurrentMap.GetNPC(npcHandle)
	if !found || npc == nil {
		log.Printf("[Zone][%s] Clicked unknown NPC handle: %d", client.RemoteAddr(), npcHandle)
		return
	}

	log.Printf("[Zone][%s] Player '%s' clicked NPC '%s' (Handle: %d, MobID: %d, Role: %s, RoleArg: %s)",
		client.RemoteAddr(), player.Char.Name, npc.Name, npc.Handle, npc.MobID, npc.Role, npc.RoleArg)

	player.OpenNPCMobID = npc.MobID

	// Si es un Gate / Portal, ejecutar la transición
	if npc.IsGate && npc.GateToWhere != "" {
		h.triggerGateWarp(player, npc.GateToWhere)
		return
	}

	// PROTO_NC_ACT_NPCMENUOPEN_REQ (0x201C, wire len = 4, payload = 2)
	// struct PROTO_NC_ACT_NPCMENUOPEN_REQ { unsigned short mobid; };
	openReq := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NPCMENUOPEN_REQ)
	openReq.WriteUint16(npc.MobID)
	_ = client.SendPacket(openReq)
}

// handleNPCMenuOpenAck procesa 0x201D (PROTO_NC_ACT_NPCMENUOPEN_ACK -> sp_NC_ACT_NPCMENUOPEN_ACK en Zone.pdb).
func (h *Handler) handleNPCMenuOpenAck(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	ack, _ := pkt.Reader.ReadByte()
	log.Printf("[Zone][%s] Player '%s' opened NPC menu (Ack: %d)",
		client.RemoteAddr(), player.Char.Name, ack)
}

// handleNPCActionCmd procesa 0x2072 (PROTO_NC_ACT_NPC_ACTION_CMD -> sp_NC_ACT_NPC_ACTION_CMD en Zone.pdb).
func (h *Handler) handleNPCActionCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	nType, _ := pkt.Reader.ReadByte()
	nNPCHandle, _ := pkt.Reader.ReadUint16()
	nECode, _ := pkt.Reader.ReadUint32()
	log.Printf("[Zone][%s] NPC action command received: Type=%d, NPCHandle=%d, ECode=%d",
		client.RemoteAddr(), nType, nNPCHandle, nECode)
}

// handleNPCMenuCmd procesa 0x2073 (PROTO_NC_ACT_NPC_MENU_CMD -> sp_NC_ACT_NPC_MENU_CMD en Zone.pdb).
func (h *Handler) handleNPCMenuCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}
	ack, _ := pkt.Reader.ReadByte()
	log.Printf("[Zone][%s] NPC menu response received: Ack=%d", client.RemoteAddr(), ack)
}

// handleUseGate procesa 0x2070 (PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ -> sp_NC_ACT_AUTO_WAY_FINDING_USE_GATE_REQ en Zone.pdb).
func (h *Handler) handleUseGate(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.CurrentMap == nil {
		return
	}

	// PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK (0x2071, nError: 0)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_AUTO_WAY_FINDING_USE_GATE_ACK)
	ack.WriteUint16(0) // OK
	_ = client.SendPacket(ack)

	// Buscar el gate más cercano
	gate, found := player.CurrentMap.FindGateNear(player.Char.XPos, player.Char.YPos, 350.0)
	if found && gate != nil && gate.GateToWhere != "" {
		h.triggerGateWarp(player, gate.GateToWhere)
	}
}

// triggerGateWarp consulta LinkTable.json y ejecuta la teletransportación oficial entre mapas.
func (h *Handler) triggerGateWarp(player *ZonePlayer, gateArg string) {
	if player == nil || player.Client == nil || h.dataStore == nil || h.dataStore.NPCStore == nil {
		return
	}

	if time.Since(player.LastWarpTime) < 3*time.Second {
		return
	}

	link, found := h.dataStore.NPCStore.GetLinkByArg(gateArg)
	if !found {
		log.Printf("[Zone][%s] Gate '%s' not found in LinkTable", player.Client.RemoteAddr(), gateArg)
		return
	}

	log.Printf("[Zone][%s] Gate '%s' activated -> Warping to %s (%d, %d, Dir: %d)",
		player.Client.RemoteAddr(), gateArg, link.MapServer, link.X, link.Y, link.Direct)

	h.WarpPlayer(player, link.MapServer, int(link.X), int(link.Y), byte(link.Direct))
}

// WarpPlayer traslada al jugador entre mapas oficiales, actualizando mapas, coordenadas y persistencia.
// Si el mapa pertenece a otra zona según FieldList (fiesta != h.zoneNo), realiza un traspaso inter-zona
// invocando usp_Session_IssueZoneTransfer y enviando PROTO_NC_CHAR_LOGIN_ACK.
func (h *Handler) WarpPlayer(player *ZonePlayer, targetMapID string, x, y int, rotation byte) {
	if player == nil || player.Client == nil {
		return
	}

	targetMapID = strings.TrimSpace(targetMapID)
	if targetMapID == "" && h.fieldStore != nil {
		if f0, ok := h.fieldStore.GetFieldBySerial(0); ok && f0 != nil {
			targetMapID = f0.MapID
		}
	}

	var targetZoneNo uint16 = 1
	if h.multiNode {
		targetZoneNo = h.zoneNo
	}
	if h.fieldStore != nil && targetMapID != "" {
		if f, ok := h.fieldStore.GetField(targetMapID); ok && f != nil {
			targetZoneNo = uint16(f.Fiesta)
			if h.multiNode {
				if targetZoneNo == 99 {
					targetZoneNo = 0
				}
			} else if targetZoneNo == 0 {
				targetZoneNo = 1
			}
		}
	}

	// Si el mapa destino está en otra zona (inter-zone transfer)
	if (h.multiNode || h.zoneNo > 0) && targetZoneNo != h.zoneNo && h.charRepo != nil {
		ip, port, isOnline, err := h.charRepo.ResolveMapEndpoint(h.worldID, targetZoneNo)
		if err == nil && isOnline && ip != "" && port > 0 {
			randomID := uint16(time.Now().UnixNano() & 0xFFFF)
			if randomID == 0 {
				randomID = 1
			}
			userNo := player.AccountID
			if userNo == 0 && player.Char != nil {
				userNo = player.Char.AccountID
			}
			userName := player.Username
			if userName == "" && player.Char != nil {
				userName = player.Char.Name
			}
			_ = h.charRepo.SaveCharLink(player.Char.ID, targetMapID, x, y)
			_ = h.charRepo.IssueZoneTransfer(
				h.worldID, targetZoneNo, userNo, userName,
				player.Char.ID, player.Char.Name, randomID, player.Admin,
				player.Client.RemoteIP(), targetMapID, x, y,
			)

			log.Printf("[Zone][%s] Inter-zone warp for '%s' to zone %d (%s:%d) map '%s' (%d, %d)",
				player.Client.RemoteAddr(), player.Char.Name, targetZoneNo, ip, port, targetMapID, x, y)

			ack := protocol.New_PROTO_NC_CHAR_LOGIN_ACK(ip, port)
			_ = player.Client.SendPacket(ack)

			oldMap := player.CurrentMap
			if oldMap != nil {
				oldMap.RemovePlayer(player)
			}
			player.Client.Close()
			return
		}
	}

	// Mismo nodo de zona: warp local
	oldMap := player.CurrentMap
	if oldMap != nil {
		oldMap.RemovePlayer(player)
	}

	newMap := h.mapManager.GetOrCreateMap(targetMapID)
	player.CurrentMap = newMap
	player.Char.TownName = newMap.MapID
	player.Char.XPos = x
	player.Char.YPos = y
	player.Rotation = rotation
	player.LastWarpTime = time.Now()

	if h.charRepo != nil && player.Char != nil {
		_ = h.charRepo.SaveCharLocation(player.Char.ID, newMap.MapID, x, y, rotation)
	}

	log.Printf("[Zone][%s] Player '%s' entering new map '%s' (%d, %d, Dir: %d)",
		player.Client.RemoteAddr(), player.Char.Name, newMap.MapID, x, y, rotation)

	newMap.AssignObjectID(player)
	h.sendInGameChunk(player)
}

// handleMultyLinkSelect procesa 0x181F (PROTO_NC_MAP_MULTY_LINK_SELECT_REQ -> sp_NC_MAP_MULTY_LINK_SELECT_REQ en Zone.pdb).
func (h *Handler) handleMultyLinkSelect(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	linkMapName, err := pkt.Reader.ReadString(12)
	if err != nil {
		return
	}
	destMap := strings.TrimRight(linkMapName, "\x00")

	// PROTO_NC_MAP_MULTY_LINK_SELECT_ACK (0x1820, wire len = 4, payload = 2)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_MAP_MULTY_LINK_SELECT_ACK)
	ack.WriteUint16(0) // OK
	_ = client.SendPacket(ack)

	spawnX := 6208
	spawnY := 6405
	if h.fieldStore != nil {
		if f, ok := h.fieldStore.GetField(destMap); ok && f != nil && len(f.Respawn.Spots) > 0 {
			spawnX = f.Respawn.Spots[0].X
			spawnY = f.Respawn.Spots[0].Y
		}
	}

	h.WarpPlayer(player, destMap, spawnX, spawnY, 0x55)
}

// handleTownPortal procesa 0x181A (PROTO_NC_MAP_TOWNPORTAL_REQ -> sp_NC_MAP_TOWNPORTAL_REQ en Zone.pdb).
func (h *Handler) handleTownPortal(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	portalIndex, _ := pkt.Reader.ReadByte()
	log.Printf("[Zone][%s] Town portal request (Index: %d) from '%s'",
		client.RemoteAddr(), portalIndex, player.Char.Name)

	targetMap := ""
	if player.Char != nil {
		targetMap = player.Char.TownName
	}
	if targetMap == "" && h.fieldStore != nil {
		if f0, ok := h.fieldStore.GetFieldBySerial(0); ok && f0 != nil {
			targetMap = f0.MapID
		}
	}
	spawnX := 6208
	spawnY := 6405
	if h.fieldStore != nil && targetMap != "" {
		if f, ok := h.fieldStore.GetField(targetMap); ok && f != nil && len(f.Respawn.Spots) > 0 {
			spawnX = f.Respawn.Spots[0].X
			spawnY = f.Respawn.Spots[0].Y
		}
	}
	h.WarpPlayer(player, targetMap, spawnX, spawnY, 0x55)
}

// handleWingFly procesa 0x1816 (PROTO_NC_MAP_WING_FLY_REQ -> sp_NC_MAP_WING_FLY_REQ en Zone.pdb).
func (h *Handler) handleWingFly(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	where, _ := pkt.Reader.ReadByte()
	log.Printf("[Zone][%s] Wing fly request (Where: %d) from '%s'",
		client.RemoteAddr(), where, player.Char.Name)

	targetMap := ""
	if player.Char != nil {
		targetMap = player.Char.TownName
	}
	if targetMap == "" && h.fieldStore != nil {
		if f0, ok := h.fieldStore.GetFieldBySerial(0); ok && f0 != nil {
			targetMap = f0.MapID
		}
	}
	spawnX := 6208
	spawnY := 6405
	if h.fieldStore != nil && targetMap != "" {
		if f, ok := h.fieldStore.GetField(targetMap); ok && f != nil && len(f.Respawn.Spots) > 0 {
			spawnX = f.Respawn.Spots[0].X
			spawnY = f.Respawn.Spots[0].Y
		}
	}
	h.WarpPlayer(player, targetMap, spawnX, spawnY, 0x55)
}

// sendCellChange implementa so_ply_ChangeCell (Zone.exe VA 0x518260):
//   - Si box == 8 ((location & 0xFC00) == 0x2000): envía PROTO_NC_ITEM_EQUIPCHANGE_CMD (0x3002)
//     (exchange: uint16, location: uint8, itemid: uint16 + attr si itemid != 0xFFFF)
//   - Si box != 8: envía PROTO_NC_ITEM_CELLCHANGE_CMD (0x3001)
//     (exchange: uint16, location: uint16, itemid: uint16 + attr si itemid != 0xFFFF)
func (h *Handler) sendCellChange(client *network.Client, box uint8, slot int16, it *database.Item, exchange uint16) {
	if box == 8 {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIPCHANGE_CMD)
		pw.WriteUint16(exchange)
		pw.WriteByte(byte(slot))
		if it == nil || it.ItemID == 0xFFFF {
			pw.WriteUint16(0xFFFF)
		} else {
			pw.WriteUint16(it.ItemID)
			var attrBuf []byte
			if h.itemStore != nil {
				attrBuf = h.itemStore.BuildItemAttr(it.ItemID, uint32(it.Flags), nil, 0)
			} else {
				attrBuf = []byte{byte(it.Flags)}
			}
			pw.WriteBytes(attrBuf)
		}
		_ = client.SendPacket(pw)
	} else {
		loc := (uint16(box) << 10) | (uint16(slot) & 0x03FF)
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_CELLCHANGE_CMD)
		pw.WriteUint16(exchange)
		pw.WriteUint16(loc)
		if it == nil || it.ItemID == 0xFFFF {
			pw.WriteUint16(0xFFFF)
		} else {
			pw.WriteUint16(it.ItemID)
			var attrBuf []byte
			if h.itemStore != nil {
				attrBuf = h.itemStore.BuildItemAttr(it.ItemID, uint32(it.Flags), nil, 0)
			} else {
				attrBuf = []byte{byte(it.Flags)}
			}
			pw.WriteBytes(attrBuf)
		}
		_ = client.SendPacket(pw)
	}
}

func (h *Handler) findItemAt(charID int, box uint8, slot int16) *database.Item {
	if h.charRepo == nil {
		return nil
	}
	items, err := h.charRepo.GetCharacterItems(charID)
	if err != nil {
		return nil
	}
	for _, it := range items {
		st := it.StorageType
		if st == 0 {
			st = 8
		}
		if st == box && it.Storage == slot {
			return it
		}
	}
	return nil
}

func (h *Handler) findFreeInvenSlot(charID int) (int16, bool) {
	if h.charRepo == nil {
		return 0, false
	}
	items, err := h.charRepo.GetCharacterItems(charID)
	if err != nil {
		return 0, false
	}
	used := make(map[int16]bool)
	for _, it := range items {
		if it.StorageType == 9 {
			used[it.Storage] = true
		}
	}
	for slot := int16(0); slot < 96; slot++ {
		if !used[slot] {
			return slot, true
		}
	}
	return 0, false
}

// syncPlayerEquipsAndStats actualiza el mapa de equipos del jugador, recalcula stats de equipo
// y envía PROTO_NC_CHAR_STATE_CHANGE_CMD (0x104E) si los stats cambiaron.
func (h *Handler) syncPlayerEquipsAndStats(player *ZonePlayer) {
	if player == nil || player.Char == nil || h.charRepo == nil {
		return
	}
	items, err := h.charRepo.GetCharacterItems(player.Char.ID)
	if err != nil {
		return
	}
	equips := make(map[byte]uint16)
	for _, it := range items {
		st := it.StorageType
		if st == 0 {
			st = 8
		}
		if st == 8 {
			equips[byte(it.Storage)] = it.ItemID
		}
	}
	player.Char.Equips = equips
	player.RecalculateEquipmentStats(h.itemStore)
}

// broadcastEquipAppearance transmite el cambio visual de equipamiento al mapa según Zone.exe iac_equipbriefinfo:
//   - Slot 12 (Weapon):        PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD (0x1C05, 9B)  [0x0063F820]
//   - Slot 11 (CostumeShield): PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD (0x1C02, 5B) [0x0063F910]
//   - Otros slots:             PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD (0x1C03, 6B)  [0x0063F8D0]
func (h *Handler) broadcastEquipAppearance(player *ZonePlayer, equipSlot byte, itemID uint16, upgrade byte) {
	if player == nil || player.CurrentMap == nil {
		return
	}
	if equipSlot == slotWeapon {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD)
		pw.WriteUint16(player.MapObjectID) // handle
		pw.WriteUint16(itemID)             // item
		pw.WriteByte(upgrade)              // upgrade
		pw.WriteByte(equipSlot)            // nSlotNum
		pw.WriteUint16(0xFFFF)             // currentmobid
		pw.WriteByte(0)                    // currentkilllevel
		player.CurrentMap.Broadcast(pw, 0)
	} else if equipSlot == Slot_CostumeShield {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD)
		pw.WriteUint16(player.MapObjectID) // handle
		pw.WriteUint16(itemID)             // item
		pw.WriteByte(equipSlot)            // nSlotNum
		player.CurrentMap.Broadcast(pw, 0)
	} else {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD)
		pw.WriteUint16(player.MapObjectID) // handle
		pw.WriteUint16(itemID)             // item
		pw.WriteByte(upgrade)              // upgrade
		pw.WriteByte(equipSlot)            // nSlotNum
		player.CurrentMap.Broadcast(pw, 0)
	}
}

// broadcastUnequipAppearance transmite PROTO_NC_BRIEFINFO_UNEQUIP_CMD (0x1C04, 3B) al mapa.
func (h *Handler) broadcastUnequipAppearance(player *ZonePlayer, equipSlot byte) {
	if player == nil || player.CurrentMap == nil {
		return
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_UNEQUIP_CMD)
	pw.WriteUint16(player.MapObjectID) // handle
	pw.WriteByte(equipSlot)         // slot
	player.CurrentMap.Broadcast(pw, 0)
}

// handleItemReloc implementa sp_NC_ITEM_RELOC_REQ (0x300B) y gds_NC_ITEMDB_RELOCSUC_ACK / EXCHANGESUC_ACK (Zone.exe VA 0x533070 / 0x533350):
//   - Lee from (ITEM_INVEN uint16) y to (ITEM_INVEN uint16)
//   - Si el destino está vacío: mueve el ítem (RelocateItem), envía 0x3001/0x3002 para origen (vacío) y destino (ocupado)
//   - Si el destino tiene un ítem del mismo ID y apilable (< MaxLot): combina stacks
//   - Si el destino tiene otro ítem: intercambia posiciones (ExchangeItem)
//   - Envía PROTO_NC_ITEM_RELOC_ACK (0x300C) con err = 0x0241 (577)
func (h *Handler) handleItemReloc(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	fromLoc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	toLoc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	fromBox := uint8(fromLoc >> 10)
	fromSlot := int16(fromLoc & 0x03FF)
	toBox := uint8(toLoc >> 10)
	toSlot := int16(toLoc & 0x03FF)

	if fromBox == toBox && fromSlot == toSlot {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_ACK)
		ack.WriteUint16(0x0241)
		_ = client.SendPacket(ack)
		return
	}

	srcItem := h.findItemAt(player.Char.ID, fromBox, fromSlot)
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_ACK)
		ack.WriteUint16(0x0242) // Reloc fail
		_ = client.SendPacket(ack)
		return
	}

	dstItem := h.findItemAt(player.Char.ID, toBox, toSlot)
	if dstItem == nil {
		// Mover a celda vacía
		_ = h.charRepo.SetItemStorage(srcItem.ItemKey, toBox, toSlot)
		srcItem.StorageType = toBox
		srcItem.Storage = toSlot

		h.sendCellChange(client, fromBox, fromSlot, nil, 0xFFFF)
		h.sendCellChange(client, toBox, toSlot, srcItem, 0xFFFF)
	} else if dstItem.ItemID == srcItem.ItemID && toBox == 9 && fromBox == 9 {
		// Verificar si es apilable (MaxLot > 1)
		maxLot := uint32(1)
		if h.itemStore != nil {
			if info, ok := h.itemStore.GetItemInfo(srcItem.ItemID); ok && info != nil && info.MaxLot > 1 {
				maxLot = info.MaxLot
			}
		}
		srcLot := uint32(srcItem.Flags)
		if srcLot == 0 {
			srcLot = 1
		}
		dstLot := uint32(dstItem.Flags)
		if dstLot == 0 {
			dstLot = 1
		}

		if maxLot > 1 && dstLot < maxLot {
			total := srcLot + dstLot
			if total <= maxLot {
				dstItem.Flags = int32(total)
				_ = h.charRepo.UpdateItemFlags(dstItem.ItemKey, int(dstItem.Flags))
				_ = h.charRepo.DeleteItem(srcItem.ItemKey)
				h.sendCellChange(client, fromBox, fromSlot, nil, 0xFFFF)
				h.sendCellChange(client, toBox, toSlot, dstItem, 0xFFFF)
			} else {
				dstItem.Flags = int32(maxLot)
				srcItem.Flags = int32(total - maxLot)
				_ = h.charRepo.UpdateItemFlags(dstItem.ItemKey, int(dstItem.Flags))
				_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(srcItem.Flags))
				h.sendCellChange(client, fromBox, fromSlot, srcItem, 0xFFFF)
				h.sendCellChange(client, toBox, toSlot, dstItem, 0xFFFF)
			}
		} else {
			_ = h.charRepo.ExchangeItemStorage(srcItem.ItemKey, toBox, toSlot, dstItem.ItemKey, fromBox, fromSlot)
			srcItem.StorageType = toBox
			srcItem.Storage = toSlot
			dstItem.StorageType = fromBox
			dstItem.Storage = fromSlot
			h.sendCellChange(client, toBox, toSlot, srcItem, fromLoc)
			h.sendCellChange(client, fromBox, fromSlot, dstItem, toLoc)
		}
	} else {
		// Intercambio de celdas (Exchange)
		_ = h.charRepo.ExchangeItemStorage(srcItem.ItemKey, toBox, toSlot, dstItem.ItemKey, fromBox, fromSlot)
		srcItem.StorageType = toBox
		srcItem.Storage = toSlot
		dstItem.StorageType = fromBox
		dstItem.Storage = fromSlot
		h.sendCellChange(client, toBox, toSlot, srcItem, fromLoc)
		h.sendCellChange(client, fromBox, fromSlot, dstItem, toLoc)
	}

	if fromBox == 8 || toBox == 8 {
		h.syncPlayerEquipsAndStats(player)
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_RELOC_ACK)
	ack.WriteUint16(0x0241) // Éxito oficial en Zone.exe (0x0241)
	_ = client.SendPacket(ack)
}

// handleItemSplit implementa sp_NC_ITEM_SPLIT_REQ (0x300D) y gds_NC_ITEMDB_SPLITSUC_ACK (Zone.exe VA 0x539740):
//   - Lee from (uint16), to (uint16), lot (uint8)
//   - Reduce el stack origen en `lot` y crea un nuevo ítem en `to` con `Flags = lot`
//   - Envía 0x3001 para ambas celdas y PROTO_NC_ITEM_SPLIT_ACK (0x300E) con err = 0x0441 (1089)
func (h *Handler) handleItemSplit(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	fromLoc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	toLoc, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}
	splitLot, err := pkt.Reader.ReadByte()
	if err != nil || splitLot == 0 {
		return
	}

	fromBox := uint8(fromLoc >> 10)
	fromSlot := int16(fromLoc & 0x03FF)
	toBox := uint8(toLoc >> 10)
	toSlot := int16(toLoc & 0x03FF)

	srcItem := h.findItemAt(player.Char.ID, fromBox, fromSlot)
	dstItem := h.findItemAt(player.Char.ID, toBox, toSlot)
	if srcItem == nil || dstItem != nil || srcItem.Flags <= int32(splitLot) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SPLIT_ACK)
		ack.WriteUint16(0x0442) // Split fail
		_ = client.SendPacket(ack)
		return
	}

	srcItem.Flags -= int32(splitLot)
	_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(srcItem.Flags))

	newItem := &database.Item{StorageType: toBox, Owner: player.Char.ID, Storage: int16(toSlot), ItemID: srcItem.ItemID, Flags: int32(int(splitLot))}
	err = h.charRepo.CreateItem(newItem)
	if err != nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SPLIT_ACK)
		ack.WriteUint16(0x0442)
		_ = client.SendPacket(ack)
		return
	}


	h.sendCellChange(client, fromBox, fromSlot, srcItem, 0xFFFF)
	h.sendCellChange(client, toBox, toSlot, newItem, 0xFFFF)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SPLIT_ACK)
	ack.WriteUint16(0x0441) // Éxito oficial en Zone.exe (0x0441)
	_ = client.SendPacket(ack)
}

// handleItemEquip implementa sp_NC_ITEM_EQUIP_REQ (0x300F) y gds_NC_ITEMDB_EQUIPSUC_ACK (Zone.exe VA 0x52CBA0):
//   - Lee invenslot (uint8) del inventario (Box 9)
//   - Determina el slot de equipamiento desde ItemInfo.shn (info.Equip)
//   - Mueve o intercambia el ítem con Box 8 (Equip), recalcula stats y transmite cambio visual al mapa
//   - Envía PROTO_NC_ITEM_EQUIP_ACK (0x3011) con err = 0x0281 (641)
func (h *Handler) handleItemEquip(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	invenSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	srcItem := h.findItemAt(player.Char.ID, 9, int16(invenSlot))
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
		ack.WriteUint16(0x0282)
		_ = client.SendPacket(ack)
		return
	}

	equipSlot := slotWeapon
	if h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(srcItem.ItemID); ok && info != nil {
			if info.Equip > 0 && info.Equip != 255 {
				equipSlot = byte(info.Equip)
			}
		}
	}

	h.equipItemToSlot(player, srcItem, invenSlot, equipSlot)
}

// handleItemRingEquip implementa sp_NC_ITEM_RINGEQUIP_REQ (0x3010):
//   - Lee invenslot (uint8) y ringslot (uint8)
func (h *Handler) handleItemRingEquip(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	invenSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	ringSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	srcItem := h.findItemAt(player.Char.ID, 9, int16(invenSlot))
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
		ack.WriteUint16(0x0282)
		_ = client.SendPacket(ack)
		return
	}

	h.equipItemToSlot(player, srcItem, invenSlot, ringSlot)
}

func (h *Handler) equipItemToSlot(player *ZonePlayer, srcItem *database.Item, invenSlot byte, equipSlot byte) {
	client := player.Client
	if !h.canEquipItem(player, srcItem) {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
		ack.WriteUint16(0x0282)
		_ = client.SendPacket(ack)
		return
	}

	dstEquip := h.findItemAt(player.Char.ID, 8, int16(equipSlot))
	invenLoc := (uint16(9) << 10) | (uint16(invenSlot) & 0x03FF)
	equipLoc := (uint16(8) << 10) | (uint16(equipSlot) & 0x03FF)

	if dstEquip == nil {
		_ = h.charRepo.SetItemStorage(srcItem.ItemKey, 8, int16(equipSlot))
		srcItem.StorageType = 8
		srcItem.Storage = int16(equipSlot)

		h.sendCellChange(client, 9, int16(invenSlot), nil, 0xFFFF)
		h.sendCellChange(client, 8, int16(equipSlot), srcItem, 0xFFFF)
	} else {
		_ = h.charRepo.ExchangeItemStorage(srcItem.ItemKey, 8, int16(equipSlot), dstEquip.ItemKey, 9, int16(invenSlot))
		srcItem.StorageType = 8
		srcItem.Storage = int16(equipSlot)
		dstEquip.StorageType = 9
		dstEquip.Storage = int16(invenSlot)

		h.sendCellChange(client, 8, int16(equipSlot), srcItem, invenLoc)
		h.sendCellChange(client, 9, int16(invenSlot), dstEquip, equipLoc)
	}

	h.syncPlayerEquipsAndStats(player)
	h.broadcastEquipAppearance(player, equipSlot, srcItem.ItemID, 0)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
	ack.WriteUint16(0x0281) // Éxito oficial en Zone.exe (0x0281)
	_ = client.SendPacket(ack)
}

func (h *Handler) canEquipItem(player *ZonePlayer, item *database.Item) bool {
	if player == nil || player.Char == nil || item == nil {
		return false
	}
	if h.itemStore == nil {
		if player.Client != nil {
			log.Printf("[Zone][%s] cannot validate equipment item %d: official item data is not loaded",
				player.Client.RemoteAddr(), item.ItemID)
		}
		return false
	}
	info, ok := h.itemStore.GetItemInfo(item.ItemID)
	if !ok || info == nil || uint32(player.Char.Level) < info.DemandLv {
		return false
	}
	allowed, known := h.itemStore.CanUseItemByClassID(item.ItemID, player.Char.Job)
	return known && allowed
}

func sortInventoryItems(items []*database.Item, store *item.Store) {
	sort.SliceStable(items, func(i, j int) bool {
		leftOrder, leftGrade, leftID := int(^uint(0)>>1), uint32(0), items[i].ItemID
		rightOrder, rightGrade, rightID := int(^uint(0)>>1), uint32(0), items[j].ItemID
		if store != nil {
			if info, ok := store.GetItemInfo(items[i].ItemID); ok && info != nil {
				if info.SortOrder >= 0 {
					leftOrder = info.SortOrder
				}
				leftGrade = info.Grade
			}
			if info, ok := store.GetItemInfo(items[j].ItemID); ok && info != nil {
				if info.SortOrder >= 0 {
					rightOrder = info.SortOrder
				}
				rightGrade = info.Grade
			}
		}
		if leftOrder != rightOrder {
			return leftOrder < rightOrder
		}
		if leftGrade != rightGrade {
			return leftGrade < rightGrade
		}
		if leftID != rightID {
			return leftID < rightID
		}
		return items[i].Storage < items[j].Storage
	})
}

// handleItemUnequip implementa sp_NC_ITEM_UNEQUIP_REQ (0x3012) y gds_NC_ITEMDB_UNEQUIPSUC_ACK (Zone.exe VA 0x53B0A0):
//   - Lee equipslot (uint8) e invenslot (uint8)
//   - Mueve el ítem de Box 8 a Box 9, recalcula stats, transmite PROTO_NC_BRIEFINFO_UNEQUIP_CMD (0x1C04)
//   - Envía PROTO_NC_ITEM_EQUIP_ACK (0x3011) con err = 0x0281 (641)
func (h *Handler) handleItemUnequip(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	equipSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	invenSlot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}

	srcEquip := h.findItemAt(player.Char.ID, 8, int16(equipSlot))
	if srcEquip == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
		ack.WriteUint16(0x0282)
		_ = client.SendPacket(ack)
		return
	}

	targetSlot := int16(invenSlot)
	if h.findItemAt(player.Char.ID, 9, targetSlot) != nil {
		if freeSlot, found := h.findFreeInvenSlot(player.Char.ID); found {
			targetSlot = freeSlot
		} else {
			ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
			ack.WriteUint16(0x0282)
			_ = client.SendPacket(ack)
			return
		}
	}

	_ = h.charRepo.SetItemStorage(srcEquip.ItemKey, 9, targetSlot)
	srcEquip.StorageType = 9
	srcEquip.Storage = targetSlot

	h.sendCellChange(client, 9, targetSlot, srcEquip, 0xFFFF)
	h.sendCellChange(client, 8, int16(equipSlot), nil, 0xFFFF)

	h.syncPlayerEquipsAndStats(player)
	h.broadcastUnequipAppearance(player, equipSlot)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_EQUIP_ACK)
	ack.WriteUint16(0x0281) // Éxito oficial en Zone.exe (0x0281)
	_ = client.SendPacket(ack)
}

// handleItemDrop implementa sp_NC_ITEM_DROP_REQ (0x3007) y gds_NC_ITEMDB_DROPALLSUC_ACK (Zone.exe VA 0x52AB80):
//   - Lee slot (uint8), lot (uint8), loc (SHINE_XY_TYPE: x int32, y int32)
//   - Elimina o reduce el ítem en Box 9, envía 0x3001, spawnea el ítem en el suelo del mapa (0x1C0A)
//   - Envía PROTO_NC_ITEM_DROP_ACK (0x3008) con err = 0x0301
func (h *Handler) handleItemDrop(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || player.CurrentMap == nil || h.charRepo == nil {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	dropLot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	dropX, err := pkt.Reader.ReadInt32()
	if err != nil {
		dropX = int32(player.Char.XPos)
	}
	dropY, err := pkt.Reader.ReadInt32()
	if err != nil {
		dropY = int32(player.Char.YPos)
	}

	srcItem := h.findItemAt(player.Char.ID, 9, int16(slot))
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DROP_ACK)
		ack.WriteUint16(0x0304) // Drop fail en Zone.exe
		_ = client.SendPacket(ack)
		return
	}

	curLot := int(srcItem.Flags)
	if curLot <= 0 {
		curLot = 1
	}
	actualDrop := int(dropLot)
	if actualDrop <= 0 || actualDrop >= curLot {
		actualDrop = curLot
		_ = h.charRepo.DeleteItem(srcItem.ItemKey)
		h.sendCellChange(client, 9, int16(slot), nil, 0xFFFF)
	} else {
		srcItem.Flags = int32(curLot - actualDrop)
		_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(srcItem.Flags))
		h.sendCellChange(client, 9, int16(slot), srcItem, 0xFFFF)
	}

	player.CurrentMap.SpawnDroppedItem(srcItem.ItemID, uint32(actualDrop), int(srcItem.Flags), dropX, dropY, player.MapObjectID, player.Char.ID)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_DROP_ACK)
	ack.WriteUint16(0x0301) // Éxito oficial en Zone.exe (0x0301)
	_ = client.SendPacket(ack)
}

// handleItemPick implementa sp_NC_ITEM_PICK_REQ (0x3009) y so_ply_itempickresult (Zone.exe VA 0x5308D0):
//   - Lee itemhandle (uint16)
//   - Consume el ítem del mapa (enviando 0x1C0E PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
//   - Inserta en Box 9 del jugador, envía 0x3001 PROTO_NC_ITEM_CELLCHANGE_CMD
//   - Envía PROTO_NC_ITEM_PICK_ACK (0x300A, 10 bytes): itemid(2), lot(4), error(2 = 0x0341), itemhandle(2)
func (h *Handler) handleItemPick(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || player.CurrentMap == nil || h.charRepo == nil {
		return
	}

	itemHandle, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	freeSlot, hasFree := h.findFreeInvenSlot(player.Char.ID)
	if !hasFree {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_ACK)
		ack.WriteUint16(0xFFFF)
		ack.WriteUint32(0)
		ack.WriteUint16(0x0344) // Inventory full
		ack.WriteUint16(itemHandle)
		_ = client.SendPacket(ack)
		return
	}

	// Verificar propiedad temporal del drop antes de consumirlo del suelo
	if ent := player.CurrentMap.GetEntity(itemHandle); ent != nil {
		if pre, ok := ent.(*ZoneDroppedItem); ok && pre != nil {
			if pre.OwnerCharID != 0 && !pre.AllCanLoot && pre.OwnerCharID != player.Char.ID {
				if pre.OwnerExpiresAt.IsZero() || time.Now().Before(pre.OwnerExpiresAt) {
					ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_ACK)
					ack.WriteUint16(0xFFFF)
					ack.WriteUint32(0)
					ack.WriteUint16(0x0343) // Ownership locked by another player
					ack.WriteUint16(itemHandle)
					_ = client.SendPacket(ack)
					return
				}
			}
		}
	}

	dropped, found := player.CurrentMap.ConsumeDroppedItem(itemHandle)
	if !found || dropped == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_ACK)
		ack.WriteUint16(0xFFFF)
		ack.WriteUint32(0)
		ack.WriteUint16(0x0342) // Item not found on field
		ack.WriteUint16(itemHandle)
		_ = client.SendPacket(ack)
		return
	}

	lot := int(dropped.Lot)
	if lot <= 0 {
		lot = 1
	}

	if dropped.IsCenDrop {
		beforeMoney := player.Char.TCharacterState.Money
		_ = h.charRepo.AddCharMoney(player.Char.ID, int64(lot))
		if player.Char.TCharacterState.Money == beforeMoney {
			player.Char.TCharacterState.Money += int64(lot)
		}
		cenCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CENCHANGE_CMD)
		cenCmd.WriteUint64(uint64(player.Char.TCharacterState.Money))
		_ = client.SendPacket(cenCmd)

		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_ACK)
		ack.WriteUint16(dropped.ItemID)
		ack.WriteUint32(uint32(lot))
		ack.WriteUint16(0)
		ack.WriteUint16(itemHandle)
		_ = client.SendPacket(ack)
		return
	}

	newItem := &database.Item{StorageType: 9, Owner: player.Char.ID, Storage: int16(freeSlot), ItemID: dropped.ItemID, Flags: int32(lot)}
	err = h.charRepo.CreateItem(newItem)
	if err != nil {
		return
	}
	for _, opt := range dropped.StatOptions {
		_ = h.charRepo.SetItemOption(newItem.ItemKey, int16(opt.Type), int64(opt.Value))
	}

	h.sendCellChangeWithOptions(client, 9, freeSlot, newItem, dropped.StatOptions, 0xFFFF)

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_PICK_ACK)
	ack.WriteUint16(dropped.ItemID) // itemid (2B)
	ack.WriteUint32(dropped.Lot)    // lot (4B)
	ack.WriteUint16(0x0341)         // error = 0x0341 (éxito oficial en so_ply_itempickresult)
	ack.WriteUint16(itemHandle)     // itemhandle (2B)
	_ = client.SendPacket(ack)
}

// handleItemUse implementa sp_NC_ITEM_USE_REQ (0x3015), so_NC_ITEMDB_USE_ACK (VA 0x53FB50) y so_ply_ItemEffect (VA 0x59C730):
//   - Lee slot (uint8) e invenType (uint8)
//   - Aplica los efectos de ItemUseEffect.shn (0 = HP, 1 = SP, 2 = Full HP/SP) y envía 0x240E / 0x240F
//   - Consume 1 unidad del stack (o destruye el ítem si Lot == 1) y envía 0x3001 PROTO_NC_ITEM_CELLCHANGE_CMD
//   - Transmite PROTO_NC_ITEM_SOMEONEUSE_CMD (0x3029, 12B) al mapa y envía PROTO_NC_ITEM_USE_ACK (0x3016, 5B) con error = 0x0700
func (h *Handler) handleItemUse(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		return
	}
	invenType, err := pkt.Reader.ReadByte()
	if err != nil {
		invenType = 9
	}
	if invenType == 0 {
		invenType = 9
	}

	srcItem := h.findItemAt(player.Char.ID, invenType, int16(slot))
	if srcItem == nil {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_USE_ACK)
		ack.WriteUint16(0x0701) // Use fail
		ack.WriteUint16(0xFFFF)
		ack.WriteByte(invenType)
		_ = client.SendPacket(ack)
		return
	}

	usedItemID := srcItem.ItemID
	usedItemName := ""
	if h.itemStore != nil {
		if info, ok := h.itemStore.GetItemInfo(usedItemID); ok && info != nil {
			usedItemName = info.InxName
		}
	}

	// Si el ítem usado es un cofre de tesoro / MysteryVault, abrir recompensa oficial
	if _, opened := h.openTreasureChestItem(player, srcItem); opened {
		ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_USE_ACK)
		ack.WriteUint16(0x0700)
		ack.WriteUint16(usedItemID)
		ack.WriteByte(invenType)
		_ = client.SendPacket(ack)
		return
	}

	// Aplicar efectos desde ItemUseEffect.shn (so_ply_ItemEffect VA 0x59C730)
	hpChanged := false
	spChanged := false
	if h.itemStore != nil {
		if eff, hasEff := h.itemStore.GetItemUseEffect(usedItemName); hasEff && eff != nil {
			applyEff := func(effectType uint32, v uint16) {
				val := int(v)
				if val == 0 {
					return
				}
				switch effectType {
				case 0: // Recover HP
					player.Char.CurHP += val
					if player.Char.CurHP > player.MaxHP {
						player.Char.CurHP = player.MaxHP
					}
					hpChanged = true
				case 1: // Recover SP
					player.Char.CurSP += val
					if player.Char.CurSP > player.MaxSP {
						player.Char.CurSP = player.MaxSP
					}
					spChanged = true
				case 2: // Recover Full HP & SP (99999)
					player.Char.CurHP = player.MaxHP
					player.Char.CurSP = player.MaxSP
					hpChanged = true
					spChanged = true
				}
			}
			applyEff(eff.UseEffectA, eff.UseValueA)
			applyEff(eff.UseEffectB, eff.UseValueB)
			applyEff(eff.UseEffectC, eff.UseValueC)
		}
	}

	if hpChanged {
		hpPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
		hpPkt.WriteUint32(uint32(player.Char.CurHP))
		hpPkt.WriteUint16(0) // hpchangeorder
		_ = client.SendPacket(hpPkt)
	}
	if spChanged {
		spPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
		spPkt.WriteUint32(uint32(player.Char.CurSP))
		_ = client.SendPacket(spPkt)
	}
	if hpChanged || spChanged {
		_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)
	}

	// Consumir 1 unidad del ítem
	curLot := srcItem.Flags
	if curLot <= 1 {
		_ = h.charRepo.DeleteItem(srcItem.ItemKey)
		h.sendCellChange(client, invenType, int16(slot), nil, 0xFFFF)
	} else {
		srcItem.Flags = curLot - 1
		_ = h.charRepo.UpdateItemFlags(srcItem.ItemKey, int(srcItem.Flags))
		h.sendCellChange(client, invenType, int16(slot), srcItem, 0xFFFF)
	}

	// Broadcast PROTO_NC_ITEM_SOMEONEUSE_CMD (0x3029, 12B)
	if player.CurrentMap != nil {
		someUse := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_SOMEONEUSE_CMD)
		someUse.WriteUint16(player.MapObjectID)
		someUse.WriteInt32(int32(player.Char.XPos))
		someUse.WriteInt32(int32(player.Char.YPos))
		someUse.WriteUint16(usedItemID)
		player.CurrentMap.Broadcast(someUse, 0)
	}

	// Enviar PROTO_NC_ITEM_USE_ACK (0x3016, 5B)
	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_USE_ACK)
	ack.WriteUint16(0x0700)   // error = 0x0700 (éxito oficial en so_ply_ItemUseResultPacket)
	ack.WriteUint16(usedItemID)
	ack.WriteByte(invenType)
	_ = client.SendPacket(ack)
}

// handleAutoArrangeInven implementa sp_NC_ITEM_AUTO_ARRANGE_INVEN_REQ (0x304A, Zone.exe VA 0x51C360):
//   - Ordena los ítems de Box 9 por la categoría oficial ItemSort.shn y los criterios de ItemInfo.shn
//   - Compacta las posiciones desde el slot 0 y envía el ACK oficial 0x30A8
func (h *Handler) handleAutoArrangeInven(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	items, err := h.charRepo.GetCharacterItems(player.Char.ID)
	if err != nil {
		log.Printf("[Zone][%s] auto-arrange: failed to load inventory for character %d: %v",
			client.RemoteAddr(), player.Char.ID, err)
		return
	}
	inventory := make([]*database.Item, 0, len(items))
	for _, it := range items {
		if it != nil && it.StorageType == 9 {
			inventory = append(inventory, it)
		}
	}
	sortInventoryItems(inventory, h.itemStore)

	itemKeys := make([]int64, len(inventory))
	oldSlots := make([]int16, len(inventory))
	for i, it := range inventory {
		itemKeys[i] = it.ItemKey
		oldSlots[i] = it.Storage
	}
	if err := h.charRepo.ArrangeInventory(player.Char.ID, itemKeys); err != nil {
		log.Printf("[Zone][%s] auto-arrange: failed to persist inventory for character %d: %v",
			client.RemoteAddr(), player.Char.ID, err)
		return
	}

	for _, slot := range oldSlots {
		h.sendCellChange(client, 9, slot, nil, 0xFFFF)
	}
	for slot, it := range inventory {
		it.Storage = int16(slot)
		it.StorageType = 9
		h.sendCellChange(client, 9, int16(slot), it, 0xFFFF)
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_AUTO_ARRANGE_INVEN_ACK)
	ack.WriteUint16(0x30A8) // Código oficial en sp_NC_ITEM_AUTO_ARRANGE_INVEN_REQ (VA 0x51C360)
	_ = client.SendPacket(ack)
}

// handleSoulStoneHPBuy implementa sp_NC_SOULSTONE_HP_BUY_REQ (0x5001) y gds_NC_ITEMDB_BUYHPSOULSTONESUC_ACK (VA 0x52D820):
//   - Lee number (uint16)
//   - Descuenta el precio en Cen (10 Cen por piedra en Nivel 1 o según Level), incrementa StonesHP
//   - Envía PROTO_NC_SOULSTONE_HP_BUY_ACK (0x5003), PROTO_NC_ITEM_BUY_ACK (0x3004 = 0x0201) y PROTO_NC_CHAR_CENCHANGE_CMD (0x1033)
func (h *Handler) handleSoulStoneHPBuy(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	buyNum, err := pkt.Reader.ReadUint16()
	if err != nil || buyNum == 0 {
		return
	}

	maxStone := int(player.MaxStonesHP)
	if player.LevelParam.MAXSoulHP > 0 {
		maxStone = player.LevelParam.MAXSoulHP
	}
	if int(player.StonesHP)+int(buyNum) > maxStone {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0742) // Max stones exceeded
		_ = client.SendPacket(fail)
		return
	}

	// Precio oficial por piedra desde ClassParam (LevelParam.PriceHPStone)
	totalCost := int64(player.LevelParam.PriceHPStone) * int64(buyNum)
	if player.Char.TCharacterState.Money < totalCost {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0743) // Insufficient money
		_ = client.SendPacket(fail)
		return
	}
	if okSub, err := h.charRepo.SubCharMoney(player.Char.ID, totalCost); err != nil || !okSub {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0743)
		_ = client.SendPacket(fail)
		return
	}

	player.Char.TCharacterState.Money -= totalCost
	player.StonesHP += int16(buyNum)
	player.Char.TCharacterState.HPStones = player.StonesHP

	_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)

	// 1. PROTO_NC_SOULSTONE_HP_BUY_ACK (0x5003, 2B: totalnumber)
	buyAck := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_BUY_ACK)
	buyAck.WriteUint16(uint16(player.StonesHP))
	_ = client.SendPacket(buyAck)

	// 2. PROTO_NC_ITEM_BUY_ACK (0x3004, 2B: err = 0x0201)
	itemBuyAck := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_BUY_ACK)
	itemBuyAck.WriteUint16(0x0201)
	_ = client.SendPacket(itemBuyAck)

	// 3. PROTO_NC_CHAR_CENCHANGE_CMD (0x1033, 8B: cen uint64)
	cenCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CENCHANGE_CMD)
	cenCmd.WriteUint64(uint64(player.Char.TCharacterState.Money))
	_ = client.SendPacket(cenCmd)
}

// handleSoulStoneSPBuy implementa sp_NC_SOULSTONE_SP_BUY_REQ (0x5002) y gds_NC_ITEMDB_BUYSPSOULSTONESUC_ACK (VA 0x52ED00):
func (h *Handler) handleSoulStoneSPBuy(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	buyNum, err := pkt.Reader.ReadUint16()
	if err != nil || buyNum == 0 {
		return
	}

	maxStone := int(player.MaxStonesSP)
	if player.LevelParam.MAXSoulSP > 0 {
		maxStone = player.LevelParam.MAXSoulSP
	}
	if int(player.StonesSP)+int(buyNum) > maxStone {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0742) // Max stones exceeded
		_ = client.SendPacket(fail)
		return
	}

	// Precio oficial por piedra desde ClassParam (LevelParam.PriceSPStone)
	totalCost := int64(player.LevelParam.PriceSPStone) * int64(buyNum)
	if player.Char.TCharacterState.Money < totalCost {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0743) // Insufficient money
		_ = client.SendPacket(fail)
		return
	}
	if okSub, err := h.charRepo.SubCharMoney(player.Char.ID, totalCost); err != nil || !okSub {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_BUYFAIL_ACK)
		fail.WriteUint16(0x0743)
		_ = client.SendPacket(fail)
		return
	}

	player.Char.TCharacterState.Money -= totalCost
	player.StonesSP += int16(buyNum)
	player.Char.TCharacterState.SPStones = player.StonesSP

	_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)

	// 1. PROTO_NC_SOULSTONE_SP_BUY_ACK (0x5004, 2B: totalnumber)
	buyAck := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_SP_BUY_ACK)
	buyAck.WriteUint16(uint16(player.StonesSP))
	_ = client.SendPacket(buyAck)

	// 2. PROTO_NC_ITEM_BUY_ACK (0x3004, 2B: err = 0x0201)
	itemBuyAck := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_BUY_ACK)
	itemBuyAck.WriteUint16(0x0201)
	_ = client.SendPacket(itemBuyAck)

	// 3. PROTO_NC_CHAR_CENCHANGE_CMD (0x1033, 8B: cen uint64)
	cenCmd := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CENCHANGE_CMD)
	cenCmd.WriteUint64(uint64(player.Char.TCharacterState.Money))
	_ = client.SendPacket(cenCmd)
}

// handleSoulStoneHPUse implementa sp_NC_SOULSTONE_HP_USE_REQ (0x5007, VA 0x588A40):
//   - Verifica StonesHP > 0, decrementa StonesHP en 1, recupera LevelParam.SoulHP de HP
//   - Transmite PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD (0x500B, 2B) al mapa
//   - Envía PROTO_NC_SOULSTONE_HP_USESUC_ACK (0x5008, 0B) y PROTO_NC_BAT_HPCHANGE_CMD (0x240E, 6B)
func (h *Handler) handleSoulStoneHPUse(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	if player.State == 3 || player.Char.CurHP <= 0 || player.StonesHP <= 0 {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_USEFAIL_ACK)
		_ = client.SendPacket(fail)
		return
	}

	player.StonesHP--
	player.Char.TCharacterState.HPStones = player.StonesHP

	recoverHP := player.LevelParam.SoulHP
	player.Char.CurHP += recoverHP
	if player.Char.CurHP > player.MaxHP {
		player.Char.CurHP = player.MaxHP
	}

	_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)

	if player.CurrentMap != nil {
		someCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD)
		someCmd.WriteUint16(player.MapObjectID)
		player.CurrentMap.Broadcast(someCmd, 0)
	}

	sucAck := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_USESUC_ACK)
	_ = client.SendPacket(sucAck)

	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(uint32(player.Char.CurHP))
	hpCmd.WriteUint16(0) // hpchangeorder
	_ = client.SendPacket(hpCmd)
}

// handleSoulStoneSPUse implementa sp_NC_SOULSTONE_SP_USE_REQ (0x5009, VA 0x589180):
//   - Verifica StonesSP > 0, decrementa StonesSP en 1, recupera LevelParam.SoulSP de SP
//   - Transmite PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD (0x500C, 2B) al mapa
//   - Envía PROTO_NC_SOULSTONE_SP_USESUC_ACK (0x500A, 0B) y PROTO_NC_BAT_SPCHANGE_CMD (0x240F, 4B)
func (h *Handler) handleSoulStoneSPUse(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil || h.charRepo == nil {
		return
	}

	if player.State == 3 || player.Char.CurHP <= 0 || player.StonesSP <= 0 {
		fail := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_USEFAIL_ACK)
		_ = client.SendPacket(fail)
		return
	}

	player.StonesSP--
	player.Char.TCharacterState.SPStones = player.StonesSP

	recoverSP := player.LevelParam.SoulSP
	player.Char.CurSP += recoverSP
	if player.Char.CurSP > player.MaxSP {
		player.Char.CurSP = player.MaxSP
	}

	_ = h.charRepo.SaveCharStats(player.Char.ID, int16(player.Char.CurHP), int16(player.Char.CurSP), player.StonesHP, player.StonesSP)

	if player.CurrentMap != nil {
		someCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD)
		someCmd.WriteUint16(player.MapObjectID)
		player.CurrentMap.Broadcast(someCmd, 0)
	}

	sucAck := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_SP_USESUC_ACK)
	_ = client.SendPacket(sucAck)

	spCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
	spCmd.WriteUint32(uint32(player.Char.CurSP))
	_ = client.SendPacket(spCmd)
}


const slotWeapon byte = 12 // ItemEquipEnum::ITEMEQUIP_RIGHTHAND (Slot_Weapon)
