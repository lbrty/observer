import asyncio

import pytest
import pytest_asyncio
from dotenv import load_dotenv
from fastapi import FastAPI


@pytest.fixture(scope="session")
def event_loop():
    policy = asyncio.get_event_loop_policy()
    loop = policy.new_event_loop()
    yield loop
    loop.close()


@pytest.fixture(scope="session")
def env_settings():
    return load_dotenv(".env.test")


@pytest_asyncio.fixture(scope="session")
async def app(env_settings) -> FastAPI:
    pass


@pytest_asyncio.fixture(scope="session")
async def db(env_settings):
    pass


@pytest_asyncio.fixture(scope="session")
async def client(env_settings, db, app):
    pass
