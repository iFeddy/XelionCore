package zone

import (
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/abstate"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// ActiveAbState modela una instancia de estado anormal activo con sus atributos oficiales de combate.
type ActiveAbState struct {
	ID            uint16
	Name          string
	Level         int
	Duration      time.Duration
	ExpiresAt     time.Time
	DotDamage     int
	DotInterval   time.Duration
	LastDotTick   time.Time
	IsBuff        bool
	IsDebuff      bool
	IsStun        bool
	IsRoot        bool
	IsSilence     bool
	ShieldHP      int
	SpeedDownRate int
	AtkMod        int
	DefMod        int
	WCMod         int
	ACMod         int
	MAMod         int
	MRMod         int
	THMod         int
	TBMod         int
	MaxHPMod      int
	MaxSPMod      int
	PolymorphID   uint16
}

// applySubAbStateActions traduce los 4 pares ActionIndex/ActionArg oficiales de SubAbstate.shn a modificadores de combate.
func applySubAbStateActions(st *ActiveAbState, actions [4]SubAbStateAction) {
	if st == nil {
		return
	}
	for _, a := range actions {
		switch a.Index {
		case 1: // WCPlus (daño físico)
			st.WCMod += a.Arg
		case 3, 74: // ACPlus (defensa física)
			st.ACMod += a.Arg
		case 7: // TBPlus (evasión)
			st.TBMod += a.Arg
		case 9: // MaxHPPlus
			st.MaxHPMod += a.Arg
		case 11: // MaxSPPlus
			st.MaxSPMod += a.Arg
		case 13: // MAPlus (daño mágico, ej. SubStaMentalShout +11)
			st.MAMod += a.Arg
		case 15: // MRPlus (defensa mágica, ej. SubStaMentalShout +11)
			st.MRMod += a.Arg
		case 17: // Shield absorb (ej. SubStaDivineShield Invincible=1040, SubStaNeglect=1e9)
			st.ShieldHP += a.Arg
		case 19: // Stun (ej. SubStaBattleBlowStun)
			st.IsStun = true
		case 21: // THMinus (ej. SubStaSeverBone)
			st.THMod -= a.Arg
		case 25: // Root / Hold
			st.IsRoot = true
		case 30: // HP Heal Over Time (ej. SubStaRestore)
			if a.Arg > 0 {
				st.DotDamage = -a.Arg
				st.DotInterval = 2 * time.Second
			}
		case 73: // Bleed/Poison DoT (ej. SubStaRedSlash)
			if a.Arg > 0 {
				st.DotDamage = a.Arg
				st.DotInterval = 2 * time.Second
			}
		case 94: // WCMinus
			st.WCMod -= a.Arg
		}
	}
}

// AbsorbDamageWithShield absorbe daño entrante usando escudos activos (ActionIndex=17: Invincible, Neglect, etc.).
// Si un escudo se agota, remueve el AbState y notifica 0x2429(0) + 0x2428. Retorna el daño restante (>= 0).
func (p *ZonePlayer) AbsorbDamageWithShield(zMap *ZoneMap, dmg int) int {
	if p == nil || dmg <= 0 || len(p.ActiveStates) == 0 {
		return dmg
	}
	var broken []uint16
	for id, st := range p.ActiveStates {
		if st == nil || st.ShieldHP <= 0 {
			continue
		}
		if st.ShieldHP >= dmg {
			st.ShieldHP -= dmg
			dmg = 0
			if st.ShieldHP == 0 {
				broken = append(broken, id)
			}
			break
		}
		dmg -= st.ShieldHP
		st.ShieldHP = 0
		broken = append(broken, id)
	}
	for _, id := range broken {
		p.RemoveAbState(id)
		if p.Client != nil {
			inform := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATEINFORM_CMD)
			inform.WriteUint32(uint32(id))
			inform.WriteUint32(0)
			_ = p.Client.SendPacket(inform)
		}
		if zMap != nil {
			resetCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATERESET_CMD)
			resetCmd.WriteUint16(p.MapObjectID)
			resetCmd.WriteUint32(uint32(id))
			zMap.Broadcast(resetCmd, 0)
		}
	}
	return dmg
}

// RecalculateAbStateStats recalcula los bonos acumulados de todos los AbStates activos del jugador.
func (p *ZonePlayer) RecalculateAbStateStats() {
	if p == nil {
		return
	}
	var wc, ac, ma, mr, th, tb, maxHP, maxSP int
	for _, st := range p.ActiveStates {
		if st == nil {
			continue
		}
		wc += st.WCMod
		ac += st.ACMod
		ma += st.MAMod
		mr += st.MRMod
		th += st.THMod
		tb += st.TBMod
		maxHP += st.MaxHPMod
		maxSP += st.MaxSPMod
	}
	p.AbWC = wc
	p.AbAC = ac
	p.AbMA = ma
	p.AbMR = mr
	p.AbTH = th
	p.AbTB = tb
	p.AbMaxHP = maxHP
	p.AbMaxSP = maxSP
}

func clampStatUint32(base uint32, mod int) uint32 {
	v := int(base) + mod
	if v < 0 {
		return 0
	}
	return uint32(v)
}

// syncPlayerAbStateParams recalcula los bonos de AbState y notifica al cliente vía PROTO_NC_CHAR_BASEPARAMCHANGE_CMD
// cualquier cambio en WC, AC, TH, TB, MA, MR, MaxHP o MaxSP.
func (h *Handler) syncPlayerAbStateParams(p *ZonePlayer) {
	if p == nil {
		return
	}
	oldMinWC := clampStatUint32(p.EquipMinWC, p.AbWC)
	oldMaxWC := clampStatUint32(p.EquipMaxWC, p.AbWC)
	oldAC := clampStatUint32(p.EquipAC, p.AbAC)
	oldTH := clampStatUint32(p.EquipTH, p.AbTH)
	oldTB := clampStatUint32(p.EquipTB, p.AbTB)
	oldMinMA := clampStatUint32(p.EquipMinMA, p.AbMA)
	oldMaxMA := clampStatUint32(p.EquipMaxMA, p.AbMA)
	oldMR := clampStatUint32(p.EquipMR, p.AbMR)
	oldMaxHP := clampStatUint32(uint32(p.MaxHP), p.AbMaxHP)
	oldMaxSP := clampStatUint32(uint32(p.MaxSP), p.AbMaxSP)

	p.RecalculateAbStateStats()

	if p.Client == nil {
		return
	}
	type paramDelta struct {
		what uint8
		val  uint32
	}
	candidates := []struct {
		what uint8
		oldV uint32
		newV uint32
	}{
		{0x06, oldMinWC, clampStatUint32(p.EquipMinWC, p.AbWC)},
		{0x07, oldMaxWC, clampStatUint32(p.EquipMaxWC, p.AbWC)},
		{0x08, oldAC, clampStatUint32(p.EquipAC, p.AbAC)},
		{0x09, oldTH, clampStatUint32(p.EquipTH, p.AbTH)},
		{0x0A, oldTB, clampStatUint32(p.EquipTB, p.AbTB)},
		{0x0B, oldMinMA, clampStatUint32(p.EquipMinMA, p.AbMA)},
		{0x0C, oldMaxMA, clampStatUint32(p.EquipMaxMA, p.AbMA)},
		{0x0D, oldMR, clampStatUint32(p.EquipMR, p.AbMR)},
		{0x10, oldMaxHP, clampStatUint32(uint32(p.MaxHP), p.AbMaxHP)},
		{0x11, oldMaxSP, clampStatUint32(uint32(p.MaxSP), p.AbMaxSP)},
	}
	var diffs []paramDelta
	for _, c := range candidates {
		if c.oldV != c.newV {
			diffs = append(diffs, paramDelta{what: c.what, val: c.newV})
		}
	}
	if len(diffs) > 0 {
		pw := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_BASEPARAMCHANGE_CMD)
		pw.WriteByte(byte(len(diffs)))
		for _, d := range diffs {
			pw.WriteByte(d.what)
			pw.WriteUint32(d.val)
		}
		_ = p.Client.SendPacket(pw)
	}
}

// ComputeAbStateBit genera el bitfield oficial de 99 bytes (792 bits) a partir de los estados activos.
func ComputeAbStateBit(states map[uint16]*ActiveAbState) [99]byte {
	var bitfield [99]byte
	for id := range states {
		if id >= 792 {
			continue
		}
		byteIdx := id / 8
		bit := id % 8
		bitfield[byteIdx] |= (1 << bit)
	}
	return bitfield
}

// AddAbState incorpora un estado anormal activo al jugador y recalcula flags.
func (p *ZonePlayer) AddAbState(st *ActiveAbState) {
	if st == nil {
		return
	}
	if p.ActiveStates == nil {
		p.ActiveStates = make(map[uint16]*ActiveAbState)
	}
	p.ActiveStates[st.ID] = st
}

// RemoveAbState elimina un estado anormal por su ID oficial.
func (p *ZonePlayer) RemoveAbState(id uint16) {
	if p.ActiveStates != nil {
		delete(p.ActiveStates, id)
	}
}

// GetAbState obtiene un estado anormal por su ID si está presente.
func (p *ZonePlayer) GetAbState(id uint16) *ActiveAbState {
	if p.ActiveStates == nil {
		return nil
	}
	return p.ActiveStates[id]
}

// GetAbStateBit devuelve el bitfield oficial de 99 bytes correspondiente a sus estados activos.
func (p *ZonePlayer) GetAbStateBit() [99]byte {
	return ComputeAbStateBit(p.ActiveStates)
}

// CanMove verifica si el jugador puede desplazarse (falso si está aturdido o inmovilizado).
func (p *ZonePlayer) CanMove() bool {
	if p.State == 3 { // Dead
		return false
	}
	for _, st := range p.ActiveStates {
		if st.IsStun || st.IsRoot {
			return false
		}
	}
	return true
}

// CanAttack verifica si el jugador puede ejecutar golpes físicos (falso si está aturdido).
func (p *ZonePlayer) CanAttack() bool {
	if p.State == 3 {
		return false
	}
	for _, st := range p.ActiveStates {
		if st.IsStun {
			return false
		}
	}
	return true
}

// CanCast verifica si el jugador puede canalizar habilidades (falso si está aturdido o silenciado).
func (p *ZonePlayer) CanCast() bool {
	if p.State == 3 {
		return false
	}
	for _, st := range p.ActiveStates {
		if st.IsStun || st.IsSilence {
			return false
		}
	}
	return true
}

// IsSilenced verifica si el jugador tiene un estado de silencio activo (Zone.exe smo_SkillBlastOption == 2 -> 0x0FC7).
func (p *ZonePlayer) IsSilenced() bool {
	for _, st := range p.ActiveStates {
		if st != nil && st.IsSilence {
			return true
		}
	}
	return false
}

// IsStunned verifica si el jugador tiene un estado de aturdimiento activo.
func (p *ZonePlayer) IsStunned() bool {
	for _, st := range p.ActiveStates {
		if st != nil && st.IsStun {
			return true
		}
	}
	return false
}

// HasAbState verifica si el jugador tiene activo un AbState por su ID oficial.
func (p *ZonePlayer) HasAbState(id uint16) bool {
	return p.GetAbState(id) != nil
}

// AbStateEngine gobierna la aplicación de estados, cálculo de ticks DoT/HoT y limpieza de expirados.
type AbStateEngine struct {
	mu sync.RWMutex
	db *abstate.Database
}

// NewAbStateEngine inicializa el motor de estados anormales.
func NewAbStateEngine(db *abstate.Database) *AbStateEngine {
	return &AbStateEngine{
		db: db,
	}
}

// CreateActiveState instancia un ActiveAbState a partir de la definición JSON o SHN.
func (e *AbStateEngine) CreateActiveState(stateID uint16, level int) *ActiveAbState {
	st := &ActiveAbState{
		ID:          stateID,
		Level:       level,
		DotInterval: 1 * time.Second,
	}

	if e.db != nil {
		if cfg, ok := e.db.GetByID(stateID); ok && cfg != nil {
			st.Name = cfg.Name
			if lvlCfg, err := cfg.GetLevelConfig(level); err == nil && lvlCfg != nil {
				st.Duration = time.Duration(lvlCfg.DurationMs) * time.Millisecond
				if st.Duration == 0 {
					st.Duration = 10 * time.Second
				}
				st.ExpiresAt = time.Now().Add(st.Duration)
				st.DotDamage = int(lvlCfg.DotDamage)
				if lvlCfg.DotIntervalMs > 0 {
					st.DotInterval = time.Duration(lvlCfg.DotIntervalMs) * time.Millisecond
				}
				st.SpeedDownRate = int(lvlCfg.SpeedDownRate)
				st.ShieldHP = int(lvlCfg.MaxShield)
			}
		}
	}

	if st.Duration == 0 {
		st.Duration = 30 * time.Second
		st.ExpiresAt = time.Now().Add(st.Duration)
	}

	return st
}

// handleDepolymorphCmd implementa PROTO_NC_CHAR_DEPOLYMORPH_CMD (0x10B0).
func (h *Handler) handleDepolymorphCmd(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	// Remueve estado de polimorfismo si existe
	var removedID uint16
	for id, st := range player.ActiveStates {
		if st.PolymorphID > 0 {
			removedID = id
			player.RemoveAbState(id)
			break
		}
	}

	// Notificar actualización de estado según Zone.exe (0x2429 ABSTATEINFORM_CMD 8B + 0x2428 ABSTATERESET_CMD 6B)
	if removedID != 0 {
		if player.Client != nil {
			inform := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATEINFORM_CMD)
			inform.WriteUint32(uint32(removedID))
			inform.WriteUint32(0)
			_ = player.Client.SendPacket(inform)
		}
		if player.CurrentMap != nil {
			resetCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATERESET_CMD)
			resetCmd.WriteUint16(player.MapObjectID)
			resetCmd.WriteUint32(uint32(removedID))
			player.CurrentMap.Broadcast(resetCmd, 0)
		}
	}
}

// handleChargedBuffEraseReq implementa PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ (0x112E).
// Struct oficial PDB: PROTO_NC_CHAR_CHARGEDBUFF_ERASE_REQ (4 bytes: unsigned int nChargedBuffKey).
func (h *Handler) handleChargedBuffEraseReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil || player.Char == nil {
		return
	}

	var key uint32
	if pkt.Reader.Remaining() >= 4 {
		v, err := pkt.Reader.ReadUint32()
		if err != nil {
			return
		}
		key = v
	} else if pkt.Reader.Remaining() >= 2 {
		v, err := pkt.Reader.ReadUint16()
		if err != nil {
			return
		}
		key = uint32(v)
	} else {
		return
	}

	if h.charRepo != nil {
		_ = h.charRepo.DelChargedBuff(player.Char.ID, int(key))
	}

	player.RemoveAbState(uint16(key))

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_CHARGEDBUFF_ERASE_ACK)
	ack.WriteUint16(uint16(key))
	_ = client.SendPacket(ack)
}

// handleAbStateEraseReq implementa PROTO_NC_BAT_ABSTATE_ERASE_REQ (0x2454).
func (h *Handler) handleAbStateEraseReq(client *network.Client, pkt *protocol.Packet) {
	player, ok := h.getPlayer(client)
	if !ok || player == nil {
		return
	}

	stateIndex, err := pkt.Reader.ReadUint16()
	if err != nil {
		return
	}

	st := player.GetAbState(stateIndex)

	// Verificado en Zone.exe sp_NC_BAT_ABSTATE_ERASE_REQ (0x0042B900) / EraseGoodBuff (0x0042ACC0):
	// PROTO_NC_BAT_ABSTATE_ERASE_ACK (0x2455, 2B):
	//   - 0x1C70: borrado exitoso (buff válido dentro de rango < 792)
	//   - 0x1C71: el estado no está activo en el personaje
	//   - 0x1C72: ID fuera de bitfield (>= 0x318 / 792) o debuff no removible manualmente
	var errCode uint16
	removed := false
	switch {
	case stateIndex >= 792:
		errCode = 0x1C72
	case st == nil:
		errCode = 0x1C71
	case st.IsDebuff && !st.IsBuff:
		errCode = 0x1C72
	default:
		player.RemoveAbState(stateIndex)
		errCode = 0x1C70
		removed = true
	}

	ack := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATE_ERASE_ACK)
	ack.WriteUint16(errCode)
	_ = client.SendPacket(ack)

	if removed {
		h.syncPlayerAbStateParams(player)
		inform := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATEINFORM_CMD)
		inform.WriteUint32(uint32(stateIndex))
		inform.WriteUint32(0)
		_ = client.SendPacket(inform)

		if player.CurrentMap != nil {
			resetCmd := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATERESET_CMD)
			resetCmd.WriteUint16(player.MapObjectID)
			resetCmd.WriteUint32(uint32(stateIndex))
			player.CurrentMap.Broadcast(resetCmd, 0)
		}
	}
}
