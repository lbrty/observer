from typing import Protocol

from structlog import get_logger

from observer.services.keys import KeychainLoader

logger = get_logger(service="crypto")


class CryptoServiceInterface(Protocol):
    keys: KeychainLoader = None

    async def encrypt(self, data: bytes, key_hash: str | None = None) -> bytes:
        raise NotImplementedError

    async def decrypt(self, data: bytes, key_hash: str | None = None) -> bytes:
        raise NotImplementedError


class CryptoService(CryptoServiceInterface):
    def __init__(self, keychain: KeychainLoader):
        self.keys = keychain

    async def encrypt(self, data: bytes, key_hash: str | None = None) -> bytes:
        ...

    async def decrypt(self, data: bytes, key_hash: str | None = None) -> bytes:
        ...
