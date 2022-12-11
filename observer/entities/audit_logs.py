from datetime import datetime

from pydantic import BaseModel

from observer.common.types import Identifier


class AuditLog(BaseModel):
    id: Identifier
    ref: str  # format - origin=<user_id...>;source=services:users;action=create:user;
    data: dict | None
    created_at: datetime
    expires_at: datetime | None
