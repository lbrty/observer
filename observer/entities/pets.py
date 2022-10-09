from datetime import datetime

from pydantic import BaseModel

from observer.common.types import Identifier, PetStatus


class Pet(BaseModel):
    id: Identifier
    name: str
    notes: str | None
    status: PetStatus
    registration_id: str | None
    owner_id: Identifier
    created_at: datetime
