package shn_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

func TestChecksumEngine_Verification(t *testing.T) {
	dataDir := filepath.Join("..", "..", "..", "data")
	viewDir := filepath.Join(dataDir, "View")
	engine, err := shn.NewChecksumEngine(dataDir, viewDir)
	if err != nil {
		t.Fatalf("Failed to initialize ChecksumEngine: %v", err)
	}

	checksums := engine.GetChecksums()

	// Verificar que el hash de Abstate.shn coincide con el hash del cliente oficial desencriptado
	expectedAbstateHash := "2667922d726861481bab726e9f1dd76c"
	if engine.GetRawHash(0) != expectedAbstateHash {
		t.Fatalf("Expected Abstate.shn hash to be %s, got %s", expectedAbstateHash, engine.GetRawHash(0))
	}

	// 1. Verificación exitosa con los 49 hashes legítimos
	valid, wrongIndex := engine.VerifyDataFiles(checksums)
	if !valid {
		t.Fatalf("Expected valid checksums, but failed at index %d (%s)", wrongIndex, engine.GetFileName(wrongIndex))
	}

	// 2. Simular alteración del archivo Abstate.shn (índice 0)
	tamperedChecksums := checksums
	tamperedChecksums[0][0] ^= 0xFF

	valid, wrongIndex = engine.VerifyDataFiles(tamperedChecksums)
	if valid {
		t.Fatalf("Expected verification to fail for tampered Abstate.shn, but passed")
	}
	if wrongIndex != 0 {
		t.Errorf("Expected wrong index 0, got %d (%s)", wrongIndex, engine.GetFileName(wrongIndex))
	}
	if engine.GetFileName(wrongIndex) != "Abstate.shn" {
		t.Errorf("Expected file name 'Abstate.shn', got '%s'", engine.GetFileName(wrongIndex))
	}

	// 3. Simular alteración del archivo ChargedDeletableBuff.shn (índice 48)
	tamperedChecksums = checksums
	tamperedChecksums[48][5] ^= 0xAA

	valid, wrongIndex = engine.VerifyDataFiles(tamperedChecksums)
	if valid {
		t.Fatalf("Expected verification to fail for tampered ChargedDeletableBuff.shn, but passed")
	}
	if wrongIndex != 48 {
		t.Errorf("Expected wrong index 48, got %d (%s)", wrongIndex, engine.GetFileName(wrongIndex))
	}
	if engine.GetFileName(wrongIndex) != "ChargedDeletableBuff.shn" {
		t.Errorf("Expected file name 'ChargedDeletableBuff.shn', got '%s'", engine.GetFileName(wrongIndex))
	}
}
