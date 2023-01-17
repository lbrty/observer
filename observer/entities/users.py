from datetime import datetime
from typing import Optional

from pydantic import BaseModel, EmailStr

from observer.common.types import Identifier, Role, SomeBool, SomeStr


class User(BaseModel):
    id: Identifier
    ref_id: Identifier
    email: EmailStr
    full_name: SomeStr
    password_hash: str
    role: Role
    is_active: bool
    is_confirmed: bool
    mfa_enabled: bool
    mfa_encrypted_secret: SomeStr
    mfa_encrypted_backup_codes: SomeStr


class NewUser(BaseModel):
    ref_id: Identifier
    email: EmailStr
    full_name: SomeStr
    password_hash: str
    role: Role
    is_active: bool
    is_confirmed: bool


class UserUpdate(BaseModel):
    email: Optional[EmailStr]
    full_name: SomeStr
    role: Optional[Role]
    is_active: SomeBool
    mfa_enabled: SomeBool
    mfa_encrypted_secret: SomeStr
    mfa_encrypted_backup_codes: SomeStr


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
