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
	numShares int
	threshold int
	gpgKeyDir string
}

// Name of the `split` command
func (*SplitArgs) Name() string { return "split" }

// Synopsis of the `split` command
func (*SplitArgs) Synopsis() string { return "shamir split a file" }

// Usage of the `split` command
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

// SetFlags initializes split command flags
func (args *SplitArgs) SetFlags(flagSet *flag.FlagSet) {
	flagSet.IntVar(&args.numShares, "num-shares", 0, "Total number of shares (The 'n' in 'k of n')")
	flagSet.IntVar(&args.threshold, "threshold", 2, "Required parts to reconstruct (The 'k' of 'k of n')")
	flagSet.StringVar(&args.gpgKeyDir, "gpg-keys", "", "Directory with .asc files of public keys")
}

// Execute runs the split command
func (args *SplitArgs) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	remaining := flagSet.Args()
	if len(remaining) <= 0 {
		log.Fatal("At least one file required to fracture.")
	}
	if len(args.gpgKeyDir) == 0 && args.numShares == 0 {
		log.Fatal("One of --num-shares or --gpg-keys is required.")
	}
	filename := remaining[0]
	err := horcrux.Split(filename, args.numShares, args.threshold, args.gpgKeyDir)
	if err != nil {
		log.Fatal("Error!", err)
	}
	return subcommands.ExitSuccess
}
