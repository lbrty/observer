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

from observer.db.tables import metadata

vulnerability_categories = Table(
    "vulnerability_categories",
    metadata,
    Column("id", UUID(), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("name", Text(), nullable=False),
    Index("ux_vulnerability_categories_name", text("lower(name)"), unique=True),
)

# Phone numbers are encrypted if provided
displaced_persons = Table(
    "displaced_persons",
    metadata,
    Column("id", UUID(), primary_key=True, server_default=text("gen_random_uuid()")),
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
    Column("consultant_id", UUID(), ForeignKey("users.id"), nullable=True),
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
    Index("ix_displaced_persons_current_city_id", "current_city_id"),
    Index("ix_displaced_persons_from_state_id", "from_state_id"),
    Index("ix_displaced_persons_from_city_id", "from_city_id"),
    Index("ix_displaced_persons_project_id", "project_id"),
    Index("ix_displaced_persons_tags", "tags"),
)
