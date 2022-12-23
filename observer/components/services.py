from observer.api.exceptions import InternalError
from observer.context import ctx
from observer.services.audit_logs import AuditServiceInterface
from observer.services.auth import AuthServiceInterface
from observer.services.categories import CategoryServiceInterface
from observer.services.crypto import CryptoServiceInterface
from observer.services.idp import IDPServiceInterface
from observer.services.jwt import JWTService
from observer.services.keys import Keychain
from observer.services.mailer import MailerInterface
from observer.services.permissions import PermissionsServiceInterface
from observer.services.projects import ProjectsServiceInterface
from observer.services.users import UsersServiceInterface
from observer.services.world import WorldServiceInterface


async def users_service() -> UsersServiceInterface:
    if ctx.users_service:
        return ctx.users_service

    raise RuntimeError("UsersService is None")


async def projects_service() -> ProjectsServiceInterface:
    if ctx.projects_service:
        return ctx.projects_service

    raise RuntimeError("ProjectsService is None")


async def permissions_service() -> PermissionsServiceInterface:
    if ctx.permissions_service:
        return ctx.permissions_service

    raise RuntimeError("PermissionsService is None")


async def auth_service() -> AuthServiceInterface:
    if ctx.auth_service:
        return ctx.auth_service

    raise RuntimeError("AuthService is None")


async def jwt_service() -> JWTService:
    if ctx.jwt_service:
        return ctx.jwt_service

    raise RuntimeError("JWTService is None")


async def crypto_service() -> CryptoServiceInterface:
    if ctx.crypto_service:
        return ctx.crypto_service

    raise RuntimeError("CryptoService is None")


async def keychain() -> Keychain:
    if not ctx.keychain:
        raise RuntimeError("Keychain is None")

    if len(ctx.keychain.keys) == 0:
        raise InternalError(message="private keys not found")

    return ctx.keychain


async def audit_service() -> AuditServiceInterface:
    if ctx.audit_service:
        return ctx.audit_service

    raise RuntimeError("AuditService is None")


async def world_service() -> WorldServiceInterface:
    if ctx.world_service:
        return ctx.world_service

    raise RuntimeError("PlacesService is None")


async def idp_service() -> IDPServiceInterface:
    if ctx.idp_service:
        return ctx.idp_service

    raise RuntimeError("IDPService is None")


async def category_service() -> CategoryServiceInterface:
    if ctx.category_service:
        return ctx.category_service

    raise RuntimeError("CategoryService is None")


async def mailer() -> MailerInterface:
    if ctx.mailer:
        return ctx.mailer

    raise RuntimeError("Mailer is None")
