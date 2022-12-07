from starlette import status

from tests.helpers.auth import get_auth_tokens


async def test_mfa_configure_request(client, ensure_db, app_context, consultant_user):
    auth_token = await get_auth_tokens(app_context, consultant_user)
    resp = await client.post(
        "/mfa/configure",
        json=dict(
            email=consultant_user.email,
            password="secret",
        ),
        cookies=dict(
            access_token=auth_token.access_token,
            refresh_token=auth_token.refresh_token,
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    resp_json = resp.json()
    assert "secret" in resp_json
    assert "qr_image" in resp_json


async def test_mfa_setup_request(client, ensure_db, app_context, consultant_user):
    auth_token = await get_auth_tokens(app_context, consultant_user)
    resp = await client.post(
        "/mfa/configure",
        json=dict(
            email=consultant_user.email,
            password="secret",
        ),
        cookies=dict(
            access_token=auth_token.access_token,
            refresh_token=auth_token.refresh_token,
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    resp.json()
