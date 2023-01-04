import sys

from observer.app import create_app
from observer.context import Repositories, ctx
from observer.db import PoolOptions, connect, disconnect
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
from observer.services.audit_logs import AuditService
from observer.services.auth import AuthService
from observer.services.categories import CategoryService
from observer.services.crypto import CryptoService
from observer.services.idp import IDPService
from observer.services.jwt import JWTService
from observer.services.keys import Keychain
from observer.services.mailer import Mailer
from observer.services.mfa import MFAService
from observer.services.migration_history import MigrationService
from observer.services.permissions import PermissionsService
from observer.services.projects import ProjectsService
from observer.services.secrets import SecretsService
from observer.services.storage import init_storage
from observer.services.users import UsersService
from observer.services.world import WorldService
from observer.settings import db_settings, settings

app = create_app(settings)


@app.on_event("startup")
async def on_startup():
    ctx.db = await connect(
        db_settings.db_uri,
        PoolOptions(
            pool_size=db_settings.pool_size,
            pool_timeout=db_settings.pool_timeout,
            echo=db_settings.echo,
            echo_pool=db_settings.echo_pool,
            max_overflow=db_settings.max_overflow,
        ),
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

    ctx.storage = init_storage(settings.storage_kind, settings)
    ctx.keychain = Keychain(ctx.storage)
    await ctx.keychain.load(settings.keystore_path)
    num_keys = len(ctx.keychain.keys)
    if num_keys == 0:
        print(f"No keys found, please generate new keys and move to {settings.keystore_path}")
        sys.exit(1)

    print(f"Key loader: {settings.key_loader_type}, Keystore: {settings.keystore_path}, Keys loaded: {num_keys}")
    ctx.jwt_service = JWTService(ctx.keychain.keys[0])
    ctx.mailer = Mailer()
    ctx.audit_service = AuditService(ctx.repos.audit)
    ctx.crypto_service = CryptoService(ctx.keychain)
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
    ctx.migrations_service = MigrationService(ctx.repos.migrations, ctx.world_service)


@app.on_event("shutdown")
async def on_shutdown():
    await disconnect(ctx.db.engine)
