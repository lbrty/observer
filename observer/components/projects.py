from fastapi import Depends, Path

from observer.api.exceptions import ForbiddenError
from observer.common.permissions import (
    assert_can_invite,
    assert_deletable,
    assert_updatable,
    assert_viewable,
)
from observer.common.types import Identifier
from observer.components.auth import current_user
from observer.components.services import permissions_service, projects_service
from observer.entities.projects import Project
from observer.entities.users import User
from observer.services.permissions import IPermissionsService
from observer.services.projects import IProjectsService


async def current_project(
    project_id: Identifier = Path(..., description="Project ID path param"),
    projects: IProjectsService = Depends(projects_service),
) -> Project:
    return await projects.get_by_id(project_id)


async def viewable_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
    permissions: IPermissionsService = Depends(permissions_service),
) -> Project:
    """Returns project instance if user is admin or has `can_read=True` permission"""
    is_owner = project.owner_id == user.id
    if not is_owner:
        permission = await permissions.find(project.id, user.id)
        assert_viewable(user, permission)

    return project


async def updatable_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
    permissions: IPermissionsService = Depends(permissions_service),
) -> Project:
    """Returns project instance if user is admin or has `can_update=True` permission"""
    is_owner = project.owner_id == user.id
    if not is_owner:
        permission = await permissions.find(project.id, user.id)
        assert_updatable(user, permission)

    return project


async def deletable_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
    permissions: IPermissionsService = Depends(permissions_service),
) -> Project:
    """Returns project instance if user is admin or has `can_delete=True` permission"""
    is_owner = project.owner_id == user.id
    if not is_owner:
        permission = await permissions.find(project.id, user.id)
        assert_deletable(user, permission)

    return project


async def invitable_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
    permissions: IPermissionsService = Depends(permissions_service),
) -> Project:
    """Returns project instance if user is admin or has `can_invite_members=True` permission"""
    is_owner = project.owner_id == user.id
    if not is_owner:
        permission = await permissions.find(project.id, user.id)
        assert_can_invite(user, permission)

    return project


async def owned_project(
    user: User = Depends(current_user),
    project: Project = Depends(current_project),
) -> Project:
    """Returns project instance if user is admin or has `can_delete=True` permission"""
    is_owner = project.owner_id == user.id
    if is_owner:
        return project

    raise ForbiddenError(message="Action not permitted")
