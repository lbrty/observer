import base64
import os
from dataclasses import dataclass
from typing import Protocol

import shortuuid
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives.hashes import SHA256
from structlog import get_logger

from observer.common.types import SomeStr
from observer.services.keys import IKeychain

logger = get_logger(service="crypto")


@dataclass
class AESCipherOptions:
    secret: bytes
    iv: bytes


class ICryptoService(Protocol):
    keychain: IKeychain | None
    padding: padding.OAEP | None

    async def encrypt(self, key_hash: SomeStr, data: bytes) -> bytes:
        raise NotImplementedError

    async def decrypt(self, key_hash: SomeStr, data: bytes) -> bytes:
        raise NotImplementedError

    async def aes_cipher_options(self, key_bits: int) -> AESCipherOptions:
        raise NotImplementedError

    async def gen_key(self, key_bits: int) -> bytes:
        raise NotImplementedError

    async def aes_encrypt(self, secret: bytes, iv: bytes, data: bytes) -> bytes:
        raise NotImplementedError

    async def aes_decrypt(self, secret: bytes, iv: bytes, ciphertext: bytes) -> bytes:
        raise NotImplementedError


class CryptoService(ICryptoService):
    def __init__(self, keychain: IKeychain):
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
        key = await self.keychain.find(key_hash)
        return base64.b64encode(key.private_key.public_key().encrypt(data, self.padding))

    async def decrypt(self, key_hash: SomeStr, data: bytes) -> bytes:
        key = await self.keychain.find(key_hash)
        return key.private_key.decrypt(base64.b64decode(data), self.padding)

    # AES encryption and decryption is based on official docs and recommendations
    # https://cryptography.io/en/latest/hazmat/primitives/symmetric-encryption/
    async def aes_cipher_options(self, key_bits: int) -> AESCipherOptions:
        secret = await self.gen_key(key_bits)
        # 96-bit IV values can be processed more efficiently, so is recommended
        # for situations in which efficiency is more desired.
        # See for more https://crypto.stackexchange.com/questions/41601/aes-gcm-recommended-iv-size-why-12-bytes
        iv = await self.gen_key(12)
        return AESCipherOptions(
            secret=secret,
            iv=iv,
        )

    async def gen_key(self, key_bits: int) -> bytes:
        random_bytes = base64.b64encode(os.urandom(key_bits))
        return bytes(shortuuid.uuid(name=random_bytes.decode()))

    async def aes_encrypt(self, secret: bytes, iv: bytes, data: bytes) -> bytes:
        cipher = Cipher(
            algorithms.AES(key=secret),
            modes.GCM(iv),
        )
        encryptor = cipher.encryptor()
        return encryptor.update(data) + encryptor.finalize()

    async def aes_decrypt(self, secret: bytes, iv: bytes, ciphertext: bytes) -> bytes:
        cipher = Cipher(
            algorithms.AES(key=secret),
            modes.GCM(iv),
        )
        decryptor = cipher.decryptor()
        return decryptor.update(ciphertext) + decryptor.finalize()
