from io import BytesIO
from typing import Generator

from observer.entities.documents import Document
from observer.services.crypto import AESCipherOptions, ICryptoService
from observer.services.storage import IStorage

CHUNK_SIZE = 512


class DownloadHandler:
    def __init__(self, storage: IStorage, crypto: ICryptoService):
        self.storage = storage
        self.crypto = crypto

    async def stream(self, document: Document) -> Generator:
        fd = await self.storage.open(document.path)
        secrets = await self.get_encryption_secrets(document)
        encrypted_contents = await fd.read()
        decrypted_contents = BytesIO(
            await self.crypto.aes_decrypt(
                secrets.secret,
                secrets.iv,
                encrypted_contents,
            )
        )
        while (chunk := decrypted_contents.read(CHUNK_SIZE)) is not None:
            yield chunk

    async def get_encryption_secrets(self, document: Document) -> AESCipherOptions:
        key_hash, secrets = document.encryption_key.split(":", maxsplit=1)
        decrypted_secrets = await self.crypto.decrypt(key_hash, secrets.encode())
        encoded_secrets = decrypted_secrets.decode()
        return await self.crypto.parse_aes_secrets(encoded_secrets)
