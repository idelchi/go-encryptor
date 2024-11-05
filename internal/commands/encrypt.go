package commands

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/gocry/internal/config"
	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/internal/logic"
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
			cfg.Operation = encrypt.Encrypt

			return setFileAndValidate(cfg, args)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return logic.Run(cfg)
		},
	}

	return cmd
}
