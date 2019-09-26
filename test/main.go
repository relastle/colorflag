package main

import (
	"flag"

	"github.com/relastle/colorflag"
)

func main() {
	// main command
	flag.String("arg1", "", "Description of arg1")
	flag.String("arg2", "", "Description of arg2")

	// sub command 1
	sub1FlagSet := flag.NewFlagSet("sub1", flag.ExitOnError)
	sub1FlagSet.String("sub1-option-string", "", "Description of string option")
	sub1FlagSet.Bool("sub1-flag", false, "Description of flag")

	// sub command 2
	sub2FlagSet := flag.NewFlagSet("sub2", flag.ExitOnError)
	sub2FlagSet.String("subarg2", "", "Description of subarg2")

	flagSets := []*flag.FlagSet{sub1FlagSet, sub2FlagSet}
	subCommand := colorflag.Parse(flagSets)

	switch subCommand {
	case "sub1":
		// Handle sub1 sub-command
	case "sub2":
		// Handle sub2 sub-command
	}
}
