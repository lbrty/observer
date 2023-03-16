from dataclasses import dataclass

from observer.db import Database
from observer.repositories.audit_logs import IAuditRepository
from observer.repositories.categories import ICategoryRepository
from observer.repositories.documents import IDocumentsRepository
from observer.repositories.family_members import IFamilyRepository
from observer.repositories.migration_history import IMigrationRepository
from observer.repositories.offices import IOfficesRepository
from observer.repositories.people import IPeopleRepository
from observer.repositories.permissions import IPermissionsRepository
from observer.repositories.pets import IPetsRepository
from observer.repositories.projects import IProjectsRepository
from observer.repositories.support_records import ISupportRecordsRepository
from observer.repositories.users import IUsersRepository
from observer.repositories.world import IWorldRepository
from observer.services.audit_logs import IAuditService
from observer.services.auth import IAuthService
from observer.services.categories import ICategoryService
from observer.services.crypto import ICryptoService, IKeychain
from observer.services.documents import IDocumentsService
from observer.services.downloads import DownloadHandler
from observer.services.family_members import IFamilyService
from observer.services.jwt import JWTService
from observer.services.mailer import IMailer
from observer.services.mfa import IMFAService
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


@dataclass
class Repositories:
    audit: IAuditRepository = None
    users: IUsersRepository = None
    projects: IProjectsRepository = None
    offices: IOfficesRepository = None
    permissions: IPermissionsRepository = None
    world: IWorldRepository = None
    category: ICategoryRepository = None
    people: IPeopleRepository = None
    family: IFamilyRepository = None
    pets: IPetsRepository = None
    documents: IDocumentsRepository = None
    support: ISupportRecordsRepository = None
    migrations: IMigrationRepository = None


@dataclass
class Context:
    db: Database = None
    uploads: UploadHandler = None
    downloads: DownloadHandler = None
    repos: Repositories = None
    storage: IStorage = None
    keychain: IKeychain = None
    mailer: IMailer = None
    audit_service: IAuditService = None
    jwt_service: JWTService = None
    auth_service: IAuthService = None
    crypto_service: ICryptoService = None
    mfa_service: IMFAService = None
    users_service: IUsersService = None
    offices_service: IOfficesService = None
    projects_service: IProjectsService = None
    permissions_service: IPermissionsService = None
    world_service: IWorldService = None
    category_service: ICategoryService = None
    people_service: IPeopleService = None
    family_service: IFamilyService = None
    pets_service: IPetsService = None
    documents_service: IDocumentsService = None
    support_service: ISupportRecordsService = None
    migrations_service: IMigrationService = None
    secrets_service: ISecretsService = None


ctx = Context()
