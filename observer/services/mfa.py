import io
from dataclasses import dataclass
from typing import Protocol

from pyotp import TOTP, random_base32
from qrcode import QRCode, constants
from qrcode.image.styledpil import StyledPilImage
from qrcode.image.styles.colormasks import RadialGradiantColorMask
from qrcode.image.styles.moduledrawers import RoundedModuleDrawer

from observer.common.types import Identifier
from observer.services.crypto import CryptoServiceInterface


@dataclass
class MFASecret:
    secret: str
    totp_uri: str


class MFAServiceInterface(Protocol):
    crypto_service: CryptoServiceInterface
    totp_leeway: int = 0

    async def into_qr(self, secret: MFASecret) -> bytes:
        raise NotImplementedError

    async def create(self, app_name: str, ref_id: Identifier) -> MFASecret:
        raise NotImplementedError

    async def valid(self, totp_code: str, secret: str) -> bool:
        raise NotImplementedError


class MFAService(MFAServiceInterface):
    def __init__(self, totp_leeway: int, crypto_service: CryptoServiceInterface):
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

    async def create(self, app_name: str, ref_id: Identifier) -> MFASecret:
        secret = random_base32()
        totp_uri = TOTP(secret).provisioning_uri(ref_id, app_name)
        return MFASecret(
            secret=secret,
            totp_uri=totp_uri,
        )

    async def valid(self, totp_code: str, secret: str) -> bool:
        totp = TOTP(secret)
        return totp.verify(totp_code, valid_window=self.totp_leeway)
