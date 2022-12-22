from sqlalchemy import (
    CheckConstraint,
    Column,
    DateTime,
    ForeignKey,
    Index,
    Table,
    Text,
    func,
    text,
)
from sqlalchemy.dialects.postgresql import ARRAY, DATE, UUID

from observer.db import metadata

categories = Table(
    "categories",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("name", Text(), nullable=False),
    Index("ux_categories_name", text("lower(name)"), unique=True),
)

# Phone numbers are encrypted if provided
displaced_persons = Table(
    "displaced_persons",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
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
    Column("from_place_id", UUID(as_uuid=True), ForeignKey("places.id", ondelete="SET NULL"), nullable=True),
    Column("from_state_id", UUID(as_uuid=True), ForeignKey("states.id", ondelete="SET NULL"), nullable=True),
    Column("current_place_id", UUID(as_uuid=True), ForeignKey("places.id", ondelete="SET NULL"), nullable=True),
    Column("current_state_id", UUID(as_uuid=True), ForeignKey("states.id", ondelete="SET NULL"), nullable=True),
    Column("project_id", UUID(as_uuid=True), ForeignKey("projects.id", ondelete="SET NULL"), nullable=True),
    Column("category_id", UUID(as_uuid=True), ForeignKey("categories.id", ondelete="SET NULL"), nullable=True),
    # User's id who registered
    Column("consultant_id", UUID(as_uuid=True), ForeignKey("users.id", ondelete="SET NULL"), nullable=True),
    Column("tags", ARRAY(Text()), nullable=True),
    Column("created_at", DateTime(timezone=True), server_default=func.now()),
    Column("updated_at", DateTime(timezone=True), server_default=func.now(), onupdate=func.now()),
    # Constraints
    CheckConstraint(
        """status IN (
            'consulted',
            'needs_call_back',
            'needs_legal_support',
            'needs_social_support',
            'needs_monitoring',
            'registered',
            'unknown'
        )""",
        name="displaced_persons_status",
    ),
    # Indexes
    Index("ix_displaced_persons_full_name", text("lower(full_name)")),
    Index("ix_displaced_persons_reference_id", "reference_id"),
    Index("ix_displaced_persons_status", "status"),
    Index("ix_displaced_persons_email", text("lower(email)")),
    Index("ix_displaced_persons_birth_date", "birth_date"),
    Index("ix_displaced_persons_category_id", "category_id"),
    Index("ix_displaced_persons_consultant_id", "consultant_id"),
    Index("ix_displaced_persons_current_state_id", "current_state_id"),
    Index("ix_displaced_persons_current_place_id", "current_place_id"),
    Index("ix_displaced_persons_from_state_id", "from_state_id"),
    Index("ix_displaced_persons_from_place_id", "from_place_id"),
    Index("ix_displaced_persons_project_id", "project_id"),
    Index("ix_displaced_persons_tags", "tags"),
)
