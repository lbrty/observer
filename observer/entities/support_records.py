from datetime import datetime
from typing import Optional
from uuid import UUID

from observer.common.types import BeneficiaryAge, Identifier, SupportType
from observer.entities.base import ModelBase


class SupportRecord(ModelBase):
    id: Identifier
    description: str
    type: SupportType
    consultant_id: UUID
    beneficiary_age: Optional[BeneficiaryAge]
    created_at: datetime
