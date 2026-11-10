package zone

import (
	"sync"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
)

// CellCoord representa las coordenadas discretas de un sector espacial en la grilla AOI.
type CellCoord struct {
	X int
	Y int
}

const (
	// DefaultCellSize define el tamaño del sector en unidades de mundo oficiales (1000 unidades, ~80 celdas tile).
	DefaultCellSize = 1000
)

// GridEntityType discrimina el tipo de entidad alojada en la grilla espacial.
type GridEntityType byte

const (
	GridEntityPlayer GridEntityType = 1
	GridEntityNPC    GridEntityType = 2
	GridEntityItem   GridEntityType = 3
	GridEntityTent   GridEntityType = 4
	GridEntityMover  GridEntityType = 5
)

// GridEntity es la interfaz de abstracción para todas las entidades físicas en el mapa.
type GridEntity interface {
	GetHandle() uint16
	GetPosition() (int, int)
	GetEntityType() GridEntityType
	GetClient() *network.Client
}

// SpatialGrid gestiona la partición espacial en sectores y suscripciones dinámicas AOI.
type SpatialGrid struct {
	mu         sync.RWMutex
	cellSize   int
	cells      map[CellCoord]map[uint16]GridEntity
	entityCell map[uint16]CellCoord
}

// NewSpatialGrid inicializa una nueva grilla espacial con el tamaño de celda especificado.
func NewSpatialGrid(cellSize int) *SpatialGrid {
	if cellSize <= 0 {
		cellSize = DefaultCellSize
	}
	return &SpatialGrid{
		cellSize:   cellSize,
		cells:      make(map[CellCoord]map[uint16]GridEntity),
		entityCell: make(map[uint16]CellCoord),
	}
}

// GetCellCoord calcula la coordenada de cuadrante (CellCoord) para una posición de mundo (x, y).
func (g *SpatialGrid) GetCellCoord(x, y int) CellCoord {
	cx := x / g.cellSize
	cy := y / g.cellSize
	if x < 0 && x%g.cellSize != 0 {
		cx--
	}
	if y < 0 && y%g.cellSize != 0 {
		cy--
	}
	return CellCoord{X: cx, Y: cy}
}

// AddEntity registra una entidad en la grilla espacial.
func (g *SpatialGrid) AddEntity(e GridEntity) {
	if e == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	h := e.GetHandle()
	x, y := e.GetPosition()
	coord := g.GetCellCoord(x, y)

	if _, ok := g.cells[coord]; !ok {
		g.cells[coord] = make(map[uint16]GridEntity)
	}
	g.cells[coord][h] = e
	g.entityCell[h] = coord
}

// RemoveEntity elimina una entidad de la grilla espacial.
func (g *SpatialGrid) RemoveEntity(handle uint16) (GridEntity, CellCoord, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	coord, ok := g.entityCell[handle]
	if !ok {
		return nil, CellCoord{}, false
	}
	delete(g.entityCell, handle)

	var ent GridEntity
	if cell, exists := g.cells[coord]; exists {
		ent = cell[handle]
		delete(cell, handle)
		if len(cell) == 0 {
			delete(g.cells, coord)
		}
	}
	return ent, coord, true
}

// UpdateEntityPosition actualiza la posición de la entidad y retorna las celdas nuevas ingresadas y abandonadas.
func (g *SpatialGrid) UpdateEntityPosition(e GridEntity, newX, newY int) (crossCell bool, enteredCells, leftCells []CellCoord) {
	if e == nil {
		return false, nil, nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	h := e.GetHandle()
	oldCoord, ok := g.entityCell[h]
	newCoord := g.GetCellCoord(newX, newY)

	if !ok {
		if _, exists := g.cells[newCoord]; !exists {
			g.cells[newCoord] = make(map[uint16]GridEntity)
		}
		g.cells[newCoord][h] = e
		g.entityCell[h] = newCoord
		return true, g.getSurroundingCells(newCoord), nil
	}

	if oldCoord == newCoord {
		return false, nil, nil
	}

	// Mover entidad de celda
	if cell, exists := g.cells[oldCoord]; exists {
		delete(cell, h)
		if len(cell) == 0 {
			delete(g.cells, oldCoord)
		}
	}

	if _, exists := g.cells[newCoord]; !exists {
		g.cells[newCoord] = make(map[uint16]GridEntity)
	}
	g.cells[newCoord][h] = e
	g.entityCell[h] = newCoord

	oldSurr := g.getSurroundingCells(oldCoord)
	newSurr := g.getSurroundingCells(newCoord)

	oldSet := make(map[CellCoord]bool, len(oldSurr))
	for _, c := range oldSurr {
		oldSet[c] = true
	}
	newSet := make(map[CellCoord]bool, len(newSurr))
	for _, c := range newSurr {
		newSet[c] = true
	}

	for _, c := range newSurr {
		if !oldSet[c] {
			enteredCells = append(enteredCells, c)
		}
	}
	for _, c := range oldSurr {
		if !newSet[c] {
			leftCells = append(leftCells, c)
		}
	}

	return true, enteredCells, leftCells
}

func (g *SpatialGrid) getSurroundingCells(center CellCoord) []CellCoord {
	cells := make([]CellCoord, 0, 9)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			cells = append(cells, CellCoord{X: center.X + dx, Y: center.Y + dy})
		}
	}
	return cells
}

// GetSurroundingCells retorna las 9 celdas que componen la vecindad 3x3 del centro dado.
func (g *SpatialGrid) GetSurroundingCells(coord CellCoord) []CellCoord {
	return g.getSurroundingCells(coord)
}

// GetEntitiesInCells retorna todas las entidades contenidas en las celdas provistas.
func (g *SpatialGrid) GetEntitiesInCells(coords []CellCoord) []GridEntity {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var result []GridEntity
	for _, c := range coords {
		if cell, ok := g.cells[c]; ok {
			for _, e := range cell {
				result = append(result, e)
			}
		}
	}
	return result
}

// GetSurroundingEntities retorna todas las entidades activas en la vecindad 3x3 de (x, y).
func (g *SpatialGrid) GetSurroundingEntities(x, y int) []GridEntity {
	coord := g.GetCellCoord(x, y)
	surr := g.getSurroundingCells(coord)
	return g.GetEntitiesInCells(surr)
}

// GetSurroundingPlayers retorna todos los jugadores activos en la vecindad 3x3 de (x, y).
func (g *SpatialGrid) GetSurroundingPlayers(x, y int) []*ZonePlayer {
	entities := g.GetSurroundingEntities(x, y)
	var players []*ZonePlayer
	for _, e := range entities {
		if zp, ok := e.(*ZonePlayer); ok {
			players = append(players, zp)
		}
	}
	return players
}

// BroadcastAOI transmite un paquete únicamente a los jugadores en la vecindad 3x3 alrededor de (x, y).
func (g *SpatialGrid) BroadcastAOI(pw *protocol.PacketWriter, x, y int, excludeHandle uint16) {
	players := g.GetSurroundingPlayers(x, y)
	for _, p := range players {
		if p.MapObjectID == excludeHandle || p.Client == nil || p.HideExceptMe {
			continue
		}
		_ = p.Client.SendPacket(pw)
	}
}

