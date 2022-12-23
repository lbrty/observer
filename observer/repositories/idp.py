from typing import Protocol

from observer.db import Database


class IDPRepositoryInterface(Protocol):
    ...


class IDPRepository(IDPRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db
