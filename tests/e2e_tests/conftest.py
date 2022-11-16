import asyncio
import os

import httpx
import pytest
from dotenv import load_dotenv
from fastapi import FastAPI
from pydantic.networks import PostgresDsn

from observer.app import create_app
from observer.context import ctx
from observer.db import PoolOptions, connect, disconnect, metadata
from observer.settings import db_settings, settings


@pytest.fixture(scope="session")
def event_loop():
    policy = asyncio.get_event_loop_policy()
    loop = policy.new_event_loop()
    yield loop
    loop.close()


@pytest.fixture(scope="session")
async def env_settings():
    load_dotenv(".env.test")
    settings.debug = True
    db_settings.db_uri = PostgresDsn(os.getenv("DB_URI"), scheme="postgresql+asyncpg")
    yield settings


@pytest.fixture(scope="session")
async def db(env_settings):
    ctx.db = await connect(db_settings.db_uri, PoolOptions())
    async with ctx.db.engine.begin() as conn:
        try:
            await conn.run_sync(metadata.create_all)
            yield ctx
        finally:
            await conn.run_sync(metadata.drop_all)
            await disconnect(ctx.db.engine)


@pytest.fixture(scope="session")
async def test_app(env_settings, db) -> FastAPI:
    yield create_app(env_settings)


@pytest.fixture(scope="function")
async def client(test_app):
    app_client = httpx.AsyncClient(app=test_app)
    yield app_client
