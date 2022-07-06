package horcrux

import (
	"encoding/json"
	"hash/crc32"
	"log"
	"os"

	"gitlab.com/unit410/vault-shamir/shamir"
)

func Restore(shareFiles []string) ([]byte, error) {
	// if any of the records have a checksum, we will compare it
	// to a checksum of the assembled shares
	var checksum *uint32

	var shares [][]byte
	// ask gpg to decrypt the share files
	for _, shareFileName := range shareFiles {
		jsonBytes, err := os.ReadFile(shareFileName)
		if err != nil {
			return nil, err
		}

		var record Record
		if err := json.Unmarshal(jsonBytes, &record); err != nil {
			log.Fatalf("Error unmarshalling json bytes: %s", err)
		}

		if record.Checksum != nil {
			checksum = record.Checksum
		}

		if record.Threshold != 0 && record.Threshold > len(shareFiles) {
			log.Fatalf("Error: The threshold requires %d shares but only %d were provided.", record.Threshold, len(shareFiles))
		}

		if len(record.Pubkey) > 0 {
			share := record.Decrypt()
			if share == nil {
				continue
			}
			shares = append(shares, share)
		} else {
			shares = append(shares, record.Payload)
		}
	}

	original, err := shamir.Combine(shares)
	if err != nil {
		return nil, err
	}

	if checksum != nil && crc32.ChecksumIEEE(original) != *checksum {
		log.Fatalf("Error: Checksum of assembled shares does not match original")
	}
	return original, nil
}
