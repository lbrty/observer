from datetime import date
from typing import List, Optional, Protocol

from pydantic import BaseModel
from sqlalchemy import and_, delete, desc, insert, select

from observer.common.sqlalchemy import parse_order_by
from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.migration_history import migration_history
from observer.entities.migration_history import MigrationHistory, NewMigrationHistory


class MigrationFilters(BaseModel):
    idp_id: Identifier
    migration_date_from: Optional[date]
    migration_date_to: Optional[date]
    project_id: Identifier
    from_place_id: Optional[Identifier]
    current_place_id: Optional[Identifier]
    order_by: Optional[str]


class IMigrationRepository(Protocol):
    async def add_record(self, new_record: NewMigrationHistory) -> MigrationHistory:
        raise NotImplementedError

    async def get_record(self, record_id: Identifier) -> Optional[MigrationHistory]:
        raise NotImplementedError

    async def get_records(self, filters: MigrationFilters) -> List[MigrationHistory]:
        raise NotImplementedError

    async def get_records_by_person_id(self, idp_id: Identifier) -> List[MigrationHistory]:
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

    async def get_record(self, record_id: Identifier) -> Optional[MigrationHistory]:
        query = select(migration_history).where(migration_history.c.id == record_id)
        if result := await self.db.fetchone(query):
            return MigrationHistory(**result)

        return None

    async def get_records_by_person_id(self, idp_id: Identifier) -> List[MigrationHistory]:
        query = (
            select(migration_history)
            .where(migration_history.c.idp_id == idp_id)
            .order_by(desc(migration_history.c.created_at))
        )
        rows = await self.db.fetchall(query)
        return [MigrationHistory(**row) for row in rows]

    async def get_records(self, filters: MigrationFilters) -> List[MigrationHistory]:
        wheres = []
        if filters.idp_id:
            wheres.append(migration_history.c.idp_id == filters.idp_id)

        if filters.project_id:
            wheres.append(migration_history.c.project_id == filters.project_id)

        if filters.current_place_id:
            wheres.append(migration_history.c.current_place_id == filters.current_place_id)

        if filters.from_place_id:
            wheres.append(migration_history.c.from_place_id == filters.from_place_id)

        if filters.migration_date_from and filters.migration_date_to:
            if filters.migration_date_from >= filters.migration_date_to:
                cond = migration_history.c.migration_date <= filters.migration_date_to
            else:
                cond = and_(
                    migration_history.c.migration_date >= filters.migration_date_from,
                    migration_history.c.migration_date <= filters.migration_date_to,
                )

            if cond:
                wheres.append(cond)

        query = select(migration_history).where(*wheres)
        if order_by := parse_order_by(filters.order_by, migration_history):
            query = query.order_by(order_by)

        rows = await self.db.fetchall(query)
        return [MigrationHistory(**row) for row in rows]

    async def delete_record(self, record_id: Identifier) -> MigrationHistory:
        query = delete(migration_history).where(migration_history.c.id == record_id).returning("*")
        result = await self.db.fetchone(query)
        return MigrationHistory(**result)
