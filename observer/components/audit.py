from datetime import datetime, timezone
from typing import Any

from box import Box
from fastapi import Header
from fastapi.encoders import jsonable_encoder

from observer.schemas.audit_logs import NewAuditLog


class Props:
    def __init__(self, data: Box):
        self.data = data

    def new_event(self, refs: str, data: Any) -> NewAuditLog:
        now = datetime.now(tz=timezone.utc)
        expires_at = None
        if self.data.expires_in:
            expires_at = now + self.data.expires_in

        return NewAuditLog(
            ref=f"{self.data.tag},{refs}",
            data=jsonable_encoder(data) if data else None,
            expires_at=expires_at,
        )


class Tracked:
    def __init__(self, **data: Any):
        self.tracker = Props(Box(data))

    async def __call__(self) -> Props:
        return self.tracker


async def client_ip(ip_address: str = Header("", alias="X-Forwarded-For")) -> str:
    return ip_address
