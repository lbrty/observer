import asyncio
import os

import pytest
import pytest_asyncio
from dotenv import load_dotenv
from fastapi import FastAPI
from pydantic.networks import PostgresDsn

from observer.app import create_app
from observer.context import ctx
from observer.db import connect, PoolOptions
from observer.settings import db_settings, settings


@pytest.fixture(scope="session")
def event_loop():
    policy = asyncio.get_event_loop_policy()
    loop = policy.new_event_loop()
    yield loop
    loop.close()


@pytest.fixture(scope="session")
def env_settings():
    load_dotenv(".env.test")
    settings.debug = True
    db_settings.db_uri = PostgresDsn(os.getenv("DB_URI"), scheme="postgresql+asyncpg")


@pytest_asyncio.fixture(scope="session")
async def app(env_settings) -> FastAPI:
    return create_app(env_settings)


@pytest_asyncio.fixture(scope="session")
async def db(env_settings):
    ctx.db = await connect(
        db_settings.db_uri,
        PoolOptions()
    )


@pytest_asyncio.fixture(scope="session")
async def client(event_loop, env_settings, db, app):
    pass
