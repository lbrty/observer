from sqlalchemy import Column, Table, Text
from sqlalchemy.dialects.postgresql import TIMESTAMP, UUID

from observer.db.tables import metadata
from observer.db.util import utcnow

projects = Table(
    "projects",
    metadata,
    Column("id", UUID(), primary_key=True),
    Column("description", Text(), nullable=False),
    Column("type", Text(), nullable=False),
    Column("consultant_id", UUID(), nullable=False),
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow, nullable=True),
)
