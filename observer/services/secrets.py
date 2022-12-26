from typing import Protocol, IO, Any

from observer.entities.idp import PersonalInfo
from observer.services.crypto import CryptoServiceInterface


class SecretsServiceInterface(Protocol):
    tag: str
    crypto_service: CryptoServiceInterface

    async def encrypt_personal_info(self, personal_info: PersonalInfo) -> PersonalInfo:
        raise NotImplementedError

    async def decrypt_personal_info(self, personal_info: PersonalInfo) -> PersonalInfo:
        raise NotImplementedError

    async def encrypt_document(self, secret: str, stream: IO[Any]) -> PersonalInfo:
        raise NotImplementedError

    async def decrypt_document(self, secret: str, stream: IO[Any]) -> PersonalInfo:
        raise NotImplementedError


class SecretsService(SecretsServiceInterface):
    tag: str = "service=secrets"
    crypto_service: CryptoServiceInterface

    def __init__(self, crypto_service: CryptoServiceInterface):
        self.crypto_service = crypto_service

    async def encrypt_personal_info(self, personal_info: PersonalInfo) -> PersonalInfo:
        raise NotImplementedError

    async def decrypt_personal_info(self, personal_info: PersonalInfo) -> PersonalInfo:
        raise NotImplementedError

    async def encrypt_document(self, secret: str, stream: IO[Any]) -> PersonalInfo:
        raise NotImplementedError

    async def decrypt_document(self, secret: str, stream: IO[Any]) -> PersonalInfo:
        raise NotImplementedError
