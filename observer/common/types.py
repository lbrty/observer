from datetime import date, datetime
from enum import Enum
from pathlib import Path
from typing import Dict, List, Optional, Tuple, TypeAlias
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


class BeneficiaryAge(str, Enum):
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


BeneficiaryAgeLabels: Dict[BeneficiaryAge, str] = {
    BeneficiaryAge.infant: "0-1",
    BeneficiaryAge.toddler: "1-3",
    BeneficiaryAge.pre_school: "4-5",
    BeneficiaryAge.middle_childhood: "6-11",
    BeneficiaryAge.young_teen: "12-14",
    BeneficiaryAge.teenager: "15-17",
    BeneficiaryAge.young_adult: "18-25",
    BeneficiaryAge.early_adult: "26-34",
    BeneficiaryAge.middle_aged_adult: "35-59",
    BeneficiaryAge.old_adult: "60-100+",
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
    name: SomeStr
    code: SomeStr
    country_id: SomeIdentifier


class PlaceFilters(BaseModel):
    name: SomeStr
    code: SomeStr
    place_type: PlaceType | None
    country_id: SomeIdentifier
    state_id: SomeIdentifier


class StorageKind(str, Enum):
    fs: str = "fs"
    s3: str = "s3"


FileInfo: TypeAlias = Tuple[datetime, str | Path]
