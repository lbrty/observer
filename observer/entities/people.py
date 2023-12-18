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
    status: Optional[DisplacedPersonStatus] = None
    external_id: Optional[str] = None
    reference_id: Optional[str] = None
    email: Optional[str] = None
    full_name: str
    sex: Optional[Sex] = None
    pronoun: Optional[str] = None
    birth_date: Optional[date] = None
    notes: Optional[str] = None
    phone_number: Optional[str] = None
    phone_number_additional: Optional[str] = None
    project_id: Optional[Identifier] = None
    category_id: Optional[Identifier] = None
    office_id: Optional[Identifier] = None
    tags: Optional[List[str]] = None


class Person(BasePerson):
    id: Identifier
    consultant_id: Optional[Identifier] = None
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


class NewPerson(BasePerson):
    ...


class UpdatePerson(BasePerson):
    full_name: Optional[str] = None


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
