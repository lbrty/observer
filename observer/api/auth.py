from fastapi import APIRouter, Depends
from starlette import status

from observer.components.auth import refresh_token_cookie
from observer.context import ctx
from observer.schemas.auth import LoginPayload, RegistrationPayload, TokenResponse

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
