from sqlalchemy import CheckConstraint, Column, DateTime, Index, Table, Text, func, text
from sqlalchemy.dialects.postgresql import UUID

from observer.db import metadata

projects = Table(
    "support_records",
    metadata,
    Column("id", UUID(), primary_key=True, server_default=text("gen_random_uuid()")),
    Column("description", Text(), nullable=True),
    Column("type", Text(), nullable=False),
    Column("consultant_id", UUID(), nullable=False),
    Column("beneficiary_age", Text(), nullable=True),
    Column("owner_id", UUID(), nullable=False),
    Column("created_at", DateTime(timezone=True), server_default=func.now(), nullable=True),
    Index("ix_projects_type", "type"),
    Index("ix_projects_description", "description"),
    Index("ix_projects_consultant_id", "consultant_id"),
    Index("ix_projects_beneficiary_age", "beneficiary_age"),
    Index("ix_projects_owner_id", "owner_id"),
    CheckConstraint("type IN ('humanitarian', 'legal', 'medical', 'general')", name="support_records_types"),
    CheckConstraint(
        "beneficiary_age IN ('0-1', '1-3', '4-5', '6-11', '12-14', '15-17', '18-25', '26-34', '35-59', '60-100+')",
        name="support_records_beneficiary_ages",
    ),
)
