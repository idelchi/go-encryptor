# gocry

[![Go Reference](https://pkg.go.dev/badge/github.com/idelchi/gocry.svg)](https://pkg.go.dev/github.com/idelchi/gocry)
[![Go Report Card](https://goreportcard.com/badge/github.com/idelchi/gocry)](https://goreportcard.com/report/github.com/idelchi/gocry)
[![Build Status](https://github.com/idelchi/gocry/actions/workflows/github-actions.yml/badge.svg)](https://github.com/idelchi/gocry/actions/workflows/github-actions.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`gocry` is a command-line utility for encrypting and decrypting files using a specified key.

It supports both file encryption and line-by-line encryption based on directives within the file.

The tool can read the file from stdin and write the encrypted/decrypted content to stdout,
making it suitable for use as a filter in git.

## Installation

### From source

```sh
go install github.com/idelchi/gocry@latest
```

### From installation script

```sh
curl -sSL https://raw.githubusercontent.com/idelchi/gocry/refs/heads/dev/install.sh | sh -s -- -d ~/.local/bin
```

## Usage

```sh
gocry [flags] command [flags]
```

### Global Flags and Environment Variables

| Flag             | Environment Variable      | Description                         | Default                  |
| ---------------- | ------------------------- | ----------------------------------- | ------------------------ |
| `-j, --parallel` | `GOCRY_PARALLEL`          | Number of parallel workers          | `runtime.NumCPU()`       |
| `-k, --key`      | `GOCRY_KEY`               | Key for encryption/decryption       | -                        |
| `-f, --key-file` | `GOCRY_KEY_FILE`          | Path to the key file                | -                        |
| `-m, --mode`     | `GOCRY_MODE`              | Mode of operation: `file` or `line` | `file`                   |
| `--encrypt`      | `GOCRY_ENCRYPT_DIRECTIVE` | Directive for encryption            | `### DIRECTIVE: ENCRYPT` |
| `--decrypt`      | `GOCRY_DECRYPT_DIRECTIVE` | Directive for decryption            | `### DIRECTIVE: DECRYPT` |
| `-s, --show`     | `GOCRY_SHOW`              | Show the configuration and exit     | `false`                  |
| `-h, --help`     | -                         | Help for `gocry`                    | -                        |
| `-v, --version`  | -                         | Version for `gocry`                 | -                        |

### Commands

#### `encrypt` - Encrypt content

Encrypt a file or specific lines within a file.

Examples:

```sh
# Encrypt an entire file
gocry -f path/to/keyfile encrypt input.txt > encrypted.txt.enc

# Encrypt specific lines in a file
gocry -f path/to/keyfile -m line encrypt input.txt > encrypted.txt
```

#### `decrypt` - Decrypt content

Decrypt a file or specific lines within a file.

Examples:

```sh
# Decrypt an entire file
gocry -f path/to/keyfile decrypt encrypted.txt.enc > decrypted.txt.dec

# Decrypt specific lines in a file
gocry -f path/to/keyfile -m line decrypt encrypted.txt > decrypted.txt
```

### Git Integration

gocry can be used as a filter in git for automatic encryption/decryption of files.

When stdin is given, gocry reads the file from stdin (using the file arguments for logging),
and writes the encrypted/decrypted content to stdout.

**.gitconfig:**

```gitconfig
[filter "encrypt:line"]
    clean = "gocry -f ~/.secrets/key -m line encrypt %f"
    smudge = "gocry -f ~/.secrets/key  -m line decrypt %f"
    required = true

[filter "encrypt:file"]
    clean = "gocry -f ~/.secrets/key -m file encrypt  %f"
    smudge = "gocry -f ~/.secrets/key -m file decrypt %f"
    required = true
```

**.gitattributes:**

```gitattributes
*                       filter=encrypt:line
**/secrets/*            filter=encrypt:file
```

### Line-by-Line Encryption

When using `--mode line`, gocry processes only lines containing specific directives:

**Input Example:**

```text
This is a normal line.
This line will be encrypted. ### DIRECTIVE: ENCRYPT
Another normal line.
```

**After Encryption:**

```text
This is a normal line.
### DIRECTIVE: DECRYPT: VGhpcyBsaW5lIHdpbGwgYmUgZW5jcnlwdGVkLiBPRmx2eGZpRk9GMkF3PT0=
Another normal line.
```

**After Decryption:**

```text
This is a normal line.
This line will be encrypted. ### DIRECTIVE: ENCRYPT
Another normal line.
```

For detailed help on any command:

```sh
gocry --help
gocry <command> --help
```
