import sys

from observer.app import create_app
from observer.context import ctx
from observer.db import PoolOptions, connect, disconnect
from observer.repositories.audit_logs import AuditRepository
from observer.repositories.users import UsersRepository
from observer.services.audit_logs import AuditService
from observer.services.auth import AuthService
from observer.services.crypto import CryptoService
from observer.services.jwt import JWTService
from observer.services.keys import get_key_loader
from observer.services.mfa import MFAService
from observer.services.users import UsersService
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

    ctx.keychain = get_key_loader(settings.key_loader_type)
    await ctx.keychain.load(settings.keystore_path)
    num_keys = len(ctx.keychain.keys)
    if num_keys == 0:
        print(f"No keys found, please generate new keys and move to {settings.keystore_path}")
        sys.exit(1)

    print(f"Key loader: {settings.keychain}, Keystore: {settings.keystore_path}, Keys loaded: {num_keys}")
    ctx.jwt_service = JWTService(ctx.keychain.keys[0])
    ctx.audit_repo = AuditRepository(ctx.db)
    ctx.audit_service = AuditService(ctx.audit_repo)
    ctx.crypto_service = CryptoService(ctx.keychain)
    ctx.mfa_service = MFAService(settings.totp_leeway, ctx.crypto_service)
    ctx.users_repo = UsersRepository(ctx.db)
    ctx.users_service = UsersService(ctx.users_repo, ctx.crypto_service)
    ctx.auth_service = AuthService(ctx.crypto_service, ctx.mfa_service, ctx.jwt_service, ctx.users_service)


@app.on_event("shutdown")
async def on_shutdown():
    await disconnect(ctx.db.engine)
