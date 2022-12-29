from dataclasses import dataclass
from typing import Optional

from observer.db import Database
from observer.repositories.audit_logs import IAuditRepository
from observer.repositories.categories import ICategoryRepository
from observer.repositories.idp import IIDPRepository
from observer.repositories.migration_history import IMigrationRepository
from observer.repositories.permissions import IPermissionsRepository
from observer.repositories.projects import IProjectsRepository
from observer.repositories.users import IUsersRepository
from observer.repositories.world import IWorldRepository
from observer.services.audit_logs import IAuditService
from observer.services.auth import IAuthService
from observer.services.categories import ICategoryService
from observer.services.crypto import ICryptoService, Keychain
from observer.services.idp import IIDPService
from observer.services.jwt import JWTService
from observer.services.mailer import IMailer
from observer.services.mfa import IMFAService
from observer.services.migration_history import IMigrationService
from observer.services.permissions import IPermissionsService
from observer.services.projects import IProjectsService
from observer.services.secrets import ISecretsService
from observer.services.users import IUsersService
from observer.services.world import IWorldService


@dataclass
class Context:
    db: Optional[Database] = None
    keychain: Optional[Keychain] = None
    mailer: Optional[IMailer] = None
    audit_service: Optional[IAuditService] = None
    audit_repo: Optional[IAuditRepository] = None
    jwt_service: Optional[JWTService] = None
    auth_service: Optional[IAuthService] = None
    crypto_service: Optional[ICryptoService] = None
    mfa_service: Optional[IMFAService] = None
    users_repo: Optional[IUsersRepository] = None
    users_service: Optional[IUsersService] = None
    projects_repo: Optional[IProjectsRepository] = None
    projects_service: Optional[IProjectsService] = None
    permissions_repo: Optional[IPermissionsRepository] = None
    permissions_service: Optional[IPermissionsService] = None
    world_repo: Optional[IWorldRepository] = None
    world_service: Optional[IWorldService] = None
    category_repo: Optional[ICategoryRepository] = None
    category_service: Optional[ICategoryService] = None
    idp_repo: Optional[IIDPRepository] = None
    idp_service: Optional[IIDPService] = None
    migrations_repo: Optional[IMigrationRepository] = None
    migrations_service: Optional[IMigrationService] = None
    secrets_service: Optional[ISecretsService] = None


ctx = Context()
