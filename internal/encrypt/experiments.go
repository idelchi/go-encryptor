package encrypt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// processLines processes each line of the input data sequentially.
// It maintains the original line order in the output.
// Returns a boolean indicating if any encryption/decryption was performed and any error encountered.
func (e *Encryptor) processLinesExperiments(reader io.Reader, writer io.Writer) (bool, error) {
	// Read all lines first to maintain output order
	var lines []string

	var anyFound bool

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())

		if strings.Contains(scanner.Text(), e.Directives.Encrypt) || strings.Contains(scanner.Text(), e.Directives.Decrypt) {
			anyFound = true
		}
	}

	if !anyFound {
		// No directives found, return input as-is
		for _, line := range lines {
			if _, err := fmt.Fprintln(writer, line); err != nil {
				return false, fmt.Errorf("%w: writing error: %w", ErrProcessing, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("%w: scanning error: %w", ErrProcessing, err)
	}

	wasPartOfPrevious := make([]bool, len(lines))

	var anyProcessed bool

	for idx := range lines {
		// Skip if this line was processed as part of previous directive
		if wasPartOfPrevious[idx] {
			continue
		}

		var (
			result       string
			wasProcessed bool
		)

		line := lines[idx]

		// Process based on operation type and directives
		switch {
		case e.Operation == Encrypt:
			switch {
			case strings.TrimSpace(line) == e.Directives.Encrypt && idx+1 < len(lines):
				nextLine := lines[idx+1]
				contentToEncrypt := line + "\n" + nextLine
				encryptedLine, err := e.encryptData([]byte(contentToEncrypt))
				if err != nil {
					return false, err
				}
				result = fmt.Sprintf("%s: %s", e.Directives.Decrypt, string(encryptedLine))
				wasProcessed = true
				// Mark next line as processed
				wasPartOfPrevious[idx+1] = true

			case strings.HasSuffix(line, e.Directives.Encrypt):
				encryptedLine, err := e.encryptData([]byte(line))
				if err != nil {
					return false, err
				}
				result = fmt.Sprintf("%s: %s", e.Directives.Decrypt, string(encryptedLine))
				wasProcessed = true

			default:
				result = line
			}

		case e.Operation == Decrypt && strings.HasPrefix(line, e.Directives.Decrypt+": "):
			encryptedData := strings.TrimPrefix(line, e.Directives.Decrypt+": ")
			decryptedLine, err := e.decryptData([]byte(encryptedData))
			if err != nil {
				return false, err
			}
			result = string(decryptedLine)
			wasProcessed = true

		default:
			result = line
		}

		if _, err := fmt.Fprintln(writer, result); err != nil {
			return false, fmt.Errorf("%w: writing error: %w", ErrProcessing, err)
		}

		if wasProcessed {
			anyProcessed = true
		}
	}

	return anyProcessed, nil
}
