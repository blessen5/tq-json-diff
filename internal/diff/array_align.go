package diff

import (
	"fmt"
)

// ArrayMatchStrategy defines how array elements are aligned during comparison.
type ArrayMatchStrategy string

const (
	// MatchIndex compares array elements by ordinal index (default).
	MatchIndex ArrayMatchStrategy = "index"
	// MatchAuto automatically detects identity keys in object arrays.
	MatchAuto ArrayMatchStrategy = "auto"
	// MatchKey uses an explicitly provided key field.
	MatchKey ArrayMatchStrategy = "key"
)

// Standard identity key candidate list for automatic detection.
var autoKeyCandidates = []string{"id", "_id", "uuid", "key", "name", "slug", "code", "email", "username"}

// DetectArrayKey scans array objects to detect a common unique identifier key.
func DetectArrayKey(slice []any) string {
	if len(slice) == 0 {
		return ""
	}

	for _, cand := range autoKeyCandidates {
		matchedAll := true
		for _, item := range slice {
			m, ok := item.(map[string]any)
			if !ok {
				matchedAll = false
				break
			}
			if val, exists := m[cand]; !exists || val == nil {
				matchedAll = false
				break
			}
		}
		if matchedAll {
			return cand
		}
	}

	return ""
}

type matchedPair struct {
	oldIdx  int
	newIdx  int
	oldItem any
	newItem any
}

type arrayAlignmentResult struct {
	matchedPairs []matchedPair
	removedOld   []int
	addedNew     []int
	keyField     string
}

// alignArrayObjects matches array objects using the given key field.
func alignArrayObjects(oldSlice, newSlice []any, keyField string) arrayAlignmentResult {
	res := arrayAlignmentResult{
		keyField: keyField,
	}

	oldKeyMap := make(map[string]int)
	for i, item := range oldSlice {
		if m, ok := item.(map[string]any); ok {
			if val, exists := m[keyField]; exists && val != nil {
				oldKeyMap[fmt.Sprintf("%v", val)] = i
			}
		}
	}

	matchedOldIndices := make(map[int]bool)
	matchedNewIndices := make(map[int]bool)

	for j, item := range newSlice {
		if m, ok := item.(map[string]any); ok {
			if val, exists := m[keyField]; exists && val != nil {
				keyStr := fmt.Sprintf("%v", val)
				if oldIdx, found := oldKeyMap[keyStr]; found {
					if !matchedOldIndices[oldIdx] {
						matchedOldIndices[oldIdx] = true
						matchedNewIndices[j] = true
						res.matchedPairs = append(res.matchedPairs, matchedPair{
							oldIdx:  oldIdx,
							newIdx:  j,
							oldItem: oldSlice[oldIdx],
							newItem: newSlice[j],
						})
					}
				}
			}
		}
	}

	for i := range oldSlice {
		if !matchedOldIndices[i] {
			res.removedOld = append(res.removedOld, i)
		}
	}

	for j := range newSlice {
		if !matchedNewIndices[j] {
			res.addedNew = append(res.addedNew, j)
		}
	}

	return res
}
