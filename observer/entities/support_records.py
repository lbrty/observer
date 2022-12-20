from datetime import datetime
from typing import Optional

from observer.common.types import BeneficiaryAge, Identifier, SomeStr, SupportType
from observer.entities.base import ModelBase


class SupportRecord(ModelBase):
    id: Identifier
    description: SomeStr
    type: SupportType
    consultant_id: Identifier
    beneficiary_age: Optional[BeneficiaryAge]
    owner_id: Identifier
    created_at: datetime
