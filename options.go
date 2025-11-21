package konghelp

import (
	"github.com/alecthomas/kong"
)

type Options struct {
}

func Help(options ...Options) kong.Option {
	return kong.Help(PrettyHelpPrinter)
}
