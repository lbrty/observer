from pathlib import Path

from pydantic.env_settings import BaseSettings
from pydantic.networks import PostgresDsn

here = Path(__file__).parent.parent


class Settings(BaseSettings):
    debug: bool = False
    base_path: Path = here
    db_uri: PostgresDsn

    # Keystore and RSA key settings
    key_store_path: Path = here / "keys"
    key_size: int = 2048
    public_exponent: int = 65537

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


settings = Settings()
