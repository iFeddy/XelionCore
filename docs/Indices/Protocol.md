# Ãndice: pkg/protocol

> Generado automÃ¡ticamente por scripts/gen_indices.ps1. Solo firmas con nÂ° de lÃ­nea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. 

## pkg/protocol/opcodes.go
- L12: `func MakeOpCode(dept uint8, cmd uint16) uint16`
- L17: `func ExtractHeader(opcode uint16) uint8`
- L22: `func ExtractType(opcode uint16) uint16`
- L2962: `type ServerError uint16`
- L2978: `type WorldStatus uint8`
- L2993: `type CreateCharError uint16`
- L3004: `type Job byte`

## pkg/protocol/packet.go
- L6: `type Packet struct`
- L15: `func NewPacket(data []byte) *Packet`

## pkg/protocol/packet_reader.go
- L11: `type PacketReader struct`
- L17: `func NewPacketReader(data []byte) *PacketReader`
- L25: `func (r *PacketReader) Remaining() int`
- L30: `func (r *PacketReader) ReadByte() (byte, error)`
- L40: `func (r *PacketReader) ReadBool() (bool, error)`
- L49: `func (r *PacketReader) ReadUint16() (uint16, error)`
- L59: `func (r *PacketReader) ReadInt16() (int16, error)`
- L65: `func (r *PacketReader) ReadUint32() (uint32, error)`
- L75: `func (r *PacketReader) ReadInt32() (int32, error)`
- L81: `func (r *PacketReader) ReadUint64() (uint64, error)`
- L91: `func (r *PacketReader) ReadInt64() (int64, error)`
- L97: `func (r *PacketReader) ReadBytes(n int) ([]byte, error)`
- L108: `func (r *PacketReader) ReadString(fixedLen int) (string, error)`
- L121: `func (r *PacketReader) Skip(n int) error`

## pkg/protocol/packet_test.go
- L8: `func TestPacketSerializationAndCrypto(t *testing.T)`
- L42: `func TestOpCodePacking(t *testing.T)`
- L55: `func TestStringPadding(t *testing.T)`

## pkg/protocol/packet_writer.go
- L8: `type PacketWriter struct`
- L14: `func NewPacketWriter(opCode uint16) *PacketWriter`
- L24: `func (pw *PacketWriter) OpCode() uint16`
- L29: `func (pw *PacketWriter) WriteByte(b byte)`
- L34: `func (pw *PacketWriter) WriteBool(v bool)`
- L43: `func (pw *PacketWriter) WriteUint16(v uint16)`
- L50: `func (pw *PacketWriter) WriteInt16(v int16)`
- L55: `func (pw *PacketWriter) WriteUint32(v uint32)`
- L62: `func (pw *PacketWriter) WriteInt32(v int32)`
- L67: `func (pw *PacketWriter) WriteUint64(v uint64)`
- L74: `func (pw *PacketWriter) WriteInt64(v int64)`
- L79: `func (pw *PacketWriter) WriteBytes(b []byte)`
- L84: `func (pw *PacketWriter) WriteString(s string, fixedLen int)`
- L96: `func (pw *PacketWriter) Fill(count int, b byte)`
- L103: `func (pw *PacketWriter) Data() []byte`
- L112: `func (pw *PacketWriter) BuildWirePacket() []byte`

## pkg/protocol/proto_nc_char.go
- L6: `func New_PROTO_NC_CHAR_CREATE_FAIL_ACK(errorCode uint16) *PacketWriter`
- L16: `func New_PROTO_NC_CHAR_CREATE_ACK(charInfoWriter func(pw *PacketWriter)) *PacketWriter`
- L28: `func New_PROTO_NC_CHAR_DELETE_ACK(slot byte) *PacketWriter`
- L38: `func New_PROTO_NC_CHAR_LOGIN_ACK(ip string, port uint16) *PacketWriter`
- L48: `func New_PROTO_NC_CHAR_LOGINFAIL_ACK(errCode uint16) *PacketWriter`
- L56: `func New_PROTO_NC_USER_WILL_WORLD_SELECT_ACK() *PacketWriter`
- L64: `func New_PROTO_NC_USER_BACKTOWORLDLIST_ACK() *PacketWriter`
- L70: `func New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK(data []byte) *PacketWriter`
- L89: `func New_PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK(data []byte) *PacketWriter`
- L108: `func New_PROTO_NC_CHAR_OPTION_GET_GAME_ACK(data []byte) *PacketWriter`
- L127: `func New_PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK(data []byte) *PacketWriter`
- L146: `func New_PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK(data []byte) *PacketWriter`
- L165: `type SHORT_CUT_DATA struct`
- L173: `type KEY_MAP_DATA struct`
- L181: `type GAME_OPTION_DATA struct`
- L190: `func New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD(entries []SHORT_CUT_DATA, rawData []byte) *PacketWriter`
- L213: `func New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD(entries []KEY_MAP_DATA, rawData []byte) *PacketWriter`
- L236: `func New_PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD(entries []GAME_OPTION_DATA, rawData []byte) *PacketWriter`
- L259: `func New_PROTO_NC_MISC_GAMETIME_ACK(hour, minute, second byte) *PacketWriter`

## pkg/protocol/proto_nc_friend.go
- L8: `func New_PROTO_NC_FRIEND_SET_CONFIRM_REQ(requesterName string) *PacketWriter`
- L16: `func New_PROTO_NC_FRIEND_SET_ACK(targetName string, errCode uint16) *PacketWriter`
- L25: `type FriendEntry struct`
- L34: `func New_PROTO_NC_FRIEND_LIST_CMD(friends []FriendEntry) *PacketWriter`
- L48: `func New_PROTO_NC_FRIEND_ADD_CMD(f FriendEntry) *PacketWriter`
- L59: `func New_PROTO_NC_FRIEND_DEL_ACK(friendName string, errCode uint16) *PacketWriter`
- L68: `func New_PROTO_NC_FRIEND_LOGIN_CMD(charName, mapName string) *PacketWriter`
- L77: `func New_PROTO_NC_FRIEND_LOGOUT_CMD(charName string) *PacketWriter`
- L84: `func New_PROTO_NC_FRIEND_POINT_ACK(points uint32) *PacketWriter`

## pkg/protocol/proto_nc_guild.go
- L17: `func packLogonInfo(t time.Time, isOnline bool) uint32`
- L36: `func packTm(t time.Time) []byte`
- L55: `type GuildListItemEntry struct`
- L66: `func New_PROTO_NC_GUILD_LIST_ACK(errCode uint16, partMask byte, totalCount, listCount uint16, guilds []GuildListItemEntry) *PacketWriter`
- L86: `func New_PROTO_NC_GUILD_MAKE_ACK(errCode uint16, guildNo uint32, guildName string, guildType byte, createDate time.Time) *PacketWriter`
- L98: `func New_PROTO_NC_GUILD_DISMISS_ACK(errCode uint16) *PacketWriter`
- L106: `func New_PROTO_NC_GUILD_DISMISS_CMD(guildName string, dismissStatus byte, dismissDate time.Time) *PacketWriter`
- L116: `func New_PROTO_NC_GUILD_MEMBER_INVITE_REQ(targetName string) *PacketWriter`
- L124: `func New_PROTO_NC_GUILD_MEMBER_INVITE_ACK(targetName string, errCode uint16) *PacketWriter`
- L133: `func New_PROTO_NC_GUILD_MEMBER_JOIN_REQ(guildName, fromCharID string) *PacketWriter`
- L142: `func New_PROTO_NC_GUILD_MEMBER_JOIN_ACK(guildName string, joinAck byte) *PacketWriter`
- L150: `type GuildMemberClientEntry struct`
- L165: `func writeGuildMemberClient(pw *PacketWriter, m GuildMemberClientEntry)`
- L180: `func New_PROTO_NC_GUILD_MEMBER_JOIN_CMD(member GuildMemberClientEntry) *PacketWriter`
- L187: `func New_PROTO_NC_GUILD_MEMBER_LEAVE_ACK(errCode uint16) *PacketWriter`
- L195: `func New_PROTO_NC_GUILD_MEMBER_LEAVE_CMD(charID string) *PacketWriter`
- L203: `func New_PROTO_NC_GUILD_MEMBER_VANISH_ACK(charID string, errCode uint16) *PacketWriter`
- L212: `func New_PROTO_NC_GUILD_MEMBER_VANISH_CMD(charID string) *PacketWriter`
- L220: `func New_PROTO_NC_GUILD_MEMBER_GRADE_ACK(charID string, grade byte, errCode uint16) *PacketWriter`
- L230: `func New_PROTO_NC_GUILD_MEMBER_GRADE_CMD(charID string, grade byte) *PacketWriter`
- L238: `func New_PROTO_NC_GUILD_MEMBER_LOGON_CMD(charID string) *PacketWriter`
- L245: `func New_PROTO_NC_GUILD_MEMBER_LOGOFF_CMD(charID string) *PacketWriter`
- L252: `func New_PROTO_NC_GUILD_MEMBER_MAP_CMD(charID string, mapName string) *PacketWriter`
- L260: `func New_PROTO_NC_GUILD_MEMBER_LEVEL_CMD(charID string, level byte) *PacketWriter`
- L269: `func New_PROTO_NC_GUILD_NOTIFY_ACK(errCode uint16) *PacketWriter`
- L277: `func New_PROTO_NC_GUILD_NOTIFY_CMD(notifyDate time.Time, notifyCharID string, notify string) *PacketWriter`
- L289: `type GuildClientDetails struct`
- L312: `func New_PROTO_NC_GUILD_MY_GUILD_INFO_ACK(errCode uint16, g *GuildClientDetails) *PacketWriter`
- L371: `func New_PROTO_NC_GUILD_MEMBER_LIST_ACK(totalCount, memberIndex, count uint16, members []GuildMemberClientEntry) *PacketWriter`
- L383: `func New_PROTO_NC_GUILD_CHAT_ACK(errCode uint16) *PacketWriter`
- L391: `func New_PROTO_NC_GUILD_CHAT_CMD(guildNo uint32, talker string, msg string) *PacketWriter`
- L403: `func New_PROTO_NC_GUILD_NAME_ACK(guildNo uint32, guildName string) *PacketWriter`
- L415: `type AcademyListItemEntry struct`
- L425: `func New_PROTO_NC_GUILD_ACADEMY_LIST_ACK(errCode uint16, partMask byte, totalCount, count uint32, list []AcademyListItemEntry) *PacketWriter`
- L442: `type AcademyClientDetails struct`
- L457: `func New_PROTO_NC_GUILD_ACADEMY_MY_GUILD_ACADEMY_INFO_ACK(errCode uint16, a *AcademyClientDetails) *PacketWriter`
- L480: `type AcademyMemberClientEntry struct`
- L494: `func writeAcademyMemberClient(pw *PacketWriter, m AcademyMemberClientEntry)`
- L509: `func New_PROTO_NC_GUILD_ACADEMY_MEMBER_LIST_ACK(totalCount, memberIndex, count uint16, members []AcademyMemberClientEntry) *PacketWriter`
- L522: `func New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_ACK(guildName string, errCode uint16) *PacketWriter`
- L531: `func New_PROTO_NC_GUILD_ACADEMY_MEMBER_JOIN_CMD(member AcademyMemberClientEntry) *PacketWriter`
- L538: `func New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_ACK(errCode uint16) *PacketWriter`
- L546: `func New_PROTO_NC_GUILD_ACADEMY_MEMBER_LEAVE_CMD(charID string) *PacketWriter`
- L553: `func New_PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_ACK(charID string, errCode uint16) *PacketWriter`
- L561: `func New_PROTO_NC_GUILD_ACADEMY_MEMBER_VANISH_CMD(charID string) *PacketWriter`
- L568: `func New_PROTO_NC_GUILD_ACADEMY_NOTIFY_ACK(errCode uint16) *PacketWriter`
- L575: `func New_PROTO_NC_GUILD_ACADEMY_NOTIFY_CMD(notifyDate time.Time, notifyCharID, notify string) *PacketWriter`
- L587: `func New_PROTO_NC_GUILD_ACADEMY_CHAT_ACK(errCode uint16) *PacketWriter`
- L594: `func New_PROTO_NC_GUILD_ACADEMY_CHAT_CMD(guildNo uint32, talker string, msg string) *PacketWriter`

## pkg/protocol/proto_nc_kq.go
- L11: `type ProtoKQInfoClient struct`
- L34: `func writeProtoKQInfoClient(pw *PacketWriter, q ProtoKQInfoClient)`
- L60: `func New_PROTO_NC_KQ_LIST_TIME_ACK(serverTime time.Time) *PacketWriter`
- L69: `func New_PROTO_NC_KQ_LIST_ACK(serverTime time.Time, startHandle, endHandle uint32, quests []ProtoKQInfoClient) *PacketWriter`
- L84: `func New_PROTO_NC_KQ_LIST_ADD_ACK(quests []ProtoKQInfoClient) *PacketWriter`
- L94: `type KQUpdateItem struct`
- L101: `func New_PROTO_NC_KQ_LIST_UPDATE_ACK(updates []KQUpdateItem) *PacketWriter`
- L113: `func New_PROTO_NC_KQ_LIST_DELETE_ACK(deleteHandles []uint32) *PacketWriter`
- L124: `func New_PROTO_NC_KQ_STATUS_ACK(handle uint32, status byte, joiners []string) *PacketWriter`
- L137: `func New_PROTO_NC_KQ_JOIN_ACK(handle uint32, errCode uint16) *PacketWriter`
- L146: `func New_PROTO_NC_KQ_JOIN_CANCEL_ACK(handle uint32, errCode uint16) *PacketWriter`
- L154: `type KQJoinCharInfo struct`
- L163: `func New_PROTO_NC_KQ_JOIN_LIST_ACK(errCode uint16, chars []KQJoinCharInfo) *PacketWriter`
- L178: `func New_PROTO_NC_KQ_TEAM_SELECT_ACK(errCode uint16, teamType byte) *PacketWriter`
- L187: `func New_PROTO_NC_KQ_TEAM_SELECT_CMD(charName string, teamType byte) *PacketWriter`
- L196: `func New_PROTO_NC_KQ_JOINING_ALARM_CMD(handle uint32, id uint16, minLev, maxLev byte, msg string) *PacketWriter`
- L210: `func New_PROTO_NC_KQ_JOINING_ALARM_END_CMD(handle uint32, id uint16) *PacketWriter`
- L218: `func New_PROTO_NC_KQ_NOTIFY_CMD(msg string) *PacketWriter`

## pkg/protocol/proto_nc_party.go
- L8: `func New_PROTO_NC_PARTY_JOINPROPOSE_REQ(masterName string) *PacketWriter`
- L16: `func New_PROTO_NC_PARTY_JOIN_ACK(memberID string, charClass byte, errCode uint16) *PacketWriter`
- L26: `type PartyMemberInfo struct`
- L34: `func New_PROTO_NC_PARTY_MEMBER_LIST_CMD(members []PartyMemberInfo) *PacketWriter`
- L46: `func New_PROTO_NC_PARTY_LEAVE_ACK(errCode uint16) *PacketWriter`
- L54: `func New_PROTO_NC_PARTY_LEAVE_CMD(partyID uint16, charID uint32) *PacketWriter`
- L63: `func New_PROTO_NC_PARTY_KICKOFF_ACK(targetName string, errCode uint16) *PacketWriter`
- L72: `func New_PROTO_NC_PARTY_DISMISS_CMD(partyID uint16) *PacketWriter`
- L79: `func New_PROTO_NC_PARTY_CHANGEMASTER_ACK(errCode uint16) *PacketWriter`
- L87: `func New_PROTO_NC_PARTY_ITEM_LOOTING_CMD(lootType, itemGrade byte) *PacketWriter`
- L95: `type PartyFinderEntry struct`
- L103: `func New_PROTO_NC_PARTY_FINDER_LIST_ACK(entries []PartyFinderEntry) *PacketWriter`

## pkg/protocol/proto_nc_report.go
- L7: `func New_PROTO_NC_DATA_REPORT_ADD_ACK(errCode uint16, hReporter uint16, harmer string, regNumber int32) *PacketWriter`
- L17: `func New_PROTO_NC_DATA_REPORT_CANCEL_ACK(regNumber int32, hReporter uint16, errCode uint16) *PacketWriter`
- L26: `func New_PROTO_NC_DATA_REPORT_GET_ACK(hReporter uint16, errCode uint16, regNumber int32, result uint8, remark string) *PacketWriter`

## pkg/protocol/proto_nc_user.go
- L7: `type PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ_DATA struct`
- L12: `type PROTO_NC_USER_US_LOGIN_REQ_DATA struct`
- L19: `type PROTO_NC_USER_LOGIN_ACK_WorldInfo struct`
- L26: `type PROTO_NC_USER_LOGINFAIL_ACK_DATA struct`
- L31: `type PROTO_NC_USER_WORLDSELECT_REQ_DATA struct`
- L36: `type PROTO_NC_USER_WORLDSELECT_ACK_DATA struct`
- L46: `func New_PROTO_NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK() *PacketWriter`
- L54: `func New_PROTO_NC_MISC_HEARTBEAT_ACK() *PacketWriter`
- L59: `func New_PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK() *PacketWriter`
- L65: `func New_PROTO_NC_USER_LOGINFAIL_ACK(errCode uint16) *PacketWriter`
- L74: `func New_PROTO_NC_USER_LOGIN_ACK(opCode uint16, worlds []PROTO_NC_USER_LOGIN_ACK_WorldInfo) *PacketWriter`
- L90: `func New_PROTO_NC_USER_WORLDSELECT_ACK(status uint8, ip string, port uint16, token string) *PacketWriter`
- L104: `func New_PROTO_NC_USER_FILE_CHECK_ACK(allow bool) *PacketWriter`

## pkg/protocol/proto_nc_zone.go
- L9: `func New_PROTO_NC_MAP_LOGINFAIL_ACK(errCode uint16, wrongDataFileIndex byte) *PacketWriter`
- L18: `func New_PROTO_NC_CHAR_INGAME_ACK() *PacketWriter`
- L24: `func New_PROTO_NC_CHAR_POINTS_ACK(usablePoints byte) *PacketWriter`
- L31: `type ShineCharStatVar struct`
- L37: `func (s ShineCharStatVar) WriteToPacket(pw *PacketWriter)`
- L46: `type PROTO_NC_CHAR_ZONE_CHARDATA_REQ_DATA struct`
- L55: `type PROTO_NC_MAP_LOGIN_REQ_DATA struct`
- L62: `func Read_PROTO_NC_MAP_LOGIN_REQ(pr *PacketReader) (*PROTO_NC_MAP_LOGIN_REQ_DATA, error)`

