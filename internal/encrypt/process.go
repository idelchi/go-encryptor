package encrypt

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

// processLines processes each line of the input data in parallel when possible.
// It maintains the original line order in the output while leveraging parallel processing.
// Returns a boolean indicating if any encryption/decryption was performed and any error encountered.
func (e *Encryptor) processLines(reader io.Reader, writer io.Writer, parallel int) (bool, error) {
	var processed bool
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Channel for ordered output
	type lineOutput struct {
		line  string
		index int
	}
	outputChan := make(chan lineOutput)

	// Read all lines first to maintain output order
	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return processed, fmt.Errorf("scanning error: %v", err)
	}

	// Initialize result storage and channels
	results := make([]string, len(lines))
	numWorkers := parallel
	workChan := make(chan int)
	errChan := make(chan error)

	// Start worker goroutines for parallel processing
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range workChan {
				line := lines[idx]
				var result string

				// Process each line based on operation type and directives
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

	// Distribute work to workers
	go func() {
		for i := range lines {
			workChan <- i
		}
		close(workChan)
	}()

	// Wait for completion and close channels
	go func() {
		wg.Wait()
		close(errChan)
		close(outputChan)
	}()

	// Check for processing errors
	if err := <-errChan; err != nil {
		return processed, err
	}

	// Write results maintaining original order
	for _, line := range results {
		if _, err := fmt.Fprintln(writer, line); err != nil {
			return processed, err
		}
	}

	return processed, nil
}

// processWholeFile processes the entire input as a single block of data.
// It's used when line-by-line processing is not required.
// Returns true if processing was performed and any error encountered.
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
