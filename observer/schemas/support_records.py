from datetime import datetime

from pydantic import Field

from observer.common.types import BeneficiaryAge, Identifier, SupportType
from observer.schemas.base import SchemaBase


class SupportRecord(SchemaBase):
    id: Identifier = Field(..., description="Support record ID")
    description: str | None = Field(None, description="Description of support record")
    type: SupportType = Field(..., description="Type of support")
    consultant_id: Identifier = Field(..., description="Consultan ID")
    beneficiary_age: BeneficiaryAge | None = Field(..., description="Benefiary age")
    owner_id: Identifier = Field(..., description="Owner of support humans or pets")
    created_at: datetime = Field(..., description="Creation date")


class SupportRecordsResponse(SchemaBase):
    total: int = Field(..., description="Total count of support records")
    items: list[SupportRecord] = Field(..., description="List of support records")
