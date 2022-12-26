from typing import Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.idp import IDP, NewIDP, PersonalInfo
from observer.repositories.idp import IDPRepositoryInterface
from observer.schemas.idp import NewIDPRequest, UpdateIDPRequest
from observer.services.categories import CategoryServiceInterface
from observer.services.crypto import CryptoServiceInterface
from observer.services.projects import ProjectsServiceInterface
from observer.services.secrets import SecretsServiceInterface
from observer.services.world import WorldServiceInterface


class IDPServiceInterface(Protocol):
    tag: str
    repo: IDPRepositoryInterface
    crypto_service: CryptoServiceInterface
    categories_service: CategoryServiceInterface
    projects_service: ProjectsServiceInterface
    world_service: WorldServiceInterface
    secrets_service: SecretsServiceInterface

    async def create_idp(self, new_idp: NewIDPRequest) -> IDP:
        raise NotImplementedError

    async def get_idp(self, idp_id: Identifier) -> IDP:
        raise NotImplementedError

    async def get_personal_info(self, idp_id: Identifier) -> PersonalInfo:
        raise NotImplementedError

    async def update_idp(self, idp_id: Identifier, updates: UpdateIDPRequest) -> IDP:
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
        secrets: SecretsServiceInterface
    ):
        self.repo = idp_repository
        self.crypto_service = crypto_service
        self.categories_service = categories
        self.projects_service = projects
        self.world_service = world
        self.secrets_service = secrets

    async def create_idp(self, new_idp: NewIDPRequest) -> IDP:
        new_idp = NewIDP(**new_idp.dict())
        if new_idp.project_id:
            await self.projects_service.get_by_id(new_idp.project_id)

        if new_idp.category_id:
            await self.categories_service.get_category(new_idp.category_id)

        pi = PersonalInfo(
            email=new_idp.email,
            phone_number=new_idp.phone_number,
            phone_number_additional=new_idp.phone_number_additional,
        )

        pi = await self.secrets_service.encrypt_personal_info(pi)
        new_idp.email = pi.email
        new_idp.phone_number = pi.phone_number
        new_idp.phone_number_additional = pi.phone_number_additional
        return await self.repo.create_idp(new_idp)

    async def get_idp(self, idp_id: Identifier) -> IDP:
        if idp := await self.repo.get_idp(idp_id):
            return idp

        raise NotFoundError(message="IDP record not found")

    async def get_personal_info(self, idp_id: Identifier) -> PersonalInfo:
        idp = await self.get_idp(idp_id)
        personal_info = PersonalInfo(
            full_name=idp.full_name,
            email=idp.email,
            phone_number=idp.phone_number,
            phone_number_additional=idp.phone_number_additional,
        )

        return personal_info

    async def update_idp(self, idp_id: Identifier, updates: UpdateIDPRequest) -> IDP:
        idp = await self.get_idp(idp_id)
        return idp
