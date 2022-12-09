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


async def test_configured_mfa_requires_totp_code_during_login(
    authorized_client,
    client,
    ensure_db,
    app_context,
    consultant_user,
):
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

    resp = await client.post(
        "/auth/token",
        json=dict(
            email=consultant_user.email,
            password="secret",
        ),
    )
    assert resp.status_code == status.HTTP_417_EXPECTATION_FAILED
    assert resp.json() == {
        "code": "totp_required_error",
        "data": None,
        "message": "totp required",
        "status_code": 417,
    }


async def test_configured_mfa_works_as_expected_when_correct_credentials_given(
    authorized_client,
    client,
    ensure_db,
    app_context,
    consultant_user,
):
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

    resp = await client.post(
        "/auth/token",
        json=dict(
            email=consultant_user.email,
            password="secret",
            totp_code=totp.now(),
        ),
    )
    assert resp.status_code == status.HTTP_200_OK

    resp_json = resp.json()
    token_data, _ = await app_context.jwt_service.decode(resp_json["refresh_token"])
    assert token_data.ref_id == consultant_user.ref_id


async def test_configured_mfa_works_as_expected_when_incorrect_totp_code_given(
    authorized_client,
    client,
    ensure_db,
    app_context,
    consultant_user,
):
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

    resp = await client.post(
        "/auth/token",
        json=dict(
            email=consultant_user.email,
            password="secret",
            totp_code="wrong_totp",
        ),
    )
    assert resp.status_code == status.HTTP_401_UNAUTHORIZED
    assert resp.json() == {
        "code": "totp_error",
        "data": None,
        "message": "invalid totp code",
        "status_code": 401,
    }


async def test_mfa_reset_request_works_as_expected(authorized_client, client, ensure_db, app_context, consultant_user):
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
    backup_codes = resp_json["backup_codes"]
    assert len(backup_codes) == 6

    resp = await client.post(
        "/mfa/reset",
        json=dict(
            email=consultant_user.email,
            backup_code=backup_codes[0],
        ),
    )
    assert resp.status_code == status.HTTP_204_NO_CONTENT
    user = await app_context.users_repo.get_by_id(consultant_user.id)
    assert user.mfa_enabled is False
    assert user.mfa_encrypted_secret is None
    assert user.mfa_encrypted_backup_codes is None
