package matcher

import (
	"fmt"
	"strconv"
	"strings"

	"jdiff/internal/diff"
)

// MatcherKind represents the type of segment matcher.
type MatcherKind int

const (
	// MatcherExactKey matches an exact object property key.
	MatcherExactKey MatcherKind = iota
	// MatcherWildcardKey matches any object property key (*).
	MatcherWildcardKey
	// MatcherExactIndex matches an exact array index ([0], [1]).
	MatcherExactIndex
	// MatcherWildcardIndex matches any array index ([*]).
	MatcherWildcardIndex
)

// SegmentMatcher matches a single diff.Segment.
type SegmentMatcher struct {
	Kind  MatcherKind
	Key   string
	Index int
}

// Match checks if a diff.Segment satisfies this segment matcher.
func (sm SegmentMatcher) Match(seg diff.Segment) bool {
	switch sm.Kind {
	case MatcherExactKey:
		return seg.Kind == diff.SegmentKindKey && seg.Key == sm.Key
	case MatcherWildcardKey:
		return seg.Kind == diff.SegmentKindKey
	case MatcherExactIndex:
		return seg.Kind == diff.SegmentKindIndex && seg.Index == sm.Index
	case MatcherWildcardIndex:
		return seg.Kind == diff.SegmentKindIndex
	default:
		return false
	}
}

// Pattern represents a compiled ignore rule made of segment matchers.
type Pattern struct {
	Raw      string
	Segments []SegmentMatcher
}

// Match checks if the given Path matches this pattern or is a descendant of it (subtree match).
func (p Pattern) Match(path diff.Path) bool {
	pathSegs := path.Segments()
	if len(pathSegs) < len(p.Segments) || len(p.Segments) == 0 {
		return false
	}

	for i, sm := range p.Segments {
		if !sm.Match(pathSegs[i]) {
			return false
		}
	}
	return true
}

// ParsePattern parses an ignore rule string into a compiled Pattern.
// Supported syntax:
//   - "timestamp" -> [ExactKey("timestamp")]
//   - "*.timestamp" -> [WildcardKey, ExactKey("timestamp")]
//   - "users[0].name" -> [ExactKey("users"), ExactIndex(0), ExactKey("name")]
//   - "users[*].id" -> [ExactKey("users"), WildcardIndex, ExactKey("id")]
//   - "[*].id" -> [WildcardIndex, ExactKey("id")]
func ParsePattern(raw string) (Pattern, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Pattern{}, fmt.Errorf("empty ignore pattern")
	}

	var matchers []SegmentMatcher
	tokens := splitPathTokens(raw)

	for _, token := range tokens {
		if token == "" {
			continue
		}

		if strings.HasPrefix(token, "[") && strings.HasSuffix(token, "]") {
			inner := token[1 : len(token)-1]
			if inner == "*" {
				matchers = append(matchers, SegmentMatcher{Kind: MatcherWildcardIndex})
			} else {
				idx, err := strconv.Atoi(inner)
				if err != nil || idx < 0 {
					return Pattern{}, fmt.Errorf("invalid array index in pattern %q: %s", raw, token)
				}
				matchers = append(matchers, SegmentMatcher{Kind: MatcherExactIndex, Index: idx})
			}
		} else if token == "*" {
			matchers = append(matchers, SegmentMatcher{Kind: MatcherWildcardKey})
		} else {
			if strings.Contains(token, "*") {
				return Pattern{}, fmt.Errorf("partial wildcards are not supported in pattern %q: %s", raw, token)
			}
			matchers = append(matchers, SegmentMatcher{Kind: MatcherExactKey, Key: token})
		}
	}

	if len(matchers) == 0 {
		return Pattern{}, fmt.Errorf("invalid ignore pattern %q", raw)
	}

	return Pattern{
		Raw:      raw,
		Segments: matchers,
	}, nil
}

// splitPathTokens splits a path string into dot and bracket tokens.
// E.g. "users[0].name" -> ["users", "[0]", "name"]
// E.g. "[*].id" -> ["[foo]" -> "[*]", "id"]
func splitPathTokens(s string) []string {
	var tokens []string
	var current strings.Builder

	flush := func() {
		if current.Len() > 0 {
			tokens = append(tokens, current.String())
			current.Reset()
		}
	}

	inBracket := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch ch {
		case '.':
			if inBracket {
				current.WriteByte(ch)
			} else {
				flush()
			}
		case '[':
			flush()
			inBracket = true
			current.WriteByte(ch)
		case ']':
			current.WriteByte(ch)
			inBracket = false
			flush()
		default:
			current.WriteByte(ch)
		}
	}
	flush()

	return tokens
}

// Matcher aggregates multiple ignore patterns and matches paths against them.
type Matcher struct {
	Patterns []Pattern
}

// New creates a Matcher from a list of raw rule strings.
func New(rules []string) (*Matcher, error) {
	var patterns []Pattern
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		p, err := ParsePattern(rule)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, p)
	}
	return &Matcher{Patterns: patterns}, nil
}

// Matches returns true if any configured pattern matches the given path.
func (m *Matcher) Matches(path diff.Path) bool {
	if m == nil || len(m.Patterns) == 0 {
		return false
	}
	for _, p := range m.Patterns {
		if p.Match(path) {
			return true
		}
	}
	return false
}

// Rules returns the raw pattern strings configured in this matcher.
func (m *Matcher) Rules() []string {
	if m == nil {
		return nil
	}
	res := make([]string, len(m.Patterns))
	for i, p := range m.Patterns {
		res[i] = p.Raw
	}
	return res
}
