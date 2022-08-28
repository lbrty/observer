from pathlib import Path
from typing import Optional, List

from pydantic.env_settings import BaseSettings
from pydantic.networks import PostgresDsn

from observer.services.crypto import KeyLoaderTypes

here = Path(__file__).parent.parent


class Settings(BaseSettings):
    debug: bool = False
    port: int = 3000
    base_path: Path = here
    db_uri: PostgresDsn

    # OpenAPI
    title: str = "Observer API"
    description: str = "Observer API server"

    # Keystore and RSA key settings
    key_loader: KeyLoaderTypes = KeyLoaderTypes.fs
    key_store_path: Path = here / "keys"
    key_size: int = 2048
    key_passwords: Optional[str] = None
    public_exponent: int = 65537

    # CORS settings
    cors_origins: List[str] = ["*"]
    cors_allow_credentials: bool = True

    # gzip settings
    gzip_level: int = 8
    gzip_after_bytes: int = 1024

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


settings = Settings()
