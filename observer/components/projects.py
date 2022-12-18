from fastapi import Depends

from observer.api.exceptions import ForbiddenError, NotFoundError
from observer.common.types import Identifier, Role
from observer.components.auth import current_user
from observer.components.services import permissions_service, projects_service
from observer.entities.projects import Project
from observer.entities.users import User
from observer.services.permissions import PermissionsServiceInterface
from observer.services.projects import ProjectsServiceInterface


async def current_project(
    project_id: Identifier,
    projects: ProjectsServiceInterface = Depends(projects_service),
) -> Project:
    if project := await projects.get_by_id(project_id):
        return project
    raise NotFoundError(message="Project not found")


async def viewable_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
) -> Project:
    """Returns project instance if user is admin or has `can_read=True` permission"""
    permission = await permissions.find(project.id, user.id)
    can_read = permission and permission.can_read
    is_admin = user.role == Role.admin
    is_owner = project.owner_id == str(user.id)
    if is_admin or can_read or is_owner:
        return project

    raise ForbiddenError(message="You can not view this project")


async def updatable_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
) -> Project:
    """Returns project instance if user is admin or has `can_update=True` permission"""
    permission = await permissions.find(project.id, user.id)
    can_update = permission and permission.can_update
    is_admin = user.role == Role.admin
    if is_admin or can_update:
        return project

    raise ForbiddenError(message="You can not update this project")


async def deletable_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
) -> Project:
    """Returns project instance if user is admin or has `can_delete=True` permission"""
    permission = await permissions.find(project.id, user.id)
    can_delete = permission and permission.can_delete
    is_admin = user.role == Role.admin
    is_owner = project.owner_id == str(user.id)
    if is_admin or can_delete or is_owner:
        return project

    raise ForbiddenError(message="You can not delete this project")


async def invitable_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
    permissions: PermissionsServiceInterface = Depends(permissions_service),
) -> Project:
    """Returns project instance if user is admin or has `can_invite_members=True` permission"""
    permission = await permissions.find(project.id, user.id)
    can_invite = permission and permission.can_invite_members
    is_admin = user.role == Role.admin
    if is_admin or can_invite:
        return project

    raise ForbiddenError(message="You can not invite members in this project")


async def owned_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
) -> Project:
    """Returns project instance if user is admin or has `can_delete=True` permission"""
    is_owner = project.owner_id == str(user.id)
    is_admin = user.role == Role.admin
    if is_admin or is_owner:
        return project

    raise ForbiddenError(message="Action not permitted")
