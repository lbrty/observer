from typing import TypeAlias

from pydantic import BaseModel

from observer.entities.permissions import Permission
from observer.entities.projects import Project
from observer.entities.users import User
from observer.entities.world import Country, Place, State

SomeUser: TypeAlias = User | None
SomePermission: TypeAlias = Permission | None
SomeProject: TypeAlias = Project | None
SomeCountry: TypeAlias = Country | None
SomeState: TypeAlias = State | None
SomePlace: TypeAlias = Place | None


class ModelBase(BaseModel):
    class Config:
        frozen = True
        use_enum_values = True
