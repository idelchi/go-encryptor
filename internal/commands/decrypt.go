package commands

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/gocry/internal/config"
	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/internal/logic"
)

// NewDecryptCommand creates a new cobra command for the decrypt operation.
func NewDecryptCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decrypt file",
		Aliases: []string{"dec"},
		Short:   "Decrypt files",
		Long:    "Decrypt a file using the specified key. Output is printed to stdout.",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(_ *cobra.Command, args []string) error {
			cfg.Operation = encrypt.Decrypt

			return setFileAndValidate(cfg, args)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return logic.Run(cfg)
		},
	}

	return cmd
}
