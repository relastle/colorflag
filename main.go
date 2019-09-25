package colorflag

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// OutputFormatter is a formatter that constructs help
// messages in a structured way.
type OutputFormatter struct {
	Indent        int
	currentIndent int
	result        string
}

func newOutputFormatter(indent int) *OutputFormatter {
	return &OutputFormatter{
		Indent: indent,
	}
}

func (o *OutputFormatter) addIndent() {
	for i := 0; i < o.currentIndent; i++ {
		o.result += " "
	}
}

// AddGroup adds group name which is followd by
// multiple options or flags
func (o *OutputFormatter) AddGroup(groupName string) {
	o.addIndent()
	o.result += color.YellowString(groupName) + "\n"
	o.currentIndent += o.Indent
}

// AddFlag adds group name which is followd by
// multiple options or flags
func (o *OutputFormatter) AddFlag(flagName string, flagType string, flagUsage string) {
	o.addIndent()
	o.result += fmt.Sprintf(
		"%v (%v) %v\n",
		color.GreenString(flagName),
		flagType,
		flagUsage,
	)
}

// Print prints constructed help message
func (o *OutputFormatter) Print() {
	fmt.Printf(o.result)
}

func printUsage(flagSet *flag.FlagSet) {
	outputFormatter := newOutputFormatter(2)
	outputFormatter.AddGroup(flagSet.Name())
	flagSet.VisitAll(func(flg *flag.Flag) {
		outputFormatter.AddFlag(flg.Name, "", flg.Usage)
	})
	outputFormatter.Print()
}

func fetchFlagSet(flagSets []*flag.FlagSet, firstArg string) *flag.FlagSet {
	for _, flagSet := range flagSets {
		if flagSet.Name() == firstArg {
			return flagSet
		}
	}
	return nil
}

func overrideUsages(flagSets []*flag.FlagSet) {
	// Override main help
	flag.CommandLine.Usage = func() {
		flag.CommandLine.VisitAll(func(flg *flag.Flag) {
			fmt.Println(flg.Name)

		})

		for _, flagSet := range flagSets {
			printUsage(flagSet)
		}
	}

	// Override sub command help
	for _, flagSet := range flagSets {
		flagSet.Usage = func() {
			printUsage(flagSet)
		}
	}
}

// Parse parse subcommands and override usage
func Parse(flagSets []*flag.FlagSet) {
	overrideUsages(flagSets)

	firstArg := os.Args[1]
	fetchedFlagSet := fetchFlagSet(flagSets, firstArg)
	if fetchedFlagSet != nil {
		fetchedFlagSet.Parse(os.Args[2:])
	} else {
		flag.Parse()
	}
}
