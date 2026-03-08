package handler_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	domainreport "github.com/lbrty/observer/internal/domain/report"
	"github.com/lbrty/observer/internal/handler"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucreport "github.com/lbrty/observer/internal/usecase/report"
)

func newReportHandler(ctrl *gomock.Controller) (*handler.ReportHandler, *repomock.MockReportRepository) {
	repo := repomock.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)
	return handler.NewReportHandler(uc), repo
}

func expectAllReportMethodsEmpty(repo *repomock.MockReportRepository) {
	empty := []domainreport.CountResult{}
	repo.EXPECT().CountConsultations(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountBySex(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByIDPStatus(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByCategory(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByCurrentRegion(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountBySphere(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountPeopleBySphere(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByOffice(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByAgeGroup(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountConsultationsByAgeGroup(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByTag(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountFamilyUnits(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByCaseStatus(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().StatusFlowReport(gomock.Any(), gomock.Any()).Return([]domainreport.StatusFlow{}, nil)
}

func TestReportHandler_Generate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newReportHandler(ctrl)

	projectID := testID().String()
	expectAllReportMethodsEmpty(repo)

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/reports", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Generate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseResponse[map[string]any](w)
	assert.Contains(t, resp, "consultations")
	assert.Contains(t, resp, "by_sex")
	assert.Contains(t, resp, "status_flow")
}

func TestReportHandler_Generate_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newReportHandler(ctrl)

	projectID := testID().String()
	repo.EXPECT().CountConsultations(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("db error"))

	c, w := newTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/reports", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Generate(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
