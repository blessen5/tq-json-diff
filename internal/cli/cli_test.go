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
			if !strings.Contains(out, "--stats") || !strings.Contains(out, "--max-file-size") || !strings.Contains(out, "--quiet") {
				t.Errorf("expected stdout to document new options, got: %s", out)
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
			if out != "jdiff v0.9.0" {
				t.Errorf("expected stdout to be 'jdiff v0.9.0', got: %q", out)
			}
			if stderr.Len() > 0 {
				t.Errorf("expected stderr to be empty, got: %s", stderr.String())
			}
		})
	}
}

func TestCLIExitCodes(t *testing.T) {
	tmpDir := t.TempDir()
	fileA := filepath.Join(tmpDir, "a.json")
	fileB := filepath.Join(tmpDir, "b.json")
	fileC := filepath.Join(tmpDir, "c.json")

	_ = os.WriteFile(fileA, []byte(`{"v": 1}`), 0644)
	_ = os.WriteFile(fileB, []byte(`{"v": 1}`), 0644)
	_ = os.WriteFile(fileC, []byte(`{"v": 2}`), 0644)

	var stdout, stderr bytes.Buffer

	// Identical files -> ExitCodeOK (0)
	c1 := New(&stdout, &stderr)
	if code := c1.Run([]string{fileA, fileB}); code != ExitCodeOK {
		t.Errorf("expected ExitCodeOK (0), got %d", code)
	}

	// Different files -> ExitCodeDiff (1)
	stdout.Reset()
	c2 := New(&stdout, &stderr)
	if code := c2.Run([]string{fileA, fileC}); code != ExitCodeDiff {
		t.Errorf("expected ExitCodeDiff (1), got %d", code)
	}

	// Non-existent file -> ExitCodeError (2)
	stderr.Reset()
	c3 := New(&stdout, &stderr)
	if code := c3.Run([]string{fileA, "missing.json"}); code != ExitCodeError {
		t.Errorf("expected ExitCodeError (2), got %d", code)
	}
}

func TestCLIQuietMode(t *testing.T) {
	tmpDir := t.TempDir()
	fileA := filepath.Join(tmpDir, "a.json")
	fileB := filepath.Join(tmpDir, "b.json")
	fileC := filepath.Join(tmpDir, "c.json")

	_ = os.WriteFile(fileA, []byte(`{"v": 1}`), 0644)
	_ = os.WriteFile(fileB, []byte(`{"v": 1}`), 0644)
	_ = os.WriteFile(fileC, []byte(`{"v": 2}`), 0644)

	var stdout, stderr bytes.Buffer

	c1 := New(&stdout, &stderr)
	code1 := c1.Run([]string{"--quiet", fileA, fileB})
	if code1 != ExitCodeOK || stdout.Len() > 0 {
		t.Errorf("expected ExitCodeOK (0) and empty stdout in quiet mode, got code: %d, stdout: %s", code1, stdout.String())
	}

	stdout.Reset()
	c2 := New(&stdout, &stderr)
	code2 := c2.Run([]string{"-q", fileA, fileC})
	if code2 != ExitCodeDiff || stdout.Len() > 0 {
		t.Errorf("expected ExitCodeDiff (1) and empty stdout in quiet mode, got code: %d, stdout: %s", code2, stdout.String())
	}
}

func TestCLIStats(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"v": 1}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"v": 2}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--stats", "--output", "json", oldFile, newFile})
	if code != ExitCodeDiff {
		t.Fatalf("expected ExitCodeDiff (1), got %d", code)
	}

	var parsed map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\nOutput: %s", err, stdout.String())
	}

	statsMap, ok := parsed["statistics"].(map[string]any)
	if !ok {
		t.Fatalf("expected statistics in JSON output: %+v", parsed)
	}

	if _, hasOld := statsMap["old_size"]; !hasOld {
		t.Errorf("missing old_size in statistics: %+v", statsMap)
	}
}

func TestCLIMaxFileSize(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"large_payload": "this is more than twenty bytes"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"large_payload": "this is also more than twenty bytes"}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--max-file-size", "20B", oldFile, newFile})
	if code != ExitCodeError {
		t.Fatalf("expected ExitCodeError (2) for file exceeding max-file-size, got %d", code)
	}

	if !strings.Contains(stderr.String(), "exceeds maximum allowed file size") {
		t.Errorf("expected error message in stderr, got: %s", stderr.String())
	}
}

func TestCLIMaxChanges(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"a": 10, "b": 20, "c": 30, "d": 40, "e": 50}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--max-changes", "2", "--no-color", oldFile, newFile})
	if code != ExitCodeDiff {
		t.Fatalf("expected ExitCodeDiff (1), got %d", code)
	}

	out := stdout.String()
	if !strings.Contains(out, "maximum changes limit 2 reached") {
		t.Errorf("expected truncation message, got: %s", out)
	}
}

func TestCLIOutputPatch(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"title": "old"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"title": "new"}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--output", "patch", oldFile, newFile})
	if code != ExitCodeDiff {
		t.Fatalf("expected exit code %d (diff), got %d, stderr: %s", ExitCodeDiff, code, stderr.String())
	}

	var patchOps []map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &patchOps); err != nil {
		t.Fatalf("expected valid JSON Patch array: %v", err)
	}

	if len(patchOps) != 1 || patchOps[0]["op"] != "replace" || patchOps[0]["path"] != "/title" || patchOps[0]["value"] != "new" {
		t.Errorf("unexpected patch op: %+v", patchOps)
	}
}

func TestCLIApplyCommand(t *testing.T) {
	tmpDir := t.TempDir()
	patchFile := filepath.Join(tmpDir, "patch.json")
	inputFile := filepath.Join(tmpDir, "input.json")
	outputFile := filepath.Join(tmpDir, "patched.json")

	_ = os.WriteFile(patchFile, []byte(`[{"op": "replace", "path": "/val", "value": 100}]`), 0644)
	_ = os.WriteFile(inputFile, []byte(`{"val": 50}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"apply", "--output-file", outputFile, patchFile, inputFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d, stderr: %s", ExitCodeOK, code, stderr.String())
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read patched output file: %v", err)
	}

	var parsed map[string]any
	_ = json.Unmarshal(data, &parsed)
	numVal, ok := parsed["val"].(json.Number)
	if !ok || numVal.String() != "100" {
		if floatVal, ok := parsed["val"].(float64); !ok || floatVal != 100 {
			t.Errorf("expected val to be 100, got: %v", parsed["val"])
		}
	}
}

func TestCLIVerifyPatch(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"a": 1, "b": [1, 2]}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"a": 2, "b": [1, 2, 3]}`), 0644)

	var stdout, stderr bytes.Buffer
	c := New(&stdout, &stderr)

	code := c.Run([]string{"--verify-patch", oldFile, newFile})
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d, stderr: %s", ExitCodeOK, code, stderr.String())
	}

	if !strings.Contains(stdout.String(), "Patch verification successful.") {
		t.Errorf("expected success message, got: %s", stdout.String())
	}
}
