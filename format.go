package main

import (
	"github.com/alecthomas/kong"
	"strings"
)

func PrettyValueFormatter(formatter kong.HelpValueFormatter) kong.HelpValueFormatter {
	return func(value *kong.Value) string {
		parts := []string{formatter(value)}
		if value.Tag != nil && value.Tag.Required {
			parts = append(parts, ColorRequired("[required]"))
		}
		return strings.Join(parts, " ")
	}
}
