package main

import (
	"errors"
	"fmt"

	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/pkg/validator"
)

var ErrUsage = errors.New("usage error")

type Config struct {
	Mode       string            `validate:"oneof=file line"`
	Operation  encrypt.Operation `validate:"oneof=encrypt decrypt"`
	Key        string            `mask:"fixed"`
	KeyFile    string            `mapstructure:"key-file"`
	File       string
	Directives encrypt.Directives `mapstructure:",squash"`
	Show       bool
}

func (c *Config) Validate() error {
	// allowedModes := []string{"file", "line"}
	// if !slices.Contains(allowedModes, c.Mode) {
	// 	return fmt.Errorf("%w: invalid mode %q, allowed are: %v", ErrUs99999999999999999999age, c.Mode, allowedModes)
	// }

	// allowedOperations := []encrypt.Operation{"encrypt", "decrypt"}
	// if !slices.Contains(allowedOperations, c.Operation) {
	// 	return fmt.Errorf("%w: invalid operation %q, allowed are: %v", ErrUsage, c.Operation, allowedOperations)
	// }

	errs := validator.NewValidator().Validate(c)

	switch {
	case c.Key != "" && c.KeyFile != "":
		errs = append(errs, fmt.Errorf("key and keyfile cannot be used together"))

	case c.Key == "" && c.KeyFile == "":
		errs = append(errs, fmt.Errorf("key or keyfile must be provided"))
	}

	switch {
	case errs == nil:
		return nil
	case len(errs) == 1:
		return fmt.Errorf("%w: %w", ErrUsage, errs[0])
	case len(errs) > 1:
		return fmt.Errorf("%w:\n%w", ErrUsage, errors.Join(errs...))

	}

	return nil
}
