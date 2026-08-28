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
	if !strings.Contains(out, "jdiff apply [options] <patch.json> <input.json>") {
		t.Errorf("expected stdout to document apply command, got: %s", out)
	}
	if !strings.Contains(out, "--stats") || !strings.Contains(out, "--max-file-size") {
		t.Errorf("expected stdout to document performance options, got: %s", out)
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
	if out != "jdiff v0.9.0" {
		t.Errorf("expected stdout to be 'jdiff v0.9.0', got: %q", out)
	}
}

func TestBinaryExitCodesAndQuiet(t *testing.T) {
	tmpDir := t.TempDir()
	f1 := filepath.Join(tmpDir, "f1.json")
	f2 := filepath.Join(tmpDir, "f2.json")
	f3 := filepath.Join(tmpDir, "f3.json")

	_ = os.WriteFile(f1, []byte(`{"a": 1}`), 0644)
	_ = os.WriteFile(f2, []byte(`{"a": 1}`), 0644)
	_ = os.WriteFile(f3, []byte(`{"a": 2}`), 0644)

	// Case 1: Identical -> exit 0
	cmd1 := exec.Command(binPath, "--quiet", f1, f2)
	if err := cmd1.Run(); err != nil {
		t.Fatalf("expected identical files to return exit code 0, got: %v", err)
	}

	// Case 2: Different -> exit 1
	cmd2 := exec.Command(binPath, "--quiet", f1, f3)
	err2 := cmd2.Run()
	if exitErr, ok := err2.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
		t.Fatalf("expected diff to return exit code 1, got: %v", err2)
	}
}

func TestBinaryStats(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"name": "v1"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"name": "v2"}`), 0644)

	cmd := exec.Command(binPath, "--stats", oldFile, newFile)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	_ = cmd.Run() // diff exists so exit code 1 is expected

	out := stdout.String()
	if !strings.Contains(out, "Comparison Statistics") || !strings.Contains(out, "Old size:") || !strings.Contains(out, "Total time:") {
		t.Errorf("expected statistics in output, got:\n%s", out)
	}
}

func TestBinaryMaxFileSize(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	_ = os.WriteFile(oldFile, []byte(`{"name": "this is more than ten bytes"}`), 0644)
	_ = os.WriteFile(newFile, []byte(`{"name": "this is also more than ten bytes"}`), 0644)

	cmd := exec.Command(binPath, "--max-file-size", "10B", oldFile, newFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 2 {
		t.Fatalf("expected exit code 2 on resource limit, got: %v", err)
	}
	if !strings.Contains(stderr.String(), "exceeds maximum allowed file size") {
		t.Errorf("expected error message in stderr, got: %s", stderr.String())
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
	_ = cmd1.Run() // exit code 1 because diff exists

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

	// 4. Compare patched result against new.json (should return 0)
	cmdDiff := exec.Command(binPath, patchedResultFile, newFile)
	out, err := cmdDiff.CombinedOutput()
	if err != nil {
		t.Fatalf("diff comparison between patched result and new.json failed: %v, output: %s", err, string(out))
	}

	if strings.TrimSpace(string(out)) != "No differences found." {
		t.Errorf("expected patched document to equal new.json exactly, but got diff:\n%s", string(out))
	}
}
