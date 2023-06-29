package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/csunibo/synta"
	"github.com/google/subcommands"
)

func parseFile(p subcommands.Command, f *flag.FlagSet) (*synta.Synta, subcommands.ExitStatus) {
	filename := f.Arg(0)
	if filename == "" {
		fmt.Println(p.Usage())
		return nil, subcommands.ExitUsageError
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error while reading file: %s\n%v\n", filename, err)
		return nil, subcommands.ExitFailure
	}

	synta, err := synta.ParseSynta(string(contents))
	if err != nil {
		fmt.Printf("Invalid syuntax:\n%v", err)
		return nil, subcommands.ExitFailure
	}
	return &synta, subcommands.ExitSuccess
}
