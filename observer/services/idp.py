from typing import Optional, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import EncryptedFieldValue, Identifier
from observer.entities.idp import IDP, NewIDP, PersonalInfo, UpdateIDP
from observer.repositories.idp import IIDPRepository
from observer.schemas.idp import NewIDPRequest, UpdateIDPRequest
from observer.services.categories import ICategoryService
from observer.services.crypto import ICryptoService
from observer.services.projects import IProjectsService
from observer.services.secrets import ISecretsService
from observer.services.world import IWorldService


class IIDPService(Protocol):
    tag: str
    repo: IIDPRepository
    crypto_service: ICryptoService
    categories_service: ICategoryService
    projects_service: IProjectsService
    world_service: IWorldService
    secrets_service: ISecretsService

    async def create_idp(self, new_idp: NewIDPRequest) -> IDP:
        raise NotImplementedError

    async def get_idp(self, idp_id: Identifier) -> IDP:
        raise NotImplementedError

    async def update_idp(self, idp_id: Identifier, updates: UpdateIDPRequest) -> IDP:
        raise NotImplementedError

    async def delete_idp(self, idp_id: Identifier) -> Optional[IDP]:
        raise NotImplementedError


class IDPService(IIDPService):
    tag: str = "source=service:idp"

    def __init__(
        self,
        idp_repository: IIDPRepository,
        crypto_service: ICryptoService,
        categories: ICategoryService,
        projects: IProjectsService,
        world: IWorldService,
        secrets: ISecretsService,
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

    async def update_idp(self, idp_id: Identifier, updates: UpdateIDPRequest) -> IDP:
        """Update IDP record

        NOTES:
            Since we return IDP records with encrypted fields which contain `********`
            instead of real encrypted value we need to check if field does not have
            the value above we can update these field otherwise we need to skip updating them.
            So for this reason we initialize `PersonalInfo` instance which is then populated
            and encrypted and later assigned to relevant `idp_updates` fields.
        """
        idp_updates = UpdateIDP(**updates.dict())
        pi = PersonalInfo()
        if updates.email != EncryptedFieldValue:
            pi.email = updates.email

        if updates.phone_number != EncryptedFieldValue:
            pi.phone_number = updates.phone_number

        if updates.phone_number_additional != EncryptedFieldValue:
            pi.phone_number_additional = updates.phone_number_additional

        pi = await self.secrets_service.encrypt_personal_info(pi)
        if pi.email:
            idp_updates.email = pi.email

        if pi.phone_number:
            idp_updates.phone_number = pi.phone_number

        if pi.phone_number_additional:
            idp_updates.phone_number_additional = pi.phone_number_additional

        updated = await self.repo.update_idp(idp_id, idp_updates)
        return updated

    async def delete_idp(self, idp_id: Identifier) -> Optional[IDP]:
        return await self.repo.delete_idp(idp_id)
