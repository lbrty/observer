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


class BaseCategory(BaseModel):
    name: str


class Category(BaseCategory):
    id: Identifier


class NewCategory(BaseCategory):
    ...


class UpdateCategory(BaseCategory):
    ...


class DisplacedPerson(BaseModel):
    id: Identifier
    encryption_key: SomeStr
    status: DisplacedPersonStatus
    external_id: SomeStr
    reference_id: SomeStr
    email: SomeStr
    full_name: str
    birth_date: SomeDate
    notes: SomeStr
    phone_number: SomeStr
    phone_number_additional: SomeStr
    migration_date: SomeDate
    # Location info
    from_place_id: SomeIdentifier
    current_place_id: SomeIdentifier
    project_id: SomeIdentifier
    category_id: SomeIdentifier
    # User's id who registered
    consultant_id: SomeIdentifier
    tags: SomeList
    created_at: SomeDatetime
    updated_at: SomeDatetime
