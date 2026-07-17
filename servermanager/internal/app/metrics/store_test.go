package metrics

import (
	"errors"
	"sync"
	"testing"
	"time"

	"servermanager/internal/domain"
)

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

var testNow = time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)

func testStore() (*Store, time.Time) {
	return NewStore(30*time.Second, 24*time.Hour, fixedClock{t: testNow}), testNow
}

func TestParseRange(t *testing.T) {
	for raw, want := range map[string]Range{
		"":    {Label: "24h", Duration: 24 * time.Hour},
		"24h": {Label: "24h", Duration: 24 * time.Hour},
		"1h":  {Label: "1h", Duration: time.Hour},
	} {
		got, err := ParseRange(raw)
		if err != nil || got != want {
			t.Errorf("ParseRange(%q) = %+v, %v; want %+v", raw, got, err, want)
		}
	}
	for _, raw := range []string{"7d", "60", "1h ", "24H", "-1h"} {
		if _, err := ParseRange(raw); !errors.Is(err, domain.ErrInvalid) {
			t.Errorf("ParseRange(%q) = %v, want ErrInvalid", raw, err)
		}
	}
}

func TestQueryEmptyStoreHasEveryKey(t *testing.T) {
	s, _ := testStore()

	series := s.Query("stt-app-unknown", DBSeriesKeys, 24*time.Hour)
	if len(series) != len(DBSeriesKeys) {
		t.Fatalf("got %d keys, want %d", len(series), len(DBSeriesKeys))
	}
	for _, key := range DBSeriesKeys {
		points, ok := series[key]
		if !ok || points == nil {
			t.Errorf("series[%q] = %v, want empty non-nil slice", key, points)
		}
		if len(points) != 0 {
			t.Errorf("series[%q] has %d points, want 0", key, len(points))
		}
	}
}

func TestQuerySparseWindowReturnsRawPoints(t *testing.T) {
	s, now := testStore()

	// Outside the 1h window: must not appear.
	s.Record("app", now.Add(-2*time.Hour), map[string]float64{"cpu": 99})
	// A handful of in-window samples (≤ maxQueryPoints) come back verbatim —
	// a fresh container must chart after two ticks, not after a 24h/48 bucket
	// fills (the review's phase-6 blocker).
	times := []time.Time{
		now.Add(-10 * time.Minute),
		now.Add(-10 * time.Minute).Add(30 * time.Second),
		now.Add(-1 * time.Minute),
	}
	for i, at := range times {
		s.Record("app", at, map[string]float64{"cpu": float64(10 * (i + 1))})
	}

	points := s.Query("app", []string{"cpu"}, time.Hour)["cpu"]
	if len(points) != len(times) {
		t.Fatalf("got %d points, want %d raw: %+v", len(points), len(times), points)
	}
	for i, p := range points {
		if !p.Time.Equal(times[i]) || p.Value != float64(10*(i+1)) {
			t.Errorf("point %d = %+v, want raw sample at %v value %d", i, p, times[i], 10*(i+1))
		}
	}
}

func TestQueryDenseWindowBucketAverages(t *testing.T) {
	s, now := testStore()

	// Outside the 1h window: must not appear.
	s.Record("app", now.Add(-2*time.Hour), map[string]float64{"cpu": 99})
	// One point per 30s over the hour (120 > maxQueryPoints) triggers
	// bucketing; a constant value pins mean-not-sum averaging.
	for i := 0; i < 120; i++ {
		s.Record("app", now.Add(-time.Duration(i)*30*time.Second), map[string]float64{"cpu": 7})
	}

	points := s.Query("app", []string{"cpu"}, time.Hour)["cpu"]
	if len(points) == 0 || len(points) > maxQueryPoints {
		t.Fatalf("got %d points, want 1..%d", len(points), maxQueryPoints)
	}

	// Bucket time is the bucket END, on the window's grid; averages of a
	// constant stay that constant.
	width := time.Hour / maxQueryPoints
	start := now.Add(-time.Hour)
	for i, p := range points {
		if p.Value != 7 {
			t.Errorf("bucket %d average = %v, want 7 (window leak or sum-not-mean)", i, p.Value)
		}
		offset := p.Time.Sub(start)
		if offset%width != 0 || offset <= 0 || offset > time.Hour {
			t.Errorf("point time %v is not a bucket end", p.Time)
		}
		if i > 0 && !points[i-1].Time.Before(p.Time) {
			t.Errorf("points out of order at %d: %+v", i, points)
		}
	}
}

func TestQueryCapsAt48Points(t *testing.T) {
	s, now := testStore()
	// One point per 30s over a full hour = 120 samples.
	for i := 0; i < 120; i++ {
		s.Record("app", now.Add(-time.Duration(i)*30*time.Second), map[string]float64{"cpu": 1})
	}
	if points := s.Query("app", []string{"cpu"}, time.Hour)["cpu"]; len(points) > maxQueryPoints {
		t.Errorf("got %d points, want <= %d", len(points), maxQueryPoints)
	}
}

func TestRingWraparound(t *testing.T) {
	r := newRing(4)
	for i := 0; i < 10; i++ {
		r.push(domain.MetricPoint{Time: testNow.Add(time.Duration(i) * time.Second), Value: float64(i)})
	}

	var got []float64
	r.forEach(func(p domain.MetricPoint) { got = append(got, p.Value) })
	want := []float64{6, 7, 8, 9} // only the newest capacity survive, oldest-first
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestDropStale(t *testing.T) {
	s, now := testStore()
	s.Record("app", now.Add(-25*time.Hour), map[string]float64{"cpu": 1})
	s.Record("app", now.Add(-1*time.Minute), map[string]float64{"cpu": 2})
	s.Record("gone", now.Add(-25*time.Hour), map[string]float64{"cpu": 3})

	s.DropStale(now.Add(-24 * time.Hour))

	if points := s.Query("app", []string{"cpu"}, 24*time.Hour)["cpu"]; len(points) != 1 || points[0].Value != 2 {
		t.Errorf("app series after DropStale = %+v, want only the fresh point", points)
	}
	s.mu.RLock()
	_, staleKept := s.series["gone"]
	s.mu.RUnlock()
	if staleKept {
		t.Error("fully stale container name not forgotten")
	}
}

// TestStoreConcurrency is the -race hammer: recorders, queriers, and the
// stale sweep on the same store.
func TestStoreConcurrency(t *testing.T) {
	s, now := testStore()

	var wg sync.WaitGroup
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				s.Record("app", now.Add(-time.Duration(i)*time.Second), map[string]float64{"cpu": 1, "mem": 2})
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			s.Query("app", AppSeriesKeys, time.Hour)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			s.DropStale(now.Add(-24 * time.Hour))
		}
	}()
	wg.Wait()
}
