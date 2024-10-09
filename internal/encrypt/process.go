package encrypt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// processLines processes each line of the input data, encrypting or decrypting lines
// that contain the specific directive.
func (e *Encryptor) processLines(reader io.Reader, writer io.Writer) (bool, error) {
	var processed bool

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case e.Operation == Encrypt && strings.HasSuffix(line, e.Directives.Encrypt):
			encryptedLine, err := e.encryptData([]byte(line))
			if err != nil {
				return processed, err
			}

			processed = true
			fmt.Fprintf(writer, "%s: %s\n", e.Directives.Decrypt, encryptedLine)
		case e.Operation == Decrypt && strings.HasPrefix(line, fmt.Sprintf("%s: ", e.Directives.Decrypt)):
			encryptedData := strings.TrimPrefix(line, fmt.Sprintf("%s: ", e.Directives.Decrypt))
			decryptedLine, err := e.decryptData([]byte(encryptedData))
			if err != nil {
				return processed, err
			}

			processed = true
			fmt.Fprintln(writer, string(decryptedLine))
		default:
			fmt.Fprintln(writer, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return processed, fmt.Errorf("scanning error: %v", err)
	}
	return processed, nil
}

// processWholeFile processes the entire input data as a single encrypted or decrypted block.
func (e *Encryptor) processWholeFile(reader io.Reader, writer io.Writer) (bool, error) {
	switch e.Operation {
	case Encrypt:
		return true, e.encryptStream(reader, writer)
	case Decrypt:
		return true, e.decryptStream(reader, writer)
	default:
		return false, fmt.Errorf("invalid operation")
	}
}
