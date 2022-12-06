from starlette import status


async def test_mfa_setup_request(client, ensure_db, app_context, consultant_user):
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
        "/mfa/configure",
        json=dict(
            email=consultant_user.email,
            password="secret",
        ),
        cookies=resp_json,
    )
    assert resp.status_code == status.HTTP_200_OK

    resp_json = resp.json()
    assert "secret" in resp_json
    assert "qr_image" in resp_json
