package main

import (
	"flag"

	"github.com/relastle/colorflag"
)

func main() {
	// main command
	flag.String("main-opt1-string", "", "Description of `quoted target`")
	flag.Int("main-opt2-integer", 0, "Description of integer option")

	// sub command 1
	sub1FlagSet := flag.NewFlagSet("sub1", flag.ExitOnError)
	sub1FlagSet.String("sub1-string-option", "default string value", "Description of string option of sub command 1")
	sub1FlagSet.Int("sub1-integer-option", 2, "Description of integer option for sub command")
	sub1FlagSet.Bool("sub1-bool-flag", false, "Description of flag for sub command")

	// sub command 2
	sub2FlagSet := flag.NewFlagSet("sub2", flag.ExitOnError)
	sub2FlagSet.String("sub2-string option", "", "Description of string option of sub command 2")

	// Optional
	colorflag.Indent = 4                // default is 2
	colorflag.ExpandsSubCommand = false // default is true

	// Parse (and return selected sub command name)
	subCommand := colorflag.Parse([]*flag.FlagSet{
		sub1FlagSet,
		sub2FlagSet,
	})

	switch subCommand {
	case "sub1":
		// Handle sub1 sub-command
	case "sub2":
		// Handle sub2 sub-command
	}
}
