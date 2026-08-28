package patch

import (
	"sort"

	"jdiff/internal/diff"
)

// GenerateRollback generates an inverse RFC 6902 JSON Patch that mutates the new document back to the old document.
func GenerateRollback(result *diff.DiffResult) Patch {
	if result == nil || len(result.Changes) == 0 {
		return Patch{}
	}

	var replacements []diff.Change
	var additions []diff.Change // Inverted removals
	var removals []diff.Change  // Inverted additions

	for _, c := range result.Changes {
		switch c.Type {
		case diff.ChangeModified:
			replacements = append(replacements, c)
		case diff.ChangeRemoved:
			// In old but not in new -> rollback must ADD it back
			additions = append(additions, c)
		case diff.ChangeAdded:
			// In new but not in old -> rollback must REMOVE it
			removals = append(removals, c)
		}
	}

	// Removals (on array indices) should be sorted in descending order
	sort.SliceStable(removals, func(i, j int) bool {
		p1 := removals[i].Path.Segments()
		p2 := removals[j].Path.Segments()
		if len(p1) > 0 && len(p2) > 0 && p1[len(p1)-1].Kind == diff.SegmentKindIndex && p2[len(p2)-1].Kind == diff.SegmentKindIndex {
			return p1[len(p1)-1].Index > p2[len(p2)-1].Index
		}
		return !removals[i].Path.Less(removals[j].Path)
	})

	// Additions sorted ascending
	sort.SliceStable(additions, func(i, j int) bool {
		return additions[i].Path.Less(additions[j].Path)
	})

	// Replacements sorted deterministically
	sort.SliceStable(replacements, func(i, j int) bool {
		return replacements[i].Path.Less(replacements[j].Path)
	})

	var ops []Operation

	// 1. Replacements restore old value
	for _, c := range replacements {
		ops = append(ops, Operation{
			Op:    OpReplace,
			Path:  FromPath(c.Path),
			Value: c.OldValue,
		})
	}

	// 2. Additions restore removed old value
	for _, c := range additions {
		ops = append(ops, Operation{
			Op:    OpAdd,
			Path:  FromPath(c.Path),
			Value: c.OldValue,
		})
	}

	// 3. Removals strip added new value
	for _, c := range removals {
		ops = append(ops, Operation{
			Op:   OpRemove,
			Path: FromPath(c.Path),
		})
	}

	return ops
}
