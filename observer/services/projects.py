from typing import Protocol

from observer.common.types import Identifier
from observer.entities.base import SomeProject
from observer.entities.projects import NewProject, Project, ProjectUpdate
from observer.entities.users import User
from observer.repositories.projects import ProjectsRepositoryInterface
from observer.schemas.projects import ProjectResponse, ProjectsResponse


class ProjectsServiceInterface(Protocol):
    tag: str
    repo: ProjectsRepositoryInterface

    async def get_by_id(self, project_id: Identifier) -> SomeProject:
        raise NotImplementedError

    async def create_project(self, new_project: NewProject) -> Project:
        raise NotImplementedError

    async def update_project(self, project_id: Identifier, updates: ProjectUpdate) -> Project:
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

    async def create_project(self, new_project: NewProject) -> Project:
        project = await self.repo.create_project(new_project)
        return project

    async def update_project(self, project_id: Identifier, updates: ProjectUpdate) -> Project:
        project = await self.repo.update_project(project_id, updates)
        return project

    @staticmethod
    async def to_response(project: Project) -> ProjectResponse:
        return ProjectResponse(**project.dict())

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> ProjectsResponse:
        return ProjectsResponse(
            total=total,
            items=[ProjectResponse(**user.dict()) for user in user_list],
        )
