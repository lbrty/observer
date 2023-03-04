from typing import List, Optional, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.permissions import NewPermission, Permission, UpdatePermission
from observer.repositories.permissions import IPermissionsRepository
from observer.schemas.permissions import NewPermissionRequest, UpdatePermissionRequest


class IPermissionsService(Protocol):
    repo: IPermissionsRepository

    async def get_by_id(self, permission_id: Identifier) -> Optional[Permission]:
        raise NotImplementedError

    # For admins to see the list of users with access to projects
    async def get_by_project_id(self, project_id: Identifier) -> List[Permission]:
        raise NotImplementedError

    # For all users to see which projects they have access to
    async def get_by_user_id(self, user_id: Identifier) -> List[Permission]:
        raise NotImplementedError

    async def find(self, project_id: Optional[Identifier], user_id: Identifier) -> Permission:
        raise NotImplementedError

    async def create_permission(self, new_permission: NewPermissionRequest) -> Permission:
        raise NotImplementedError

    async def update_permission(self, permission_id: Identifier, updates: UpdatePermissionRequest) -> Permission:
        raise NotImplementedError


class PermissionsService(IPermissionsService):
    def __init__(self, repo: IPermissionsRepository):
        self.repo = repo

    async def get_by_id(self, permission_id: Identifier) -> Optional[Permission]:
        return await self.repo.get_by_id(permission_id)

    async def get_by_project_id(self, project_id: Identifier) -> List[Permission]:
        return await self.repo.get_by_project_id(project_id)

    async def get_by_user_id(self, user_id: Identifier) -> List[Permission]:
        return await self.repo.get_by_user_id(user_id)

    async def find(self, project_id: Optional[Identifier], user_id: Identifier) -> Permission:
        if permission := await self.repo.find(project_id, user_id):
            return permission

        raise NotFoundError(message="Project or member not found")

    async def create_permission(self, new_permission: NewPermissionRequest) -> Permission:
        permission = await self.repo.create_permission(NewPermission(**new_permission.dict()))
        return permission

    async def update_permission(self, permission_id: Identifier, updates: UpdatePermissionRequest) -> Permission:
        updates = UpdatePermission(**updates.dict())
        if permission := await self.repo.update_permission(permission_id, updates):
            return permission

        raise NotFoundError(message="Project or member not found")
