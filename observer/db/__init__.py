from contextlib import asynccontextmanager
from dataclasses import dataclass
from typing import Any, Protocol

from sqlalchemy import MetaData
from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import AsyncAdaptedQueuePool

metadata = MetaData()


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


@dataclass
class PoolOptions:
    pool_size: int = 5
    max_overflow: int = 10
    pool_timeout: int = 30
    echo: bool = False
    echo_pool: bool = False
    pool_class: Any = AsyncAdaptedQueuePool


async def connect(uri: str, pool_options: PoolOptions) -> Database:
    engine = create_async_engine(
        uri,
        echo=pool_options.echo,
        echo_pool=pool_options.echo_pool,
        pool_size=pool_options.pool_size,
        pool_timeout=pool_options.pool_timeout,
        max_overflow=pool_options.max_overflow,
        poolclass=pool_options.pool_class,
    )

    return Database(
        engine=engine,
        session=sessionmaker(engine, class_=AsyncSession),
    )


async def disconnect(engine: AsyncEngine):
    await engine.dispose()
