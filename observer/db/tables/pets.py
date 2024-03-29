from sqlalchemy import Column, DateTime, ForeignKey, Index, Table, Text, func, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

pets = Table(
    "pets",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("name", Text(), nullable=False),
    Column("notes", Text(), nullable=True),
    Column("status", Text(), nullable=False),
    Column("registration_id", Text(), nullable=True),  # registration document id
    Column("created_at", DateTime(timezone=True), server_default=func.now(), nullable=True),
    Column("owner_id", UUID(as_uuid=True), ForeignKey("people.id", ondelete="SET NULL"), nullable=True),
    Column("project_id", UUID(as_uuid=True), ForeignKey("projects.id", ondelete="SET NULL"), nullable=True),
    Index("ix_pets_name", text("lower(name)")),
    Index("ix_pets_status", text("status")),
    Index("ix_pets_registration_id", text("status")),
    Index("ix_pets_owner_id", text("owner_id")),
    Index("ix_pets_project_id", text("project_id")),
)
