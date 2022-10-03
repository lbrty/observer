from sqlalchemy import Column, ForeignKey, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import TIMESTAMP, UUID

from observer.db import metadata
from observer.db.util import utcnow

pets = Table(
    "pets",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", Text(), nullable=False),
    Column("notes", Text(), nullable=True),
    Column("status", Text(), nullable=False),
    Column("registration_id", Text(), nullable=True),  # registration document id
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow),
    Column("updated_at", TIMESTAMP(timezone=True), default=utcnow, onupdate=utcnow),
    Column("owner_id", UUID, ForeignKey("displaced_persons.id"), nullable=True),
    Index("ix_pets_name", text("lower(name)")),
    Index("ix_pets_status", text("status")),
    Index("ix_pets_registration_id", text("status")),
)
