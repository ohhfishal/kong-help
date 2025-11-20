package main

import (
	"errors"
	"fmt"
	"io"
	"github.com/alecthomas/kong"
)


var _ kong.HelpPrinter = PrettyHelpPrinter

func PrettyHelpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	if ctx.Empty() {
		options.Summary = false
	}
	w := newHelpWriter(ctx, options)
	selected := ctx.Selected()
	if selected == nil {
		printApp(w, ctx.Model)
	} else {
		return errors.New("command help not supported")
		// printCommand(w, ctx.Model, selected)
	}
	return w.Write(ctx.Stdout)
	return kong.DefaultHelpPrinter(options, ctx)
}

func printApp(w *helpWriter, app *kong.Application) {
	if !w.Options.NoAppSummary {
		w.Printf("Usage: %s%s", app.Name, app.Summary())
	}
	// printNodeDetail(w, app.Node, true)
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


type helpWriter struct {
	Options kong.HelpOptions
}
func newHelpWriter(ctx *kong.Context, options kong.HelpOptions) *helpWriter {
	return &helpWriter{
		Options: options,
	}
}

func (w *helpWriter) Print(line string) {
	fmt.Println(line)
}
func (w *helpWriter) Printf(format string, args ...any) {
	fmt.Printf(format, args...)
	fmt.Println("")
}

func (w *helpWriter) Write(writer io.Writer) error {
	return nil
}
