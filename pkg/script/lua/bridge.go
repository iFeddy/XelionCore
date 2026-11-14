package lua

// WorldBridge desacopla el runtime de Lua del motor de ZoneService, proveyendo
// la interfaz que las 274 funciones C-API de Zone.pdb invocan para consultar o
// mutar el estado del mundo, jugadores, monstruos, puertas y eventos.
type WorldBridge interface {
	// Entidades / Objetos
	GetObjectHP(handle uint32) int32
	SetObjectHP(handle uint32, hp int32)
	KillObject(handle uint32)
	GetObjectType(handle uint32) int
	GetObjectCoord(handle uint32) (x, y float64, mapName string, ok bool)
	GetMaster(handle uint32) uint32
	GetAdminLevel(handle uint32) int
	GetPlayerName(handle uint32) string
	GetCharNo(handle uint32) int64
	GetLevel(handle uint32) int
	GetBaseClass(handle uint32) int
	IsPlayerExist(handle uint32) bool
	IsObjectDead(handle uint32) bool
	IsMovable(handle uint32) bool
	GetWhoKillMe(handle uint32) uint32

	// Movimiento / Animación / Efectos
	MoveTo(handle uint32, x, y float64, speed float64)
	MoveStop(handle uint32)
	CastTeleport(handle uint32, mapName string, x, y float64)
	CameraMove(handle uint32, cameraID int)
	Animate(handle uint32, animID int, forced bool)
	Emotion(handle uint32, emotionID int)
	PlaySound(soundID int, x, y float64)
	ObjectEffect(handle uint32, effectID int)

	// Spawns, Monstruos y NPCs
	MobRegenXY(mobID uint32, mapName string, x, y float64, layer int) uint32
	MobRegenCircle(mobID uint32, mapName string, cx, cy float64, radius float64) uint32
	MobRegenRectangle(mobID uint32, mapName string, x1, y1, x2, y2 float64) uint32
	NPCRegen(npcID uint32, mapName string, x, y float64) uint32
	NPCVanish(handle uint32)
	GroupRegen(groupID uint32, mapName string) int
	AggroSet(mobHandle, targetHandle uint32, amount int)
	AggroReset(mobHandle uint32)

	// Mapa, Geometría y Puertas (.sbi)
	DoorAction(doorName string, action int) bool
	DoorBuild(doorName string) bool
	CloseAllDoors()
	GetCurMapIndex() int
	GetMapName() string
	IsInArea(areaName string, x, y float64) bool
	IsInMap(handle uint32, mapName string) bool

	// Combate, Skills y AbState
	SetAbstate(handle uint32, stateID uint16, strength uint16, durationSec uint32)
	ResetAbstate(handle uint32, stateID uint16)
	GetAbstate(handle uint32, stateID uint16) (int, bool)
	Heal(handle uint32, amount int32)
	Damaged(handle uint32, attackerHandle uint32, damage int32)
	SkillBlast(handle uint32, skillID uint16, x, y float64)

	// Kingdom Quest, Misiones y Rankings
	EndOfKingdomQuest(success bool)
	TimerStart(timerID int, seconds int)
	TimerEnd(timerID int)
	ScoreBoard(recordType int, charNo int64, score int)
	SendRankingList()

	// Avisos, Mensajes y Chat
	Notice(msg string)
	NoticeString(msg string)
	ScriptMsg(handle uint32, msg string)
	SysMsg(handle uint32, code int)
	SomebodyShout(handle uint32, msg string)

	// Auditoría GameLog
	SendGameLogDataType4(charNo int64, mapName string, x, y, z int, int1, int2, int3 int)
	SendGameLogDataType5(charNo int64, int1, int2 int)

	// Mascotas (Pet)
	PetGetMind(petHandle uint32) int
	PetChangeMind(petHandle uint32, newMind int)
	PetUnsummon(petHandle uint32)
	PetGetActionRecord(petHandle uint32, action string) (aiType, smInx, hairEffect, soundFile string, ok bool)
	PetGetRunSpeed(petHandle uint32) int
	PetGetWalkSpeed(petHandle uint32) int
	PetGetStress(petHandle uint32) int
	PetIsMasterCalling(petHandle uint32) bool
	PetSaveTendency(petHandle uint32) bool

	// Misiones, Menús Servidor, Tutorial y Recompensas Lua (Fase 15)
	ServerMenu(playerHandle, npcHandle uint32, title string, options []MenuOptionEntry)
	NPCMenuOpen(npcHandle, playerHandle uint32)
	QuestResult(mapName, result string)
	QuestResultIndividual(playerHandle uint32, result string)
	QuestMobKill(playerHandle uint32, questID uint16, mobID uint32, count int)
	QuestMobKillAllInMap(mapName string, questID uint16, mobID uint32, maxCount int)
	GetTutorialInfo(playerHandle uint32) int
	SaveTutorialInfo(playerHandle uint32, saveType, step int)
	ProgressTutorial(playerHandle uint32, step int)

	// Kingdom Quest / Instance Dungeon / Items Lua (Fase 15)
	ShowKQTimerWithLife(mapName string, seconds int)
	ShowKQTimerWithLifeObj(playerHandle uint32, seconds int)
	GetKQLimitSecond(mapName string) int
	GetKQTeamType(handle uint32) int
	IsKQJoiner(handle uint32) bool
	KQRewardIndex(playerHandle uint32, rewardIndex string)
	ScoreBoardAllInMap(mapName string, show bool, round, redWin, redScore, blueWin, blueScore int)
	IndunRankResult(mapName string, args []int)
	RewardAll(mapName, rewardType string)
	RewardItem(playerHandle uint32, itemInx string, count int)
	InvenItemDestroy(playerHandle uint32, itemInx string, count int) int
	GetItemLot(playerHandle uint32, itemInx string) (int, bool)
	DropItem(itemInx string, mobHandle, playerHandle uint32, rate int)

	// Helpers de Escenario Wedding & GuildTournament Lua (Fase 15)
	SpawnNamedNPC(mapName, mobInx string, x, y float64, dir, scale int) uint32
	ScenarioSetTimeLimit(field string, seconds int)
	ScenarioSetInterrupt(field, eventName, callback string)
	ScenarioSetScriptFile(field, scriptName string)
	ScenarioSendChatWindow(field, speaker, dialogKey string)
	ScenarioEnd(field string)
	ScenarioFindPlayer(field, charName, mode string) uint32
	ScenarioCompleteWedding(groomHandle, brideHandle uint32)
	ScenarioApplyAbnormalState(handle uint32, staName string, level, durationMs int)
	ScenarioSetNPCInteractHandler(npcHandle uint32, callback string)
	ScenarioRegisterCaptureFlag(field string, npcHandle uint32, redEffect, blueEffect string)
	ScenarioRegisterGTPortal(field string, npcHandle, flagHandle uint32, redName, blueName string)
	ScenarioRegisterGuildGate(field string, npcHandle uint32, guildName string, flagHandle uint32, gateName string)
	ScenarioBroadcastToField(field, msg string)
	ScenarioEvaluateMatchScore(field string) string
	ScenarioRewardGuildTournament(field, winningGuild string)
	ScenarioWait(seconds int)
}

// MenuOptionEntry representa una opción de menú enviada por cServerMenu desde Lua.
type MenuOptionEntry struct {
	Text     string
	Callback string
}

// DefaultBridge provee una implementación base no nula y segura para entornos de test
// o inicializaciones previas al enlazado del WorldMap.
type DefaultBridge struct{}

func (b *DefaultBridge) GetObjectHP(handle uint32) int32                                  { return 100 }
func (b *DefaultBridge) SetObjectHP(handle uint32, hp int32)                              {}
func (b *DefaultBridge) KillObject(handle uint32)                                         {}
func (b *DefaultBridge) GetObjectType(handle uint32) int                                  { return 2 } // Player
func (b *DefaultBridge) GetObjectCoord(handle uint32) (float64, float64, string, bool)   { return 6200, 6400, "RouN", true }
func (b *DefaultBridge) GetMaster(handle uint32) uint32                                   { return 0 }
func (b *DefaultBridge) GetAdminLevel(handle uint32) int                                  { return 0 }
func (b *DefaultBridge) GetPlayerName(handle uint32) string                               { return "Player" }
func (b *DefaultBridge) GetCharNo(handle uint32) int64                                   { return int64(handle) }
func (b *DefaultBridge) GetLevel(handle uint32) int                                       { return 1 }
func (b *DefaultBridge) GetBaseClass(handle uint32) int                                   { return 1 }
func (b *DefaultBridge) IsPlayerExist(handle uint32) bool                                 { return true }
func (b *DefaultBridge) IsObjectDead(handle uint32) bool                                  { return false }
func (b *DefaultBridge) IsMovable(handle uint32) bool                                     { return true }
func (b *DefaultBridge) GetWhoKillMe(handle uint32) uint32                                { return 0 }
func (b *DefaultBridge) MoveTo(handle uint32, x, y float64, speed float64)                {}
func (b *DefaultBridge) MoveStop(handle uint32)                                           {}
func (b *DefaultBridge) CastTeleport(handle uint32, mapName string, x, y float64)         {}
func (b *DefaultBridge) CameraMove(handle uint32, cameraID int)                           {}
func (b *DefaultBridge) Animate(handle uint32, animID int, forced bool)                   {}
func (b *DefaultBridge) Emotion(handle uint32, emotionID int)                             {}
func (b *DefaultBridge) PlaySound(soundID int, x, y float64)                              {}
func (b *DefaultBridge) ObjectEffect(handle uint32, effectID int)                         {}
func (b *DefaultBridge) MobRegenXY(mobID uint32, mapName string, x, y float64, layer int) uint32 { return 1000 }
func (b *DefaultBridge) MobRegenCircle(mobID uint32, mapName string, cx, cy float64, radius float64) uint32 { return 1001 }
func (b *DefaultBridge) MobRegenRectangle(mobID uint32, mapName string, x1, y1, x2, y2 float64) uint32 { return 1002 }
func (b *DefaultBridge) NPCRegen(npcID uint32, mapName string, x, y float64) uint32      { return 2000 }
func (b *DefaultBridge) NPCVanish(handle uint32)                                          {}
func (b *DefaultBridge) GroupRegen(groupID uint32, mapName string) int                    { return 1 }
func (b *DefaultBridge) AggroSet(mobHandle, targetHandle uint32, amount int)              {}
func (b *DefaultBridge) AggroReset(mobHandle uint32)                                      {}
func (b *DefaultBridge) DoorAction(doorName string, action int) bool                      { return true }
func (b *DefaultBridge) DoorBuild(doorName string) bool                                   { return true }
func (b *DefaultBridge) CloseAllDoors()                                                   {}
func (b *DefaultBridge) GetCurMapIndex() int                                              { return 1 }
func (b *DefaultBridge) GetMapName() string                                               { return "RouN" }
func (b *DefaultBridge) IsInArea(areaName string, x, y float64) bool                      { return true }
func (b *DefaultBridge) IsInMap(handle uint32, mapName string) bool                       { return true }
func (b *DefaultBridge) SetAbstate(handle uint32, stateID uint16, strength uint16, durationSec uint32) {}
func (b *DefaultBridge) ResetAbstate(handle uint32, stateID uint16)                       {}
func (b *DefaultBridge) GetAbstate(handle uint32, stateID uint16) (int, bool)             { return 0, false }
func (b *DefaultBridge) Heal(handle uint32, amount int32)                                 {}
func (b *DefaultBridge) Damaged(handle uint32, attackerHandle uint32, damage int32)       {}
func (b *DefaultBridge) SkillBlast(handle uint32, skillID uint16, x, y float64)           {}
func (b *DefaultBridge) EndOfKingdomQuest(success bool)                                   {}
func (b *DefaultBridge) TimerStart(timerID int, seconds int)                              {}
func (b *DefaultBridge) TimerEnd(timerID int)                                             {}
func (b *DefaultBridge) ScoreBoard(recordType int, charNo int64, score int)               {}
func (b *DefaultBridge) SendRankingList()                                                 {}
func (b *DefaultBridge) Notice(msg string)                                                {}
func (b *DefaultBridge) NoticeString(msg string)                                          {}
func (b *DefaultBridge) ScriptMsg(handle uint32, msg string)                              {}
func (b *DefaultBridge) SysMsg(handle uint32, code int)                                   {}
func (b *DefaultBridge) SomebodyShout(handle uint32, msg string)                          {}
func (b *DefaultBridge) SendGameLogDataType4(charNo int64, mapName string, x, y, z int, int1, int2, int3 int) {}
func (b *DefaultBridge) SendGameLogDataType5(charNo int64, int1, int2 int)                {}
func (b *DefaultBridge) PetGetMind(petHandle uint32) int                                  { return 50 }
func (b *DefaultBridge) PetChangeMind(petHandle uint32, newMind int)                      {}
func (b *DefaultBridge) PetUnsummon(petHandle uint32)                                     {}
func (b *DefaultBridge) PetGetActionRecord(petHandle uint32, action string) (string, string, string, string, bool) {
	return "Idle", "SM_PetIdle", "", "", true
}
func (b *DefaultBridge) PetGetRunSpeed(petHandle uint32) int                               { return 115 }
func (b *DefaultBridge) PetGetWalkSpeed(petHandle uint32) int                              { return 40 }
func (b *DefaultBridge) PetGetStress(petHandle uint32) int                                 { return 0 }
func (b *DefaultBridge) PetIsMasterCalling(petHandle uint32) bool                          { return false }
func (b *DefaultBridge) PetSaveTendency(petHandle uint32) bool                             { return true }
func (b *DefaultBridge) ServerMenu(playerHandle, npcHandle uint32, title string, options []MenuOptionEntry) {}
func (b *DefaultBridge) NPCMenuOpen(npcHandle, playerHandle uint32)                        {}
func (b *DefaultBridge) QuestResult(mapName, result string)                                {}
func (b *DefaultBridge) QuestResultIndividual(playerHandle uint32, result string)          {}
func (b *DefaultBridge) QuestMobKill(playerHandle uint32, questID uint16, mobID uint32, count int) {}
func (b *DefaultBridge) QuestMobKillAllInMap(mapName string, questID uint16, mobID uint32, maxCount int) {}
func (b *DefaultBridge) GetTutorialInfo(playerHandle uint32) int                           { return 0 }
func (b *DefaultBridge) SaveTutorialInfo(playerHandle uint32, saveType, step int)          {}
func (b *DefaultBridge) ProgressTutorial(playerHandle uint32, step int)                    {}
func (b *DefaultBridge) ShowKQTimerWithLife(mapName string, seconds int)                   {}
func (b *DefaultBridge) ShowKQTimerWithLifeObj(playerHandle uint32, seconds int)           {}
func (b *DefaultBridge) GetKQLimitSecond(mapName string) int                               { return 1800 }
func (b *DefaultBridge) GetKQTeamType(handle uint32) int                                   { return 0 }
func (b *DefaultBridge) IsKQJoiner(handle uint32) bool                                     { return true }
func (b *DefaultBridge) KQRewardIndex(playerHandle uint32, rewardIndex string)             {}
func (b *DefaultBridge) ScoreBoardAllInMap(mapName string, show bool, round, redWin, redScore, blueWin, blueScore int) {}
func (b *DefaultBridge) IndunRankResult(mapName string, args []int)                        {}
func (b *DefaultBridge) RewardAll(mapName, rewardType string)                              {}
func (b *DefaultBridge) RewardItem(playerHandle uint32, itemInx string, count int)         {}
func (b *DefaultBridge) InvenItemDestroy(playerHandle uint32, itemInx string, count int) int { return 1 }
func (b *DefaultBridge) GetItemLot(playerHandle uint32, itemInx string) (int, bool)        { return 1, false }
func (b *DefaultBridge) DropItem(itemInx string, mobHandle, playerHandle uint32, rate int) {}
func (b *DefaultBridge) SpawnNamedNPC(mapName, mobInx string, x, y float64, dir, scale int) uint32 { return 2000 }
func (b *DefaultBridge) ScenarioSetTimeLimit(field string, seconds int)                    {}
func (b *DefaultBridge) ScenarioSetInterrupt(field, eventName, callback string)            {}
func (b *DefaultBridge) ScenarioSetScriptFile(field, scriptName string)                    {}
func (b *DefaultBridge) ScenarioSendChatWindow(field, speaker, dialogKey string)           {}
func (b *DefaultBridge) ScenarioEnd(field string)                                          {}
func (b *DefaultBridge) ScenarioFindPlayer(field, charName, mode string) uint32            { return 1 }
func (b *DefaultBridge) ScenarioCompleteWedding(groomHandle, brideHandle uint32)           {}
func (b *DefaultBridge) ScenarioApplyAbnormalState(handle uint32, staName string, level, durationMs int) {}
func (b *DefaultBridge) ScenarioSetNPCInteractHandler(npcHandle uint32, callback string)   {}
func (b *DefaultBridge) ScenarioRegisterCaptureFlag(field string, npcHandle uint32, redEffect, blueEffect string) {}
func (b *DefaultBridge) ScenarioRegisterGTPortal(field string, npcHandle, flagHandle uint32, redName, blueName string) {}
func (b *DefaultBridge) ScenarioRegisterGuildGate(field string, npcHandle uint32, guildName string, flagHandle uint32, gateName string) {}
func (b *DefaultBridge) ScenarioBroadcastToField(field, msg string)                        {}
func (b *DefaultBridge) ScenarioEvaluateMatchScore(field string) string                    { return "" }
func (b *DefaultBridge) ScenarioRewardGuildTournament(field, winningGuild string)          {}
func (b *DefaultBridge) ScenarioWait(seconds int)                                          {}

