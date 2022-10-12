"""support_records

Revision ID: acefbe2abd79
Revises: 98cf562f08e2
Create Date: 2022-10-09 22:08:36.117371
"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql

from observer.db.util import utcnow

# revision identifiers, used by Alembic.
revision = 'acefbe2abd79'
down_revision = '98cf562f08e2'
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "support_records",
        sa.Column("id", postgresql.UUID(), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("description", sa.Text(), nullable=False),
        sa.Column("type", sa.Text(), nullable=False),
        sa.Column("consultant_id", postgresql.UUID(), nullable=False),
        sa.Column("created_at", postgresql.TIMESTAMP(timezone=True), default=utcnow, nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ix_support_records_type"), "support_records", ["type"])
    op.create_index(op.f("ix_support_records_description"), "support_records", [sa.text("lower(description)")])
    op.create_index(op.f("ix_support_records_consultant_id"), "support_records", ["consultant_id"])


def downgrade():
    op.drop_table("support_records")
