package diff

import (
	"strings"
	"testing"
)

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

func TestDeeplyNestedAdditionAndRemoval(t *testing.T) {
	oldJSON := []byte(`{
		"app": {
			"server": {
				"legacy_port": 80
			}
		}
	}`)
	newJSON := []byte(`{
		"app": {
			"server": {
				"ssl_port": 443
			}
		}
	}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Added()) != 1 || res.Added()[0].Path.String() != "app.server.ssl_port" {
		t.Errorf("expected added app.server.ssl_port, got %+v", res.Added())
	}
	if len(res.Removed()) != 1 || res.Removed()[0].Path.String() != "app.server.legacy_port" {
		t.Errorf("expected removed app.server.legacy_port, got %+v", res.Removed())
	}
}

func TestNewAndRemovedNestedObjects(t *testing.T) {
	oldJSON := []byte(`{
		"user": {
			"name": "John",
			"old_section": {
				"foo": "bar"
			}
		}
	}`)
	newJSON := []byte(`{
		"user": {
			"name": "John",
			"address": {
				"city": "Kochi"
			}
		}
	}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Added()) != 1 || res.Added()[0].Path.String() != "user.address" {
		t.Errorf("expected added user.address, got %+v", res.Added())
	}
	if len(res.Removed()) != 1 || res.Removed()[0].Path.String() != "user.old_section" {
		t.Errorf("expected removed user.old_section, got %+v", res.Removed())
	}
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

func TestRootJSONValues(t *testing.T) {
	t.Run("root primitives identical", func(t *testing.T) {
		res, err := CompareBytes([]byte(`"hello"`), []byte(`"hello"`))
		if err != nil {
			t.Fatal(err)
		}
		if res.HasChanges() {
			t.Errorf("expected no changes for identical root strings")
		}
	})

	t.Run("root primitive modified", func(t *testing.T) {
		res, err := CompareBytes([]byte(`10`), []byte(`20`))
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Modified()) != 1 || res.Modified()[0].Path.String() != "(root)" {
			t.Errorf("unexpected root diff: %+v", res.Modified())
		}
	})

	t.Run("root primitive vs object", func(t *testing.T) {
		res, err := CompareBytes([]byte(`"plain"`), []byte(`{"key": "val"}`))
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Modified()) != 1 || res.Modified()[0].Path.String() != "(root)" {
			t.Errorf("unexpected root diff: %+v", res.Modified())
		}
	})

	t.Run("root array modified", func(t *testing.T) {
		res, err := CompareBytes([]byte(`[1, 2]`), []byte(`[1, 3]`))
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Modified()) != 1 || res.Modified()[0].Path.String() != "(root)" {
			t.Errorf("unexpected root diff: %+v", res.Modified())
		}
	})

	t.Run("root array identical", func(t *testing.T) {
		res, err := CompareBytes([]byte(`[1, 2, "a"]`), []byte(`[1, 2, "a"]`))
		if err != nil {
			t.Fatal(err)
		}
		if res.HasChanges() {
			t.Errorf("expected no diff for identical root arrays")
		}
	})

	t.Run("root null vs value", func(t *testing.T) {
		res, err := CompareBytes([]byte(`null`), []byte(`true`))
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Modified()) != 1 {
			t.Errorf("expected 1 modified change")
		}
	})
}

func TestSummaryAndFormatting(t *testing.T) {
	oldJSON := []byte(`{
		"name": "Blessen",
		"age": 19,
		"city": "Kochi",
		"old_flag": true
	}`)
	newJSON := []byte(`{
		"name": "Blessen",
		"age": 20,
		"city": "Bengaluru",
		"country": "India"
	}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	summary := res.Summary()
	if summary.Added != 1 || summary.Removed != 1 || summary.Modified != 2 {
		t.Errorf("unexpected summary: %+v", summary)
	}

	out := res.String()
	expectedSubstrings := []string{
		"MODIFIED\n  age\n    - 19\n    + 20",
		"city\n    - \"Kochi\"\n    + \"Bengaluru\"",
		"ADDED\n  country\n    + \"India\"",
		"REMOVED\n  old_flag\n    - true",
		"Summary:\n  Added:     1\n  Removed:   1\n  Modified:  2",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(out, sub) {
			t.Errorf("expected output to contain %q, but got:\n%s", sub, out)
		}
	}
}

func TestDeterministicOrdering(t *testing.T) {
	oldJSON := []byte(`{"z": 1, "a": 2, "m": 3}`)
	newJSON := []byte(`{"z": 10, "a": 20, "m": 30}`)

	for i := 0; i < 5; i++ {
		res, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			t.Fatal(err)
		}
		mods := res.Modified()
		if len(mods) != 3 || mods[0].Path.String() != "a" || mods[1].Path.String() != "m" || mods[2].Path.String() != "z" {
			t.Fatalf("expected deterministic alphabetical sorting (a, m, z), got: %v, %v, %v",
				mods[0].Path.String(), mods[1].Path.String(), mods[2].Path.String())
		}
	}
}

func TestEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name string
		json string
	}{
		{"empty object", `{}`},
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

	t.Run("false vs true", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`{"a": false}`), []byte(`{"a": true}`))
		if len(res.Modified()) != 1 {
			t.Errorf("expected false vs true to be detected as modified")
		}
	})

	t.Run("zero vs empty string", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`{"a": 0}`), []byte(`{"a": ""}`))
		if len(res.Modified()) != 1 {
			t.Errorf("expected 0 vs empty string to be detected as modified type change")
		}
	})
}
