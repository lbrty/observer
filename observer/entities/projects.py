from pydantic import BaseModel

from observer.common.types import Identifier, SomeStr


class Project(BaseModel):
    id: Identifier
    name: str
    description: SomeStr


class NewProject(BaseModel):
    name: str
    description: SomeStr


class ProjectUpdate(BaseModel):
    name: SomeStr
    description: SomeStr
