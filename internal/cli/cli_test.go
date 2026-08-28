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
			if !strings.Contains(out, "jdiff - JSON Structural Diff") {
				t.Errorf("expected stdout to contain help title, got: %s", out)
			}
			if !strings.Contains(out, "Usage:\n  jdiff <old.json> <new.json>") {
				t.Errorf("expected stdout to contain usage, got: %s", out)
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
			if out != "jdiff v0.2.0" {
				t.Errorf("expected stdout to be 'jdiff v0.2.0', got: %q", out)
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
	if !strings.Contains(stderr.String(), "jdiff - JSON Structural Diff") {
		t.Errorf("expected stderr to contain help text, got: %s", stderr.String())
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

func TestCLIExecutionWithRealFiles(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	oldContent := `{"name": "Blessen", "age": 19, "city": "Kochi"}`
	newContent := `{"name": "Blessen", "age": 20, "city": "Bengaluru", "country": "India"}`

	if err := os.WriteFile(oldFile, []byte(oldContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newFile, []byte(newContent), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d. stderr: %s", ExitCodeOK, code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "MODIFIED:\n    age\n        - 19\n        + 20") {
		t.Errorf("expected modified age in output, got:\n%s", out)
	}
	if !strings.Contains(out, "city\n        - \"Kochi\"\n        + \"Bengaluru\"") {
		t.Errorf("expected modified city in output, got:\n%s", out)
	}
	if !strings.Contains(out, "ADDED:\n    country\n        + \"India\"") {
		t.Errorf("expected added country in output, got:\n%s", out)
	}
	if strings.Contains(out, "name") {
		t.Errorf("unchanged field name should not appear in output, got:\n%s", out)
	}
}

func TestCLIMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "existing.json")
	if err := os.WriteFile(existingFile, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	nonExistent := filepath.Join(tmpDir, "non_existent.json")
	code := c.Run([]string{nonExistent, existingFile})
	if code != ExitCodeError {
		t.Fatalf("expected exit code %d for missing file, got %d", ExitCodeError, code)
	}

	if !strings.Contains(stderr.String(), "failed to read") {
		t.Errorf("expected 'failed to read' in stderr, got: %s", stderr.String())
	}
}

func TestCLIInvalidJSONFile(t *testing.T) {
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "bad.json")
	validFile := filepath.Join(tmpDir, "good.json")

	if err := os.WriteFile(invalidFile, []byte(`{invalid`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(validFile, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{invalidFile, validFile})
	if code != ExitCodeError {
		t.Fatalf("expected exit code %d for invalid JSON, got %d", ExitCodeError, code)
	}

	if !strings.Contains(stderr.String(), "invalid JSON in") {
		t.Errorf("expected 'invalid JSON in' error in stderr, got: %s", stderr.String())
	}
}
