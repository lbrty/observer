from typing import Tuple

from fastapi import Depends

from observer.api.exceptions import ForbiddenError, NotFoundError
from observer.common.permissions import permission_matrix
from observer.common.types import Identifier, Role
from observer.components.auth import current_user
from observer.components.services import permissions_service, projects_service
from observer.entities.permissions import BasePermission
from observer.entities.projects import Project, ProjectMember
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


async def project_with_member(
    project_id: Identifier,
    user: User = Depends(current_user),
    projects: ProjectsServiceInterface = Depends(projects_service),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
) -> Tuple[Project, ProjectMember]:
    """We check if user is admin user or the one with `can_update` permission"""
    project = await projects.get_by_id(project_id)
    if project is None:
        raise NotFoundError(message="Project not found")

    """
    NOTE:
        Since both `permission` instances intersect and are sub-classes of
        Pydantic models both will have `.dict()` method.
    """
    if user.role == Role.admin:
        permission = permission_matrix[Role.admin]
    else:
        permission = await permissions.find(project_id, user.id)

    if permission and permission.can_update:
        return (
            project,
            ProjectMember(
                **dict(
                    **user.dict(),
                    permissions=BasePermission(**permission.dict()),
                )
            ),
        )
    else:
        raise ForbiddenError(message="Yoo can not view this project")
