from starlette import status


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
    token_data, _ = await app_context.jwt_service.decode(resp_json["refresh_token"])
    assert token_data.ref_id == consultant_user.ref_id

    audit_log = await app_context.audit_service.find_by_ref(
        f"origin=auth,source=service:auth,action=token:login,ref_id={consultant_user.ref_id}"
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
        "data": None,
        "message": "Wrong email or password",
        "status_code": status.HTTP_401_UNAUTHORIZED,
    }


async def test_token_refresh_works_as_expected(authorized_client, ensure_db, app_context, consultant_user):
    resp = await authorized_client.post("/auth/token/refresh")
    assert resp.status_code == status.HTTP_200_OK
    resp_json = resp.json()
    token_data, _ = await app_context.jwt_service.decode(resp_json["refresh_token"])
    assert token_data.ref_id == consultant_user.ref_id
    audit_log = await app_context.audit_service.find_by_ref(
        f"origin=auth,source=service:auth,action=token:refresh,ref_id={consultant_user.ref_id}"
    )
    assert audit_log.data["ref_id"] == consultant_user.ref_id


async def test_token_refresh_works_as_expected_when_refresh_token_is_invalid(client, ensure_db, app_context):
    resp = await client.post("/auth/token/refresh", cookies=dict(refresh_token="INVALID-TOKEN"))
    assert resp.status_code == status.HTTP_403_FORBIDDEN
    assert resp.json() == {
        "code": "unauthorized",
        "data": None,
        "message": "Invalid refresh token",
        "status_code": 403,
    }

    audit_log = await app_context.audit_service.find_by_ref(
        "origin=auth,source=service:auth,action=token:refresh,kind=error"
    )
    assert audit_log.data == dict(refresh_token="INVALID-TOKEN", notice="invalid refresh token")


async def test_registration_works_as_expected(client, ensure_db, app_context):
    resp = await client.post(
        "/auth/register",
        json=dict(
            email="email@example.com",
            password="!@1StronKPassw0rd#",
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED
    resp_json = resp.json()
    token_data, _ = await app_context.jwt_service.decode(resp_json["refresh_token"])
    user = await app_context.users_service.get_by_email("email@example.com")
    assert token_data.ref_id == user.ref_id

    audit_log = await app_context.audit_service.find_by_ref(
        f"origin=auth,source=service:auth,action=token:register,ref_id={user.ref_id}"
    )
    assert audit_log.data == dict(ref_id=user.ref_id, role=user.role.value)
