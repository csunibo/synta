package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type formatCommand struct {
	clear bool
	write bool
}

func (*formatCommand) Name() string     { return "format" }
func (*formatCommand) Synopsis() string { return "Format a synta file in the standard style." }
func (*formatCommand) Usage() string {
	return `format [-clear] [-write] <fle>:
  Format a synta file in the standard style.
`
}

func (p *formatCommand) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.clear, "clear", false, "Remove unused definitions")
	f.BoolVar(&p.write, "write", false, "Write changes to the input file")
}

func (p *formatCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	synta, status := parseFile(p, f)
	if status != subcommands.ExitSuccess {
		return status
	}

	fmt.Println(synta)
	return subcommands.ExitSuccess
}
