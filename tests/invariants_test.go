package tests

import (
	"testing"

	"jdiff/internal/diff"
	"jdiff/internal/matcher"
	"jdiff/internal/patch"
)

// Invariant 1: Comparing identical JSON documents must produce zero changes.
func TestInvariantIdenticalDocumentsZeroChanges(t *testing.T) {
	cases := [][]byte{
		[]byte(`{}`),
		[]byte(`[]`),
		[]byte(`{"a": 1, "b": [1, 2, 3], "c": {"nested": true}}`),
		[]byte(`"string value"`),
		[]byte(`12345.67`),
		[]byte(`null`),
		[]byte(`false`),
	}

	for _, doc := range cases {
		res, err := diff.CompareBytes(doc, doc)
		if err != nil {
			t.Fatalf("diff error: %v", err)
		}
		if res.HasChanges() {
			t.Errorf("expected 0 changes for identical document, got %d", len(res.Changes))
		}
	}
}

// Invariant 2: Applying a generated patch transforms old document into new document.
func TestInvariantPatchRoundTrip(t *testing.T) {
	oldJSON := []byte(`{
		"title": "v1",
		"items": ["a", "b", "c"],
		"config": {"debug": false, "retries": 3}
	}`)
	newJSON := []byte(`{
		"title": "v2",
		"items": ["a", "X", "c", "d"],
		"config": {"debug": true, "timeout": 30}
	}`)

	res, err := diff.CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatalf("diff error: %v", err)
	}

	p := patch.Generate(res)
	ok, err := patch.Verify(oldJSON, newJSON, p)
	if err != nil || !ok {
		t.Fatalf("patch round-trip verification failed: err=%v, ok=%v", err, ok)
	}
}

// Invariant 3: Changing an ignored field must never create a reported change.
func TestInvariantIgnoredFieldsIsolated(t *testing.T) {
	oldJSON := []byte(`{"id": 1, "timestamp": "2026-01-01T00:00:00Z", "name": "App"}`)
	newJSON := []byte(`{"id": 1, "timestamp": "2026-08-28T12:00:00Z", "name": "App"}`)

	m, err := matcher.New([]string{"timestamp"})
	if err != nil {
		t.Fatal(err)
	}

	res, err := diff.CompareBytesWithOptions(oldJSON, newJSON, diff.Options{Matcher: m})
	if err != nil {
		t.Fatal(err)
	}

	if res.HasChanges() {
		t.Errorf("expected 0 reported changes when only ignored fields change, got: %d", len(res.Changes))
	}
	if res.Ignored != 1 {
		t.Errorf("expected ignored count to be 1, got %d", res.Ignored)
	}
}

// Invariant 4: Formatting, whitespace, and key order variations produce zero changes.
func TestInvariantFormattingAndKeyOrderIgnored(t *testing.T) {
	doc1 := []byte("{\n  \"b\": 2,\n  \"a\": 1\n}")
	doc2 := []byte("{\"a\":1,\"b\":2}")

	res, err := diff.CompareBytes(doc1, doc2)
	if err != nil {
		t.Fatal(err)
	}
	if res.HasChanges() {
		t.Errorf("expected zero changes for differently formatted documents, got: %d", len(res.Changes))
	}
}
