from datetime import timedelta

from fastapi import APIRouter, BackgroundTasks, Depends, Request, Response
from starlette import status

from observer.api.exceptions import ForbiddenError
from observer.components.auth import refresh_token_cookie
from observer.components.services import audit_service, auth_service
from observer.context import ctx
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
    audit_logs: AuditServiceInterface = Depends(audit_service),
) -> TokenResponse:
    """Login using email and password"""
    user, auth_token = await ctx.auth_service.token_login(login_payload)

    # Now we need to save login event
    audit_log = await ctx.auth_service.create_log(
        f"action=token:login,ref_id={user.ref_id}",
        timedelta(days=settings.auth_audit_event_login_days),
        data=dict(ref_id=user.ref_id),
    )
    tasks.add_task(audit_logs.add_event, audit_log)
    return auth_token


@router.post(
    "/token/refresh",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
)
async def token_refresh(
    tasks: BackgroundTasks,
    refresh_token: str = Depends(refresh_token_cookie),
    audit_logs: AuditServiceInterface = Depends(audit_service),
) -> TokenResponse:
    """Refresh access token using refresh token"""
    try:
        token_data, result = await ctx.auth_service.refresh_token(refresh_token)
        audit_log = await ctx.auth_service.create_log(
            f"action=token:refresh,ref_id={token_data.ref_id}",
            timedelta(days=settings.auth_audit_event_refresh_days),
            data=dict(ref_id=token_data.ref_id),
        )
        tasks.add_task(audit_logs.add_event, audit_log)
        return result
    except ForbiddenError:
        audit_log = await ctx.auth_service.create_log(
            f"{ctx.auth_service.tag},action=token:refresh,kind=error",
            timedelta(days=settings.auth_audit_event_lifetime_days),
            data=dict(refresh_token=refresh_token, notice="invalid refresh token"),
        )
        tasks.add_task(audit_logs.add_event, audit_log)
        raise


@router.post(
    "/register",
    response_model=TokenResponse,
    status_code=status.HTTP_201_CREATED,
)
async def token_register(registration_payload: RegistrationPayload) -> TokenResponse:
    """Register using email and password"""
    result = await ctx.auth_service.register(registration_payload)
    return result


@router.post(
    "/change-password",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def change_password(change_password_payload: ChangePasswordRequest) -> Response:
    """Reset password for user using email"""
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/reset-password",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def reset_password_request(
    request: Request,
    reset_password_payload: ResetPasswordRequest,
    auth: AuthServiceInterface = Depends(auth_service),
) -> Response:
    """Reset password for user using email"""
    await auth.reset_password(
        reset_password_payload.email,
        metadata=dict(
            host=request.client.host,
        ),
    )
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/reset-password/{code}",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def reset_password_with_code(code: str, new_password_payload: NewPasswordRequest) -> Response:
    """Reset password using reset code"""
    return Response(status_code=status.HTTP_204_NO_CONTENT)
