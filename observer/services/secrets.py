from typing import IO, Any, Protocol

from observer.entities.idp import PersonalInfo
from observer.services.crypto import CryptoServiceInterface


class SecretsServiceInterface(Protocol):
    tag: str
    crypto_service: CryptoServiceInterface

    async def encrypt_personal_info(self, pi: PersonalInfo) -> PersonalInfo:
        raise NotImplementedError

    async def decrypt_personal_info(self, pi: PersonalInfo) -> PersonalInfo:
        raise NotImplementedError

    async def encrypt_document(self, secret: str, stream: IO[Any]) -> bytes:
        raise NotImplementedError

    async def decrypt_document(self, secret: str, stream: IO[Any]) -> bytes:
        raise NotImplementedError


class SecretsService(SecretsServiceInterface):
    tag: str = "service=secrets"
    crypto_service: CryptoServiceInterface

    def __init__(self, crypto_service: CryptoServiceInterface):
        self.crypto_service = crypto_service

    async def encrypt_personal_info(self, pi: PersonalInfo) -> PersonalInfo:
        key_hash = self.crypto_service.keychain.keys[0].hash
        if pi.email:
            encrypted_email = await self.crypto_service.encrypt(
                key_hash,
                pi.email.encode(),
            )
            pi.email = f"{key_hash}:{encrypted_email.decode()}"

        if pi.phone_number and pi.phone_number.strip():
            encrypted_phone_number = await self.crypto_service.encrypt(
                key_hash,
                pi.phone_number.encode(),
            )
            pi.phone_number = f"{key_hash}:{encrypted_phone_number.decode()}"

        if pi.phone_number_additional and pi.phone_number_additional.strip():
            encrypted_phone_number = await self.crypto_service.encrypt(
                key_hash,
                pi.phone_number_additional.encode(),
            )
            pi.phone_number_additional = f"{key_hash}:{encrypted_phone_number.decode()}"

        return pi

    async def decrypt_personal_info(self, pi: PersonalInfo) -> PersonalInfo:
        if ":" in str(pi.email):
            key_hash, data = pi.email.split(":", maxsplit=1)
            decrypted_email = await self.crypto_service.decrypt(
                key_hash,
                data.encode(),
            )
            pi.email = decrypted_email.decode()

        if ":" in str(pi.phone_number):
            key_hash, data = pi.phone_number.split(":", maxsplit=1)
            decrypted_phone_number = await self.crypto_service.decrypt(
                key_hash,
                data.encode(),
            )
            pi.phone_number = decrypted_phone_number.decode()

        if ":" in str(pi.phone_number_additional):
            key_hash, data = pi.phone_number_additional.split(":", maxsplit=1)
            decrypted_phone_number_additional = await self.crypto_service.decrypt(
                key_hash,
                data.encode(),
            )
            pi.phone_number_additional = decrypted_phone_number_additional.decode()

        return pi

    async def encrypt_document(self, secret: str, stream: IO[Any]) -> bytes:
        raise NotImplementedError

    async def decrypt_document(self, secret: str, stream: IO[Any]) -> bytes:
        raise NotImplementedError
