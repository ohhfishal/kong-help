# kong-help

Library to generate pretty help output for Go CLI's that use [alecthomas/kong](https://github.com/alecthomas/kong).


## Basic Example

<!-- TODO: Make sure the import is correct -->

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

## Upcoming
- [ ] Better examples/docs
- [ ] Adjustable width (and improved wrapping)
- [ ] Support for command grouping
- [ ] Support for all `kong.HelpOption`'s
- [ ] Help formatting options? (IE control of styling)
- [ ] Tests
