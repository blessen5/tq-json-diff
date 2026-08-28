package diff

import (
	"encoding/json"
	"fmt"
	"testing"
)

func BenchmarkSmallJSON(b *testing.B) {
	oldJSON := []byte(`{"name": "Alice", "age": 30, "active": true, "city": "London"}`)
	newJSON := []byte(`{"name": "Alice", "age": 31, "active": false, "city": "Paris"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMediumJSON(b *testing.B) {
	oldDoc := map[string]any{
		"server": map[string]any{
			"host": "localhost",
			"port": 8080,
			"tls":  true,
		},
		"database": map[string]any{
			"primary": map[string]any{
				"host": "db.prod.internal",
				"port": 5432,
				"pool": 20,
			},
		},
		"features": []any{"auth", "metrics", "logging", "tracing"},
	}
	newDoc := map[string]any{
		"server": map[string]any{
			"host": "0.0.0.0",
			"port": 8080,
			"tls":  true,
		},
		"database": map[string]any{
			"primary": map[string]any{
				"host": "db.prod.internal",
				"port": 5432,
				"pool": 50,
			},
		},
		"features": []any{"auth", "metrics", "profiling"},
	}

	oldJSON, _ := json.Marshal(oldDoc)
	newJSON, _ := json.Marshal(newDoc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeeplyNestedJSON(b *testing.B) {
	buildNested := func(depth int, leafVal string) []byte {
		m := map[string]any{"leaf": leafVal}
		for i := 0; i < depth; i++ {
			m = map[string]any{fmt.Sprintf("level_%d", i): m}
		}
		data, _ := json.Marshal(m)
		return data
	}

	oldJSON := buildNested(20, "old")
	newJSON := buildNested(20, "new")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLargeArray(b *testing.B) {
	oldArr := make([]any, 1000)
	newArr := make([]any, 1000)
	for i := 0; i < 1000; i++ {
		oldArr[i] = fmt.Sprintf("item_%d", i)
		if i%2 == 0 {
			newArr[i] = fmt.Sprintf("modified_%d", i)
		} else {
			newArr[i] = fmt.Sprintf("item_%d", i)
		}
	}

	oldJSON, _ := json.Marshal(map[string]any{"items": oldArr})
	newJSON, _ := json.Marshal(map[string]any{"items": newArr})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLargeObject(b *testing.B) {
	oldObj := make(map[string]any)
	newObj := make(map[string]any)
	for i := 0; i < 1000; i++ {
		k := fmt.Sprintf("key_%d", i)
		oldObj[k] = i
		if i%3 == 0 {
			newObj[k] = i * 10
		} else {
			newObj[k] = i
		}
	}

	oldJSON, _ := json.Marshal(oldObj)
	newJSON, _ := json.Marshal(newObj)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFewChangesInLargeInput(b *testing.B) {
	oldObj := make(map[string]any)
	newObj := make(map[string]any)
	for i := 0; i < 1000; i++ {
		k := fmt.Sprintf("key_%d", i)
		oldObj[k] = i
		newObj[k] = i
	}
	newObj["key_999"] = 999999

	oldJSON, _ := json.Marshal(oldObj)
	newJSON, _ := json.Marshal(newObj)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CompareBytes(oldJSON, newJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}
