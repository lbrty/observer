from observer.context import Context
from observer.entities.users import User
from observer.schemas.auth import TokenResponse


async def get_auth_tokens(ctx: Context, user: User) -> TokenResponse:
    return await ctx.auth_service.create_token(user.ref_id)
