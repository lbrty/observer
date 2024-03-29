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
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

support_records = Table(
    "support_records",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("description", Text(), nullable=True),
    Column("type", Text(), nullable=False),
    Column("consultant_id", UUID(as_uuid=True), nullable=False),
    Column("age_group", Text(), nullable=True),
    Column("record_for", Text(), nullable=False),
    Column("owner_id", UUID(as_uuid=True), nullable=False),
    Column(
        "project_id",
        UUID(as_uuid=True),
        ForeignKey("projects.id", ondelete="SET NULL"),
        nullable=True,
    ),
    Column("created_at", DateTime(timezone=True), server_default=func.now(), nullable=True),
    Index("ix_support_records_type", "type"),
    Index("ix_support_records_description", "description"),
    Index("ix_support_records_consultant_id", "consultant_id"),
    Index("ix_support_records_age_group", "age_group"),
    Index("ix_support_records_owner_id", "owner_id"),
    Index("ix_support_records_project_id", "project_id"),
    CheckConstraint("type IN ('humanitarian', 'legal', 'medical', 'general')", name="support_records_types"),
    CheckConstraint("record_for IN ('person', 'pet')", name="support_records_record_for"),
    CheckConstraint(
        (
            "age_group IN ("
            "'infant', 'toddler', 'pre_school', "
            "'middle_childhood', 'young_teen', "
            "'teenager', 'young_adult', 'early_adult', "
            "'middle_aged_adult', 'old_adult'"
            ")"
        ),
        name="support_records_age_group",
    ),
)
