from typing import List, Protocol

from sqlalchemy import delete, insert, select

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.migration_history import migration_history
from observer.entities.migration_history import (
    MigrationHistory,
    NewMigrationHistory,
)


class IMigrationRepository(Protocol):
    async def add_record(self, new_record: NewMigrationHistory) -> MigrationHistory:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError

    async def get_records(self, record_id: Identifier) -> List[MigrationHistory]:
        raise NotImplementedError

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        raise NotImplementedError


class MigrationRepository(IMigrationRepository):
    def __init__(self, db: Database):
        self.db = db

    async def add_record(self, new_record: NewMigrationHistory) -> MigrationHistory:
        query = insert(migration_history).values(**new_record.dict()).returning("*")
        result = await self.db.fetchone(query)
        return MigrationHistory(**result)

    async def get_record(self, record_id: Identifier) -> MigrationHistory:
        query = select(migration_history).where(migration_history.c.id == record_id)
        result = await self.db.fetchone(query)
        return MigrationHistory(**result)

    async def get_records(self, record_id: Identifier) -> List[MigrationHistory]:
        query = select(migration_history).where(migration_history.c.id == record_id)
        rows = await self.db.fetchall(query)
        return [MigrationHistory(**row) for row in rows]

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        query = delete(migration_history).where(migration_history.c.id == record_id).returning("*")
        result = await self.db.fetchone(query)
        return MigrationHistory(**result)
