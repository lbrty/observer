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
from observer.entities.users import NewUser
from observer.entities.world import NewCountry, NewState
from observer.repositories.audit_logs import AuditRepository
from observer.repositories.categories import CategoryRepository
from observer.repositories.idp import IDPRepository
from observer.repositories.migration_history import MigrationRepository
from observer.repositories.permissions import PermissionsRepository
from observer.repositories.projects import ProjectsRepository
from observer.repositories.users import UsersRepository
from observer.repositories.world import WorldRepository
from observer.schemas.crypto import PrivateKey
from observer.services.audit_logs import AuditService
from observer.services.auth import AuthService
from observer.services.categories import CategoryService
from observer.services.crypto import CryptoService
from observer.services.idp import IDPService
from observer.services.jwt import JWTService
from observer.services.keys import FS
from observer.services.mfa import MFAService
from observer.services.migration_history import MigrationService
from observer.services.permissions import PermissionsService
from observer.services.projects import ProjectsService
from observer.services.secrets import SecretsService
from observer.services.users import UsersService
from observer.services.world import WorldService
from observer.settings import db_settings, settings
from tests.mocks.mailer import MockMailer


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

    ctx.keychain = FS()
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
    ctx.keychain.keys = [
        PrivateKey(
            hash=str(h.hexdigest())[:16].upper(),
            private_key=load_pem_private_key(
                private_key_bytes,
                password=None,
            ),
        )
    ]
    ctx.jwt_service = JWTService(ctx.keychain.keys[0])
    ctx.audit_repo = AuditRepository(ctx.db)
    ctx.mailer = MockMailer()
    ctx.audit_service = AuditService(ctx.audit_repo)
    ctx.crypto_service = CryptoService(ctx.keychain)
    ctx.mfa_service = MFAService(settings.totp_leeway, ctx.crypto_service)
    ctx.users_repo = UsersRepository(ctx.db)
    ctx.users_service = UsersService(ctx.users_repo, ctx.crypto_service)
    ctx.auth_service = AuthService(
        ctx.crypto_service,
        ctx.mfa_service,
        ctx.jwt_service,
        ctx.users_service,
    )
    ctx.secrets_service = SecretsService(ctx.crypto_service)
    ctx.projects_repo = ProjectsRepository(ctx.db)
    ctx.projects_service = ProjectsService(ctx.projects_repo)
    ctx.permissions_repo = PermissionsRepository(ctx.db)
    ctx.permissions_service = PermissionsService(ctx.permissions_repo)
    ctx.world_repo = WorldRepository(ctx.db)
    ctx.world_service = WorldService(ctx.world_repo)
    ctx.category_repo = CategoryRepository(ctx.db)
    ctx.category_service = CategoryService(ctx.category_repo)
    ctx.idp_repo = IDPRepository(ctx.db)
    ctx.idp_service = IDPService(
        ctx.idp_repo,
        ctx.crypto_service,
        ctx.category_service,
        ctx.projects_service,
        ctx.world_service,
        ctx.secrets_service,
    )
    ctx.migrations_repo = MigrationRepository(ctx.db)
    ctx.migrations_service = MigrationService(ctx.migrations_repo, ctx.world_service)

    yield ctx


@pytest.fixture(scope="function")
async def ensure_db(env_settings, db_engine):
    async with db_engine.begin() as conn:
        await conn.run_sync(metadata.drop_all)
        await conn.run_sync(metadata.create_all)


@pytest.fixture(scope="session")
def test_app(env_settings) -> FastAPI:
    return create_app(env_settings)


@pytest.fixture(scope="function")
async def consultant_user(ensure_db, app_context):  # type:ignore
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


@pytest.fixture(scope="function")
async def guest_user(ensure_db, app_context):  # type:ignore
    user = await app_context.users_service.repo.create_user(
        NewUser(
            ref_id="ref-guest-1",
            email="guest-1@example.com",
            full_name="full name",
            password_hash=hash_password("secret"),
            role=Role.guest,
            is_active=True,
            is_confirmed=True,
        )
    )

    yield user


@pytest.fixture(scope="function")
async def admin_user(ensure_db, app_context):  # type:ignore
    user = await app_context.users_service.repo.create_user(
        NewUser(
            ref_id="ref-admin-1",
            email="admin-1@example.com",
            full_name="full name",
            password_hash=hash_password("secret"),
            role=Role.admin,
            is_active=True,
            is_confirmed=True,
        )
    )

    yield user


@pytest.fixture(scope="function")
async def staff_user(ensure_db, app_context):  # type:ignore
    user = await app_context.users_service.repo.create_user(
        NewUser(
            ref_id="ref-staff-1",
            email="staff-1@example.com",
            full_name="full name",
            password_hash=hash_password("secret"),
            role=Role.staff,
            is_active=True,
            is_confirmed=True,
        )
    )

    yield user


@pytest.fixture(scope="function")
async def default_country(ensure_db, app_context):  # type:ignore
    country = await app_context.world_repo.create_country(
        NewCountry(
            name="No Stack Country",
            code="nsc",
        )
    )

    yield country


@pytest.fixture(scope="function")
async def default_state(ensure_db, app_context, default_country):  # type:ignore
    state = await app_context.world_repo.create_state(
        NewState(
            name="No Code State",
            code="ncs",
            country_id=default_country.id,
        )
    )

    yield state


@pytest.fixture(scope="session")
async def client(test_app):
    app_client = httpx.AsyncClient(app=test_app, base_url="http://test")
    yield app_client


@pytest.fixture(scope="function")
async def authorized_client(test_app, app_context, consultant_user):
    token = await app_context.auth_service.create_token(consultant_user.ref_id)  # noqa
    app_client = httpx.AsyncClient(
        app=test_app,
        base_url="http://test",
        cookies=token.dict(),
    )
    yield app_client


@pytest.fixture(scope="function")
async def clean_mailbox(app_context):
    app_context.mailer.messages = []
