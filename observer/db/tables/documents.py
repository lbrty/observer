from sqlalchemy import Column, DateTime, ForeignKey, Index, Table, Text, func, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

documents = Table(
    "documents",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("encryption_key", Text(), nullable=True),
    Column("name", Text(), nullable=False),
    Column("path", Text(), nullable=False),
    Column("mimetype", Text(), nullable=False),
    Column(
        "project_id",
        UUID(as_uuid=True),
        ForeignKey("projects.id", ondelete="CASCADE"),
        nullable=False,
    ),
    Column("owner_id", UUID(as_uuid=True), nullable=False),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
    Index("ix_documents_name", text("lower(name)")),
    Index("ix_documents_owner_id", "owner_id"),
    Index("ix_documents_project", "project_id"),
)
