import base64
from dataclasses import dataclass

from fastapi import UploadFile

from observer.api.exceptions import TooLargeDocumentError, UnsupportedDocumentError
from observer.entities.documents import AllowedDocumentTypes
from observer.services.crypto import ICryptoService
from observer.settings import settings


@dataclass
class SealedFile:
    encryption_key: str
    encrypted_file: bytes


class UploadHandler:
    def __init__(self, crypto: ICryptoService):
        self.crypto = crypto

    async def process_upload(self, file: UploadFile) -> SealedFile:
        """Validate and encrypt uploaded file.

        To encrypt data we do the following things

            1. Generate AES secret and IV,
            2. Encrypt file contents,
            3. Create pair `secret:iv` and encrypt them using RSA key,
            4. B64 encode AES secrets pair `secret:iv` and encrypted file contents,
            5. B64 encode secrets and encrypted contents in the format `B64(secret:iv):B64(encrypted contents)`
        """
        if file.content_type not in AllowedDocumentTypes:
            raise UnsupportedDocumentError

        contents = await file.read()
        if len(contents) > settings.max_upload_size:
            raise TooLargeDocumentError

        seal_options = await self.crypto.aes_cipher_options(settings.aes_key_bits)
        encrypted_contents = base64.b64encode(
            await self.crypto.aes_encrypt(
                seal_options.secret,
                seal_options.iv,
                contents,
            )
        )
        secret = base64.b64encode(seal_options.secret).decode()
        iv = base64.b64encode(seal_options.iv).decode()
        key_hash = self.crypto.keychain.keys[0].hash
        secrets = await self.crypto.encrypt(key_hash, f"{secret}:{iv}".encode())
        encrypted_secrets = base64.b64encode(secrets).decode()
        encryption_key = f"{key_hash}:{encrypted_secrets}"
        return SealedFile(
            encryption_key=encryption_key,
            encrypted_file=encrypted_contents,
        )
