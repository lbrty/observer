package admin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/user"
	iulid "github.com/lbrty/observer/internal/ulid"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func TestResetPassword_InvalidatesSessions(t *testing.T) {
	uc, d := newUserUCDeps(t)
	ctx := context.Background()
	uid := iulid.New()

	d.credRepo.EXPECT().GetByUserID(ctx, uid).Return(&user.Credentials{
		UserID:       uid,
		PasswordHash: "oldhash",
		Salt:         "oldsalt",
	}, nil)
	d.hasher.EXPECT().Hash("newPassword1!").Return("newhash", "newsalt", nil)
	d.credRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
	d.sessionRepo.EXPECT().DeleteByUserID(ctx, uid).Return(nil)

	err := uc.ResetPassword(ctx, uid, ucadmin.ResetPasswordInput{NewPassword: "newPassword1!"})
	require.NoError(t, err)
}
