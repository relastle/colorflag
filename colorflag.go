package colorflag

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
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
func (o *OutputFormatter) AddFlag(flg *flag.Flag) {
	name, usage := flag.UnquoteUsage(flg)
	o.addIndent()
	o.result += fmt.Sprintf(
		"%v <%v> %v\n",
		color.GreenString("-"+flg.Name),
		name,
		usage,
	)
}

// AddSubCommand adds subcommand
func (o *OutputFormatter) AddSubCommand(subCommand string) {
	o.addIndent()
	o.result += fmt.Sprintf(
		"%v\n",
		color.GreenString(subCommand),
	)
}

// CloseGroup closes one group.
// which break line and unshift indent
func (o *OutputFormatter) CloseGroup() {
	o.result += "\n"
	o.currentIndent -= o.Indent
}

// Print prints constructed help message
func (o *OutputFormatter) Print() {
	fmt.Printf(o.result)
}

func overrideSubCommandUsage(flagSet *flag.FlagSet) {
	flagSet.Usage = func() {
		outputFormatter := newOutputFormatter(2)
		outputFormatter.AddGroup(flagSet.Name())
		flagSet.VisitAll(func(flg *flag.Flag) {
			outputFormatter.AddFlag(flg)
		})
		outputFormatter.Print()
	}
}

func fetchFlagSet(flagSets []*flag.FlagSet, firstArg string) *flag.FlagSet {
	for _, flagSet := range flagSets {
		if flagSet.Name() == firstArg {
			return flagSet
		}
	}
	return nil
}

// overrideUsages overrides usage help massege of
// main command and sub commands
func overrideUsages(flagSets []*flag.FlagSet) {
	// Override main help
	flag.Usage = func() {
		outputFormatter := newOutputFormatter(2)

		outputFormatter.AddGroup(flag.CommandLine.Name())
		flag.CommandLine.VisitAll(func(flg *flag.Flag) {
			outputFormatter.AddFlag(flg)
		})
		outputFormatter.CloseGroup()

		outputFormatter.AddGroup("Sub Commands")
		for _, flagSet := range flagSets {
			outputFormatter.AddSubCommand(flagSet.Name())
		}
		outputFormatter.Print()
	}
	// set colorable stderr
	flag.CommandLine.SetOutput(colorable.NewColorableStderr())

	// Override sub command help
	for _, flagSet := range flagSets {
		overrideSubCommandUsage(flagSet)
		// set colorable stderr
		flagSet.SetOutput(colorable.NewColorableStderr())
	}
}

// Parse parse subcommands and override usage
func Parse(flagSets []*flag.FlagSet) string {
	overrideUsages(flagSets)

	if len(os.Args) == 1 {
		flag.Parse()
		return ""
	}

	firstArg := os.Args[1]
	fetchedFlagSet := fetchFlagSet(flagSets, firstArg)

	if fetchedFlagSet == nil {
		flag.Parse()
		return ""
	}
	fetchedFlagSet.Parse(os.Args[2:])
	return fetchedFlagSet.Name()
}
