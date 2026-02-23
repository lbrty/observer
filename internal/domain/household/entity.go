package household

import "time"

// Relationship describes a household member's relation to the head.
type Relationship string

const (
	RelationshipHead          Relationship = "head"
	RelationshipSpouse        Relationship = "spouse"
	RelationshipChild         Relationship = "child"
	RelationshipParent        Relationship = "parent"
	RelationshipSibling       Relationship = "sibling"
	RelationshipGrandchild    Relationship = "grandchild"
	RelationshipGrandparent   Relationship = "grandparent"
	RelationshipOtherRelative Relationship = "other_relative"
	RelationshipNonRelative   Relationship = "non_relative"
)

// Household represents a family unit within a project.
type Household struct {
	ID              string
	ProjectID       string
	ReferenceNumber *string
	HeadPersonID    *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Member represents a person's membership in a household.
type Member struct {
	HouseholdID  string
	PersonID     string
	Relationship Relationship
}
