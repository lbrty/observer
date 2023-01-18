from datetime import datetime
from typing import Optional

from observer.common.types import (
    AgeGroup,
    Identifier,
    SupportRecordSubject,
    SupportType,
)
from observer.entities.base import ModelBase


class BaseSupportRecord(ModelBase):
    description: Optional[str]
    type: SupportType
    consultant_id: Identifier
    age_group: Optional[AgeGroup]
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
