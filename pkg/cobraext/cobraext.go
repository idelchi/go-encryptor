package cobraext

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewDefaultRootCommand creates the root command with common configuration.
// It sets up environment variable binding and flag handling.
func NewDefaultRootCommand(version string) *cobra.Command {
	root := &cobra.Command{
		Version:          version,
		SilenceUsage:     true,
		SilenceErrors:    true,
		Use:              "gocry [flags] command [flags]",
		Short:            "File/line encryption utility",
		Long:             "gocry is a utility for encrypting and decrypting files or lines of text.",
		TraverseChildren: true, // TODO(Idelchi): Breaks suggestions, see below.
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
		RunE: UnknownSubcommandAction,
	}

	root.CompletionOptions.DisableDefaultCmd = true
	root.Flags().SortFlags = false

	root.SetVersionTemplate("{{ .Version }}\n")

	return root
}

// https://github.com/containerd/nerdctl/blob/242e6fc6e861b61b878bd7df8bf25e95674c036d/cmd/nerdctl/main.go#L401-L418
func UnknownSubcommandAction(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	err := fmt.Sprintf("unknown subcommand %q for %q", args[0], cmd.Name())
	if suggestions := cmd.SuggestionsFor(args[0]); len(suggestions) > 0 {
		err += "\n\nDid you mean this?\n"
		for _, s := range suggestions {
			err += fmt.Sprintf("\t%v\n", s)
		}
	}
	return fmt.Errorf(err)
}
