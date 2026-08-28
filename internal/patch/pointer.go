package patch

import (
	"fmt"
	"strconv"
	"strings"

	"jdiff/internal/diff"
)

// Escape encodes special characters in a JSON Pointer token according to RFC 6901.
// '~' is encoded as '~0' and '/' is encoded as '~1'.
func Escape(s string) string {
	s = strings.ReplaceAll(s, "~", "~0")
	s = strings.ReplaceAll(s, "/", "~1")
	return s
}

// Unescape decodes RFC 6901 escaped characters in a JSON Pointer token.
// '~1' is decoded to '/' and '~0' is decoded to '~'.
func Unescape(s string) (string, error) {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '~' {
			if i+1 >= len(s) {
				return "", fmt.Errorf("invalid escape sequence at end of pointer token: %q", s)
			}
			switch s[i+1] {
			case '0':
				sb.WriteByte('~')
				i++
			case '1':
				sb.WriteByte('/')
				i++
			default:
				return "", fmt.Errorf("invalid escape sequence ~%c in pointer token: %q", s[i+1], s)
			}
		} else {
			sb.WriteByte(s[i])
		}
	}
	return sb.String(), nil
}

// FromPath converts an internal diff.Path to an RFC 6901 JSON Pointer string.
func FromPath(p diff.Path) string {
	segs := p.Segments()
	if len(segs) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, seg := range segs {
		sb.WriteString("/")
		switch seg.Kind {
		case diff.SegmentKindKey:
			sb.WriteString(Escape(seg.Key))
		case diff.SegmentKindIndex:
			sb.WriteString(strconv.Itoa(seg.Index))
		}
	}
	return sb.String()
}

// ParsePointer splits an RFC 6901 JSON Pointer into unescaped constituent tokens.
func ParsePointer(ptr string) ([]string, error) {
	if ptr == "" {
		return nil, nil
	}

	if !strings.HasPrefix(ptr, "/") {
		return nil, fmt.Errorf("invalid JSON Pointer %q: must start with '/'", ptr)
	}

	rawTokens := strings.Split(ptr[1:], "/")
	tokens := make([]string, len(rawTokens))
	for i, raw := range rawTokens {
		unescaped, err := Unescape(raw)
		if err != nil {
			return nil, err
		}
		tokens[i] = unescaped
	}
	return tokens, nil
}
