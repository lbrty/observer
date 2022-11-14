import pytest
from fastapi import FastAPI
from dotenv import load_dotenv


@pytest.fixture(scope="session")
def env_settings():
    return load_dotenv(".env.test")


@pytest.fixture(scope="session")
async def app(env_settings) -> FastAPI:
    pass


@pytest.fixture(scope="session")
async def db(env_settings):
    pass


@pytest.fixture(scope="session")
async def client(env_settings, db, app):
    pass
