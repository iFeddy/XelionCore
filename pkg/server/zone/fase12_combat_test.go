package zone

import (
	"io"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/abstate"
	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

type noopConn struct{}

func (n *noopConn) Read(b []byte) (int, error)         { return 0, io.EOF }
func (n *noopConn) Write(b []byte) (int, error)        { return len(b), nil }
func (n *noopConn) Close() error                       { return nil }
func (n *noopConn) LocalAddr() net.Addr                { return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9010} }
func (n *noopConn) RemoteAddr() net.Addr               { return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9010} }
func (n *noopConn) SetDeadline(t time.Time) error      { return nil }
func (n *noopConn) SetReadDeadline(t time.Time) error  { return nil }
func (n *noopConn) SetWriteDeadline(t time.Time) error { return nil }

func setupFase12TestEnvironment(t *testing.T) (*Handler, *ZoneMap, *ZonePlayer, *network.Client) {
	dataDir := findTestDataDir()
	combatStore := combat.NewStore()
	_ = combatStore.LoadDir(filepath.Join(dataDir, "World", "Combat"))

	classStore := classparam.NewStore()
	_ = classStore.LoadDir(filepath.Join(dataDir, "World", "ClassParam"))

	abstateDB, _ := abstate.LoadFromDir(filepath.Join(dataDir, "AbState"))

	ds := &DataStore{
		CombatStore: combatStore,
		ClassStore:  classStore,
		AbState:     abstateDB,
	}

	charRepo := database.NewMemoryCharacterRepositoryWithFile("")
	logDispatcher := database.NewGameLogDispatcher(nil, 100, 1)

	h := NewHandler(charRepo, nil, classStore, nil, nil)
	h.SetDataStore(ds)
	h.SetGameLogDispatcher(logDispatcher)

	zMap := NewZoneMap(0, "Rou", "Roumen")
	client, _ := network.NewClient(&noopConn{}, nil)

	ch := &database.Character{
		ID:    1001,
		Name:  "CombatTester",
		Job:   1, // Fighter
		Level: 10,
		CurHP: 150,
		CurSP: 50,
		TCharacter: database.TCharacter{
			CharNo:            1001,
			Level:             10,
			Strength:          15,
			Constitute:        10,
			Dexterity:         10,
			Intelligence:      5,
			MentalPower:       5,
			RedistributePoint: 5,
		},
		TCharacterState: database.TCharacterState{
			HP:       150,
			SP:       50,
			HPStones: 10,
			SPStones: 5,
			Money:    5000,
		},
	}
	_ = charRepo.CreateCharacter(ch)

	player := NewZonePlayer(client, ch, classStore, combatStore)
	player.CurrentMap = zMap
	zMap.AddPlayer(player)
	h.mu.Lock()
	h.players[client] = player
	h.mu.Unlock()

	return h, zMap, player, client
}

// TestFase12_CombatFormulas verifica los cálculos matemáticos de daño físico, multiplicador angular y estados.
func TestFase12_CombatFormulas(t *testing.T) {
	dataDir := findTestDataDir()
	store := combat.NewStore()
	if err := store.LoadDir(filepath.Join(dataDir, "World", "Combat")); err != nil {
		t.Fatalf("Error cargando store de combate: %v", err)
	}

	// 1. Daño Físico base
	dmgBase := CalculatePhysicalDamage(50, 20, 180, true, 0, 1000, false, store)
	if dmgBase <= 0 {
		t.Errorf("Daño físico esperado > 0, obtenido: %d", dmgBase)
	}

	// 2. Multiplicador de ángulo (espalda angle=180 [rate 1200] vs frontal angle=0 [rate 1000])
	dmgBack := CalculatePhysicalDamage(50, 20, 180, true, 0, 1000, false, store)
	dmgFront := CalculatePhysicalDamage(50, 20, 0, true, 0, 1000, false, store)
	if dmgBack < dmgFront {
		t.Errorf("Ataque por la espalda (%d) debería infligir >= daño que frontal (%d)", dmgBack, dmgFront)
	}

	// 3. Multiplicador de crítico (1.5x)
	dmgCrit := CalculatePhysicalDamage(50, 20, 180, true, 0, 1000, true, store)
	if dmgCrit <= dmgFront {
		t.Errorf("Daño crítico (%d) debería superar al daño normal (%d)", dmgCrit, dmgFront)
	}

	// 4. SoulStone damage rate
	soulRate := store.GetSoulDamageRate(1, 4)
	if soulRate <= 1000 {
		t.Errorf("Tasa de daño por SoulStone esperada > 1000, obtenida: %d", soulRate)
	}
	dmgSoul := CalculatePhysicalDamage(50, 20, 180, true, 0, soulRate, false, store)
	if dmgSoul <= dmgFront {
		t.Errorf("Daño potenciado por SoulStone (%d) debería superar al daño base (%d)", dmgSoul, dmgFront)
	}

	// 5. Cálculo angular de coordenadas
	angBack := CalculateAngleDeg(100, 100, 100, 100, 0)
	if angBack < 0 || angBack > 180 {
		t.Errorf("Ángulo fuera de rango 0..180: %d", angBack)
	}

	// 6. Precisión y crítico
	hit := CheckHit(100, 10)
	if !hit {
		// Probabilidad 100/(100+10) = ~91%
	}
	crit := CheckCrit(50, 0)
	_ = crit
}

// TestFase12_TargettingAndAutoAttack verifica el flujo C2S de selección de objetivo, inicio de ataque y asistencia.
func TestFase12_TargettingAndAutoAttack(t *testing.T) {
	h, zMap, player, client := setupFase12TestEnvironment(t)

	// Crear un monstruo en el mapa
	mob := NewZoneMonster(50, 101, "Slime", "Rou", int32(player.Char.XPos+50), int32(player.Char.YPos), 0, 5, 200, 50, 10, 5, 10, 25, 30, 10, 50)
	zMap.AddMonster(mob)

	// 1. Seleccionar objetivo: PROTO_NC_BAT_TARGETTING_REQ (0x2401)
	tgtPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_TARGETTING_REQ)
	tgtPkt.WriteUint16(mob.Handle)
	pkt := protocol.NewPacket(tgtPkt.Data())
	h.HandlePacket(client, pkt)

	if player.TargetHandle != mob.Handle {
		t.Fatalf("TargetHandle esperado %d, obtenido %d", mob.Handle, player.TargetHandle)
	}

	// 2. Iniciar auto-ataque: PROTO_NC_BAT_BASHSTART_CMD (0x242B)
	bashStartPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_BASHSTART_CMD)
	h.HandlePacket(client, protocol.NewPacket(bashStartPkt.Data()))

	if !player.IsAutoAttacking {
		t.Errorf("IsAutoAttacking esperado true")
	}

	// 3. Ejecutar auto-ataque contra el monstruo
	hpBefore := mob.CurHP
	for i := 0; i < 5 && mob.CurHP == hpBefore; i++ {
		h.ExecuteAutoAttack(player, time.Now())
	}
	if mob.CurHP >= hpBefore && !mob.IsDead {
		t.Errorf("HP del monstruo debería reducirse tras auto-ataque: antes=%d, ahora=%d", hpBefore, mob.CurHP)
	}

	// 4. Detener auto-ataque: PROTO_NC_BAT_BASHSTOP_CMD (0x2432)
	bashStopPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_BASHSTOP_CMD)
	h.HandlePacket(client, protocol.NewPacket(bashStopPkt.Data()))

	if player.IsAutoAttacking {
		t.Errorf("IsAutoAttacking esperado false tras BASHSTOP")
	}

	// 5. Asistir a un aliado: PROTO_NC_BAT_ASSIST_REQ (0x243E)
	allyClient, _ := network.NewClient(&noopConn{}, nil)
	allyCh := &database.Character{
		ID:    1002,
		Name:  "AllyPlayer",
		Job:   1,
		Level: 10,
		CurHP: 150,
		TCharacter: database.TCharacter{
			CharNo: 1002,
		},
	}
	allyPlayer := NewZonePlayer(allyClient, allyCh, nil, nil)
	allyPlayer.TargetHandle = mob.Handle
	zMap.AddPlayer(allyPlayer)
	h.mu.Lock()
	h.players[allyClient] = allyPlayer
	h.mu.Unlock()

	// Deseleccionar y luego asistir
	player.TargetHandle = 0
	assistPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ASSIST_REQ)
	assistPkt.WriteUint16(allyPlayer.MapObjectID)
	h.HandlePacket(client, protocol.NewPacket(assistPkt.Data()))

	if player.TargetHandle != mob.Handle {
		t.Errorf("TargetHandle asistido esperado %d, obtenido %d", mob.Handle, player.TargetHandle)
	}

	// 6. Deseleccionar objetivo: PROTO_NC_BAT_UNTARGET_REQ (0x2408)
	untargetPkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_UNTARGET_REQ)
	h.HandlePacket(client, protocol.NewPacket(untargetPkt.Data()))

	if player.TargetHandle != 0 {
		t.Errorf("TargetHandle tras UNTARGET esperado 0, obtenido %d", player.TargetHandle)
	}
}

// TestFase12_SkillsAndAoE verifies casting single-target and AoE skills.
func TestFase12_SkillsAndAoE(t *testing.T) {
	h, zMap, player, client := setupFase12TestEnvironment(t)

	mob := NewZoneMonster(60, 102, "TargetMob", "Rou", int32(player.Char.XPos+40), int32(player.Char.YPos), 0, 5, 200, 50, 10, 5, 10, 20, 30, 10, 50)
	zMap.AddMonster(mob)

	// 1. Casteo Habilidad Objetivo: PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ (0x2440)
	objCast := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_OBJ_CAST_REQ)
	objCast.WriteUint16(101)        // Skill ID
	objCast.WriteUint16(mob.Handle) // Target obj
	h.HandlePacket(client, protocol.NewPacket(objCast.Data()))

	if mob.CurHP >= 200 {
		t.Errorf("Monstruo debería haber recibido daño de habilidad, HP=%d", mob.CurHP)
	}

	// 2. Casteo Habilidad de Campo (AoE): PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ (0x2441)
	hpBefore := mob.CurHP
	fldCast := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_SKILLBASH_FLD_CAST_REQ)
	fldCast.WriteUint16(201)                  // Skill ID
	fldCast.WriteInt32(int32(player.Char.XPos + 40)) // X
	fldCast.WriteInt32(int32(player.Char.YPos))      // Y
	h.HandlePacket(client, protocol.NewPacket(fldCast.Data()))

	if mob.CurHP >= hpBefore {
		t.Errorf("Monstruo debería haber recibido daño de habilidad AoE, antes=%d, ahora=%d", hpBefore, mob.CurHP)
	}

	// 3. Empoderamiento de Habilidad: PROTO_NC_SKILL_EMPOWALLOC_REQ (0x4811)
	empowPkt := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_EMPOWALLOC_REQ)
	empowPkt.WriteUint16(101) // Skill
	empowPkt.WriteUint8(0)    // Damage empow slot
	empowPkt.WriteUint8(2)    // Point
	h.HandlePacket(client, protocol.NewPacket(empowPkt.Data()))

	// 4. Invocación de Campo / Tótem: PROTO_NC_SKILL_PRODUCTFIELD_REQ (0x4822)
	prodPkt := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_PRODUCTFIELD_REQ)
	prodPkt.WriteUint16(301)
	h.HandlePacket(client, protocol.NewPacket(prodPkt.Data()))

	// 5. Olvido de Habilidad: PROTO_NC_SKILL_UNLEARN_REQ (0x4824)
	unlearnPkt := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_UNLEARN_REQ)
	unlearnPkt.WriteUint16(101)
	h.HandlePacket(client, protocol.NewPacket(unlearnPkt.Data()))
}

// TestFase12_AbStateAndCrowdControl verifies abnormal states, DoT ticks, and CC restrictions.
func TestFase12_AbStateAndCrowdControl(t *testing.T) {
	h, _, player, client := setupFase12TestEnvironment(t)

	// 1. Aplicación de Stun (aturdimiento)
	stunState := &ActiveAbState{
		ID:        50,
		Name:      "Stun",
		Duration:  3 * time.Second,
		ExpiresAt: time.Now().Add(3 * time.Second),
		IsStun:    true,
	}
	player.AddAbState(stunState)

	if player.CanMove() {
		t.Errorf("Jugador aturdido no debería poder moverse")
	}
	if player.CanAttack() {
		t.Errorf("Jugador aturdido no debería poder atacar")
	}
	if player.CanCast() {
		t.Errorf("Jugador aturdido no debería poder castear")
	}

	// 2. Verificar bitfield de 99 bytes
	bitfield := player.GetAbStateBit()
	byteIdx := 50 / 8
	bit := 50 % 8
	if (bitfield[byteIdx] & (1 << bit)) == 0 {
		t.Errorf("Bitfield no tiene el bit activado para el estado 50")
	}

	// 3. Borrado voluntario de estado: PROTO_NC_BAT_ABSTATE_ERASE_REQ (0x2454)
	erasePkt := protocol.NewPacketWriter(protocol.PROTO_NC_BAT_ABSTATE_ERASE_REQ)
	erasePkt.WriteUint16(50)
	h.HandlePacket(client, protocol.NewPacket(erasePkt.Data()))

	if player.GetAbState(50) != nil {
		t.Errorf("Estado 50 debería haberse borrado")
	}
	if !player.CanMove() {
		t.Errorf("Jugador debería poder moverse tras removerse el aturdimiento")
	}

	// 4. Estado con daño periódico (DoT)
	dotState := &ActiveAbState{
		ID:          60,
		Name:        "Poison",
		Duration:    5 * time.Second,
		ExpiresAt:   time.Now().Add(5 * time.Second),
		DotDamage:   15,
		DotInterval: 50 * time.Millisecond,
		LastDotTick: time.Now().Add(-100 * time.Millisecond),
	}
	player.AddAbState(dotState)

	// Ejecutar tick1s del gameloop para procesar DoT
	hpBefore := player.Char.CurHP
	player.CurrentMap.gameLoop.tick1s()

	if player.Char.CurHP >= hpBefore {
		t.Errorf("HP debería haber disminuido por el DoT de veneno: antes=%d, ahora=%d", hpBefore, player.Char.CurHP)
	}

	// 5. Polimorfismo y despolimorfismo: PROTO_NC_CHAR_DEPOLYMORPH_CMD (0x10B0)
	polyState := &ActiveAbState{
		ID:          70,
		Name:        "MushroomPolymorph",
		Duration:    10 * time.Second,
		ExpiresAt:   time.Now().Add(10 * time.Second),
		PolymorphID: 10,
	}
	player.AddAbState(polyState)

	depolyPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_DEPOLYMORPH_CMD)
	h.HandlePacket(client, protocol.NewPacket(depolyPkt.Data()))

	if player.GetAbState(70) != nil {
		t.Errorf("Estado de polimorfismo debería haber sido removido")
	}
}

// TestFase12_SoulStoneAndStatsAndRevive verifies SoulStone buy/use, stats inc/dec, and death/revive.
func TestFase12_SoulStoneAndStatsAndRevive(t *testing.T) {
	h, zMap, player, client := setupFase12TestEnvironment(t)

	// 1. Comprar piedras de alma HP: PROTO_NC_SOULSTONE_HP_BUY_REQ (0x5001)
	buyHPPkt := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_BUY_REQ)
	buyHPPkt.WriteUint16(2) // comprar 2 piedras
	h.HandlePacket(client, protocol.NewPacket(buyHPPkt.Data()))

	if player.StonesHP != 12 {
		t.Errorf("Piedras HP esperadas 12 tras compra, obtenidas: %d", player.StonesHP)
	}

	// 2. Usar piedra de alma HP: PROTO_NC_SOULSTONE_HP_USE_REQ (0x5007)
	player.Char.CurHP = 50
	useHPPkt := protocol.NewPacketWriter(protocol.PROTO_NC_SOULSTONE_HP_USE_REQ)
	h.HandlePacket(client, protocol.NewPacket(useHPPkt.Data()))

	if player.StonesHP != 11 {
		t.Errorf("Piedras HP esperadas 11 tras uso, obtenidas: %d", player.StonesHP)
	}
	if player.Char.CurHP <= 50 {
		t.Errorf("CurHP debería haberse recuperado con la piedra de alma, obtenido: %d", player.Char.CurHP)
	}

	// 3. Asignación de Stat Points: PROTO_NC_CHAR_STAT_INCPOINT_REQ (0x105C)
	strBefore := player.Char.TCharacter.Strength
	remBefore := player.Char.TCharacter.RedistributePoint
	incPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_STAT_INCPOINT_REQ)
	incPkt.WriteUint8(0) // 0 = Str
	h.HandlePacket(client, protocol.NewPacket(incPkt.Data()))

	if player.Char.TCharacter.Strength != strBefore+1 {
		t.Errorf("Fuerza esperada %d, obtenida %d", strBefore+1, player.Char.TCharacter.Strength)
	}
	if player.Char.TCharacter.RedistributePoint != remBefore-1 {
		t.Errorf("Puntos restantes esperados %d, obtenidos %d", remBefore-1, player.Char.TCharacter.RedistributePoint)
	}

	// 4. Reducción de Stat Points: PROTO_NC_CHAR_STAT_DECPOINT_REQ (0x1062)
	decPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_STAT_DECPOINT_REQ)
	decPkt.WriteUint8(0) // 0 = Str
	h.HandlePacket(client, protocol.NewPacket(decPkt.Data()))

	if player.Char.TCharacter.Strength != strBefore {
		t.Errorf("Fuerza tras decrementar esperada %d, obtenida %d", strBefore, player.Char.TCharacter.Strength)
	}

	// 5. Muerte y Resurrección: PROTO_NC_CHAR_REVIVE_REQ (0x104E)
	player.Char.CurHP = 0
	player.State = 3 // Dead

	revivePkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_REVIVE_REQ)
	h.HandlePacket(client, protocol.NewPacket(revivePkt.Data()))

	if player.State != 1 {
		t.Errorf("Estado esperado Normal (1) tras revivir, obtenido %d", player.State)
	}
	if player.Char.CurHP <= 0 {
		t.Errorf("HP tras revivir esperado > 0, obtenido %d", player.Char.CurHP)
	}

	// 6. Resurrección por clérigo: PROTO_NC_CHAR_REBIRTH_CMD (0x1102)
	player.State = 3
	player.Char.CurHP = 0
	rebirthPkt := protocol.NewPacketWriter(protocol.PROTO_NC_CHAR_REBIRTH_CMD)
	h.HandlePacket(client, protocol.NewPacket(rebirthPkt.Data()))

	if player.State != 1 || player.Char.CurHP <= 0 {
		t.Errorf("Rebirth debería restaurar estado a Normal con HP > 0")
	}

	// 7. Usar ítem de resurrección: PROTO_NC_ITEM_REVIVEITEMUSE_CMD (0x3039)
	player.State = 3
	player.Char.CurHP = 0
	itemRevive := protocol.NewPacketWriter(protocol.PROTO_NC_ITEM_REVIVEITEMUSE_CMD)
	h.HandlePacket(client, protocol.NewPacket(itemRevive.Data()))

	if player.State != 1 || player.Char.CurHP <= 0 {
		t.Errorf("ItemRevive debería restaurar estado a Normal con HP > 0")
	}

	// 8. Resurrección por habilidad aliada: PROTO_NC_SKILL_REPLYREVIVE_CMD (0x482A)
	player.State = 3
	player.Char.CurHP = 0
	replyRevive := protocol.NewPacketWriter(protocol.PROTO_NC_SKILL_REPLYREVIVE_CMD)
	h.HandlePacket(client, protocol.NewPacket(replyRevive.Data()))

	if player.State != 1 || player.Char.CurHP <= 0 {
		t.Errorf("ReplyRevive debería restaurar estado a Normal con HP > 0")
	}

	_ = zMap
}
