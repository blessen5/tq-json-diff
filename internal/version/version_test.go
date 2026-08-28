package version

import (
	"strings"
	"testing"
)

func TestVersionShort(t *testing.T) {
	if Short() != "v1.0.0" {
		t.Errorf("expected Short() == %q, got %q", "v1.0.0", Short())
	}
}

func TestVersionInfo(t *testing.T) {
	info := Info()
	expectedSubstrings := []string{
		"jdiff v1.0.0",
		"commit:",
		"built at:",
		"go version:",
		"os/arch:",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(info, sub) {
			t.Errorf("expected Info() to contain %q, but got:\n%s", sub, info)
		}
	}
}
