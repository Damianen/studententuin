package user

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

func TestCreateUser_Execute(t *testing.T) {
	tests := []struct {
		name    string
		input   UserInput
		setup   func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, c *mocks.MockClock)
		wantErr string
	}{
		{
			name:  "success",
			input: UserInput{Email: "Test@Example.com", Password: "secret123", Name: "Alice"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				h.EXPECT().Hash("secret123").Return("hashed", nil)
				pr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				ur.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(u *domain.User, ctx context.Context) error {
					if u.Email != "test@example.com" {
						t.Errorf("expected lowercased/trimmed email, got %q", u.Email)
					}
					if u.DisplayName != "Alice" {
						t.Errorf("expected name Alice, got %q", u.DisplayName)
					}
					if u.Status != "active" {
						t.Errorf("expected status active, got %q", u.Status)
					}
					return nil
				})
			},
		},
		{
			name:    "empty email",
			input:   UserInput{Email: "  ", Password: "secret123"},
			setup:   func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "email and password required",
		},
		{
			name:    "empty password",
			input:   UserInput{Email: "a@b.com", Password: ""},
			setup:   func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "email and password required",
		},
		{
			name:  "hasher error",
			input: UserInput{Email: "a@b.com", Password: "secret"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				h.EXPECT().Hash("secret").Return("", errors.New("hash failed"))
			},
			wantErr: "hash failed",
		},
		{
			name:  "password repo error",
			input: UserInput{Email: "a@b.com", Password: "secret"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				h.EXPECT().Hash("secret").Return("hashed", nil)
				ur.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				pr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("pw repo error"))
			},
			wantErr: "pw repo error",
		},
		{
			name:  "user repo error (email conflict)",
			input: UserInput{Email: "a@b.com", Password: "secret"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				h.EXPECT().Hash("secret").Return("hashed", nil)
				ur.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("email already exists"))
			},
			wantErr: "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mocks.NewMockUserRepo(ctrl)
			pr := mocks.NewMockPasswordRepo(ctrl)
			h := mocks.NewMockPasswordHasher(ctrl)
			c := mocks.NewMockClock(ctrl)

			tt.setup(ur, pr, h, c)

			svc := NewService(Dependencies{
				UserRepo:     ur,
				PasswordRepo: pr,
				Clock:        c,
				Hasher:       h,
			})

			err := svc.Create.Execute(context.Background(), tt.input)
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

func TestGetUser_Execute(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		setup   func(ur *mocks.MockUserRepo)
		wantErr string
	}{
		{
			name: "success",
			setup: func(ur *mocks.MockUserRepo) {
				ur.EXPECT().FindByID(userID.String(), gomock.Any()).Return(&domain.User{ID: userID, Email: "a@b.com"}, nil)
			},
		},
		{
			name: "not found",
			setup: func(ur *mocks.MockUserRepo) {
				ur.EXPECT().FindByID(userID.String(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mocks.NewMockUserRepo(ctrl)

			tt.setup(ur)

			svc := NewService(Dependencies{UserRepo: ur})
			result, err := svc.Get.Execute(context.Background(), userID.String())

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
			if result.ID != userID {
				t.Fatalf("expected user ID %s, got %s", userID, result.ID)
			}
		})
	}
}

func TestUpdateUser_Execute(t *testing.T) {
	email := "new@example.com"
	name := "Bob"
	emptyStr := ""

	tests := []struct {
		name    string
		input   UserUpdateInput
		setup   func(ur *mocks.MockUserRepo, c *mocks.MockClock)
		wantErr string
	}{
		{
			name:  "success email only",
			input: UserUpdateInput{ID: "123", Email: &email},
			setup: func(ur *mocks.MockUserRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				ur.EXPECT().Update("123", gomock.Any(), gomock.Any()).DoAndReturn(func(id string, updates map[string]any, ctx context.Context) error {
					if updates["email"] != "new@example.com" {
						t.Errorf("expected email in updates, got %v", updates)
					}
					return nil
				})
			},
		},
		{
			name:  "success name only",
			input: UserUpdateInput{ID: "123", Name: &name},
			setup: func(ur *mocks.MockUserRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				ur.EXPECT().Update("123", gomock.Any(), gomock.Any()).DoAndReturn(func(id string, updates map[string]any, ctx context.Context) error {
					if updates["display_name"] != "Bob" {
						t.Errorf("expected display_name in updates, got %v", updates)
					}
					return nil
				})
			},
		},
		{
			name:  "success both fields",
			input: UserUpdateInput{ID: "123", Email: &email, Name: &name},
			setup: func(ur *mocks.MockUserRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				ur.EXPECT().Update("123", gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:    "empty email error",
			input:   UserUpdateInput{ID: "123", Email: &emptyStr},
			setup: func(ur *mocks.MockUserRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "email cannot be empty",
		},
		{
			name:    "empty name error",
			input:   UserUpdateInput{ID: "123", Name: &emptyStr},
			setup: func(ur *mocks.MockUserRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "name cannot be empty",
		},
		{
			name:    "no fields error",
			input:   UserUpdateInput{ID: "123"},
			setup: func(ur *mocks.MockUserRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			wantErr: "no fields to update",
		},
		{
			name:  "repo error",
			input: UserUpdateInput{ID: "123", Email: &email},
			setup: func(ur *mocks.MockUserRepo, c *mocks.MockClock) {
				c.EXPECT().Now().Return(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
				ur.EXPECT().Update("123", gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mocks.NewMockUserRepo(ctrl)
			c := mocks.NewMockClock(ctrl)

			tt.setup(ur, c)

			svc := NewService(Dependencies{UserRepo: ur, Clock: c})
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

func TestDeleteUser_Execute(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		setup   func(ur *mocks.MockUserRepo)
		wantErr string
	}{
		{
			name: "success",
			setup: func(ur *mocks.MockUserRepo) {
				user := &domain.User{ID: userID}
				ur.EXPECT().FindByID(userID.String(), gomock.Any()).Return(user, nil)
				ur.EXPECT().Delete(user, gomock.Any()).Return(nil)
			},
		},
		{
			name: "find error",
			setup: func(ur *mocks.MockUserRepo) {
				ur.EXPECT().FindByID(userID.String(), gomock.Any()).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(ur *mocks.MockUserRepo) {
				user := &domain.User{ID: userID}
				ur.EXPECT().FindByID(userID.String(), gomock.Any()).Return(user, nil)
				ur.EXPECT().Delete(user, gomock.Any()).Return(errors.New("delete failed"))
			},
			wantErr: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mocks.NewMockUserRepo(ctrl)

			tt.setup(ur)

			svc := NewService(Dependencies{UserRepo: ur})
			err := svc.Delete.Execute(context.Background(), userID.String())

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
