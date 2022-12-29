from fastapi import APIRouter, BackgroundTasks, Depends, Response
from starlette import status

from observer.components.auth import authenticated_user, current_user
from observer.components.services import (
    audit_service,
    auth_service,
    mailer,
    users_service,
)
from observer.entities.base import SomeUser
from observer.entities.users import User
from observer.schemas.users import UserResponse
from observer.services.audit_logs import IAuditService
from observer.services.auth import IAuthService
from observer.services.mailer import EmailMessage, IMailer
from observer.services.users import IUsersService
from observer.settings import settings

router = APIRouter(prefix="/account")


@router.get(
    "/me",
    response_model=UserResponse,
    status_code=status.HTTP_200_OK,
)
async def get_me(
    user: User = Depends(authenticated_user),
) -> UserResponse:
    return UserResponse(**user.dict())


@router.get("/confirm/{code}", status_code=status.HTTP_204_NO_CONTENT)
async def confirm_account(
    tasks: BackgroundTasks,
    code: str,
    user: SomeUser = Depends(current_user),
    audits: IAuditService = Depends(audit_service),
    auth: IAuthService = Depends(auth_service),
    users: IUsersService = Depends(users_service),
):
    user = await users.confirm_user(user.id if user else None, code)
    audit_log = await auth.create_log(
        f"endpoint=confirm_account,action=confirm:account,ref_id={user.ref_id}",
        None,
        data=dict(
            code=code,
            ref_id=user.ref_id,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post("/confirmation/resend", status_code=status.HTTP_204_NO_CONTENT)
async def resend_confirmation(
    tasks: BackgroundTasks,
    user: User = Depends(authenticated_user),
    audits: IAuditService = Depends(audit_service),
    auth: IAuthService = Depends(auth_service),
    users: IUsersService = Depends(users_service),
    mail: IMailer = Depends(mailer),
):
    confirmation = await users.create_confirmation(user.id)
    link = f"{settings.app_domain}{settings.confirmation_url.format(code=confirmation.code)}"
    tasks.add_task(
        mail.send,
        EmailMessage(
            to_email=user.email,
            from_email=settings.from_email,
            subject=settings.mfa_reset_subject,
            body=f"To confirm your email please use the following link {link}",
        ),
    )
    audit_log = await auth.create_log(
        f"endpoint=resend_confirmation,action=resend:confirmation,ref_id={user.ref_id}",
        None,
        data=dict(ref_id=user.ref_id),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
