package cli

import (
	"bytes"
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

			out := stdout.String()
			if !strings.HasPrefix(out, "jdiff version") {
				t.Errorf("expected stdout to start with 'jdiff version', got: %s", out)
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

	if !strings.Contains(stderr.String(), "Error: two JSON file paths are required") {
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

	if !strings.Contains(stderr.String(), "Error: too many arguments") {
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

	if !strings.Contains(stderr.String(), "Error: unknown flag \"--foo\"") {
		t.Errorf("expected unknown flag error in stderr, got: %s", stderr.String())
	}
}

func TestCLIPositionalArgsPhase1(t *testing.T) {
	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"left.json", "right.json"})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
	}

	out := stdout.String()
	if !strings.Contains(out, "Comparing left.json and right.json...") {
		t.Errorf("expected comparison acknowledgment in stdout, got: %s", out)
	}
}
