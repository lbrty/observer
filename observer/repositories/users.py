from typing import Protocol

from sqlalchemy import insert, select

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.users import users
from observer.entities.users import NewUser, User


class UsersRepositoryInterface(Protocol):
    async def get_by_id(self, user_id: Identifier) -> User | None:
        raise NotImplementedError

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        raise NotImplementedError

    async def get_by_email(self, email: str) -> User | None:
        raise NotImplementedError

    async def create_user(self, new_user: NewUser) -> User:
        raise NotImplementedError


class UsersRepository(UsersRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def get_by_id(self, user_id: Identifier) -> User | None:
        query = select(users).where(users.c.id == user_id)
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        query = select(users).where(users.c.ref_id == ref_id)
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def get_by_email(self, email: str) -> User | None:
        query = select(users).where(users.c.email == email)
        if result := await self.db.fetchone(query):
            return User(**result)

        return None

    async def create_user(self, new_user: NewUser) -> User:
        query = insert(users).values(**new_user.dict()).returning("*")
        if result := await self.db.fetchone(query):
            return User(**result)
