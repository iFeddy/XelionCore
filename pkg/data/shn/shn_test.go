package shn_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

func TestSHNParser_ClassName(t *testing.T) {
	path := filepath.Join("..", "..", "..", "data", "ClassName.shn")
	file, err := shn.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read ClassName.shn: %v", err)
	}

	if file.RecordCount != 28 {
		t.Errorf("Expected 28 records, got %d", file.RecordCount)
	}

	if len(file.Columns) != 4 {
		t.Fatalf("Expected 4 columns, got %d", len(file.Columns))
	}

	expectedCols := []string{"ClassID", "acPrefix", "acEngName", "acLocalName"}
	for i, name := range expectedCols {
		if file.Columns[i].Name != name {
			t.Errorf("Column %d: expected %s, got %s", i, name, file.Columns[i].Name)
		}
	}

	// Fila 1 (índice 1 en SHN): Fighter
	row1 := file.Rows[1]
	classID := row1.GetByte("ClassID")
	prefix := row1.GetString("acPrefix")
	engName := row1.GetString("acEngName")

	if classID != 1 {
		t.Errorf("Expected ClassID 1 (Fighter), got %d", classID)
	}
	if prefix != "Fig" {
		t.Errorf("Expected acPrefix 'Fig', got '%s'", prefix)
	}
	if engName != "Fighter" {
		t.Errorf("Expected acEngName 'Fighter', got '%s'", engName)
	}
}

func TestSHNParser_BadNameFilter(t *testing.T) {
	path := filepath.Join("..", "..", "..", "data", "BadNameFilter.shn")
	file, err := shn.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read BadNameFilter.shn: %v", err)
	}

	if file.RecordCount != 1175 {
		t.Errorf("Expected 1175 filtered names, got %d", file.RecordCount)
	}

	if len(file.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(file.Columns))
	}
}

func TestSHNParser_MapInfo(t *testing.T) {
	path := filepath.Join("..", "..", "..", "data", "MapInfo.shn")
	file, err := shn.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read MapInfo.shn: %v", err)
	}

	// Buscar mapa inicial 'RouN' (Roumen)
	rounRow := file.FindRow("MapName", "RouN")
	if rounRow == nil {
		t.Fatalf("Could not find 'RouN' map in MapInfo.shn")
	}

	mapID := rounRow.GetInt("ID")
	if mapID == 0 {
		t.Errorf("Expected non-zero ID for RouN, got %d", mapID)
	}
}

func TestSHNParser_SingleData(t *testing.T) {
	path := filepath.Join("..", "..", "..", "data", "SingleData.shn")
	file, err := shn.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read SingleData.shn: %v", err)
	}

	if file.RecordCount != 41 {
		t.Errorf("Expected 41 records, got %d", file.RecordCount)
	}

	voteLimitRow := file.FindRow("SingleDataIDX", "KQVote_VoteLimitTime")
	if voteLimitRow == nil {
		t.Fatalf("Could not find KQVote_VoteLimitTime row in SingleData.shn")
	}

	val := voteLimitRow.GetUint16("SingleDataValue")
	if val != 60 {
		t.Errorf("Expected KQVote_VoteLimitTime = 60, got %d", val)
	}
}

func TestSHNParser_CharacterAppearance(t *testing.T) {
	// FaceInfo.shn
	facePath := filepath.Join("..", "..", "..", "data", "FaceInfo.shn")
	faceFile, err := shn.ReadFile(facePath)
	if err != nil {
		t.Fatalf("Failed to read FaceInfo.shn: %v", err)
	}
	if faceFile.RecordCount != 48 {
		t.Errorf("Expected 48 face records, got %d", faceFile.RecordCount)
	}
	for i := 0; i < 5 && i < len(faceFile.Rows); i++ {
		t.Logf("FaceRow %d: ID=%d, Name=%s, Grade=%d", i, faceFile.Rows[i].GetByte("ID"), faceFile.Rows[i].GetString("FaceName"), faceFile.Rows[i].GetByte("Grade"))
	}

	// HairInfo.shn
	hairPath := filepath.Join("..", "..", "..", "data", "HairInfo.shn")
	hairFile, err := shn.ReadFile(hairPath)
	if err != nil {
		t.Fatalf("Failed to read HairInfo.shn: %v", err)
	}
	if hairFile.RecordCount != 70 {
		t.Errorf("Expected 70 hair records, got %d", hairFile.RecordCount)
	}
	for _, col := range hairFile.Columns {
		t.Logf("HairInfo col: %s (%d)", col.Name, col.Type)
	}

	// HairColorInfo.shn
	colorPath := filepath.Join("..", "..", "..", "data", "HairColorInfo.shn")
	colorFile, err := shn.ReadFile(colorPath)
	if err != nil {
		t.Fatalf("Failed to read HairColorInfo.shn: %v", err)
	}
	for _, col := range colorFile.Columns {
		t.Logf("HairColorInfo col: %s (%d)", col.Name, col.Type)
	}
}


type legacySHNRow struct {
	values map[string][]byte
	cols   [][]byte
}

type legacySHNFile struct {
	recordCount int
	rowSize     int
	columns     []shn.Column
	rows        []legacySHNRow
}

// parseSHNLegacyReference replica el parser original (triple copia + map[string][]byte + [][]byte)
// exclusivamente como oráculo de paridad binaria 1:1 para todas las tablas .shn del juego.
func parseSHNLegacyReference(data []byte) *legacySHNFile {
	if binary.LittleEndian.Uint16(data[:2]) == 6 && len(data) >= 684 {
		count := int(binary.LittleEndian.Uint16(data[2:4]))
		offset := 4
		cols := []shn.Column{
			{Name: "RecordSize", Type: shn.TypeListWord, Size: 2},
			{Name: "QuestID", Type: shn.TypeListDword, Size: 4},
			{Name: "TotalSize", Type: shn.TypeListDword, Size: 4},
			{Name: "HeaderData", Type: shn.TypeListRemark, Size: 680},
			{Name: "ScriptStart", Type: shn.TypeListStrAuto, Size: 0},
			{Name: "ScriptDoing", Type: shn.TypeListStrAuto, Size: 0},
			{Name: "ScriptEnd", Type: shn.TypeListStrAuto, Size: 0},
		}
		rows := make([]legacySHNRow, 0, count)
		for i := 0; i < count; i++ {
			if offset+680 > len(data) {
				break
			}
			rec := data[offset : offset+680]
			offset += 680
			len1 := int(binary.LittleEndian.Uint16(rec[0x294:0x296]))
			len3 := int(binary.LittleEndian.Uint16(rec[0x296:0x298]))
			len2 := int(binary.LittleEndian.Uint16(rec[0x298:0x29A]))
			var s1, s2, s3 []byte
			if offset+len1 <= len(data) {
				s1 = data[offset : offset+len1]
				offset += len1
			}
			if offset+len2 <= len(data) {
				s2 = data[offset : offset+len2]
				offset += len2
			}
			if offset+len3 <= len(data) {
				s3 = data[offset : offset+len3]
				offset += len3
			}
			totalSize := uint32(680 + len1 + len2 + len3)
			totalBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(totalBytes, totalSize)
			vals := map[string][]byte{
				"RecordSize":  rec[:2],
				"QuestID":     rec[4:8],
				"TotalSize":   totalBytes,
				"HeaderData":  rec,
				"ScriptStart": s1,
				"ScriptEnd":   s2,
				"ScriptDoing": s3,
			}
			rowCols := [][]byte{rec[:2], rec[4:8], totalBytes, rec, s1, s3, s2}
			rows = append(rows, legacySHNRow{values: vals, cols: rowCols})
		}
		return &legacySHNFile{recordCount: len(rows), rowSize: 680, columns: cols, rows: rows}
	}

	dataMode := binary.LittleEndian.Uint32(data[28:32])
	var full []byte
	if dataMode == 1 {
		dec := shn.Decrypt(data[36:])
		full = make([]byte, 36+len(dec))
		copy(full[:36], data[:36])
		copy(full[36:], dec)
	} else {
		full = data
	}

	recordCount := int(binary.LittleEndian.Uint32(full[40:44]))
	rowSize := int(binary.LittleEndian.Uint32(full[44:48]))
	colCount := int(binary.LittleEndian.Uint32(full[48:52]))
	headerEnd := 52 + colCount*56
	columns := make([]shn.Column, colCount)
	colOffset := 52
	hasVarCol := false
	for i := 0; i < colCount; i++ {
		desc := full[colOffset : colOffset+56]
		columns[i] = shn.Column{
			Name: string(bytes.TrimRight(desc[:48], "\x00")),
			Type: binary.LittleEndian.Uint32(desc[48:52]),
			Size: binary.LittleEndian.Uint32(desc[52:56]),
		}
		if columns[i].Type == shn.TypeListVarStr {
			hasVarCol = true
		}
		colOffset += 56
	}

	rows := make([]legacySHNRow, 0, recordCount)
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
		rowBytes := full[rowOffset : rowOffset+stepLen]
		rowVals := make(map[string][]byte, colCount)
		rowCols := make([][]byte, colCount)
		curr := 2
		for idx, col := range columns {
			if col.Type == shn.TypeListVarStr && idx == len(columns)-1 {
				if curr <= len(rowBytes) {
					val := make([]byte, len(rowBytes)-curr)
					copy(val, rowBytes[curr:])
					rowCols[idx] = val
					rowVals[col.Name] = val
				}
				curr = len(rowBytes)
				continue
			}
			end := curr + int(col.Size)
			if end <= len(rowBytes) {
				val := make([]byte, col.Size)
				copy(val, rowBytes[curr:end])
				rowCols[idx] = val
				rowVals[col.Name] = val
			}
			curr = end
		}
		rows = append(rows, legacySHNRow{values: rowVals, cols: rowCols})
		rowOffset += stepLen
	}
	return &legacySHNFile{recordCount: recordCount, rowSize: rowSize, columns: columns, rows: rows}
}

func TestSHN_BinaryParity_AllFiles(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	dirs := []string{dataDir, filepath.Join(dataDir, "View")}

	var shnPaths []string
	for _, d := range dirs {
		entries, err := os.ReadDir(d)
		if err != nil {
			t.Fatalf("failed to read dir %s: %v", d, err)
		}
		for _, e := range entries {
			if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".shn") {
				shnPaths = append(shnPaths, filepath.Join(d, e.Name()))
			}
		}
	}

	if len(shnPaths) < 200 {
		t.Fatalf("expected >= 200 SHN files in data/ and data/View/, got %d", len(shnPaths))
	}

	var totalRows, repeatedColFiles int
	var verifiedQuestData bool

	for _, p := range shnPaths {
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("ReadFile %s: %v", p, err)
		}
		ref := parseSHNLegacyReference(raw)
		opt, err := shn.ReadFileUncached(p)
		if err != nil {
			t.Fatalf("shn.ReadFileUncached(%s): %v", p, err)
		}

		if opt.RecordCount != ref.recordCount || opt.RowSize != ref.rowSize || len(opt.Columns) != len(ref.columns) {
			t.Fatalf("%s metadata mismatch: got (%d,%d,%d), want (%d,%d,%d)",
				filepath.Base(p), opt.RecordCount, opt.RowSize, len(opt.Columns),
				ref.recordCount, ref.rowSize, len(ref.columns))
		}
		if len(opt.Rows) != len(ref.rows) {
			t.Fatalf("%s row count mismatch: got %d, want %d", filepath.Base(p), len(opt.Rows), len(ref.rows))
		}

		// Verificar ColumnIndex (firstColIndex) y columnas repetidas
		seenCols := make(map[string]int, len(ref.columns))
		hasRepeated := false
		for i, c := range ref.columns {
			if opt.Columns[i] != c {
				t.Fatalf("%s col %d mismatch: got %+v, want %+v", filepath.Base(p), i, opt.Columns[i], c)
			}
			if _, ok := seenCols[c.Name]; !ok {
				seenCols[c.Name] = i
				if got := opt.ColumnIndex(c.Name); got != i {
					t.Fatalf("%s ColumnIndex(%q): got %d, want first=%d", filepath.Base(p), c.Name, got, i)
				}
			} else {
				hasRepeated = true
			}
		}
		if hasRepeated {
			repeatedColFiles++
		}

		for r := range ref.rows {
			optRow := &opt.Rows[r]
			refRow := &ref.rows[r]
			if optRow.ColumnCount() != len(ref.columns) {
				t.Fatalf("%s row %d ColumnCount=%d, want %d", filepath.Base(p), r, optRow.ColumnCount(), len(ref.columns))
			}
			if len(optRow.RawRecord()) == 0 {
				t.Fatalf("%s row %d RawRecord is empty", filepath.Base(p), r)
			}
			for c := range ref.columns {
				got := optRow.GetBytesAt(c)
				want := refRow.cols[c]
				if !bytes.Equal(got, want) {
					t.Fatalf("%s row %d col %d (%s) mismatch: got %x, want %x",
						filepath.Base(p), r, c, ref.columns[c].Name, got, want)
				}
			}
			for colName, want := range refRow.values {
				got := optRow.GetBytes(colName)
				if !bytes.Equal(got, want) {
					t.Fatalf("%s row %d GetBytes(%q) mismatch: got %x, want %x",
						filepath.Base(p), r, colName, got, want)
				}
			}
		}

		if strings.EqualFold(filepath.Base(p), "QuestData.shn") {
			verifiedQuestData = true
			if len(opt.Rows) == 0 || opt.Rows[0].ColumnCount() != 7 {
				t.Fatalf("QuestData.shn invalid parsed layout: rows=%d", len(opt.Rows))
			}
		}
		totalRows += len(opt.Rows)
	}

	if !verifiedQuestData {
		t.Fatalf("QuestData.shn was not verified")
	}
	t.Logf("Verified 1:1 binary parity across %d SHN files (%d rows, %d files with repeated columns, QuestData v6 verified).",
		len(shnPaths), totalRows, repeatedColFiles)
}

func TestChecksumEngine_FromFiles_Parity(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	viewDir := filepath.Join(dataDir, "View")

	files := make(map[string]*shn.File)
	for _, d := range []string{dataDir, viewDir} {
		entries, err := os.ReadDir(d)
		if err != nil {
			t.Fatalf("ReadDir %s: %v", d, err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".shn") {
				continue
			}
			f, err := shn.ReadFile(filepath.Join(d, e.Name()))
			if err != nil {
				t.Fatalf("ReadFile %s: %v", e.Name(), err)
			}
			base := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
			files[base] = f
			files[strings.ToLower(base)] = f
		}
	}

	ref, err := shn.NewChecksumEngine(dataDir, viewDir)
	if err != nil {
		t.Fatalf("NewChecksumEngine: %v", err)
	}
	opt, err := shn.NewChecksumEngineFromFiles(files, dataDir, viewDir)
	if err != nil {
		t.Fatalf("NewChecksumEngineFromFiles: %v", err)
	}

	for i := 0; i < 49; i++ {
		if opt.GetRawHash(i) != ref.GetRawHash(i) {
			t.Fatalf("Checksum mismatch at index %d (%s): got %s, want %s",
				i, shn.ChecksumFileList[i], opt.GetRawHash(i), ref.GetRawHash(i))
		}
	}
	if ok, badIdx := opt.VerifyDataFiles(ref.GetChecksums()); !ok {
		t.Fatalf("VerifyDataFiles failed at index %d", badIdx)
	}
}
