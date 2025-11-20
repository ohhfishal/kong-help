package main

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"reflect"
	"strings"
)

var _ kong.HelpPrinter = PrettyHelpPrinter

// TODO: Make these configurable
var ColorExample = color.New(color.FgYellow).SprintFunc()
var ColorRequired = color.New(color.FgRed).SprintFunc()
var ColorDefault = color.New(color.FgMagenta).SprintFunc()
var ColorLow = color.HiBlackString
var ColorType = ColorExample

func PrettyHelpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	if ctx.Empty() {
		options.Summary = false
	}

	// TODO: Have this controlled via an option
	options.ValueFormatter = PrettyValueFormatter(options.ValueFormatter)

	w := newHelpWriter(ctx, options)
	selected := ctx.Selected()
	if selected == nil {
		printApp(w, ctx.Model)
	} else {
		return errors.New("command help not supported")
		// printCommand(w, ctx.Model, selected)
	}
	if _, err := w.WriteTo(ctx.Stdout); err != nil {
		return err
	}
	// return nil
	return kong.DefaultHelpPrinter(options, ctx)
}

func printApp(w *helpWriter, app *kong.Application) {
	if !w.Options.NoAppSummary {
		w.Print("")
		w.Printf("  %s: %s%s", ColorExample("Usuage"), app.Name, app.Summary())
	}
	printNodeDetail(w, app.Node, true)
	cmds := app.Leaves(true)
	if len(cmds) > 0 && app.HelpFlag != nil {
		w.Print("")
		if w.Options.Summary {
			w.Printf(`Run "%s --help" for more information.`, app.Name)
		} else {
			w.Printf(`Run "%s <command> --help" for more information on a command.`, app.Name)
		}
	}
}

func printCardHeader(w *helpWriter, title string) {
	w.Printf("╭─ %s ──────────────────────────────────────────────────────────────╮", ColorLow(title))
}
func printCardFooter(w *helpWriter) {
	w.Print("╰──────────────────────────────────────────────────────────────────────────╯")
}

func printPositionals(w *helpWriter, args []*kong.Positional) {
	rows := [][]string{}
	for _, arg := range args {
		var prefix = "  "
		if arg.Tag != nil && arg.Tag.Required {
			prefix = ColorRequired("* ")
		}
		components := []string{
			prefix,
			arg.Name,
			ColorType(strings.ToUpper(arg.Target.Kind().String())),
			// TODO: Write a custom ValueFormatter to do: "Help Message. [required] [default=""] etc
			w.Options.ValueFormatter(arg),
		}
		rows = append(rows, components)
	}
	w.PrintColumns(rows)
}

func printFlags(w *helpWriter, flags [][]*kong.Flag) {
	printCardHeader(w, "Options")

	rows := [][]string{}
	for _, groups := range flags {
		// TODO: Support groups
		for _, flag := range groups {
			value := flag.Value
			if value == nil {
				continue
			}

			var prefix = "  "
			if value.Tag != nil && value.Tag.Required {
				prefix = ColorRequired("* ")
			}
			var flagStr = "  "
			if flag.Short != 0 {
				flagStr = fmt.Sprintf("-%c", flag.Short)
			}

			if value.Name != "" {
				if flagStr == "  " {
					flagStr += "  --" + value.Name
				} else {
					flagStr += ", --" + value.Name
				}
				if tag := value.Tag; tag != nil && tag.HasDefault {
					var q string
					if value.Target.Kind() == reflect.String {
						q = `"`
					}
					flagStr += fmt.Sprintf(`=%s`, ColorDefault(q+tag.Default+q))
				}
			}

			components := []string{
				prefix,
				flagStr,
				ColorType(strings.ToUpper(value.Target.Kind().String())),
				// TODO: Write a custom ValueFormatter to do: "Help Message. [required] [default=""] etc
				w.Options.ValueFormatter(value),
			}
			rows = append(rows, components)
		}
	}
	w.CardSection().PrintColumns(rows)

	printCardFooter(w)
}

func printNodeDetail(w *helpWriter, node *kong.Node, hide bool) {
	if node.Help != "" {
		w.Print("")
		w.Wrap(node.Help)
	}
	if w.Options.Summary {
		return
	}
	if node.Detail != "" {
		w.Print("")
		w.Wrap(node.Detail)
	}
	if len(node.Positional) > 0 {
		w.Print("")
		printCardHeader(w, "Arguments")
		printPositionals(w.CardSection(), node.Positional)
		printCardFooter(w)
	}
	if !w.Options.FlagsLast {
		printFlags(w, node.AllFlags(true))
	}
	// TODO: Print the commands here
	if w.Options.FlagsLast {
		printFlags(w, node.AllFlags(true))
	}
	// printFlags := func() {
	// 	if flags := node.AllFlags(true); len(flags) > 0 {
	// 		groupedFlags := collectFlagGroups(flags)
	// 		for _, group := range groupedFlags {
	// 			w.Print("")
	// 			if group.Metadata.Title != "" {
	// 				w.Wrap(group.Metadata.Title)
	// 			}
	// 			if group.Metadata.Description != "" {
	// 				w.Indent().Wrap(group.Metadata.Description)
	// 				w.Print("")
	// 			}
	// 			writeFlags(w.Indent(), group.Flags)
	// 		}
	// 	}
	// }
	// if !w.FlagsLast {
	// 	printFlags()
	// }
	// var cmds []*kong.Node
	// if w.Options.NoExpandSubcommands {
	// 	cmds = node.Children
	// } else {
	// 	cmds = node.Leaves(hide)
	// }
	// if len(cmds) > 0 {
	// 	iw := w.Indent()
	// 	if w.Options.Tree {
	// 		w.Print("")
	// 		w.Print("Commands:")
	// 		writeCommandTree(iw, node)
	// 	} else {
	// 		groupedCmds := collectCommandGroups(cmds)
	// 		for _, group := range groupedCmds {
	// 			w.Print("")
	// 			if group.Metadata.Title != "" {
	// 				w.Wrap(group.Metadata.Title)
	// 			}
	// 			if group.Metadata.Description != "" {
	// 				w.Indent().Wrap(group.Metadata.Description)
	// 				w.Print("")
	// 			}
	//
	// 			if w.Compact {
	// 				writeCompactCommandList(group.Commands, iw)
	// 			} else {
	// 				writeCommandList(group.Commands, iw)
	// 			}
	// 		}
	// 	}
	// }
	// if w.FlagsLast {
	// 	printFlags()
	// }
}
