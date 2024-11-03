package commands

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/gocry/internal/config"
	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/gocry/internal/logic"
)

func NewEncryptCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encrypt file",
		Aliases: []string{"enc"},
		Short:   "Encrypt files",
		Long:    "Encrypt a file using the specified key. Output is printed to stdout.",
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cfg.Operation = encrypt.Encrypt
			cfg.File = args[0]

			if err := validate(cfg, cfg); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic.Run(cfg)
		},
	}

	return cmd
}
