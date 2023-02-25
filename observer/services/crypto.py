import base64
import os
from dataclasses import dataclass
from typing import Optional, Protocol, Tuple

from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives.hashes import SHA256
from structlog import get_logger

from observer.services.keychain import IKeychain

logger = get_logger(service="crypto")


@dataclass
class AESCipherOptions:
    secret: bytes
    iv: bytes
    tag: bytes


class ICryptoService(Protocol):
    keychain: IKeychain | None
    padding: padding.OAEP | None

    async def encrypt(self, key_hash: Optional[str], data: bytes) -> bytes:
        """Encrypt data and return in Base64 representation"""
        raise NotImplementedError

    async def decrypt(self, key_hash: Optional[str], data: bytes) -> bytes:
        """Decrypt data, encrypted data should have Base64 representation"""
        raise NotImplementedError

    async def aes_cipher_options(self, key_bits: int) -> AESCipherOptions:
        raise NotImplementedError

    async def gen_key(self, key_bits: int) -> bytes:
        raise NotImplementedError

    async def aes_encrypt(self, secrets: AESCipherOptions, data: bytes) -> Tuple[bytes, bytes]:
        raise NotImplementedError

    async def aes_decrypt(self, secrets: AESCipherOptions, ciphertext: bytes) -> bytes:
        raise NotImplementedError

    async def parse_aes_secrets(self, encoded_secrets: str) -> AESCipherOptions:
        raise NotImplementedError

    async def format_aes_secrets(self, secrets: AESCipherOptions) -> str:
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

    async def encrypt(self, key_hash: Optional[str], data: bytes) -> bytes:
        key = await self.keychain.find(key_hash)
        return base64.b64encode(key.private_key.public_key().encrypt(data, self.padding))

    async def decrypt(self, key_hash: Optional[str], data: bytes) -> bytes:
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
        tag = await self.gen_key(16)
        return AESCipherOptions(
            secret=secret,
            iv=iv,
            tag=tag,
        )

    async def gen_key(self, key_bits: int) -> bytes:
        return os.urandom(key_bits)

    async def aes_encrypt(self, secrets: AESCipherOptions, data: bytes) -> Tuple[bytes, bytes]:
        cipher = Cipher(
            algorithms.AES(key=secrets.secret),
            modes.GCM(secrets.iv),
        )
        encryptor = cipher.encryptor()
        ciphertext = encryptor.update(data) + encryptor.finalize()
        return encryptor.tag, ciphertext

    async def aes_decrypt(self, secrets: AESCipherOptions, ciphertext: bytes) -> bytes:
        cipher = Cipher(
            algorithms.AES(key=secrets.secret),
            modes.GCM(secrets.iv, secrets.tag),
        )
        decryptor = cipher.decryptor()
        return decryptor.update(base64.b64decode(ciphertext)) + decryptor.finalize()

    async def parse_aes_secrets(self, encoded_secrets: str) -> AESCipherOptions:
        secret, iv, tag = encoded_secrets.split(":", maxsplit=2)
        return AESCipherOptions(
            secret=base64.b64decode(secret),
            iv=base64.b64decode(iv),
            tag=base64.b64decode(tag),
        )

    async def format_aes_secrets(self, secrets: AESCipherOptions) -> str:
        """Returns AES secrets in the following format

            `<secret>:<iv>:<tag>`
        """
        secret = base64.b64encode(secrets.secret).decode()
        iv = base64.b64encode(secrets.iv).decode()
        tag = base64.b64encode(secrets.tag).decode()
        return f"{secret}:{iv}:{tag}"
