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
    Column("beneficiary_age", Text(), nullable=True),
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
    Index("ix_support_records_beneficiary_age", "beneficiary_age"),
    Index("ix_support_records_owner_id", "owner_id"),
    Index("ix_support_records_project_id", "project_id"),
    CheckConstraint("type IN ('humanitarian', 'legal', 'medical', 'general')", name="support_records_types"),
    CheckConstraint(
        "beneficiary_age IN ('0-1', '1-3', '4-5', '6-11', '12-14', '15-17', '18-25', '26-34', '35-59', '60-100+')",
        name="support_records_beneficiary_ages",
    ),
)
