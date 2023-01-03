from typing import Optional, Protocol

from observer.api.exceptions import NotFoundError
from observer.common.types import Identifier
from observer.entities.documents import Document, NewDocument
from observer.repositories.documents import IDocumentsRepository
from observer.schemas.documents import NewDocumentRequest
from observer.services.crypto import ICryptoService


class IDocumentsService(Protocol):
    repo: IDocumentsRepository
    crypto: ICryptoService

    async def create_document(self, new_document: NewDocumentRequest) -> Document:
        raise NotImplementedError

    async def get_document(self, doc_id: Identifier) -> Optional[Document]:
        raise NotImplementedError

    async def delete_document(self, doc_id: Identifier) -> Document:
        raise NotImplementedError


class DocumentsService(IDocumentsService):
    def __init__(self, repo: IDocumentsRepository, crypto: ICryptoService):
        self.repo = repo
        self.crypto = crypto

    async def create_document(self, new_document: NewDocumentRequest) -> Document:
        return await self.repo.create_document(NewDocument(**new_document.dict()))

    async def get_document(self, doc_id: Identifier) -> Optional[Document]:
        if document := await self.repo.get_document(doc_id):
            return document

        raise NotFoundError(message="Document not found")

    async def delete_document(self, doc_id: Identifier) -> Document:
        return await self.repo.delete_document(doc_id)
