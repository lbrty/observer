from typing import Optional

from pydantic import Field

from observer.common.types import (
    DisplacedPersonStatus,
    Identifier,
    SomeDate,
    SomeDatetime,
    SomeIdentifier,
    SomeList,
    SomeStr,
)
from observer.schemas.base import SchemaBase


class BaseCategory(SchemaBase):
    name: str = Field(..., description="Category name")


class CategoryResponse(BaseCategory):
    id: Identifier


class NewCategoryRequest(BaseCategory):
    ...


class UpdateCategoryRequest(BaseCategory):
    ...


class BaseIDP(SchemaBase):
    status: Optional[DisplacedPersonStatus] = Field(DisplacedPersonStatus.registered, description="Current status")
    reference_id: SomeStr = Field(None, description="Reference ID, maybe some of state issued IDs etc.")
    email: SomeStr = Field(None, description="Contact email")
    full_name: str = Field(None, description="Full name")
    birth_date: SomeDatetime = Field(None, description="Birth date")
    notes: SomeStr = Field(None, description="Additional notes")
    phone_number: SomeStr = Field(None, description="Primary phone number")
    phone_number_additional: SomeStr = Field(None, description="Displaced person ID")
    migration_date: SomeDate = Field(None, description="Date when person has moved")
    # Location info
    from_place_id: SomeIdentifier = Field(None, description="Place of origin city/town/village")
    current_place_id: SomeIdentifier = Field(None, description="Current or destination city/town/village")
    project_id: SomeIdentifier = Field(None, description="Related project ID")
    category_id: SomeIdentifier = Field(None, description="Vulnerability category ID")
    tags: SomeList = Field(None, description="List of tags")


class NewIDPRequest(BaseIDP):
    ...


class UpdateIDPRequest(BaseIDP):
    ...


class IDPResponse(BaseIDP):
    id: Identifier = Field(..., description="Displaced person ID")
    external_id: SomeStr = Field(None, description="External identifier")
    # User's id who registered
    consultant_id: SomeIdentifier = Field(..., description="Consultant ID")
    created_at: SomeDatetime = Field(None, description="Creation date")
    updated_at: SomeDatetime = Field(None, description="Update date")
