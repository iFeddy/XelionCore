package zone

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/misc"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/script/lua"
	glua "github.com/yuin/gopher-lua"
)

// ZoneWorldBridge implementa lua.WorldBridge conectando scripts Lua AIScript / MobBehaviorDescript
// directamente al estado vivo de un ZoneMap.
type ZoneWorldBridge struct {
	zMap *ZoneMap
}

func NewZoneWorldBridge(zMap *ZoneMap) *ZoneWorldBridge {
	return &ZoneWorldBridge{zMap: zMap}
}

func (b *ZoneWorldBridge) GetObjectHP(handle uint32) int32 {
	if b.zMap == nil {
		return 0
	}
	if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
		return int32(p.Char.CurHP)
	}
	if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		mob.mu.RLock()
		defer mob.mu.RUnlock()
		return int32(mob.CurHP)
	}
	return 0
}

func (b *ZoneWorldBridge) SetObjectHP(handle uint32, hp int32) {
	if b.zMap == nil {
		return
	}
	if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
		p.Char.CurHP = int(hp)
	}
	if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		mob.mu.Lock()
		mob.CurHP = int(hp)
		if mob.CurHP <= 0 {
			mob.CurHP = 0
			mob.IsDead = true
		}
		mob.mu.Unlock()
	}
}

func (b *ZoneWorldBridge) KillObject(handle uint32) {
	if b.zMap == nil {
		return
	}
	if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		mob.mu.Lock()
		mob.CurHP = 0
		alreadyDead := mob.IsDead
		mob.IsDead = true
		mob.mu.Unlock()
		if !alreadyDead {
			b.zMap.NotifyMonsterDeath(mob, nil)
		}
		return
	}
	if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
		p.Char.CurHP = 0
		p.State = 3
		emitReallyKill(b.zMap, p.MapObjectID, p.MapObjectID)
	}
}

func (b *ZoneWorldBridge) GetObjectType(handle uint32) int {
	if b.zMap == nil {
		return 0
	}
	if b.zMap.GetPlayer(uint16(handle)) != nil {
		return 2
	}
	if b.zMap.GetMonster(uint16(handle)) != nil {
		return 6
	}
	return 0
}

func (b *ZoneWorldBridge) GetObjectCoord(handle uint32) (float64, float64, string, bool) {
	if b.zMap == nil {
		return 0, 0, "", false
	}
	if ent := b.zMap.GetEntity(uint16(handle)); ent != nil {
		x, y := ent.GetPosition()
		return float64(x), float64(y), b.zMap.MapID, true
	}
	return 0, 0, b.zMap.MapID, false
}

func (b *ZoneWorldBridge) GetMaster(handle uint32) uint32 { return 0 }
func (b *ZoneWorldBridge) GetAdminLevel(handle uint32) int {
	if b.zMap != nil {
		if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
			return int(p.Char.AdminLevel)
		}
	}
	return 0
}
func (b *ZoneWorldBridge) GetPlayerName(handle uint32) string {
	if b.zMap != nil {
		if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
			return p.Char.Name
		}
		if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
			return mob.Name
		}
	}
	return ""
}
func (b *ZoneWorldBridge) GetCharNo(handle uint32) int64 {
	if b.zMap != nil {
		if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
			return int64(p.Char.ID)
		}
	}
	return 0
}
func (b *ZoneWorldBridge) GetLevel(handle uint32) int {
	if b.zMap != nil {
		if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
			return int(p.Char.Level)
		}
		if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
			return int(mob.Level)
		}
	}
	return 1
}
func (b *ZoneWorldBridge) GetBaseClass(handle uint32) int {
	if b.zMap != nil {
		if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
			return int(p.Char.Job)
		}
	}
	return 1
}
func (b *ZoneWorldBridge) IsPlayerExist(handle uint32) bool {
	return b.zMap != nil && b.zMap.GetPlayer(uint16(handle)) != nil
}
func (b *ZoneWorldBridge) IsObjectDead(handle uint32) bool {
	if b.zMap == nil {
		return true
	}
	if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		mob.mu.RLock()
		defer mob.mu.RUnlock()
		return mob.IsDead
	}
	if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
		return p.Char.CurHP <= 0 || p.State == 3
	}
	return true
}
func (b *ZoneWorldBridge) IsMovable(handle uint32) bool {
	return !b.IsObjectDead(handle)
}
func (b *ZoneWorldBridge) GetWhoKillMe(handle uint32) uint32 {
	if b.zMap != nil {
		if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
			mob.mu.RLock()
			defer mob.mu.RUnlock()
			return uint32(mob.LastAttacker)
		}
	}
	return 0
}

func (b *ZoneWorldBridge) MoveTo(handle uint32, x, y float64, speed float64) {
	if b.zMap == nil {
		return
	}
	if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		if !b.zMap.IsPositionWalkable(x, y) {
			return
		}
		spd := uint16(speed)
		if spd == 0 {
			spd = mob.WalkSpeed
		}
		b.zMap.moveMonsterStep(mob, int32(x), int32(y), spd, false)
	}
}

func (b *ZoneWorldBridge) MoveStop(handle uint32) {
	if b.zMap == nil {
		return
	}
	if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		mob.mu.Lock()
		mob.IsMoving = false
		mx, my := mob.X, mob.Y
		mob.mu.Unlock()
		stop := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONESTOP_CMD)
		stop.WriteUint16(mob.Handle)
		stop.WriteInt32(mx)
		stop.WriteInt32(my)
		b.zMap.grid.BroadcastAOI(stop, int(mx), int(my), 0)
	}
}

func (b *ZoneWorldBridge) CastTeleport(handle uint32, mapName string, x, y float64) {
	if b.zMap != nil {
		b.zMap.bridgeCastTeleport(handle, mapName, x, y)
	}
}
func (b *ZoneWorldBridge) CameraMove(handle uint32, cameraID int)                   {}
func (b *ZoneWorldBridge) Animate(handle uint32, animID int, forced bool)           {}
func (b *ZoneWorldBridge) Emotion(handle uint32, emotionID int)                     {}
func (b *ZoneWorldBridge) PlaySound(soundID int, x, y float64)                      {}
func (b *ZoneWorldBridge) ObjectEffect(handle uint32, effectID int)                 {}

func (b *ZoneWorldBridge) MobRegenXY(mobID uint32, mapName string, x, y float64, layer int) uint32 {
	if b.zMap == nil {
		return 0
	}
	b.zMap.mu.RLock()
	cat := b.zMap.mobCatalog
	b.zMap.mu.RUnlock()
	if cat == nil {
		return 0
	}
	tmpl := cat.ByID[uint16(mobID)]
	if tmpl == nil {
		return 0
	}
	mob := b.zMap.SpawnMobFromTemplate(tmpl, int32(x), int32(y), 0, "", tmpl.RegenIntervalSec, true)
	if mob != nil {
		return uint32(mob.Handle)
	}
	return 0
}

func (b *ZoneWorldBridge) MobRegenCircle(mobID uint32, mapName string, cx, cy float64, radius float64) uint32 {
	angle := rand.Float64() * 2 * math.Pi
	dist := rand.Float64() * radius
	return b.MobRegenXY(mobID, mapName, cx+math.Cos(angle)*dist, cy+math.Sin(angle)*dist, 0)
}

func (b *ZoneWorldBridge) MobRegenRectangle(mobID uint32, mapName string, x1, y1, x2, y2 float64) uint32 {
	rx := x1
	if x2 > x1 {
		rx = x1 + rand.Float64()*(x2-x1)
	}
	ry := y1
	if y2 > y1 {
		ry = y1 + rand.Float64()*(y2-y1)
	}
	return b.MobRegenXY(mobID, mapName, rx, ry, 0)
}

func (b *ZoneWorldBridge) NPCRegen(npcID uint32, mapName string, x, y float64) uint32 {
	if b.zMap != nil {
		return b.zMap.bridgeNPCRegen(npcID, mapName, x, y)
	}
	return 0
}
func (b *ZoneWorldBridge) NPCVanish(handle uint32) {
	if b.zMap != nil {
		b.zMap.bridgeNPCVanish(handle)
	}
}
func (b *ZoneWorldBridge) GroupRegen(groupID uint32, mapName string) int {
	if b.zMap != nil {
		return b.zMap.bridgeGroupRegen(groupID, mapName)
	}
	return 0
}

func (b *ZoneWorldBridge) AggroSet(mobHandle, targetHandle uint32, amount int) {
	if b.zMap == nil {
		return
	}
	if mob := b.zMap.GetMonster(uint16(mobHandle)); mob != nil {
		mob.AddAggro(uint16(targetHandle), amount)
	}
}

func (b *ZoneWorldBridge) AggroReset(mobHandle uint32) {
	if b.zMap == nil {
		return
	}
	if mob := b.zMap.GetMonster(uint16(mobHandle)); mob != nil {
		mob.ResetAggro()
	}
}

func (b *ZoneWorldBridge) DoorAction(doorName string, action int) bool {
	if b.zMap != nil && b.zMap.doors != nil {
		return b.zMap.doors.DoorAction(doorName, action)
	}
	return false
}
func (b *ZoneWorldBridge) DoorBuild(doorName string) bool {
	if b.zMap != nil && b.zMap.doors != nil {
		return b.zMap.doors.DoorAction(doorName, 1)
	}
	return false
}
func (b *ZoneWorldBridge) CloseAllDoors() {
	if b.zMap != nil && b.zMap.doors != nil {
		b.zMap.doors.CloseAllDoors()
	}
}
func (b *ZoneWorldBridge) GetCurMapIndex() int {
	if b.zMap != nil {
		return int(b.zMap.ID)
	}
	return 0
}
func (b *ZoneWorldBridge) GetMapName() string {
	if b.zMap != nil {
		return b.zMap.MapID
	}
	return ""
}
func (b *ZoneWorldBridge) IsInArea(areaName string, x, y float64) bool {
	if b.zMap != nil && b.zMap.dataStore != nil && b.zMap.dataStore.BlockInfo != nil {
		if aid := b.zMap.dataStore.BlockInfo.GetAID(b.zMap.MapID); aid != nil {
			return aid.Contains(areaName, x, y)
		}
	}
	return false
}
func (b *ZoneWorldBridge) IsInMap(handle uint32, mapName string) bool {
	if b.zMap == nil || !strings.EqualFold(b.zMap.MapID, mapName) {
		return false
	}
	return b.zMap.GetEntity(uint16(handle)) != nil
}
func (b *ZoneWorldBridge) SetAbstate(handle uint32, stateID uint16, strength uint16, durationSec uint32) {
	if b.zMap == nil {
		return
	}
	st := &ActiveAbState{ID: stateID, Level: int(strength), ExpiresAt: time.Now().Add(time.Duration(durationSec) * time.Second)}
	if p := b.zMap.GetPlayer(uint16(handle)); p != nil {
		p.AddAbState(st)
	} else if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		mob.AddAbState(st)
	}
}
func (b *ZoneWorldBridge) ResetAbstate(handle uint32, stateID uint16) {
	if b.zMap == nil {
		return
	}
	if p := b.zMap.GetPlayer(uint16(handle)); p != nil {
		p.RemoveAbState(stateID)
	} else if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		mob.RemoveAbState(stateID)
	}
}
func (b *ZoneWorldBridge) GetAbstate(handle uint32, stateID uint16) (int, bool) {
	if b.zMap == nil {
		return 0, false
	}
	if p := b.zMap.GetPlayer(uint16(handle)); p != nil {
		if st := p.GetAbState(stateID); st != nil {
			return st.Level, true
		}
	} else if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		if st := mob.GetAbState(stateID); st != nil {
			return st.Level, true
		}
	}
	return 0, false
}
func (b *ZoneWorldBridge) Heal(handle uint32, amount int32) {
	if b.zMap == nil || amount <= 0 {
		return
	}
	if p := b.zMap.GetPlayer(uint16(handle)); p != nil && p.Char != nil {
		p.Char.CurHP += int(amount)
		if p.Char.CurHP > p.MaxHP {
			p.Char.CurHP = p.MaxHP
		}
	} else if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		mob.mu.Lock()
		mob.CurHP += int(amount)
		if mob.CurHP > mob.MaxHP {
			mob.CurHP = mob.MaxHP
		}
		mob.mu.Unlock()
	}
}
func (b *ZoneWorldBridge) Damaged(handle uint32, attackerHandle uint32, damage int32) {
	if b.zMap == nil || damage <= 0 {
		return
	}
	if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
		var attacker *ZonePlayer
		if attackerHandle != 0 {
			attacker = b.zMap.GetPlayer(uint16(attackerHandle))
		}
		b.zMap.ApplyDamageToMonster(mob, attacker, int(damage))
	} else if p := b.zMap.GetPlayer(uint16(handle)); p != nil {
		b.zMap.applyMonsterDamageToPlayer(nil, p, int(damage), false)
	}
}
func (b *ZoneWorldBridge) SkillBlast(handle uint32, skillID uint16, x, y float64) {}
func (b *ZoneWorldBridge) EndOfKingdomQuest(success bool) {
	if b.zMap != nil {
		b.zMap.bridgeEndOfKingdomQuest(success)
	}
}
func (b *ZoneWorldBridge) TimerStart(timerID int, seconds int) {
	if b.zMap != nil {
		b.zMap.mu.Lock()
		if b.zMap.aiTimersSec == nil {
			b.zMap.aiTimersSec = make(map[int]int64)
		}
		b.zMap.aiTimersSec[timerID] = time.Now().Unix() + int64(seconds)
		b.zMap.mu.Unlock()
	}
}
func (b *ZoneWorldBridge) TimerEnd(timerID int) {
	if b.zMap != nil {
		b.zMap.mu.Lock()
		delete(b.zMap.aiTimersSec, timerID)
		b.zMap.mu.Unlock()
	}
}
func (b *ZoneWorldBridge) ScoreBoard(recordType int, charNo int64, score int) {}
func (b *ZoneWorldBridge) SendRankingList()                                   {}
func (b *ZoneWorldBridge) Notice(msg string) {
	if b.zMap != nil {
		b.zMap.bridgeBroadcastNotice(msg)
	}
}
func (b *ZoneWorldBridge) NoticeString(msg string) {
	if b.zMap != nil {
		b.zMap.bridgeBroadcastNotice(msg)
	}
}
func (b *ZoneWorldBridge) ScriptMsg(handle uint32, msg string) {
	if b.zMap != nil {
		b.zMap.bridgeScriptMsg(handle, msg)
	}
}
func (b *ZoneWorldBridge) SysMsg(handle uint32, code int) {
	if b.zMap != nil {
		b.zMap.bridgeSysMsg(handle, uint16(code))
	}
}
func (b *ZoneWorldBridge) SomebodyShout(handle uint32, msg string) {}
func (b *ZoneWorldBridge) SendGameLogDataType4(charNo int64, mapName string, x, y, z int, int1, int2, int3 int) {
	if b.zMap != nil && b.zMap.handler != nil && b.zMap.handler.gameLogDispatcher != nil {
		b.zMap.handler.gameLogDispatcher.EnqueueType4(1, int(charNo), mapName, x, y, z, int1, int2, int3)
	}
}
func (b *ZoneWorldBridge) SendGameLogDataType5(charNo int64, int1, int2 int) {
	if b.zMap != nil && b.zMap.handler != nil && b.zMap.handler.gameLogDispatcher != nil {
		b.zMap.handler.gameLogDispatcher.EnqueueType5(1, int(charNo), int1, int2)
	}
}
func (b *ZoneWorldBridge) PetGetMind(petHandle uint32) int {
	if b.zMap != nil {
		return b.zMap.bridgePetGetMind(petHandle)
	}
	return 50
}
func (b *ZoneWorldBridge) PetChangeMind(petHandle uint32, newMind int) {
	if b.zMap != nil {
		b.zMap.bridgePetChangeMind(petHandle, newMind)
	}
}
func (b *ZoneWorldBridge) PetUnsummon(petHandle uint32) {
	if b.zMap != nil {
		b.zMap.bridgePetUnsummon(petHandle)
	}
}
func (b *ZoneWorldBridge) PetGetActionRecord(petHandle uint32, action string) (string, string, string, string, bool) {
	if b.zMap != nil {
		return b.zMap.bridgePetGetActionRecord(petHandle, action)
	}
	return "Idle", "SM_PetIdle", "", "", true
}
func (b *ZoneWorldBridge) PetGetRunSpeed(petHandle uint32) int  { return 115 }
func (b *ZoneWorldBridge) PetGetWalkSpeed(petHandle uint32) int { return 40 }
func (b *ZoneWorldBridge) PetGetStress(petHandle uint32) int {
	if b.zMap != nil {
		return b.zMap.bridgePetGetStress(petHandle)
	}
	return 0
}
func (b *ZoneWorldBridge) PetIsMasterCalling(petHandle uint32) bool { return false }
func (b *ZoneWorldBridge) PetSaveTendency(petHandle uint32) bool {
	if b.zMap != nil {
		return b.zMap.bridgePetSaveTendency(petHandle)
	}
	return true
}
func (b *ZoneWorldBridge) ServerMenu(playerHandle, npcHandle uint32, title string, options []lua.MenuOptionEntry) {
	if b.zMap != nil {
		b.zMap.bridgeServerMenu(playerHandle, npcHandle, title, options)
	}
}
func (b *ZoneWorldBridge) NPCMenuOpen(npcHandle, playerHandle uint32) {
	if b.zMap != nil {
		b.zMap.bridgeNPCMenuOpen(npcHandle, playerHandle)
	}
}
func (b *ZoneWorldBridge) QuestResult(mapName, result string) {
	if b.zMap != nil {
		b.zMap.bridgeQuestResult(mapName, result)
	}
}
func (b *ZoneWorldBridge) QuestResultIndividual(playerHandle uint32, result string) {
	if b.zMap != nil {
		b.zMap.bridgeQuestResultIndividual(playerHandle, result)
	}
}
func (b *ZoneWorldBridge) QuestMobKill(playerHandle uint32, questID uint16, mobID uint32, count int) {
	if b.zMap != nil {
		b.zMap.bridgeQuestMobKill(playerHandle, questID, mobID, count)
	}
}
func (b *ZoneWorldBridge) QuestMobKillAllInMap(mapName string, questID uint16, mobID uint32, maxCount int) {
	if b.zMap != nil {
		b.zMap.bridgeQuestMobKillAllInMap(mapName, questID, mobID, maxCount)
	}
}
func (b *ZoneWorldBridge) GetTutorialInfo(playerHandle uint32) int {
	if b.zMap != nil {
		return b.zMap.bridgeGetTutorialInfo(playerHandle)
	}
	return 0
}
func (b *ZoneWorldBridge) SaveTutorialInfo(playerHandle uint32, saveType, step int) {
	if b.zMap != nil {
		b.zMap.bridgeSaveTutorialInfo(playerHandle, saveType, step)
	}
}
func (b *ZoneWorldBridge) ProgressTutorial(playerHandle uint32, step int) {
	if b.zMap != nil {
		b.zMap.bridgeProgressTutorial(playerHandle, step)
	}
}
func (b *ZoneWorldBridge) ShowKQTimerWithLife(mapName string, seconds int) {
	if b.zMap != nil {
		b.zMap.bridgeShowKQTimerWithLife(seconds)
	}
}
func (b *ZoneWorldBridge) ShowKQTimerWithLifeObj(playerHandle uint32, seconds int) {
	if b.zMap != nil {
		b.zMap.bridgeShowKQTimerWithLifeObj(playerHandle, seconds)
	}
}
func (b *ZoneWorldBridge) GetKQLimitSecond(mapName string) int {
	if b.zMap != nil {
		return b.zMap.bridgeGetKQLimitSecond()
	}
	return 1800
}
func (b *ZoneWorldBridge) GetKQTeamType(handle uint32) int {
	if b.zMap != nil {
		return b.zMap.bridgeGetKQTeamType(handle)
	}
	return 0
}
func (b *ZoneWorldBridge) IsKQJoiner(handle uint32) bool {
	return b.zMap != nil && b.zMap.GetPlayer(uint16(handle)) != nil
}
func (b *ZoneWorldBridge) KQRewardIndex(playerHandle uint32, rewardIndex string) {
	if b.zMap != nil {
		b.zMap.bridgeKQRewardIndex(playerHandle, rewardIndex)
	}
}
func (b *ZoneWorldBridge) ScoreBoardAllInMap(mapName string, show bool, round, redWin, redScore, blueWin, blueScore int) {
	if b.zMap != nil {
		b.zMap.bridgeScoreBoardAllInMap(show, round, redWin, redScore, blueWin, blueScore)
	}
}
func (b *ZoneWorldBridge) IndunRankResult(mapName string, args []int) {
	if b.zMap != nil {
		b.zMap.bridgeIndunRankResult(args)
	}
}
func (b *ZoneWorldBridge) RewardAll(mapName, rewardType string) {
	if b.zMap != nil {
		b.zMap.bridgeRewardAll(rewardType)
	}
}
func (b *ZoneWorldBridge) RewardItem(playerHandle uint32, itemInx string, count int) {
	if b.zMap != nil {
		b.zMap.bridgeRewardItem(playerHandle, itemInx, count)
	}
}
func (b *ZoneWorldBridge) InvenItemDestroy(playerHandle uint32, itemInx string, count int) int {
	if b.zMap != nil {
		return b.zMap.bridgeInvenItemDestroy(playerHandle, itemInx, count)
	}
	return 0
}
func (b *ZoneWorldBridge) GetItemLot(playerHandle uint32, itemInx string) (int, bool) {
	if b.zMap != nil {
		return b.zMap.bridgeGetItemLot(playerHandle, itemInx)
	}
	return 0, false
}
func (b *ZoneWorldBridge) DropItem(itemInx string, mobHandle, playerHandle uint32, rate int) {
	if b.zMap != nil {
		b.zMap.bridgeDropItem(itemInx, mobHandle, playerHandle, rate)
	}
}
func (b *ZoneWorldBridge) SpawnNamedNPC(mapName, mobInx string, x, y float64, dir, scale int) uint32 {
	if b.zMap != nil {
		return b.zMap.bridgeSpawnNamedNPC(mobInx, x, y, dir)
	}
	return 0
}
func (b *ZoneWorldBridge) ScenarioSetTimeLimit(field string, seconds int) {
	if b.zMap != nil {
		b.zMap.mu.Lock()
		b.zMap.InstanceLimitSec = seconds
		b.zMap.mu.Unlock()
	}
}
func (b *ZoneWorldBridge) ScenarioSetInterrupt(field, eventName, callback string) {
	if b.zMap != nil {
		b.zMap.mu.Lock()
		if b.zMap.ScenarioInterrupts == nil {
			b.zMap.ScenarioInterrupts = make(map[string]string)
		}
		b.zMap.ScenarioInterrupts[eventName] = callback
		b.zMap.mu.Unlock()
	}
}
func (b *ZoneWorldBridge) ScenarioSetScriptFile(field, scriptName string) {}
func (b *ZoneWorldBridge) ScenarioSendChatWindow(field, speaker, dialogKey string) {
	if b.zMap != nil {
		b.zMap.bridgeScriptMsg(0, speaker+": "+dialogKey)
	}
}
func (b *ZoneWorldBridge) ScenarioEnd(field string) {
	if b.zMap != nil {
		b.zMap.bridgeScenarioEnd()
	}
}
func (b *ZoneWorldBridge) ScenarioFindPlayer(field, charName, mode string) uint32 {
	if b.zMap != nil {
		for _, p := range b.zMap.GetPlayers() {
			if p.Char != nil && strings.EqualFold(p.Char.Name, charName) {
				return uint32(p.MapObjectID)
			}
		}
	}
	return 0
}
func (b *ZoneWorldBridge) ScenarioCompleteWedding(groomHandle, brideHandle uint32) {
	if b.zMap != nil {
		b.zMap.bridgeCompleteWedding(groomHandle, brideHandle)
	}
}
func (b *ZoneWorldBridge) ScenarioApplyAbnormalState(handle uint32, staName string, level, durationMs int) {
	if b.zMap != nil {
		b.zMap.bridgeApplyAbstateByName(handle, staName, level, durationMs)
	}
}
func (b *ZoneWorldBridge) ScenarioSetNPCInteractHandler(npcHandle uint32, callback string) {
	if b.zMap != nil {
		b.zMap.mu.Lock()
		if b.zMap.NPCInteractCB == nil {
			b.zMap.NPCInteractCB = make(map[uint32]string)
		}
		b.zMap.NPCInteractCB[npcHandle] = callback
		b.zMap.mu.Unlock()
	}
}
func (b *ZoneWorldBridge) ScenarioRegisterCaptureFlag(field string, npcHandle uint32, redEffect, blueEffect string) {
	if b.zMap != nil {
		b.zMap.bridgeRegisterGTFlag(npcHandle, redEffect, blueEffect)
	}
}
func (b *ZoneWorldBridge) ScenarioRegisterGTPortal(field string, npcHandle, flagHandle uint32, redName, blueName string) {
	if b.zMap != nil {
		b.zMap.bridgeRegisterGTPortal(npcHandle, flagHandle)
	}
}
func (b *ZoneWorldBridge) ScenarioRegisterGuildGate(field string, npcHandle uint32, guildName string, flagHandle uint32, gateName string) {
	if b.zMap != nil {
		b.zMap.bridgeRegisterGTGate(npcHandle, guildName, flagHandle)
	}
}
func (b *ZoneWorldBridge) ScenarioBroadcastToField(field, msg string) {
	if b.zMap != nil {
		b.zMap.bridgeBroadcastNotice(msg)
	}
}
func (b *ZoneWorldBridge) ScenarioEvaluateMatchScore(field string) string {
	if b.zMap != nil {
		return b.zMap.bridgeEvaluateGTMatchScore()
	}
	return ""
}
func (b *ZoneWorldBridge) ScenarioRewardGuildTournament(field, winningGuild string) {
	if b.zMap != nil {
		b.zMap.bridgeRewardGuildTournament(winningGuild)
	}
}
func (b *ZoneWorldBridge) ScenarioWait(seconds int) {
	// No bloquea invocaciones síncronas dentro del GameLoop o tests unitarios
}

// ============================================================================
// Implementación LuaFieldRuntimeBridge (firmas oficiales Zone.exe)
// ============================================================================

func (b *ZoneWorldBridge) MobRegenXYByName(mapName, mobInx string, x, y float64, dir, layer int) uint32 {
	if b.zMap != nil {
		return b.zMap.bridgeMobRegenXYByName(mapName, mobInx, x, y, dir, layer)
	}
	return 0
}
func (b *ZoneWorldBridge) MobRegenCircleByName(mapName, mobInx string, x, y, radius float64) uint32 {
	if b.zMap != nil {
		angle := rand.Float64() * 2 * math.Pi
		dist := rand.Float64() * radius
		return b.zMap.bridgeMobRegenXYByName(mapName, mobInx, x+math.Cos(angle)*dist, y+math.Sin(angle)*dist, 0, 0)
	}
	return 0
}
func (b *ZoneWorldBridge) MobRegenRectByName(mapName, mobInx string, x1, y1, x2, y2 float64) uint32 {
	if b.zMap != nil {
		rx := x1
		if x2 > x1 {
			rx = x1 + rand.Float64()*(x2-x1)
		}
		ry := y1
		if y2 > y1 {
			ry = y1 + rand.Float64()*(y2-y1)
		}
		return b.zMap.bridgeMobRegenXYByName(mapName, mobInx, rx, ry, 0, 0)
	}
	return 0
}
func (b *ZoneWorldBridge) GroupRegenByName(mapName, groupInx string) int {
	if b.zMap != nil {
		return b.zMap.bridgeGroupRegenByName(mapName, groupInx, 0, 0, 0, 0)
	}
	return 0
}
func (b *ZoneWorldBridge) GroupRegenInstanceXY(mapName, groupInx string, x, y, radius float64, maxCount int) int {
	if b.zMap != nil {
		return b.zMap.bridgeGroupRegenByName(mapName, groupInx, x, y, radius, maxCount)
	}
	return 0
}
func (b *ZoneWorldBridge) NPCRegenByName(mapName, npcInx string, x, y float64, dir int) uint32 {
	if b.zMap != nil {
		return b.zMap.bridgeSpawnNamedNPC(npcInx, x, y, dir)
	}
	return 0
}
func (b *ZoneWorldBridge) MobSuicideByMap(mapName string, handle uint32) {
	if b.zMap != nil {
		b.zMap.bridgeMobSuicideByMap(mapName, handle)
	}
}
func (b *ZoneWorldBridge) VanishAll(mapName, mobInx string) {
	if b.zMap != nil {
		b.zMap.bridgeVanishAll(mapName, mobInx)
	}
}
func (b *ZoneWorldBridge) VanishReserv(mapName, mobInx string, seconds int) {
	if b.zMap != nil {
		b.zMap.mu.Lock()
		if b.zMap.VanishTimers == nil {
			b.zMap.VanishTimers = make(map[string]int64)
		}
		b.zMap.VanishTimers[mobInx] = time.Now().Unix() + int64(seconds)
		b.zMap.mu.Unlock()
	}
}
func (b *ZoneWorldBridge) VanishTimer(mapName string, seconds int) {
	if b.zMap != nil {
		b.zMap.mu.Lock()
		if b.zMap.VanishTimers == nil {
			b.zMap.VanishTimers = make(map[string]int64)
		}
		b.zMap.VanishTimers["*"] = time.Now().Unix() + int64(seconds)
		b.zMap.mu.Unlock()
	}
}
func (b *ZoneWorldBridge) ObjectCount(mapName string, objType int) int {
	if b.zMap != nil {
		return b.zMap.bridgeObjectCount(mapName, objType)
	}
	return 0
}
func (b *ZoneWorldBridge) GetPlayerListInMap(mapName string) []uint32 {
	if b.zMap != nil {
		return b.zMap.bridgeGetPlayerList(mapName)
	}
	return nil
}
func (b *ZoneWorldBridge) LinkToAll(srcMap, destMap string, x, y float64) {
	if b.zMap != nil {
		b.zMap.bridgeLinkToAll(srcMap, destMap, x, y)
	}
}
func (b *ZoneWorldBridge) DoorBuildByField(fieldName, mobInx string, x, y float64, dir int, scale float64) uint32 {
	if b.zMap != nil {
		return b.zMap.bridgeDoorBuildByField(fieldName, mobInx, x, y, dir, scale)
	}
	return 0
}
func (b *ZoneWorldBridge) DoorActionByField(fieldOrHandle interface{}, blockName, action string) bool {
	if b.zMap != nil {
		return b.zMap.bridgeDoorActionByField(fieldOrHandle, blockName, action)
	}
	return false
}
func (b *ZoneWorldBridge) IsInAreaByMap(handle uint32, mapName, areaName string) bool {
	if b.zMap != nil {
		return b.zMap.bridgeIsInAreaByMap(handle, mapName, areaName)
	}
	return false
}
func (b *ZoneWorldBridge) SetAbstateByName(handle uint32, staName string, level, durationMs int) {
	if b.zMap != nil {
		b.zMap.bridgeApplyAbstateByName(handle, staName, level, durationMs)
	}
}
func (b *ZoneWorldBridge) SetAbstateInArea(mapName, staName string, level, durationMs int, areaName string, cx, cy, radius float64) {
	if b.zMap != nil {
		b.zMap.bridgeSetAbstateInArea(mapName, staName, level, durationMs, areaName, cx, cy, radius)
	}
}
func (b *ZoneWorldBridge) EndOfKingdomQuestMap(mapName string) bool {
	if b.zMap != nil {
		b.zMap.bridgeEndOfKingdomQuest(true)
		return true
	}
	return false
}
func (b *ZoneWorldBridge) RegisterFieldScriptFunc(mapName, eventName, funcName string) bool {
	if b.zMap != nil {
		b.zMap.mu.Lock()
		if b.zMap.FieldScriptFuncs == nil {
			b.zMap.FieldScriptFuncs = make(map[string]string)
		}
		b.zMap.FieldScriptFuncs[eventName] = funcName
		b.zMap.mu.Unlock()
		return true
	}
	return false
}
func (b *ZoneWorldBridge) RegisterAIScriptFunc(handle uint32, eventName, funcName string) bool {
	if b.zMap != nil && handle != 0 {
		b.zMap.mu.Lock()
		if b.zMap.AIScriptFuncs == nil {
			b.zMap.AIScriptFuncs = make(map[uint32]map[string]string)
		}
		if b.zMap.AIScriptFuncs[handle] == nil {
			b.zMap.AIScriptFuncs[handle] = make(map[string]string)
		}
		b.zMap.AIScriptFuncs[handle][eventName] = funcName
		b.zMap.mu.Unlock()
		if strings.EqualFold(eventName, "Entrance") || strings.EqualFold(eventName, "Routine") {
			if mob := b.zMap.GetMonster(uint16(handle)); mob != nil {
				mob.mu.Lock()
				mob.AIScriptFunc = funcName
				mob.mu.Unlock()
			}
		}
		return true
	}
	return false
}
func (b *ZoneWorldBridge) TimerStartMap(mapName string, timerID, seconds int) {
	b.TimerStart(timerID, seconds)
}
func (b *ZoneWorldBridge) TimerEndMap(mapName string, timerID int) {
	b.TimerEnd(timerID)
}

func resolveZoneDataDir() string {
	for _, cand := range []string{"data", filepath.Join("..", "..", "..", "data"), filepath.Join("..", "..", "data")} {
		if fi, err := os.Stat(filepath.Join(cand, "MobBehaviorDescript")); err == nil && fi.IsDir() {
			return cand
		}
		if fi, err := os.Stat(filepath.Join(cand, "LuaScript")); err == nil && fi.IsDir() {
			return cand
		}
	}
	return "data"
}

// initMapLuaAI carga DefaultBehavior.lua y los scripts de MobBehaviorDescript / AIScript para el mapa.
func (m *ZoneMap) initMapLuaAI() {
	m.mu.Lock()
	if m.luaEngine != nil {
		m.mu.Unlock()
		return
	}
	dataRoot := resolveZoneDataDir()
	bridge := NewZoneWorldBridge(m)
	eng, err := lua.NewEngine(lua.Config{
		BasePath: filepath.Join(dataRoot, "LuaScript"),
		Bridge:   bridge,
	})
	if err != nil {
		m.mu.Unlock()
		return
	}
	m.worldBridge = bridge
	m.luaEngine = eng
	m.mu.Unlock()

	defaultLua := filepath.Join(dataRoot, "MobBehaviorDescript", "DefaultBehavior.lua")
	if _, err := os.Stat(defaultLua); err == nil {
		_ = eng.DoFile(defaultLua)
	}
	kqLua := filepath.Join(dataRoot, "MobBehaviorDescript", "KQ", m.MapID+".lua")
	if _, err := os.Stat(kqLua); err == nil {
		_ = eng.DoFile(kqLua)
	}
	// Precargar scripts oficiales de AIScript comunes si están presentes
	if matches, err := filepath.Glob(filepath.Join(dataRoot, "LuaScript", "AIScript", "*.lua")); err == nil {
		for _, p := range matches {
			_ = eng.DoFile(p)
		}
	}
}

// RunMobBehaviorLua ejecuta el script Lua oficial de comportamiento para un monstruo (MobBehaviorDescript / AIScript).
func (m *ZoneMap) RunMobBehaviorLua(mob *ZoneMonster, funcName string) (handled bool) {
	if m == nil || mob == nil || funcName == "" {
		return false
	}
	m.mu.RLock()
	eng := m.luaEngine
	m.mu.RUnlock()
	if eng == nil {
		return false
	}

	rets, err := eng.CallFunction(funcName, glua.LNumber(mob.Handle), glua.LNumber(m.ID))
	if err != nil {
		return false
	}
	if len(rets) > 0 {
		if num, ok := rets[0].(glua.LNumber); ok && int(num) == 1 {
			// ReturnAI.END = 1
			return true
		}
	}
	return false
}

// BroadcastMobChat emite un diálogo de monstruo desde data/World/Misc/MobChat.json si existe para ese evento.
func (m *ZoneMap) BroadcastMobChat(mob *ZoneMonster, category string) {
	if m == nil || mob == nil {
		return
	}
	m.mu.RLock()
	ds := m.dataStore
	m.mu.RUnlock()
	if ds == nil || ds.MiscStore == nil {
		return
	}
	chats := ds.MiscStore.MobChats()
	var entries []misc.MobChatEntry
	switch strings.ToLower(category) {
	case "peace":
		entries = chats.Peace
	case "attack":
		entries = chats.Attack
	case "damaged":
		entries = chats.Damaged
	case "dead":
		entries = chats.Dead
	case "helpmain":
		entries = chats.HelpMain
	case "helpsub":
		entries = chats.HelpSub
	case "racerescue":
		entries = chats.RaceRescue
	}
	var matched []misc.MobChatEntry
	for _, e := range entries {
		if strings.EqualFold(e.Mob, mob.InxName) && len(e.Messages) > 0 {
			matched = append(matched, e)
		}
	}
	if len(matched) == 0 {
		return
	}
	chosen := matched[rand.Intn(len(matched))]
	text := chosen.Messages[rand.Intn(len(chosen.Messages))]
	if text == "" {
		return
	}
	mx, my := mob.GetPosition()
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONECHAT_CMD)
	pw.WriteByte(0) // itemLinkDataCount
	pw.WriteUint16(mob.Handle)
	pw.WriteByte(byte(len(text)))
	pw.WriteByte(0x02) // flag: chatwin=1, GMColor=0
	pw.WriteByte(0)    // nChatFontColorID
	pw.WriteByte(0)    // nChatBalloonColorID
	pw.WriteString(text, len(text))
	m.grid.BroadcastAOI(pw, mx, my, 0)
}

// TickMonsters100ms ejecuta la máquina de estados AI de todos los monstruos del mapa en cada tick de 100ms.
func (m *ZoneMap) TickMonsters100ms(now time.Time) {
	mobs := m.GetAllMonsters()
	if len(mobs) == 0 {
		return
	}

	m.mu.RLock()
	ds := m.dataStore
	cat := m.mobCatalog
	m.mu.RUnlock()

	for _, mob := range mobs {
		if mob == nil {
			continue
		}

		mob.mu.RLock()
		isDead := mob.IsDead
		respawnAt := mob.RespawnTime
		aiScript := mob.AIScriptFunc
		mob.mu.RUnlock()

		// 1. Ciclo de Respawn cuando está muerto
		if isDead {
			if !respawnAt.IsZero() && !now.Before(respawnAt) {
				m.respawnMonster(mob)
			}
			continue
		}

		// 2. Hook opcional de AIScript Lua oficial
		if aiScript != "" {
			if m.RunMobBehaviorLua(mob, aiScript) {
				continue
			}
		}

		// 3. Evaluar tabla de aggro y triggers
		mob.mu.RLock()
		hasAggro := len(mob.Aggro) > 0
		aiState := mob.AIState
		mob.mu.RUnlock()

		if !hasAggro && aiState != MobAIReturn {
			// A. Detección agresiva respetando StealthMatrix.json + LoS .shbd
			if m.scanAggroTarget(mob, ds) {
				continue
			}
			// B. Patrullaje MobRoam o desplazamiento pasivo
			m.tickMobRoam(mob, now)
			continue
		}

		if aiState == MobAIReturn {
			m.tickMobReturn(mob, now)
			continue
		}

		// 4. Combate y persecución (Chase / Attack)
		m.tickMobCombat(mob, now, ds, cat)
	}
}

// TickMonsters1s ejecuta recuperación de HP fuera de combate y ticks de AbState de monstruos.
func (m *ZoneMap) TickMonsters1s(now time.Time) {
	mobs := m.GetAllMonsters()
	for _, mob := range mobs {
		if mob == nil {
			continue
		}
		mob.mu.Lock()
		if mob.IsDead {
			mob.mu.Unlock()
			continue
		}

		// Recuperación de HP en paz (PceHPRcv cada PceHPRcvDly ms)
		if len(mob.Aggro) == 0 && mob.PeaceHPRecov > 0 && mob.CurHP < mob.MaxHP {
			dly := time.Duration(mob.PeaceHPRecovDly) * time.Millisecond
			if dly <= 0 {
				dly = 3 * time.Second
			}
			if mob.LastPeaceHPAt.IsZero() || now.Sub(mob.LastPeaceHPAt) >= dly {
				mob.CurHP += mob.PeaceHPRecov
				if mob.CurHP > mob.MaxHP {
					mob.CurHP = mob.MaxHP
				}
				mob.LastPeaceHPAt = now
			}
		}

		// Ticks de DoT y expiración de AbState
		var dotTotal int
		var expiredAbStates []uint16
		for id, st := range mob.ActiveStates {
			if st.DotDamage > 0 && (st.LastDotTick.IsZero() || now.Sub(st.LastDotTick) >= st.DotInterval) {
				st.LastDotTick = now
				dotTotal += st.DotDamage
			}
			if !st.ExpiresAt.IsZero() && now.After(st.ExpiresAt) {
				delete(mob.ActiveStates, id)
				expiredAbStates = append(expiredAbStates, id)
			}
		}
		lastAttackerHandle := mob.LastAttacker
		mob.mu.Unlock()

		for _, id := range expiredAbStates {
			resetCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATERESET_CMD)
			resetCmd.WriteUint16(mob.Handle)
			resetCmd.WriteUint32(uint32(id))
			m.Broadcast(resetCmd, 0)
		}

		if dotTotal > 0 {
			var killer *ZonePlayer
			if lastAttackerHandle != 0 {
				killer = m.GetPlayer(lastAttackerHandle)
			}
			m.ApplyDamageToMonster(mob, killer, dotTotal)
		}
	}
}

func (m *ZoneMap) respawnMonster(mob *ZoneMonster) {
	m.mu.RLock()
	rg := m.regenGroups[mob.RegenGroupID]
	m.mu.RUnlock()

	sx, sy := mob.SpawnX, mob.SpawnY
	if rg != nil {
		sx, sy = m.pickWalkableCoordInGroup(rg.Def)
	}

	mob.mu.Lock()
	mob.IsDead = false
	mob.CurHP = mob.MaxHP
	mob.CurSP = mob.MaxSP
	mob.X = sx
	mob.Y = sy
	mob.Direct = mob.SpawnDirect
	mob.RespawnTime = time.Time{}
	mob.Aggro = make(map[uint16]int)
	mob.DamageContrib = make(map[int]int)
	mob.FiredTriggers = make(map[int]bool)
	mob.ActiveStates = make(map[uint16]*ActiveAbState)
	mob.TargetHandle = 0
	mob.LastAttacker = 0
	mob.LastKillerChar = 0
	mob.AttackSeqIndex = 0
	mob.HelpCalled = false
	mob.IsMoving = false
	if mob.Route != nil || mob.IsRoaming {
		mob.AIState = MobAIRoam
	} else {
		mob.AIState = MobAIIdle
	}
	mob.mu.Unlock()

	if m.grid != nil {
		m.grid.AddEntity(mob)
	}
	m.BroadcastMobSpawn(mob, true)
}

func (m *ZoneMap) scanAggroTarget(mob *ZoneMonster, ds *DataStore) bool {
	mob.mu.RLock()
	aggressive := mob.EnemyDetectType > 0
	detectRange := mob.DetectRange
	mx, my := mob.X, mob.Y
	mob.mu.RUnlock()

	if !aggressive || detectRange <= 0 {
		return false
	}

	players := m.grid.GetSurroundingPlayers(int(mx), int(my))
	var chosen *ZonePlayer
	bestDist := float64(detectRange + 1)

	for _, p := range players {
		if p == nil || p.Char == nil || p.Char.CurHP <= 0 || p.State == 3 {
			continue
		}
		px, py := p.GetPosition()
		dist := math.Hypot(float64(px)-float64(mx), float64(py)-float64(my))
		if dist > float64(detectRange) {
			continue
		}
		// Validar sigilo mediante StealthMatrix.json
		if ds != nil && ds.MiscStore != nil {
			// Detecta estados de sigilo en jugador si posee Hide/Stealth
			if _, ok := p.ActiveStates[1]; ok {
				if !ds.MiscStore.CanSee("Mob", "Hide") {
					continue
				}
			}
		}
		// Validar LoS en .shbd
		if m.blockInfo != nil && m.blockInfo.IsWorldAttackBlocked(float64(mx), float64(my), float64(px), float64(py)) {
			continue
		}
		if dist < bestDist {
			bestDist = dist
			chosen = p
		}
	}

	if chosen != nil {
		mob.AddAggro(chosen.MapObjectID, 10)
		m.BroadcastMobChat(mob, "Attack")
		return true
	}
	return false
}

func (m *ZoneMap) tickMobRoam(mob *ZoneMonster, now time.Time) {
	mob.mu.Lock()
	if !mob.RoamWaitUntil.IsZero() && now.Before(mob.RoamWaitUntil) {
		mob.mu.Unlock()
		return
	}
	if !mob.LastMoveTick.IsZero() && now.Sub(mob.LastMoveTick) < 300*time.Millisecond {
		mob.mu.Unlock()
		return
	}
	mob.LastMoveTick = now

	var targetX, targetY int32
	if mob.Route != nil && len(mob.Route.Points) > 0 {
		if mob.RoamIndex < 0 || mob.RoamIndex >= len(mob.Route.Points) {
			mob.RoamIndex = 0
		}
		wp := mob.Route.Points[mob.RoamIndex]
		targetX, targetY = int32(wp.X), int32(wp.Y)
		dist := math.Hypot(float64(targetX-mob.X), float64(targetY-mob.Y))
		if dist <= 25.0 {
			// Llegó al waypoint: avanzar índice según LoopMode (rotate vs return)
			restMs := mob.RoamRestTimeMs
			if restMs <= 0 {
				restMs = 800
			}
			mob.RoamWaitUntil = now.Add(time.Duration(restMs) * time.Millisecond)

			mode := strings.ToLower(mob.Route.LoopMode)
			if mode == "return" {
				if !mob.RoamReverse {
					if mob.RoamIndex+1 < len(mob.Route.Points) {
						mob.RoamIndex++
					} else if len(mob.Route.Points) > 1 {
						mob.RoamReverse = true
						mob.RoamIndex--
					}
				} else {
					if mob.RoamIndex-1 >= 0 {
						mob.RoamIndex--
					} else if len(mob.Route.Points) > 1 {
						mob.RoamReverse = false
						mob.RoamIndex++
					}
				}
			} else {
				mob.RoamIndex = (mob.RoamIndex + 1) % len(mob.Route.Points)
			}
			mob.mu.Unlock()
			return
		}
	} else if mob.IsRoaming && mob.RoamDistance > 0 {
		if !mob.IsMoving {
			angle := rand.Float64() * 2 * math.Pi
			r := float64(mob.RoamDistance) * (0.3 + 0.7*rand.Float64())
			candX := mob.SpawnX + int32(math.Cos(angle)*r)
			candY := mob.SpawnY + int32(math.Sin(angle)*r)
			if !m.IsPositionWalkable(float64(candX), float64(candY)) {
				mob.mu.Unlock()
				return
			}
			mob.DestX = candX
			mob.DestY = candY
			mob.IsMoving = true
		}
		targetX, targetY = mob.DestX, mob.DestY
		if math.Hypot(float64(targetX-mob.X), float64(targetY-mob.Y)) <= 20.0 {
			mob.IsMoving = false
			restMs := mob.RoamRestTimeMs
			if restMs <= 0 {
				restMs = 1500
			}
			mob.RoamWaitUntil = now.Add(time.Duration(restMs) * time.Millisecond)
			mob.mu.Unlock()
			return
		}
	} else {
		mob.mu.Unlock()
		return
	}

	nx, ny := stepToward(mob.X, mob.Y, targetX, targetY, float64(mob.WalkSpeed)*0.35)
	walkSpeed := mob.WalkSpeed
	mob.mu.Unlock()

	if !m.IsPositionWalkable(float64(nx), float64(ny)) {
		mob.mu.Lock()
		mob.IsMoving = false
		mob.mu.Unlock()
		return
	}
	m.moveMonsterStep(mob, nx, ny, walkSpeed, false)
}

func (m *ZoneMap) tickMobReturn(mob *ZoneMonster, now time.Time) {
	mob.mu.Lock()
	if !mob.LastMoveTick.IsZero() && now.Sub(mob.LastMoveTick) < 200*time.Millisecond {
		mob.mu.Unlock()
		return
	}
	mob.LastMoveTick = now
	sx, sy := mob.SpawnX, mob.SpawnY
	mx, my := mob.X, mob.Y
	speed := mob.RunSpeed
	if speed == 0 {
		speed = 115
	}
	dist := math.Hypot(float64(sx-mx), float64(sy-my))
	if dist <= 30.0 {
		mob.X = sx
		mob.Y = sy
		mob.CurHP = mob.MaxHP
		mob.CurSP = mob.MaxSP
		mob.Aggro = make(map[uint16]int)
		mob.DamageContrib = make(map[int]int)
		mob.FiredTriggers = make(map[int]bool)
		mob.TargetHandle = 0
		mob.HelpCalled = false
		mob.IsMoving = false
		if mob.Route != nil || mob.IsRoaming {
			mob.AIState = MobAIRoam
		} else {
			mob.AIState = MobAIIdle
		}
		mob.mu.Unlock()
		stop := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_SOMEONESTOP_CMD)
		stop.WriteUint16(mob.Handle)
		stop.WriteInt32(sx)
		stop.WriteInt32(sy)
		m.grid.BroadcastAOI(stop, int(sx), int(sy), 0)
		return
	}
	nx, ny := stepToward(mx, my, sx, sy, float64(speed)*0.35)
	mob.mu.Unlock()

	m.moveMonsterStep(mob, nx, ny, speed, true)
}

func (m *ZoneMap) tickMobCombat(mob *ZoneMonster, now time.Time, ds *DataStore, cat *MobCatalog) {
	mob.mu.RLock()
	targetH := mob.TargetHandle
	mx, my := mob.X, mob.Y
	sx, sy := mob.SpawnX, mob.SpawnY
	followRange := mob.FollowRange
	returnToRegen := mob.ReturnToRegen
	atkRange := mob.AttackRange
	atkDelay := time.Duration(mob.AttackDelayMs) * time.Millisecond
	lastAtk := mob.LastAttackAt
	runSpeed := mob.RunSpeed
	hp := mob.CurHP
	maxHP := mob.MaxHP
	helpCalled := mob.HelpCalled
	familyCount := mob.FamilyRescueCnt
	familyArea := mob.FamilyRescueArea
	triggers := mob.Triggers
	mob.mu.RUnlock()

	if atkRange <= 0 {
		atkRange = 120
	}
	if atkDelay <= 0 {
		atkDelay = 1200 * time.Millisecond
	}
	if runSpeed == 0 {
		runSpeed = 105
	}

	target := m.GetPlayer(targetH)
	if target == nil || target.Char == nil || target.Char.CurHP <= 0 || target.State == 3 {
		nextH := mob.RemoveAggroTarget(targetH)
		if nextH == 0 {
			mob.mu.Lock()
			if returnToRegen {
				mob.AIState = MobAIReturn
			} else {
				mob.AIState = MobAIIdle
			}
			mob.mu.Unlock()
		}
		return
	}

	px, py := target.GetPosition()
	leashDist := math.Hypot(float64(mx-sx), float64(my-sy))
	targetLeashDist := math.Hypot(float64(int32(px)-sx), float64(int32(py)-sy))
	if returnToRegen && followRange > 0 && (leashDist > float64(followRange) || targetLeashDist > float64(followRange*2)) {
		mob.ResetAggro()
		mob.mu.Lock()
		mob.AIState = MobAIReturn
		mob.mu.Unlock()
		return
	}

	// Llamada de auxilio de familia si HP < 50% o según trigger
	// Verificado en Zone.exe sm_DuringBattle (0x00568E36): umbral ((MaxHP * 0x19a) >> 10) (~40% HP).
	if !helpCalled && familyCount > 0 && maxHP > 0 && hp <= (maxHP*410)/1024 {
		mob.mu.Lock()
		mob.HelpCalled = true
		mob.mu.Unlock()
		m.callFamilyHelp(mob, target, familyArea, familyCount)
	}

	// Evaluar triggers de MobSetting.json
	if len(triggers) > 0 && maxHP > 0 {
		hpPct := hp * 100 / maxHP
		for idx, tr := range triggers {
			mob.mu.RLock()
			fired := mob.FiredTriggers[idx]
			mob.mu.RUnlock()
			if fired {
				continue
			}
			// Disparar llamada de auxilio cuando ConditionID está activo y HP cae debajo de 50%
			if tr.ConditionID > 0 && hpPct <= 50 {
				mob.mu.Lock()
				mob.FiredTriggers[idx] = true
				mob.mu.Unlock()
				m.callFamilyHelp(mob, target, 400, 1)
			}
		}
	}

	distToTarget := math.Hypot(float64(px)-float64(mx), float64(py)-float64(my))
	if distToTarget > float64(atkRange) {
		mob.mu.Lock()
		if !mob.LastMoveTick.IsZero() && now.Sub(mob.LastMoveTick) < 200*time.Millisecond {
			mob.mu.Unlock()
			return
		}
		mob.LastMoveTick = now
		mob.AIState = MobAIChase
		nx, ny := stepToward(mx, my, int32(px), int32(py), float64(runSpeed)*0.3)
		mob.mu.Unlock()

		if m.IsPositionWalkable(float64(nx), float64(ny)) {
			m.moveMonsterStep(mob, nx, ny, runSpeed, true)
		}
		return
	}

	// En rango de ataque: respetar AttackDelay
	if !lastAtk.IsZero() && now.Sub(lastAtk) < atkDelay {
		return
	}
	m.executeMonsterAttack(mob, target, now, ds, cat)
}

func (m *ZoneMap) callFamilyHelp(caller *ZoneMonster, target *ZonePlayer, radius, maxAllies int) {
	if caller == nil || target == nil || maxAllies <= 0 {
		return
	}
	if radius <= 0 {
		radius = 450
	}
	m.BroadcastMobChat(caller, "HelpMain")
	cx, cy := caller.GetPosition()
	recruited := 0
	for _, ally := range m.GetAllMonsters() {
		if recruited >= maxAllies {
			break
		}
		if ally == nil || ally.Handle == caller.Handle {
			continue
		}
		ally.mu.RLock()
		dead := ally.IsDead
		busy := len(ally.Aggro) > 0
		sameFamily := strings.EqualFold(ally.InxName, caller.InxName) || (ally.RaceType != 0 && ally.RaceType == caller.RaceType)
		ax, ay := ally.X, ally.Y
		ally.mu.RUnlock()

		if dead || busy || !sameFamily {
			continue
		}
		if math.Hypot(float64(ax)-float64(cx), float64(ay)-float64(cy)) <= float64(radius) {
			ally.AddAggro(target.MapObjectID, 5)
			m.BroadcastMobChat(ally, "HelpSub")
			recruited++
		}
	}
}

func asPermilleInt(v any) int {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case int:
		return t
	case float64:
		return int(t)
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(t))
		if err == nil {
			return n
		}
	}
	n, _ := strconv.Atoi(strings.TrimSpace(fmt.Sprint(v)))
	return n
}

func (m *ZoneMap) executeMonsterAttack(mob *ZoneMonster, target *ZonePlayer, now time.Time, ds *DataStore, cat *MobCatalog) {
	mob.mu.Lock()
	mob.LastAttackAt = now
	mob.AIState = MobAIAttack
	skillName := ""
	if mob.MobAI != nil {
		// Verificar cambios de skill condicionales (por HP permille/porcentaje)
		if mob.MaxHP > 0 {
			hpPermille := mob.CurHP * 1000 / mob.MaxHP
			for _, sc := range mob.MobAI.SkillChanges {
				thresh := asPermilleInt(sc.Value)
				if thresh > 0 && thresh <= 100 {
					thresh *= 10
				}
				if strings.Contains(strings.ToLower(sc.Condition), "hp") && thresh > 0 && hpPermille <= thresh && sc.To != "" && sc.To != "-" {
					skillName = sc.To
					break
				}
			}
		}
		// Si no aplicó SkillChange, rotar por Sequence
		if skillName == "" && len(mob.MobAI.Sequence) > 0 {
			ord := mob.MobAI.Sequence[mob.AttackSeqIndex%len(mob.MobAI.Sequence)]
			mob.AttackSeqIndex++
			if ord.Attack != "" && ord.Attack != "-" && !strings.EqualFold(ord.Attack, "Attack") {
				skillName = ord.Attack
			}
		}
	}
	minWC := mob.MinWC
	maxWC := mob.MaxWC
	th := mob.TH
	mob.mu.Unlock()

	var skillDef *MobSkillDef
	if skillName != "" && cat != nil {
		skillDef = cat.SkillsByInx[strings.ToLower(skillName)]
	}

	// Emisión de inicio de swing o skill
	swingStart := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SWING_START_CMD)
	swingStart.WriteUint16(mob.Handle)
	swingStart.WriteUint16(target.MapObjectID)
	swingStart.WriteUint8(1)
	swingStart.WriteUint16(uint16(mob.AttackDelayMs))
	swingStart.WriteUint8(1)
	swingStart.WriteUint8(1)
	m.Broadcast(swingStart, 0)

	// Evaluar impacto vs Evasión (TB) del jugador
	hitBonus := int(th) + 30
	if skillDef != nil && skillDef.PyHitRate > 0 {
		hitBonus += skillDef.PyHitRate / 100
	}
	if !CheckHit(hitBonus, int(target.EquipTB)) {
		var curHP uint32
		if target.Char != nil {
			curHP = uint32(target.Char.CurHP)
		}
		emitSwingDamage(m, mob, target, mob.Handle, target.MapObjectID, 0, curHP, false, true, false, 1)
		return
	}

	rawMin := int(minWC)
	rawMax := int(maxWC)
	if skillDef != nil && skillDef.MaxWC > 0 {
		rawMin += int(skillDef.MinWC)
		rawMax += int(skillDef.MaxWC)
	}
	if rawMax <= rawMin {
		rawMax = rawMin + 2
	}
	atkRoll := rawMin + rand.Intn(rawMax-rawMin+1)
	px, py := target.GetPosition()
	mx, my := mob.GetPosition()
	angle := CalculateAngleDeg(float64(mx), float64(my), float64(px), float64(py), 0)
	isCrit := CheckCrit(5, 0)
	var combatStore = dsCombatStore(ds)
	dmg := CalculatePhysicalDamage(atkRoll, int(target.EquipAC), angle, false, 0, 1000, isCrit, combatStore)
	if dmg < 1 {
		dmg = 1
	}

	m.applyMonsterDamageToPlayer(mob, target, dmg, isCrit)
}

func dsCombatStore(ds *DataStore) interface {
	GetAngleDamageRate(isMonster bool, angle int) int
	GetSoulDamageRate(demandSoul int, soulCount int) int
} {
	if ds != nil && ds.CombatStore != nil {
		return ds.CombatStore
	}
	return nil
}

// Reemplazamos con tipado directo en la llamada si CombatStore es nil-safe; usemos helper con *combat.Store.
func (m *ZoneMap) applyMonsterDamageToPlayer(mob *ZoneMonster, target *ZonePlayer, dmg int, isCrit bool) {
	if target == nil || target.Char == nil || dmg <= 0 {
		return
	}
	if target.Char.CurHP <= dmg {
		target.Char.CurHP = 0
		target.State = 3
		target.IsAutoAttacking = false
		target.TargetHandle = 0
	} else {
		target.Char.CurHP -= dmg
	}
	restHP := uint32(target.Char.CurHP)
	isDead := target.Char.CurHP <= 0

	var attackerH uint16
	var mobID uint16
	if mob != nil {
		attackerH = mob.Handle
		mobID = mob.MobID
	}

	emitSwingDamage(m, mob, target, attackerH, target.MapObjectID, uint16(dmg), restHP, isCrit, false, isDead, 1)

	if target.Client != nil {
		hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
		hpCmd.WriteUint32(restHP)
		hpCmd.WriteUint16(0)
		_ = target.Client.SendPacket(hpCmd)
	}

	if isDead {
		if mob != nil {
			mob.RemoveAggroTarget(target.MapObjectID)
		}

		m.mu.RLock()
		h := m.handler
		m.mu.RUnlock()
		if h != nil && h.gameLogDispatcher != nil {
			px, py := target.GetPosition()
			h.gameLogDispatcher.EnqueueType8(8, target.Char.ID, m.MapID, px, py, 0, int(mobID))
		}
	}
}

func stepToward(fx, fy, tx, ty int32, stepDist float64) (int32, int32) {
	dx := float64(tx - fx)
	dy := float64(ty - fy)
	dist := math.Hypot(dx, dy)
	if dist <= stepDist || dist == 0 {
		return tx, ty
	}
	ratio := stepDist / dist
	return fx + int32(math.Round(dx*ratio)), fy + int32(math.Round(dy*ratio))
}

func (m *ZoneMap) moveMonsterStep(mob *ZoneMonster, nx, ny int32, speed uint16, running bool) {
	mob.mu.Lock()
	fromX, fromY := mob.X, mob.Y
	mob.X = nx
	mob.Y = ny
	mob.IsMoving = true
	mob.mu.Unlock()

	if m.grid != nil {
		m.grid.UpdateEntityPosition(mob, int(nx), int(ny))
	}

	op := protocol.PROTO_NC_ACT_SOMEONEMOVEWALK_CMD
	if running {
		op = protocol.PROTO_NC_ACT_SOMEONEMOVERUN_CMD
	}
	pw := protocol.NewPacketWriter(op)
	pw.WriteUint16(mob.Handle)
	pw.WriteInt32(fromX)
	pw.WriteInt32(fromY)
	pw.WriteInt32(nx)
	pw.WriteInt32(ny)
	pw.WriteUint16(speed)
	pw.WriteUint16(0)
	if m.grid != nil {
		m.grid.BroadcastAOI(pw, int(nx), int(ny), 0)
	}
}

// ApplyDamageToMonster aplica daño unificado al monstruo, registra aggro/contribución y dispara su muerte oficial.
func (m *ZoneMap) ApplyDamageToMonster(mob *ZoneMonster, attacker *ZonePlayer, dmg int) (restHP uint32, isDead bool) {
	if mob == nil || dmg <= 0 {
		return 0, false
	}
	var atkHandle uint16
	var charID int
	if attacker != nil {
		atkHandle = attacker.MapObjectID
		if attacker.Char != nil {
			charID = attacker.Char.ID
		}
	}

	mob.mu.Lock()
	if mob.IsDead {
		mob.mu.Unlock()
		return 0, false
	}
	if mob.CurHP <= dmg {
		mob.CurHP = 0
		mob.IsDead = true
		isDead = true
	} else {
		mob.CurHP -= dmg
	}
	restHP = uint32(mob.CurHP)
	mob.mu.Unlock()

	mob.RecordDamage(atkHandle, charID, dmg)

	if isDead {
		m.NotifyMonsterDeath(mob, attacker)
	} else {
		m.BroadcastMobChat(mob, "Damaged")
	}
	return restHP, isDead
}

// NotifyMonsterDeath delega la muerte del monstruo al pipeline completo de EXP, drops, anuncios y títulos.
func (m *ZoneMap) NotifyMonsterDeath(mob *ZoneMonster, killer *ZonePlayer) {
	if m == nil || mob == nil {
		return
	}
	if m.grid != nil {
		m.grid.RemoveEntity(mob.Handle)
	}
	m.BroadcastMobChat(mob, "Dead")

	m.mu.RLock()
	h := m.handler
	m.mu.RUnlock()

	// Fijar timer de respawn
	regenSec := mob.RegenIntervalSec
	if regenSec <= 0 {
		regenSec = 25
	}
	mob.mu.Lock()
	mob.AIState = MobAIDead
	mob.RespawnTime = time.Now().Add(time.Duration(regenSec) * time.Second)
	mob.IsMoving = false
	mob.mu.Unlock()

	if h != nil {
		h.HandleMonsterDeath(m, mob, killer)
	}
}
