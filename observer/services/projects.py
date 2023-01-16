from typing import List, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.permissions import NewPermission, Permission
from observer.entities.projects import NewProject, Project, ProjectMember, ProjectUpdate
from observer.entities.users import User
from observer.repositories.projects import IProjectsRepository
from observer.schemas.permissions import NewPermissionRequest
from observer.schemas.projects import (
    NewProjectRequest,
    ProjectResponse,
    ProjectsResponse,
    UpdateProjectRequest,
)


class IProjectsService(Protocol):
    repo: IProjectsRepository

    async def get_by_id(self, project_id: Identifier) -> Project:
        raise NotImplementedError

    async def create_project(self, new_project: NewProjectRequest) -> Project:
        raise NotImplementedError

    async def update_project(self, project_id: Identifier, updates: UpdateProjectRequest) -> Project:
        raise NotImplementedError

    async def delete_project(self, project_id: Identifier) -> Project:
        raise NotImplementedError

    async def get_members(self, project_id: Identifier, offset: int, limit: int) -> List[ProjectMember]:
        raise NotImplementedError

    async def add_member(self, project_id: Identifier, new_permission: NewPermissionRequest) -> Permission:
        raise NotImplementedError

    async def delete_member(self, project_id: Identifier, user_id: Identifier) -> Permission:
        raise NotImplementedError

    @staticmethod
    async def to_response(project: Project) -> ProjectResponse:
        raise NotImplementedError

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> ProjectsResponse:
        raise NotImplementedError


class ProjectsService(IProjectsService):
    def __init__(self, repo: IProjectsRepository):
        self.repo = repo

    async def get_by_id(self, project_id: Identifier) -> Project:
        if project := await self.repo.get_by_id(project_id):
            return project

        raise NotFoundError(message="Project not found")

    async def create_project(self, new_project: NewProjectRequest) -> Project:
        project = await self.repo.create_project(NewProject(**new_project.dict()))
        return project

    async def update_project(self, project_id: Identifier, updates: UpdateProjectRequest) -> Project:
        project = await self.repo.update_project(project_id, ProjectUpdate(**updates.dict()))
        return project

    async def delete_project(self, project_id: Identifier) -> Project:
        project = await self.repo.delete_project(project_id)
        return project

    async def get_members(self, project_id: Identifier, offset: int, limit: int) -> List[ProjectMember]:
        members = await self.repo.get_members(project_id, offset, limit)
        return members

    async def add_member(self, project_id: Identifier, new_permission: NewPermissionRequest) -> Permission:
        permission = await self.repo.add_member(project_id, NewPermission(**new_permission.dict()))
        return permission

    async def delete_member(self, project_id: Identifier, user_id: Identifier) -> Permission:
        permission = await self.repo.delete_member(project_id, user_id)
        return permission

    @staticmethod
    async def to_response(project: Project) -> ProjectResponse:
        return ProjectResponse(**project.dict())

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> ProjectsResponse:
        return ProjectsResponse(
            total=total,
            items=[ProjectResponse(**user.dict()) for user in user_list],
        )
