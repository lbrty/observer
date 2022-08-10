from pathlib import Path

from pydantic.env_settings import BaseSettings
from pydantic.networks import PostgresDsn

here = Path(__file__).parent.parent


class Settings(BaseSettings):
    debug: bool = False
    base_path: Path = here
    db_uri: PostgresDsn

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


settings = Settings()
