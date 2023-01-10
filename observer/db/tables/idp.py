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
people = Table(
    "people",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("status", Text(), nullable=True),
    Column("external_id", Text(), nullable=True),
    Column("reference_id", Text(), nullable=True),
    Column("email", Text(), nullable=True),
    Column("full_name", Text(), nullable=False),
    Column("sex", Text(), nullable=True),
    Column("pronoun", Text(), nullable=True),
    Column("birth_date", DATE(), nullable=True),
    Column("notes", Text(), nullable=True),
    Column("phone_number", Text(), nullable=True),
    Column("phone_number_additional", Text(), nullable=True),
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
        name="people_status",
    ),
    CheckConstraint(
        """sex IN (
            'male',
            'female',
            'unknown'
        )""",
        name="people_sex",
    ),
    # Indexes
    Index("ix_people_full_name", text("lower(full_name)")),
    Index("ix_people_reference_id", "reference_id"),
    Index("ix_people_status", "status"),
    Index("ix_people_email", text("lower(email)")),
    Index("ix_people_birth_date", "birth_date"),
    Index("ix_people_category_id", "category_id"),
    Index("ix_people_consultant_id", "consultant_id"),
    Index("ix_people_project_id", "project_id"),
    Index("ix_people_tags", "tags"),
)
