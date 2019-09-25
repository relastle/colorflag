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
	sub1FlagSet.String("subarg1", "", "Description of subarg1")

	// sub command 2
	sub2FlagSet := flag.NewFlagSet("sub2", flag.ExitOnError)
	sub2FlagSet.String("subarg2", "", "Description of subarg2")

	flagSets := []*flag.FlagSet{sub1FlagSet, sub2FlagSet}
	colorflag.Parse(flagSets)
}
