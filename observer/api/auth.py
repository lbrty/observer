from fastapi import APIRouter

from observer.context import ctx
from observer.schemas.auth import LoginPayload, TokenResponse
from observer.schemas.health import HealthResponse

router = APIRouter(prefix="/auth")


@router.post("/login", response_model=HealthResponse)
async def login(login_payload: LoginPayload) -> TokenResponse:
    """Login using email and password"""
    result = await ctx.auth_service.token_login(login_payload)
    return result
