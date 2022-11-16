import asyncio
import hashlib
import os

import httpx
import pytest
from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key
from cryptography.hazmat.primitives.serialization import (
    Encoding,
    NoEncryption,
    PrivateFormat,
    load_pem_private_key,
)
from dotenv import load_dotenv
from fastapi import FastAPI
from pydantic.networks import PostgresDsn

from observer.app import create_app
from observer.context import ctx
from observer.db import PoolOptions, connect, disconnect, metadata
from observer.schemas.crypto import PrivateKey
from observer.services.crypto import FSLoader
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
    settings.key_loader = FSLoader()
    priv_key = generate_private_key(
        public_exponent=settings.public_exponent,
        key_size=2048,
    )

    priv_key_bytes = priv_key.private_bytes(
        encoding=Encoding.PEM,
        format=PrivateFormat.PKCS8,
        encryption_algorithm=NoEncryption(),
    )

    h = hashlib.new("sha256", priv_key_bytes)
    settings.key_loader.keys = [
        PrivateKey(
            hash=str(h.hexdigest())[:16].upper(),
            key=load_pem_private_key(
                priv_key_bytes,
                password=None,
            ),
        )
    ]

    return settings


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
def test_app(env_settings, db) -> FastAPI:
    return create_app(env_settings)


@pytest.fixture(scope="function")
async def client(test_app):
    app_client = httpx.AsyncClient(app=test_app, base_url="http://test")
    yield app_client
