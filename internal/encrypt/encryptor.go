package encrypt

import (
	"fmt"
	"io"
)

// Directives defines the markers used to identify content for encryption/decryption.
// The markers are configurable via mapstructure tags for external configuration.
type Directives struct {
	// Encrypt specifies the suffix that marks a line for encryption
	Encrypt string `mapstructure:"encrypt"`

	// Decrypt specifies the prefix that marks encrypted content
	Decrypt string `mapstructure:"decrypt"`
}

// Encryptor handles encryption and decryption operations.
type Encryptor struct {
	// Key is the encryption key used for AES cipher operations
	Key []byte

	// Operation specifies whether to encrypt or decrypt
	Operation Operation

	// Mode determines whether to process the input line-by-line or as a whole file
	Mode Mode

	// Directives contains the markers used to identify content for processing
	Directives Directives

	// Parallel specifies the number of goroutines to use for parallel processing
	Parallel int
}

// Process handles encryption and decryption based on the provided configuration.
// It returns a boolean indicating if any processing was performed and any error encountered.
// The processing mode (Line or File) determines how the input is handled:
//   - Line mode processes the input line-by-line, maintaining order
//   - File mode treats the entire input as a single block of data
func (e *Encryptor) Process(reader io.Reader, writer io.Writer) (bool, error) {
	switch e.Mode {
	case Line:
		return e.processLines(reader, writer, e.Parallel)
	case File:
		return e.processWholeFile(reader, writer)
	default:
		return false, fmt.Errorf("invalid mode: %s", e.Mode) //nolint: err113
	}
}
