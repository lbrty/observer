from typing import Optional, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.support_records import support_records
from observer.entities.support_records import (
    NewSupportRecord,
    SupportRecord,
    UpdateSupportRecord,
)


class ISupportRecordsRepository(Protocol):
    async def create_record(self, new_record: NewSupportRecord) -> SupportRecord:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> Optional[SupportRecord]:
        raise NotImplementedError

    async def update_record(self, record_id: Identifier, updates: UpdateSupportRecord) -> Optional[SupportRecord]:
        raise NotImplementedError

    async def delete_record(self, record_id: Identifier) -> Optional[SupportRecord]:
        raise NotImplementedError


class SupportRecordsRepository(ISupportRecordsRepository):
    def __init__(self, db: Database):
        self.db = db

    async def create_record(self, new_record: NewSupportRecord) -> SupportRecord:
        query = insert(support_records).values(**new_record.dict()).returning("*")
        row = await self.db.fetchone(query)
        return SupportRecord(**row)

    async def get_record(self, record_id: Identifier) -> Optional[SupportRecord]:
        query = select(support_records).where(support_records.c.id == record_id)
        if row := await self.db.fetchone(query):
            return SupportRecord(**row)

        return None

    async def update_record(self, record_id: Identifier, updates: UpdateSupportRecord) -> Optional[SupportRecord]:
        query = update(support_records).values(**updates.dict()).returning("*")
        if row := await self.db.fetchone(query):
            return SupportRecord(**row)

        return None

    async def delete_record(self, record_id: Identifier) -> Optional[SupportRecord]:
        query = delete(support_records).where(support_records.c.id == record_id).returning("*")
        if row := await self.db.fetchone(query):
            return SupportRecord(**row)

        return None
