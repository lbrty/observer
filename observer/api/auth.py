from fastapi import APIRouter, Depends, Response
from starlette import status

from observer.components.auth import refresh_token_cookie
from observer.context import ctx
from observer.schemas.auth import (
    ChangePasswordRequest,
    LoginPayload,
    NewPasswordRequest,
    RegistrationPayload,
    ResetPasswordRequest,
    TokenResponse,
)

router = APIRouter(prefix="/auth")


@router.post(
    "/token",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
)
async def token_login(login_payload: LoginPayload) -> TokenResponse:
    """Login using email and password"""
    result = await ctx.auth_service.token_login(login_payload)
    return result


@router.post(
    "/token/refresh",
    response_model=TokenResponse,
    status_code=status.HTTP_200_OK,
)
async def token_refresh(refresh_token: str = Depends(refresh_token_cookie)) -> TokenResponse:
    """Refresh access token using refresh token"""
    result = await ctx.auth_service.refresh_token(refresh_token)
    return result


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
async def reset_password(change_password_payload: ChangePasswordRequest) -> Response:
    """Reset password for user using email"""
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/reset-password",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def reset_password(reset_password_payload: ResetPasswordRequest) -> Response:
    """Reset password for user using email"""
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post(
    "/reset-password/{code}",
    status_code=status.HTTP_204_NO_CONTENT,
)
async def reset_password(code: str, new_password_payload: NewPasswordRequest) -> Response:
    """Reset password using reset code"""
    return Response(status_code=status.HTTP_204_NO_CONTENT)
