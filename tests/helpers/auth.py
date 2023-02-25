from observer.context import Context
from observer.entities.users import User
from observer.schemas.auth import TokenResponse


async def get_auth_tokens(ctx: Context, user: User) -> TokenResponse:
    if ctx.auth_service:
        return await ctx.auth_service.create_token(user.id)

    raise RuntimeError("AuthService is None")
