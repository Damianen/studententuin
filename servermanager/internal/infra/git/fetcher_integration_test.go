//go:build integration

package git

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// fixtureRepo builds a local repository with one commit on the default
// branch and one on "feature", entirely with go-git.
func fixtureRepo(t *testing.T) (dir, mainSHA, featureSHA string) {
	t.Helper()
	dir = t.TempDir()

	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("init fixture repo: %v", err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("worktree: %v", err)
	}

	commit := func(name, content string) string {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := wt.Add(name); err != nil {
			t.Fatal(err)
		}
		sha, err := wt.Commit("add "+name, &gogit.CommitOptions{
			Author: &object.Signature{Name: "fixture", Email: "fixture@test.local", When: time.Now()},
		})
		if err != nil {
			t.Fatalf("commit: %v", err)
		}
		return sha.String()
	}

	mainSHA = commit("README.md", "hello\n")

	if err := wt.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("feature"),
		Create: true,
	}); err != nil {
		t.Fatalf("checkout feature: %v", err)
	}
	featureSHA = commit("feature.txt", "feature\n")

	// Leave HEAD on the default branch so that is what a plain clone gets.
	if err := wt.Checkout(&gogit.CheckoutOptions{Branch: "refs/heads/master"}); err != nil {
		t.Fatalf("checkout back to master: %v", err)
	}

	return dir, mainSHA, featureSHA
}

func TestIntegrationFetchLocalRepo(t *testing.T) {
	dir, mainSHA, featureSHA := fixtureRepo(t)
	f := NewFetcher(testHosts, 10<<20)
	ctx := context.Background()

	t.Run("default branch", func(t *testing.T) {
		co, err := f.fetchInto(ctx, dir, "")
		if err != nil {
			t.Fatalf("fetchInto: %v", err)
		}
		defer co.Cleanup()

		if co.CommitSHA != mainSHA {
			t.Errorf("CommitSHA = %s, want %s", co.CommitSHA, mainSHA)
		}
		if _, err := os.Stat(filepath.Join(co.Dir, "README.md")); err != nil {
			t.Errorf("README.md missing in checkout: %v", err)
		}
		if _, err := os.Stat(filepath.Join(co.Dir, "feature.txt")); err == nil {
			t.Error("feature.txt present in default-branch checkout")
		}

		co.Cleanup()
		if _, err := os.Stat(co.Dir); !os.IsNotExist(err) {
			t.Errorf("clone dir still exists after Cleanup: %v", err)
		}
	})

	t.Run("named branch", func(t *testing.T) {
		co, err := f.fetchInto(ctx, dir, "feature")
		if err != nil {
			t.Fatalf("fetchInto(feature): %v", err)
		}
		defer co.Cleanup()

		if co.CommitSHA != featureSHA {
			t.Errorf("CommitSHA = %s, want %s", co.CommitSHA, featureSHA)
		}
		if _, err := os.Stat(filepath.Join(co.Dir, "feature.txt")); err != nil {
			t.Errorf("feature.txt missing in feature checkout: %v", err)
		}
	})

	t.Run("missing branch", func(t *testing.T) {
		if _, err := f.fetchInto(ctx, dir, "no-such-branch"); err == nil || !strings.Contains(err.Error(), "not found") {
			t.Errorf("fetchInto(missing branch) = %v, want branch-not-found error", err)
		}
	})

	t.Run("size cap", func(t *testing.T) {
		tiny := NewFetcher(testHosts, 1)
		if _, err := tiny.fetchInto(ctx, dir, ""); err == nil || !strings.Contains(err.Error(), "size limit") {
			t.Errorf("fetchInto(over cap) = %v, want size limit error", err)
		}
	})

	t.Run("unknown path", func(t *testing.T) {
		if _, err := f.fetchInto(ctx, filepath.Join(t.TempDir(), "missing"), ""); err == nil {
			t.Error("fetchInto(unknown path) = nil, want error")
		}
	})
}
