package app

import (
	"api/internal/domain"
	"api/internal/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
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
