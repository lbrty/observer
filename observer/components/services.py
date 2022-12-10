from observer.api.exceptions import InternalError
from observer.context import ctx
from observer.services.audit_logs import AuditServiceInterface
from observer.services.auth import AuthServiceInterface
from observer.services.crypto import CryptoServiceInterface
from observer.services.jwt import JWTService
from observer.services.keys import Keychain
from observer.services.mailer import MailerInterface
from observer.services.users import UsersServiceInterface


async def users_service() -> UsersServiceInterface:
    return ctx.users_service


async def auth_service() -> AuthServiceInterface:
    return ctx.auth_service


async def jwt_service() -> JWTService:
    return ctx.jwt_service


async def crypto_service() -> CryptoServiceInterface:
    return ctx.crypto_service


async def keychain() -> Keychain:
    if len(ctx.keychain.keys) == 0:
        raise InternalError(message="private keys not found")

    return ctx.keychain


async def audit_service() -> AuditServiceInterface:
    return ctx.audit_service


async def mailer() -> MailerInterface:
    return ctx.mailer
