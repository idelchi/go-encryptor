package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// encryptBytes encrypts the given byte slice using AES-CFB mode.
// It prepends a random IV to the ciphertext and returns the complete encrypted block.
// The returned format is: [16 bytes IV][variable-length ciphertext].
func (e *Encryptor) encryptBytes(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	// Allocate space for IV and ciphertext in a single slice
	ciphertext := make([]byte, aes.BlockSize+len(data))
	initializationVector := ciphertext[:aes.BlockSize]

	// Generate random IV using crypto/rand
	if _, err := io.ReadFull(rand.Reader, initializationVector); err != nil {
		return nil, fmt.Errorf("generating IV: %w", err)
	}

	// Encrypt data using CFB mode
	stream := cipher.NewCFBEncrypter(block, initializationVector)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// decryptBytes decrypts the given ciphertext using AES-CFB mode.
// It expects the input to be in the format: [16 bytes IV][variable-length ciphertext].
// Returns the original plaintext on success.
func (e *Encryptor) decryptBytes(ciphertext []byte) ([]byte, error) {
	// Verify minimum length requirement for IV
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("%w: ciphertext too short", ErrProcessing)
	}

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	// Extract IV and actual ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt data using CFB mode
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext) // Decryption happens in-place

	return ciphertext, nil
}
