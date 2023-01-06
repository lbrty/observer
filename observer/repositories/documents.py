from typing import List, Optional, Protocol

from sqlalchemy import delete, insert, select

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.documents import documents
from observer.entities.documents import Document, NewDocument


class IDocumentsRepository(Protocol):
    async def create_document(self, new_document: NewDocument) -> Document:
        raise NotImplementedError

    async def get_document(self, doc_id: Identifier) -> Optional[Document]:
        raise NotImplementedError

    async def get_by_owner_id(self, owner_id: Identifier) -> List[Document]:
        raise NotImplementedError

    async def delete_document(self, doc_id: Identifier) -> Document:
        raise NotImplementedError


class DocumentsRepository(IDocumentsRepository):
    def __init__(self, db: Database):
        self.db = db

    async def create_document(self, new_document: NewDocument) -> Document:
        query = insert(documents).values(**new_document.dict()).returning("*")
        row = await self.db.fetchone(query)
        return Document(**row)

    async def get_document(self, doc_id: Identifier) -> Optional[Document]:
        query = select(documents).where(documents.c.id == doc_id)
        if row := await self.db.fetchone(query):
            return Document(**row)

        return None

    async def get_by_owner_id(self, owner_id: Identifier) -> List[Document]:
        query = select(documents).where(documents.c.owner_id == owner_id)
        rows = await self.db.fetchall(query)
        return [Document(**row) for row in rows]

    async def delete_document(self, doc_id: Identifier) -> Document:
        query = delete(documents).where(documents.c.id == doc_id).returning("*")
        row = await self.db.fetchone(query)
        return Document(**row)
