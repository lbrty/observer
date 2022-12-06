from typing import List

from pydantic import Field, SecretStr

from observer.schemas.base import SchemaBase


class MFAActivationRequest(SchemaBase):
    secret: SecretStr = Field(..., description="TOTP secret")
    totp_code: SecretStr = Field(..., description="TOTP code to validate and save for user")


class MFAActivationResponse(SchemaBase):
    secret: str = Field(..., description="TOTP secret")
    qr_image: str = Field(..., description="Base64 QR code image")


class MFABackupCodesResponse(SchemaBase):
    backup_codes: List[str] = Field(..., description="List of backup codes")
