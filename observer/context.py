from dataclasses import dataclass
from typing import Optional

from observer.db import Database
from observer.repositories.audit_logs import IAuditRepository
from observer.repositories.categories import ICategoryRepository
from observer.repositories.documents import IDocumentsRepository
from observer.repositories.family_members import IFamilyRepository
from observer.repositories.idp import IIDPRepository
from observer.repositories.migration_history import IMigrationRepository
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
from observer.services.idp import IIDPService
from observer.services.jwt import JWTService
from observer.services.mailer import IMailer
from observer.services.mfa import IMFAService
from observer.services.migration_history import IMigrationService
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
    audit: Optional[IAuditRepository] = None
    users: Optional[IUsersRepository] = None
    projects: Optional[IProjectsRepository] = None
    permissions: Optional[IPermissionsRepository] = None
    world: Optional[IWorldRepository] = None
    category: Optional[ICategoryRepository] = None
    idp: Optional[IIDPRepository] = None
    family: Optional[IFamilyRepository] = None
    pets: Optional[IPetsRepository] = None
    documents: Optional[IDocumentsRepository] = None
    support: Optional[ISupportRecordsRepository] = None
    migrations: Optional[IMigrationRepository] = None


@dataclass
class Context:
    db: Optional[Database] = None
    uploads: Optional[UploadHandler] = None
    downloads: Optional[DownloadHandler] = None
    repos: Optional[Repositories] = None
    storage: Optional[IStorage] = None
    keychain: Optional[IKeychain] = None
    mailer: Optional[IMailer] = None
    audit_service: Optional[IAuditService] = None
    jwt_service: Optional[JWTService] = None
    auth_service: Optional[IAuthService] = None
    crypto_service: Optional[ICryptoService] = None
    mfa_service: Optional[IMFAService] = None
    users_service: Optional[IUsersService] = None
    projects_service: Optional[IProjectsService] = None
    permissions_service: Optional[IPermissionsService] = None
    world_service: Optional[IWorldService] = None
    category_service: Optional[ICategoryService] = None
    idp_service: Optional[IIDPService] = None
    family_service: Optional[IFamilyService] = None
    pets_service: Optional[IPetsService] = None
    documents_service: Optional[IDocumentsService] = None
    support_service: Optional[ISupportRecordsService] = None
    migrations_service: Optional[IMigrationService] = None
    secrets_service: Optional[ISecretsService] = None


ctx = Context()
