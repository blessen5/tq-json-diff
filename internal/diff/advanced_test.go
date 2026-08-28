package diff

import (
	"testing"
	"time"
)

func TestNumericTolerance(t *testing.T) {
	oldJSON := []byte(`{"price": 10.0001, "rate": 0.051}`)
	newJSON := []byte(`{"price": 10.0003, "rate": 0.052}`)

	// Without tolerance -> 2 changes
	res1, err := CompareBytes(oldJSON, newJSON)
	if err != nil || len(res1.Changes) != 2 {
		t.Fatalf("expected 2 changes without tolerance, got %d", len(res1.Changes))
	}

	// With tolerance 0.01 -> 0 changes
	res2, err := CompareBytesWithOptions(oldJSON, newJSON, Options{
		Tolerance: ToleranceOptions{
			NumericDelta: 0.01,
		},
	})
	if err != nil || res2.HasChanges() {
		t.Errorf("expected 0 changes within numeric tolerance, got: %d", len(res2.Changes))
	}
}

func TestTimestampTolerance(t *testing.T) {
	oldJSON := []byte(`{"updated_at": "2026-08-28T19:00:00Z"}`)
	newJSON := []byte(`{"updated_at": "2026-08-28T19:00:04Z"}`)

	// With 5s tolerance -> 0 changes
	res, err := CompareBytesWithOptions(oldJSON, newJSON, Options{
		Tolerance: ToleranceOptions{
			TimeDelta: 5 * time.Second,
		},
	})
	if err != nil || res.HasChanges() {
		t.Errorf("expected 0 changes within timestamp tolerance, got: %d", len(res.Changes))
	}
}

func TestSmartArrayAlignment(t *testing.T) {
	// Two objects in array with reordered IDs
	oldJSON := []byte(`{
		"users": [
			{"id": 1, "name": "Alice", "role": "admin"},
			{"id": 2, "name": "Bob", "role": "user"}
		]
	}`)
	newJSON := []byte(`{
		"users": [
			{"id": 2, "name": "Bob", "role": "user"},
			{"id": 1, "name": "Alice Cooper", "role": "admin"}
		]
	}`)

	res, err := CompareBytesWithOptions(oldJSON, newJSON, Options{
		ArrayStrategy: MatchAuto,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Smart alignment should match Alice to Alice (and find 1 change: name modified), and Bob to Bob (0 changes)
	if len(res.Changes) != 1 {
		t.Fatalf("expected exactly 1 change with auto array alignment, got %d: %+v", len(res.Changes), res.Changes)
	}
	if res.Changes[0].Path.String() != "users[1].name" || res.Changes[0].NewValue != "Alice Cooper" {
		t.Errorf("unexpected change: %+v", res.Changes[0])
	}
}
