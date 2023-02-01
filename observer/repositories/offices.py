from typing import List, Optional, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.offices import offices
from observer.entities.offices import Office


class IOfficesRepository(Protocol):
    async def create_office(self, name: str) -> Office:
        raise NotImplementedError

    async def get_office(self, office_id: Identifier) -> Optional[Office]:
        raise NotImplementedError

    async def get_offices(self, name: Optional[str], offset: int, limit: int) -> List[Office]:
        raise NotImplementedError

    async def update_office(self, office_id: Identifier, new_name: str) -> Office:
        raise NotImplementedError

    async def delete_office(self, office_id: Identifier) -> Office:
        raise NotImplementedError


class OfficesRepository(IOfficesRepository):
    def __init__(self, db: Database):
        self.db = db

    async def create_office(self, name: str) -> Office:
        query = insert(offices).values(name=name).returning("*")
        result = await self.db.fetchone(query)
        return Office(**result)

    async def get_office(self, office_id: Identifier) -> Optional[Office]:
        query = select(offices).where(offices.c.id == office_id)
        if result := await self.db.fetchone(query):
            return Office(**result)

        return None

    async def get_offices(self, name: Optional[str], offset: int, limit: int) -> List[Office]:
        query = select(offices)
        if name:
            query = query.where(offices.c.name.ilike(f"%{name}%"))

        query = query.offset(offset).limit(limit)
        result = await self.db.fetchall(query)
        return [Office(**row) for row in result]

    async def update_office(self, office_id: Identifier, new_name: str) -> Office:
        query = (
            update(offices)
            .values(name=new_name)
            .where(
                offices.c.id == office_id,
            )
            .returning("*")
        )
        result = await self.db.fetchone(query)
        return Office(**result)

    async def delete_office(self, office_id: Identifier) -> Office:
        query = (
            delete(offices)
            .where(
                offices.c.id == office_id,
            )
            .returning("*")
        )
        result = await self.db.fetchone(query)
        return Office(**result)
