from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.types import Role
from observer.schemas.permissions import InvitePermissionRequest
from observer.schemas.users import UserInviteRequest
from tests.helpers.auth import get_auth_tokens


async def test_admins_and_staff_can_invite_users_to_projects(
    client,
    app_context,
    admin_user,
    staff_user,
    default_project,
):
    for user in [admin_user, staff_user]:
        tokens = await get_auth_tokens(app_context, user)
        email = f"{user.id}-new-consultant@examples.com"
        invite_request = UserInviteRequest(
            email=email,
            role=Role.consultant,
            permissions=[
                InvitePermissionRequest(
                    project_id=default_project.id,
                    can_create=True,
                    can_read=True,
                    can_update=True,
                    can_delete=True,
                    can_create_projects=True,
                    can_read_documents=True,
                    can_read_personal_info=True,
                    can_invite_members=True,
                )
            ],
        )

        resp = await client.post(
            "/admin/invites",
            json=jsonable_encoder(invite_request),
            cookies=tokens.dict(),
        )
        assert resp.status_code == status.HTTP_201_CREATED
        new_user = await app_context.users_service.get_by_email(email)
        assert new_user is not None

    resp = await client.get("/admin/invites", cookies=tokens.dict())
    assert resp.status_code == status.HTTP_200_OK
    assert len(resp.json()) == 2


async def test_admins_and_staff_can_delete_invites(
    client,
    app_context,
    admin_user,
    staff_user,
    default_project,
):
    for user in [admin_user, staff_user]:
        tokens = await get_auth_tokens(app_context, user)
        email = f"{user.id}-new-consultant@examples.com"
        invite_request = UserInviteRequest(
            email=email,
            role=Role.consultant,
            permissions=[
                InvitePermissionRequest(
                    project_id=default_project.id,
                    can_create=True,
                    can_read=True,
                    can_update=True,
                    can_delete=True,
                    can_create_projects=True,
                    can_read_documents=True,
                    can_read_personal_info=True,
                    can_invite_members=True,
                )
            ],
        )

        resp = await client.post(
            "/admin/invites",
            json=jsonable_encoder(invite_request),
            cookies=tokens.dict(),
        )
        assert resp.status_code == status.HTTP_201_CREATED
        new_user = await app_context.users_service.get_by_email(email)
        assert new_user is not None

        invite_code = resp.json()["code"]
        resp = await client.delete(
            f"/admin/invites/{invite_code}?delete_user=true",
            cookies=tokens.dict(),
        )
        assert resp.status_code == status.HTTP_204_NO_CONTENT
        new_user = await app_context.repos.users.get_by_id(new_user.id)
        assert new_user is None
        invite = await app_context.repos.users.get_invite(invite_code)
        assert invite is None


async def test_non_admins_non_staffs_can_not_invite_users(authorized_client, app_context, default_project):
    email = "new-consultant@examples.com"
    invite_request = UserInviteRequest(
        email=email,
        role=Role.consultant,
        permissions=[
            InvitePermissionRequest(
                project_id=default_project.id,
                can_create=True,
                can_read=True,
                can_update=True,
                can_delete=True,
                can_create_projects=True,
                can_read_documents=True,
                can_read_personal_info=True,
                can_invite_members=True,
            )
        ],
    )

    resp = await authorized_client.post(
        "/admin/invites",
        json=jsonable_encoder(invite_request),
    )
    assert resp.status_code == status.HTTP_403_FORBIDDEN
