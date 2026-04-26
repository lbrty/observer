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

func newPetReportHandler(ctrl *gomock.Controller) (*report.PetReportHandler, *repomock.MockPetReportRepository) {
	repo := repomock.NewMockPetReportRepository(ctrl)
	uc := ucreport.NewPetReportUseCase(repo)
	return report.NewPetReportHandler(uc), repo
}

func TestPetReportHandler_Generate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	h, repo := newPetReportHandler(ctrl)

	projectID := handlertest.TestID().String()
	empty := []domainreport.CountResult{}
	repo.EXPECT().CountByStatus(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByOwnership(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByMonth(gomock.Any(), gomock.Any()).Return(empty, nil)
	repo.EXPECT().CountByStatusByMonth(gomock.Any(), gomock.Any()).Return([]domainreport.MonthlyStatusCount{}, nil)

	c, w := handlertest.NewTestContextWithParams(http.MethodGet, "/projects/"+projectID+"/reports/pets", nil, gin.Params{
		{Key: "project_id", Value: projectID},
	})
	h.Generate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := handlertest.ParseResponse[map[string]any](w)
	assert.Contains(t, resp, "by_status")
	assert.Contains(t, resp, "by_ownership")
	assert.Contains(t, resp, "by_month")
	assert.Contains(t, resp, "by_status_by_month")
}
