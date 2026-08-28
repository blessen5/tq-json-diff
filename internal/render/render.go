package render

import (
	"fmt"
	"io"
	"strings"

	"jdiff/internal/diff"
)

// ANSI color escape sequences.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Options configures terminal rendering presentation behavior.
type Options struct {
	Color        bool
	Compact      bool
	Verbose      bool
	SummaryOnly  bool
	OldPath      string
	NewPath      string
	IgnoredRules []string
}

// Render formats a DiffResult to the provided writer according to Options.
func Render(w io.Writer, result *diff.DiffResult, opts Options) error {
	// Case 1: Summary-Only presentation mode
	if opts.SummaryOnly {
		return renderSummaryOnly(w, result, opts)
	}

	// Case 2: No differences detected
	if !result.HasChanges() {
		if opts.Verbose {
			renderVerboseHeader(w, opts)
		}
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	var sb strings.Builder

	// Verbose Header
	if opts.Verbose {
		renderVerboseHeaderToString(&sb, opts)
		sb.WriteString("Changes:\n")
	}

	// Case 3: Compact presentation mode
	if opts.Compact {
		renderCompactToString(&sb, result, opts)
	} else {
		// Case 4: Standard hierarchical presentation mode
		renderStandardToString(&sb, result, opts)
	}

	// Summary Footer
	renderSummaryFooterToString(&sb, result, opts)

	_, err := fmt.Fprint(w, sb.String())
	return err
}

func renderVerboseHeader(w io.Writer, opts Options) {
	var sb strings.Builder
	renderVerboseHeaderToString(&sb, opts)
	_, _ = fmt.Fprint(w, sb.String())
}

func renderVerboseHeaderToString(sb *strings.Builder, opts Options) {
	if len(opts.IgnoredRules) > 0 {
		sb.WriteString("Ignoring:\n")
		for _, r := range opts.IgnoredRules {
			sb.WriteString(fmt.Sprintf("  %s\n", r))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("Comparing:\n")
	sb.WriteString(fmt.Sprintf("  Old: %s\n", opts.OldPath))
	sb.WriteString(fmt.Sprintf("  New: %s\n\n", opts.NewPath))
}

func renderCompactToString(sb *strings.Builder, result *diff.DiffResult, opts Options) {
	for _, c := range result.Changes {
		pathStr := c.Path.String()
		if opts.Color {
			pathStr = colorCyan + pathStr + colorReset
		}

		switch c.Type {
		case diff.ChangeModified:
			tag := "MODIFIED"
			if opts.Color {
				tag = colorYellow + colorBold + tag + colorReset
			}
			sb.WriteString(fmt.Sprintf("%s %s: %s → %s\n", tag, pathStr, diff.FormatValue(c.OldValue), diff.FormatValue(c.NewValue)))

		case diff.ChangeAdded:
			tag := "ADDED"
			if opts.Color {
				tag = colorGreen + colorBold + tag + colorReset
			}
			sb.WriteString(fmt.Sprintf("%s %s: %s\n", tag, pathStr, diff.FormatValue(c.NewValue)))

		case diff.ChangeRemoved:
			tag := "REMOVED"
			if opts.Color {
				tag = colorRed + colorBold + tag + colorReset
			}
			sb.WriteString(fmt.Sprintf("%s %s: %s\n", tag, pathStr, diff.FormatValue(c.OldValue)))
		}
	}
	sb.WriteString("\n")
}

func renderStandardToString(sb *strings.Builder, result *diff.DiffResult, opts Options) {
	modified := result.Modified()
	if len(modified) > 0 {
		tag := "MODIFIED"
		if opts.Color {
			tag = colorYellow + colorBold + tag + colorReset
		}
		sb.WriteString(tag + "\n")
		for i, c := range modified {
			pathStr := c.Path.String()
			if opts.Color {
				pathStr = colorCyan + pathStr + colorReset
			}
			sb.WriteString(fmt.Sprintf("  %s\n", pathStr))
			if opts.Color {
				sb.WriteString(fmt.Sprintf("    %s- %s%s\n", colorRed, diff.FormatValue(c.OldValue), colorReset))
				sb.WriteString(fmt.Sprintf("    %s+ %s%s", colorGreen, diff.FormatValue(c.NewValue), colorReset))
			} else {
				sb.WriteString(fmt.Sprintf("    - %s\n", diff.FormatValue(c.OldValue)))
				sb.WriteString(fmt.Sprintf("    + %s", diff.FormatValue(c.NewValue)))
			}
			if i < len(modified)-1 {
				sb.WriteString("\n\n")
			}
		}
		sb.WriteString("\n\n")
	}

	added := result.Added()
	if len(added) > 0 {
		tag := "ADDED"
		if opts.Color {
			tag = colorGreen + colorBold + tag + colorReset
		}
		sb.WriteString(tag + "\n")
		for i, c := range added {
			pathStr := c.Path.String()
			if opts.Color {
				pathStr = colorCyan + pathStr + colorReset
			}
			sb.WriteString(fmt.Sprintf("  %s\n", pathStr))
			if opts.Color {
				sb.WriteString(fmt.Sprintf("    %s+ %s%s", colorGreen, diff.FormatValue(c.NewValue), colorReset))
			} else {
				sb.WriteString(fmt.Sprintf("    + %s", diff.FormatValue(c.NewValue)))
			}
			if i < len(added)-1 {
				sb.WriteString("\n\n")
			}
		}
		sb.WriteString("\n\n")
	}

	removed := result.Removed()
	if len(removed) > 0 {
		tag := "REMOVED"
		if opts.Color {
			tag = colorRed + colorBold + tag + colorReset
		}
		sb.WriteString(tag + "\n")
		for i, c := range removed {
			pathStr := c.Path.String()
			if opts.Color {
				pathStr = colorCyan + pathStr + colorReset
			}
			sb.WriteString(fmt.Sprintf("  %s\n", pathStr))
			if opts.Color {
				sb.WriteString(fmt.Sprintf("    %s- %s%s", colorRed, diff.FormatValue(c.OldValue), colorReset))
			} else {
				sb.WriteString(fmt.Sprintf("    - %s", diff.FormatValue(c.OldValue)))
			}
			if i < len(removed)-1 {
				sb.WriteString("\n\n")
			}
		}
		sb.WriteString("\n\n")
	}
}

func renderSummaryFooterToString(sb *strings.Builder, result *diff.DiffResult, opts Options) {
	summary := result.Summary()
	sb.WriteString("Summary:\n")
	sb.WriteString(fmt.Sprintf("  Added:     %d\n", summary.Added))
	sb.WriteString(fmt.Sprintf("  Removed:   %d\n", summary.Removed))
	sb.WriteString(fmt.Sprintf("  Modified:  %d", summary.Modified))
	if summary.Ignored > 0 {
		sb.WriteString(fmt.Sprintf("\n  Ignored:   %d", summary.Ignored))
	}
	sb.WriteString("\n")
}

func renderSummaryOnly(w io.Writer, result *diff.DiffResult, opts Options) error {
	summary := result.Summary()
	var sb strings.Builder
	title := "JSON Diff Summary"
	if opts.Color {
		title = colorBold + title + colorReset
	}
	sb.WriteString(title + "\n\n")
	sb.WriteString(fmt.Sprintf("Added:     %d\n", summary.Added))
	sb.WriteString(fmt.Sprintf("Removed:   %d\n", summary.Removed))
	sb.WriteString(fmt.Sprintf("Modified:  %d\n", summary.Modified))
	if summary.Ignored > 0 {
		sb.WriteString(fmt.Sprintf("Ignored:   %d\n", summary.Ignored))
	}
	sb.WriteString(fmt.Sprintf("Total:     %d\n", summary.Total()))

	_, err := fmt.Fprint(w, sb.String())
	return err
}
