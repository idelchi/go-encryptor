package logic

import (
	"fmt"
	"os"

	"github.com/idelchi/go-next-tag/pkg/stdin"
	"github.com/idelchi/gocry/internal/config"
	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/internal/printer"
	"github.com/idelchi/gogen/pkg/key"
)

func Run(cfg *config.Config) error {
	var (
		encryptionKey []byte
		err           error
	)

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
		Parallel:   cfg.Parallel,
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
