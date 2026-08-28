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
	// Build binary for integration testing
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
	if !strings.Contains(out, "jdiff - JSON Structural Diff") {
		t.Errorf("expected stdout to contain help title, got: %s", out)
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
	if out != "jdiff v0.3.0" {
		t.Errorf("expected stdout to be 'jdiff v0.3.0', got: %q", out)
	}
}

func TestBinaryNoArgsError(t *testing.T) {
	cmd := exec.Command(binPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected command to exit with non-zero code")
	}

	if !strings.Contains(stderr.String(), "jdiff - JSON Structural Diff") {
		t.Errorf("expected help output in stderr, got: %s", stderr.String())
	}
}

func TestBinaryDiffExecution(t *testing.T) {
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "old.json")
	newFile := filepath.Join(tmpDir, "new.json")

	oldJSON := `{
		"user": {
			"profile": {
				"name": "John",
				"phone": "123456789"
			}
		}
	}`
	newJSON := `{
		"user": {
			"profile": {
				"name": "James",
				"email": "james@example.com"
			}
		}
	}`

	if err := os.WriteFile(oldFile, []byte(oldJSON), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(newFile, []byte(newJSON), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binPath, oldFile, newFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected binary to exit 0, got error: %v, stderr: %s", err, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "MODIFIED\n  user.profile.name\n    - \"John\"\n    + \"James\"") {
		t.Errorf("expected user.profile.name modification in output, got:\n%s", out)
	}
	if !strings.Contains(out, "ADDED\n  user.profile.email\n    + \"james@example.com\"") {
		t.Errorf("expected user.profile.email added in output, got:\n%s", out)
	}
	if !strings.Contains(out, "REMOVED\n  user.profile.phone\n    - \"123456789\"") {
		t.Errorf("expected user.profile.phone removed in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Summary:\n  Added:     1\n  Removed:   1\n  Modified:  1") {
		t.Errorf("expected summary block in output, got:\n%s", out)
	}
}

func TestBinaryIdenticalDocuments(t *testing.T) {
	tmpDir := t.TempDir()
	f1 := filepath.Join(tmpDir, "f1.json")
	f2 := filepath.Join(tmpDir, "f2.json")

	content := `{"service": "auth", "version": 1}`
	if err := os.WriteFile(f1, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binPath, f1, f2)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected binary to exit 0 for identical docs, got error: %v, stderr: %s", err, stderr.String())
	}

	if strings.TrimSpace(stdout.String()) != "No differences found." {
		t.Errorf("expected 'No differences found.', got: %q", stdout.String())
	}
}
