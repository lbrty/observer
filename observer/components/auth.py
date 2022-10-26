from fastapi import Cookie, Depends

from observer.components.jwt import jwt_handler
from observer.context import ctx
from observer.entities.users import User
from observer.services.jwt import JWTHandler


async def current_user(
    auth_token: str | None = Cookie("auth_token"), jwt_service: JWTHandler = Depends(jwt_handler)
) -> User | None:
    if auth_token:
        token_data, _ = await jwt_service.decode(auth_token)
        return await ctx.users_service.get_by_ref_id(token_data.ref_id)
    else:
        return None
