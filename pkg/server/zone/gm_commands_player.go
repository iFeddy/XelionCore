package zone

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func init() {
	registerGMCommand([]string{"addexp"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&addexp <amount>",
		LogType: 6,
		Exec: func(ctx *gmContext) error {
			expGain, err := strconv.ParseUint(ctx.Args[0], 10, 64)
			if err != nil || expGain == 0 {
				return errors.New("Invalid EXP amount.")
			}

			mapID := ""
			if ctx.Player.CurrentMap != nil {
				mapID = ctx.Player.CurrentMap.MapID
			}

			beforeLv := ctx.Player.Char.Level
			ctx.Handler.AddPlayerExpAndCheckLevelUp(
				ctx.Player,
				expGain,
				0,
				0,
				int32(ctx.Player.Char.XPos),
				int32(ctx.Player.Char.YPos),
				mapID,
			)

			ctx.BigInt1 = int64(expGain)
			ctx.Int1 = int(ctx.Player.Char.Level)
			ctx.Reply = fmt.Sprintf("[GM] Added %d EXP (Lv.%d -> Lv.%d, Total EXP=%d)",
				expGain, beforeLv, ctx.Player.Char.Level, ctx.Player.Char.TCharacterState.Exp)
			return nil
		},
	})

	registerGMCommand([]string{"setlevel"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&setlevel <level>",
		LogType: 7,
		Exec: func(ctx *gmContext) error {
			lvVal, err := strconv.Atoi(ctx.Args[0])
			if err != nil || lvVal < 1 {
				return errors.New("Invalid level.")
			}

			var cs *combat.Store
			if ctx.Handler.dataStore != nil {
				cs = ctx.Handler.dataStore.CombatStore
			}
			maxLv := 125
			if cs != nil {
				if m := cs.MaxLevel(); m > 0 {
					maxLv = m
				}
				if c := cs.Common(); c.LevelLimit > 0 && c.LevelLimit < maxLv {
					maxLv = c.LevelLimit
				}
			}
			if lvVal > maxLv {
				return fmt.Errorf("Level cannot exceed max level (%d).", maxLv)
			}

			targetLv := byte(lvVal)
			curLv := ctx.Player.Char.Level
			if targetLv == curLv {
				ctx.Reply = fmt.Sprintf("[GM] Already at Lv.%d.", curLv)
				return nil
			}

			ctx.Player.Char.Level = targetLv
			ctx.Player.Char.TCharacter.Level = targetLv
			if cs != nil {
				prevExp, newNext := cs.GetExpProgression(int(targetLv))
				ctx.Player.Char.TCharacterState.Exp = int64(prevExp)
				ctx.Player.PrevExp = uint64(prevExp)
				ctx.Player.NextExp = uint64(newNext)
			}

			if ctx.Handler.classStore != nil {
				if lp, ok := ctx.Handler.classStore.GetParam(JobToClassName(ctx.Player.Char.Job), int(targetLv)); ok && lp != nil {
					ctx.Player.LevelParam = *lp
					if cs != nil {
						free := cs.GetFreeStatBonus(
							ctx.Player.Char.TCharacter.Strength,
							ctx.Player.Char.TCharacter.Intelligence,
							ctx.Player.Char.TCharacter.Dexterity,
							ctx.Player.Char.TCharacter.Constitute,
							ctx.Player.Char.TCharacter.MentalPower,
						)
						if lp.MaxHP+free.MaxHP > 0 {
							ctx.Player.MaxHP = lp.MaxHP + free.MaxHP
						}
						if lp.MaxSP+free.MaxSP > 0 {
							ctx.Player.MaxSP = lp.MaxSP + free.MaxSP
						}
					}
				}
			}
			ctx.Player.RecalculateEquipmentStats(ctx.Handler.itemStore)
			ctx.Player.Char.CurHP = ctx.Player.MaxHP
			ctx.Player.Char.CurSP = ctx.Player.MaxSP

			if ctx.Player.Client != nil {
				lvCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_LEVELUP_CMD)
				lvCmd.WriteByte(ctx.Player.Char.Level)
				lvCmd.WriteUint16(0)
				ctx.Player.WriteDetailedInfoExtra(lvCmd, true)
				_ = ctx.Player.Client.SendPacket(lvCmd)

				ctx.Handler.syncPlayerVitals(ctx.Player)
			}

			ctx.Int1 = int(targetLv)
			ctx.BigInt1 = ctx.Player.Char.TCharacterState.Exp
			ctx.Reply = fmt.Sprintf("[GM] Level set: Lv.%d -> Lv.%d", curLv, targetLv)
			return nil
		},
	})

	registerGMCommand([]string{"money"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&money <amount[G|g|s|c]...> | set <amount[G|g|s|c]...>",
		LogType: 5,
		Exec: func(ctx *gmContext) error {
			setExact := false
			rawArgs := ctx.Args
			if strings.EqualFold(ctx.Args[0], "set") {
				if len(ctx.Args) < 2 {
					return errors.New("Usage: &money set <amount[G|g|s|c]...>")
				}
				setExact = true
				rawArgs = ctx.Args[1:]
			}

			delta, err := parseFiestaMoney(rawArgs)
			if err != nil {
				return err
			}

			cur := ctx.Player.Char.TCharacterState.Money
			var next int64
			if setExact {
				if delta < 0 {
					return errors.New("Money amount cannot be negative.")
				}
				next = delta
			} else {
				next = cur + delta
				if next < 0 {
					next = 0
				}
			}

			ctx.Player.Char.TCharacterState.Money = next
			ctx.Handler.syncPlayerMoney(ctx.Player)
			ctx.BigInt1 = next

			if setExact {
				ctx.Reply = fmt.Sprintf("[GM] Money set to %s (%d c)", formatFiestaMoney(next), next)
			} else {
				sign := "+"
				if delta < 0 {
					sign = ""
				}
				ctx.Reply = fmt.Sprintf("[GM] Money %s%s -> Total: %s", sign, formatFiestaMoney(delta), formatFiestaMoney(next))
			}
			return nil
		},
	})

	registerGMCommand([]string{"makeitem"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&makeitem <itemID|InxName> [count]",
		LogType: 2,
		Exec: func(ctx *gmContext) error {
			rawItem := strings.TrimSpace(ctx.Args[0])
			var itemID int
			var itemInx string

			if id, err := strconv.Atoi(rawItem); err == nil && id > 0 {
				itemID = id
				if ctx.Handler.itemStore != nil {
					if info, ok := ctx.Handler.itemStore.GetItemInfo(uint16(itemID)); ok && info != nil {
						itemInx = info.InxName
					}
				}
			} else {
				if ctx.Handler.itemStore == nil {
					return errors.New("ItemStore unavailable.")
				}
				info, ok := ctx.Handler.itemStore.GetItemInfoByName(rawItem)
				if !ok || info == nil || info.ID == 0 {
					return fmt.Errorf("Item '%s' not found.", rawItem)
				}
				itemID = int(info.ID)
				itemInx = info.InxName
			}

			count := 1
			if len(ctx.Args) >= 2 {
				if c, err := strconv.Atoi(ctx.Args[1]); err == nil && c > 0 {
					count = c
				}
			}
			ctx.TargetID = itemID
			ctx.Int1 = count

			ip := "127.0.0.1"
			if ctx.Player.Client != nil {
				ip = ctx.Player.Client.RemoteIP()
			}
			if ctx.Handler.charRepo != nil {
				groupNo, grpErr := ctx.Handler.charRepo.CreateGMItemRewardGroup(
					ctx.Player.Char.Name,
					ip,
					fmt.Sprintf("GM makeitem %d x%d", itemID, count),
				)
				if grpErr != nil {
					return fmt.Errorf("GM reward group failed: %v", grpErr)
				}
				ctx.BigInt1 = groupNo
				ctx.ItemKey = gmItemKeySeq.Add(1)
				if _, grantErr := ctx.Handler.charRepo.GrantGMItemReward(groupNo, ctx.Player.Char.ID, 9, itemID, ctx.ItemKey); grantErr != nil {
					return fmt.Errorf("GM reward grant failed: %v", grantErr)
				}
				if count > 1 {
					_ = ctx.Handler.charRepo.UpdateItemFlags(ctx.ItemKey, count)
				}
				if invItems, invErr := ctx.Handler.charRepo.GetItemsByStorage(ctx.Player.Char.ID, 9); invErr == nil && len(invItems) > 0 {
					for _, it := range invItems {
						if it != nil {
							ctx.Handler.sendCellChangeFull(ctx.Player.Client, 9, it.Storage, it, 0xFFFF)
						}
					}
				} else {
					ctx.Handler.sendCellChangeFull(ctx.Player.Client, 9, 0, ctx.Handler.findItemAt(ctx.Player.Char.ID, 9, 0), 0xFFFF)
				}
			}
			if itemInx != "" {
				ctx.Reply = fmt.Sprintf("[GM] Granted item %s (%d) x%d (group=%d, key=%d)", itemInx, itemID, count, ctx.BigInt1, ctx.ItemKey)
			} else {
				ctx.Reply = fmt.Sprintf("[GM] Granted item %d x%d (group=%d, key=%d)", itemID, count, ctx.BigInt1, ctx.ItemKey)
			}
			return nil
		},
	})

	registerGMCommand([]string{"adminlevel"}, gmCommandSpec{
		Usage:   "&adminlevel [<charName> <level>]",
		LogType: 3,
		Exec: func(ctx *gmContext) error {
			// Zone.exe ac_AdminLevel (0x0041A650): sin argumentos reporta "Admin level is %d".
			if len(ctx.Args) == 0 {
				ctx.Int1 = int(ctx.Player.Char.AdminLevel)
				ctx.Reply = fmt.Sprintf("Admin level is %d", ctx.Player.Char.AdminLevel)
				return nil
			}
			if len(ctx.Args) < 2 {
				return errors.New("Usage: &adminlevel [<charName> <level>]")
			}
			lvVal, err := strconv.Atoi(ctx.Args[1])
			if err != nil || lvVal < 0 || lvVal > 255 {
				return errors.New("Invalid admin level.")
			}
			targetName := ctx.Args[0]
			newLv := byte(lvVal)
			if tp, ok := ctx.Handler.GetPlayerByName(targetName); ok && tp != nil && tp.Char != nil {
				ctx.TargetCharID = tp.Char.ID
				tp.Char.AdminLevel = newLv
				if ctx.Handler.charRepo != nil {
					_, _ = ctx.Handler.charRepo.SetAdminLevel(tp.Char.ID, newLv)
				}
				if tp.Client != nil {
					admPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD)
					admPkt.WriteByte(newLv)
					_ = tp.Client.SendPacket(admPkt)
				}
			} else if ctx.Handler.charRepo != nil {
				if ch, err := ctx.Handler.charRepo.GetCharacterByName(targetName); err == nil && ch != nil {
					ctx.TargetCharID = ch.ID
					_, _ = ctx.Handler.charRepo.SetAdminLevel(ch.ID, newLv)
				}
			}
			ctx.Int1 = int(newLv)
			ctx.Reply = fmt.Sprintf("[GM] Set admin level for '%s' to %d", targetName, newLv)
			return nil
		},
	})

	registerGMCommand([]string{"hide"}, gmCommandSpec{
		Usage: "&hide",
		Exec: func(ctx *gmContext) error {
			ctx.Player.IsGMHidden = true
			if ctx.Player.CurrentMap != nil {
				del := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD)
				del.WriteUint16(ctx.Player.MapObjectID)
				ctx.Player.CurrentMap.Broadcast(del, ctx.Player.MapObjectID)
			}
			ctx.Reply = "[GM] Hidden mode enabled."
			return nil
		},
	})

	registerGMCommand([]string{"show"}, gmCommandSpec{
		Usage: "&show",
		Exec: func(ctx *gmContext) error {
			ctx.Player.IsGMHidden = false
			ctx.Reply = "[GM] Visible mode restored."
			return nil
		},
	})

	registerGMCommand([]string{"immortal"}, gmCommandSpec{
		Usage: "&immortal",
		Exec: func(ctx *gmContext) error {
			ctx.Player.IsGMImmortal = !ctx.Player.IsGMImmortal
			ctx.Int1 = boolToInt(ctx.Player.IsGMImmortal)
			ctx.Reply = fmt.Sprintf("[GM] Immortal mode: %v", ctx.Player.IsGMImmortal)
			return nil
		},
	})

	registerGMCommand([]string{"sethp"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&sethp <amount>",
		Exec: func(ctx *gmContext) error {
			hpVal, err := strconv.Atoi(ctx.Args[0])
			if err != nil || hpVal < 0 {
				return errors.New("Invalid HP amount.")
			}
			ctx.Player.Char.CurHP = hpVal
			ctx.Int1 = hpVal
			ctx.Handler.sendPlayerHPUpdate(ctx.Player)
			ctx.Reply = fmt.Sprintf("[GM] HP set to %d", hpVal)
			return nil
		},
	})

	registerGMCommand([]string{"setsp"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&setsp <amount>",
		Exec: func(ctx *gmContext) error {
			spVal, err := strconv.Atoi(ctx.Args[0])
			if err != nil || spVal < 0 {
				return errors.New("Invalid SP amount.")
			}
			ctx.Player.Char.CurSP = spVal
			ctx.Int1 = spVal
			ctx.Handler.sendPlayerSPUpdate(ctx.Player)
			ctx.Reply = fmt.Sprintf("[GM] SP set to %d", spVal)
			return nil
		},
	})

	registerGMCommand([]string{"fullpoint"}, gmCommandSpec{
		Usage: "&fullpoint",
		Exec: func(ctx *gmContext) error {
			ctx.Player.Char.CurHP = ctx.Player.MaxHP
			ctx.Player.Char.CurSP = ctx.Player.MaxSP
			ctx.Handler.sendPlayerHPUpdate(ctx.Player)
			ctx.Handler.sendPlayerSPUpdate(ctx.Player)
			ctx.Handler.syncPlayerVitals(ctx.Player)
			ctx.Reply = "[GM] Full HP/SP restored."
			return nil
		},
	})

	registerGMCommand([]string{"linkto", "templink"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&linkto <mapID> [x] [y]",
		Exec: func(ctx *gmContext) error {
			destMap := ctx.Args[0]
			x := ctx.Player.Char.XPos
			y := ctx.Player.Char.YPos
			if len(ctx.Args) >= 3 {
				if px, err := strconv.Atoi(ctx.Args[1]); err == nil {
					x = px
				}
				if py, err := strconv.Atoi(ctx.Args[2]); err == nil {
					y = py
				}
			} else if ctx.Handler.fieldStore != nil {
				if entry, ok := ctx.Handler.fieldStore.GetField(destMap); ok && entry != nil && len(entry.Respawn.Spots) > 0 {
					x = entry.Respawn.Spots[0].X
					y = entry.Respawn.Spots[0].Y
				}
			}
			ctx.Int1 = x
			ctx.Int2 = y
			ctx.Handler.WarpPlayer(ctx.Player, destMap, x, y, ctx.Player.Rotation)
			ctx.Reply = fmt.Sprintf("[GM] Warped to %s (%d, %d)", destMap, x, y)
			return nil
		},
	})

	registerGMCommand([]string{"z"}, gmCommandSpec{
		MinArgs: 2,
		Usage:   "&z <x> <y>",
		Exec: func(ctx *gmContext) error {
			if ctx.Player.CurrentMap == nil {
				return errors.New("[GM] No active map.")
			}
			px, errX := strconv.Atoi(ctx.Args[0])
			py, errY := strconv.Atoi(ctx.Args[1])
			if errX != nil || errY != nil {
				return errors.New("Invalid coordinates.")
			}
			ctx.Int1 = px
			ctx.Int2 = py
			ctx.Handler.WarpPlayer(ctx.Player, ctx.Player.CurrentMap.MapID, px, py, ctx.Player.Rotation)
			ctx.Reply = fmt.Sprintf("[GM] Teleported to (%d, %d)", px, py)
			return nil
		},
	})

	registerGMCommand([]string{"psummon", "summonplayer"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&psummon <charName>",
		Exec: func(ctx *gmContext) error {
			if ctx.Player.CurrentMap == nil {
				return errors.New("[GM] No active map.")
			}
			tp, ok := ctx.Handler.GetPlayerByName(ctx.Args[0])
			if !ok || tp == nil || tp.Char == nil {
				return fmt.Errorf("[GM] Player '%s' not found online.", ctx.Args[0])
			}
			ctx.TargetCharID = tp.Char.ID
			ctx.Handler.WarpPlayer(tp, ctx.Player.CurrentMap.MapID, ctx.Player.Char.XPos, ctx.Player.Char.YPos, tp.Rotation)
			ctx.Reply = fmt.Sprintf("[GM] Summoned '%s'", tp.Char.Name)
			return nil
		},
	})

	registerGMCommand([]string{"kickout"}, gmCommandSpec{
		MinArgs: 1,
		Usage:   "&kickout <charName>",
		LogType: 4,
		Exec: func(ctx *gmContext) error {
			tp, ok := ctx.Handler.GetPlayerByName(ctx.Args[0])
			if !ok || tp == nil || tp.Char == nil {
				return fmt.Errorf("[GM] Player '%s' not found online.", ctx.Args[0])
			}
			ctx.TargetCharID = tp.Char.ID
			if tp.Client != nil {
				tp.Client.Close()
			}
			ctx.Reply = fmt.Sprintf("[GM] Kicked out '%s'", tp.Char.Name)
			return nil
		},
	})

	registerGMCommand([]string{"coord", "pfind", "charinform"}, gmCommandSpec{
		Usage: "&coord [charName]",
		Exec: func(ctx *gmContext) error {
			if len(ctx.Args) >= 1 {
				tp, ok := ctx.Handler.GetPlayerByName(ctx.Args[0])
				if !ok || tp == nil || tp.Char == nil || tp.CurrentMap == nil {
					return fmt.Errorf("[GM] Player '%s' not found online.", ctx.Args[0])
				}
				ctx.TargetCharID = tp.Char.ID
				ctx.Reply = fmt.Sprintf("[GM] Target '%s' Map=%s X=%d Y=%d Lv=%d", tp.Char.Name, tp.CurrentMap.MapID, tp.Char.XPos, tp.Char.YPos, tp.Char.Level)
				return nil
			}
			if ctx.Player.CurrentMap == nil {
				return errors.New("[GM] No active map.")
			}
			ctx.Reply = fmt.Sprintf("[GM] Map=%s X=%d Y=%d AdminLv=%d", ctx.Player.CurrentMap.MapID, ctx.Player.Char.XPos, ctx.Player.Char.YPos, ctx.Player.Char.AdminLevel)
			return nil
		},
	})

	registerGMCommand([]string{"showid"}, gmCommandSpec{
		Usage: "&showid",
		Exec: func(ctx *gmContext) error {
			ctx.Int1, ctx.Reply = toggleGMPlayerFlag(ctx.Player.Char.ID, globalGMState.showID, "ShowID")
			if ctx.Int1 == 1 {
				ctx.Reply = fmt.Sprintf("[GM] ShowID enabled (charID=%d handle=%d acc=%d)", ctx.Player.Char.ID, ctx.Player.MapObjectID, ctx.Player.Char.AccountID)
			}
			return nil
		},
	})

	registerGMCommand([]string{"horsespeed"}, gmCommandSpec{
		Usage: "&horsespeed [0-2000]",
		Exec: func(ctx *gmContext) error {
			spd := uint16(0)
			if len(ctx.Args) >= 1 {
				v, err := strconv.Atoi(ctx.Args[0])
				if err != nil || v < 0 || v > 2000 {
					return errors.New("Invalid horse speed (0-2000).")
				}
				spd = uint16(v)
			}
			globalGMState.mu.Lock()
			if spd == 0 {
				delete(globalGMState.horseSpeed, ctx.Player.Char.ID)
			} else {
				globalGMState.horseSpeed[ctx.Player.Char.ID] = spd
			}
			globalGMState.mu.Unlock()
			ctx.Int1 = int(spd)
			ctx.Reply = fmt.Sprintf("[GM] HorseSpeed set to %d", spd)
			return nil
		},
	})
}

func (h *Handler) sendPlayerHPUpdate(player *ZonePlayer) {
	if player == nil || player.Client == nil || player.Char == nil {
		return
	}
	hpCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_HPCHANGE_CMD)
	hpCmd.WriteUint32(uint32(player.Char.CurHP))
	hpCmd.WriteUint16(player.NextHPChangeOrder())
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

// syncPlayerVitals actualiza primero los techos MaxHP/MaxSP (índices 32 y 33 de CHAR_PARAMETER_DATA)
// y luego rellena las barras actuales de HP y SP con orden incremental.
func (h *Handler) syncPlayerVitals(player *ZonePlayer) {
	if player == nil || player.Client == nil || player.Char == nil {
		return
	}

	pkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_BASEPARAMCHANGE_CMD)
	pkt.WriteByte(2)
	pkt.WriteByte(0x10) // MaxHP
	pkt.WriteUint32(uint32(player.MaxHP))
	pkt.WriteByte(0x11) // MaxSP
	pkt.WriteUint32(uint32(player.MaxSP))
	_ = player.Client.SendPacket(pkt)

	h.sendPlayerHPUpdate(player)
	h.sendPlayerSPUpdate(player)
}

const (
	fiestaCopperPerSilver int64 = 1_000
	fiestaCopperPerGold   int64 = 1_000_000
	fiestaCopperPerGem    int64 = 100_000_000
)

// parseFiestaMoney convierte argumentos como ["2G", "50g", "234s566c"] o ["234566"] a Copper total.
// Respeta case-sensitive: 'G' = Gem (100.000.000), 'g' = Gold (1.000.000), 's' = Silver (1.000), 'c' = Copper (1).
func parseFiestaMoney(tokens []string) (int64, error) {
	var total int64
	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		sign := int64(1)
		switch tok[0] {
		case '-':
			sign = -1
			tok = tok[1:]
		case '+':
			tok = tok[1:]
		}
		if tok == "" {
			return 0, fmt.Errorf("Invalid money token '%s'.", tok)
		}

		// Permite tanto tokens simples ("100s", "2G", "234566") como encadenados ("2g500s25c").
		i := 0
		for i < len(tok) {
			start := i
			for i < len(tok) && tok[i] >= '0' && tok[i] <= '9' {
				i++
			}
			if start == i {
				return 0, fmt.Errorf("Invalid money format in '%s'.", tok)
			}
			val, err := strconv.ParseInt(tok[start:i], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("Money number overflow in '%s'.", tok)
			}

			var mult int64 = 1
			if i < len(tok) {
				switch tok[i] {
				case 'G':
					mult = fiestaCopperPerGem
				case 'g':
					mult = fiestaCopperPerGold
				case 's', 'S':
					mult = fiestaCopperPerSilver
				case 'c', 'C':
					mult = 1
				default:
					return 0, fmt.Errorf("Unknown money unit '%c' (valid: G, g, s, c).", tok[i])
				}
				i++
			}
			total += sign * val * mult
		}
	}
	return total, nil
}

// formatFiestaMoney formatea un valor en Copper a la notación legible de Fiesta (ej: 234566 -> "234s 566c").
func formatFiestaMoney(copper int64) string {
	if copper == 0 {
		return "0c"
	}
	neg := ""
	if copper < 0 {
		neg = "-"
		copper = -copper
	}

	gems := copper / fiestaCopperPerGem
	rem := copper % fiestaCopperPerGem
	gold := rem / fiestaCopperPerGold
	rem %= fiestaCopperPerGold
	silver := rem / fiestaCopperPerSilver
	cop := rem % fiestaCopperPerSilver

	var parts []string
	if gems > 0 {
		parts = append(parts, fmt.Sprintf("%dG", gems))
	}
	if gold > 0 {
		parts = append(parts, fmt.Sprintf("%dg", gold))
	}
	if silver > 0 {
		parts = append(parts, fmt.Sprintf("%ds", silver))
	}
	if cop > 0 {
		parts = append(parts, fmt.Sprintf("%dc", cop))
	}
	return neg + strings.Join(parts, " ")
}

func init() {
	registerGMCommand([]string{"learnskill"}, gmCommandSpec{
		MinArgs:           1,
		Usage:             "&learnskill <id|all>",
		LogType:           7,
		DefaultAdminLevel: 50,
		Exec: func(ctx *gmContext) error {
			if ctx.Handler == nil || ctx.Handler.dataStore == nil || ctx.Handler.dataStore.SkillCatalog == nil {
				return errors.New("SkillCatalog unavailable.")
			}
			cat := ctx.Handler.dataStore.SkillCatalog
			player := ctx.Player
			job := player.Char.Job
			className := JobToClassName(job)
			playerLv := int(player.Char.Level)

			learnOne := func(sk *SkillDef) {
				newSkill := &database.CharacterSkill{
					CharNo:     player.Char.ID,
					SkillNo:    sk.ID,
					SkillLevel: 1,
				}
				player.SetLearnedSkill(newSkill)
				if ctx.Handler.charRepo != nil {
					_ = ctx.Handler.charRepo.SaveCharacterSkill(player.Char.ID, newSkill)
				}
				if player.Client != nil {
					sucCmd := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_SKILL_LEARNSUC_CMD)
					sucCmd.WriteUint16(sk.ID)
					_ = player.Client.SendPacket(sucCmd)
				}
			}

			arg := strings.TrimSpace(ctx.Args[0])
			if strings.EqualFold(arg, "all") {
				// Recolectar y ordenar por ID todas las skills de la clase del GM con 0 < DemandLv <= player.Level
				var eligible []*SkillDef
				for _, sk := range cat.ByID {
					if sk == nil {
						continue
					}
					if !cat.CanClassUseSkill(sk.UseClass, job) {
						continue
					}
					if sk.DemandLv <= 0 || sk.DemandLv > playerLv {
						continue
					}
					eligible = append(eligible, sk)
				}
				sort.Slice(eligible, func(i, j int) bool {
					return eligible[i].ID < eligible[j].ID
				})

				learned := 0
				skipped := 0
				for _, sk := range eligible {
					if player.HasSkill(sk.ID) {
						skipped++
						continue
					}
					learnOne(sk)
					learned++
				}

				ctx.Int1 = learned
				ctx.Reply = fmt.Sprintf("[GM] Learned %d skills for %s (Lv.%d) (skipped %d already known).",
					learned, className, playerLv, skipped)
				return nil
			}

			rawID, err := strconv.ParseUint(arg, 10, 16)
			if err != nil {
				return errors.New("Invalid skill ID. Usage: &learnskill <id|all>")
			}
			skillID := uint16(rawID)
			sk := cat.Get(skillID)
			if sk == nil {
				return fmt.Errorf("Skill ID %d not found in ActiveSkill.shn.", skillID)
			}
			if !cat.CanClassUseSkill(sk.UseClass, job) {
				return fmt.Errorf("Skill %d (%s) cannot be learned by class %s.", sk.ID, sk.Name, className)
			}
			if player.HasSkill(sk.ID) {
				ctx.Reply = fmt.Sprintf("[GM] Skill %d (%s) is already learned.", sk.ID, sk.Name)
				return nil
			}

			learnOne(sk)
			ctx.Int1 = int(sk.ID)
			ctx.Reply = fmt.Sprintf("[GM] Learned skill %d (%s, DemandLv=%d).", sk.ID, sk.Name, sk.DemandLv)
			return nil
		},
	})
}
