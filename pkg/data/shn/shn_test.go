package shn_test

import (
	"path/filepath"
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
