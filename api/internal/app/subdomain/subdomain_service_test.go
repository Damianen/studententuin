package subdomain

import (
	"api/internal/app/ports"
	"api/internal/domain"
	"api/internal/mocks"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestCreateSubdomain_Execute(t *testing.T) {
	validUUID := uuid.New().String()

	tests := []struct {
		name    string
		input   SubdomainInput
		setup   func(sr *mocks.MockSubdomainRepo)
		wantErr bool
	}{
		{
			name:  "success",
			input: SubdomainInput{Name: "mysite", FullDomain: "mysite.example.com", UserID: validUUID},
			setup: func(sr *mocks.MockSubdomainRepo) {
				sr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:    "invalid UUID",
			input:   SubdomainInput{Name: "mysite", FullDomain: "mysite.example.com", UserID: "not-a-uuid"},
			setup:   func(sr *mocks.MockSubdomainRepo) {},
			wantErr: true,
		},
		{
			name:  "repo error",
			input: SubdomainInput{Name: "mysite", FullDomain: "mysite.example.com", UserID: validUUID},
			setup: func(sr *mocks.MockSubdomainRepo) {
				sr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sr := mocks.NewMockSubdomainRepo(ctrl)

			tt.setup(sr)

			svc := NewService(Dependencies{SubdomainRepo: sr})
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

func TestGetSubdomain_Execute(t *testing.T) {
	subID := uuid.New()

	tests := []struct {
		name    string
		setup   func(sr *mocks.MockSubdomainRepo)
		wantErr string
	}{
		{
			name: "success",
			setup: func(sr *mocks.MockSubdomainRepo) {
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(&domain.Subdomain{ID: subID}, nil)
			},
		},
		{
			name: "not found",
			setup: func(sr *mocks.MockSubdomainRepo) {
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sr := mocks.NewMockSubdomainRepo(ctrl)

			tt.setup(sr)

			svc := NewService(Dependencies{SubdomainRepo: sr})
			result, err := svc.Get.Execute(context.Background(), subID.String())

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
			if result.ID != subID {
				t.Fatalf("expected subdomain ID %s, got %s", subID, result.ID)
			}
		})
	}
}

func TestUpdateSubdomain_Execute(t *testing.T) {
	domainStr := "new.example.com"
	nameStr := "newname"
	emptyStr := ""

	tests := []struct {
		name    string
		input   SubdomainUpdateInput
		setup   func(sr *mocks.MockSubdomainRepo, c *mocks.MockClock)
		wantErr string
	}{
		{
			name:  "success",
			input: SubdomainUpdateInput{ID: "123", FullDomain: &domainStr, Name: &nameStr},
			setup: func(sr *mocks.MockSubdomainRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				sr.EXPECT().Update("123", gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:  "empty domain error",
			input: SubdomainUpdateInput{ID: "123", FullDomain: &emptyStr},
			setup: func(sr *mocks.MockSubdomainRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "domain cannot be empty",
		},
		{
			name:  "empty name error",
			input: SubdomainUpdateInput{ID: "123", Name: &emptyStr},
			setup: func(sr *mocks.MockSubdomainRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "name cannot be empty",
		},
		{
			name:  "no fields error",
			input: SubdomainUpdateInput{ID: "123"},
			setup: func(sr *mocks.MockSubdomainRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "no fields to update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sr := mocks.NewMockSubdomainRepo(ctrl)
			c := mocks.NewMockClock(ctrl)

			tt.setup(sr, c)

			svc := NewService(Dependencies{SubdomainRepo: sr, Clock: c})
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

func TestDeleteSubdomain_Execute(t *testing.T) {
	subID := uuid.New()
	appID := uuid.New()
	dbID := uuid.New()

	// bare wires a subdomain with no application and no database.
	bare := func(ar *mocks.MockApplicationRepo, dr *mocks.MockDatabaseRepo) {
		ar.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
		dr.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound).AnyTimes()
	}

	tests := []struct {
		name    string
		setup   func(sr *mocks.MockSubdomainRepo, ar *mocks.MockApplicationRepo, dr *mocks.MockDatabaseRepo, sm *mocks.MockServerManagerClient)
		wantErr string
	}{
		{
			name: "success without children",
			setup: func(sr *mocks.MockSubdomainRepo, ar *mocks.MockApplicationRepo, dr *mocks.MockDatabaseRepo, sm *mocks.MockServerManagerClient) {
				sub := &domain.Subdomain{ID: subID}
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(sub, nil)
				bare(ar, dr)
				sr.EXPECT().Delete(sub, gomock.Any()).Return(nil)
			},
		},
		{
			name: "removes app and database containers before the row",
			setup: func(sr *mocks.MockSubdomainRepo, ar *mocks.MockApplicationRepo, dr *mocks.MockDatabaseRepo, sm *mocks.MockServerManagerClient) {
				sub := &domain.Subdomain{ID: subID}
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(sub, nil)
				ar.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(&domain.Application{ID: appID, SubdomainID: subID}, nil)
				sm.EXPECT().Remove(gomock.Any(), appID.String()).Return(nil)
				dr.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(&domain.Database{ID: dbID, SubdomainID: subID}, nil)
				sm.EXPECT().RemoveDatabase(gomock.Any(), dbID.String()).Return(nil)
				sr.EXPECT().Delete(sub, gomock.Any()).Return(nil)
			},
		},
		{
			name: "never-deployed app is not an error",
			setup: func(sr *mocks.MockSubdomainRepo, ar *mocks.MockApplicationRepo, dr *mocks.MockDatabaseRepo, sm *mocks.MockServerManagerClient) {
				sub := &domain.Subdomain{ID: subID}
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(sub, nil)
				ar.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(&domain.Application{ID: appID, SubdomainID: subID}, nil)
				sm.EXPECT().Remove(gomock.Any(), appID.String()).Return(fmt.Errorf("servermanager remove: %w", ports.ErrAppNotDeployed))
				dr.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
				sr.EXPECT().Delete(sub, gomock.Any()).Return(nil)
			},
		},
		{
			name: "database removal failure blocks the row delete",
			setup: func(sr *mocks.MockSubdomainRepo, ar *mocks.MockApplicationRepo, dr *mocks.MockDatabaseRepo, sm *mocks.MockServerManagerClient) {
				sub := &domain.Subdomain{ID: subID}
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(sub, nil)
				ar.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
				dr.EXPECT().FindBySubdomainID(subID.String(), gomock.Any()).Return(&domain.Database{ID: dbID, SubdomainID: subID}, nil)
				sm.EXPECT().RemoveDatabase(gomock.Any(), dbID.String()).Return(errors.New("sweep failed"))
			},
			wantErr: "sweep failed",
		},
		{
			name: "find error",
			setup: func(sr *mocks.MockSubdomainRepo, ar *mocks.MockApplicationRepo, dr *mocks.MockDatabaseRepo, sm *mocks.MockServerManagerClient) {
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(sr *mocks.MockSubdomainRepo, ar *mocks.MockApplicationRepo, dr *mocks.MockDatabaseRepo, sm *mocks.MockServerManagerClient) {
				sub := &domain.Subdomain{ID: subID}
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(sub, nil)
				bare(ar, dr)
				sr.EXPECT().Delete(sub, gomock.Any()).Return(errors.New("delete failed"))
			},
			wantErr: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sr := mocks.NewMockSubdomainRepo(ctrl)
			ar := mocks.NewMockApplicationRepo(ctrl)
			dr := mocks.NewMockDatabaseRepo(ctrl)
			sm := mocks.NewMockServerManagerClient(ctrl)

			tt.setup(sr, ar, dr, sm)

			svc := NewService(Dependencies{SubdomainRepo: sr, ApplicationRepo: ar, DatabaseRepo: dr, ServerManager: sm})
			err := svc.Delete.Execute(context.Background(), subID.String())

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

func TestCheckUser_Execute(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	subID := uuid.New()

	tests := []struct {
		name      string
		setup     func(sr *mocks.MockSubdomainRepo)
		wantOwner bool
		wantErr   string
	}{
		{
			name: "owner match",
			setup: func(sr *mocks.MockSubdomainRepo) {
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(&domain.Subdomain{ID: subID, UserID: userID}, nil)
			},
			wantOwner: true,
		},
		{
			name: "non-owner",
			setup: func(sr *mocks.MockSubdomainRepo) {
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(&domain.Subdomain{ID: subID, UserID: otherUserID}, nil)
			},
			wantOwner: false,
		},
		{
			name: "find error",
			setup: func(sr *mocks.MockSubdomainRepo) {
				sr.EXPECT().FindByID(subID.String(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sr := mocks.NewMockSubdomainRepo(ctrl)
			ur := mocks.NewMockUserRepo(ctrl)

			tt.setup(sr)

			svc := NewService(Dependencies{SubdomainRepo: sr, UserRepo: ur})
			isOwner, err := svc.CheckUser.Execute(context.Background(), userID.String(), subID.String())

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
			if isOwner != tt.wantOwner {
				t.Fatalf("expected isOwner=%v, got %v", tt.wantOwner, isOwner)
			}
		})
	}
}
