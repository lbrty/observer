from datetime import datetime
from typing import List

from pydantic import EmailStr, Field, SecretStr

from observer.common.types import Identifier, Role, SomeStr
from observer.schemas.base import SchemaBase


class BaseUser(SchemaBase):
    email: EmailStr = Field(..., description="E-mail address of a user")
    full_name: SomeStr = Field(None, description="Full name of a user")
    role: Role = Field(..., description="Role of a user")


class NewUserRequest(BaseUser):
    password: SecretStr = Field(..., description="Password which user has provided")


class UpdateUserRequest(BaseUser):
    is_active: bool = Field(True, description="Is user active?")
    is_confirmed: bool = Field(False, description="Is user confirmed?")


class UserPasswordUpdate(BaseUser):
    old_password: SecretStr = Field(..., description="Valid old password")
    new_password: SecretStr = Field(..., description="New password")


class UserMFAUpdateRequest(SchemaBase):
    mfa_enabled: bool = Field(False, description="Is MFA enabled for user?")
    mfa_encrypted_secret: SomeStr = Field(None, description="Secret value for MFA")
    mfa_encrypted_backup_codes: SomeStr = Field(None, description="Backup codes for MFA")


class UserResponse(BaseUser):
    id: Identifier = Field(..., description="ID of user")
    ref_id: Identifier = Field(..., description="Reference ID generated using short uuid format")
    is_active: bool = Field(True, description="Is user active?")
    is_confirmed: bool = Field(False, description="Is user confirmed?")
    mfa_enabled: bool = Field(False, description="Is MFA enabled for user?")


class UsersResponse(SchemaBase):
    total: int = Field(..., description="Total count of users")
    items: List[UserResponse] = Field(..., description="List of users")


class UserInviteRequest(SchemaBase):
    email: EmailStr = Field(..., description="User email to send a new invite")


class UserInviteResponse(SchemaBase):
    code: str = Field(..., description="Invite code")
    user_id: Identifier = Field(..., description="User ID")
    expires_at: datetime = Field(..., description="Expiration datetime")
