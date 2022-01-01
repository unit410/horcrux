package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"hash/crc32"
	"horcrux/pkg/horcrux"
	"io/ioutil"
	"log"

	"github.com/google/subcommands"
	"gitlab.com/unit410/vault-shamir/shamir"
)

// RestoreArgs models command line arguments for the `restore` command
type RestoreArgs struct {
	output string
}

// Name of the restore command
func (*RestoreArgs) Name() string { return "restore" }

// Synopsis of the restore command
func (*RestoreArgs) Synopsis() string { return "restore a file from shamir shares" }

// Usage of the restore command
func (*RestoreArgs) Usage() string {
	return `restore <share file> <share file> <share file>

  Restore a file from shamir shares. Specify at least as many
  shares as required by the threshold.

  By default the output is written to STDOUT. Use --output to
  write to a file.

`
}

// SetFlags initializes restore command flags
func (args *RestoreArgs) SetFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&args.output, "output", "", "Output to this file (default: stdout)")
}

// Execute runs the restore command
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
		horcrux.Assert(err)

		var record horcrux.Record
		json.Unmarshal(jsonBytes, &record)

		if record.Checksum != nil {
			checksum = record.Checksum
		}

		if record.Threshold != 0 && record.Threshold > len(shareFiles) {
			log.Fatalf("Error: The threshold requires %d shares but only %d were provided.", record.Threshold, len(shareFiles))
		}

		if len(record.Pubkey) > 0 {
			share := horcrux.DecryptPayload(record.Payload, record.Pubkey)
			if share == nil {
				continue
			}
			shares = append(shares, share)
		} else {
			shares = append(shares, record.Payload)
		}
	}

	original, err := shamir.Combine(shares)
	horcrux.Assert(err)

	if checksum != nil && crc32.ChecksumIEEE(original) != *checksum {
		log.Fatalf("Error: Checksum of assembled shares does not match original")
	}

	// if output is not defined, we output to stdout
	if args.output == "" {
		fmt.Printf("%s", original)
		return subcommands.ExitSuccess
	}

	err = ioutil.WriteFile(args.output, original, 0644)
	horcrux.Assert(err)

	return subcommands.ExitSuccess
}
