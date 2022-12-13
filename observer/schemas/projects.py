from pydantic import Field

from observer.common.types import Identifier
from observer.schemas.base import SchemaBase


class ProjectFilters(SchemaBase):
    id: Identifier | None = Field(None, description="ID of project")
    name: str | None = Field(None, description="Name of project")
    description: str | None = Field(None, description="Description of project")


class BaseProject(SchemaBase):
    name: str = Field(..., description="Name of project")
    description: str | None = Field(None, description="Description of project")


class ProjectResponse(SchemaBase):
    id: Identifier = Field(..., description="ID of project")
    name: str = Field(..., description="Name of project")
    description: str | None = Field(None, description="Description of project")


class NewProjectRequest(BaseProject):
    ...


class UpdateProjectRequest(BaseProject):
    ...


class ProjectsResponse(SchemaBase):
    total: int = Field(..., description="Total count of projects")
    items: list[ProjectResponse] = Field(..., description="List of projects")
