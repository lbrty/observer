package admin_test

import (
	"testing"

	"go.uber.org/mock/gomock"

	cryptomock "github.com/lbrty/observer/internal/crypto/mock"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
)

type userUCDeps struct {
	userRepo      *mock_repo.MockUserRepository
	credRepo      *mock_repo.MockCredentialsRepository
	sessionRepo   *mock_repo.MockSessionRepository
	loginAttempts *mock_repo.MockLoginAttemptStore
	hasher        *cryptomock.MockPasswordHasher
	auditRepo     *mock_repo.MockAuditLogRepository
}

func newUserUCDeps(t *testing.T) (*ucadmin.UserUseCase, *userUCDeps) {
	t.Helper()
	ctrl := gomock.NewController(t)
	d := &userUCDeps{
		userRepo:      mock_repo.NewMockUserRepository(ctrl),
		credRepo:      mock_repo.NewMockCredentialsRepository(ctrl),
		sessionRepo:   mock_repo.NewMockSessionRepository(ctrl),
		loginAttempts: mock_repo.NewMockLoginAttemptStore(ctrl),
		hasher:        cryptomock.NewMockPasswordHasher(ctrl),
		auditRepo:     mock_repo.NewMockAuditLogRepository(ctrl),
	}
	d.auditRepo.EXPECT().Log(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	auditUC := ucaudit.NewAuditUseCase(d.auditRepo)
	uc := ucadmin.NewUserUseCase(d.userRepo, d.credRepo, d.hasher, d.sessionRepo, d.loginAttempts, auditUC)
	return uc, d
}

func ptr[T any](v T) *T { return &v }
