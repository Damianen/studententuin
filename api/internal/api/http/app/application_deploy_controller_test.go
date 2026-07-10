package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"api/internal/app/ports"
	"api/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func deployReadyApp(id, subID uuid.UUID) *domain.Application {
	repo := "https://github.com/user/repo"
	return &domain.Application{
		ID:            id,
		SubdomainID:   subID,
		RepositoryURL: &repo,
		Type:          domain.ApplicationTypeNodejs,
	}
}

func TestDeployEndpoint(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	appID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	path := "/subdomain/" + subID.String() + "/application/" + appID.String() + "/deploy"

	t.Run("queues and returns the deployment id", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployReadyApp(appID, subID), nil)
		m.serverManager.EXPECT().Deploy(gomock.Any(), appID.String(), gomock.Any()).Return("dep-123", nil)
		m.clock.EXPECT().Now().Return(time.Now()).AnyTimes()
		m.applicationRepo.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).Return(nil)
		// The spawned poller keeps polling until terminal; park then finish.
		m.serverManager.EXPECT().DeploymentStatus(gomock.Any(), "dep-123").
			Return(&ports.DeploymentStatus{Status: "failed"}, nil).AnyTimes()
		m.applicationRepo.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		w := doJSON(r, http.MethodPost, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusAccepted {
			t.Fatalf("expected 202, got %d (body %s)", w.Code, w.Body.String())
		}
		var env envelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatal(err)
		}
		var resp struct {
			DeploymentID string `json:"deployment_id"`
		}
		if err := json.Unmarshal(env.Data, &resp); err != nil || resp.DeploymentID != "dep-123" {
			t.Errorf("data = %s, err = %v", env.Data, err)
		}
	})

	t.Run("no repository url is a 400", func(t *testing.T) {
		r, m := newAppRouter(t)
		application := deployReadyApp(appID, subID)
		application.RepositoryURL = nil
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(application, nil)

		if w := doJSON(r, http.MethodPost, path, "", authCookie(t, ownerID.String())); w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("deploy in flight is a 409", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployReadyApp(appID, subID), nil)
		m.serverManager.EXPECT().Deploy(gomock.Any(), appID.String(), gomock.Any()).
			Return("", fmt.Errorf("servermanager deploy: %w", ports.ErrDeployInFlight))

		if w := doJSON(r, http.MethodPost, path, "", authCookie(t, ownerID.String())); w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("manager down is a 502", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployReadyApp(appID, subID), nil)
		m.serverManager.EXPECT().Deploy(gomock.Any(), appID.String(), gomock.Any()).
			Return("", fmt.Errorf("servermanager deploy: connection refused"))

		if w := doJSON(r, http.MethodPost, path, "", authCookie(t, ownerID.String())); w.Code != http.StatusBadGateway {
			t.Fatalf("expected 502, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("foreign subdomain is a 404", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployReadyApp(appID, uuid.New()), nil)

		if w := doJSON(r, http.MethodPost, path, "", authCookie(t, ownerID.String())); w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d (body %s)", w.Code, w.Body.String())
		}
	})
}

func TestGetDeploymentEndpoint(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	appID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	path := "/subdomain/" + subID.String() + "/application/" + appID.String() + "/deployment/dep-123"

	t.Run("returns the job status", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployReadyApp(appID, subID), nil)
		m.serverManager.EXPECT().DeploymentStatus(gomock.Any(), "dep-123").Return(&ports.DeploymentStatus{
			ID: "dep-123", AppID: appID.String(), Status: "failed",
			Error: "build failed: exit 1", BuildLog: "npm ERR! boom\n",
			CreatedAt: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2026, 7, 10, 12, 1, 0, 0, time.UTC),
		}, nil)

		w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String()))
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
		var env envelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatal(err)
		}
		var resp struct {
			Status   string `json:"status"`
			Error    string `json:"error"`
			BuildLog string `json:"build_log"`
		}
		if err := json.Unmarshal(env.Data, &resp); err != nil ||
			resp.Status != "failed" || resp.Error != "build failed: exit 1" || resp.BuildLog == "" {
			t.Errorf("data = %s, err = %v", env.Data, err)
		}
	})

	t.Run("unknown deployment is a 404", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployReadyApp(appID, subID), nil)
		m.serverManager.EXPECT().DeploymentStatus(gomock.Any(), "dep-123").
			Return(nil, fmt.Errorf("status: %w", ports.ErrDeploymentNotFound))

		if w := doJSON(r, http.MethodGet, path, "", authCookie(t, ownerID.String())); w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d (body %s)", w.Code, w.Body.String())
		}
	})
}

func TestStartStopEndpoints(t *testing.T) {
	ownerID := uuid.New()
	subID := uuid.New()
	appID := uuid.New()
	owned := &domain.Subdomain{ID: subID, UserID: ownerID}
	base := "/subdomain/" + subID.String() + "/application/" + appID.String()

	t.Run("stop stops the container", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployReadyApp(appID, subID), nil)
		m.serverManager.EXPECT().Stop(gomock.Any(), appID.String()).Return(nil)
		m.clock.EXPECT().Now().Return(time.Now())
		m.applicationRepo.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).Return(nil)

		if w := doJSON(r, http.MethodPost, base+"/stop", "", authCookie(t, ownerID.String())); w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("start starts the container", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployReadyApp(appID, subID), nil)
		m.serverManager.EXPECT().Start(gomock.Any(), appID.String()).Return(nil)
		m.clock.EXPECT().Now().Return(time.Now())
		m.applicationRepo.EXPECT().Update(appID.String(), gomock.Any(), gomock.Any()).Return(nil)

		if w := doJSON(r, http.MethodPost, base+"/start", "", authCookie(t, ownerID.String())); w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("never deployed is a 400", func(t *testing.T) {
		r, m := newAppRouter(t)
		m.subdomainRepo.EXPECT().FindByID(subID.String(), gomock.Any()).Return(owned, nil)
		m.applicationRepo.EXPECT().FindByID(appID.String(), gomock.Any()).Return(deployReadyApp(appID, subID), nil)
		m.serverManager.EXPECT().Stop(gomock.Any(), appID.String()).
			Return(fmt.Errorf("servermanager stop: %w", ports.ErrAppNotDeployed))

		if w := doJSON(r, http.MethodPost, base+"/stop", "", authCookie(t, ownerID.String())); w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d (body %s)", w.Code, w.Body.String())
		}
	})
}
