package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"io"
	"strings"
)

// Ensure we implement the correct interfaces
var _ io.WriterTo = &helpWriter{}

type helpWriter struct {
	Options kong.HelpOptions
	lines   *[]string
	prefix  string
	suffix  string
	width   int
}

func newHelpWriter(ctx *kong.Context, options kong.HelpOptions) *helpWriter {
	lines := []string{}
	return &helpWriter{
		Options: options,
		// TODO: Get this from somewhere else
		lines: &lines,
		width: 100,
	}
}

func (w *helpWriter) PrintColumns(lines [][]string) {
	for _, line := range lines {
		// TODO: Use the width
		// TODO: Auto format width
		// TODO: Support wrapping of elements
		// TODO: Fix and make more robust
		// TODO: Ignore color strings in count...
		w.Print(fmt.Sprintf("%-72s", strings.Join(line, " ")))
	}
}

func (w *helpWriter) Print(line string) {
	*w.lines = append(*w.lines, fmt.Sprintf("%s%s%s", w.prefix, line, w.suffix))
}
func (w *helpWriter) Printf(format string, args ...any) {
	w.Print(fmt.Sprintf(format, args...))
}

func (w *helpWriter) WriteTo(writer io.Writer) (int64, error) {
	var count int64
	for _, line := range *w.lines {
		n, err := fmt.Fprintln(writer, line)
		if err != nil {
			return -1, err
		}
		count = count + (int64)(n)
	}
	return count, nil
}

// TODO: Have this wrap the line and print them
func (h *helpWriter) Wrap(text string) {
	h.Print(text)
	// w := bytes.NewBuffer(nil)
	// doc.ToText(w, strings.TrimSpace(text), "", "    ", h.width) //nolint:staticcheck // cross-package links not possible
	// for _, line := range strings.Split(strings.TrimSpace(w.String()), "\n") {
	// 	h.Print(line)
	// }
}

func (h *helpWriter) Indent() *helpWriter {
	return &helpWriter{
		Options: h.Options,
		prefix:  h.prefix + "  ",
		lines:   h.lines,
		width:   h.width - 2,
	}
}

func (h *helpWriter) CardSection() *helpWriter {
	return &helpWriter{
		Options: h.Options,
		prefix:  "│ " + h.prefix,
		suffix:  h.suffix + " │",
		lines:   h.lines,
		width:   h.width - 2,
	}
}
