# gocry

gocry is a command-line utility for encrypting and decrypting files using a specified key.

It supports file encryption and line-by-line encryption based on directives within the file.

The program outputs the processed content to standard output (stdout).

Can be used as filters in git.

`.gitconfig`

```toml
[filter "encrypt:line"]
    clean = "gocry -k ~/.secrets/key -m lines encrypt %f"
    smudge = "gocry -k ~/.secrets/key  -m lines decrypt %f"
    required = true

[filter "encrypt:file"]
    clean = "gocry -f ~/.secrets/key -m file encrypt  %f"
    smudge = "gocry -f ~/.secrets/key -m file decrypt %f"
    required = true
```

`.gitattributes`

```gitattributes
*                       filter=encrypt:line
**/secrets/*            filter=encrypt:file
```

## Installation

### From source

```sh
go install github.com/idelchi/gocry@latest
```

## Usage

```sh
gocry [flags] [file]
```

The available flags include:

- `-s, --show`: Show the configuration and exit
- `-m, --mode`: Mode of operation: "file" or "line" (default "file")
- `-k, --key`: Key for encryption/decryption
- `-f, --key-file`: Path to the key file
- `--encrypt`: Directives for encryption (default `### DIRECTIVE: ENCRYPT`)
- `--decrypt`: Directives for decryption (default `### DIRECTIVE: DECRYPT`)
- `--version`: Show the version information and exit
- `-h, --help`: Show the help information and exit
- `-s, --show`: Show the configuration and exit

### Examples

#### Encrypt a File

Encrypt `input.txt` output the result to `encrypted.txt`:

```sh
gocry -f path/to/keyfile encrypt input.txt > encrypted.txt.enc
```

#### Decrypt a File

Decrypt `encrypted.txt` using the same key and output the result to `decrypted.txt`:

```sh
gocry -k path/to/keyfile decrypt encrypted.txt.enc > decrypted.txt.dec
```

#### Encrypt Specific Lines in a File

Encrypt lines in `input.txt` that contain the directive `### DIRECTIVE: ENCRYPT` and output the result to `encrypted.txt`:

```sh
gocry -m line -k path/to/keyfile encrypt input.txt > encrypted.txt
```

#### Show the Configuration

Display the current configuration based on the provided flags:

```sh
gocry -s -k path/to/keyfile encrypt input.txt
```

#### Display Help Information

Show detailed help information:

```sh
gocry --help
```

## Directives for Line-by-Line Encryption

When using `--mode line`, `gocry` processes only the lines that contain specific directives:

- To encrypt a line, append `### DIRECTIVE: ENCRYPT` to the line.
- To decrypt a line, it should start with `### DIRECTIVE: DECRYPT:` followed by the encrypted content.

The directives themselves can be customized using the `--directives.encrypt` and `--directives.decrypt` flags.

### Example Input File (input.txt):

```
This is a normal line.
This line will be encrypted. ### DIRECTIVE: ENCRYPT
Another normal line.
```

### Encrypting the File:

```sh
gocry -m line -k path/to/keyfile encrypt input.txt > encrypted.txt
```

### Resulting Output (encrypted.txt):

```
This is a normal line.
### DIRECTIVE: DECRYPT: VGhpcyBsaW5lIHdpbGwgYmUgZW5jcnlwdGVkLiBPRmx2eGZpRk9GMkF3PT0=
Another normal line.
```

### Decrypting the File:

```sh
gocry -k path/to/keyfile -m line decrypt encrypted.txt > decrypted.txt
```

### Resulting Output (decrypted.txt):

```
This is a normal line.
This line will be encrypted. ### DIRECTIVE: ENCRYPT
Another normal line.
```

## For More Details

To display a comprehensive list of flags and their descriptions, run:

```sh
gocry --help
```
