package auth

import (
	"api/internal/domain"
	"api/internal/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestLoginAuth_Execute(t *testing.T) {
	userID := uuid.New()
	user := &domain.User{ID: userID, Email: "test@example.com"}
	cred := &domain.PasswordCredential{UserId: userID, PasswordHash: "hashed"}

	tests := []struct {
		name    string
		input   LoginInput
		setup   func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, jwt *mocks.MockJwtTokenizer)
		want    string
		wantErr string
	}{
		{
			name:  "success",
			input: LoginInput{Email: "Test@Example.com", Password: "secret"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, jwt *mocks.MockJwtTokenizer) {
				ur.EXPECT().FindByEmail("test@example.com", gomock.Any()).Return(user, nil)
				pr.EXPECT().FindById(userID.String(), gomock.Any()).Return(cred, nil)
				h.EXPECT().Compare("secret", "hashed").Return(true, nil)
				jwt.EXPECT().CreateToken(userID.String()).Return("jwt-token", nil)
			},
			want: "jwt-token",
		},
		{
			name:  "user not found",
			input: LoginInput{Email: "nope@example.com", Password: "secret"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, jwt *mocks.MockJwtTokenizer) {
				ur.EXPECT().FindByEmail("nope@example.com", gomock.Any()).Return(nil, errors.New("user not found"))
			},
			wantErr: "user not found",
		},
		{
			name:  "password cred not found",
			input: LoginInput{Email: "test@example.com", Password: "secret"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, jwt *mocks.MockJwtTokenizer) {
				ur.EXPECT().FindByEmail("test@example.com", gomock.Any()).Return(user, nil)
				pr.EXPECT().FindById(userID.String(), gomock.Any()).Return(nil, errors.New("credential not found"))
			},
			wantErr: "credential not found",
		},
		{
			name:  "wrong password",
			input: LoginInput{Email: "test@example.com", Password: "wrong"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, jwt *mocks.MockJwtTokenizer) {
				ur.EXPECT().FindByEmail("test@example.com", gomock.Any()).Return(user, nil)
				pr.EXPECT().FindById(userID.String(), gomock.Any()).Return(cred, nil)
				h.EXPECT().Compare("wrong", "hashed").Return(false, nil)
			},
			wantErr: "invalid credentials",
		},
		{
			name:  "hasher error",
			input: LoginInput{Email: "test@example.com", Password: "secret"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, jwt *mocks.MockJwtTokenizer) {
				ur.EXPECT().FindByEmail("test@example.com", gomock.Any()).Return(user, nil)
				pr.EXPECT().FindById(userID.String(), gomock.Any()).Return(cred, nil)
				h.EXPECT().Compare("secret", "hashed").Return(false, errors.New("compare failed"))
			},
			wantErr: "compare failed",
		},
		{
			name:  "token creation error",
			input: LoginInput{Email: "test@example.com", Password: "secret"},
			setup: func(ur *mocks.MockUserRepo, pr *mocks.MockPasswordRepo, h *mocks.MockPasswordHasher, jwt *mocks.MockJwtTokenizer) {
				ur.EXPECT().FindByEmail("test@example.com", gomock.Any()).Return(user, nil)
				pr.EXPECT().FindById(userID.String(), gomock.Any()).Return(cred, nil)
				h.EXPECT().Compare("secret", "hashed").Return(true, nil)
				jwt.EXPECT().CreateToken(userID.String()).Return("", errors.New("token error"))
			},
			wantErr: "token error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mocks.NewMockUserRepo(ctrl)
			pr := mocks.NewMockPasswordRepo(ctrl)
			h := mocks.NewMockPasswordHasher(ctrl)
			jwt := mocks.NewMockJwtTokenizer(ctrl)
			c := mocks.NewMockClock(ctrl)

			tt.setup(ur, pr, h, jwt)

			svc := NewService(Dependencies{
				UserRepo:     ur,
				PasswordRepo: pr,
				Clock:        c,
				Hasher:       h,
				JwtTokenizer: jwt,
			})

			token, err := svc.Login.Execute(context.Background(), tt.input)
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
			if token != tt.want {
				t.Fatalf("expected token %q, got %q", tt.want, token)
			}
		})
	}
}
