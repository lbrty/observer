from typing import Protocol

from sqlalchemy import insert, select, update

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.users import users
from observer.entities.base import SomeUser
from observer.entities.users import NewUser, User, UserUpdate


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


class UsersRepository(UsersRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def get_by_id(self, user_id: Identifier) -> SomeUser:
        query = select(users).where(users.c.id == user_id)
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

        query = update(users).values(update_values).where(users.c.id == user_id).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)
