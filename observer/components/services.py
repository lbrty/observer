from observer.context import ctx
from observer.services.users import UsersServiceInterface


async def users_service() -> UsersServiceInterface:
    return ctx.users_service
