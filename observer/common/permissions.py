from typing import Dict

from observer.common.types import Role
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
