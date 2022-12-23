from datetime import datetime, timedelta, timezone
from typing import Protocol

from observer.repositories.idp import IDPRepositoryInterface
from observer.schemas.audit_logs import NewAuditLog


class IDPServiceInterface(Protocol):
    tag: str
    repo: IDPRepositoryInterface

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        raise NotImplementedError


class IDPService(IDPServiceInterface):
    tag: str = "source=service:idp"

    def __init__(self, idp_repository: IDPRepositoryInterface):
        self.repo = idp_repository

    async def create_log(self, ref: str, expires_in: timedelta | None, data: dict | None = None) -> NewAuditLog:
        now = datetime.now(tz=timezone.utc)
        expires_at = None
        if expires_in:
            expires_at = now + expires_in

        return NewAuditLog(
            ref=f"{self.tag},{ref}",
            data=data,
            expires_at=expires_at,
        )
