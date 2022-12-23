from datetime import datetime, timezone
from typing import List, Optional, Protocol

from sqlalchemy import delete, insert, select

from observer.db import Database
from observer.db.tables.audit_logs import audit_logs
from observer.entities.audit_logs import AuditLog
from observer.schemas.audit_logs import NewAuditLog


class AuditRepositoryInterface(Protocol):
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


class AuditRepository(AuditRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def add_event(self, new_event: NewAuditLog) -> AuditLog:
        values = dict(
            **new_event.dict(),
            created_at=datetime.now(tz=timezone.utc),
        )
        query = insert(audit_logs).values(**values).returning("*")
        result = await self.db.fetchone(query)
        return AuditLog(**result)

    async def find_by_ref(self, ref: str) -> Optional[AuditLog]:
        query = select(audit_logs).where(audit_logs.c.ref == ref)
        if result := await self.db.fetchone(query):
            return AuditLog(**result)

        return None

    async def find_by_key(self, key: str) -> Optional[AuditLog]:
        query = select(audit_logs).where(audit_logs.c.ref.ilike(f"%{key}%"))
        if result := await self.db.fetchone(query):
            return AuditLog(**result)

        return None

    async def find_expired_events(self, expiration: datetime) -> List[AuditLog]:
        query = select(audit_logs).where(audit_logs.c.expires_at < expiration)
        rows = await self.db.fetchone(query)
        return [AuditLog(**row) for row in rows]

    async def delete_event(self, ref: str) -> Optional[AuditLog]:
        query = delete(audit_logs).where(audit_logs.c.ref == ref).returning("*")
        if result := await self.db.fetchone(query):
            return AuditLog(**result)

        return None
