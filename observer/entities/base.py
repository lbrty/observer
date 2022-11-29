from typing import TypeAlias

from pydantic import BaseModel

from observer.entities.users import User

SomeUser: TypeAlias = User | None


class ModelBase(BaseModel):
    class Config:
        frozen = True
        use_enum_values = True
