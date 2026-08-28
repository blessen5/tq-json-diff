package patch

import (
	"encoding/json"
	"testing"
)

func FuzzPointer(f *testing.F) {
	f.Add("/users/0/name")
	f.Add("/a~1b/c~0d")
	f.Add("")
	f.Add("/")
	f.Add("invalid_no_slash")
	f.Add("/~")
	f.Add("/~2")

	f.Fuzz(func(t *testing.T, ptr string) {
		tokens, err := ParsePointer(ptr)
		if err == nil && len(tokens) > 0 {
			for _, tok := range tokens {
				_ = Escape(tok)
			}
		}
	})
}

func FuzzPatchApply(f *testing.F) {
	f.Add([]byte(`{"a": 1}`), []byte(`[{"op": "replace", "path": "/a", "value": 2}]`))
	f.Add([]byte(`[1, 2]`), []byte(`[{"op": "add", "path": "/1", "value": 99}]`))
	f.Add([]byte(`{"items": ["a"]}`), []byte(`[{"op": "remove", "path": "/items/0"}]`))
	f.Add([]byte(`{}`), []byte(`invalid json patch`))

	f.Fuzz(func(t *testing.T, docBytes, patchBytes []byte) {
		var doc any
		if err := json.Unmarshal(docBytes, &doc); err != nil {
			return
		}

		var patchDoc Patch
		if err := json.Unmarshal(patchBytes, &patchDoc); err != nil {
			return
		}

		// Apply must never panic on arbitrary patches and documents
		_, _ = Apply(doc, patchDoc)
	})
}
