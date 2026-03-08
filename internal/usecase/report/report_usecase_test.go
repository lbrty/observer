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

func TestReportUseCase_GenerateCustom_SingleDimension(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	repo.EXPECT().CustomQuery(
		gomock.Any(), "proj1", "people", []string{"sex"}, gomock.Any(),
	).Return([]domainreport.CustomResult{
		{Dimensions: map[string]string{"sex": "male"}, Count: 10},
		{Dimensions: map[string]string{"sex": "female"}, Count: 15},
	}, 25, nil)

	input := ucreport.CustomReportInput{
		Metric:  "people",
		GroupBy: []string{"sex"},
	}
	out, err := uc.GenerateCustom(context.Background(), "proj1", input)
	require.NoError(t, err)
	assert.Len(t, out.Rows, 2)
	assert.Equal(t, 25, out.Total)
	assert.Equal(t, "people", out.Metric)
	assert.Equal(t, []string{"sex"}, out.GroupBy)
	assert.Equal(t, "male", out.Rows[0].Dimensions["sex"])
	assert.Equal(t, 10, out.Rows[0].Count)
}

func TestReportUseCase_GenerateCustom_TwoDimensions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	repo.EXPECT().CustomQuery(
		gomock.Any(), "proj1", "events", []string{"sex", "age_group"}, gomock.Any(),
	).Return([]domainreport.CustomResult{
		{Dimensions: map[string]string{"sex": "male", "age_group": "child"}, Count: 5},
		{Dimensions: map[string]string{"sex": "female", "age_group": "adult"}, Count: 8},
		{Dimensions: map[string]string{"sex": "male", "age_group": "adult"}, Count: 12},
	}, 25, nil)

	input := ucreport.CustomReportInput{
		Metric:  "events",
		GroupBy: []string{"sex", "age_group"},
	}
	out, err := uc.GenerateCustom(context.Background(), "proj1", input)
	require.NoError(t, err)
	assert.Len(t, out.Rows, 3)
	assert.Equal(t, 25, out.Total)
	assert.Equal(t, []string{"sex", "age_group"}, out.GroupBy)
	assert.Equal(t, "child", out.Rows[0].Dimensions["age_group"])
}

func TestReportUseCase_GenerateCustom_InvalidDimension(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	input := ucreport.CustomReportInput{
		Metric:  "people",
		GroupBy: []string{"nonexistent"},
	}
	out, err := uc.GenerateCustom(context.Background(), "proj1", input)
	assert.Nil(t, out)
	assert.ErrorContains(t, err, "unknown dimension: nonexistent")
}

func TestReportUseCase_GenerateCustom_MoreThanTwoDimensions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	input := ucreport.CustomReportInput{
		Metric:  "people",
		GroupBy: []string{"sex", "age_group", "region"},
	}
	out, err := uc.GenerateCustom(context.Background(), "proj1", input)
	assert.Nil(t, out)
	assert.ErrorContains(t, err, "at most 2 dimensions")
}

func TestReportUseCase_GenerateCustom_EmptyGroupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	input := ucreport.CustomReportInput{
		Metric:  "people",
		GroupBy: []string{},
	}
	out, err := uc.GenerateCustom(context.Background(), "proj1", input)
	assert.Nil(t, out)
	assert.ErrorContains(t, err, "at least one dimension required")
}

func TestReportUseCase_GenerateCustom_NilGroupBy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	input := ucreport.CustomReportInput{
		Metric: "people",
	}
	out, err := uc.GenerateCustom(context.Background(), "proj1", input)
	assert.Nil(t, out)
	assert.ErrorContains(t, err, "at least one dimension required")
}

func TestReportUseCase_GenerateCustom_WithFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	repo.EXPECT().CustomQuery(
		gomock.Any(), "proj1", "people", []string{"office"}, gomock.Any(),
	).Return([]domainreport.CustomResult{
		{Dimensions: map[string]string{"office": "Bishkek"}, Count: 20},
	}, 20, nil)

	input := ucreport.CustomReportInput{
		Metric:   "people",
		GroupBy:  []string{"office"},
		DateFrom: "2024-01-01",
		DateTo:   "2024-12-31",
		Sex:      "female",
	}
	out, err := uc.GenerateCustom(context.Background(), "proj1", input)
	require.NoError(t, err)
	assert.Len(t, out.Rows, 1)
	assert.Equal(t, 20, out.Total)
}

func TestReportUseCase_GenerateCustom_InvalidDateFrom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	input := ucreport.CustomReportInput{
		Metric:   "people",
		GroupBy:  []string{"sex"},
		DateFrom: "not-a-date",
	}
	out, err := uc.GenerateCustom(context.Background(), "proj1", input)
	assert.Nil(t, out)
	assert.ErrorContains(t, err, "parse date_from")
}

func TestReportUseCase_GenerateCustom_InvalidDateTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_repo.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	input := ucreport.CustomReportInput{
		Metric:  "people",
		GroupBy: []string{"sex"},
		DateTo:  "bad-date",
	}
	out, err := uc.GenerateCustom(context.Background(), "proj1", input)
	assert.Nil(t, out)
	assert.ErrorContains(t, err, "parse date_to")
}

func setupReportMocks(ctrl *gomock.Controller) *mock_repo.MockReportRepository {
	return mock_repo.NewMockReportRepository(ctrl)
}

// expectAllReportCalls sets up expectations for all 14 repo methods returning empty slices.
func expectAllReportCalls(repo *mock_repo.MockReportRepository) {
	repo.EXPECT().CountConsultations(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountBySex(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByIDPStatus(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByCategory(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByCurrentRegion(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountBySphere(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountPeopleBySphere(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByOffice(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByAgeGroup(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountConsultationsByAgeGroup(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByTag(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountFamilyUnits(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().CountByCaseStatus(gomock.Any(), gomock.Any()).Return([]domainreport.CountResult{}, nil)
	repo.EXPECT().StatusFlowReport(gomock.Any(), gomock.Any()).Return([]domainreport.StatusFlow{}, nil)
}

func TestReportUseCase_Generate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := setupReportMocks(ctrl)
	expectAllReportCalls(repo)

	uc := ucreport.NewReportUseCase(repo)
	out, err := uc.Generate(context.Background(), "proj-1", ucreport.ReportInput{
		DateFrom: "2025-01-01",
		DateTo:   "2025-12-31",
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "consultations", out.Consultations.Group)
	assert.Equal(t, "by_sex", out.BySex.Group)
	assert.Equal(t, "by_idp_status", out.ByIDPStatus.Group)
	assert.Equal(t, "by_category", out.ByCategory.Group)
	assert.Equal(t, "by_region", out.ByRegion.Group)
	assert.Equal(t, "by_sphere", out.BySphere.Group)
	assert.Equal(t, "people_by_sphere", out.PeopleBySphere.Group)
	assert.Equal(t, "by_office", out.ByOffice.Group)
	assert.Equal(t, "by_age_group", out.ByAgeGroup.Group)
	assert.Equal(t, "consultations_by_age_group", out.ConsultationsByAgeGroup.Group)
	assert.Equal(t, "by_tag", out.ByTag.Group)
	assert.Equal(t, "family_units", out.FamilyUnits.Group)
	assert.Equal(t, "by_case_status", out.ByCaseStatus.Group)
	assert.Empty(t, out.StatusFlow)
}

func TestReportUseCase_Generate_InvalidDateFrom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := setupReportMocks(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	_, err := uc.Generate(context.Background(), "proj-1", ucreport.ReportInput{
		DateFrom: "not-a-date",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse date_from")
}

func TestReportUseCase_Generate_InvalidDateTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := setupReportMocks(ctrl)
	uc := ucreport.NewReportUseCase(repo)

	_, err := uc.Generate(context.Background(), "proj-1", ucreport.ReportInput{
		DateTo: "bad",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse date_to")
}

func TestReportUseCase_Generate_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := setupReportMocks(ctrl)
	repoErr := errors.New("db connection lost")
	repo.EXPECT().CountConsultations(gomock.Any(), gomock.Any()).Return(nil, repoErr)

	uc := ucreport.NewReportUseCase(repo)

	_, err := uc.Generate(context.Background(), "proj-1", ucreport.ReportInput{})

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
