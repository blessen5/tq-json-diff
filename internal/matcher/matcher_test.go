package matcher

import (
	"testing"

	"jdiff/internal/diff"
)

func TestExactPathMatching(t *testing.T) {
	m, err := New([]string{"timestamp", "metadata.created_at"})
	if err != nil {
		t.Fatalf("failed to create matcher: %v", err)
	}

	p1 := diff.NewPath().AppendKey("timestamp")
	if !m.Matches(p1) {
		t.Errorf("expected %s to match", p1.String())
	}

	p2 := diff.NewPath().AppendKey("metadata").AppendKey("created_at")
	if !m.Matches(p2) {
		t.Errorf("expected %s to match", p2.String())
	}

	p3 := diff.NewPath().AppendKey("name")
	if m.Matches(p3) {
		t.Errorf("expected %s NOT to match", p3.String())
	}
}

func TestSubtreeMatching(t *testing.T) {
	m, err := New([]string{"metadata"})
	if err != nil {
		t.Fatalf("failed to create matcher: %v", err)
	}

	pRoot := diff.NewPath().AppendKey("metadata")
	if !m.Matches(pRoot) {
		t.Errorf("expected root metadata to match")
	}

	pChild := diff.NewPath().AppendKey("metadata").AppendKey("nested").AppendKey("key")
	if !m.Matches(pChild) {
		t.Errorf("expected metadata.nested.key to match subtree rule")
	}

	pOther := diff.NewPath().AppendKey("data").AppendKey("metadata")
	if m.Matches(pOther) {
		t.Errorf("expected data.metadata NOT to match top-level metadata rule")
	}
}

func TestWildcardKeyMatching(t *testing.T) {
	m, err := New([]string{"*.timestamp", "users.*.email"})
	if err != nil {
		t.Fatalf("failed to create matcher: %v", err)
	}

	p1 := diff.NewPath().AppendKey("user").AppendKey("timestamp")
	if !m.Matches(p1) {
		t.Errorf("expected %s to match *.timestamp", p1.String())
	}

	p2 := diff.NewPath().AppendKey("order").AppendKey("timestamp")
	if !m.Matches(p2) {
		t.Errorf("expected %s to match *.timestamp", p2.String())
	}

	p3 := diff.NewPath().AppendKey("users").AppendKey("profile").AppendKey("email")
	if !m.Matches(p3) {
		t.Errorf("expected %s to match users.*.email", p3.String())
	}

	p4 := diff.NewPath().AppendKey("timestamp")
	if m.Matches(p4) {
		t.Errorf("expected %s NOT to match *.timestamp (needs 2 levels)", p4.String())
	}
}

func TestArrayExactAndWildcardMatching(t *testing.T) {
	m, err := New([]string{"users[0].session_id", "users[*].temporary_token", "[*].id"})
	if err != nil {
		t.Fatalf("failed to create matcher: %v", err)
	}

	p1 := diff.NewPath().AppendKey("users").AppendIndex(0).AppendKey("session_id")
	if !m.Matches(p1) {
		t.Errorf("expected %s to match users[0].session_id", p1.String())
	}

	p2 := diff.NewPath().AppendKey("users").AppendIndex(1).AppendKey("session_id")
	if m.Matches(p2) {
		t.Errorf("expected %s NOT to match users[0].session_id", p2.String())
	}

	p3 := diff.NewPath().AppendKey("users").AppendIndex(0).AppendKey("temporary_token")
	if !m.Matches(p3) {
		t.Errorf("expected %s to match users[*].temporary_token", p3.String())
	}

	p4 := diff.NewPath().AppendKey("users").AppendIndex(99).AppendKey("temporary_token")
	if !m.Matches(p4) {
		t.Errorf("expected %s to match users[*].temporary_token", p4.String())
	}

	p5 := diff.NewPath().AppendIndex(3).AppendKey("id")
	if !m.Matches(p5) {
		t.Errorf("expected %s to match [*].id", p5.String())
	}
}

func TestInvalidPatterns(t *testing.T) {
	invalid := []string{
		"",
		"   ",
		"users[abc]",
		"users[-1]",
		"prefix*suffix",
	}

	for _, raw := range invalid {
		_, err := ParsePattern(raw)
		if err == nil {
			t.Errorf("expected pattern %q to fail parsing, but it succeeded", raw)
		}
	}
}
