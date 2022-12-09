from typing import Protocol

from sqlalchemy import insert, select

from observer.db import Database
from observer.db.tables.audit_logs import audit_logs
from observer.entities.audit_logs import AuditLog
from observer.schemas.audit_logs import NewAuditLog


class AuditLogsRepositoryInterface(Protocol):
    async def add_event(self, new_event: NewAuditLog) -> AuditLog:
        raise NotImplementedError

    async def find_by_ref(self, ref: str):
        raise NotImplementedError


class AuditLogsRepository(AuditLogsRepositoryInterface):
    def __init__(self, db: Database):
        self.db = db

    async def add_event(self, new_event: NewAuditLog) -> AuditLog:
        query = insert(audit_logs).values(**new_event.dict()).returning("*")
        if result := await self.db.fetchone(query):
            return AuditLog(**result)

    async def find_by_ref(self, ref: str) -> AuditLog:
        query = select(audit_logs).where(audit_logs.c.ref == ref)
        if result := await self.db.fetchone(query):
            return AuditLog(**result)
