from dataclasses import dataclass

from observer.db import Database
from observer.services.crypto import KeychainLoader
from observer.services.jwt import JWTHandler


@dataclass
class Context:
    db: Database = None
    key_loader: KeychainLoader = None
    jwt_handler: JWTHandler = None


ctx = Context()
