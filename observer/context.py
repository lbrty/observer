from dataclasses import dataclass

from observer.db import Database
from observer.repositories.audit_logs import AuditRepositoryInterface
from observer.repositories.users import UsersRepositoryInterface
from observer.services.audit_logs import AuditServiceInterface
from observer.services.auth import AuthServiceInterface
from observer.services.crypto import CryptoServiceInterface, Keychain
from observer.services.jwt import JWTService
from observer.services.mailer import MailerInterface
from observer.services.mfa import MFAServiceInterface
from observer.services.users import UsersServiceInterface


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


ctx = Context()
