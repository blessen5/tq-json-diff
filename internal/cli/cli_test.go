package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"flag --help", []string{"--help"}},
		{"flag -h", []string{"-h"}},
		{"command help", []string{"help"}},
		{"embedded flag", []string{"file1.json", "--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			c := New(&stdout, &stderr)

			code := c.Run(tt.args)
			if code != ExitCodeOK {
				t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
			}

			out := stdout.String()
			if !strings.Contains(out, "jdiff [options] <old.json> <new.json>") {
				t.Errorf("expected stdout to contain usage, got: %s", out)
			}
			if !strings.Contains(out, "--no-color") || !strings.Contains(out, "--compact") ||
				!strings.Contains(out, "--verbose") || !strings.Contains(out, "--summary") {
				t.Errorf("expected stdout to document all options, got: %s", out)
			}
			if stderr.Len() > 0 {
				t.Errorf("expected stderr to be empty, got: %s", stderr.String())
			}
		})
	}
}

func TestCLIVersion(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"flag --version", []string{"--version"}},
		{"flag -v", []string{"-v"}},
		{"command version", []string{"version"}},
		{"embedded flag", []string{"file1.json", "--version"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			c := New(&stdout, &stderr)

			code := c.Run(tt.args)
			if code != ExitCodeOK {
				t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
			}

			out := strings.TrimSpace(stdout.String())
			if out != "jdiff v0.5.0" {
				t.Errorf("expected stdout to be 'jdiff v0.5.0', got: %q", out)
			}
			if stderr.Len() > 0 {
				t.Errorf("expected stderr to be empty, got: %s", stderr.String())
			}
		})
	}
}

func TestCLINoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{})
	if code != ExitCodeError {
		t.Fatalf("expected exit code %d, got %d", ExitCodeError, code)
	}

	if stdout.Len() > 0 {
		t.Errorf("expected stdout to be empty, got: %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "jdiff [options] <old.json> <new.json>") {
		t.Errorf("expected stderr to contain help text, got: %s", stderr.String())
	}
}

func TestCLINoColorFlag(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"k": "old"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"k": "new"}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--no-color", oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
	}

	out := stdout.String()
	if strings.Contains(out, "\033[") {
		t.Errorf("expected no ANSI codes when --no-color is specified, got: %q", out)
	}
}

func TestCLICompactFlag(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"items": ["A", "B"]}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"items": ["A", "C"]}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--compact", "--no-color", oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
	}

	out := stdout.String()
	if !strings.Contains(out, "MODIFIED items[1]: \"B\" → \"C\"") {
		t.Errorf("expected compact format line in stdout, got: %s", out)
	}
}

func TestCLIVerboseFlag(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"k": "old"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"k": "new"}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--verbose", "--no-color", oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
	}

	out := stdout.String()
	if !strings.Contains(out, "Comparing:\n  Old: "+oldFile+"\n  New: "+newFile) {
		t.Errorf("expected verbose file headers in stdout, got: %s", out)
	}
}

func TestCLISummaryFlag(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"a": [1, 2]}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"a": [1, 20, 30]}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--summary", oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
	}

	out := stdout.String()
	if strings.Contains(out, "MODIFIED\n  a") {
		t.Errorf("summary mode should not show individual diffs, got: %s", out)
	}
	if !strings.Contains(out, "JSON Diff Summary\n\nAdded:     1\nRemoved:   0\nModified:  1\nTotal:     2") {
		t.Errorf("expected summary block, got: %s", out)
	}
}

func TestCLIMissingArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"only_one.json"})
	if code != ExitCodeError {
		t.Fatalf("expected exit code %d, got %d", ExitCodeError, code)
	}

	if !strings.Contains(stderr.String(), "jdiff: missing input files") {
		t.Errorf("expected error message in stderr, got: %s", stderr.String())
	}
}

func TestCLITooManyArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"f1.json", "f2.json", "f3.json"})
	if code != ExitCodeError {
		t.Fatalf("expected exit code %d, got %d", ExitCodeError, code)
	}

	if !strings.Contains(stderr.String(), "jdiff: too many arguments provided") {
		t.Errorf("expected error message in stderr, got: %s", stderr.String())
	}
}

func TestCLIUnknownFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--foo", "f1.json"})
	if code != ExitCodeError {
		t.Fatalf("expected exit code %d, got %d", ExitCodeError, code)
	}

	if !strings.Contains(stderr.String(), "jdiff: unknown flag \"--foo\"") {
		t.Errorf("expected unknown flag error in stderr, got: %s", stderr.String())
	}
}
