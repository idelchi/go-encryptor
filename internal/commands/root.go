package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/idelchi/gocry/internal/config"
)

// NewRootCommand creates the root command with common configuration.
// It sets up environment variable binding and flag handling.
func NewRootCommand(cfg *config.Config, version string) *cobra.Command {
	root := &cobra.Command{
		Version:          version,
		SilenceUsage:     true,
		SilenceErrors:    true,
		Use:              "gocry [flags] command [flags]",
		Short:            "File/line encryption utility",
		Long:             "gocry is a utility for encrypting and decrypting files or lines of text.",
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			viper.SetEnvPrefix(cmd.Root().Name())
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
			viper.AutomaticEnv()

			if err := viper.BindPFlags(cmd.Root().Flags()); err != nil {
				return fmt.Errorf("binding root flags: %w", err)
			}

			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return fmt.Errorf("binding command flags: %w", err)
			}

			return nil
		},
	}

	root.Flags().BoolP("show", "s", false, "Show the configuration and exit")
	root.Flags().StringP("key", "k", "", "Encryption key")
	root.Flags().StringP("key-file", "f", "", "Path to the key file with the encryption key")
	root.Flags().StringP("mode", "m", "file", "Mode of operation: file or line")
	root.Flags().StringP("encrypt", "e", "### DIRECTIVE: ENCRYPT", "Directives for encryption")
	root.Flags().StringP("decrypt", "d", "### DIRECTIVE: DECRYPT", "Directives for decryption")

	root.AddCommand(NewEncryptCommand(cfg), NewDecryptCommand(cfg))

	root.CompletionOptions.DisableDefaultCmd = true
	root.Flags().SortFlags = false

	root.SetVersionTemplate("{{ .Version }}\n")

	return root
}
