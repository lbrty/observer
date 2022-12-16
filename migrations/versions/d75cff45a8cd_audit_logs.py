"""audit_logs

Revision ID: d75cff45a8cd
Revises: 6e0abb758075
Create Date: 2022-10-08 22:13:49.480561
"""
import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

from observer.db.util import utcnow

# revision identifiers, used by Alembic.
revision = "d75cff45a8cd"
down_revision = "6e0abb758075"
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        "audit_logs",
        sa.Column("id", postgresql.UUID(as_uuid=True), nullable=False, server_default=sa.text("gen_random_uuid()")),
        sa.Column("ref", sa.Text(), nullable=False),
        sa.Column("data", postgresql.JSONB(), nullable=True, server_default="{}"),
        sa.Column("created_at", postgresql.TIMESTAMP(timezone=True), default=utcnow),
        sa.Column("expires_at", postgresql.TIMESTAMP(timezone=True), default=utcnow, nullable=True),
        sa.PrimaryKeyConstraint("id"),
    )

    op.create_index(op.f("ix_audit_logs_ref"), "audit_logs", [sa.text("lower(ref)")])


def downgrade():
    op.drop_table("audit_logs")
