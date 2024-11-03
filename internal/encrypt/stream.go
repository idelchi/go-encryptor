package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func (e *Encryptor) encryptStream(reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return fmt.Errorf("creating cipher: %w", err)
	}

	// Generate IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("generating IV: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	// Create base64 encoder first
	base64Encoder := base64.NewEncoder(base64.StdEncoding, writer)
	defer base64Encoder.Close()

	// Write IV through the base64 encoder
	if _, err := base64Encoder.Write(iv); err != nil {
		return fmt.Errorf("writing IV: %w", err)
	}

	// Read and encrypt data
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			encrypted := make([]byte, n)
			stream.XORKeyStream(encrypted, buf[:n])
			if _, err := base64Encoder.Write(encrypted); err != nil {
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
	// Create base64 decoder first
	base64Decoder := base64.NewDecoder(base64.StdEncoding, reader)

	// Read IV through the base64 decoder
	iv := make([]byte, aes.BlockSize)
	n, err := io.ReadFull(base64Decoder, iv)
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

	// Read and decrypt data
	buf := make([]byte, 4096)
	for {
		n, err := base64Decoder.Read(buf)
		if n > 0 {
			decrypted := make([]byte, n)
			stream.XORKeyStream(decrypted, buf[:n])
			if _, err := writer.Write(decrypted); err != nil {
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
