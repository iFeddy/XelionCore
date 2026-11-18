package zone

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
)

func init() {
	registerGMCommand([]string{"killallmob"}, gmCommandSpec{
		Usage: "&killallmob",
		Exec: func(ctx *gmContext) error {
			if ctx.Player.CurrentMap == nil {
				return errors.New("[GM] No active map.")
			}
			killed := 0
			ctx.Player.CurrentMap.mu.Lock()
			for hnd, mob := range ctx.Player.CurrentMap.mobs {
				if mob != nil && !mob.IsDead {
					mob.CurHP = 0
					mob.IsDead = true
					killed++
					delete(ctx.Player.CurrentMap.mobs, hnd)
				}
			}
			ctx.Player.CurrentMap.mu.Unlock()
			ctx.Int1 = killed
			ctx.Reply = fmt.Sprintf("[GM] Cleared %d monsters from map", killed)
			return nil
		},
	})

	registerGMCommand([]string{"mobbreed"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&mobbreed <mobID> [count 1-20]",
		Exec: func(ctx *gmContext) error {
			if ctx.Player.CurrentMap == nil {
				return errors.New("[GM] No active map.")
			}
			mobID, err := strconv.Atoi(ctx.Args[0])
			if err != nil || mobID <= 0 {
				return errors.New("Invalid mobID.")
			}
			count := 1
			if len(ctx.Args) >= 2 {
				if c, err := strconv.Atoi(ctx.Args[1]); err == nil && c > 0 && c <= 20 {
					count = c
				}
			}
			ctx.TargetID = mobID
			ctx.Int1 = count
			var tmpl *MobTemplate
			if ctx.Handler.dataStore != nil && ctx.Handler.dataStore.MobCatalog != nil {
				tmpl = ctx.Handler.dataStore.MobCatalog.ByID[uint16(mobID)]
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
				_ = ctx.Player.CurrentMap.SpawnMobFromTemplate(tmpl, int32(ctx.Player.Char.XPos), int32(ctx.Player.Char.YPos), 0, "", 0, true)
			}
			ctx.Reply = fmt.Sprintf("[GM] Spawned mob %d x%d", mobID, count)
			return nil
		},
	})

	registerGMCommand([]string{"gmchat"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&gmchat <message>",
		Exec: func(ctx *gmContext) error {
			msg := strings.Join(ctx.Args, " ")
			ctx.Handler.DispatchWMS("wms_NC_CT_CHARTTING_CMD", map[string]interface{}{
				"sender":  ctx.Player.Char.Name,
				"message": msg,
			})
			ctx.Reply = fmt.Sprintf("[GMChat] %s", msg)
			return nil
		},
	})

	registerGMCommand([]string{"mfind"}, gmCommandSpec{
		Usage: "&mfind [mobID|name]",
		Exec: func(ctx *gmContext) error {
			ctx.Reply, ctx.Int1 = ctx.Handler.gmFindMob(ctx.Player, ctx.Args)
			return nil
		},
	})

	registerGMCommand([]string{"expbonus", "itembonus", "fieldbonus"}, gmCommandSpec{
		Usage: "&expbonus [rate%]",
		Exec: func(ctx *gmContext) error {
			ctx.Reply, ctx.Int1 = ctx.Handler.gmUpdateBonusRate(ctx.Cmd, ctx.Player, ctx.Args)
			return nil
		},
	})

	registerGMCommand([]string{"distribute", "mobdist"}, gmCommandSpec{
		Usage: "&distribute",
		Exec: func(ctx *gmContext) error {
			ctx.Reply, ctx.Int1 = ctx.Handler.gmMapDistribution(ctx.Player)
			return nil
		},
	})

	registerGMCommand([]string{"performance"}, gmCommandSpec{
		Usage: "&performance",
		Exec: func(ctx *gmContext) error {
			ctx.Reply, ctx.Int1 = ctx.Handler.gmPerformanceInfo(ctx.Player)
			return nil
		},
	})

	registerGMCommand([]string{"dist"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&dist <x> <y> | <charName>",
		Exec: func(ctx *gmContext) error {
			ctx.Reply, ctx.Int1 = ctx.Handler.gmDistance(ctx.Player, ctx.Args)
			return nil
		},
	})

	registerGMCommand([]string{"scrint"}, gmCommandSpec{
		Usage: "&scrint",
		Exec: func(ctx *gmContext) error {
			ctx.Reply, ctx.Int1 = ctx.Handler.gmScriptOrTimer(ctx.Cmd, ctx.Player, ctx.Args)
			return nil
		},
	})

	registerGMCommand([]string{"timer"}, gmCommandSpec{
		Usage: "&timer [seconds]",
		Exec: func(ctx *gmContext) error {
			ctx.Reply, ctx.Int1 = ctx.Handler.gmScriptOrTimer(ctx.Cmd, ctx.Player, ctx.Args)
			return nil
		},
	})

	registerGMCommand([]string{"oxblock"}, gmCommandSpec{
		Usage: "&oxblock",
		Exec: func(ctx *gmContext) error {
			globalGMState.mu.Lock()
			globalGMState.oxBlocked = !globalGMState.oxBlocked
			state := globalGMState.oxBlocked
			globalGMState.mu.Unlock()
			ctx.Int1 = boolToInt(state)
			ctx.Reply = fmt.Sprintf("[GM] OXBlock toggled=%v", state)
			return nil
		},
	})

	registerGMCommand([]string{"staticdamage"}, gmCommandSpec{
		Usage: "&staticdamage",
		Exec: func(ctx *gmContext) error {
			ctx.Int1, ctx.Reply = toggleGMPlayerFlag(ctx.Player.Char.ID, globalGMState.staticDamage, "StaticDamage")
			return nil
		},
	})

	registerGMCommand([]string{"deathaura"}, gmCommandSpec{
		Usage: "&deathaura",
		Exec: func(ctx *gmContext) error {
			ctx.Int1, ctx.Reply = toggleGMPlayerFlag(ctx.Player.Char.ID, globalGMState.deathAura, "DeathAura")
			if ctx.Int1 == 1 && ctx.Player.CurrentMap != nil {
				killed := 0
				ctx.Player.CurrentMap.mu.Lock()
				for hnd, m := range ctx.Player.CurrentMap.mobs {
					if m != nil && !m.IsDead {
						mx, my := m.GetPosition()
						dx := float64(mx) - float64(ctx.Player.Char.XPos)
						dy := float64(my) - float64(ctx.Player.Char.YPos)
						if math.Hypot(dx, dy) <= 500 {
							m.CurHP = 0
							m.IsDead = true
							delete(ctx.Player.CurrentMap.mobs, hnd)
							killed++
						}
					}
				}
				ctx.Player.CurrentMap.mu.Unlock()
				ctx.Int1 = killed
				ctx.Reply = fmt.Sprintf("[GM] DeathAura enabled (cleared %d nearby mobs)", killed)
			}
			return nil
		},
	})

	registerGMCommand([]string{"partyinfo"}, gmCommandSpec{
		Usage: "&partyinfo [charName]",
		Exec: func(ctx *gmContext) error {
			targetName := ctx.Player.Char.Name
			if len(ctx.Args) >= 1 {
				targetName = ctx.Args[0]
			}
			ctx.Handler.DispatchWMS("wms_NC_PARTY_INFORM_REQ", map[string]interface{}{
				"requester": ctx.Player.Char.Name,
				"target":    targetName,
			})
			ctx.Reply = fmt.Sprintf("[GM] Requested PartyInfo for '%s'", targetName)
			return nil
		},
	})

	registerGMCommand([]string{"propose", "proposeres", "weddingreset", "reserve", "enterready", "enter"}, gmCommandSpec{
		Usage: "&propose [args]",
		Exec: func(ctx *gmContext) error {
			ctx.Handler.DispatchWMS("wms_NC_WEDDING_GM_CMD", map[string]interface{}{
				"cmd":  ctx.Cmd,
				"char": ctx.Player.Char.Name,
				"args": strings.Join(ctx.Args, " "),
			})
			ctx.Reply = fmt.Sprintf("[GM] Dispatched wedding command '%s'", ctx.Cmd)
			return nil
		},
	})
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
