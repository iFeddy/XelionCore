package shn_test

import (
	"path/filepath"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

func TestChecksumEngine_Verification(t *testing.T) {
	shineDir := filepath.Join("..", "..", "..", "data")
	engine, err := shn.NewChecksumEngine(shineDir)
	if err != nil {
		t.Fatalf("Failed to initialize ChecksumEngine: %v", err)
	}

	checksums := engine.GetChecksums()

	// 1. Verificación exitosa con los 49 hashes legítimos
	valid, wrongIndex := engine.VerifyDataFiles(checksums)
	if !valid {
		t.Fatalf("Expected valid checksums, but failed at index %d (%s)", wrongIndex, engine.GetFileName(wrongIndex))
	}

	// 2. Simular alteración del archivo ActiveSkill.shn (índice 47)
	tamperedChecksums := checksums
	tamperedChecksums[47][0] ^= 0xFF // Corromper un byte

	valid, wrongIndex = engine.VerifyDataFiles(tamperedChecksums)
	if valid {
		t.Fatalf("Expected verification to fail for tampered ActiveSkill.shn, but passed")
	}
	if wrongIndex != 47 {
		t.Errorf("Expected wrong index 47, got %d (%s)", wrongIndex, engine.GetFileName(wrongIndex))
	}
	if engine.GetFileName(wrongIndex) != "ActiveSkill.shn" {
		t.Errorf("Expected file name 'ActiveSkill.shn', got '%s'", engine.GetFileName(wrongIndex))
	}

	// 3. Simular alteración del archivo ChargedDeletableBuff.shn (índice 0)
	tamperedChecksums = checksums
	tamperedChecksums[0][5] ^= 0xAA

	valid, wrongIndex = engine.VerifyDataFiles(tamperedChecksums)
	if valid {
		t.Fatalf("Expected verification to fail for tampered ChargedDeletableBuff.shn, but passed")
	}
	if wrongIndex != 0 {
		t.Errorf("Expected wrong index 0, got %d", wrongIndex)
	}
}
