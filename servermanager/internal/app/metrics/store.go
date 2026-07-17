// Package metrics holds the in-memory metrics pipeline (§6 phase 6): a
// background collector samples managed containers into a ring-buffer store
// the apps/databases metrics endpoints query. History is lost on restart by
// design — the MVP keeps the manager stateless on disk.
package metrics

import (
	"fmt"
	"sync"
	"time"

	"servermanager/internal/domain"
	"servermanager/internal/ports"
)

// maxQueryPoints caps a Query response: the window is averaged into at most
// this many buckets, which is plenty for a dashboard chart.
const maxQueryPoints = 48

// Series keys are the frontend contract. Every container gets cpu/mem;
// databases add postgres vitals.
var (
	AppSeriesKeys = []string{"cpu", "mem"}
	DBSeriesKeys  = []string{"conn", "qps", "cpu", "disk"}
)

// Range is a validated metrics window; Label is the canonical wire form.
type Range struct {
	Label    string
	Duration time.Duration
}

// ParseRange validates ?range= against the allowlist. Empty defaults to 24h,
// anything else is domain.ErrInvalid — the window sizes bucket math and must
// never be free text.
func ParseRange(raw string) (Range, error) {
	switch raw {
	case "", "24h":
		return Range{Label: "24h", Duration: 24 * time.Hour}, nil
	case "1h":
		return Range{Label: "1h", Duration: time.Hour}, nil
	default:
		return Range{}, fmt.Errorf("range %q must be one of 1h, 24h: %w", raw, domain.ErrInvalid)
	}
}

// Store keeps per-series ring buffers keyed by container NAME — names are
// stable across redeploys while container IDs change at every cutover, so a
// redeployed app keeps its history.
type Store struct {
	mu       sync.RWMutex
	capacity int
	clock    ports.Clock
	series   map[string]map[string]*ring // container name → series key → ring
}

// NewStore sizes the rings for retention/interval samples plus slack for
// jittered ticks.
func NewStore(interval, retention time.Duration, clock ports.Clock) *Store {
	capacity := int(retention/interval) + 2
	return &Store{
		capacity: capacity,
		clock:    clock,
		series:   map[string]map[string]*ring{},
	}
}

// Record appends one value per series key for the container.
func (s *Store) Record(name string, at time.Time, values map[string]float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyed, ok := s.series[name]
	if !ok {
		keyed = map[string]*ring{}
		s.series[name] = keyed
	}
	for key, value := range values {
		r, ok := keyed[key]
		if !ok {
			r = newRing(s.capacity)
			keyed[key] = r
		}
		r.push(domain.MetricPoint{Time: at, Value: value})
	}
}

// Query returns the requested series over the window ending now. Sparse
// windows (≤ maxQueryPoints samples) return the raw points, so a fresh
// container gets a chart after two ticks instead of waiting out a bucket;
// denser windows are averaged into fixed-width buckets whose time is the
// bucket END, with empty buckets omitted (chart gaps). Every requested key is
// present in the result — as an empty slice when there is no data, never nil.
func (s *Store) Query(name string, keys []string, window time.Duration) map[string][]domain.MetricPoint {
	now := s.clock.Now()
	start := now.Add(-window)
	width := window / maxQueryPoints

	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string][]domain.MetricPoint, len(keys))
	keyed := s.series[name]
	for _, key := range keys {
		points := []domain.MetricPoint{}
		if r, ok := keyed[key]; ok {
			points = r.window(start, now)
			if len(points) > maxQueryPoints {
				points = bucketAverage(points, start, width)
			}
		}
		out[key] = points
	}
	return out
}

// DropStale removes points older than cutoff and forgets containers whose
// series all emptied (removed apps/databases age out of memory).
func (s *Store) DropStale(cutoff time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, keyed := range s.series {
		for key, r := range keyed {
			r.dropBefore(cutoff)
			if r.size == 0 {
				delete(keyed, key)
			}
		}
		if len(keyed) == 0 {
			delete(s.series, name)
		}
	}
}

// bucketAverage folds window-filtered, oldest-first points into fixed-width
// buckets.
func bucketAverage(points []domain.MetricPoint, start time.Time, width time.Duration) []domain.MetricPoint {
	var sums [maxQueryPoints]float64
	var counts [maxQueryPoints]int
	for _, p := range points {
		idx := int(p.Time.Sub(start) / width)
		if idx >= maxQueryPoints { // p.Time == end lands past the last bucket
			idx = maxQueryPoints - 1
		}
		sums[idx] += p.Value
		counts[idx]++
	}

	buckets := []domain.MetricPoint{}
	for i := range counts {
		if counts[i] == 0 {
			continue
		}
		buckets = append(buckets, domain.MetricPoint{
			Time:  start.Add(time.Duration(i+1) * width),
			Value: sums[i] / float64(counts[i]),
		})
	}
	return buckets
}

// ring is a fixed-capacity buffer of chronologically pushed points.
type ring struct {
	buf  []domain.MetricPoint
	next int // write position
	size int
}

func newRing(capacity int) *ring {
	return &ring{buf: make([]domain.MetricPoint, capacity)}
}

func (r *ring) push(p domain.MetricPoint) {
	r.buf[r.next] = p
	r.next = (r.next + 1) % len(r.buf)
	if r.size < len(r.buf) {
		r.size++
	}
}

// window copies out the points inside (start, end], oldest-first.
func (r *ring) window(start, end time.Time) []domain.MetricPoint {
	points := []domain.MetricPoint{}
	r.forEach(func(p domain.MetricPoint) {
		if !p.Time.After(start) || p.Time.After(end) {
			return
		}
		points = append(points, p)
	})
	return points
}

// forEach visits the retained points oldest-first.
func (r *ring) forEach(fn func(domain.MetricPoint)) {
	start := r.next - r.size
	if start < 0 {
		start += len(r.buf)
	}
	for i := 0; i < r.size; i++ {
		fn(r.buf[(start+i)%len(r.buf)])
	}
}

// dropBefore removes points older than cutoff. Points are pushed in time
// order, so this only trims from the oldest end.
func (r *ring) dropBefore(cutoff time.Time) {
	start := r.next - r.size
	if start < 0 {
		start += len(r.buf)
	}
	for r.size > 0 {
		oldest := r.buf[start%len(r.buf)]
		if !oldest.Time.Before(cutoff) {
			return
		}
		r.buf[start%len(r.buf)] = domain.MetricPoint{}
		start++
		r.size--
	}
}
