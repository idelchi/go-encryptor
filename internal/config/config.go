package config

import (
	"errors"
	"fmt"

	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/pkg/validator"
)

// ErrUsage indicates an error in command-line usage or configuration.
var ErrUsage = errors.New("usage error")

type Key struct {
	// `validate:"hexadecimal,len=64|hexadecimal,len=32"`
	String string `mask:"fixed" validate:"exclusive=File" mapstructure:"key" label:"--key"`
	File   string `validate:"exclusive=String" mapstructure:"key-file" label:"--key-file"`
}

type Config struct {
	// Show enables output display
	Show bool

	Mode       string             `validate:"oneof=file line"`
	Operation  encrypt.Operation  `validate:"oneof=encrypt decrypt"`
	Key        Key                `mapstructure:",squash"`
	File       string             `validate:"required"`
	Directives encrypt.Directives `mapstructure:",squash"`
}

// Validate performs configuration validation using the validator package.
// It returns a wrapped ErrUsage if any validation rules are violated.
func Validate(config any) error {
	validator := validator.NewValidator()

	if err := registerExclusive(validator); err != nil {
		return fmt.Errorf("registering exclusive: %w", err)
	}

	errs := validator.Validate(config)

	switch {
	case errs == nil:
		return nil
	case len(errs) == 1:
		return fmt.Errorf("%w: %w", ErrUsage, errs[0])
	case len(errs) > 1:
		return fmt.Errorf("%ws:\n%w", ErrUsage, errors.Join(errs...))
	}

	return nil
}
