from datetime import datetime, timezone
from typing import Optional, Protocol

from sqlalchemy import delete, insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.users import confirmations, invites, password_resets, users
from observer.entities.base import SomeUser
from observer.entities.users import (
    Confirmation,
    Invite,
    NewUser,
    PasswordReset,
    User,
    UserUpdate,
)


class IUsersRepository(Protocol):
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

    async def get_password_reset(self, code: str) -> Optional[PasswordReset]:
        raise NotImplementedError

    async def create_confirmation(self, user_id: Identifier, code: str, expires_at: datetime) -> Confirmation:
        raise NotImplementedError

    async def get_confirmation(self, code: str) -> Optional[Confirmation]:
        raise NotImplementedError

    async def confirm_user(self, user_id: Identifier) -> User:
        raise NotImplementedError

    async def create_invite(self, user_id: Identifier, code: str, expires_at: datetime) -> Invite:
        raise NotImplementedError

    async def get_invite(self, code: str) -> Optional[Invite]:
        raise NotImplementedError

    async def delete_invite(self, code: str) -> Invite:
        raise NotImplementedError


class UsersRepository(IUsersRepository):
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
        result = await self.db.fetchone(query)
        return User(**result)

    async def create_confirmation(self, user_id: Identifier, code: str, expires_at: datetime) -> Confirmation:
        values = dict(
            code=code,
            user_id=str(user_id),
            expires_at=expires_at,
        )
        query = insert(confirmations).values(**values).returning("*")
        result = await self.db.fetchone(query)
        return Confirmation(**result)

    async def get_confirmation(self, code: str) -> Optional[Confirmation]:
        query = select(confirmations).where(confirmations.c.code == code)
        if result := await self.db.fetchone(query):
            return Confirmation(**result)

        return None

    async def create_invite(self, user_id: Identifier, code: str, expires_at: datetime) -> Invite:
        values = dict(
            code=code,
            user_id=str(user_id),
            expires_at=expires_at,
        )
        query = insert(invites).values(**values).returning("*")
        result = await self.db.fetchone(query)
        return Invite(**result)

    async def get_invite(self, code: str) -> Optional[Invite]:
        query = select(invites).where(invites.c.code == code)
        if result := await self.db.fetchone(query):
            return Invite(**result)

        return None

    async def delete_invite(self, code: str) -> Invite:
        query = delete(invites).where(invites.c.code == code).returning("*")
        result = await self.db.fetchone(query)
        return Invite(**result)

    async def confirm_user(self, user_id: Identifier) -> User:
        query = update(users).values(dict(is_confirmed=True)).where(users.c.id == str(user_id)).returning("*")
        result = await self.db.fetchone(query)
        return User(**result)

    async def update_user(self, user_id: Identifier, updates: UserUpdate) -> User:
        update_values = {}
        if updates.email:
            update_values["email"] = updates.email

        if updates.full_name:
            update_values["full_name"] = updates.full_name  # type:ignore

        if updates.role:
            update_values["role"] = updates.role.value  # type:ignore

        if updates.is_active:
            update_values["is_active"] = updates.is_active  # type:ignore

        if updates.mfa_enabled is not None and updates.mfa_encrypted_secret and updates.mfa_encrypted_backup_codes:
            update_values["mfa_enabled"] = updates.mfa_enabled  # type:ignore
            update_values["mfa_encrypted_secret"] = updates.mfa_encrypted_secret  # type:ignore
            update_values["mfa_encrypted_backup_codes"] = updates.mfa_encrypted_backup_codes  # type:ignore

        query = update(users).values(update_values).where(users.c.id == str(user_id)).returning("*")
        result = await self.db.fetchone(query)
        return User(**result)

    async def update_password(self, user_id: Identifier, new_password_hash: str) -> User:
        update_values = dict(password_hash=new_password_hash)
        query = update(users).values(update_values).where(users.c.id == str(user_id)).returning("*")
        result = await self.db.fetchone(query)
        return User(**result)

    async def reset_mfa(self, user_id: Identifier) -> User:
        update_values = dict(
            mfa_enabled=False,
            mfa_encrypted_secret=None,
            mfa_encrypted_backup_codes=None,
        )

        query = update(users).values(update_values).where(users.c.id == str(user_id)).returning("*")
        result = await self.db.fetchone(query)
        return User(**result)

    async def create_password_reset_code(self, user_id: Identifier, code: str) -> PasswordReset:
        values = dict(
            code=code,
            user_id=str(user_id),
            created_at=datetime.now(tz=timezone.utc),
        )
        query = insert(password_resets).values(**values).returning("*")
        result = await self.db.fetchone(query)
        return PasswordReset(**result)

    async def get_password_reset(self, code: str) -> Optional[PasswordReset]:
        query = select(password_resets).where(password_resets.c.code == code)
        if result := await self.db.fetchone(query):
            return PasswordReset(**result)

        return None
