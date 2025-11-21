package main

import (
	"fmt"
	// "math"
	"github.com/alecthomas/kong"
	"io"
	"log/slog"
	"regexp"
	"strconv"
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
		lines:   &lines,
		// TODO: Get this from somewhere else
		// TODO: Use $COLUMNS? golang.org/x/term .IsTerminal(os.Stdout.Fd())?
		width: 80,
	}
}

// NOTE: The first part of each row must have an identical length.
// Only the lasts part is truncated by word.
func (w *helpWriter) PrintColumns(lines [][]string) {
	if len(lines) == 0 {
		return
	}

	// Calculate the max width of each part (besides the last one)
	maxes := []int{-1}
	for _, parts := range lines {
		for i, part := range parts {
			if i == len(parts) || i == 0 {
				continue
			}
			part = stripAnsiCodes(part)
			if len(maxes) == i {
				maxes = append(maxes, len(part))
			}
			if maxes[i] < len(part) {
				maxes[i] = len(part)
			}
		}
	}

	// Pad out the columns
	for j, parts := range lines {
		for i, part := range parts {
			if i == len(parts)-1 || i == 0 {
				continue
			}
			visible := stripAnsiCodes(part)
			if width := maxes[i] - len(visible); width >= 1 {
				padded := part + strings.Repeat(" ", width)
				lines[j][i] = padded
			}
		}
	}

	for _, parts := range lines {
		// Checking if the whole line fits
		line := strings.Join(parts, " ")
		visible := stripAnsiCodes(line)
		if len(visible) <= w.width {
			padding := strings.Repeat(" ", w.width-len(visible))
			w.Print(line + padding)
			continue
		}

		// TODO: Support wrapping of elements
		// Checking if we just need to wrap the description
		line = strings.Join(parts[:len(parts)-1], " ")
		visible = stripAnsiCodes(line)
		padding := strings.Repeat(" ", w.width-len(visible)-len(" ..."))
		w.Print(line + " ..." + padding)
		slog.Warn("spliting line not implemented", "line", line)
		continue

		// TODO: Worst case
		// Checking if we need to aggresively truncate
		// slog.Error("terminal is way too small!", "line", line)
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
func (h *helpWriter) Wrap(text string) string {
	return text
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
		width:   h.width - 4,
	}
}

func VisibleMap(r rune) rune {
	switch {
	case strconv.IsPrint(r):
		return r
	default:
		return -1
	}
}

func stripAnsiCodes(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(s, "")
}
