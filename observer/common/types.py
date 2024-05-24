from datetime import datetime
from enum import Enum
from pathlib import Path, PosixPath
from typing import Dict, Optional, Tuple, TypeAlias
from uuid import UUID

from pydantic import BaseModel, EmailStr, Field

Identifier: TypeAlias = UUID | str
EncryptedFieldValue: str = "*" * 8


class Role(str, Enum):
    admin: str = "admin"
    staff: str = "staff"
    consultant: str = "consultant"
    guest: str = "guest"


class Sex(str, Enum):
    male: str = "male"
    female: str = "female"
    unknown: str = "unknown"


class DisplacedPersonStatus(str, Enum):
    consulted: str = "consulted"
    needs_call_back: str = "needs_call_back"
    needs_legal_support: str = "needs_legal_support"
    needs_social_support: str = "needs_social_support"
    needs_monitoring: str = "needs_monitoring"
    registered: str = "registered"
    unknown: str = "unknown"


class PetStatus(str, Enum):
    registered: str = "registered"
    adopted: str = "adopted"
    owner_found: str = "owner_found"
    needs_shelter: str = "needs_shelter"
    unknown: str = "unknown"


class AgeGroup(str, Enum):
    infant: str = "infant"
    toddler: str = "toddler"
    pre_school: str = "pre_school"
    middle_childhood: str = "middle_childhood"
    young_teen: str = "young_teen"
    teenager: str = "teenager"
    young_adult: str = "young_adult"
    early_adult: str = "early_adult"
    middle_aged_adult: str = "middle_aged_adult"
    old_adult: str = "old_adult"
    unknown: str = "unknown"


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
    humanitarian: str = "humanitarian"
    legal: str = "legal"
    medical: str = "medical"
    general: str = "general"


class SupportRecordSubject(str, Enum):
    person: str = "person"
    pet: str = "pet"


class StateFilters(BaseModel):
    name: Optional[str]
    code: Optional[str]
    country_id: Optional[Identifier]


class PlaceFilters(BaseModel):
    name: Optional[str]
    code: Optional[str]
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


FileInfo: TypeAlias = Tuple[datetime, str | Path | PosixPath]


class Pagination(BaseModel):
    limit: int = Field(20, description="How many items to show?")
    offset: int = Field(0, description="What is the starting point?")
