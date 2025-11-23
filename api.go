package konghelp

import (
	"github.com/alecthomas/kong"
)

type Options struct {
	// TODO: Expose ShowBoolTypes to show BOOL on boolean flags/args (Currently false)
}

func Help(options ...Options) kong.Option {
	return kong.Help(PrettyHelpPrinter)
}
