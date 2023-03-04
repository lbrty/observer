from typing import List, Optional, Protocol, Sequence

from sqlalchemy import delete, insert, select, text

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

    async def get_by_project_id(self, project_id: Identifier) -> List[Document]:
        raise NotImplementedError

    async def delete_document(self, doc_id: Identifier) -> Optional[Document]:
        raise NotImplementedError

    async def bulk_delete(self, doc_ids: Sequence[Identifier]) -> List[Identifier]:
        raise NotImplementedError


class DocumentsRepository(IDocumentsRepository):
    def __init__(self, db: Database):
        self.db = db

    async def create_document(self, new_document: NewDocument) -> Document:
        values = new_document.dict()
        # values["created_at"]
        query = insert(documents).values(**values).returning("*")
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

    async def get_by_project_id(self, project_id: Identifier) -> List[Document]:
        query = select(documents).where(documents.c.project_id == project_id)
        rows = await self.db.fetchall(query)
        return [Document(**row) for row in rows]

    async def delete_document(self, doc_id: Identifier) -> Optional[Document]:
        query = delete(documents).where(documents.c.id == doc_id).returning("*")
        if row := await self.db.fetchone(query):
            return Document(**row)

        return None

    async def bulk_delete(self, doc_ids: Sequence[Identifier]) -> List[Identifier]:
        query = delete(documents).where(documents.c.id.in_(doc_ids)).returning(text("id"))
        rows = await self.db.fetchall(query)
        return [row["id"] for row in rows]
