from fastapi import APIRouter, BackgroundTasks, Depends, Query, Response
from fastapi.encoders import jsonable_encoder
from starlette import status

from observer.common.exceptions import get_api_errors
from observer.common.types import Identifier, Role
from observer.components.audit import Props, Tracked, client_ip
from observer.components.auth import RequiresRoles
from observer.components.services import (
    audit_service,
    auth_service,
    mailer,
    users_service,
)
from observer.entities.base import SomeUser
from observer.schemas.users import NewUserRequest, UserInviteRequest, UserInviteResponse
from observer.services.audit_logs import IAuditService
from observer.services.auth import IAuthService
from observer.services.mailer import EmailMessage, IMailer
from observer.services.users import IUsersService
from observer.settings import settings

router = APIRouter(prefix="/invites")


@router.post(
    "",
    response_model=UserInviteResponse,
    status_code=status.HTTP_201_CREATED,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
        status.HTTP_409_CONFLICT,
    ),
    tags=["admin", "invites"],
)
async def create_invite(
    tasks: BackgroundTasks,
    invite_request: UserInviteRequest,
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.staff]),
    ),
    users: IUsersService = Depends(users_service),
    auth: IAuthService = Depends(auth_service),
    audits: IAuditService = Depends(audit_service),
    ip_address: str = Depends(client_ip),
    mail: IMailer = Depends(mailer),
    props: Props = Depends(
        Tracked(
            tag="endpoint=create_invite,action=create:invite",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> UserInviteResponse:
    random_password = await auth.gen_password()
    new_user = await users.create_user(
        NewUserRequest(
            email=invite_request.email,
            role=invite_request.role,
            password=random_password,
        ),
        is_active=False,
    )
    invite = await users.create_invite(new_user.id)
    link = f"https://{settings.app_domain}{settings.invite_url.format(code=invite.code)}"
    tasks.add_task(
        mail.send,
        EmailMessage(
            to_email=new_user.email,
            from_email=settings.from_email,
            subject=settings.mfa_reset_subject,
            body=f"You are invited to join {settings.invite_subject} please use the following {link}",
        ),
    )
    audit_log = props.new_event(
        f"action=send:invite,ref_id={user.ref_id}",
        data=dict(ref_id=new_user.ref_id),
    )
    tasks.add_task(audits.add_event, audit_log)
    audit_log = props.new_event(
        f"ref_id={user.ref_id}",
        data=dict(
            ref_id=new_user.ref_id,
            role=new_user.role.value,
            ip_address=ip_address,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return UserInviteResponse(**invite.dict())


@router.delete(
    "/{code}",
    status_code=status.HTTP_204_NO_CONTENT,
    responses=get_api_errors(
        status.HTTP_401_UNAUTHORIZED,
        status.HTTP_403_FORBIDDEN,
        status.HTTP_404_NOT_FOUND,
    ),
    tags=["admin", "invites"],
)
async def delete_invite(
    code: Identifier,
    tasks: BackgroundTasks,
    delete_user: bool = Query(False, description="If set then related user will be deleted"),
    user: SomeUser = Depends(
        RequiresRoles([Role.admin, Role.staff]),
    ),
    users: IUsersService = Depends(users_service),
    audits: IAuditService = Depends(audit_service),
    ip_address: str = Depends(client_ip),
    props: Props = Depends(
        Tracked(
            tag="endpoint=delete_invite",
            expires_in=None,
        ),
        use_cache=False,
    ),
) -> Response:
    invite = await users.get_invite(code, validate=False)
    await users.delete_invite(code)
    if delete_user:
        deleted_user = await users.delete_user(user.id)
        audit_log = props.new_event(
            f"action=delete:user,ref_id={user.ref_id}",
            data=dict(
                invite=jsonable_encoder(deleted_user),
                ip_address=ip_address,
            ),
        )
        tasks.add_task(audits.add_event, audit_log)

    audit_log = props.new_event(
        f"action=delete:invite,ref_id={user.ref_id}",
        data=dict(
            invite=jsonable_encoder(invite),
            ip_address=ip_address,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
