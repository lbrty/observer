from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.types import Role
from observer.schemas.users import (
    AdminUpdateUserRequest,
    CreateUserRequest,
    UserResponse,
)
from tests.helpers.auth import get_auth_tokens


async def test_admins_can_create_users(
    client,
    app_context,
    admin_user,
    secure_password,
):
    tokens = await get_auth_tokens(app_context, admin_user)
    email = f"{admin_user.id}-new-guest@examples.com"
    create_user_request = CreateUserRequest(
        email=email,
        full_name=f"Guest {email}",
        password=secure_password,
        role=Role.guest,
        is_active=True,
    )
    payload = jsonable_encoder(create_user_request)
    payload["password"] = secure_password
    resp = await client.post(
        "/admin/users",
        json=payload,
        cookies=tokens.dict(),
    )
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    new_user = await app_context.users_service.get_by_email(email)
    assert new_user is not None
    assert new_user.is_active
    assert new_user.is_confirmed


async def test_admins_can_filter_users(
    client,
    app_context,
    admin_user,
    staff_user,
    consultant_user,
    guest_user,
):
    tokens = await get_auth_tokens(app_context, admin_user)
    resp = await client.get(
        "/admin/users",
        cookies=tokens.dict(),
    )
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == {
        "total": 4,
        "items": [
            jsonable_encoder(UserResponse(**admin_user.dict())),
            jsonable_encoder(UserResponse(**staff_user.dict())),
            jsonable_encoder(UserResponse(**consultant_user.dict())),
            jsonable_encoder(UserResponse(**guest_user.dict())),
        ],
    }

    resp = await client.get(
        f"/admin/users?offset=2&limit=2",
        cookies=tokens.dict(),
    )

    assert resp.json() == {
        "total": 4,
        "items": [
            jsonable_encoder(UserResponse(**consultant_user.dict())),
            jsonable_encoder(UserResponse(**guest_user.dict())),
        ],
    }

    resp = await client.get(
        f"/admin/users?email={guest_user.email}",
        cookies=tokens.dict(),
    )

    assert resp.json() == {
        "total": 1,
        "items": [
            jsonable_encoder(UserResponse(**guest_user.dict())),
        ],
    }


async def test_admins_can_get_user_details(
    client,
    app_context,
    admin_user,
    guest_user,
    secure_password,
):
    tokens = await get_auth_tokens(app_context, admin_user)
    resp = await client.get(
        f"/admin/users/{guest_user.id}",
        cookies=tokens.dict(),
    )
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == jsonable_encoder(UserResponse(**guest_user.dict()))


async def test_admins_can_update_users(
    client,
    app_context,
    admin_user,
    guest_user,
):
    tokens = await get_auth_tokens(app_context, admin_user)
    guest_user.role = Role.consultant
    guest_user.full_name = "FN Guest"
    updates = AdminUpdateUserRequest(**guest_user.dict())
    resp = await client.put(
        f"/admin/users/{guest_user.id}",
        json=jsonable_encoder(updates),
        cookies=tokens.dict(),
    )
    updated_guest_user = await app_context.repos.users.get_by_id(guest_user.id)
    assert resp.status_code == status.HTTP_200_OK
    assert resp.json() == jsonable_encoder(UserResponse(**updated_guest_user.dict()))


async def test_admins_can_delete_users(
    client,
    app_context,
    admin_user,
    consultant_user,
    default_person,
    secure_password,
):
    tokens = await get_auth_tokens(app_context, admin_user)
    resp = await client.delete(
        f"/admin/users/{consultant_user.id}",
        cookies=tokens.dict(),
    )
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    user = await app_context.repos.users.get_by_id(consultant_user.id)
    assert user is None

    # Check if FK on delete = set null worked out
    person = await app_context.repos.people.get_person(default_person.id)
    assert person.consultant_id is None
