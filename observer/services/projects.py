from datetime import datetime, timedelta, timezone
from typing import Protocol

from observer.common.types import Identifier
from observer.entities.base import SomeProject
from observer.entities.projects import NewProject, Project, ProjectUpdate
from observer.entities.users import User
from observer.repositories.projects import ProjectsRepositoryInterface
from observer.schemas.audit_logs import NewAuditLog
from observer.schemas.projects import (
    NewProjectRequest,
    ProjectResponse,
    ProjectsResponse,
    UpdateProjectRequest,
)


class ProjectsServiceInterface(Protocol):
    tag: str
    repo: ProjectsRepositoryInterface

    async def get_by_id(self, project_id: Identifier) -> SomeProject:
        raise NotImplementedError

    async def create_project(self, new_project: NewProjectRequest) -> Project:
        raise NotImplementedError

    async def update_project(self, project_id: Identifier, updates: UpdateProjectRequest) -> Project:
        raise NotImplementedError

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        raise NotImplementedError

    @staticmethod
    async def to_response(project: Project) -> ProjectResponse:
        raise NotImplementedError

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> ProjectsResponse:
        raise NotImplementedError


class ProjectsService(ProjectsServiceInterface):
    tag: str = "source=service:projects"

    def __init__(self, repo: ProjectsRepositoryInterface):
        self.repo = repo

    async def get_by_id(self, project_id: Identifier) -> SomeProject:
        if project := await self.repo.get_by_id(project_id):
            return project

        return None

    async def create_project(self, new_project: NewProjectRequest) -> Project:
        project = await self.repo.create_project(NewProject(**new_project.dict()))
        return project

    async def update_project(self, project_id: Identifier, updates: UpdateProjectRequest) -> Project:
        project = await self.repo.update_project(project_id, ProjectUpdate(**updates.dict()))
        return project

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        now = datetime.now(tz=timezone.utc)
        expires_at = None
        if expires_in:
            expires_at = now + expires_in

        return NewAuditLog(
            ref=f"{self.tag},{ref}",
            data=data,
            expires_at=expires_at,
        )

    @staticmethod
    async def to_response(project: Project) -> ProjectResponse:
        return ProjectResponse(**project.dict())

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> ProjectsResponse:
        return ProjectsResponse(
            total=total,
            items=[ProjectResponse(**user.dict()) for user in user_list],
        )
