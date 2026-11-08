package world

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Handler despacha los paquetes de red del World Server (WorldManager + Character unificados).
type Handler struct {
	worldID       byte
	charRepo      database.CharacterRepository
	accountRepo   database.AccountRepository
	shopRepo      database.AccountShopRepository
	sessionRepo   database.LoginSessionRepository
	logRepo       database.AccountLogRepository
	transfers     *database.TransferManager
	dataStore      *WorldDataStore
	partyMgr       *PartyManager
	friendMgr      *FriendManager
	guildMgr       *GuildManager
	kqMgr          *KQManager
	raidMgr        *RaidManager
	holyPromiseMgr *HolyPromiseManager
	zoneIP         string
	zonePort       uint16
	activeClients  map[int]*network.Client
	mu             sync.Mutex
}

// NewHandler inicializa un nuevo despachador de paquetes para WorldService.
func NewHandler(
	charRepo database.CharacterRepository,
	transfers *database.TransferManager,
	dataStore *WorldDataStore,
	zoneIP string,
	zonePort uint16,
) *Handler {
	if zoneIP == "" {
		zoneIP = "127.0.0.1"
	}
	if zonePort == 0 {
		zonePort = 9210
	}
	pm := NewPartyManager(charRepo)
	fm := NewFriendManager(charRepo, pm)
	gm := NewGuildManager(charRepo)
	km := NewKQManager(dataStore, charRepo)
	rm := NewRaidManager(charRepo, pm)
	hpm := NewHolyPromiseManager(charRepo)
	return &Handler{
		charRepo:       charRepo,
		transfers:      transfers,
		dataStore:      dataStore,
		partyMgr:       pm,
		friendMgr:      fm,
		guildMgr:       gm,
		kqMgr:          km,
		raidMgr:        rm,
		holyPromiseMgr: hpm,
		zoneIP:         zoneIP,
		zonePort:       zonePort,
		activeClients:  make(map[int]*network.Client),
	}
}

// SetAccountRepo vincula los repositorios de cuenta, logs y sesiones.
func (h *Handler) SetAccountRepo(repo database.AccountRepository) {
	h.accountRepo = repo
	if s, ok := repo.(database.LoginSessionRepository); ok {
		h.sessionRepo = s
	}
	if l, ok := repo.(database.AccountLogRepository); ok {
		h.logRepo = l
	}
	if sh, ok := repo.(database.AccountShopRepository); ok {
		h.shopRepo = sh
	}
}

// SetWorldID asigna el ID de mundo correspondiente a este WorldServer.
func (h *Handler) SetWorldID(id byte) {
	h.worldID = id
}

// ActivePlayersCount retorna la cantidad de clientes activos en WorldService.
func (h *Handler) ActivePlayersCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.activeClients)
}

// OnClientDisconnect limpia la sesión activa en el mapa de concurrencia.
func (h *Handler) OnClientDisconnect(client *network.Client) {
	if client.AccountID > 0 {
		h.mu.Lock()
		if h.activeClients[client.AccountID] == client {
			delete(h.activeClients, client.AccountID)
		}
		h.mu.Unlock()
	}
}

// HandlePacket despacha paquetes basándose exclusivamente en OpCodes oficiales PROTO_NC_*.
func (h *Handler) HandlePacket(client *network.Client, pkt *protocol.Packet) {
	switch pkt.OpCode {
	case protocol.PROTO_NC_MISC_HEARTBEAT_ACK,
		protocol.PROTO_NC_MISC_HEARTBEAT_REQ:
		h.handleHeartbeat(client, pkt)

	case protocol.PROTO_NC_MISC_GAMETIME_REQ:
		h.handleGameTime(client, pkt)

	case protocol.PROTO_NC_MISC_CLIENT_DEBUG_MSG_CMD:
		h.handleClientDebug(client, pkt)

	case protocol.PROTO_NC_MISC_PINGTEST_TOOL_WM_CLIENT_ZONE_DB:
		h.handlePing(client, pkt)

	case protocol.PROTO_NC_MISC_APEX_CLIENT_DATA_CMD,
		protocol.PROTO_NC_MISC_APEX_CLIENT_CHCSTART_CMD:
		h.handleApexServerData(client, pkt)

	case protocol.PROTO_NC_MISC_WEB_KEY_REQ:
		h.handleWebKey(client, pkt)

	case protocol.PROTO_NC_MISC_CS_ACK:
		h.handleCSAck(client, pkt)

	case protocol.PROTO_NC_MISC_HACK_SCAN_STORE_CMD:
		h.handleHackScan(client, pkt)

	case protocol.PROTO_NC_MISC_CLIENT_LOADING_BUG_DETECT_CMD:
		h.handleClientLoadingBug(client, pkt)

	case protocol.PROTO_NC_MISC_SET_CHAT_BLOCK_SPAM_FILTER_CMD:
		h.handleSpamFilterSet(client, pkt)

	case protocol.PROTO_NC_MISC_SPAMMER_REPORT_REQ:
		h.handleSpammerReport(client, pkt)

	case protocol.PROTO_NC_MISC_SPAMMER_CHAT_BAN_REQ:
		h.handleChatBan(client, pkt)

	case protocol.PROTO_NC_MISC_ITEMSHOP_URL_REQ:
		h.handleItemShopURL(client, pkt)

	case protocol.PROTO_NC_USER_LOGINWORLD_REQ:
		h.handleWorldClientKey(client, pkt)

	case protocol.PROTO_NC_USER_NORMALLOGOUT_CMD:
		h.handleNormalLogout(client, pkt)

	case protocol.PROTO_NC_USER_AVATAR_LIST_REQ:
		h.handleAvatarList(client, pkt)

	case protocol.PROTO_NC_USER_WILL_WORLD_SELECT_REQ:
		h.handleBackToWorldSelect(client, pkt)

	case protocol.PROTO_NC_USER_CONNECTION_CHANGE_CHAR_ID_REQ:
		h.handleChangeCharID(client, pkt)

	case protocol.PROTO_NC_USER_RETURN_CHECK_REQ:
		h.handleReturnCheck(client, pkt)

	case protocol.PROTO_NC_USER_SET_RETURN_CMD:
		h.handleSetReturn(client, pkt)

	case protocol.PROTO_NC_USER_LOGOUT_LAST_DAY_REQ:
		h.handleLogoutLastDay(client, pkt)

	case protocol.PROTO_NC_MISC_WEB_DB_KEY_REQ:
		h.handleWebDBKey(client, pkt)

	case protocol.PROTO_NC_MISC_CONNECTFROMWHERE_DB_REQ:
		h.handleConnectFromWhere(client, pkt)

	case protocol.PROTO_NC_MISC_EVENT_HIT3_ADD_CASH:
		h.handleEventHit3AddCash(client, pkt)

	case protocol.PROTO_NC_MISC_EVENT_DONE_MUNSANG_WM2ACC:
		h.handleEventDoneMunsang(client, pkt)

	case protocol.PROTO_NC_MISC_EVENT_L20_DB_REQ:
		h.handleEventL20(client, pkt)

	case protocol.PROTO_NC_MISC_EVENT_LONG_TIME_PLAY_CMD:
		h.handleEventLongTimePlay(client, pkt)

	case protocol.PROTO_NC_CHAR_LOGIN_REQ:
		h.handleCharSelect(client, pkt)

	case protocol.PROTO_NC_CHAR_SET_STYLE_REQ:
		h.handleSetStyle(client, pkt)

	case protocol.PROTO_NC_CHAR_SET_STYLE_GET_INFO_REQ:
		h.handleSetStyleGetInfo(client, pkt)

	case protocol.PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD:
		h.handleNewbieGuideViewSet(client, pkt)

	case protocol.PROTO_NC_CHAR_EVENT_ATTENDANCE_REWARD_REQ:
		h.handleAttendanceReward(client, pkt)

	case protocol.PROTO_NC_CHAR_TUTORIAL_POPUP_ACK:
		h.handleTutorialPopup(client, pkt)

	case protocol.PROTO_NC_AVATAR_CREATE_REQ:
		h.handleCreateCharacter(client, pkt)

	case protocol.PROTO_NC_AVATAR_ERASE_REQ:
		h.handleDeleteCharacter(client, pkt)

	case protocol.PROTO_NC_AVATAR_RENAME_REQ:
		h.handleRenameCharacter(client, pkt)

	case protocol.PROTO_NC_KQ_LIST_REFRESH_REQ:
		h.handleGotIngame(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_ALL_REQ:
		h.handleGetAllSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_REQ:
		h.handleGetQuickBar(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ:
		h.handleGetQuickBarState(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_VIDEO_REQ:
		h.handleGetVideoSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_SOUND_REQ:
		h.handleGetSoundSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_GAME_REQ:
		h.handleGetGameSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_REQ:
		h.handleGetClientSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_REQ:
		h.handleGetShortCuts(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_SHORTCUTSIZE_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_VIDEO_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_SOUND_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_GAME_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_WINDOWPOS_CMD,
		protocol.PROTO_NC_CHAR_OPTION_SET_KEYMAPPING_CMD:
		h.handleSaveSettings(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ:
		h.handleImproveSetShortcut(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_IMPROVE_SET_KEYMAP_REQ:
		h.handleImproveSetKeymap(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_IMPROVE_SET_GAMEOPTION_REQ:
		h.handleImproveSetGameOption(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_IMPROVE_INIT_KEYMAP_REQ:
		h.handleImproveInitKeymap(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_IMPROVE_INIT_GAMEOPTION_REQ:
		h.handleImproveInitGameOption(client, pkt)

	case protocol.PROTO_NC_CHAR_OPTION_LOGIN_BLOCKDATA_ERR:
		h.handleLoginBlockDataErr(client, pkt)

	// Subsistema Oficial de Chat (Act Chat - Dept: 8)
	case protocol.PROTO_NC_ACT_WHISPER_REQ:
		h.handleWhisper(client, pkt)

	case protocol.PROTO_NC_ACT_NOTICE_REQ:
		h.handleNotice(client, pkt)

	case protocol.PROTO_NC_ACT_PARTYCHAT_REQ:
		h.handlePartyChat(client, pkt)

	// Subsistema Oficial de Grupos (Party - Dept: 14)
	case protocol.PROTO_NC_PARTY_JOIN_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.partyMgr.InviteMember(client, targetName)

	case protocol.PROTO_NC_PARTY_JOINPROPOSE_ALLOW_ACK:
		h.partyMgr.AcceptInvite(client)

	case protocol.PROTO_NC_PARTY_JOINPROPOSE_REJECT_ACK:
		h.partyMgr.RejectInvite(client)

	case protocol.PROTO_NC_PARTY_LEAVE_REQ:
		h.partyMgr.LeaveParty(client)

	case protocol.PROTO_NC_PARTY_KICKOFF_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.partyMgr.KickMember(client, targetName)

	case protocol.PROTO_NC_PARTY_CHANGEMASTER_REQ:
		newMaster, _ := pkt.Reader.ReadString(20)
		h.partyMgr.ChangeMaster(client, newMaster)

	case protocol.PROTO_NC_PARTY_ITEM_LOOTING_SET:
		lootType, _ := pkt.Reader.ReadByte()
		itemGrade, _ := pkt.Reader.ReadByte()
		h.partyMgr.SetLooting(client, lootType, itemGrade)

	case protocol.PROTO_NC_PARTY_FINDER_ADD_REQ:
		purpose, _ := pkt.Reader.ReadString(pkt.Reader.Remaining())
		h.partyMgr.AddFinderAd(client, purpose)

	case protocol.PROTO_NC_PARTY_FINDER_DELETE_REQ:
		h.partyMgr.DeleteFinderAd(client)

	case protocol.PROTO_NC_PARTY_FINDER_LIST_REQ:
		h.partyMgr.ListFinderAds(client)

	case protocol.PROTO_NC_PARTY_SET_LOOTER_REQ:
		looterName, _ := pkt.Reader.ReadString(20)
		h.partyMgr.SetLooter(client, looterName)

	// Subsistema Oficial de Amigos (Friend - Dept: 21)
	case protocol.PROTO_NC_FRIEND_SET_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.friendMgr.RequestFriend(client, targetName)

	case protocol.PROTO_NC_FRIEND_SET_CONFIRM_ACK:
		h.friendMgr.AcceptFriend(client)

	case protocol.PROTO_NC_FRIEND_DEL_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.friendMgr.DeleteFriend(client, targetName)

	case protocol.PROTO_NC_FRIEND_POINT_REQ:
		h.friendMgr.SendFriendPoints(client)

	case protocol.PROTO_NC_FRIEND_FIND_FRIENDS_REQ:
		h.friendMgr.FindFriends(client, pkt)

	// Subsistema Oficial de Gremios (Guild - Dept: 29)
	case protocol.PROTO_NC_GUILD_LIST_REQ:
		h.guildMgr.HandleGuildList(client, 0)

	case protocol.PROTO_NC_GUILD_MAKE_REQ:
		name, _ := pkt.Reader.ReadString(16)
		password, _ := pkt.Reader.ReadString(12)
		guildType, _ := pkt.Reader.ReadByte()
		h.guildMgr.HandleGuildMake(client, name, password, guildType)

	case protocol.PROTO_NC_GUILD_MY_GUILD_INFO_REQ:
		h.guildMgr.HandleMyGuildInfo(client)

	case protocol.PROTO_NC_GUILD_MEMBER_LIST_REQ:
		h.guildMgr.HandleMemberList(client)

	case protocol.PROTO_NC_GUILD_MEMBER_INVITE_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.guildMgr.HandleInvite(client, targetName)

	case protocol.PROTO_NC_GUILD_MEMBER_JOIN_ACK:
		guildName, _ := pkt.Reader.ReadString(16)
		joinAck, _ := pkt.Reader.ReadByte()
		h.guildMgr.HandleJoinAck(client, guildName, joinAck)

	case protocol.PROTO_NC_GUILD_MEMBER_LEAVE_REQ:
		h.guildMgr.HandleLeave(client)

	case protocol.PROTO_NC_GUILD_MEMBER_VANISH_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.guildMgr.HandleVanish(client, targetName)

	case protocol.PROTO_NC_GUILD_MEMBER_GRADE_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		grade, _ := pkt.Reader.ReadByte()
		h.guildMgr.HandleGrade(client, targetName, grade)

	case protocol.PROTO_NC_GUILD_NOTIFY_REQ:
		nLen, _ := pkt.Reader.ReadUint16()
		notifyBytes, _ := pkt.Reader.ReadBytes(int(nLen))
		h.guildMgr.HandleNotify(client, string(notifyBytes))

	case protocol.PROTO_NC_GUILD_DISMISS_REQ:
		h.guildMgr.HandleDismiss(client)

	case protocol.PROTO_NC_GUILD_CHAT_REQ:
		_, _ = pkt.Reader.ReadByte() // itemLinkCount
		lenByte, _ := pkt.Reader.ReadByte()
		msgBytes, _ := pkt.Reader.ReadBytes(int(lenByte))
		h.guildMgr.HandleChat(client, string(msgBytes))

	case protocol.PROTO_NC_GUILD_WAR_REQ:
		targetName, _ := pkt.Reader.ReadString(16)
		h.guildMgr.HandleWarReq(client, targetName)

	case protocol.PROTO_NC_GUILD_WAR_ABLE_LIST_REQ:
		h.guildMgr.HandleWarAbleListReq(client)

	case protocol.PROTO_NC_GUILD_WAR_LIST_REQ:
		h.guildMgr.HandleWarListReq(client)

	case protocol.PROTO_NC_GUILD_NAME_REQ:
		newName, _ := pkt.Reader.ReadString(16)
		h.guildMgr.HandleGuildNameReq(client, newName)

	case protocol.PROTO_NC_GUILD_TOURNAMENT_JOIN_REQ:
		gtNo, _ := pkt.Reader.ReadUint16()
		h.guildMgr.HandleTournamentJoinReq(client, int(gtNo))

	case protocol.PROTO_NC_GUILD_TOURNAMENT_LIST_REQ:
		h.guildMgr.HandleTournamentListReq(client)

	case protocol.PROTO_NC_GUILD_MY_GUILD_TOKEN_INFO_REQ:
		h.guildMgr.HandleTokenInfoReq(client)

	case protocol.PROTO_NC_GUILD_MY_GUILD_TOURNAMENT_MATCH_TIME_REQ:
		h.guildMgr.HandleMyTournamentMatchTimeReq(client)

	case protocol.PROTO_NC_GUILD_TOURNAMENT_JOIN_LIST_REQ:
		gtNo, _ := pkt.Reader.ReadUint16()
		h.guildMgr.HandleTournamentJoinListReq(client, int(gtNo))

	case protocol.PROTO_NC_GUILD_HISTORY_LIST_REQ:
		h.guildMgr.HandleHistoryListReq(client)

	case protocol.PROTO_NC_GUILD_EMBLEM_CHECK_AVAILABILITY_REQ:
		h.guildMgr.HandleEmblemCheckAvailabilityReq(client)

	case protocol.PROTO_NC_GUILD_EMBLEM_SAVE_REQ:
		data, _ := pkt.Reader.ReadBytes(pkt.Reader.Remaining())
		h.guildMgr.HandleEmblemSaveReq(client, data)

	// Subsistema Oficial de Academia de Gremios (Guild Academy - Dept: 38)
	case protocol.PROTO_NC_GUILD_ACADEMY_LIST_REQ:
		h.guildMgr.HandleAcademyList(client, 0)

	case protocol.PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_REQ:
		h.guildMgr.HandleAcademyMyInfo(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_MY_ACADEMY_RANK_INFO_REQ:
		h.guildMgr.HandleAcademyRankInfo(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_REQ:
		h.guildMgr.HandleAcademyMemberList(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_REQ:
		guildName, _ := pkt.Reader.ReadString(16)
		h.guildMgr.HandleAcademyJoin(client, guildName)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_REQ:
		h.guildMgr.HandleAcademyLeave(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.guildMgr.HandleAcademyVanish(client, targetName)

	case protocol.PROTO_NC_GUILD_ACADEMY_MASTER_TELEPORT_REQ:
		h.guildMgr.HandleAcademyMasterTeleport(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		min, _ := pkt.Reader.ReadInt32()
		h.guildMgr.HandleAcademyChatBan(client, targetName, int(min))

	case protocol.PROTO_NC_GUILD_ACADEMY_NOTIFY_REQ:
		nLen, _ := pkt.Reader.ReadUint16()
		notifyBytes, _ := pkt.Reader.ReadBytes(int(nLen))
		h.guildMgr.HandleAcademyNotify(client, string(notifyBytes))

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_INVITE_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.guildMgr.HandleAcademyGuildInvite(client, targetName)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_GUILD_JOIN_ACK:
		isAccept, _ := pkt.Reader.ReadByte()
		h.guildMgr.HandleAcademyGuildJoinAck(client, isAccept)

	case protocol.PROTO_NC_GUILD_ACADEMY_GET_REWARD_MONEY_REQ:
		h.guildMgr.HandleAcademyGetRewardMoney(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_GET_REWARD_ITEM_REQ:
		h.guildMgr.HandleAcademyGetRewardItem(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_SET_REWARD_MONEY_REQ:
		money, _ := pkt.Reader.ReadUint64()
		h.guildMgr.HandleAcademySetRewardMoney(client, money)

	case protocol.PROTO_NC_GUILD_ACADEMY_SET_REWARD_ITEM_REQ:
		itemID, _ := pkt.Reader.ReadUint16()
		h.guildMgr.HandleAcademySetRewardItem(client, itemID)

	case protocol.PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_MONEY_REQ:
		h.guildMgr.HandleAcademyClearRewardMoney(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_CLEAR_REWARD_ITEM_REQ:
		h.guildMgr.HandleAcademyClearRewardItem(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_GET_GRADUATE_REQ:
		h.guildMgr.HandleAcademyGetGraduate(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_GET_RANKING_LIST_REQ:
		h.guildMgr.HandleAcademyGetRankingList(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_SET_MASTER_REQ:
		newMaster, _ := pkt.Reader.ReadString(20)
		h.guildMgr.HandleAcademySetMaster(client, newMaster)

	case protocol.PROTO_NC_GUILD_ACADEMY_MEMBER_CHAT_BAN_CANCEL_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.guildMgr.HandleAcademyChatBanCancel(client, targetName)

	case protocol.PROTO_NC_GUILD_ACADEMY_HISTORY_LIST_REQ:
		h.guildMgr.HandleAcademyHistoryList(client)

	case protocol.PROTO_NC_GUILD_ACADEMY_CHAT_REQ:
		_, _ = pkt.Reader.ReadByte() // itemLinkCount
		lenByte, _ := pkt.Reader.ReadByte()
		msgBytes, _ := pkt.Reader.ReadBytes(int(lenByte))
		h.guildMgr.HandleAcademyChat(client, string(msgBytes))

	// Subsistema Oficial de Kingdom Quests (KQ - Dept: 22)
	case protocol.PROTO_NC_KQ_LIST_REQ:
		h.kqMgr.HandleList(client)

	case protocol.PROTO_NC_KQ_STATUS_REQ:
		handle, _ := pkt.Reader.ReadUint32()
		h.kqMgr.HandleStatus(client, handle)

	case protocol.PROTO_NC_KQ_JOIN_REQ:
		handle, _ := pkt.Reader.ReadUint32()
		var level, class byte = 1, 1
		ch, err := h.charRepo.GetCharacterByID(client.CharacterID)
		if err == nil && ch != nil {
			level = ch.Level
			class = ch.Job
		}
		h.kqMgr.HandleJoin(client, handle, level, class)

	case protocol.PROTO_NC_KQ_JOIN_CANCEL_REQ:
		handle, _ := pkt.Reader.ReadUint32()
		h.kqMgr.HandleJoinCancel(client, handle)

	case protocol.PROTO_NC_KQ_JOIN_LIST_REQ:
		handle, _ := pkt.Reader.ReadUint32()
		h.kqMgr.HandleJoinList(client, handle)

	case protocol.PROTO_NC_KQ_TEAM_SELECT_REQ:
		teamType, _ := pkt.Reader.ReadByte()
		h.kqMgr.HandleTeamSelect(client, teamType)

	case protocol.PROTO_NC_KQ_SCHEDULE_REQ:
		h.kqMgr.HandleSchedule(client)

	case protocol.PROTO_NC_KQ_VOTE_START_CHECK_REQ:
		h.kqMgr.HandleVoteStartCheck(client)

	case protocol.PROTO_NC_KQ_VOTE_START_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.kqMgr.HandleVoteStart(client, targetName)

	case protocol.PROTO_NC_KQ_VOTE_VOTING_REQ:
		voteVal, _ := pkt.Reader.ReadByte()
		h.kqMgr.HandleVoteCast(client, voteVal != 0)

	// Subsistema Oficial de Prisión (Prison - Dept: 31)
	case protocol.PROTO_NC_PRISON_ADD_GM_REQ:
		h.handlePrisonSendIn(client, pkt)

	case protocol.PROTO_NC_PRISON_END_REQ:
		h.handlePrisonTakeOut(client, pkt)

	case protocol.PROTO_NC_PRISON_ALTER_GM_REQ:
		h.handlePrisonAlterTime(client, pkt)

	case protocol.PROTO_NC_PRISON_GET_REQ:
		h.handlePrisonConfirm(client, pkt)

	// Subsistema Oficial de Reportes (Report - Dept: 32)
	case protocol.PROTO_NC_REPORT_ADD_REQ:
		h.handleReportAdd(client, pkt)

	case protocol.PROTO_NC_REPORT_CANCEL_REQ:
		h.handleReportCancel(client, pkt)

	case protocol.PROTO_NC_REPORT_GET_REQ:
		h.handleReportGet(client, pkt)

	// Subsistema Oficial de Promesa Sagrada (Holy Promise - Dept: 37)
	case protocol.PROTO_NC_HOLY_PROMISE_SET_UP_REQ:
		h.handleHolyPromiseSetUp(client, pkt)

	case protocol.PROTO_NC_HOLY_PROMISE_SET_UP_CONFIRM_ACK:
		h.handleHolyPromiseSetUpConfirm(client, pkt)

	case protocol.PROTO_NC_HOLY_PROMISE_DEL_UP_REQ:
		h.handleHolyPromiseDelUp(client, pkt)

	case protocol.PROTO_NC_HOLY_PROMISE_DEL_DOWN_REQ:
		h.handleHolyPromiseDelDown(client, pkt)

	// Subsistema Oficial de Mazmorras Instanciadas / BattleRank (Dept: 41)
	case protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGEDEALT_REQ:
		h.handleBattleRankClearTime(client, pkt)

	case protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_DAMAGETAKEN_REQ:
		h.handleBattleRankScore(client, pkt)

	case protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_HEALINGDONE_REQ:
		h.handleBattleRankPoint(client, pkt)

	case protocol.PROTO_NC_INSTANCE_DUNGEON_RANK_TAB_MYRANK_REQ:
		h.handleBattleRankMyScore(client, pkt)

	// Subsistema Oficial de Restricciones de Chat (Chat Restrict - Dept: 42)
	case protocol.PROTO_NC_CHAT_RESTRICT_ADD_REQ:
		h.handleChatRestrictAdd(client, pkt)

	case protocol.PROTO_NC_CHAT_RESTRICT_DEL_REQ:
		h.handleChatRestrictDel(client, pkt)

	case protocol.PROTO_NC_CHAT_RESTRICT_DEL_ALL_REQ:
		h.handleChatRestrictDelAll(client, pkt)

	// Subsistema Oficial de Bandas (Raid - Dept: 44)
	case protocol.PROTO_NC_RAID_MAKE_REQ:
		h.raidMgr.HandleMake(client)

	case protocol.PROTO_NC_RAID_LOOTING_TYPE_REQ:
		grade, _ := pkt.Reader.ReadByte()
		lootType, _ := pkt.Reader.ReadByte()
		h.raidMgr.HandleLootingType(client, grade, lootType)

	case protocol.PROTO_NC_RAID_INVITATION_REQ:
		targetName, _ := pkt.Reader.ReadString(20)
		h.raidMgr.HandleInvitation(client, targetName)

	case protocol.PROTO_NC_RAID_INVITATION_ANSWER_REQ:
		isAllow, _ := pkt.Reader.ReadByte()
		masterName, _ := pkt.Reader.ReadString(20)
		h.raidMgr.HandleInvitationAnswer(client, isAllow, masterName)

	case protocol.PROTO_NC_RAID_LEAVE_REQ:
		h.raidMgr.HandleLeave(client)

	case protocol.PROTO_NC_RAID_VANISH_REQ:
		vanishName, _ := pkt.Reader.ReadString(20)
		h.raidMgr.HandleVanish(client, vanishName)

	case protocol.PROTO_NC_RAID_SLOT_MOVE_REQ:
		prePos, _ := pkt.Reader.ReadUint32()
		postPos, _ := pkt.Reader.ReadUint32()
		preChar, _ := pkt.Reader.ReadString(20)
		postChar, _ := pkt.Reader.ReadString(20)
		h.raidMgr.HandleSlotMove(client, prePos, postPos, preChar, postChar)

	case protocol.PROTO_NC_RAID_SET_LEADER_REQ:
		newLeader, _ := pkt.Reader.ReadString(20)
		h.raidMgr.HandleSetLeader(client, newLeader)

	case protocol.PROTO_NC_RAID_SET_SUB_LEADER_REQ:
		subLeader, _ := pkt.Reader.ReadString(20)
		h.raidMgr.HandleSetSubLeader(client, subLeader)

	case protocol.PROTO_NC_RAID_DEL_SUB_LEADER_REQ:
		delSubLeader, _ := pkt.Reader.ReadString(20)
		h.raidMgr.HandleDelSubLeader(client, delSubLeader)

	case protocol.PROTO_NC_RAID_SET_LOOTER_REQ:
		looterName, _ := pkt.Reader.ReadString(20)
		h.raidMgr.HandleSetLooter(client, looterName)

	case protocol.PROTO_NC_RAID_CHAT_REQ:
		_, _ = pkt.Reader.ReadByte()
		lenByte, _ := pkt.Reader.ReadByte()
		chatBytes, _ := pkt.Reader.ReadBytes(int(lenByte))
		h.raidMgr.HandleChat(client, string(chatBytes))

	case protocol.PROTO_NC_RAID_WARNING_REQ:
		_, _ = pkt.Reader.ReadByte()
		lenByte, _ := pkt.Reader.ReadByte()
		warningBytes, _ := pkt.Reader.ReadBytes(int(lenByte))
		h.raidMgr.HandleWarning(client, string(warningBytes))

	default:
		log.Printf("[World][%s] Unhandled OpCode: 0x%04X (Dept=%d, Cmd=%d)",
			client.RemoteAddr(), pkt.OpCode, pkt.Header, pkt.Type)
	}
}

// handleGameTime responde con la hora del servidor (PROTO_NC_MISC_GAMETIME_ACK - 0x080E).
func (h *Handler) handleGameTime(client *network.Client, pkt *protocol.Packet) {
	now := time.Now()
	resp := protocol.New_PROTO_NC_MISC_GAMETIME_ACK(byte(now.Hour()), byte(now.Minute()), byte(now.Second()))
	_ = client.SendPacket(resp)
}

// handleWorldClientKey valida el token de traspaso proveniente de LoginService
// y autentica la sesión del cliente en el WorldServer.
// Soporta la estructura oficial PROTO_NC_USER_LOGINWORLD_REQ (320 bytes payload):
//   - char user[256] /* Name256Byte */
//   - unsigned short validate_new[32] (64 bytes)
func (h *Handler) handleWorldClientKey(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[World][%s] Received PROTO_NC_USER_LOGINWORLD_REQ (len %d)",
		client.RemoteAddr(), len(pkt.Raw))

	var username string
	var token64 [64]byte

	if pkt.Reader.Remaining() >= 320 {
		uBytes, _ := pkt.Reader.ReadBytes(256)
		username = strings.TrimRight(string(uBytes), "\x00")
		username = strings.TrimSpace(username)

		tBytes, _ := pkt.Reader.ReadBytes(64)
		copy(token64[:], tBytes)
	} else if pkt.Reader.Remaining() >= 64 {
		tBytes, _ := pkt.Reader.ReadBytes(64)
		copy(token64[:], tBytes)
	} else if pkt.Reader.Remaining() >= 18+32 {
		_ = pkt.Reader.Skip(18)
		tBytes, _ := pkt.Reader.ReadBytes(32)
		copy(token64[:32], tBytes)
	}

	var authenticated bool
	var userNo int
	var admin byte

	// 1. Validar token oficial de 64 bytes contra MySQL (usp_Session_ConsumeWorldToken)
	if h.sessionRepo != nil {
		uNo, uID, err := h.sessionRepo.ConsumeWorldToken(h.worldID, token64, client.RemoteIP())
		if err == nil && uNo > 0 {
			authenticated = true
			userNo = uNo
			if username == "" {
				username = uID
			}
			if h.accountRepo != nil {
				acc, err := h.accountRepo.GetAccountByUsername(username)
				if err == nil && acc != nil {
					admin = acc.Admin
				}
			}
		}
	}

	// 2. Fallback de soporte para tests en memoria (TransferManager)
	if !authenticated && h.transfers != nil {
		tokenStr := strings.TrimRight(string(token64[:32]), "\x00")
		tokenStr = strings.TrimSpace(tokenStr)
		transfer, ok := h.transfers.ConsumeTransfer(tokenStr, token64[:], username, h.worldID)
		if ok && transfer != nil {
			authenticated = true
			userNo = transfer.AccountID
			username = transfer.Username
			admin = transfer.Admin
		}
	}

	if !authenticated {
		log.Printf("[World][%s] World token verification failed for user '%s'", client.RemoteAddr(), username)
		h.sendError(client, protocol.ErrInvalidCredentials)
		client.Close()
		return
	}

	// 3. Expulsión limpia de sesión previa duplicada (KICKOFFFROMWORLD_CMD / CONNECTCUT_CMD)
	h.mu.Lock()
	if prevClient, exists := h.activeClients[userNo]; exists && prevClient != nil && prevClient != client {
		log.Printf("[World][%s] Sesión previa activa detectada para UserNo %d, expulsando de WorldServer",
			client.RemoteAddr(), userNo)
		_ = prevClient.SendPacket(protocol.New_PROTO_NC_USER_KICKOFFFROMWORLD_CMD(uint32(userNo)))
		_ = prevClient.SendPacket(protocol.New_PROTO_NC_USER_CONNECTCUT_CMD(1))
		prevClient.Close()
		delete(h.activeClients, userNo)
	}
	h.activeClients[userNo] = client
	h.mu.Unlock()

	client.AccountID = userNo
	client.Username = username
	client.Admin = admin
	client.IsAuthenticated = true

	log.Printf("[World][%s] Authenticated user '%s' (AccountID: %d)",
		client.RemoteAddr(), client.Username, client.AccountID)

	// Enviar lista de personajes disponibles
	h.sendCharacterList(client)

	// Notificar si la cuenta puede saltarse el tutorial según usp_Tutorial_getIsSkip
	if progRepo := h.getProgRepo(); progRepo != nil {
		if canSkip, err := progRepo.GetTutorialIsSkip(client.AccountID); err == nil && canSkip {
			_ = client.SendPacket(protocol.New_PROTO_NC_USER_TUTORIAL_CAN_SKIP_CMD(true))
		}
	}
}

// sendCharacterList despacha la lista oficial de avatares (PROTO_NC_USER_LOGINWORLD_ACK - 0x0C14).
func (h *Handler) sendCharacterList(client *network.Client) {
	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err != nil {
		log.Printf("[World][%s] Error loading characters for account %d: %v",
			client.RemoteAddr(), client.AccountID, err)
		h.sendError(client, protocol.ErrDatabaseError)
		return
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINWORLD_ACK)
	if client.RandomID == 0 {
		client.RandomID = uint16(rand.Intn(65534) + 1)
	}
	pw.WriteUint16(client.RandomID)
	pw.WriteByte(byte(len(chars)))

	for _, ch := range chars {
		WriteBasicCharInfo(ch, pw)
	}

	log.Printf("[World][%s] Sent %d character(s) to '%s'",
		client.RemoteAddr(), len(chars), client.Username)

	_ = client.SendPacket(pw)
}

// handleCreateCharacter procesa la creación de personaje validando contra BadNameFilter.shn
// y asignando valores base desde default_character_data.json.
func (h *Handler) handleCreateCharacter(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		log.Printf("[World][%s] Unauthenticated CreateCharacter request rejected", client.RemoteAddr())
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil {
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	name, err := pkt.Reader.ReadString(20)
	if err != nil {
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}
	name = strings.TrimSpace(name)

	jobGender, _ := pkt.Reader.ReadByte()
	hair, _ := pkt.Reader.ReadByte()
	color, _ := pkt.Reader.ReadByte()
	style, _ := pkt.Reader.ReadByte()

	if slot > 5 {
		log.Printf("[World][%s] Invalid slot %d requested", client.RemoteAddr(), slot)
		h.sendCharCreationError(client, protocol.CreateCharError_ErrorInMaxSlot)
		return
	}

	if len(name) < 3 || len(name) > 16 {
		log.Printf("[World][%s] Invalid character name length '%s'", client.RemoteAddr(), name)
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	isMale := ((jobGender >> 7) & 0x01) == 1
	classID := byte((jobGender >> 2) & 0x1F)
	job := protocol.Job(classID)

	switch job {
	case protocol.Job_Fighter, protocol.Job_Cleric, protocol.Job_Archer, protocol.Job_Mage, protocol.Job_Trickster, protocol.Job_Crusader:
		// Clase válida soportada por el cliente oficial
	default:
		log.Printf("[World][%s] Invalid job ID %d requested", client.RemoteAddr(), classID)
		h.sendCharCreationError(client, protocol.CreateCharError_WrongClass)
		return
	}

	// 1. Validación de apariencia oficial (FaceInfo, HairInfo, HairColorInfo y límites de clase)
	if h.dataStore != nil && !h.dataStore.ValidateAppearance(classID, hair, color, style) {
		log.Printf("[World][%s] Character appearance rejected by limits (Class=%d, Hair=%d, Color=%d, Face=%d)",
			client.RemoteAddr(), classID, hair, color, style)
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	// 2. Validación contra BadNameFilter.shn oficial
	if h.dataStore != nil && h.dataStore.BadNames != nil && h.dataStore.BadNames.IsBadName(name) {
		log.Printf("[World][%s] Character name '%s' rejected by BadNameFilter", client.RemoteAddr(), name)
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	// 3. Validación de unicidad de nombre en la base de datos
	if existing, _ := h.charRepo.GetCharacterByName(name); existing != nil {
		log.Printf("[World][%s] Character name '%s' already exists in database", client.RemoteAddr(), name)
		h.sendCharCreationError(client, protocol.CreateCharError_NameTaken)
		return
	}

	// 4. Cargar parámetros base desde default_character_data.json o tutorial_character_data.json
	isTutorial := client.TutorialChosen && !client.TutorialSkip
	var classData DefaultClassData
	var exists bool

	if isTutorial && h.dataStore != nil && h.dataStore.TutorialData != nil {
		classData, exists = h.dataStore.TutorialData.GetTutorialData(classID)
	}
	if !exists && h.dataStore != nil && h.dataStore.DefaultData != nil {
		classData, exists = h.dataStore.DefaultData.GetClassData(classID)
	}

	startMap := ""
	if h.dataStore != nil && h.dataStore.Fields != nil {
		if f, ok := h.dataStore.Fields.GetFieldBySerial(1); ok {
			startMap = f.MapID
		}
	}
	startX := 6888
	startY := 7045
	hp := 100
	sp := 100
	hpSoul := int16(15)
	spSoul := int16(7)
	money := int64(0)
	initLevel := byte(1)
	initExp := int64(0)
	var str, con, dex, intel, men int

	if exists {
		startMap = classData.StartMap
		if isTutorial && h.dataStore != nil && h.dataStore.Fields != nil && h.dataStore.Fields.TutorialMap() != nil {
			startMap = h.dataStore.Fields.TutorialMap().MapID
		}
		if classData.StartX != 0 {
			startX = classData.StartX
		}
		if classData.StartY != 0 {
			startY = classData.StartY
		}
		if classData.HP != 0 {
			hp = classData.HP
		}
		if classData.SP != 0 {
			sp = classData.SP
		}
		if classData.HPSoul != 0 {
			hpSoul = classData.HPSoul
		}
		if classData.SPSoul != 0 {
			spSoul = classData.SPSoul
		}
		money = classData.Money
		if classData.InitLevel > 0 {
			initLevel = classData.InitLevel
		}
		initExp = classData.InitExp

		// Si es clase avanzada (ej. Sentinel Lv 60), leer stats exactos de ClassParam/<Class>.json
		if classData.InitLevel > 1 && h.dataStore != nil && h.dataStore.ClassParams != nil {
			if param, ok := h.dataStore.ClassParams.GetParam(classData.ClassName, int(classData.InitLevel)); ok {
				hp = param.MaxHP
				sp = param.MaxSP
				hpSoul = int16(param.SoulHP)
				spSoul = int16(param.SoulSP)
				str = param.Strength
				con = param.Constitution
				dex = param.Dexterity
				intel = param.Intelligence
				men = param.MentalPower
			}
		}
	}

	race := byte(1)
	if classID == 21 || classID == 26 {
		race = 2
	}

	newChar := &database.Character{
		AccountID: client.AccountID,
		Name:      name,
		Slot:      slot,
		Job:       classID,
		Race:      race,
		Level:     initLevel,
		Male:      isMale,
		Hair:      hair,
		HairColor: color,
		Face:      style,
		TownName:  startMap,
		XPos:      startX,
		YPos:      startY,
		CurHP:     hp,
		CurSP:     sp,
		Equips:    make(map[byte]uint16),
	}
	newChar.TCharacter.Strength = str
	newChar.TCharacter.Constitute = con
	newChar.TCharacter.Dexterity = dex
	newChar.TCharacter.Intelligence = intel
	newChar.TCharacter.MentalPower = men
	newChar.TCharacterState.Money = money
	newChar.TCharacterState.HPStones = hpSoul
	newChar.TCharacterState.SPStones = spSoul
	newChar.TCharacterState.Exp = initExp

	// Distribuir ítems iniciales: los de slot de equipo van a Equips, los demás a Items (mochila)
	var invenSlot int16 = 0
	itemsToDistribute := classData.StarterItems
	if isTutorial && len(classData.TutorialItems) > 0 {
		itemsToDistribute = classData.TutorialItems
	}

	for _, sItem := range itemsToDistribute {
		var equipSlot byte
		if h.dataStore != nil {
			equipSlot = h.dataStore.GetItemEquipSlot(sItem.ItemID)
		}
		if equipSlot > 0 && newChar.Equips[equipSlot] == 0 {
			newChar.Equips[equipSlot] = sItem.ItemID
		} else {
			amount := sItem.Amount
			if amount <= 0 {
				amount = 1
			}
			newChar.Items = append(newChar.Items, &database.Item{
				StorageType: 9, // Inventario / Mochila (Box 9)
				Storage:     invenSlot,
				ItemID:      sItem.ItemID,
				Flags:       int32(amount),
			})
			invenSlot++
		}
	}

	// Fallback para entornos de pruebas unitarias sin ItemInfo.shn cargado
	if len(newChar.Equips) == 0 {
		if job == protocol.Job_Archer {
			newChar.Equips[EquipLeftHand] = 1250
		} else if job == protocol.Job_Trickster {
			newChar.Equips[EquipRightHand] = 57363
		} else if job == protocol.Job_Crusader {
			newChar.Equips[EquipRightHand] = 62923
		} else if job == protocol.Job_Cleric {
			newChar.Equips[EquipRightHand] = 750
		} else if job == protocol.Job_Mage {
			newChar.Equips[EquipRightHand] = 1750
		} else {
			newChar.Equips[EquipRightHand] = 250
		}
	}

	for _, sc := range classData.Shortcuts {
		newChar.Shortcuts = append(newChar.Shortcuts, &database.Shortcut{
			SlotNo:   sc.Slot,
			CodeType: 1,
			CodeNo:   int16(sc.CodeNo),
			Value:    int32(sc.Value),
		})
	}

	for _, sk := range classData.StarterSkills {
		newChar.Skills = append(newChar.Skills, &database.CharacterSkill{
			SkillNo:    uint16(sk),
			SkillLevel: 1,
		})
	}

	for _, qID := range classData.StarterQuests {
		newChar.Quests = append(newChar.Quests, &database.QuestData{
			CharNo:  newChar.ID,
			QuestNo: qID,
			Status:  1, // Activa
		})
	}

	if err := h.charRepo.CreateCharacter(newChar); err != nil {
		log.Printf("[World][%s] Failed to create character '%s': %v", client.RemoteAddr(), name, err)
		h.sendCharCreationError(client, protocol.CreateCharError_FailedToCreate)
		return
	}

	log.Printf("[World][%s] Created character '%s' (ID: %d, Job: %d, Slot: %d) for user '%s'",
		client.RemoteAddr(), name, newChar.ID, classID, slot, client.Username)

	// Enviar respuesta oficial de confirmación PROTO_NC_AVATAR_CREATESUCC_ACK (0x1406)
	resp := protocol.New_PROTO_NC_AVATAR_CREATESUCC_ACK(func(pw *protocol.PacketWriter) {
		WriteBasicCharInfo(newChar, pw)
	})
	_ = client.SendPacket(resp)
}

// handleDeleteCharacter ejecuta el soft-delete oficial liberando el slot y nickname.
func (h *Handler) handleDeleteCharacter(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil || slot > 5 {
		return
	}

	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err != nil {
		failResp := protocol.New_PROTO_NC_AVATAR_ERASEFAIL_ACK(0x0001)
		_ = client.SendPacket(failResp)
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
		log.Printf("[World][%s] Character in slot %d not found for deletion", client.RemoteAddr(), slot)
		failResp := protocol.New_PROTO_NC_AVATAR_ERASEFAIL_ACK(0x0001)
		_ = client.SendPacket(failResp)
		return
	}

	// 1. Validar que no sea Maestro de Gremio activo
	if avatarRepo := h.getAvatarRepo(); avatarRepo != nil {
		if isGM, err := avatarRepo.IsActiveGuildMaster(targetChar.ID); err == nil && isGM {
			log.Printf("[World][%s] Cannot delete character '%s' (ID: %d): active Guild Master",
				client.RemoteAddr(), targetChar.Name, targetChar.ID)
			failResp := protocol.New_PROTO_NC_AVATAR_ERASEFAIL_ACK(0x0002) // Guild Master error
			_ = client.SendPacket(failResp)
			return
		}
	}

	if err := h.charRepo.DeleteCharacter(client.AccountID, slot); err != nil {
		log.Printf("[World][%s] Error deleting character '%s' in slot %d: %v", client.RemoteAddr(), targetChar.Name, slot, err)
		failResp := protocol.New_PROTO_NC_AVATAR_ERASEFAIL_ACK(0x0001)
		_ = client.SendPacket(failResp)
		return
	}

	log.Printf("[World][%s] Soft-deleted character in slot %d for user '%s'", client.RemoteAddr(), slot, client.Username)

	resp := protocol.New_PROTO_NC_AVATAR_ERASESUCC_ACK(slot)
	_ = client.SendPacket(resp)
}

// handleBackToWorldSelect procesa la solicitud de volver a la selección de mundos
// generando un One-Time Password criptográfico oficial de 32 bytes (Name8 sOTP)
// y emitiendo PROTO_NC_USER_WILL_WORLD_SELECT_ACK (0x0C34).
func (h *Handler) handleBackToWorldSelect(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated || client.AccountID <= 0 {
		resp := protocol.New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK([32]byte{})
		_ = client.SendPacket(resp)
		return
	}

	var otp [32]byte
	const hexChars = "0123456789abcdef"
	for i := 0; i < 32; i++ {
		otp[i] = hexChars[rand.Intn(len(hexChars))]
	}

	if h.sessionRepo != nil {
		err := h.sessionRepo.CreateOTP(client.AccountID, otp)
		if err != nil {
			log.Printf("[World][%s] Error persistiendo sOTP en tActiveAccountSession para UserNo %d: %v",
				client.RemoteAddr(), client.AccountID, err)
		} else {
			log.Printf("[World][%s] Generado sOTP de 32 bytes para UserNo %d", client.RemoteAddr(), client.AccountID)
		}
	}

	resp := protocol.New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK(otp)
	_ = client.SendPacket(resp)
}

// handleReturnCheck procesa la consulta de usuario que retorna (0x0C60).
func (h *Handler) handleReturnCheck(client *network.Client, pkt *protocol.Packet) {
	var handle uint16
	userNo := uint32(client.AccountID)
	if pkt.Reader.Remaining() >= 6 {
		handle, _ = pkt.Reader.ReadUint16()
		u, err := pkt.Reader.ReadUint32()
		if err == nil && u > 0 {
			userNo = u
		}
	}

	var isReturn byte
	var lastSec int32
	if h.logRepo != nil && userNo > 0 {
		diff, err := h.logRepo.GetDateDiffLastLogout(int(userNo))
		if err == nil && diff >= 30*24*3600 {
			isReturn = 1
			lastSec = int32(diff)
		}
	}

	resp := protocol.New_PROTO_NC_USER_RETURN_CHECK_ACK(handle, userNo, isReturn, lastSec)
	_ = client.SendPacket(resp)
}

// handleSetReturn procesa la confirmación de reclamo de paquete de retorno (0x0C64).
func (h *Handler) handleSetReturn(client *network.Client, pkt *protocol.Packet) {
	userNo := client.AccountID
	if pkt.Reader.Remaining() >= 4 {
		u, err := pkt.Reader.ReadUint32()
		if err == nil && u > 0 {
			userNo = int(u)
		}
	}

	if h.accountRepo != nil && userNo > 0 {
		_ = h.accountRepo.SetReturnUserItem(userNo)
		log.Printf("[World][%s] Paquete de retorno acreditado para UserNo %d", client.RemoteAddr(), userNo)
	}
}

// handleLogoutLastDay procesa la consulta de días desde la última desconexión (0x0C62).
func (h *Handler) handleLogoutLastDay(client *network.Client, pkt *protocol.Packet) {
	var handle uint16
	userNo := uint32(client.AccountID)
	if pkt.Reader.Remaining() >= 6 {
		handle, _ = pkt.Reader.ReadUint16()
		u, err := pkt.Reader.ReadUint32()
		if err == nil && u > 0 {
			userNo = u
		}
	}

	var lastSec int32
	if h.logRepo != nil && userNo > 0 {
		diff, err := h.logRepo.GetDateDiffLastLogout(int(userNo))
		if err == nil {
			lastSec = int32(diff)
		}
	}

	resp := protocol.New_PROTO_NC_USER_LOGOUT_LAST_DAY_ACK(handle, userNo, lastSec)
	_ = client.SendPacket(resp)
}

// handleWebDBKey procesa la solicitud de token SSO para ItemShop Web (0x081D).
func (h *Handler) handleWebDBKey(client *network.Client, pkt *protocol.Packet) {
	var handle uint16
	userNo := uint32(client.AccountID)
	if pkt.Reader.Remaining() >= 6 {
		handle, _ = pkt.Reader.ReadUint16()
		u, _ := pkt.Reader.ReadUint32()
		if u > 0 {
			userNo = u
		}
	}

	var key string
	var errCode uint16 = 0
	if h.accountRepo != nil && userNo > 0 {
		k, err := h.accountRepo.GetOrIssueWebKey(int(userNo), "")
		if err != nil {
			errCode = 1
		} else {
			key = k
		}
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_WEB_DB_KEY_ACK)
	pw.WriteUint16(errCode)
	pw.WriteUint16(handle)
	pw.WriteUint32(userNo)
	pw.WriteString(key, 32)
	_ = client.SendPacket(pw)
}

// handleConnectFromWhere procesa la verificación de origen de conexión y PC Bang (0x082F).
func (h *Handler) handleConnectFromWhere(client *network.Client, pkt *protocol.Packet) {
	var hZone uint16
	var zoneHeader []byte
	var connectIPStr string = client.RemoteIP()

	rem := pkt.Reader.Remaining()
	if rem >= 12 {
		hZone, _ = pkt.Reader.ReadUint16()
		zoneHeader, _ = pkt.Reader.ReadBytes(6)
		ipBytes, _ := pkt.Reader.ReadBytes(4)
		if len(ipBytes) == 4 && (ipBytes[0] != 0 || ipBytes[1] != 0 || ipBytes[2] != 0 || ipBytes[3] != 0) {
			connectIPStr = fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
		}
	} else if rem >= 10 {
		zoneHeader, _ = pkt.Reader.ReadBytes(6)
		ipBytes, _ := pkt.Reader.ReadBytes(4)
		if len(ipBytes) == 4 && (ipBytes[0] != 0 || ipBytes[1] != 0 || ipBytes[2] != 0 || ipBytes[3] != 0) {
			connectIPStr = fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
		}
	}

	if len(zoneHeader) < 6 {
		zoneHeader = make([]byte, 6)
	}

	var isPcbang bool
	if h.accountRepo != nil {
		pc, _ := h.accountRepo.CheckPcbangIP(connectIPStr)
		isPcbang = pc
	}

	var connectFrom uint8 = 0
	if isPcbang {
		connectFrom = 1
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_CONNECTFROMWHERE_DB_ACK)
	pw.WriteUint16(hZone)
	pw.WriteBytes(zoneHeader)
	pw.WriteByte(connectFrom)
	_ = client.SendPacket(pw)
}

// handleEventHit3AddCash procesa la acreditación de cash por evento (0x0820).
func (h *Handler) handleEventHit3AddCash(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	if rem >= 28 {
		userIDStr, _ := pkt.Reader.ReadString(20)
		addCash, _ := pkt.Reader.ReadUint64()
		userIDStr = strings.TrimSpace(userIDStr)
		if h.shopRepo != nil && userIDStr != "" && addCash > 0 {
			_ = h.shopRepo.AddEventCashByUserID(userIDStr, int64(addCash))
		}
	} else if rem >= 8 {
		userNo, _ := pkt.Reader.ReadUint32()
		cash, _ := pkt.Reader.ReadUint32()
		if h.shopRepo != nil && userNo > 0 && cash > 0 {
			_ = h.shopRepo.AddCash(int(userNo), 0, int(cash), "EventHit3")
		}
	}
}

// handleEventDoneMunsang procesa la validación y canje de cupón cultural Munsang (0x0838).
func (h *Handler) handleEventDoneMunsang(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var (
		userNo  uint32
		charNo  uint32
		worldNo uint8
		hZone   uint16
		questID uint16
		nError  uint16
	)

	if rem >= 13 {
		userNo, _ = pkt.Reader.ReadUint32()
		charNo, _ = pkt.Reader.ReadUint32()
		worldNo, _ = pkt.Reader.ReadByte()
		hZone, _ = pkt.Reader.ReadUint16()
		questID, _ = pkt.Reader.ReadUint16()
	} else if rem >= 4 {
		userNo, _ = pkt.Reader.ReadUint32()
	}

	if h.shopRepo != nil && userNo > 0 {
		errCode, err := h.shopRepo.RedeemMunsangEvent(int(userNo), questID)
		if err == nil {
			nError = errCode
		} else {
			nError = 5445 // Código de fallo oficial de Account.exe (0x1545)
		}
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_EVENT_DONE_MUNSANG_ACC2WM)
	pw.WriteUint32(userNo)
	pw.WriteUint32(charNo)
	pw.WriteByte(worldNo)
	pw.WriteUint16(nError)
	pw.WriteUint16(hZone)
	pw.WriteUint16(questID)
	_ = client.SendPacket(pw)
}

// handleEventL20 procesa el registro y verificación del evento Nivel 20 (0x082C).
func (h *Handler) handleEventL20(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	var (
		wmHandle   uint16
		userID     string
		userNo     uint32
		eventNo    uint8 = 1
		worldNo    uint8 = h.worldID
		charNo     uint32
		eventPoint uint32
		nError     uint16
	)

	if rem >= 272 {
		wmHandle, _ = pkt.Reader.ReadUint16()
		uID, _ := pkt.Reader.ReadString(256)
		userID = strings.TrimSpace(uID)
		userNo, _ = pkt.Reader.ReadUint32()
		eventNo, _ = pkt.Reader.ReadByte()
		worldNo, _ = pkt.Reader.ReadByte()
		charNo, _ = pkt.Reader.ReadUint32()
		eventPoint, _ = pkt.Reader.ReadUint32()
	} else if rem >= 4 {
		u, _ := pkt.Reader.ReadUint32()
		userNo = u
	} else if client.AccountID > 0 {
		userNo = uint32(client.AccountID)
	}

	if h.logRepo != nil && userNo > 0 {
		err := h.logRepo.EventL20SetFull(int(userNo), userID, eventNo, worldNo, int(charNo), int(eventPoint))
		if err != nil {
			nError = 1
		}
	}

	pw := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_EVENT_L20_DB_ACK)
	pw.WriteUint16(wmHandle)
	pw.WriteUint32(userNo)
	pw.WriteByte(eventNo)
	pw.WriteUint32(charNo)
	pw.WriteUint16(nError)
	_ = client.SendPacket(pw)
}

// handleEventLongTimePlay procesa el otorgamiento de boletos por tiempo de juego (0x0843).
func (h *Handler) handleEventLongTimePlay(client *network.Client, pkt *protocol.Packet) {
	rem := pkt.Reader.Remaining()
	if rem >= 6 {
		userNo, _ := pkt.Reader.ReadUint32()
		playTime, _ := pkt.Reader.ReadUint16()
		if h.logRepo != nil && userNo > 0 {
			_ = h.logRepo.EventiPadSetTicket(int(userNo), playTime)
		}
	} else if rem >= 4 {
		userNo, _ := pkt.Reader.ReadUint32()
		if h.logRepo != nil && userNo > 0 {
			_ = h.logRepo.EventiPadSetTicket(int(userNo))
		}
	}
}

// handleCharSelect procesa la selección de un personaje y entrega el traspaso oficial hacia ZoneServer.
func (h *Handler) handleCharSelect(client *network.Client, pkt *protocol.Packet) {
	if !client.IsAuthenticated {
		return
	}

	slot, err := pkt.Reader.ReadByte()
	if err != nil || slot > 5 {
		log.Printf("[World][%s] Invalid character slot in CharSelect: %v", client.RemoteAddr(), err)
		return
	}

	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err != nil {
		h.sendError(client, protocol.ErrDatabaseError)
		return
	}

	var selected *database.Character
	for _, ch := range chars {
		if ch.Slot == slot {
			selected = ch
			break
		}
	}

	if selected == nil {
		log.Printf("[World][%s] Character in slot %d not found for account %d",
			client.RemoteAddr(), slot, client.AccountID)
		failResp := protocol.New_PROTO_NC_CHAR_LOGINFAIL_ACK(uint16(protocol.ErrInvalidCredentials))
		_ = client.SendPacket(failResp)
		return
	}

	// 1. Validar si requiere cambio de ID obligatorio (bNeedChangeID)
	if selected.ChangeID != nil && selected.ChangeID.NeedChangeID != 0 {
		log.Printf("[World][%s] Character '%s' (ID: %d) requires name change before login",
			client.RemoteAddr(), selected.Name, selected.ID)
		failResp := protocol.New_PROTO_NC_CHAR_LOGINFAIL_ACK(0x000F) // NeedChangeID
		_ = client.SendPacket(failResp)
		return
	}

	// 2. Validar si el personaje se encuentra encarcelado
	if avatarRepo := h.getAvatarRepo(); avatarRepo != nil {
		if prison, err := avatarRepo.GetPrisonInfo(selected.ID); err == nil && prison != nil && prison.Reason != "" {
			log.Printf("[World][%s] Character '%s' (ID: %d) is in prison: %s",
				client.RemoteAddr(), selected.Name, selected.ID, prison.Reason)
			failResp := protocol.New_PROTO_NC_CHAR_LOGINFAIL_ACK(uint16(protocol.ErrInvalidCredentials))
			_ = client.SendPacket(failResp)
			return
		}
	}

	client.CharacterID = selected.ID
	client.Slot = slot
	h.partyMgr.RegisterOnline(client, selected.ID, selected.Name, selected.Job, selected.Level)

	randomID := client.RandomID
	if randomID == 0 {
		randomID = uint16(rand.Intn(65534) + 1)
		client.RandomID = randomID
	}
	clientIP := client.RemoteIP()

	// 3. Resolver endpoint de zona para el mapa dinámicamente sin fallback hardcodeado
	targetIP := h.zoneIP
	targetPort := h.zonePort
	mapName := selected.TownName
	if mapName == "" && h.dataStore != nil && h.dataStore.Fields != nil {
		if f, ok := h.dataStore.Fields.GetFieldBySerial(int(selected.Map)); ok {
			mapName = f.MapID
		}
	}
	if mapName == "" && h.dataStore != nil && h.dataStore.MapInfo != nil {
		for _, row := range h.dataStore.MapInfo.Rows {
			if row.GetUint16("MapIndex") == selected.Map {
				if s := row.GetString("MapName"); s != "" {
					mapName = s
					break
				}
			}
		}
	}
	var zoneNo uint16 = 1
	if h.dataStore != nil && h.dataStore.Fields != nil {
		if f, ok := h.dataStore.Fields.GetField(mapName); ok {
			if f.Fiesta > 0 {
				zoneNo = uint16(f.Fiesta)
			}
		} else if f, ok := h.dataStore.Fields.GetFieldBySerial(int(selected.Map)); ok {
			if f.Fiesta > 0 {
				zoneNo = uint16(f.Fiesta)
			}
		}
	}
	if topo := h.getTopoRepo(); topo != nil {
		if ip, port, ok, err := topo.ResolveMapEndpoint(h.worldID, zoneNo); err == nil && ok && ip != "" {
			targetIP = ip
			targetPort = port
		}
	}

	// 4. Emitir traspaso y login
	h.transfers.CreateZoneTransfer(
		client.AccountID,
		client.Username,
		selected.Name,
		selected.ID,
		randomID,
		client.Admin,
		clientIP,
	)
	if topo := h.getTopoRepo(); topo != nil {
		_ = topo.IssueZoneTransfer(h.worldID, zoneNo, client.AccountID, client.Username, selected.ID, selected.Name, randomID, client.Admin, clientIP, mapName, selected.XPos, selected.YPos)
	}
	_ = h.charRepo.CharacterLogin(selected.ID)

	log.Printf("[World][%s] Character '%s' (ID: %d) selected. Zone transfer token generated for %s:%d",
		client.RemoteAddr(), selected.Name, selected.ID, targetIP, targetPort)

	// Enviar respuesta oficial PROTO_NC_CHAR_LOGIN_ACK (0x1003) con IP y Puerto de Zone
	resp := protocol.New_PROTO_NC_CHAR_LOGIN_ACK(targetIP, targetPort)
	_ = client.SendPacket(resp)

	// Enviar configuración de interfaz del personaje (SendToClient_CharOptionData @ 0x00431680, DataType == 2 NA2016):
	// 1. 0x7032 (PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD)
	// 2. 0x7033 (PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD)
	// 3. 0x7034 (PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD)
	h.sendCharOptionDataImprove(client, selected.ID)

	// 5. Enviar lista de guías de novato vistas si existen registros
	if progRepo := h.getProgRepo(); progRepo != nil {
		if guides, err := progRepo.GetInGameGuideViewList(selected.ID); err == nil && len(guides) > 0 {
			_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_CLIENT_NEWBIE_GUIDE_VIEW_LIST_CMD(guides))
		}
	}
}

// sendCharOptionDataImprove envía automáticamente los paquetes del sistema mejorado de opciones (DataType == 2)
// justo después de PROTO_NC_CHAR_LOGIN_ACK (0x1003), replicando CWMClientSession::SendToClient_CharOptionData.
func (h *Handler) sendCharOptionDataImprove(client *network.Client, charNo int) {
	progRepo := h.getProgRepo()
	var shortcutEntries []protocol.SHORT_CUT_DATA
	if progRepo != nil {
		if scList, err := progRepo.GetImproveShortCutList(charNo); err == nil && len(scList) > 0 {
			shortcutEntries = make([]protocol.SHORT_CUT_DATA, 0, len(scList))
			for _, sc := range scList {
				shortcutEntries = append(shortcutEntries, protocol.SHORT_CUT_DATA{
					SlotNo: sc.SlotNo,
					CodeNo: sc.CodeNo,
					Value:  sc.Value,
				})
			}
		}
	}

	if len(shortcutEntries) == 0 {
		if shortcuts, err := h.charRepo.GetCharacterShortcuts(charNo); err == nil && len(shortcuts) > 0 {
			for _, sc := range shortcuts {
				if sc == nil {
					continue
				}
				shortcutEntries = append(shortcutEntries, protocol.SHORT_CUT_DATA{
					SlotNo: sc.SlotNo,
					CodeNo: uint16(sc.CodeNo),
					Value:  sc.Value,
				})
			}
		}
	}

	var keyMapEntries []protocol.KEY_MAP_DATA
	if progRepo != nil {
		if kmList, err := progRepo.GetImproveKeyMap(charNo); err == nil && len(kmList) > 0 {
			keyMapEntries = make([]protocol.KEY_MAP_DATA, 0, len(kmList))
			for _, km := range kmList {
				keyMapEntries = append(keyMapEntries, protocol.KEY_MAP_DATA{
					FunctionNo: km.FunctionNo,
					ExtendKey:  km.ExtendKey,
					ASCIICode:  km.ASCIICode,
				})
			}
		}
	}

	var gameOptEntries []protocol.GAME_OPTION_DATA
	if progRepo != nil {
		if goList, err := progRepo.GetImproveGameOptions(charNo); err == nil && len(goList) > 0 {
			gameOptEntries = make([]protocol.GAME_OPTION_DATA, 0, len(goList))
			for _, goOpt := range goList {
				gameOptEntries = append(gameOptEntries, protocol.GAME_OPTION_DATA{
					OptionNo: goOpt.OptionNo,
					Value:    goOpt.Value,
				})
			}
		}
	}

	var rawShortcutData, rawKeyMapData, rawGameOptionData []byte
	if opts, err := h.charRepo.GetCharacterOptions(charNo); err == nil && opts != nil {
		if len(shortcutEntries) == 0 && len(opts.ShortCutData) >= 2 && (len(opts.ShortCutData)-2)%7 == 0 {
			rawShortcutData = opts.ShortCutData
		}
		if len(keyMapEntries) == 0 && len(opts.KeyMapping) >= 2 && (len(opts.KeyMapping)-2)%4 == 0 {
			rawKeyMapData = opts.KeyMapping
		}
		if len(gameOptEntries) == 0 && len(opts.Game) >= 2 && (len(opts.Game)-2)%3 == 0 {
			rawGameOptionData = opts.Game
		}
	}

	// 1. 0x7032: PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
	pktShortcut := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD(shortcutEntries, rawShortcutData)
	_ = client.SendPacket(pktShortcut)

	// 2. 0x7033: PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD
	pktKeyMap := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD(keyMapEntries, rawKeyMapData)
	_ = client.SendPacket(pktKeyMap)

	// 3. 0x7034: PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
	pktGameOpt := protocol.New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD(gameOptEntries, rawGameOptionData)
	_ = client.SendPacket(pktGameOpt)
}

func (h *Handler) sendError(client *network.Client, errCode protocol.ServerError) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_USER_LOGINFAIL_ACK)
	resp.WriteUint16(uint16(errCode))
	_ = client.SendPacket(resp)
}

func (h *Handler) sendCharCreationError(client *network.Client, errCode protocol.CreateCharError) {
	resp := protocol.New_PROTO_NC_AVATAR_CREATEFAIL_ACK(uint16(errCode))
	_ = client.SendPacket(resp)
}

func (h *Handler) getActiveCharNo(client *network.Client) int {
	if client.CharacterID > 0 {
		return client.CharacterID
	}
	chars, err := h.charRepo.GetCharactersByAccount(client.AccountID)
	if err == nil && len(chars) > 0 {
		return chars[0].ID
	}
	return 0
}

func (h *Handler) handleGetQuickBar(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.ShortCutData) > 0 {
			data = opts.ShortCutData
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK(data))
}

func (h *Handler) handleGetQuickBarState(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.ShortCutSize) > 0 {
			data = opts.ShortCutSize
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK(data))
}

func (h *Handler) handleGetGameSettings(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.Game) > 0 {
			data = opts.Game
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_GAME_ACK(data))
}

func (h *Handler) handleGetClientSettings(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.WindowsPos) > 0 {
			data = opts.WindowsPos
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK(data))
}

func (h *Handler) handleGetShortCuts(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	var data []byte
	if charNo > 0 {
		opts, err := h.charRepo.GetCharacterOptions(charNo)
		if err == nil && opts != nil && len(opts.KeyMapping) > 0 {
			data = opts.KeyMapping
		}
	}
	_ = client.SendPacket(protocol.New_PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK(data))
}

func (h *Handler) handleSaveSettings(client *network.Client, pkt *protocol.Packet) {
	charNo := h.getActiveCharNo(client)
	if charNo <= 0 {
		log.Printf("[World][%s] Cannot save settings: No active character selected", client.RemoteAddr())
		return
	}

	payload, _ := pkt.Reader.ReadBytes(pkt.Reader.Remaining())
	var optType string

	switch pkt.OpCode {
	case protocol.PROTO_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD:
		optType = "shortcut_data"
	case protocol.PROTO_NC_CHAR_OPTION_SET_SHORTCUTSIZE_CMD:
		optType = "shortcut_size"
	case protocol.PROTO_NC_CHAR_OPTION_SET_GAME_CMD:
		optType = "game"
	case protocol.PROTO_NC_CHAR_OPTION_SET_WINDOWPOS_CMD:
		optType = "window_pos"
	case protocol.PROTO_NC_CHAR_OPTION_SET_KEYMAPPING_CMD:
		optType = "key_mapping"
	default:
		log.Printf("[World][%s] Unknown save settings OpCode 0x%04X", client.RemoteAddr(), pkt.OpCode)
		return
	}

	if err := h.charRepo.SaveCharacterOption(charNo, optType, payload); err != nil {
		log.Printf("[World][%s] Error saving %s for char %d: %v", client.RemoteAddr(), optType, charNo, err)
		return
	}

	log.Printf("[World][%s] Saved %s (%d bytes) in tCharacterOptions for char %d",
		client.RemoteAddr(), optType, len(payload), charNo)
}

func (h *Handler) handlePrisonGet(client *network.Client, pkt *protocol.Packet) {
	resp := protocol.NewPacketWriter(protocol.PROTO_NC_PRISON_GET_ACK)
	resp.WriteByte(0) // bInPrison = 0 (libre / no encarcelado)
	_ = client.SendPacket(resp)
}

func (h *Handler) handleGotIngame(client *network.Client, pkt *protocol.Packet) {
	log.Printf("[World][%s] Player '%s' entering in-game zone (0x%04X). Dispatching official in-game ACKs...",
		client.RemoteAddr(), client.Username, pkt.OpCode)

	// 1. Registrar en KQManager y despachar reloj y lista oficial de Kingdom Quests
	h.kqMgr.RegisterOnline(client, client.CharacterID)
	h.kqMgr.HandleList(client)

	// 2. Despachar lista oficial de amigos del jugador desde tFriend (0x5407)
	h.friendMgr.SendFriendList(client)
	currentMap := ""
	if h.charRepo != nil && client.CharacterID > 0 {
		if ch, err := h.charRepo.GetCharacterByID(client.CharacterID); err == nil && ch != nil {
			currentMap = ch.TownName
			if currentMap == "" && h.dataStore != nil && h.dataStore.Fields != nil {
				if f, ok := h.dataStore.Fields.GetFieldBySerial(int(ch.Map)); ok {
					currentMap = f.MapID
				}
			}
		}
	}
	h.friendMgr.BroadcastLogin(client, client.Username, currentMap)

	// 3. Registrar en GuildManager y enviar información de gremio si pertenece a uno
	h.guildMgr.RegisterOnline(client, client.CharacterID, client.Username)
	h.guildMgr.HandleMyGuildInfo(client)
	h.guildMgr.HandleMemberList(client)

	// 4. Packet PROTO_NC_MISC_SERVER_TIME_NOTIFY_CMD (0x0845 / Dept:2, Cmd:69) - Sincronización oficial del reloj
	p4 := protocol.NewPacketWriter(protocol.PROTO_NC_MISC_SERVER_TIME_NOTIFY_CMD)
	now := time.Now()
	millis := uint32(now.Hour()*3600000 + now.Minute()*60000 + now.Second()*1000 + now.Nanosecond()/1000000)
	p4.WriteUint32(millis)
	_ = client.SendPacket(p4)

	log.Printf("[World][%s] Official in-game ACK sequence sent successfully for '%s'",
		client.RemoteAddr(), client.Username)
}
