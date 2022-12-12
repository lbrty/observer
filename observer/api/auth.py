from datetime import datetime, timedelta, timezone

from fastapi import APIRouter, BackgroundTasks, Depends, Response
from starlette import status

from observer.api.exceptions import ForbiddenError
from observer.components.auth import authenticated_user, refresh_token_cookie
from observer.components.services import audit_service, auth_service, mailer
from observer.entities.users import User
from observer.schemas.auth import (
    ChangePasswordRequest,
    LoginPayload,
    NewPasswordRequest,
    RegistrationPayload,
    ResetPasswordRequest,
    TokenResponse,
)
from observer.services.audit_logs import AuditServiceInterface
from observer.services.auth import AuthServiceInterface
from observer.services.mailer import EmailMessage, MailerInterface
from observer.settings import settings

router = APIRouter(prefix="/auth")


@router.post(
    "/token",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
)
async def token_login(
    tasks: BackgroundTasks,
    login_payload: LoginPayload,
    audits: AuditServiceInterface = Depends(audit_service),
    auth: AuthServiceInterface = Depends(auth_service),
) -> TokenResponse:
    """Login using email and password"""
    user, auth_token = await auth.token_login(login_payload)

    # Now we need to save login event
    audit_log = await auth.create_log(
        f"action=token:login,ref_id={user.ref_id}",
        timedelta(days=settings.auth_audit_event_login_days),
        data=dict(ref_id=user.ref_id),
    )
    tasks.add_task(audits.add_event, audit_log)
    return auth_token


@router.post(
    "/token/refresh",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
)
async def token_refresh(
    tasks: BackgroundTasks,
    refresh_token: str = Depends(refresh_token_cookie),
    audits: AuditServiceInterface = Depends(audit_service),
    auth: AuthServiceInterface = Depends(auth_service),
) -> TokenResponse:
    """Refresh access token using refresh token"""
    try:
        token_data, result = await auth.refresh_token(refresh_token)
        audit_log = await auth.create_log(
            f"action=token:refresh,ref_id={token_data.ref_id}",
            timedelta(days=settings.auth_audit_event_refresh_days),
            data=dict(ref_id=token_data.ref_id),
        )
        tasks.add_task(audits.add_event, audit_log)
        return result
    except ForbiddenError:
        # Since it is an exception we need to create audit log synchronously
        audit_log = await auth.create_log(
            "action=token:refresh,kind=error",
            timedelta(days=settings.auth_audit_event_lifetime_days),
            data=dict(refresh_token=refresh_token, notice="invalid refresh token"),
        )
        await audits.add_event(audit_log)
        raise


@router.post(
    "/register",
    response_model=TokenResponse,
    status_code=status.HTTP_201_CREATED,
)
async def token_register(
    tasks: BackgroundTasks,
    registration_payload: RegistrationPayload,
    audits: AuditServiceInterface = Depends(audit_service),
    auth: AuthServiceInterface = Depends(auth_service),
) -> TokenResponse:
    """Register using email and password"""
    user, token_response = await auth.register(registration_payload)
    audit_log = await auth.create_log(
        f"action=token:register,ref_id={user.ref_id}",
        None,
        data=dict(ref_id=user.ref_id, role=user.role.value),
    )
    tasks.add_task(audits.add_event, audit_log)
    return token_response


@router.post(
    "/change-password",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def change_password(
    tasks: BackgroundTasks,
    change_password_payload: ChangePasswordRequest,
    audits: AuditServiceInterface = Depends(audit_service),
    user: User = Depends(authenticated_user),
    auth: AuthServiceInterface = Depends(auth_service),
    mail: MailerInterface = Depends(mailer),
) -> Response:
    """Change password for user"""
    await auth.change_password(user.id, change_password_payload)
    tasks.add_task(
        mail.send,
        EmailMessage(
            to_email=user.email,
            from_email=settings.from_email,
            subject=settings.auth_password_change_subject,
            body=(
                "Your account password has been updated "
                f"at {datetime.now(tz=timezone.utc).strftime('%m/%d/%Y, %H:%M:%S')}."
            ),
        ),
    )
    audit_log = await auth.create_log(
        f"endpoint=change_password,action=change_password:request,ref_id={user.ref_id}",
        timedelta(days=settings.auth_audit_event_lifetime_days),
        data=dict(
            ref_id=user.ref_id,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/reset-password",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def reset_password_request(
    tasks: BackgroundTasks,
    reset_password_payload: ResetPasswordRequest,
    audits: AuditServiceInterface = Depends(audit_service),
    auth: AuthServiceInterface = Depends(auth_service),
    mail: MailerInterface = Depends(mailer),
) -> Response:
    """Reset password for user using email"""
    user, password_reset = await auth.reset_password_request(reset_password_payload.email)
    reset_link = f"{settings.app_domain}/{settings.password_reset_url.format(code=password_reset.code)}"
    tasks.add_task(
        mail.send,
        EmailMessage(
            to_email=user.email,
            from_email=settings.from_email,
            subject=settings.mfa_reset_subject,
            body=f"To reset you password please use the following link {reset_link}.",
        ),
    )
    audit_log = await auth.create_log(
        f"action=reset_password:request,ref_id={user.ref_id}",
        timedelta(days=settings.auth_audit_event_lifetime_days),
        data=dict(
            email=reset_password_payload.email,
            ref_id=user.ref_id,
            role=user.role.value,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/reset-password/{code}",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def reset_password_with_code(
    tasks: BackgroundTasks,
    code: str,
    new_password_payload: NewPasswordRequest,
    audits: AuditServiceInterface = Depends(audit_service),
    auth: AuthServiceInterface = Depends(auth_service),
    mail: MailerInterface = Depends(mailer),
) -> Response:
    """Reset password using reset code"""
    user = await auth.reset_password_with_code(code, new_password_payload.password.get_secret_value())
    tasks.add_task(
        mail.send,
        EmailMessage(
            to_email=user.email,
            from_email=settings.from_email,
            subject=settings.mfa_reset_subject,
            body=f"You password has been reset at {datetime.now(tz=timezone.utc).strftime('%m/%d/%Y, %H:%M:%S')}.",
        ),
    )
    audit_log = await auth.create_log(
        f"endpoint=reset_password_with_code,action=reset:password,ref_id={user.ref_id}",
        timedelta(days=settings.auth_audit_event_lifetime_days),
        data=dict(
            code=code,
            ref_id=user.ref_id,
        ),
    )
    tasks.add_task(audits.add_event, audit_log)
    return Response(status_code=status.HTTP_204_NO_CONTENT)
