package konghelp

import (
	"github.com/alecthomas/kong"
	"strings"
)

func PrettyValueFormatter(formatter kong.HelpValueFormatter) kong.HelpValueFormatter {
	return func(value *kong.Value) string {
		parts := []string{formatter(value)}

		tag := value.Tag
		if tag == nil {
			return parts[0]
		}
		if tag.Required {
			parts = append(parts, ColorRequired("[required]"))
		}
		return strings.Join(parts, " ")
	}
}
