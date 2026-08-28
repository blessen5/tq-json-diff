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
	// ChangeAdded indicates a property was added in the new document.
	ChangeAdded ChangeType = "ADDED"
	// ChangeRemoved indicates a property was removed from the old document.
	ChangeRemoved ChangeType = "REMOVED"
	// ChangeModified indicates a value was modified between old and new documents.
	ChangeModified ChangeType = "MODIFIED"
)

// JSONType represents the high-level JSON data type of a value.
type JSONType string

const (
	// JSONTypeString represents a JSON string.
	JSONTypeString JSONType = "string"
	// JSONTypeNumber represents a JSON number.
	JSONTypeNumber JSONType = "number"
	// JSONTypeBoolean represents a JSON boolean.
	JSONTypeBoolean JSONType = "boolean"
	// JSONTypeNull represents a JSON null literal.
	JSONTypeNull JSONType = "null"
	// JSONTypeObject represents a JSON object structure.
	JSONTypeObject JSONType = "object"
	// JSONTypeArray represents a JSON array structure.
	JSONTypeArray JSONType = "array"
)

// DetectType determines the JSONType of a decoded Go value.
func DetectType(v any) JSONType {
	if v == nil {
		return JSONTypeNull
	}
	switch v.(type) {
	case string:
		return JSONTypeString
	case json.Number, float64, int, int64, float32:
		return JSONTypeNumber
	case bool:
		return JSONTypeBoolean
	case map[string]any:
		return JSONTypeObject
	case []any:
		return JSONTypeArray
	default:
		return JSONTypeObject
	}
}

// Change represents a single structural delta between two JSON documents.
type Change struct {
	Path     Path
	Type     ChangeType
	OldValue any
	NewValue any
	OldType  JSONType
	NewType  JSONType
}

// Summary encapsulates the total counts of detected changes.
type Summary struct {
	Added    int
	Removed  int
	Modified int
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

// Summary calculates the summary statistics of the diff result.
func (r *DiffResult) Summary() Summary {
	if r == nil {
		return Summary{}
	}
	s := Summary{}
	for _, c := range r.Changes {
		switch c.Type {
		case ChangeAdded:
			s.Added++
		case ChangeRemoved:
			s.Removed++
		case ChangeModified:
			s.Modified++
		}
	}
	return s
}

// Format writes the formatted diff output and summary to the given writer.
func (r *DiffResult) Format(w io.Writer) error {
	if !r.HasChanges() {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	var sections []string

	modified := r.Modified()
	if len(modified) > 0 {
		var sb strings.Builder
		sb.WriteString("MODIFIED\n")
		for i, c := range modified {
			sb.WriteString(fmt.Sprintf("  %s\n", c.Path.String()))
			sb.WriteString(fmt.Sprintf("    - %s\n", FormatValue(c.OldValue)))
			sb.WriteString(fmt.Sprintf("    + %s", FormatValue(c.NewValue)))
			if i < len(modified)-1 {
				sb.WriteString("\n\n")
			}
		}
		sections = append(sections, sb.String())
	}

	added := r.Added()
	if len(added) > 0 {
		var sb strings.Builder
		sb.WriteString("ADDED\n")
		for i, c := range added {
			sb.WriteString(fmt.Sprintf("  %s\n", c.Path.String()))
			sb.WriteString(fmt.Sprintf("    + %s", FormatValue(c.NewValue)))
			if i < len(added)-1 {
				sb.WriteString("\n\n")
			}
		}
		sections = append(sections, sb.String())
	}

	removed := r.Removed()
	if len(removed) > 0 {
		var sb strings.Builder
		sb.WriteString("REMOVED\n")
		for i, c := range removed {
			sb.WriteString(fmt.Sprintf("  %s\n", c.Path.String()))
			sb.WriteString(fmt.Sprintf("    - %s", FormatValue(c.OldValue)))
			if i < len(removed)-1 {
				sb.WriteString("\n\n")
			}
		}
		sections = append(sections, sb.String())
	}

	// Change Summary block
	summary := r.Summary()
	var sumSB strings.Builder
	sumSB.WriteString("Summary:\n")
	sumSB.WriteString(fmt.Sprintf("  Added:     %d\n", summary.Added))
	sumSB.WriteString(fmt.Sprintf("  Removed:   %d\n", summary.Removed))
	sumSB.WriteString(fmt.Sprintf("  Modified:  %d", summary.Modified))
	sections = append(sections, sumSB.String())

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

	changes := compare(NewPath(), oldVal, newVal)
	sortChanges(changes)
	return &DiffResult{Changes: changes}, nil
}

// compare recursively compares old and new JSON values and returns detected changes.
func compare(path Path, oldVal, newVal any) []Change {
	var changes []Change

	oldType := DetectType(oldVal)
	newType := DetectType(newVal)

	// Case 1: Both are JSON objects (maps)
	if oldType == JSONTypeObject && newType == JSONTypeObject {
		oldMap := oldVal.(map[string]any)
		newMap := newVal.(map[string]any)

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
			childPath := path.Append(k)
			oldChild, inOld := oldMap[k]
			newChild, inNew := newMap[k]

			if inOld && !inNew {
				changes = append(changes, Change{
					Path:     childPath,
					Type:     ChangeRemoved,
					OldValue: oldChild,
					NewValue: nil,
					OldType:  DetectType(oldChild),
					NewType:  JSONTypeNull,
				})
			} else if !inOld && inNew {
				changes = append(changes, Change{
					Path:     childPath,
					Type:     ChangeAdded,
					OldValue: nil,
					NewValue: newChild,
					OldType:  JSONTypeNull,
					NewType:  DetectType(newChild),
				})
			} else {
				// Present in both - recurse deeply
				changes = append(changes, compare(childPath, oldChild, newChild)...)
			}
		}
		return changes
	}

	// Case 2: Type mismatch (e.g. object vs string, number vs string, null vs value)
	if oldType != newType {
		changes = append(changes, Change{
			Path:     path,
			Type:     ChangeModified,
			OldValue: oldVal,
			NewValue: newVal,
			OldType:  oldType,
			NewType:  newType,
		})
		return changes
	}

	// Case 3: Both are slices/arrays (v0.3.0: safe whole-value equality)
	if oldType == JSONTypeArray {
		if !reflect.DeepEqual(oldVal, newVal) {
			changes = append(changes, Change{
				Path:     path,
				Type:     ChangeModified,
				OldValue: oldVal,
				NewValue: newVal,
				OldType:  oldType,
				NewType:  newType,
			})
		}
		return changes
	}

	// Case 4: Primitive values with identical types (string, number, bool, null)
	if !valuesEqual(oldVal, newVal) {
		changes = append(changes, Change{
			Path:     path,
			Type:     ChangeModified,
			OldValue: oldVal,
			NewValue: newVal,
			OldType:  oldType,
			NewType:  newType,
		})
	}

	return changes
}

// sortChanges sorts a slice of changes deterministically by path string and change type.
func sortChanges(changes []Change) {
	sort.SliceStable(changes, func(i, j int) bool {
		if changes[i].Path.String() != changes[j].Path.String() {
			return changes[i].Path.String() < changes[j].Path.String()
		}
		return changes[i].Type < changes[j].Type
	})
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
