package version

import (
	"strings"
	"testing"
)

func TestVersionShort(t *testing.T) {
	if Short() != Version {
		t.Errorf("expected Short() == %q, got %q", Version, Short())
	}
}

func TestVersionInfo(t *testing.T) {
	info := Info()
	expectedSubstrings := []string{
		"jdiff version",
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
