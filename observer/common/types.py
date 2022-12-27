from datetime import date, datetime
from enum import Enum
from typing import List, Optional, TypeAlias
from uuid import UUID

from pydantic import BaseModel

Identifier: TypeAlias = UUID | str
SomeStr: TypeAlias = Optional[str]
SomeDate: TypeAlias = Optional[date]
SomeBool: TypeAlias = Optional[bool]
SomeDatetime: TypeAlias = Optional[datetime]
SomeIdentifier: TypeAlias = Optional[Identifier]
SomeList: TypeAlias = Optional[List]

EncryptedFieldValue: str = "*" * 8


class Role(str, Enum):
    admin = "admin"
    staff = "staff"
    consultant = "consultant"
    guest = "guest"


class DisplacedPersonStatus(str, Enum):
    consulted = "consulted"
    needs_call_back = "needs_call_back"
    needs_legal_support = "needs_legal_support"
    needs_social_support = "needs_social_support"
    needs_monitoring = "needs_monitoring"
    registered = "registered"
    unknown = "unknown"


class PetStatus(str, Enum):
    registered = "registered"
    adopted = "adopted"
    owner_found = "owner_found"
    needs_shelter = "needs_shelter"
    unknown = "unknown"


class BeneficiaryAge(str, Enum):
    infant = "0-1"
    toddler = "1-3"
    pre_school = "4-5"
    middle_childhood = "6-11"
    young_teen = "12-14"
    teenager = "15-17"
    young_adult = "18-25"
    early_adult = "26-34"
    middle_aged_adult = "35-59"
    old_adult = "60-100+"


class SupportType(str, Enum):
    humanitarian = "humanitarian"
    legal = "legal"
    medical = "medical"
    general = "general"


class PlaceType(str, Enum):
    city = "city"
    town = "town"
    village = "village"


class StateFilters(BaseModel):
    name: SomeStr
    code: SomeStr
    country_id: SomeIdentifier


class PlaceFilters(BaseModel):
    name: SomeStr
    code: SomeStr
    place_type: PlaceType | None
    country_id: SomeIdentifier
    state_id: SomeIdentifier
