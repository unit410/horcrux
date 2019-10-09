package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"hash/crc32"
	"horcrux/internal"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/subcommands"
	"gitlab.com/polychainlabs/vault-shamir/shamir"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/// Restore command line arguments
type RestoreArgs struct {
	output string
}

/// Restore command name
func (*RestoreArgs) Name() string { return "restore" }

/// Short summary for the restore command
func (*RestoreArgs) Synopsis() string { return "restore a file from shamir shares" }

/// Full usage for the restore command
func (*RestoreArgs) Usage() string {
	return `restore <share file> <share file> <share file>

  Restore a file from shamir shares. Specify at least as many
  shares as required by the threshold.

  By default the output is written to STDOUT. Use --output to
  write to a file.

`
}

/// Initialze restore command flags
func (args *RestoreArgs) SetFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&args.output, "output", "", "Output to this file (default: stdout)")
}

/// Run the restore command
func (args *RestoreArgs) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	shareFiles := flagSet.Args()
	if len(shareFiles) < 2 {
		log.Fatal("At least two files are required to assemble from a fracture.")
	}

	// if any of the records have a checksum, we will compare it to a checksum of the aseembled shares
	var checksum *uint32

	var shares [][]byte
	// ask gpg to decrypt the share files
	for _, shareFileName := range shareFiles {
		jsonBytes, err := ioutil.ReadFile(shareFileName)
		assert(err)

		var record internal.Record
		json.Unmarshal(jsonBytes, &record)

		if record.Checksum != nil {
			checksum = record.Checksum
		}

		if record.Threshold != 0 && record.Threshold > len(shareFiles) {
			log.Fatalf("Error: The threshold requires %d shares but only %d were provided.", record.Threshold, len(shareFiles))
		}

		if len(record.Pubkey) > 0 {
			share := decryptPayload(record.Payload, record.Pubkey)
			if share == nil {
				continue
			}
			shares = append(shares, share)
		} else {
			shares = append(shares, record.Payload)
		}
	}

	original, err := shamir.Combine(shares)
	assert(err)

	if checksum != nil && crc32.ChecksumIEEE(original) != *checksum {
		log.Fatalf("Error: Checksum of assembled shares does not match original")
	}

	// if output is not defined, we output to stdout
	if args.output == "" {
		fmt.Printf("%s", original)
		return subcommands.ExitSuccess
	}

	err = ioutil.WriteFile(args.output, original, 0644)
	assert(err)

	return subcommands.ExitSuccess
}

// askForConfirmation asks the user for confirmation. A user must type "y" or "yes"
// Any other input will be considered as a No.
// Return true if the user confirmed with a yes, false otherwise.
func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	log.Printf("%s [y/N]: ", prompt)

	response, err := reader.ReadString('\n')
	assert(err)

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func decryptPayload(payload []byte, pubkey []byte) (share []byte) {
	packetReader := packet.NewReader(bytes.NewReader(pubkey))
	entity, err := openpgp.ReadEntity(packetReader)
	assert(err)

	fingerprint := entity.PrimaryKey.Fingerprint

	log.Printf("----------------------------------------------------------------\n")
	log.Printf("%X\n", fingerprint)
	for _, value := range entity.Identities {
		log.Printf("    %s\n", value.Name)
	}

	userPresent := askForConfirmation("Has the above identity's private key been loaded or smartcard inserted? (No to skip)")
	if !userPresent {
		log.Printf("Skipping...\n")
		return nil
	}

	{
		// we must import the pubkey before gpg --card-status can pair with the private key
		cmd := exec.Command("gpg", "--import")
		cmd.Stdin = bytes.NewBuffer(pubkey)
		err = cmd.Run()
		assert(err)
	}

	for {
		{
			// run --card-status to associate pubkeys with any inserted smartcards.
			// Without this command - gpg doesn't yet realize the inserted card has
			// a private key for a known public key in our keyring
			cmd := exec.Command("gpg", "--card-status")
			// Execute, but ignore the output in case user is not using a smartcard.
			cmd.CombinedOutput()
		}

		// ask gpg to decrypt the file
		log.Printf("Decrypting %x share...\n", fingerprint)
		cmd := exec.Command("gpg", "--decrypt")

		var stderr bytes.Buffer
		cmd.Stdin = bytes.NewReader(payload)
		cmd.Stderr = &stderr
		stdout, err := cmd.Output()
		if err != nil {
			log.Printf("%s\n", stderr.String())
			retry := askForConfirmation("Failed to decrypt share. Retry?")
			if !retry {
				break
			}

			// try again
			continue
		}

		log.Printf("%x share decrypted with size %d\n", fingerprint, len(stdout))
		return stdout
	}

	return nil
}
