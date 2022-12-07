from observer.context import ctx
from observer.services.crypto import CryptoServiceInterface
from observer.services.jwt import JWTService
from observer.services.users import UsersServiceInterface


async def users_service() -> UsersServiceInterface:
    return ctx.users_service


async def jwt_service() -> JWTService:
    return ctx.jwt_service


async def crypto_service() -> CryptoServiceInterface:
    return ctx.crypto_service
