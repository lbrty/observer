from datetime import datetime
from typing import Optional

from observer.common.types import (
    BeneficiaryAge,
    Identifier,
    SomeStr,
    SupportRecordSubject,
    SupportType,
)
from observer.entities.base import ModelBase


class BaseSupportRecord(ModelBase):
    description: SomeStr
    type: SupportType
    consultant_id: Identifier
    beneficiary_age: Optional[BeneficiaryAge]
    record_for: SupportRecordSubject
    owner_id: Identifier
    project_id: Identifier


class SupportRecord(BaseSupportRecord):
    id: Identifier
    created_at: datetime


class NewSupportRecord(BaseSupportRecord):
    ...


class UpdateSupportRecord(BaseSupportRecord):
    ...
