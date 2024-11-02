package encrypt

// Operation represents the encryption or decryption operation.
type Operation string

const (
	// Encrypt operation.
	Encrypt Operation = "encrypt"
	// Decrypt operation.
	Decrypt Operation = "decrypt"
)

// Type represents whether encryption is deterministic or not.
type Type string

// Mode represents the mode of operation.
type Mode string

const (
	// Line mode processes each line of the input data.
	Line Mode = "line"
	// File mode processes the entire input data as a single block.
	File Mode = "file"
)
