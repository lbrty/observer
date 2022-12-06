from typing import List

from fastapi import Cookie, Depends

from observer.api.exceptions import ForbiddenError, UnauthorizedError
from observer.common.types import Role
from observer.components import services
from observer.components.jwt import jwt_handler
from observer.entities.base import SomeUser
from observer.entities.users import User
from observer.services.jwt import JWTService
from observer.services.users import UsersServiceInterface


async def current_user(
    access_token: str | None = Cookie(None),
    jwt_service: JWTService = Depends(jwt_handler),
    users_service: UsersServiceInterface = Depends(services.users_service),
) -> SomeUser:
    if access_token:
        token_data, _ = await jwt_service.decode(access_token)
        return await users_service.get_by_ref_id(token_data.ref_id)
    else:
        return None


async def authenticated_user(user: SomeUser = Depends(current_user)) -> User:
    if not user:
        raise UnauthorizedError(message="Please authenticate")

    return user


class RequiredRoles:
    """Check roles authenticated user"""

    def __init__(self, roles: List[Role]):
        self.roles = roles

    async def __call__(self, user: SomeUser = Depends(authenticated_user), **kwargs):
        if user.role not in self.roles:
            raise ForbiddenError(message="Access forbidden")


async def admin_user(user: SomeUser = Depends(RequiredRoles([Role.admin]))) -> User:
    if not user:
        raise ForbiddenError(message="Access forbidden")

    return user


async def refresh_token_cookie(refresh_token: str = Cookie(None, description="Refresh token cookie")) -> str:
    if refresh_token:
        return refresh_token

    raise UnauthorizedError(message="Refresh token cookie is not set")
