package render

import (
	"fmt"
	"io"
	"strings"

	"jdiff/internal/diff"
)

// ANSI color codes
const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
)

// Options specifies rendering configurations for diff presentation.
type Options struct {
	Color       bool
	Compact     bool
	Verbose     bool
	SummaryOnly bool
	OldFile     string
	NewFile     string
}

// colorize applies ANSI color codes if color is enabled.
func colorize(text, code string, enabled bool) string {
	if !enabled {
		return text
	}
	return code + text + ansiReset
}

// Render formats and writes the diff result according to the specified options.
func Render(w io.Writer, res *diff.DiffResult, opts Options) error {
	summary := res.Summary()
	total := summary.Added + summary.Removed + summary.Modified

	// 1. Summary-only mode
	if opts.SummaryOnly {
		var sb strings.Builder
		sb.WriteString("JSON Diff Summary\n\n")
		sb.WriteString(fmt.Sprintf("Added:     %d\n", summary.Added))
		sb.WriteString(fmt.Sprintf("Removed:   %d\n", summary.Removed))
		sb.WriteString(fmt.Sprintf("Modified:  %d\n", summary.Modified))
		sb.WriteString(fmt.Sprintf("Total:     %d\n", total))
		_, err := fmt.Fprint(w, sb.String())
		return err
	}

	var sb strings.Builder

	// 2. Verbose file comparison header
	if opts.Verbose {
		sb.WriteString("Comparing:\n")
		sb.WriteString(fmt.Sprintf("  Old: %s\n", opts.OldFile))
		sb.WriteString(fmt.Sprintf("  New: %s\n\n", opts.NewFile))
		if res.HasChanges() {
			sb.WriteString("Changes:\n")
		}
	}

	// 3. No changes case
	if !res.HasChanges() {
		sb.WriteString("No differences found.\n")
		_, err := fmt.Fprint(w, sb.String())
		return err
	}

	// 4. Compact presentation mode
	if opts.Compact {
		for _, c := range res.Changes {
			switch c.Type {
			case diff.ChangeModified:
				tag := colorize("MODIFIED", ansiYellow, opts.Color)
				path := colorize(c.Path.String(), ansiCyan, opts.Color)
				oldVal := colorize(diff.FormatValue(c.OldValue), ansiRed, opts.Color)
				newVal := colorize(diff.FormatValue(c.NewValue), ansiGreen, opts.Color)
				sb.WriteString(fmt.Sprintf("%s %s: %s → %s\n", tag, path, oldVal, newVal))
			case diff.ChangeAdded:
				tag := colorize("ADDED", ansiGreen, opts.Color)
				path := colorize(c.Path.String(), ansiCyan, opts.Color)
				newVal := colorize(diff.FormatValue(c.NewValue), ansiGreen, opts.Color)
				sb.WriteString(fmt.Sprintf("%s %s: %s\n", tag, path, newVal))
			case diff.ChangeRemoved:
				tag := colorize("REMOVED", ansiRed, opts.Color)
				path := colorize(c.Path.String(), ansiCyan, opts.Color)
				oldVal := colorize(diff.FormatValue(c.OldValue), ansiRed, opts.Color)
				sb.WriteString(fmt.Sprintf("%s %s: %s\n", tag, path, oldVal))
			}
		}
		sb.WriteString("\nSummary:\n")
		sb.WriteString(fmt.Sprintf("  Added:     %d\n", summary.Added))
		sb.WriteString(fmt.Sprintf("  Removed:   %d\n", summary.Removed))
		sb.WriteString(fmt.Sprintf("  Modified:  %d\n", summary.Modified))

		_, err := fmt.Fprint(w, sb.String())
		return err
	}

	// 5. Standard structured presentation mode
	var sections []string

	modified := res.Modified()
	if len(modified) > 0 {
		var sec strings.Builder
		sec.WriteString(colorize("MODIFIED", ansiYellow, opts.Color) + "\n")
		for i, c := range modified {
			pathStr := colorize(c.Path.String(), ansiCyan, opts.Color)
			oldStr := colorize("- "+diff.FormatValue(c.OldValue), ansiRed, opts.Color)
			newStr := colorize("+ "+diff.FormatValue(c.NewValue), ansiGreen, opts.Color)

			sec.WriteString(fmt.Sprintf("  %s\n", pathStr))
			sec.WriteString(fmt.Sprintf("    %s\n", oldStr))
			sec.WriteString(fmt.Sprintf("    %s", newStr))
			if i < len(modified)-1 {
				sec.WriteString("\n\n")
			}
		}
		sections = append(sections, sec.String())
	}

	added := res.Added()
	if len(added) > 0 {
		var sec strings.Builder
		sec.WriteString(colorize("ADDED", ansiGreen, opts.Color) + "\n")
		for i, c := range added {
			pathStr := colorize(c.Path.String(), ansiCyan, opts.Color)
			newStr := colorize("+ "+diff.FormatValue(c.NewValue), ansiGreen, opts.Color)

			sec.WriteString(fmt.Sprintf("  %s\n", pathStr))
			sec.WriteString(fmt.Sprintf("    %s", newStr))
			if i < len(added)-1 {
				sec.WriteString("\n\n")
			}
		}
		sections = append(sections, sec.String())
	}

	removed := res.Removed()
	if len(removed) > 0 {
		var sec strings.Builder
		sec.WriteString(colorize("REMOVED", ansiRed, opts.Color) + "\n")
		for i, c := range removed {
			pathStr := colorize(c.Path.String(), ansiCyan, opts.Color)
			oldStr := colorize("- "+diff.FormatValue(c.OldValue), ansiRed, opts.Color)

			sec.WriteString(fmt.Sprintf("  %s\n", pathStr))
			sec.WriteString(fmt.Sprintf("    %s", oldStr))
			if i < len(removed)-1 {
				sec.WriteString("\n\n")
			}
		}
		sections = append(sections, sec.String())
	}

	var sumSec strings.Builder
	sumSec.WriteString("Summary:\n")
	sumSec.WriteString(fmt.Sprintf("  Added:     %d\n", summary.Added))
	sumSec.WriteString(fmt.Sprintf("  Removed:   %d\n", summary.Removed))
	sumSec.WriteString(fmt.Sprintf("  Modified:  %d", summary.Modified))
	sections = append(sections, sumSec.String())

	sb.WriteString(strings.Join(sections, "\n\n") + "\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}
