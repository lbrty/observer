from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.types import Role
from observer.schemas.users import CreateUserRequest
from tests.helpers.auth import get_auth_tokens


async def test_admins_can_create_users(
    client,
    app_context,
    admin_user,
    default_project,
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
