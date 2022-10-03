from sqlalchemy import Column, ForeignKey, Index, Table, Text, text
from sqlalchemy.dialects.postgresql import ARRAY, DATE, TIMESTAMP, UUID

from observer.db.tables import metadata
from observer.db.util import utcnow

vulnerability_categories = Table(
    "vulnerability_categories",
    metadata,
    Column("id", UUID, primary_key=True),
    Column("name", Text(), nullable=False),
    Index("ux_categories_name", text("lower(name)"), unique=True),
)

# Phone numbers are encrypted if provided
displaced_persons = Table(
    "displaced_persons",
    metadata,
    Column("id", UUID(), primary_key=True),
    Column("encryption_key", Text(), nullable=True),
    Column("status", Text(), nullable=True),
    Column("external_id", Text(), nullable=True),
    Column("reference_id", Text(), nullable=True),
    Column("email", Text(), nullable=True),
    Column("full_name", Text(), nullable=False),
    Column("birth_date", DATE(), nullable=True),
    Column("notes", Text(), nullable=True),
    Column("phone_number", Text(), nullable=True),
    Column("phone_number_additional", Text(), nullable=True),
    Column("migration_date", DATE(), nullable=True),
    # Location info
    Column("from_city_id", UUID(), ForeignKey("cities.id"), nullable=True),
    Column("from_state_id", UUID(), ForeignKey("states.id"), nullable=True),
    Column("current_city_id", UUID(), ForeignKey("cities.id"), nullable=True),
    Column("current_state_id", UUID(), ForeignKey("states.id"), nullable=True),
    Column("project_id", UUID(), ForeignKey("projects.id"), nullable=True),
    Column("category_id", UUID(), ForeignKey("vulnerability_categories.id"), nullable=True),
    # User's id who registered
    Column("creator_id", UUID(), ForeignKey("users.id"), nullable=True),
    Column("tags", ARRAY(Text()), nullable=True),
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow),
    Column("updated_at", TIMESTAMP(timezone=True), default=utcnow, onupdate=utcnow),
    # Indexes
    Index("ix_displaced_persons_full_name", text("lower(full_name)")),
    Index("ix_displaced_persons_status", "status"),
    Index("ix_displaced_persons_email", text("lower(email)")),
    Index("ix_displaced_persons_birth_date", "birth_date"),
    Index("ix_displaced_persons_tags", "tags"),
)
