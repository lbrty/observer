from dataclasses import dataclass

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from pytoolz.functional import pipe

from observer.api import (
    account,
    auth,
    health,
    idp,
    mfa,
    migration_history,
    projects,
    world,
)
from observer.api.exceptions import BaseAPIException
from observer.common.exceptions import handle_api_exception
from observer.settings import Settings

__all__ = ("create_app",)


@dataclass
class Environment:
    settings: Settings
    app: FastAPI


def init_integrations(env: Environment) -> Environment:
    """Integrations can be for example Sentry, Datadog etc."""
    return env


def init_routes(env: Environment) -> Environment:
    env.app.include_router(account.router)
    env.app.include_router(auth.router)
    env.app.include_router(idp.router)
    env.app.include_router(health.router)
    env.app.include_router(migration_history.router)
    env.app.include_router(mfa.router)
    env.app.include_router(world.router)
    env.app.include_router(projects.router)
    return env


def init_middlewares(env: Environment) -> Environment:
    env.app.add_middleware(
        CORSMiddleware,
        allow_origins=env.settings.cors_origins,
        allow_credentials=env.settings.cors_allow_credentials,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    env.app.add_middleware(
        GZipMiddleware,
        minimum_size=env.settings.gzip_after_bytes,
        compresslevel=env.settings.gzip_level,
    )

    return env


def init_exception_handlers(env: Environment) -> Environment:
    env.app.add_exception_handler(BaseAPIException, handle_api_exception)
    return env


def create_app(settings: Settings) -> FastAPI:
    env: Environment = pipe(
        [
            init_integrations,
            init_routes,
            init_integrations,
            init_exception_handlers,
        ],
        Environment(
            settings=settings,
            app=FastAPI(
                debug=settings.debug,
                title=settings.title,
                description=settings.description,
                version=settings.version,
                redoc_url="/docs",
            ),
        ),
    )

    return env.app
