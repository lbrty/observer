from typing import Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import EncryptedFieldValue, Identifier
from observer.entities.people import NewPerson, Person, PersonalInfo, UpdatePerson
from observer.repositories.people import IPeopleRepository
from observer.schemas.people import NewPersonRequest, UpdatePersonRequest
from observer.services.categories import ICategoryService
from observer.services.crypto import ICryptoService
from observer.services.offices import IOfficesService
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

    async def create_person(self, new_person: NewPersonRequest) -> Person:
        raise NotImplementedError

    async def get_person(self, person_id: Identifier) -> Person:
        raise NotImplementedError

    async def update_person(self, person_id: Identifier, updates: UpdatePersonRequest) -> Person:
        raise NotImplementedError

    async def delete_person(self, person_id: Identifier) -> Person:
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
        offices: IOfficesService,
    ):
        self.repo = people_repository
        self.crypto_service = crypto_service
        self.categories_service = categories
        self.projects_service = projects
        self.world_service = world
        self.secrets_service = secrets
        self.offices_service = offices

    async def create_person(self, new_person: NewPersonRequest) -> Person:
        person = NewPerson(**new_person.model_dump())
        if person.project_id:
            await self.projects_service.get_by_id(person.project_id)

        if person.category_id:
            await self.categories_service.get_category(person.category_id)

        if person.office_id:
            await self.offices_service.get_office(person.office_id)

        pi = PersonalInfo(
            email=person.email,
            phone_number=person.phone_number,
            phone_number_additional=person.phone_number_additional,
        )

        pi = await self.secrets_service.encrypt_personal_info(pi)
        person.email = pi.email  # type: ignore
        person.phone_number = pi.phone_number  # type: ignore
        person.phone_number_additional = pi.phone_number_additional  # type: ignore
        return await self.repo.create_person(person)

    async def get_person(self, person_id: Identifier) -> Person:
        if person := await self.repo.get_person(person_id):
            return person

        raise NotFoundError(message="Person not found")

    async def update_person(self, person_id: Identifier, updates: UpdatePersonRequest) -> Person:
        """Update person

        NOTES:
            Since we return person with encrypted fields which contain `********`
            instead of real encrypted value we need to check if field does not have
            the value above we can update these field otherwise we need to skip updating them.
            So for this reason we initialize `PersonalInfo` instance which is then populated
            and encrypted and later assigned to relevant `updates` fields.
        """
        person_updates = UpdatePerson(**updates.model_dump())
        pi = PersonalInfo()
        if updates.email and updates.email != EncryptedFieldValue:
            pi.email = updates.email

        if updates.phone_number and updates.phone_number != EncryptedFieldValue:
            pi.phone_number = updates.phone_number

        if updates.phone_number_additional and updates.phone_number_additional != EncryptedFieldValue:
            pi.phone_number_additional = updates.phone_number_additional

        pi = await self.secrets_service.encrypt_personal_info(pi)
        if pi.email:
            person_updates.email = pi.email

        if pi.phone_number:
            person_updates.phone_number = pi.phone_number

        if pi.phone_number_additional:
            person_updates.phone_number_additional = pi.phone_number_additional

        if updated := await self.repo.update_person(person_id, person_updates):
            return updated

        raise NotFoundError(message="Person not found")

    async def delete_person(self, person_id: Identifier) -> Person:
        if person := await self.repo.delete_person(person_id):
            return person

        raise NotFoundError(message="Person not found")
