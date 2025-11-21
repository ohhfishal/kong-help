package konghelp

import (
	"github.com/alecthomas/kong"
	"reflect"
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

func formatValue(value reflect.Value, showBool bool) string {
	switch value.Kind() {
	case reflect.Pointer:
		return formatType(value.Type().Elem())
	case reflect.Bool:
		if !showBool {
			return ""
		}
		fallthrough
	default:
		return formatType(value.Type())
	}
}

func formatType(t reflect.Type) string {
	name := t.Name()
	// Slices, maps etc
	if name == "" {
		name = t.String()
	}
	return ColorType(strings.ToUpper(name))
}
