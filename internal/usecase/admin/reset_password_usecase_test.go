package admin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/domain/user"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	iulid "github.com/lbrty/observer/internal/ulid"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func TestResetPassword_InvalidatesSessions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	mockSessionRepo := mock_repo.NewMockSessionRepository(ctrl)
	hasher := crypto.NewArgonHasher()

	uc := ucadmin.NewUserUseCase(mockUserRepo, mockCredRepo, hasher, mockSessionRepo, nil)

	ctx := context.Background()
	uid := iulid.New()
	hash, salt, err := hasher.Hash("oldPassword1!")
	require.NoError(t, err)

	cred := &user.Credentials{UserID: uid, PasswordHash: hash, Salt: salt}

	mockCredRepo.EXPECT().GetByUserID(ctx, uid).Return(cred, nil)
	mockCredRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
	mockSessionRepo.EXPECT().DeleteByUserID(ctx, uid).Return(nil)

	err = uc.ResetPassword(ctx, uid, ucadmin.ResetPasswordInput{NewPassword: "newPassword1!"})
	require.NoError(t, err)
}
