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
	if !strings.Contains(out, "jdiff apply [options] <patch.json> <input.json>") {
		t.Errorf("expected stdout to document apply command, got: %s", out)
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
	if out != "jdiff v0.8.0" {
		t.Errorf("expected stdout to be 'jdiff v0.8.0', got: %q", out)
	}
}

func TestBinaryOutputPatch(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"title": "v1", "tags": ["a", "b"]}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"title": "v2", "tags": ["a", "c", "d"]}`), 0644)

	cmd := exec.Command(binPath, "--output", "patch", oldFile, newFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected command to exit 0, got err: %v, stderr: %s", err, stderr.String())
	}

	var parsed []map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatalf("patch output is not valid JSON array: %v\nOutput:\n%s", err, stdout.String())
	}
	if len(parsed) < 2 {
		t.Errorf("expected multiple patch operations, got: %+v", parsed)
	}
}

func TestBinaryRoundTripAndApply(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")
	patchFile := filepath.Join(tmpDir, "changes.json")
	patchedResultFile := filepath.Join(tmpDir, "result.json")

	oldJSON := `{
		"name": "Service",
		"version": 1,
		"features": ["auth", "logging", "metrics"],
		"database": {
			"host": "localhost",
			"port": 5432
		}
	}`
	newJSON := `{
		"name": "Service Pro",
		"version": 2,
		"features": ["auth", "tracing", "metrics", "caching"],
		"database": {
			"host": "db.prod.internal",
			"port": 5432,
			"pool_size": 20
		}
	}`

	_ = os.WriteFile(oldFile, []byte(oldJSON), 0644)
	_ = os.WriteFile(newFile, []byte(newJSON), 0644)

	// 1. Generate patch
	cmd1 := exec.Command(binPath, "--output", "patch", "--output-file", patchFile, oldFile, newFile)
	if out, err := cmd1.CombinedOutput(); err != nil {
		t.Fatalf("failed to generate patch: %v, output: %s", err, string(out))
	}

	// 2. Verify patch CLI command
	cmdVerify := exec.Command(binPath, "--verify-patch", oldFile, newFile)
	if out, err := cmdVerify.CombinedOutput(); err != nil {
		t.Fatalf("verify-patch failed: %v, output: %s", err, string(out))
	}

	// 3. Apply patch
	cmdApply := exec.Command(binPath, "apply", "--output-file", patchedResultFile, patchFile, oldFile)
	if out, err := cmdApply.CombinedOutput(); err != nil {
		t.Fatalf("apply command failed: %v, output: %s", err, string(out))
	}

	// 4. Compare patched result against new.json
	cmdDiff := exec.Command(binPath, patchedResultFile, newFile)
	out, err := cmdDiff.CombinedOutput()
	if err != nil {
		t.Fatalf("diff comparison between patched result and new.json failed: %v", err)
	}

	if strings.TrimSpace(string(out)) != "No differences found." {
		t.Errorf("expected patched document to equal new.json exactly, but got diff:\n%s", string(out))
	}
}
