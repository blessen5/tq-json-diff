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
	if res.String() != "" {
		t.Fatalf("expected empty string output for identical docs, got %q", res.String())
	}
}

func TestAddedProperty(t *testing.T) {
	oldJSON := []byte(`{"name": "Alice"}`)
	newJSON := []byte(`{"name": "Alice", "country": "India"}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	c := res.Changes[0]
	if c.Type != ChangeAdded || c.Path != "country" || c.NewValue != "India" {
		t.Errorf("unexpected change: %+v", c)
	}

	out := res.String()
	expected := "ADDED:\n    country\n        + \"India\"\n"
	if out != expected {
		t.Errorf("expected output:\n%s\ngot:\n%s", expected, out)
	}
}

func TestRemovedProperty(t *testing.T) {
	oldJSON := []byte(`{"name": "Alice", "city": "Kochi"}`)
	newJSON := []byte(`{"name": "Alice"}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	c := res.Changes[0]
	if c.Type != ChangeRemoved || c.Path != "city" || c.OldValue != "Kochi" {
		t.Errorf("unexpected change: %+v", c)
	}

	out := res.String()
	expected := "REMOVED:\n    city\n        - \"Kochi\"\n"
	if out != expected {
		t.Errorf("expected output:\n%s\ngot:\n%s", expected, out)
	}
}

func TestModifiedProperty(t *testing.T) {
	oldJSON := []byte(`{"age": 19}`)
	newJSON := []byte(`{"age": 20}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	c := res.Changes[0]
	if c.Type != ChangeModified || c.Path != "age" {
		t.Errorf("unexpected change: %+v", c)
	}

	out := res.String()
	if !strings.Contains(out, "MODIFIED:\n    age\n        - 19\n        + 20") {
		t.Errorf("unexpected format:\n%s", out)
	}
}

func TestMultipleChanges(t *testing.T) {
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

	if len(res.Modified()) != 2 {
		t.Errorf("expected 2 modified changes, got %d", len(res.Modified()))
	}
	if len(res.Added()) != 1 {
		t.Errorf("expected 1 added change, got %d", len(res.Added()))
	}
	if len(res.Removed()) != 1 {
		t.Errorf("expected 1 removed change, got %d", len(res.Removed()))
	}

	out := res.String()
	if strings.Contains(out, "name") {
		t.Errorf("expected unchanged field 'name' not to be present in output, got:\n%s", out)
	}
	if !strings.Contains(out, "MODIFIED:") || !strings.Contains(out, "ADDED:") || !strings.Contains(out, "REMOVED:") {
		t.Errorf("expected all 3 sections in output, got:\n%s", out)
	}
}

func TestNestedObjectModification(t *testing.T) {
	oldJSON := []byte(`{"user": {"name": "John", "age": 20}}`)
	newJSON := []byte(`{"user": {"name": "James", "age": 20}}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	c := res.Changes[0]
	if c.Type != ChangeModified || c.Path != "user.name" {
		t.Errorf("unexpected change: %+v", c)
	}

	out := res.String()
	if !strings.Contains(out, "user.name") || !strings.Contains(out, "- \"John\"") || !strings.Contains(out, "+ \"James\"") {
		t.Errorf("unexpected output:\n%s", out)
	}
}

func TestNestedPropertyAdditionAndRemoval(t *testing.T) {
	oldJSON := []byte(`{"user": {"profile": {"temp": 123}}}`)
	newJSON := []byte(`{"user": {"profile": {"email": "test@example.com"}}}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Added()) != 1 || res.Added()[0].Path != "user.profile.email" {
		t.Errorf("unexpected added: %+v", res.Added())
	}
	if len(res.Removed()) != 1 || res.Removed()[0].Path != "user.profile.temp" {
		t.Errorf("unexpected removed: %+v", res.Removed())
	}
}

func TestPrimitiveTypes(t *testing.T) {
	t.Run("string change", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`{"k": "a"}`), []byte(`{"k": "b"}`))
		if len(res.Modified()) != 1 || res.Modified()[0].Path != "k" {
			t.Errorf("unexpected result: %+v", res)
		}
	})

	t.Run("number change", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`{"k": 10.5}`), []byte(`{"k": 20.5}`))
		if len(res.Modified()) != 1 || res.Modified()[0].Path != "k" {
			t.Errorf("unexpected result: %+v", res)
		}
	})

	t.Run("boolean change", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`{"k": true}`), []byte(`{"k": false}`))
		if len(res.Modified()) != 1 || res.Modified()[0].Path != "k" {
			t.Errorf("unexpected result: %+v", res)
		}
	})

	t.Run("null change", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`{"k": null}`), []byte(`{"k": "now_a_string"}`))
		if len(res.Modified()) != 1 || res.Modified()[0].Path != "k" {
			t.Errorf("unexpected result: %+v", res)
		}
		if FormatValue(res.Modified()[0].OldValue) != "null" {
			t.Errorf("expected null formatting, got %s", FormatValue(res.Modified()[0].OldValue))
		}
	})

	t.Run("both null", func(t *testing.T) {
		res, _ := CompareBytes([]byte(`{"k": null}`), []byte(`{"k": null}`))
		if res.HasChanges() {
			t.Errorf("expected no change for both null, got: %+v", res)
		}
	})
}

func TestArrayHandling(t *testing.T) {
	t.Run("array modified", func(t *testing.T) {
		oldJSON := []byte(`{"tags": ["go", "json"]}`)
		newJSON := []byte(`{"tags": ["go", "diff"]}`)

		res, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(res.Modified()) != 1 || res.Modified()[0].Path != "tags" {
			t.Fatalf("expected 1 modified change for tags, got %+v", res.Modified())
		}
		out := res.String()
		if !strings.Contains(out, "MODIFIED:\n    tags") {
			t.Errorf("unexpected output:\n%s", out)
		}
	})

	t.Run("array unchanged", func(t *testing.T) {
		oldJSON := []byte(`{"tags": ["go", "json"]}`)
		newJSON := []byte(`{"tags": ["go", "json"]}`)

		res, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.HasChanges() {
			t.Errorf("expected unchanged arrays not to produce diff, got %+v", res.Changes)
		}
	})
}

func TestInvalidJSON(t *testing.T) {
	_, errOld := CompareBytes([]byte(`{invalid`), []byte(`{}`))
	if errOld == nil || !strings.Contains(errOld.Error(), "old document") {
		t.Errorf("expected error indicating old document JSON error, got: %v", errOld)
	}

	_, errNew := CompareBytes([]byte(`{}`), []byte(`{invalid`))
	if errNew == nil || !strings.Contains(errNew.Error(), "new document") {
		t.Errorf("expected error indicating new document JSON error, got: %v", errNew)
	}
}

func TestTypeChange(t *testing.T) {
	oldJSON := []byte(`{"field": "string_val"}`)
	newJSON := []byte(`{"field": {"nested": "obj"}}`)

	res, err := CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Modified()) != 1 || res.Modified()[0].Path != "field" {
		t.Fatalf("expected type change to be reported as modified, got: %+v", res)
	}
}
