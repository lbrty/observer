from typing import Protocol

from sqlalchemy import select

from observer.common.types import Identifier
from observer.db import Database
from observer.db.tables.users import users
from observer.entities.users import User


class UsersRepositoryInterface(Protocol):
    async def get_by_id(self, user_id: Identifier) -> User | None:
        raise NotImplementedError

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        raise NotImplementedError


class UsersRepository(UsersRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def get_by_id(self, user_id: Identifier) -> User | None:
        with self.db.session as conn:
            query = select(users).where(users.c.id == user_id)

            if result := await conn.execute(query):
                row = result.fetchone()
                return User(**row)

        return None

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        with self.db.session as conn:
            query = select(users).where(users.c.ref_id == ref_id)
            if result := await conn.execute(query):
                row = result.fetchone()
                return User(**row)

        return None
