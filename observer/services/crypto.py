from typing import Protocol

from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.hashes import SHA256
from structlog import get_logger

from observer.api.exceptions import InternalError
from observer.common.types import SomeStr
from observer.schemas.crypto import PrivateKey
from observer.services.keys import Keychain

logger = get_logger(service="crypto")


class CryptoServiceInterface(Protocol):
    keychain: Keychain = None

    async def encrypt(self, key_hash: SomeStr, data: bytes) -> bytes:
        raise NotImplementedError

    async def decrypt(self, key_hash: SomeStr, data: bytes) -> bytes:
        raise NotImplementedError

    async def aes_encrypt(self, secret: str, data: bytes) -> bytes:
        raise NotImplementedError

    async def aes_decrypt(self, secret: str, data: bytes) -> bytes:
        raise NotImplementedError


class CryptoService(CryptoServiceInterface):
    def __init__(self, keychain: Keychain):
        self.keychain = keychain

        # More about padding https://en.wikipedia.org/wiki/Optimal_asymmetric_encryption_padding
        # it is sane default to use OAEP padding which brings some randomness and has proven
        # hardening against "chose ciphertext attacks".
        self.padding = padding.OAEP(
            mgf=padding.MGF1(algorithm=SHA256()),
            algorithm=SHA256(),
            label=None,
        )

    async def encrypt(self, key_hash: SomeStr, data: bytes) -> bytes:
        key = await self.find_key(key_hash)
        return key.private_key.public_key().encrypt(data, self.padding)

    async def decrypt(self, key_hash: SomeStr, data: bytes) -> bytes:
        key = await self.find_key(key_hash)
        return key.private_key.decrypt(data, self.padding)

    async def find_key(self, key_hash: SomeStr) -> PrivateKey:
        for key in self.keychain.keys:
            if key.hash == key_hash:
                return key

        raise InternalError(message="private keys not found")

    async def aes_encrypt(self, secret: str, data: bytes) -> bytes:
        ...

    async def aes_decrypt(self, secret: str, data: bytes) -> bytes:
        ...
