from dataclasses import dataclass
from typing import Any

from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import AsyncAdaptedQueuePool
from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, create_async_engine
from sqlalchemy import MetaData

metadata = MetaData()


@dataclass
class Database:
    engine: AsyncEngine
    session: sessionmaker


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
        session=sessionmaker(engine, class_=AsyncEngine),
    )


async def disconnect(engine: AsyncEngine):
    await engine.dispose()
