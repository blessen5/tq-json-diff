package render

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"jdiff/internal/diff"
)

func sampleDiffResult() *diff.DiffResult {
	p1 := diff.NewPath().AppendKey("name")
	p2 := diff.NewPath().AppendKey("country")
	p3 := diff.NewPath().AppendKey("old_tag")

	return &diff.DiffResult{
		Changes: []diff.Change{
			{Path: p1, Type: diff.ChangeModified, OldValue: "Alice", NewValue: "Bob"},
			{Path: p2, Type: diff.ChangeAdded, NewValue: "India"},
			{Path: p3, Type: diff.ChangeRemoved, OldValue: "legacy"},
		},
		Ignored: 2,
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
		err      bool
	}{
		{"", FormatHuman, false},
		{"human", FormatHuman, false},
		{"HUMAN", FormatHuman, false},
		{"json", FormatJSON, false},
		{"JSON", FormatJSON, false},
		{"unified", FormatUnified, false},
		{"xml", "", true},
		{"csv", "", true},
	}

	for _, tt := range tests {
		f, err := ParseFormat(tt.input)
		if (err != nil) != tt.err {
			t.Errorf("input %q: expected error %v, got %v", tt.input, tt.err, err)
		}
		if f != tt.expected {
			t.Errorf("input %q: expected format %q, got %q", tt.input, tt.expected, f)
		}
	}
}

func TestRenderJSON(t *testing.T) {
	res := sampleDiffResult()
	var buf bytes.Buffer
	err := Render(&buf, res, Options{Format: FormatJSON})
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("JSON output is not valid JSON: %v\nOutput:\n%s", err, buf.String())
	}

	summary, ok := parsed["summary"].(map[string]any)
	if !ok {
		t.Fatalf("missing summary in JSON output")
	}
	if summary["modified"].(float64) != 1 || summary["added"].(float64) != 1 || summary["removed"].(float64) != 1 || summary["ignored"].(float64) != 2 || summary["total"].(float64) != 3 {
		t.Errorf("unexpected summary values: %+v", summary)
	}

	changes, ok := parsed["changes"].([]any)
	if !ok || len(changes) != 3 {
		t.Fatalf("expected 3 changes in JSON output, got: %+v", parsed["changes"])
	}

	c0 := changes[0].(map[string]any)
	if c0["path"] != "name" || c0["type"] != "modified" || c0["old"] != "Alice" || c0["new"] != "Bob" {
		t.Errorf("unexpected change[0]: %+v", c0)
	}

	c1 := changes[1].(map[string]any)
	if c1["path"] != "country" || c1["type"] != "added" || c1["new"] != "India" || c1["old"] != nil {
		t.Errorf("unexpected change[1]: %+v", c1)
	}

	c2 := changes[2].(map[string]any)
	if c2["path"] != "old_tag" || c2["type"] != "removed" || c2["old"] != "legacy" || c2["new"] != nil {
		t.Errorf("unexpected change[2]: %+v", c2)
	}
}

func TestRenderJSONSummaryOnly(t *testing.T) {
	res := sampleDiffResult()
	var buf bytes.Buffer
	err := Render(&buf, res, Options{Format: FormatJSON, SummaryOnly: true})
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("JSON summary is not valid JSON: %v", err)
	}

	if _, ok := parsed["changes"]; ok {
		t.Errorf("summary only should not contain changes field")
	}
	if _, ok := parsed["summary"]; !ok {
		t.Errorf("missing summary object in output")
	}
}

func TestRenderJSONEmptyDiff(t *testing.T) {
	res := &diff.DiffResult{}
	var buf bytes.Buffer
	err := Render(&buf, res, Options{Format: FormatJSON})
	if err != nil {
		t.Fatal(err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("empty JSON diff is not valid JSON: %v", err)
	}

	changes := parsed["changes"].([]any)
	if len(changes) != 0 {
		t.Errorf("expected empty changes array, got: %v", changes)
	}
}

func TestRenderUnified(t *testing.T) {
	res := sampleDiffResult()
	var buf bytes.Buffer
	err := Render(&buf, res, Options{
		Format:  FormatUnified,
		OldPath: "old.json",
		NewPath: "new.json",
	})
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "--- old.json\n+++ new.json") {
		t.Errorf("missing unified header in output:\n%s", out)
	}
	if !strings.Contains(out, "@@ name\n- \"Alice\"\n+ \"Bob\"") {
		t.Errorf("missing modified section in unified output:\n%s", out)
	}
	if !strings.Contains(out, "@@ country\n+ \"India\"") {
		t.Errorf("missing added section in unified output:\n%s", out)
	}
	if !strings.Contains(out, "@@ old_tag\n- \"legacy\"") {
		t.Errorf("missing removed section in unified output:\n%s", out)
	}
}

func TestRenderNoColor(t *testing.T) {
	res := sampleDiffResult()
	var buf bytes.Buffer
	err := Render(&buf, res, Options{Color: false})
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI codes in no-color mode, got: %q", out)
	}
	if !strings.Contains(out, "MODIFIED\n  name\n    - \"Alice\"\n    + \"Bob\"") {
		t.Errorf("expected modified section in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Ignored:   2") {
		t.Errorf("expected Ignored: 2 in summary, got:\n%s", out)
	}
}

func TestRenderColor(t *testing.T) {
	res := sampleDiffResult()
	var buf bytes.Buffer
	err := Render(&buf, res, Options{Color: true})
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, colorYellow) || !strings.Contains(out, colorGreen) || !strings.Contains(out, colorRed) {
		t.Errorf("expected color ANSI codes in output, got:\n%s", out)
	}
}

func TestRenderCompact(t *testing.T) {
	res := sampleDiffResult()
	var buf bytes.Buffer
	err := Render(&buf, res, Options{Compact: true, Color: false})
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "MODIFIED name: \"Alice\" → \"Bob\"") {
		t.Errorf("expected compact line for modified, got:\n%s", out)
	}
	if !strings.Contains(out, "ADDED country: \"India\"") {
		t.Errorf("expected compact line for added, got:\n%s", out)
	}
	if !strings.Contains(out, "REMOVED old_tag: \"legacy\"") {
		t.Errorf("expected compact line for removed, got:\n%s", out)
	}
}

func TestRenderVerboseWithIgnores(t *testing.T) {
	res := sampleDiffResult()
	var buf bytes.Buffer
	err := Render(&buf, res, Options{
		Verbose:      true,
		Color:        false,
		OldPath:      "file1.json",
		NewPath:      "file2.json",
		IgnoredRules: []string{"timestamp", "*.updated_at"},
	})
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "Ignoring:\n  timestamp\n  *.updated_at\n\nComparing:\n  Old: file1.json\n  New: file2.json") {
		t.Errorf("expected verbose header with ignores, got:\n%s", out)
	}
}

func TestRenderSummaryOnly(t *testing.T) {
	res := sampleDiffResult()
	var buf bytes.Buffer
	err := Render(&buf, res, Options{SummaryOnly: true, Color: false})
	if err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Contains(out, "MODIFIED\n  name") {
		t.Errorf("summary only should not contain diffs, got:\n%s", out)
	}
	if !strings.Contains(out, "JSON Diff Summary\n\nAdded:     1\nRemoved:   1\nModified:  1\nIgnored:   2\nTotal:     3") {
		t.Errorf("expected full summary block, got:\n%s", out)
	}
}
