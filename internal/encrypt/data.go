package encrypt

import (
	"encoding/base64"
	"fmt"
)

// encryptData encrypts the given data and encodes it in base64.
// This is used for line-mode encryption where the output needs to be
// safely represented as a string in the output file.
func (e *Encryptor) encryptData(data []byte) ([]byte, error) {
	ciphertext, err := e.encryptBytes(data)
	if err != nil {
		return nil, err
	}

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// decryptData decodes the base64 data and decrypts it.
// This is used for line-mode decryption where the input is expected
// to be base64 encoded ciphertext.
func (e *Encryptor) decryptData(data []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("decoding base64: %w", err)
	}

	return e.decryptBytes(ciphertext)
}
