package config

import (
	"errors"
	"fmt"

	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gogen/pkg/validator"
)

// ErrUsage indicates an error in command-line usage or configuration.
var ErrUsage = errors.New("usage error")

type Key struct {
	// String is a hexadecimal key string
	String string `mask:"fixed" validate:"hexadecimal,len=64,exclusive=File" mapstructure:"key" label:"--key"`
	// File is a path to a file containing a hexadecimal key string
	File string `validate:"exclusive=String" mapstructure:"key-file" label:"--key-file"`
}

type Config struct {
	// Show enables output display
	Show bool

	// Parallel is the number of parallel workers to use
	Parallel int `validate:"min=1" mapstructure:"parallel"`

	// Mode is the encryption mode
	Mode encrypt.Mode `validate:"oneof=file line"`
	// Operation is the encryption operation
	Operation encrypt.Operation `validate:"oneof=encrypt decrypt" mapstructure:"-"`
	// Key is the encryption key
	Key Key `mapstructure:",squash"`
	// File is the path to the input file
	File string `validate:"required" mapstructure:"-"`
	// Directives contains the markers used to identify content for processing
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
