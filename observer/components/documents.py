from fastapi import Depends

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.components.auth import authenticated_user
from observer.components.services import documents_service, permissions_service
from observer.entities.documents import Document
from observer.entities.users import User
from observer.services.documents import IDocumentsService
from observer.services.permissions import IPermissionsService


class DocumentWithTests:
    def __init__(self, *tests):
        self.tests = tests

    async def __call__(
        self,
        doc_id: Identifier,
        user: User = Depends(authenticated_user),
        documents: IDocumentsService = Depends(documents_service),
        permissions: IPermissionsService = Depends(permissions_service),
    ) -> Document:
        document = await documents.get_document(doc_id)
        permission = None
        if document.project_id:
            try:
                permission = await permissions.find(document.project_id, user.id)
            except NotFoundError:
                pass

        for test in self.tests:
            test(user, permission)

        return document
