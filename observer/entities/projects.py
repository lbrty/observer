from typing import Optional
from uuid import UUID

from pydantic import BaseModel

from observer.common.types import Identifier, Role
from observer.entities.permissions import Permission


class Project(BaseModel):
    id: Identifier
    name: str
    description: Optional[str]
    owner_id: Optional[UUID]


class NewProject(BaseModel):
    name: str
    description: Optional[str]
    owner_id: Optional[UUID]


class ProjectUpdate(BaseModel):
    name: Optional[str]
    description: Optional[str]


class ProjectMember(BaseModel):
    user_id: Identifier
    is_active: bool
    full_name: Optional[str]
    role: Role
    permissions: Permission
