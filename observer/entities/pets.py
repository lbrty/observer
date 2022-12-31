from datetime import datetime

from pydantic import BaseModel

from observer.common.types import Identifier, PetStatus, SomeStr


class Pet(BaseModel):
    id: Identifier
    name: str
    notes: SomeStr
    status: PetStatus
    registration_id: SomeStr
    owner_id: Identifier
    project_id: Identifier
    created_at: datetime
