package commands

import (
	"runtime"

	"github.com/spf13/cobra"

	"github.com/idelchi/gocry/internal/config"
	"github.com/idelchi/gogen/pkg/cobraext"
)

// NewRootCommand creates the root command with common configuration.
// It sets up environment variable binding and flag handling.
func NewRootCommand(cfg *config.Config, version string) *cobra.Command {
	root := cobraext.NewDefaultRootCommand(version)

	root.Use = "gocry [flags] command [flags]"
	root.Short = "File/line encryption utility"
	root.Long = "gocry is a utility for encrypting and decrypting files or lines of text."

	root.Flags().BoolP("show", "s", false, "Show the configuration and exit")
	root.Flags().IntP("parallel", "j", runtime.NumCPU(), "Number of parallel workers")
	root.Flags().StringP("key", "k", "", "Encryption key")
	root.Flags().StringP("key-file", "f", "", "Path to the key file with the encryption key")
	root.Flags().StringP("mode", "m", "file", "Mode of operation: file or line")
	root.Flags().StringP("encrypt", "e", "### DIRECTIVE: ENCRYPT", "Directives for encryption")
	root.Flags().StringP("decrypt", "d", "### DIRECTIVE: DECRYPT", "Directives for decryption")
	root.Flags().BoolP("experiments", "x", false, "Enable experimental features")

	root.AddCommand(NewEncryptCommand(cfg), NewDecryptCommand(cfg))

	return root
}
