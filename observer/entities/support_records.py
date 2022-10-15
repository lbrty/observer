from datetime import datetime
from typing import Optional

from observer.common.types import BeneficiaryAge, Identifier, SupportType
from observer.entities.base import ModelBase


class SupportRecord(ModelBase):
    id: Identifier
    description: str
    type: SupportType
    consultant_id: Identifier
    beneficiary_age: Optional[BeneficiaryAge]
    displaced_person_id: Identifier
    created_at: datetime
