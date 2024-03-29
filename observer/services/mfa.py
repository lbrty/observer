import io
from dataclasses import dataclass
from hashlib import sha1
from typing import List, Protocol, Set

import shortuuid
from pyotp import TOTP, random_base32
from qrcode import QRCode, constants
from qrcode.image.styledpil import StyledPilImage
from qrcode.image.styles.colormasks import RadialGradiantColorMask
from qrcode.image.styles.moduledrawers import RoundedModuleDrawer

from observer.common.types import Identifier
from observer.services.crypto import ICryptoService


@dataclass
class MFASecret:
    secret: str
    totp_uri: str


@dataclass
class MFASetupResult:
    encrypted_secret: str
    plain_backup_codes: List[str]
    encrypted_backup_codes: str


class IMFAService(Protocol):
    crypto_service: ICryptoService
    totp_leeway: int = 0

    async def into_qr(self, secret: MFASecret) -> bytes:
        raise NotImplementedError

    async def create(self, app_name: str, user_id: Identifier) -> MFASecret:
        raise NotImplementedError

    async def valid(self, totp_code: str, secret: str) -> bool:
        raise NotImplementedError

    async def create_backup_codes(self, how_many: int) -> List[str]:
        raise NotImplementedError

    async def setup_mfa(self, mfa_secret: str, key_hash: str, num_backup_codes: int) -> MFASetupResult:
        raise NotImplementedError


class MFAService(IMFAService):
    def __init__(self, totp_leeway: int, crypto_service: ICryptoService):
        self.totp_leeway = totp_leeway
        self.crypto_service = crypto_service

    async def into_qr(self, secret: MFASecret) -> bytes:
        qr = QRCode(error_correction=constants.ERROR_CORRECT_L)
        qr.add_data(secret.totp_uri)
        image = qr.make_image(
            image_factory=StyledPilImage,
            color_mask=RadialGradiantColorMask(),
            module_drawer=RoundedModuleDrawer(),
        )
        img_bytes = io.BytesIO()
        image.save(img_bytes, format="jpeg")
        return img_bytes.getvalue()

    async def create(self, app_name: str, user_id: Identifier) -> MFASecret:
        secret = random_base32()
        totp_uri = TOTP(secret).provisioning_uri(str(user_id), app_name)
        return MFASecret(
            secret=secret,
            totp_uri=totp_uri,
        )

    async def valid(self, totp_code: str, secret: str) -> bool:
        totp = TOTP(secret)
        return totp.verify(totp_code, valid_window=self.totp_leeway)

    async def create_backup_codes(self, how_many: int) -> List[str]:
        backup_codes: Set[str] = set()
        while True:
            if len(backup_codes) == how_many:
                break

            hashcode = sha1(shortuuid.uuid().encode())
            hexcode = hashcode.hexdigest()
            code = hexcode[-6:]
            if code not in backup_codes:
                backup_codes.add(code)

        return list(backup_codes)

    async def setup_mfa(self, mfa_secret: str, key_hash: str, num_backup_codes: int) -> MFASetupResult:
        """Generate TOTP backup codes and save with secret

        NOTE: Since this is rather short length values RSA public key encryption used.

        Args:
            mfa_secret(str): TOTP secret
            key_hash(str): hash of key to use
            num_backup_codes(int): number of backup keys to generate

        Returns:
            MFASetupResult: list of backup code which user needs to save.
        """
        backup_codes = await self.create_backup_codes(num_backup_codes)
        encrypted_secret = await self.crypto_service.encrypt(key_hash, mfa_secret.encode())
        encrypted_backup_codes = await self.crypto_service.encrypt(key_hash, ",".join(backup_codes).encode())
        return MFASetupResult(
            encrypted_secret=f"{key_hash}:{encrypted_secret.decode()}",
            plain_backup_codes=backup_codes,
            encrypted_backup_codes=f"{key_hash}:{encrypted_backup_codes.decode()}",
        )
