from typing import Dict

from observer.api.exceptions import ForbiddenError
from observer.common.types import Role
from observer.entities.permissions import Permission
from observer.schemas.permissions import BasePermission

FullProjectAccess = BasePermission(
    can_create=True,
    can_read=True,
    can_update=True,
    can_delete=True,
    can_create_projects=True,
    can_read_documents=True,
    can_read_personal_info=True,
    can_invite_members=True,
)

permission_matrix: Dict[Role, BasePermission] = {
    Role.admin: FullProjectAccess,
    Role.consultant: FullProjectAccess,
    Role.staff: BasePermission(
        can_create=True,
        can_read=True,
        can_update=True,
        can_delete=False,
        can_create_projects=True,
        can_read_documents=False,
        can_read_personal_info=False,
        can_invite_members=True,
    ),
    Role.guest: BasePermission(
        can_create=False,
        can_read=True,
        can_update=False,
        can_delete=False,
        can_create_projects=False,
        can_read_documents=False,
        can_read_personal_info=False,
        can_invite_members=False,
    ),
}


def assert_viewable(permission: Permission | None):
    can_do = permission and permission.can_read
    if not permission or not can_do:
        raise ForbiddenError(message="Permission denied")


def assert_writable(permission: Permission | None):
    can_do = permission and permission.can_create
    if not permission or not can_do:
        raise ForbiddenError(message="Permission denied")


def assert_deletable(permission: Permission | None):
    can_do = permission and permission.can_delete
    if not permission or not can_do:
        raise ForbiddenError(message="Permission denied")


def assert_docs_readable(permission: Permission | None):
    can_do = permission and permission.can_read_documents
    if not permission or not can_do:
        raise ForbiddenError(message="Permission denied")


def assert_docs_updatable(permission: Permission | None):
    can_do = permission and permission.can_update
    if not permission or not can_do:
        raise ForbiddenError(message="Permission denied")


def assert_can_invite(permission: Permission | None):
    can_invite = permission and permission.can_invite_members
    if not permission or not can_invite:
        raise ForbiddenError(message="Permission denied")


def assert_can_see_private_info(permission: Permission | None):
    can_do = permission and permission.can_read_personal_info
    if not permission or not can_do:
        raise ForbiddenError(message="Permission denied")


def assert_can_create_projects(permission: Permission | None):
    can_do = permission and permission.can_create_projects
    if not permission or not can_do:
        raise ForbiddenError(message="Permission denied")
