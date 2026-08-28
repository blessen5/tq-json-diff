package diff

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
)

// ToleranceOptions configures numeric and temporal fuzzy matching tolerances.
type ToleranceOptions struct {
	NumericDelta   float64       // Absolute numeric difference threshold
	NumericPercent float64       // Percentage numeric difference threshold (e.g. 0.05 for 5%)
	TimeDelta      time.Duration // Maximum allowed duration drift between timestamps
}

// IsEnabled returns true if any fuzzy tolerance rule is active.
func (t ToleranceOptions) IsEnabled() bool {
	return t.NumericDelta > 0 || t.NumericPercent > 0 || t.TimeDelta > 0
}

// NumbersWithinTolerance returns true if two numbers are within the specified numeric tolerance.
func NumbersWithinTolerance(a, b json.Number, opts ToleranceOptions) bool {
	fA, errA := a.Float64()
	fB, errB := b.Float64()
	if errA != nil || errB != nil {
		return false
	}

	diff := math.Abs(fA - fB)

	// Absolute delta check
	if opts.NumericDelta > 0 && diff <= opts.NumericDelta {
		return true
	}

	// Percentage check
	if opts.NumericPercent > 0 {
		base := math.Abs(fA)
		if base == 0 {
			base = math.Abs(fB)
		}
		if base > 0 {
			pct := diff / base
			if pct <= (opts.NumericPercent / 100.0) {
				return true
			}
		}
	}

	return false
}

var commonTimeFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// TimestampsWithinTolerance returns true if both strings are valid timestamps within TimeDelta.
func TimestampsWithinTolerance(sA, sB string, opts ToleranceOptions) bool {
	if opts.TimeDelta <= 0 {
		return false
	}

	tA, okA := parseTime(sA)
	if !okA {
		return false
	}
	tB, okB := parseTime(sB)
	if !okB {
		return false
	}

	drift := tA.Sub(tB)
	if drift < 0 {
		drift = -drift
	}

	return drift <= opts.TimeDelta
}

func parseTime(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	for _, fmtStr := range commonTimeFormats {
		if t, err := time.Parse(fmtStr, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// ParseNumericTolerance parses a numeric tolerance string (e.g. "0.01" or "5%").
func ParseNumericTolerance(s string) (delta float64, percent float64, err error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		p, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil || p < 0 {
			return 0, 0, err
		}
		return 0, p, nil
	}

	d, err := strconv.ParseFloat(s, 64)
	if err != nil || d < 0 {
		return 0, 0, err
	}
	return d, 0, nil
}
