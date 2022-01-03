package horcrux

import (
	"fmt"
	"os"
	"time"

	"horcrux/internal/gpg"
)

// DecryptPayload decrypts and returns the payload. A pubkey can optionally
// be provided to display additional data about the payload
func (r Record) Decrypt() (share []byte) {
	// Get the ID of the key that this payload is encrypted to
	keyID, err := gpg.GetEncryptionPacketKeyID(r.Payload)
	Assert(err)

	// Import the pubkey so gpg recognizes the card's keys
	if !r.PubkeyIsOmitted() {
		err = gpg.ImportPubkey(r.Pubkey)
		if err != nil {
			warnf(err.Error())
		}
	}

	// Wait for the key to be available
	logf("Waiting for key or card for %s...\n", keyID)

	// Display helpful headers
	if !r.PubkeyIsOmitted() {
		for _, name := range getEntityNames(r.Pubkey, keyID) {
			logf("- Name: %s\n", name)
		}
	}

	needsNewline := false
	for {
		availability := SecretKeyLocation(keyID)
		if availability == KeyLocationLocal {
			break
		}
		if availability == KeyLocationStub {
			if smartcardHasKey(keyID) {
				break
			}
		}
		fmt.Fprintf(os.Stderr, ".")
		time.Sleep(500 * time.Millisecond)
		needsNewline = true
	}
	if needsNewline {
		logln("")
	}

	// Decrypt the file
	for {
		logln("Decrypting share: ", keyID)
		stdout, stderr, err := gpg.Decrypt(r.Payload)
		if err != nil {
			logf("%s\n", stderr.String())
			retry := AskForConfirmation(os.Stdin, "Failed to decrypt share. Retry?")
			if !retry {
				break
			}
			// try again
			continue
		}
		return stdout
	}

	return nil
}

func getEntityNames(pubkey []byte, keyID string) []string {
	entity, err := gpg.ReadEntity(pubkey)
	if err != nil {
		warnf(err.Error())
		return nil
	}
	result := []string{}
	for _, identity := range entity.Identities {
		result = append(result, identity.Name)
	}
	return result
}
