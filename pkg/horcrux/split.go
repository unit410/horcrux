package horcrux

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strings"

	"gitlab.com/unit410/vault-shamir/shamir"
	"golang.org/x/crypto/openpgp"
)

func Split(filename string, numShares int, threshold int, gpgKeyDir string) error {

	// if we have a gpg-keys rather than num-shares, read the files and build the entity list
	var entities []*openpgp.Entity
	if len(gpgKeyDir) > 0 {
		files, err := ioutil.ReadDir(gpgKeyDir)
		Assert(err)

		// read armored file and store into entites array
		// the size of this array will be the number of shares
		for _, file := range files {
			// we are only interested in regular files not directories
			if !file.Mode().IsRegular() {
				continue
			}

			keyFile := path.Join(gpgKeyDir, file.Name())
			log.Printf("importing key %s\n", keyFile)

			// we import into gpg for use with encryption
			cmd := exec.Command("gpg", "--import", keyFile)
			_, err := cmd.CombinedOutput()
			Assert(err)

			entityList := getEntityListFromFile(keyFile)
			entities = append(entities, entityList...)
		}

		for _, entity := range entities {
			log.Printf("Identity: %x\n", entity.PrimaryKey.Fingerprint[:])
			for _, value := range entity.Identities {
				log.Printf("    %s\n", value.Name)
			}
		}

		numShares = len(entities)
	}

	// break the file into n shares
	log.Printf("Fracturing into %d shares requiring %d to assemble\n", numShares, threshold)
	original, err := ioutil.ReadFile(filename)
	Assert(err)

	checksum := crc32.ChecksumIEEE(original)

	shares, err := shamir.Split(original, numShares, threshold)
	Assert(err)

	if !CheckAllCombinations(original, shares, threshold) {
		return errors.New("combination permutation check failed")
	}

	if len(shares) != numShares {
		return errors.New("something went wrong - number of shares does not match requested num-shares")
	}

	// each encrypted part will be an output filename of the form:
	// <basename>.<part id>.json
	basename := strings.TrimSuffix(filename, path.Ext(filename))
	useGpg := len(entities) > 0

	for idx, part := range shares {
		var outputPart Record
		outputPart.Threshold = threshold
		outputPart.Payload = part
		outputPart.Checksum = &checksum

		// if we are using gpg, store the public key header
		// this allows restore to occur without needing to gather the gpg key files
		if useGpg {
			entity := entities[idx]

			var buf bytes.Buffer
			serializeWithoutSigs(entity, &buf)
			outputPart.Pubkey = buf.Bytes()

			fingerprint := entity.PrimaryKey.Fingerprint[:]

			// encrypted output buffer
			var encBuf bytes.Buffer

			// We have trust issues, so we use gpg directly to encrypt rather than golang openpgp
			cmd := exec.Command("gpg", "-r", fmt.Sprintf("%x", fingerprint), "--encrypt", "--trust-model", "always")
			cmd.Stdout = &encBuf
			cmd.Stdin = bytes.NewBuffer(outputPart.Payload)
			err = cmd.Run()
			Assert(err)

			outputPart.Payload = encBuf.Bytes()
		}

		jsonBytes, err := json.Marshal(&outputPart)
		Assert(err)

		// idx + 1 to start the count at 1
		partFilename := fmt.Sprintf("%s.%d.json", basename, idx+1)
		err = ioutil.WriteFile(partFilename, jsonBytes, 0644)
		Assert(err)
		log.Printf("Wrote %s\n", partFilename)
	}

	return nil
}
