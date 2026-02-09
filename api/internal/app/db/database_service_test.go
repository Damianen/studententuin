package db

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

func TestCreateDatabase_Execute(t *testing.T) {
	validUUID := uuid.New().String()

	tests := []struct {
		name    string
		input   DatabaseInput
		setup   func(dr *mocks.MockDatabaseRepo)
		wantErr bool
	}{
		{
			name: "success",
			input: DatabaseInput{
				SubdomainID: validUUID,
				Name:        "mydb",
				Version:     "15",
				Type:        domain.DatabaseTypePostgres,
				Status:      domain.DatabaseStatusProvisioning,
			},
			setup: func(dr *mocks.MockDatabaseRepo) {
				dr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "invalid UUID",
			input: DatabaseInput{
				SubdomainID: "not-a-uuid",
				Name:        "mydb",
			},
			setup:   func(dr *mocks.MockDatabaseRepo) {},
			wantErr: true,
		},
		{
			name: "repo error",
			input: DatabaseInput{
				SubdomainID: validUUID,
				Name:        "mydb",
				Version:     "15",
				Type:        domain.DatabaseTypePostgres,
				Status:      domain.DatabaseStatusProvisioning,
			},
			setup: func(dr *mocks.MockDatabaseRepo) {
				dr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			dr := mocks.NewMockDatabaseRepo(ctrl)

			tt.setup(dr)

			svc := NewService(Dependencies{DatabaseRepo: dr})
			// CreateDatabase also has subdomainRepo field but NewService doesn't set it,
			// so we set databaseRepo directly via the service
			svc.Create.databaseRepo = dr

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

func TestGetDatabase_Execute(t *testing.T) {
	dbID := uuid.New()

	tests := []struct {
		name    string
		setup   func(dr *mocks.MockDatabaseRepo)
		wantErr string
	}{
		{
			name: "success",
			setup: func(dr *mocks.MockDatabaseRepo) {
				dr.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(&domain.Database{ID: dbID}, nil)
			},
		},
		{
			name: "not found",
			setup: func(dr *mocks.MockDatabaseRepo) {
				dr.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			dr := mocks.NewMockDatabaseRepo(ctrl)

			tt.setup(dr)

			svc := NewService(Dependencies{DatabaseRepo: dr})
			result, err := svc.Get.Execute(context.Background(), dbID.String())

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
			if result.ID != dbID {
				t.Fatalf("expected ID %s, got %s", dbID, result.ID)
			}
		})
	}
}

func TestUpdateDatabase_Execute(t *testing.T) {
	nameStr := "newdb"
	emptyStr := ""
	versionStr := "16"
	dbType := domain.DatabaseTypeMySQL
	dbStatus := domain.DatabaseStatusRunning

	tests := []struct {
		name    string
		input   DatabaseUpdateInput
		setup   func(dr *mocks.MockDatabaseRepo, c *mocks.MockClock)
		wantErr string
	}{
		{
			name: "success with name",
			input: DatabaseUpdateInput{ID: "123", Name: &nameStr},
			setup: func(dr *mocks.MockDatabaseRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				dr.EXPECT().Update("123", gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "success with multiple fields",
			input: DatabaseUpdateInput{ID: "123", Name: &nameStr, Version: &versionStr, Type: &dbType, Status: &dbStatus},
			setup: func(dr *mocks.MockDatabaseRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				dr.EXPECT().Update("123", gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:  "empty name error",
			input: DatabaseUpdateInput{ID: "123", Name: &emptyStr},
			setup: func(dr *mocks.MockDatabaseRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "name cannot be empty",
		},
		{
			name:  "no fields error",
			input: DatabaseUpdateInput{ID: "123"},
			setup: func(dr *mocks.MockDatabaseRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "no fields to update!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			dr := mocks.NewMockDatabaseRepo(ctrl)
			c := mocks.NewMockClock(ctrl)

			tt.setup(dr, c)

			svc := NewService(Dependencies{DatabaseRepo: dr, Clock: c})
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

func TestDeleteDatabase_Execute(t *testing.T) {
	dbID := uuid.New()

	tests := []struct {
		name    string
		setup   func(dr *mocks.MockDatabaseRepo)
		wantErr string
	}{
		{
			name: "success",
			setup: func(dr *mocks.MockDatabaseRepo) {
				db := &domain.Database{ID: dbID}
				dr.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(db, nil)
				dr.EXPECT().Delete(db, gomock.Any()).Return(nil)
			},
		},
		{
			name: "find error",
			setup: func(dr *mocks.MockDatabaseRepo) {
				dr.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(dr *mocks.MockDatabaseRepo) {
				db := &domain.Database{ID: dbID}
				dr.EXPECT().FindByID(dbID.String(), gomock.Any()).Return(db, nil)
				dr.EXPECT().Delete(db, gomock.Any()).Return(errors.New("delete failed"))
			},
			wantErr: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			dr := mocks.NewMockDatabaseRepo(ctrl)

			tt.setup(dr)

			svc := NewService(Dependencies{DatabaseRepo: dr})
			err := svc.Delete.Execute(context.Background(), dbID.String())

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
