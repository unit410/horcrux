package cmd

import (
	"context"
	"flag"
	"horcrux/pkg/horcrux"
	"log"

	"github.com/google/subcommands"
)

// SplitArgs models arguments for the `split` command
type SplitArgs struct {
	numShares   int
	threshold   int
	outputDir   string
	omitPubkeys bool
}

// Name of the `split` command
func (*SplitArgs) Name() string { return "split" }

// Synopsis of the `split` command
func (*SplitArgs) Synopsis() string { return "shamir split a file" }

// Usage of the `split` command
func (*SplitArgs) Usage() string {
	return `split [options] <file.ext> <gpg-recipient-pubkey-file>, (<gpg-recipient-pubkey-file>, ...)

  Shamir split a file into shares, optionally GPG encrypting to the given recipients.

    file.0.json
    file.1.json
    ...
    file.n.json

  Options:

`
}

// SetFlags initializes split command flags
func (args *SplitArgs) SetFlags(flagSet *flag.FlagSet) {
	flagSet.IntVar(&args.numShares, "num-shares", 0, "Total number of shares (The 'n' in 'm of n')")
	flagSet.IntVar(&args.threshold, "threshold", 2, "Required parts to reconstruct (The 'm' of 'm of n')")
	flagSet.BoolVar(&args.omitPubkeys, "omit-pubkeys", false, "Omit pubkeys from the output split")
	flagSet.StringVar(&args.outputDir, "output-dir", ".", "Directory where shares should be output")
}

// Execute runs the split command
func (args *SplitArgs) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	argv := flagSet.Args()
	if len(argv) <= 0 {
		log.Fatal("At least one file required to fracture.")
	}
	filename := argv[0]
	gpgRecipientFiles := argv[1:]

	// Validate input
	numGpgRecipients := len(argv) - 1
	if args.numShares == 0 && numGpgRecipients <= 1 {
		log.Fatal("Either 2+ gpg recipients *or* --num-shares must be set")
	}
	if args.numShares > 0 && numGpgRecipients > 0 {
		log.Fatal("Cannot set both gpg recipients and --num-shares")
	}

	var err error
	if args.numShares > 0 {
		err = horcrux.Split(filename, args.numShares, args.threshold, args.outputDir)
	} else {
		err = horcrux.SplitEncrypt(filename, gpgRecipientFiles, args.threshold, args.outputDir, args.omitPubkeys)
	}
	if err != nil {
		log.Fatal("Error:", err)
	}
	return subcommands.ExitSuccess
}
