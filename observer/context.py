from dataclasses import dataclass
from typing import Optional

from observer.db import Database
from observer.repositories.audit_logs import AuditRepositoryInterface
from observer.repositories.categories import CategoryRepositoryInterface
from observer.repositories.idp import IDPRepositoryInterface
from observer.repositories.permissions import PermissionsRepositoryInterface
from observer.repositories.projects import ProjectsRepositoryInterface
from observer.repositories.users import UsersRepositoryInterface
from observer.repositories.world import WorldRepositoryInterface
from observer.services.audit_logs import AuditServiceInterface
from observer.services.auth import AuthServiceInterface
from observer.services.categories import CategoryServiceInterface
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
    db: Optional[Database] = None
    keychain: Optional[Keychain] = None
    mailer: Optional[MailerInterface] = None
    audit_service: Optional[AuditServiceInterface] = None
    audit_repo: Optional[AuditRepositoryInterface] = None
    jwt_service: Optional[JWTService] = None
    auth_service: Optional[AuthServiceInterface] = None
    crypto_service: Optional[CryptoServiceInterface] = None
    mfa_service: Optional[MFAServiceInterface] = None
    users_repo: Optional[UsersRepositoryInterface] = None
    users_service: Optional[UsersServiceInterface] = None
    projects_repo: Optional[ProjectsRepositoryInterface] = None
    projects_service: Optional[ProjectsServiceInterface] = None
    permissions_repo: Optional[PermissionsRepositoryInterface] = None
    permissions_service: Optional[PermissionsServiceInterface] = None
    world_repo: Optional[WorldRepositoryInterface] = None
    world_service: Optional[WorldServiceInterface] = None
    category_repo: Optional[CategoryRepositoryInterface] = None
    category_service: Optional[CategoryServiceInterface] = None
    idp_repo: Optional[IDPRepositoryInterface] = None
    idp_service: Optional[IDPServiceInterface] = None


ctx = Context()
