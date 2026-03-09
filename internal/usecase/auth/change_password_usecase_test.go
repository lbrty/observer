package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/user"
	iulid "github.com/lbrty/observer/internal/ulid"
	ucauth "github.com/lbrty/observer/internal/usecase/auth"
)

func TestChangePassword_InvalidatesSessions(t *testing.T) {
	uc, _, mockCredRepo, mockSessionRepo, _, hasher := setupAuthUseCase(t)

	ctx := context.Background()
	uid := iulid.New()
	password := "oldPassword1!"
	hash, salt, err := hasher.Hash(password)
	require.NoError(t, err)

	cred := &user.Credentials{UserID: uid, PasswordHash: hash, Salt: salt}

	mockCredRepo.EXPECT().GetByUserID(ctx, uid).Return(cred, nil)
	mockCredRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
	mockSessionRepo.EXPECT().DeleteByUserID(ctx, uid).Return(nil)

	err = uc.ChangePassword(ctx, uid, ucauth.ChangePasswordInput{
		CurrentPassword: password,
		NewPassword:     "newPassword1!",
	})
	require.NoError(t, err)
}
