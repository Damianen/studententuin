//go:build integration

package deployments

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"servermanager/internal/app/apps"
	"servermanager/internal/domain"
	buildinfra "servermanager/internal/infra/build"
	clockinfra "servermanager/internal/infra/clock"
	"servermanager/internal/infra/docker"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
)

// localFetcher satisfies ports.SourceFetcher against an on-disk repo: the
// pipeline test is about clone→build→cutover wiring, and the real https
// fetcher has its own integration tests in infra/git.
type localFetcher struct{ repoDir string }

func (f localFetcher) Fetch(ctx context.Context, _, _ string) (*domain.SourceCheckout, error) {
	dir, err := os.MkdirTemp("", "stt-deploy-clone-*")
	if err != nil {
		return nil, err
	}
	repo, err := gogit.PlainCloneContext(ctx, dir, false, &gogit.CloneOptions{
		URL: f.repoDir, Depth: 1, SingleBranch: true, Tags: gogit.NoTags,
	})
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	head, err := repo.Head()
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return &domain.SourceCheckout{
		Dir:       dir,
		CommitSHA: head.Hash().String(),
		Cleanup:   domain.OnceCleanup(func() { os.RemoveAll(dir) }),
	}, nil
}

// fixtureRepo commits the nodeapp build fixture into a fresh local repo.
func fixtureRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"package.json", "index.js"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "infra", "build", "testdata", "nodeapp", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := wt.Add(name); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := wt.Commit("fixture", &gogit.CommitOptions{
		Author: &object.Signature{Name: "fixture", Email: "fixture@test.local", When: time.Now()},
	}); err != nil {
		t.Fatal(err)
	}
	return dir
}

// TestIntegrationDeployPipeline is the manager-side half of the §8 e2e
// smoke: two consecutive deploys, asserting cutover and image GC.
func TestIntegrationDeployPipeline(t *testing.T) {
	ctx := context.Background()

	runtime, err := docker.New(ctx)
	if err != nil {
		t.Fatalf("docker daemon unreachable (integration tests need one): %v", err)
	}
	t.Cleanup(func() { runtime.Close() })

	builder, err := buildinfra.New(ctx, "nixpacks", runtime)
	if err != nil {
		t.Fatalf("nixpacks unavailable (integration tests need it): %v", err)
	}

	appID := uuid.NewString()
	limits := apps.Limits{
		DefaultMemoryBytes: 256 * 1024 * 1024,
		MaxMemoryBytes:     1024 * 1024 * 1024,
		DefaultNanoCPUs:    500_000_000,
		MaxNanoCPUs:        2_000_000_000,
		DefaultPidsLimit:   256,
		DefaultRuntime:     domain.RuntimeRunc,
	}
	svc := NewService(Dependencies{
		Fetcher:  localFetcher{repoDir: fixtureRepo(t)},
		Builder:  builder,
		Runtime:  runtime,
		Images:   runtime,
		Runner:   apps.NewService(apps.Dependencies{Runtime: runtime, Limits: limits}).Run,
		Clock:    clockinfra.System{},
		GitHosts: []string{"github.com"},
		Budgets: Budgets{
			CloneTimeout: time.Minute,
			BuildTimeout: 8 * time.Minute,
			HealthGrace:  2 * time.Second,
		},
	})

	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_ = runtime.Remove(cleanupCtx, domain.AppContainerName(appID))
		if refs, err := runtime.ListAppImages(cleanupCtx, appID); err == nil {
			for _, ref := range refs {
				_ = runtime.RemoveImage(cleanupCtx, ref)
			}
		}
	})

	deploy := func(t *testing.T) domain.DeploymentJob {
		t.Helper()
		id, err := svc.Deploy.Execute(ctx, DeployInput{
			AppID:         appID,
			RepositoryURL: "https://github.com/fixture/nodeapp", // validated, then redirected by localFetcher
			Port:          3000,
			Env:           map[string]string{"DEPLOY_MARKER": "integration"},
		})
		if err != nil {
			t.Fatalf("Deploy.Execute: %v", err)
		}
		deadline := time.Now().Add(9 * time.Minute)
		for {
			job, err := svc.Get.Execute(id)
			if err != nil {
				t.Fatalf("Get: %v", err)
			}
			if job.Terminal() {
				if job.Status != domain.DeploymentStatusRunning {
					t.Fatalf("deployment failed: %s\n--- build log ---\n%s", job.Error, job.BuildLog)
				}
				return job
			}
			if time.Now().After(deadline) {
				t.Fatalf("deployment stuck in %s", job.Status)
			}
			time.Sleep(time.Second)
		}
	}

	first := deploy(t)
	state, err := runtime.Inspect(ctx, first.ContainerID)
	if err != nil || !state.Running {
		t.Fatalf("container after first deploy: state=%+v err=%v", state, err)
	}
	if first.CommitSHA == "" || first.Image == "" {
		t.Errorf("job missing source info: %+v", first)
	}

	second := deploy(t)
	if second.ContainerID == first.ContainerID {
		t.Error("second deploy reused the first container")
	}
	if state, err := runtime.Inspect(ctx, second.ContainerID); err != nil || !state.Running {
		t.Fatalf("container after second deploy: state=%+v err=%v", state, err)
	}
	// The first deploy's image is GC'd; only the latest remains.
	refs, err := runtime.ListAppImages(ctx, appID)
	if err != nil {
		t.Fatalf("ListAppImages: %v", err)
	}
	if len(refs) != 1 || refs[0] != second.Image {
		t.Errorf("images after second deploy = %v, want only %s", refs, second.Image)
	}
}
