package main

import (
	"context"
	"fmt"
	"github.com/alecthomas/kong"
	"os"
	"os/signal"
	"syscall"
)

type CMD struct {
	Verbosity int    `short:"v" type:"counter" help:"Set verbosity"`
	Default   string `short:"d" enum:"a,b,c" default:"a" help:"Enum example flag (${enum})."`
	Test      string `default:"test"`
	Tree      struct {
		Left struct {
			Arg      int    `arg:"" required:"" help:"Number"`
			Filename string `arg:"" default:"-" type:"filecontent" help:"Filepath"`
		} `cmd:"" help:"Go left."`
		Right struct {
			Arg int `arg:"" required:"" help:"Number"`
		} `cmd:"" help:"Go right."`
	} `cmd:"" help:"Go down a tree"`
	Flat struct {
	} `cmd:""`
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	var cli CMD
	kongCtx := kong.Parse(
		&cli, kong.BindTo(ctx, new(context.Context)),
		kong.Help(PrettyHelpPrinter),
		kong.ConfigureHelp(kong.HelpOptions{
			WrapUpperBound: -1, // Uses terminal width
			// Tree: true,
			// NoExpandSubcommands: true,
		}),
	)

	if err := kongCtx.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
