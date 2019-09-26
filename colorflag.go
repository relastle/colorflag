package colorflag

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Indent indicates depth of one indent level
	// (the number of spaces inserted).
	Indent int = 2

	// ExpandsSubCommand defines whether options and flags
	// of sub-commands are displayed in the top level help
	// message.
	ExpandsSubCommand bool = true
)

// OutputFormatter is a formatter that constructs help
// messages in a structured way.
type OutputFormatter struct {
	Indent        int
	currentIndent int
	result        string
	currentFlags  []*flag.Flag
}

func newOutputFormatter(indent int) *OutputFormatter {
	return &OutputFormatter{
		Indent: indent,
	}
}

func (o *OutputFormatter) makeOffsetSpaces(n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += " "
	}
	return res
}

func (o *OutputFormatter) addIndent() {
	o.result += o.makeOffsetSpaces(o.currentIndent)
}

// InitGroup adds group name which is followd by
// multiple options or flags
func (o *OutputFormatter) InitGroup(groupName string) {
	o.addIndent()
	o.result += color.YellowString(groupName) + "\n"
	o.currentIndent += o.Indent
}

// AddFlag adds group name which is followd by
// multiple options or flags
func (o *OutputFormatter) AddFlag(flg *flag.Flag) {
	o.currentFlags = append(o.currentFlags, flg)
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
	flagNames := []string{}
	names := []string{}
	usages := []string{}
	defValues := []string{}
	for _, flg := range o.currentFlags {
		name, usage := flag.UnquoteUsage(flg)
		flagNames = append(flagNames, flg.Name)
		names = append(names, name)
		usages = append(usages, usage)
		defValues = append(defValues, flg.DefValue)
	}

	offsetSlices1 := makeOffsets(flagNames)
	offsetSlices2 := makeOffsets(names)

	for i := 0; i < len(o.currentFlags); i++ {
		flagName := flagNames[i]
		name := names[i]
		usage := usages[i]
		defValue := defValues[i]
		offset1 := offsetSlices1[i]
		offset2 := offsetSlices2[i]

		o.addIndent()
		o.result += fmt.Sprintf(
			"%v%v <%v>%v %v (default: %v)\n",
			color.GreenString("-"+flagName),
			o.makeOffsetSpaces(offset1),
			name,
			o.makeOffsetSpaces(offset2),
			usage,
			defValue,
		)
	}
	o.result += "\n"
	o.currentIndent -= o.Indent
	o.currentFlags = []*flag.Flag{}
}

// Print prints constructed help message
func (o *OutputFormatter) Print() {
	fmt.Printf(o.result)
}

func overrideSubCommandUsage(flagSet *flag.FlagSet) {
	flagSet.Usage = func() {
		outputFormatter := newOutputFormatter(Indent)
		outputFormatter.InitGroup(flagSet.Name())
		flagSet.VisitAll(func(flg *flag.Flag) {
			outputFormatter.AddFlag(flg)
		})
		outputFormatter.CloseGroup()
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
		outputFormatter := newOutputFormatter(Indent)

		outputFormatter.InitGroup(flag.CommandLine.Name())
		flag.CommandLine.VisitAll(func(flg *flag.Flag) {
			outputFormatter.AddFlag(flg)
		})
		outputFormatter.CloseGroup()

		outputFormatter.InitGroup("subcommands")
		for _, flagSet := range flagSets {
			if ExpandsSubCommand {
				outputFormatter.InitGroup(flagSet.Name())
				flagSet.VisitAll(func(flg *flag.Flag) {
					outputFormatter.AddFlag(flg)
				})
				outputFormatter.CloseGroup()
			} else {
				outputFormatter.AddSubCommand(flagSet.Name())
			}
		}
		outputFormatter.CloseGroup()

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
