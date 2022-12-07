from observer.context import ctx
from observer.services.jwt import JWTService


async def jwt_service() -> JWTService:
    return ctx.jwt_service
