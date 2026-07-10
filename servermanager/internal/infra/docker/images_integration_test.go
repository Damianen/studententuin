//go:build integration

package docker

import (
	"context"
	"errors"
	"testing"

	"servermanager/internal/domain"

	"github.com/google/uuid"
	"github.com/moby/moby/client"
)

func TestIntegrationImageStore(t *testing.T) {
	r := integrationRuntime(t) // pulls the busybox fixture
	ctx := context.Background()
	appID := uuid.NewString()

	tagA := domain.AppImageRef(appID, uuid.NewString())
	tagB := domain.AppImageRef(appID, uuid.NewString())
	for _, tag := range []string{tagA, tagB} {
		if _, err := r.cli.ImageTag(ctx, client.ImageTagOptions{Source: fixtureImage, Target: tag}); err != nil {
			t.Fatalf("tagging fixture as %s: %v", tag, err)
		}
	}
	t.Cleanup(func() {
		for _, tag := range []string{tagA, tagB} {
			_, _ = r.cli.ImageRemove(ctx, tag, client.ImageRemoveOptions{})
		}
	})

	refs, err := r.ListAppImages(ctx, appID)
	if err != nil {
		t.Fatalf("ListAppImages: %v", err)
	}
	got := map[string]bool{}
	for _, ref := range refs {
		got[ref] = true
	}
	if len(refs) != 2 || !got[tagA] || !got[tagB] {
		t.Fatalf("ListAppImages = %v, want [%s %s]", refs, tagA, tagB)
	}

	if err := r.RemoveImage(ctx, tagA); err != nil {
		t.Fatalf("RemoveImage: %v", err)
	}
	refs, err = r.ListAppImages(ctx, appID)
	if err != nil {
		t.Fatalf("ListAppImages after remove: %v", err)
	}
	if len(refs) != 1 || refs[0] != tagB {
		t.Errorf("ListAppImages after remove = %v, want [%s]", refs, tagB)
	}

	if err := r.RemoveImage(ctx, tagA); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("RemoveImage(missing) = %v, want ErrNotFound", err)
	}

	// Unknown app lists empty, not an error.
	if refs, err := r.ListAppImages(ctx, uuid.NewString()); err != nil || len(refs) != 0 {
		t.Errorf("ListAppImages(unknown) = %v, %v; want empty, nil", refs, err)
	}
}
