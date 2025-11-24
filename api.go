package konghelp

import (
	"github.com/alecthomas/kong"
	"golang.org/x/term"
	"log/slog"
	"os"
)

// DefaultWidth is the fallback width if it can not be automatically determined.
const DefaultWidth = 80

type Options struct {
	// UseWidth controls the max width printed.
	// If UseWidth is 0, stdout's width is used instead if a terminal, otherwise 80.
	UseWidth int
	// TODO: Expose ShowBoolTypes to show BOOL on boolean flags/args (Currently false)
	width int
}

func Help(options ...Options) kong.Option {
	var opts Options
	if len(options) > 0 {
		opts = options[0]
	}

	opts.width = DefaultWidth
	if opts.UseWidth > 0 {
		opts.width = opts.UseWidth
	} else if stdout := (int)(os.Stdout.Fd()); term.IsTerminal(stdout) {
		width, _, err := term.GetSize(stdout)
		if err != nil {
			slog.Warn(
				"could not get the terminal width, using default",
				slog.Int("default", DefaultWidth),
				slog.Any("err", err),
			)
		}
		opts.width = width
	}
	return kong.Help(NewPrettyPrinter(opts))
}
