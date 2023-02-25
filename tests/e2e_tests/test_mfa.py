from http.cookies import SimpleCookie

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

    cookies = SimpleCookie(resp.headers["set-cookie"])
    token_data, _ = await app_context.jwt_service.decode(cookies["refresh_token"].value)
    assert token_data.user_id == str(consultant_user.id)


async def test_configure_mfa_works_as_expected_when_incorrect_totp_code_given(
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
    assert resp.status_code == status.HTTP_400_BAD_REQUEST
    assert resp.json() == {
        "code": "totp_error",
        "message": "Invalid totp code",
        "status_code": 400,
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
    user = await app_context.repos.users.get_by_id(consultant_user.id)
    assert user.mfa_enabled is False
    assert user.mfa_encrypted_secret is None
    assert user.mfa_encrypted_backup_codes is None


async def test_mfa_reset_request_works_as_expected_when_invalid_backup_code_given(
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
        "/mfa/reset",
        json=dict(
            email=consultant_user.email,
            backup_code="wronk-backup-kode",
        ),
    )
    assert resp.status_code == status.HTTP_401_UNAUTHORIZED
    assert resp.json() == {
        "code": "totp_invalid_backup_code_error",
        "message": "Invalid backup code",
        "status_code": 401,
    }


async def test_mfa_reset_request_works_as_expected_when_random_email_given(
    client,
    ensure_db,
    app_context,
):
    payload = dict(
        email="some-random@email.com",
        backup_code="hacking-attempt",
    )
    resp = await client.post(
        "/mfa/reset",
        json=payload,
    )
    assert resp.status_code == status.HTTP_204_NO_CONTENT
    audit_log = await app_context.audit_service.find_by_ref(
        "endpoint=reset_mfa,action=reset:mfa,kind=error",
    )
    assert audit_log.data == payload
