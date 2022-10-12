"""support_records

Revision ID: acefbe2abd79
Revises: 98cf562f08e2
Create Date: 2022-10-09 22:08:36.117371
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

from observer.db.util import utcnow

# revision identifiers, used by Alembic.
revision = "acefbe2abd79"
down_revision = "98cf562f08e2"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "support_records",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("description", sa.Text(), nullable=False),
        sa.Column("type", sa.Text(), nullable=False),
        sa.Column("consultant_id", postgresql.UUID(), nullable=False),
        sa.Column("beneficiary_age", sa.Text(), nullable=True),
        sa.Column("owner_id", postgresql.UUID(), nullable=False),
        sa.Column("created_at", postgresql.TIMESTAMP(timezone=True), default=utcnow, nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ix_support_records_type"), "support_records", ["type"])
    op.create_index(op.f("ix_support_records_description"), "support_records", [sa.text("lower(description)")])
    op.create_index(op.f("ix_support_records_consultant_id"), "support_records", ["consultant_id"])
    op.create_index(op.f("ix_support_records_owner_id"), "support_records", ["owner_id"])

    sa.CheckConstraint("type IN ('humanitarian', 'legal', 'general')", name="support_records_types"),
    sa.CheckConstraint(
        "beneficiary_age IN ('0-1', '1-3', '4-5', '6-11', '12-14', '15-17', '18-25', '26-34', '35-59', '60-100+')",
        name="support_records_beneficiary_ages",
    ),


def downgrade():
    op.drop_table("support_records")
