from sqlalchemy import Column, DateTime, Index, Table, Text, func, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db.tables import metadata

documents = Table(
    "documents",
    metadata,
    Column("id", UUID(), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("encryption_key", Text(), nullable=True),
    Column("name", Text(), nullable=False),
    Column("path", Text(), nullable=False),
    Column("owner_id", UUID(), nullable=False),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
    Index("ix_documents_name", text("lower(name)")),
    Index("ix_documents_owner_id", "owner_id"),
)
