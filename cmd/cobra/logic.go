package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/idelchi/go-next-tag/pkg/stdin"
	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/internal/printer"
)

func processFiles(cfg *Config) error {
	var key []byte

	_, err := os.Stat(cfg.Key)
	if errors.Is(err, os.ErrNotExist) {
		printer.Stderrln("Key file %q does not exist, assuming the key is given as a string", cfg.Key)

		key = []byte(cfg.Key)
	} else {
		key, err = os.ReadFile(cfg.Key)
		if err != nil {
			return fmt.Errorf("reading key file: %w", err)
		}
	}

	key, err = encrypt.DecodeKey(string(key))

	data, err := loadData(cfg.File)
	if err != nil {
		return fmt.Errorf("loading data: %w", err)
	}
	defer data.Close()

	encryptor := &encrypt.Encryptor{
		Key:        key,
		Operation:  encrypt.Operation(cfg.Operation),
		Mode:       encrypt.Mode(cfg.Mode),
		Directives: cfg.Directives,
	}

	processed, err := encryptor.Process(data, os.Stdout)
	if err != nil {
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
