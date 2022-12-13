from datetime import datetime, timezone
from typing import Protocol

from sqlalchemy import insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.users import confirmations, password_resets, users
from observer.entities.base import SomeUser
from observer.entities.users import (
    Confirmation,
    NewUser,
    PasswordReset,
    User,
    UserUpdate,
)


class UsersRepositoryInterface(Protocol):
    async def get_by_id(self, user_id: Identifier) -> SomeUser:
        raise NotImplementedError

    async def get_by_ref_id(self, ref_id: Identifier) -> SomeUser:
        raise NotImplementedError

    async def get_by_email(self, email: str) -> SomeUser:
        raise NotImplementedError

    async def create_user(self, new_user: NewUser) -> User:
        raise NotImplementedError

    async def update_user(self, user_id: Identifier, updates: UserUpdate) -> User:
        raise NotImplementedError

    async def update_password(self, user_id: Identifier, new_password_hash: str) -> User:
        raise NotImplementedError

    async def reset_mfa(self, user_id: Identifier) -> User:
        raise NotImplementedError

    async def create_password_reset_code(self, user_id: Identifier, code: str) -> PasswordReset:
        raise NotImplementedError

    async def get_password_reset(self, code: str) -> PasswordReset | None:
        raise NotImplementedError

    async def create_confirmation(self, user_id: Identifier, code: str, expires_at: datetime) -> Confirmation:
        raise NotImplementedError

    async def get_confirmation(self, code: str) -> Confirmation | None:
        raise NotImplementedError

    async def confirm_user(self, user_id: Identifier) -> User:
        raise NotImplementedError


class UsersRepository(UsersRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def get_by_id(self, user_id: Identifier) -> SomeUser:
        query = select(users).where(users.c.id == str(user_id))
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def get_by_ref_id(self, ref_id: Identifier) -> SomeUser:
        query = select(users).where(users.c.ref_id == ref_id)
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def get_by_email(self, email: str) -> SomeUser:
        query = select(users).where(users.c.email == email)
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def create_user(self, new_user: NewUser) -> User:
        query = insert(users).values(**new_user.dict()).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)

    async def create_confirmation(self, user_id: Identifier, code: str, expires_at: datetime) -> Confirmation:
        values = dict(
            code=code,
            user_id=str(user_id),
            expires_at=expires_at,
        )
        query = insert(confirmations).values(**values).returning("*")
        if result := await self.db.fetchone(query):
            return Confirmation(**result)

    async def get_confirmation(self, code: str) -> Confirmation | None:
        query = select(confirmations).where(confirmations.c.code == code)
        if result := await self.db.fetchone(query):
            return Confirmation(**result)

        return None

    async def confirm_user(self, user_id: Identifier) -> User:
        query = update(users).values(dict(is_confirmed=True)).where(users.c.id == str(user_id)).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)

    async def update_user(self, user_id: Identifier, updates: UserUpdate) -> User:
        update_values = {}
        if updates.email:
            update_values["email"] = updates.email

        if updates.full_name:
            update_values["full_name"] = updates.full_name

        if updates.role:
            update_values["role"] = updates.role

        if updates.is_active:
            update_values["is_active"] = updates.is_active

        if updates.mfa_enabled is not None and updates.mfa_encrypted_secret and updates.mfa_encrypted_backup_codes:
            update_values["mfa_enabled"] = updates.mfa_enabled
            update_values["mfa_encrypted_secret"] = updates.mfa_encrypted_secret
            update_values["mfa_encrypted_backup_codes"] = updates.mfa_encrypted_backup_codes

        query = update(users).values(update_values).where(users.c.id == str(user_id)).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)

    async def update_password(self, user_id: Identifier, new_password_hash: str) -> User:
        update_values = dict(password_hash=new_password_hash)
        query = update(users).values(update_values).where(users.c.id == str(user_id)).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)

    async def reset_mfa(self, user_id: Identifier) -> User:
        update_values = dict(
            mfa_enabled=False,
            mfa_encrypted_secret=None,
            mfa_encrypted_backup_codes=None,
        )

        query = update(users).values(update_values).where(users.c.id == str(user_id)).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)

    async def create_password_reset_code(self, user_id: Identifier, code: str) -> PasswordReset:
        values = dict(
            code=code,
            user_id=str(user_id),
            created_at=datetime.now(tz=timezone.utc),
        )
        query = insert(password_resets).values(**values).returning("*")
        if result := await self.db.fetchone(query):
            return PasswordReset(**result)

    async def get_password_reset(self, code: str) -> PasswordReset | None:
        query = select(password_resets).where(password_resets.c.code == code)
        if result := await self.db.fetchone(query):
            return PasswordReset(**result)

        return None
