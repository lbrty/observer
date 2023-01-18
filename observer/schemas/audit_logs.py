from datetime import datetime
from typing import Dict, List, Optional

from pydantic import BaseModel, Field

from observer.common.types import Identifier


class AuditLogFilters(BaseModel):
    action: Optional[str] = Field(None, description="Type of action create,update,delete")
    origin: Optional[str] = Field(None, description="Origin something like user_id, document_id etc.")
    source: Optional[str] = Field(None, description="Which parts of system created a given audit log")


class AuditLog(BaseModel):
    id: Identifier = Field(..., description="Audit log id")
    ref: str = Field(
        ...,
        description="Reference in the following format - origin=<user_id...>;source=services:users;action=create:user;",
    )
    data: Optional[Dict] = Field(None, description="JSON slice with changes")
    created_at: datetime = Field(..., description="Creation date time of event")
    expires_at: Optional[datetime] = Field(None, description="Expiration date time of event")


class NewAuditLog(BaseModel):
    ref: str = Field(
        ...,
        description="Reference in the following format - origin=<user_id...>;source=services:users;action=create:user;",
    )
    data: Optional[Dict] = Field(None, description="JSON slice with changes")
    expires_at: Optional[datetime] = Field(None, description="Expiration date time of event")


class AuditLogsResponse(BaseModel):
    total: int = Field(..., description="Total count of records")
    items: List[AuditLog] = Field(..., description="List of audit logs")
