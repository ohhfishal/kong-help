package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	konghelp "github.com/ohhfishal/kong-help"
	"log/slog"
	"os"
	"time"
)

type CMD struct {
	LogConfig LogConfig `embed:"" group:"Logging Flags:"`

	File      *os.File      `help:"File"`
	Numbers   []int         `help:"Numbers to list"`
	Beat      time.Duration `default:"1s"`
	Verbosity int           `short:"v" type:"counter" help:"Set verbosity"`
	Default   string        `short:"d" enum:"a,b,c" default:"a" help:"Enum example flag (${enum})."`
	Bad       string        `help:"Here is a flag that is really long and has a lot of text in it." default:"bad"`
	Test      string        `default:"no help"`
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

type LogConfig struct {
	Disable     bool       `help:"Disable logging. Shorthand for handler=discard."`
	HandlerType string     `name:"handler" enum:"json,discard,text" env:"HANDLER" default:"json" help:"Handler to use (${enum}) (env=$$${env})"`
	Level       slog.Level `default:"INFO" enum:"WARN,ERROR,INFO,DEBUG" help:"Set logging level (${enum})."`
	AddSource   bool       `help:"Sets AddSource in the slog handler."`
	SetDefault  bool       `default:"true" help:"Set the global slog logger to use this config."`
}

func main() {
	var cli CMD
	kongCtx := kong.Parse(
		&cli,
		konghelp.Help(),
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
