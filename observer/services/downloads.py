from io import BytesIO
from typing import AsyncGenerator

from observer.entities.documents import Document
from observer.services.crypto import AESCipherOptions, ICryptoService
from observer.services.storage import IStorage

CHUNK_SIZE = 512


class DownloadHandler:
    def __init__(self, storage: IStorage, crypto: ICryptoService):
        self.storage = storage
        self.crypto = crypto

    async def stream(self, document: Document) -> AsyncGenerator:
        fd = await self.storage.open(document.path)
        secrets = await self.get_encryption_secrets(document)
        encrypted_contents = await fd.read()
        decrypted_contents = await self.crypto.aes_decrypt(secrets, encrypted_contents)
        data = BytesIO(decrypted_contents)
        while True:
            if chunk := data.read(CHUNK_SIZE):
                yield chunk
            else:
                break

    async def get_encryption_secrets(self, document: Document) -> AESCipherOptions:
        key_hash, secrets = (document.encryption_key or "").split(":", maxsplit=1)
        decrypted_secrets = await self.crypto.decrypt(key_hash, secrets.encode())
        encoded_secrets = decrypted_secrets.decode()
        aes_cipher_options = await self.crypto.parse_aes_secrets(encoded_secrets)
        return aes_cipher_options
