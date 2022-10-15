from sqlalchemy import CheckConstraint, Column, Table, Text
from sqlalchemy.dialects.postgresql import TIMESTAMP, UUID

from observer.db.tables import metadata
from observer.db.util import utcnow

projects = Table(
    "projects",
    metadata,
    Column("id", UUID(), primary_key=True),
    Column("description", Text(), nullable=False),
    Column("type", Text(), nullable=False),
    Column("consultant_id", UUID(), nullable=False),
    Column("beneficiary_age", Text(), nullable=True),
    Column("owner_id", UUID(), nullable=False),
    Column("created_at", TIMESTAMP(timezone=True), default=utcnow, nullable=True),
    CheckConstraint("type IN ('humanitarian', 'legal', 'medical', 'general')", name="support_records_types"),
    CheckConstraint(
        "beneficiary_age IN ('0-1', '1-3', '4-5', '6-11', '12-14', '15-17', '18-25', '26-34', '35-59', '60-100+')",
        name="support_records_beneficiary_ages",
    ),
)
