import base64

from fastapi import APIRouter, Depends
from starlette import status

from observer.components.auth import authenticated_user
from observer.components.mfa import mfa_service
from observer.entities.users import User
from observer.schemas.mfa import MFAActivationResponse
from observer.services.mfa import MFAServiceInterface
from observer.settings import settings

router = APIRouter(prefix="/mfa")


@router.post(
    "/setup",
    response_model=MFAActivationResponse,
    status_code=status.HTTP_200_OK,
)
async def setup_mfa(
    user: User = Depends(authenticated_user),
    mfa: MFAServiceInterface = Depends(mfa_service),
) -> MFAActivationResponse:
    """Setup MFA authentication"""
    mfa_secret = await mfa.create(settings.title, user.ref_id)
    qr_image = await mfa.into_qr(mfa_secret)
    return MFAActivationResponse(secret=mfa_secret.secret, qr_image=base64.b64encode(qr_image))
