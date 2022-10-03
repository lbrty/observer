from sqlalchemy import Column, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db.tables import metadata

documents = Table(
    "documents",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("encryption_key", Text()),
    Column("name", Text(), nullable=False),
    Column("path", Text(), nullable=False),
    Column("owner_id", UUID, nullable=False),
    Index("ix_documents_name", text("owner_id")),
    Index("ux_documents_name", text("lower(name)")),
)
