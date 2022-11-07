from contextlib import asynccontextmanager
from dataclasses import dataclass
from typing import Protocol

from sqlalchemy.ext.asyncio import AsyncEngine
from sqlalchemy.orm import sessionmaker


class DbProxyInterface(Protocol):
    async def fetchone(self, statement):
        raise NotImplementedError

    async def fetchall(self, statement):
        raise NotImplementedError

    async def execute(self, statement):
        raise NotImplementedError

    async def transaction(self):
        raise NotImplementedError


@dataclass
class Database(DbProxyInterface):
    engine: AsyncEngine
    session: sessionmaker

    async def fetchone(self, statement):
        result = await self.execute(statement)
        return result.fetchone()

    async def fetchall(self, statement):
        result = await self.execute(statement)
        return result.fetchall()

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
            # making sure we call __aexit__ in case an exception is raised
            #  by the code running under this contextmanager
            pass
