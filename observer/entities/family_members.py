from datetime import date
from typing import Optional

from observer.common.types import AgeGroup, Identifier, Sex
from observer.entities.base import ModelBase


class BaseFamilyMember(ModelBase):
    full_name: Optional[str]
    birth_date: Optional[date]
    sex: Optional[Sex]
    notes: Optional[str]
    age_group: AgeGroup
    person_id: Identifier


class FamilyMember(BaseFamilyMember):
    id: Identifier
    project_id: Identifier


class NewFamilyMember(BaseFamilyMember):
    project_id: Identifier


class UpdateFamilyMember(BaseFamilyMember):
    ...
