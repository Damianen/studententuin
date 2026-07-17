package git

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"servermanager/internal/domain"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

var testHosts = []string{"github.com"}

func TestFetchRejectsBadInput(t *testing.T) {
	f := NewFetcher(testHosts, 1<<20)

	if _, err := f.Fetch(context.Background(), "http://github.com/u/r", ""); !errors.Is(err, domain.ErrInvalid) {
		t.Errorf("non-https url: err = %v, want ErrInvalid", err)
	}
	if _, err := f.Fetch(context.Background(), "https://evil.io/u/r", ""); !errors.Is(err, domain.ErrInvalid) {
		t.Errorf("disallowed host: err = %v, want ErrInvalid", err)
	}
	if _, err := f.Fetch(context.Background(), "https://github.com/u/r", "bad branch"); !errors.Is(err, domain.ErrInvalid) {
		t.Errorf("bad branch: err = %v, want ErrInvalid", err)
	}
}

func TestMapCloneErr(t *testing.T) {
	bg := context.Background()
	canceled, cancel := context.WithCancel(bg)
	cancel()

	cases := []struct {
		name string
		ctx  context.Context
		err  error
		want string
	}{
		{"timeout", canceled, errors.New("whatever"), "timed out"},
		{"auth required", bg, transport.ErrAuthenticationRequired, "private"},
		{"auth failed", bg, transport.ErrAuthorizationFailed, "private"},
		{"not found", bg, transport.ErrRepositoryNotFound, "repository not found"},
		{"ref not found", bg, plumbing.ErrReferenceNotFound, `branch "b" not found`},
		{"no matching refspec", bg, gogit.NoMatchingRefSpecError{}, `branch "b" not found`},
		{"other", bg, errors.New("dial tcp: boom"), "clone failed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mapCloneErr(tc.ctx, tc.err, "b")
			if !strings.Contains(got.Error(), tc.want) {
				t.Errorf("mapCloneErr = %q, want it to contain %q", got, tc.want)
			}
		})
	}
}

func TestCommitSubject(t *testing.T) {
	long := strings.Repeat("x", 300)
	cases := []struct {
		name, in, want string
	}{
		{"first line only", "fix: water the plants\n\nlong body\n", "fix: water the plants"},
		{"trims whitespace", "  tidy \n", "tidy"},
		{"empty message", "", ""},
		{"caps at 200 runes", long + "\nbody", long[:200]},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := commitSubject(tc.in); got != tc.want {
				t.Errorf("commitSubject = %q, want %q", got, tc.want)
			}
		})
	}

	// Rune cap, not byte cap: multibyte subjects must not be split mid-rune.
	wide := strings.Repeat("é", 250)
	if got := commitSubject(wide); got != strings.Repeat("é", 200) {
		t.Errorf("multibyte subject capped to %d runes, want 200", len([]rune(got)))
	}
}

func TestCheckTreeSize(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.bin"), make([]byte, 600), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.bin"), make([]byte, 600), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := checkTreeSize(dir, 2000); err != nil {
		t.Errorf("under cap: %v, want nil", err)
	}
	if err := checkTreeSize(dir, 1000); err == nil || !strings.Contains(err.Error(), "size limit") {
		t.Errorf("over cap: %v, want size limit error", err)
	}
}
