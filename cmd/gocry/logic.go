package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/idelchi/go-encryptor/internal/encrypt"
	"github.com/idelchi/go-encryptor/pg/stdin"
)

func run(cfg Config) error {
	key, err := os.ReadFile(cfg.Key)
	if err != nil {
		return fmt.Errorf("reading key file: %w", err)
	}

	if cfg.GPG {
		key, err = deriveKeyFromGPG(string(key))
		if err != nil {
			return fmt.Errorf("deriving key from GPG: %w", err)
		}
	}

	data, err := loadData(cfg.File)
	if err != nil {
		return fmt.Errorf("deriving key from GPG: %w", err)
	}

	// Use os.Stdout as the writer
	processed, err := encrypt.Process(cfg.Mode, cfg.Operation, cfg.Type, key, data, os.Stdout)
	if err != nil {
		fmt.Errorf("error processing data: %w", err)

		return fmt.Errorf("processing data: %w", err)
	}
	if cfg.Mode == "file" {
		fmt.Fprintf(os.Stderr, "%sed file: %q\n", cfg.Operation, cfg.File)
	}

	if cfg.Mode == "line" && processed {
		fmt.Fprintf(os.Stderr, "%sed lines in: %q\n", cfg.Operation, cfg.File)
	}

	return nil
}

func loadData(file string) (data *os.File, err error) {
	if stdin.IsPiped() {
		data = os.Stdin
	} else {
		// Open the input file
		data, err = os.Open(file)
		if err != nil {
			fmt.Errorf("opening input file %q: %w", file, err)
			os.Exit(1)
		}
		defer data.Close()
	}

	return data, nil
}

func deriveKeyFromGPG(gpgKey string) ([]byte, error) {
	gpgKeyDecoded, err := base64.StdEncoding.DecodeString(gpgKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 gpg key: %w", err)
	}

	// Use SHA-256 to derive a 32-byte key for AES-256
	hash := sha256.Sum256(gpgKeyDecoded)
	return hash[:], nil
}
