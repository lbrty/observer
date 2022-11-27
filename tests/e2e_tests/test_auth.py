from starlette import status


async def test_token_login(client, ensure_db, consultant_user):
    resp = await client.post(
        "/auth/token",
        json=dict(
            email=consultant_user.email,
            password="secret",
        ),
    )
    resp_json = resp.json()
    assert len(resp_json) == 2
    assert "access_token" in resp_json
    assert "refresh_token" in resp_json
    assert resp.status_code == status.HTTP_200_OK


async def test_token_refresh(client, ensure_db, consultant_user):
    pass


async def test_invalid_token_results_in_http_403(client):
    pass
