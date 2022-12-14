from fastapi import Depends

from observer.api.exceptions import ForbiddenError, NotFoundError
from observer.common.types import Identifier, Role
from observer.components.auth import current_user
from observer.components.services import permissions_service, projects_service
from observer.entities.projects import Project
from observer.entities.users import User
from observer.services.permissions import PermissionsServiceInterface
from observer.services.projects import ProjectsServiceInterface


async def project_details(
    project_id: Identifier,
    user: User = Depends(current_user),
    projects: ProjectsServiceInterface = Depends(projects_service),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
) -> Project:
    project = await projects.get_by_id(project_id)
    if project is None:
        raise NotFoundError(message="Project not found")

    if user.role == Role.admin:
        return project

    permission = await permissions.find(project_id, user.id)
    if permission and permission.can_read:
        return project
    else:
        raise ForbiddenError(message="Yoo can not view this project")
