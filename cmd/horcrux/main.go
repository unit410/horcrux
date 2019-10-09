package main

import (
	"context"
	"flag"
	"os"

	"horcrux/internal/cmd"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(&cmd.SplitArgs{}, "")
	subcommands.Register(&cmd.RestoreArgs{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
