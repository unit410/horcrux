package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"horcrux/pkg/horcrux"
	"log"
	"os"

	"github.com/google/subcommands"
)

// DecryptArgs models command line arguments for the `decrypt` command
type DecryptArgs struct {
	output string
}

// Name of the decrypt command
func (*DecryptArgs) Name() string { return "decrypt" }

// Synopsis of the decrypt command
func (*DecryptArgs) Synopsis() string { return "decrypt a share file that was encryted with a gpg key" }

// Usage for the decrypt command
func (*DecryptArgs) Usage() string {
	return `decrypt <share file>

  Decrypt a share file that has been encrypted with a gpg key.

  By default the output is written to STDOUT. Use --output to
  write to a file.

`
}

// SetFlags initializes decrypt command flags
func (args *DecryptArgs) SetFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&args.output, "output", "", "Output to this file (default: stdout)")
}

// Execute runs the decrypt command
func (args *DecryptArgs) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	shareFiles := flagSet.Args()
	if len(shareFiles) != 1 {
		log.Fatal("A share file is required.")
	}

	shareFileName := shareFiles[0]
	cleartextRecord := horcrux.Record{}

	jsonBytes, err := os.ReadFile(shareFileName)
	horcrux.Assert(err)

	var record horcrux.Record
	if err := json.Unmarshal(jsonBytes, &record); err != nil {
		log.Fatalf("Error unmarshalling json bytes: %s", err)
	}

	cleartextRecord.Threshold = record.Threshold
	cleartextRecord.Checksum = record.Checksum

	if len(record.Pubkey) == 0 {
		log.Fatal("This share has not been encrypted with gpg. Aborting.")
	}

	// ask gpg to decrypt the share files
	share := record.Decrypt()
	if share == nil {
		log.Fatal("Failed to decrypt.")
	}

	cleartextRecord.Payload = share

	cleartextRecordJSONBytes, err := json.Marshal(&cleartextRecord)
	horcrux.Assert(err)

	// if output is not defined, we output to stdout
	if args.output == "" {
		fmt.Printf("%s", cleartextRecordJSONBytes)
	} else {
		err = os.WriteFile(args.output, cleartextRecordJSONBytes, 0644)
		horcrux.Assert(err)
	}

	return subcommands.ExitSuccess
}
