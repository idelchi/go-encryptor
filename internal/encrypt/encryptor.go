package encrypt

import (
	"fmt"
	"io"
)

type Directives struct {
	Encrypt string
	Decrypt string
}

// Encryptor handles encryption and decryption operations.
type Encryptor struct {
	Key        []byte
	Operation  Operation
	Mode       Mode
	Type       Type
	Directives Directives
}

// Process handles encryption and decryption based on the provided configuration.
// It delegates to either processLines or processWholeFile depending on the mode.
func (e *Encryptor) Process(reader io.Reader, writer io.Writer) (bool, error) {
	switch e.Mode {
	case Line:
		return e.processLines(reader, writer)
	case File:
		return e.processWholeFile(reader, writer)
	default:
		return false, fmt.Errorf("invalid mode: %s", e.Mode)
	}
}
