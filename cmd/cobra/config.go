package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/idelchi/gocry/internal/encrypt"
)

var ErrUsage = errors.New("usage error")

type Config struct {
	Mode       string
	Operation  encrypt.Operation
	Key        string `mask:"fixed"`
	File       string
	Directives encrypt.Directives
	Show       bool
}

func (c *Config) Validate() error {
	allowedModes := []string{"file", "line"}
	if !slices.Contains(allowedModes, c.Mode) {
		return fmt.Errorf("%w: invalid mode %q, allowed are: %v", ErrUsage, c.Mode, allowedModes)
	}

	allowedOperations := []encrypt.Operation{"encrypt", "decrypt"}
	if !slices.Contains(allowedOperations, c.Operation) {
		return fmt.Errorf("%w: invalid operation %q, allowed are: %v", ErrUsage, c.Operation, allowedOperations)
	}

	if c.Key == "" {
		return fmt.Errorf("%w: key file/string must be provided", ErrUsage)
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
}
