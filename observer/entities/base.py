from typing import TypeAlias

from pydantic import BaseModel

from observer.entities.permissions import Permission
from observer.entities.projects import Project
from observer.entities.users import User

SomeUser: TypeAlias = User | None
SomePermission: TypeAlias = Permission | None
SomeProject: TypeAlias = Project | None


class ModelBase(BaseModel):
    class Config:
        frozen = True
        use_enum_values = True
