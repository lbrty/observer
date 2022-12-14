from pydantic import BaseModel

from observer.common.types import Identifier


class BasePermission(BaseModel):
    can_create: bool
    can_read: bool
    can_update: bool
    can_delete: bool
    can_create_projects: bool
    can_read_documents: bool
    can_read_personal_info: bool


class Permission(BasePermission):
    id: Identifier
    user_id: Identifier
    project_id: Identifier


class UpdatePermission(BasePermission):
    ...


class NewPermission(BasePermission):
    user_id: Identifier
    project_id: Identifier
