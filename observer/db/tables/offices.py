from sqlalchemy import Column, ForeignKey, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

offices = Table(
    "offices",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("name", Text(), nullable=False),
    Column("place_id", UUID(as_uuid=True), ForeignKey("places.id", ondelete="SET NULL"), nullable=True),
    # Indexes
    Index("ix_offices_place_id", "place_id"),
)
