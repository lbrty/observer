from typing import List

from pydantic import Field

from observer.common.types import Identifier
from observer.schemas.base import SchemaBase


class OfficeResponse(SchemaBase):
    id: Identifier = Field(..., description="Office ID")
    name: str = Field(..., description="Name of the office")


class OfficesResponse(SchemaBase):
    total: int = Field(..., description="Total number of offices")
    items: List[OfficeResponse] = Field(..., description="List of offices")


class NewOfficeRequest(SchemaBase):
    name: str = Field(..., description="Name of the office")


class UpdateOfficeRequest(SchemaBase):
    name: str = Field(..., description="Name of the office")
