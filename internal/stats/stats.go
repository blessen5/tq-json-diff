package stats

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Stats holds performance, resource, and execution metrics for a comparison.
type Stats struct {
	OldSize      int64         `json:"old_size"`
	OldIsStdin   bool          `json:"old_is_stdin"`
	NewSize      int64         `json:"new_size"`
	NewIsStdin   bool          `json:"new_is_stdin"`
	ParseTime    time.Duration `json:"parse_time_ms"`
	CompareTime  time.Duration `json:"compare_time_ms"`
	TotalTime    time.Duration `json:"total_time_ms"`
	AllocBytes   uint64        `json:"alloc_bytes"`
	ChangesCount int           `json:"changes_count"`
}

// FormatBytes formats a byte count into a human-readable string (B, KB, MB, GB).
func FormatBytes(b int64) string {
	if b < 0 {
		return "0 B"
	}
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
	)

	switch {
	case b >= gb:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(gb))
	case b >= mb:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(mb))
	case b >= kb:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(kb))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

// ParseSize parses a size string with optional units (e.g. "100B", "50KB", "10MB", "1GB").
func ParseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, fmt.Errorf("empty size string")
	}

	var multiplier int64 = 1
	var numStr string

	switch {
	case strings.HasSuffix(s, "GB"):
		multiplier = 1024 * 1024 * 1024
		numStr = strings.TrimSpace(strings.TrimSuffix(s, "GB"))
	case strings.HasSuffix(s, "G"):
		multiplier = 1024 * 1024 * 1024
		numStr = strings.TrimSpace(strings.TrimSuffix(s, "G"))
	case strings.HasSuffix(s, "MB"):
		multiplier = 1024 * 1024
		numStr = strings.TrimSpace(strings.TrimSuffix(s, "MB"))
	case strings.HasSuffix(s, "M"):
		multiplier = 1024 * 1024
		numStr = strings.TrimSpace(strings.TrimSuffix(s, "M"))
	case strings.HasSuffix(s, "KB"):
		multiplier = 1024
		numStr = strings.TrimSpace(strings.TrimSuffix(s, "KB"))
	case strings.HasSuffix(s, "K"):
		multiplier = 1024
		numStr = strings.TrimSpace(strings.TrimSuffix(s, "K"))
	case strings.HasSuffix(s, "B"):
		multiplier = 1
		numStr = strings.TrimSpace(strings.TrimSuffix(s, "B"))
	default:
		multiplier = 1
		numStr = s
	}

	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil || val < 0 {
		return 0, fmt.Errorf("invalid size format: %q", s)
	}

	return int64(val * float64(multiplier)), nil
}
