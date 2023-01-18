from typing import List, Optional

from fastapi import Cookie, Depends

from observer.api.exceptions import (
    ForbiddenError,
    RegistrationsClosedError,
    UnauthorizedError,
)
from observer.common.types import Role
from observer.components import services
from observer.components.services import jwt_service
from observer.entities.base import SomeUser
from observer.entities.users import User
from observer.services.jwt import JWTService
from observer.services.users import IUsersService
from observer.settings import settings


async def current_user(
    access_token: Optional[str] = Cookie(None),
    jwt: JWTService = Depends(jwt_service),
    users_service: IUsersService = Depends(services.users_service),
) -> SomeUser:
    if access_token:
        token_data, _ = await jwt.decode(access_token)
        return await users_service.get_by_ref_id(token_data.ref_id)
    else:
        return None


async def authenticated_user(user: SomeUser = Depends(current_user)) -> User:
    if not user:
        raise UnauthorizedError(message="Please authenticate")

    return user


class RequiresRoles:
    """Check roles authenticated user"""

    def __init__(self, roles: List[Role]):
        self.roles = roles

    async def __call__(self, user: SomeUser = Depends(authenticated_user)) -> User:
        if user is None:
            raise ForbiddenError(message="Access forbidden")

        if user.role not in self.roles:
            raise ForbiddenError(message="Access forbidden")

        return user


async def admin_user(user: SomeUser = Depends(RequiresRoles([Role.admin]))) -> User:
    if not user:
        raise ForbiddenError(message="Access forbidden")

    return user


async def refresh_token_cookie(refresh_token: str = Cookie(None, description="Refresh token cookie")) -> str:
    if refresh_token:
        return refresh_token

    raise UnauthorizedError(message="Refresh token cookie is not set")


async def allow_registrations():
    if settings.invite_only:
        raise RegistrationsClosedError(message="Registrations are not allowed")
