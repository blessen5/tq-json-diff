package schema

import (
	"fmt"
	"io"
	"strings"

	"jdiff/internal/diff"
)

// Classification categorizes the semantic severity of a change for API/schema compatibility.
type Classification string

const (
	// SeverityBreaking indicates an incompatible schema modification (removed field or mutated type).
	SeverityBreaking Classification = "BREAKING"
	// SeverityAdditive indicates a backward-compatible addition (new field or array element).
	SeverityAdditive Classification = "ADDITIVE"
	// SeverityValueChange indicates a value change preserving identical data types.
	SeverityValueChange Classification = "VALUE_CHANGE"
)

// Finding represents a single classified schema delta.
type Finding struct {
	Path           diff.Path
	Classification Classification
	Description    string
	Change         diff.Change
}

// Report encapsulates breaking change and schema evolution metrics.
type Report struct {
	BreakingCount int
	AdditiveCount int
	ValueCount    int
	Findings      []Finding
}

// HasBreaking returns true if any breaking changes were discovered.
func (r *Report) HasBreaking() bool {
	return r != nil && r.BreakingCount > 0
}

// Analyze evaluates a DiffResult for schema evolution and breaking API changes.
func Analyze(result *diff.DiffResult) *Report {
	if result == nil || len(result.Changes) == 0 {
		return &Report{}
	}

	rep := &Report{}

	for _, c := range result.Changes {
		var f Finding
		f.Path = c.Path
		f.Change = c

		switch c.Type {
		case diff.ChangeRemoved:
			f.Classification = SeverityBreaking
			f.Description = fmt.Sprintf("Field removed (type: %s) — clients depending on this field will break", c.OldType)
			rep.BreakingCount++

		case diff.ChangeAdded:
			f.Classification = SeverityAdditive
			f.Description = fmt.Sprintf("New field added (type: %s) — backward-compatible extension", c.NewType)
			rep.AdditiveCount++

		case diff.ChangeModified:
			if c.OldType != c.NewType {
				f.Classification = SeverityBreaking
				f.Description = fmt.Sprintf("Type mutated from %s to %s — incompatible deserialization contract", c.OldType, c.NewType)
				rep.BreakingCount++
			} else {
				f.Classification = SeverityValueChange
				f.Description = fmt.Sprintf("Value modified (%s) — schema structure preserved", c.OldType)
				rep.ValueCount++
			}
		}

		rep.Findings = append(rep.Findings, f)
	}

	return rep
}

// Format writes a human-readable schema compatibility report.
func (r *Report) Format(w io.Writer, color bool) error {
	const (
		colorReset  = "\033[0m"
		colorRed    = "\033[31m"
		colorGreen  = "\033[32m"
		colorYellow = "\033[33m"
		colorCyan   = "\033[36m"
		colorBold   = "\033[1m"
	)

	var sb strings.Builder
	title := "Schema Compatibility & Breaking Change Analysis"
	if color {
		title = colorBold + title + colorReset
	}
	sb.WriteString(title + "\n\n")

	if len(r.Findings) == 0 {
		sb.WriteString("No schema modifications detected (100% backward-compatible).\n")
		_, err := fmt.Fprint(w, sb.String())
		return err
	}

	for _, f := range r.Findings {
		tag := string(f.Classification)
		pathStr := f.Path.String()
		if color {
			pathStr = colorCyan + pathStr + colorReset
			switch f.Classification {
			case SeverityBreaking:
				tag = colorRed + colorBold + tag + colorReset
			case SeverityAdditive:
				tag = colorGreen + colorBold + tag + colorReset
			case SeverityValueChange:
				tag = colorYellow + tag + colorReset
			}
		}
		sb.WriteString(fmt.Sprintf("[%s] %s\n  %s\n", tag, pathStr, f.Description))
	}

	sb.WriteString("\nCompatibility Summary:\n")
	if r.BreakingCount > 0 {
		breakTag := fmt.Sprintf("  Breaking Changes:  %d (INCOMPATIBLE)", r.BreakingCount)
		if color {
			breakTag = colorRed + colorBold + breakTag + colorReset
		}
		sb.WriteString(breakTag + "\n")
	} else {
		passTag := "  Breaking Changes:  0 (BACKWARD-COMPATIBLE)"
		if color {
			passTag = colorGreen + passTag + colorReset
		}
		sb.WriteString(passTag + "\n")
	}
	sb.WriteString(fmt.Sprintf("  Additive Changes:  %d\n", r.AdditiveCount))
	sb.WriteString(fmt.Sprintf("  Value Mutations:   %d\n", r.ValueCount))

	_, err := fmt.Fprint(w, sb.String())
	return err
}
