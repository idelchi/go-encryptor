package encrypt

import (
	"encoding/base64"
	"fmt"
)

// encryptData encrypts the given data and encodes it in base64.
func (e *Encryptor) encryptData(data []byte) ([]byte, error) {
	ciphertext, err := e.encryptBytes(data)
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// decryptData decodes the base64 data and decrypts it.
func (e *Encryptor) decryptData(data []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("decoding base64: %w", err)
	}
	return e.decryptBytes(ciphertext)
}
