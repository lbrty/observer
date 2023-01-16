from datetime import datetime
from typing import List

from pydantic import BaseModel, Field

from observer.common.types import Identifier, PetStatus
from observer.schemas.base import SchemaBase


class PetFilter(SchemaBase):
    name: str | None = Field(None, description="Pet's name")
    notes: str | None = Field(None, description="Additional notes")
    status: PetStatus | None = Field(None, description="Pet status")
    registration_id: str | None = Field(None, description="Pet's registration ID, passport ID etc.")
    owner_id: Identifier | None = Field(None, description="Pet's owner ID")


class BasePet(SchemaBase):
    name: str = Field(..., description="Pet's name")
    notes: str | None = Field(None, description="Additional notes")
    status: PetStatus = Field(..., description="Pet status")
    registration_id: str | None = Field(None, description="Pet's registration ID, passport ID etc.")
    owner_id: Identifier = Field(..., description="Pet's owner ID (Person ID)")
    project_id: Identifier = Field(..., description="Project owner ID")


class PetResponse(BaseModel):
    id: Identifier = Field(..., description="Pet ID")
    created_at: datetime = Field(..., description="Creation datetime")


class NewPetRequest(BasePet):
    ...


class UpdatePetRequest(BasePet):
    ...


class PetsResponse(BaseModel):
    total: int = Field(..., description="Total count of pets")
    items: List[PetResponse] = Field(..., description="List of pets")
