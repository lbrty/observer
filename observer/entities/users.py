from datetime import datetime
from typing import Optional

from pydantic import BaseModel, EmailStr

from observer.common.types import Identifier, Role


class User(BaseModel):
    id: Identifier
    email: EmailStr
    full_name: Optional[str] = None
    password_hash: str
    role: Role
    is_active: bool
    is_confirmed: bool
    office_id: Optional[Identifier] = None
    mfa_enabled: bool
    mfa_encrypted_secret: Optional[str] = None
    mfa_encrypted_backup_codes: Optional[str] = None


class NewUser(BaseModel):
    email: EmailStr
    full_name: Optional[str] = None
    password_hash: str
    role: Role
    is_active: bool
    is_confirmed: bool
    office_id: Optional[Identifier] = None


class UserUpdate(BaseModel):
    email: Optional[EmailStr] = None
    full_name: Optional[str] = None
    role: Optional[Role] = None
    is_active: Optional[bool] = None
    office_id: Optional[Identifier] = None
    mfa_enabled: Optional[bool] = None
    mfa_encrypted_secret: Optional[str] = None
    mfa_encrypted_backup_codes: Optional[str] = None


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
