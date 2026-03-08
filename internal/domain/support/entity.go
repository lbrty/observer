package support

import (
	"time"
)

// SupportType represents the kind of support provided.
type SupportType string

const (
	SupportTypeHumanitarian  SupportType = "humanitarian"
	SupportTypeLegal         SupportType = "legal"
	SupportTypeSocial        SupportType = "social"
	SupportTypePsychological SupportType = "psychological"
	SupportTypeMedical       SupportType = "medical"
	SupportTypeGeneral       SupportType = "general"
)

// SupportSphere represents the topic area of support.
type SupportSphere string

const (
	SphereHousingAssistance SupportSphere = "housing_assistance"
	SphereDocumentRecovery  SupportSphere = "document_recovery"
	SphereSocialBenefits    SupportSphere = "social_benefits"
	SpherePropertyRights    SupportSphere = "property_rights"
	SphereEmploymentRights  SupportSphere = "employment_rights"
	SphereFamilyLaw         SupportSphere = "family_law"
	SphereHealthcareAccess  SupportSphere = "healthcare_access"
	SphereEducationAccess   SupportSphere = "education_access"
	SphereFinancialAid      SupportSphere = "financial_aid"
	SpherePsychSupport      SupportSphere = "psychological_support"
	SphereOther             SupportSphere = "other"
)

// ReferralStatus represents the state of a referral.
type ReferralStatus string

const (
	ReferralPending    ReferralStatus = "pending"
	ReferralAccepted   ReferralStatus = "accepted"
	ReferralCompleted  ReferralStatus = "completed"
	ReferralDeclined   ReferralStatus = "declined"
	ReferralNoResponse ReferralStatus = "no_response"
)

// Record represents a single support interaction.
type Record struct {
	ID               string
	PersonID         string
	ProjectID        string
	ConsultantID     *string
	RecordedBy       *string
	OfficeID         *string
	ReferredToOffice *string
	Type             SupportType
	Sphere           *SupportSphere
	ReferralStatus   *ReferralStatus
	ProvidedAt       *time.Time
	Notes            *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// RecordListFilter holds optional filters for listing support records.
type RecordListFilter struct {
	ProjectID      string
	PersonID       *string
	ConsultantID   *string
	OfficeID       *string
	Type           *SupportType
	Sphere         *SupportSphere
	ReferralStatus *ReferralStatus
	DateFrom       *time.Time
	DateTo         *time.Time
	Page           int
	PerPage        int
}
