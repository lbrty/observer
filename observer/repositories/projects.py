from typing import List, Protocol

from sqlalchemy import and_, delete, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.permissions import permissions
from observer.db.tables.projects import projects
from observer.db.tables.users import users
from observer.entities.base import SomeProject
from observer.entities.permissions import BasePermission
from observer.entities.projects import NewProject, Project, ProjectMember, ProjectUpdate


class ProjectsRepositoryInterface(Protocol):
    async def get_by_id(self, project_id: Identifier) -> SomeProject:
        raise NotImplementedError

    async def create_project(self, new_user: NewProject) -> Project:
        raise NotImplementedError

    async def update_project(self, project_id: Identifier, updates: ProjectUpdate) -> Project:
        raise NotImplementedError

    async def delete_project(self, project_id: Identifier) -> Project:
        raise NotImplementedError

    async def get_members(self, project_id: Identifier, offset: int, limit: int) -> List[ProjectMember]:
        raise NotImplementedError


class ProjectsRepository(ProjectsRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def get_by_id(self, project_id: Identifier) -> SomeProject:
        query = select(projects).where(projects.c.id == str(project_id))
        if result := await self.db.fetchone(query):
            return Project(**result)

        return None

    async def create_project(self, new_project: NewProject) -> Project:
        query = insert(projects).values(**new_project.dict()).returning("*")
        result = await self.db.fetchone(query)
        return Project(**result)

    async def update_project(self, project_id: Identifier, updates: ProjectUpdate) -> Project:
        update_values = {}
        if updates.name:
            update_values["name"] = updates.name

        if updates.description:
            update_values["description"] = updates.description  # type:ignore

        query = update(projects).values(update_values).where(projects.c.id == project_id).returning("*")
        result = await self.db.fetchone(query)
        return Project(**result)

    async def delete_project(self, project_id: Identifier) -> Project:
        query = delete(projects).where(projects.c.id == project_id).returning("*")
        result = await self.db.fetchone(query)
        return Project(**result)

    async def get_members(self, project_id: Identifier, offset: int, limit: int) -> List[ProjectMember]:
        join_stmt = permissions.join(
            users,
            and_(
                users.c.id == permissions.c.user_id,
                permissions.c.project_id == project_id,
            ),
        )
        query = (
            select(
                [
                    users.c.ref_id,
                    users.c.full_name,
                    users.c.is_active,
                    users.c.role,
                    permissions.c.can_create,
                    permissions.c.can_read,
                    permissions.c.can_update,
                    permissions.c.can_delete,
                    permissions.c.can_create_projects,
                    permissions.c.can_read_documents,
                    permissions.c.can_read_personal_info,
                ]
            )
            .select_from(join_stmt)
            .offset(offset)
            .limit(limit)
        )

        rows = await self.db.fetchall(query)
        return [self.row_to_member(dict(row)) for row in rows]

    def row_to_member(self, data: dict) -> ProjectMember:
        permission = BasePermission(**data)
        return ProjectMember(
            **dict(
                **data,
                permissions=permission,
            )
        )
