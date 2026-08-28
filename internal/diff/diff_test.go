package diff

import (
	"strings"
	"testing"
)

// mockMatcher is a simple matcher for testing diff ignore options directly
type mockMatcher struct {
	ignored map[string]bool
}

func (m *mockMatcher) Matches(path Path) bool {
	return m.ignored[path.String()]
}

func TestIdenticalObjects(t *testing.T) {
	oldJSON := []byte(`{"name": "Alice", "age": 30, "active": true}`)
	newJSON := []byte(`{"name": "Alice", "age": 30, "active": true}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.HasChanges() {
		t.Fatalf("expected no changes, got %d changes", len(res.Changes))
	}
	out := res.String()
	if strings.TrimSpace(out) != "No differences found." {
		t.Fatalf("expected 'No differences found.', got %q", out)
	}
}

func TestIgnoreRulesInDiff(t *testing.T) {
	oldJSON := []byte(`{
		"name": "Project",
		"timestamp": "2026-01-01",
		"request_id": "req-123",
		"metadata": {
			"created": "2025",
			"updated": "2026"
		},
		"version": 1
	}`)
	newJSON := []byte(`{
		"name": "Project",
		"timestamp": "2026-08-01",
		"request_id": "req-456",
		"metadata": {
			"created": "2024",
			"updated": "2027"
		},
		"version": 2
	}`)

	matcher := &mockMatcher{
		ignored: map[string]bool{
			"timestamp":  true,
			"request_id": true,
			"metadata":   true,
		},
	}

	res, err := CompareBytesWithOptions(oldJSON, newJSON, Options{Matcher: matcher})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Modified()) != 1 {
		t.Fatalf("expected 1 modified change (version), got %d", len(res.Modified()))
	}
	if res.Modified()[0].Path.String() != "version" {
		t.Errorf("expected version to be modified, got: %s", res.Modified()[0].Path.String())
	}

	// Ignored: timestamp (1) + request_id (1) + metadata (2) = 4
	if res.Ignored != 4 {
		t.Errorf("expected 4 ignored changes, got %d", res.Ignored)
	}
}

func TestIgnoreAllChanges(t *testing.T) {
	oldJSON := []byte(`{"timestamp": "2026-01-01"}`)
	newJSON := []byte(`{"timestamp": "2026-08-01"}`)

	matcher := &mockMatcher{ignored: map[string]bool{"timestamp": true}}
	res, err := CompareBytesWithOptions(oldJSON, newJSON, Options{Matcher: matcher})
	if err != nil {
		t.Fatal(err)
	}

	if res.HasChanges() {
		t.Errorf("expected no changes, got %d", len(res.Changes))
	}
	if res.Ignored != 1 {
		t.Errorf("expected 1 ignored change, got %d", res.Ignored)
	}
	if strings.TrimSpace(res.String()) != "No differences found." {
		t.Errorf("expected 'No differences found.', got %q", res.String())
	}
}

func TestDeeplyNestedModification(t *testing.T) {
	oldJSON := []byte(`{
		"user": {
			"profile": {
				"contact": {
					"email": "old@example.com"
				}
			}
		}
	}`)
	newJSON := []byte(`{
		"user": {
			"profile": {
				"contact": {
					"email": "new@example.com"
				}
			}
		}
	}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Modified()) != 1 {
		t.Fatalf("expected 1 modified change, got %d", len(res.Modified()))
	}
	c := res.Modified()[0]
	if c.Path.String() != "user.profile.contact.email" {
		t.Errorf("expected path 'user.profile.contact.email', got %q", c.Path.String())
	}
	if c.OldValue != "old@example.com" || c.NewValue != "new@example.com" {
		t.Errorf("unexpected values: old=%v, new=%v", c.OldValue, c.NewValue)
	}
}

func TestArrayElementModification(t *testing.T) {
	oldJSON := []byte(`{"items": ["A", "B", "C"]}`)
	newJSON := []byte(`{"items": ["A", "X", "C"]}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Modified()) != 1 {
		t.Fatalf("expected 1 modified change, got %d", len(res.Modified()))
	}
	c := res.Modified()[0]
	if c.Path.String() != "items[1]" {
		t.Errorf("expected path 'items[1]', got %q", c.Path.String())
	}
	if c.OldValue != "B" || c.NewValue != "X" {
		t.Errorf("expected B -> X, got %v -> %v", c.OldValue, c.NewValue)
	}
}

func TestArrayAdditionAndRemoval(t *testing.T) {
	t.Run("root array addition", func(t *testing.T) {
		oldJSON := []byte(`["A", "B"]`)
		newJSON := []byte(`["A", "B", "C"]`)

		res, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Added()) != 1 || res.Added()[0].Path.String() != "[2]" {
			t.Errorf("expected added [2], got %+v", res.Added())
		}
	})

	t.Run("root array removal", func(t *testing.T) {
		oldJSON := []byte(`["A", "B", "C"]`)
		newJSON := []byte(`["A", "B"]`)

		res, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Removed()) != 1 || res.Removed()[0].Path.String() != "[2]" {
			t.Errorf("expected removed [2], got %+v", res.Removed())
		}
	})

	t.Run("empty to populated", func(t *testing.T) {
		oldJSON := []byte(`[]`)
		newJSON := []byte(`["A", "B"]`)

		res, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Added()) != 2 {
			t.Fatalf("expected 2 additions, got %d", len(res.Added()))
		}
		if res.Added()[0].Path.String() != "[0]" || res.Added()[1].Path.String() != "[1]" {
			t.Errorf("expected [0] and [1], got %s and %s", res.Added()[0].Path.String(), res.Added()[1].Path.String())
		}
	})

	t.Run("populated to empty", func(t *testing.T) {
		oldJSON := []byte(`["A", "B"]`)
		newJSON := []byte(`[]`)

		res, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Removed()) != 2 {
			t.Fatalf("expected 2 removals, got %d", len(res.Removed()))
		}
	})
}

func TestArrayOfObjects(t *testing.T) {
	oldJSON := []byte(`{
		"users": [
			{"id": 1, "name": "John"},
			{"id": 2, "name": "Mary"}
		]
	}`)
	newJSON := []byte(`{
		"users": [
			{"id": 1, "name": "James"},
			{"id": 2, "name": "Mary"},
			{"id": 3, "name": "David"}
		]
	}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Modified()) != 1 || res.Modified()[0].Path.String() != "users[0].name" {
		t.Errorf("expected modified users[0].name, got %+v", res.Modified())
	}
	if len(res.Added()) != 1 || res.Added()[0].Path.String() != "users[2]" {
		t.Errorf("expected added users[2], got %+v", res.Added())
	}
}

func TestNestedArraysAndObjects(t *testing.T) {
	oldJSON := []byte(`{
		"data": {
			"groups": [
				{
					"values": [10, 20, 30]
				}
			]
		}
	}`)
	newJSON := []byte(`{
		"data": {
			"groups": [
				{
					"values": [10, 25, 30]
				}
			]
		}
	}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Modified()) != 1 {
		t.Fatalf("expected 1 modified change, got %d", len(res.Modified()))
	}
	c := res.Modified()[0]
	if c.Path.String() != "data.groups[0].values[1]" {
		t.Errorf("expected path 'data.groups[0].values[1]', got %q", c.Path.String())
	}
}

func TestArrayElementTypes(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`["a", "b"]`), []byte(`["a", "c"]`))
		if len(res.Modified()) != 1 || res.Modified()[0].Path.String() != "[1]" {
			t.Errorf("unexpected: %+v", res.Modified())
		}
	})

	t.Run("numbers", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`[10, 20]`), []byte(`[10, 30]`))
		if len(res.Modified()) != 1 || res.Modified()[0].Path.String() != "[1]" {
			t.Errorf("unexpected: %+v", res.Modified())
		}
	})

	t.Run("booleans", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`[true, false]`), []byte(`[true, true]`))
		if len(res.Modified()) != 1 || res.Modified()[0].Path.String() != "[1]" {
			t.Errorf("unexpected: %+v", res.Modified())
		}
	})

	t.Run("null in array", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`[null, "val"]`), []byte(`["val", "val"]`))
		if len(res.Modified()) != 1 || res.Modified()[0].Path.String() != "[0]" {
			t.Errorf("unexpected: %+v", res.Modified())
		}
	})

	t.Run("nested matrix array", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`[[1, 2], [3, 4]]`), []byte(`[[1, 2], [3, 5]]`))
		if len(res.Modified()) != 1 || res.Modified()[0].Path.String() != "[1][1]" {
			t.Errorf("expected [1][1] modified, got: %+v", res.Modified())
		}
	})
}

func TestTypeConversions(t *testing.T) {
	tests := []struct {
		name    string
		oldJSON string
		newJSON string
		oldType JSONType
		newType JSONType
	}{
		{"string to number", `{"val": "20"}`, `{"val": 20}`, JSONTypeString, JSONTypeNumber},
		{"number to string", `{"val": 20}`, `{"val": "20"}`, JSONTypeNumber, JSONTypeString},
		{"boolean to string", `{"val": true}`, `{"val": "true"}`, JSONTypeBoolean, JSONTypeString},
		{"null to value", `{"val": null}`, `{"val": "hello"}`, JSONTypeNull, JSONTypeString},
		{"value to null", `{"val": 123}`, `{"val": null}`, JSONTypeNumber, JSONTypeNull},
		{"object to primitive", `{"val": {"k": "v"}}`, `{"val": "flat"}`, JSONTypeObject, JSONTypeString},
		{"primitive to object", `{"val": "flat"}`, `{"val": {"k": "v"}}`, JSONTypeString, JSONTypeObject},
		{"array to primitive", `{"val": [1, 2]}`, `{"val": 123}`, JSONTypeArray, JSONTypeNumber},
		{"primitive to array", `{"val": 123}`, `{"val": [1, 2]}`, JSONTypeNumber, JSONTypeArray},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CompareBytes([]byte(tt.oldJSON), []byte(tt.newJSON))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res.Modified()) != 1 {
				t.Fatalf("expected 1 modified change, got %d", len(res.Modified()))
			}
			c := res.Modified()[0]
			if c.OldType != tt.oldType || c.NewType != tt.newType {
				t.Errorf("expected types %s -> %s, got %s -> %s", tt.oldType, tt.newType, c.OldType, c.NewType)
			}
		})
	}
}

func TestSummaryAndFormatting(t *testing.T) {
	oldJSON := []byte(`{
		"name": "Project",
		"tags": ["go", "json", "cli"]
	}`)
	newJSON := []byte(`{
		"name": "Project",
		"tags": ["go", "diff", "cli", "opensource"]
	}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	summary := res.Summary()
	if summary.Added != 1 || summary.Removed != 0 || summary.Modified != 1 {
		t.Errorf("unexpected summary: %+v", summary)
	}

	out := res.String()
	if !strings.Contains(out, "MODIFIED\n  tags[1]\n    - \"json\"\n    + \"diff\"") {
		t.Errorf("expected tags[1] modification, got:\n%s", out)
	}
	if !strings.Contains(out, "ADDED\n  tags[3]\n    + \"opensource\"") {
		t.Errorf("expected tags[3] addition, got:\n%s", out)
	}
}

func TestDeterministicOrdering(t *testing.T) {
	oldJSON := []byte(`{"items": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}`)
	newJSON := []byte(`{"items": [10, 20, 30, 40, 50, 60, 70, 80, 90, 100]}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	mods := res.Modified()
	if len(mods) != 10 {
		t.Fatalf("expected 10 modifications, got %d", len(mods))
	}

	if mods[2].Path.String() != "items[2]" || mods[9].Path.String() != "items[9]" {
		t.Errorf("unexpected natural sorting order: %s, %s", mods[2].Path.String(), mods[9].Path.String())
	}
}

func TestEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name string
		json string
	}{
		{"empty object", `{}`},
		{"empty array", `[]`},
		{"nested empty array", `{"a": []}`},
		{"nested empty object", `{"a": {}}`},
		{"double nested empty object", `{"a": {"b": {}}}`},
		{"null field", `{"a": null}`},
		{"false field", `{"a": false}`},
		{"zero field", `{"a": 0}`},
		{"empty string field", `{"a": ""}`},
	}

	for _, ec := range edgeCases {
		t.Run(ec.name, func(t *testing.T) {
			res, err := CompareBytes([]byte(ec.json), []byte(ec.json))
			if err != nil {
				t.Fatalf("unexpected error comparing %s: %v", ec.name, err)
			}
			if res.HasChanges() {
				t.Errorf("expected no changes when comparing %s against itself, got: %+v", ec.name, res.Changes)
			}
		})
	}
}
