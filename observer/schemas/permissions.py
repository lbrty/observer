from pydantic import Field

from observer.common.types import Identifier
from observer.schemas.base import SchemaBase


class BasePermission(SchemaBase):
    can_create: bool = Field(False, description="User can create records")
    can_read: bool = Field(True, description="User can create records")
    can_update: bool = Field(False, description="User can update records")
    can_delete: bool = Field(False, description="User can delete records")
    can_read_documents: bool = Field(False, description="User can read documents")
    can_read_personal_info: bool = Field(False, description="User can read personal info")
    can_invite_members: bool = Field(False, description="User can invite new members to projects")


class PermissionResponse(BasePermission):
    id: Identifier = Field(..., description="Permission ID")
    user_id: Identifier = Field(..., description="User can create records")
    project_id: Identifier = Field(..., description="User can create records")


class NewPermissionRequest(BasePermission):
    user_id: Identifier = Field(..., description="User can create records")
    project_id: Identifier = Field(..., description="User can create records")


class InvitePermissionRequest(BasePermission):
    project_id: Identifier = Field(..., description="User can create records")


class UpdatePermissionRequest(BasePermission):
    ...


class PermissionsResponse(SchemaBase):
    total: int = Field(..., description="Total count of permissions")
    items: list[PermissionResponse] = Field(..., description="List of permissions")
