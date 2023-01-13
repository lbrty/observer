from datetime import date
from typing import Optional

from observer.common.types import AgeGroup, Identifier, Sex
from observer.entities.base import ModelBase


class BaseFamilyMember(ModelBase):
    full_name: Optional[str]
    birth_date: Optional[date]
    sex: Optional[Sex]
    notes: str
    age_group: AgeGroup
    idp_id: Identifier
    project_id: Identifier


class FamilyMember(BaseFamilyMember):
    id: Identifier


class NewFamilyMember(BaseFamilyMember):
    ...


class UpdateFamilyMember(BaseFamilyMember):
    ...
