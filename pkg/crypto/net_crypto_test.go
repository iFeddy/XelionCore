package crypto

import (
	"bytes"
	"testing"
)

func TestNetCryptoSymmetry(t *testing.T) {
	initialOffset := int16(123)
	c1, err := NewNetCrypto(initialOffset)
	if err != nil {
		t.Fatalf("unexpected error creating NetCrypto: %v", err)
	}
	c2, err := NewNetCrypto(initialOffset)
	if err != nil {
		t.Fatalf("unexpected error creating NetCrypto: %v", err)
	}

	original := []byte("Hello Fiesta Online Login Server!")
	encrypted := make([]byte, len(original))
	copy(encrypted, original)

	c1.Crypt(encrypted)
	if bytes.Equal(encrypted, original) {
		t.Fatalf("encrypted buffer should differ from original")
	}

	decrypted := make([]byte, len(encrypted))
	copy(decrypted, encrypted)

	c2.Crypt(decrypted)
	if !bytes.Equal(decrypted, original) {
		t.Fatalf("decrypted buffer '%s' does not match original '%s'", decrypted, original)
	}
}

func TestNetCryptoWrapAround(t *testing.T) {
	c, err := NewNetCrypto(TableSize - 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf := []byte{0x01, 0x02}
	c.Crypt(buf)

	if c.XorPos() != 1 {
		t.Fatalf("expected XorPos to wrap around to 1, got %d", c.XorPos())
	}
}
