from typing import TypeAlias

from pydantic import BaseModel

from observer.entities.projects import Project
from observer.entities.users import User

SomeUser: TypeAlias = User | None
SomeProject: TypeAlias = Project | None


class ModelBase(BaseModel):
    class Config:
        frozen = True
        use_enum_values = True
