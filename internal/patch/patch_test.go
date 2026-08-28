package patch

import (
	"encoding/json"
	"testing"

	"jdiff/internal/diff"
)

func TestPointerEscaping(t *testing.T) {
	tests := []struct {
		input    string
		escaped  string
		unescape string
	}{
		{"a/b", "a~1b", "a/b"},
		{"m~n", "m~0n", "m~n"},
		{"a/b~c", "a~1b~0c", "a/b~c"},
		{"simple", "simple", "simple"},
	}

	for _, tt := range tests {
		esc := Escape(tt.input)
		if esc != tt.escaped {
			t.Errorf("Escape(%q) = %q, expected %q", tt.input, esc, tt.escaped)
		}
		unesc, err := Unescape(esc)
		if err != nil || unesc != tt.unescape {
			t.Errorf("Unescape(%q) = %q, err: %v, expected %q", esc, unesc, err, tt.unescape)
		}
	}
}

func TestFromPath(t *testing.T) {
	p1 := diff.NewPath().AppendKey("users").AppendIndex(0).AppendKey("profile/name")
	ptr := FromPath(p1)
	expected := "/users/0/profile~1name"
	if ptr != expected {
		t.Errorf("expected pointer %q, got %q", expected, ptr)
	}

	pRoot := diff.NewPath()
	if FromPath(pRoot) != "" {
		t.Errorf("expected empty string for root pointer, got %q", FromPath(pRoot))
	}
}

func TestPatchGenerationAndApplication(t *testing.T) {
	oldJSON := []byte(`{
		"title": "Old Title",
		"count": 10,
		"active": true,
		"metadata": {
			"author": "Alice"
		},
		"items": ["A", "B", "C"]
	}`)
	newJSON := []byte(`{
		"title": "New Title",
		"count": 20,
		"metadata": {
			"author": "Alice Cooper",
			"editor": "Bob"
		},
		"items": ["A", "X", "C", "D"]
	}`)

	res, err := diff.CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}

	patch := Generate(res)
	if len(patch) == 0 {
		t.Fatalf("expected non-empty patch")
	}

	// Verify patch
	ok, err := Verify(oldJSON, newJSON, patch)
	if err != nil {
		t.Fatalf("verification error: %v", err)
	}
	if !ok {
		t.Errorf("expected verification to succeed")
	}
}

func TestMultiArrayRemovalsAndAdditions(t *testing.T) {
	oldJSON := []byte(`{"tags": ["zero", "one", "two", "three", "four"]}`)
	newJSON := []byte(`{"tags": ["zero", "ONE", "four", "five"]}`)

	res, err := diff.CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	patch := Generate(res)
	ok, err := Verify(oldJSON, newJSON, patch)
	if err != nil {
		t.Fatalf("verification failed: %v", err)
	}
	if !ok {
		t.Errorf("expected multi-item array patch application to match target exactly")
	}
}

func TestEmptyPatch(t *testing.T) {
	oldJSON := []byte(`{"a": 1}`)
	newJSON := []byte(`{"a": 1}`)

	res, err := diff.CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	patch := Generate(res)
	if len(patch) != 0 {
		t.Errorf("expected 0 patch operations for identical docs, got %d", len(patch))
	}

	ok, err := Verify(oldJSON, newJSON, patch)
	if err != nil || !ok {
		t.Errorf("expected verification to succeed on empty patch")
	}
}

func TestInvalidPatchApplication(t *testing.T) {
	doc := map[string]any{"name": "test"}

	t.Run("replace non-existent key", func(t *testing.T) {
		p := Patch{{Op: OpReplace, Path: "/missing", Value: 123}}
		_, err := Apply(doc, p)
		if err == nil {
			t.Errorf("expected error replacing missing key")
		}
	})

	t.Run("remove non-existent key", func(t *testing.T) {
		p := Patch{{Op: OpRemove, Path: "/missing"}}
		_, err := Apply(doc, p)
		if err == nil {
			t.Errorf("expected error removing missing key")
		}
	})

	t.Run("array index out of bounds", func(t *testing.T) {
		arrDoc := []any{"A", "B"}
		p := Patch{{Op: OpReplace, Path: "/10", Value: "Z"}}
		_, err := Apply(arrDoc, p)
		if err == nil {
			t.Errorf("expected error with out of bounds array index")
		}
	})
}

func TestMarshalJSONPatch(t *testing.T) {
	p := Patch{
		{Op: OpAdd, Path: "/key", Value: "val"},
		{Op: OpRemove, Path: "/old"},
		{Op: OpReplace, Path: "/title", Value: "new"},
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	var parsed []map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatal(err)
	}

	if len(parsed) != 3 {
		t.Fatalf("expected 3 operations, got %d", len(parsed))
	}
	if _, hasVal := parsed[1]["value"]; hasVal {
		t.Errorf("remove operation should not have value field: %+v", parsed[1])
	}
}
