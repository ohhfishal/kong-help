package konghelp

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"strings"
)

var _ kong.HelpPrinter = PrettyHelpPrinter

// TODO: Make these configurable
var ColorExample = color.New(color.FgYellow).SprintFunc()
var ColorRequired = color.New(color.FgRed).SprintFunc()
var ColorDefault = color.New(color.FgMagenta).SprintFunc()
var ColorCommand = color.New(color.FgCyan).SprintFunc()
var ColorLow = color.HiBlackString
var ColorType = ColorExample
var ColorGroup = color.New(color.FgBlue).Add(color.Underline).SprintFunc()

func PrettyHelpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	if ctx.Empty() {
		options.Summary = false
	}

	// TODO: Have this controlled via an option
	options.ValueFormatter = PrettyValueFormatter(options.ValueFormatter)

	w := newHelpWriter(ctx, options)
	if selected := ctx.Selected(); selected == nil {
		app := ctx.Model
		if !w.Options.NoAppSummary {
			w.Print("")
			w.Printf("  %s: %s%s", ColorExample("Usuage"), app.Name, app.Summary())
		}
		printNodeDetail(w, app.Node, true)
	} else {
		if !w.Options.NoAppSummary {
			w.Print("")
			w.Printf("  %s: %s", ColorExample("Usuage"), selected.Summary())
		}
		printNodeDetail(w, selected, true)
	}
	// TODO: Handle Run %s --help for more info lines
	if _, err := w.WriteTo(ctx.Stdout); err != nil {
		return err
	}
	return nil
}

func printNodeDetail(w *helpWriter, node *kong.Node, hide bool) {
	if node.Help != "" {
		w.Print("")
		w.Print(w.Wrap(node.Help))
	}
	if w.Options.Summary {
		return
	}
	if node.Detail != "" {
		w.Print("")
		w.Print(w.Wrap(node.Detail))
	}
	if len(node.Positional) > 0 {
		w.Print("")
		printPositionals(w, node.Positional)
	}
	if !w.Options.FlagsLast {
		printFlags(w, node.AllFlags(true))
	}
	// TODO: Print the commands here
	var cmds []*kong.Node
	if w.Options.NoExpandSubcommands {
		cmds = node.Children
	} else {
		cmds = node.Leaves(hide)
	}
	if len(cmds) > 0 {
		if w.Options.Tree {
			// TODO: Fix
			// TODO: Make it look nice with characters in the tree command
			panic("Options.Tree not supported")
		} else {
			printCommands(w, cmds)
		}
	}
	if w.Options.FlagsLast {
		printFlags(w, node.AllFlags(true))
	}
}

func printPositionals(w *helpWriter, args []*kong.Positional) {
	lines := [][]string{}
	for _, arg := range args {
		var prefix = "  "
		if arg.Tag != nil && arg.Tag.Required {
			prefix = ColorRequired("* ")
		}
		components := []string{
			prefix,
			arg.Name,
			// TODO: Parse format/enum to do along the line of PATH[existing file] or STRING[enum]
			formatValue(arg.Target, false),
			// TODO: Write a custom ValueFormatter to do: "Help Message. [required] [default=""] etc
			w.Options.ValueFormatter(arg),
		}
		lines = append(lines, components)
	}
	printCard(w, "Arguments", lines)
}

func printFlags(w *helpWriter, flags [][]*kong.Flag) {
	lines := [][]string{}
	for _, collection := range collectFlagGroups(flags) {
		lines = append(lines, formatGroup(collection.Metadata)...)
		for _, flagset := range collection.Flags {
			for _, flag := range flagset {
				line := formatFlag(flag, w.Options.ValueFormatter)
				lines = append(lines, line)
			}
		}
	}
	printCard(w, "Options", lines)
}

func printCommands(w *helpWriter, cmds []*kong.Command) {
	if w.Options.Compact {
		// TODO: Fix
		panic("compact not supported")
	}
	// TODO: Handle groups
	lines := [][]string{}
	for _, cmd := range cmds {
		if cmd.Hidden {
			continue
		}
		lines = append(lines, []string{
			"  ",
			ColorCommand(cmd.Path()),
			cmd.Help,
		})
	}
	printCard(w, "Commands", lines)

	// groupedCmds := collectCommandGroups(cmds)
	// for _, group := range groupedCmds {
	// 	w.Print("")
	// 	if group.Metadata.Title != "" {
	// 		w.Wrap(group.Metadata.Title)
	// 	}
	// 	if group.Metadata.Description != "" {
	// 		w.Indent().Wrap(group.Metadata.Description)
	// 		w.Print("")
	// 	}
	//
	// 	if w.Compact {
	// 		writeCompactCommandList(group.Commands, iw)
	// 	} else {
	// 		writeCommandList(group.Commands, iw)
	// 	}
	// }
}

func printCard(w *helpWriter, header string, lines [][]string) {
	printCardHeader(w, header)
	w.CardSection().PrintColumns(lines)
	printCardFooter(w)
}

func printCardHeader(w *helpWriter, title string) {
	padding := strings.Repeat("─", w.width-len(title)-7)
	w.Printf("╭─ %s ─%s─╮", ColorLow(title), padding)
}

func printCardFooter(w *helpWriter) {
	padding := strings.Repeat("─", w.width-2)
	w.Print(fmt.Sprintf("╰%s╯", padding))
}

// Directly from kong source code:

type helpFlagGroup struct {
	Metadata *kong.Group
	Flags    [][]*kong.Flag
}

func collectFlagGroups(flags [][]*kong.Flag) []helpFlagGroup {
	// Group keys in order of appearance.
	groups := []*kong.Group{}
	// Flags grouped by their group key.
	flagsByGroup := map[string][][]*kong.Flag{}

	for _, levelFlags := range flags {
		levelFlagsByGroup := map[string][]*kong.Flag{}

		for _, flag := range levelFlags {
			key := ""
			if flag.Group != nil {
				key = flag.Group.Key
				groupAlreadySeen := false
				for _, group := range groups {
					if key == group.Key {
						groupAlreadySeen = true
						break
					}
				}
				if !groupAlreadySeen {
					groups = append(groups, flag.Group)
				}
			}

			levelFlagsByGroup[key] = append(levelFlagsByGroup[key], flag)
		}

		for key, flags := range levelFlagsByGroup {
			flagsByGroup[key] = append(flagsByGroup[key], flags)
		}
	}

	out := []helpFlagGroup{}
	// Ungrouped flags are always displayed first.
	if ungroupedFlags, ok := flagsByGroup[""]; ok {
		out = append(out, helpFlagGroup{
			Metadata: &kong.Group{Title: "Flags:"},
			Flags:    ungroupedFlags,
		})
	}
	for _, group := range groups {
		out = append(out, helpFlagGroup{Metadata: group, Flags: flagsByGroup[group.Key]})
	}
	return out
}
