import shortuuid
from fastapi import APIRouter, BackgroundTasks, Depends, Response
from starlette import status

from observer.components.auth import authenticated_user
from observer.components.services import (
    audit_service,
    auth_service,
    mailer,
    users_service,
)
from observer.entities.base import SomeUser
from observer.entities.users import User
from observer.services.audit_logs import AuditServiceInterface
from observer.services.auth import AuthServiceInterface
from observer.services.mailer import EmailMessage, MailerInterface
from observer.services.users import UsersServiceInterface
from observer.settings import settings

router = APIRouter(prefix="/account")


@router.get("/confirmation/{code}", status_code=status.HTTP_204_NO_CONTENT)
async def confirm_account(
    tasks: BackgroundTasks,
    code: str,
    user: SomeUser = Depends(authenticated_user),
    audits: AuditServiceInterface = Depends(audit_service),
    auth: AuthServiceInterface = Depends(auth_service),
    users: UsersServiceInterface = Depends(users_service),
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


@router.get("/confirmation/resend", status_code=status.HTTP_204_NO_CONTENT)
async def resend_confirmation(
    tasks: BackgroundTasks,
    user: User = Depends(authenticated_user),
    audits: AuditServiceInterface = Depends(audit_service),
    auth: AuthServiceInterface = Depends(auth_service),
    users: UsersServiceInterface = Depends(users_service),
    mail: MailerInterface = Depends(mailer),
):
    confirmation = await users.create_confirmation(user.id, shortuuid.uuid())
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
