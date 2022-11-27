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
        "status_code": 401,
    }


async def test_token_refresh(client, ensure_db, consultant_user):
    pass


async def test_invalid_token_results_in_http_403(client):
    pass
