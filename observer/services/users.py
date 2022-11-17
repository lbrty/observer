from typing import Protocol

from observer.common.types import Identifier
from observer.entities.users import User
from observer.repositories.users import UsersRepositoryInterface
from observer.schemas import users


class UsersServiceInterface(Protocol):
    repo: UsersRepositoryInterface

    async def get_by_id(self, user_id: Identifier) -> User | None:
        raise NotImplementedError

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        raise NotImplementedError

    async def get_by_email(self, email: str) -> User | None:
        raise NotImplementedError

    @staticmethod
    async def to_response(user: User) -> users.User:
        raise NotImplementedError

    @staticmethod
    async def list_to_response(user_list: list[User]) -> list[users.User]:
        raise NotImplementedError


class UsersService(Protocol):
    def __init__(self, users_repository: UsersRepositoryInterface):
        self.repo = users_repository

    async def get_by_id(self, user_id: Identifier) -> User | None:
        ...

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        return await self.repo.get_by_ref_id(ref_id)

    async def get_by_email(self, email: str) -> User | None:
        return await self.repo.get_by_email(email)

    @staticmethod
    async def to_response(user: User) -> users.User:
        return users.User(**user.dict())

    @staticmethod
    async def list_to_response(user_list: list[User]) -> list[users.User]:
        return [users.User(**user.dict()) for user in user_list]
