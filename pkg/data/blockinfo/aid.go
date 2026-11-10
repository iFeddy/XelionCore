package blockinfo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

// Tipos de área oficial de AreaInfoData en Zone.pdb
const (
	AreaTypeCircle = 0 // AreaInfoData::AE_Circle
	AreaTypeSquare = 1 // AreaInfoData::AE_Square
)

// IArea define la interfaz común de consulta espacial de áreas.
type IArea interface {
	AreaName() string
	Type() uint32
	Contains(x, y float64) bool
}

// AreaCircle representa un área circular de evento (AE_Circle).
type AreaCircle struct {
	Name    string
	CenterX float32
	CenterY float32
	Radius  float32
}

func (a *AreaCircle) AreaName() string { return a.Name }
func (a *AreaCircle) Type() uint32     { return AreaTypeCircle }
func (a *AreaCircle) Contains(x, y float64) bool {
	dx := x - float64(a.CenterX)
	dy := y - float64(a.CenterY)
	distSq := dx*dx + dy*dy
	r := float64(a.Radius)
	return distSq <= (r * r)
}

// AreaSquare representa un área cuadrada o rectangular de evento con rotación opcional (AE_Square).
type AreaSquare struct {
	Name     string
	CenterX  float32
	CenterY  float32
	ExtentX  float32 // Ancho total
	ExtentY  float32 // Alto total
	Rotation float32 // Ángulo en radianes
}

func (a *AreaSquare) AreaName() string { return a.Name }
func (a *AreaSquare) Type() uint32     { return AreaTypeSquare }
func (a *AreaSquare) Contains(x, y float64) bool {
	dx := x - float64(a.CenterX)
	dy := y - float64(a.CenterY)

	// Si hay rotación, rotar inversamente el vector relativo
	if a.Rotation != 0 {
		cosA := math.Cos(-float64(a.Rotation))
		sinA := math.Sin(-float64(a.Rotation))
		rx := dx*cosA - dy*sinA
		ry := dx*sinA + dy*cosA
		dx = rx
		dy = ry
	}

	halfW := float64(a.ExtentX) / 2.0
	halfH := float64(a.ExtentY) / 2.0

	return math.Abs(dx) <= halfW && math.Abs(dy) <= halfH
}

// AreaInfoData administra las áreas de un mapa cargadas desde archivos .aid
// (equivalente a AreaInfoData en Zone.pdb).
type AreaInfoData struct {
	Areas      []IArea
	ByName     map[string][]IArea
}

// LoadAID deserializa un archivo .aid desde un io.Reader.
func LoadAID(r io.Reader) (*AreaInfoData, error) {
	var countBuf [4]byte
	if _, err := io.ReadFull(r, countBuf[:]); err != nil {
		return nil, fmt.Errorf("error leyendo contador aid: %w", err)
	}

	count := binary.LittleEndian.Uint32(countBuf[0:4])
	aid := &AreaInfoData{
		Areas:  make([]IArea, 0, count),
		ByName: make(map[string][]IArea, count),
	}

	for i := uint32(0); i < count; i++ {
		var nameBuf [32]byte
		if _, err := io.ReadFull(r, nameBuf[:]); err != nil {
			return nil, fmt.Errorf("error leyendo nombre de área %d: %w", i, err)
		}

		idx := bytes.IndexByte(nameBuf[:], 0)
		if idx == -1 {
			idx = 32
		}
		name := string(nameBuf[:idx])

		var typeBuf [4]byte
		if _, err := io.ReadFull(r, typeBuf[:]); err != nil {
			return nil, fmt.Errorf("error leyendo tipo de área %s: %w", name, err)
		}
		areaType := binary.LittleEndian.Uint32(typeBuf[0:4])

		if areaType == AreaTypeCircle {
			var circBuf [12]byte
			if _, err := io.ReadFull(r, circBuf[:]); err != nil {
				return nil, fmt.Errorf("error leyendo parámetros circulares de área %s: %w", name, err)
			}
			cx := math.Float32frombits(binary.LittleEndian.Uint32(circBuf[0:4]))
			cy := math.Float32frombits(binary.LittleEndian.Uint32(circBuf[4:8]))
			rad := math.Float32frombits(binary.LittleEndian.Uint32(circBuf[8:12]))

			circle := &AreaCircle{
				Name:    name,
				CenterX: cx,
				CenterY: cy,
				Radius:  rad,
			}
			aid.Areas = append(aid.Areas, circle)
			aid.ByName[name] = append(aid.ByName[name], circle)
		} else if areaType == AreaTypeSquare {
			var sqBuf [20]byte
			if _, err := io.ReadFull(r, sqBuf[:]); err != nil {
				return nil, fmt.Errorf("error leyendo parámetros rectangulares de área %s: %w", name, err)
			}
			cx := math.Float32frombits(binary.LittleEndian.Uint32(sqBuf[0:4]))
			cy := math.Float32frombits(binary.LittleEndian.Uint32(sqBuf[4:8]))
			ex := math.Float32frombits(binary.LittleEndian.Uint32(sqBuf[8:12]))
			ey := math.Float32frombits(binary.LittleEndian.Uint32(sqBuf[12:16]))
			rot := math.Float32frombits(binary.LittleEndian.Uint32(sqBuf[16:20]))

			square := &AreaSquare{
				Name:     name,
				CenterX:  cx,
				CenterY:  cy,
				ExtentX:  ex,
				ExtentY:  ey,
				Rotation: rot,
			}
			aid.Areas = append(aid.Areas, square)
			aid.ByName[name] = append(aid.ByName[name], square)
		} else {
			return nil, fmt.Errorf("tipo de área desconocido %d en área %s", areaType, name)
		}
	}

	return aid, nil
}

// LoadAIDFile lee un archivo .aid desde disco.
func LoadAIDFile(path string) (*AreaInfoData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadAID(f)
}

// Contains comprueba si la coordenada (x, y) está contenida dentro de alguna de las áreas con el nombre dado.
func (aid *AreaInfoData) Contains(areaName string, x, y float64) bool {
	if aid == nil {
		return false
	}
	areas, ok := aid.ByName[areaName]
	if !ok {
		return false
	}
	for _, a := range areas {
		if a.Contains(x, y) {
			return true
		}
	}
	return false
}

// AreasAt retorna la lista de nombres de todas las áreas que contienen la coordenada (x, y).
func (aid *AreaInfoData) AreasAt(x, y float64) []string {
	if aid == nil {
		return nil
	}
	var res []string
	seen := make(map[string]bool)
	for _, a := range aid.Areas {
		if a.Contains(x, y) {
			if !seen[a.AreaName()] {
				seen[a.AreaName()] = true
				res = append(res, a.AreaName())
			}
		}
	}
	return res
}

