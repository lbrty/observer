from typing import Protocol

from sqlalchemy import insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.projects import projects
from observer.entities.base import SomeProject
from observer.entities.projects import NewProject, Project, ProjectUpdate


class ProjectsRepositoryInterface(Protocol):
    async def get_by_id(self, project_id: Identifier) -> SomeProject:
        raise NotImplementedError

    async def create_project(self, new_user: NewProject) -> Project:
        raise NotImplementedError

    async def update_project(self, project_id: Identifier, updates: ProjectUpdate) -> Project:
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
            update_values["full_name"] = updates.description  # type:ignore

        query = update(projects).values(update_values).where(projects.c.id == str(project_id)).returning("*")
        result = await self.db.fetchone(query)
        return Project(**result)
