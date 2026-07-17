package cgroup

import (
	"fmt"
	"strconv"
	"strings"
)

// parseCPUStat extracts usage_usec from a cpu.stat flat-keyed file.
func parseCPUStat(data string) (uint64, error) {
	v, ok := flatKeyValue(data, "usage_usec")
	if !ok {
		return 0, fmt.Errorf("missing usage_usec")
	}
	return strconv.ParseUint(v, 10, 64)
}

// parseCPUMax turns cpu.max ("<quota> <period>" or "max <period>") into a
// core count. "max" means unlimited and maps to 0 — the collector falls back
// to the host's CPU count for the percentage denominator.
func parseCPUMax(data string) (float64, error) {
	fields := strings.Fields(data)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty cpu.max")
	}
	if fields[0] == "max" {
		return 0, nil
	}
	quota, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, fmt.Errorf("quota: %w", err)
	}
	period := 100000.0 // kernel default when the file has one field
	if len(fields) > 1 {
		if period, err = strconv.ParseFloat(fields[1], 64); err != nil {
			return 0, fmt.Errorf("period: %w", err)
		}
	}
	if quota <= 0 || period <= 0 {
		return 0, fmt.Errorf("non-positive cpu.max %q", strings.TrimSpace(data))
	}
	return quota / period, nil
}

// parseMemoryStat extracts inactive_file from memory.stat. A missing key
// degrades to 0 (all memory counts as working set) rather than failing the
// sample.
func parseMemoryStat(data string) (uint64, error) {
	v, ok := flatKeyValue(data, "inactive_file")
	if !ok {
		return 0, nil
	}
	return strconv.ParseUint(v, 10, 64)
}

// parseUintValue parses a single-value file like memory.current.
func parseUintValue(data string) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(data), 10, 64)
}

// flatKeyValue finds "key value" in a cgroup flat-keyed file.
func flatKeyValue(data, key string) (string, bool) {
	for _, line := range strings.Split(data, "\n") {
		if v, ok := strings.CutPrefix(line, key+" "); ok {
			return strings.TrimSpace(v), true
		}
	}
	return "", false
}
