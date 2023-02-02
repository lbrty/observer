from typing import List, Optional, Protocol

from observer.common.types import Identifier
from observer.entities.permissions import NewPermission, Permission, UpdatePermission
from observer.repositories.permissions import IPermissionsRepository
from observer.schemas.permissions import (
    NewPermissionRequest,
    PermissionResponse,
    PermissionsResponse,
    UpdatePermissionRequest,
)


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

    async def find(self, project_id: Identifier, user_id: Identifier) -> Optional[Permission]:
        raise NotImplementedError

    async def create_permission(self, new_permission: NewPermissionRequest) -> Permission:
        raise NotImplementedError

    async def update_permission(
        self,
        permission_id: Identifier,
        updates: UpdatePermissionRequest,
    ) -> Optional[Permission]:
        raise NotImplementedError

    @staticmethod
    async def to_response(permission: Permission) -> PermissionResponse:
        raise NotImplementedError

    @staticmethod
    async def list_to_response(total: int, permission_list: List[Permission]) -> PermissionsResponse:
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

    async def find(self, project_id: Identifier, user_id: Identifier) -> Optional[Permission]:
        return await self.repo.find(project_id, user_id)

    async def create_permission(self, new_permission: NewPermissionRequest) -> Permission:
        permission = await self.repo.create_permission(NewPermission(**new_permission.dict()))
        return permission

    async def update_permission(
        self,
        permission_id: Identifier,
        updates: UpdatePermissionRequest,
    ) -> Optional[Permission]:
        permission = await self.repo.update_permission(permission_id, UpdatePermission(**updates.dict()))
        return permission

    @staticmethod
    async def to_response(permission: Permission) -> PermissionResponse:
        return PermissionResponse(**permission.dict())

    @staticmethod
    async def list_to_response(total: int, permission_list: List[Permission]) -> PermissionsResponse:
        return PermissionsResponse(
            total=total,
            items=[PermissionResponse(**user.dict()) for user in permission_list],
        )
