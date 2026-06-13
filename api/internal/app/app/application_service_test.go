package app

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"api/internal/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestCreateApplication_Execute(t *testing.T) {
	validUUID := uuid.New().String()

	tests := []struct {
		name    string
		input   ApplicationInput
		setup   func(ar *mocks.MockApplicationRepo)
		wantErr bool
	}{
		{
			name: "success",
			input: ApplicationInput{
				SubdomainID: validUUID,
				Name:        "myapp",
				Type:        domain.ApplicationTypeNodejs,
				Status:      domain.ApplicationStatusPending,
			},
			setup: func(ar *mocks.MockApplicationRepo) {
				ar.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "invalid UUID",
			input: ApplicationInput{
				SubdomainID: "not-a-uuid",
				Name:        "myapp",
			},
			setup:   func(ar *mocks.MockApplicationRepo) {},
			wantErr: true,
		},
		{
			name: "repo error",
			input: ApplicationInput{
				SubdomainID: validUUID,
				Name:        "myapp",
				Type:        domain.ApplicationTypeNodejs,
				Status:      domain.ApplicationStatusPending,
			},
			setup: func(ar *mocks.MockApplicationRepo) {
				ar.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ar := mocks.NewMockApplicationRepo(ctrl)

			tt.setup(ar)

			svc := NewService(Dependencies{ApplicationRepo: ar})
			svc.Create.applicationRepo = ar

			err := svc.Create.Execute(context.Background(), tt.input)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetApplication_Execute(t *testing.T) {
	appID := uuid.New()

	tests := []struct {
		name    string
		setup   func(ar *mocks.MockApplicationRepo)
		wantErr string
	}{
		{
			name: "success",
			setup: func(ar *mocks.MockApplicationRepo) {
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(&domain.Application{ID: appID}, nil)
			},
		},
		{
			name: "not found",
			setup: func(ar *mocks.MockApplicationRepo) {
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ar := mocks.NewMockApplicationRepo(ctrl)

			tt.setup(ar)

			svc := NewService(Dependencies{ApplicationRepo: ar})
			result, err := svc.Get.Execute(context.Background(), appID.String())

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.ID != appID {
				t.Fatalf("expected ID %s, got %s", appID, result.ID)
			}
		})
	}
}

func TestUpdateApplication_Execute(t *testing.T) {
	nameStr := "newapp"
	emptyStr := ""
	repoURL := "https://github.com/example/repo"
	appType := domain.ApplicationTypeNodejs
	appStatus := domain.ApplicationStatusRunning

	tests := []struct {
		name    string
		input   ApplicationUpdateInput
		setup   func(ar *mocks.MockApplicationRepo, c *mocks.MockClock)
		wantErr string
	}{
		{
			name: "success with name",
			input: ApplicationUpdateInput{ID: "123", Name: &nameStr},
			setup: func(ar *mocks.MockApplicationRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				ar.EXPECT().Update("123", gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "success with multiple fields",
			input: ApplicationUpdateInput{ID: "123", Name: &nameStr, RepoUrl: &repoURL, Type: &appType, Status: &appStatus},
			setup: func(ar *mocks.MockApplicationRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				ar.EXPECT().Update("123", gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:  "empty name error",
			input: ApplicationUpdateInput{ID: "123", Name: &emptyStr},
			setup: func(ar *mocks.MockApplicationRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "name cannot be empty",
		},
		{
			name:  "no fields error",
			input: ApplicationUpdateInput{ID: "123"},
			setup: func(ar *mocks.MockApplicationRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "no fields to update!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ar := mocks.NewMockApplicationRepo(ctrl)
			c := mocks.NewMockClock(ctrl)

			tt.setup(ar, c)

			svc := NewService(Dependencies{ApplicationRepo: ar, Clock: c})
			err := svc.Update.Execute(context.Background(), tt.input)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetApplicationLogs_Execute(t *testing.T) {
	appID := uuid.New()
	subdomainID := uuid.New()
	opts := ports.LogOptions{Tail: 200}
	ownedApp := &domain.Application{ID: appID, SubdomainID: subdomainID}

	tests := []struct {
		name     string
		setup    func(ar *mocks.MockApplicationRepo, sm *mocks.MockServerManagerClient)
		wantErr  error
		wantLogs int
	}{
		{
			name: "success",
			setup: func(ar *mocks.MockApplicationRepo, sm *mocks.MockServerManagerClient) {
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
				sm.EXPECT().Logs(gomock.Any(), appID.String(), opts).Return([]ports.LogEntry{
					{ID: "1", Level: "info", Message: "tick"},
				}, nil)
			},
			wantLogs: 1,
		},
		{
			// The manager is never called for an app in someone else's
			// subdomain — no sm.EXPECT() proves it.
			name: "app in different subdomain",
			setup: func(ar *mocks.MockApplicationRepo, sm *mocks.MockServerManagerClient) {
				foreign := &domain.Application{ID: appID, SubdomainID: uuid.New()}
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(foreign, nil)
			},
			wantErr: ErrNotInSubdomain,
		},
		{
			name: "app not found",
			setup: func(ar *mocks.MockApplicationRepo, sm *mocks.MockServerManagerClient) {
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name: "never deployed means empty logs, not an error",
			setup: func(ar *mocks.MockApplicationRepo, sm *mocks.MockServerManagerClient) {
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
				sm.EXPECT().Logs(gomock.Any(), appID.String(), opts).Return(nil, ports.ErrAppNotDeployed)
			},
			wantLogs: 0,
		},
		{
			name: "manager failure propagates",
			setup: func(ar *mocks.MockApplicationRepo, sm *mocks.MockServerManagerClient) {
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
				sm.EXPECT().Logs(gomock.Any(), appID.String(), opts).Return(nil, errors.New("connection refused"))
			},
			wantErr: errors.New("connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ar := mocks.NewMockApplicationRepo(ctrl)
			sm := mocks.NewMockServerManagerClient(ctrl)

			tt.setup(ar, sm)

			svc := NewService(Dependencies{ApplicationRepo: ar, ServerManager: sm})
			logs, err := svc.Logs.Execute(context.Background(), subdomainID.String(), appID.String(), opts)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if logs == nil {
				t.Fatal("logs is nil, want a (possibly empty) slice")
			}
			if len(logs) != tt.wantLogs {
				t.Fatalf("len(logs) = %d, want %d", len(logs), tt.wantLogs)
			}
		})
	}
}

func TestStreamApplicationLogs_Execute(t *testing.T) {
	appID := uuid.New()
	subdomainID := uuid.New()
	opts := ports.LogOptions{Tail: 200}
	ownedApp := &domain.Application{ID: appID, SubdomainID: subdomainID}

	t.Run("stream passes through", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ar := mocks.NewMockApplicationRepo(ctrl)
		sm := mocks.NewMockServerManagerClient(ctrl)
		ch := make(chan ports.LogEntry, 1)
		ch <- ports.LogEntry{ID: "1", Level: "info", Message: "live"}
		close(ch)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
		sm.EXPECT().StreamLogs(gomock.Any(), appID.String(), opts).Return((<-chan ports.LogEntry)(ch), nil)

		svc := NewService(Dependencies{ApplicationRepo: ar, ServerManager: sm})
		got, err := svc.StreamLogs.Execute(context.Background(), subdomainID.String(), appID.String(), opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if entry := <-got; entry.Message != "live" {
			t.Errorf("entry = %+v", entry)
		}
	})

	t.Run("foreign subdomain never reaches the manager", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ar := mocks.NewMockApplicationRepo(ctrl)
		sm := mocks.NewMockServerManagerClient(ctrl)
		foreign := &domain.Application{ID: appID, SubdomainID: uuid.New()}
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(foreign, nil)

		svc := NewService(Dependencies{ApplicationRepo: ar, ServerManager: sm})
		if _, err := svc.StreamLogs.Execute(context.Background(), subdomainID.String(), appID.String(), opts); !errors.Is(err, ErrNotInSubdomain) {
			t.Errorf("err = %v, want ErrNotInSubdomain", err)
		}
	})

	t.Run("not deployed propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ar := mocks.NewMockApplicationRepo(ctrl)
		sm := mocks.NewMockServerManagerClient(ctrl)
		ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(ownedApp, nil)
		sm.EXPECT().StreamLogs(gomock.Any(), appID.String(), opts).Return(nil, ports.ErrAppNotDeployed)

		svc := NewService(Dependencies{ApplicationRepo: ar, ServerManager: sm})
		if _, err := svc.StreamLogs.Execute(context.Background(), subdomainID.String(), appID.String(), opts); !errors.Is(err, ports.ErrAppNotDeployed) {
			t.Errorf("err = %v, want ErrAppNotDeployed", err)
		}
	})
}

func TestDeleteApplication_Execute(t *testing.T) {
	appID := uuid.New()

	tests := []struct {
		name    string
		setup   func(ar *mocks.MockApplicationRepo)
		wantErr string
	}{
		{
			name: "success",
			setup: func(ar *mocks.MockApplicationRepo) {
				a := &domain.Application{ID: appID}
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(a, nil)
				ar.EXPECT().Delete(a, gomock.Any()).Return(nil)
			},
		},
		{
			name: "find error",
			setup: func(ar *mocks.MockApplicationRepo) {
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(ar *mocks.MockApplicationRepo) {
				a := &domain.Application{ID: appID}
				ar.EXPECT().FindByID(appID.String(), gomock.Any()).Return(a, nil)
				ar.EXPECT().Delete(a, gomock.Any()).Return(errors.New("delete failed"))
			},
			wantErr: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ar := mocks.NewMockApplicationRepo(ctrl)

			tt.setup(ar)

			svc := NewService(Dependencies{ApplicationRepo: ar})
			err := svc.Delete.Execute(context.Background(), appID.String())

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
