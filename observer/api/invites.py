from datetime import datetime, timezone

from fastapi import APIRouter, BackgroundTasks, Depends
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.api.exceptions import WeakPasswordError
from observer.common import bcrypt
from observer.common.bcrypt import is_strong_password
from observer.common.exceptions import get_api_errors
from observer.components.audit import Props, Tracked, client_ip
from observer.components.services import (
    audit_service,
    auth_service,
    mailer,
    users_service,
)
from observer.schemas.auth import LoginPayload, TokenResponse
from observer.schemas.users import InviteJoinRequest
from observer.services.audit_logs import IAuditService
from observer.services.auth import IAuthService
from observer.services.mailer import EmailMessage, IMailer
from observer.services.users import IUsersService
from observer.settings import settings

router = APIRouter(prefix="/invites")


@router.post(
    "/join/{code}",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
    responses=get_api_errors(
        status.HTTP_404_NOT_FOUND,
        status.HTTP_409_CONFLICT,
    ),
    tags=["invites"],
)
async def join_with_invite(
    tasks: BackgroundTasks,
    code: str,
    join_request: InviteJoinRequest,
    users: IUsersService = Depends(users_service),
    audits: IAuditService = Depends(audit_service),
    auth: IAuthService = Depends(auth_service),
    ip_address: str = Depends(client_ip),
    mail: IMailer = Depends(mailer),
    props: Props = Depends(
        Tracked(
            tag="endpoint=join_with_invite",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> TokenResponse:
    """Once users receive invites we want them to hit this endpoint with new account password.
    Given invite is valid and password is strong enough we do the following things:

        1. Hash password and update user's password,
        2. Confirm user automatically,
        3. Issue session token and refresh token (aka authenticate users).
    """
    invite = await users.get_invite(code)
    user = await users.get_by_id(invite.user_id)
    if not is_strong_password(join_request.password.get_secret_value(), settings.password_policy):
        raise WeakPasswordError(message="Given password is weak")

    password_hash = bcrypt.hash_password(join_request.password.get_secret_value())
    await users.update_password(user.id, password_hash)
    tasks.add_task(
        mail.send,
        EmailMessage(
            to_email=user.email,
            from_email=settings.from_email,
            subject=settings.auth_password_change_subject,
            body=f"Your password has been updated at {datetime.now(tz=timezone.utc).strftime('%m/%d/%Y, %H:%M:%S')}.",
        ),
    )
    await users.just_confirm_user(user.id)
    await users.delete_invite(code)
    audit_log = props.new_event(
        f"action=delete:invite,ref_id={user.ref_id}",
        data=dict(
            invite=jsonable_encoder(invite),
            ip_address=ip_address,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    user, auth_token = await auth.token_login(
        LoginPayload(
            email=user.email,
            password=join_request.password,
        )
    )

    # Now we need to save login event
    audit_log = props.new_event(
        f"action=token:login,source=invite,ref_id={user.ref_id}",
        data=dict(ref_id=user.ref_id, ip_address=ip_address),
    )
    tasks.add_task(audits.add_event, audit_log)
    return auth_token
