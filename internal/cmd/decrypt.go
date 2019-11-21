package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"horcrux/internal"
	"io/ioutil"
	"log"

	"github.com/google/subcommands"
)

/// Decrypt command line arguments
type DecryptArgs struct {
	output string
}

/// Decrypt command name
func (*DecryptArgs) Name() string { return "decrypt" }

/// Short summary for the decrypt command
func (*DecryptArgs) Synopsis() string { return "decrypt a share file that was encryted with a gpg key" }

/// Full usage for the decrypt command
func (*DecryptArgs) Usage() string {
	return `decrypt <share file>

  Decrypt a share file that has been encrypted with a gpg key.

  By default the output is written to STDOUT. Use --output to
  write to a file.

`
}

/// Initialze decrypt command flags
func (args *DecryptArgs) SetFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&args.output, "output", "", "Output to this file (default: stdout)")
}

/// Run the decrypt command
func (args *DecryptArgs) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	shareFiles := flagSet.Args()
	if len(shareFiles) != 1 {
		log.Fatal("A share file is required.")
	}

	shareFileName := shareFiles[0]
	cleartextRecord := internal.Record{}

	jsonBytes, err := ioutil.ReadFile(shareFileName)
	assert(err)

	var record internal.Record
	json.Unmarshal(jsonBytes, &record)

	cleartextRecord.Threshold = record.Threshold
	cleartextRecord.Checksum = record.Checksum

	if len(record.Pubkey) == 0 {
		log.Fatal("This share has not been encrypted with gpg. Aborting.")
	}

	// ask gpg to decrypt the share files
	share := decryptPayload(record.Payload, record.Pubkey)
	if share == nil {
		log.Fatal("Failed to decrypt.")
	}

	cleartextRecord.Payload = share

	cleartextRecordJSONBytes, err := json.Marshal(&cleartextRecord)
	assert(err)

	// if output is not defined, we output to stdout
	if args.output == "" {
		fmt.Printf("%s", cleartextRecordJSONBytes)
	} else {
		err = ioutil.WriteFile(args.output, cleartextRecordJSONBytes, 0644)
		assert(err)
	}

	return subcommands.ExitSuccess
}
