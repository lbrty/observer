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
    db: Database = None
    keychain: Keychain = None
    mailer: MailerInterface = None
    audit_service: AuditServiceInterface = None
    audit_repo: AuditRepositoryInterface = None
    jwt_service: JWTService = None
    auth_service: AuthServiceInterface = None
    crypto_service: CryptoServiceInterface = None
    mfa_service: MFAServiceInterface = None
    users_repo: UsersRepositoryInterface = None
    users_service: UsersServiceInterface = None


ctx = Context()
