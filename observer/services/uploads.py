import base64
import hashlib
import os
from dataclasses import dataclass
from typing import Tuple

from fastapi import UploadFile

from observer.api.exceptions import TooLargeDocumentError, UnsupportedDocumentError
from observer.entities.documents import AllowedDocumentTypes
from observer.services.crypto import ICryptoService
from observer.services.storage import IStorage
from observer.settings import settings


@dataclass
class SealedFile:
    path: str
    encryption_key: str
    encrypted_file: bytes


class UploadHandler:
    def __init__(self, storage: IStorage, crypto: ICryptoService):
        self.storage = storage
        self.crypto = crypto

    async def process_upload(self, file: UploadFile, path: str) -> Tuple[int, SealedFile]:
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
        content_size = len(contents)
        if content_size > settings.max_upload_size:
            raise TooLargeDocumentError

        seal_options = await self.crypto.aes_cipher_options(settings.aes_key_bits)
        tag, ciphertext = await self.crypto.aes_encrypt(seal_options, contents)
        seal_options.tag = tag
        encrypted_contents = base64.b64encode(ciphertext)

        # Save file with hashed name and extension
        extension = AllowedDocumentTypes[file.content_type]
        filename = f"{hashlib.md5(file.filename.encode()).hexdigest()}.{extension}"
        file_path = os.path.join(path, filename)
        await self.storage.save(file_path, encrypted_contents)
        full_path = os.path.join(self.storage.root, file_path)

        # Encode and encrypt AES options
        key_hash = self.crypto.keychain.keys[0].hash
        secrets = await self.crypto.format_aes_secrets(seal_options)
        encrypted_secrets = await self.crypto.encrypt(key_hash, secrets.encode())
        encryption_key = f"{key_hash}:{encrypted_secrets.decode()}"
        return (
            content_size,
            SealedFile(
                path=full_path,
                encryption_key=encryption_key,
                encrypted_file=encrypted_contents,
            ),
        )
