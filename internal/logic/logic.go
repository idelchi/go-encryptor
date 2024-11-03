// Package logic implements the core business logic for the encryption/decryption service.
package logic

import (
	"fmt"
	"os"

	"github.com/idelchi/gocry/internal/config"
	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/internal/printer"
	"github.com/idelchi/gogen/pkg/key"
)

// Run executes the main encryption/decryption logic based on the provided configuration.
// It handles key loading, input data loading, and processes the data according to the
// specified mode and operation.
func Run(cfg *config.Config) error {
	var (
		encryptionKey []byte
		err           error
	)

	// Load encryption key either from hex string or file
	switch {
	case cfg.Key.String != "":
		encryptionKey, err = key.FromHex(cfg.Key.String)
	case cfg.Key.File != "":
		encryptionKey, err = os.ReadFile(cfg.Key.File)
		if err != nil {
			return fmt.Errorf("reading key file: %w", err)
		}

		encryptionKey, err = key.FromHex(string(encryptionKey))
	}

	if err != nil {
		return fmt.Errorf("reading key: %w", err)
	}

	// Ensure key meets AES-256 requirement
	if len(encryptionKey) != 32 {
		return fmt.Errorf("invalid key length: got %d bytes, want 32", len(encryptionKey))
	}

	// Load input data from stdin or file
	data, err := loadData(cfg.File)
	if err != nil {
		return fmt.Errorf("loading data: %w", err)
	}
	defer data.Close()

	// Initialize encryptor with configuration
	encryptor := &encrypt.Encryptor{
		Key:        encryptionKey,
		Operation:  encrypt.Operation(cfg.Operation),
		Mode:       encrypt.Mode(cfg.Mode),
		Directives: cfg.Directives,
		Parallel:   cfg.Parallel,
	}

	// Process data and handle any errors
	processed, err := encryptor.Process(data, os.Stdout)
	if err != nil {
		return fmt.Errorf("processing data: %w", err)
	}

	// Print operation summary based on mode
	if cfg.Mode == "file" {
		printer.Stderrln("\n%sed file: %q", cfg.Operation, cfg.File)
	}

	if cfg.Mode == "line" && processed {
		printer.Stderrln("\n%sed lines in: %q", cfg.Operation, cfg.File)
	}

	return nil
}

// loadData returns a file handle for the input data.
func loadData(file string) (*os.File, error) {
	data, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("opening input file %q: %w", file, err)
	}

	return data, nil
}
