from typing import Optional, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import EncryptedFieldValue, Identifier
from observer.entities.people import IDP, NewIDP, PersonalInfo, UpdateIDP
from observer.repositories.people import IPeopleRepository
from observer.schemas.people import NewIDPRequest, UpdateIDPRequest
from observer.services.categories import ICategoryService
from observer.services.crypto import ICryptoService
from observer.services.projects import IProjectsService
from observer.services.secrets import ISecretsService
from observer.services.world import IWorldService


class IPeopleService(Protocol):
    repo: IPeopleRepository
    crypto_service: ICryptoService
    categories_service: ICategoryService
    projects_service: IProjectsService
    world_service: IWorldService
    secrets_service: ISecretsService

    async def create_person(self, new_person: NewIDPRequest) -> IDP:
        raise NotImplementedError

    async def get_person(self, person_id: Identifier) -> IDP:
        raise NotImplementedError

    async def update_person(self, person_id: Identifier, updates: UpdateIDPRequest) -> IDP:
        raise NotImplementedError

    async def delete_person(self, person_id: Identifier) -> Optional[IDP]:
        raise NotImplementedError


class PeopleService(IPeopleService):
    def __init__(
        self,
        people_repository: IPeopleRepository,
        crypto_service: ICryptoService,
        categories: ICategoryService,
        projects: IProjectsService,
        world: IWorldService,
        secrets: ISecretsService,
    ):
        self.repo = people_repository
        self.crypto_service = crypto_service
        self.categories_service = categories
        self.projects_service = projects
        self.world_service = world
        self.secrets_service = secrets

    async def create_person(self, new_person: NewIDPRequest) -> IDP:
        new_person = NewIDP(**new_person.dict())
        if new_person.project_id:
            await self.projects_service.get_by_id(new_person.project_id)

        if new_person.category_id:
            await self.categories_service.get_category(new_person.category_id)

        pi = PersonalInfo(
            email=new_person.email,
            phone_number=new_person.phone_number,
            phone_number_additional=new_person.phone_number_additional,
        )

        pi = await self.secrets_service.encrypt_personal_info(pi)
        new_person.email = pi.email
        new_person.phone_number = pi.phone_number
        new_person.phone_number_additional = pi.phone_number_additional
        return await self.repo.create_person(new_person)

    async def get_person(self, person_id: Identifier) -> IDP:
        if idp := await self.repo.get_person(person_id):
            return idp

        raise NotFoundError(message="Person not found")

    async def update_person(self, person_id: Identifier, updates: UpdateIDPRequest) -> IDP:
        """Update IDP record

        NOTES:
            Since we return IDP records with encrypted fields which contain `********`
            instead of real encrypted value we need to check if field does not have
            the value above we can update these field otherwise we need to skip updating them.
            So for this reason we initialize `PersonalInfo` instance which is then populated
            and encrypted and later assigned to relevant `idp_updates` fields.
        """
        person_updates = UpdateIDP(**updates.dict())
        pi = PersonalInfo()
        if updates.email != EncryptedFieldValue:
            pi.email = updates.email

        if updates.phone_number != EncryptedFieldValue:
            pi.phone_number = updates.phone_number

        if updates.phone_number_additional != EncryptedFieldValue:
            pi.phone_number_additional = updates.phone_number_additional

        pi = await self.secrets_service.encrypt_personal_info(pi)
        if pi.email:
            person_updates.email = pi.email

        if pi.phone_number:
            person_updates.phone_number = pi.phone_number

        if pi.phone_number_additional:
            person_updates.phone_number_additional = pi.phone_number_additional

        updated = await self.repo.update_person(person_id, person_updates)
        return updated

    async def delete_person(self, person_id: Identifier) -> Optional[IDP]:
        return await self.repo.delete_person(person_id)
