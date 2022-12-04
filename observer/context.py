from dataclasses import dataclass

from observer.db import Database
from observer.repositories.users import UsersRepositoryInterface
from observer.services.auth import AuthServiceInterface
from observer.services.crypto import CryptoServiceInterface, KeychainLoader
from observer.services.jwt import JWTService
from observer.services.mfa import MFAServiceInterface
from observer.services.users import UsersServiceInterface


@dataclass
class Context:
    db: Database = None
    key_loader: KeychainLoader = None
    jwt_service: JWTService = None
    auth_service: AuthServiceInterface = None
    crypto_service: CryptoServiceInterface = None
    mfa_service: MFAServiceInterface = None
    users_repo: UsersRepositoryInterface = None
    users_service: UsersServiceInterface = None


ctx = Context()
