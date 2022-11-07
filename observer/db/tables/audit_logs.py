from sqlalchemy import Column, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import JSONB, TIMESTAMP, UUID

from observer.db.tables import metadata
from observer.db.util import utcnow

audit_logs = Table(
    "audit_logs",
    metadata,
    Column("id", UUID(), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("ref", Text(), nullable=False),  # format - origin=<user_id...>;source=services:users;action=create:user;
    Column("data", JSONB(), nullable=True, default={}),
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow),
    Column("expires_at", TIMESTAMP(timezone=True), default=utcnow, nullable=True),
    Index("ix_audit_logs_ref", text("lower(ref)")),
)
