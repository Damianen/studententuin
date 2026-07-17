package metrics

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"servermanager/internal/domain"
	"servermanager/internal/mocks"

	"go.uber.org/mock/gomock"
)

const (
	testAppUUID = "8b9f5e9e-9a3a-4b7e-9a59-1d2f3a4b5c6d"
	testDBUUID  = "6f1d2c3b-4a5e-4f60-8b7a-9c0d1e2f3a4b"
)

// stepClock lets tests advance time between ticks; shared by collector and
// store so query windows track the samples.
type stepClock struct{ t time.Time }

func (c *stepClock) Now() time.Time          { return c.t }
func (c *stepClock) advance(d time.Duration) { c.t = c.t.Add(d) }

type collectorHarness struct {
	runtime *mocks.MockContainerRuntime
	source  *mocks.MockMetricsSource
	exec    *mocks.MockContainerExec
	store   *Store
	clock   *stepClock
	c       *Collector
}

func newCollectorHarness(t *testing.T) *collectorHarness {
	t.Helper()
	ctrl := gomock.NewController(t)
	h := &collectorHarness{
		runtime: mocks.NewMockContainerRuntime(ctrl),
		source:  mocks.NewMockMetricsSource(ctrl),
		exec:    mocks.NewMockContainerExec(ctrl),
		clock:   &stepClock{t: testNow},
	}
	h.store = NewStore(30*time.Second, 24*time.Hour, h.clock)
	h.c = NewCollector(Dependencies{
		Runtime:   h.runtime,
		Source:    h.source,
		Exec:      h.exec,
		Store:     h.store,
		Clock:     h.clock,
		Interval:  30 * time.Second,
		Retention: 24 * time.Hour,
	})
	return h
}

func appContainer() domain.ManagedContainer {
	return domain.ManagedContainer{ID: "cid-app-1", Name: "stt-app-" + testAppUUID, Kind: domain.KindApp, OwnerID: testAppUUID}
}

func dbContainer() domain.ManagedContainer {
	return domain.ManagedContainer{ID: "cid-db-1", Name: "stt-db-" + testDBUUID, Kind: domain.KindDB, OwnerID: testDBUUID}
}

func TestCollectorCPUAcrossTicks(t *testing.T) {
	h := newCollectorHarness(t)
	ctx := context.Background()
	ctr := appContainer()

	h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return([]domain.ManagedContainer{ctr}, nil).Times(2)
	// 15s of cpu time over 30s wall on a 1-core limit = 50%.
	gomock.InOrder(
		h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{
			CPUUsageUsec: 1_000_000, CPULimitCores: 1, MemWorkingSetBytes: 512 * 1024 * 1024,
		}, nil),
		h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{
			CPUUsageUsec: 16_000_000, CPULimitCores: 1, MemWorkingSetBytes: 512 * 1024 * 1024,
		}, nil),
	)

	h.c.tick(ctx)
	// First sight: mem flows, cpu is baseline-only.
	series := h.store.Query(ctr.Name, AppSeriesKeys, time.Hour)
	if len(series["cpu"]) != 0 {
		t.Errorf("cpu after first tick = %+v, want no points (baseline only)", series["cpu"])
	}
	if len(series["mem"]) != 1 || series["mem"][0].Value != 512 {
		t.Errorf("mem after first tick = %+v, want one 512MB point", series["mem"])
	}

	h.clock.advance(30 * time.Second)
	h.c.tick(ctx)
	series = h.store.Query(ctr.Name, AppSeriesKeys, time.Hour)
	if len(series["cpu"]) != 1 || series["cpu"][0].Value != 50 {
		t.Errorf("cpu after second tick = %+v, want one 50%% point", series["cpu"])
	}
}

func TestCollectorCPUResetRebaselines(t *testing.T) {
	h := newCollectorHarness(t)
	ctx := context.Background()
	ctr := appContainer()

	h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return([]domain.ManagedContainer{ctr}, nil).Times(3)
	gomock.InOrder(
		h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{CPUUsageUsec: 90_000_000, CPULimitCores: 1}, nil),
		// Counter went backwards: restart. Skip, rebaseline.
		h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{CPUUsageUsec: 1_000_000, CPULimitCores: 1}, nil),
		h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{CPUUsageUsec: 16_000_000, CPULimitCores: 1}, nil),
	)

	for i := 0; i < 3; i++ {
		h.c.tick(ctx)
		h.clock.advance(30 * time.Second)
	}

	points := h.store.Query(ctr.Name, []string{"cpu"}, time.Hour)["cpu"]
	if len(points) != 1 || points[0].Value != 50 {
		t.Errorf("cpu points = %+v, want exactly the one post-reset 50%% point", points)
	}
}

func TestCollectorUnlimitedCPUFallsBackToHostCores(t *testing.T) {
	h := newCollectorHarness(t)
	ctx := context.Background()
	ctr := appContainer()

	h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return([]domain.ManagedContainer{ctr}, nil).Times(2)
	// 30s of cpu over 30s wall, unlimited limit → 100% of one host core; the
	// denominator becomes runtime.NumCPU(), so the value is 100/NumCPU — the
	// test only pins it is in (0, 100].
	gomock.InOrder(
		h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{CPUUsageUsec: 0}, nil),
		h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{CPUUsageUsec: 30_000_000}, nil),
	)

	h.c.tick(ctx)
	h.clock.advance(30 * time.Second)
	h.c.tick(ctx)

	points := h.store.Query(ctr.Name, []string{"cpu"}, time.Hour)["cpu"]
	if len(points) != 1 || points[0].Value <= 0 || points[0].Value > 100 {
		t.Errorf("cpu points = %+v, want one clamped point in (0,100]", points)
	}
}

func TestCollectorOneBadContainerDoesNotStopOthers(t *testing.T) {
	h := newCollectorHarness(t)
	ctx := context.Background()
	bad := appContainer()
	good := domain.ManagedContainer{ID: "cid-app-2", Name: "stt-app-11111111-2222-4333-8444-555555555555", Kind: domain.KindApp}

	h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return([]domain.ManagedContainer{bad, good}, nil)
	h.source.EXPECT().Sample(gomock.Any(), bad.ID).Return(domain.StatsSample{}, domain.ErrNotFound)
	h.source.EXPECT().Sample(gomock.Any(), good.ID).Return(domain.StatsSample{MemWorkingSetBytes: 1024 * 1024}, nil)

	h.c.tick(ctx)

	if points := h.store.Query(bad.Name, []string{"mem"}, time.Hour)["mem"]; len(points) != 0 {
		t.Errorf("bad container recorded %+v, want nothing", points)
	}
	if points := h.store.Query(good.Name, []string{"mem"}, time.Hour)["mem"]; len(points) != 1 || points[0].Value != 1 {
		t.Errorf("good container mem = %+v, want one 1MB point", points)
	}
}

func TestCollectorListFailureSkipsTick(t *testing.T) {
	h := newCollectorHarness(t)
	h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return(nil, errors.New("daemon down"))
	h.c.tick(context.Background()) // must not panic or record
}

func TestCollectorDBStats(t *testing.T) {
	h := newCollectorHarness(t)
	ctx := context.Background()
	ctr := dbContainer()

	h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return([]domain.ManagedContainer{ctr}, nil).Times(2)
	h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{CPULimitCores: 1, MemWorkingSetBytes: 64 * 1024 * 1024}, nil).Times(2)
	// Identity is fetched once and cached.
	h.exec.EXPECT().ContainerEnv(gomock.Any(), ctr.ID, "POSTGRES_USER", "POSTGRES_DB").
		Return(map[string]string{"POSTGRES_USER": "app", "POSTGRES_DB": "appdb"}, nil)

	var gotArgv []string
	gomock.InOrder(
		h.exec.EXPECT().ExecCapture(gomock.Any(), ctr.ID, gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, argv []string) (string, error) {
				gotArgv = argv
				return "3|100|10485760\n", nil
			}),
		h.exec.EXPECT().ExecCapture(gomock.Any(), ctr.ID, gomock.Any()).
			Return("5|400|10485760\n", nil),
	)

	h.c.tick(ctx)
	series := h.store.Query(ctr.Name, DBSeriesKeys, time.Hour)
	if len(series["conn"]) != 1 || series["conn"][0].Value != 3 {
		t.Errorf("conn = %+v, want one 3 point", series["conn"])
	}
	if len(series["disk"]) != 1 || series["disk"][0].Value != 10 {
		t.Errorf("disk = %+v, want one 10MB point", series["disk"])
	}
	if len(series["qps"]) != 0 {
		t.Errorf("qps after first tick = %+v, want none (baseline only)", series["qps"])
	}

	// Argv-shape pin (security regression): fixed flags, re-validated user,
	// the maintenance db, and the constant SQL — nothing else.
	wantArgv := []string{"psql", "-X", "-A", "-t", "-F", "|", "-U", "app", "-d", "postgres", "-c", pgStatsSQL}
	if !reflect.DeepEqual(gotArgv, wantArgv) {
		t.Errorf("psql argv = %q, want %q", gotArgv, wantArgv)
	}

	h.clock.advance(30 * time.Second)
	h.c.tick(ctx)
	series = h.store.Query(ctr.Name, DBSeriesKeys, time.Hour)
	// (400-100) transactions over 30s = 10 tps.
	if len(series["qps"]) != 1 || series["qps"][0].Value != 10 {
		t.Errorf("qps after second tick = %+v, want one 10 point", series["qps"])
	}
	if len(series["cpu"]) != 1 { // cgroup cpu flows for dbs too
		t.Errorf("db cpu = %+v, want one point", series["cpu"])
	}
}

func TestCollectorDBExecFailureDropsQPSBaseline(t *testing.T) {
	h := newCollectorHarness(t)
	ctx := context.Background()
	ctr := dbContainer()

	h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return([]domain.ManagedContainer{ctr}, nil).Times(3)
	h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{}, nil).Times(3)
	h.exec.EXPECT().ContainerEnv(gomock.Any(), ctr.ID, "POSTGRES_USER", "POSTGRES_DB").
		Return(map[string]string{"POSTGRES_USER": "app", "POSTGRES_DB": "appdb"}, nil)
	gomock.InOrder(
		h.exec.EXPECT().ExecCapture(gomock.Any(), ctr.ID, gomock.Any()).Return("1|100|1048576\n", nil),
		// Restarting container: exec 409s. qps baseline must drop.
		h.exec.EXPECT().ExecCapture(gomock.Any(), ctr.ID, gomock.Any()).Return("", domain.ErrConflict),
		h.exec.EXPECT().ExecCapture(gomock.Any(), ctr.ID, gomock.Any()).Return("1|50|1048576\n", nil),
	)

	for i := 0; i < 3; i++ {
		h.c.tick(ctx)
		h.clock.advance(30 * time.Second)
	}

	// Tick 3's xacts (50) is below tick 1's (100); had the baseline survived
	// the failed tick, that would be a negative delta — either way no qps
	// point may exist.
	if points := h.store.Query(ctr.Name, []string{"qps"}, time.Hour)["qps"]; len(points) != 0 {
		t.Errorf("qps = %+v, want none across a failure + counter reset", points)
	}
	// conn kept flowing on the good ticks.
	if points := h.store.Query(ctr.Name, []string{"conn"}, time.Hour)["conn"]; len(points) != 2 {
		t.Errorf("conn = %+v, want two points", points)
	}
}

func TestCollectorDBInvalidIdentitySkipsPsqlForever(t *testing.T) {
	h := newCollectorHarness(t)
	ctx := context.Background()
	ctr := dbContainer()

	h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return([]domain.ManagedContainer{ctr}, nil).Times(2)
	gomock.InOrder(
		h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{MemWorkingSetBytes: 1024 * 1024}, nil),
		h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{MemWorkingSetBytes: 3 * 1024 * 1024}, nil),
	)
	// The bad identifier fails the db-identifier rule. Fetched once, verdict
	// cached; ExecCapture has no expectation — any call fails the test.
	h.exec.EXPECT().ContainerEnv(gomock.Any(), ctr.ID, "POSTGRES_USER", "POSTGRES_DB").
		Return(map[string]string{"POSTGRES_USER": "Robert'); DROP", "POSTGRES_DB": "appdb"}, nil)

	h.c.tick(ctx)
	h.clock.advance(30 * time.Second)
	h.c.tick(ctx)

	// cpu/mem still flow for the container: both ticks recorded, and a sparse
	// window returns them raw.
	if points := h.store.Query(ctr.Name, []string{"mem"}, time.Hour)["mem"]; len(points) != 2 || points[0].Value != 1 || points[1].Value != 3 {
		t.Errorf("mem = %+v, want both raw samples despite the identity failure", points)
	}
	if points := h.store.Query(ctr.Name, []string{"conn"}, time.Hour)["conn"]; len(points) != 0 {
		t.Errorf("conn = %+v, want none", points)
	}
}

func TestCollectorPrunesBaselinesForGoneContainers(t *testing.T) {
	h := newCollectorHarness(t)
	ctx := context.Background()
	ctr := appContainer()

	gomock.InOrder(
		h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return([]domain.ManagedContainer{ctr}, nil),
		h.runtime.EXPECT().ListManagedContainers(gomock.Any()).Return(nil, nil),
	)
	h.source.EXPECT().Sample(gomock.Any(), ctr.ID).Return(domain.StatsSample{}, nil)

	h.c.tick(ctx)
	if _, ok := h.c.cpuBaselines[ctr.ID]; !ok {
		t.Fatal("baseline not recorded on first sight")
	}
	h.c.tick(ctx)
	if _, ok := h.c.cpuBaselines[ctr.ID]; ok {
		t.Error("baseline for a gone container not pruned")
	}
}

func TestParsePGStats(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		conn, xacts, bytes, err := parsePGStats("3|1234|567890\n")
		if err != nil || conn != 3 || xacts != 1234 || bytes != 567890 {
			t.Errorf("parsePGStats = %v %v %v %v", conn, xacts, bytes, err)
		}
	})

	t.Run("rejects NULLs and garbage", func(t *testing.T) {
		for _, out := range []string{
			"",
			"3|1234",
			"3|1234|567890|extra",
			"3||567890",
			"||",
			"psql: error: connection to server failed",
			"three|four|five",
		} {
			if _, _, _, err := parsePGStats(out); err == nil {
				t.Errorf("parsePGStats(%q) = nil error, want parse failure", out)
			}
		}
	})
}
