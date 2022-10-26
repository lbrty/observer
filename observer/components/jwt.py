from observer.context import ctx
from observer.services.jwt import JWTHandler


async def jwt_handler() -> JWTHandler:
    return ctx.jwt_handler
