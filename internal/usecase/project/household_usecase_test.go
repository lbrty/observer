package project_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/lbrty/observer/internal/domain/household"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

func TestHouseholdUseCase_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	now := time.Now().UTC()
	mockRepo.EXPECT().List(gomock.Any(), "proj1", gomock.Any(), gomock.Any()).Return([]*household.Household{
		{ID: "h1", ProjectID: "proj1", MemberCount: 3, CreatedAt: now, UpdatedAt: now},
		{ID: "h2", ProjectID: "proj1", MemberCount: 1, CreatedAt: now, UpdatedAt: now},
	}, 2, nil)

	out, err := uc.List(context.Background(), "proj1", ucproject.ListHouseholdsInput{Page: 1, PerPage: 20})
	require.NoError(t, err)
	assert.Len(t, out.Households, 2)
	assert.Equal(t, 2, out.Total)
	assert.Equal(t, "h1", out.Households[0].ID)
	assert.Equal(t, 3, out.Households[0].MemberCount)
}

func TestHouseholdUseCase_List_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	repoErr := errors.New("db error")
	mockRepo.EXPECT().List(gomock.Any(), "proj1", gomock.Any(), gomock.Any()).Return(nil, 0, repoErr)

	_, err := uc.List(context.Background(), "proj1", ucproject.ListHouseholdsInput{})
	assert.ErrorIs(t, err, repoErr)
}

func TestHouseholdUseCase_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	now := time.Now().UTC()
	headID := "p1"
	mockRepo.EXPECT().GetByID(gomock.Any(), "h1").Return(&household.Household{
		ID:           "h1",
		ProjectID:    "proj1",
		HeadPersonID: &headID,
		MemberCount:  2,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)
	mockMemberRepo.EXPECT().List(gomock.Any(), "h1").Return([]*household.Member{
		{HouseholdID: "h1", PersonID: "p1", Relationship: household.RelationshipHead},
		{HouseholdID: "h1", PersonID: "p2", Relationship: household.RelationshipSpouse},
	}, nil)

	dto, err := uc.Get(context.Background(), "h1")
	require.NoError(t, err)
	assert.Equal(t, "h1", dto.ID)
	require.NotNil(t, dto.HeadPersonID)
	assert.Equal(t, "p1", *dto.HeadPersonID)
	assert.Len(t, dto.Members, 2)
	assert.Equal(t, "head", dto.Members[0].Relationship)
	assert.Equal(t, "spouse", dto.Members[1].Relationship)
}

func TestHouseholdUseCase_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, household.ErrHouseholdNotFound)

	_, err := uc.Get(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, household.ErrHouseholdNotFound)
}

func TestHouseholdUseCase_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, h *household.Household) error {
		assert.NotEmpty(t, h.ID)
		assert.Equal(t, "proj1", h.ProjectID)
		require.NotNil(t, h.ReferenceNumber)
		assert.Equal(t, "REF-001", *h.ReferenceNumber)
		return nil
	})

	dto, err := uc.Create(context.Background(), "proj1", ucproject.CreateHouseholdInput{
		ReferenceNumber: ptr("REF-001"),
		HeadPersonID:    ptr("p1"),
	})
	require.NoError(t, err)
	assert.Equal(t, "proj1", dto.ProjectID)
	assert.NotEmpty(t, dto.ID)
	assert.Equal(t, []ucproject.HouseholdMemberDTO{}, dto.Members)
}

func TestHouseholdUseCase_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	repoErr := errors.New("insert failed")
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

	_, err := uc.Create(context.Background(), "proj1", ucproject.CreateHouseholdInput{})
	assert.ErrorIs(t, err, repoErr)
}

func TestHouseholdUseCase_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	now := time.Now().UTC()
	existing := &household.Household{
		ID:        "h1",
		ProjectID: "proj1",
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockRepo.EXPECT().GetByID(gomock.Any(), "h1").Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, h *household.Household) error {
		require.NotNil(t, h.ReferenceNumber)
		assert.Equal(t, "REF-002", *h.ReferenceNumber)
		return nil
	})

	dto, err := uc.Update(context.Background(), "h1", ucproject.UpdateHouseholdInput{
		ReferenceNumber: ptr("REF-002"),
	})
	require.NoError(t, err)
	assert.Equal(t, "h1", dto.ID)
}

func TestHouseholdUseCase_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, household.ErrHouseholdNotFound)

	_, err := uc.Update(context.Background(), "nonexistent", ucproject.UpdateHouseholdInput{})
	assert.ErrorIs(t, err, household.ErrHouseholdNotFound)
}

func TestHouseholdUseCase_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "h1").Return(nil)

	err := uc.Delete(context.Background(), "h1")
	require.NoError(t, err)
}

func TestHouseholdUseCase_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	mockRepo.EXPECT().Delete(gomock.Any(), "nonexistent").Return(household.ErrHouseholdNotFound)

	err := uc.Delete(context.Background(), "nonexistent")
	assert.ErrorIs(t, err, household.ErrHouseholdNotFound)
}

func TestHouseholdUseCase_AddMember_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	mockMemberRepo.EXPECT().Add(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, m *household.Member) error {
		assert.Equal(t, "h1", m.HouseholdID)
		assert.Equal(t, "p3", m.PersonID)
		assert.Equal(t, household.RelationshipChild, m.Relationship)
		return nil
	})

	dto, err := uc.AddMember(context.Background(), "h1", ucproject.AddMemberInput{
		PersonID:     "p3",
		Relationship: "child",
	})
	require.NoError(t, err)
	assert.Equal(t, "p3", dto.PersonID)
	assert.Equal(t, "child", dto.Relationship)
}

func TestHouseholdUseCase_AddMember_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	mockMemberRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(household.ErrMemberExists)

	_, err := uc.AddMember(context.Background(), "h1", ucproject.AddMemberInput{
		PersonID:     "p1",
		Relationship: "head",
	})
	assert.ErrorIs(t, err, household.ErrMemberExists)
}

func TestHouseholdUseCase_RemoveMember_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	mockMemberRepo.EXPECT().Remove(gomock.Any(), "h1", "p3").Return(nil)

	err := uc.RemoveMember(context.Background(), "h1", "p3")
	require.NoError(t, err)
}

func TestHouseholdUseCase_RemoveMember_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockHouseholdRepository(ctrl)
	mockMemberRepo := mock_repo.NewMockHouseholdMemberRepository(ctrl)
	uc := ucproject.NewHouseholdUseCase(mockRepo, mockMemberRepo)

	mockMemberRepo.EXPECT().Remove(gomock.Any(), "h1", "nonexistent").Return(household.ErrMemberNotFound)

	err := uc.RemoveMember(context.Background(), "h1", "nonexistent")
	assert.ErrorIs(t, err, household.ErrMemberNotFound)
}
