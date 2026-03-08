package report_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	domainreport "github.com/lbrty/observer/internal/domain/report"
	mock_repo "github.com/lbrty/observer/internal/repository/mock"
	ucreport "github.com/lbrty/observer/internal/usecase/report"
)

func expectAllPetReportCalls(repo *mock_repo.MockPetReportRepository) {
	repo.EXPECT().CountByStatus(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByOwnership(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByMonth(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByStatusByMonth(gomock.Any(), gomock.Any()).Return([]domainreport.MonthlyStatusCount{}, nil)
}

func TestPetReportUseCase_Generate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockPetReportRepository(ctrl)
	expectAllPetReportCalls(repo)

	uc := ucreport.NewPetReportUseCase(repo)
	out, err := uc.Generate(context.Background(), "proj-1", ucreport.PetReportInput{
		DateFrom: "2025-01-01",
		DateTo:   "2025-12-31",
		Status:   "healthy",
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "by_status", out.ByStatus.Group)
	assert.Equal(t, "by_ownership", out.ByOwnership.Group)
	assert.Equal(t, "by_month", out.ByMonth.Group)
	assert.Empty(t, out.ByStatusByMonth)
}

func TestPetReportUseCase_Generate_InvalidDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockPetReportRepository(ctrl)
	uc := ucreport.NewPetReportUseCase(repo)

	_, err := uc.Generate(context.Background(), "proj-1", ucreport.PetReportInput{
		DateFrom: "not-a-date",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse date_from")
}

func TestPetReportUseCase_Generate_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockPetReportRepository(ctrl)
	repoErr := errors.New("db connection lost")
	repo.EXPECT().CountByStatus(gomock.Any(), gomock.Any()).Return(nil, repoErr)

	uc := ucreport.NewPetReportUseCase(repo)

	_, err := uc.Generate(context.Background(), "proj-1", ucreport.PetReportInput{})

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
