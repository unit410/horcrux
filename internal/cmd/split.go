package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/subcommands"
	"gitlab.com/polychainlabs/vault-shamir/shamir"
	"golang.org/x/crypto/openpgp"

	"horcrux/internal"
)

/// Arguments for the `split` command
type SplitArgs struct {
	numShares int
	threshold int
	gpgKeyDir string
}

/// Split command Name
func (*SplitArgs) Name() string { return "split" }

/// Split command synopsis (short summary)
func (*SplitArgs) Synopsis() string { return "shamir split a file" }

/// Split command full usage
func (*SplitArgs) Usage() string {
	return `split [options] <file.ext>

  Shamir split a file into shares.

  The output share content will be written to files in the current working directory.

    file.0.json
    file.1.json
    ...
    file.n.json

  Options:

`
}

/// Initialize flag reading for split command
func (args *SplitArgs) SetFlags(flagSet *flag.FlagSet) {
	flagSet.IntVar(&args.numShares, "num-shares", 0, "Total number of shares (The 'n' in 'k of n')")
	flagSet.IntVar(&args.threshold, "threshold", 2, "Required parts to reconstruct (The 'k' of 'k of n')")
	flagSet.StringVar(&args.gpgKeyDir, "gpg-keys", "", "Directory with .asc files of public keys")
}

// Serialize the public part of the given Entity to w, excluding signatures from other entities
func serializeWithoutSigs(entity *openpgp.Entity, w io.Writer) error {
	err := entity.PrimaryKey.Serialize(w)
	if err != nil {
		return err
	}
	for _, ident := range entity.Identities {
		err = ident.UserId.Serialize(w)
		if err != nil {
			return err
		}
		err = ident.SelfSignature.Serialize(w)
		if err != nil {
			return err
		}
	}
	for _, subkey := range entity.Subkeys {
		err = subkey.PublicKey.Serialize(w)
		if err != nil {
			return err
		}
		err = subkey.Sig.Serialize(w)
		if err != nil {
			return err
		}
	}
	return nil
}

/// Run the split command
func (args *SplitArgs) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	remaining := flagSet.Args()
	if len(remaining) <= 0 {
		log.Fatal("At least one file required to fracture.")
	}
	filename := remaining[0]
	numShares := args.numShares

	if len(args.gpgKeyDir) == 0 && numShares == 0 {
		log.Fatal("One of --num-shares or --gpg-keys is required.")
	}

	// if we have a gpg-keys rather than num-shares, read the files and build the entity list
	var entities []*openpgp.Entity
	if len(args.gpgKeyDir) > 0 {
		files, err := ioutil.ReadDir(args.gpgKeyDir)
		internal.Assert(err)

		// read armored file and store into entites array
		// the size of this array will be the number of shares
		for _, file := range files {
			// we are only interested in regular files not directories
			if !file.Mode().IsRegular() {
				continue
			}

			keyringFilename := path.Join(args.gpgKeyDir, file.Name())
			log.Printf("importing key %s\n", keyringFilename)

			// we import into gpg for use with encryption
			cmd := exec.Command("gpg", "--import", keyringFilename)
			_, err := cmd.CombinedOutput()
			internal.Assert(err)

			keyringReader, err := os.Open(keyringFilename)
			internal.Assert(err)
			defer keyringReader.Close()

			entityList, err := openpgp.ReadArmoredKeyRing(keyringReader)
			internal.Assert(err)

			keyringReader.Close()
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
	log.Printf("Fracturing into %d shares requiring %d to assemble\n", numShares, args.threshold)
	original, err := ioutil.ReadFile(filename)
	internal.Assert(err)

	checksum := crc32.ChecksumIEEE(original)

	shares, err := shamir.Split(original, numShares, args.threshold)
	internal.Assert(err)

	if !internal.CheckAllCombinations(original, shares, args.threshold) {
		log.Fatal("Combination permutation check failed")
	}

	if len(shares) != numShares {
		log.Fatal("Something went wrong - number of shares does not match requested num-shares")
	}

	// each encrypted part will be an output filename of the form:
	// <basename>.<part id>.json
	basename := strings.TrimSuffix(filename, path.Ext(filename))
	useGpg := len(entities) > 0

	for idx, part := range shares {
		var outputPart internal.Record
		outputPart.Threshold = args.threshold
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
			internal.Assert(err)

			outputPart.Payload = encBuf.Bytes()
		}

		jsonBytes, err := json.Marshal(&outputPart)
		internal.Assert(err)

		// idx + 1 to start the count at 1
		partFilename := fmt.Sprintf("%s.%d.json", basename, idx+1)
		err = ioutil.WriteFile(partFilename, jsonBytes, 0644)
		internal.Assert(err)
		log.Printf("Wrote %s\n", partFilename)
	}

	return subcommands.ExitSuccess
}
