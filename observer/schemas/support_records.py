from datetime import datetime
from typing import Optional

from pydantic import Field

from observer.common.types import (
    AgeGroup,
    Identifier,
    SupportRecordSubject,
    SupportType,
)
from observer.schemas.base import SchemaBase


class BaseSupportRecord(SchemaBase):
    description: Optional[str] = Field(None, description="Description of support record")
    type: SupportType = Field(..., description="Type of support")
    consultant_id: Identifier = Field(..., description="Consultant ID")
    age_group: Optional[AgeGroup] = Field(None, description="Beneficiary age group")
    record_for: SupportRecordSubject = Field(..., description="Record subject")
    owner_id: Identifier = Field(..., description="Owner of support humans or pets")
    project_id: Identifier = Field(..., description="Project ID")


class SupportRecordResponse(BaseSupportRecord):
    id: Identifier = Field(..., description="Support record ID")
    created_at: datetime = Field(..., description="Creation date")


class SupportRecordsResponse(SchemaBase):
    total: int = Field(..., description="Total count of support records")
    items: list[SupportRecordResponse] = Field(..., description="List of support records")


class NewSupportRecordRequest(BaseSupportRecord):
    ...


class UpdateSupportRecordRequest(BaseSupportRecord):
    ...
