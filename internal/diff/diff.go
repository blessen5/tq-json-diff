package diff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

// ChangeType represents the type of difference detected.
type ChangeType string

const (
	// ChangeAdded indicates a value was added in the new document.
	ChangeAdded ChangeType = "ADDED"
	// ChangeRemoved indicates a value was removed from the old document.
	ChangeRemoved ChangeType = "REMOVED"
	// ChangeModified indicates a value was modified between old and new.
	ChangeModified ChangeType = "MODIFIED"
)

// Change represents a single structural delta between two JSON documents.
type Change struct {
	Path     string
	Type     ChangeType
	OldValue any
	NewValue any
}

// DiffResult holds all changes detected between two JSON documents.
type DiffResult struct {
	Changes []Change
}

// HasChanges returns true if any differences were detected.
func (r *DiffResult) HasChanges() bool {
	return r != nil && len(r.Changes) > 0
}

// Added returns all added changes.
func (r *DiffResult) Added() []Change {
	var list []Change
	for _, c := range r.Changes {
		if c.Type == ChangeAdded {
			list = append(list, c)
		}
	}
	return list
}

// Removed returns all removed changes.
func (r *DiffResult) Removed() []Change {
	var list []Change
	for _, c := range r.Changes {
		if c.Type == ChangeRemoved {
			list = append(list, c)
		}
	}
	return list
}

// Modified returns all modified changes.
func (r *DiffResult) Modified() []Change {
	var list []Change
	for _, c := range r.Changes {
		if c.Type == ChangeModified {
			list = append(list, c)
		}
	}
	return list
}

// Format writes the formatted diff output to the given writer.
func (r *DiffResult) Format(w io.Writer) error {
	if !r.HasChanges() {
		return nil
	}

	var sections []string

	modified := r.Modified()
	if len(modified) > 0 {
		var sb strings.Builder
		sb.WriteString("MODIFIED:\n")
		for i, c := range modified {
			sb.WriteString(fmt.Sprintf("    %s\n", c.Path))
			sb.WriteString(fmt.Sprintf("        - %s\n", FormatValue(c.OldValue)))
			sb.WriteString(fmt.Sprintf("        + %s", FormatValue(c.NewValue)))
			if i < len(modified)-1 {
				sb.WriteString("\n\n")
			}
		}
		sections = append(sections, sb.String())
	}

	added := r.Added()
	if len(added) > 0 {
		var sb strings.Builder
		sb.WriteString("ADDED:\n")
		for i, c := range added {
			sb.WriteString(fmt.Sprintf("    %s\n", c.Path))
			sb.WriteString(fmt.Sprintf("        + %s", FormatValue(c.NewValue)))
			if i < len(added)-1 {
				sb.WriteString("\n\n")
			}
		}
		sections = append(sections, sb.String())
	}

	removed := r.Removed()
	if len(removed) > 0 {
		var sb strings.Builder
		sb.WriteString("REMOVED:\n")
		for i, c := range removed {
			sb.WriteString(fmt.Sprintf("    %s\n", c.Path))
			sb.WriteString(fmt.Sprintf("        - %s", FormatValue(c.OldValue)))
			if i < len(removed)-1 {
				sb.WriteString("\n\n")
			}
		}
		sections = append(sections, sb.String())
	}

	output := strings.Join(sections, "\n\n") + "\n"
	_, err := fmt.Fprint(w, output)
	return err
}

// String returns the formatted diff as a string.
func (r *DiffResult) String() string {
	var buf bytes.Buffer
	_ = r.Format(&buf)
	return buf.String()
}

// CompareBytes parses two JSON byte slices and computes their structural differences.
func CompareBytes(oldJSON, newJSON []byte) (*DiffResult, error) {
	var oldVal any
	var newVal any

	decOld := json.NewDecoder(bytes.NewReader(oldJSON))
	decOld.UseNumber()
	if err := decOld.Decode(&oldVal); err != nil {
		return nil, fmt.Errorf("old document: %w", err)
	}

	decNew := json.NewDecoder(bytes.NewReader(newJSON))
	decNew.UseNumber()
	if err := decNew.Decode(&newVal); err != nil {
		return nil, fmt.Errorf("new document: %w", err)
	}

	changes := compare("", oldVal, newVal)
	return &DiffResult{Changes: changes}, nil
}

// compare recursively compares old and new JSON values and returns detected changes.
func compare(path string, oldVal, newVal any) []Change {
	var changes []Change

	// Case 1: Both are JSON objects (maps)
	oldMap, oldIsMap := oldVal.(map[string]any)
	newMap, newIsMap := newVal.(map[string]any)

	if oldIsMap && newIsMap {
		// Collect all unique keys in deterministic sorted order
		keySet := make(map[string]struct{})
		for k := range oldMap {
			keySet[k] = struct{}{}
		}
		for k := range newMap {
			keySet[k] = struct{}{}
		}

		keys := make([]string, 0, len(keySet))
		for k := range keySet {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			childPath := joinPath(path, k)
			oldChild, inOld := oldMap[k]
			newChild, inNew := newMap[k]

			if inOld && !inNew {
				changes = append(changes, Change{
					Path:     childPath,
					Type:     ChangeRemoved,
					OldValue: oldChild,
					NewValue: nil,
				})
			} else if !inOld && inNew {
				changes = append(changes, Change{
					Path:     childPath,
					Type:     ChangeAdded,
					OldValue: nil,
					NewValue: newChild,
				})
			} else {
				// Present in both - recurse
				changes = append(changes, compare(childPath, oldChild, newChild)...)
			}
		}
		return changes
	}

	// Case 2: One is an object and the other is not -> modification
	if oldIsMap != newIsMap {
		changes = append(changes, Change{
			Path:     path,
			Type:     ChangeModified,
			OldValue: oldVal,
			NewValue: newVal,
		})
		return changes
	}

	// Case 3: Both are slices/arrays (v0.2.0: safe whole-value equality)
	oldSlice, oldIsSlice := oldVal.([]any)
	newSlice, newIsSlice := newVal.([]any)

	if oldIsSlice || newIsSlice {
		if !reflect.DeepEqual(oldVal, newVal) {
			changes = append(changes, Change{
				Path:     path,
				Type:     ChangeModified,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
		return changes
	}

	// Case 4: Primitive values (string, number, bool, nil)
	if !valuesEqual(oldVal, newVal) {
		changes = append(changes, Change{
			Path:     path,
			Type:     ChangeModified,
			OldValue: oldVal,
			NewValue: newVal,
		})
	}

	return changes
}

// valuesEqual checks equality between two JSON primitive values.
func valuesEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare json.Number or other primitives
	numA, isNumA := a.(json.Number)
	numB, isNumB := b.(json.Number)
	if isNumA && isNumB {
		return numA.String() == numB.String()
	}

	return reflect.DeepEqual(a, b)
}

// joinPath creates a dot-separated JSON path.
func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

// FormatValue formats a Go value into a readable JSON representation for terminal display.
func FormatValue(v any) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case json.Number:
		return val.String()
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(b)
	}
}
