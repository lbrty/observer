from fastapi import APIRouter

from observer.context import ctx
from observer.schemas.auth import LoginPayload, TokenResponse

router = APIRouter(prefix="/auth")


@router.post("/token", response_model=TokenResponse)
async def token_login(login_payload: LoginPayload) -> TokenResponse:
    """Login using email and password"""
    result = await ctx.auth_service.token_login(login_payload)
    return result
