from datetime import datetime

from observer.common.types import Identifier, PetStatus, SomeStr
from observer.entities.base import ModelBase


class BasePet(ModelBase):
    name: str
    notes: SomeStr
    status: PetStatus
    registration_id: SomeStr
    owner_id: Identifier
    project_id: Identifier


class Pet(BasePet):
    id: Identifier
    created_at: datetime


class NewPet(BasePet):
    ...


class UpdatePet(BasePet):
    ...
