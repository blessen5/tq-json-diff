package tests

import (
	"strings"
	"testing"

	"jdiff/internal/diff"
)

func TestRegressionNullVsFalseVsZero(t *testing.T) {
	// null vs missing
	res1, _ := diff.CompareBytes([]byte(`{"k": null}`), []byte(`{}`))
	if len(res1.Changes) != 1 || res1.Changes[0].Type != diff.ChangeRemoved {
		t.Errorf("expected 1 removed change for null vs missing, got: %+v", res1.Changes)
	}

	// false vs zero
	res2, _ := diff.CompareBytes([]byte(`{"k": false}`), []byte(`{"k": 0}`))
	if len(res2.Changes) != 1 || res2.Changes[0].Type != diff.ChangeModified {
		t.Errorf("expected 1 modified change for false vs 0, got: %+v", res2.Changes)
	}

	// empty object vs empty array
	res3, _ := diff.CompareBytes([]byte(`{"k": {}}`), []byte(`{"k": []}`))
	if len(res3.Changes) != 1 || res3.Changes[0].Type != diff.ChangeModified {
		t.Errorf("expected 1 modified change for empty object vs empty array, got: %+v", res3.Changes)
	}
}

func TestRegressionUnicodeAndEscapes(t *testing.T) {
	oldDoc := []byte(`{"greeting": "Hello 🌍", "path/with/slash": 1, "path~with~tilde": 2}`)
	newDoc := []byte(`{"greeting": "Hello 世界", "path/with/slash": 10, "path~with~tilde": 20}`)

	res, err := diff.CompareBytes(oldDoc, newDoc)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Changes) != 3 {
		t.Errorf("expected 3 changes for Unicode and special character paths, got: %d", len(res.Changes))
	}
}

func TestRegressionMaxDepthEnforcement(t *testing.T) {
	// 5 nested levels with MaxDepth = 3
	doc1 := []byte(`{"l1": {"l2": {"l3": {"l4": {"l5": 1}}}}}`)
	doc2 := []byte(`{"l1": {"l2": {"l3": {"l4": {"l5": 2}}}}}`)

	_, err := diff.CompareBytesWithOptions(doc1, doc2, diff.Options{MaxDepth: 3})
	if err == nil {
		t.Fatalf("expected max recursion depth error, got nil")
	}
	if !strings.Contains(err.Error(), "maximum recursion depth (3) exceeded") {
		t.Errorf("unexpected error message: %v", err)
	}
}
