from fastapi import Cookie

from observer.entities.users import User

Bearer = "Bearer"


async def current_user(auth_token: str | None = Cookie("auth_token")) -> User | None:
    if auth_token:
        # TODO: decode jwt token
        ...
    else:
        return None
