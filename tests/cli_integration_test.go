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
	if !strings.Contains(out, "--ignore") || !strings.Contains(out, "--config") || !strings.Contains(out, "--show-config") {
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
	if out != "jdiff v0.6.0" {
		t.Errorf("expected stdout to be 'jdiff v0.6.0', got: %q", out)
	}
}

func TestBinaryIgnoreExecution(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	oldJSON := `{
		"name": "Project",
		"timestamp": "2026-01-01",
		"version": 1
	}`
	newJSON := `{
		"name": "Project",
		"timestamp": "2026-08-01",
		"version": 2
	}`

	_ = os.WriteFile(oldFile, []byte(oldJSON), 0644)
	_ = os.WriteFile(newFile, []byte(newJSON), 0644)

	cmd := exec.Command(binPath, "--ignore", "timestamp", "--no-color", oldFile, newFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	out := stdout.String()
	if strings.Contains(out, "timestamp") {
		t.Errorf("timestamp should be ignored, got:\n%s", out)
	}
	if !strings.Contains(out, "MODIFIED\n  version\n    - 1\n    + 2") {
		t.Errorf("expected version modification, got:\n%s", out)
	}
	if !strings.Contains(out, "Ignored:   1") {
		t.Errorf("expected Ignored: 1 in summary, got:\n%s", out)
	}
}

func TestBinaryConfigFileExecution(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "my-config.json")
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(cfgFile, []byte(`{"ignore": ["*.updated_at", "request_id"]}`), 0644)
	_ = os.WriteFile(oldFile, []byte(`{"request_id": "r1", "user": {"name": "A", "updated_at": "t1"}}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"request_id": "r2", "user": {"name": "B", "updated_at": "t2"}}`), 0644)

	cmd := exec.Command(binPath, "--config", cfgFile, "--no-color", oldFile, newFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	out := stdout.String()
	if strings.Contains(out, "request_id") || strings.Contains(out, "updated_at") {
		t.Errorf("expected request_id and updated_at to be ignored, got:\n%s", out)
	}
	if !strings.Contains(out, "MODIFIED\n  user.name\n    - \"A\"\n    + \"B\"") {
		t.Errorf("expected user.name modification, got:\n%s", out)
	}
	if !strings.Contains(out, "Ignored:   2") {
		t.Errorf("expected Ignored: 2 in summary, got:\n%s", out)
	}
}

func TestBinaryShowConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, "custom.json")
	_ = os.WriteFile(cfgFile, []byte(`{"ignore": ["timestamp", "metadata.*"]}`), 0644)

	cmd := exec.Command(binPath, "--config", cfgFile, "--ignore", "extra_rule", "--show-config")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "Ignore rules:\n  extra_rule\n  timestamp\n  metadata.*") {
		t.Errorf("unexpected --show-config output:\n%s", out)
	}
}
