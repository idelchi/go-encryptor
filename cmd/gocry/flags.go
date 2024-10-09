package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/idelchi/go-next-tag/pkg/stdin"
	"github.com/idelchi/godyl/pkg/flagexp"
)

func flags() {
	pflag.StringP("mode", "m", "file", "Mode of operation: file or line")
	pflag.StringP("operation", "o", "encrypt", "Operation to perform: encrypt or decrypt")
	pflag.StringP("key", "k", "", "Path to the key file")
	pflag.StringP("type", "t", "nondeterministic", "Encryption type: deterministic or nondeterministic")
	pflag.Bool("gpg", true, "Whether a GPG key is used for encryption/decryption")

	pflag.String("directives.encrypt", "### DIRECTIVE: ENCRYPT", "Directives for encryption")
	pflag.String("directives.decrypt", "### DIRECTIVE: DECRYPT", "Directives for decryption")

	pflag.Bool("version", false, "Show the version information and exit")
	pflag.BoolP("help", "h", false, "Show the help information and exit")
	pflag.BoolP("show", "s", false, "Show the configuration and exit")

	pflag.CommandLine.SortFlags = false

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [file]\n\n", "gocry")
		fmt.Fprintf(os.Stderr, "Encrypt or decrypt a file using a key. Output is printed to stdout.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		pflag.PrintDefaults()
	}
}

func parseFlags() (cfg Config, err error) {
	flags()

	// Parse the command-line flags with suggestions enabled
	if err := flagexp.ParseWithSuggestions(os.Args[1:]); err != nil {
		return cfg, fmt.Errorf("parsing flags: %w", err)
	}

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return Config{}, fmt.Errorf("binding flags: %w", err)
	}

	viper.SetEnvPrefix("gocry")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unmarshalling config: %w", err)
	}

	handleExitFlags(cfg)

	// Validate the input
	if err := validateInput(&cfg); err != nil {
		return cfg, fmt.Errorf("validating input: %w", err)
	}

	return cfg, nil
}

func validateInput(cfg *Config) error {
	switch args, isPiped := pflag.NArg(), stdin.IsPiped(); {
	case args > 1:
		return fmt.Errorf("too many arguments: %d", pflag.NArg())
	case args == 0 && !isPiped:
		return errors.New("input must be provided either via stdin and/or as a positional argument")
	default:
		cfg.File = pflag.Arg(0)
	}

	return nil
}

func handleExitFlags(cfg Config) {
	if viper.GetBool("version") {
		fmt.Println(version)
		os.Exit(0)
	}

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	if viper.GetBool("show") {
		fmt.Println(PrintJSON(cfg))

		os.Exit(0)
	}
}

// PrintJSON returns a pretty-printed JSON representation of the provided object.
func PrintJSON(obj any) string {
	bytes, err := json.MarshalIndent(obj, "  ", "    ")
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}
