import re
from datetime import datetime, timedelta, timezone

from sqlalchemy import update
from starlette import status

from observer.db.tables.users import confirmations

SECURE_PASSWORD = "!@1StronKPassw0rd#Accounts"


async def test_registration_creates_account_confirmation(
    client,
    ensure_db,
    app_context,
    clean_mailbox,
):
    app_context.mailer.messages = []
    resp = await client.post(
        "/auth/register",
        json=dict(
            email="emails@examples.com",
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    message = app_context.mailer.messages[0]
    [confirmation_link] = re.findall(r"https?://\S+", message.body)
    parts = confirmation_link.split("/")
    await app_context.users_service.get_confirmation(parts[-1])


async def test_it_is_possible_to_resend_account_confirmation_email(
    client,
    ensure_db,
    app_context,
    clean_mailbox,
):
    app_context.mailer.messages = []
    resp = await client.post(
        "/auth/register",
        json=dict(
            email="emails@examples.com",
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED
    resp_json = resp.json()
    message = app_context.mailer.messages[0]
    [confirmation_link] = re.findall(r"https?://\S+", message.body)
    parts = confirmation_link.split("/")
    confirmation = await app_context.users_service.get_confirmation(parts[-1])
    # assert confirmation is not None
    resp = await client.post(
        "/account/confirmation/resend",
        cookies=resp_json,
    )
    assert resp.status_code == status.HTTP_204_NO_CONTENT
    assert len(app_context.mailer.messages) == 2
    message = app_context.mailer.messages[1]
    [confirmation_link] = re.findall(r"https?://\S+", message.body)
    parts = confirmation_link.split("/")
    recent_confirmation = await app_context.users_service.get_confirmation(parts[-1])
    assert confirmation is not None
    assert confirmation.expires_at < recent_confirmation.expires_at


async def test_account_confirmation_works_as_expected_for_same_authenticated_user(
    client,
    ensure_db,
    app_context,
    clean_mailbox,
):
    app_context.mailer.messages = []
    resp = await client.post(
        "/auth/register",
        json=dict(
            email="emails@examples.com",
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    message = app_context.mailer.messages[0]
    [confirmation_link] = re.findall(r"https?://\S+", message.body)
    parts = confirmation_link.split("/")
    confirmation = await app_context.users_service.get_confirmation(parts[-1])

    assert confirmation is not None
    resp = await client.get(f"/account/confirm/{confirmation.code}", cookies=resp.json())
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    user = await app_context.users_service.get_by_id(confirmation.user_id)
    assert user.is_confirmed


async def test_account_confirmation_works_as_expected_if_user_is_not_authenticated(
    client,
    ensure_db,
    app_context,
    clean_mailbox,
):
    app_context.mailer.messages = []
    resp = await client.post(
        "/auth/register",
        json=dict(
            email="emails@examples.com",
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    message = app_context.mailer.messages[0]
    [confirmation_link] = re.findall(r"https?://\S+", message.body)
    parts = confirmation_link.split("/")
    confirmation = await app_context.users_service.get_confirmation(parts[-1])

    assert confirmation is not None
    resp = await client.get(
        f"/account/confirm/{confirmation.code}",
    )
    assert resp.status_code == status.HTTP_204_NO_CONTENT

    user = await app_context.users_service.get_by_id(confirmation.user_id)
    assert user.is_confirmed


async def test_account_confirmation_works_as_expected_when_other_authenticated_user_tries_to_confirm_other_account(
    client,
    authorized_client,
    ensure_db,
    app_context,
    clean_mailbox,
):
    resp = await client.post(
        "/auth/register",
        json=dict(
            email="emails@examples.com",
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    message = app_context.mailer.messages[0]
    [confirmation_link] = re.findall(r"https?://\S+", message.body)
    parts = confirmation_link.split("/")
    confirmation = await app_context.users_service.get_confirmation(parts[-1])
    assert confirmation is not None

    resp = await authorized_client.get(
        f"/account/confirm/{confirmation.code}",
    )
    assert resp.status_code == status.HTTP_404_NOT_FOUND
    assert resp.json() == {
        "code": "not_found",
        "data": None,
        "message": "Confirmation code not found",
        "status_code": 404,
    }
    user = await app_context.users_service.get_by_id(confirmation.user_id)
    assert user.is_confirmed is False


async def test_account_confirmation_works_as_expected_when_it_has_expired(
    client,
    ensure_db,
    app_context,
    clean_mailbox,
):
    resp = await client.post(
        "/auth/register",
        json=dict(
            email="emails@examples.com",
            password=SECURE_PASSWORD,
        ),
    )
    assert resp.status_code == status.HTTP_201_CREATED

    message = app_context.mailer.messages[0]
    [confirmation_link] = re.findall(r"https?://\S+", message.body)
    parts = confirmation_link.split("/")
    confirmation = await app_context.users_service.get_confirmation(parts[-1])
    assert confirmation is not None

    query = update(confirmations).values(expires_at=datetime.now(tz=timezone.utc) - timedelta(hours=10)).returning("*")
    await app_context.db.execute(query)
    resp = await client.get(
        f"/account/confirm/{confirmation.code}",
    )
    assert resp.status_code == status.HTTP_409_CONFLICT
    assert resp.json() == {
        "code": "confirmation_code_expired_error",
        "status_code": 409,
        "message": "Confirmation code has already expired",
        "data": None,
    }
    user = await app_context.users_service.get_by_id(confirmation.user_id)
    assert user.is_confirmed is False
