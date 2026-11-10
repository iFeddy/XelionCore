package blockinfo

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

// MapBlockInformation representa la matriz de bloqueo y transitabilidad cargada desde archivos .shbd oficiales
// (equivalente exacto en C++ a MapBlock::MapBlockInformation en Zone.pdb).
type MapBlockInformation struct {
	RowBytes    uint32
	HeightTiles uint32
	WidthTiles  uint32
	Data        []byte
}

// LoadSHBD deserializa un archivo .shbd desde un io.Reader.
// Cabecera oficial: uint32 rowBytes, uint32 heightTiles, seguido de rowBytes * heightTiles bytes de datos.
func LoadSHBD(r io.Reader) (*MapBlockInformation, error) {
	var header [8]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, fmt.Errorf("error leyendo cabecera shbd: %w", err)
	}

	rowBytes := binary.LittleEndian.Uint32(header[0:4])
	heightTiles := binary.LittleEndian.Uint32(header[4:8])
	widthTiles := rowBytes * 8

	totalBytes := int(rowBytes * heightTiles)
	data := make([]byte, totalBytes)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("error leyendo payload shbd (%d bytes): %w", totalBytes, err)
	}

	return &MapBlockInformation{
		RowBytes:    rowBytes,
		HeightTiles: heightTiles,
		WidthTiles:  widthTiles,
		Data:        data,
	}, nil
}

// LoadSHBDFile lee un archivo .shbd desde disco.
func LoadSHBDFile(path string) (*MapBlockInformation, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadSHBD(f)
}

// IsMoveBlocked verifica en O(1) si una celda en coordenadas de tile está bloqueada para el movimiento.
// En Zone.exe (0x0049cf70): 1 bit por celda. 1 = bloqueado, 0 = transitable. Fuera de límites se considera bloqueado.
func (mbi *MapBlockInformation) IsMoveBlocked(tileX, tileY uint32) bool {
	if mbi == nil || len(mbi.Data) == 0 {
		return false
	}
	if tileX >= mbi.WidthTiles || tileY >= mbi.HeightTiles {
		return true
	}

	byteIdx := tileY*mbi.RowBytes + (tileX / 8)
	if int(byteIdx) >= len(mbi.Data) {
		return true
	}

	bit := tileX & 7
	return (mbi.Data[byteIdx] & (1 << bit)) != 0
}

// IsAttackBlocked verifica si la línea de visión (LoS) entre dos celdas de tile está obstruida.
// Aplica el algoritmo determinístico de trazado de líneas de Bresenham verificando IsMoveBlocked en cada paso.
func (mbi *MapBlockInformation) IsAttackBlocked(x0, y0, x1, y1 uint32) bool {
	if mbi == nil || len(mbi.Data) == 0 {
		return false
	}

	// Si el origen o destino están fuera de rango
	if x0 >= mbi.WidthTiles || y0 >= mbi.HeightTiles || x1 >= mbi.WidthTiles || y1 >= mbi.HeightTiles {
		return true
	}

	dx := int(x1) - int(x0)
	dy := int(y1) - int(y0)
	sx := 1
	if dx < 0 {
		sx = -1
		dx = -dx
	}
	sy := 1
	if dy < 0 {
		sy = -1
		dy = -dy
	}

	err := dx - dy
	curX := int(x0)
	curY := int(y0)

	for {
		// Validar celda intermedia (excluyendo el origen exacto)
		if (uint32(curX) != x0 || uint32(curY) != y0) && (uint32(curX) != x1 || uint32(curY) != y1) {
			if mbi.IsMoveBlocked(uint32(curX), uint32(curY)) {
				return true
			}
		}

		if curX == int(x1) && curY == int(y1) {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			curX += sx
		}
		if e2 < dx {
			err += dx
			curY += sy
		}
	}

	return false
}

// FindNearestWalkable localiza en espiral la celda transitable más cercana dentro de maxRadius tiles.
func (mbi *MapBlockInformation) FindNearestWalkable(tileX, tileY uint32, maxRadius uint32) (uint32, uint32, bool) {
	if !mbi.IsMoveBlocked(tileX, tileY) {
		return tileX, tileY, true
	}

	for r := uint32(1); r <= maxRadius; r++ {
		// Inspeccionar el perímetro del cuadrado de radio r
		minX := int(tileX) - int(r)
		maxX := int(tileX) + int(r)
		minY := int(tileY) - int(r)
		maxY := int(tileY) + int(r)

		// Lados horizontales
		for x := minX; x <= maxX; x++ {
			if x >= 0 && x < int(mbi.WidthTiles) {
				if minY >= 0 && minY < int(mbi.HeightTiles) && !mbi.IsMoveBlocked(uint32(x), uint32(minY)) {
					return uint32(x), uint32(minY), true
				}
				if maxY >= 0 && maxY < int(mbi.HeightTiles) && !mbi.IsMoveBlocked(uint32(x), uint32(maxY)) {
					return uint32(x), uint32(maxY), true
				}
			}
		}

		// Lados verticales
		for y := minY + 1; y < maxY; y++ {
			if y >= 0 && y < int(mbi.HeightTiles) {
				if minX >= 0 && minX < int(mbi.WidthTiles) && !mbi.IsMoveBlocked(uint32(minX), uint32(y)) {
					return uint32(minX), uint32(y), true
				}
				if maxX >= 0 && maxX < int(mbi.WidthTiles) && !mbi.IsMoveBlocked(uint32(maxX), uint32(y)) {
					return uint32(maxX), uint32(y), true
				}
			}
		}
	}

	return tileX, tileY, false
}

// WorldToTile convierte coordenadas continuas de mundo al espacio discreto de celdas tile según la escala oficial de Zone.exe.
// Fórmula binaria de Zone.exe (0x0049cf70): tileX = (worldX * 8) / 100
func WorldToTile(worldX, worldY float64) (uint32, uint32) {
	tx := uint32(math.Floor(worldX * 8.0 / 100.0))
	ty := uint32(math.Floor(worldY * 8.0 / 100.0))
	return tx, ty
}

// TileToWorld convierte coordenadas de celda tile a coordenadas continuas de mundo (centro de celda).
func TileToWorld(tileX, tileY uint32) (float64, float64) {
	wx := (float64(tileX) + 0.5) * 100.0 / 8.0
	wy := (float64(tileY) + 0.5) * 100.0 / 8.0
	return wx, wy
}

// IsWorldMoveBlocked consulta si una coordenada de mundo está bloqueada.
func (mbi *MapBlockInformation) IsWorldMoveBlocked(worldX, worldY float64) bool {
	tx, ty := WorldToTile(worldX, worldY)
	return mbi.IsMoveBlocked(tx, ty)
}

// IsWorldAttackBlocked consulta si una trayectoria en coordenadas de mundo está bloqueada.
func (mbi *MapBlockInformation) IsWorldAttackBlocked(wx0, wy0, wx1, wy1 float64) bool {
	tx0, ty0 := WorldToTile(wx0, wy0)
	tx1, ty1 := WorldToTile(wx1, wy1)
	return mbi.IsAttackBlocked(tx0, ty0, tx1, ty1)
}

