from dataclasses import dataclass

from observer.db import Database
from observer.repositories.audit_logs import AuditRepositoryInterface
from observer.repositories.idp import IDPRepositoryInterface
from observer.repositories.permissions import PermissionsRepositoryInterface
from observer.repositories.projects import ProjectsRepositoryInterface
from observer.repositories.users import UsersRepositoryInterface
from observer.repositories.world import WorldRepositoryInterface
from observer.services.audit_logs import AuditServiceInterface
from observer.services.auth import AuthServiceInterface
from observer.services.crypto import CryptoServiceInterface, Keychain
from observer.services.idp import IDPServiceInterface
from observer.services.jwt import JWTService
from observer.services.mailer import MailerInterface
from observer.services.mfa import MFAServiceInterface
from observer.services.permissions import PermissionsServiceInterface
from observer.services.projects import ProjectsServiceInterface
from observer.services.users import UsersServiceInterface
from observer.services.world import WorldServiceInterface


@dataclass
class Context:
    db: Database | None = None
    keychain: Keychain | None = None
    mailer: MailerInterface | None = None
    audit_service: AuditServiceInterface | None = None
    audit_repo: AuditRepositoryInterface | None = None
    jwt_service: JWTService | None = None
    auth_service: AuthServiceInterface | None = None
    crypto_service: CryptoServiceInterface | None = None
    mfa_service: MFAServiceInterface | None = None
    users_repo: UsersRepositoryInterface | None = None
    users_service: UsersServiceInterface | None = None
    projects_repo: ProjectsRepositoryInterface | None = None
    projects_service: ProjectsServiceInterface | None = None
    permissions_repo: PermissionsRepositoryInterface | None = None
    permissions_service: PermissionsServiceInterface | None = None
    world_repo: WorldRepositoryInterface | None = None
    world_service: WorldServiceInterface | None = None
    idp_repo: IDPRepositoryInterface | None = None
    idp_service: IDPServiceInterface | None = None


ctx = Context()
