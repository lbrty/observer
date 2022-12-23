from typing import Protocol

from observer.entities.idp import IDP, NewIDP
from observer.repositories.idp import IDPRepositoryInterface
from observer.schemas.idp import NewIDPRequest
from observer.services.crypto import CryptoServiceInterface


class IDPServiceInterface(Protocol):
    tag: str
    repo: IDPRepositoryInterface
    crypto_service: CryptoServiceInterface

    async def create_idp(self, new_idp: NewIDPRequest) -> IDP:
        raise NotImplementedError


class IDPService(IDPServiceInterface):
    tag: str = "source=service:idp"

    def __init__(self, idp_repository: IDPRepositoryInterface, crypto_service: CryptoServiceInterface):
        self.repo = idp_repository
        self.crypto_service = crypto_service

    async def create_idp(self, new_idp: NewIDPRequest) -> IDP:
        # TODO: Generate RSA key and assign as encryption key
        #       Then encrypt fields we need to encrypt
        return await self.repo.create_idp(NewIDP(**new_idp.dict()))
