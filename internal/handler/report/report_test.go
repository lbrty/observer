package report_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	domainreport "github.com/lbrty/observer/internal/domain/report"
	"github.com/lbrty/observer/internal/handler/handlertest"
	"github.com/lbrty/observer/internal/handler/report"
	repomock "github.com/lbrty/observer/internal/repository/mock"
	ucreport "github.com/lbrty/observer/internal/usecase/report"
)

func newReportHandler(ctrl *gomock.Controller) (*report.ReportHandler, *repomock.MockReportRepository) {
	repo := repomock.NewMockReportRepository(ctrl)
	uc := ucreport.NewReportUseCase(repo)
	return report.NewReportHandler(uc), repo
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

	projectID := handlertest.TestID().String()
	expectAllReportMethodsEmpty(repo)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/reports", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Generate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Contains(t, resp, "consultations")
	assert.Contains(t, resp, "by_sex")
	assert.Contains(t, resp, "status_flow")
}
