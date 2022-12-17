from pydantic import BaseModel

from observer.common.types import Identifier, Role, SomeStr
from observer.entities.permissions import BasePermission


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


class ProjectMember(BaseModel):
    ref_id: Identifier
    is_active: bool
    full_name: SomeStr
    role: Role
    permissions: BasePermission
