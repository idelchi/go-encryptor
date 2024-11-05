package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/gocry/internal/config"
	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/internal/logic"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// NewEncryptCommand creates a new cobra command for the encrypt operation.
func NewEncryptCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encrypt file",
		Aliases: []string{"enc"},
		Short:   "Encrypt files",
		Long:    "Encrypt a file using the specified key. Output is printed to stdout.",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, args []string) error {
			arg, err := cobraext.PipeOrArg(args)
			if err != nil {
				return fmt.Errorf("reading password: %w", err)
			}

			cfg.File = arg
			cfg.Operation = encrypt.Encrypt

			return cobraext.Validate(cfg, cfg)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return logic.Run(cfg)
		},
	}

	return cmd
}
