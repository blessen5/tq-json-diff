package patch

import (
	"testing"

	"jdiff/internal/diff"
)

func TestRollbackPatchApplication(t *testing.T) {
	oldJSON := []byte(`{
		"title": "Version 1",
		"items": ["alpha", "beta", "gamma"],
		"config": {"cache": true}
	}`)
	newJSON := []byte(`{
		"title": "Version 2",
		"items": ["alpha", "BETA", "delta"],
		"config": {"cache": false, "extra": "yes"}
	}`)

	res, err := diff.CompareBytes(oldJSON, newJSON)
	if err != nil {
		t.Fatal(err)
	}

	// Generate rollback patch (should transform newJSON -> oldJSON)
	rollbackPatch := GenerateRollback(res)
	if len(rollbackPatch) == 0 {
		t.Fatalf("expected non-empty rollback patch")
	}

	ok, err := Verify(newJSON, oldJSON, rollbackPatch)
	if err != nil || !ok {
		t.Fatalf("rollback verification failed: err=%v, ok=%v", err, ok)
	}
}
