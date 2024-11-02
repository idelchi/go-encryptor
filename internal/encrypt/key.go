package encrypt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// DecodeKey decodes a hex-encoded key string into bytes
func DecodeKey(keyStr string) ([]byte, error) {
	key, err := hex.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid hex key: %w", err)
	}

	// Validate key length
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: got %d bytes, want 32", len(key))
	}

	return key, nil
}

// EncodeKey encodes a key as a hex string
func EncodeKey(key []byte) string {
	return hex.EncodeToString(key)
}

// GenerateKey generates a 32-byte key for AES-256 encryption.
func GenerateKey() ([]byte, error) {
	const length = 32

	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("error generating key: %w", err)
	}
	return key, nil
}

// GenerateKeyString generates and returns a hex-encoded key
func GenerateKeyString() (string, error) {
	key, err := GenerateKey()
	if err != nil {
		return "", err
	}
	return EncodeKey(key), nil
}
