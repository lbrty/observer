from fastapi import Cookie, Depends

from observer.components import services
from observer.components.jwt import jwt_handler
from observer.entities.users import User
from observer.services.jwt import JWTHandler
from observer.services.users import UsersServiceInterface


async def current_user(
    auth_token: str | None = Cookie("auth_token"),
    jwt_service: JWTHandler = Depends(jwt_handler),
    users_service: UsersServiceInterface = Depends(services.users_service),
) -> User | None:
    if auth_token:
        token_data, _ = await jwt_service.decode(auth_token)
        return await users_service.get_by_ref_id(token_data.ref_id)
    else:
        return None
