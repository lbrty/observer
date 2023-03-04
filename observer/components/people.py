from fastapi import Depends

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.components.auth import authenticated_user
from observer.components.services import (
    documents_service,
    people_service,
    permissions_service,
)
from observer.entities.people import Person
from observer.entities.users import User
from observer.services.documents import IDocumentsService
from observer.services.people import IPeopleService
from observer.services.permissions import IPermissionsService


class PersonWithTests:
    def __init__(self, *tests):
        self.tests = tests

    async def __call__(
        self,
        person_id: Identifier,
        user: User = Depends(authenticated_user),
        documents: IDocumentsService = Depends(documents_service),
        permissions: IPermissionsService = Depends(permissions_service),
        people: IPeopleService = Depends(people_service),
    ) -> Person:
        person = await people.get_person(person_id)
        permission = None
        if person.project_id:
            try:
                permission = await permissions.find(person.project_id, user.id)
            except NotFoundError:
                pass

        for test in self.tests:
            test(user, permission)

        return person
