from pathlib import Path
from typing import List, Optional

from password_strength import PasswordPolicy
from pydantic.env_settings import BaseSettings
from pydantic.networks import PostgresDsn

from observer.schemas.crypto import KeyLoaderTypes

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
    version: str = "0.1.0"
    app_domain: str = "https://observer.app"

    # Keystore and RSA key settings
    key_loader_type: KeyLoaderTypes = KeyLoaderTypes.fs
    keystore_path: Path = here / "keys"
    key_size: int = 2048
    key_passwords: Optional[str] = None
    public_exponent: int = 65537
    aes_key_bits: int = 32
    # Password strength constraints
    password_policy: PasswordPolicy = PasswordPolicy.from_names(
        length=8,
        uppercase=1,
        numbers=1,
        nonletters=1,
        strength=0.68,
    )
    password_reset_url: str = "/reset-password/{code}"
    password_reset_expiration_minutes: int = 15
    confirmation_expiration_minutes: int = 20
    confirmation_url: str = "/account/confirm/{code}"

    # Allow 10 seconds more for otp codes
    totp_leeway: int = 10
    num_backup_codes: int = 6

    # CORS settings
    cors_origins: List[str] = ["*"]
    cors_allow_credentials: bool = True

    # gzip settings
    gzip_level: int = 8
    gzip_after_bytes: int = 1024

    # swagger settings
    swagger_output_file: Path = here / "docs/openapi.yml"

    # Mailer settings
    from_email: str = "no-reply@email.com"
    mfa_reset_subject: str = "MFA has been reset"
    auth_password_change_subject: str = "Your password has been updated"

    # Audit log settings
    mfa_audit_event_lifetime_days: int = 365
    auth_audit_event_lifetime_days: int = 365
    auth_audit_event_login_days: int = 7
    auth_audit_event_refresh_days: int = 7

    # Storage options
    # Values below are optional exception
    # is for storage backend type.
    # Other settings must be checked manually.
    storage_backend: str = "fs"
    # Local storage
    fs_keys_path: Optional[str] = str(here / "keys")
    fs_documents_path: Optional[str] = str(here / "documents")
    # Block storage
    s3_endpoint: Optional[str] = "https://s3.aws.amazon.com/observer"
    s3_keys_path: Optional[str] = "keys"
    s3_documents_path: Optional[str] = "documents"


class DatabaseSettings(SettingsBase):
    db_uri: PostgresDsn
    pool_size: int = 5
    max_overflow: int = 10
    pool_timeout: int = 30
    echo: bool = False
    echo_pool: bool = False


db_settings = DatabaseSettings()
settings = Settings()
