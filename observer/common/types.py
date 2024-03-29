from datetime import datetime
from enum import Enum
from pathlib import Path
from typing import Dict, Optional, Tuple, TypeAlias
from uuid import UUID

from pydantic import BaseModel, EmailStr, Field

from observer.schemas.base import SchemaBase

Identifier: TypeAlias = UUID | str
EncryptedFieldValue: str = "*" * 8


class Role(str, Enum):
    admin = "admin"
    staff = "staff"
    consultant = "consultant"
    guest = "guest"


class Sex(str, Enum):
    male = "male"
    female = "female"
    unknown = "unknown"


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


class AgeGroup(str, Enum):
    infant = "infant"
    toddler = "toddler"
    pre_school = "pre_school"
    middle_childhood = "middle_childhood"
    young_teen = "young_teen"
    teenager = "teenager"
    young_adult = "young_adult"
    early_adult = "early_adult"
    middle_aged_adult = "middle_aged_adult"
    old_adult = "old_adult"
    unknown = "unknown"


AgeGroupLabels: Dict[AgeGroup, str] = {
    AgeGroup.infant: "0-1",
    AgeGroup.toddler: "1-3",
    AgeGroup.pre_school: "4-5",
    AgeGroup.middle_childhood: "6-11",
    AgeGroup.young_teen: "12-14",
    AgeGroup.teenager: "15-17",
    AgeGroup.young_adult: "18-25",
    AgeGroup.early_adult: "26-34",
    AgeGroup.middle_aged_adult: "35-59",
    AgeGroup.old_adult: "60-100+",
}


class SupportType(str, Enum):
    humanitarian = "humanitarian"
    legal = "legal"
    medical = "medical"
    general = "general"


class SupportRecordSubject(str, Enum):
    person = "person"
    pet = "pet"


class PlaceType(str, Enum):
    city = "city"
    town = "town"
    village = "village"


class StateFilters(BaseModel):
    name: Optional[str]
    code: Optional[str]
    country_id: Optional[Identifier]


class PlaceFilters(BaseModel):
    name: Optional[str]
    code: Optional[str]
    place_type: PlaceType | None
    country_id: Optional[Identifier]
    state_id: Optional[Identifier]


class UserFilters(BaseModel):
    email: Optional[EmailStr]
    full_name: Optional[str]
    role: Optional[Role]
    office_id: Optional[Identifier]
    is_active: Optional[bool]


class StorageKind(str, Enum):
    fs: str = "fs"
    s3: str = "s3"


FileInfo: TypeAlias = Tuple[datetime, str | Path]


class Pagination(SchemaBase):
    limit: int = Field(100, description="How many items to show?")
    offset: int = Field(0, description="What is the starting point?")
