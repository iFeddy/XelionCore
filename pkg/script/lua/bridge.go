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

