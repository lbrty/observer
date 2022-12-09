from typing import Protocol

from observer.schemas.audit_logs import NewAuditLog


class AuditLogsServiceInterface(Protocol):
    async def add_event(self, event: NewAuditLog):
        raise NotImplementedError

    async def find_by_ref(self, ref_id: str):
        raise NotImplementedError


class AuditLogsService(AuditLogsServiceInterface):
    def __init__(self, repo):
        self.repo = repo

    async def add_event(self, event: NewAuditLog):
        pass

    async def find_by_ref(self, ref_id: str):
        pass
