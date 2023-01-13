from datetime import date
from typing import Optional

from pydantic import Field

from observer.common.types import AgeGroup, Identifier, Sex
from observer.entities.base import ModelBase


class BaseFamilyMember(ModelBase):
    full_name: Optional[str] = Field(None, description="Full name of family member")
    birth_date: Optional[date] = Field(None, description="Full name of family member")
    sex: Optional[Sex] = Field(None, description="Sex of family member")
    notes: Optional[str] = Field(None, description="Notes")
    age_group: AgeGroup = Field(..., description="Age group of family member")
    idp_id: Identifier = Field(..., description="Person ID")
    project_id: Identifier = Field(..., description="Project ID")


class FamilyMemberResponse(BaseFamilyMember):
    id: Identifier = Field(..., description="Family member ID")


class NewFamilyMemberRequest(BaseFamilyMember):
    ...


class UpdateFamilyMemberRequest(BaseFamilyMember):
    ...
