from typing import List

from pydantic import Field

from observer.common.types import Identifier, SomeStr
from observer.schemas.base import SchemaBase
from observer.schemas.permissions import BasePermission
from observer.schemas.users import UserResponse


class ProjectFilters(SchemaBase):
    id: Identifier | None = Field(None, description="ID of project")
    name: SomeStr = Field(None, description="Name of project")
    description: SomeStr = Field(None, description="Description of project")


class BaseProject(SchemaBase):
    name: str = Field(..., description="Name of project")
    description: SomeStr = Field(None, description="Description of project")


class ProjectResponse(SchemaBase):
    id: Identifier = Field(..., description="ID of project")
    name: str = Field(..., description="Name of project")
    description: SomeStr = Field(None, description="Description of project")


class NewProjectRequest(BaseProject):
    ...


class UpdateProjectRequest(BaseProject):
    ...


class ProjectsResponse(SchemaBase):
    total: int = Field(..., description="Total count of projects")
    items: List[ProjectResponse] = Field(..., description="List of projects")


class ProjectMemberResponse(SchemaBase):
    user: UserResponse
    permissions: BasePermission


# TODO: Add pagination
class ProjectMembersResponse(SchemaBase):
    # total: int = Field(..., description="Total count of members")
    items: List[ProjectMemberResponse] = Field(..., description="List of project member")
