from starlette import status


async def test_token_login(client, ensure_db, app_context, consultant_user):
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


async def test_token_refresh_works_as_expected(client, ensure_db, app_context, consultant_user):
    resp = await client.post(
        "/auth/token",
        json=dict(
            email=consultant_user.email,
            password="secret",
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    resp_json = resp.json()
    resp = await client.post(
        "/auth/token/refresh",
        cookies={"refresh_token": resp_json["refresh_token"]},
    )
    assert resp.status_code == status.HTTP_200_OK
    resp_json = resp.json()
    token_data, _ = await app_context.jwt_service.decode(resp_json["refresh_token"])
    assert token_data.ref_id == consultant_user.ref_id


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
