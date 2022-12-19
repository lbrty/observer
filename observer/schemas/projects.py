from typing import List

from pydantic import Field

from observer.common.types import Identifier, Role, SomeIdentifier, SomeStr
from observer.schemas.base import SchemaBase
from observer.schemas.permissions import PermissionResponse


class ProjectFilters(SchemaBase):
    id: SomeIdentifier = Field(None, description="ID of project")
    name: SomeStr = Field(None, description="Name of project")
    description: SomeStr = Field(None, description="Description of project")


class BaseProject(SchemaBase):
    name: str = Field(..., description="Name of project")
    description: SomeStr = Field(None, description="Description of project")


class ProjectResponse(SchemaBase):
    id: Identifier = Field(..., description="ID of project")
    name: str = Field(..., description="Name of project")
    description: SomeStr = Field(None, description="Description of project")
    owner_id: SomeStr = Field(None, description="ID of creator")


class NewProjectRequest(BaseProject):
    owner_id: SomeStr = Field(None, description="ID of creator it overridden currently active user")


class UpdateProjectRequest(BaseProject):
    ...


class ProjectsResponse(SchemaBase):
    total: int = Field(..., description="Total count of projects")
    items: List[ProjectResponse] = Field(..., description="List of projects")


class ProjectMemberResponse(SchemaBase):
    ref_id: Identifier
    is_active: bool
    full_name: SomeStr
    role: Role
    permissions: PermissionResponse


# TODO: Add pagination
class ProjectMembersResponse(SchemaBase):
    # total: int = Field(..., description="Total count of members")
    items: List[ProjectMemberResponse] = Field(..., description="List of project member")
