package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIBreakingChanges(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"id": 1, "token": "abc"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"id": "one"}`), 0644) // token removed + id changed type

	var stdout, stderr bytes.Buffer
	app := New(&stdout, &stderr)

	code := app.Run([]string{"--breaking", oldFile, newFile})
	if code != ExitCodeDiff {
		t.Errorf("expected exit code 1 (breaking changes present), got %d", code)
	}

	out := stdout.String()
	if !strings.Contains(out, "INCOMPATIBLE") || !strings.Contains(out, "token") {
		t.Errorf("expected breaking change warning in output, got:\n%s", out)
	}
}

func TestCLITolerance(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"rate": 0.051}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"rate": 0.052}`), 0644)

	var stdout, stderr bytes.Buffer
	app := New(&stdout, &stderr)

	code := app.Run([]string{"--numeric-tolerance", "0.01", oldFile, newFile})
	if code != ExitCodeOK {
		t.Errorf("expected exit code 0 within tolerance, got %d", code)
	}
}

func TestCLIRollbackAndHTMLFormats(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"v": 1}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"v": 2}`), 0644)

	// Test rollback output
	var stdout1, stderr1 bytes.Buffer
	app1 := New(&stdout1, &stderr1)
	code1 := app1.Run([]string{"--output", "rollback", oldFile, newFile})
	if code1 != ExitCodeDiff {
		t.Errorf("expected exit code 1 for rollback output, got %d", code1)
	}
	if !strings.Contains(stdout1.String(), `"value": 1`) {
		t.Errorf("expected rollback patch restoring old value, got: %s", stdout1.String())
	}

	// Test HTML output
	var stdout2, stderr2 bytes.Buffer
	app2 := New(&stdout2, &stderr2)
	code2 := app2.Run([]string{"--output", "html", oldFile, newFile})
	if code2 != ExitCodeDiff {
		t.Errorf("expected exit code 1 for HTML output, got %d", code2)
	}
	if !strings.Contains(stdout2.String(), "<!DOCTYPE html>") {
		t.Errorf("expected valid HTML output, got: %s", stdout2.String())
	}
}
