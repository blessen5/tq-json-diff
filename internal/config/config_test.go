package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, ".jdiff.json")
	content := `{
		"ignore": [
			"timestamp",
			"request_id",
			"*.updated_at"
		]
	}`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"timestamp", "request_id", "*.updated_at"}
	if !reflect.DeepEqual(cfg.Ignore, expected) {
		t.Errorf("expected %+v, got %+v", expected, cfg.Ignore)
	}
}

func TestLoadMissingConfig(t *testing.T) {
	_, err := Load("nonexistent-config-file.json")
	if err == nil {
		t.Errorf("expected error loading nonexistent config, got nil")
	}
}

func TestLoadInvalidJSONConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(cfgPath, []byte(`{invalid json`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(cfgPath)
	if err == nil {
		t.Errorf("expected parse error for invalid JSON, got nil")
	}
}

func TestMergeRules(t *testing.T) {
	cli := []string{"timestamp", "request_id"}
	file := []string{"request_id", "*.updated_at", "users[*].id"}

	merged := Merge(cli, file)
	expected := []string{"timestamp", "request_id", "*.updated_at", "users[*].id"}

	if !reflect.DeepEqual(merged, expected) {
		t.Errorf("expected %v, got %v", expected, merged)
	}
}
