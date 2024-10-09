package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// encryptStream reads from the reader, encrypts the data, and writes to the writer.
func (e *Encryptor) encryptStream(reader io.Reader, writer io.Writer) error {
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return fmt.Errorf("creating cipher: %w", err)
	}

	var iv []byte
	if e.Type == Deterministic {
		iv = e.Key[:aes.BlockSize]
	} else {
		iv = make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return fmt.Errorf("generating IV: %w", err)
		}
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	base64Encoder := base64.NewEncoder(base64.StdEncoding, writer)
	defer base64Encoder.Close()

	// Write IV to the output (unencoded) if non-deterministic
	if e.Type == NonDeterministic {
		if _, err := writer.Write(iv); err != nil {
			return fmt.Errorf("writing IV: %w", err)
		}
	}

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

// decryptStream reads from the reader, decrypts the data, and writes to the writer.
func (e *Encryptor) decryptStream(reader io.Reader, writer io.Writer) error {
	var iv []byte
	if e.Type == Deterministic {
		iv = e.Key[:aes.BlockSize]
	} else {
		// Read IV from the input (if not deterministic)
		iv = make([]byte, aes.BlockSize)
		n, err := io.ReadFull(reader, iv)
		if err != nil {
			return fmt.Errorf("reading IV: %w", err)
		}
		if n < aes.BlockSize {
			return fmt.Errorf("IV too short")
		}
	}

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return fmt.Errorf("creating cipher: %w", err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	base64Decoder := base64.NewDecoder(base64.StdEncoding, reader)

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
