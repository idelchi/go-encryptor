package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/idelchi/gocry/internal/encrypt"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func flags() *cobra.Command {
	cfg := &Config{}
	root := newRootCmd(cfg)

	// Persistent flags shared across encrypt/decrypt commands
	root.Flags().StringP("key", "k", "", "Path to the key file")
	root.Flags().StringP("mode", "m", "file", "Mode of operation: file or line")
	root.Flags().BoolP("show", "s", false, "Show the configuration and exit")
	root.Flags().StringP("encrypt", "e", "### DIRECTIVE: ENCRYPT", "Directives for encryption")
	root.Flags().StringP("decrypt", "d", "### DIRECTIVE: DECRYPT", "Directives for decryption")

	encrypt := newEncryptCmd(cfg)
	decrypt := newDecryptCmd(cfg)
	generate := newGenerateCmd()

	root.AddCommand(encrypt, decrypt, generate)

	root.CompletionOptions.DisableDefaultCmd = true
	root.SetVersionTemplate("{{ .Version }}\n")
	root.Flags().SortFlags = false

	generate.SetHelpFunc(func(command *cobra.Command, strings []string) {
		command.Flags().MarkHidden("key")
		command.Flags().MarkHidden("mode")
		command.Flags().MarkHidden("encrypt")
		command.Flags().MarkHidden("decrypt")
		command.Parent().HelpFunc()(command, strings)
	})

	return root
}

func validate(cfg *Config) error {
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	if cfg.Show {
		pretty.PrintJSONMasked(cfg)

		os.Exit(0)

		return nil
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
}

func newRootCmd(_ *Config) *cobra.Command {
	root := &cobra.Command{
		Version:          version,
		Use:              "gonc [flags] command [flags]",
		Short:            "File/line encryption utility",
		Long:             "gonc is a utility for encrypting and decrypting files or lines of text.",
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip config processing for generate command
			if cmd.Name() == "generate" {
				return nil
			}

			viper.SetEnvPrefix(cmd.Root().Name())
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viper.AutomaticEnv()

			if err := viper.BindPFlags(cmd.Root().Flags()); err != nil {
				return fmt.Errorf("failed to bind flags: %w", err)
			}

			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				return fmt.Errorf("failed to bind persistent flags: %w", err)
			}

			return nil
		},
	}

	root.CompletionOptions.DisableDefaultCmd = true
	root.SetVersionTemplate("{{ .Version }}\n")

	return root
}

func newEncryptCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:     "encrypt file",
		Aliases: []string{"enc"},
		Short:   "Encrypt files",
		Long:    "Encrypt a file using the specified key. Output is printed to stdout.",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.Operation = encrypt.Encrypt

			if err := validate(cfg); err != nil {
				return err
			}

			cfg.File = args[0]
			return processFiles(cfg)
		},
	}
}

func newDecryptCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:     "decrypt file",
		Aliases: []string{"dec"},
		Short:   "Decrypt files",
		Long:    "Decrypt a file using the specified key. Output is printed to stdout.",
		// Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cfg.Operation = encrypt.Decrypt

			if err := validate(cfg); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.File = args[0]

			return processFiles(cfg)
		},
	}
}

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate a new encryption key",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := encrypt.GenerateKeyString()
			if err != nil {
				return fmt.Errorf("generating key: %w", err)
			}
			fmt.Println(key)

			return nil
		},
	}
	return cmd
}
