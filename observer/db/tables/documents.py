from sqlalchemy import Column, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import TIMESTAMP, UUID

from observer.db.tables import metadata
from observer.db.util import utcnow

documents = Table(
    "documents",
    metadata,
    Column("id", UUID(), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("encryption_key", Text(), nullable=True),
    Column("name", Text(), nullable=False),
    Column("path", Text(), nullable=False),
    Column("owner_id", UUID(), nullable=False),
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow),
    Index("ix_documents_name", text("lower(name)")),
)
