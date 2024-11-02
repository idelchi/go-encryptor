package key

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// Decode decodes a hex-encoded key string into bytes
func Decode(hexKey []byte) ([]byte, error) {
	hexString := strings.TrimSpace(string(hexKey))

	key, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, fmt.Errorf("invalid hex key: %w", err)
	}

	return key, nil
}

// Encode encodes a key as a hex string
func Encode(key []byte) string {
	return hex.EncodeToString(key)
}

// generate generates a 32-byte key for AES-256 encryption.
func generate(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("error generating key: %w", err)
	}
	return key, nil
}

// GenerateString generates and returns a hex-encoded key
func GenerateHex(length int) (string, error) {
	key, err := generate(length)
	if err != nil {
		return "", err
	}

	return Encode(key), nil
}
