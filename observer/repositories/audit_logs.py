from typing import Protocol

from observer.db import Database
from observer.schemas.audit_logs import NewAuditLog


class AuditLogsRepositoryInterface(Protocol):
    async def add_event(self, event: NewAuditLog):
        raise NotImplementedError

    async def find_by_ref(self, ref_id: str):
        raise NotImplementedError


class AuditLogsRepository(AuditLogsRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def add_event(self, event: NewAuditLog):
        pass

    async def find_by_ref(self, ref_id: str):
        pass
