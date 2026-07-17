package cgroup

import (
	"context"
	"errors"
	"strings"
	"testing"

	"servermanager/internal/domain"
)

var (
	systemdID  = strings.Repeat("a", 64)
	cgroupfsID = strings.Repeat("b", 64)
)

func TestSampleSystemdLayout(t *testing.T) {
	s := New("testdata/systemd")

	sample, err := s.Sample(context.Background(), systemdID)
	if err != nil {
		t.Fatalf("Sample: %v", err)
	}
	if sample.CPUUsageUsec != 1234567 {
		t.Errorf("CPUUsageUsec = %d, want 1234567", sample.CPUUsageUsec)
	}
	if sample.CPULimitCores != 0.5 {
		t.Errorf("CPULimitCores = %v, want 0.5", sample.CPULimitCores)
	}
	// 100 MiB current − 4 MiB inactive_file = 96 MiB working set.
	if want := uint64(104857600 - 4194304); sample.MemWorkingSetBytes != want {
		t.Errorf("MemWorkingSetBytes = %d, want %d", sample.MemWorkingSetBytes, want)
	}
}

func TestSampleCgroupfsLayout(t *testing.T) {
	s := New("testdata/cgroupfs")

	sample, err := s.Sample(context.Background(), cgroupfsID)
	if err != nil {
		t.Fatalf("Sample: %v", err)
	}
	if sample.CPUUsageUsec != 42 {
		t.Errorf("CPUUsageUsec = %d, want 42", sample.CPUUsageUsec)
	}
	// "max <period>" = unlimited = 0 cores (collector falls back to host CPUs).
	if sample.CPULimitCores != 0 {
		t.Errorf("CPULimitCores = %v, want 0 for unlimited", sample.CPULimitCores)
	}
	// inactive_file above memory.current clamps to 0, never negative.
	if sample.MemWorkingSetBytes != 0 {
		t.Errorf("MemWorkingSetBytes = %d, want clamped 0", sample.MemWorkingSetBytes)
	}
}

func TestSampleMissingContainer(t *testing.T) {
	s := New("testdata/systemd")
	if _, err := s.Sample(context.Background(), strings.Repeat("c", 64)); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Sample(missing) = %v, want ErrNotFound", err)
	}
}

func TestSampleRejectsBadContainerID(t *testing.T) {
	s := New("testdata/systemd")
	for _, id := range []string{
		"",
		"short",
		"../../../../etc/passwd",
		strings.Repeat("a", 63) + "/", // path meta after valid hex
		strings.Repeat("A", 64),       // uppercase is not a docker id
		strings.Repeat("a", 65),       // too long
	} {
		if _, err := s.Sample(context.Background(), id); !errors.Is(err, domain.ErrInvalid) {
			t.Errorf("Sample(%q) = %v, want ErrInvalid", id, err)
		}
	}
}

// TestNewToleratesMissingRoot pins that a bad root warns but never fails —
// metrics are non-critical and must not block manager startup.
func TestNewToleratesMissingRoot(t *testing.T) {
	s := New(t.TempDir())
	if _, err := s.Sample(context.Background(), systemdID); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("Sample on empty root = %v, want ErrNotFound", err)
	}
}

func TestParseCPUMax(t *testing.T) {
	cases := []struct {
		in      string
		want    float64
		wantErr bool
	}{
		{"50000 100000\n", 0.5, false},
		{"200000 100000\n", 2, false},
		{"max 100000\n", 0, false},
		{"max\n", 0, false},
		{"100000\n", 1, false}, // single field uses the kernel default period
		{"", 0, true},
		{"abc 100000\n", 0, true},
		{"50000 0\n", 0, true},
		{"-1 100000\n", 0, true},
	}
	for _, tc := range cases {
		got, err := parseCPUMax(tc.in)
		if (err != nil) != tc.wantErr {
			t.Errorf("parseCPUMax(%q) err = %v, wantErr %v", tc.in, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && got != tc.want {
			t.Errorf("parseCPUMax(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestParseCPUStat(t *testing.T) {
	if _, err := parseCPUStat("user_usec 1\nsystem_usec 2\n"); err == nil {
		t.Error("parseCPUStat without usage_usec = nil, want error")
	}
	if got, err := parseCPUStat("usage_usec 99\nuser_usec 1\n"); err != nil || got != 99 {
		t.Errorf("parseCPUStat = %d, %v; want 99, nil", got, err)
	}
}

func TestParseMemoryStatMissingKeyIsZero(t *testing.T) {
	if got, err := parseMemoryStat("anon 5\n"); err != nil || got != 0 {
		t.Errorf("parseMemoryStat without inactive_file = %d, %v; want 0, nil", got, err)
	}
}
