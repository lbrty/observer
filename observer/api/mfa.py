import base64

from fastapi import APIRouter, Depends
from starlette import status

from observer.api.exceptions import TOTPError, TOTPExistsError
from observer.components.auth import authenticated_user
from observer.components.mfa import mfa_service
from observer.entities.users import User
from observer.schemas.mfa import (
    MFAActivationRequest,
    MFAActivationResponse,
    MFABackupCodesResponse,
)
from observer.services.mfa import MFAServiceInterface
from observer.settings import settings

router = APIRouter(prefix="/mfa")


@router.post(
    "/configure",
    response_model=MFAActivationResponse,
    status_code=status.HTTP_200_OK,
)
async def configure_mfa(
    user: User = Depends(authenticated_user),
    mfa: MFAServiceInterface = Depends(mfa_service),
) -> MFAActivationResponse:
    """Setup MFA authentication"""
    if user.mfa_enabled and user.mfa_encrypted_secret is not None:
        raise TOTPExistsError

    mfa_secret = await mfa.create(settings.title, user.ref_id)
    qr_image = await mfa.into_qr(mfa_secret)
    return MFAActivationResponse(
        secret=mfa_secret.secret,
        qr_image=base64.b64encode(qr_image),
    )


@router.post(
    "/setup",
    response_model=MFABackupCodesResponse,
    status_code=status.HTTP_200_OK,
)
async def setup_mfa(
    activation_request: MFAActivationRequest,
    user: User = Depends(authenticated_user),
    mfa: MFAServiceInterface = Depends(mfa_service),
) -> MFABackupCodesResponse:
    """Save MFA configuration and create backup codes"""
    if user.mfa_enabled and user.mfa_encrypted_secret is not None:
        raise TOTPExistsError

    if await mfa.valid(activation_request.totp_code.get_secret_value(), activation_request.secret.get_secret_value()):
        return MFABackupCodesResponse(backup_codes=[])

    raise TOTPError