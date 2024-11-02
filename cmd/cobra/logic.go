package main

import (
	"fmt"
	"os"

	"github.com/idelchi/go-next-tag/pkg/stdin"
	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/internal/printer"
	"github.com/idelchi/gocry/pkg/key"
)

func processFiles(cfg *Config) error {
	var encryptionKey []byte
	var err error

	switch {
	case cfg.Key != "":
		encryptionKey, err = key.Decode([]byte(cfg.Key))
	case cfg.KeyFile != "":
		encryptionKey, err = os.ReadFile(cfg.KeyFile)
		if err != nil {
			return fmt.Errorf("reading key file: %w", err)
		}

		encryptionKey, err = key.Decode(encryptionKey)
	}

	// Validate key length
	if len(encryptionKey) != 32 {
		return fmt.Errorf("invalid key length: got %d bytes, want 32", len(encryptionKey))
	}

	if err != nil {
		return fmt.Errorf("reading key: %w", err)
	}

	data, err := loadData(cfg.File)
	if err != nil {
		return fmt.Errorf("loading data: %w", err)
	}
	defer data.Close()

	encryptor := &encrypt.Encryptor{
		Key:        encryptionKey,
		Operation:  encrypt.Operation(cfg.Operation),
		Mode:       encrypt.Mode(cfg.Mode),
		Directives: cfg.Directives,
	}

	processed, err := encryptor.Process(data, os.Stdout)
	if err != nil {
		return fmt.Errorf("processing data: %w", err)
	}

	if cfg.Mode == "file" {
		printer.Stderrln("\n%sed file: %q", cfg.Operation, cfg.File)
	}

	if cfg.Mode == "line" && processed {
		printer.Stderrln("\n%sed lines in: %q", cfg.Operation, cfg.File)
	}

	return nil
}

func loadData(file string) (*os.File, error) {
	if stdin.IsPiped() {
		return os.Stdin, nil
	}

	data, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("opening input file %q: %w", file, err)
	}

	return data, nil
}
