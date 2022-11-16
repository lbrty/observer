from pathlib import Path
from typing import List, Optional

from pydantic.env_settings import BaseSettings
from pydantic.networks import PostgresDsn

from observer.schemas.crypto import KeyLoaderTypes
from observer.services.crypto import FSLoader, KeychainLoader

here = Path(__file__).parent.parent


class SettingsBase(BaseSettings):
    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


class Settings(SettingsBase):
    debug: bool = False
    port: int = 3000
    base_path: Path = here

    # OpenAPI
    title: str = "Observer API"
    description: str = "Observer API server"

    # Keystore and RSA key settings
    key_loader: KeychainLoader = FSLoader
    key_loader_type: KeyLoaderTypes = KeyLoaderTypes.fs
    keystore_path: Path = here / "keys"
    key_size: int = 2048
    key_passwords: Optional[str] = None
    public_exponent: int = 65537

    # CORS settings
    cors_origins: List[str] = ["*"]
    cors_allow_credentials: bool = True

    # gzip settings
    gzip_level: int = 8
    gzip_after_bytes: int = 1024

    # swagger settings
    swagger_output_file: Path = here / "docs/openapi.yml"


class DatabaseSettings(SettingsBase):
    db_uri: PostgresDsn
    pool_size: int = 5
    max_overflow: int = 10
    pool_timeout: int = 30
    echo: bool = False
    echo_pool: bool = False


db_settings = DatabaseSettings()
settings = Settings()
