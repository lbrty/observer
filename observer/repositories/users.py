from typing import Protocol

from observer.common.types import Identifier
from observer.entities.users import User


class UsersRepositoryInterface(Protocol):
    async def get_by_id(self, user_id: Identifier) -> User | None:
        raise NotImplementedError

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        raise NotImplementedError


class UsersRepository(UsersRepositoryInterface):
    async def get_by_id(self, user_id: Identifier) -> User | None:
        raise NotImplementedError

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        raise NotImplementedError
