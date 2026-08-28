package diff

import (
	"fmt"
	"strings"
)

// SegmentKind indicates whether a path segment is an object key or an array index.
type SegmentKind int

const (
	// SegmentKindKey represents an object property key.
	SegmentKindKey SegmentKind = iota
	// SegmentKindIndex represents an array element index.
	SegmentKindIndex
)

// Segment represents a single step in a JSON path.
type Segment struct {
	Kind  SegmentKind
	Key   string
	Index int
}

// Path represents a structured, immutable JSON path composed of key and index segments.
type Path struct {
	segments []Segment
}

// NewPath creates a new empty Path (representing the root document).
func NewPath() Path {
	return Path{segments: nil}
}

// AppendKey returns a new Path with the object property key appended.
func (p Path) AppendKey(key string) Path {
	newSegs := make([]Segment, len(p.segments)+1)
	copy(newSegs, p.segments)
	newSegs[len(p.segments)] = Segment{
		Kind: SegmentKindKey,
		Key:  key,
	}
	return Path{segments: newSegs}
}

// Append is a convenience alias for AppendKey to maintain compatibility.
func (p Path) Append(key string) Path {
	return p.AppendKey(key)
}

// AppendIndex returns a new Path with the array index appended.
func (p Path) AppendIndex(index int) Path {
	newSegs := make([]Segment, len(p.segments)+1)
	copy(newSegs, p.segments)
	newSegs[len(p.segments)] = Segment{
		Kind:  SegmentKindIndex,
		Index: index,
	}
	return Path{segments: newSegs}
}

// String returns the formatted JSON path representation (e.g. "users[0].name", "[2]", "data.groups[0].values[1]").
// If the path represents the root document (0 segments), it returns "(root)".
func (p Path) String() string {
	if len(p.segments) == 0 {
		return "(root)"
	}

	var sb strings.Builder
	for i, seg := range p.segments {
		switch seg.Kind {
		case SegmentKindKey:
			if i > 0 {
				sb.WriteString(".")
			}
			sb.WriteString(seg.Key)
		case SegmentKindIndex:
			sb.WriteString(fmt.Sprintf("[%d]", seg.Index))
		}
	}
	return sb.String()
}

// IsRoot returns true if this path represents the document root.
func (p Path) IsRoot() bool {
	return len(p.segments) == 0
}

// Segments returns a copy of the path's constituent segments.
func (p Path) Segments() []Segment {
	res := make([]Segment, len(p.segments))
	copy(res, p.segments)
	return res
}

// Less compares two paths deterministically with natural numeric index ordering.
func (p Path) Less(other Path) bool {
	minLen := len(p.segments)
	if len(other.segments) < minLen {
		minLen = len(other.segments)
	}

	for i := 0; i < minLen; i++ {
		s1 := p.segments[i]
		s2 := other.segments[i]

		if s1.Kind != s2.Kind {
			return s1.Kind < s2.Kind
		}

		if s1.Kind == SegmentKindKey {
			if s1.Key != s2.Key {
				return s1.Key < s2.Key
			}
		} else {
			if s1.Index != s2.Index {
				return s1.Index < s2.Index
			}
		}
	}

	return len(p.segments) < len(other.segments)
}
