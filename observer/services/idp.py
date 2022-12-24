from typing import Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.idp import IDP, NewIDP
from observer.repositories.idp import IDPRepositoryInterface
from observer.schemas.idp import NewIDPRequest
from observer.services.categories import CategoryServiceInterface
from observer.services.crypto import CryptoServiceInterface
from observer.services.projects import ProjectsServiceInterface
from observer.services.world import WorldServiceInterface


class IDPServiceInterface(Protocol):
    tag: str
    repo: IDPRepositoryInterface
    crypto_service: CryptoServiceInterface

    async def create_idp(self, new_idp: NewIDPRequest) -> IDP:
        raise NotImplementedError

    async def get_idp(self, idp_id: Identifier) -> IDP:
        raise NotImplementedError


class IDPService(IDPServiceInterface):
    tag: str = "source=service:idp"

    def __init__(
        self,
        idp_repository: IDPRepositoryInterface,
        crypto_service: CryptoServiceInterface,
        categories: CategoryServiceInterface,
        projects: ProjectsServiceInterface,
        world: WorldServiceInterface,
    ):
        self.repo = idp_repository
        self.crypto_service = crypto_service
        self.categories_service = categories
        self.projects_service = projects
        self.world_service = world

    async def create_idp(self, new_idp: NewIDPRequest) -> IDP:
        new_idp = NewIDP(**new_idp.dict())
        key_hash = self.crypto_service.keychain.keys[0].hash
        if new_idp.project_id:
            await self.projects_service.get_by_id(new_idp.project_id)

        if new_idp.current_place_id:
            await self.world_service.get_place(new_idp.current_place_id)

        if new_idp.from_place_id:
            await self.world_service.get_place(new_idp.from_place_id)

        if new_idp.category_id:
            await self.categories_service.get_category(new_idp.category_id)

        if new_idp.email:
            encrypted_email = await self.crypto_service.encrypt(
                key_hash,
                new_idp.email.encode(),
            )
            new_idp.email = f"{key_hash}:{encrypted_email.decode()}"

        if new_idp.phone_number and new_idp.phone_number.strip():
            encrypted_phone_number = await self.crypto_service.encrypt(
                key_hash,
                new_idp.phone_number.encode(),
            )
            new_idp.phone_number = f"{key_hash}:{encrypted_phone_number.decode()}"

        if new_idp.phone_number_additional and new_idp.phone_number_additional.strip():
            encrypted_phone_number = await self.crypto_service.encrypt(
                key_hash,
                new_idp.phone_number_additional.encode(),
            )
            new_idp.phone_number_additional = f"{key_hash}:{encrypted_phone_number.decode()}"

        return await self.repo.create_idp(new_idp)

    async def get_idp(self, idp_id: Identifier) -> IDP:
        if idp := await self.repo.get_idp(idp_id):
            return idp

        raise NotFoundError(message="IDP record not found")
