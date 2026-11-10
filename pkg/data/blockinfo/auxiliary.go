package blockinfo

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// ArenaAttribute representa matrices zonales de atributos de arena cargadas desde archivos .shab y .shad
// (ej. KDArena.shab, UrgDragon.shad).
type ArenaAttribute struct {
	RowBytes    uint32
	HeightTiles uint32
	WidthTiles  uint32
	Data        []byte
}

// LoadArenaAttribute lee archivos .shab o .shad.
func LoadArenaAttribute(r io.Reader) (*ArenaAttribute, error) {
	var header [8]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, fmt.Errorf("error leyendo cabecera arena attribute: %w", err)
	}

	rowBytes := binary.LittleEndian.Uint32(header[0:4])
	heightTiles := binary.LittleEndian.Uint32(header[4:8])
	widthTiles := rowBytes * 8

	totalBytes := int(rowBytes * heightTiles)
	data := make([]byte, totalBytes)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("error leyendo payload arena attribute (%d bytes): %w", totalBytes, err)
	}

	return &ArenaAttribute{
		RowBytes:    rowBytes,
		HeightTiles: heightTiles,
		WidthTiles:  widthTiles,
		Data:        data,
	}, nil
}

// LoadArenaAttributeFile lee un archivo .shab o .shad desde disco.
func LoadArenaAttributeFile(path string) (*ArenaAttribute, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadArenaAttribute(f)
}

// IsBlocked verifica si una celda de la arena tiene el bit de atributo activado.
func (aa *ArenaAttribute) IsBlocked(tileX, tileY uint32) bool {
	if aa == nil || len(aa.Data) == 0 {
		return false
	}
	if tileX >= aa.WidthTiles || tileY >= aa.HeightTiles {
		return true
	}
	byteIdx := tileY*aa.RowBytes + (tileX / 8)
	if int(byteIdx) >= len(aa.Data) {
		return true
	}
	bit := tileX & 7
	return (aa.Data[byteIdx] & (1 << bit)) != 0
}

// LoadBDTFile lee tablas binarias de detalle de bloques y estructuras (.bdt).
func LoadBDTFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// LoadSHMDFile lee descriptores de escena (.shmd).
func LoadSHMDFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

