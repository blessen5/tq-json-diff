package cli

import (
	"bytes"
	"encoding/json"
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
			if !strings.Contains(out, "--output") || !strings.Contains(out, "--output-file") {
				t.Errorf("expected stdout to document output options, got: %s", out)
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
			if out != "jdiff v0.7.0" {
				t.Errorf("expected stdout to be 'jdiff v0.7.0', got: %q", out)
			}
			if stderr.Len() > 0 {
				t.Errorf("expected stderr to be empty, got: %s", stderr.String())
			}
		})
	}
}

func TestCLIOutputJSON(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"name": "Alice", "age": 20}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"name": "Bob", "age": 20}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--output", "json", oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d, stderr: %s", ExitCodeOK, code, stderr.String())
	}

	var parsed map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}

	summary := parsed["summary"].(map[string]any)
	if summary["modified"].(float64) != 1 || summary["total"].(float64) != 1 {
		t.Errorf("unexpected summary: %+v", summary)
	}
}

func TestCLIOutputUnified(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"name": "Alice"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"name": "Bob"}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--output", "unified", "--no-color", oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
	}

	out := stdout.String()
	if !strings.Contains(out, "--- "+oldFile+"\n+++ "+newFile) {
		t.Errorf("missing unified header, got:\n%s", out)
	}
	if !strings.Contains(out, "@@ name\n- \"Alice\"\n+ \"Bob\"") {
		t.Errorf("missing diff hunk, got:\n%s", out)
	}
}

func TestCLIOutputUnsupported(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--output", "xml", "f1.json", "f2.json"})
	if code != ExitCodeError {
		t.Fatalf("expected error, got %d", code)
	}

	if !strings.Contains(stderr.String(), "unsupported output format: xml") {
		t.Errorf("expected error message in stderr, got: %s", stderr.String())
	}
}

func TestCLIOutputFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")
	outDest := filepath.Join(tmpDir, "result.json")

	_ = os.WriteFile(oldFile, []byte(`{"k": "v1"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"k": "v2"}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--output", "json", "--output-file", outDest, oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d, stderr: %s", ExitCodeOK, code, stderr.String())
	}

	if stdout.Len() > 0 {
		t.Errorf("expected stdout to be empty when writing to output-file, got: %s", stdout.String())
	}

	content, err := os.ReadFile(outDest)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("output file content is not valid JSON: %v", err)
	}
}

func TestCLIStdin(t *testing.T) {
	tmpDir := t.TempDir()
	otherFile := filepath.Join(tmpDir, "other.json")
	_ = os.WriteFile(otherFile, []byte(`{"val": 20}`), 0644)

	t.Run("stdin as old", func(t *testing.T) {
		stdin := strings.NewReader(`{"val": 10}`)
		var stdout, stderr bytes.Buffer
		c := NewWithStdin(stdin, &stdout, &stderr)

		code := c.Run([]string{"--no-color", "-", otherFile})
		if code != ExitCodeOK {
			t.Fatalf("expected exit code %d, got %d, stderr: %s", ExitCodeOK, code, stderr.String())
		}
		if !strings.Contains(stdout.String(), "MODIFIED\n  val\n    - 10\n    + 20") {
			t.Errorf("unexpected stdout:\n%s", stdout.String())
		}
	})

	t.Run("stdin as new", func(t *testing.T) {
		stdin := strings.NewReader(`{"val": 30}`)
		var stdout, stderr bytes.Buffer
		c := NewWithStdin(stdin, &stdout, &stderr)

		code := c.Run([]string{"--no-color", otherFile, "-"})
		if code != ExitCodeOK {
			t.Fatalf("expected exit code %d, got %d, stderr: %s", ExitCodeOK, code, stderr.String())
		}
		if !strings.Contains(stdout.String(), "MODIFIED\n  val\n    - 20\n    + 30") {
			t.Errorf("unexpected stdout:\n%s", stdout.String())
		}
	})

	t.Run("both inputs stdin error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		c := New(&stdout, &stderr)

		code := c.Run([]string{"-", "-"})
		if code != ExitCodeError {
			t.Fatalf("expected exit code %d, got %d", ExitCodeError, code)
		}
		if !strings.Contains(stderr.String(), "cannot read both inputs from stdin") {
			t.Errorf("expected error message in stderr, got: %s", stderr.String())
		}
	})
}

func TestCLIIgnoreFlag(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"name": "App", "timestamp": "2026-01-01"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"name": "App", "timestamp": "2026-08-01"}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--ignore", "timestamp", oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
	}

	out := strings.TrimSpace(stdout.String())
	if out != "No differences found." {
		t.Errorf("expected 'No differences found.', got: %q", out)
	}
}

func TestCLIShowConfig(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		c := New(&stdout, &stderr)

		code := c.Run([]string{"--show-config"})
		if code != ExitCodeOK {
			t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
		}
		if strings.TrimSpace(stdout.String()) != "No ignore rules configured." {
			t.Errorf("unexpected output: %s", stdout.String())
		}
	})

	t.Run("with cli rules and config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		cfgFile := filepath.Join(tmpDir, "custom.json")
		_ = os.WriteFile(cfgFile, []byte(`{"ignore": ["*.updated_at", "users[*].id"]}`), 0644)

		var stdout, stderr bytes.Buffer
		c := New(&stdout, &stderr)

		code := c.Run([]string{"--config", cfgFile, "--ignore", "timestamp", "--show-config"})
		if code != ExitCodeOK {
			t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
		}
		out := stdout.String()
		if !strings.Contains(out, "Ignore rules:\n  timestamp\n  *.updated_at\n  users[*].id") {
			t.Errorf("expected merged active rules in stdout, got:\n%s", out)
		}
	})
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
