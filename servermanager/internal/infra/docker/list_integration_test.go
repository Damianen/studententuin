//go:build integration

package docker

import (
	"context"
	"testing"

	"servermanager/internal/domain"

	"github.com/google/uuid"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func TestIntegrationListManagedContainers(t *testing.T) {
	r := integrationRuntime(t)
	ctx := context.Background()

	// Running managed app container: must be listed.
	runningID := uuid.NewString()
	cleanupApp(t, r, runningID)
	cid, err := r.Create(ctx, integrationSpec(runningID))
	if err != nil {
		t.Fatalf("Create running: %v", err)
	}
	if err := r.Start(ctx, cid); err != nil {
		t.Fatalf("Start: %v", err)
	}

	// Created-but-stopped managed container: excluded (running only).
	stoppedID := uuid.NewString()
	cleanupApp(t, r, stoppedID)
	if _, err := r.Create(ctx, integrationSpec(stoppedID)); err != nil {
		t.Fatalf("Create stopped: %v", err)
	}

	// Name-filter lookalike (substring match, non-UUID remainder): excluded.
	decoyName := "stt-app-decoy-" + uuid.NewString()[:8]
	decoy, err := r.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name:   decoyName,
		Config: &container.Config{Image: fixtureImage, Cmd: []string{"sleep", "3600"}},
	})
	if err != nil {
		t.Fatalf("creating decoy: %v", err)
	}
	t.Cleanup(func() {
		_, _ = r.cli.ContainerRemove(context.Background(), decoy.ID,
			client.ContainerRemoveOptions{Force: true, RemoveVolumes: true})
	})
	if _, err := r.cli.ContainerStart(ctx, decoy.ID, client.ContainerStartOptions{}); err != nil {
		t.Fatalf("starting decoy: %v", err)
	}

	managed, err := r.ListManagedContainers(ctx)
	if err != nil {
		t.Fatalf("ListManagedContainers: %v", err)
	}

	byOwner := map[string]domain.ManagedContainer{}
	for _, m := range managed {
		byOwner[m.OwnerID] = m
		if m.Name == decoyName {
			t.Errorf("decoy %q listed as managed", decoyName)
		}
	}

	got, ok := byOwner[runningID]
	if !ok {
		t.Fatalf("running app %s missing from %+v", runningID, managed)
	}
	if got.ID != cid || got.Kind != domain.KindApp || got.Name != domain.AppContainerName(runningID) {
		t.Errorf("managed = %+v, want id %s kind app name %s", got, cid, domain.AppContainerName(runningID))
	}
	if _, ok := byOwner[stoppedID]; ok {
		t.Errorf("stopped container %s listed, want running only", stoppedID)
	}
}
