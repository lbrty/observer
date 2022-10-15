from pydantic import BaseModel, Field

from observer.common.types import Identifier
from observer.schemas.base import SchemaBase


class BasePermission(SchemaBase):
    can_create: bool = Field(False, description="User can create records")
    can_read: bool = Field(True, description="User can create records")
    can_update: bool = Field(False, description="User can update records")
    can_delete: bool = Field(False, description="User can delete records")
    can_create_projects: bool = Field(False, description="User can create projects")
    can_read_documents: bool = Field(False, description="User can read documents")
    can_read_personal_info: bool = Field(False, description="User can read personal info")
    user_id: Identifier = Field(..., description="User can create records")
    project_id: Identifier = Field(..., description="User can create records")


class Permission(BasePermission):
    id: Identifier = Field(..., description="Permission ID")


class UpdatePermission(BasePermission):
    ...


class PermissionsResponse(BaseModel):
    total: int = Field(..., description="Total count of permissions")
    items: list[Permission] = Field(..., description="List of permissions")
