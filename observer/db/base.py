from contextlib import asynccontextmanager
from dataclasses import dataclass
from typing import Protocol

from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, async_sessionmaker


class DbProxy(Protocol):
    async def fetchone(self, statement):
        raise NotImplementedError

    async def fetchall(self, statement):
        raise NotImplementedError

    async def execute(self, statement):
        raise NotImplementedError

    async def transaction(self):
        raise NotImplementedError


@dataclass
class Database(DbProxy):
    engine: AsyncEngine
    session: async_sessionmaker[AsyncSession]

    async def fetchone(self, statement):
        result = await self.execute(statement)
        return result.mappings().fetchone()

    async def fetchall(self, statement):
        result = await self.execute(statement)
        return result.mappings().fetchall()

    async def execute(self, statement):
        async with self.session.begin() as conn:
            try:
                result = await conn.execute(statement)
                await conn.commit()
                return result
            except Exception:
                await conn.rollback()
                raise

    @asynccontextmanager
    async def transaction(self):
        try:
            async with self.session.begin() as conn:
                yield
                await conn.commit()
        finally:
            pass
