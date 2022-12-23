from datetime import datetime
from typing import List, Optional, Protocol

from observer.entities.audit_logs import AuditLog
from observer.repositories.audit_logs import AuditRepositoryInterface
from observer.schemas.audit_logs import NewAuditLog


class AuditServiceInterface(Protocol):
    async def add_event(self, new_event: NewAuditLog) -> AuditLog:
        raise NotImplementedError

    async def find_by_ref(self, ref: str) -> Optional[AuditLog]:
        raise NotImplementedError

    async def find_by_key(self, key: str) -> Optional[AuditLog]:
        raise NotImplementedError

    async def find_expired_events(self, expiration: datetime) -> List[AuditLog]:
        raise NotImplementedError

    async def delete_event(self, ref: str) -> Optional[AuditLog]:
        raise NotImplementedError


class AuditService(AuditServiceInterface):
    def __init__(self, repo: AuditRepositoryInterface):
        self.repo = repo

    async def add_event(self, new_event: NewAuditLog) -> AuditLog:
        return await self.repo.add_event(new_event)

    async def find_by_ref(self, ref: str) -> Optional[AuditLog]:
        return await self.repo.find_by_ref(ref)

    async def find_by_key(self, key: str) -> Optional[AuditLog]:
        return await self.repo.find_by_key(key)

    async def find_expired_events(self, expiration: datetime) -> List[AuditLog]:
        return await self.repo.find_expired_events(expiration)

    async def delete_event(self, ref: str) -> Optional[AuditLog]:
        return await self.repo.delete_event(ref)
