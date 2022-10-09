from pydantic import BaseModel, EmailStr

from observer.common.types import Identifier, Role


class User(BaseModel):
    id: Identifier
    email: EmailStr
    full_name: str | None
    password_hash: str
    role: Role
    is_active: bool
    is_confirmed: bool
    mfa_enabled: bool
    mfa_encrypted_secret: str | None
    mfa_encrypted_backup_codes: str | None
