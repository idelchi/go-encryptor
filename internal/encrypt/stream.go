package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func (e *Encryptor) encryptStream(reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return fmt.Errorf("creating cipher: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("generating IV: %w", err)
	}

	// Write IV directly
	if _, err := writer.Write(iv); err != nil {
		return fmt.Errorf("writing IV: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	buf := make([]byte, 4096)
	encrypted := make([]byte, 4096) // Pre-allocate encryption buffer

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

func (e *Encryptor) decryptStream(reader io.Reader, writer io.Writer) error {
	// Read IV directly (no base64 decoder)
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
	buf := make([]byte, 4096)
	decrypted := make([]byte, 4096) // Pre-allocate decryption buffer

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
