# gocry

gocry is a command-line utility for encrypting and decrypting files using a specified key.

It supports file encryption and line-by-line encryption based on directives within the file.

The program outputs the processed content to standard output (stdout).

Can be used as filters in git.

`.gitconfig`

```toml
[filter "encrypt:line"]
    clean = "gocry -m lines -o encrypt -k ~/.secrets/key %f"
    smudge = "gocry -m lines -o decrypt -k ~/.secrets/key %f %f"
    required = true

[filter "encrypt:file"]
    clean = "gocry -m file -o encrypt -k ~/.secrets/key %f"
    smudge = "gocry -m file -o decrypt -k ~/.secrets/key %f %f"
    required = true
```

`.gitattributes`

```toml
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

- `-m, --mode`: Mode of operation: "file" or "line" (default "file")
- `-o, --operation`: Operation to perform: "encrypt" or "decrypt" (default "encrypt")
- `-k, --key`: Path to the key file (required)
- `-t, --type`: Encryption type: "deterministic" or "nondeterministic" (default "nondeterministic")
- `--gpg`: Whether a GPG key is used for encryption/decryption (default true)
- `--directives.encrypt`: Directives for encryption (default "### DIRECTIVE: ENCRYPT")
- `--directives.decrypt`: Directives for decryption (default "### DIRECTIVE: DECRYPT")
- `--version`: Show the version information and exit
- `-h, --help`: Show the help information and exit
- `-s, --show`: Show the configuration and exit

### Examples

#### Encrypt a File

Encrypt `input.txt` output the result to `encrypted.txt`:

```sh
gocry -k path/to/keyfile input.txt > encrypted.txt
```

#### Decrypt a File

Decrypt `encrypted.txt` using the same key and output the result to `decrypted.txt`:

```sh
gocry -o decrypt -k path/to/keyfile encrypted.txt > decrypted.txt
```

#### Encrypt Specific Lines in a File

Encrypt lines in `input.txt` that contain the directive `### DIRECTIVE: ENCRYPT` and output the result to `encrypted.txt`:

```sh
gocry -m line -k path/to/keyfile input.txt > encrypted.txt
```

#### Show the Configuration

Display the current configuration based on the provided flags:

```sh
gocry -s -k path/to/keyfile input.txt
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
gocry -m line -k path/to/keyfile input.txt > encrypted.txt
```

### Resulting Output (encrypted.txt):

```
This is a normal line.
### DIRECTIVE: DECRYPT: VGhpcyBsaW5lIHdpbGwgYmUgZW5jcnlwdGVkLiBPRmx2eGZpRk9GMkF3PT0=
Another normal line.
```

### Decrypting the File:

```sh
gocry -m line -o decrypt -k path/to/keyfile encrypted.txt > decrypted.txt
```

### Resulting Output (decrypted.txt):

```
This is a normal line.
This line will be encrypted. ### DIRECTIVE: ENCRYPT
Another normal line.
```

## Notes on Encryption Types

### Deterministic Encryption (`--type deterministic`):

- Uses a fixed initialization vector (IV) derived from the key.
- Encrypting the same content multiple times produces the same ciphertext.
- Useful when you need consistent encryption results for the same input.

### Nondeterministic Encryption (`--type nondeterministic`):

- Uses a randomly generated IV.
- Encrypting the same content multiple times produces different ciphertexts.
- Provides better security by ensuring that identical plaintexts encrypt to different ciphertexts.

## For More Details

To display a comprehensive list of flags and their descriptions, run:

```sh
gocry --help
```

## Under the Hood

gocry uses AES encryption in CFB (Cipher Feedback) mode. The key provided should be of appropriate length for AES (16, 24, or 32 bytes for AES-128, AES-192, or AES-256 respectively).

## Generating a Key

You can generate a key using OpenSSL or any other cryptographic library. Here's an example using OpenSSL to generate a 256-bit (32-byte) key:

```sh
openssl rand -out keyfile -base64 32
```
