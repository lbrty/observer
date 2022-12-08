import base64

from fastapi import APIRouter, Depends
from starlette import status

from observer.api.exceptions import TOTPError
from observer.components.mfa import mfa_service, user_with_no_mfa
from observer.components.services import crypto_service, keychain, users_service
from observer.entities.users import User
from observer.schemas.mfa import (
    MFAActivationRequest,
    MFAActivationResponse,
    MFABackupCodesResponse,
)
from observer.schemas.users import UserMFAUpdateRequest
from observer.services.crypto import CryptoServiceInterface
from observer.services.keys import Keychain
from observer.services.mfa import MFAServiceInterface
from observer.services.users import UsersServiceInterface
from observer.settings import settings

router = APIRouter(prefix="/mfa")


@router.post(
    "/configure",
    response_model=MFAActivationResponse,
    status_code=status.HTTP_200_OK,
)
async def configure_mfa(
    user: User = Depends(user_with_no_mfa),
    mfa: MFAServiceInterface = Depends(mfa_service),
) -> MFAActivationResponse:
    """Setup MFA authentication"""
    mfa_secret = await mfa.create(settings.title, user.ref_id)
    qr_image = await mfa.into_qr(mfa_secret)
    return MFAActivationResponse(
        secret=mfa_secret.secret,
        qr_image=base64.b64encode(qr_image),
    )


@router.post(
    "/setup",
    response_model=MFABackupCodesResponse,
    status_code=status.HTTP_201_CREATED,
)
async def setup_mfa(
    activation_request: MFAActivationRequest,
    user: User = Depends(user_with_no_mfa),
    mfa: MFAServiceInterface = Depends(mfa_service),
    crypto: CryptoServiceInterface = Depends(crypto_service),
    user_service: UsersServiceInterface = Depends(users_service),
    key_chain: Keychain = Depends(keychain),
) -> MFABackupCodesResponse:
    """Save MFA configuration and create backup codes"""
    if await mfa.valid(activation_request.totp_code.get_secret_value(), activation_request.secret.get_secret_value()):
        key_hash = key_chain.keys[0].hash
        backup_codes = await mfa.create_backup_codes(settings.num_backup_codes)
        encrypted_secret = crypto.encrypt(key_hash, activation_request.secret.get_secret_value().encode())
        encrypted_backup_codes = crypto.encrypt(key_hash, ",".join(backup_codes).encode())
        mfa_update_request = UserMFAUpdateRequest(
            mfa_enabled=True,
            mfa_encrypted_secret=encrypted_secret,
            mfa_encrypted_backup_codes=encrypted_backup_codes,
        )
        await user_service.update_mfa(
            user.id,
            mfa_update_request,
        )
        return MFABackupCodesResponse(backup_codes=[])

    raise TOTPError
