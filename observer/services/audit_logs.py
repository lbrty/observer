from typing import Protocol

from observer.entities.audit_logs import AuditLog
from observer.repositories.audit_logs import AuditRepositoryInterface
from observer.schemas.audit_logs import NewAuditLog


class AuditServiceInterface(Protocol):
    async def add_event(self, event: NewAuditLog):
        raise NotImplementedError

    async def find_by_ref(self, ref_id: str) -> AuditLog:
        raise NotImplementedError


class AuditService(AuditServiceInterface):
    def __init__(self, repo: AuditRepositoryInterface):
        self.repo = repo

    async def add_event(self, new_event: NewAuditLog) -> AuditLog:
        return await self.repo.add_event(new_event)

    async def find_by_ref(self, ref: str) -> AuditLog:
        return await self.repo.find_by_ref(ref)
