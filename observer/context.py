from dataclasses import dataclass

from observer.db import Database
from observer.services.crypto import KeychainLoader
from observer.services.jwt import JWTHandler
from observer.services.users import UsersServiceInterface


@dataclass
class Context:
    db: Database = None
    key_loader: KeychainLoader = None
    jwt_handler: JWTHandler = None
    users_service: UsersServiceInterface = None


ctx = Context()
