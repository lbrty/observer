from fastapi import Depends

from observer.common.types import Identifier
from observer.components.auth import current_user
from observer.components.services import projects_service
from observer.entities.projects import Project
from observer.entities.users import User
from observer.services.projects import ProjectsServiceInterface


async def project_details(
    project_id: Identifier,
    user: User = Depends(current_user),
    projects: ProjectsServiceInterface = Depends(projects_service),
) -> Project:
    project = await projects.get_by_id(project_id)
    return project
