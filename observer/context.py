from dataclasses import dataclass

from observer.db import Database
from observer.services.crypto import KeychainLoader


@dataclass
class Context:
    db: Database = None
    key_loader: KeychainLoader = None


ctx = Context()
