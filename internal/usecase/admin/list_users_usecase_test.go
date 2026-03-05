package admin_test

import (
	"context"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/crypto"
	"github.com/lbrty/observer/internal/domain/user"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucadmin "github.com/lbrty/observer/internal/usecase/admin"
)

func newUserUC(t *testing.T) (*ucadmin.UserUseCase, *mock_repo.MockUserRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)
	mockCredRepo := mock_repo.NewMockCredentialsRepository(ctrl)
	hasher := crypto.NewArgonHasher()
	uc := ucadmin.NewUserUseCase(mockUserRepo, mockCredRepo, hasher)
	return uc, mockUserRepo
}

func TestUserUseCase_List_Success(t *testing.T) {
	uc, mockRepo := newUserUC(t)

	ctx := context.Background()
	now := time.Now().UTC()

	users := []*user.User{
		{ID: ulid.MustNew(ulid.Now(), nil), FirstName: "Alice", LastName: "A", Email: "a@test.com", Phone: "+1", Role: user.RoleAdmin, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: ulid.MustNew(ulid.Now(), nil), FirstName: "Bob", LastName: "B", Email: "b@test.com", Phone: "+2", Role: user.RoleStaff, IsActive: true, CreatedAt: now, UpdatedAt: now},
	}

	mockRepo.EXPECT().
		List(ctx, user.UserListFilter{Page: 1, PerPage: 10}).
		Return(users, 2, nil)

	out, err := uc.List(ctx, ucadmin.ListUsersInput{Page: 1, PerPage: 10})
	require.NoError(t, err)
	assert.Equal(t, 2, out.Total)
	assert.Len(t, out.Users, 2)
	assert.Equal(t, "Alice", out.Users[0].FirstName)
	assert.Equal(t, 1, out.Page)
	assert.Equal(t, 10, out.PerPage)
}

func TestUserUseCase_List_DefaultPagination(t *testing.T) {
	uc, mockRepo := newUserUC(t)

	mockRepo.EXPECT().
		List(gomock.Any(), user.UserListFilter{Page: 1, PerPage: 20}).
		Return(nil, 0, nil)

	out, err := uc.List(context.Background(), ucadmin.ListUsersInput{})
	require.NoError(t, err)
	assert.Equal(t, 1, out.Page)
	assert.Equal(t, 20, out.PerPage)
}

func TestUserUseCase_List_ClampPerPage(t *testing.T) {
	uc, mockRepo := newUserUC(t)

	mockRepo.EXPECT().
		List(gomock.Any(), user.UserListFilter{Page: 1, PerPage: 100}).
		Return(nil, 0, nil)

	out, err := uc.List(context.Background(), ucadmin.ListUsersInput{Page: 1, PerPage: 500})
	require.NoError(t, err)
	assert.Equal(t, 100, out.PerPage)
}

func TestUserUseCase_List_WithFilters(t *testing.T) {
	uc, mockRepo := newUserUC(t)

	active := true
	mockRepo.EXPECT().
		List(gomock.Any(), user.UserListFilter{
			Page:     1,
			PerPage:  20,
			Search:   "alice",
			Role:     "admin",
			IsActive: &active,
		}).
		Return(nil, 0, nil)

	_, err := uc.List(context.Background(), ucadmin.ListUsersInput{
		Search:   "alice",
		Role:     "admin",
		IsActive: &active,
	})
	require.NoError(t, err)
}
