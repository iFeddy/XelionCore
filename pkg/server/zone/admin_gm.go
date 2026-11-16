package zone

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// gmRuntimeState mantiene multiplicadores y flags GM globales/por jugador verificados en AdminLvSet.shn.
type gmRuntimeState struct {
	mu           sync.RWMutex
	expBonus     int
	itemBonus    int
	fieldBonus   map[string]int
	oxBlocked    bool
	scriptPaused bool
	staticDamage map[int]bool
	deathAura    map[int]bool
	showID       map[int]bool
	horseSpeed   map[int]uint16
}

var globalGMState = &gmRuntimeState{
	expBonus:     100,
	itemBonus:    100,
	fieldBonus:   make(map[string]int),
	staticDamage: make(map[int]bool),
	deathAura:    make(map[int]bool),
	showID:       make(map[int]bool),
	horseSpeed:   make(map[int]uint16),
}

// AdminLvCatalog indexa los comandos GM oficiales cargados desde data/AdminLvSet.shn.
type AdminLvCatalog struct {
	Levels map[string]byte
}

var gmItemKeySeq atomic.Int64

func init() {
	gmItemKeySeq.Store(time.Now().UnixMicro())
}

// BuildAdminLvCatalog construye el catálogo oficial de permisos GM leyendo AdminLvSet.shn.
func BuildAdminLvCatalog(ds *DataStore) *AdminLvCatalog {
	cat := &AdminLvCatalog{
		Levels: make(map[string]byte),
	}
	if ds == nil {
		return cat
	}
	file := ds.GetSHN("AdminLvSet.shn")
	if file == nil {
		return cat
	}
	for _, row := range file.Rows {
		cmd := strings.ToLower(strings.TrimSpace(row.GetString("ALS_Cmd")))
		if cmd == "" {
			continue
		}
		lv := row.GetByte("ALS_Lv")
		cat.Levels[cmd] = lv
	}
	log.Printf("[Zone][GM] Indexed %d official GM commands from AdminLvSet.shn.", len(cat.Levels))
	return cat
}

// RequiredLevel consulta el nivel mínimo exigido por AdminLvSet.shn para un comando.
func (c *AdminLvCatalog) RequiredLevel(cmd string) (byte, bool) {
	if c == nil || c.Levels == nil {
		return 0, false
	}
	lv, ok := c.Levels[strings.ToLower(strings.TrimSpace(cmd))]
	return lv, ok
}

func (h *Handler) sendPlayerNotice(player *ZonePlayer, msg string) {
	if player == nil || player.Client == nil {
		return
	}
	msgBytes := []byte(msg)
	if len(msgBytes) > 250 {
		msgBytes = msgBytes[:250]
	}
	pw := protocol.NewPacketWriter(protocol.PROTO_NC_ACT_NOTICE_CMD)
	pw.WriteByte(0)
	pw.WriteByte(byte(len(msgBytes)))
	pw.WriteBytes(msgBytes)
	_ = player.Client.SendPacket(pw)
}

// handleGMCommand intercepta y despacha comandos administrativos gobernados 100% por AdminLvSet.shn + tCharacter.nAdminLevel.
func (h *Handler) handleGMCommand(player *ZonePlayer, rawText string) bool {
	if player == nil || player.Char == nil {
		return false
	}
	trimmed := strings.TrimSpace(rawText)
	if !strings.HasPrefix(trimmed, "&") && !strings.HasPrefix(trimmed, "/") {
		return false
	}
	body := strings.TrimSpace(trimmed[1:])
	if body == "" {
		return false
	}
	fields := strings.Fields(body)
	cmd := strings.ToLower(fields[0])
	args := fields[1:]

	if h.dataStore == nil || h.dataStore.AdminLvCatalog == nil {
		if strings.HasPrefix(trimmed, "&") {
			h.sendPlayerNotice(player, "GM catalog (AdminLvSet.shn) unavailable.")
			return true
		}
		return false
	}

	reqLv, exists := h.dataStore.AdminLvCatalog.RequiredLevel(cmd)
	if !exists {
		if strings.HasPrefix(trimmed, "&") {
			h.sendPlayerNotice(player, fmt.Sprintf("Unknown GM command: %s", cmd))
			return true
		}
		return false
	}

	// Validación estricta: tCharacter.nAdminLevel >= AdminLvSet.shn[cmd].ALS_Lv
	if player.Char.AdminLevel < reqLv {
		if player.Char.AdminLevel == 0 {
			return false
		}
		if h.gameLogDispatcher != nil {
			h.gameLogDispatcher.EnqueueTypeL(
				99,
				player.Char.ID,
				player.Char.ID,
				0,
				0,
				int(player.Char.AdminLevel),
				int(reqLv),
				0,
				0,
			)
		}
		h.sendPlayerNotice(player, fmt.Sprintf("Insufficient admin level for '%s' (have=%d, required=%d).", cmd, player.Char.AdminLevel, reqLv))
		return true
	}

	var (
		logType      int   = 1
		targetCharID int   = player.Char.ID
		targetID     int
		itemKey      int64
		int1         int
		int2         int
		int3         int = int(player.Char.AdminLevel)
		bigint1      int64
		replyMsg     = fmt.Sprintf("[GM] Executed '%s'", cmd)
	)

	switch cmd {
	case "makeitem":
		if len(args) < 1 {
			h.sendPlayerNotice(player, "Usage: &makeitem <itemID> [count]")
			return true
		}
		itemID, err := strconv.Atoi(args[0])
		if err != nil || itemID <= 0 {
			h.sendPlayerNotice(player, "Invalid itemID.")
			return true
		}
		count := 1
		if len(args) >= 2 {
			if c, err := strconv.Atoi(args[1]); err == nil && c > 0 {
				count = c
			}
		}
		targetID = itemID
		int1 = count
		int2 = int(reqLv)

		ip := "127.0.0.1"
		if player.Client != nil {
			ip = player.Client.RemoteIP()
		}
		if h.charRepo != nil {
			groupNo, grpErr := h.charRepo.CreateGMItemRewardGroup(player.Char.Name, ip, fmt.Sprintf("GM makeitem %d x%d", itemID, count))
			if grpErr != nil {
				h.sendPlayerNotice(player, fmt.Sprintf("GM reward group failed: %v", grpErr))
				return true
			}
			bigint1 = groupNo
			itemKey = gmItemKeySeq.Add(1)
			if _, grantErr := h.charRepo.GrantGMItemReward(groupNo, player.Char.ID, 9, itemID, itemKey); grantErr != nil {
				h.sendPlayerNotice(player, fmt.Sprintf("GM reward grant failed: %v", grantErr))
				return true
			}
			if count > 1 {
				_ = h.charRepo.UpdateItemFlags(itemKey, count)
			}
			if invItems, invErr := h.charRepo.GetItemsByStorage(player.Char.ID, 9); invErr == nil && len(invItems) > 0 {
				for _, it := range invItems {
					if it != nil {
						h.sendCellChangeFull(player.Client, 9, it.Storage, it, 0xFFFF)
					}
				}
			} else {
				h.sendCellChangeFull(player.Client, 9, 0, h.findItemAt(player.Char.ID, 9, 0), 0xFFFF)
			}
		}
		logType = 2
		replyMsg = fmt.Sprintf("[GM] Granted item %d x%d (group=%d, key=%d)", itemID, count, bigint1, itemKey)

	case "adminlevel":
		// Zone.exe ac_AdminLevel (0x0041A650): sin argumentos reporta "Admin level is %d".
		if len(args) == 0 {
			logType = 3
			int1 = int(player.Char.AdminLevel)
			int2 = int(reqLv)
			replyMsg = fmt.Sprintf("Admin level is %d", player.Char.AdminLevel)
			break
		}
		if len(args) < 2 {
			h.sendPlayerNotice(player, "Usage: &adminlevel [<charName> <level>]")
			return true
		}
		lvVal, err := strconv.Atoi(args[1])
		if err != nil || lvVal < 0 || lvVal > 255 {
			h.sendPlayerNotice(player, "Invalid admin level.")
			return true
		}
		targetName := args[0]
		newLv := byte(lvVal)
		if tp, ok := h.GetPlayerByName(targetName); ok && tp != nil && tp.Char != nil {
			targetCharID = tp.Char.ID
			tp.Char.AdminLevel = newLv
			if h.charRepo != nil {
				_, _ = h.charRepo.SetAdminLevel(tp.Char.ID, newLv)
			}
			if tp.Client != nil {
				admPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD)
				admPkt.WriteByte(newLv)
				_ = tp.Client.SendPacket(admPkt)
			}
		} else if h.charRepo != nil {
			if ch, err := h.charRepo.GetCharacterByName(targetName); err == nil && ch != nil {
				targetCharID = ch.ID
				_, _ = h.charRepo.SetAdminLevel(ch.ID, newLv)
			}
		}
		logType = 3
		int1 = int(newLv)
		int2 = int(reqLv)
		replyMsg = fmt.Sprintf("[GM] Set admin level for '%s' to %d", targetName, newLv)

	case "hide":
		player.IsGMHidden = true
		if player.CurrentMap != nil {
			del := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
			del.WriteUint16(player.MapObjectID)
			player.CurrentMap.Broadcast(del, player.MapObjectID)
		}
		replyMsg = "[GM] Hidden mode enabled."

	case "show":
		player.IsGMHidden = false
		replyMsg = "[GM] Visible mode restored."

	case "immortal":
		player.IsGMImmortal = !player.IsGMImmortal
		int1 = boolToInt(player.IsGMImmortal)
		replyMsg = fmt.Sprintf("[GM] Immortal mode: %v", player.IsGMImmortal)

	case "sethp":
		if len(args) >= 1 {
			if hpVal, err := strconv.Atoi(args[0]); err == nil && hpVal >= 0 {
				player.Char.CurHP = hpVal
				int1 = hpVal
				h.sendPlayerHPUpdate(player)
			}
		}

	case "setsp":
		if len(args) >= 1 {
			if spVal, err := strconv.Atoi(args[0]); err == nil && spVal >= 0 {
				player.Char.CurSP = spVal
				int1 = spVal
				h.sendPlayerSPUpdate(player)
			}
		}

	case "fullpoint":
		player.Char.CurHP = player.MaxHP
		player.Char.CurSP = player.MaxSP
		h.sendPlayerHPUpdate(player)
		h.sendPlayerSPUpdate(player)
		replyMsg = "[GM] Full HP/SP restored."

	case "linkto", "templink":
		if len(args) >= 1 {
			destMap := args[0]
			x := player.Char.XPos
			y := player.Char.YPos
			if len(args) >= 3 {
				if px, err := strconv.Atoi(args[1]); err == nil {
					x = px
				}
				if py, err := strconv.Atoi(args[2]); err == nil {
					y = py
				}
			} else if h.fieldStore != nil {
				if entry, ok := h.fieldStore.GetField(destMap); ok && entry != nil && len(entry.Respawn.Spots) > 0 {
					x = entry.Respawn.Spots[0].X
					y = entry.Respawn.Spots[0].Y
				}
			}
			int1 = x
			int2 = y
			h.WarpPlayer(player, destMap, x, y, player.Rotation)
			replyMsg = fmt.Sprintf("[GM] Warped to %s (%d, %d)", destMap, x, y)
		}

	case "z":
		if len(args) >= 2 && player.CurrentMap != nil {
			px, _ := strconv.Atoi(args[0])
			py, _ := strconv.Atoi(args[1])
			int1 = px
			int2 = py
			h.WarpPlayer(player, player.CurrentMap.MapID, px, py, player.Rotation)
		}

	case "psummon", "summonplayer":
		if len(args) >= 1 && player.CurrentMap != nil {
			if tp, ok := h.GetPlayerByName(args[0]); ok && tp != nil && tp.Char != nil {
				targetCharID = tp.Char.ID
				h.WarpPlayer(tp, player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, tp.Rotation)
				replyMsg = fmt.Sprintf("[GM] Summoned '%s'", tp.Char.Name)
			}
		}

	case "kickout":
		if len(args) >= 1 {
			if tp, ok := h.GetPlayerByName(args[0]); ok && tp != nil && tp.Char != nil {
				targetCharID = tp.Char.ID
				logType = 4
				if tp.Client != nil {
					tp.Client.Close()
				}
				replyMsg = fmt.Sprintf("[GM] Kicked out '%s'", tp.Char.Name)
			}
		}

	case "killallmob":
		if player.CurrentMap != nil {
			killed := 0
			player.CurrentMap.mu.Lock()
			for hnd, mob := range player.CurrentMap.mobs {
				if mob != nil && !mob.IsDead {
					mob.CurHP = 0
					mob.IsDead = true
					killed++
					delete(player.CurrentMap.mobs, hnd)
				}
			}
			player.CurrentMap.mu.Unlock()
			int1 = killed
			replyMsg = fmt.Sprintf("[GM] Cleared %d monsters from map", killed)
		}

	case "mobbreed":
		if len(args) >= 1 && player.CurrentMap != nil {
			mobID, _ := strconv.Atoi(args[0])
			count := 1
			if len(args) >= 2 {
				if c, err := strconv.Atoi(args[1]); err == nil && c > 0 && c <= 20 {
					count = c
				}
			}
			targetID = mobID
			int1 = count
			var tmpl *MobTemplate
			if h.dataStore != nil && h.dataStore.MobCatalog != nil {
				tmpl = h.dataStore.MobCatalog.ByID[uint16(mobID)]
			}
			if tmpl == nil {
				tmpl = &MobTemplate{
					ID:    uint16(mobID),
					Name:  fmt.Sprintf("Mob%d", mobID),
					Level: 1,
					MaxHP: 100,
				}
			}
			for i := 0; i < count; i++ {
				_ = player.CurrentMap.SpawnMobFromTemplate(tmpl, int32(player.Char.XPos), int32(player.Char.YPos), 0, "", 0, true)
			}
			replyMsg = fmt.Sprintf("[GM] Spawned mob %d x%d", mobID, count)
		}

	case "coord", "pfind", "charinform":
		if player.CurrentMap != nil {
			replyMsg = fmt.Sprintf("[GM] Map=%s X=%d Y=%d AdminLv=%d", player.CurrentMap.MapID, player.Char.XPos, player.Char.YPos, player.Char.AdminLevel)
		}
		if len(args) >= 1 {
			if tp, ok := h.GetPlayerByName(args[0]); ok && tp != nil && tp.Char != nil && tp.CurrentMap != nil {
				targetCharID = tp.Char.ID
				replyMsg = fmt.Sprintf("[GM] Target '%s' Map=%s X=%d Y=%d Lv=%d", tp.Char.Name, tp.CurrentMap.MapID, tp.Char.XPos, tp.Char.YPos, tp.Char.Level)
			}
		}

	case "gmchat":
		if len(args) > 0 {
			h.DispatchWMS("wms_NC_CT_CHARTTING_CMD", map[string]interface{}{
				"sender":  player.Char.Name,
				"message": strings.Join(args, " "),
			})
		}

	case "mfind":
		replyMsg, int1 = h.gmFindMob(player, args)

	case "expbonus", "itembonus", "fieldbonus":
		replyMsg, int1 = h.gmUpdateBonusRate(cmd, player, args)

	case "distribute", "mobdist":
		replyMsg, int1 = h.gmMapDistribution(player)

	case "performance":
		replyMsg, int1 = h.gmPerformanceInfo(player)

	case "dist":
		replyMsg, int1 = h.gmDistance(player, args)

	case "scrint", "timer":
		replyMsg, int1 = h.gmScriptOrTimer(cmd, player, args)

	case "oxblock":
		globalGMState.mu.Lock()
		globalGMState.oxBlocked = !globalGMState.oxBlocked
		state := globalGMState.oxBlocked
		globalGMState.mu.Unlock()
		int1 = boolToInt(state)
		replyMsg = fmt.Sprintf("[GM] OXBlock toggled=%v", state)

	case "staticdamage":
		int1, replyMsg = toggleGMPlayerFlag(player.Char.ID, globalGMState.staticDamage, "StaticDamage")

	case "deathaura":
		int1, replyMsg = toggleGMPlayerFlag(player.Char.ID, globalGMState.deathAura, "DeathAura")
		if int1 == 1 && player.CurrentMap != nil {
			killed := 0
			player.CurrentMap.mu.Lock()
			for hnd, m := range player.CurrentMap.mobs {
				if m != nil && !m.IsDead {
					mx, my := m.GetPosition()
					dx := float64(mx) - float64(player.Char.XPos)
					dy := float64(my) - float64(player.Char.YPos)
					if math.Hypot(dx, dy) <= 500 {
						m.CurHP = 0
						m.IsDead = true
						delete(player.CurrentMap.mobs, hnd)
						killed++
					}
				}
			}
			player.CurrentMap.mu.Unlock()
			int1 = killed
			replyMsg = fmt.Sprintf("[GM] DeathAura enabled (cleared %d nearby mobs)", killed)
		}

	case "showid":
		int1, replyMsg = toggleGMPlayerFlag(player.Char.ID, globalGMState.showID, "ShowID")
		if int1 == 1 {
			replyMsg = fmt.Sprintf("[GM] ShowID enabled (charID=%d handle=%d acc=%d)", player.Char.ID, player.MapObjectID, player.Char.AccountID)
		}

	case "horsespeed":
		spd := uint16(0)
		if len(args) >= 1 {
			if v, err := strconv.Atoi(args[0]); err == nil && v >= 0 && v <= 2000 {
				spd = uint16(v)
			}
		}
		globalGMState.mu.Lock()
		if spd == 0 {
			delete(globalGMState.horseSpeed, player.Char.ID)
		} else {
			globalGMState.horseSpeed[player.Char.ID] = spd
		}
		globalGMState.mu.Unlock()
		int1 = int(spd)
		replyMsg = fmt.Sprintf("[GM] HorseSpeed set to %d", spd)

	case "partyinfo":
		targetName := player.Char.Name
		if len(args) >= 1 {
			targetName = args[0]
		}
		h.DispatchWMS("wms_NC_PARTY_INFORM_REQ", map[string]interface{}{
			"requester": player.Char.Name,
			"target":    targetName,
		})
		replyMsg = fmt.Sprintf("[GM] Requested PartyInfo for '%s'", targetName)

	case "propose", "proposeres", "weddingreset", "reserve", "enterready", "enter":
		h.DispatchWMS("wms_NC_WEDDING_GM_CMD", map[string]interface{}{
			"cmd":  cmd,
			"char": player.Char.Name,
			"args": strings.Join(args, " "),
		})
		replyMsg = fmt.Sprintf("[GM] Dispatched wedding command '%s'", cmd)

	default:
		// Comandos adicionales indexados en AdminLvSet.shn (incluidos alias localizados).
		int2 = int(reqLv)
	}

	if h.gameLogDispatcher != nil {
		h.gameLogDispatcher.EnqueueTypeL(
			logType,
			player.Char.ID,
			targetCharID,
			targetID,
			itemKey,
			int1,
			int2,
			int3,
			bigint1,
		)
	}

	h.sendPlayerNotice(player, replyMsg)
	return true
}

func toggleGMPlayerFlag(charID int, flagMap map[int]bool, label string) (int, string) {
	globalGMState.mu.Lock()
	next := !flagMap[charID]
	flagMap[charID] = next
	globalGMState.mu.Unlock()
	return boolToInt(next), fmt.Sprintf("[GM] %s toggled=%v", label, next)
}

func (h *Handler) gmFindMob(player *ZonePlayer, args []string) (string, int) {
	if player.CurrentMap == nil {
		return "[GM] No active map.", 0
	}
	query := ""
	if len(args) > 0 {
		query = strings.ToLower(strings.TrimSpace(args[0]))
	}
	mobs := player.CurrentMap.GetAllMonsters()
	matches := make([]string, 0, 5)
	count := 0
	for _, m := range mobs {
		if m == nil || m.CurHP <= 0 {
			continue
		}
		idStr := strconv.Itoa(int(m.MobID))
		if query != "" && idStr != query && !strings.Contains(strings.ToLower(m.Name), query) && !strings.Contains(strings.ToLower(m.InxName), query) {
			continue
		}
		count++
		if len(matches) < 4 {
			mx, my := m.GetPosition()
			matches = append(matches, fmt.Sprintf("%s#%d@(%d,%d)", m.InxName, m.Handle, mx, my))
		}
	}
	if count == 0 {
		return fmt.Sprintf("[GM] mfind '%s': 0 mobs found on %s", query, player.CurrentMap.MapID), 0
	}
	return fmt.Sprintf("[GM] mfind (%d total): %s", count, strings.Join(matches, ", ")), count
}

func (h *Handler) gmUpdateBonusRate(cmd string, player *ZonePlayer, args []string) (string, int) {
	globalGMState.mu.Lock()
	defer globalGMState.mu.Unlock()
	mapID := ""
	if player.CurrentMap != nil {
		mapID = player.CurrentMap.MapID
	}
	if len(args) == 0 {
		switch cmd {
		case "expbonus":
			return fmt.Sprintf("[GM] ExpBonus=%d%%", globalGMState.expBonus), globalGMState.expBonus
		case "itembonus":
			return fmt.Sprintf("[GM] ItemBonus=%d%%", globalGMState.itemBonus), globalGMState.itemBonus
		default:
			cur := globalGMState.fieldBonus[mapID]
			return fmt.Sprintf("[GM] FieldBonus(%s)=%d%%", mapID, cur), cur
		}
	}
	rate, err := strconv.Atoi(args[0])
	if err != nil || rate < 0 || rate > 10000 {
		return fmt.Sprintf("[GM] Invalid bonus rate for %s", cmd), 0
	}
	switch cmd {
	case "expbonus":
		globalGMState.expBonus = rate
	case "itembonus":
		globalGMState.itemBonus = rate
	default:
		globalGMState.fieldBonus[mapID] = rate
	}
	return fmt.Sprintf("[GM] %s set to %d%%", cmd, rate), rate
}

func (h *Handler) gmMapDistribution(player *ZonePlayer) (string, int) {
	if player.CurrentMap == nil {
		return "[GM] No active map.", 0
	}
	mobs := player.CurrentMap.GetAllMonsters()
	alive := 0
	for _, m := range mobs {
		if m != nil && m.CurHP > 0 {
			alive++
		}
	}
	players := len(player.CurrentMap.GetPlayers())
	return fmt.Sprintf("[GM] Dist[%s]: players=%d mobs_alive=%d mobs_total=%d", player.CurrentMap.MapID, players, alive, len(mobs)), alive
}

func (h *Handler) gmPerformanceInfo(player *ZonePlayer) (string, int) {
	players, mobs := 0, 0
	mapID := "none"
	if player.CurrentMap != nil {
		mapID = player.CurrentMap.MapID
		players = len(player.CurrentMap.GetPlayers())
		mobs = len(player.CurrentMap.GetAllMonsters())
	}
	g := runtime.NumGoroutine()
	return fmt.Sprintf("[GM] Perf[%s]: goroutines=%d players=%d mobs=%d", mapID, g, players, mobs), g
}

func (h *Handler) gmDistance(player *ZonePlayer, args []string) (string, int) {
	if len(args) >= 2 {
		tx, ex := strconv.ParseFloat(args[0], 64)
		ty, ey := strconv.ParseFloat(args[1], 64)
		if ex == nil && ey == nil {
			d := math.Hypot(tx-float64(player.Char.XPos), ty-float64(player.Char.YPos))
			return fmt.Sprintf("[GM] Dist to (%.0f,%.0f) = %.2f", tx, ty, d), int(d)
		}
	}
	if len(args) == 1 {
		if tp, ok := h.GetPlayerByName(args[0]); ok && tp != nil && tp.Char != nil {
			d := math.Hypot(float64(tp.Char.XPos)-float64(player.Char.XPos), float64(tp.Char.YPos)-float64(player.Char.YPos))
			return fmt.Sprintf("[GM] Dist to %s = %.2f", tp.Char.Name, d), int(d)
		}
	}
	return "Usage: &dist <x> <y> | <charName>", 0
}

func (h *Handler) gmScriptOrTimer(cmd string, player *ZonePlayer, args []string) (string, int) {
	if cmd == "scrint" {
		globalGMState.mu.Lock()
		globalGMState.scriptPaused = !globalGMState.scriptPaused
		st := globalGMState.scriptPaused
		globalGMState.mu.Unlock()
		return fmt.Sprintf("[GM] ScriptInterrupt paused=%v", st), boolToInt(st)
	}
	secs := 0
	if len(args) >= 1 {
		if v, err := strconv.Atoi(args[0]); err == nil && v >= 0 {
			secs = v
		}
	}
	return fmt.Sprintf("[GM] Timer configured=%ds", secs), secs
}

func (h *Handler) sendPlayerHPUpdate(player *ZonePlayer) {
	if player == nil || player.Client == nil || player.Char == nil {
		return
	}
	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(uint32(player.Char.CurHP))
	hpCmd.WriteUint16(0)
	_ = player.Client.SendPacket(hpCmd)
}

func (h *Handler) sendPlayerSPUpdate(player *ZonePlayer) {
	if player == nil || player.Client == nil || player.Char == nil {
		return
	}
	spCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SPCHANGE_CMD)
	spCmd.WriteUint32(uint32(player.Char.CurSP))
	_ = player.Client.SendPacket(spCmd)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
