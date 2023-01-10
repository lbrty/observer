from datetime import datetime
from typing import List, Optional

from pydantic import Field

from observer.common.types import DisplacedPersonStatus, Identifier, Sex
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
    reference_id: Optional[str] = Field(None, description="Reference ID, maybe some of state issued IDs etc.")
    email: Optional[str] = Field(None, description="Contact email")
    full_name: str = Field(None, description="Full name")
    sex: Optional[Sex] = Field(None, description="Person's sex")
    pronoun: Optional[str] = Field(None, description="Person's pronouns")
    birth_date: Optional[datetime] = Field(None, description="Birth date")
    notes: Optional[str] = Field(None, description="Additional notes")
    phone_number: Optional[str] = Field(None, description="Primary phone number")
    phone_number_additional: Optional[str] = Field(None, description="Displaced person ID")
    project_id: Optional[Identifier] = Field(None, description="Related project ID")
    category_id: Optional[Identifier] = Field(None, description="Vulnerability category ID")
    tags: Optional[List[str]] = Field(None, description="List of tags")


class NewIDPRequest(BaseIDP):
    ...


class UpdateIDPRequest(BaseIDP):
    ...


class IDPResponse(BaseIDP):
    id: Identifier = Field(..., description="Displaced person ID")
    external_id: Optional[str] = Field(None, description="External identifier")
    # User's id who registered
    consultant_id: Optional[Identifier] = Field(..., description="Consultant ID")
    created_at: Optional[datetime] = Field(None, description="Creation date")
    updated_at: Optional[datetime] = Field(None, description="Update date")


class PersonalInfoResponse(SchemaBase):
    full_name: str = Field(None, description="Full name")
    sex: Optional[Sex] = Field(None, description="Person's sex")
    pronoun: Optional[str] = Field(None, description="Person's pronouns")
    email: Optional[str] = Field(None, description="Contact email")
    phone_number: Optional[str] = Field(None, description="Primary phone number")
    phone_number_additional: Optional[str] = Field(None, description="Displaced person ID")
