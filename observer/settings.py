from pathlib import Path
from typing import List, Optional

from password_strength import PasswordPolicy
from pydantic import validator
from pydantic.env_settings import BaseSettings
from pydantic.networks import PostgresDsn

from observer.common.types import StorageKind

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
    app_domain: str = "observer.app"

    # Invite only mode
    invite_only: bool = False

    # Looks something like "admin@examples.com,email2@example.com"
    admin_emails: Optional[List[str]] = None
    # Keystore and RSA key settings
    keystore_path: str = "keys"
    key_size: int = 2048
    public_exponent: int = 65537
    aes_key_bits: int = 32

    # User management related settings
    password_reset_url: str = "/reset-password/{code}"
    password_reset_expiration_minutes: int = 15
    confirmation_expiration_minutes: int = 20
    invite_expiration_minutes: int = 15
    invite_url: str = "/account/invites/{code}"
    invite_subject: str = "You are invited to join Observer"
    confirmation_url: str = "/account/confirm/{code}"

    # JWT token expiration durations
    access_token_expiration_minutes: int = 15
    refresh_token_expiration_days: int = 180

    # Allow 10 seconds more for otp codes
    totp_leeway: int = 10
    num_backup_codes: int = 6

    # CORS settings
    cors_origins: Optional[List[str]] = None
    cors_allow_credentials: bool = True

    # gzip settings
    gzip_level: int = 8
    gzip_after_bytes: int = 1024

    # swagger settings
    swagger_output_file: Path = here / "docs/openapi.yml"

    # Mailer settings
    mailer_type: str = "dummy"
    from_email: str = "no-reply@email.com"
    mfa_reset_subject: str = "MFA has been reset"
    password_change_subject: str = "Your password has been updated"

    # Audit log settings
    mfa_event_expiration_days: int = 365
    audit_event_expiration_days: int = 365
    login_event_expiration_days: int = 7
    token_refresh_event_expiration_days: int = 7

    # Storage options
    # Values below are optional exception
    # is for storage backend type.
    # Other settings must be checked manually.
    storage_kind: StorageKind = StorageKind.fs
    storage_root: str = str(here / "uploads")
    # Maximum upload file size is 5Mb
    max_upload_size: int = 1024 * 1024 * 5
    # Local storage uses the same path
    documents_path: str = "documents"
    # Block storage
    s3_endpoint: Optional[str] = "https://s3.aws.amazon.com/observer"
    s3_region: Optional[str] = "eu-central-1"
    s3_bucket: Optional[str] = "observer-keys"

    @property
    def password_policy(self) -> PasswordPolicy:
        """Password strength constraints"""
        return PasswordPolicy.from_names(
            length=8,
            uppercase=1,
            numbers=1,
            nonletters=1,
            strength=0.68,
        )

    @validator("admin_emails", pre=True)
    def validate_admin_emails(cls, val):  # noqa
        return (val or "").split(",")

    @validator("cors_origins", pre=True)
    def validate_cors_origins(cls, val):  # noqa
        if val:
            return val.split(",")

        return ["*"]


class DatabaseSettings(SettingsBase):
    db_uri: PostgresDsn
    pool_size: int = 5
    max_overflow: int = 10
    pool_timeout: int = 30
    echo: bool = False
    echo_pool: bool = False


db_settings = DatabaseSettings()
settings = Settings()
