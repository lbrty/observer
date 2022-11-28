from typing import Protocol

import shortuuid

from observer.common import bcrypt
from observer.common.types import Identifier
from observer.entities.users import NewUser, User
from observer.repositories.users import UsersRepositoryInterface
from observer.schemas.users import NewUserRequest, UserResponse, UsersResponse


class UsersServiceInterface(Protocol):
    repo: UsersRepositoryInterface

    async def get_by_id(self, user_id: Identifier) -> User | None:
        raise NotImplementedError

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        raise NotImplementedError

    async def get_by_email(self, email: str) -> User | None:
        raise NotImplementedError

    async def create_user(self, new_user: NewUserRequest) -> User:
        raise NotImplementedError

    @staticmethod
    async def to_response(user: User) -> UserResponse:
        raise NotImplementedError

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> UsersResponse:
        raise NotImplementedError


class UsersService(UsersServiceInterface):
    def __init__(self, users_repository: UsersRepositoryInterface):
        self.repo = users_repository

    async def get_by_id(self, user_id: Identifier) -> User | None:
        ...

    async def get_by_ref_id(self, ref_id: Identifier) -> User | None:
        return await self.repo.get_by_ref_id(ref_id)

    async def get_by_email(self, email: str) -> User | None:
        return await self.repo.get_by_email(email)

    async def create_user(self, new_user: NewUserRequest) -> User:
        ref_id = shortuuid.uuid(name=new_user.email)
        password_hash = bcrypt.hash_password(new_user.password.get_secret_value())
        user = NewUser(
            ref_id=ref_id,
            email=new_user.email,
            full_name=new_user.full_name,
            password_hash=password_hash,
            role=new_user.role,
            is_active=True,
            is_confirmed=False,
        )
        return await self.repo.create_user(user)

    @staticmethod
    async def to_response(user: User) -> UserResponse:
        return UserResponse(**user.dict())

    @staticmethod
    async def list_to_response(total: int, user_list: list[User]) -> UsersResponse:
        return UsersResponse(
            total=total,
            items=[UserResponse(**user.dict()) for user in user_list],
        )
