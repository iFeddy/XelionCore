package zone

import (
	"fmt"
	"log"
	"slices"
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

var gmItemKeySeq atomic.Int64

func init() {
	gmItemKeySeq.Store(time.Now().UnixMicro())
}

// AdminLvCatalog indexa los comandos GM oficiales cargados desde data/AdminLvSet.shn.
type AdminLvCatalog struct {
	Levels map[string]byte
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
	cat.LogUnimplementedCommands()
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

// LogUnimplementedCommands reporta al arranque los comandos presentes en AdminLvSet.shn sin handler registrado.
func (c *AdminLvCatalog) LogUnimplementedCommands() {
	if c == nil || len(c.Levels) == 0 {
		return
	}
	var missing []string
	for cmd, lv := range c.Levels {
		if _, ok := gmRegistry[cmd]; !ok {
			missing = append(missing, fmt.Sprintf("%s(lv=%d)", cmd, lv))
		}
	}
	if len(missing) == 0 {
		return
	}
	slices.Sort(missing)
	log.Printf("[Zone][GM] Unimplemented commands in AdminLvSet.shn (%d/%d): %s",
		len(missing), len(c.Levels), strings.Join(missing, ", "))
}

// gmContext encapsula el estado de ejecución y los campos de auditoría TypeL de un comando GM.
type gmContext struct {
	Handler      *Handler
	Player       *ZonePlayer
	Cmd          string
	Args         []string
	ReqLv        byte
	LogType      int
	TargetCharID int
	TargetID     int
	ItemKey      int64
	Int1         int
	Int2         int
	Int3         int
	BigInt1      int64
	Reply        string
}

// gmCommandSpec define la firma, validación mínima y ejecución de un comando GM.
type gmCommandSpec struct {
	MinArgs int
	Usage   string
	LogType int
	Exec    func(ctx *gmContext) error
}

var gmRegistry = map[string]gmCommandSpec{}

func registerGMCommand(names []string, spec gmCommandSpec) {
	if spec.LogType == 0 {
		spec.LogType = 1
	}
	for _, name := range names {
		gmRegistry[strings.ToLower(strings.TrimSpace(name))] = spec
	}
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

	// Comando de Ayuda
	if cmd == "help" {
		if player.Char.AdminLevel == 0 {
			return false
		}
		h.sendGMHelp(player, args)
		return true
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

	spec, ok := gmRegistry[cmd]
	if !ok {
		h.sendPlayerNotice(player, fmt.Sprintf("[GM] Command '%s' is registered in AdminLvSet.shn but not implemented.", cmd))
		return true
	}

	if len(args) < spec.MinArgs && spec.Usage != "" {
		h.sendPlayerNotice(player, "Usage: "+formatCommandUsage(cmd, spec))
		return true
	}

	ctx := &gmContext{
		Handler:      h,
		Player:       player,
		Cmd:          cmd,
		Args:         args,
		ReqLv:        reqLv,
		LogType:      spec.LogType,
		TargetCharID: player.Char.ID,
		Int2:         int(reqLv),
		Int3:         int(player.Char.AdminLevel),
		Reply:        fmt.Sprintf("[GM] Executed '%s'", cmd),
	}

	if err := spec.Exec(ctx); err != nil {
		h.sendPlayerNotice(player, err.Error())
		return true
	}

	if h.gameLogDispatcher != nil {
		h.gameLogDispatcher.EnqueueTypeL(
			ctx.LogType,
			player.Char.ID,
			ctx.TargetCharID,
			ctx.TargetID,
			ctx.ItemKey,
			ctx.Int1,
			ctx.Int2,
			ctx.Int3,
			ctx.BigInt1,
		)
	}

	if ctx.Reply != "" {
		h.sendPlayerNotice(player, ctx.Reply)
	}
	return true
}

func toggleGMPlayerFlag(charID int, flagMap map[int]bool, label string) (int, string) {
	globalGMState.mu.Lock()
	next := !flagMap[charID]
	flagMap[charID] = next
	globalGMState.mu.Unlock()
	return boolToInt(next), fmt.Sprintf("[GM] %s toggled=%v", label, next)
}

// formatCommandUsage adapta la firma Usage al alias exacto invocado (ej. "&expbonus [rate%]" -> "&itembonus [rate%]").
func formatCommandUsage(cmd string, spec gmCommandSpec) string {
	if spec.Usage == "" {
		return "&" + cmd
	}
	if _, rest, hasArgs := strings.Cut(spec.Usage, " "); hasArgs {
		return "&" + cmd + " " + rest
	}
	return "&" + cmd
}

// sendGMHelp lista los comandos disponibles para el AdminLevel del jugador o el detalle de uno puntual.
func (h *Handler) sendGMHelp(player *ZonePlayer, args []string) {
	if player == nil || player.Char == nil || h.dataStore == nil || h.dataStore.AdminLvCatalog == nil {
		return
	}
	adminLv := player.Char.AdminLevel

	// Modo detalle: &help <comando>
	if len(args) > 0 {
		target := strings.ToLower(strings.TrimSpace(args[0]))
		reqLv, inSHN := h.dataStore.AdminLvCatalog.RequiredLevel(target)
		spec, inReg := gmRegistry[target]
		if !inSHN && !inReg {
			h.sendPlayerNotice(player, fmt.Sprintf("[GM] Unknown command '%s'.", target))
			return
		}
		if inSHN && adminLv < reqLv {
			h.sendPlayerNotice(player, fmt.Sprintf("[GM] '%s' requires AdminLevel %d (yours: %d).", target, reqLv, adminLv))
			return
		}
		usage := formatCommandUsage(target, spec)
		status := ""
		if !inReg {
			status = " [unimplemented]"
		}
		h.sendPlayerNotice(player, fmt.Sprintf("[GM] %s (Lv.%d)%s", usage, reqLv, status))
		return
	}

	// Modo listado: muestra todos los comandos habilitados para su AdminLevel actual
	var available []string
	for cmd, reqLv := range h.dataStore.AdminLvCatalog.Levels {
		if adminLv < reqLv {
			continue
		}
		spec, ok := gmRegistry[cmd]
		if !ok {
			continue
		}
		available = append(available, formatCommandUsage(cmd, spec))
	}

	if len(available) == 0 {
		h.sendPlayerNotice(player, fmt.Sprintf("[GM] No commands available for AdminLevel %d.", adminLv))
		return
	}

	slices.Sort(available)
	h.sendPlayerNotice(player, fmt.Sprintf("[GM] Commands for AdminLevel %d (%d total):", adminLv, len(available)))

	// Paginar en notices de <= 230 bytes para respetar el tope de PROTO_NC_ACT_NOTICE_CMD
	var chunk strings.Builder
	for _, entry := range available {
		if chunk.Len() > 0 && chunk.Len()+len(entry)+3 > 230 {
			h.sendPlayerNotice(player, chunk.String())
			chunk.Reset()
		}
		if chunk.Len() > 0 {
			chunk.WriteString(" | ")
		}
		chunk.WriteString(entry)
	}
	if chunk.Len() > 0 {
		h.sendPlayerNotice(player, chunk.String())
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
