from datetime import datetime
from typing import Optional

from observer.common.types import Identifier, PetStatus
from observer.entities.base import ModelBase


class BasePet(ModelBase):
    name: str
    notes: Optional[str]
    status: PetStatus
    registration_id: Optional[str]
    owner_id: Identifier
    project_id: Identifier


class Pet(BasePet):
    id: Identifier
    created_at: datetime


class NewPet(BasePet):
    ...


class UpdatePet(BasePet):
    ...
