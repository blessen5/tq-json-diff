package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var binPath string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "jdiff-test-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	binName := "jdiff"
	if runtime.GOOS == "windows" {
		binName = "jdiff.exe"
	}
	binPath = filepath.Join(tmpDir, binName)

	cmd := exec.Command("go", "build", "-o", binPath, "jdiff")
	if err := cmd.Run(); err != nil {
		panic("failed to build jdiff binary for integration tests: " + err.Error())
	}

	os.Exit(m.Run())
}

func TestBinaryHelpFlag(t *testing.T) {
	cmd := exec.Command(binPath, "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "jdiff [options] <old.json> <new.json>") {
		t.Errorf("expected stdout to contain usage, got: %s", out)
	}
}

func TestBinaryVersionFlag(t *testing.T) {
	cmd := exec.Command(binPath, "--version")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	out := strings.TrimSpace(stdout.String())
	if out != "jdiff v0.5.0" {
		t.Errorf("expected stdout to be 'jdiff v0.5.0', got: %q", out)
	}
}

func TestBinaryArrayDiffExecution(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	oldJSON := `{
		"languages": ["Python", "Java", "C"]
	}`
	newJSON := `{
		"languages": ["Python", "Go", "C", "Rust"]
	}`

	_ = os.WriteFile(oldFile, []byte(oldJSON), 0644)
	_ = os.WriteFile(newFile, []byte(newJSON), 0644)

	cmd := exec.Command(binPath, "--no-color", oldFile, newFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "MODIFIED\n  languages[1]\n    - \"Java\"\n    + \"Go\"") {
		t.Errorf("expected languages[1] modified in output, got:\n%s", out)
	}
	if !strings.Contains(out, "ADDED\n  languages[3]\n    + \"Rust\"") {
		t.Errorf("expected languages[3] added in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Summary:\n  Added:     1\n  Removed:   0\n  Modified:  1") {
		t.Errorf("expected summary in output, got:\n%s", out)
	}
}

func TestBinaryCompactArrayDiff(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	oldJSON := `{"tags": ["go", "v1"]}`
	newJSON := `{"tags": ["go", "v2", "beta"]}`

	_ = os.WriteFile(oldFile, []byte(oldJSON), 0644)
	_ = os.WriteFile(newFile, []byte(newJSON), 0644)

	cmd := exec.Command(binPath, "--compact", "--no-color", oldFile, newFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "MODIFIED tags[1]: \"v1\" → \"v2\"") {
		t.Errorf("expected compact diff line, got: %s", out)
	}
	if !strings.Contains(out, "ADDED tags[2]: \"beta\"") {
		t.Errorf("expected compact added line, got: %s", out)
	}
}
