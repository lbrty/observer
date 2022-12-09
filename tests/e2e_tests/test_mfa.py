from pyotp import TOTP
from starlette import status


async def test_mfa_configure_request(authorized_client, ensure_db, app_context, consultant_user):
    resp = await authorized_client.post(
        "/mfa/configure",
        json=dict(
            email=consultant_user.email,
            password="secret",
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    resp_json = resp.json()
    assert "secret" in resp_json
    assert "qr_image" in resp_json


async def test_mfa_setup_request(authorized_client, ensure_db, app_context, consultant_user):
    resp = await authorized_client.post(
        "/mfa/configure",
        json=dict(
            email=consultant_user.email,
            password="secret",
        ),
    )
    assert resp.status_code == status.HTTP_200_OK
    resp_json = resp.json()
    assert "secret" in resp_json
    secret = resp_json["secret"]
    totp = TOTP(secret)

    resp = await authorized_client.post(
        "/mfa/setup",
        json=dict(
            totp_code=totp.now(),
            secret=secret,
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    resp_json = resp.json()
    assert len(resp_json["backup_codes"]) == 6
