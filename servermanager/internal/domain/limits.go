package domain

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ParseMemoryLimit converts a human-readable memory limit ("256m", "1g",
// "1048576") into bytes. Suffixes are 1024-based, matching Docker's units.
func ParseMemoryLimit(s string) (int64, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	if v == "" {
		return 0, fmt.Errorf("memory limit is empty")
	}

	multiplier := int64(1)
	switch v[len(v)-1] {
	case 'b':
		v = v[:len(v)-1]
	case 'k':
		multiplier = 1024
		v = v[:len(v)-1]
	case 'm':
		multiplier = 1024 * 1024
		v = v[:len(v)-1]
	case 'g':
		multiplier = 1024 * 1024 * 1024
		v = v[:len(v)-1]
	}

	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory limit %q", s)
	}
	if n <= 0 {
		return 0, fmt.Errorf("memory limit must be positive, got %q", s)
	}
	if n > math.MaxInt64/multiplier {
		return 0, fmt.Errorf("memory limit %q overflows", s)
	}
	return n * multiplier, nil
}

// MaxCPULimit is a sanity bound on parsed CPU limits, far above any real host.
const MaxCPULimit = 1024

// ParseCPULimit converts a CPU limit ("0.5", "2") into NanoCPUs as used by
// Docker's HostConfig (1 CPU = 1e9).
func ParseCPULimit(s string) (int64, error) {
	v := strings.TrimSpace(s)
	if v == "" {
		return 0, fmt.Errorf("cpu limit is empty")
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("invalid cpu limit %q", s)
	}
	if f <= 0 {
		return 0, fmt.Errorf("cpu limit must be positive, got %q", s)
	}
	if f > MaxCPULimit {
		return 0, fmt.Errorf("cpu limit %q exceeds %d CPUs", s, MaxCPULimit)
	}
	return int64(math.Round(f * 1e9)), nil
}
