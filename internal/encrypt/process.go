package encrypt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

// processLines processes each line of the input data in parallel when possible
func (e *Encryptor) processLines(reader io.Reader, writer io.Writer) (bool, error) {
	var processed bool
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Channel for ordered output
	type lineOutput struct {
		line  string
		index int
	}
	outputChan := make(chan lineOutput)

	// Read all lines first since we need to maintain order
	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return processed, fmt.Errorf("scanning error: %v", err)
	}

	// Process lines in parallel
	results := make([]string, len(lines))
	numWorkers := 4 // Adjust based on system capabilities

	// Create work channel
	workChan := make(chan int)
	errChan := make(chan error)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range workChan {
				line := lines[idx]
				var result string

				switch {
				case e.Operation == Encrypt && strings.HasSuffix(line, e.Directives.Encrypt):
					encryptedLine, err := e.encryptData([]byte(line))
					if err != nil {
						errChan <- err
						return
					}
					mu.Lock()
					processed = true
					mu.Unlock()
					result = fmt.Sprintf("%s: %s", e.Directives.Decrypt, string(encryptedLine))

				case e.Operation == Decrypt && strings.HasPrefix(line, fmt.Sprintf("%s: ", e.Directives.Decrypt)):
					encryptedData := strings.TrimPrefix(line, fmt.Sprintf("%s: ", e.Directives.Decrypt))
					decryptedLine, err := e.decryptData([]byte(encryptedData))
					if err != nil {
						errChan <- err
						return
					}
					mu.Lock()
					processed = true
					mu.Unlock()
					result = string(decryptedLine)

				default:
					result = line
				}

				results[idx] = result
			}
		}()
	}

	// Send work
	go func() {
		for i := range lines {
			workChan <- i
		}
		close(workChan)
	}()

	// Wait for completion or error
	go func() {
		wg.Wait()
		close(errChan)
		close(outputChan)
	}()

	// Check for errors
	if err := <-errChan; err != nil {
		return processed, err
	}

	// Write results in order
	for _, line := range results {
		if _, err := fmt.Fprintln(writer, line); err != nil {
			return processed, err
		}
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
