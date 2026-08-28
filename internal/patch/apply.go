package patch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"jdiff/internal/diff"
)

// Apply applies an RFC 6902 JSON Patch document to an in-memory JSON data structure.
// Returns the resulting modified JSON data structure or an error if any operation fails.
func Apply(doc any, patch Patch) (any, error) {
	cloned, err := cloneJSON(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to clone input document: %w", err)
	}

	root := cloned

	for opIdx, op := range patch {
		tokens, err := ParsePointer(op.Path)
		if err != nil {
			return nil, fmt.Errorf("operation [%d]: %w", opIdx, err)
		}

		if len(tokens) == 0 {
			switch op.Op {
			case OpAdd, OpReplace:
				root = op.Value
				continue
			case OpRemove:
				root = nil
				continue
			default:
				return nil, fmt.Errorf("operation [%d]: unsupported op %q", opIdx, op.Op)
			}
		}

		newRoot, err := applyAtPath(root, tokens, op)
		if err != nil {
			return nil, fmt.Errorf("operation [%d] (%s %s): %w", opIdx, op.Op, op.Path, err)
		}
		root = newRoot
	}

	return root, nil
}

func applyAtPath(root any, tokens []string, op Operation) (any, error) {
	if len(tokens) == 1 {
		return applyLeaf(root, tokens[0], op)
	}

	token := tokens[0]
	rest := tokens[1:]

	switch container := root.(type) {
	case map[string]any:
		child, exists := container[token]
		if !exists {
			return nil, fmt.Errorf("path segment %q not found in object", token)
		}
		updatedChild, err := applyAtPath(child, rest, op)
		if err != nil {
			return nil, err
		}
		container[token] = updatedChild
		return container, nil

	case []any:
		idx, err := strconv.Atoi(token)
		if err != nil || idx < 0 || idx >= len(container) {
			return nil, fmt.Errorf("invalid array index %q (array length %d)", token, len(container))
		}
		updatedChild, err := applyAtPath(container[idx], rest, op)
		if err != nil {
			return nil, err
		}
		container[idx] = updatedChild
		return container, nil

	default:
		return nil, fmt.Errorf("cannot traverse into primitive value of type %T at token %q", root, token)
	}
}

func applyLeaf(container any, key string, op Operation) (any, error) {
	switch c := container.(type) {
	case map[string]any:
		switch op.Op {
		case OpAdd, OpReplace:
			if op.Op == OpReplace {
				if _, exists := c[key]; !exists {
					return nil, fmt.Errorf("cannot replace non-existent object key %q", key)
				}
			}
			c[key] = op.Value
			return c, nil

		case OpRemove:
			if _, exists := c[key]; !exists {
				return nil, fmt.Errorf("cannot remove non-existent object key %q", key)
			}
			delete(c, key)
			return c, nil

		default:
			return nil, fmt.Errorf("unsupported op %q", op.Op)
		}

	case []any:
		switch op.Op {
		case OpAdd:
			if key == "-" {
				return append(c, op.Value), nil
			}
			idx, err := strconv.Atoi(key)
			if err != nil || idx < 0 || idx > len(c) {
				return nil, fmt.Errorf("index out of range for array add: %q (array length %d)", key, len(c))
			}
			res := make([]any, 0, len(c)+1)
			res = append(res, c[:idx]...)
			res = append(res, op.Value)
			res = append(res, c[idx:]...)
			return res, nil

		case OpReplace:
			idx, err := strconv.Atoi(key)
			if err != nil || idx < 0 || idx >= len(c) {
				return nil, fmt.Errorf("index out of range for array replace: %q (array length %d)", key, len(c))
			}
			c[idx] = op.Value
			return c, nil

		case OpRemove:
			idx, err := strconv.Atoi(key)
			if err != nil || idx < 0 || idx >= len(c) {
				return nil, fmt.Errorf("index out of range for array remove: %q (array length %d)", key, len(c))
			}
			res := make([]any, 0, len(c)-1)
			res = append(res, c[:idx]...)
			res = append(res, c[idx+1:]...)
			return res, nil

		default:
			return nil, fmt.Errorf("unsupported op %q", op.Op)
		}

	default:
		return nil, fmt.Errorf("cannot apply operation %q on non-container %T", op.Op, container)
	}
}

// Verify applies the patch to oldJSON and verifies whether the patched document matches newJSON.
func Verify(oldJSON, newJSON []byte, patch Patch) (bool, error) {
	var oldVal any
	decOld := json.NewDecoder(bytes.NewReader(oldJSON))
	decOld.UseNumber()
	if err := decOld.Decode(&oldVal); err != nil {
		return false, fmt.Errorf("failed to decode old document: %w", err)
	}

	patchedVal, err := Apply(oldVal, patch)
	if err != nil {
		return false, fmt.Errorf("patch application failed: %w", err)
	}

	patchedBytes, err := json.Marshal(patchedVal)
	if err != nil {
		return false, fmt.Errorf("failed to marshal patched value: %w", err)
	}

	diffResult, err := diff.CompareBytes(patchedBytes, newJSON)
	if err != nil {
		return false, fmt.Errorf("verification comparison error: %w", err)
	}

	return !diffResult.HasChanges(), nil
}

func cloneJSON(v any) (any, error) {
	if v == nil {
		return nil, nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var res any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&res); err != nil {
		return nil, err
	}
	return res, nil
}
