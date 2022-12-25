from typing import List, Protocol

from sqlalchemy import and_, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.permissions import permissions
from observer.entities.base import SomePermission
from observer.entities.permissions import NewPermission, Permission, UpdatePermission


class PermissionsRepositoryInterface(Protocol):
    async def get_by_id(self, permission_id: Identifier) -> SomePermission:
        raise NotImplementedError

    # For admins to see the list of users with access to projects
    async def get_by_project_id(self, project_id: Identifier) -> List[Permission]:
        raise NotImplementedError

    # For all users to see which projects they have access to
    async def get_by_user_id(self, user_id: Identifier) -> List[Permission]:
        raise NotImplementedError

    async def find(self, project_id: Identifier, user_id: Identifier) -> SomePermission:
        raise NotImplementedError

    async def create_permission(self, new_permission: NewPermission) -> Permission:
        raise NotImplementedError

    async def update_permission(self, permission_id: Identifier, updates: UpdatePermission) -> Permission:
        raise NotImplementedError


class PermissionsRepository(PermissionsRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def get_by_id(self, permission_id: Identifier) -> SomePermission:
        query = select(permissions).where(permissions.c.id == str(permission_id))
        if result := await self.db.fetchone(query):
            return Permission(**result)

        return None

    # For admins to see the list of users with access to projects
    async def get_by_project_id(self, project_id: Identifier) -> List[Permission]:
        query = select(permissions).where(permissions.c.project_id == str(project_id))
        return [Permission(**row) for row in await self.db.fetchall(query)]

    # For all users to see which projects they have access to
    async def get_by_user_id(self, user_id: Identifier) -> List[Permission]:
        query = select(permissions).where(permissions.c.user_id == str(user_id))
        return [Permission(**row) for row in await self.db.fetchall(query)]

    async def find(self, project_id: Identifier, user_id: Identifier) -> SomePermission:
        query = select(permissions).where(
            and_(
                permissions.c.user_id == user_id,
                permissions.c.project_id == project_id,
            )
        )
        if result := await self.db.fetchone(query):
            return Permission(**result)
        return None

    async def create_permission(self, new_permission: NewPermission) -> Permission:
        query = insert(permissions).values(**new_permission.dict()).returning("*")
        result = await self.db.fetchone(query)
        return Permission(**result)

    async def update_permission(self, permission_id: Identifier, updates: UpdatePermission) -> Permission:
        update_values = {}
        if updates.can_create is not None:
            update_values["can_create"] = updates.can_create

        if updates.can_read is not None:
            update_values["can_read"] = updates.can_read

        if updates.can_update is not None:
            update_values["can_update"] = updates.can_update

        if updates.can_delete is not None:
            update_values["can_delete"] = updates.can_delete

        if updates.can_create_projects is not None:
            update_values["can_create_projects"] = updates.can_create_projects

        if updates.can_read_documents is not None:
            update_values["can_read_documents"] = updates.can_read_documents

        if updates.can_read_personal_info is not None:
            update_values["can_read_personal_info"] = updates.can_read_personal_info

        if updates.can_invite_members is not None:
            update_values["can_invite_members"] = updates.can_invite_members

        query = update(permissions).values(update_values).where(permissions.c.id == str(permission_id)).returning("*")
        result = await self.db.fetchone(query)
        return Permission(**result)
