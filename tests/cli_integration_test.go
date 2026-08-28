package tests

import (
	"bytes"
	"encoding/json"
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
	if !strings.Contains(out, "--output") || !strings.Contains(out, "--output-file") {
		t.Errorf("expected stdout to document new options, got: %s", out)
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
	if out != "jdiff v0.7.0" {
		t.Errorf("expected stdout to be 'jdiff v0.7.0', got: %q", out)
	}
}

func TestBinaryOutputJSON(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"tags": ["go", "v1"]}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"tags": ["go", "v2"]}`), 0644)

	cmd := exec.Command(binPath, "--output", "json", oldFile, newFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	var parsed map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v, stdout: %s", err, stdout.String())
	}
}

func TestBinaryOutputUnified(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"tags": ["go", "v1"]}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"tags": ["go", "v2"]}`), 0644)

	cmd := exec.Command(binPath, "--output", "unified", "--no-color", oldFile, newFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "@@ tags[1]\n- \"v1\"\n+ \"v2\"") {
		t.Errorf("expected unified diff output, got:\n%s", out)
	}
}

func TestBinaryStdinPiping(t *testing.T) {
	tmpDir := t.TempDir()
	otherFile := filepath.Join(tmpDir, "target.json")
	_ = os.WriteFile(otherFile, []byte(`{"count": 100}`), 0644)

	cmd := exec.Command(binPath, "--output", "json", "-", otherFile)
	cmd.Stdin = strings.NewReader(`{"count": 50}`)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v, stderr: %s", err, stderr.String())
	}

	var parsed map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON from stdin piping: %v", err)
	}
}

func TestBinaryOutputFile(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")
	destFile := filepath.Join(tmpDir, "diff_out.json")

	_ = os.WriteFile(oldFile, []byte(`{"a": 1}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"a": 2}`), 0644)

	cmd := exec.Command(binPath, "--output", "json", "--output-file", destFile, oldFile, newFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit 0, got: %v, stderr: %s", err, stderr.String())
	}

	if stdout.Len() > 0 {
		t.Errorf("expected empty stdout, got: %s", stdout.String())
	}

	data, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("failed to read dest file: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("invalid json written to file: %v", err)
	}
}
