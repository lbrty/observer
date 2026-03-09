package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lbrty/observer/internal/middleware"
	ucaudit "github.com/lbrty/observer/internal/usecase/audit"
	ucproject "github.com/lbrty/observer/internal/usecase/project"
)

// ExportHandler streams filtered data as CSV downloads.
type ExportHandler struct {
	personUC  *ucproject.PersonUseCase
	supportUC *ucproject.SupportRecordUseCase
	petUC     *ucproject.PetUseCase
	householdUC *ucproject.HouseholdUseCase
	auditUC   *ucaudit.AuditUseCase
}

// NewExportHandler creates an ExportHandler.
func NewExportHandler(
	personUC *ucproject.PersonUseCase,
	supportUC *ucproject.SupportRecordUseCase,
	petUC *ucproject.PetUseCase,
	householdUC *ucproject.HouseholdUseCase,
	auditUC *ucaudit.AuditUseCase,
) *ExportHandler {
	return &ExportHandler{
		personUC:    personUC,
		supportUC:   supportUC,
		petUC:       petUC,
		householdUC: householdUC,
		auditUC:     auditUC,
	}
}

func ptrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ExportPeople streams people as a CSV download.
func (h *ExportHandler) ExportPeople(c *gin.Context) {
	projectID := c.Param("project_id")

	var input ucproject.ListPeopleInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	input.PerPage = 10000

	canContact := middleware.CanViewContactFrom(c)
	canPersonal := middleware.CanViewPersonalFrom(c)

	out, err := h.personUC.List(c.Request.Context(), projectID, input, canContact, canPersonal)
	if err != nil {
		internalError(c, "export people", err)
		return
	}

	filename := fmt.Sprintf("people-%s.csv", projectID)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{
		"id", "first_name", "last_name", "patronymic", "email",
		"sex", "age_group", "case_status", "primary_phone",
		"registered_at", "created_at",
	})

	for _, p := range out.People {
		_ = w.Write([]string{
			p.ID,
			p.FirstName,
			ptrStr(p.LastName),
			ptrStr(p.Patronymic),
			ptrStr(p.Email),
			p.Sex,
			ptrStr(p.AgeGroup),
			p.CaseStatus,
			ptrStr(p.PrimaryPhone),
			ptrStr(p.RegisteredAt),
			p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	w.Flush()

	h.auditUC.Record(c.Request.Context(), &projectID, "export", "person", nil, fmt.Sprintf("exported %d people", len(out.People)))
}

// ExportSupportRecords streams support records as a CSV download.
func (h *ExportHandler) ExportSupportRecords(c *gin.Context) {
	projectID := c.Param("project_id")

	var input ucproject.ListSupportRecordsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	input.PerPage = 10000

	out, err := h.supportUC.List(c.Request.Context(), projectID, input)
	if err != nil {
		internalError(c, "export support records", err)
		return
	}

	filename := fmt.Sprintf("support-records-%s.csv", projectID)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{
		"id", "person_id", "type", "sphere", "office_id",
		"referral_status", "provided_at", "created_at",
	})

	for _, r := range out.Records {
		_ = w.Write([]string{
			r.ID,
			r.PersonID,
			r.Type,
			ptrStr(r.Sphere),
			ptrStr(r.OfficeID),
			ptrStr(r.ReferralStatus),
			ptrStr(r.ProvidedAt),
			r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	w.Flush()

	h.auditUC.Record(c.Request.Context(), &projectID, "export", "support_record", nil, fmt.Sprintf("exported %d support records", len(out.Records)))
}

// ExportPets streams pets as a CSV download.
func (h *ExportHandler) ExportPets(c *gin.Context) {
	projectID := c.Param("project_id")

	var input ucproject.ListPetsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	input.PerPage = 10000

	out, err := h.petUC.List(c.Request.Context(), projectID, input)
	if err != nil {
		internalError(c, "export pets", err)
		return
	}

	filename := fmt.Sprintf("pets-%s.csv", projectID)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{
		"id", "name", "status", "registration_id", "owner_id", "created_at",
	})

	for _, p := range out.Pets {
		_ = w.Write([]string{
			p.ID,
			p.Name,
			p.Status,
			ptrStr(p.RegistrationID),
			ptrStr(p.OwnerID),
			p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	w.Flush()

	h.auditUC.Record(c.Request.Context(), &projectID, "export", "pet", nil, fmt.Sprintf("exported %d pets", len(out.Pets)))
}

// ExportHouseholds streams households as a CSV download.
func (h *ExportHandler) ExportHouseholds(c *gin.Context) {
	projectID := c.Param("project_id")

	var input ucproject.ListHouseholdsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, errJSON("errors.validation", err.Error()))
		return
	}
	input.PerPage = 10000

	out, err := h.householdUC.List(c.Request.Context(), projectID, input)
	if err != nil {
		internalError(c, "export households", err)
		return
	}

	filename := fmt.Sprintf("households-%s.csv", projectID)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{
		"id", "reference_number", "head_person_id", "member_count", "created_at",
	})

	for _, h := range out.Households {
		_ = w.Write([]string{
			h.ID,
			ptrStr(h.ReferenceNumber),
			ptrStr(h.HeadPersonID),
			strconv.Itoa(h.MemberCount),
			h.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	w.Flush()

	h.auditUC.Record(c.Request.Context(), &projectID, "export", "household", nil, fmt.Sprintf("exported %d households", len(out.Households)))
}
