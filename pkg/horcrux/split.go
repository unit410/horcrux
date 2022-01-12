package horcrux

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"horcrux/internal/gpg"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"gitlab.com/unit410/vault-shamir/shamir"
)

func SplitEncrypt(filename string, gpgRecipientFiles []string, threshold int, outputDir string, omitPubkeys bool) error {
	numShares := len(gpgRecipientFiles)
	if threshold > numShares {
		log.Fatal("Threshold cannot exceept number of shares")
	}
	if numShares < 2 {
		log.Fatal("Must split to at 2+ recipients")
	}

	records, err := splitFile(filename, numShares, threshold)
	if err != nil {
		return err
	}

	records, err = encryptRecords(records, gpgRecipientFiles, omitPubkeys)
	if err != nil {
		return err
	}

	err = writeRecords(getSplitFilebase(filename, outputDir), records)
	if err != nil {
		return err
	}

	return nil
}

func Split(filename string, numShares int, threshold int, outputDir string) error {
	if threshold > numShares {
		log.Fatal("Threshold cannot exceept number of shares")
	}
	if numShares < 2 {
		log.Fatal("Must split to at 2+ recipients")
	}

	records, err := splitFile(filename, numShares, threshold)
	if err != nil {
		return err
	}

	err = writeRecords(getSplitFilebase(filename, outputDir), records)
	if err != nil {
		return err
	}

	return nil
}

func encryptRecords(records []Record, gpgRecipientFiles []string, omitPubkeys bool) ([]Record, error) {
	if len(records) != len(gpgRecipientFiles) {
		log.Fatal("Error: Can only encrypt shares if the same number of gpg files are provided")
	}

	logln("Importing keys:")

	entities, err := gpg.GetEntities(gpgRecipientFiles)
	if err != nil {
		return nil, err
	}

	logln("Found identities:")
	for _, entity := range entities {
		logf("- %x\n", entity.PrimaryKey.Fingerprint[:])
		for _, identity := range entity.Identities {
			logf("  - %s\n", identity.Name)
		}
	}

	if len(entities) != len(records) {
		log.Fatalf("%d entities != %d records", len(entities), len(records))
	}
	for i := range records {
		entity := entities[i]

		// Serialize the pubkey so we can restore without needing to gather the gpg key files
		var buf bytes.Buffer
		err := gpg.SerializeWithoutSigs(entity, &buf)
		if err != nil {
			return nil, err
		}
		if omitPubkeys {
			records[i].OmitPubkey()
		} else {
			records[i].Pubkey = buf.Bytes()
		}

		// Encrypt the payload
		var encBuf bytes.Buffer

		// We have trust issues, so we use gpg directly to encrypt rather than golang openpgp
		fingerprint := entity.PrimaryKey.Fingerprint[:]
		cmd := exec.Command("gpg", "-r", fmt.Sprintf("%x", fingerprint), "--encrypt", "--trust-model", "always")
		cmd.Stdout = &encBuf
		cmd.Stdin = bytes.NewBuffer(records[i].Payload)
		err = cmd.Run()
		if err != nil {
			return nil, err
		}

		records[i].Payload = encBuf.Bytes()
	}
	return records, nil
}

func splitFile(filename string, numShares int, threshold int) ([]Record, error) {
	// break the file into n shares
	logf("Fracturing into %d shares requiring %d to assemble\n", numShares, threshold)
	original, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	checksum := crc32.ChecksumIEEE(original)

	shares, err := shamir.Split(original, numShares, threshold)
	if err != nil {
		return nil, err
	}

	if !CheckAllCombinations(original, shares, threshold) {
		return nil, errors.New("combination permutation check failed")
	}

	if len(shares) != numShares {
		return nil, errors.New("something went wrong - number of shares does not match requested num-shares")
	}

	logln("Writing output files")
	var records []Record
	for _, part := range shares {
		var outputPart Record
		outputPart.Threshold = threshold
		outputPart.Payload = part
		outputPart.Checksum = &checksum
		records = append(records, outputPart)
	}

	return records, nil
}

// Each encrypted part will be an output filename of the form:
// <basename>.<part id>.json
func getSplitFilebase(filename, outputDir string) string {
	basename := strings.TrimSuffix(filename, path.Ext(filename))
	// If defined, move file into outputDir
	if len(outputDir) > 0 {
		_, filename := path.Split(basename)
		basename = path.Join(outputDir, filename)
	}
	return basename
}

func writeRecords(outputBase string, records []Record) error {
	for i, r := range records {
		jsonBytes, err := json.Marshal(&r)
		if err != nil {
			return err
		}

		outputFile := fmt.Sprintf("%s.%d.json", outputBase, i+1)
		err = os.WriteFile(outputFile, jsonBytes, 0644)
		if err != nil {
			return err
		}
		logf("- Wrote %s\n", outputFile)
	}
	return nil
}
