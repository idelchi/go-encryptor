// Package encrypt provides functions to encrypt and decrypt data, both for entire files and on a per-line basis.
package encrypt

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// Operation represents the encryption or decryption operation.
type Operation string

const (
	// Encrypt operation.
	Encrypt Operation = "encrypt"
	// Decrypt operation.
	Decrypt Operation = "decrypt"
)

// Type represents whether encryption is deterministic or not.
type Type string

const (
	// Deterministic encryption uses a fixed IV derived from the key.
	Deterministic Type = "deterministic"
	// NonDeterministic encryption uses a random IV.
	NonDeterministic Type = "nondeterministic"
)

// Mode represents the mode of operation.
type Mode string

const (
	// Line mode processes each line of the input data.
	Line Mode = "line"
	// File mode processes the entire input data as a single block.
	File Mode = "file"
)

// Encryptor handles encryption and decryption operations.
type Encryptor struct {
	Key       []byte
	Operation Operation
	Mode      Mode
	Type      Type
}

// Process handles encryption and decryption based on the provided configuration.
// It delegates to either processLines or processWholeFile depending on the mode.
func (e *Encryptor) Process(reader io.Reader, writer io.Writer) (bool, error) {
	switch e.Mode {
	case Line:
		return e.processLines(reader, writer)
	case File:
		return e.processWholeFile(reader, writer)
	default:
		return false, fmt.Errorf("invalid mode: %s", e.Type)
	}
}

// processLines processes each line of the input data, encrypting or decrypting lines
// that contain the specific directive.
func (e *Encryptor) processLines(reader io.Reader, writer io.Writer) (bool, error) {
	var processed bool

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case e.Operation == Encrypt && strings.HasSuffix(line, "### DIRECTIVE: ENCRYPT"):
			encryptedLine, err := e.encryptData([]byte(line))
			if err != nil {
				return processed, err
			}

			processed = true
			fmt.Fprintf(writer, "### DIRECTIVE: DECRYPT: %s\n", encryptedLine)
		case e.Operation == Decrypt && strings.HasPrefix(line, "### DIRECTIVE: DECRYPT: "):
			encryptedData := strings.TrimPrefix(line, "### DIRECTIVE: DECRYPT: ")
			decryptedLine, err := e.decryptData([]byte(encryptedData))
			if err != nil {
				return processed, err
			}

			processed = true
			fmt.Fprintln(writer, string(decryptedLine))
		default:
			fmt.Fprintln(writer, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return processed, fmt.Errorf("scanning error: %v", err)
	}
	return processed, nil
}

// processWholeFile processes the entire input data as a single encrypted or decrypted block.
func (e *Encryptor) processWholeFile(reader io.Reader, writer io.Writer) (bool, error) {
	switch e.Operation {
	case Encrypt:
		return true, e.encryptStream(reader, writer)
	case Decrypt:
		return true, e.decryptStream(reader, writer)
	default:
		return false, fmt.Errorf("invalid operation")
	}
}

// encryptData encrypts the given data using AES in CFB mode.
func (e *Encryptor) encryptData(data []byte) ([]byte, error) {
	ciphertext, err := e.encryptBytes(data)
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// decryptData decrypts the given data using AES in CFB mode.
func (e *Encryptor) decryptData(data []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("decoding base64: %w", err)
	}
	return e.decryptBytes(ciphertext)
}

// encryptBytes encrypts the given byte slice and returns the ciphertext with IV prepended.
func (e *Encryptor) encryptBytes(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	if e.Type == Deterministic {
		copy(iv, e.Key[:aes.BlockSize])
	} else {
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, fmt.Errorf("generating IV: %w", err)
		}
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// decryptBytes decrypts the given ciphertext (with IV prepended) and returns the plaintext.
func (e *Encryptor) decryptBytes(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

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
