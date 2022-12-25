from typing import Optional

from pydantic import BaseModel

from observer.common.types import (
    DisplacedPersonStatus,
    Identifier,
    SomeDate,
    SomeDatetime,
    SomeIdentifier,
    SomeList,
    SomeStr,
)
from observer.entities.base import SomePlace


class BaseCategory(BaseModel):
    name: str


class Category(BaseCategory):
    id: Identifier


class NewCategory(BaseCategory):
    ...


class UpdateCategory(BaseCategory):
    ...


class BaseIDP(BaseModel):
    status: Optional[DisplacedPersonStatus]
    external_id: SomeStr
    reference_id: SomeStr
    email: SomeStr
    full_name: str
    birth_date: SomeDate
    notes: SomeStr
    phone_number: SomeStr
    phone_number_additional: SomeStr
    project_id: SomeIdentifier
    category_id: SomeIdentifier
    tags: SomeList


class IDP(BaseIDP):
    id: Identifier
    # User's id who registered
    consultant_id: SomeIdentifier
    created_at: SomeDatetime
    updated_at: SomeDatetime


class NewIDP(BaseIDP):
    ...


class UpdateIDP(BaseIDP):
    ...


class PersonalInfo(BaseModel):
    full_name: str
    email: SomeStr
    phone_number: SomeStr
    phone_number_additional: SomeStr
    from_place: SomePlace
    current_place: SomePlace


class BaseMigrationHistory(BaseModel):
    idp_id: Identifier
    migration_date: SomeDate
    from_place_id: SomeIdentifier
    current_place_id: SomeIdentifier


class MigrationHistory(BaseModel):
    id: Identifier
    created_at: SomeDatetime
