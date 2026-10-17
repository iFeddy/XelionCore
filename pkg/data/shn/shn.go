package shn

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

// Constantes de tipos de columnas del motor CDataReader oficial
const (
	TypeListByte     uint32 = 1
	TypeListWord     uint32 = 2
	TypeListDword    uint32 = 3
	TypeListQword    uint32 = 4
	TypeListFloat    uint32 = 5
	TypeListFileName uint32 = 6
	TypeListFileAuto uint32 = 7
	TypeListRemark   uint32 = 8
	TypeListStr      uint32 = 9
	TypeListStrAuto  uint32 = 10
	TypeListInx      uint32 = 11
	TypeListInxByte  uint32 = 12
	TypeListInxWord  uint32 = 13
	TypeListInxDword uint32 = 14
	TypeListInxQword uint32 = 15
)

// Decrypt implementa el algoritmo simétrico oficial CDataReader::Encription
// extraído directamente del flujo de instrucciones x86 de WorldManager.exe.
func Decrypt(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}

	buf := make([]byte, length)
	copy(buf, data)

	ecx := length - 1
	key := byte(length & 0xFF)
	al := key

	for ecx >= 0 {
		buf[ecx] ^= al

		al = byte(ecx & 0xFF)
		al = byte((uint16(al) * 11) & 0xFF)
		dl := byte((ecx & 0x0F) + 0x55)
		al ^= dl
		al ^= key
		al ^= 0xAA

		key = al
		ecx--
	}

	return buf
}

// Column representa un descriptor de columna CDataReader::FIELD (56 bytes exactos).
type Column struct {
	Name string
	Type uint32
	Size uint32
}

// Row representa una fila de datos parseada de un archivo .shn.
type Row struct {
	values map[string][]byte
}

// GetByte retorna el valor como byte o 0 si no existe.
func (r *Row) GetByte(colName string) byte {
	v, ok := r.values[colName]
	if !ok || len(v) == 0 {
		return 0
	}
	return v[0]
}

// GetUint16 retorna el valor como uint16 en LittleEndian.
func (r *Row) GetUint16(colName string) uint16 {
	v, ok := r.values[colName]
	if !ok || len(v) < 2 {
		return 0
	}
	return binary.LittleEndian.Uint16(v[:2])
}

// GetInt16 retorna el valor como int16 en LittleEndian.
func (r *Row) GetInt16(colName string) int16 {
	return int16(r.GetUint16(colName))
}

// GetUint32 retorna el valor como uint32 en LittleEndian.
func (r *Row) GetUint32(colName string) uint32 {
	v, ok := r.values[colName]
	if !ok || len(v) < 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(v[:4])
}

// GetInt retorna el valor numérico como int.
func (r *Row) GetInt(colName string) int {
	v, ok := r.values[colName]
	if !ok {
		return 0
	}
	switch len(v) {
	case 1:
		return int(v[0])
	case 2:
		return int(int16(binary.LittleEndian.Uint16(v)))
	case 4:
		return int(int32(binary.LittleEndian.Uint32(v)))
	default:
		return 0
	}
}

// GetInt64 retorna el valor como int64.
func (r *Row) GetInt64(colName string) int64 {
	v, ok := r.values[colName]
	if !ok || len(v) < 8 {
		return int64(r.GetInt(colName))
	}
	return int64(binary.LittleEndian.Uint64(v[:8]))
}

// GetString retorna la cadena ASCII/ANSI terminada en nulo.
func (r *Row) GetString(colName string) string {
	v, ok := r.values[colName]
	if !ok {
		return ""
	}
	idx := bytes.IndexByte(v, 0x00)
	if idx != -1 {
		return string(v[:idx])
	}
	return string(v)
}

// RawBytes retorna los bytes originales de la celda.
func (r *Row) RawBytes(colName string) []byte {
	return r.values[colName]
}

// File representa un archivo de datos binario .shn completamente cargado en memoria.
type File struct {
	Version     string
	DataMode    uint32
	RecordCount int
	RowSize     int
	Columns     []Column
	Rows        []Row
}

// ReadFile carga y parsea un archivo .shn desde el disco.
func ReadFile(filePath string) (*File, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SHN file %s: %w", filePath, err)
	}
	return Parse(raw)
}

// Parse parsea los bytes de un archivo .shn respetando la especificación CDataReader oficial.
func Parse(data []byte) (*File, error) {
	if len(data) < 52 {
		return nil, errors.New("invalid SHN file: header is too small")
	}

	dataMode := binary.LittleEndian.Uint32(data[28:32])
	versionStr := string(bytes.TrimRight(data[4:28], "\x00"))

	var full []byte
	if dataMode == 1 {
		// Descifrar cuerpo desde offset 36
		decrypted := Decrypt(data[36:])
		full = make([]byte, 36+len(decrypted))
		copy(full[:36], data[:36])
		copy(full[36:], decrypted)
	} else {
		full = data
	}

	recordCount := int(binary.LittleEndian.Uint32(full[40:44]))
	rowSize := int(binary.LittleEndian.Uint32(full[44:48]))
	colCount := int(binary.LittleEndian.Uint32(full[48:52]))

	if colCount < 0 || colCount > 500 {
		return nil, fmt.Errorf("invalid column count in SHN: %d", colCount)
	}

	headerEnd := 52 + colCount*56
	expectedTotal := headerEnd + recordCount*rowSize
	if len(full) < expectedTotal {
		return nil, fmt.Errorf("SHN data length mismatch: expected %d bytes, got %d", expectedTotal, len(full))
	}

	// Leer descriptores de columnas (56 bytes cada uno)
	columns := make([]Column, colCount)
	colOffset := 52
	for i := 0; i < colCount; i++ {
		desc := full[colOffset : colOffset+56]
		colName := string(bytes.TrimRight(desc[:48], "\x00"))
		colType := binary.LittleEndian.Uint32(desc[48:52])
		colSize := binary.LittleEndian.Uint32(desc[52:56])

		columns[i] = Column{
			Name: colName,
			Type: colType,
			Size: colSize,
		}
		colOffset += 56
	}

	// Leer filas de registros
	rows := make([]Row, recordCount)
	rowOffset := headerEnd
	for r := 0; r < recordCount; r++ {
		rowBytes := full[rowOffset : rowOffset+rowSize]
		rowVals := make(map[string][]byte, colCount)

		// En Fiesta SHN los primeros 2 bytes de cada fila son uint16 row_len
		curr := 2
		for _, col := range columns {
			end := curr + int(col.Size)
			if end <= len(rowBytes) {
				val := make([]byte, col.Size)
				copy(val, rowBytes[curr:end])
				rowVals[col.Name] = val
			}
			curr = end
		}

		rows[r] = Row{values: rowVals}
		rowOffset += rowSize
	}

	return &File{
		Version:     versionStr,
		DataMode:    dataMode,
		RecordCount: recordCount,
		RowSize:     rowSize,
		Columns:     columns,
		Rows:        rows,
	}, nil
}

// FindRows busca todas las filas cuyo valor en colName sea igual a targetVal.
func (f *File) FindRows(colName string, targetVal string) []*Row {
	var result []*Row
	for i := range f.Rows {
		if f.Rows[i].GetString(colName) == targetVal {
			result = append(result, &f.Rows[i])
		}
	}
	return result
}

// FindRow busca la primera fila cuyo valor en colName sea igual a targetVal.
func (f *File) FindRow(colName string, targetVal string) *Row {
	for i := range f.Rows {
		if f.Rows[i].GetString(colName) == targetVal {
			return &f.Rows[i]
		}
	}
	return nil
}
