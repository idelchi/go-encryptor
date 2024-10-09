package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/idelchi/gocry/internal/encrypt"
)

// Config represents the configuration for the go-encryptor application.
type Config struct {
	Mode       string `validate:"oneof=file line"`
	Operation  string `validate:"oneof=encrypt decrypt"`
	Key        string `validate:"required"`
	File       string
	Type       string `validate:"oneof=deterministic nondeterministic"`
	GPG        bool
	Directives encrypt.Directives
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	allowedModes := []string{"file", "line"}

	if !slices.Contains(allowedModes, c.Mode) {
		return fmt.Errorf("%w: invalid mode %q, allowed are: %v", ErrUsage, c.Mode, allowedModes)
	}

	allowedOperations := []string{"encrypt", "decrypt"}
	if !slices.Contains(allowedOperations, c.Operation) {
		return fmt.Errorf("%w: invalid operation %q, allowed are: %v", ErrUsage, c.Operation, allowedOperations)
	}

	if c.Key == "" {
		return fmt.Errorf("%w: key file must be provided", ErrUsage)
	}

	allowedTypes := []string{"deterministic", "nondeterministic"}
	if !slices.Contains(allowedTypes, c.Type) {
		return fmt.Errorf("%w: invalid encryption type %q, allowed are: %v", ErrUsage, c.Type, allowedTypes)
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
}

// ErrUsage represents an error due to incorrect usage.
var ErrUsage = errors.New("usage error")
