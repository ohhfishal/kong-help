# kong-help

Library to generate pretty help output for Go CLI's that use [alecthomas/kong](https://github.com/alecthomas/kong). Help output based heavily on the format used by [rendercv](https://github.com/rendercv/rendercv).

## Basic Example

```go
package main

import (
	"github.com/alecthomas/kong"
	konghelp "github.com/ohhfishal/kong-help"
	"os"
)

type CMD struct {
    // ...
}

func main() {
	var cli CMD
	kongCtx := kong.Parse(
		&cli,
        /* Just add this new option to kong.Parse */
		konghelp.Help(),
	)

	if err := kongCtx.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
```

A more detail example can be found in [_examples/readme/main.go](_examples/readme/main.go).

## Upcoming
- [ ] Better examples/docs
- [ ] Support for command grouping
- [ ] Support for all `kong.HelpOption`'s
- [ ] Help formatting options? (IE control of styling)
- [ ] Tests
