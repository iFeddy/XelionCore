package areabmp

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CFsAreaBMP representa una máscara espacial de impacto AoE cargada desde archivos .bmp
// (equivalente oficial a la clase C++ CFsAreaBMP de Zone.exe).
type CFsAreaBMP struct {
	Name     string
	Width    int
	Height   int
	Bpp      int
	RowBytes int
	// Data almacena los píxeles en un único bloque plano 1D contiguo en orden top-down: Data[y*Width+x]
	Data []byte
	// Pixels expone vistas sub-slice por fila respaldadas por el mismo backing array Data (0 bytes extra de píxeles)
	Pixels [][]byte
}

// LoadBMP deserializa una imagen BMP de 8 bpp o 1 bpp desde un io.Reader usando
// un único buffer de fila reutilizable y un backing array plano 1D contiguo.
func LoadBMP(r io.Reader, name string) (*CFsAreaBMP, error) {
	var fileHdr [14]byte
	if _, err := io.ReadFull(r, fileHdr[:]); err != nil {
		return nil, fmt.Errorf("error leyendo cabecera de archivo BMP: %w", err)
	}

	if fileHdr[0] != 'B' || fileHdr[1] != 'M' {
		return nil, fmt.Errorf("formato no reconocido: no comienza con 'BM'")
	}

	pixelOffset := binary.LittleEndian.Uint32(fileHdr[10:14])

	var dibHdr [40]byte
	if _, err := io.ReadFull(r, dibHdr[:]); err != nil {
		return nil, fmt.Errorf("error leyendo cabecera DIB BMP: %w", err)
	}

	dibSize := binary.LittleEndian.Uint32(dibHdr[0:4])
	if dibSize < 40 {
		return nil, fmt.Errorf("cabecera DIB no soportada: tamaño %d", dibSize)
	}

	width := int(int32(binary.LittleEndian.Uint32(dibHdr[4:8])))
	height := int(int32(binary.LittleEndian.Uint32(dibHdr[8:12])))
	bpp := int(binary.LittleEndian.Uint16(dibHdr[14:16]))

	if width <= 0 {
		return nil, fmt.Errorf("ancho de imagen inválido: %d", width)
	}

	isTopDown := false
	if height < 0 {
		isTopDown = true
		height = -height
	}
	if height <= 0 {
		return nil, fmt.Errorf("alto de imagen inválido: %d", height)
	}

	// Saltar paleta y bytes de cabecera adicionales hasta pixelOffset
	currentPos := uint32(14 + 40)
	if pixelOffset > currentPos {
		skip := pixelOffset - currentPos
		if _, err := io.CopyN(io.Discard, r, int64(skip)); err != nil {
			return nil, fmt.Errorf("error avanzando hasta datos de píxeles: %w", err)
		}
	}

	// En BMP cada fila está alineada a múltiplos de 4 bytes
	rowBytes := ((width*bpp + 31) / 32) * 4

	data := make([]byte, width*height)
	pixels := make([][]byte, height)
	for y := 0; y < height; y++ {
		pixels[y] = data[y*width : (y+1)*width]
	}

	// Único buffer de fila reutilizable (elimina el doble buffer rawRows completo)
	rowBuf := make([]byte, rowBytes)
	for y := 0; y < height; y++ {
		if _, err := io.ReadFull(r, rowBuf); err != nil {
			return nil, fmt.Errorf("error leyendo fila %d de BMP: %w", y, err)
		}

		dstY := y
		if !isTopDown {
			dstY = height - 1 - y // En BMP estándar, la primera fila en el archivo es la inferior
		}
		dstRow := data[dstY*width : (dstY+1)*width]

		if bpp == 8 {
			copy(dstRow, rowBuf[:width])
		} else if bpp == 1 {
			for x := 0; x < width; x++ {
				byteVal := rowBuf[x/8]
				bit := 7 - (x % 8)
				if (byteVal & (1 << bit)) != 0 {
					dstRow[x] = 1
				}
			}
		} else {
			limit := len(rowBuf)
			if limit > width {
				limit = width
			}
			copy(dstRow, rowBuf[:limit])
		}
	}

	return &CFsAreaBMP{
		Name:     name,
		Width:    width,
		Height:   height,
		Bpp:      bpp,
		RowBytes: rowBytes,
		Data:     data,
		Pixels:   pixels,
	}, nil
}

// LoadBMPFile lee una máscara .bmp desde disco.
func LoadBMPFile(path string) (*CFsAreaBMP, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return LoadBMP(f, base)
}

// PixelAt retorna el valor del píxel (x, y) en orden top-down (0 si está fuera de rango).
func (bmp *CFsAreaBMP) PixelAt(x, y int) byte {
	if bmp == nil || x < 0 || x >= bmp.Width || y < 0 || y >= bmp.Height {
		return 0
	}
	if len(bmp.Data) == bmp.Width*bmp.Height {
		return bmp.Data[y*bmp.Width+x]
	}
	if y < len(bmp.Pixels) && x < len(bmp.Pixels[y]) {
		return bmp.Pixels[y][x]
	}
	return 0
}

// IsHit comprueba si la celda (x, y) de la máscara contiene un valor de impacto activo (> 0).
func (bmp *CFsAreaBMP) IsHit(x, y int) bool {
	return bmp.PixelAt(x, y) > 0
}

// ContainsTarget evalúa si una posición objetivo (targetX, targetY) es alcanzada por el AoE
// proyectado desde el lanzador (casterX, casterY) orientado con ángulo dirDeg (en grados 0..360).
func (bmp *CFsAreaBMP) ContainsTarget(casterX, casterY, dirDeg, targetX, targetY float64) bool {
	if bmp == nil || (len(bmp.Data) == 0 && len(bmp.Pixels) == 0) {
		return false
	}

	dx := targetX - casterX
	dy := targetY - casterY

	// Rotación relativa según orientación del lanzador
	rad := dirDeg * math.Pi / 180.0
	cosA := math.Cos(-rad)
	sinA := math.Sin(-rad)

	rotX := dx*cosA - dy*sinA
	rotY := dx*sinA + dy*cosA

	// El origen de la habilidad en el BMP se centra en el punto medio del bitmap
	px := int(math.Round(rotX + float64(bmp.Width)/2.0))
	py := int(math.Round(rotY + float64(bmp.Height)/2.0))

	return bmp.IsHit(px, py)
}

// Store cataloga las máscaras .bmp de data/AreaBMP/.
type Store struct {
	mu      sync.RWMutex
	BMPs    map[string]*CFsAreaBMP
	byLower map[string]*CFsAreaBMP
}

// NewStore inicializa el catálogo de máscaras AoE.
func NewStore() *Store {
	return &Store{
		BMPs:    make(map[string]*CFsAreaBMP),
		byLower: make(map[string]*CFsAreaBMP),
	}
}

// LoadAll carga todas las máscaras .bmp encontradas en el directorio indicado.
func LoadAll(areaBMPDir string) (*Store, error) {
	entries, err := os.ReadDir(areaBMPDir)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio AreaBMP '%s': %w", areaBMPDir, err)
	}

	store := NewStore()

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.ToLower(filepath.Ext(name)) != ".bmp" {
			continue
		}

		fullPath := filepath.Join(areaBMPDir, name)
		bmp, err := LoadBMPFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("error cargando máscara %s: %w", name, err)
		}

		base := strings.TrimSuffix(name, filepath.Ext(name))
		lowerBase := strings.ToLower(base)
		lowerName := strings.ToLower(name)
		store.BMPs[base] = bmp
		store.BMPs[lowerName] = bmp
		store.byLower[lowerBase] = bmp
		store.byLower[lowerName] = bmp
	}

	return store, nil
}

// Count retorna la cantidad de máscaras .bmp únicas cargadas.
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	seen := make(map[*CFsAreaBMP]bool)
	for _, b := range s.BMPs {
		seen[b] = true
	}
	return len(seen)
}

// GetBMP busca una máscara por nombre base o nombre de archivo en O(1) (insensible a mayúsculas/minúsculas).
func (s *Store) GetBMP(name string) *CFsAreaBMP {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cleaned := strings.TrimSuffix(name, filepath.Ext(name))
	if b, ok := s.BMPs[cleaned]; ok {
		return b
	}
	if s.byLower != nil {
		if b, ok := s.byLower[strings.ToLower(cleaned)]; ok {
			return b
		}
		if b, ok := s.byLower[strings.ToLower(name)]; ok {
			return b
		}
	}
	for k, v := range s.BMPs {
		if strings.EqualFold(k, cleaned) || strings.EqualFold(k, name) {
			return v
		}
	}
	return nil
}
