package internal

// A record stores a single share and accompanying metadata
type Record struct {
	// The required number of shares to complete a successful restore
	Threshold int
	// If set, contains the gpg public key used to encrypt the payload
	Pubkey []byte
	// The share contents; possibly gpg encrypted.
	Payload []byte
	// A sha256 checksum of the original content
	// Used to check the restore operation for success
	// Use a pointer here because early version of the program did not have a checksum field.
	// This allows the field to be ignored if not present.
	Checksum *uint32
}
