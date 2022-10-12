from enum import Enum
from typing import TypeAlias
from uuid import UUID

Identifier: TypeAlias = UUID | str


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


class 		SupportType(str, Enum):
    humanitarian = "humanitarian"
    legal = "legal"
