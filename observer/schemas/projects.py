from pydantic import Field

from observer.common.types import Identifier
from observer.schemas.base import SchemaBase


class BaseProject(SchemaBase):
    name: str = Field(..., description="Name of project")
    description: str | None = Field(None, description="Description of project")


class Project(SchemaBase):
    id: Identifier = Field(..., description="ID of project")
    name: str = Field(..., description="Name of project")
    description: str | None = Field(None, description="Description of project")


class NewProject(BaseProject):
    ...


class UpdateProject(BaseProject):
    ...


class ProjectsResponse(SchemaBase):
    total: int = Field(..., description="Total count of projects")
    items: list[Project] = Field(..., description="List of projects")
