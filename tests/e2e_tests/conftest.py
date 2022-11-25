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
from observer.common.bcrypt import hash_password
from observer.common.types import Role
from observer.context import ctx
from observer.db import PoolOptions, connect, disconnect, metadata
from observer.entities.users import NewUser
from observer.repositories.users import UsersRepository
from observer.schemas.crypto import PrivateKey
from observer.services.crypto import FSLoader
from observer.services.jwt import JWTService
from observer.services.users import UsersService
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
    print(os.getenv("DB_URI"))
    db_settings.db_uri = PostgresDsn(os.getenv("DB_URI"), scheme="postgresql+asyncpg")
    return settings


@pytest.fixture(scope="session")
async def encryption_keys(env_settings):
    ctx.key_loader = FSLoader()
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
    ctx.key_loader.keys = [
        PrivateKey(
            hash=str(h.hexdigest())[:16].upper(),
            key=load_pem_private_key(
                priv_key_bytes,
                password=None,
            ),
        )
    ]
    return ctx.key_loader


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
def test_app(env_settings, db, encryption_keys) -> FastAPI:
    ctx.jwt_service = JWTService(ctx.key_loader.keys[0])
    users_repo = UsersRepository(ctx.db)
    ctx.users_service = UsersService(users_repo)
    return create_app(env_settings)


@pytest.fixture(scope="function")
async def client(test_app):
    app_client = httpx.AsyncClient(app=test_app, base_url="http://test")
    yield app_client


@pytest.fixture(scope="function")
async def consultant_user(db):
    user = await ctx.users_service.repo.create_user(
        NewUser(
            ref_id="ref-consultant-1",
            email="consultant-1@example.com",
            full_name="full name",
            password_hash=hash_password("secret"),
            role=Role.consultant,
            is_active=True,
            is_confirmed=True,
        )
    )
    return user
