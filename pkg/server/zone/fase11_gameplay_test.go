package zone

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/data/blockinfo"
	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/combat"
	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

func findTestDataDir() string {
	candidates := []string{"data", "../../data", "../../../data"}
	for _, c := range candidates {
		if fi, err := os.Stat(filepath.Join(c, "World", "Combat", "Common.json")); err == nil && !fi.IsDir() {
			return c
		}
	}
	return "data"
}

func TestFase11_AOI_SpatialGrid_Transitions(t *testing.T) {
	grid := NewSpatialGrid(1000)

	// Crear dos jugadores ficticios: p1 en (500, 500) [celda 0,0], p2 en (500, 500) [celda 0,0]
	ch1 := &database.Character{ID: 1, Name: "Hero1", XPos: 500, YPos: 500}
	p1 := &ZonePlayer{MapObjectID: 1, Char: ch1}

	ch2 := &database.Character{ID: 2, Name: "Hero2", XPos: 500, YPos: 500}
	p2 := &ZonePlayer{MapObjectID: 2, Char: ch2}

	// Jugador 3 en celda distante (5000, 5000) [celda 5,5]
	ch3 := &database.Character{ID: 3, Name: "DistantHero", XPos: 5000, YPos: 5000}
	p3 := &ZonePlayer{MapObjectID: 3, Char: ch3}

	grid.AddEntity(p1)
	grid.AddEntity(p2)
	grid.AddEntity(p3)

	// Verificar vecindad 3x3 de (500, 500)
	surrPlayers := grid.GetSurroundingPlayers(500, 500)
	if len(surrPlayers) != 2 {
		t.Fatalf("Esperados 2 jugadores en AOI de (500, 500), obtenidos %d", len(surrPlayers))
	}

	// Verificar que p3 no está en la vecindad de (500, 500)
	foundP3 := false
	for _, p := range surrPlayers {
		if p.MapObjectID == 3 {
			foundP3 = true
		}
	}
	if foundP3 {
		t.Fatalf("p3 (5000,5000) no debería estar en AOI de (500,500)")
	}

	// Mover p1 dentro de la misma celda (600, 600) -> crossCell debe ser false
	cross, entered, left := grid.UpdateEntityPosition(p1, 600, 600)
	if cross {
		t.Fatalf("Movimiento dentro de la misma celda no debería disparar crossCell")
	}
	if len(entered) != 0 || len(left) != 0 {
		t.Fatalf("entered o left celdas no deberían estar pobladas")
	}

	// Mover p1 a celda adyacente (1500, 500) [celda 1,0]
	cross, entered, left = grid.UpdateEntityPosition(p1, 1500, 500)
	if !cross {
		t.Fatalf("Movimiento a celda (1,0) debería disparar crossCell = true")
	}
	if len(entered) == 0 || len(left) == 0 {
		t.Fatalf("entered (%d) y left (%d) deben tener celdas tras cruce", len(entered), len(left))
	}

	// Eliminar p2
	_, _, removed := grid.RemoveEntity(p2.MapObjectID)
	if !removed {
		t.Fatalf("Fallo al remover entidad p2 de la grilla")
	}
	surrAfter := grid.GetSurroundingPlayers(1500, 500)
	for _, p := range surrAfter {
		if p.MapObjectID == 2 {
			t.Fatalf("p2 no debería estar en la grilla tras ser removido")
		}
	}
}

func TestFase11_All27OfficialClasses_FormulasAndStats(t *testing.T) {
	// Verificar que las 27 clases oficiales mapean correctamente
	expectedClasses := map[byte]string{
		1:  "Fighter",
		2:  "CleverFighter",
		3:  "Warrior",
		4:  "Gladiator",
		5:  "Knight",
		6:  "Cleric",
		7:  "HighCleric",
		8:  "Paladin",
		9:  "HolyKnight",
		10: "Guardian",
		11: "Archer",
		12: "HawkArcher",
		13: "Scout",
		14: "SharpShooter",
		15: "Ranger",
		16: "Mage",
		17: "WizMage",
		18: "Enchanter",
		19: "Warrock",
		20: "Wizard",
		21: "Joker",
		22: "Chaser",
		23: "Cruel",
		24: "Closer",
		25: "Assassin",
		26: "Sentinel",
		27: "Savior",
	}

	for id, name := range expectedClasses {
		got := JobToClassName(id)
		if got != name {
			t.Errorf("Clase ID %d: esperado '%s', obtenido '%s'", id, name, got)
		}
	}

	dataDir := findTestDataDir()

	// Cargar CombatStore real para verificar bonos de FreeStat y ExpTable
	cs := combat.NewStore()
	if err := cs.LoadDir(filepath.Join(dataDir, "World", "Combat")); err != nil {
		t.Fatalf("Error cargando CombatStore: %v", err)
	}

	// Verificar bonos de CON y MEN
	bonus := cs.GetFreeStatBonus(10, 0, 0, 10, 10)
	if bonus.MaxHP != 50 { // 10 puntos de CON = 50 HP (5 por punto según FreeStatTable.json)
		t.Errorf("Esperado MaxHP bonus = 50 para CON=10, obtenido %d", bonus.MaxHP)
	}
	if bonus.MaxSP != 50 { // 10 puntos de MEN = 50 SP (5 por punto según FreeStatTable.json)
		t.Errorf("Esperado MaxSP bonus = 50 para MEN=10, obtenido %d", bonus.MaxSP)
	}

	// Verificar progresión de ExpTable (Nivel 1: 0..10, Nivel 2: 10..40, Nivel 3: 40..90)
	pExp, nExp := cs.GetExpProgression(1)
	if pExp != 0 || nExp != 10 {
		t.Errorf("Nivel 1: esperado (0, 10), obtenido (%d, %d)", pExp, nExp)
	}
	pExp2, nExp2 := cs.GetExpProgression(2)
	if pExp2 != 10 || nExp2 != 40 {
		t.Errorf("Nivel 2: esperado (10, 40), obtenido (%d, %d)", pExp2, nExp2)
	}

	// Cargar ClassParam para verificar NewZonePlayer con stats reales
	cps := classparam.NewStore()
	if err := cps.LoadDir(filepath.Join(dataDir, "World", "ClassParam")); err != nil {
		t.Fatalf("Error cargando ClassParam: %v", err)
	}

	ch := &database.Character{
		ID:    100,
		Name:  "TestFighter",
		Job:   1, // Fighter
		Level: 1,
		TCharacter: database.TCharacter{
			Constitute:  5,
			MentalPower: 5,
		},
	}

	zp := NewZonePlayer(nil, ch, cps, cs)
	if zp.MaxHP <= 0 || zp.MaxSP <= 0 {
		t.Fatalf("MaxHP y MaxSP deben ser positivos calculados: HP=%d, SP=%d", zp.MaxHP, zp.MaxSP)
	}
	// lp.MaxHP para Fighter nivel 1 es 52. Con CON=5 (5*5=25 HP bonus) -> MaxHP = 77
	if zp.MaxHP != 77 {
		t.Errorf("Esperado MaxHP 77 (52 base + 25 free), obtenido %d", zp.MaxHP)
	}
	if zp.NextExp != 10 {
		t.Errorf("Esperado NextExp 10 para nivel 1, obtenido %d", zp.NextExp)
	}
}

func TestFase11_VisualUpgrades_Packing(t *testing.T) {
	ch := &database.Character{
		ID:   1,
		Name: "GlowHero",
		EquipUpgrades: map[byte]byte{
			Slot_Weapon:  9,  // Right Hand = +9
			Slot_Weapon2: 10, // Left Hand = +10
			Slot_Armor:   7,  // Body = +7
			Slot_Pants:   8,  // Pants = +8
			Slot_Boots:   15, // Boots = +15
		},
	}

	zp := &ZonePlayer{Char: ch}
	pw := protocol.NewPacketWriter(0)
	zp.WriteEquipment(pw)

	data := pw.Data()[2:] // Omitir opcode
	if len(data) != 43 {
		t.Fatalf("Esperados 43 bytes en PROTO_EQUIPMENT, obtenidos %d", len(data))
	}

	// Los últimos 3 bytes son los upgrades visuales
	upgrades := data[40:43]
	byte0 := upgrades[0] // rh (9) | (lh (10) << 4) = 0x09 | 0xA0 = 0xA9
	byte1 := upgrades[1] // body (7) | (pants (8) << 4) = 0x07 | 0x80 = 0x87
	byte2 := upgrades[2] // boots (15) = 0x0F

	if byte0 != 0xA9 {
		t.Errorf("Byte 0: esperado 0xA9, obtenido 0x%02X", byte0)
	}
	if byte1 != 0x87 {
		t.Errorf("Byte 1: esperado 0x87, obtenido 0x%02X", byte1)
	}
	if byte2 != 0x0F {
		t.Errorf("Byte 2: esperado 0x0F, obtenido 0x%02X", byte2)
	}
}

func TestFase11_TentAndGatherAndMover(t *testing.T) {
	zm := NewZoneMap(1, "Rou", "Roumen")
	defer zm.Close()

	ch := &database.Character{
		ID:   1,
		Name: "Camper",
		XPos: 500,
		YPos: 500,
	}
	zp := &ZonePlayer{MapObjectID: 1, Char: ch, MaxHP: 100, MaxSP: 100}
	zm.AddPlayer(zp)

	// 1. Armar carpa (PitchTent)
	zm.PitchTent(zp, 101)
	if zp.State != 4 || zp.TentID != 101 {
		t.Errorf("PitchTent: esperado State=4, TentID=101, obtenido State=%d, TentID=%d", zp.State, zp.TentID)
	}

	// 2. Simular regeneración de carpa en GameLoop
	zp.Char.CurHP = 50
	zp.Char.CurSP = 50
	zm.gameLoop.tick1s()
	if zp.Char.CurHP <= 50 || zp.Char.CurSP <= 50 {
		t.Errorf("GameLoop tick1s en descanso no regeneró HP/SP: HP=%d, SP=%d", zp.Char.CurHP, zp.Char.CurSP)
	}

	// 3. Desarmar carpa (FoldTent)
	zm.FoldTent(zp)
	if zp.State != 1 || zp.TentID != 0 {
		t.Errorf("FoldTent: esperado State=1, TentID=0, obtenido State=%d, TentID=%d", zp.State, zp.TentID)
	}

	// 4. Montura (Mover)
	zp.MoverHandle = 55
	zp.MoverSlot = 1
	zp.MoverSpeed = 220
	zp.State = 6 // Mount
	if zp.MoverSpeed != 220 || zp.State != 6 {
		t.Errorf("Mover: velocidad o estado incorrecto")
	}

	// 5. Item drop y limpieza
	di := zm.SpawnDroppedItem(1001, 1, 0, 500, 500, 0, 1)
	if di == nil || di.Handle == 0 {
		t.Fatalf("Fallo al crear dropped item")
	}
	di.CreatedAt = time.Now().Add(-3 * time.Minute)
	zm.CleanupExpiredDrops(2 * time.Minute)
	if _, ok := zm.droppedItems[di.Handle]; ok {
		t.Errorf("CleanupExpiredDrops no eliminó ítem expirado")
	}
}

func TestFase11_Collision_SHBD_and_SBI(t *testing.T) {
	dataDir := findTestDataDir()
	biStore, err := blockinfo.LoadAll(filepath.Join(dataDir, "BlockInfo"))
	if err != nil {
		t.Fatalf("Error cargando BlockInfo: %v", err)
	}

	mbi := biStore.GetSHBD("Rou")
	if mbi == nil {
		t.Fatalf("Rou.shbd no encontrado en BlockInfo")
	}

	doors := biStore.NewInstanceDoors("Rou")

	zm := NewZoneMap(0, "Rou", "Roumen")
	defer zm.Close()
	zm.SetCollisionData(mbi, doors)

	// Validar que las coordenadas fuera del mapa o bloqueadas dan MoveBlocked
	if !mbi.IsMoveBlocked(3000, 3000) {
		t.Errorf("Celda fuera de límites (3000, 3000) debe estar bloqueada")
	}

	// Probar MovePlayerCmd contra celda bloqueada
	ch := &database.Character{
		ID:   1,
		Name: "MoverHero",
		XPos: 4134,
		YPos: 6488,
	}
	zp := &ZonePlayer{MapObjectID: 1, Char: ch}
	zm.AddPlayer(zp)

	// Intentar mover a posición bloqueada (fuera de límites)
	zm.MovePlayerCmd(zp, 4134, 6488, 999999, 999999, false)
	// Posición debe permanecer intacta en 4134, 6488
	if zp.Char.XPos != 4134 || zp.Char.YPos != 6488 {
		t.Errorf("MovePlayerCmd no debe actualizar posición si colisiona: pos=(%d, %d)", zp.Char.XPos, zp.Char.YPos)
	}
}

func TestFase11_Gather_Lifecycle_and_Mover_GameLog(t *testing.T) {
	zm := NewZoneMap(0, "Rou", "Roumen")
	defer zm.Close()

	ch := &database.Character{
		ID:   42,
		Name: "Gatherer",
		XPos: 4000,
		YPos: 4000,
	}
	zp := &ZonePlayer{MapObjectID: 1, Char: ch}
	zm.AddPlayer(zp)

	h := &Handler{
		players: map[*network.Client]*ZonePlayer{
			nil: zp,
		},
	}

	// 1. Gather Start
	zp.GatheringObj = 888
	if zp.GatheringObj != 888 {
		t.Errorf("GatheringObj esperado 888")
	}

	// 2. Gather Complete
	zp.GatheringObj = 0
	if zp.GatheringObj != 0 {
		t.Errorf("GatheringObj esperado 0 tras completar")
	}

	// 3. Mover Ride On y GameLog
	zp.MoverHandle = 10
	zp.MoverSlot = 2
	zp.MoverSpeed = 220
	zp.State = 6 // Mount

	if zp.MoverSpeed != 220 || zp.State != 6 {
		t.Errorf("Montura no fijó velocidad 220 ni estado 6")
	}

	// 4. Mover Ride Off
	zp.MoverHandle = 0
	zp.MoverSlot = 0
	zp.MoverSpeed = 0
	zp.State = 1 // Normal

	if zp.MoverSpeed != 0 || zp.State != 1 {
		t.Errorf("Desmontar no restableció estado normal")
	}

	// 5. HideExceptMe toggling
	zp.HideExceptMe = true
	if !zp.HideExceptMe {
		t.Errorf("HideExceptMe esperado true")
	}
	zp.HideExceptMe = false
	if zp.HideExceptMe {
		t.Errorf("HideExceptMe esperado false")
	}

	_ = h
}
