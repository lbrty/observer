from datetime import datetime
from typing import Optional

from pydantic import BaseModel, EmailStr

from observer.common.types import Identifier, Role


class User(BaseModel):
    id: Identifier
    email: EmailStr
    full_name: Optional[str]
    password_hash: str
    role: Role
    is_active: bool
    is_confirmed: bool
    office_id: Optional[Identifier]
    mfa_enabled: bool
    mfa_encrypted_secret: Optional[str]
    mfa_encrypted_backup_codes: Optional[str]


class NewUser(BaseModel):
    email: EmailStr
    full_name: Optional[str]
    password_hash: str
    role: Role
    is_active: bool
    is_confirmed: bool
    office_id: Optional[Identifier]


class UserUpdate(BaseModel):
    email: Optional[EmailStr]
    full_name: Optional[str]
    role: Optional[Role]
    is_active: Optional[bool]
    office_id: Optional[Identifier]
    mfa_enabled: Optional[bool]
    mfa_encrypted_secret: Optional[str]
    mfa_encrypted_backup_codes: Optional[str]


class PasswordReset(BaseModel):
    code: str
    user_id: Identifier
    created_at: datetime


class Confirmation(BaseModel):
    code: str
    user_id: Identifier
    expires_at: datetime


class Invite(BaseModel):
    code: str
    user_id: Identifier
    expires_at: datetime
