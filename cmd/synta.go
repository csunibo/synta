package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

// We are using the simple google/subcommands library to make handle subcommands
func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&formatCommand{}, "")
	subcommands.Register(&checkCommand{}, "")
	subcommands.Register(&jsonSchemaCommand{}, "")
	subcommands.Register(&regexpCommand{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
