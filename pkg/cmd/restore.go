package cmd

import (
	"context"
	"flag"
	"fmt"
	"horcrux/pkg/horcrux"
	"log"
	"os"

	"github.com/google/subcommands"
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

	original, err := horcrux.Restore(shareFiles)
	if err != nil {
		log.Fatal("Unable to decrypt with error: ", err)
	}

	// if output is not defined, we output to stdout
	if args.output == "" {
		fmt.Printf("%s", original)
		return subcommands.ExitSuccess
	}

	err = os.WriteFile(args.output, original, 0644)
	horcrux.Assert(err)

	return subcommands.ExitSuccess
}
