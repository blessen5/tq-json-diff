package diff

import "strings"

// Path represents a structured, immutable JSON path composed of key segments.
type Path struct {
	segments []string
}

// NewPath creates a new Path from zero or more string segments.
func NewPath(segments ...string) Path {
	return Path{segments: segments}
}

// Append returns a new Path with the given segment appended to the end.
func (p Path) Append(segment string) Path {
	newSegs := make([]string, len(p.segments)+1)
	copy(newSegs, p.segments)
	newSegs[len(p.segments)] = segment
	return Path{segments: newSegs}
}

// String returns the dot-separated representation of the path (e.g. "user.profile.email").
// If the path represents the root document (0 segments), it returns "(root)".
func (p Path) String() string {
	if len(p.segments) == 0 {
		return "(root)"
	}
	return strings.Join(p.segments, ".")
}

// IsRoot returns true if this path represents the document root.
func (p Path) IsRoot() bool {
	return len(p.segments) == 0
}

// Segments returns a copy of the path's constituent segments.
func (p Path) Segments() []string {
	res := make([]string, len(p.segments))
	copy(res, p.segments)
	return res
}
