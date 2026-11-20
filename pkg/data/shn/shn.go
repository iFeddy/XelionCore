package shn

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
	TypeListVarStr   uint32 = 26
)

// DecryptInPlace implementa el algoritmo simétrico oficial CDataReader::Encription
// in-place sobre buf (0 alocaciones), extraído directamente de WorldManager.exe / Zone.exe.
func DecryptInPlace(buf []byte) {
	length := len(buf)
	if length == 0 {
		return
	}

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
}

// Decrypt implementa el algoritmo simétrico oficial CDataReader::Encription
// devolviendo un slice nuevo sin mutar el buffer de entrada.
func Decrypt(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}

	buf := make([]byte, length)
	copy(buf, data)
	DecryptInPlace(buf)
	return buf
}

// Column representa un descriptor de columna CDataReader::FIELD (56 bytes exactos).
type Column struct {
	Name string
	Type uint32
	Size uint32
}

// tableSchema almacena los metadatos inmutables compartidos por todas las filas de una tabla .shn
// (paridad con CDataReader: 1 instancia por archivo, 0 mapas ni slices de celdas por fila).
type tableSchema struct {
	columns       []Column
	firstColIndex map[string]int // Para File.ColumnIndex(name): primera aparición (0..N-1)
	lastColIndex  map[string]int // Para Row.Get*(colName) / GetBytes(colName): última aparición
	fixedOffsets  []uint32       // Offset fijo en bytes desde el inicio de row.raw (tras los 2B de row_len)
	allFixed      bool           // true si ninguna columna es TypeListVarStr (0x1a)
	isQuestData   bool           // true para el formato nativo CQuestData (version == 6)
}

func newTableSchema(columns []Column) *tableSchema {
	n := len(columns)
	s := &tableSchema{
		columns:       columns,
		firstColIndex: make(map[string]int, n),
		lastColIndex:  make(map[string]int, n),
		fixedOffsets:  make([]uint32, n),
		allFixed:      true,
	}
	var curr uint32 = 2 // En Fiesta SHN los primeros 2 bytes de cada fila son uint16 row_len
	for i, col := range columns {
		if _, exists := s.firstColIndex[col.Name]; !exists {
			s.firstColIndex[col.Name] = i
		}
		s.lastColIndex[col.Name] = i
		s.fixedOffsets[i] = curr
		if col.Type == TypeListVarStr {
			s.allFixed = false
		}
		curr += col.Size
	}
	return s
}

var questDataColumns = []Column{
	{Name: "RecordSize", Type: TypeListWord, Size: 2},
	{Name: "QuestID", Type: TypeListDword, Size: 4},
	{Name: "TotalSize", Type: TypeListDword, Size: 4},
	{Name: "HeaderData", Type: TypeListRemark, Size: 680},
	{Name: "ScriptStart", Type: TypeListStrAuto, Size: 0},
	{Name: "ScriptDoing", Type: TypeListStrAuto, Size: 0},
	{Name: "ScriptEnd", Type: TypeListStrAuto, Size: 0},
}

var questDataSchema = func() *tableSchema {
	s := &tableSchema{
		columns:       questDataColumns,
		firstColIndex: make(map[string]int, len(questDataColumns)),
		lastColIndex:  make(map[string]int, len(questDataColumns)),
		isQuestData:   true,
	}
	for i, c := range questDataColumns {
		s.firstColIndex[c.Name] = i
		s.lastColIndex[c.Name] = i
	}
	return s
}()

// Row representa una fila zero-copy respaldada directamente por el buffer contiguo desencriptado (paridad CDataReader).
type Row struct {
	schema *tableSchema
	raw    []byte
}

// RawRecord retorna el slice crudo completo del registro tal como reside en el buffer desencriptado.
func (r *Row) RawRecord() []byte {
	if r == nil {
		return nil
	}
	return r.raw
}

// ColumnCount retorna la cantidad de columnas posicionales definidas en la tabla.
func (r *Row) ColumnCount() int {
	if r == nil || r.schema == nil {
		return 0
	}
	return len(r.schema.columns)
}

func (r *Row) getQuestBytesAt(idx int) []byte {
	if len(r.raw) < 680 {
		return nil
	}
	rec := r.raw[:680]
	len1 := int(binary.LittleEndian.Uint16(rec[0x294:0x296]))
	len3 := int(binary.LittleEndian.Uint16(rec[0x296:0x298]))
	len2 := int(binary.LittleEndian.Uint16(rec[0x298:0x29A]))

	switch idx {
	case 0: // RecordSize
		return rec[:2]
	case 1: // QuestID
		return rec[4:8]
	case 2: // TotalSize
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], uint32(680+len1+len2+len3))
		return b[:]
	case 3: // HeaderData
		return rec
	case 4: // ScriptStart (s1)
		off := 680
		if off+len1 <= len(r.raw) {
			return r.raw[off : off+len1]
		}
		return nil
	case 5: // ScriptDoing (s3)
		off := 680
		if off+len1 <= len(r.raw) {
			off += len1
		}
		if off+len2 <= len(r.raw) {
			off += len2
		}
		if off+len3 <= len(r.raw) {
			return r.raw[off : off+len3]
		}
		return nil
	case 6: // ScriptEnd (s2)
		off := 680
		if off+len1 <= len(r.raw) {
			off += len1
		}
		if off+len2 <= len(r.raw) {
			return r.raw[off : off+len2]
		}
		return nil
	default:
		return nil
	}
}

// GetBytesAt retorna el slice crudo zero-copy de una columna por índice posicional en O(1).
func (r *Row) GetBytesAt(idx int) []byte {
	if r == nil || r.schema == nil || idx < 0 || idx >= len(r.schema.columns) {
		return nil
	}
	if r.schema.isQuestData {
		return r.getQuestBytesAt(idx)
	}
	col := r.schema.columns[idx]
	start := int(r.schema.fixedOffsets[idx])
	if col.Type == TypeListVarStr && idx == len(r.schema.columns)-1 {
		if start <= len(r.raw) {
			return r.raw[start:]
		}
		return nil
	}
	end := start + int(col.Size)
	if end <= len(r.raw) {
		return r.raw[start:end]
	}
	return nil
}

// GetBytes retorna el slice crudo de una columna por nombre (última aparición si el nombre se repite).
func (r *Row) GetBytes(colName string) []byte {
	if r == nil || r.schema == nil {
		return nil
	}
	idx, ok := r.schema.lastColIndex[colName]
	if !ok {
		return nil
	}
	return r.GetBytesAt(idx)
}

// GetByteAt retorna el byte en la columna posicional idx o 0.
func (r *Row) GetByteAt(idx int) byte {
	v := r.GetBytesAt(idx)
	if len(v) == 0 {
		return 0
	}
	return v[0]
}

// GetUint16At retorna uint16 en LittleEndian para la columna posicional idx.
func (r *Row) GetUint16At(idx int) uint16 {
	v := r.GetBytesAt(idx)
	if len(v) < 2 {
		return 0
	}
	return binary.LittleEndian.Uint16(v[:2])
}

// GetInt16At retorna int16 en LittleEndian para la columna posicional idx.
func (r *Row) GetInt16At(idx int) int16 {
	return int16(r.GetUint16At(idx))
}

// GetUint32At retorna uint32 en LittleEndian para la columna posicional idx.
func (r *Row) GetUint32At(idx int) uint32 {
	v := r.GetBytesAt(idx)
	if len(v) < 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(v[:4])
}

// GetInt32At retorna int32 en LittleEndian para la columna posicional idx.
func (r *Row) GetInt32At(idx int) int32 {
	return int32(r.GetUint32At(idx))
}

// GetIntAt retorna int para la columna posicional idx independientemente del ancho (1, 2 o 4 bytes).
func (r *Row) GetIntAt(idx int) int {
	v := r.GetBytesAt(idx)
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

// GetStringAt retorna la cadena ASCII/ANSI terminada en nulo de la columna posicional idx.
func (r *Row) GetStringAt(idx int) string {
	v := r.GetBytesAt(idx)
	if len(v) == 0 {
		return ""
	}
	if z := bytes.IndexByte(v, 0x00); z != -1 {
		return string(v[:z])
	}
	return string(v)
}

// ColumnIndex retorna el índice posicional de la primera columna con el nombre dado (-1 si no existe).
func (f *File) ColumnIndex(name string) int {
	if f == nil {
		return -1
	}
	if f.schema != nil {
		if idx, ok := f.schema.firstColIndex[name]; ok {
			return idx
		}
		return -1
	}
	for i, c := range f.Columns {
		if c.Name == name {
			return i
		}
	}
	return -1
}

// GetByte retorna el valor como byte o 0 si no existe.
func (r *Row) GetByte(colName string) byte {
	v := r.GetBytes(colName)
	if len(v) == 0 {
		return 0
	}
	return v[0]
}

// GetUint16 retorna el valor como uint16 en LittleEndian.
func (r *Row) GetUint16(colName string) uint16 {
	v := r.GetBytes(colName)
	if len(v) < 2 {
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
	v := r.GetBytes(colName)
	if len(v) < 4 {
		return 0
	}
	return binary.LittleEndian.Uint32(v[:4])
}

// GetInt retorna el valor numérico como int.
func (r *Row) GetInt(colName string) int {
	v := r.GetBytes(colName)
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
	v := r.GetBytes(colName)
	if len(v) < 8 {
		return int64(r.GetInt(colName))
	}
	return int64(binary.LittleEndian.Uint64(v[:8]))
}

// GetString retorna la cadena ASCII/ANSI terminada en nulo.
func (r *Row) GetString(colName string) string {
	v := r.GetBytes(colName)
	if len(v) == 0 {
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
	return r.GetBytes(colName)
}

// File representa un archivo de datos binario .shn completamente cargado en memoria.
type File struct {
	Version     string
	DataMode    uint32
	RecordCount int
	RowSize     int
	Columns     []Column
	Rows        []Row
	schema      *tableSchema
	md5Hash     Name8
	hasMD5      bool
}

// MD5Checksum retorna el hash MD5 canónico de 32 bytes (Name8) calculado sobre [cabecera 36B + cuerpo desencriptado].
func (f *File) MD5Checksum() (Name8, bool) {
	if f == nil || !f.hasMD5 {
		return Name8{}, false
	}
	return f.md5Hash, true
}

type cachedSHNEntry struct {
	size    int64
	modTime int64
	file    *File
}

var (
	sharedSHNCacheMu sync.RWMutex
	sharedSHNCache   = make(map[string]cachedSHNEntry)
)

// ClearReadCache limpia el registro compartido inmutable de archivos .shn en memoria.
func ClearReadCache() {
	sharedSHNCacheMu.Lock()
	sharedSHNCache = make(map[string]cachedSHNEntry)
	sharedSHNCacheMu.Unlock()
}

// ReadFileUncached carga y parsea un archivo .shn directamente desde disco sin consultar el cache compartido.
func ReadFileUncached(filePath string) (*File, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SHN file %s: %w", filePath, err)
	}
	return parseFile(raw, true)
}

// ReadFile carga y parsea un archivo .shn desde el disco reutilizando la instancia inmutable
// en el registro compartido si el archivo (ruta absoluta, tamaño y mtime) ya fue cargado.
func ReadFile(filePath string) (*File, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		absPath = filePath
	}

	fi, statErr := os.Stat(absPath)
	if statErr == nil {
		size := fi.Size()
		modTime := fi.ModTime().UnixNano()

		sharedSHNCacheMu.RLock()
		if entry, ok := sharedSHNCache[absPath]; ok && entry.size == size && entry.modTime == modTime {
			sharedSHNCacheMu.RUnlock()
			return entry.file, nil
		}
		sharedSHNCacheMu.RUnlock()

		f, err := ReadFileUncached(absPath)
		if err != nil {
			return nil, err
		}

		sharedSHNCacheMu.Lock()
		sharedSHNCache[absPath] = cachedSHNEntry{
			size:    size,
			modTime: modTime,
			file:    f,
		}
		sharedSHNCacheMu.Unlock()
		return f, nil
	}

	return ReadFileUncached(filePath)
}

// Parse parsea los bytes de un archivo .shn respetando la especificación CDataReader oficial
// y soportando el formato binario nativo de CQuestData (QuestData.shn), sin mutar el slice de entrada.
func Parse(data []byte) (*File, error) {
	return parseFile(data, false)
}

func computeMD5Name8(buf []byte) Name8 {
	sum := md5.Sum(buf)
	var name8 Name8
	hex.Encode(name8[:], sum[:])
	return name8
}

func parseFile(data []byte, owned bool) (*File, error) {
	if len(data) < 4 {
		return nil, errors.New("invalid SHN file: data is too small")
	}

	// Formato CQuestData::ReadModeRead de QuestData.shn (versión uint16 = 6)
	if binary.LittleEndian.Uint16(data[:2]) == 6 && len(data) >= 684 {
		count := int(binary.LittleEndian.Uint16(data[2:4]))
		offset := 4
		rows := make([]Row, 0, count)

		for i := 0; i < count; i++ {
			if offset+680 > len(data) {
				break
			}
			startOff := offset
			rec := data[offset : offset+680]
			offset += 680
			len1 := int(binary.LittleEndian.Uint16(rec[0x294:0x296]))
			len3 := int(binary.LittleEndian.Uint16(rec[0x296:0x298]))
			len2 := int(binary.LittleEndian.Uint16(rec[0x298:0x29A]))

			if offset+len1 <= len(data) {
				offset += len1
			}
			if offset+len2 <= len(data) {
				offset += len2
			}
			if offset+len3 <= len(data) {
				offset += len3
			}

			rows = append(rows, Row{
				schema: questDataSchema,
				raw:    data[startOff:offset],
			})
		}

		return &File{
			Version:     "6",
			DataMode:    0,
			RecordCount: len(rows),
			RowSize:     680,
			Columns:     questDataColumns,
			Rows:        rows,
			schema:      questDataSchema,
			md5Hash:     computeMD5Name8(data),
			hasMD5:      true,
		}, nil
	}

	if len(data) < 52 {
		return nil, errors.New("invalid SHN file: header is too small")
	}

	dataMode := binary.LittleEndian.Uint32(data[28:32])
	versionStr := string(bytes.TrimRight(data[4:28], "\x00"))

	var full []byte
	if dataMode == 1 {
		if owned {
			DecryptInPlace(data[36:])
			full = data
		} else {
			full = make([]byte, len(data))
			copy(full, data)
			DecryptInPlace(full[36:])
		}
	} else {
		full = data
	}

	md5Hash := computeMD5Name8(full)

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

	schema := newTableSchema(columns)
	hasVarCol := !schema.allFixed

	rows := make([]Row, 0, recordCount)
	rowOffset := headerEnd
	for r := 0; r < recordCount; r++ {
		if rowOffset+2 > len(full) {
			break
		}
		stepLen := rowSize
		if hasVarCol {
			stepLen = int(binary.LittleEndian.Uint16(full[rowOffset : rowOffset+2]))
			if stepLen < 2 || rowOffset+stepLen > len(full) {
				break
			}
		} else if rowOffset+stepLen > len(full) {
			break
		}

		rows = append(rows, Row{
			schema: schema,
			raw:    full[rowOffset : rowOffset+stepLen],
		})
		rowOffset += stepLen
	}

	return &File{
		Version:     versionStr,
		DataMode:    dataMode,
		RecordCount: recordCount,
		RowSize:     rowSize,
		Columns:     columns,
		Rows:        rows,
		schema:      schema,
		md5Hash:     md5Hash,
		hasMD5:      true,
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

// FindRowInt busca la primera fila cuyo valor numérico entero en colName sea igual a target.
func (f *File) FindRowInt(colName string, target int64) *Row {
	for i := range f.Rows {
		if int64(f.Rows[i].GetInt(colName)) == target {
			return &f.Rows[i]
		}
	}
	return nil
}
