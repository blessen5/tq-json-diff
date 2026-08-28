package patch

import (
	"encoding/json"
	"sort"

	"jdiff/internal/diff"
)

// OpType represents the RFC 6902 JSON Patch operation verb.
type OpType string

const (
	// OpAdd adds a value to an object or array.
	OpAdd OpType = "add"
	// OpRemove removes a value from an object or array.
	OpRemove OpType = "remove"
	// OpReplace replaces an existing value in an object or array.
	OpReplace OpType = "replace"
)

// Operation represents a single RFC 6902 JSON Patch operation.
type Operation struct {
	Op    OpType `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
}

// MarshalJSON marshals an Operation ensuring remove operations omit the value field entirely.
func (op Operation) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"op":   string(op.Op),
		"path": op.Path,
	}
	if op.Op != OpRemove {
		m["value"] = op.Value
	}
	return json.Marshal(m)
}

// Patch represents an RFC 6902 JSON Patch document (an array of Operations).
type Patch []Operation

// Generate creates a Patch document from a DiffResult, applying optimal ordering for sequential patch execution.
func Generate(result *diff.DiffResult) Patch {
	if result == nil || len(result.Changes) == 0 {
		return Patch{}
	}

	// 1. Separate changes into categories
	var replacements []diff.Change
	var removals []diff.Change
	var additions []diff.Change

	for _, c := range result.Changes {
		switch c.Type {
		case diff.ChangeModified:
			replacements = append(replacements, c)
		case diff.ChangeRemoved:
			removals = append(removals, c)
		case diff.ChangeAdded:
			additions = append(additions, c)
		}
	}

	// 2. Sort removals in descending order for array indices so index shifting does not break sequential execution
	sort.SliceStable(removals, func(i, j int) bool {
		// If both are array items under the same parent prefix, descending index first
		p1 := removals[i].Path.Segments()
		p2 := removals[j].Path.Segments()
		if len(p1) > 0 && len(p2) > 0 && p1[len(p1)-1].Kind == diff.SegmentKindIndex && p2[len(p2)-1].Kind == diff.SegmentKindIndex {
			return p1[len(p1)-1].Index > p2[len(p2)-1].Index
		}
		return !removals[i].Path.Less(removals[j].Path)
	})

	// 3. Sort additions in ascending order
	sort.SliceStable(additions, func(i, j int) bool {
		return additions[i].Path.Less(additions[j].Path)
	})

	// 4. Sort replacements deterministically
	sort.SliceStable(replacements, func(i, j int) bool {
		return replacements[i].Path.Less(replacements[j].Path)
	})

	var ops []Operation

	// Execute replacements first
	for _, c := range replacements {
		ops = append(ops, Operation{
			Op:    OpReplace,
			Path:  FromPath(c.Path),
			Value: c.NewValue,
		})
	}

	// Execute removals next (tail-to-head for arrays)
	for _, c := range removals {
		ops = append(ops, Operation{
			Op:   OpRemove,
			Path: FromPath(c.Path),
		})
	}

	// Execute additions last
	for _, c := range additions {
		ops = append(ops, Operation{
			Op:    OpAdd,
			Path:  FromPath(c.Path),
			Value: c.NewValue,
		})
	}

	return ops
}
