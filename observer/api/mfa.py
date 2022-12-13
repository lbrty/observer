import base64
from datetime import timedelta

from fastapi import APIRouter, BackgroundTasks, Depends, Response
from starlette import status

from observer.api.exceptions import BadRequestError, TOTPError
from observer.components.mfa import mfa_service, user_with_no_mfa
from observer.components.services import audit_service, keychain, mailer, users_service
from observer.entities.users import User
from observer.schemas.mfa import (
    MFAActivationRequest,
    MFAActivationResponse,
    MFAAResetRequest,
    MFABackupCodesResponse,
)
from observer.schemas.users import UserMFAUpdateRequest
from observer.services.audit_logs import AuditServiceInterface
from observer.services.keys import Keychain
from observer.services.mailer import EmailMessage, MailerInterface
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
    tasks: BackgroundTasks,
    user: User = Depends(user_with_no_mfa),
    mfa: MFAServiceInterface = Depends(mfa_service),
    user_service: UsersServiceInterface = Depends(users_service),
    audits: AuditServiceInterface = Depends(audit_service),
    key_chain: Keychain = Depends(keychain),
) -> MFABackupCodesResponse:
    """Save MFA configuration and create backup codes.

    NOTE:
        Only the latest private key is used to encrypt secret and backup codes.
    """
    if await mfa.valid(activation_request.totp_code.get_secret_value(), activation_request.secret.get_secret_value()):
        key_hash = key_chain.keys[0].hash
        mfa_setup_result = await mfa.setup_mfa(
            activation_request.secret.get_secret_value(),
            key_hash,
            settings.num_backup_codes,
        )
        mfa_update_request = UserMFAUpdateRequest(
            mfa_enabled=True,
            mfa_encrypted_secret=mfa_setup_result.encrypted_secret,
            mfa_encrypted_backup_codes=mfa_setup_result.encrypted_backup_codes,
        )
        await user_service.update_mfa(
            user.id,
            mfa_update_request,
        )
        audit_log = await user_service.create_log(
            "endpoint=setup_mfa,action=setup:mfa",
            timedelta(days=settings.mfa_audit_event_lifetime_days),
            dict(setup_totp_code=activation_request.totp_code.get_secret_value()),
        )
        tasks.add_task(audits.add_event, audit_log)
        return MFABackupCodesResponse(backup_codes=list(mfa_setup_result.plain_backup_codes))

    raise TOTPError


@router.post("/reset", status_code=status.HTTP_204_NO_CONTENT)
async def reset_mfa(
    reset_request: MFAAResetRequest,
    tasks: BackgroundTasks,
    user_service: UsersServiceInterface = Depends(users_service),
    audits: AuditServiceInterface = Depends(audit_service),
    mail: MailerInterface = Depends(mailer),
) -> Response:
    """Reset MFA using one of backup codes

    NOTE:
        HTTP 204 returned anyway to prevent user email brute forcing  because we only
        want exact matches to check and reset if given backup code is valid.
    """
    if user := await user_service.get_by_email(reset_request.email):
        if user.mfa_encrypted_backup_codes:
            await user_service.check_backup_code(user.mfa_encrypted_backup_codes, reset_request.backup_code)
            await user_service.reset_mfa(user.id)
            audit_log = await user_service.create_log(
                "endpoint=reset_mfa,action=reset:mfa",
                timedelta(days=settings.mfa_audit_event_lifetime_days),
                reset_request.dict(),
            )
            tasks.add_task(audits.add_event, audit_log)
            tasks.add_task(
                mail.send,
                EmailMessage(
                    to_email=user.email,
                    from_email=settings.from_email,
                    subject=settings.mfa_reset_subject,
                    body="Your MFA was reset, you can login using your credentials.",
                ),
            )
        else:
            raise BadRequestError(message="Backup codes not found")
    else:
        audit_log = await user_service.create_log(
            "endpoint=reset_mfa,action=reset:mfa,kind=error",
            timedelta(days=settings.mfa_audit_event_lifetime_days),
            reset_request.dict(),
        )
        tasks.add_task(audits.add_event, audit_log)

    return Response(status_code=status.HTTP_204_NO_CONTENT)
