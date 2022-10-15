from pydantic import EmailStr, Field, SecretStr

from observer.common.types import Identifier, Role
from observer.schemas.base import SchemaBase


class BaseUser(SchemaBase):
    email: EmailStr = Field(..., description="E-mail address of a user")
    full_name: str | None = Field(None, description="Full name of a user")
    role: Role = Field(..., description="Role of a user")


class NewUser(BaseUser):
    password: SecretStr = Field(..., description="Password which user has provided")


class UpdateUser(BaseUser):
    is_active: bool = Field(True, description="Is user active?")
    is_confirmed: bool = Field(False, description="Is user confirmed?")


class UserPasswordUpdate(BaseUser):
    old_password: SecretStr = Field(..., description="Valid old password")
    new_password: SecretStr = Field(..., description="New password")


class UserMFAUpdate(SchemaBase):
    mfa_enabled: bool = Field(False, description="Is MFA enabled for user?")
    mfa_encrypted_secret: SecretStr | None = Field(None, description="Secret value for MFA")
    mfa_encrypted_backup_codes: SecretStr | None = Field(None, description="Backup codes for MFA")


class User(BaseUser):
    id: Identifier = Field(..., description="ID of user")
    is_active: bool = Field(True, description="Is user active?")
    is_confirmed: bool = Field(False, description="Is user confirmed?")
    mfa_enabled: bool = Field(False, description="Is MFA enabled for user?")
    mfa_encrypted_secret: SecretStr | None = Field(None, description="Secret value for MFA")
    mfa_encrypted_backup_codes: SecretStr | None = Field(None, description="Backup codes for MFA")


class UsersResponse(SchemaBase):
    total: int = Field(..., description="Total count of users")
    items: list[User] = Field(..., description="List of users")
