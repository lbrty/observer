from datetime import datetime
from uuid import UUID

from observer.common.types import Identifier, SupportType
from observer.entities.base import ModelBase


class SupportRecord(ModelBase):
    id: Identifier
    description: str
    type: SupportType
    consultant_id: UUID
    created_at: datetime
