from datetime import datetime
from typing import List, Optional

from pydantic import EmailStr, Field, SecretStr

from observer.common.types import Identifier, Role
from observer.schemas.base import SchemaBase
from observer.schemas.permissions import InvitePermissionRequest


class BaseUser(SchemaBase):
    email: EmailStr = Field(..., description="E-mail address of a user")
    full_name: Optional[str] = Field(None, description="Full name of a user")
    role: Role = Field(..., description="Role of a user")
    office_id: Optional[Identifier] = Field(None, description="Office ID to which user belongs")


class NewUserRequest(BaseUser):
    password: SecretStr = Field(..., description="Password which user has provided")


class UpdateUserRequest(BaseUser):
    is_active: Optional[bool] = Field(..., description="Is user active?")
    is_confirmed: Optional[bool] = Field(..., description="Is user confirmed?")


class UserPasswordUpdate(BaseUser):
    old_password: SecretStr = Field(..., description="Valid old password")
    new_password: SecretStr = Field(..., description="New password")


class UserMFAUpdateRequest(SchemaBase):
    mfa_enabled: bool = Field(False, description="Is MFA enabled for user?")
    mfa_encrypted_secret: Optional[str] = Field(None, description="Secret value for MFA")
    mfa_encrypted_backup_codes: Optional[str] = Field(None, description="Backup codes for MFA")


class UserResponse(BaseUser):
    id: Identifier = Field(..., description="ID of user")
    is_active: bool = Field(True, description="Is user active?")
    is_confirmed: bool = Field(False, description="Is user confirmed?")
    mfa_enabled: bool = Field(False, description="Is MFA enabled for user?")


class UsersResponse(SchemaBase):
    total: int = Field(..., description="Total count of users")
    items: List[UserResponse] = Field(..., description="List of users")


class UserInviteRequest(SchemaBase):
    email: EmailStr = Field(..., description="User email to send a new invite")
    role: Role = Field(..., description="Role of a user")
    permissions: Optional[List[InvitePermissionRequest]] = Field(
        ..., description="List of permissions to different projects"
    )


class UserInviteResponse(SchemaBase):
    code: str = Field(..., description="Invite code")
    user_id: Identifier = Field(..., description="User ID")
    expires_at: datetime = Field(..., description="Expiration datetime")


class UserInvitesResponse(SchemaBase):
    total: int = Field(..., description="Total amount of invites")
    items: List[UserInviteResponse] = Field(..., description="List of invites DESC ordered by expiration date")


class InviteJoinRequest(SchemaBase):
    password: SecretStr = Field(..., description="Password which user has provided")


# Admin schemas
class CreateUserRequest(BaseUser):
    password: SecretStr = Field(..., description="Password which user has provided")
    is_active: bool = Field(True, description="Is user active?")


class AdminUpdateUserRequest(BaseUser):
    is_active: Optional[bool] = Field(..., description="Is user active?")
