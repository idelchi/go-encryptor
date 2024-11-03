package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// encryptStream encrypts data from reader to writer using AES-CFB mode.
// It prepends the randomly generated IV to the encrypted output.
// The encryption is done in chunks to maintain constant memory usage.
func (e *Encryptor) encryptStream(reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return fmt.Errorf("creating cipher: %w", err)
	}

	// Generate a random IV (Initialization Vector)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("generating IV: %w", err)
	}

	// Write IV directly
	if _, err := writer.Write(iv); err != nil {
		return fmt.Errorf("writing IV: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	// Use fixed-size buffers for reading and encryption
	buf := make([]byte, 4096)
	encrypted := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			stream.XORKeyStream(encrypted[:n], buf[:n])
			if _, err := writer.Write(encrypted[:n]); err != nil {
				return fmt.Errorf("writing encrypted data: %w", err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading data: %w", err)
		}
	}
	return nil
}

// decryptStream decrypts data from reader to writer using AES-CFB mode.
// It expects the IV to be prepended to the encrypted data.
// The decryption is done in chunks to maintain constant memory usage.
func (e *Encryptor) decryptStream(reader io.Reader, writer io.Writer) error {
	// Read the prepended IV
	iv := make([]byte, aes.BlockSize)
	n, err := io.ReadFull(reader, iv)
	if err != nil {
		return fmt.Errorf("reading IV: %w", err)
	}
	if n < aes.BlockSize {
		return fmt.Errorf("IV too short")
	}

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return fmt.Errorf("creating cipher: %w", err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	// Use fixed-size buffers for reading and decryption
	buf := make([]byte, 4096)
	decrypted := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			stream.XORKeyStream(decrypted[:n], buf[:n])
			if _, err := writer.Write(decrypted[:n]); err != nil {
				return fmt.Errorf("writing decrypted data: %w", err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading encrypted data: %w", err)
		}
	}
	return nil
}
