package protocol

// MakeOpCode combines a 6-bit Header and a 10-bit Type into a 16-bit OpCode.
// Formula: (Header << 10) | (Type & 0x3FF)
func MakeOpCode(header uint8, typ uint16) uint16 {
	return (uint16(header) << 10) | (typ & 0x03FF)
}

// ExtractHeader extracts the 6-bit header from a 16-bit OpCode.
func ExtractHeader(opCode uint16) uint8 {
	return uint8(opCode >> 10)
}

// ExtractType extracts the 10-bit type from a 16-bit OpCode.
func ExtractType(opCode uint16) uint16 {
	return opCode & 0x03FF
}

// Server -> Client OpCodes (SH = Server Header)
var (
	// Header 2: System / Handshake
	SH2Type_Ping              = MakeOpCode(2, 4)
	SH2Type_SetXorKeyPosition = MakeOpCode(2, 7) // Handshake offset
	SH2Type_Unk1              = MakeOpCode(2, 14)
	SH2Type_UnkTimePacket     = MakeOpCode(2, 69)
	SH2Type_UpdateClientTime  = MakeOpCode(2, 73)

	// Header 3: Login / World Selection
	SH3Type_IncorrectVersion        = MakeOpCode(3, 2)
	SH3Type_FilecheckAllow          = MakeOpCode(3, 5)
	SH3Type_Error                   = MakeOpCode(3, 9)
	SH3Type_WorldlistNew            = MakeOpCode(3, 10)
	SH3Type_WorldServerIP           = MakeOpCode(3, 12)
	SH3Type_CharacterList           = MakeOpCode(3, 20)
	SH3Type_WorldistResend          = MakeOpCode(3, 28)
	SH3Type_BackToWorldListFromChar = MakeOpCode(3, 52)
	SH3Type_VersionAllowed          = MakeOpCode(3, 103)

	// Header 4: Character Selection / Zone Transfer
	SH4Type_ConnectError              = MakeOpCode(4, 2)
	SH4Type_ServerIP                  = MakeOpCode(4, 3)
	SH4Type_CharacterInfo             = MakeOpCode(4, 56)
	SH4Type_CharacterLook             = MakeOpCode(4, 57)
	SH4Type_CharacterQuestsBusy       = MakeOpCode(4, 58)
	SH4Type_CharacterQuestsDone       = MakeOpCode(4, 59)
	SH4Type_CharacterActiveSkillList  = MakeOpCode(4, 61)
	SH4Type_CharacterPassiveSkillList = MakeOpCode(4, 62)
	SH4Type_CharacterItemList         = MakeOpCode(4, 71)
	SH4Type_CharacterInfoEnd          = MakeOpCode(4, 72)
	SH4Type_CharacterTitles           = MakeOpCode(4, 73)
	SH4Type_CharacterPoints           = MakeOpCode(4, 91)

	// Header 5: Character Management
	SH5Type_CharCreationError          = MakeOpCode(5, 4)
	SH5Type_CharCreationOK             = MakeOpCode(5, 6)
	SH5Type_CharDeleteOK               = MakeOpCode(5, 12)
	SH5Type_SendCharacterChangeNewName = MakeOpCode(5, 16)

	// Header 6: Zone Transfer & Detailed Info
	SH6Type_DetailedCharacterInfo = MakeOpCode(6, 2)
	SH6Type_Error                 = MakeOpCode(6, 3)
	SH6Type_ChangeMap             = MakeOpCode(6, 9)
	SH6Type_ChangeZone            = MakeOpCode(6, 10)
	SH6Type_Teleporter            = MakeOpCode(6, 27)

	// Header 7: Map Spawns
	SH7Type_ShowUnequip       = MakeOpCode(7, 4)
	SH7Type_ShowEquip         = MakeOpCode(7, 5)
	SH7Type_SpawnSinglePlayer = MakeOpCode(7, 6)
	SH7Type_SpawnMultiPlayer  = MakeOpCode(7, 7)
	SH7Type_RemoveObject      = MakeOpCode(7, 14)

	// Header 8: Map Movement & Chat
	SH8Type_ChatNormal = MakeOpCode(8, 2)
	SH8Type_GmNotice   = MakeOpCode(8, 17)
	SH8Type_StopTele   = MakeOpCode(8, 19)
	SH8Type_Walk       = MakeOpCode(8, 24)
	SH8Type_Move       = MakeOpCode(8, 26)
	SH8Type_Teleport   = MakeOpCode(8, 27)
	SH8Type_Jump       = MakeOpCode(8, 37)
	SH8Type_BeginRest  = MakeOpCode(8, 40)
	SH8Type_EndRest    = MakeOpCode(8, 43)
	SH8Type_BlockWalk  = MakeOpCode(8, 74)

	// Header 12: Inventory & Items
	SH12Type_ModifyItemSlot           = MakeOpCode(12, 1)
	SH12Type_ModifyEquipSlot          = MakeOpCode(12, 2)
	SH12Type_InventoryFull            = MakeOpCode(12, 4)
	SH12Type_ObtainedItem             = MakeOpCode(12, 10)
	SH12Type_MoveIteminContaInComplet = MakeOpCode(12, 12)
	SH12Type_FailedEquip              = MakeOpCode(12, 17)
	SH12Type_FailedUnequip            = MakeOpCode(12, 19)
	SH12Type_ItemUseEffect            = MakeOpCode(12, 22)
	SH12Type_ItemUpgrade              = MakeOpCode(12, 24)
	SH12Type_ItemUsedOk               = MakeOpCode(12, 26)
	SH12Type_SendPremiumItemList      = MakeOpCode(12, 33)
	SH12Type_SendRewardList           = MakeOpCode(12, 45)

	// Header 21: Social / Friends
	SH21Type_FriendList = MakeOpCode(21, 7)

	// Header 22: In-Game Status / Kingdom Quest
	SH22Type_GotIngameAck      = MakeOpCode(22, 28) // 0x581C
	SH22Type_KingdomQuestCount = MakeOpCode(22, 29) // 0x581D

	// Header 28: Settings & QuickBar
	SH28Type_LoadQuickBar      = MakeOpCode(28, 3)
	SH28Type_LoadQuickBarState = MakeOpCode(28, 5)
	SH28Type_LoadGameSettings  = MakeOpCode(28, 11)
	SH28Type_LoadClientSettings = MakeOpCode(28, 13)
	SH28Type_LoadShortCuts     = MakeOpCode(28, 15)

	// Header 31: Unknown System Block
	SH31Type_LoadUnknown       = MakeOpCode(31, 7)
)

// Client -> Server OpCodes (CH = Client Header)
var (
	// Header 2: System / Handshake
	CH2Type_Pong = MakeOpCode(2, 5)
	CH2Type_Unk1 = MakeOpCode(2, 13)

	// Header 3: Login / World Selection
	CH3Type_FileHash          = MakeOpCode(3, 4)
	CH3Type_WorldSelect       = MakeOpCode(3, 11)
	CH3Type_WorldClientKey    = MakeOpCode(3, 15)
	CH3Type_BackToCharSelect  = MakeOpCode(3, 24)
	CH3Type_WorldReRequest    = MakeOpCode(3, 27)
	CH3Type_BackToWorldSelect = MakeOpCode(3, 51)
	CH3Type_LoginLegacy55     = MakeOpCode(3, 55)
	CH3Type_LoginLegacy56     = MakeOpCode(3, 56)
	CH3Type_Login             = MakeOpCode(3, 90) // Corrected for NA2016 client
	CH3Type_Version           = MakeOpCode(3, 101)

	// Header 4: Character Selection
	CH4Type_CharSelect = MakeOpCode(4, 1)

	// Header 5: Character Management
	CH5Type_CreateCharacter     = MakeOpCode(5, 1)
	CH5Type_ChangeCharacterName = MakeOpCode(5, 5)
	CH5Type_DeleteCharacter     = MakeOpCode(5, 7)

	// Header 6: Zone Transfer
	CH6Type_TransferKey = MakeOpCode(6, 1)
	CH6Type_ClientReady = MakeOpCode(6, 3)
	CH6Type_Teleporter  = MakeOpCode(6, 26)

	// Header 8: Map Movement & Chat
	CH8Type_ChatNormal = MakeOpCode(8, 1)
	CH8Type_Stop       = MakeOpCode(8, 18)
	CH8Type_Walk       = MakeOpCode(8, 23)
	CH8Type_Run        = MakeOpCode(8, 25)
	CH8Type_Emote      = MakeOpCode(8, 32)
	CH8Type_Jump       = MakeOpCode(8, 36)
	CH8Type_BeginRest  = MakeOpCode(8, 39)
	CH8Type_EndRest    = MakeOpCode(8, 42)

	// Header 22: In-Game Status / Heartbeat
	CH22Type_GotIngame = MakeOpCode(22, 27) // 0x581B, matching fc_NC_CHAR_LOGIN_REQ

	// Header 28: Settings & QuickBar
	CH28Type_GetQuickBar       = MakeOpCode(28, 2)
	CH28Type_GetQuickBarState  = MakeOpCode(28, 4)
	CH28Type_GetGameSettings   = MakeOpCode(28, 10)
	CH28Type_GetClientSettings = MakeOpCode(28, 12)
	CH28Type_GetShortCuts      = MakeOpCode(28, 14)
	CH28Type_SaveQuickBar      = MakeOpCode(28, 16)
	CH28Type_SaveQuickBarState = MakeOpCode(28, 17)
	CH28Type_SaveGameSettings  = MakeOpCode(28, 20)
	CH28Type_SaveClientSettings = MakeOpCode(28, 21)
	CH28Type_SaveShortCuts     = MakeOpCode(28, 22)

	// Header 31: Unknown System Block
	CH31Type_GetUnknown        = MakeOpCode(31, 6)
)

// ServerError codes returned in SH3Type_Error packets.
type ServerError uint16

const (
	ErrException          ServerError = 0x42 // Unknown Error / Exception
	ErrDatabase           ServerError = 0x43 // DB Error
	ErrWrongRegion        ServerError = 0x44 // Auth Failed / Wrong Region
	ErrInvalidCredentials ServerError = 0x45 // Please check ID or password
	ErrDatabaseError      ServerError = 0x46 // DB Error
	ErrBlocked            ServerError = 0x47 // The ID has been blocked
	ErrServerMaintenance  ServerError = 0x48 // World server maintenance
	ErrTimeout            ServerError = 0x49 // Timeout
	ErrLoginFailed        ServerError = 0x4A // Login Failed
	ErrAgreementMissing   ServerError = 0x4B // Please accept the agreement
)

// WorldStatus states for server selection display.
type WorldStatus uint8

const (
	WorldStatus_Offline         WorldStatus = 0
	WorldStatus_Maintenance     WorldStatus = 1
	WorldStatus_EmptyServer     WorldStatus = 2
	WorldStatus_Reserved        WorldStatus = 3
	WorldStatus_OfflineUnkError WorldStatus = 4
	WorldStatus_Full            WorldStatus = 5
	WorldStatus_Low             WorldStatus = 6
	WorldStatus_Medium          WorldStatus = 9
	WorldStatus_High            WorldStatus = 10
)

// CreateCharError codes returned in SH5Type_CharCreationError packets.
type CreateCharError uint16

const (
	CreateCharError_FailedToCreate CreateCharError = 130 // Failed to create the char
	CreateCharError_WrongClass     CreateCharError = 131 // You have chosen wrong class for your avatar
	CreateCharError_NameTaken      CreateCharError = 132 // The name is already taken by another user
	CreateCharError_ErrorInMaxSlot CreateCharError = 133 // There is an error in the Max slot Number
	CreateCharError_NameInUse      CreateCharError = 385 // The name is already in use
)

// Job identifies character classes in Fiesta Online.
type Job byte

const (
	Job_Fighter   Job = 1
	Job_Cleric    Job = 6
	Job_Archer    Job = 11
	Job_Mage      Job = 16
	Job_Trickster Job = 21
)
