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
    pass


class PetStatus(str, Enum):
    registered = "registered"
    adopted = "adopted"
    owner_found = "owner_found"
    need_shelter = "need_shelter"
