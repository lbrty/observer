from datetime import date, datetime
from typing import List, Optional

from pydantic import BaseModel

from observer.common.types import DisplacedPersonStatus, Identifier, Sex
from observer.entities.world import Place


class BaseCategory(BaseModel):
    name: str


class Category(BaseCategory):
    id: Identifier


class NewCategory(BaseCategory):
    ...


class UpdateCategory(BaseCategory):
    ...


class BasePerson(BaseModel):
    status: Optional[DisplacedPersonStatus]
    external_id: Optional[str]
    reference_id: Optional[str]
    email: Optional[str]
    full_name: str
    sex: Optional[Sex]
    pronoun: Optional[str]
    birth_date: Optional[date]
    notes: Optional[str]
    phone_number: Optional[str]
    phone_number_additional: Optional[str]
    project_id: Optional[Identifier]
    category_id: Optional[Identifier]
    office_id: Optional[Identifier]
    tags: Optional[List[str]]


class Person(BasePerson):
    id: Identifier
    consultant_id: Optional[Identifier]
    created_at: Optional[datetime]
    updated_at: Optional[datetime]


class NewPerson(BasePerson):
    ...


class UpdatePerson(BasePerson):
    full_name: Optional[str]  # type: ignore


class PersonalInfo(BaseModel):
    full_name: Optional[str] = None
    sex: Optional[Sex] = None
    pronoun: Optional[str] = None
    email: Optional[str] = None
    phone_number: Optional[str] = None
    phone_number_additional: Optional[str] = None
    from_place: Optional[Place] = None
    current_place: Optional[Place] = None


class BaseMigrationHistory(BaseModel):
    person_id: Identifier
    migration_date: Optional[date] = None
    from_place_id: Optional[Identifier] = None
    current_place_id: Optional[Identifier] = None


class MigrationHistory(BaseModel):
    id: Identifier
    created_at: Optional[datetime] = None
