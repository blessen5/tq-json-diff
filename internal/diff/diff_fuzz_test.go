package diff

import (
	"testing"
)

func FuzzCompareBytes(f *testing.F) {
	// Seed corpus
	f.Add([]byte(`{"a": 1}`), []byte(`{"a": 2}`))
	f.Add([]byte(`[1, 2, 3]`), []byte(`[1, 3]`))
	f.Add([]byte(`{"nested": {"val": true}}`), []byte(`{"nested": {"val": false}}`))
	f.Add([]byte(`"hello"`), []byte(`"world"`))
	f.Add([]byte(`null`), []byte(`123`))
	f.Add([]byte(`invalid json`), []byte(`{}`))
	f.Add([]byte(``), []byte(``))

	f.Fuzz(func(t *testing.T, oldJSON, newJSON []byte) {
		// Calling CompareBytes must never panic on arbitrary binary or malformed JSON input
		_, _ = CompareBytes(oldJSON, newJSON)
	})
}
