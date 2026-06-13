package domain

import "testing"

func TestParseMemoryLimit(t *testing.T) {
	valid := []struct {
		in   string
		want int64
	}{
		{"256m", 256 * 1024 * 1024},
		{"128M", 128 * 1024 * 1024},
		{"1g", 1024 * 1024 * 1024},
		{"512k", 512 * 1024},
		{"100b", 100},
		{"1048576", 1048576},
		{" 256m ", 256 * 1024 * 1024},
	}
	for _, tc := range valid {
		got, err := ParseMemoryLimit(tc.in)
		if err != nil {
			t.Errorf("ParseMemoryLimit(%q) returned error: %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseMemoryLimit(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}

	invalid := []string{"", "0", "-1", "-256m", "abc", "1t", "0.5m", "m", "256 m", "9999999999999g"}
	for _, in := range invalid {
		if got, err := ParseMemoryLimit(in); err == nil {
			t.Errorf("ParseMemoryLimit(%q) = %d, want error", in, got)
		}
	}
}

func TestParseCPULimit(t *testing.T) {
	valid := []struct {
		in   string
		want int64
	}{
		{"0.5", 500_000_000},
		{"0.25", 250_000_000},
		{"1", 1_000_000_000},
		{"1.5", 1_500_000_000},
		{"2", 2_000_000_000},
		{" 0.5 ", 500_000_000},
	}
	for _, tc := range valid {
		got, err := ParseCPULimit(tc.in)
		if err != nil {
			t.Errorf("ParseCPULimit(%q) returned error: %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseCPULimit(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}

	invalid := []string{"", "0", "-1", "abc", "NaN", "Inf", "-Inf", "1025", "1e30"}
	for _, in := range invalid {
		if got, err := ParseCPULimit(in); err == nil {
			t.Errorf("ParseCPULimit(%q) = %d, want error", in, got)
		}
	}
}
