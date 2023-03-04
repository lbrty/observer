from typing import List, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.documents import Document, NewDocument
from observer.repositories.documents import IDocumentsRepository
from observer.schemas.documents import NewDocumentRequest
from observer.services.crypto import ICryptoService


class IDocumentsService(Protocol):
    repo: IDocumentsRepository
    crypto: ICryptoService

    async def create_document(self, encryption_key: str, new_document: NewDocumentRequest) -> Document:
        raise NotImplementedError

    async def get_document(self, doc_id: Identifier) -> Document:
        raise NotImplementedError

    async def get_by_owner_id(self, owner_id: Identifier) -> List[Document]:
        raise NotImplementedError

    async def get_by_project_id(self, project_id: Identifier) -> List[Document]:
        raise NotImplementedError

    async def delete_document(self, doc_id: Identifier) -> Document:
        raise NotImplementedError

    async def bulk_delete(self, doc_ids: List[Identifier]) -> List[Identifier]:
        raise NotImplementedError


class DocumentsService(IDocumentsService):
    def __init__(self, repo: IDocumentsRepository, crypto: ICryptoService):
        self.repo = repo
        self.crypto = crypto

    async def create_document(self, encryption_key: str, new_document: NewDocumentRequest) -> Document:
        return await self.repo.create_document(
            NewDocument(
                **dict(**new_document.dict(), encryption_key=encryption_key),
            )
        )

    async def get_by_owner_id(self, owner_id: Identifier) -> List[Document]:
        return await self.repo.get_by_owner_id(owner_id)

    async def get_by_project_id(self, project_id: Identifier) -> List[Document]:
        return await self.repo.get_by_project_id(project_id)

    async def get_document(self, doc_id: Identifier) -> Document:
        if document := await self.repo.get_document(doc_id):
            return document

        raise NotFoundError(message="Document not found")

    async def delete_document(self, doc_id: Identifier) -> Document:
        if document := await self.repo.delete_document(doc_id):
            return document

        raise NotFoundError(message="Document not found")

    async def bulk_delete(self, doc_ids: List[Identifier]) -> List[Identifier]:
        return await self.repo.bulk_delete(doc_ids)
