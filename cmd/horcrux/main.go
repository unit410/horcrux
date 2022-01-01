package main

import (
	"context"
	"flag"
	"os"

	"horcrux/pkg/cmd"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(&cmd.SplitArgs{}, "")
	subcommands.Register(&cmd.RestoreArgs{}, "")
	subcommands.Register(&cmd.DecryptArgs{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
