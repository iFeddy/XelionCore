package zone

import (
	"sync"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/abstate"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// Constantes oficiales de enum SubAbstateAction (0..120) extraídas directamente de re/pdo/Zone.pdb (TPI LF_ENUMERATE).
const (
	SAA_NONE                = 0
	SAA_STRRATE             = 1
	SAA_STRPLUS             = 2
	SAA_WCPLUS              = 3
	SAA_WCRATE              = 4
	SAA_ACPLUS              = 5
	SAA_ACRATE              = 6
	SAA_DEXPLUS             = 7
	SAA_TBPLUS              = 8
	SAA_TBRATE              = 9
	SAA_THPLUS              = 10
	SAA_THRATE              = 11
	SAA_INTPLUS             = 12
	SAA_MAPLUS              = 13
	SAA_MENTALPLUS          = 14
	SAA_MRPLUS              = 15
	SAA_MRRATE              = 16
	SAA_SHIELDAMOUNT        = 17
	SAA_SHIELDACRATE        = 18
	SAA_NOMOVE              = 19
	SAA_SPEEDRATE           = 20
	SAA_ATTACKSPEEDRATE     = 21
	SAA_MAXHPRATE           = 22
	SAA_MAXSPRATE           = 23
	SAA_DEADHPSPRECOVRATE   = 24
	SAA_NOATTACK            = 25
	SAA_TICK                = 26
	SAA_DOTDAMAGE           = 27
	SAA_CONHEAL             = 28
	SAA_CASTINGTIMEPLUS     = 29
	SAA_HEALAMOUNT          = 30
	SAA_POISONRESISTRATE    = 31
	SAA_DISEASERESISTRATE   = 32
	SAA_CURSERESISTRATE     = 33
	SAA_CRITICALRATE        = 34
	SAA_MAXHPPLUS           = 35
	SAA_MAXSPPLUS           = 36
	SAA_INTRATE             = 37
	SAA_FEAR                = 38
	SAA_ALLSTATEPLUS        = 39
	SAA_REVIVEHEALRATE      = 40
	SAA_COUNT               = 41
	SAA_SILIENCE            = 42
	SAA_DEADLYBLESSING      = 43
	SAA_DAMAGERATE          = 44
	SAA_MARATE              = 46
	SAA_HEALRATE            = 47
	SAA_TOTALDAMAGERATE     = 50
	SAA_REFLECTDAMAGE       = 61
	SAA_HIDEENEMY           = 65
	SAA_EVASIONAMOUNT       = 71
	SAA_USESPRATE           = 72
	SAA_ACMINUS             = 73
	SAA_ACDOWNRATE          = 74
	SAA_ATKSPEEDDOWNRATE    = 78
	SAA_CRITICALDOWNRATE    = 80
	SAA_DEXMINUS            = 81
	SAA_HEALAMOUNTMINUS     = 82
	SAA_MAMINUS             = 83
	SAA_MADOWNRATE          = 84
	SAA_MAXHPDOWNRATE       = 85
	SAA_MRMINUS             = 86
	SAA_MRDOWNRATE          = 87
	SAA_SPEEDDOWNRATE       = 88
	SAA_STRMINUS            = 89
	SAA_TBMINUS             = 90
	SAA_TBDOWNRATE          = 91
	SAA_THMINUS             = 92
	SAA_THDOWNRATE          = 93
	SAA_WCMINUS             = 94
	SAA_WCDOWNRATE          = 95
	SAA_MENDOWNRATE         = 99
	SAA_USESPDOWN           = 100
	SAA_CRIUPRATE           = 101
	SAA_MRSHIELDRATE        = 102
	SAA_ACSHIELDRATE        = 103
	SAA_SHIELDRATE          = 112
	SAA_MINHP               = 114
	SAA_DMGDOWNRATE         = 115
	SAA_ALLSTATPLUS         = 119
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
	IsNoAttack    bool
	IsSilence     bool
	IsHide        bool
	ShieldHP      int
	SpeedDownRate int
	AtkMod        int
	DefMod        int
	WCMod         int
	WCRate        int
	ACMod         int
	ACRate        int
	MAMod         int
	MARate        int
	MRMod         int
	MRRate        int
	THMod         int
	THRate        int
	TBMod         int
	TBRate        int
	MaxHPMod      int
	MaxHPRate     int
	MaxSPMod      int
	MaxSPRate     int
	CritRate      int
	DamageRate    int
	DmgDownRate   int
	ReflectRate   int
	HealRate      int
	CastTimeModMs int
	MinHPFloor    int
	PolymorphID   uint16
}

// applySubAbStateActions traduce los 4 pares ActionIndex/ActionArg oficiales de SubAbstate.shn
// usando el enum oficial SubAbstateAction (SAA_*, 0..120) extraído de re/pdo/Zone.pdb.
func applySubAbStateActions(st *ActiveAbState, actions [4]SubAbStateAction, targetMaxHP int) {
	if st == nil {
		return
	}
	hasNoMove := false
	hasNoAttack := false
	for _, a := range actions {
		switch a.Index {
		case SAA_NONE:
			continue
		case SAA_STRRATE:
			st.WCRate += a.Arg
		case SAA_STRPLUS:
			st.WCMod += a.Arg
		case SAA_WCPLUS:
			st.WCMod += a.Arg
		case SAA_WCRATE:
			st.WCRate += a.Arg
		case SAA_ACPLUS:
			st.ACMod += a.Arg
		case SAA_ACRATE, SAA_SHIELDACRATE, SAA_ACSHIELDRATE:
			st.ACRate += a.Arg
		case SAA_DEXPLUS:
			st.TBMod += a.Arg
			st.THMod += a.Arg
		case SAA_TBPLUS, SAA_EVASIONAMOUNT:
			st.TBMod += a.Arg
		case SAA_TBRATE:
			st.TBRate += a.Arg
		case SAA_THPLUS:
			st.THMod += a.Arg
		case SAA_THRATE:
			st.THRate += a.Arg
		case SAA_INTPLUS:
			st.MAMod += a.Arg
		case SAA_MAPLUS:
			st.MAMod += a.Arg
		case SAA_MENTALPLUS:
			st.MRMod += a.Arg
		case SAA_MRPLUS:
			st.MRMod += a.Arg
		case SAA_MRRATE, SAA_MRSHIELDRATE:
			st.MRRate += a.Arg
		case SAA_SHIELDAMOUNT:
			st.ShieldHP += a.Arg
		case SAA_SHIELDRATE:
			if targetMaxHP > 0 && a.Arg > 0 {
				st.ShieldHP += (targetMaxHP * a.Arg) / 1000
			}
		case SAA_NOMOVE:
			st.IsRoot = true
			hasNoMove = true
		case SAA_SPEEDDOWNRATE:
			st.SpeedDownRate += a.Arg
		case SAA_MAXHPRATE:
			st.MaxHPRate += a.Arg
		case SAA_MAXSPRATE:
			st.MaxSPRate += a.Arg
		case SAA_NOATTACK:
			st.IsNoAttack = true
			hasNoAttack = true
		case SAA_TICK:
			if a.Arg > 0 {
				st.DotInterval = time.Duration(a.Arg) * time.Millisecond
			}
		case SAA_DOTDAMAGE:
			if a.Arg > 0 {
				st.DotDamage += a.Arg
				if st.DotInterval == 0 {
					st.DotInterval = 2 * time.Second
				}
			}
		case SAA_CONHEAL, SAA_HEALAMOUNT:
			if a.Arg > 0 {
				st.DotDamage -= a.Arg
				if st.DotInterval == 0 {
					st.DotInterval = 2 * time.Second
				}
			}
		case SAA_CASTINGTIMEPLUS:
			st.CastTimeModMs += a.Arg
		case SAA_CRITICALRATE, SAA_CRIUPRATE:
			st.CritRate += a.Arg
		case SAA_MAXHPPLUS:
			st.MaxHPMod += a.Arg
		case SAA_MAXSPPLUS:
			st.MaxSPMod += a.Arg
		case SAA_INTRATE:
			st.MARate += a.Arg
		case SAA_FEAR:
			st.IsStun = true
		case SAA_ALLSTATEPLUS, SAA_ALLSTATPLUS:
			st.WCMod += a.Arg
			st.ACMod += a.Arg
			st.MAMod += a.Arg
			st.MRMod += a.Arg
			st.THMod += a.Arg
			st.TBMod += a.Arg
		case SAA_SILIENCE:
			st.IsSilence = true
		case SAA_DAMAGERATE, SAA_TOTALDAMAGERATE:
			st.DamageRate += a.Arg
		case SAA_MARATE:
			st.MARate += a.Arg
		case SAA_HEALRATE:
			st.HealRate += a.Arg
		case SAA_REFLECTDAMAGE:
			st.ReflectRate += a.Arg
		case SAA_HIDEENEMY:
			st.IsHide = true
		case SAA_ACMINUS:
			st.ACMod -= a.Arg
		case SAA_ACDOWNRATE:
			st.ACRate -= a.Arg
		case SAA_CRITICALDOWNRATE:
			st.CritRate -= a.Arg
		case SAA_DEXMINUS:
			st.TBMod -= a.Arg
			st.THMod -= a.Arg
		case SAA_MAMINUS:
			st.MAMod -= a.Arg
		case SAA_MADOWNRATE:
			st.MARate -= a.Arg
		case SAA_MAXHPDOWNRATE:
			st.MaxHPRate -= a.Arg
		case SAA_MRMINUS:
			st.MRMod -= a.Arg
		case SAA_MRDOWNRATE, SAA_MENDOWNRATE:
			st.MRRate -= a.Arg
		case SAA_STRMINUS:
			st.WCMod -= a.Arg
		case SAA_TBMINUS:
			st.TBMod -= a.Arg
		case SAA_TBDOWNRATE:
			st.TBRate -= a.Arg
		case SAA_THMINUS:
			st.THMod -= a.Arg
		case SAA_THDOWNRATE:
			st.THRate -= a.Arg
		case SAA_WCMINUS:
			st.WCMod -= a.Arg
		case SAA_WCDOWNRATE:
			st.WCRate -= a.Arg
		case SAA_MINHP:
			if a.Arg > st.MinHPFloor {
				st.MinHPFloor = a.Arg
			}
		case SAA_DMGDOWNRATE:
			st.DmgDownRate += a.Arg
		}
	}
	if hasNoMove && hasNoAttack {
		st.IsStun = true
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

// RecalculateAbStateStats recalcula los bonos planos y porcentuales (permille) de todos los AbStates activos del jugador.
func (p *ZonePlayer) RecalculateAbStateStats() {
	if p == nil {
		return
	}
	var wc, wcR, ac, acR, ma, maR, mr, mrR, th, thR, tb, tbR, maxHP, maxHPR, maxSP, maxSPR int
	var critR, dmgR, dmgDownR, reflR, healR, castMod, minHP int
	for _, st := range p.ActiveStates {
		if st == nil {
			continue
		}
		wc += st.WCMod
		wcR += st.WCRate
		ac += st.ACMod
		acR += st.ACRate
		ma += st.MAMod
		maR += st.MARate
		mr += st.MRMod
		mrR += st.MRRate
		th += st.THMod
		thR += st.THRate
		tb += st.TBMod
		tbR += st.TBRate
		maxHP += st.MaxHPMod
		maxHPR += st.MaxHPRate
		maxSP += st.MaxSPMod
		maxSPR += st.MaxSPRate
		critR += st.CritRate
		dmgR += st.DamageRate
		dmgDownR += st.DmgDownRate
		reflR += st.ReflectRate
		healR += st.HealRate
		castMod += st.CastTimeModMs
		if st.MinHPFloor > minHP {
			minHP = st.MinHPFloor
		}
	}
	p.AbWC, p.AbWCRate = wc, wcR
	p.AbAC, p.AbACRate = ac, acR
	p.AbMA, p.AbMARate = ma, maR
	p.AbMR, p.AbMRRate = mr, mrR
	p.AbTH, p.AbTHRate = th, thR
	p.AbTB, p.AbTBRate = tb, tbR
	p.AbMaxHP, p.AbMaxHPRate = maxHP, maxHPR
	p.AbMaxSP, p.AbMaxSPRate = maxSP, maxSPR
	p.AbCritRate = critR
	p.AbDamageRate = dmgR
	p.AbDmgDownRate = dmgDownR
	p.AbReflectRate = reflR
	p.AbHealRate = healR
	p.AbCastTimeMod = castMod
	p.AbMinHPFloor = minHP
}

func clampStatUint32(base uint32, mod int) uint32 {
	v := int(base) + mod
	if v < 0 {
		return 0
	}
	return uint32(v)
}

func applyStatWithRate(base uint32, flatMod, ratePermille int) uint32 {
	v := int(base) + flatMod
	if ratePermille != 0 {
		v += (v * ratePermille) / 1000
	}
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
	oldMinWC := applyStatWithRate(p.EquipMinWC, p.AbWC, p.AbWCRate)
	oldMaxWC := applyStatWithRate(p.EquipMaxWC, p.AbWC, p.AbWCRate)
	oldAC := applyStatWithRate(p.EquipAC, p.AbAC, p.AbACRate)
	oldTH := applyStatWithRate(p.EquipTH, p.AbTH, p.AbTHRate)
	oldTB := applyStatWithRate(p.EquipTB, p.AbTB, p.AbTBRate)
	oldMinMA := applyStatWithRate(p.EquipMinMA, p.AbMA, p.AbMARate)
	oldMaxMA := applyStatWithRate(p.EquipMaxMA, p.AbMA, p.AbMARate)
	oldMR := applyStatWithRate(p.EquipMR, p.AbMR, p.AbMRRate)
	oldMaxHP := applyStatWithRate(uint32(p.MaxHP), p.AbMaxHP, p.AbMaxHPRate)
	oldMaxSP := applyStatWithRate(uint32(p.MaxSP), p.AbMaxSP, p.AbMaxSPRate)

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
		{0x06, oldMinWC, applyStatWithRate(p.EquipMinWC, p.AbWC, p.AbWCRate)},
		{0x07, oldMaxWC, applyStatWithRate(p.EquipMaxWC, p.AbWC, p.AbWCRate)},
		{0x08, oldAC, applyStatWithRate(p.EquipAC, p.AbAC, p.AbACRate)},
		{0x09, oldTH, applyStatWithRate(p.EquipTH, p.AbTH, p.AbTHRate)},
		{0x0A, oldTB, applyStatWithRate(p.EquipTB, p.AbTB, p.AbTBRate)},
		{0x0B, oldMinMA, applyStatWithRate(p.EquipMinMA, p.AbMA, p.AbMARate)},
		{0x0C, oldMaxMA, applyStatWithRate(p.EquipMaxMA, p.AbMA, p.AbMARate)},
		{0x0D, oldMR, applyStatWithRate(p.EquipMR, p.AbMR, p.AbMRRate)},
		{0x10, oldMaxHP, applyStatWithRate(uint32(p.MaxHP), p.AbMaxHP, p.AbMaxHPRate)},
		{0x11, oldMaxSP, applyStatWithRate(uint32(p.MaxSP), p.AbMaxSP, p.AbMaxSPRate)},
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
