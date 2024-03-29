from observer.api.exceptions import InternalError
from observer.context import ctx
from observer.services.audit_logs import IAuditService
from observer.services.auth import IAuthService
from observer.services.categories import ICategoryService
from observer.services.crypto import ICryptoService
from observer.services.documents import IDocumentsService
from observer.services.downloads import DownloadHandler
from observer.services.family_members import IFamilyService
from observer.services.jwt import JWTService
from observer.services.keychain import IKeychain
from observer.services.mailer import IMailer
from observer.services.migration_history import IMigrationService
from observer.services.offices import IOfficesService
from observer.services.people import IPeopleService
from observer.services.permissions import IPermissionsService
from observer.services.pets import IPetsService
from observer.services.projects import IProjectsService
from observer.services.secrets import ISecretsService
from observer.services.storage import IStorage
from observer.services.support_records import ISupportRecordsService
from observer.services.uploads import UploadHandler
from observer.services.users import IUsersService
from observer.services.world import IWorldService


async def users_service() -> IUsersService:
    if ctx.users_service:
        return ctx.users_service

    raise RuntimeError("UsersService is None")


async def projects_service() -> IProjectsService:
    if ctx.projects_service:
        return ctx.projects_service

    raise RuntimeError("ProjectsService is None")


async def permissions_service() -> IPermissionsService:
    if ctx.permissions_service:
        return ctx.permissions_service

    raise RuntimeError("PermissionsService is None")


async def auth_service() -> IAuthService:
    if ctx.auth_service:
        return ctx.auth_service

    raise RuntimeError("AuthService is None")


async def jwt_service() -> JWTService:
    if ctx.jwt_service:
        return ctx.jwt_service

    raise RuntimeError("JWTService is None")


async def crypto_service() -> ICryptoService:
    if ctx.crypto_service:
        return ctx.crypto_service

    raise RuntimeError("CryptoService is None")


async def secrets_service() -> ISecretsService:
    if ctx.secrets_service:
        return ctx.secrets_service

    raise RuntimeError("CryptoService is None")


async def migrations_service() -> IMigrationService:
    if ctx.migrations_service:
        return ctx.migrations_service

    raise RuntimeError("MigrationService is None")


async def keychain() -> IKeychain:
    if not ctx.keychain:
        raise RuntimeError("Keychain is None")

    if len(ctx.keychain.keys) == 0:
        raise InternalError(message="private keys not found")

    return ctx.keychain


async def audit_service() -> IAuditService:
    if ctx.audit_service:
        return ctx.audit_service

    raise RuntimeError("AuditService is None")


async def office_service() -> IOfficesService:
    if ctx.offices_service:
        return ctx.offices_service

    raise RuntimeError("OfficesService is None")


async def storage_service() -> IStorage:
    if ctx.storage:
        return ctx.storage

    raise RuntimeError("Storage is None")


async def world_service() -> IWorldService:
    if ctx.world_service:
        return ctx.world_service

    raise RuntimeError("PlacesService is None")


async def people_service() -> IPeopleService:
    if ctx.people_service:
        return ctx.people_service

    raise RuntimeError("IPeopleService is None")


async def family_service() -> IFamilyService:
    if ctx.family_service:
        return ctx.family_service

    raise RuntimeError("FamilyService is None")


async def pets_service() -> IPetsService:
    if ctx.pets_service:
        return ctx.pets_service

    raise RuntimeError("PetsService is None")


async def documents_service() -> IDocumentsService:
    if ctx.documents_service:
        return ctx.documents_service

    raise RuntimeError("DocumentsService is None")


async def support_records_service() -> ISupportRecordsService:
    if ctx.support_service:
        return ctx.support_service

    raise RuntimeError("SupportService is None")


async def documents_upload() -> UploadHandler:
    if ctx.uploads:
        return ctx.uploads

    raise RuntimeError("UploadHandler is None")


async def documents_download() -> DownloadHandler:
    if ctx.downloads:
        return ctx.downloads

    raise RuntimeError("DownloadHandler is None")


async def category_service() -> ICategoryService:
    if ctx.category_service:
        return ctx.category_service

    raise RuntimeError("CategoryService is None")


async def mailer() -> IMailer:
    if ctx.mailer:
        return ctx.mailer

    raise RuntimeError("Mailer is None")
