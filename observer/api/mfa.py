from fastapi import APIRouter
from starlette import status

router = APIRouter(prefix="/mfa")


@router.post(
    "/setup",
    status_code=status.HTTP_200_OK,
)
async def setup_mfa():
    """Setup MFA authentication"""
