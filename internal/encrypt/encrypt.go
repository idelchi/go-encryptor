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

// Process handles encryption and decryption based on the provided configuration.
// It delegates to either processLines or processWholeFile depending on the mode.
func Process(mode, operation, encryption string, key []byte, reader io.Reader, writer io.Writer) (bool, error) {
	switch mode {
	case "line":
		return processLines(operation, encryption, key, reader, writer)
	case "file":
		return processWholeFile(operation, encryption, key, reader, writer)
	default:
		return false, fmt.Errorf("invalid mode: %s", mode)
	}
}

// processLines processes each line of the input data, encrypting or decrypting lines
// that contain the specific directive.
func processLines(operation, encryption string, key []byte, reader io.Reader, writer io.Writer) (bool, error) {
	var processed bool

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case operation == "encrypt" && strings.HasSuffix(line, "### DIRECTIVE: ENCRYPT"):
			encryptedLine, err := encryptData([]byte(line), key, encryption == "deterministic")
			if err != nil {
				return processed, err
			}

			processed = true

			fmt.Fprintf(writer, "### DIRECTIVE: DECRYPT: %s\n", encryptedLine)
		case operation == "decrypt" && strings.HasPrefix(line, "### DIRECTIVE: DECRYPT: "):

			encryptedData := strings.TrimPrefix(line, "### DIRECTIVE: DECRYPT: ")
			decryptedLine, err := decryptData([]byte(encryptedData), key)
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
func processWholeFile(operation, encryption string, key []byte, reader io.Reader, writer io.Writer) (bool, error) {
	switch operation {
	case "encrypt":
		return true, encryptStream(reader, writer, key, encryption == "deterministic")
	case "decrypt":
		return true, decryptStream(reader, writer, key)
	default:
		return false, fmt.Errorf("invalid operation: %s", operation)
	}
}

// encryptData encrypts the given data using AES in CFB mode.
// If deterministic is true, it uses a deterministic IV derived from the key;
// otherwise, it uses a randomly generated IV.
func encryptData(data, key []byte, deterministic bool) ([]byte, error) {
	ciphertext, err := encryptBytes(data, key, deterministic)
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// decryptData decrypts the given data using AES in CFB mode.
// It expects the IV to be prepended to the ciphertext.
func decryptData(data, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("decoding base64: %w", err)
	}
	return decryptBytes(ciphertext, key)
}

// encryptBytes encrypts the given byte slice and returns the ciphertext with IV prepended.
func encryptBytes(data, key []byte, deterministic bool) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]

	if deterministic {
		copy(iv, key[:aes.BlockSize])
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
func decryptBytes(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
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
func encryptStream(reader io.Reader, writer io.Writer, key []byte, deterministic bool) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("creating cipher: %w", err)
	}

	var iv []byte
	if deterministic {
		iv = key[:aes.BlockSize]
	} else {
		iv = make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return fmt.Errorf("generating IV: %w", err)
		}
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	base64Encoder := base64.NewEncoder(base64.StdEncoding, writer)
	defer base64Encoder.Close()

	// Write IV to the output (unencoded)
	if !deterministic {
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
func decryptStream(reader io.Reader, writer io.Writer, key []byte) error {
	// Read IV from the input (if not deterministic)
	iv := make([]byte, aes.BlockSize)
	n, err := io.ReadFull(reader, iv)
	if err != nil {
		return fmt.Errorf("reading IV: %w", err)
	}
	if n < aes.BlockSize {
		return fmt.Errorf("IV too short")
	}

	block, err := aes.NewCipher(key)
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
