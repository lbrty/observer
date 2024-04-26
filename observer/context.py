from dataclasses import dataclass, field
from typing import Type

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
    audit: IAuditRepository = field(init=False)
    users: IUsersRepository = field(init=False)
    projects: IProjectsRepository = field(init=False)
    offices: IOfficesRepository = field(init=False)
    permissions: IPermissionsRepository = field(init=False)
    world: IWorldRepository = field(init=False)
    category: ICategoryRepository = field(init=False)
    people: IPeopleRepository = field(init=False)
    family: IFamilyRepository = field(init=False)
    pets: IPetsRepository = field(init=False)
    documents: IDocumentsRepository = field(init=False)
    support: ISupportRecordsRepository = field(init=False)
    migrations: IMigrationRepository = field(init=False)


@dataclass
class Context:
    db: Database = field(init=False)
    uploads: UploadHandler = field(init=False)
    downloads: DownloadHandler = field(init=False)
    repos: Repositories = field(init=False)
    storage: IStorage = field(init=False)
    keychain: IKeychain = field(init=False)
    mailer: IMailer = field(init=False)
    audit_service: IAuditService = field(init=False)
    jwt_service: JWTService = field(init=False)
    auth_service: IAuthService = field(init=False)
    crypto_service: ICryptoService = field(init=False)
    mfa_service: IMFAService = field(init=False)
    users_service: IUsersService = field(init=False)
    offices_service: IOfficesService = field(init=False)
    projects_service: IProjectsService = field(init=False)
    permissions_service: IPermissionsService = field(init=False)
    world_service: IWorldService = field(init=False)
    category_service: ICategoryService = field(init=False)
    people_service: IPeopleService = field(init=False)
    family_service: IFamilyService = field(init=False)
    pets_service: IPetsService = field(init=False)
    documents_service: IDocumentsService = field(init=False)
    support_service: ISupportRecordsService = field(init=False)
    migrations_service: IMigrationService = field(init=False)
    secrets_service: ISecretsService = field(init=False)


ctx = Context()
