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
	// Pixels almacena la matriz en orden top-down: [y][x]
	Pixels [][]byte
}

// LoadBMP deserializa una imagen BMP de 8 bpp desde un io.Reader.
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

	rawRows := make([][]byte, height)
	for y := 0; y < height; y++ {
		row := make([]byte, rowBytes)
		if _, err := io.ReadFull(r, row); err != nil {
			return nil, fmt.Errorf("error leyendo fila %d de BMP: %w", y, err)
		}
		rawRows[y] = row
	}

	// Normalizar a matriz top-down [y][x]
	pixels := make([][]byte, height)
	for y := 0; y < height; y++ {
		srcY := y
		if !isTopDown {
			srcY = height - 1 - y // En BMP estándar, la primera fila en el archivo es la inferior
		}
		row := make([]byte, width)
		if bpp == 8 {
			copy(row, rawRows[srcY][:width])
		} else if bpp == 1 {
			for x := 0; x < width; x++ {
				byteVal := rawRows[srcY][x/8]
				bit := 7 - (x % 8)
				if (byteVal & (1 << bit)) != 0 {
					row[x] = 1
				}
			}
		} else {
			// Para otros formatos no indexados, truncar o copiar hasta el ancho disponible
			limit := len(rawRows[srcY])
			if limit > width {
				limit = width
			}
			copy(row, rawRows[srcY][:limit])
		}
		pixels[y] = row
	}

	return &CFsAreaBMP{
		Name:     name,
		Width:    width,
		Height:   height,
		Bpp:      bpp,
		RowBytes: rowBytes,
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

// IsHit comprueba si la celda (x, y) de la máscara contiene un valor de impacto activo (> 0).
func (bmp *CFsAreaBMP) IsHit(x, y int) bool {
	if bmp == nil || len(bmp.Pixels) == 0 {
		return false
	}
	if x < 0 || x >= bmp.Width || y < 0 || y >= bmp.Height {
		return false
	}
	return bmp.Pixels[y][x] > 0
}

// ContainsTarget evalúa si una posición objetivo (targetX, targetY) es alcanzada por el AoE
// proyectado desde el lanzador (casterX, casterY) orientado con ángulo dirDeg (en grados 0..360).
func (bmp *CFsAreaBMP) ContainsTarget(casterX, casterY, dirDeg, targetX, targetY float64) bool {
	if bmp == nil || len(bmp.Pixels) == 0 {
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
	mu   sync.RWMutex
	BMPs map[string]*CFsAreaBMP
}

// NewStore inicializa el catálogo de máscaras AoE.
func NewStore() *Store {
	return &Store{
		BMPs: make(map[string]*CFsAreaBMP),
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
		store.BMPs[base] = bmp
		// También registrar con nombre de archivo completo en minúsculas para resolución robusta
		store.BMPs[strings.ToLower(name)] = bmp
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

// GetBMP busca una máscara por nombre base o nombre de archivo (búsqueda insensible a mayúsculas/minúsculas).
func (s *Store) GetBMP(name string) *CFsAreaBMP {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cleaned := strings.TrimSuffix(name, filepath.Ext(name))
	if b, ok := s.BMPs[cleaned]; ok {
		return b
	}
	for k, v := range s.BMPs {
		if strings.EqualFold(k, cleaned) || strings.EqualFold(k, name) {
			return v
		}
	}
	return nil
}
