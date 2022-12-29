from typing import IO, Any, Protocol

from observer.entities.idp import IDP, PersonalInfo
from observer.services.crypto import ICryptoService


class ISecretsService(Protocol):
    tag: str
    crypto_service: ICryptoService

    async def encrypt_personal_info(self, pi: PersonalInfo) -> PersonalInfo:
        raise NotImplementedError

    async def decrypt_personal_info(self, pi: PersonalInfo) -> PersonalInfo:
        raise NotImplementedError

    async def encrypt_document(self, secret: str, stream: IO[Any]) -> bytes:
        raise NotImplementedError

    async def decrypt_document(self, secret: str, stream: IO[Any]) -> bytes:
        raise NotImplementedError

    async def anonymize_idp(self, idp: IDP) -> IDP:
        if idp.email:
            idp.email = "*" * 8

        if idp.phone_number:
            idp.phone_number = "*" * 8

        if idp.phone_number_additional:
            idp.phone_number_additional = "*" * 8

        return idp


class SecretsService(ISecretsService):
    tag: str = "service=secrets"
    crypto_service: ICryptoService

    def __init__(self, crypto_service: ICryptoService):
        self.crypto_service = crypto_service

    async def encrypt_personal_info(self, pi: PersonalInfo) -> PersonalInfo:
        key_hash = self.crypto_service.keychain.keys[0].hash
        personal_info = PersonalInfo()
        if pi.email:
            encrypted_email = await self.crypto_service.encrypt(
                key_hash,
                pi.email.encode(),
            )
            personal_info.email = f"{key_hash}:{encrypted_email.decode()}"

        if pi.phone_number and pi.phone_number.strip():
            encrypted_phone_number = await self.crypto_service.encrypt(
                key_hash,
                pi.phone_number.encode(),
            )
            personal_info.phone_number = f"{key_hash}:{encrypted_phone_number.decode()}"

        if pi.phone_number_additional and pi.phone_number_additional.strip():
            encrypted_phone_number = await self.crypto_service.encrypt(
                key_hash,
                pi.phone_number_additional.encode(),
            )
            personal_info.phone_number_additional = f"{key_hash}:{encrypted_phone_number.decode()}"

        return personal_info

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
