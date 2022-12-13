from pydantic import BaseModel, Field

from observer.common.types import (
    DisplacedPersonStatus,
    Identifier,
    SomeDatetime,
    SomeIdentifier,
    SomeList,
    SomeStr,
)


class VulnerabilityCategory(BaseModel):
    id: Identifier
    name: str


class DisplacedPersonResponse(BaseModel):
    id: Identifier = Field(..., description="Displaced person ID")
    encryption_key: SomeStr = Field(None, description="Encrypted encryption key")
    status: DisplacedPersonStatus = Field(DisplacedPersonStatus.registered, description="Current status")
    external_id: SomeStr = Field(None, description="External identifier")
    reference_id: SomeStr = Field(None, description="Reference ID, maybe some of state issued IDs etc.")
    email: SomeStr = Field(None, description="Contact email")
    full_name: str = Field(None, description="Full name")
    birth_date: SomeDatetime = Field(None, description="Birth date")
    notes: SomeStr = Field(None, description="Additional notes")
    phone_number: SomeStr = Field(None, description="Primary phone number")
    phone_number_additional: SomeStr = Field(None, description="Displaced person ID")
    migration_date: SomeDatetime = Field(None, description="Date when person has moved")
    # Location info
    from_city_id: SomeIdentifier = Field(None, description="City of origin")
    from_state_id: SomeIdentifier = Field(None, description="State of origin")
    current_city_id: SomeIdentifier = Field(None, description="Current or destination city")
    current_state_id: SomeIdentifier = Field(None, description="Current or destination state")
    project_id: SomeIdentifier = Field(None, description="Related project ID")
    category_id: SomeIdentifier = Field(None, description="Vulnerability category ID")
    # User's id who registered
    consultant_id: SomeIdentifier = Field(..., description="Consultant ID")
    tags: SomeList = Field(None, description="List of tags")
    created_at: SomeDatetime = Field(None, description="Creation date")
    updated_at: SomeDatetime = Field(None, description="Update date")
