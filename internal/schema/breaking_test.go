package schema

import (
	"bytes"
	"testing"

	"jdiff/internal/diff"
)

func TestAnalyzeBreakingChanges(t *testing.T) {
	oldJSON := []byte(`{
		"id": 100,
		"name": "Service",
		"deprecated_key": "legacy",
		"port": 8080,
		"active": true
	}`)
	newJSON := []byte(`{
		"id": 100,
		"name": "Service v2",
		"port": "8080",
		"active": true,
		"new_feature": true
	}`)

	res, err := diff.CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	report := Analyze(res)
	if !report.HasBreaking() {
		t.Errorf("expected breaking changes to be detected")
	}

	if report.BreakingCount != 2 { // deprecated_key removed (1), port type changed int->string (1)
		t.Errorf("expected 2 breaking changes, got %d", report.BreakingCount)
	}

	if report.AdditiveCount != 1 { // new_feature added (1)
		t.Errorf("expected 1 additive change, got %d", report.AdditiveCount)
	}

	if report.ValueCount != 1 { // name modified (1)
		t.Errorf("expected 1 value mutation, got %d", report.ValueCount)
	}

	var buf bytes.Buffer
	if err := report.Format(&buf, false); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("[BREAKING] deprecated_key")) || !bytes.Contains(buf.Bytes(), []byte("[ADDITIVE] new_feature")) {
		t.Errorf("unexpected report output:\n%s", out)
	}
}

func TestAnalyzeNonBreakingAdditiveChanges(t *testing.T) {
	oldJSON := []byte(`{"version": 1}`)
	newJSON := []byte(`{"version": 1, "extra": "compatible"}`)

	res, err := diff.CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	report := Analyze(res)
	if report.HasBreaking() {
		t.Errorf("expected 0 breaking changes for additive mutation, got %d", report.BreakingCount)
	}
	if report.AdditiveCount != 1 {
		t.Errorf("expected 1 additive change, got %d", report.AdditiveCount)
	}
}
