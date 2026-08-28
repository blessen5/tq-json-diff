package render

import (
	"bytes"
	"strings"
	"testing"

	"jdiff/internal/diff"
)

func sampleDiff(t *testing.T) *diff.DiffResult {
	oldJSON := []byte(`{
		"name": "Blessen",
		"age": 19,
		"city": "Kochi",
		"old_flag": true
	}`)
	newJSON := []byte(`{
		"name": "Blessen",
		"age": 20,
		"city": "Bengaluru",
		"country": "India"
	}`)

	res, err := diff.CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("failed to create sample diff: %v", err)
	}
	return res
}

func emptyDiff(t *testing.T) *diff.DiffResult {
	res, err := diff.CompareBytes([]byte(`{"a": 1}`), []byte(`{"a": 1}`))
	if err != nil {
		t.Fatalf("failed to create empty diff: %v", err)
	}
	return res
}

func TestRenderNoColor(t *testing.T) {
	res := sampleDiff(t)
	var buf bytes.Buffer

	opts := Options{Color: false}
	if err := Render(&buf, res, opts); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI escape sequences in no-color mode, got: %q", out)
	}

	expectedSubstrings := []string{
		"MODIFIED\n  age\n    - 19\n    + 20",
		"ADDED\n  country\n    + \"India\"",
		"REMOVED\n  old_flag\n    - true",
		"Summary:\n  Added:     1\n  Removed:   1\n  Modified:  2",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(out, sub) {
			t.Errorf("expected output to contain %q, but got:\n%s", sub, out)
		}
	}
}

func TestRenderColor(t *testing.T) {
	res := sampleDiff(t)
	var buf bytes.Buffer

	opts := Options{Color: true}
	if err := Render(&buf, res, opts); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape sequences in colored mode, got: %q", out)
	}
	if !strings.Contains(out, "\033[32mADDED\033[0m") {
		t.Errorf("expected green ADDED tag, got: %q", out)
	}
	if !strings.Contains(out, "\033[31mREMOVED\033[0m") {
		t.Errorf("expected red REMOVED tag, got: %q", out)
	}
	if !strings.Contains(out, "\033[33mMODIFIED\033[0m") {
		t.Errorf("expected yellow MODIFIED tag, got: %q", out)
	}
}

func TestRenderCompact(t *testing.T) {
	res := sampleDiff(t)
	var buf bytes.Buffer

	opts := Options{Color: false, Compact: true}
	if err := Render(&buf, res, opts); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	expectedLines := []string{
		"MODIFIED age: 19 → 20",
		"MODIFIED city: \"Kochi\" → \"Bengaluru\"",
		"ADDED country: \"India\"",
		"REMOVED old_flag: true",
	}

	for _, line := range expectedLines {
		if !strings.Contains(out, line) {
			t.Errorf("expected compact output to contain line %q, but got:\n%s", line, out)
		}
	}
	if !strings.Contains(out, "Summary:\n  Added:     1\n  Removed:   1\n  Modified:  2") {
		t.Errorf("expected compact output to contain summary, got:\n%s", out)
	}
}

func TestRenderVerbose(t *testing.T) {
	res := sampleDiff(t)
	var buf bytes.Buffer

	opts := Options{
		Color:   false,
		Verbose: true,
		OldFile: "old.json",
		NewFile: "new.json",
	}
	if err := Render(&buf, res, opts); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "Comparing:\n  Old: old.json\n  New: new.json\n\nChanges:") {
		t.Errorf("expected verbose header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "MODIFIED") || !strings.Contains(out, "Summary:") {
		t.Errorf("expected changes and summary in verbose output, got:\n%s", out)
	}
}

func TestRenderSummaryOnly(t *testing.T) {
	res := sampleDiff(t)
	var buf bytes.Buffer

	opts := Options{Color: false, SummaryOnly: true}
	if err := Render(&buf, res, opts); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Contains(out, "MODIFIED\n  age") {
		t.Errorf("summary-only output should suppress individual changes, got:\n%s", out)
	}

	expected := "JSON Diff Summary\n\nAdded:     1\nRemoved:   1\nModified:  2\nTotal:     4\n"
	if out != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestRenderSummaryOnlyNoDiff(t *testing.T) {
	res := emptyDiff(t)
	var buf bytes.Buffer

	opts := Options{Color: false, SummaryOnly: true}
	if err := Render(&buf, res, opts); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	expected := "JSON Diff Summary\n\nAdded:     0\nRemoved:   0\nModified:  0\nTotal:     0\n"
	if out != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, out)
	}
}

func TestRenderEmptyDiff(t *testing.T) {
	res := emptyDiff(t)
	var buf bytes.Buffer

	opts := Options{Color: false}
	if err := Render(&buf, res, opts); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if out != "No differences found.\n" {
		t.Errorf("expected 'No differences found.', got: %q", out)
	}
}
