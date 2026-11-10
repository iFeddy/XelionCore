package blockinfo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
)

// Constantes de acción de puerta compatibles con DoorAction de Zone.exe
const (
	DoorActionClose = 0
	DoorActionOpen  = 1
)

// MapDoorBlockElement representa un elemento de puerta o puente dinámico cargado desde archivos .sbi
// (equivalente a MapBlock::MapDoorBlockElement en Zone.pdb).
type MapDoorBlockElement struct {
	Name        string
	Left        uint32
	Top         uint32
	Right       uint32
	Bottom      uint32
	DataSize    uint32
	Extra       uint32
	MoveBlock   []byte
	AttackBlock []byte
	IsOpen      bool
}

// MapDoorArray gestiona el conjunto de puertas dinámicas de un mapa (.sbi)
// (equivalente a MapBlock::MapDoorArray en Zone.pdb).
type MapDoorArray struct {
	mu       sync.RWMutex
	Doors    map[string]*MapDoorBlockElement
	DoorList []*MapDoorBlockElement
}

// LoadSBI deserializa un archivo .sbi desde un io.Reader.
func LoadSBI(r io.Reader) (*MapDoorArray, error) {
	var countBuf [4]byte
	if _, err := io.ReadFull(r, countBuf[:]); err != nil {
		return nil, fmt.Errorf("error leyendo contador sbi: %w", err)
	}

	count := binary.LittleEndian.Uint32(countBuf[0:4])
	doors := make([]*MapDoorBlockElement, count)
	doorMap := make(map[string]*MapDoorBlockElement, count)

	// Leer las cabeceras de cada puerta (56 bytes c/u)
	for i := uint32(0); i < count; i++ {
		var hdr [56]byte
		if _, err := io.ReadFull(r, hdr[:]); err != nil {
			return nil, fmt.Errorf("error leyendo cabecera de puerta %d: %w", i, err)
		}

		nameRaw := hdr[0:32]
		idx := bytes.IndexByte(nameRaw, 0)
		if idx == -1 {
			idx = 32
		}
		name := string(nameRaw[:idx])

		left := binary.LittleEndian.Uint32(hdr[32:36])
		top := binary.LittleEndian.Uint32(hdr[36:40])
		right := binary.LittleEndian.Uint32(hdr[40:44])
		bottom := binary.LittleEndian.Uint32(hdr[44:48])
		dataSize := binary.LittleEndian.Uint32(hdr[48:52])
		extra := binary.LittleEndian.Uint32(hdr[52:56])

		door := &MapDoorBlockElement{
			Name:     name,
			Left:     left,
			Top:      top,
			Right:    right,
			Bottom:   bottom,
			DataSize: dataSize,
			Extra:    extra,
			IsOpen:   false, // Cerrada por defecto
		}

		doors[i] = door
		doorMap[name] = door
	}

	// Leer los buffers de mapa de bits (2 layers por puerta: MoveBlock y AttackBlock)
	for i := uint32(0); i < count; i++ {
		door := doors[i]
		if door.DataSize > 0 {
			door.MoveBlock = make([]byte, door.DataSize)
			if _, err := io.ReadFull(r, door.MoveBlock); err != nil {
				return nil, fmt.Errorf("error leyendo MoveBlock de puerta %s: %w", door.Name, err)
			}

			door.AttackBlock = make([]byte, door.DataSize)
			if _, err := io.ReadFull(r, door.AttackBlock); err != nil {
				return nil, fmt.Errorf("error leyendo AttackBlock de puerta %s: %w", door.Name, err)
			}
		}
	}

	return &MapDoorArray{
		Doors:    doorMap,
		DoorList: doors,
	}, nil
}

// LoadSBIFile lee un archivo .sbi desde disco.
func LoadSBIFile(path string) (*MapDoorArray, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadSBI(f)
}

// Clone produce una copia independiente de las puertas para una instancia dinámica de mapa,
// permitiendo abrir o cerrar puertas sin mutar la plantilla base compartida en RAM.
func (mda *MapDoorArray) Clone() *MapDoorArray {
	if mda == nil {
		return nil
	}
	mda.mu.RLock()
	defer mda.mu.RUnlock()

	clone := &MapDoorArray{
		Doors:    make(map[string]*MapDoorBlockElement, len(mda.Doors)),
		DoorList: make([]*MapDoorBlockElement, len(mda.DoorList)),
	}

	for i, d := range mda.DoorList {
		doorCopy := &MapDoorBlockElement{
			Name:        d.Name,
			Left:        d.Left,
			Top:         d.Top,
			Right:       d.Right,
			Bottom:      d.Bottom,
			DataSize:    d.DataSize,
			Extra:       d.Extra,
			MoveBlock:   d.MoveBlock,   // Referencia compartida inmutable a los bits
			AttackBlock: d.AttackBlock, // Referencia compartida inmutable a los bits
			IsOpen:      d.IsOpen,
		}
		clone.DoorList[i] = doorCopy
		clone.Doors[d.Name] = doorCopy
	}

	return clone
}

// DoorAction abre o cierra una puerta por nombre (equivalente a MapBlock::MapDoorArray::mda_DoorAction).
func (mda *MapDoorArray) DoorAction(doorName string, action int) bool {
	if mda == nil {
		return false
	}
	mda.mu.Lock()
	defer mda.mu.Unlock()

	door, ok := mda.Doors[doorName]
	if !ok {
		return false
	}

	door.IsOpen = (action == DoorActionOpen)
	return true
}

// CloseAllDoors cierra todas las puertas del mapa (equivalente a mda_CloseAllDoor).
func (mda *MapDoorArray) CloseAllDoors() {
	if mda == nil {
		return
	}
	mda.mu.Lock()
	defer mda.mu.Unlock()

	for _, d := range mda.DoorList {
		d.IsOpen = false
	}
}

// IsDoorMoveBlocked comprueba si una celda tile está bloqueada por alguna puerta cerrada en esta instancia.
func (mda *MapDoorArray) IsDoorMoveBlocked(tileX, tileY uint32) bool {
	if mda == nil {
		return false
	}
	mda.mu.RLock()
	defer mda.mu.RUnlock()

	for _, d := range mda.DoorList {
		if d.IsOpen || len(d.MoveBlock) == 0 {
			continue
		}
		if tileX >= d.Left && tileX <= d.Right && tileY >= d.Top && tileY <= d.Bottom {
			lx := tileX - d.Left
			ly := tileY - d.Top
			height := d.Bottom - d.Top + 1
			if height == 0 {
				continue
			}
			rowBytes := d.DataSize / height
			byteIdx := ly*rowBytes + (lx / 8)
			if int(byteIdx) < len(d.MoveBlock) {
				bit := lx & 7
				if (d.MoveBlock[byteIdx] & (1 << bit)) != 0 {
					return true
				}
			}
		}
	}
	return false
}

// IsDoorAttackBlocked comprueba si una celda tile tiene obstruida la línea de visión por una puerta cerrada.
func (mda *MapDoorArray) IsDoorAttackBlocked(tileX, tileY uint32) bool {
	if mda == nil {
		return false
	}
	mda.mu.RLock()
	defer mda.mu.RUnlock()

	for _, d := range mda.DoorList {
		if d.IsOpen || len(d.AttackBlock) == 0 {
			continue
		}
		if tileX >= d.Left && tileX <= d.Right && tileY >= d.Top && tileY <= d.Bottom {
			lx := tileX - d.Left
			ly := tileY - d.Top
			height := d.Bottom - d.Top + 1
			if height == 0 {
				continue
			}
			rowBytes := d.DataSize / height
			byteIdx := ly*rowBytes + (lx / 8)
			if int(byteIdx) < len(d.AttackBlock) {
				bit := lx & 7
				if (d.AttackBlock[byteIdx] & (1 << bit)) != 0 {
					return true
				}
			}
		}
	}
	return false
}

