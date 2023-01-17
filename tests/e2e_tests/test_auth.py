from datetime import datetime, timedelta, timezone
from http.cookies import SimpleCookie

from sqlalchemy import select, update
from starlette import status

from observer.db.tables.users import password_resets
from observer.entities.users import PasswordReset

SECURE_PASSWORD = "!@1StronKPassw0rd#"


async def test_token_login_works_as_expected(client, ensure_db, app_context, consultant_user):
    resp = await client.post(
        "/auth/token",
        json=dict(
            email=consultant_user.email,
            password="secret",
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    resp_json = resp.json()
    token_data, _ = await app_context.jwt_service.decode(resp_json["access_token"])
    assert token_data.ref_id == consultant_user.ref_id

    cookies = SimpleCookie(resp.headers["set-cookie"])
    token_data, _ = await app_context.jwt_service.decode(cookies["refresh_token"].value)
    assert token_data.ref_id == consultant_user.ref_id

    audit_log = await app_context.audit_service.find_by_ref(
        f"endpoint=token_login,action=token:login,ref_id={consultant_user.ref_id}"
    )
    assert audit_log.data["ref_id"] == consultant_user.ref_id


async def test_token_login_fails_if_credentials_are_wrong(client, ensure_db, consultant_user):
    resp = await client.post(
        "/auth/token",
        json=dict(
            email=consultant_user.email,
            password="wronk passwort",
        ),
    )
    assert resp.status_code == status.HTTP_401_UNAUTHORIZED
    assert resp.json() == {
        "code": "unauthorized",
        "message": "Wrong email or password",
        "status_code": status.HTTP_401_UNAUTHORIZED,
    }


async def test_token_refresh_works_as_expected(
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post("/auth/token/refresh")
    assert resp.status_code == status.HTTP_200_OK

    cookies = SimpleCookie(resp.headers["set-cookie"])
    token_data, _ = await app_context.jwt_service.decode(cookies["refresh_token"].value)
    assert token_data.ref_id == consultant_user.ref_id

    audit_log = await app_context.audit_service.find_by_ref(
        f"endpoint=token_refresh,action=token:refresh,ref_id={consultant_user.ref_id}"
    )
    assert audit_log.data["ref_id"] == consultant_user.ref_id


async def test_token_refresh_works_as_expected_when_refresh_token_is_invalid(client, ensure_db, app_context):
    resp = await client.post("/auth/token/refresh", cookies=dict(refresh_token="INVALID-TOKEN"))
    assert resp.status_code == status.HTTP_403_FORBIDDEN
    assert resp.json() == {
        "code": "unauthorized",
        "message": "Invalid refresh token",
        "status_code": 403,
    }

    audit_log = await app_context.audit_service.find_by_ref("endpoint=token_refresh,action=token:refresh,kind=error")
    assert audit_log.data == dict(refresh_token="INVALID-TOKEN", notice="invalid refresh token")


async def test_registration_works_as_expected(client, ensure_db, app_context):
    resp = await client.post(
        "/auth/register",
        json=dict(
            email="email@example.com",
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    cookies = SimpleCookie(resp.headers["set-cookie"])
    token_data, _ = await app_context.jwt_service.decode(cookies["refresh_token"].value)
    user = await app_context.users_service.get_by_email("email@example.com")
    assert token_data.ref_id == user.ref_id

    audit_log = await app_context.audit_service.find_by_ref(
        f"endpoint=token_register,action=token:register,ref_id={user.ref_id}"
    )
    assert audit_log.data == dict(ref_id=user.ref_id, role=user.role.value)


async def test_registration_returns_error_when_system_works_in_invite_only_mode(
    client,
    ensure_db,
    app_context,
    env_settings,
):
    env_settings.invite_only = True
    resp = await client.post(
        "/auth/register",
        json=dict(
            email="email@example.com",
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_400_BAD_REQUEST
    assert resp.json() == {
        "code": "registrations_closed_error",
        "message": "Registrations are not allowed",
        "status_code": 400,
    }


async def test_password_reset_request_works_as_expected(
    client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await client.post("/auth/reset-password", json=dict(email=consultant_user.email))
    assert resp.status_code == status.HTTP_204_NO_CONTENT
    query = select(password_resets).where(password_resets.c.user_id == str(consultant_user.id))

    result = await app_context.db.fetchone(query)
    password_reset = PasswordReset(**result)
    assert password_reset is not None


async def test_password_reset_with_valid_code_works_as_expected(
    client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await client.post("/auth/reset-password", json=dict(email=consultant_user.email))
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    query = select(password_resets).where(password_resets.c.user_id == str(consultant_user.id))
    result = await app_context.db.fetchone(query)
    password_reset = PasswordReset(**result)
    resp = await client.post(
        f"/auth/reset-password/{password_reset.code}",
        json=dict(
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_204_NO_CONTENT
    audit_log = await app_context.audit_service.find_by_ref(
        f"endpoint=reset_password_with_code,action=reset:password,ref_id={consultant_user.ref_id}"
    )
    assert audit_log.data == dict(code=password_reset.code, ref_id=consultant_user.ref_id)


async def test_password_reset_works_as_expected_when_expired_reset_code_is_used(
    client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await client.post("/auth/reset-password", json=dict(email=consultant_user.email))
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    query = (
        update(password_resets)
        .values(dict(created_at=datetime.now(tz=timezone.utc) - timedelta(hours=2)))
        .where(password_resets.c.user_id == str(consultant_user.id))
        .returning("*")
    )
    result = await app_context.db.fetchone(query)
    password_reset = PasswordReset(**result)
    resp = await client.post(
        f"/auth/reset-password/{password_reset.code}",
        json=dict(
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_401_UNAUTHORIZED
    assert resp.json() == {
        "code": "password_reset_code_expired_error",
        "message": "password reset code hash expired",
        "status_code": 401,
    }


async def test_password_reset_works_as_expected_when_unknown_reset_code_is_used(
    client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await client.post("/auth/reset-password", json=dict(email=consultant_user.email))
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    resp = await client.post(
        "/auth/reset-password/random-code",
        json=dict(
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {"code": "not_found", "status_code": 404, "message": "not found"}


async def test_password_change_works_as_expected(
    client,
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/auth/change-password",
        json=dict(
            old_password="secret",
            new_password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_204_NO_CONTENT
    audit_log = await app_context.audit_service.find_by_ref(
        f"endpoint=change_password,action=change:password,ref_id={consultant_user.ref_id}"
    )
    assert audit_log.data["ref_id"] == consultant_user.ref_id
    resp = await client.post(
        "/auth/token",
        json=dict(
            email=consultant_user.email,
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    cookies = SimpleCookie(resp.headers["set-cookie"])
    token_data, _ = await app_context.jwt_service.decode(cookies["refresh_token"].value)
    assert token_data.ref_id == consultant_user.ref_id


async def test_password_change_works_as_expected_when_new_password_is_weak(
    client,
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/auth/change-password",
        json=dict(
            old_password="secret",
            new_password="weak-pass",
        ),
    )
    assert resp.status_code == status.HTTP_400_BAD_REQUEST
    assert resp.json() == {
        "code": "weak_password_error",
        "message": "Given password is weak",
        "status_code": 400,
    }


async def test_password_change_works_as_expected_when_old_password_is_invalid(
    client,
    authorized_client,
    ensure_db,
    app_context,
    consultant_user,
):
    resp = await authorized_client.post(
        "/auth/change-password",
        json=dict(
            old_password="invalid-secret",
            new_password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_403_FORBIDDEN
    assert resp.json() == {
        "code": "invalid_password_error",
        "message": "Invalid password",
        "status_code": 403,
    }
