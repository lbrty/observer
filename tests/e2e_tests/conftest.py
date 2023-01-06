import asyncio
import hashlib
import os
from itertools import chain
from pathlib import Path

import aiofiles
import httpx
import pytest
from aiobotocore.config import AioConfig
from aiobotocore.session import AioSession
from aiofiles.tempfile import TemporaryDirectory
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
from observer.context import Repositories, ctx
from observer.db import Database, metadata
from observer.entities.users import NewUser
from observer.entities.world import NewCountry, NewState
from observer.repositories.audit_logs import AuditRepository
from observer.repositories.categories import CategoryRepository
from observer.repositories.documents import DocumentsRepository
from observer.repositories.idp import IDPRepository
from observer.repositories.migration_history import MigrationRepository
from observer.repositories.permissions import PermissionsRepository
from observer.repositories.pets import PetsRepository
from observer.repositories.projects import ProjectsRepository
from observer.repositories.support_records import SupportRecordsRepository
from observer.repositories.users import UsersRepository
from observer.repositories.world import WorldRepository
from observer.schemas.crypto import PrivateKey
from observer.services.audit_logs import AuditService
from observer.services.auth import AuthService
from observer.services.categories import CategoryService
from observer.services.crypto import CryptoService
from observer.services.documents import DocumentsService
from observer.services.idp import IDPService
from observer.services.jwt import JWTService
from observer.services.keychain import Keychain
from observer.services.mfa import MFAService
from observer.services.migration_history import MigrationService
from observer.services.permissions import PermissionsService
from observer.services.pets import PetsService
from observer.services.projects import ProjectsService
from observer.services.secrets import SecretsService
from observer.services.storage import FSStorage
from observer.services.support_records import SupportRecordsService
from observer.services.uploads import UploadHandler
from observer.services.users import UsersService
from observer.services.world import WorldService
from observer.settings import db_settings, settings
from tests.e2e_tests.moto_server import MotoService
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


# TODO: Cleanup moto fixtures
@pytest.fixture(scope="function")
def aws_credentials():
    """Mocked AWS Credentials for moto."""
    os.environ["AWS_ACCESS_KEY_ID"] = "testing"
    os.environ["AWS_SECRET_ACCESS_KEY"] = "testing"
    os.environ["AWS_SECURITY_TOKEN"] = "testing"
    os.environ["AWS_SESSION_TOKEN"] = "testing"
    os.environ["AWS_DEFAULT_REGION"] = "us-central-1"


@pytest.fixture(scope="function")
def aio_session():
    session = AioSession()
    return session


def moto_config():
    return {"aws_secret_access_key": "xxx", "aws_access_key_id": "xxx"}


@pytest.fixture(scope="function")
def region() -> str:
    return "eu-central-1"


@pytest.fixture(scope="function")
def signature_version() -> str:
    return "v4"


@pytest.fixture(scope="function")
def aio_config(signature_version) -> AioConfig:
    return AioConfig(signature_version=signature_version, read_timeout=5, connect_timeout=5)


@pytest.fixture(scope="function")
async def s3_server():
    async with MotoService("s3", ssl=False) as svc:
        yield svc.endpoint_url


@pytest.fixture(scope="function")
async def s3_client(region, aio_session, aio_config, s3_server):
    async with aio_session.create_client("s3", region_name=region, endpoint_url=s3_server, config=aio_config) as client:
        yield client


async def recursive_delete(s3_client, bucket_name):
    # Recursively deletes a bucket and all of its contents.
    paginator = s3_client.get_paginator("list_object_versions")
    async for n in paginator.paginate(Bucket=bucket_name, Prefix=""):
        for obj in chain(
            n.get("Versions", []),
            n.get("DeleteMarkers", []),
            n.get("Contents", []),
            n.get("CommonPrefixes", []),
        ):
            kwargs = dict(Bucket=bucket_name, Key=obj["Key"])
            if "VersionId" in obj:
                kwargs["VersionId"] = obj["VersionId"]
            resp = await s3_client.delete_object(**kwargs)
            assert resp["ResponseMetadata"]["HTTPStatusCode"] == 204

    resp = await s3_client.delete_bucket(Bucket=bucket_name)
    assert resp["ResponseMetadata"]["HTTPStatusCode"] == 204


@pytest.fixture(scope="function")
async def bucket_name(s3_client, create_bucket):
    name = await create_bucket()
    yield name


@pytest.fixture(scope="function")
async def create_bucket(s3_client):
    region_name = "eu-central-1"
    bucket_name = "test-buck"
    await recursive_delete(s3_client, bucket_name)

    async def _f():
        bucket_kwargs = {"Bucket": bucket_name}
        if region_name != "us-east-1":
            bucket_kwargs["CreateBucketConfiguration"] = {
                "LocationConstraint": region_name,
            }
        response = await s3_client.create_bucket(**bucket_kwargs)
        assert response["ResponseMetadata"]["HTTPStatusCode"] == 200
        return bucket_name

    try:
        yield _f
    finally:
        await recursive_delete(s3_client, bucket_name)


@pytest.fixture(scope="function")
async def temp_storage(env_settings):
    async with TemporaryDirectory() as temp_dir:
        yield temp_dir


@pytest.fixture(scope="function")
async def temp_keystore(env_settings, temp_storage):
    pth = Path(temp_storage) / "keys"
    pth.mkdir(parents=True, exist_ok=True)
    for n in range(5):
        private_key = generate_private_key(
            public_exponent=settings.public_exponent,
            key_size=settings.key_size,
        )

        private_key_bytes = private_key.private_bytes(
            encoding=Encoding.PEM,
            format=PrivateFormat.PKCS8,
            encryption_algorithm=NoEncryption(),
        )
        async with aiofiles.open(pth / f"key-{n}.pem", "wb") as fp:
            await fp.write(private_key_bytes)
            await asyncio.sleep(0.1)

    yield pth


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
    ctx.repos = Repositories(
        audit=AuditRepository(ctx.db),
        users=UsersRepository(ctx.db),
        projects=ProjectsRepository(ctx.db),
        permissions=PermissionsRepository(ctx.db),
        world=WorldRepository(ctx.db),
        category=CategoryRepository(ctx.db),
        idp=IDPRepository(ctx.db),
        pets=PetsRepository(ctx.db),
        documents=DocumentsRepository(ctx.db),
        support=SupportRecordsRepository(ctx.db),
        migrations=MigrationRepository(ctx.db),
    )

    ctx.storage = FSStorage(settings.documents_path)
    ctx.keychain = Keychain()
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
            filename="key1.pem",
            hash=str(h.hexdigest())[:16].upper(),
            private_key=load_pem_private_key(
                private_key_bytes,
                password=None,
            ),
        )
    ]
    ctx.jwt_service = JWTService(ctx.keychain.keys[0])
    ctx.mailer = MockMailer()
    ctx.audit_service = AuditService(ctx.repos.audit)
    ctx.crypto_service = CryptoService(ctx.keychain)
    ctx.uploads = UploadHandler(ctx.crypto_service)
    ctx.mfa_service = MFAService(settings.totp_leeway, ctx.crypto_service)
    ctx.users_service = UsersService(ctx.repos.users, ctx.crypto_service)
    ctx.auth_service = AuthService(
        ctx.crypto_service,
        ctx.mfa_service,
        ctx.jwt_service,
        ctx.users_service,
    )
    ctx.secrets_service = SecretsService(ctx.crypto_service)
    ctx.projects_service = ProjectsService(ctx.repos.projects)
    ctx.permissions_service = PermissionsService(ctx.repos.permissions)
    ctx.world_service = WorldService(ctx.repos.world)
    ctx.category_service = CategoryService(ctx.repos.category)
    ctx.idp_service = IDPService(
        ctx.repos.idp,
        ctx.crypto_service,
        ctx.category_service,
        ctx.projects_service,
        ctx.world_service,
        ctx.secrets_service,
    )

    ctx.pets_service = PetsService(ctx.repos.pets)
    ctx.documents_service = DocumentsService(ctx.repos.documents, ctx.crypto_service)
    ctx.support_service = SupportRecordsService(ctx.repos.support)
    ctx.migrations_service = MigrationService(ctx.repos.migrations, ctx.world_service)

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
    country = await app_context.repos.world.create_country(
        NewCountry(
            name="No Stack Country",
            code="nsc",
        )
    )

    yield country


@pytest.fixture(scope="function")
async def default_state(ensure_db, app_context, default_country):  # type:ignore
    state = await app_context.repos.world.create_state(
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
