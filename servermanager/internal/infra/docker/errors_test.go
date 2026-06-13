package docker

import (
	"errors"
	"fmt"
	"testing"

	"servermanager/internal/domain"

	cerrdefs "github.com/containerd/errdefs"
)

func TestMapErr(t *testing.T) {
	cases := []struct {
		name string
		in   error
		want error
	}{
		{"nil", nil, nil},
		{"not found", fmt.Errorf("wrapped: %w", cerrdefs.ErrNotFound), domain.ErrNotFound},
		{"conflict", fmt.Errorf("wrapped: %w", cerrdefs.ErrConflict), domain.ErrConflict},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mapErr("op", tc.in)
			if tc.want == nil {
				if got != nil {
					t.Fatalf("mapErr = %v, want nil", got)
				}
				return
			}
			if !errors.Is(got, tc.want) {
				t.Errorf("mapErr(%v) = %v, want errors.Is %v", tc.in, got, tc.want)
			}
		})
	}

	plain := errors.New("daemon exploded")
	got := mapErr("op", plain)
	if !errors.Is(got, plain) {
		t.Errorf("mapErr should wrap unknown errors, got %v", got)
	}
	if errors.Is(got, domain.ErrNotFound) || errors.Is(got, domain.ErrConflict) {
		t.Errorf("mapErr mapped unknown error to a sentinel: %v", got)
	}
}
