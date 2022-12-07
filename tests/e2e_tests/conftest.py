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
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.pool import NullPool

from observer.app import create_app
from observer.common.bcrypt import hash_password
from observer.common.types import Role
from observer.context import ctx
from observer.db import Database, metadata
from observer.entities.users import NewUser, User
from observer.repositories.users import UsersRepository
from observer.schemas.crypto import PrivateKey
from observer.services.auth import AuthService
from observer.services.crypto import CryptoService
from observer.services.jwt import JWTService
from observer.services.keys import FSLoader
from observer.services.mfa import MFAService
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
    db_settings.db_uri = PostgresDsn(os.getenv("DB_URI"), scheme="postgresql+asyncpg")
    return settings


@pytest.fixture(scope="session")
async def db_engine(env_settings):
    opts = dict(
        echo=False,
        echo_pool=False,
        isolation_level="AUTOCOMMIT",
        poolclass=NullPool,
    )
    engine = create_async_engine(db_settings.db_uri, **opts)
    yield engine
    await engine.dispose()


@pytest.fixture(scope="session")
async def app_context(db_engine):
    ctx.db = Database(
        engine=db_engine,
        session=sessionmaker(db_engine, class_=AsyncSession),
    )

    ctx.key_loader = FSLoader()
    private_key = generate_private_key(
        public_exponent=settings.public_exponent,
        key_size=settings.key_size,
    )

    private_key_bytes = private_key.private_bytes(
        encoding=Encoding.PEM,
        format=PrivateFormat.PKCS8,
        encryption_algorithm=NoEncryption(),
    )

    h = hashlib.new("sha256", private_key_bytes)
    ctx.key_loader.keys = [
        PrivateKey(
            hash=str(h.hexdigest())[:16].upper(),
            private_key=load_pem_private_key(
                private_key_bytes,
                password=None,
            ),
        )
    ]
    ctx.jwt_service = JWTService(ctx.key_loader.keys[0])
    ctx.crypto_service = CryptoService(ctx.key_loader)
    ctx.mfa_service = MFAService(settings.totp_leeway, ctx.crypto_service)
    ctx.users_repo = UsersRepository(ctx.db)
    ctx.users_service = UsersService(ctx.users_repo)
    ctx.auth_service = AuthService(ctx.jwt_service, ctx.users_service)

    yield ctx


@pytest.fixture(scope="function")
async def ensure_db(env_settings, db_engine):
    async with db_engine.begin() as conn:
        await conn.run_sync(metadata.drop_all)
        await conn.run_sync(metadata.create_all)


@pytest.fixture(scope="session")
def test_app(env_settings) -> FastAPI:
    return create_app(env_settings)


@pytest.fixture(scope="session")
async def client(test_app):
    app_client = httpx.AsyncClient(app=test_app, base_url="http://test")
    yield app_client


@pytest.fixture(scope="function")
async def consultant_user(ensure_db, app_context) -> User:
    user = await app_context.users_service.repo.create_user(
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

    yield user
