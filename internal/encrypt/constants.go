package encrypt

// Operation represents the encryption or decryption operation.
// It determines whether the Encryptor will encrypt or decrypt the input data.
type Operation string

const (
	// Encrypt specifies that data should be encrypted.
	// This operation takes plaintext input and produces encrypted output.
	Encrypt Operation = "encrypt"

	// Decrypt specifies that data should be decrypted.
	// This operation takes encrypted input and produces plaintext output.
	Decrypt Operation = "decrypt"
)

// Type represents whether encryption is deterministic or not.
// Note: This type is prepared for future extensions of the encryption package.
type Type string

// Mode represents the mode of operation for processing input data.
// It determines how the input data is handled during encryption/decryption.
type Mode string

const (
	// Line mode processes each line of the input data separately.
	// It allows for selective encryption/decryption of individual lines
	// based on markers and maintains the original file structure.
	Line Mode = "line"

	// File mode processes the entire input data as a single block.
	// It treats the entire input as one piece of data to be
	// encrypted/decrypted, suitable for binary files or whole-file encryption.
	File Mode = "file"
)
