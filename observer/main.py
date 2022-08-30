from observer.app import create_app
from observer.context import ctx
from observer.db import connect, PoolOptions, disconnect
from observer.services.crypto import get_key_loader
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

    ctx.key_loader_type = get_key_loader(settings.key_loader_type)
    await ctx.key_loader_type.load(settings.keystore_path)
    print(
        f"Key loader: {settings.key_loader_type}, Keystore: {settings.keystore_path}, Keys loaded:"
        f" {len(ctx.key_loader_type.keys)}"
    )


@app.on_event("shutdown")
async def on_shutdown():
    await disconnect(ctx.db.engine)
