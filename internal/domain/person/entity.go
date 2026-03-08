package person

import (
	"encoding/json"
	"time"
)

// CaseStatus represents the lifecycle state of a person's case.
type CaseStatus string

const (
	CaseStatusNew      CaseStatus = "new"
	CaseStatusActive   CaseStatus = "active"
	CaseStatusClosed   CaseStatus = "closed"
	CaseStatusArchived CaseStatus = "archived"
)

// Sex represents a person's sex.
type Sex string

const (
	SexMale    Sex = "male"
	SexFemale  Sex = "female"
	SexOther   Sex = "other"
	SexUnknown Sex = "unknown"
)

// AgeGroup represents a person's age bracket.
type AgeGroup string

const (
	AgeGroupInfant          AgeGroup = "infant"
	AgeGroupToddler         AgeGroup = "toddler"
	AgeGroupPreSchool       AgeGroup = "pre_school"
	AgeGroupMiddleChildhood AgeGroup = "middle_childhood"
	AgeGroupYoungTeen       AgeGroup = "young_teen"
	AgeGroupTeenager        AgeGroup = "teenager"
	AgeGroupYoungAdult      AgeGroup = "young_adult"
	AgeGroupEarlyAdult      AgeGroup = "early_adult"
	AgeGroupMiddleAged      AgeGroup = "middle_aged_adult"
	AgeGroupOldAdult        AgeGroup = "old_adult"
	AgeGroupUnknown         AgeGroup = "unknown"
)

// Person represents an individual registered in a project.
type Person struct {
	ID             string
	ProjectID      string
	ConsultantID   *string
	OfficeID       *string
	CurrentPlaceID *string
	OriginPlaceID  *string
	ExternalID     *string
	FirstName      string
	LastName       *string
	Patronymic     *string
	Email          *string
	BirthDate      *time.Time
	Sex            Sex
	AgeGroup       *AgeGroup
	PrimaryPhone   *string
	PhoneNumbers   json.RawMessage
	CaseStatus     CaseStatus
	ConsentGiven   bool
	ConsentDate    *time.Time
	RegisteredAt   *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// PersonListFilter holds optional filters for listing people.
type PersonListFilter struct {
	ProjectID    string
	ConsultantID *string
	OfficeID     *string
	CaseStatus   *CaseStatus
	Sex          *Sex
	AgeGroup     *AgeGroup
	CategoryID   *string
	RegionID     *string
	HasPets      *bool
	Search       *string
	TagIDs       []string
	Page         int
	PerPage      int
}
