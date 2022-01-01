package horcrux

// Record stores a single share and accompanying metadata
type Record struct {
	// Threshold is the required number of shares to complete a successful restore
	Threshold int
	// Pubkey if set, contains the gpg public key used to encrypt the payload
	Pubkey []byte
	// Payload contains the share contents; possibly gpg encrypted.
	Payload []byte
	// Checksum sha256 checksum of the original content
	// Used to check the restore operation for success
	// Use a pointer here because early version of the program did not have a checksum field.
	// This allows the field to be ignored if not present.
	Checksum *uint32
}
