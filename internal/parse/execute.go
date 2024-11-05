// Package parse generates and executes the command-line interface for the application.
package parse

import (
	"errors"
	"fmt"

	"github.com/idelchi/gocry/internal/commands"
	"github.com/idelchi/gocry/internal/config"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// Execute creates and configures the command-line interface.
// It runs the root command with all subcommands and flags configured.
func Execute(version string) error {
	cfg := &config.Config{}
	root := commands.NewRootCommand(cfg, version)

	switch err := root.Execute(); {
	case errors.Is(err, cobraext.ErrExitGracefully):
		return nil
	case err != nil:
		return fmt.Errorf("executing command: %w", err)
	default:
		return nil
	}
}
