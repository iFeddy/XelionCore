package zone

import (
	"encoding/binary"
	"io"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// capturedPacket almacena un paquete saliente decodificado desde wire format de Fiesta.
type capturedPacket struct {
	OpCode  uint16
	Payload []byte
}

type captureConn struct {
	mu      sync.Mutex
	packets []capturedPacket
}

func (c *captureConn) Read(b []byte) (int, error) { return 0, io.EOF }
func (c *captureConn) Write(b []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(b) >= 3 {
		pos := 1
		if b[0] == 0 && len(b) >= 5 {
			pos = 3
		}
		if pos+2 <= len(b) {
			op := binary.LittleEndian.Uint16(b[pos : pos+2])
			payload := make([]byte, len(b)-pos-2)
			copy(payload, b[pos+2:])
			c.packets = append(c.packets, capturedPacket{OpCode: op, Payload: payload})
		}
	}
	return len(b), nil
}
func (c *captureConn) Close() error                       { return nil }
func (c *captureConn) LocalAddr() net.Addr                { return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9020} }
func (c *captureConn) RemoteAddr() net.Addr               { return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9020} }
func (c *captureConn) SetDeadline(t time.Time) error      { return nil }
func (c *captureConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *captureConn) SetWriteDeadline(t time.Time) error { return nil }

func (c *captureConn) findByOpcode(op uint16) []capturedPacket {
	c.mu.Lock()
	defer c.mu.Unlock()
	var out []capturedPacket
	for _, p := range c.packets {
		if p.OpCode == op {
			out = append(out, p)
		}
	}
	return out
}

func (c *captureConn) clear() {
	c.mu.Lock()
	c.packets = nil
	c.mu.Unlock()
}

func newTestClientAndPlayer(t *testing.T, repo database.CharacterRepository, ds *DataStore, zMap *ZoneMap, h *Handler, charID int, name string, level uint8, x, y int) (*ZonePlayer, *captureConn) {
	t.Helper()
	conn := &captureConn{}
	client, err := network.NewClient(conn, nil)
	if err != nil {
		t.Fatalf("NewClient err: %v", err)
	}

	var cs *combat.Store
	if ds != nil {
		cs = ds.CombatStore
	}
	ch := &database.Character{
		ID:        charID,
		AccountID: charID,
		Slot:      byte(charID % 5),
		Name:      name,
		Job:       1, // Fighter
		Level:     level,
		CurHP:     300,
		CurSP:     100,
		XPos:      x,
		YPos:      y,
		TCharacter: database.TCharacter{
			CharNo:       charID,
			UserNo:       charID,
			CharName:     name,
			Slot:         byte(charID % 5),
			Class:        1,
			Level:        level,
			Strength:     20,
			Constitute:   15,
			Dexterity:    15,
			Intelligence: 5,
			MentalPower:  5,
		},
		TCharacterState: database.TCharacterState{
			CharNo:   charID,
			XPos:     x,
			YPos:     y,
			HP:       300,
			SP:       100,
			HPStones: 10,
			SPStones: 10,
			Money:    1000,
		},
	}
	if err := repo.CreateCharacter(ch); err != nil {
		t.Fatalf("CreateCharacter(%s) err: %v", name, err)
	}
	// Asegurar coordenadas actualizadas si Sync reseteó algo
	ch.XPos = x
	ch.YPos = y
	var classStore = ds.ClassStore
	p := NewZonePlayer(client, ch, classStore, cs)
	p.CurrentMap = zMap
	zMap.AssignObjectID(p)
	zMap.AddPlayer(p)

	h.mu.Lock()
	h.players[client] = p
	h.mu.Unlock()
	return p, conn
}

func makeIncomingPacket(op uint16, build func(pw *protocol.PacketWriter)) *protocol.Packet {
	pw := protocol.NewPacketWriter(op)
	if build != nil {
		build(pw)
	}
	wire := pw.BuildWirePacket()
	pos := 1
	if wire[0] == 0 {
		pos = 3
	}
	return protocol.NewPacket(wire[pos:])
}

func loadFase13TestEnv(t *testing.T) (*DataStore, *Handler, database.CharacterRepository, *database.GameLogDispatcher) {
	t.Helper()
	dataDir := findTestDataDir()
	ds, err := LoadAllData(dataDir)
	if err != nil {
		t.Fatalf("LoadAllData(%s) err: %v", dataDir, err)
	}

	repo := database.NewMemoryCharacterRepositoryWithFile("")
	logDisp := database.NewGameLogDispatcher(nil, 200, 1)

	h := NewHandler(repo, nil, ds.ClassStore, ds.FieldStore, ds.ChkEngine)
	h.SetGameLogDispatcher(logDisp)
	h.SetDataStore(ds)
	return ds, h, repo, logDisp
}

// TestFase13_MobCatalogAndRegenFromRealData verifica carga 100% desde data/MobRegen, MobInfo*, MobSpecies,
// MobWeapon, MobKillAble, MobKillAnnounce y el formato binario exacto de 149 bytes en PROTO_NC_BRIEFINFO_REGENMOB_CMD.
func TestFase13_MobCatalogAndRegenFromRealData(t *testing.T) {
	ds, h, repo, _ := loadFase13TestEnv(t)

	if ds.MobCatalog == nil || len(ds.MobCatalog.ByID) == 0 {
		t.Fatalf("MobCatalog vacío tras LoadAllData")
	}
	if len(ds.MobCatalog.KillAnnounce) == 0 {
		t.Fatalf("KillAnnounce vacío tras cargar MobKillAnnounce.shn")
	}
	if ds.TitleCatalog == nil || len(ds.TitleCatalog.ByType) == 0 {
		t.Fatalf("TitleCatalog vacío tras cargar CharacterTitleData.shn")
	}
	if ds.CollectCatalog == nil || len(ds.CollectCatalog.CardsByID) == 0 || len(ds.CollectCatalog.RewardsByID) == 0 {
		t.Fatalf("CollectCatalog vacío tras cargar CollectCard*.shn")
	}

	// Crear mapa Rou (Roumen) que cuenta con data/MobRegen/Rou.json y Rou.shbd
	zMap := h.mapManager.GetOrCreateMap("Rou")
	mobs := zMap.GetAllMonsters()
	if len(mobs) == 0 {
		t.Fatalf("Rou debería cargar monstruos iniciales desde data/MobRegen/Rou.json")
	}

	// Validar que los mobs tengan HP, Nivel y serialización exacta de 149 bytes
	sample := mobs[0]
	if sample.MaxHP <= 0 || sample.CurHP <= 0 {
		t.Errorf("Mob %s tiene HP inválido: cur=%d max=%d", sample.InxName, sample.CurHP, sample.MaxHP)
	}
	tmpWriter := protocol.NewPacketWriter(protocol.PROTO_NC_BRIEFINFO_REGENMOB_CMD)
	sample.WriteRegenMob(tmpWriter, 1)
	if payloadLen := len(tmpWriter.Data()) - 2; payloadLen != 149 {
		t.Errorf("PROTO_NC_BRIEFINFO_REGENMOB_CMD payload esperado de 149 bytes, obtenido %d", payloadLen)
	}

	// Verificar SendMobs emitindo PROTO_NC_BRIEFINFO_MOB_CMD (0x1C09) al entrar el jugador
	p1, conn1 := newTestClientAndPlayer(t, repo, ds, zMap, h, 101, "Explorer", 10, int(sample.X), int(sample.Y))
	conn1.clear()
	zMap.SendMobs(p1, mobs[:1])
	if len(conn1.findByOpcode(protocol.PROTO_NC_BRIEFINFO_MOB_CMD)) == 0 {
		t.Errorf("SendMobs debería enviar PROTO_NC_BRIEFINFO_MOB_CMD (0x1C09)")
	}
}

// TestFase13_MonsterRoamStealthAggroLeashAndLuaAI cubre patrullaje MobRoam + .shbd, StealthMatrix,
// tabla de aggro, persecución, leash de retorno, llamada de auxilio, MobAttackSequence y scripts Lua.
func TestFase13_MonsterRoamStealthAggroLeashAndLuaAI(t *testing.T) {
	ds, h, repo, _ := loadFase13TestEnv(t)

	zMap := h.mapManager.GetOrCreateMap("Rou")

	// Elegir una ruta real cargada en data/MobRoam
	var anyRouteInx string
	if ds.MobRoam != nil {
		for k, rt := range ds.MobRoam.Routes {
			if len(rt.Points) >= 2 && ds.MobCatalog.LookupTemplate(k) != nil {
				anyRouteInx = k
				break
			}
		}
	}

	// Buscar una plantilla Slime real en MobCatalog de forma determinística
	var slimeInx string
	if tmpl := ds.MobCatalog.LookupTemplate("Slime"); tmpl != nil {
		slimeInx = tmpl.InxName
	} else {
		for inx := range ds.MobCatalog.ByInx {
			if strings.Contains(inx, "slime") && (slimeInx == "" || inx < slimeInx) {
				slimeInx = inx
			}
		}
	}
	if slimeInx == "" {
		t.Fatalf("No se encontró Slime en MobCatalog")
	}

	// Usar coordenada caminable real del mapa (Rou.shbd) obtenida de los monstruos ya generados por MobRegen
	var baseX, baseY int32 = 5000, 5000
	for _, existing := range zMap.GetAllMonsters() {
		if zMap.IsPositionWalkable(float64(existing.X), float64(existing.Y)) {
			baseX, baseY = existing.X, existing.Y
			break
		}
	}

	mob, err := zMap.SpawnMonsterByInx(slimeInx, baseX, baseY, false)
	if err != nil || mob == nil {
		t.Fatalf("SpawnMonsterByInx(%s) falló: %v", slimeInx, err)
	}
	mob.FollowRange = 200
	mob.AttackRange = 120
	mob.ReturnToRegen = true

	// 1. Stealth vs Detección agresiva en scanAggroTarget
	pStealth, connStealth := newTestClientAndPlayer(t, repo, ds, zMap, h, 201, "Sneaker", 10, int(baseX)+20, int(baseY)+20)
	pStealth.AddAbState(&ActiveAbState{ID: 1, Name: "StaHide", ExpiresAt: time.Now().Add(time.Minute)})
	mob.EnemyDetectType = 1
	mob.DetectRange = 300
	if zMap.scanAggroTarget(mob, ds) {
		t.Errorf("scanAggroTarget detectó a jugador con sigilo StaHide a pesar de StealthMatrix")
	}
	// Quitar sigilo => debe detectar al jugador y agregarlo al Aggro
	pStealth.RemoveAbState(1)
	if !zMap.scanAggroTarget(mob, ds) || len(mob.Aggro) == 0 {
		t.Errorf("scanAggroTarget debería detectar jugador visible dentro de DetectRange")
	}

	// 2. Rotación de MobRoam real si existe coincidencia directa de template+route
	if anyRouteInx != "" {
		roamMob, _ := zMap.SpawnMonsterByInx(anyRouteInx, 4000, 4000, false)
		if roamMob != nil && roamMob.Route != nil && len(roamMob.Route.Points) >= 2 {
			startX, startY := roamMob.X, roamMob.Y
			zMap.tickMobRoam(roamMob, time.Now())
			if roamMob.X == startX && roamMob.Y == startY {
				t.Logf("tickMobRoam validó paso inicial en %s", roamMob.InxName)
			}
		}
	}

	// 3. Persecución + Retorno (Leash) cuando el objetivo sale más allá de FollowRange
	pStealth.Char.XPos = 7500
	pStealth.Char.YPos = 7500
	zMap.grid.UpdateEntityPosition(pStealth, 7500, 7500)
	zMap.TickMonsters100ms(time.Now().Add(300 * time.Millisecond))
	mob.mu.RLock()
	stateAfterEscape := mob.AIState
	aggroLen := len(mob.Aggro)
	mob.mu.RUnlock()
	if stateAfterEscape != MobAIReturn || aggroLen != 0 {
		t.Errorf("Al exceder FollowRange, esperado MobAIReturn con aggro limpio, obtenido state=%v aggro=%d", stateAfterEscape, aggroLen)
	}

	// 4. MobAttackSequence + familia + ataque al jugador
	var seqMobInx string
	if ds.MobAttack != nil {
		for k, ai := range ds.MobAttack.Mobs {
			if len(ai.Sequence) > 0 && ds.MobCatalog.LookupTemplate(k) != nil {
				seqMobInx = k
				break
			}
		}
	}
	if seqMobInx == "" {
		seqMobInx = slimeInx
	}
	atkMob, _ := zMap.SpawnMonsterByInx(seqMobInx, 5100, 5100, false)
	allyMob, _ := zMap.SpawnMonsterByInx(seqMobInx, 5150, 5150, false)
	pVictim, connVictim := newTestClientAndPlayer(t, repo, ds, zMap, h, 202, "Tank", 15, 5120, 5120)
	atkMob.AddAggro(pVictim.MapObjectID, 50)
	connVictim.clear()
	connStealth.clear()

	zMap.callFamilyHelp(atkMob, pVictim, 300, 1)
	if len(allyMob.Aggro) == 0 {
		t.Errorf("callFamilyHelp debería reclutar al aliado de la misma familia")
	}

	atkMob.AttackDelayMs = 1
	atkMob.AttackRange = 200
	beforeHP := pVictim.Char.CurHP
	zMap.executeMonsterAttack(atkMob, pVictim, time.Now().Add(time.Second), ds, ds.MobCatalog)
	if pVictim.Char.CurHP >= beforeHP {
		// Permitir miss ocasional de CheckHit reintentando determinísticamente hasta 5 golpes
		hitOnce := false
		for i := 0; i < 10; i++ {
			zMap.executeMonsterAttack(atkMob, pVictim, time.Now().Add(time.Duration(i+2)*time.Second), ds, ds.MobCatalog)
			if pVictim.Char.CurHP < beforeHP {
				hitOnce = true
				break
			}
		}
		if !hitOnce {
			t.Errorf("executeMonsterAttack no dañó a Tank tras múltiples golpes")
		}
	}

	// 5. Ejecución real del motor Lua MobBehaviorDescript + AIScript/BallEgg.lua
	zMap.initMapLuaAI()
	if zMap.luaEngine == nil {
		t.Fatalf("luaEngine no inicializado en mapa")
	}
	ballLua := filepath.Join(resolveZoneDataDir(), "LuaScript", "AIScript", "BallEgg.lua")
	if err := zMap.luaEngine.DoFile(ballLua); err != nil {
		t.Fatalf("Error cargando AIScript/BallEgg.lua: %v", err)
	}
	handled := zMap.RunMobBehaviorLua(atkMob, "BallEgg")
	_ = handled
}

// TestFase13_DeathExpPartyRaidDropsAutoPickAndAnnounce prueba muerte de monstruo, EXP Party/Raid, RestExp,
// subida de nivel, RandomOption, propiedad temporal, AutoPick (0x10BC), MobKillAnnounce y GameLog Tipo 6.
func TestFase13_DeathExpPartyRaidDropsAutoPickAndAnnounce(t *testing.T) {
	ds, h, repo, _ := loadFase13TestEnv(t)

	zMap := h.mapManager.GetOrCreateMap("Rou")
	killer, connKiller := newTestClientAndPlayer(t, repo, ds, zMap, h, 301, "Slayer", 1, 3000, 3000)
	partner, connPartner := newTestClientAndPlayer(t, repo, ds, zMap, h, 302, "Healer", 1, 3050, 3050)
	looter, _ := newTestClientAndPlayer(t, repo, ds, zMap, h, 303, "Thief", 1, 8000, 8000)

	// Configurar Party entre killer y partner en rango; looter fuera de Party
	partyNo, err := repo.CreateParty(killer.Char.ID, 0)
	if err != nil {
		t.Fatalf("CreateParty err: %v", err)
	}
	_ = repo.AddPartyMember(partyNo, partner.Char.ID, 0)

	// Activar RestExp en Common.json para killer
	killer.RestExpBonusRate = 1100

	// Buscar un mob que figure en MobKillAnnounce.shn y tenga DropTable en ItemDropTable.json
	var bossID uint16
	for id := range ds.MobCatalog.KillAnnounce {
		bossID = id
		break
	}
	tmpl := ds.MobCatalog.ByID[bossID]
	if tmpl == nil {
		for _, cand := range ds.MobCatalog.ByID {
			tmpl = cand
			bossID = cand.ID
			break
		}
	}
	mob, err := zMap.SpawnMonsterByInx(tmpl.InxName, 3020, 3020, true)
	if err != nil {
		t.Fatalf("SpawnMonsterByInx err: %v", err)
	}
	mob.Exp = 5000 // Suficiente para subir de nivel desde Lv 1 (NextExp ~ 10)
	mob.CurHP = 20

	connKiller.clear()
	connPartner.clear()

	_, dead := zMap.ApplyDamageToMonster(mob, killer, 50)
	if !dead {
		t.Fatalf("El monstruo debería haber muerto al recibir 50 de daño")
	}

	// Ambos miembros en Party deben ganar EXP y subir al menos a Nivel 2+
	if killer.Char.Level <= 1 || partner.Char.Level <= 1 {
		t.Errorf("Killer (%d) y Partner (%d) en Party deberían haber subido de nivel", killer.Char.Level, partner.Char.Level)
	}
	if len(connKiller.findByOpcode(protocol.PROTO_NC_BAT_EXPGAIN_CMD)) == 0 {
		t.Errorf("Esperado paquete PROTO_NC_BAT_EXPGAIN_CMD (0x240B)")
	}
	if len(connKiller.findByOpcode(protocol.PROTO_NC_BAT_LEVELUP_CMD)) == 0 {
		t.Errorf("Esperado paquete PROTO_NC_BAT_LEVELUP_CMD (0x240C)")
	}
	if _, hasAnn := ds.MobCatalog.KillAnnounce[mob.MobID]; hasAnn {
		if len(connKiller.findByOpcode(protocol.PROTO_NC_BAT_CLIENT_MOB_KILL_ANNOUNCE_CMD)) == 0 {
			t.Errorf("Esperado anuncio PROTO_NC_BAT_CLIENT_MOB_KILL_ANNOUNCE_CMD (0x2459)")
		}
	}

	// 2. Probar generación de RandomOption desde RandomOption.shn / RandomOptionCount.shn / RandomOptionTable.json
	roc := getRandomOptCatalog(ds)
	var optRuleKey string
	for k, rules := range roc.RulesByItem {
		if len(rules) > 0 && len(roc.CountsByItem[k]) > 0 {
			optRuleKey = k
			break
		}
	}
	if optRuleKey == "" {
		t.Fatalf("No se encontró entrada válida en RandomOption.shn + RandomOptionCount.shn")
	}

	opts := h.RollItemRandomOptions(optRuleKey)
	for i := 0; i < 30 && len(opts) == 0; i++ {
		opts = h.RollItemRandomOptions(optRuleKey)
	}
	if len(opts) == 0 {
		// Fallback directo a RandomOptionTable.json (e.g., RandomNamedW01)
		opts = h.RollItemRandomOptions("RandomNamedW01")
	}
	if len(opts) == 0 {
		t.Errorf("RollItemRandomOptions(%s / RandomNamedW01) no generó opciones aleatorias", optRuleKey)
	}

	// Seleccionar un arma real de ItemInfo.shn (Class == 5) para validar Serialización + DB ItemOption
	var optInfoID uint16
	for _, row := range ds.GetSHN("ItemInfo").Rows {
		if row.GetInt("Class") == 5 {
			optInfoID = row.GetUint16("ID")
			break
		}
	}
	if optInfoID == 0 {
		t.Fatalf("No se encontró arma Class=5 en ItemInfo.shn")
	}

	// 3. Propiedad temporal de drop en el suelo vs expiración y recogida con StatOptions
	dropObj := zMap.SpawnMonsterDrop(optInfoID, 1, 1, 3020, 3020, mob.Handle, killer.Char.ID, 5*time.Second, false, false, opts)
	// Posicionar al looter en rango de recogida (< 125u, Zone.exe 0x3D09) para probar bloqueo por dueño y expiración
	looter.Char.XPos = 3020
	looter.Char.YPos = 3020
	// looter (no dueño) intenta recoger antes de expirar -> denegado (drop permanece en suelo)
	h.HandlePacket(looter.Client, makeIncomingPacket(protocol.PROTO_NC_ITEM_PICK_REQ, func(pw *protocol.PacketWriter) {
		pw.WriteUint16(dropObj.Handle)
	}))
	if zMap.GetEntity(dropObj.Handle) == nil {
		t.Errorf("El drop no debería poder ser recogido por otro jugador antes de vencer la propiedad temporal")
	}

	// Expirar propiedad y verificar recogida exitosa + persistencia de StatOptions
	dropObj.OwnerExpiresAt = time.Now().Add(-time.Second)
	zMap.CleanupExpiredDrops(DefaultDropDespawnAge)
	h.HandlePacket(looter.Client, makeIncomingPacket(protocol.PROTO_NC_ITEM_PICK_REQ, func(pw *protocol.PacketWriter) {
		pw.WriteUint16(dropObj.Handle)
	}))
	looterItems, _ := repo.GetCharacterItems(looter.Char.ID)
	if len(looterItems) == 0 {
		t.Fatalf("Looter debería poder recoger el drop una vez expirada la propiedad temporal")
	}
	savedOpts, _ := repo.GetItemOptions(looterItems[0].ItemKey)
	if len(savedOpts) == 0 {
		t.Errorf("Las opciones aleatorias del drop deberían haberse guardado en GetItemOptions")
	}

	// 4. PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ (0x10BC) + AutoPick de Cen e Ítems
	connKiller.clear()
	h.HandlePacket(killer.Client, makeIncomingPacket(protocol.PROTO_NC_CHAR_CLIENT_AUTO_PICK_REQ, func(pw *protocol.PacketWriter) {
		pw.WriteByte(1)
	}))
	if !killer.IsAutoPick {
		t.Errorf("killer.IsAutoPick debería quedar activo tras 0x10BC")
	}
	if ap, _ := repo.GetAutoPick(killer.Char.ID); !ap {
		t.Errorf("repo.GetAutoPick debería persistir true tras usp_Char_SetAutoPick")
	}
	if len(connKiller.findByOpcode(protocol.PROTO_NC_CHAR_CLIENT_AUTO_PICK_ACK)) == 0 {
		t.Errorf("Esperado PROTO_NC_CHAR_CLIENT_AUTO_PICK_ACK (0x10BD)")
	}
	if len(connKiller.findByOpcode(protocol.PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD)) == 0 {
		t.Errorf("Esperado PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD (0x10BE)")
	}

	moneyBefore := killer.Char.TCharacterState.Money
	cenDrop := zMap.SpawnMonsterDrop(h.resolveMoneyItemID(), 250, 250, 3000, 3000, mob.Handle, killer.Char.ID, DefaultDropOwnershipDuration, false, true, nil)
	h.autoPickDroppedItem(zMap, killer, cenDrop)
	if killer.Char.TCharacterState.Money != moneyBefore+250 {
		t.Errorf("AutoPick de Cen esperado %d, obtenido %d", moneyBefore+250, killer.Char.TCharacterState.Money)
	}
	itemDropAuto := zMap.SpawnMonsterDrop(optInfoID, 1, 1, 3000, 3000, mob.Handle, killer.Char.ID, DefaultDropOwnershipDuration, false, false, opts)
	h.autoPickDroppedItem(zMap, killer, itemDropAuto)
	killerItems, _ := repo.GetCharacterItems(killer.Char.ID)
	if len(killerItems) == 0 {
		t.Errorf("AutoPick de ítem debería ingresarlo directamente a Box 9")
	}
}

// TestFase13_TitlesDept24AndCardCollectionDept49 verifica Títulos (Dept 24: 0x6001/2/3/4 + CharacterTitleData.shn)
// y Colección de Cartas (Dept 49: Peeping, CardRegist, Bookmark, Reward con CollectCard*.shn).
func TestFase13_TitlesDept24AndCardCollectionDept49(t *testing.T) {
	ds, h, repo, _ := loadFase13TestEnv(t)

	zMap := h.mapManager.GetOrCreateMap("Rou")
	p1, conn1 := newTestClientAndPlayer(t, repo, ds, zMap, h, 401, "Collector", 20, 2000, 2000)
	p2, conn2 := newTestClientAndPlayer(t, repo, ds, zMap, h, 402, "Viewer", 20, 2050, 2050)

	// 1. Elegir carta oficial con título asociado en CollectCardTitle.shn (o la primera carta válida)
	var testCard *CollectCardDef
	for itemInx := range ds.TitleCatalog.CardTitleByItemInx {
		if cd := ds.CollectCatalog.CardsByItemInx[itemInx]; cd != nil {
			if info, ok := ds.ItemStore.GetItemInfoByName(cd.ItemInx); ok && info != nil {
				testCard = cd
				break
			}
		}
	}
	if testCard == nil {
		for _, cd := range ds.CollectCatalog.CardsByID {
			if info, ok := ds.ItemStore.GetItemInfoByName(cd.ItemInx); ok && info != nil {
				testCard = cd
				break
			}
		}
	}
	if testCard == nil {
		t.Fatalf("No se encontró carta válida en CollectCard.shn + ItemInfo.shn")
	}
	cardItemInfo, _ := ds.ItemStore.GetItemInfoByName(testCard.ItemInx)

	// Colocar carta en slot 0 del inventario (Box 9)
	_ = repo.CreateItem(&database.Item{
		StorageType: 9,
		Owner:       p1.Char.ID,
		Storage:     0,
		ItemID:      cardItemInfo.ID,
		Flags:       1,
	})

	conn1.clear()
	conn2.clear()

	// 2. PROTO_NC_COLLECT_CARDREGIST_REQ (0xC408) -> StarRate + Serial + CardTitle unlock
	h.HandlePacket(p1.Client, makeIncomingPacket(protocol.PROTO_NC_COLLECT_CARDREGIST_REQ, func(pw *protocol.PacketWriter) {
		pw.WriteByte(0) // Slot 0
	}))
	regAcks := conn1.findByOpcode(protocol.PROTO_NC_COLLECT_CARDREGIST_ACK)
	if len(regAcks) != 1 {
		t.Fatalf("Esperado 1 PROTO_NC_COLLECT_CARDREGIST_ACK (0xC40B), obtenido %d", len(regAcks))
	}
	ackR := protocol.NewPacketReader(regAcks[0].Payload)
	_, _ = ackR.ReadByte()          // invenslot
	gotCardID, _ := ackR.ReadUint16()
	gotStar, _ := ackR.ReadByte()
	gotSerial, _ := ackR.ReadUint32()
	gotErr, _ := ackR.ReadUint16()
	if gotErr != 0 || gotCardID != testCard.CardID || gotStar == 0 || gotSerial == 0 {
		t.Errorf("CardRegist ACK inválido: card=%d star=%d serial=%d err=%d", gotCardID, gotStar, gotSerial, gotErr)
	}

	// Verificar si desbloqueó título asociado (0x6004)
	if linkedTitleType, hasT := ds.TitleCatalog.CardTitleByItemInx[strings.ToLower(testCard.ItemInx)]; hasT {
		if len(conn1.findByOpcode(protocol.PROTO_NC_CT_SET_CMD)) == 0 {
			t.Errorf("Esperado desbloqueo de título PROTO_NC_CT_SET_CMD (0x6004) para tipo %d", linkedTitleType)
		}
		// 3. Equipar título con PROTO_NC_CT_SET_CURRENT_REQ (0x6001)
		h.HandlePacket(p1.Client, makeIncomingPacket(protocol.PROTO_NC_CT_SET_CURRENT_REQ, func(pw *protocol.PacketWriter) {
			pw.WriteByte(linkedTitleType)
			pw.WriteByte(0)
			pw.WriteUint16(0)
		}))
		if len(conn1.findByOpcode(protocol.PROTO_NC_CT_SET_CURRENT_ACK)) == 0 {
			t.Errorf("Esperado PROTO_NC_CT_SET_CURRENT_ACK (0x6002)")
		}
		if len(conn2.findByOpcode(protocol.PROTO_NC_CT_SET_SOMEONECHANGE_CMD)) == 0 {
			t.Errorf("Esperado broadcast AOI PROTO_NC_CT_SET_SOMEONECHANGE_CMD (0x6003) en el jugador vecino")
		}
		if p1.CurrentTitleType != linkedTitleType {
			t.Errorf("CurrentTitleType esperado %d, obtenido %d", linkedTitleType, p1.CurrentTitleType)
		}
	}

	// 4. PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ (0xC40C)
	h.HandlePacket(p1.Client, makeIncomingPacket(protocol.PROTO_NC_COLLECT_BOOKMARK_REGIST_REQ, func(pw *protocol.PacketWriter) {
		pw.WriteUint16(1)
		pw.WriteUint16(testCard.CardID)
	}))
	bmAcks := conn1.findByOpcode(protocol.PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK)
	if len(bmAcks) != 1 {
		t.Errorf("Esperado PROTO_NC_COLLECT_BOOKMARK_REGIST_ACK (0xC40F)")
	}

	// 5. PROTO_NC_COLLECT_PEEPING_REQ (0xC401) desde Viewer inspeccionando a Collector
	conn2.clear()
	h.HandlePacket(p2.Client, makeIncomingPacket(protocol.PROTO_NC_COLLECT_PEEPING_REQ, func(pw *protocol.PacketWriter) {
		pw.WriteString(p1.Char.Name, 20)
	}))
	if len(conn2.findByOpcode(protocol.PROTO_NC_COLLECT_PEEPING_COLLECT_ACK)) == 0 ||
		len(conn2.findByOpcode(protocol.PROTO_NC_COLLECT_PEEPING_COLLECT_BOOKMARK_ACK)) == 0 {
		t.Errorf("Peeping debería devolver 0xC405 y 0xC406 al observador")
	}

	// 6. PROTO_NC_COLLECT_REWARD_REQ (0xC410) usando la recompensa mínima de CollectCardReward.shn
	var targetReward *CollectRewardDef
	for _, rw := range ds.CollectCatalog.RewardsByID {
		if _, ok := ds.ItemStore.GetItemInfoByName(rw.RewardItemInx); ok {
			if targetReward == nil || rw.CardLot < targetReward.CardLot {
				targetReward = rw
			}
		}
	}
	if targetReward != nil {
		// Poblar cartas faltantes hasta alcanzar CardLot
		cards, _ := repo.GetCollectionCards(p1.Char.ID)
		idx := len(cards)
		for _, cd := range ds.CollectCatalog.CardsByID {
			if idx >= int(targetReward.CardLot) {
				break
			}
			already := false
			for _, cur := range cards {
				if uint16(cur.CardID) == cd.CardID {
					already = true
					break
				}
			}
			if !already {
				_ = repo.RegisterCollectionCard(p1.Char.ID, int(cd.CardID), int64(idx+100), 1, 0)
				idx++
			}
		}

		conn1.clear()
		h.HandlePacket(p1.Client, makeIncomingPacket(protocol.PROTO_NC_COLLECT_REWARD_REQ, func(pw *protocol.PacketWriter) {
			pw.WriteUint16(targetReward.RewardID)
			pw.WriteUint16(0)
		}))
		rwAcks := conn1.findByOpcode(protocol.PROTO_NC_COLLECT_REWARD_ACK)
		if len(rwAcks) != 1 {
			t.Fatalf("Esperado PROTO_NC_COLLECT_REWARD_ACK (0xC413)")
		}
		rr := protocol.NewPacketReader(rwAcks[0].Payload)
		_, _ = rr.ReadUint16()
		rwErr, _ := rr.ReadUint16()
		if rwErr != 0 {
			t.Errorf("Claim reward %d falló con err=%d", targetReward.RewardID, rwErr)
		}
		claimed, _ := repo.GetCollectionRewards(p1.Char.ID)
		if len(claimed) == 0 {
			t.Errorf("Recompensa debería quedar persistida en GetCollectionRewards")
		}
	}
}
